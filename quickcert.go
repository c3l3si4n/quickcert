package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5"
)

// psql -h crt.sh -p 5432 -U guest certwatch
var CRTSH_DATABASE_URL = "postgres://guest@crt.sh:5432/certwatch?sslmode=disable&default_query_exec_mode=simple_protocol"

func IterStdin() []string {
	out := []string{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		out = append(out, scanner.Text())
	}

	return out
}

var Limit = 15000

func main() {
	uniqueMap := make(map[string]bool)
	uniqueMapMutex := sync.Mutex{}

	query := `
WITH ci AS (
    SELECT 
           x509_commonName(sub.CERTIFICATE) COMMON_NAME

         
        FROM (SELECT cai.CERTIFICATE CERTIFICATE
                  FROM certificate_and_identities cai
                  WHERE plainto_tsquery('certwatch', '%s') @@ identities(cai.CERTIFICATE)
                      AND cai.NAME_VALUE ILIKE ('%%' || '%s')
                  LIMIT %d OFFSET %d
) sub

        GROUP BY sub.CERTIFICATE
)
SELECT 

        ci.COMMON_NAME COMMON_NAME
        
    FROM ci
            
            
         
    WHERE COMMON_NAME IS NOT NULL
`
	// iterate lines in stdin
	// for each line, prepare a query
	stdin := IterStdin()
	for _, line := range stdin {
		page := 0
		var routineQueue = make(chan bool, 10)
		stop := false
		var wait sync.WaitGroup
		lineCopy := line
		for {
			if stop {
				break
			}
			offset := Limit * page

			line := strings.ToLower(lineCopy)
			preparedQuery := fmt.Sprintf(query, line, line, Limit, offset)
			wait.Add(1)
			routineQueue <- true
			pageTmp := page
			go func(page int) {
				retries := 0
				success := false
				for !success && retries <= 5 {

					conn, err := pgx.Connect(context.Background(), CRTSH_DATABASE_URL)

					if err != nil {
						fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
						os.Exit(1)
					}
					rows, err := conn.Query(context.Background(), preparedQuery)
					if err != nil {
						log.Println(err)

						rows, err = conn.Query(context.Background(), preparedQuery)
						if err != nil {
							retries += 1
							conn.Close(context.Background())

							continue
						}
					}
					subdomains, err := pgx.CollectRows(rows, pgx.RowTo[string])
					if err != nil {
						log.Println(err)

						subdomains, err = pgx.CollectRows(rows, pgx.RowTo[string])
						if err != nil {
							retries += 1
							conn.Close(context.Background())

							continue
						}
					}
					if len(subdomains) == 0 {
						stop = true
						<-routineQueue
						wait.Done()
						success = true
						conn.Close(context.Background())

						return
					}
					for _, subdomain := range subdomains {
						subdomain = strings.ToLower(subdomain)
						if !strings.Contains(subdomain, line) {
							continue
						}
						if strings.HasPrefix(subdomain, "*.") {
							subdomain = strings.Replace(subdomain, "*.", "", 1)
						}
						uniqueMapMutex.Lock()
						if _, ok := uniqueMap[subdomain]; !ok {
							fmt.Println(subdomain)
							uniqueMap[subdomain] = true

						}
						uniqueMapMutex.Unlock()

					}
					conn.Close(context.Background())
					success = true

				}
				wait.Done()
				<-routineQueue

			}(pageTmp)
			page += 1

		}

		wait.Wait()

	}
}
