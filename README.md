# MOTD Client

A Go client for interacting with Message of the Day (MOTD) services.

- [MOTD Client](#motd-client)
  - [Description](#description)
  - [Installation](#installation)
  - [Usage](#usage)
  - [Configuration](#configuration)
  - [Project Structure](#project-structure)
    - [Architecture Benefits](#architecture-benefits)
  - [Testing](#testing)
    - [Test Coverage](#test-coverage)
  - [Dependencies](#dependencies)
  - [License](#license)
  - [Author](#author)

## Description

MOTD Client is a Go-based client application designed to fetch and display Message of the Day content from various services. This client provides a simple and efficient way to retrieve MOTD information.

## Installation

```bash
# Clone the repository
git clone https://github.com/stevielcb/motd-client.git

# Navigate to the project directory
cd motd-client

# Install dependencies
go mod download
```

## Usage

```bash
# Build the project
go build

# Run the client
./motd-client

# Run with debug logging
MOTD_LOG_LEVEL=debug ./motd-client
```

## Configuration

The client can be configured using environment variables with the `MOTD_` prefix:

| Variable | Default | Description |
|----------|---------|-------------|
| `MOTD_HOST` | `localhost` | Server hostname |
| `MOTD_PORT` | `4200` | Server port |
| `MOTD_TIMEOUT_MS` | `100` | Connection timeout in milliseconds |
| `MOTD_LOG_LEVEL` | `info` | Log level (debug, info, warn, error) |

Example:

```bash
MOTD_HOST=example.com MOTD_PORT=8080 MOTD_LOG_LEVEL=debug ./motd-client
```

## Project Structure

The project follows a clean architecture pattern with proper separation of concerns:

```plaintext
motd-client/
├── main.go                    # Main application entry point
├── go.mod                     # Go module file
├── go.sum                     # Dependency checksums
├── README.md                  # This file
└── internal/                  # Internal packages
    ├── app/                   # Application orchestration
    │   ├── app.go            # Main application logic
    │   └── app_test.go       # Unit tests for application logic
    ├── config/               # Configuration management
    │   ├── config.go         # Configuration loading and validation
    │   └── config_test.go    # Unit tests for configuration
    ├── logger/               # Logging setup
    │   ├── logger.go         # Structured logging configuration
    │   └── logger_test.go    # Unit tests for logging
    ├── network/              # Network communication
    │   ├── client.go         # TCP client for server communication
    │   └── client_test.go    # Unit tests for network client
    └── terminal/             # Terminal environment handling
        ├── terminal.go       # Terminal detection and formatting
        └── terminal_test.go  # Unit tests for terminal package
```

### Architecture Benefits

- **Separation of Concerns**: Each package has a single responsibility
- **Testability**: Components can be tested in isolation with comprehensive unit tests
- **Maintainability**: Clear boundaries between different functionalities
- **Reusability**: Packages can be reused in other contexts
- **Dependency Injection**: Dependencies are explicitly passed rather than using globals
- **Interface-based Design**: Uses interfaces for better testability and flexibility

## Testing

The project includes comprehensive unit tests for all packages:

```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal/terminal/...
go test ./internal/network/...
go test ./internal/config/...
go test ./internal/logger/...
go test ./internal/app/...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...
```

### Test Coverage

- **Terminal Package**: Tests terminal environment detection and message formatting
- **Network Package**: Tests TCP client functionality with mock servers
- **Config Package**: Tests configuration loading and validation
- **Logger Package**: Tests logging setup and configuration
- **App Package**: Tests application orchestration with mocked dependencies

## Dependencies

This project uses Go modules for dependency management. Dependencies are listed in the `go.mod` file.

- **Go 1.24** - Latest Go version
- **github.com/kelseyhightower/envconfig** - Environment variable configuration

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Author

- Stevie LCB (@stevielcb)
