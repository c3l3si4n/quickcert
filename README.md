# quickcert
This is a simple tool that will directly connect to crt.sh's PostgreSQL server and dump all subdomains related to a company through that. This has some benefits over the traditional method of querying crt.sh's API, such as not having random timeouts and having no limits on output size.

## Features
- Subdomains are automatically ordered by oldest to newest.
- Uses fast postgresql client library (pgx)

## Installation
```
go install github.com/c3l3si4n/quickcert@HEAD
```

## Usage
```
echo att.com | quickcert
```
