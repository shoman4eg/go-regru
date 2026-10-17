package regru

import (
	"context"
	"encoding/json"
	"testing"
)

func TestZoneService(t *testing.T) {
	const twoAnswer = `{"domains":[{"dname":"test.ru","result":"success","service_id":"12345"},{"dname":"test.com","result":"success","service_id":"12346"}]}`
	const oneAnswer = `{"domains":[{"dname":"test.ru","result":"success","service_id":"12345"}]}`
	const twoIn = `{"domains":[{"dname":"test.ru"},{"dname":"test.com"}]}`
	two := ZoneDomainsOptions{Domains: []ServiceIdent{{DName: "test.ru"}, {DName: "test.com"}}}
	one := ZoneDomainsOptions{Domains: []ServiceIdent{{DName: "test.ru"}}}
	one1, zero := 1, 0

	ctx := context.Background()
	tests := []struct {
		path, input, answer string
		call                func(*ZoneService) (any, error)
		check               func(any) bool
	}{
		{"zone/nop", twoIn, twoAnswer,
			func(s *ZoneService) (any, error) { a, _, err := s.Nop(ctx, &two); return a, err }, nil},
		{"zone/add_alias", `{"domains":[{"dname":"test.ru"},{"dname":"test.com"}],"subdomain":"@","ipaddr":"111.111.111.111"}`, twoAnswer,
			func(s *ZoneService) (any, error) {
				a, _, err := s.AddAlias(ctx, &ZoneAddAliasOptions{two, "@", "111.111.111.111"})
				return a, err
			}, nil},
		{"zone/add_aaaa", `{"domains":[{"dname":"test.ru"},{"dname":"test.com"}],"subdomain":"@","ipaddr":"aa11::a111:11aa:aaa1:aa1a"}`, twoAnswer,
			func(s *ZoneService) (any, error) {
				a, _, err := s.AddAAAA(ctx, &ZoneAddAliasOptions{two, "@", "aa11::a111:11aa:aaa1:aa1a"})
				return a, err
			}, nil},
		{"zone/add_caa", `{"domains":[{"dname":"test.ru"}],"subdomain":"@","flags":0,"tag":"issuewild","value":"ca.example.com"}`, oneAnswer,
			func(s *ZoneService) (any, error) {
				a, _, err := s.AddCAA(ctx, &ZoneAddCAAOptions{one, "@", 0, "issuewild", "ca.example.com"})
				return a, err
			}, nil},
		{"zone/add_cname", `{"domains":[{"dname":"test.ru"},{"dname":"test.com"}],"subdomain":"mail","canonical_name":"mx10.test.ru"}`, twoAnswer,
			func(s *ZoneService) (any, error) {
				a, _, err := s.AddCNAME(ctx, &ZoneAddCNAMEOptions{two, "mail", "mx10.test.ru"})
				return a, err
			}, nil},
		{"zone/add_https", `{"domains":[{"dname":"test.ru"}],"subdomain":"@","priority":1,"target":".","value":"alpn=h3 ipv4hint=193.12.11.10,193.12.11.11"}`, oneAnswer,
			func(s *ZoneService) (any, error) {
				a, _, err := s.AddHTTPS(ctx, &ZoneAddHTTPSOptions{one, "@", 1, ".", "alpn=h3 ipv4hint=193.12.11.10,193.12.11.11"})
				return a, err
			}, nil},
		{"zone/add_mx", `{"domains":[{"dname":"test.ru"},{"dname":"test.com"}],"subdomain":"@","mail_server":"mail"}`, twoAnswer,
			func(s *ZoneService) (any, error) {
				a, _, err := s.AddMX(ctx, &ZoneAddMXOptions{ZoneDomainsOptions: two, Subdomain: "@", MailServer: "mail"})
				return a, err
			}, nil},
		// The doc sends numbers as strings ("10", "0", "5060"); we send JSON numbers.
		{"zone/add_ns", `{"domains":[{"dname":"test.ru"},{"dname":"test.com"}],"subdomain":"tt","dns_server":"ns1.test.ru","record_number":10}`, twoAnswer,
			func(s *ZoneService) (any, error) {
				a, _, err := s.AddNS(ctx, &ZoneAddNSOptions{two, "tt", "ns1.test.ru", 10})
				return a, err
			}, nil},
		{"zone/add_txt", `{"domains":[{"dname":"test.ru"},{"dname":"test.com"}],"subdomain":"mail","text":"testmail"}`, twoAnswer,
			func(s *ZoneService) (any, error) {
				a, _, err := s.AddTXT(ctx, &ZoneAddTXTOptions{two, "mail", "testmail"})
				return a, err
			}, nil},
		{"zone/add_srv", `{"domains":[{"dname":"test.ru"},{"dname":"test.com"}],"service":"_sip._udp","priority":0,"port":5060,"target":"sip.test.ru"}`, twoAnswer,
			func(s *ZoneService) (any, error) {
				a, _, err := s.AddSRV(ctx, &ZoneAddSRVOptions{ZoneDomainsOptions: two, Service: "_sip._udp", Target: "sip.test.ru", Port: 5060})
				return a, err
			}, nil},
		{"zone/get_resource_records", twoIn,
			`{"domains":[{"dname":"test.ru","result":"success","rrs":[{"content":"111.222.111.222","prio":"0","rectype":"A","state":"A","subname":"www"}],"service_id":"12345","soa":{"minimum_ttl":"12h","ttl":"1d"}}]}`,
			func(s *ZoneService) (any, error) { a, _, err := s.GetResourceRecords(ctx, &two); return a, err },
			func(a any) bool {
				d := a.(*ZoneRecordsAnswer).Domains[0]
				return d.ServiceID == 12345 && d.RRs[0].Content == "111.222.111.222" && d.RRs[0].RecType == "A" && d.SOA.MinimumTTL == "12h"
			}},
		{"zone/update_records",
			`{"domains":[{"dname":"test.ru","action_list":[{"action":"add_alias","subdomain":"www","ipaddr":"11.22.33.44"},{"action":"add_cname","subdomain":"@","canonical_name":"www.test.ru"},{"action":"add_https","subdomain":"@","priority":0,"target":"x.test.ru"},{"action":"add_caa","subdomain":"@","flags":1,"tag":"issue","value":"ca.example.com"}]}]}`,
			`{"domains":[{"dname":"test.ru","service_id":12345,"action_list":[{"action":"add_alias","result":"success"},{"action":"add_cname","result":"success"}],"result":"success"},{"dname":"test.ru","service_id":12345,"error_params":{"action":"add_alias"},"error_code":"INVALID_ACTION","error_text":"Invalid action: add_alias"}]}`,
			func(s *ZoneService) (any, error) {
				a, _, err := s.UpdateRecords(ctx, &ZoneUpdateRecordsOptions{Domains: []ZoneUpdateRecordsDomain{{
					ServiceIdent: ServiceIdent{DName: "test.ru"},
					ActionList: []ZoneAction{
						{Action: "add_alias", Subdomain: "www", IPAddr: "11.22.33.44"},
						{Action: "add_cname", Subdomain: "@", CanonicalName: "www.test.ru"},
						{Action: "add_https", Subdomain: "@", Priority: &zero, Target: "x.test.ru"},
						{Action: "add_caa", Subdomain: "@", Flags: &one1, Tag: "issue", Value: "ca.example.com"},
					},
				}}})
				return a, err
			},
			func(a any) bool {
				d := a.(*ZoneUpdateRecordsAnswer).Domains
				return d[0].ActionList[1].Action == "add_cname" && d[0].ActionList[1].Err() == nil &&
					ErrorCode(d[1].Err()) == "INVALID_ACTION" && d[1].ErrorParams["action"] == "add_alias"
			}},
		{"zone/update_soa", `{"domains":[{"dname":"test.ru"},{"dname":"test.com"}],"ttl":"1d","minimum_ttl":"4h"}`, twoAnswer,
			func(s *ZoneService) (any, error) {
				a, _, err := s.UpdateSOA(ctx, &ZoneUpdateSOAOptions{two, "1d", "4h"})
				return a, err
			}, nil},
		{"zone/tune_forwarding", twoIn, twoAnswer,
			func(s *ZoneService) (any, error) { a, _, err := s.TuneForwarding(ctx, &two); return a, err }, nil},
		{"zone/clear_forwarding", twoIn, twoAnswer,
			func(s *ZoneService) (any, error) { a, _, err := s.ClearForwarding(ctx, &two); return a, err }, nil},
		{"zone/tune_parking", twoIn, twoAnswer,
			func(s *ZoneService) (any, error) { a, _, err := s.TuneParking(ctx, &two); return a, err }, nil},
		{"zone/clear_parking", twoIn, twoAnswer,
			func(s *ZoneService) (any, error) { a, _, err := s.ClearParking(ctx, &two); return a, err }, nil},
		// The doc example puts both dnames into one object; we use a proper list.
		{"zone/remove_record", `{"domains":[{"dname":"test.ru"},{"dname":"test.com"}],"subdomain":"@","content":"111.111.111.111","record_type":"A"}`, twoAnswer,
			func(s *ZoneService) (any, error) {
				a, _, err := s.RemoveRecord(ctx, &ZoneRemoveRecordOptions{ZoneDomainsOptions: two, Subdomain: "@", RecordType: "A", Content: "111.111.111.111"})
				return a, err
			}, nil},
		{"zone/clear", `{"service_id":12345}`, oneAnswer,
			func(s *ZoneService) (any, error) {
				a, _, err := s.Clear(ctx, &ZoneDomainsOptions{ServiceIdent: ServiceIdent{ServiceID: 12345}})
				return a, err
			}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			c, mux := setup(t)
			handle(t, mux, "/"+tt.path, json.RawMessage(tt.input), tt.answer)
			got, err := tt.call(c.Zone)
			if err != nil {
				t.Fatal(err)
			}
			check := tt.check
			if check == nil {
				check = func(a any) bool {
					d := a.(*ZoneAnswer).Domains
					return d[0].DName == "test.ru" && d[0].ServiceID == 12345 && d[0].Err() == nil
				}
			}
			if !check(got) {
				t.Errorf("unexpected answer %+v", got)
			}
		})
	}
}

// The test account and real accounts swap number/string encodings of prio and service_id.
func TestZoneRecordsTestVsRealAccount(t *testing.T) {
	for name, answer := range map[string]string{
		"real": `{"domains":[{"dname":"example.xyz","result":"success","service_id":"12345","servtype":"domain",
			"rrs":[{"content":"192.0.2.1","prio":10,"rectype":"MX","state":"A","subname":"@"}],"soa":{"minimum_ttl":"10m","ttl":"10m"}}]}`,
		"test": `{"domains":[{"dname":"example.xyz","result":"success","service_id":12345,"servtype":"domain",
			"rrs":[{"content":"192.0.2.1","prio":"10","rectype":"MX","state":"A","subname":"@"}],"soa":{"minimum_ttl":"12h","ttl":"1d"}}]}`,
	} {
		t.Run(name, func(t *testing.T) {
			var got ZoneRecordsAnswer
			if err := json.Unmarshal([]byte(answer), &got); err != nil {
				t.Fatal(err)
			}
			d := got.Domains[0]
			if d.ServiceID != 12345 || d.RRs[0].Prio != 10 || d.Err() != nil {
				t.Errorf("got %+v", d)
			}
		})
	}
}
