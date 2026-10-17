package regru

import (
	"context"
	"encoding/json"
)

// ServiceService handles functions of the "service" category.
type ServiceService service

// ServiceListOptions identifies one service or a list of services (sections 2.2.4.3, 2.2.4.4).
type ServiceListOptions struct {
	ServiceIdent
	Domains  []ServiceIdent `json:"domains,omitempty"`
	Services []ServiceIdent `json:"services,omitempty"`
}

// ServiceListAnswer is returned by functions that answer with a per-service result list.
type ServiceListAnswer struct {
	Services []ServiceItem `json:"services"`
}

// Nop checks availability of the services and returns their IDs.
func (s *ServiceService) Nop(ctx context.Context, opts *ServiceListOptions) (*ServiceListAnswer, *Response, error) {
	return call[ServiceListAnswer](ctx, s.client, "service/nop", opts)
}

// ServiceGetPricesOptions are the parameters of GetPrices.
type ServiceGetPricesOptions struct {
	ShowRenewData Bool   `json:"show_renew_data,omitempty"`
	Currency      string `json:"currency,omitempty"` // rur, uah, usd, eur
}

// ServicePricesAnswer is returned by GetPrices.
type ServicePricesAnswer struct {
	Currency      string         `json:"currency"`
	PriceGroup    FlexString     `json:"pricegroup"`
	ShowRenewData Bool           `json:"show_renew_data"`
	Prices        []ServicePrice `json:"prices"`
}

// ServicePrice is a single price entry of GetPrices.
type ServicePrice struct {
	ServOrTLD string    `json:"servortld"`
	Unit      string    `json:"unit"`
	SubType   string    `json:"subtype"`
	Price     FlexFloat `json:"price"`
}

// GetPrices returns registration/renewal prices of services (service/get_prices).
func (s *ServiceService) GetPrices(ctx context.Context, opts *ServiceGetPricesOptions) (*ServicePricesAnswer, *Response, error) {
	return call[ServicePricesAnswer](ctx, s.client, "service/get_prices", opts)
}

// ServiceGetServtypeDetailsOptions are the parameters of GetServtypeDetails.
// ServType may hold several comma-separated types; SubType is then ignored.
type ServiceGetServtypeDetailsOptions struct {
	ServType     string `json:"servtype,omitempty"`
	SubType      string `json:"subtype,omitempty"`
	UnrollPrices Bool   `json:"unroll_prices,omitempty"`
	OnlyActual   Bool   `json:"only_actual,omitempty"`
	ShowHidden   Bool   `json:"show_hidden,omitempty"`
}

// ServiceServtypeDetailsAnswer is returned by GetServtypeDetails.
type ServiceServtypeDetailsAnswer []ServiceServtypeDetail

// ServiceServtypeDetail describes one service subtype.
// PricesNew/PricesRenew are set only with UnrollPrices; their format is undocumented.
type ServiceServtypeDetail struct {
	ServType     string          `json:"servtype"`
	SubType      string          `json:"subtype"`
	Unit         string          `json:"unit"`
	ExtParams    map[string]any  `json:"extparams"`
	IsRenewable  Bool            `json:"is_renewable"`
	CommonName   any             `json:"commonname"`
	PeriodsNew   string          `json:"periods_new"`
	PeriodsRenew string          `json:"periods_renew"`
	PriceNew     FlexFloat       `json:"price_new"`
	PriceRenew   FlexFloat       `json:"price_renew"`
	PricesNew    json.RawMessage `json:"prices_new,omitempty"`
	PricesRenew  json.RawMessage `json:"prices_renew,omitempty"`
}

// GetServtypeDetails returns prices and common data of a service type (service/get_servtype_details).
func (s *ServiceService) GetServtypeDetails(ctx context.Context, opts *ServiceGetServtypeDetailsOptions) (*ServiceServtypeDetailsAnswer, *Response, error) {
	return call[ServiceServtypeDetailsAnswer](ctx, s.client, "service/get_servtype_details", opts)
}

// ServiceCreateOptions are the parameters of Create: domain_name, servtype, period, subtype,
// user_servid, pay_type, ok_if_no_money, folder_name/folder_id, no_new_folder, comment,
// admin_comment plus the servtype-specific fields from section 7.4.
// shortcut: untyped map because each servtype has its own field set; add typed structs per servtype if needed.
type ServiceCreateOptions map[string]any

// ServiceCreateAnswer is returned by Create.
type ServiceCreateAnswer struct {
	Descr           string    `json:"descr"`
	BillID          FlexInt   `json:"bill_id"`
	Payment         FlexFloat `json:"payment"`
	PayNotes        string    `json:"pay_notes"`
	ServiceID       FlexInt   `json:"service_id"`
	UplinkServiceID FlexInt   `json:"uplink_service_id"`
}

// Create orders a new service (service/create).
func (s *ServiceService) Create(ctx context.Context, opts ServiceCreateOptions) (*ServiceCreateAnswer, *Response, error) {
	return call[ServiceCreateAnswer](ctx, s.client, "service/create", opts)
}

// ServiceCheckCreateOptions are the parameters of CheckCreate.
type ServiceCheckCreateOptions struct {
	DomainName string `json:"domain_name,omitempty"`
	ServType   string `json:"servtype,omitempty"`
	SubType    string `json:"subtype,omitempty"`
	Period     int    `json:"period,omitempty"`
}

// ServiceCheckCreateAnswer is returned by CheckCreate.
type ServiceCheckCreateAnswer struct {
	Descr string         `json:"descr"`
	Param map[string]any `json:"param"`
}

// CheckCreate validates service order parameters (service/check_create).
func (s *ServiceService) CheckCreate(ctx context.Context, opts *ServiceCheckCreateOptions) (*ServiceCheckCreateAnswer, *Response, error) {
	return call[ServiceCheckCreateAnswer](ctx, s.client, "service/check_create", opts)
}

// Delete deletes a service (service/delete).
func (s *ServiceService) Delete(ctx context.Context, opts *ServiceIdent) (*Response, error) {
	_, resp, err := call[struct{}](ctx, s.client, "service/delete", opts)
	return resp, err
}

// ServiceGetInfoOptions are the parameters of GetInfo.
type ServiceGetInfoOptions struct {
	ServiceListOptions
	ShowFolders Bool `json:"show_folders,omitempty"`
}

// ServiceInfoAnswer is returned by GetInfo and GetList.
type ServiceInfoAnswer struct {
	Services []ServiceInfo `json:"services"`
}

// ServiceInfo describes a service. User2Login is returned to partners only;
// Folders (undocumented format) only with ShowFolders.
type ServiceInfo struct {
	ServiceItem
	CreationDate    string          `json:"creation_date"`
	ExpirationDate  string          `json:"expiration_date"`
	UplinkServiceID FlexInt         `json:"uplink_service_id"`
	User2Login      string          `json:"user2_login"`
	Folders         json.RawMessage `json:"folders,omitempty"`
}

// GetInfo returns information about services (service/get_info).
func (s *ServiceService) GetInfo(ctx context.Context, opts *ServiceGetInfoOptions) (*ServiceInfoAnswer, *Response, error) {
	return call[ServiceInfoAnswer](ctx, s.client, "service/get_info", opts)
}

// ServiceGetListOptions are the parameters of GetList.
type ServiceGetListOptions struct {
	ServType string `json:"servtype,omitempty"`
}

// GetList returns the list of active services (service/get_list).
func (s *ServiceService) GetList(ctx context.Context, opts *ServiceGetListOptions) (*ServiceInfoAnswer, *Response, error) {
	return call[ServiceInfoAnswer](ctx, s.client, "service/get_list", opts)
}

// ServiceFoldersAnswer is returned by GetFolders.
type ServiceFoldersAnswer struct {
	Folders []ServiceFolder `json:"folders"`
}

// ServiceFolder is a folder the service belongs to.
type ServiceFolder struct {
	FolderID   FlexInt `json:"folder_id"`
	FolderName string  `json:"folder_name"`
}

// GetFolders returns the folders the service belongs to (service/get_folders).
func (s *ServiceService) GetFolders(ctx context.Context, opts *ServiceIdent) (*ServiceFoldersAnswer, *Response, error) {
	return call[ServiceFoldersAnswer](ctx, s.client, "service/get_folders", opts)
}

// ServiceGetDetailsOptions are the parameters of GetDetails.
type ServiceGetDetailsOptions struct {
	ServiceListOptions
	SeparateGroups   Bool `json:"separate_groups,omitempty"`
	ShowContactsOnly Bool `json:"show_contacts_only,omitempty"`
	SplittedContacts Bool `json:"splitted_contacts,omitempty"`
}

// ServiceDetailsAnswer is returned by GetDetails.
type ServiceDetailsAnswer struct {
	Services []ServiceDetails `json:"services"`
}

// ServiceDetails holds servtype-specific service data.
// shortcut: Details/Contacts are untyped since fields differ per servtype and zone; type them when a caller needs it.
type ServiceDetails struct {
	ServiceItem
	Details  map[string]any `json:"details,omitempty"`
	Contacts map[string]any `json:"contacts,omitempty"`
}

// GetDetails returns additional service data, incl. domain contacts (service/get_details).
func (s *ServiceService) GetDetails(ctx context.Context, opts *ServiceGetDetailsOptions) (*ServiceDetailsAnswer, *Response, error) {
	return call[ServiceDetailsAnswer](ctx, s.client, "service/get_details", opts)
}

// ServiceDedicatedServerListAnswer is returned by GetDedicatedServerList.
type ServiceDedicatedServerListAnswer struct {
	ServerList []ServiceDedicatedServer `json:"server_list"`
}

// ServiceDedicatedServer is a dedicated server available for order.
type ServiceDedicatedServer struct {
	ServerID    FlexInt   `json:"server_id"`
	CPUContent  string    `json:"cpu_content"`
	CPUCount    FlexInt   `json:"cpu_count"`
	CPUCore     FlexInt   `json:"cpu_core"`
	RAMContent  string    `json:"ram_content"`
	RAMCount    FlexInt   `json:"ram_count"`
	RAMSize     string    `json:"ram_size"`
	HDDContent  string    `json:"hdd_content"`
	HDDCount    FlexInt   `json:"hdd_count"`
	HDDSize     string    `json:"hdd_size"`
	MonthTraf   FlexInt   `json:"month_traf"`
	PriceRetail FlexFloat `json:"price_retail"`
}

// GetDedicatedServerList returns dedicated servers available for order (service/get_dedicated_server_list).
func (s *ServiceService) GetDedicatedServerList(ctx context.Context, opts *ServiceIdent) (*ServiceDedicatedServerListAnswer, *Response, error) {
	return call[ServiceDedicatedServerListAnswer](ctx, s.client, "service/get_dedicated_server_list", opts)
}

// ServiceUpdateOptions are the parameters of Update: dname, servtype plus the
// servtype-specific fields from section 7.12 (e.g. subtask, fwdfrom, fwdto, webfwd_type).
// shortcut: untyped map because each servtype has its own field set; add typed structs per servtype if needed.
type ServiceUpdateOptions map[string]any

// ServiceUpdateAnswer is returned by Update.
type ServiceUpdateAnswer struct {
	Descr string `json:"descr"`
}

// Update configures a service (service/update).
func (s *ServiceService) Update(ctx context.Context, opts ServiceUpdateOptions) (*ServiceUpdateAnswer, *Response, error) {
	return call[ServiceUpdateAnswer](ctx, s.client, "service/update", opts)
}

// ServiceRenewOptions are the parameters of Renew.
type ServiceRenewOptions struct {
	ServiceListOptions
	PaymentOptions
	Period int `json:"period,omitempty"`
}

// ServiceRenewAnswer is returned by Renew. Status is renew_success or only_bill_created.
// Services is set for list requests, DName/ServType for single-service ones.
type ServiceRenewAnswer struct {
	Period   FlexInt       `json:"period"`
	BillID   FlexInt       `json:"bill_id"`
	Payment  FlexFloat     `json:"payment"`
	Currency string        `json:"currency"`
	Status   string        `json:"status"`
	DName    string        `json:"dname"`
	ServType string        `json:"servtype"`
	Services []ServiceItem `json:"services"`
}

// Renew renews domains or services (service/renew).
func (s *ServiceService) Renew(ctx context.Context, opts *ServiceRenewOptions) (*ServiceRenewAnswer, *Response, error) {
	return call[ServiceRenewAnswer](ctx, s.client, "service/renew", opts)
}

// ServiceBillsAnswer is returned by GetBills.
type ServiceBillsAnswer struct {
	Services []ServiceBills `json:"services"`
}

// ServiceBills lists bill IDs of a service.
type ServiceBills struct {
	ServiceItem
	Bills []FlexInt `json:"bills"`
}

// GetBills returns bills related to the services (service/get_bills).
func (s *ServiceService) GetBills(ctx context.Context, opts *ServiceListOptions) (*ServiceBillsAnswer, *Response, error) {
	return call[ServiceBillsAnswer](ctx, s.client, "service/get_bills", opts)
}

// ServiceReadonlyAnswer is returned by CheckReadonly.
type ServiceReadonlyAnswer struct {
	Services []ServiceReadonly `json:"services"`
}

// ServiceReadonly tells whether a service may be edited.
type ServiceReadonly struct {
	ServiceItem
	ErrorDetail string `json:"error_detail"`
	IsReadonly  Bool   `json:"is_readonly"`
}

// CheckReadonly checks whether the services may be edited (service/check_readonly).
func (s *ServiceService) CheckReadonly(ctx context.Context, opts *ServiceListOptions) (*ServiceReadonlyAnswer, *Response, error) {
	return call[ServiceReadonlyAnswer](ctx, s.client, "service/check_readonly", opts)
}

// ServiceSetAutorenewFlagOptions are the parameters of SetAutorenewFlag.
type ServiceSetAutorenewFlagOptions struct {
	ServiceIdent
	FlagValue Bool `json:"flag_value"` // no omitempty: false must be sent as 0
}

// SetAutorenewFlag sets or clears the auto-renewal flag (service/set_autorenew_flag).
func (s *ServiceService) SetAutorenewFlag(ctx context.Context, opts *ServiceSetAutorenewFlagOptions) (*Response, error) {
	_, resp, err := call[struct{}](ctx, s.client, "service/set_autorenew_flag", opts)
	return resp, err
}

// Suspend suspends a service; for a domain, removes delegation (service/suspend).
func (s *ServiceService) Suspend(ctx context.Context, opts *ServiceIdent) (*Response, error) {
	_, resp, err := call[struct{}](ctx, s.client, "service/suspend", opts)
	return resp, err
}

// Resume resumes a suspended service; for a domain, delegates it (service/resume).
func (s *ServiceService) Resume(ctx context.Context, opts *ServiceIdent) (*Response, error) {
	_, resp, err := call[struct{}](ctx, s.client, "service/resume", opts)
	return resp, err
}

// ServiceDepreciatedPeriodAnswer is returned by GetDepreciatedPeriod.
type ServiceDepreciatedPeriodAnswer struct {
	DepreciatedPeriod FlexFloat `json:"depreciated_period"`
}

// GetDepreciatedPeriod returns the number of periods left until expiration (service/get_depreciated_period).
func (s *ServiceService) GetDepreciatedPeriod(ctx context.Context, opts *ServiceIdent) (*ServiceDepreciatedPeriodAnswer, *Response, error) {
	return call[ServiceDepreciatedPeriodAnswer](ctx, s.client, "service/get_depreciated_period", opts)
}

// ServiceUpgradeOptions are the parameters of Upgrade. ServiceIdent.SubType is the new subtype (tariff).
type ServiceUpgradeOptions struct {
	ServiceIdent
	Period int `json:"period,omitempty"`
}

// ServiceUpgradeAnswer is returned by Upgrade.
type ServiceUpgradeAnswer struct {
	BillID           FlexInt   `json:"bill_id"`
	NewServiceID     FlexInt   `json:"new_service_id"`
	ReturnedAmount   FlexFloat `json:"returned_amount"`
	WithdrawedAmount FlexFloat `json:"withdrawed_amount"`
}

// Upgrade changes the service subtype (tariff) (service/upgrade).
func (s *ServiceService) Upgrade(ctx context.Context, opts *ServiceUpgradeOptions) (*ServiceUpgradeAnswer, *Response, error) {
	return call[ServiceUpgradeAnswer](ctx, s.client, "service/upgrade", opts)
}

// ServicePartcontrolGrantOptions are the parameters of PartcontrolGrant.
type ServicePartcontrolGrantOptions struct {
	ServiceIdent
	NewLogin string `json:"newlogin,omitempty"`
}

// ServicePartcontrolAnswer is returned by PartcontrolGrant and PartcontrolRevoke.
type ServicePartcontrolAnswer struct {
	NewLogin  string  `json:"newlogin"`
	ServiceID FlexInt `json:"service_id"`
}

// PartcontrolGrant grants partial control of a service to another user (service/partcontrol_grant).
func (s *ServiceService) PartcontrolGrant(ctx context.Context, opts *ServicePartcontrolGrantOptions) (*ServicePartcontrolAnswer, *Response, error) {
	return call[ServicePartcontrolAnswer](ctx, s.client, "service/partcontrol_grant", opts)
}

// PartcontrolRevoke revokes partial control of a service (service/partcontrol_revoke).
func (s *ServiceService) PartcontrolRevoke(ctx context.Context, opts *ServiceIdent) (*ServicePartcontrolAnswer, *Response, error) {
	return call[ServicePartcontrolAnswer](ctx, s.client, "service/partcontrol_revoke", opts)
}

// ServiceResendMailOptions are the parameters of ResendMail.
type ServiceResendMailOptions struct {
	ServiceIdent
	MailType string `json:"mailtype,omitempty"` // SSL only: approver_email or certificate_email
}

// ResendMail resends the service email message (service/resend_mail).
func (s *ServiceService) ResendMail(ctx context.Context, opts *ServiceResendMailOptions) (*ServiceItem, *Response, error) {
	return call[ServiceItem](ctx, s.client, "service/resend_mail", opts)
}

// ServiceRefillOptions are the parameters of Refill (amount/currency appear only in the doc example).
type ServiceRefillOptions struct {
	ServiceIdent
	Amount   float64 `json:"amount,omitempty"`
	Currency string  `json:"currency,omitempty"`
}

// ServiceRefillAnswer is returned by Refill.
type ServiceRefillAnswer struct {
	BillID    FlexInt   `json:"bill_id"`
	PayType   string    `json:"pay_type"`
	Currency  string    `json:"currency"`
	Payment   FlexFloat `json:"payment"`
	PayNotes  string    `json:"pay_notes"`
	ServiceID FlexInt   `json:"service_id"`
}

// Refill tops up the balance of an external service, e.g. srv_seowizard (service/refill).
func (s *ServiceService) Refill(ctx context.Context, opts *ServiceRefillOptions) (*ServiceRefillAnswer, *Response, error) {
	return call[ServiceRefillAnswer](ctx, s.client, "service/refill", opts)
}

// ServiceBalanceAnswer is returned by GetBalance.
type ServiceBalanceAnswer struct {
	Balance  FlexFloat `json:"balance"`
	Currency string    `json:"currency"`
}

// GetBalance returns the balance of an external service, e.g. srv_seowizard (service/get_balance).
func (s *ServiceService) GetBalance(ctx context.Context, opts *ServiceIdent) (*ServiceBalanceAnswer, *Response, error) {
	return call[ServiceBalanceAnswer](ctx, s.client, "service/get_balance", opts)
}

// ServiceLinkAnswer is returned by SeowizardManageLink and GetWebsitebuilderLink.
type ServiceLinkAnswer struct {
	Link string `json:"link"`
}

// SeowizardManageLink returns the control panel link of srv_seowizard (service/seowizard_manage_link).
func (s *ServiceService) SeowizardManageLink(ctx context.Context, opts *ServiceIdent) (*ServiceLinkAnswer, *Response, error) {
	return call[ServiceLinkAnswer](ctx, s.client, "service/seowizard_manage_link", opts)
}

// GetWebsitebuilderLink returns the website builder publishing link (service/get_websitebuilder_link).
func (s *ServiceService) GetWebsitebuilderLink(ctx context.Context, opts *ServiceIdent) (*ServiceLinkAnswer, *Response, error) {
	return call[ServiceLinkAnswer](ctx, s.client, "service/get_websitebuilder_link", opts)
}
