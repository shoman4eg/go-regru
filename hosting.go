package regru

import "context"

// HostingService handles functions of the "hosting" category.
type HostingService service

// Nop tests access to the hosting category (hosting/nop).
// hosting/get_os_templates is listed in section 3.3 but undocumented, so it is not implemented.
func (s *HostingService) Nop(ctx context.Context) (*map[string]any, *Response, error) {
	return call[map[string]any](ctx, s.client, "hosting/nop", nil)
}
