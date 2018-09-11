package errors

// APIError represents an error that can be sent in an error response.
type APIError struct {
	Status           int         `json:"-"`
	ErrorCode        string      `json:"error_code"`
	Message          string      `json:"message"`
	DeveloperMessage string      `json:"developer_message,omitempty"`
	Details          interface{} `json:"details,omitempty"`
}

// Error returns the error message.
func (e APIError) Error() string {
	return e.Message
}

// StatusCode returns the HTTP status code.
func (e APIError) StatusCode() int {
	return e.Status
}
