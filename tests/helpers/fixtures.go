package helpers

import (
	"encoding/base64"

	"github.com/sachin-duhan/postal-go/common/types"
)

// MessageFixtures provides predefined message structures for testing
type MessageFixtures struct{}

// NewMessageFixtures creates a new instance of message fixtures
func NewMessageFixtures() *MessageFixtures {
	return &MessageFixtures{}
}

// BasicMessage returns a simple valid message for testing
func (f *MessageFixtures) BasicMessage() *types.Message {
	return &types.Message{
		To:       []string{"recipient@example.com"},
		From:     "sender@example.com",
		Subject:  "Test Subject",
		Body:     "This is a test message body.",
	}
}

// HTMLMessage returns a message with HTML content
func (f *MessageFixtures) HTMLMessage() *types.Message {
	return &types.Message{
		To:       []string{"recipient@example.com"},
		From:     "sender@example.com",
		Subject:  "HTML Test Subject",
		HTMLBody: "<h1>Test HTML Message</h1><p>This is a test message with HTML content.</p>",
	}
}

// ComplexMessage returns a message with all fields populated
func (f *MessageFixtures) ComplexMessage() *types.Message {
	return &types.Message{
		To:      []string{"recipient1@example.com", "recipient2@example.com"},
		CC:      []string{"cc@example.com"},
		BCC:     []string{"bcc@example.com"},
		From:    "sender@example.com",
		Sender:  "actual-sender@example.com",
		Subject: "Complex Test Subject",
		Tag:     "test-tag",
		ReplyTo: "reply@example.com",
		Body:    "Plain text content",
		HTMLBody: "<h1>HTML Content</h1><p>Rich HTML content with <strong>formatting</strong>.</p>",
		Headers: map[string]string{
			"X-Custom-Header": "custom-value",
			"X-Priority":      "high",
			"X-Category":      "test",
		},
		Attachments: []types.Attachment{
			{
				Name:        "document.pdf",
				ContentType: "application/pdf",
				Data:        base64.StdEncoding.EncodeToString([]byte("fake PDF content")),
			},
			{
				Name:        "image.png",
				ContentType: "image/png",
				Data:        base64.StdEncoding.EncodeToString([]byte("fake PNG content")),
			},
		},
	}
}

// MessageWithAttachment returns a message with a single attachment
func (f *MessageFixtures) MessageWithAttachment() *types.Message {
	return &types.Message{
		To:      []string{"recipient@example.com"},
		From:    "sender@example.com",
		Subject: "Message with Attachment",
		Body:    "Please find the attached file.",
		Attachments: []types.Attachment{
			{
				Name:        "test-file.txt",
				ContentType: "text/plain",
				Data:        base64.StdEncoding.EncodeToString([]byte("This is test file content")),
			},
		},
	}
}

// MultipleRecipientsMessage returns a message with multiple recipients
func (f *MessageFixtures) MultipleRecipientsMessage() *types.Message {
	return &types.Message{
		To: []string{
			"recipient1@example.com",
			"recipient2@example.com",
			"recipient3@example.com",
		},
		CC: []string{
			"cc1@example.com",
			"cc2@example.com",
		},
		BCC:      []string{"bcc@example.com"},
		From:     "sender@example.com",
		Subject:  "Multiple Recipients Test",
		HTMLBody: "<h1>Broadcast Message</h1><p>This message is sent to multiple recipients.</p>",
	}
}

// InvalidMessage returns a message with validation errors
func (f *MessageFixtures) InvalidMessage() *types.Message {
	return &types.Message{
		// Missing To field
		From:    "invalid-email", // Invalid email format
		Subject: "",              // Empty subject
		// Missing body content
	}
}

// EmptyMessage returns a completely empty message
func (f *MessageFixtures) EmptyMessage() *types.Message {
	return &types.Message{}
}

// RawMessageFixtures provides predefined raw message structures for testing
type RawMessageFixtures struct{}

// NewRawMessageFixtures creates a new instance of raw message fixtures
func NewRawMessageFixtures() *RawMessageFixtures {
	return &RawMessageFixtures{}
}

// BasicRawMessage returns a simple valid raw message
func (f *RawMessageFixtures) BasicRawMessage() *types.RawMessage {
	return &types.RawMessage{
		Mail: "From: sender@example.com\r\n" +
			"To: recipient@example.com\r\n" +
			"Subject: Raw Message Test\r\n" +
			"\r\n" +
			"This is a raw message body.",
		To:   []string{"recipient@example.com"},
		From: "sender@example.com",
	}
}

// ComplexRawMessage returns a raw message with headers
func (f *RawMessageFixtures) ComplexRawMessage() *types.RawMessage {
	return &types.RawMessage{
		Mail: "From: sender@example.com\r\n" +
			"To: recipient@example.com\r\n" +
			"Subject: Complex Raw Message\r\n" +
			"Content-Type: text/html\r\n" +
			"X-Custom-Header: custom-value\r\n" +
			"\r\n" +
			"<h1>HTML Raw Message</h1><p>This is HTML content in a raw message.</p>",
		To:   []string{"recipient@example.com"},
		From: "sender@example.com",
		Headers: map[string]string{
			"X-Custom-Header": "custom-value",
			"X-Priority":      "high",
		},
	}
}

// MultipartRawMessage returns a raw message with multipart content
func (f *RawMessageFixtures) MultipartRawMessage() *types.RawMessage {
	boundary := "boundary123"
	return &types.RawMessage{
		Mail: "From: sender@example.com\r\n" +
			"To: recipient@example.com\r\n" +
			"Subject: Multipart Raw Message\r\n" +
			"Content-Type: multipart/mixed; boundary=" + boundary + "\r\n" +
			"\r\n" +
			"--" + boundary + "\r\n" +
			"Content-Type: text/plain\r\n" +
			"\r\n" +
			"This is the plain text part.\r\n" +
			"--" + boundary + "\r\n" +
			"Content-Type: text/html\r\n" +
			"\r\n" +
			"<h1>HTML Part</h1><p>This is the HTML part.</p>\r\n" +
			"--" + boundary + "--\r\n",
		To:   []string{"recipient@example.com"},
		From: "sender@example.com",
	}
}

// InvalidRawMessage returns a raw message with validation errors
func (f *RawMessageFixtures) InvalidRawMessage() *types.RawMessage {
	return &types.RawMessage{
		// Missing Mail content
		To:   []string{"invalid-email"}, // Invalid email
		From: "invalid-sender",          // Invalid email
	}
}

// ResultFixtures provides predefined result structures for testing
type ResultFixtures struct{}

// NewResultFixtures creates a new instance of result fixtures
func NewResultFixtures() *ResultFixtures {
	return &ResultFixtures{}
}

// SuccessResult returns a successful result
func (f *ResultFixtures) SuccessResult() *types.Result {
	return &types.Result{
		MessageID: "msg_12345",
		Status:    "success",
		Data: map[string]interface{}{
			"queue_id":   "queue_67890",
			"priority":   "normal",
			"scheduled":  false,
		},
	}
}

// FailedResult returns a failed result with errors
func (f *ResultFixtures) FailedResult() *types.Result {
	return &types.Result{
		MessageID: "msg_12346",
		Status:    "failed",
		Errors: []string{
			"invalid recipient email",
			"subject too long",
		},
	}
}

// PartialSuccessResult returns a partial success result
func (f *ResultFixtures) PartialSuccessResult() *types.Result {
	return &types.Result{
		MessageID: "msg_12347",
		Status:    "partial_success",
		Data: map[string]interface{}{
			"sent_count":   2,
			"failed_count": 1,
		},
		Errors: []string{
			"one recipient failed validation",
		},
	}
}

// ErrorFixtures provides predefined error structures for testing
type ErrorFixtures struct{}

// NewErrorFixtures creates a new instance of error fixtures
func NewErrorFixtures() *ErrorFixtures {
	return &ErrorFixtures{}
}

// ValidationError returns a validation error
func (f *ErrorFixtures) ValidationError() *types.PostalError {
	return &types.PostalError{
		Code:       "validation_error",
		Message:    "Invalid request data",
		StatusCode: 400,
		Details: map[string]interface{}{
			"field": "email",
			"value": "invalid@",
		},
	}
}

// UnauthorizedError returns an unauthorized error
func (f *ErrorFixtures) UnauthorizedError() *types.PostalError {
	return &types.PostalError{
		Code:       "unauthorized",
		Message:    "Invalid API key",
		StatusCode: 401,
	}
}

// RateLimitError returns a rate limit error
func (f *ErrorFixtures) RateLimitError() *types.PostalError {
	return &types.PostalError{
		Code:       "rate_limit",
		Message:    "Rate limit exceeded",
		StatusCode: 429,
		Details: map[string]interface{}{
			"limit":     100,
			"remaining": 0,
			"reset_at":  "2023-12-01T10:00:00Z",
		},
	}
}

// ServerError returns a server error
func (f *ErrorFixtures) ServerError() *types.PostalError {
	return &types.PostalError{
		Code:       "server_error",
		Message:    "Internal server error",
		StatusCode: 500,
	}
}

// AttachmentFixtures provides predefined attachment structures for testing
type AttachmentFixtures struct{}

// NewAttachmentFixtures creates a new instance of attachment fixtures
func NewAttachmentFixtures() *AttachmentFixtures {
	return &AttachmentFixtures{}
}

// TextAttachment returns a simple text attachment
func (f *AttachmentFixtures) TextAttachment() types.Attachment {
	content := "This is a test text file content.\nLine 2\nLine 3"
	return types.Attachment{
		Name:        "test.txt",
		ContentType: "text/plain",
		Data:        base64.StdEncoding.EncodeToString([]byte(content)),
	}
}

// PDFAttachment returns a fake PDF attachment
func (f *AttachmentFixtures) PDFAttachment() types.Attachment {
	content := "%PDF-1.4\nFake PDF content for testing"
	return types.Attachment{
		Name:        "document.pdf",
		ContentType: "application/pdf",
		Data:        base64.StdEncoding.EncodeToString([]byte(content)),
	}
}

// ImageAttachment returns a fake image attachment
func (f *AttachmentFixtures) ImageAttachment() types.Attachment {
	content := "PNG\r\nFake PNG content for testing"
	return types.Attachment{
		Name:        "image.png",
		ContentType: "image/png",
		Data:        base64.StdEncoding.EncodeToString([]byte(content)),
	}
}

// LargeAttachment returns a large attachment for testing size limits
func (f *AttachmentFixtures) LargeAttachment() types.Attachment {
	// Create a large content (1MB of 'A' characters)
	content := make([]byte, 1024*1024)
	for i := range content {
		content[i] = 'A'
	}
	
	return types.Attachment{
		Name:        "large-file.txt",
		ContentType: "text/plain",
		Data:        base64.StdEncoding.EncodeToString(content),
	}
}

// InvalidAttachment returns an attachment with missing fields
func (f *AttachmentFixtures) InvalidAttachment() types.Attachment {
	return types.Attachment{
		// Missing Name, ContentType, and Data
	}
}

// EmailFixtures provides various email addresses for testing
type EmailFixtures struct{}

// NewEmailFixtures creates a new instance of email fixtures
func NewEmailFixtures() *EmailFixtures {
	return &EmailFixtures{}
}

// ValidEmails returns a list of valid email addresses
func (f *EmailFixtures) ValidEmails() []string {
	return []string{
		"user@example.com",
		"user.name@example.com",
		"user+tag@example.com",
		"user_name@example.com",
		"user123@example.com",
		"user@subdomain.example.com",
		"user@example.co.uk",
		"123456@example.com",
		"test.email+tag@long-domain-name.co.uk",
	}
}

// InvalidEmails returns a list of invalid email addresses
func (f *EmailFixtures) InvalidEmails() []string {
	return []string{
		"",
		"plaintext",
		"@example.com",
		"user@",
		"user@example",
		"user@@example.com",
		"user@example..com",
		"user example@example.com",
		"user",
		"@",
		"user@.com",
		"user@com",
		"user@example.",
		"user@.example.com",
	}
}

// SpecialCaseEmails returns edge case email addresses
func (f *EmailFixtures) SpecialCaseEmails() []string {
	return []string{
		"a@b.co",                    // Minimal valid email
		"very.long.email.address@very.long.domain.name.example.com", // Long email
		"user+tag+another@example.com", // Multiple plus signs
		"user.with.many.dots@example.com", // Multiple dots
	}
}

// GetAllFixtures returns all fixture instances
func GetAllFixtures() (
	*MessageFixtures,
	*RawMessageFixtures,
	*ResultFixtures,
	*ErrorFixtures,
	*AttachmentFixtures,
	*EmailFixtures,
) {
	return NewMessageFixtures(),
		NewRawMessageFixtures(),
		NewResultFixtures(),
		NewErrorFixtures(),
		NewAttachmentFixtures(),
		NewEmailFixtures()
}