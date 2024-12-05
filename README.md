# QuickCert

QuickCert is a high-performance tool for extracting subdomains from SSL/TLS certificate transparency logs using crt.sh's PostgreSQL database. Unlike traditional HTTP API methods, QuickCert offers improved reliability and unlimited result retrieval by directly connecting to the certificate transparency database.

## Features

- 🚀 Direct PostgreSQL connection to crt.sh database
- 💪 Multi-threaded processing (10 concurrent connections)
- 🔄 Automatic retry mechanism for failed queries
- 🎯 Smart duplicate filtering
- ⚡ High-performance using pgx driver
- 📝 Case-insensitive matching
- 🧹 Automatic wildcard certificate handling

## Installation

### Using Go Install
```bash
go install github.com/c3l3si4n/quickcert@HEAD
```

### Building from Source
```bash
git clone https://github.com/c3l3si4n/quickcert.git
cd quickcert
go build
```

## Usage

### Basic Usage
```bash
echo "example.com" | quickcert
```

### Multiple Domains
```bash
cat domains.txt | quickcert
```

### Combining with Other Tools
```bash
echo "example.com" | quickcert | tee subdomains.txt
```

## Technical Details

- Database: Connects to crt.sh PostgreSQL database (certwatch)
- Connection String: `postgres://guest@crt.sh:5432/certwatch`
- Query Limit: 15,000 records per page
- Retry Mechanism: Up to 5 retries per failed query
- Concurrent Connections: 10 parallel queries

## Features in Detail

1. **Duplicate Handling**
   - Automatically removes duplicate subdomains
   - Converts all domains to lowercase for consistent matching

2. **Wildcard Certificate Processing**
   - Automatically strips `*.` from wildcard certificates
   - Ensures proper subdomain formatting

3. **Error Handling**
   - Graceful handling of database connection issues
   - Automatic query retries on failure
   - Concurrent connection management

## Limitations

- Fixed number of concurrent connections (10)
- Dependent on crt.sh database availability

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Acknowledgments

- crt.sh for providing public access to their certificate transparency database
- The Go community for excellent database drivers and tools

## Author

[c3l3si4n](https://github.com/c3l3si4n)
