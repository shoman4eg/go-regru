package regru

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
)

// DomainService handles functions of the "domain" category.
type DomainService service

// DomainContacts are domain contact fields; keys depend on the TLD and registrant type
// (e.g. person, person_r, passport, phone, e_mail, country, o_company, a_email). See domain/create.
type DomainContacts map[string]string

// DomainNSS are name servers: ns0..ns3 and ns0ip..ns3ip (several IPs comma-separated).
type DomainNSS map[string]string

// DomainListOptions identifies one domain (embedded ServiceIdent) or a list of domains (Domains).
type DomainListOptions struct {
	ServiceIdent
	Domains []ServiceIdent `json:"domains,omitempty"`
}

// DomainServicesAnswer is the per-domain result list. Some functions return it as
// "services" instead of "domains".
type DomainServicesAnswer struct {
	Domains  []ServiceItem `json:"domains,omitempty"`
	Services []ServiceItem `json:"services,omitempty"`
}

// DomainOrderAnswer is returned by functions that create a bill.
type DomainOrderAnswer struct {
	BillID   FlexInt       `json:"bill_id"`
	Payment  FlexFloat     `json:"payment"`
	PayType  string        `json:"pay_type"`
	PayNotes string        `json:"pay_notes"`
	Domains  []ServiceItem `json:"domains"`
}

// Nop checks API access; with a domain name set it returns the domain's service_id.
//
// API: domain/nop
func (s *DomainService) Nop(ctx context.Context, opts *ServiceIdent) (*ServiceIDAnswer, *Response, error) {
	return call[ServiceIDAnswer](ctx, s.client, "domain/nop", opts)
}

// DomainGetPricesOptions are the parameters of GetPrices.
type DomainGetPricesOptions struct {
	ShowRenewData  Bool   `json:"show_renew_data,omitempty"`
	ShowUpdateData Bool   `json:"show_update_data,omitempty"`
	Currency       string `json:"currency,omitempty"` // RUR (default), UAH, USD, EUR
}

// DomainPrice is the per-year price of a zone.
// shortcut: update prices (show_update_data) are not modelled, add fields once their names are known.
type DomainPrice struct {
	RegPrice         FlexFloat `json:"reg_price"`
	RetailRegPrice   FlexFloat `json:"retail_reg_price"`
	RenewPrice       FlexFloat `json:"renew_price"`
	RetailRenewPrice FlexFloat `json:"retail_renew_price"`
	RegMinPeriod     FlexInt   `json:"reg_min_period"`
	RegMaxPeriod     FlexInt   `json:"reg_max_period"`
}

// DomainGetPricesAnswer is returned by GetPrices.
type DomainGetPricesAnswer struct {
	Currency       string `json:"currency"`
	PriceGroup     string `json:"price_group"`
	LastUpdateTime string `json:"last_update_time"`
	// Prices by zone; IDN prices that differ use the "__idn." key prefix.
	Prices map[string]DomainPrice `json:"prices"`
}

// GetPrices returns registration/renewal prices in all zones.
//
// API: domain/get_prices
func (s *DomainService) GetPrices(ctx context.Context, opts *DomainGetPricesOptions) (*DomainGetPricesAnswer, *Response, error) {
	return call[DomainGetPricesAnswer](ctx, s.client, "domain/get_prices", opts)
}

// DomainGetSuggestOptions are the parameters of GetSuggest.
type DomainGetSuggestOptions struct {
	Word           string   `json:"word"`
	AdditionalWord string   `json:"additional_word,omitempty"`
	Category       string   `json:"category,omitempty"` // pattern, search_trends, all
	TLDs           []string `json:"tlds,omitempty"`
	UseHyphen      Bool     `json:"use_hyphen,omitempty"`
}

// DomainSuggestion is a suggested name and the zones where it is available.
type DomainSuggestion struct {
	Name    string   `json:"name"`
	AvailIn []string `json:"avail_in"`
}

// DomainGetSuggestAnswer is returned by GetSuggest.
type DomainGetSuggestAnswer struct {
	Suggestions []DomainSuggestion `json:"suggestions"`
}

// GetSuggest suggests domain names by keywords.
//
// API: domain/get_suggest
func (s *DomainService) GetSuggest(ctx context.Context, opts *DomainGetSuggestOptions) (*DomainGetSuggestAnswer, *Response, error) {
	return call[DomainGetSuggestAnswer](ctx, s.client, "domain/get_suggest", opts)
}

// DomainGetPremiumPricesOptions are the parameters of GetPremiumPrices.
type DomainGetPremiumPricesOptions struct {
	Domains  []string `json:"domains"` // at most 100
	Currency string   `json:"currency,omitempty"`
}

// DomainPremiumPrice is the price of a premium domain.
type DomainPremiumPrice struct {
	DomainName     string    `json:"domain_name"`
	RegisterPrice  FlexFloat `json:"register_price"`
	RegisterPeriod FlexInt   `json:"register_period"`
	RenewPrice     FlexFloat `json:"renew_price"`
	RenewPeriod    FlexInt   `json:"renew_period"`
}

// DomainGetPremiumPricesAnswer is returned by GetPremiumPrices.
type DomainGetPremiumPricesAnswer struct {
	Currency string               `json:"currency"`
	Domains  []DomainPremiumPrice `json:"domains"`
}

// GetPremiumPrices returns prices of premium domains.
//
// API: domain/get_premium_prices
func (s *DomainService) GetPremiumPrices(ctx context.Context, opts *DomainGetPremiumPricesOptions) (*DomainGetPremiumPricesAnswer, *Response, error) {
	return call[DomainGetPremiumPricesAnswer](ctx, s.client, "domain/get_premium_prices", opts)
}

// DomainGetDeletedOptions are the parameters of GetDeleted. Dates are YYYY-MM-DD.
type DomainGetDeletedOptions struct {
	TLDs        []string `json:"tlds,omitempty"`
	DeletedFrom string   `json:"deleted_from,omitempty"`
	DeletedTo   string   `json:"deleted_to,omitempty"`
	CreatedFrom string   `json:"created_from,omitempty"`
	CreatedTo   string   `json:"created_to,omitempty"`
	HideReg     Bool     `json:"hidereg,omitempty"`
	MinPR       int      `json:"min_pr,omitempty"`
	MinCY       int      `json:"min_cy,omitempty"`
}

// DomainDeleted is a released domain.
type DomainDeleted struct {
	DomainName      string  `json:"domain_name"`
	DateDelete      string  `json:"date_delete"`
	Registered      string  `json:"registered"` // "NOT REGISTERED" or "REGISTERED"
	FirstCreateDate string  `json:"first_create_date"`
	YandexTIC       FlexInt `json:"yandex_tic"`
	GooglePR        FlexInt `json:"google_pr"`
}

// DomainGetDeletedAnswer is returned by GetDeleted.
type DomainGetDeletedAnswer struct {
	Domains []DomainDeleted `json:"domains"`
}

// GetDeleted returns released domains.
//
// API: domain/get_deleted
func (s *DomainService) GetDeleted(ctx context.Context, opts *DomainGetDeletedOptions) (*DomainGetDeletedAnswer, *Response, error) {
	return call[DomainGetDeletedAnswer](ctx, s.client, "domain/get_deleted", opts)
}

// DomainCheckOptions are the parameters of Check. Set DomainName or Domains.
type DomainCheckOptions struct {
	DomainName     string         `json:"domain_name,omitempty"`
	Domains        []ServiceIdent `json:"domains,omitempty"`
	IsTransfer     Bool           `json:"is_transfer,omitempty"`
	PremiumAsTaken Bool           `json:"premium_as_taken,omitempty"`
	Currency       string         `json:"currency,omitempty"`
}

// DomainCheckItem is the availability of one domain; Result is "Available" if it can be registered.
type DomainCheckItem struct {
	ItemResult
	DName      string    `json:"dname"`
	IsPremium  Bool      `json:"is_premium"`
	Price      FlexFloat `json:"price"`
	RenewPrice FlexFloat `json:"renew_price"`
	Currency   string    `json:"currency"`
}

// DomainCheckAnswer is returned by Check.
type DomainCheckAnswer struct {
	Domains []DomainCheckItem `json:"domains"`
}

// Check checks whether domains can be registered (or transferred).
//
// API: domain/check
func (s *DomainService) Check(ctx context.Context, opts *DomainCheckOptions) (*DomainCheckAnswer, *Response, error) {
	return call[DomainCheckAnswer](ctx, s.client, "domain/check", opts)
}

// DomainCreateItem is a domain in a bulk registration with optional extra services.
type DomainCreateItem struct {
	DName          string `json:"dname"`
	SrvCertificate string `json:"srv_certificate,omitempty"`
	SrvParking     string `json:"srv_parking,omitempty"`
	SrvWebfwd      string `json:"srv_webfwd,omitempty"`
}

// DomainCreateOptions are the parameters of Create. Set DomainName or Domains, and
// Contacts or ProfileType+ProfileName.
// shortcut: period 0 (transfer without renewal) is dropped by omitempty, make Period *int if needed.
type DomainCreateOptions struct {
	DomainName   string             `json:"domain_name,omitempty"`
	Domains      []DomainCreateItem `json:"domains,omitempty"`
	Period       int                `json:"period,omitempty"`
	EndUserIP    string             `json:"enduser_ip,omitempty"`
	Contacts     DomainContacts     `json:"contacts,omitempty"`
	ProfileType  string             `json:"profile_type,omitempty"`
	ProfileName  string             `json:"profile_name,omitempty"`
	NSS          DomainNSS          `json:"nss,omitempty"`
	NotDelegated Bool               `json:"not_delegated,omitempty"`
	UserServID   string             `json:"user_servid,omitempty"`
	Comment      string             `json:"comment,omitempty"`
	AdminComment string             `json:"admin_comment,omitempty"`
	SubType      string             `json:"subtype,omitempty"` // "preorder" for .РФ preorders
	RegPremium   Bool               `json:"reg_premium,omitempty"`
	PaymentOptions
	FolderIdent
	NoNewFolder Bool `json:"no_new_folder,omitempty"`
}

// Create orders domain registration.
//
// API: domain/create
func (s *DomainService) Create(ctx context.Context, opts *DomainCreateOptions) (*DomainOrderAnswer, *Response, error) {
	return call[DomainOrderAnswer](ctx, s.client, "domain/create", opts)
}

// DomainTransferOptions are the parameters of Transfer.
type DomainTransferOptions struct {
	DomainCreateOptions
	AuthInfo string `json:"authinfo,omitempty"`
}

// Transfer orders a domain transfer from another registrar.
//
// API: domain/transfer
func (s *DomainService) Transfer(ctx context.Context, opts *DomainTransferOptions) (*DomainOrderAnswer, *Response, error) {
	return call[DomainOrderAnswer](ctx, s.client, "domain/transfer", opts)
}

// DomainTransferStatus is the transfer status of one domain (pendingTransfer, authinfoError, ...).
type DomainTransferStatus struct {
	ServiceItem
	Status string `json:"status"`
}

// DomainGetTransferStatusAnswer is returned by GetTransferStatus.
type DomainGetTransferStatusAnswer struct {
	Domains []DomainTransferStatus `json:"domains"`
}

// GetTransferStatus returns the transfer status of domains.
//
// API: domain/get_transfer_status
func (s *DomainService) GetTransferStatus(ctx context.Context, opts *DomainListOptions) (*DomainGetTransferStatusAnswer, *Response, error) {
	return call[DomainGetTransferStatusAnswer](ctx, s.client, "domain/get_transfer_status", opts)
}

// DomainSetNewAuthInfoOptions are the parameters of SetNewAuthInfo.
type DomainSetNewAuthInfoOptions struct {
	DomainListOptions
	AuthInfo string `json:"authinfo"`
}

// DomainSetNewAuthInfoAnswer is returned by SetNewAuthInfo.
type DomainSetNewAuthInfoAnswer struct {
	DName     string        `json:"dname"`
	ServiceID FlexInt       `json:"service_id"`
	Domains   []ServiceItem `json:"domains"`
}

// SetNewAuthInfo replaces a wrong transfer key before the request is sent to the registry.
//
// API: domain/set_new_authinfo
func (s *DomainService) SetNewAuthInfo(ctx context.Context, opts *DomainSetNewAuthInfoOptions) (*DomainSetNewAuthInfoAnswer, *Response, error) {
	return call[DomainSetNewAuthInfoAnswer](ctx, s.client, "domain/set_new_authinfo", opts)
}

// CancelTransfer cancels a domain transfer.
// It is listed in section 3.3 but undocumented, hence the untyped answer.
//
// API: domain/cancel_transfer
func (s *DomainService) CancelTransfer(ctx context.Context, opts *DomainListOptions) (*map[string]any, *Response, error) {
	return call[map[string]any](ctx, s.client, "domain/cancel_transfer", opts)
}

// DomainGetReregDataOptions are the parameters of GetReregData.
type DomainGetReregDataOptions struct {
	Domains       []ServiceIdent `json:"domains,omitempty"`
	DNameMatching string         `json:"dname_matching,omitempty"`
	SearchQuery   string         `json:"search_query,omitempty"`
	MaxChars      int            `json:"max_chars,omitempty"`
	MinPR         int            `json:"min_pr,omitempty"`
	MinCY         int            `json:"min_cy,omitempty"`
	Kley          Bool           `json:"kley,omitempty"`
	Zone          string         `json:"zone,omitempty"`
	Views         int            `json:"views,omitempty"`
	Vis           int            `json:"vis,omitempty"`
	Traf          int            `json:"traf,omitempty"`
	Price         float64        `json:"price,omitempty"`
	Registrar     string         `json:"registrar,omitempty"`
	DeleteDate    string         `json:"delete_date,omitempty"`
	SortCol       string         `json:"sortcol,omitempty"` // name, price, fdate, date, pr, cy, views, vis, sch, reg
	Premium       Bool           `json:"premium,omitempty"`
	Limit         int            `json:"limit,omitempty"` // default and max 1000
	LimitFrom     int            `json:"limit_from,omitempty"`
}

// DomainReregData describes an expiring domain.
type DomainReregData struct {
	DName            string    `json:"dname"`
	TLD              string    `json:"tld"`
	FirstCreateDate  string    `json:"first_create_date"`
	DeleteDate       string    `json:"delete_date"`
	StartPrice       FlexFloat `json:"start_price"`
	Bid              FlexFloat `json:"bid"`
	Price            FlexFloat `json:"price"`
	BlitzPrice       FlexFloat `json:"blitz_price"`
	UniAvgAttendance FlexInt   `json:"uni_avg_attendance"`
	AvgViewings      FlexInt   `json:"avg_viewings"`
	AllAvgTraffic    FlexInt   `json:"all_avg_traffic"`
	SearchQueryList  string    `json:"search_query_list"`
	YandexTIC        FlexInt   `json:"yandex_tic"`
	TICMirrors       string    `json:"tic_mirrors"`
	GooglePR         FlexInt   `json:"google_pr"`
	Registrar        string    `json:"registrar"`
	IsRecommended    Bool      `json:"is_recommended"`
}

// GetReregData returns expiring domains.
//
// API: domain/get_rereg_data
func (s *DomainService) GetReregData(ctx context.Context, opts *DomainGetReregDataOptions) (*[]DomainReregData, *Response, error) {
	return call[[]DomainReregData](ctx, s.client, "domain/get_rereg_data", opts)
}

// DomainReregBid is a bid on an expiring domain.
type DomainReregBid struct {
	DName      string  `json:"dname"`
	Price      float64 `json:"price"`
	Instalment Bool    `json:"instalment,omitempty"`
	Autopay    Bool    `json:"autopay,omitempty"`
}

// DomainSetReregBidsOptions are the parameters of SetReregBids.
type DomainSetReregBidsOptions struct {
	Contacts DomainContacts   `json:"contacts,omitempty"`
	NSS      DomainNSS        `json:"nss,omitempty"`
	Domains  []DomainReregBid `json:"domains"`
}

// SetReregBids places bids on expiring domains.
//
// API: domain/set_rereg_bids
func (s *DomainService) SetReregBids(ctx context.Context, opts *DomainSetReregBidsOptions) (*DomainOrderAnswer, *Response, error) {
	return call[DomainOrderAnswer](ctx, s.client, "domain/set_rereg_bids", opts)
}

// DomainUserReregBid is an expiring domain with the user's bid.
type DomainUserReregBid struct {
	DomainReregData
	IsActive  Bool      `json:"is_active"`
	UserBid   FlexFloat `json:"user_bid"`
	MaxBid    FlexFloat `json:"max_bid"`
	NextPrice FlexFloat `json:"next_price"`
}

// DomainGetUserReregBidsAnswer is returned by GetUserReregBids.
type DomainGetUserReregBidsAnswer struct {
	Domains []DomainUserReregBid `json:"domains"`
	// ReregBids lists the allowed bids by zone.
	ReregBids map[string][]FlexFloat `json:"rereg_bids"`
}

// GetUserReregBids returns expiring domains the user has bid on.
//
// API: domain/get_user_rereg_bids
func (s *DomainService) GetUserReregBids(ctx context.Context) (*DomainGetUserReregBidsAnswer, *Response, error) {
	return call[DomainGetUserReregBidsAnswer](ctx, s.client, "domain/get_user_rereg_bids", nil)
}

// DomainGetDocsUploadURIAnswer is returned by GetDocsUploadURI.
type DomainGetDocsUploadURIAnswer struct {
	DocsUploadSID FlexString `json:"docs_upload_sid"`
	URL           string     `json:"url"`
}

// GetDocsUploadURI returns a link for uploading documents for a .RU/.SU/.РФ domain.
//
// API: domain/get_docs_upload_uri
func (s *DomainService) GetDocsUploadURI(ctx context.Context, opts *ServiceIdent) (*DomainGetDocsUploadURIAnswer, *Response, error) {
	return call[DomainGetDocsUploadURIAnswer](ctx, s.client, "domain/get_docs_upload_uri", opts)
}

// DomainUpdateContactsOptions are the parameters of UpdateContacts.
type DomainUpdateContactsOptions struct {
	DomainListOptions
	Contacts        DomainContacts `json:"contacts"`
	PayoutAgreement Bool           `json:"payout_agreement,omitempty"`
}

// UpdateContacts changes domain contacts.
//
// API: domain/update_contacts
func (s *DomainService) UpdateContacts(ctx context.Context, opts *DomainUpdateContactsOptions) (*DomainServicesAnswer, *Response, error) {
	return call[DomainServicesAnswer](ctx, s.client, "domain/update_contacts", opts)
}

// DomainUpdatePrivatePersonFlagOptions are the parameters of UpdatePrivatePersonFlag.
type DomainUpdatePrivatePersonFlagOptions struct {
	DomainListOptions
	PrivatePersonFlag Bool `json:"private_person_flag"` // no omitempty: 0 clears the flag
}

// DomainUpdatePrivatePersonFlagAnswer is returned by UpdatePrivatePersonFlag.
type DomainUpdatePrivatePersonFlagAnswer struct {
	Domains []ServiceItem `json:"domains"`
	PPFlag  string        `json:"pp_flag"` // "is set" or "is cleared"
}

// UpdatePrivatePersonFlag hides or shows contacts in whois.
//
// API: domain/update_private_person_flag
func (s *DomainService) UpdatePrivatePersonFlag(ctx context.Context, opts *DomainUpdatePrivatePersonFlagOptions) (*DomainUpdatePrivatePersonFlagAnswer, *Response, error) {
	return call[DomainUpdatePrivatePersonFlagAnswer](ctx, s.client, "domain/update_private_person_flag", opts)
}

// DomainNSOptions are the parameters of RegisterNS and DeleteNS.
type DomainNSOptions struct {
	ServiceIdent
	NS0   string `json:"ns0"`
	NS0IP string `json:"ns0ip,omitempty"`
}

// DomainNSAnswer is returned by RegisterNS and DeleteNS.
type DomainNSAnswer struct {
	Resp map[string]any `json:"resp"` // registrar's detailed answer
}

// RegisterNS registers a child name server in the NSI registry (international domains only).
//
// API: domain/register_ns
func (s *DomainService) RegisterNS(ctx context.Context, opts *DomainNSOptions) (*DomainNSAnswer, *Response, error) {
	return call[DomainNSAnswer](ctx, s.client, "domain/register_ns", opts)
}

// DeleteNS deletes a child name server from the NSI registry (international domains only).
//
// API: domain/delete_ns
func (s *DomainService) DeleteNS(ctx context.Context, opts *DomainNSOptions) (*DomainNSAnswer, *Response, error) {
	return call[DomainNSAnswer](ctx, s.client, "domain/delete_ns", opts)
}

// DomainNS is a name server; IP is set only for in-zone servers (a string or an array of strings).
type DomainNS struct {
	NS string `json:"ns"`
	IP any    `json:"ip,omitempty"`
}

// DomainNSSItem is the name servers of one domain.
type DomainNSSItem struct {
	ServiceItem
	NSS []DomainNS `json:"nss"`
}

// DomainGetNSSAnswer is returned by GetNSS.
type DomainGetNSSAnswer struct {
	Domains []DomainNSSItem `json:"domains"`
}

// GetNSS returns name servers of domains.
//
// API: domain/get_nss
func (s *DomainService) GetNSS(ctx context.Context, opts *DomainListOptions) (*DomainGetNSSAnswer, *Response, error) {
	return call[DomainGetNSSAnswer](ctx, s.client, "domain/get_nss", opts)
}

// DomainUpdateNSSOptions are the parameters of UpdateNSS.
type DomainUpdateNSSOptions struct {
	DomainListOptions
	NSS             DomainNSS `json:"nss,omitempty"`
	Undelegate      Bool      `json:"undelegate,omitempty"`
	PayoutAgreement Bool      `json:"payout_agreement,omitempty"`
}

// UpdateNSS changes name servers of domains.
//
// API: domain/update_nss
func (s *DomainService) UpdateNSS(ctx context.Context, opts *DomainUpdateNSSOptions) (*DomainServicesAnswer, *Response, error) {
	return call[DomainServicesAnswer](ctx, s.client, "domain/update_nss", opts)
}

// Delegate sets the delegation flag of domains.
//
// API: domain/delegate
func (s *DomainService) Delegate(ctx context.Context, opts *DomainListOptions) (*DomainServicesAnswer, *Response, error) {
	return call[DomainServicesAnswer](ctx, s.client, "domain/delegate", opts)
}

// Undelegate clears the delegation flag of domains.
//
// API: domain/undelegate
func (s *DomainService) Undelegate(ctx context.Context, opts *DomainListOptions) (*DomainServicesAnswer, *Response, error) {
	return call[DomainServicesAnswer](ctx, s.client, "domain/undelegate", opts)
}

// DomainTransferToAnotherAccountOptions are the parameters of TransferToAnotherAccount.
type DomainTransferToAnotherAccountOptions struct {
	DomainListOptions
	NewUserName     string `json:"new_user_name"`
	SetMeAsReferrer Bool   `json:"set_me_as_referrer,omitempty"`
}

// TransferToAnotherAccount moves domains to another account; Result is "request_is_sent" on success.
//
// API: domain/transfer_to_another_account
func (s *DomainService) TransferToAnotherAccount(ctx context.Context, opts *DomainTransferToAnotherAccountOptions) (*DomainServicesAnswer, *Response, error) {
	return call[DomainServicesAnswer](ctx, s.client, "domain/transfer_to_another_account", opts)
}

// DomainEnteringMessage is an incoming domain transfer to this account.
type DomainEnteringMessage struct {
	ID          FlexInt  `json:"id"`
	DomainNames []string `json:"domain_names"`
}

// DomainLookAtEnteringListAnswer is returned by LookAtEnteringList.
type DomainLookAtEnteringListAnswer struct {
	Messages []DomainEnteringMessage `json:"messages"`
}

// LookAtEnteringList lists domains being transferred to this account.
//
// API: domain/look_at_entering_list
func (s *DomainService) LookAtEnteringList(ctx context.Context) (*DomainLookAtEnteringListAnswer, *Response, error) {
	return call[DomainLookAtEnteringListAnswer](ctx, s.client, "domain/look_at_entering_list", nil)
}

// DomainEnteringAction accepts or refuses one incoming transfer message.
type DomainEnteringAction struct {
	ID         int64  `json:"id,omitempty"`
	ActionType string `json:"action_type,omitempty"` // accept/refuse, yes/no, 1/0
}

// DomainAcceptOrRefuseEnteringListOptions are the parameters of AcceptOrRefuseEnteringList.
// Set ID+ActionType or Domains.
type DomainAcceptOrRefuseEnteringListOptions struct {
	DomainEnteringAction
	Domains []DomainEnteringAction `json:"domains,omitempty"`
}

// DomainEnteringResult is the result for one transfer message ("accepted" or "refused").
type DomainEnteringResult struct {
	ItemResult
	ID         FlexInt `json:"id"`
	ActionType string  `json:"action_type"`
}

// DomainAcceptOrRefuseEnteringListAnswer is returned by AcceptOrRefuseEnteringList.
type DomainAcceptOrRefuseEnteringListAnswer struct {
	Domains []DomainEnteringResult `json:"domains"`
}

// AcceptOrRefuseEnteringList accepts or refuses domains transferred to this account.
//
// API: domain/accept_or_refuse_entering_list
func (s *DomainService) AcceptOrRefuseEnteringList(ctx context.Context, opts *DomainAcceptOrRefuseEnteringListOptions) (*DomainAcceptOrRefuseEnteringListAnswer, *Response, error) {
	return call[DomainAcceptOrRefuseEnteringListAnswer](ctx, s.client, "domain/accept_or_refuse_entering_list", opts)
}

// RequestToTransfer asks for REG.RU domains to be moved to this account.
//
// API: domain/request_to_transfer
func (s *DomainService) RequestToTransfer(ctx context.Context, opts *DomainListOptions) (*DomainServicesAnswer, *Response, error) {
	return call[DomainServicesAnswer](ctx, s.client, "domain/request_to_transfer", opts)
}

// DomainGetTLDInfoOptions are the parameters of GetTLDInfo.
type DomainGetTLDInfoOptions struct {
	TLD string `json:"tld"`
	DomainListOptions
}

// DomainTLDInfo describes a zone.
type DomainTLDInfo struct {
	TLD                           string     `json:"tld"`
	EmailVerification             Bool       `json:"email_verification"`
	IDNAllowed                    Bool       `json:"idn_allowed"`
	MinDNameLength                FlexInt    `json:"min_dname_length"`
	MaxDNameLength                FlexInt    `json:"max_dname_length"`
	MinNSSCount                   FlexInt    `json:"min_nss_count"`
	MaxNSSCount                   FlexInt    `json:"max_nss_count"`
	PrivacyProtectionAllowed      Bool       `json:"privacy_protection_allowed"`
	RegPeriods                    FlexString `json:"regperiods"`
	RenewPeriods                  FlexString `json:"renewperiods"`
	RejectDaysOfRenewAfterCreate  FlexInt    `json:"reject_days_of_renew_after_create"`
	RenewGracePeriodBeforeExpDate FlexInt    `json:"renew_grace_period_before_expdate"`
	RenewGracePeriodAfterExpDate  FlexInt    `json:"renew_grace_period_after_expdate"`
	RestorationGracePeriod        FlexInt    `json:"restoration_grace_period"`
}

// DomainGetTLDInfoAnswer is returned by GetTLDInfo.
type DomainGetTLDInfoAnswer struct {
	TLDInfo DomainTLDInfo `json:"tld_info"`
}

// GetTLDInfo returns information about a zone.
//
// API: domain/get_tld_info
func (s *DomainService) GetTLDInfo(ctx context.Context, opts *DomainGetTLDInfoOptions) (*DomainGetTLDInfoAnswer, *Response, error) {
	return call[DomainGetTLDInfoAnswer](ctx, s.client, "domain/get_tld_info", opts)
}

// DomainSendEmailVerificationLetterOptions are the parameters of SendEmailVerificationLetter.
type DomainSendEmailVerificationLetterOptions struct {
	DomainListOptions
	CheckOnly Bool `json:"check_only,omitempty"`
}

// DomainStatusAnswer is a single status answer.
type DomainStatusAnswer struct {
	Status string `json:"status"`
}

// SendEmailVerificationLetter sends (or with CheckOnly only checks) registrant email verification.
// Status is ALREADY_VERIFIED, NOT_VERIFIED, SENT or NOT_REQUIRED.
//
// API: domain/send_email_verification_letter
func (s *DomainService) SendEmailVerificationLetter(ctx context.Context, opts *DomainSendEmailVerificationLetterOptions) (*DomainStatusAnswer, *Response, error) {
	return call[DomainStatusAnswer](ctx, s.client, "domain/send_email_verification_letter", opts)
}

// DownloadCertificate returns the PDF domain ownership certificate (service srv_voucher).
// A JSON reply instead of the file is returned as *ErrorResponse; its Code is the API
// error code or the answer status (DOMAIN_NOT_FOUND, SERVICE_NOT_FOUND).
//
// API: domain/download_certificate
func (s *DomainService) DownloadCertificate(ctx context.Context, opts *ServiceIdent) ([]byte, *Response, error) {
	req, err := s.client.NewRequest(ctx, "domain/download_certificate", opts)
	if err != nil {
		return nil, nil, err
	}
	resp, err := s.client.client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	response := &Response{Response: resp}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, response, err
	}
	var env envelope
	if bytes.HasPrefix(bytes.TrimSpace(body), []byte("{")) && json.Unmarshal(body, &env) == nil && env.Result != "" {
		if env.Result != "success" {
			return nil, response, &ErrorResponse{Response: resp, Code: env.ErrorCode, Text: env.ErrorText, Params: env.ErrorParams}
		}
		var a DomainStatusAnswer
		_ = json.Unmarshal(env.Answer, &a)
		return nil, response, &ErrorResponse{Response: resp, Code: a.Status, Text: "no certificate returned"}
	}
	return body, response, nil
}

// OrderTrustee orders the Trustee service; Status is OK, ALREADY_ORDERED or NOT_AVAILABLE.
//
// API: domain/order_trustee
func (s *DomainService) OrderTrustee(ctx context.Context, opts *DomainListOptions) (*DomainStatusAnswer, *Response, error) {
	return call[DomainStatusAnswer](ctx, s.client, "domain/order_trustee", opts)
}
