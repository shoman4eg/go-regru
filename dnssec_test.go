package regru

import (
	"context"
	"testing"
)

func TestDNSSECService(t *testing.T) {
	ctx := context.Background()
	three := &DNSSECOptions{Domains: []DNSSECDomain{{ServiceIdent: ServiceIdent{DName: "test.ru"}}, {ServiceIdent: ServiceIdent{DName: "test.com"}}, {ServiceIdent: ServiceIdent{DName: "test.net"}}}}
	threeIn := map[string]any{"domains": []map[string]string{{"dname": "test.ru"}, {"dname": "test.com"}, {"dname": "test.net"}}}
	changed := `{"domains":[{"dname":"test.ru","result":"success","service_id":"12345"},{"dname":"test.com","error_params":{"domain_name":"test.com"},"error_text":"DNSSEC updating for domain already in progress","result":"error"}]}`
	firstOK := func(a *DNSSECAnswer) bool {
		return len(a.Domains) >= 2 && a.Domains[0].ServiceID == 12345 && a.Domains[0].Result == "success" && a.Domains[1].Result == "error"
	}
	tests := []struct {
		path, answer string
		opts         *DNSSECOptions
		want         any
		fn           func(*DNSSECService, context.Context, *DNSSECOptions) (*DNSSECAnswer, *Response, error)
		ok           func(*DNSSECAnswer) bool
	}{
		{"/dnssec/nop", `{"domains":[{"dname":"test.ru","result":"success","service_id":"12345"},{"dname":"test.net","domain_name":"test.net","error_params":null,"error_text":"This domain doesn't support DNSSEC","result":"error"}]}`,
			&DNSSECOptions{Domains: []DNSSECDomain{{ServiceIdent: ServiceIdent{DName: "test.ru"}}, {ServiceIdent: ServiceIdent{DName: "test.net"}}}},
			map[string]any{"domains": []map[string]string{{"dname": "test.ru"}, {"dname": "test.net"}}},
			(*DNSSECService).Nop, firstOK},
		{"/dnssec/get_status", `{"domains":[{"dname":"test.ru","result":"success","service_id":"12345","status":"enabled"},{"dname":"test.com","result":"success","service_id":"6789","status":"updating"}]}`,
			three, threeIn, (*DNSSECService).GetStatus,
			func(a *DNSSECAnswer) bool {
				return a.Domains[0].Status == "enabled" && a.Domains[1].Status == "updating"
			}},
		{"/dnssec/enable", changed, three, threeIn, (*DNSSECService).Enable, firstOK},
		{"/dnssec/disable", changed, three, threeIn, (*DNSSECService).Disable, firstOK},
		{"/dnssec/renew_ksk", changed, three, threeIn, (*DNSSECService).RenewKSK, firstOK},
		{"/dnssec/renew_zsk", changed, three, threeIn, (*DNSSECService).RenewZSK, firstOK},
		{"/dnssec/get_records", `{"domains":[{"dname":"test.ru","records":[{"DNSKEY":"test.ru. IN DNSKEY 256 3 13 CPaF"},{"DNSKEY":"test.ru. IN DNSKEY 257 3 13 X2eh","DS":"test.ru. IN DS 55528 13 2 cbb7"}],"result":"success","service_id":"12345"},{"dname":"test.com","records":[{"alg":"13","digest":"17F3","digtype":"2","keytag":"2371"}],"result":"success","service_id":"6789"}]}`,
			three, threeIn, (*DNSSECService).GetRecords,
			func(a *DNSSECAnswer) bool {
				return a.Domains[0].Records[1]["DS"] == "test.ru. IN DS 55528 13 2 cbb7" && a.Domains[1].Records[0]["keytag"] == "2371"
			}},
		{"/dnssec/add_keys", `{"domains":[{"dname":"test.ru","result":"success","service_id":"12345"},{"dname":"test.com","result":"success","service_id":"6789"}]}`,
			&DNSSECOptions{Domains: []DNSSECDomain{{ServiceIdent: ServiceIdent{DName: "test.ru"}, Records: []string{"test.ru. 3600 IN DS 2371 13 2 4508"}}}},
			map[string]any{"domains": []map[string]any{{"dname": "test.ru", "records": []string{"test.ru. 3600 IN DS 2371 13 2 4508"}}}},
			(*DNSSECService).AddKeys,
			func(a *DNSSECAnswer) bool { return a.Domains[1].ServiceID == 6789 }},
	}
	c, mux := setup(t)
	for _, tt := range tests {
		handle(t, mux, tt.path, tt.want, tt.answer)
		if got, _, err := tt.fn(c.DNSSEC, ctx, tt.opts); err != nil {
			t.Errorf("%s: %v", tt.path, err)
		} else if !tt.ok(got) {
			t.Errorf("%s: got %+v", tt.path, got)
		}
	}
}
