# Ditto - Simple File Transfer Utility

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Ditto is a lightweight TCP-based file transfer utility written in Go. It allows for efficient file transfer between computers over a network connection using a client-server architecture.

## Features

- **Client-Server Architecture**: Transfer files between machines with different roles
- **Duplicate Detection**: Automatically skips files that already exist at the destination
- **Multi-File Transfer**: Send multiple files in a single command
- **Progress Reporting**: View transfer progress and file sizes
- **Cross-Platform**: Works on Windows, Linux, and macOS
- **No External Dependencies**: Pure Go implementation with no third-party packages

## Installation

### From Source

Requires Go 1.15 or higher.

```bash
# Clone repository
git clone https://github.com/yourusername/ditto.git
cd ditto

# Build for your current platform
make build

# Or build for all supported platforms
make cross-build
```

### Pre-built Binaries

Download the appropriate binary for your platform from the [Releases](https://github.com/omkar-sutar/ditto/releases/) page.

## Usage

### Starting a Server

```bash
# Basic usage
./ditto server 8080 /path/to/destination

# Examples
./ditto server 8080 /home/user/downloads
./ditto server 9000 C:/Users/user/Downloads
```

The server will listen on the specified port and save received files to the destination directory.

### Sending Files (Client)

```bash
# Basic usage
./ditto client address:port file1 [file2 file3 ...]

# Examples
./ditto client localhost:8080 document.pdf image.jpg
./ditto client 192.168.1.100:8080 *.txt
./ditto client server.example.com:9000 backup.zip
```

## Building from Source

### Prerequisites

- Go 1.15 or higher
- GNU Make (for using the Makefile)

### Build Commands

```bash
# Update dependencies
make deps

# Build for current platform
make build

# Build for specific platforms
make linux
make windows

# Build for all platforms
make cross-build

# Create distribution packages
make dist

# Clean build artifacts
make clean
```

Build outputs are placed in the `bin/` directory, organized by platform and architecture.

## Network Protocol

Ditto uses a simple binary protocol for file transfers:

1. **Connection Start**: Client sends a connection start byte
2. **Per-File Transfer**:
   - Send file name (with length header)
   - Send file size (with length header)
   - Server checks for file existence
   - If file doesn't exist, client sends file data
3. **Connection End**: Client sends a connection end byte

## Development

### Project Structure

```
ditto/
├── main.go      # Entry point and CLI handling
├── client.go    # Client implementation
├── server.go    # Server implementation 
├── control.go   # Protocol constants and shared code
├── Makefile     # Build configuration
└── README.md    # This file
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Acknowledgments

- Built with [Go](https://golang.org/)
