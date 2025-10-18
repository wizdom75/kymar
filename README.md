# Kymar - Database Client Pro

A professional cross-platform database client built with Go and Fyne, supporting MySQL and PostgreSQL.

## Features

- 🎨 Beautiful dark theme with modern UI
- 🔐 Support for TCP/IP and SSH tunnel connections
- 🗄️ MySQL and PostgreSQL support
- ⚡ Fast query execution with keyboard shortcuts (Cmd+Enter)
- 📊 Automatic table browsing and data preview
- 🔍 Intelligent column width adjustment
- 💾 Save and manage connection credentials

## Project Structure

```
kymar/
├── cmd/
│   └── kymar/             # Application entry point
│       └── main.go
├── internal/              # Private application code
│   ├── config/           # Configuration management
│   │   └── config.go
│   ├── db/               # Database connection logic
│   │   ├── connection.go
│   │   └── models.go
│   ├── ssh/              # SSH tunnel support
│   │   └── tunnel.go
│   └── ui/               # User interface components
│       ├── theme.go
│       ├── login.go
│       └── main_interface.go
├── go.mod
├── go.sum
└── README.md
```

## Installation

### Prerequisites

- Go 1.24.5 or later
- C compiler (for CGO dependencies)

### Build

```bash
# Using Makefile (recommended)
make build

# Or manually
go build -o kymar ./cmd/kymar
```

### Run

```bash
# For development (recommended)
make run

# Or run the binary
./kymar

# Or run directly with go
go run ./cmd/kymar
```

## Development

### Available Make Commands

```bash
make help          # Show all available commands
make run           # Run the app (for development)
make build         # Build the binary
make clean         # Remove build artifacts
make test          # Run tests
make fmt           # Format code
make vet           # Run go vet
make check         # Run fmt, vet, and build
make install       # Install dependencies
make info          # Show project information
make size          # Show binary size
```

For a complete list of commands, run `make help`.

## Usage

### Quick Start

1. Launch the application
2. Select connection type (TCP/IP or SSH)
3. Enter database credentials
4. Click "Connect"
5. Browse tables in the sidebar
6. Click a table to automatically load its data
7. Write custom queries in the SQL editor
8. Press Cmd+Enter or click "Run Query" to execute

### Keyboard Shortcuts

- `Cmd+Enter` - Execute query

## Configuration

### Saved Connections

Connection credentials are automatically saved to `~/.kymar/connections.json` when you check "Save this connection" during login.

**Security Note**: Passwords are currently stored in plain text in the configuration file. This is similar to many database clients (Sequel Ace, MySQL Workbench, etc.) but means you should ensure proper file permissions on your system. The configuration file is created with `0600` permissions (owner read/write only).

Future improvements may include:
- OS keychain integration
- Password encryption
- SSH key authentication support

## Dependencies

- [Fyne](https://fyne.io/) - Cross-platform GUI toolkit
- [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) - MySQL driver
- [lib/pq](https://github.com/lib/pq) - PostgreSQL driver
- [golang.org/x/crypto/ssh](https://pkg.go.dev/golang.org/x/crypto/ssh) - SSH client

## Development

### Code Organization

The project follows Go best practices:

- `cmd/` - Application entry points
- `internal/` - Private application code (not importable by external projects)
  - `db/` - Database connection and query logic
  - `ssh/` - SSH tunnel implementation
  - `ui/` - User interface components and screens

### Package Structure

- **internal/db**: Database connection management, DSN building, connection pooling
- **internal/ssh**: SSH tunnel dialer for secure database connections
- **internal/ui**: All UI components including theme, login screen, and main interface
- **internal/config**: Configuration and saved connections management

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

