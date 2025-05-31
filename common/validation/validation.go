package validation

import (
	"fmt"
	"strings"

	"github.com/sachin-duhan/postal-go/common/types"
)

// ValidateMessage validates a message before sending
func ValidateMessage(msg *types.Message) error {
	var errors []string

	// Required fields
	if len(msg.To) == 0 {
		errors = append(errors, "recipient (To) is required")
	}

	if msg.From == "" {
		errors = append(errors, "sender (From) is required")
	}

	if msg.Subject == "" {
		errors = append(errors, "subject is required")
	}

	// Content validation
	if msg.Body == "" && msg.HTMLBody == "" {
		errors = append(errors, "either plain body or HTML body is required")
	}

	// Email format validation
	for _, to := range msg.To {
		if !isValidEmail(to) {
			errors = append(errors, fmt.Sprintf("invalid recipient email: %s", to))
		}
	}

	if !isValidEmail(msg.From) {
		errors = append(errors, fmt.Sprintf("invalid sender email: %s", msg.From))
	}

	// Attachment validation
	for _, att := range msg.Attachments {
		if att.Name == "" {
			errors = append(errors, "attachment name is required")
		}
		if att.ContentType == "" {
			errors = append(errors, "attachment content type is required")
		}
		if att.Data == "" {
			errors = append(errors, "attachment data is required")
		}
	}

	if len(errors) > 0 {
		return types.NewPostalError("validation_error", strings.Join(errors, "; "), 400)
	}

	return nil
}

// ValidateRawMessage validates a raw message before sending
func ValidateRawMessage(msg *types.RawMessage) error {
	var errors []string

	if msg.Mail == "" {
		errors = append(errors, "raw mail content is required")
	}

	if len(msg.To) == 0 {
		errors = append(errors, "recipient (To) is required")
	}

	if msg.From == "" {
		errors = append(errors, "sender (From) is required")
	}

	// Email format validation
	for _, to := range msg.To {
		if !isValidEmail(to) {
			errors = append(errors, fmt.Sprintf("invalid recipient email: %s", to))
		}
	}

	if !isValidEmail(msg.From) {
		errors = append(errors, fmt.Sprintf("invalid sender email: %s", msg.From))
	}

	if len(errors) > 0 {
		return types.NewPostalError("validation_error", strings.Join(errors, "; "), 400)
	}

	return nil
}

// isValidEmail performs basic email format validation
func isValidEmail(email string) bool {
	// Basic email validation
	if email == "" {
		return false
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	if len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}

	return strings.Contains(parts[1], ".")
}
