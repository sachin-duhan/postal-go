package types

// Result represents the response from the Postal API
type Result struct {
	MessageID string                 `json:"message_id"`
	Status    string                 `json:"status"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Errors    []string               `json:"errors,omitempty"`
}

// Success returns true if the API call was successful
func (r *Result) Success() bool {
	return r.Status == "success"
}

// Failed returns true if the API call failed
func (r *Result) Failed() bool {
	return !r.Success()
}
