package validation

import (
	"strings"
	"testing"

	"github.com/sachin-duhan/postal-go/common/types"
)

func TestValidateMessage(t *testing.T) {
	tests := []struct {
		name        string
		message     *types.Message
		wantErr     bool
		errContains []string
	}{
		{
			name: "valid message with HTML body",
			message: &types.Message{
				To:       []string{"recipient@example.com"},
				From:     "sender@example.com",
				Subject:  "Test Subject",
				HTMLBody: "<p>Test Body</p>",
			},
			wantErr: false,
		},
		{
			name: "valid message with plain body",
			message: &types.Message{
				To:      []string{"recipient@example.com"},
				From:    "sender@example.com",
				Subject: "Test Subject",
				Body:    "Test Body",
			},
			wantErr: false,
		},
		{
			name: "valid message with both bodies",
			message: &types.Message{
				To:       []string{"recipient@example.com"},
				From:     "sender@example.com",
				Subject:  "Test Subject",
				Body:     "Test Body",
				HTMLBody: "<p>Test Body</p>",
			},
			wantErr: false,
		},
		{
			name: "valid message with multiple recipients",
			message: &types.Message{
				To:       []string{"recipient1@example.com", "recipient2@example.com"},
				CC:       []string{"cc@example.com"},
				BCC:      []string{"bcc@example.com"},
				From:     "sender@example.com",
				Subject:  "Test Subject",
				HTMLBody: "Test Body",
			},
			wantErr: false,
		},
		{
			name: "valid message with attachments",
			message: &types.Message{
				To:      []string{"recipient@example.com"},
				From:    "sender@example.com",
				Subject: "Test Subject",
				Body:    "Test Body",
				Attachments: []types.Attachment{
					{
						Name:        "test.txt",
						ContentType: "text/plain",
						Data:        "base64data",
					},
				},
			},
			wantErr: false,
		},
		{
			name:        "missing all required fields",
			message:     &types.Message{},
			wantErr:     true,
			errContains: []string{"recipient (To) is required", "sender (From) is required", "subject is required", "either plain body or HTML body is required"},
		},
		{
			name: "missing recipient",
			message: &types.Message{
				From:    "sender@example.com",
				Subject: "Test Subject",
				Body:    "Test Body",
			},
			wantErr:     true,
			errContains: []string{"recipient (To) is required"},
		},
		{
			name: "empty recipient list",
			message: &types.Message{
				To:      []string{},
				From:    "sender@example.com",
				Subject: "Test Subject",
				Body:    "Test Body",
			},
			wantErr:     true,
			errContains: []string{"recipient (To) is required"},
		},
		{
			name: "missing sender",
			message: &types.Message{
				To:      []string{"recipient@example.com"},
				Subject: "Test Subject",
				Body:    "Test Body",
			},
			wantErr:     true,
			errContains: []string{"sender (From) is required"},
		},
		{
			name: "missing subject",
			message: &types.Message{
				To:   []string{"recipient@example.com"},
				From: "sender@example.com",
				Body: "Test Body",
			},
			wantErr:     true,
			errContains: []string{"subject is required"},
		},
		{
			name: "missing body",
			message: &types.Message{
				To:      []string{"recipient@example.com"},
				From:    "sender@example.com",
				Subject: "Test Subject",
			},
			wantErr:     true,
			errContains: []string{"either plain body or HTML body is required"},
		},
		{
			name: "invalid recipient email",
			message: &types.Message{
				To:      []string{"invalid-email"},
				From:    "sender@example.com",
				Subject: "Test Subject",
				Body:    "Test Body",
			},
			wantErr:     true,
			errContains: []string{"invalid recipient email: invalid-email"},
		},
		{
			name: "multiple invalid emails",
			message: &types.Message{
				To:      []string{"invalid1", "valid@example.com", "invalid2"},
				From:    "sender@example.com",
				Subject: "Test Subject",
				Body:    "Test Body",
			},
			wantErr:     true,
			errContains: []string{"invalid recipient email: invalid1", "invalid recipient email: invalid2"},
		},
		{
			name: "invalid sender email",
			message: &types.Message{
				To:      []string{"recipient@example.com"},
				From:    "invalid-sender",
				Subject: "Test Subject",
				Body:    "Test Body",
			},
			wantErr:     true,
			errContains: []string{"invalid sender email: invalid-sender"},
		},
		{
			name: "attachment missing name",
			message: &types.Message{
				To:      []string{"recipient@example.com"},
				From:    "sender@example.com",
				Subject: "Test Subject",
				Body:    "Test Body",
				Attachments: []types.Attachment{
					{
						ContentType: "text/plain",
						Data:        "base64data",
					},
				},
			},
			wantErr:     true,
			errContains: []string{"attachment name is required"},
		},
		{
			name: "attachment missing content type",
			message: &types.Message{
				To:      []string{"recipient@example.com"},
				From:    "sender@example.com",
				Subject: "Test Subject",
				Body:    "Test Body",
				Attachments: []types.Attachment{
					{
						Name: "test.txt",
						Data: "base64data",
					},
				},
			},
			wantErr:     true,
			errContains: []string{"attachment content type is required"},
		},
		{
			name: "attachment missing data",
			message: &types.Message{
				To:      []string{"recipient@example.com"},
				From:    "sender@example.com",
				Subject: "Test Subject",
				Body:    "Test Body",
				Attachments: []types.Attachment{
					{
						Name:        "test.txt",
						ContentType: "text/plain",
					},
				},
			},
			wantErr:     true,
			errContains: []string{"attachment data is required"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMessage(tt.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && len(tt.errContains) > 0 {
				errStr := err.Error()
				for _, contains := range tt.errContains {
					if !strings.Contains(errStr, contains) {
						t.Errorf("ValidateMessage() error = %v, want error containing %v", err, contains)
					}
				}
			}
		})
	}
}

func TestValidateRawMessage(t *testing.T) {
	tests := []struct {
		name        string
		message     *types.RawMessage
		wantErr     bool
		errContains []string
	}{
		{
			name: "valid raw message",
			message: &types.RawMessage{
				Mail: "From: sender@example.com\r\nTo: recipient@example.com\r\nSubject: Test\r\n\r\nBody",
				To:   []string{"recipient@example.com"},
				From: "sender@example.com",
			},
			wantErr: false,
		},
		{
			name: "valid raw message with headers",
			message: &types.RawMessage{
				Mail: "From: sender@example.com\r\nTo: recipient@example.com\r\nSubject: Test\r\n\r\nBody",
				To:   []string{"recipient@example.com"},
				From: "sender@example.com",
				Headers: map[string]string{
					"X-Custom": "value",
				},
			},
			wantErr: false,
		},
		{
			name:        "missing all fields",
			message:     &types.RawMessage{},
			wantErr:     true,
			errContains: []string{"raw mail content is required", "recipient (To) is required", "sender (From) is required"},
		},
		{
			name: "missing mail content",
			message: &types.RawMessage{
				To:   []string{"recipient@example.com"},
				From: "sender@example.com",
			},
			wantErr:     true,
			errContains: []string{"raw mail content is required"},
		},
		{
			name: "missing recipient",
			message: &types.RawMessage{
				Mail: "From: sender@example.com\r\nSubject: Test\r\n\r\nBody",
				From: "sender@example.com",
			},
			wantErr:     true,
			errContains: []string{"recipient (To) is required"},
		},
		{
			name: "missing sender",
			message: &types.RawMessage{
				Mail: "To: recipient@example.com\r\nSubject: Test\r\n\r\nBody",
				To:   []string{"recipient@example.com"},
			},
			wantErr:     true,
			errContains: []string{"sender (From) is required"},
		},
		{
			name: "invalid recipient email",
			message: &types.RawMessage{
				Mail: "From: sender@example.com\r\nTo: invalid\r\nSubject: Test\r\n\r\nBody",
				To:   []string{"invalid"},
				From: "sender@example.com",
			},
			wantErr:     true,
			errContains: []string{"invalid recipient email: invalid"},
		},
		{
			name: "invalid sender email",
			message: &types.RawMessage{
				Mail: "From: invalid\r\nTo: recipient@example.com\r\nSubject: Test\r\n\r\nBody",
				To:   []string{"recipient@example.com"},
				From: "invalid",
			},
			wantErr:     true,
			errContains: []string{"invalid sender email: invalid"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRawMessage(tt.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRawMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && len(tt.errContains) > 0 {
				errStr := err.Error()
				for _, contains := range tt.errContains {
					if !strings.Contains(errStr, contains) {
						t.Errorf("ValidateRawMessage() error = %v, want error containing %v", err, contains)
					}
				}
			}
		})
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		// Valid emails
		{"user@example.com", true},
		{"user.name@example.com", true},
		{"user+tag@example.com", true},
		{"user_name@example.com", true},
		{"user123@example.com", true},
		{"user@subdomain.example.com", true},
		{"user@example.co.uk", true},
		{"1234567890@example.com", true},
		{"user@123.123.123.123", true},

		// Invalid emails
		{"", false},
		{"plaintext", false},
		{"@example.com", false},
		{"user@", false},
		{"user@example", false},
		{"user@@example.com", false},
		{"user@example..com", false},
		{"user example@example.com", false},
		{"user", false},
		{"@", false},
		{"user@.com", false},
		{"user@com", false},
		{"user@example.", false},
		{"user@.example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			if got := isValidEmail(tt.email); got != tt.valid {
				t.Errorf("isValidEmail(%q) = %v, want %v", tt.email, got, tt.valid)
			}
		})
	}
}

func TestValidationErrorAggregation(t *testing.T) {
	// Test that multiple validation errors are properly aggregated
	message := &types.Message{
		To:      []string{"invalid1", "invalid2"},
		From:    "invalid-sender",
		Subject: "", // missing
		// missing body
		Attachments: []types.Attachment{
			{
				// missing all fields
			},
		},
	}

	err := ValidateMessage(message)
	if err == nil {
		t.Fatal("expected validation error")
	}

	expectedErrors := []string{
		"subject is required",
		"either plain body or HTML body is required",
		"invalid recipient email: invalid1",
		"invalid recipient email: invalid2",
		"invalid sender email: invalid-sender",
		"attachment name is required",
		"attachment content type is required",
		"attachment data is required",
	}

	errStr := err.Error()
	for _, expected := range expectedErrors {
		if !strings.Contains(errStr, expected) {
			t.Errorf("error %q does not contain expected message %q", errStr, expected)
		}
	}
}

func BenchmarkValidateMessage(b *testing.B) {
	message := &types.Message{
		To:      []string{"recipient1@example.com", "recipient2@example.com", "recipient3@example.com"},
		CC:      []string{"cc1@example.com", "cc2@example.com"},
		BCC:     []string{"bcc1@example.com", "bcc2@example.com"},
		From:    "sender@example.com",
		Subject: "Test Subject",
		Body:    "Test Body",
		HTMLBody: "<p>Test Body</p>",
		Attachments: []types.Attachment{
			{
				Name:        "test1.txt",
				ContentType: "text/plain",
				Data:        "base64data1",
			},
			{
				Name:        "test2.pdf",
				ContentType: "application/pdf",
				Data:        "base64data2",
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ValidateMessage(message)
	}
}

func BenchmarkIsValidEmail(b *testing.B) {
	emails := []string{
		"user@example.com",
		"user.name+tag@subdomain.example.co.uk",
		"invalid-email",
		"another.user@domain.com",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, email := range emails {
			_ = isValidEmail(email)
		}
	}
}