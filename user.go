package regru

import (
	"context"
	"encoding/json"
)

// UserService handles functions of the "user" category.
type UserService service

// Nop tests access to the user category (user/nop).
func (s *UserService) Nop(ctx context.Context) (*map[string]any, *Response, error) {
	return call[map[string]any](ctx, s.client, "user/nop", nil)
}

// UserCreateOptions are the parameters of user/create.
// Set Login and Password, or EmailOnly for registration by email only.
type UserCreateOptions struct {
	Login              string `json:"user_login,omitempty"`
	Password           string `json:"user_password,omitempty"`
	EmailOnly          Bool   `json:"email_only,omitempty"`
	Email              string `json:"user_email,omitempty"`
	IP                 string `json:"user_ip,omitempty"`
	DefaultCountryCode string `json:"default_country_code,omitempty"`
	CountryCode        string `json:"user_country_code,omitempty"` // used in the doc example, not in the field table

	FirstName string `json:"user_first_name,omitempty"`
	LastName  string `json:"user_last_name,omitempty"`
	Company   string `json:"user_company,omitempty"`
	JabberID  string `json:"user_jabber_id,omitempty"`
	ICQ       string `json:"user_icq,omitempty"`
	Phone     string `json:"user_phone,omitempty"` // +7.4951234567
	Fax       string `json:"user_fax,omitempty"`
	Addr      string `json:"user_addr,omitempty"`
	City      string `json:"user_city,omitempty"`
	State     string `json:"user_state,omitempty"`
	Postcode  string `json:"user_postcode,omitempty"`
	WMID      string `json:"user_wmid,omitempty"`
	Website   string `json:"user_website,omitempty"`
	Language  string `json:"user_language,omitempty"` // ru or en

	Subscribe       Bool     `json:"user_subsribe,omitempty"`   // sic: API spelling
	MailNotify      *Bool    `json:"user_mailnotify,omitempty"` // API default is 1
	SetMeAsReferrer Bool     `json:"set_me_as_referrer,omitempty"`
	CheckOnly       Bool     `json:"check_only,omitempty"`
	WhiteListIPs    []string `json:"white_list_ips,omitempty"`
}

// UserCreateAnswer is returned by Create.
type UserCreateAnswer struct {
	UserID FlexInt `json:"user_id"`
}

// Create registers a new user (user/create). Resellers only.
func (s *UserService) Create(ctx context.Context, opts *UserCreateOptions) (*UserCreateAnswer, *Response, error) {
	return call[UserCreateAnswer](ctx, s.client, "user/create", opts)
}

// UserGetStatisticsOptions are the parameters of user/get_statistics.
type UserGetStatisticsOptions struct {
	DateFrom string `json:"date_from,omitempty"`
	DateTill string `json:"date_till,omitempty"`
}

// UserStatistics is returned by GetStatistics.
type UserStatistics struct {
	CostsForPeriod          FlexFloat `json:"costs_for_period"`
	ActiveDomainsCnt        FlexInt   `json:"active_domains_cnt"`
	ActiveDomainsGetCtrlCnt FlexInt   `json:"active_domains_get_ctrl_cnt"`
	RenewDomainsCnt         FlexInt   `json:"renew_domains_cnt"`
	RenewDomainsGetCtrlCnt  FlexInt   `json:"renew_domains_get_ctrl_cnt"`
	TransInDomainsCnt       FlexInt   `json:"trans_in_domains_cnt"`
	UndelegatedDomainsCnt   FlexInt   `json:"undelegated_domains_cnt"`
	RegDomainsCnt           FlexInt   `json:"reg_domains_cnt"`
	DomainFoldersCnt        FlexInt   `json:"domain_folders_cnt"`
}

// GetStatistics returns account statistics (user/get_statistics).
func (s *UserService) GetStatistics(ctx context.Context, opts *UserGetStatisticsOptions) (*UserStatistics, *Response, error) {
	return call[UserStatistics](ctx, s.client, "user/get_statistics", opts)
}

// UserGetBalanceOptions are the parameters of user/get_balance.
type UserGetBalanceOptions struct {
	Currency string `json:"currency,omitempty"` // RUR (default), USD, EUR, UAH
}

// UserBalance is returned by GetBalance.
type UserBalance struct {
	Currency string    `json:"currency"`
	Prepay   FlexFloat `json:"prepay"`
	Blocked  FlexFloat `json:"blocked"`
	Credit   FlexFloat `json:"credit"`
}

// GetBalance returns the account balance (user/get_balance).
func (s *UserService) GetBalance(ctx context.Context, opts *UserGetBalanceOptions) (*UserBalance, *Response, error) {
	return call[UserBalance](ctx, s.client, "user/get_balance", opts)
}

// UserResellerURLOptions are the parameters of user/set_reseller_url and user/get_reseller_url.
type UserResellerURLOptions struct {
	ServType string `json:"servtype"`      // srv_seowizard
	URLType  string `json:"url_type"`      // refill or refund
	URL      string `json:"url,omitempty"` // set only; may contain <service_id> and <email> tokens
}

// UserResellerURL is returned by GetResellerURL.
type UserResellerURL struct {
	URL string `json:"url"`
}

// SetResellerURL sets the redirect URL for an external service (user/set_reseller_url). Resellers only.
func (s *UserService) SetResellerURL(ctx context.Context, opts *UserResellerURLOptions) (*map[string]any, *Response, error) {
	return call[map[string]any](ctx, s.client, "user/set_reseller_url", opts)
}

// GetResellerURL returns the redirect URL for an external service (user/get_reseller_url). Resellers only.
func (s *UserService) GetResellerURL(ctx context.Context, opts *UserResellerURLOptions) (*UserResellerURL, *Response, error) {
	return call[UserResellerURL](ctx, s.client, "user/get_reseller_url", opts)
}

// UserPerson is an element of UserPersons.
type UserPerson struct {
	PersonID       FlexInt `json:"person_id"`
	UserID         FlexInt `json:"user_id"`
	Type           string  `json:"type"` // pp or org
	IsEntrepreneur Bool    `json:"is_entrepreneur"`
	ESIAIdentified Bool    `json:"esia_identified"`
	// shortcut: raw contact fields, upgrade to a typed struct when the contact set is shared with domain contacts.
	PersonContacts json.RawMessage `json:"person_contacts"`
}

// UserPersons is returned by GetPersons.
type UserPersons struct {
	Persons []UserPerson `json:"persons"`
}

// GetPersons returns the user's persons and their ESIA identification status (user/get_persons).
func (s *UserService) GetPersons(ctx context.Context) (*UserPersons, *Response, error) {
	return call[UserPersons](ctx, s.client, "user/get_persons", nil)
}
