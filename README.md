# MOTD Client

A Go client for interacting with Message of the Day (MOTD) services.

- [MOTD Client](#motd-client)
  - [Description](#description)
  - [Installation](#installation)
  - [Usage](#usage)
  - [Configuration](#configuration)
  - [Project Structure](#project-structure)
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

```plaintext
motd-client/
├── main.go         # Main application entry point
├── config.go       # Configuration management
├── go.mod          # Go module file
├── go.sum          # Dependency checksums
└── README.md       # This file
```

## Dependencies

This project uses Go modules for dependency management. Dependencies are listed in the `go.mod` file.

- **Go 1.24** - Latest Go version
- **github.com/kelseyhightower/envconfig** - Environment variable configuration

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Author

- Stevie LCB (@stevielcb)
