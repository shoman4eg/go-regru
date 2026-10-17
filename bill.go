package regru

import "context"

// BillService handles functions of the "bill" category.
type BillService service

// BillOptions identify one bill (BillID) or a list of bills (Bills).
type BillOptions struct {
	BillID int64   `json:"bill_id,omitempty"`
	Bills  []int64 `json:"bills,omitempty"`
}

// BillItem is a line of a bill.
type BillItem struct {
	ItemType  string  `json:"itemtype"` // prepayment or service
	DName     string  `json:"dname"`
	ServType  string  `json:"servtype"`
	ServiceID FlexInt `json:"service_id"`
	Action    string  `json:"action"` // new or renew
}

// Bill is an entry of BillsAnswer; which fields are set depends on the function.
type Bill struct {
	ItemResult
	BillID       FlexInt    `json:"bill_id"`
	BillDate     string     `json:"bill_date"`
	PayDate      string     `json:"pay_date"`
	Currency     string     `json:"currency"`
	Payment      FlexFloat  `json:"payment"`
	OrigPayment  FlexFloat  `json:"orig_payment"`
	TotalPayment FlexFloat  `json:"total_payment"`
	PayType      string     `json:"pay_type"`
	OldPayType   string     `json:"old_pay_type"`
	PayStatus    string     `json:"pay_status"` // notpayed, confirmed, payed, cancelled
	Status       string     `json:"status"`     // delete: deleted or active
	Items        []BillItem `json:"items"`
}

// BillsAnswer is returned by all bill functions.
type BillsAnswer struct {
	Bills []Bill `json:"bills"`
}

// Nop tests access to bills (bill/nop).
func (s *BillService) Nop(ctx context.Context, opts *BillOptions) (*BillsAnswer, *Response, error) {
	return call[BillsAnswer](ctx, s.client, "bill/nop", opts)
}

// BillListOptions are the paging parameters of bill/get_not_payed.
type BillListOptions struct {
	Limit  int `json:"limit,omitempty"` // default 100, max 1024
	Offset int `json:"offset,omitempty"`
}

// GetNotPayed returns unpaid bills (bill/get_not_payed).
func (s *BillService) GetNotPayed(ctx context.Context, opts *BillListOptions) (*BillsAnswer, *Response, error) {
	return call[BillsAnswer](ctx, s.client, "bill/get_not_payed", opts)
}

// BillGetForPeriodOptions are the parameters of bill/get_for_period.
type BillGetForPeriodOptions struct {
	StartDate string `json:"start_date"` // YYYY-MM-DD
	EndDate   string `json:"end_date"`
	PayType   string `json:"pay_type,omitempty"`
	All       Bool   `json:"all,omitempty"`
	BillListOptions
}

// GetForPeriod returns bills for a period (bill/get_for_period). Resellers only.
func (s *BillService) GetForPeriod(ctx context.Context, opts *BillGetForPeriodOptions) (*BillsAnswer, *Response, error) {
	return call[BillsAnswer](ctx, s.client, "bill/get_for_period", opts)
}

// BillChangePayTypeOptions are the parameters of bill/change_pay_type.
type BillChangePayTypeOptions struct {
	BillOptions
	PayType  string `json:"pay_type"` // prepay, yamoney, bank
	Currency string `json:"currency"` // RUR or USD
}

// ChangePayType changes the payment method of bills (bill/change_pay_type); prepay pays immediately.
func (s *BillService) ChangePayType(ctx context.Context, opts *BillChangePayTypeOptions) (*BillsAnswer, *Response, error) {
	return call[BillsAnswer](ctx, s.client, "bill/change_pay_type", opts)
}

// Delete deletes unpaid bills (bill/delete).
func (s *BillService) Delete(ctx context.Context, opts *BillOptions) (*BillsAnswer, *Response, error) {
	return call[BillsAnswer](ctx, s.client, "bill/delete", opts)
}
