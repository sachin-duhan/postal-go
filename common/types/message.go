package types

// Message represents an email message with builder pattern
type Message struct {
	To          []string          `json:"to"`
	CC          []string          `json:"cc"`
	BCC         []string          `json:"bcc"`
	From        string            `json:"from"`
	Sender      string            `json:"sender,omitempty"`
	Subject     string            `json:"subject"`
	Tag         string            `json:"tag,omitempty"`
	ReplyTo     string            `json:"reply_to,omitempty"`
	Body        string            `json:"plain_body,omitempty"`
	HTMLBody    string            `json:"html_body,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	Attachments []Attachment      `json:"attachments,omitempty"`
}

// Attachment represents an email attachment
type Attachment struct {
	Name        string `json:"name"`
	ContentType string `json:"content_type"`
	Data        string `json:"data"` // Base64 encoded
}

// RawMessage represents a pre-formatted email message
type RawMessage struct {
	Mail    string            `json:"mail"`
	To      []string          `json:"to"`
	From    string            `json:"from"`
	Headers map[string]string `json:"headers,omitempty"`
}
