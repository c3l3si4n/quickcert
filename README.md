# quickcert
This is a simple tool that will directly connect to crt.sh's PostgreSQL server and dump all subdomains related to a company through parsing of certificate transparency logs. This has some benefits over the traditional method of querying crt.sh's HTTP API, such as being more resilient to random timeouts and having no limits on output quantity.

## Features
- Subdomains are automatically ordered by oldest to newest.
- Multi-threaded (number of threads is fixed)
- Uses fast postgresql client library (pgx)

## Installation
```
go install github.com/c3l3si4n/quickcert@HEAD
```

## Usage
```
echo att.com | quickcert
```
