package regru

import "context"

// DNSSECService handles functions of the "dnssec" category.
type DNSSECService service

// DNSSECDomain identifies a domain in a DNSSEC request.
type DNSSECDomain struct {
	ServiceIdent
	Records []string `json:"records,omitempty"` // add_keys only: DNSKEY or DS records
}

// DNSSECOptions identify one domain (embedded fields) or a list of domains (Domains).
type DNSSECOptions struct {
	DNSSECDomain
	Domains []DNSSECDomain `json:"domains,omitempty"`
}

// DNSSECDomainResult is a per-domain entry of DNSSECAnswer.
type DNSSECDomainResult struct {
	ServiceItem
	Status string `json:"status,omitempty"` // get_status: enabled, disabled, updating
	// get_records: DS/DNSKEY for REG.RU name servers, free-form (alg, digest, ...) otherwise.
	Records []map[string]FlexString `json:"records,omitempty"`
}

// DNSSECAnswer is returned by all dnssec functions.
type DNSSECAnswer struct {
	Domains []DNSSECDomainResult `json:"domains"`
}

// Nop checks whether DNSSEC can be managed for domains (dnssec/nop).
func (s *DNSSECService) Nop(ctx context.Context, opts *DNSSECOptions) (*DNSSECAnswer, *Response, error) {
	return call[DNSSECAnswer](ctx, s.client, "dnssec/nop", opts)
}

// GetStatus returns the DNSSEC status of domains (dnssec/get_status).
func (s *DNSSECService) GetStatus(ctx context.Context, opts *DNSSECOptions) (*DNSSECAnswer, *Response, error) {
	return call[DNSSECAnswer](ctx, s.client, "dnssec/get_status", opts)
}

// Enable enables DNSSEC for domains on REG.RU name servers (dnssec/enable).
func (s *DNSSECService) Enable(ctx context.Context, opts *DNSSECOptions) (*DNSSECAnswer, *Response, error) {
	return call[DNSSECAnswer](ctx, s.client, "dnssec/enable", opts)
}

// Disable disables DNSSEC for domains (dnssec/disable).
func (s *DNSSECService) Disable(ctx context.Context, opts *DNSSECOptions) (*DNSSECAnswer, *Response, error) {
	return call[DNSSECAnswer](ctx, s.client, "dnssec/disable", opts)
}

// RenewKSK renews the KSK key of domains on REG.RU name servers (dnssec/renew_ksk).
func (s *DNSSECService) RenewKSK(ctx context.Context, opts *DNSSECOptions) (*DNSSECAnswer, *Response, error) {
	return call[DNSSECAnswer](ctx, s.client, "dnssec/renew_ksk", opts)
}

// RenewZSK renews the ZSK key of domains on REG.RU name servers (dnssec/renew_zsk).
func (s *DNSSECService) RenewZSK(ctx context.Context, opts *DNSSECOptions) (*DNSSECAnswer, *Response, error) {
	return call[DNSSECAnswer](ctx, s.client, "dnssec/renew_zsk", opts)
}

// GetRecords returns the DNSSEC records of domains (dnssec/get_records).
func (s *DNSSECService) GetRecords(ctx context.Context, opts *DNSSECOptions) (*DNSSECAnswer, *Response, error) {
	return call[DNSSECAnswer](ctx, s.client, "dnssec/get_records", opts)
}

// AddKeys passes KSK keys to the parent zone for domains NOT on REG.RU name servers (dnssec/add_keys).
func (s *DNSSECService) AddKeys(ctx context.Context, opts *DNSSECOptions) (*DNSSECAnswer, *Response, error) {
	return call[DNSSECAnswer](ctx, s.client, "dnssec/add_keys", opts)
}
