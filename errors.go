package regru

import (
	"errors"
	"fmt"
	"net/http"
)

// Error codes worth handling programmatically. See section 2.3.3 of the API docs.
const (
	ErrCodeIPRateLimit      = "IP_EXCEEDED_ALLOWED_CONNECTION_RATE"
	ErrCodeAccountRateLimit = "ACCOUNT_EXCEEDED_ALLOWED_CONNECTION_RATE"
	ErrCodePasswordAuth     = "PASSWORD_AUTH_FAILED"
	// API access is allowed only from IPs listed in the account's API settings.
	ErrCodeAccessDeniedFromIP = "ACCESS_DENIED_FROM_IP"
)

// ErrorResponse is a top-level API error ("result": "error").
type ErrorResponse struct {
	Response *http.Response `json:"-"`
	Code     string         `json:"error_code"`
	Text     string         `json:"error_text"`
	Params   map[string]any `json:"error_params,omitempty"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("regru: %s: %s", e.Code, e.Text)
}

// ErrorCode returns the API error code of err, or "" if err is not an *ErrorResponse.
func ErrorCode(err error) string {
	var e *ErrorResponse
	if errors.As(err, &e) {
		return e.Code
	}
	return ""
}

// ItemResult is the per-item status in answers that process lists of domains or services.
// For failed items Result holds the error text (or "error" with ErrorText set).
type ItemResult struct {
	Result    string `json:"result,omitempty"`
	ErrorCode string `json:"error_code,omitempty"`
	ErrorText string `json:"error_text,omitempty"`
}

// Err returns a non-nil *ErrorResponse if the item failed.
func (r ItemResult) Err() error {
	// DNSSEC items report failures as "result":"error" without error_code.
	if r.ErrorCode == "" && r.Result != "error" {
		return nil
	}
	text := r.ErrorText
	if text == "" {
		text = r.Result
	}
	return &ErrorResponse{Code: r.ErrorCode, Text: text}
}
