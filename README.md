# 📮 Postal-Go

A powerful and flexible Go client library for the [Postal](https://github.com/postalserver/postal) email server. This library provides a simple, middleware-based approach to interact with Postal's API, making it easy to send emails with advanced features and customizations.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/doc/go1.21) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT) [![Test Coverage](https://img.shields.io/badge/coverage-95.5%25-brightgreen.svg)](coverage.html)

## About Postal
Postal is a complete and fully featured mail server for use by websites & web servers. Think Sendgrid, Mailgun or Postmark but open source and ready for you to run on your own servers. For more information, visit the [Postal GitHub repository](https://github.com/postalserver/postal).

## ✨ Features

- 🚀 **Simple and intuitive API** - Clean, idiomatic Go interface
- 🔌 **Middleware support** - Extensible request/response processing pipeline
- 📧 **Multiple message types** - Support for both simple and raw MIME messages
- 🔄 **Built-in retry mechanism** - Automatic retry with configurable backoff
- ⚡ **Concurrent request handling** - Thread-safe for high-throughput applications
- ✅ **Comprehensive validation** - Input validation with detailed error messages
- 🔍 **Debug mode** - Detailed logging for troubleshooting
- 📊 **High test coverage** - 95.5% test coverage with extensive test suite
- 🎨 **HTML email support** - Send rich HTML emails with attachments
- 🔧 **Highly configurable** - Flexible configuration through functional options
- 🏃 **Minimal dependencies** - Only uses `golang.org/x/time` for rate limiting

## 📦 Installation

```bash
go get github.com/sachin-duhan/postal-go
```

**Requirements:**
- Go 1.21 or higher
- A Postal server instance (for production use)
- Docker and Docker Compose (for integration testing)

## 🚀 Quick Start

```go
package main

import (
    "context"
    "log"
    
    postal "github.com/sachin-duhan/postal-go"
    "github.com/sachin-duhan/postal-go/common/types"
)

func main() {
    // Initialize client with functional options
    client := postal.NewClient(
        postal.WithAPIKey("your-api-key"),
        postal.WithBaseURL("https://postal.example.com"),
        postal.WithDebug(true), // Enable debug logging
    )

    // Create a message
    message := &types.Message{
        To:       []string{"recipient@example.com"},
        From:     "sender@yourdomain.com",
        Subject:  "Hello from Postal-Go!",
        Body:     "This is a plain text email.",
        HTMLBody: `<h1>Hello!</h1><p>This is an <strong>HTML</strong> email.</p>`,
    }

    // Send the message
    result, err := client.SendMessage(context.Background(), message)
    if err != nil {
        log.Fatal(err)
    }
    
    if result.Success() {
        log.Printf("Message sent successfully! ID: %s", result.MessageID)
    }
}
```

## 📚 Documentation

### Examples
Explore our comprehensive examples:
- [Simple Example](./examples/simple/main.go) - Basic message sending
- [Advanced Example](./examples/advanced/main.go) - Middleware, attachments, and advanced features

### Advanced Usage

#### Configuration Options
```go
client := postal.NewClient(
    postal.WithAPIKey("your-api-key"),
    postal.WithBaseURL("https://postal.example.com"),
    postal.WithTimeout(30 * time.Second),
    postal.WithMaxRetries(3),
    postal.WithRetryInterval(time.Second),
    postal.WithMaxConcurrency(10),
    postal.WithDebug(true),
    postal.WithMiddleware(customMiddleware),
)
```

#### Sending Messages with Attachments
```go
message := &types.Message{
    To:       []string{"recipient@example.com"},
    From:     "sender@yourdomain.com",
    Subject:  "Invoice Attached",
    HTMLBody: "<p>Please find the invoice attached.</p>",
    Attachments: []types.Attachment{
        {
            Name:        "invoice.pdf",
            ContentType: "application/pdf",
            Data:        base64EncodedData,
        },
    },
}
```

#### Using Middleware
```go
// Create a logging middleware
loggingMiddleware := func(next http.RoundTripper) http.RoundTripper {
    return RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
        log.Printf("Request: %s %s", req.Method, req.URL)
        resp, err := next.RoundTrip(req)
        if resp != nil {
            log.Printf("Response: %d", resp.StatusCode)
        }
        return resp, err
    })
}

client := postal.NewClient(
    postal.WithAPIKey("your-api-key"),
    postal.WithBaseURL("https://postal.example.com"),
    postal.WithMiddleware(loggingMiddleware),
)
```

#### Error Handling
```go
result, err := client.SendMessage(ctx, message)
if err != nil {
    if postalErr, ok := err.(*types.PostalError); ok {
        if postalErr.IsRateLimit() {
            log.Println("Rate limited, retry later")
        } else if postalErr.IsUnauthorized() {
            log.Println("Invalid API key")
        } else if postalErr.IsServerError() {
            log.Println("Server error:", postalErr.Message)
        }
    }
    return err
}
```

## 🧪 Testing

The library includes a comprehensive test suite with excellent coverage:

### Test Coverage
- **Client Package**: 95.5% coverage
- **Types Package**: 100% coverage
- **Validation Package**: 98.4% coverage
- **Transport Package**: 79.5% coverage

### Running Tests
```bash
# Run all tests
make test

# Run tests with coverage report
make test-coverage

# Run specific package tests
go test ./common/validation/...

# Run benchmarks
go test -bench=. -benchmem

# Run tests with race detection
go test -race ./...
```

### Integration Testing
```bash
# Run integration tests with Docker
make integration-test

# Or manually
cd tests/integration
docker-compose up -d
go test ./tests/integration/...
docker-compose down
```

## 🛠️ Development

### Prerequisites
```bash
# Install development tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install mvdan.cc/gofumpt@latest
go install gotest.tools/gotestsum@latest
```

### Development Commands
```bash
# Format code
gofumpt -l -w .

# Run linting
golangci-lint run

# Build
make build

# Run all checks
make lint test
```

### Project Structure
```
postal-go/
├── client.go              # Main client interface
├── options.go             # Configuration options
├── common/                # Shared types and utilities
│   ├── config/            # Configuration management
│   ├── types/             # Message, Result, Error types
│   ├── utils/             # URL handling utilities
│   └── validation/        # Input validation
├── internal/              # Internal packages
│   ├── middleware/        # Built-in middleware
│   └── transport/         # HTTP transport layer
├── examples/              # Usage examples
├── tests/                 # Test suites
└── scripts/               # Development scripts
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Postal](https://github.com/postalserver/postal) - The amazing open-source mail server
- All contributors who have helped improve this library

---

Made with ❤️ by [Sachin Duhan](https://github.com/sachin-duhan)