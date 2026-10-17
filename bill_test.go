package regru

import (
	"context"
	"testing"
)

func TestBillService(t *testing.T) {
	ctx := context.Background()
	notPayed := `{"bills":[{"payment":"100","pay_status":"notpayed","currency":"RUR","bill_date":"1917-10-26","pay_type":"bank","pay_date":null,"bill_id":"12345","orig_payment":"100","items":[{"dname":"october.org","itemtype":"service","action":"new","service_id":"12345","servtype":"domain"}]}]}`
	tests := []struct {
		path, answer string
		want         any
		call         func(*Client) (*BillsAnswer, error)
		ok           func(*BillsAnswer) bool
	}{
		{"/bill/nop", `{"bills":[{"bill_id":"12345","pay_status":"notpayed"},{"bill_id":"12346","pay_status":"payed"}]}`, map[string]any{"bills": []int{12345, 12346}},
			func(c *Client) (*BillsAnswer, error) {
				r, _, err := c.Bill.Nop(ctx, &BillOptions{Bills: []int64{12345, 12346}})
				return r, err
			},
			func(a *BillsAnswer) bool { return len(a.Bills) == 2 && a.Bills[1].PayStatus == "payed" }},
		{"/bill/get_not_payed", notPayed, map[string]any{},
			func(c *Client) (*BillsAnswer, error) {
				r, _, err := c.Bill.GetNotPayed(ctx, &BillListOptions{})
				return r, err
			},
			func(a *BillsAnswer) bool {
				b := a.Bills[0]
				return b.BillID == 12345 && b.Payment == 100 && b.Items[0].ServiceID == 12345 && b.Items[0].DName == "october.org"
			}},
		{"/bill/get_for_period", notPayed, map[string]any{"start_date": "1917-10-26", "end_date": "1917-11-07"},
			func(c *Client) (*BillsAnswer, error) {
				r, _, err := c.Bill.GetForPeriod(ctx, &BillGetForPeriodOptions{StartDate: "1917-10-26", EndDate: "1917-11-07"})
				return r, err
			},
			func(a *BillsAnswer) bool { return a.Bills[0].PayType == "bank" }},
		{"/bill/change_pay_type", `{"bills":[{"bill_id":"123456","currency":"RUR","old_pay_type":"bank","pay_status":"payed","pay_type":"prepay","payment":"100","result":"success","total_payment":"100"}]}`,
			map[string]any{"bills": []int{123456}, "pay_type": "prepay", "currency": "RUR"},
			func(c *Client) (*BillsAnswer, error) {
				r, _, err := c.Bill.ChangePayType(ctx, &BillChangePayTypeOptions{BillOptions: BillOptions{Bills: []int64{123456}}, PayType: "prepay", Currency: "RUR"})
				return r, err
			},
			func(a *BillsAnswer) bool {
				b := a.Bills[0]
				return b.OldPayType == "bank" && b.TotalPayment == 100 && b.Err() == nil
			}},
		{"/bill/delete", `{"bills":[{"bill_id":"12345","pay_status":"notpayed","result":"success","status":"deleted"},{"bill_id":"12346","error_code":"BILL_CAN_NOT_REMOVED","pay_status":"payed","status":"active"},{"bill_id":"12347","error_code":"BILL_ID_NOT_FOUND"}]}`,
			map[string]any{"bill_id": 12345},
			func(c *Client) (*BillsAnswer, error) {
				r, _, err := c.Bill.Delete(ctx, &BillOptions{BillID: 12345})
				return r, err
			},
			func(a *BillsAnswer) bool {
				return a.Bills[0].Status == "deleted" && ErrorCode(a.Bills[1].Err()) == "BILL_CAN_NOT_REMOVED" && a.Bills[2].BillID == 12347
			}},
	}
	c, mux := setup(t)
	for _, tt := range tests {
		handle(t, mux, tt.path, tt.want, tt.answer)
		if got, err := tt.call(c); err != nil {
			t.Errorf("%s: %v", tt.path, err)
		} else if !tt.ok(got) {
			t.Errorf("%s: got %+v", tt.path, got)
		}
	}
}
