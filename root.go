package regru

import "context"

// NopAnswer is returned by Nop and ResellerNop.
type NopAnswer struct {
	Login  string  `json:"login"`
	UserID FlexInt `json:"user_id"`
}

// Nop checks access to the API and returns the login.
func (c *Client) Nop(ctx context.Context) (*NopAnswer, *Response, error) {
	return call[NopAnswer](ctx, c, "nop", nil)
}

// ResellerNop checks reseller (partner) access to the API.
func (c *Client) ResellerNop(ctx context.Context) (*NopAnswer, *Response, error) {
	return call[NopAnswer](ctx, c, "reseller_nop", nil)
}

// UserIDAnswer is returned by GetUserID.
type UserIDAnswer struct {
	UserID FlexInt `json:"user_id"`
}

// GetUserID returns the ID of the authenticated user.
func (c *Client) GetUserID(ctx context.Context) (*UserIDAnswer, *Response, error) {
	return call[UserIDAnswer](ctx, c, "get_user_id", nil)
}

// ServiceIDAnswer is returned by GetServiceID.
type ServiceIDAnswer struct {
	ServiceID FlexInt `json:"service_id"`
}

// GetServiceID returns service_id of a domain or service.
func (c *Client) GetServiceID(ctx context.Context, opts *ServiceIdent) (*ServiceIDAnswer, *Response, error) {
	return call[ServiceIDAnswer](ctx, c, "get_service_id", opts)
}
