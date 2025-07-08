package types

import (
	"encoding/json"
	"testing"
)

func TestMessageJSONMarshaling(t *testing.T) {
	tests := []struct {
		name    string
		message *Message
	}{
		{
			name: "basic message",
			message: &Message{
				To:      []string{"recipient@example.com"},
				From:    "sender@example.com",
				Subject: "Test Subject",
				Body:    "Plain text body",
			},
		},
		{
			name: "message with HTML body",
			message: &Message{
				To:       []string{"recipient@example.com"},
				From:     "sender@example.com",
				Subject:  "Test Subject",
				HTMLBody: "<h1>HTML Body</h1>",
			},
		},
		{
			name: "message with both bodies",
			message: &Message{
				To:       []string{"recipient@example.com"},
				From:     "sender@example.com",
				Subject:  "Test Subject",
				Body:     "Plain text body",
				HTMLBody: "<h1>HTML Body</h1>",
			},
		},
		{
			name: "complex message with all fields",
			message: &Message{
				To:      []string{"recipient1@example.com", "recipient2@example.com"},
				CC:      []string{"cc@example.com"},
				BCC:     []string{"bcc@example.com"},
				From:    "sender@example.com",
				Sender:  "actual-sender@example.com",
				Subject: "Complex Test Subject",
				Tag:     "test-tag",
				ReplyTo: "reply@example.com",
				Body:    "Plain text body",
				HTMLBody: "<h1>HTML Body</h1><p>Content</p>",
				Headers: map[string]string{
					"X-Custom-Header": "custom-value",
					"X-Priority":      "high",
				},
				Attachments: []Attachment{
					{
						Name:        "document.pdf",
						ContentType: "application/pdf",
						Data:        "base64encodeddata",
					},
					{
						Name:        "image.png",
						ContentType: "image/png",
						Data:        "anotherbas64string",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON
			jsonData, err := json.Marshal(tt.message)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}

			// Unmarshal back to Message
			var unmarshaled Message
			err = json.Unmarshal(jsonData, &unmarshaled)
			if err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}

			// Verify all fields
			if !slicesEqual(unmarshaled.To, tt.message.To) {
				t.Errorf("To = %v, want %v", unmarshaled.To, tt.message.To)
			}
			if !slicesEqual(unmarshaled.CC, tt.message.CC) {
				t.Errorf("CC = %v, want %v", unmarshaled.CC, tt.message.CC)
			}
			if !slicesEqual(unmarshaled.BCC, tt.message.BCC) {
				t.Errorf("BCC = %v, want %v", unmarshaled.BCC, tt.message.BCC)
			}
			if unmarshaled.From != tt.message.From {
				t.Errorf("From = %v, want %v", unmarshaled.From, tt.message.From)
			}
			if unmarshaled.Sender != tt.message.Sender {
				t.Errorf("Sender = %v, want %v", unmarshaled.Sender, tt.message.Sender)
			}
			if unmarshaled.Subject != tt.message.Subject {
				t.Errorf("Subject = %v, want %v", unmarshaled.Subject, tt.message.Subject)
			}
			if unmarshaled.Tag != tt.message.Tag {
				t.Errorf("Tag = %v, want %v", unmarshaled.Tag, tt.message.Tag)
			}
			if unmarshaled.ReplyTo != tt.message.ReplyTo {
				t.Errorf("ReplyTo = %v, want %v", unmarshaled.ReplyTo, tt.message.ReplyTo)
			}
			if unmarshaled.Body != tt.message.Body {
				t.Errorf("Body = %v, want %v", unmarshaled.Body, tt.message.Body)
			}
			if unmarshaled.HTMLBody != tt.message.HTMLBody {
				t.Errorf("HTMLBody = %v, want %v", unmarshaled.HTMLBody, tt.message.HTMLBody)
			}

			// Verify Headers map
			if !mapsEqual(unmarshaled.Headers, tt.message.Headers) {
				t.Errorf("Headers = %v, want %v", unmarshaled.Headers, tt.message.Headers)
			}

			// Verify Attachments
			if !attachmentsEqual(unmarshaled.Attachments, tt.message.Attachments) {
				t.Errorf("Attachments = %v, want %v", unmarshaled.Attachments, tt.message.Attachments)
			}
		})
	}
}

func TestAttachmentJSONMarshaling(t *testing.T) {
	attachment := Attachment{
		Name:        "test-file.txt",
		ContentType: "text/plain",
		Data:        "VGVzdCBjb250ZW50", // Base64 for "Test content"
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(attachment)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	// Unmarshal back to Attachment
	var unmarshaled Attachment
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if unmarshaled.Name != attachment.Name {
		t.Errorf("Name = %v, want %v", unmarshaled.Name, attachment.Name)
	}
	if unmarshaled.ContentType != attachment.ContentType {
		t.Errorf("ContentType = %v, want %v", unmarshaled.ContentType, attachment.ContentType)
	}
	if unmarshaled.Data != attachment.Data {
		t.Errorf("Data = %v, want %v", unmarshaled.Data, attachment.Data)
	}
}

func TestRawMessageJSONMarshaling(t *testing.T) {
	tests := []struct {
		name    string
		rawMsg  *RawMessage
	}{
		{
			name: "basic raw message",
			rawMsg: &RawMessage{
				Mail: "From: sender@example.com\r\nTo: recipient@example.com\r\nSubject: Test\r\n\r\nBody content",
				To:   []string{"recipient@example.com"},
				From: "sender@example.com",
			},
		},
		{
			name: "raw message with headers",
			rawMsg: &RawMessage{
				Mail: "From: sender@example.com\r\nTo: recipient@example.com\r\nSubject: Test\r\n\r\nBody content",
				To:   []string{"recipient1@example.com", "recipient2@example.com"},
				From: "sender@example.com",
				Headers: map[string]string{
					"X-Custom":   "value",
					"X-Priority": "high",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON
			jsonData, err := json.Marshal(tt.rawMsg)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}

			// Unmarshal back to RawMessage
			var unmarshaled RawMessage
			err = json.Unmarshal(jsonData, &unmarshaled)
			if err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}

			if unmarshaled.Mail != tt.rawMsg.Mail {
				t.Errorf("Mail = %v, want %v", unmarshaled.Mail, tt.rawMsg.Mail)
			}
			if !slicesEqual(unmarshaled.To, tt.rawMsg.To) {
				t.Errorf("To = %v, want %v", unmarshaled.To, tt.rawMsg.To)
			}
			if unmarshaled.From != tt.rawMsg.From {
				t.Errorf("From = %v, want %v", unmarshaled.From, tt.rawMsg.From)
			}
			if !mapsEqual(unmarshaled.Headers, tt.rawMsg.Headers) {
				t.Errorf("Headers = %v, want %v", unmarshaled.Headers, tt.rawMsg.Headers)
			}
		})
	}
}

func TestJSONOmitEmpty(t *testing.T) {
	// Test that empty/nil fields are omitted from JSON
	message := &Message{
		To:      []string{"recipient@example.com"},
		From:    "sender@example.com",
		Subject: "Test Subject",
		Body:    "Test body",
		// All other fields are empty and should be omitted
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	jsonStr := string(jsonData)

	// Should contain required fields
	requiredFields := []string{"to", "from", "subject", "plain_body"}
	for _, field := range requiredFields {
		if !messageContains(jsonStr, field) {
			t.Errorf("JSON should contain %s field", field)
		}
	}

	// Should NOT contain omitted fields (check for JSON field names with quotes)
	omittedFields := []string{"\"sender\":", "\"tag\":", "\"reply_to\":", "\"html_body\":", "\"headers\":", "\"attachments\":"}
	for _, field := range omittedFields {
		if messageContains(jsonStr, field) {
			t.Errorf("JSON should not contain %s field when empty", field)
		}
	}
}

func TestJSONFieldNames(t *testing.T) {
	// Test that JSON field names match expected API format
	message := &Message{
		To:       []string{"recipient@example.com"},
		From:     "sender@example.com",
		Subject:  "Test Subject",
		Body:     "Plain text",
		HTMLBody: "<p>HTML text</p>",
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	jsonStr := string(jsonData)

	// Test specific field names that might be different from Go field names
	expectedMappings := map[string]string{
		"Body":     "plain_body",
		"HTMLBody": "html_body",
		"ReplyTo":  "reply_to",
	}

	for goField, _ := range expectedMappings {
		if messageContains(jsonStr, goField) {
			t.Errorf("JSON should not contain Go field name %s", goField)
		}
		// We only check for the JSON field if the Go field has a value
		// since we're testing the field naming, not the omitEmpty behavior
	}

	// Should contain the correct JSON field names
	if !messageContains(jsonStr, "plain_body") {
		t.Error("JSON should contain plain_body field")
	}
	if !messageContains(jsonStr, "html_body") {
		t.Error("JSON should contain html_body field")
	}
}

func BenchmarkMessageJSONMarshal(b *testing.B) {
	message := &Message{
		To:      []string{"recipient1@example.com", "recipient2@example.com", "recipient3@example.com"},
		CC:      []string{"cc1@example.com", "cc2@example.com"},
		BCC:     []string{"bcc1@example.com"},
		From:    "sender@example.com",
		Subject: "Benchmark Test Subject",
		Body:    "This is a benchmark test message body",
		HTMLBody: "<h1>Benchmark Test</h1><p>This is a benchmark test message body</p>",
		Headers: map[string]string{
			"X-Custom-Header": "custom-value",
			"X-Priority":      "high",
		},
		Attachments: []Attachment{
			{
				Name:        "document.pdf",
				ContentType: "application/pdf",
				Data:        "base64encodeddata",
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(message)
		if err != nil {
			b.Fatalf("json.Marshal() error = %v", err)
		}
	}
}

func BenchmarkRawMessageJSONMarshal(b *testing.B) {
	rawMsg := &RawMessage{
		Mail: "From: sender@example.com\r\nTo: recipient@example.com\r\nSubject: Benchmark Test\r\n\r\nThis is a benchmark test message body",
		To:   []string{"recipient@example.com"},
		From: "sender@example.com",
		Headers: map[string]string{
			"X-Custom": "value",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(rawMsg)
		if err != nil {
			b.Fatalf("json.Marshal() error = %v", err)
		}
	}
}

// Helper functions for comparison
func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func mapsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

func attachmentsEqual(a, b []Attachment) bool {
	if len(a) != len(b) {
		return false
	}
	for i, att := range a {
		if att.Name != b[i].Name || att.ContentType != b[i].ContentType || att.Data != b[i].Data {
			return false
		}
	}
	return true
}

func messageContains(s, substr string) bool {
	return len(s) >= len(substr) && substr != "" && 
		(s == substr || len(s) > len(substr) && messageContainsHelper(s, substr))
}

func messageContainsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}