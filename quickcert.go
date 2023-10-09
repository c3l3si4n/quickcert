package main 

import (
	"bufio"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"os"
	"strings"
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

func main() {
	uniqueMap := make(map[string]bool)

	conn, err := pgx.Connect(context.Background(), CRTSH_DATABASE_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	query := `
WITH ci AS (
    SELECT min(sub.CERTIFICATE_ID) ID,
           x509_commonName(sub.CERTIFICATE) COMMON_NAME,
           min(sub.ISSUER_CA_ID) ISSUER_CA_ID,
           array_agg(DISTINCT sub.NAME_VALUE) NAME_VALUES,

           x509_notBefore(sub.CERTIFICATE) NOT_BEFORE,
           x509_notAfter(sub.CERTIFICATE) NOT_AFTER,
           encode(x509_serialNumber(sub.CERTIFICATE), 'hex') SERIAL_NUMBER
        FROM (SELECT cai.*
                  FROM certificate_and_identities cai
                  WHERE plainto_tsquery('certwatch', '%s') @@ identities(cai.CERTIFICATE)
                      AND cai.NAME_VALUE ILIKE ('%%' || '%s')
                  LIMIT 1000 OFFSET %d
) sub
        GROUP BY sub.CERTIFICATE
)
SELECT 

        ci.COMMON_NAME
        
    FROM ci
            LEFT JOIN LATERAL (
                SELECT min(ctle.ENTRY_TIMESTAMP) ENTRY_TIMESTAMP
                    FROM ct_log_entry ctle
                    WHERE ctle.CERTIFICATE_ID = ci.ID
            ) le ON TRUE,
         ca
    WHERE ci.ISSUER_CA_ID = ca.ID
    ORDER BY le.ENTRY_TIMESTAMP NULLS LAST;
`
	// iterate lines in stdin
	// for each line, prepare a query
	stdin := IterStdin()
	for _, line := range stdin {
		page := 0
		for {
			offset := 1000 * page
			preparedQuery := fmt.Sprintf(query, line, line, offset)

			rows, err := conn.Query(context.Background(), preparedQuery)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
				os.Exit(1)
			}
			subdomains, err := pgx.CollectRows(rows, pgx.RowTo[string])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
				os.Exit(1)
			}
			if len(subdomains) == 0 {
				break
			}
			for _, subdomain := range subdomains {
				if !strings.Contains(subdomain, line) {
					continue
				}
				if strings.HasPrefix(subdomain, "*.") {
					subdomain = strings.Replace(subdomain, "*.", "", 1)
				}
				if _, ok := uniqueMap[subdomain]; !ok {
					fmt.Println(subdomain)
					uniqueMap[subdomain] = true

				}

			}
			page += 1
		}

	}
}
