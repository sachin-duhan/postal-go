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

	// Check for spaces
	if strings.Contains(email, " ") {
		return false
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	localPart := parts[0]
	domain := parts[1]

	if len(localPart) == 0 || len(domain) == 0 {
		return false
	}

	// Domain must contain at least one dot and not start/end with dot
	if !strings.Contains(domain, ".") {
		return false
	}

	// Domain cannot start or end with dot
	if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
		return false
	}

	// Domain cannot have consecutive dots
	if strings.Contains(domain, "..") {
		return false
	}

	// Domain parts cannot be empty
	domainParts := strings.Split(domain, ".")
	for _, part := range domainParts {
		if len(part) == 0 {
			return false
		}
	}

	return true
}
