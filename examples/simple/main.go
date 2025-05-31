package main

import (
	"context"
	"log"
	"time"

	client "github.com/sachin-duhan/postal-go"
	"github.com/sachin-duhan/postal-go/common/types"
)

func main() {
	// Initialize the client
	postalClient, err := client.NewClient(
		"https://postal.example.com", // Replace with your Postal server URL
		"your-api-key",               // Replace with your API key
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Configure the client with custom settings
	config := &client.Config{
		Timeout:        10 * time.Second,
		MaxRetries:     3,
		RetryInterval:  time.Second,
		MaxConcurrency: 5,
		Debug:          true,
	}
	postalClient = postalClient.WithConfig(config)

	// Create a simple message
	message := &types.Message{
		To:      []string{"recipient@example.com"},
		From:    "sender@yourdomain.com",
		Subject: "Test Email from Postal-Go",
		Body:    "This is a test email sent using the Postal-Go client library.",
		HTMLBody: `
			<html>
				<body>
					<h1>Test Email</h1>
					<p>This is a test email sent using the <strong>Postal-Go</strong> client library.</p>
				</body>
			</html>
		`,
	}

	// Send the message
	ctx := context.Background()
	result, err := postalClient.SendMessage(ctx, message)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	log.Printf("Message sent successfully! Result: %+v", result)
}
