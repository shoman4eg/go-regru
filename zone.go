package regru

import "context"

// ZoneService handles functions of the "zone" category.
type ZoneService service

// ZoneDomainsOptions identifies one domain (embedded ServiceIdent) or a list of domains (Domains).
// DNS zone management works only for domains delegated to ns1.reg.ru and ns2.reg.ru.
type ZoneDomainsOptions struct {
	ServiceIdent
	Domains []ServiceIdent `json:"domains,omitempty"`
}

// ZoneAnswer is the per-domain result list returned by most zone functions.
type ZoneAnswer struct {
	Domains []ServiceItem `json:"domains"`
}

// Nop checks that DNS zone management is available for the domains.
//
// API: zone/nop
func (s *ZoneService) Nop(ctx context.Context, opts *ZoneDomainsOptions) (*ZoneAnswer, *Response, error) {
	return call[ZoneAnswer](ctx, s.client, "zone/nop", opts)
}

// ZoneAddAliasOptions are the parameters of AddAlias and AddAAAA.
type ZoneAddAliasOptions struct {
	ZoneDomainsOptions
	Subdomain string `json:"subdomain"` // "@" for the domain itself, "*" for wildcard
	IPAddr    string `json:"ipaddr"`
}

// AddAlias adds an A record binding a subdomain to an IPv4 address.
//
// API: zone/add_alias
func (s *ZoneService) AddAlias(ctx context.Context, opts *ZoneAddAliasOptions) (*ZoneAnswer, *Response, error) {
	return call[ZoneAnswer](ctx, s.client, "zone/add_alias", opts)
}

// AddAAAA adds an AAAA record binding a subdomain to an IPv6 address.
//
// API: zone/add_aaaa
func (s *ZoneService) AddAAAA(ctx context.Context, opts *ZoneAddAliasOptions) (*ZoneAnswer, *Response, error) {
	return call[ZoneAnswer](ctx, s.client, "zone/add_aaaa", opts)
}

// ZoneAddCAAOptions are the parameters of AddCAA.
type ZoneAddCAAOptions struct {
	ZoneDomainsOptions
	Subdomain string `json:"subdomain"`
	Flags     int    `json:"flags"` // 0 or 128
	Tag       string `json:"tag"`   // issue, issuewild, iodef
	Value     string `json:"value"`
}

// AddCAA adds a CAA record restricting certificate issuance for a subdomain.
//
// API: zone/add_caa
func (s *ZoneService) AddCAA(ctx context.Context, opts *ZoneAddCAAOptions) (*ZoneAnswer, *Response, error) {
	return call[ZoneAnswer](ctx, s.client, "zone/add_caa", opts)
}

// ZoneAddCNAMEOptions are the parameters of AddCNAME.
type ZoneAddCNAMEOptions struct {
	ZoneDomainsOptions
	Subdomain     string `json:"subdomain"`
	CanonicalName string `json:"canonical_name"`
}

// AddCNAME adds a CNAME record pointing a subdomain to another domain.
//
// API: zone/add_cname
func (s *ZoneService) AddCNAME(ctx context.Context, opts *ZoneAddCNAMEOptions) (*ZoneAnswer, *Response, error) {
	return call[ZoneAnswer](ctx, s.client, "zone/add_cname", opts)
}

// ZoneAddHTTPSOptions are the parameters of AddHTTPS.
type ZoneAddHTTPSOptions struct {
	ZoneDomainsOptions
	Subdomain string `json:"subdomain"`
	Priority  int    `json:"priority"` // 0 is AliasMode, >0 is ServiceMode
	Target    string `json:"target"`
	Value     string `json:"value,omitempty"` // space-separated key=value pairs, e.g. "alpn=h3 port=443"
}

// AddHTTPS adds an HTTPS (SVCB) record for a subdomain.
//
// API: zone/add_https
func (s *ZoneService) AddHTTPS(ctx context.Context, opts *ZoneAddHTTPSOptions) (*ZoneAnswer, *Response, error) {
	return call[ZoneAnswer](ctx, s.client, "zone/add_https", opts)
}

// ZoneAddMXOptions are the parameters of AddMX.
type ZoneAddMXOptions struct {
	ZoneDomainsOptions
	Subdomain  string `json:"subdomain,omitempty"` // default "@"
	Priority   int    `json:"priority,omitempty"`  // 0 (highest, default) to 10
	MailServer string `json:"mail_server"`
}

// AddMX adds an MX record.
//
// API: zone/add_mx
func (s *ZoneService) AddMX(ctx context.Context, opts *ZoneAddMXOptions) (*ZoneAnswer, *Response, error) {
	return call[ZoneAnswer](ctx, s.client, "zone/add_mx", opts)
}

// ZoneAddNSOptions are the parameters of AddNS.
type ZoneAddNSOptions struct {
	ZoneDomainsOptions
	Subdomain    string `json:"subdomain"`
	DNSServer    string `json:"dns_server"`
	RecordNumber int    `json:"record_number,omitempty"`
}

// AddNS adds an NS record delegating a subdomain to other DNS servers.
//
// API: zone/add_ns
func (s *ZoneService) AddNS(ctx context.Context, opts *ZoneAddNSOptions) (*ZoneAnswer, *Response, error) {
	return call[ZoneAnswer](ctx, s.client, "zone/add_ns", opts)
}

// ZoneAddTXTOptions are the parameters of AddTXT.
type ZoneAddTXTOptions struct {
	ZoneDomainsOptions
	Subdomain string `json:"subdomain"`
	Text      string `json:"text"`
}

// AddTXT adds a TXT record.
//
// API: zone/add_txt
func (s *ZoneService) AddTXT(ctx context.Context, opts *ZoneAddTXTOptions) (*ZoneAnswer, *Response, error) {
	return call[ZoneAnswer](ctx, s.client, "zone/add_txt", opts)
}

// ZoneAddSRVOptions are the parameters of AddSRV.
type ZoneAddSRVOptions struct {
	ZoneDomainsOptions
	Service  string `json:"service"` // e.g. "_sip._udp"
	Priority int    `json:"priority"`
	Weight   int    `json:"weight,omitempty"`
	Target   string `json:"target"`
	Port     int    `json:"port"`
}

// AddSRV adds an SRV record.
//
// API: zone/add_srv
func (s *ZoneService) AddSRV(ctx context.Context, opts *ZoneAddSRVOptions) (*ZoneAnswer, *Response, error) {
	return call[ZoneAnswer](ctx, s.client, "zone/add_srv", opts)
}

// ZoneRecord is a resource record in GetResourceRecords answers.
type ZoneRecord struct {
	Subname string  `json:"subname"`
	RecType string  `json:"rectype"`
	State   string  `json:"state"`
	Prio    FlexInt `json:"prio"`
	Content string  `json:"content"`
}

// ZoneSOA holds zone cache lifetimes, e.g. "1d", "12h".
type ZoneSOA struct {
	TTL        string `json:"ttl"`
	MinimumTTL string `json:"minimum_ttl"`
}

// ZoneRecordsDomain is a per-domain entry of ZoneRecordsAnswer.
type ZoneRecordsDomain struct {
	ServiceItem
	RRs []ZoneRecord `json:"rrs"`
	SOA ZoneSOA      `json:"soa"`
}

// ZoneRecordsAnswer is returned by GetResourceRecords.
type ZoneRecordsAnswer struct {
	Domains []ZoneRecordsDomain `json:"domains"`
}

// GetResourceRecords returns the resource records of each domain zone.
//
// API: zone/get_resource_records
func (s *ZoneService) GetResourceRecords(ctx context.Context, opts *ZoneDomainsOptions) (*ZoneRecordsAnswer, *Response, error) {
	return call[ZoneRecordsAnswer](ctx, s.client, "zone/get_resource_records", opts)
}

// ZoneAction is one step of UpdateRecords. Action is one of add_alias, add_aaaa, add_caa,
// add_cname, add_https, add_mx, add_ns, add_txt, add_srv, add_spf, remove_record;
// set the fields of the matching function.
type ZoneAction struct {
	Action        string `json:"action"`
	Subdomain     string `json:"subdomain,omitempty"`
	IPAddr        string `json:"ipaddr,omitempty"`
	CanonicalName string `json:"canonical_name,omitempty"`
	MailServer    string `json:"mail_server,omitempty"`
	DNSServer     string `json:"dns_server,omitempty"`
	RecordNumber  int    `json:"record_number,omitempty"`
	Text          string `json:"text,omitempty"`
	Service       string `json:"service,omitempty"`
	Priority      *int   `json:"priority,omitempty"` // pointer: 0 is meaningful for https/srv
	Weight        int    `json:"weight,omitempty"`
	Target        string `json:"target,omitempty"`
	Port          int    `json:"port,omitempty"`
	Flags         *int   `json:"flags,omitempty"` // pointer: 0 is meaningful for caa
	Tag           string `json:"tag,omitempty"`
	Value         string `json:"value,omitempty"`
	RecordType    string `json:"record_type,omitempty"`
	Content       string `json:"content,omitempty"`
}

// ZoneUpdateRecordsDomain is a domain with its own action list.
type ZoneUpdateRecordsDomain struct {
	ServiceIdent
	ActionList []ZoneAction `json:"action_list,omitempty"`
}

// ZoneUpdateRecordsOptions are the parameters of UpdateRecords. ActionList applies to
// all domains; alternatively set ActionList per domain in Domains.
type ZoneUpdateRecordsOptions struct {
	ServiceIdent
	Domains    []ZoneUpdateRecordsDomain `json:"domains,omitempty"`
	ActionList []ZoneAction              `json:"action_list,omitempty"`
}

// ZoneActionResult is the result of one ZoneAction.
type ZoneActionResult struct {
	ItemResult
	Action string `json:"action"`
}

// ZoneUpdateRecordsResult is a per-domain entry of ZoneUpdateRecordsAnswer.
type ZoneUpdateRecordsResult struct {
	ServiceItem
	ActionList  []ZoneActionResult `json:"action_list"`
	ErrorParams map[string]any     `json:"error_params,omitempty"`
}

// ZoneUpdateRecordsAnswer is returned by UpdateRecords.
type ZoneUpdateRecordsAnswer struct {
	Domains []ZoneUpdateRecordsResult `json:"domains"`
}

// UpdateRecords adds and/or removes several records in one request; actions run in order
// and stop at the first error.
//
// API: zone/update_records
func (s *ZoneService) UpdateRecords(ctx context.Context, opts *ZoneUpdateRecordsOptions) (*ZoneUpdateRecordsAnswer, *Response, error) {
	return call[ZoneUpdateRecordsAnswer](ctx, s.client, "zone/update_records", opts)
}

// ZoneUpdateSOAOptions are the parameters of UpdateSOA. Values are seconds or a number
// with suffix m (months), w, d or h.
type ZoneUpdateSOAOptions struct {
	ZoneDomainsOptions
	TTL        string `json:"ttl,omitempty"`
	MinimumTTL string `json:"minimum_ttl,omitempty"`
}

// UpdateSOA changes the zone cache lifetimes.
//
// API: zone/update_soa
func (s *ZoneService) UpdateSOA(ctx context.Context, opts *ZoneUpdateSOAOptions) (*ZoneAnswer, *Response, error) {
	return call[ZoneAnswer](ctx, s.client, "zone/update_soa", opts)
}

// TuneForwarding adds the resource records needed for web forwarding.
//
// API: zone/tune_forwarding
func (s *ZoneService) TuneForwarding(ctx context.Context, opts *ZoneDomainsOptions) (*ZoneAnswer, *Response, error) {
	return call[ZoneAnswer](ctx, s.client, "zone/tune_forwarding", opts)
}

// ClearForwarding removes the resource records needed for web forwarding.
//
// API: zone/clear_forwarding
func (s *ZoneService) ClearForwarding(ctx context.Context, opts *ZoneDomainsOptions) (*ZoneAnswer, *Response, error) {
	return call[ZoneAnswer](ctx, s.client, "zone/clear_forwarding", opts)
}

// TuneParking adds the resource records needed for domain parking.
//
// API: zone/tune_parking
func (s *ZoneService) TuneParking(ctx context.Context, opts *ZoneDomainsOptions) (*ZoneAnswer, *Response, error) {
	return call[ZoneAnswer](ctx, s.client, "zone/tune_parking", opts)
}

// ClearParking removes the resource records needed for domain parking.
//
// API: zone/clear_parking
func (s *ZoneService) ClearParking(ctx context.Context, opts *ZoneDomainsOptions) (*ZoneAnswer, *Response, error) {
	return call[ZoneAnswer](ctx, s.client, "zone/clear_parking", opts)
}

// ZoneRemoveRecordOptions are the parameters of RemoveRecord.
type ZoneRemoveRecordOptions struct {
	ZoneDomainsOptions
	Subdomain  string `json:"subdomain"`
	RecordType string `json:"record_type"`
	Priority   int    `json:"priority,omitempty"` // default 0; not for A-like records
	Content    string `json:"content,omitempty"`  // empty removes all matching records
}

// RemoveRecord removes a resource record.
//
// API: zone/remove_record
func (s *ZoneService) RemoveRecord(ctx context.Context, opts *ZoneRemoveRecordOptions) (*ZoneAnswer, *Response, error) {
	return call[ZoneAnswer](ctx, s.client, "zone/remove_record", opts)
}

// Clear removes all resource records of the zone.
//
// API: zone/clear
func (s *ZoneService) Clear(ctx context.Context, opts *ZoneDomainsOptions) (*ZoneAnswer, *Response, error) {
	return call[ZoneAnswer](ctx, s.client, "zone/clear", opts)
}
