# 📮 Postal-Go

A powerful and flexible Go client library for the Postal email server. This library provides a simple, middleware-based approach to interact with Postal's API, making it easy to send emails with advanced features and customizations.

## Features

- 🚀 Simple and intuitive API
- 🔌 Middleware support for custom request handling
- 📧 Support for both simple and raw email messages
- 🔄 Built-in retry mechanism
- ⚡ Concurrent request handling
- ✅ Request validation
- 🔍 Debug mode
- 📝 Comprehensive logging
- 🎨 HTML email support
- 🔧 Highly configurable

## Installation

```bash
go get github.com/sachin-duhan/postal-go
```

## Quick Start

```go
package main

import (
    "context"
    "log"
    client "github.com/sachin-duhan/postal-go"
    "github.com/sachin-duhan/postal-go/common/types"
)

func main() {
    // Initialize client
    postalClient, err := client.NewClient(
        "https://postal.example.com",
        "your-api-key",
    )
    if err != nil {
        log.Fatal(err)
    }

    // Create a message
    message := &types.Message{
        To:      []string{"recipient@example.com"},
        From:    "sender@yourdomain.com",
        Subject: "Hello from Postal-Go!",
        Body:    "This is a test email.",
        HTMLBody: `<h1>Hello!</h1><p>This is a test email.</p>`,
    }

    // Send the message
    result, err := postalClient.SendMessage(context.Background(), message)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Message sent: %+v", result)
}
```

## Documentation

Check out our [examples](./examples) directory for detailed examples:
- [Simple Example](./examples/simple/main.go): Basic usage
- [Advanced Example](./examples/advanced/main.go): Advanced features

## Features

### Configuration
```go
config := &client.Config{
    Timeout:        30 * time.Second,
    MaxRetries:     3,
    RetryInterval:  time.Second,
    MaxConcurrency: 10,
    Debug:         true,
}
```

### Middleware Support
- Logging middleware
- Retry middleware
- Header middleware
- Custom middleware support

### Message Types
- Standard messages with HTML support
- Raw messages with MIME support
- Custom headers
- Multiple recipients

### Error Handling
- Validation errors
- API errors
- Retry logic
- Context support

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.