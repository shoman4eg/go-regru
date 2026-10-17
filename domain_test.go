package regru

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

func TestDomainService(t *testing.T) {
	ctx := context.Background()
	m := func(kv ...any) map[string]any {
		r := map[string]any{}
		for i := 0; i < len(kv); i += 2 {
			r[kv[i].(string)] = kv[i+1]
		}
		return r
	}
	dnames := func(names ...string) []any {
		var r []any
		for _, n := range names {
			r = append(r, m("dname", n))
		}
		return r
	}
	idents := func(names ...string) []ServiceIdent {
		var r []ServiceIdent
		for _, n := range names {
			r = append(r, ServiceIdent{DName: n})
		}
		return r
	}
	contacts := DomainContacts{"country": "RU", "e_mail": "test@test.ru", "person": "Svyatoslav V Ryurik"}
	nss := DomainNSS{"ns0": "ns1.reg.ru", "ns1": "ns2.reg.ru"}
	order := `{"bill_id":"1234","domains":[{"dname":"vschizh.ru","result":"success","service_id":"12345"}],"pay_notes":"Amount successfully charged","pay_type":"prepay","payment":"600"}`
	twoDomains := `{"domains":[{"dname":"test.ru","result":"success","service_id":"12345"},{"dname":"test.su","result":"success","service_id":"12346"}]}`
	nsResp := `{"resp":{"actionstatus":"success","actiontype":"AddCns","status":"success"}}`

	type C = *Client
	tests := []struct {
		path   string
		want   any
		answer string
		call   func(C) (any, error)
		check  func(any) bool
	}{
		{"nop", m("dname", "qqq.ru"), `{"service_id":"123"}`,
			func(c C) (any, error) { a, _, err := c.Domain.Nop(ctx, &ServiceIdent{DName: "qqq.ru"}); return a, err },
			func(a any) bool { return a.(*ServiceIDAnswer).ServiceID == 123 }},
		{"get_prices", m("currency", "RUR"), `{"currency":"RUR","price_group":"Retail","prices":{"__idn.com":{"reg_max_period":"10","reg_min_period":"1","reg_price":"960"},"рф":{"reg_max_period":"1","reg_min_period":"1","reg_price":"1200"}}}`,
			func(c C) (any, error) {
				a, _, err := c.Domain.GetPrices(ctx, &DomainGetPricesOptions{Currency: "RUR"})
				return a, err
			},
			func(a any) bool {
				p := a.(*DomainGetPricesAnswer).Prices
				return p["рф"].RegPrice == 1200 && p["__idn.com"].RegMaxPeriod == 10
			}},
		{"get_suggest", m("word", "дом", "additional_word", "новый", "tlds", []any{"рф"}, "use_hyphen", 1), `{"suggestions":[{"avail_in":["рф","su"],"name":"дом"}]}`,
			func(c C) (any, error) {
				a, _, err := c.Domain.GetSuggest(ctx, &DomainGetSuggestOptions{Word: "дом", AdditionalWord: "новый", TLDs: []string{"рф"}, UseHyphen: true})
				return a, err
			},
			func(a any) bool {
				s := a.(*DomainGetSuggestAnswer).Suggestions
				return s[0].Name == "дом" && s[0].AvailIn[1] == "su"
			}},
		{"get_premium_prices", m("currency", "USD", "domains", []any{"ua.ae.org", "007.academy"}), `{"currency":"USD","domains":[{"domain_name":"007.academy","register_period":"1","register_price":"238.16","renew_period":"1","renew_price":"238.16"}]}`,
			func(c C) (any, error) {
				a, _, err := c.Domain.GetPremiumPrices(ctx, &DomainGetPremiumPricesOptions{Currency: "USD", Domains: []string{"ua.ae.org", "007.academy"}})
				return a, err
			},
			func(a any) bool {
				d := a.(*DomainGetPremiumPricesAnswer).Domains[0]
				return d.RegisterPrice == 238.16 && d.RegisterPeriod == 1
			}},
		{"get_deleted", m("min_pr", 1, "tlds", []any{"ru"}), `{"domains":[{"domain_name":"apmatypa.ru","registered":"NOT REGISTERED","yandex_tic":0,"first_create_date":"2006-05-18","google_pr":"2","date_delete":"2013-09-26"}]}`,
			func(c C) (any, error) {
				a, _, err := c.Domain.GetDeleted(ctx, &DomainGetDeletedOptions{MinPR: 1, TLDs: []string{"ru"}})
				return a, err
			},
			func(a any) bool {
				d := a.(*DomainGetDeletedAnswer).Domains[0]
				return d.GooglePR == 2 && d.Registered == "NOT REGISTERED"
			}},
		{"check", m("domains", dnames("ya.ru", "rr.ru.com")), `{"domains":[{"dname":"ya.ru","error_code":"DOMAIN_ALREADY_EXISTS","result":"Domain already exists, use whois service"},{"currency":"RUR","dname":"rr.ru.com","is_premium":"1","price":"60000","renew_price":"600","result":"Available"}]}`,
			func(c C) (any, error) {
				a, _, err := c.Domain.Check(ctx, &DomainCheckOptions{Domains: idents("ya.ru", "rr.ru.com")})
				return a, err
			},
			func(a any) bool {
				d := a.(*DomainCheckAnswer).Domains
				return ErrorCode(d[0].Err()) == "DOMAIN_ALREADY_EXISTS" && d[1].Result == "Available" && bool(d[1].IsPremium) && d[1].Price == 60000
			}},
		{"create", m("profile_type", "RU.PP", "profile_name", "my_like_ru_profile", "domain_name", "vschizh.su", "nss", m("ns0", "ns1.reg.ru", "ns1", "ns2.reg.ru")), order,
			func(c C) (any, error) {
				a, _, err := c.Domain.Create(ctx, &DomainCreateOptions{ProfileType: "RU.PP", ProfileName: "my_like_ru_profile", DomainName: "vschizh.su", NSS: nss})
				return a, err
			},
			func(a any) bool {
				o := a.(*DomainOrderAnswer)
				return o.BillID == 1234 && o.Payment == 600 && o.Domains[0].ServiceID == 12345
			}},
		{"transfer", m("domain_name", "test.com", "authinfo", "secret", "period", 1, "pay_type", "prepay"), order,
			func(c C) (any, error) {
				a, _, err := c.Domain.Transfer(ctx, &DomainTransferOptions{
					DomainCreateOptions: DomainCreateOptions{DomainName: "test.com", Period: 1, PaymentOptions: PaymentOptions{PayType: "prepay"}},
					AuthInfo:            "secret",
				})
				return a, err
			},
			func(a any) bool { return a.(*DomainOrderAnswer).PayType == "prepay" }},
		{"get_transfer_status", m("domains", dnames("ya.ru", "qwerty.ooo")), `{"domains":[{"dname":"ya.ru","error_code":"DOMAIN_NOT_FOUND","error_params":{"domain_name":"ya.ru"},"error_text":"Domain ya.ru not found or not owned by You","result":"error"},{"dname":"qwerty.ooo","result":"success","service_id":"12345","status":"notAvailable"}]}`,
			func(c C) (any, error) {
				a, _, err := c.Domain.GetTransferStatus(ctx, &DomainListOptions{Domains: idents("ya.ru", "qwerty.ooo")})
				return a, err
			},
			func(a any) bool {
				d := a.(*DomainGetTransferStatusAnswer).Domains
				return d[0].Err() != nil && d[1].Status == "notAvailable" && d[1].ServiceID == 12345
			}},
		{"set_new_authinfo", m("authinfo", "hjgf7jfj8f8", "dname", "test.com"), `{"dname":"test.com","service_id":"12345"}`,
			func(c C) (any, error) {
				a, _, err := c.Domain.SetNewAuthInfo(ctx, &DomainSetNewAuthInfoOptions{DomainListOptions: DomainListOptions{ServiceIdent: ServiceIdent{DName: "test.com"}}, AuthInfo: "hjgf7jfj8f8"})
				return a, err
			},
			func(a any) bool { return a.(*DomainSetNewAuthInfoAnswer).ServiceID == 12345 }},
		{"cancel_transfer", m("dname", "test.com"), `{"domains":[]}`,
			func(c C) (any, error) {
				a, _, err := c.Domain.CancelTransfer(ctx, &DomainListOptions{ServiceIdent: ServiceIdent{DName: "test.com"}})
				return a, err
			},
			func(a any) bool { _, ok := (*a.(*map[string]any))["domains"]; return ok }},
		// The doc example sends "min_pr":"2"; the API takes numbers or strings.
		{"get_rereg_data", m("domains", dnames("qqq.ru", "eee.ru"), "min_pr", 2, "sortcol", "vis"), `[{"dname":"test.ru","tld":"ru","yandex_tic":"1000","tic_mirrors":"null","uni_avg_attendance":"100","search_query_list":"слово дело","start_price":"2500.00","bid":"2500.00","avg_viewings":"2","google_pr":"9","price":"2500.00","blitz_price":"2500.00","first_create_date":"2001-01-01","delete_date":"2020-01-01","is_recommended":"1","all_avg_traffic":"3","registrar":"REGRU-RU"},{"dname":"test.su","search_query_list":null}]`,
			func(c C) (any, error) {
				a, _, err := c.Domain.GetReregData(ctx, &DomainGetReregDataOptions{Domains: idents("qqq.ru", "eee.ru"), MinPR: 2, SortCol: "vis"})
				return a, err
			},
			func(a any) bool {
				d := *a.(*[]DomainReregData)
				return len(d) == 2 && d[0].YandexTIC == 1000 && d[0].BlitzPrice == 2500 && bool(d[0].IsRecommended)
			}},
		{"set_rereg_bids", m("contacts", m("country", "RU", "e_mail", "test@test.ru", "person", "Svyatoslav V Ryurik"), "nss", m("ns0", "ns1.reg.ru", "ns1", "ns2.reg.ru"), "domains", []any{m("dname", "vschizh.ru", "price", 225), m("dname", "vschizh.su", "price", 400)}), order,
			func(c C) (any, error) {
				a, _, err := c.Domain.SetReregBids(ctx, &DomainSetReregBidsOptions{Contacts: contacts, NSS: nss, Domains: []DomainReregBid{{DName: "vschizh.ru", Price: 225}, {DName: "vschizh.su", Price: 400}}})
				return a, err
			},
			func(a any) bool { return a.(*DomainOrderAnswer).BillID == 1234 }},
		{"get_user_rereg_bids", nil, `{"domains":[{"delete_date":"2099-12-31","yandex_tic":null,"max_bid":225.00,"next_price":590.00,"blitz_price":8000.00,"user_bid":225.00,"tld":"ru","dname":"test.ru","registrar":"REGRU","is_active":1,"first_create_date":"2000-01-01","google_pr":null,"is_recommended":1}],"rereg_bids":{"рф":[0,225,590],"ru":[0,225,590]}}`,
			func(c C) (any, error) { a, _, err := c.Domain.GetUserReregBids(ctx); return a, err },
			func(a any) bool {
				r := a.(*DomainGetUserReregBidsAnswer)
				return r.Domains[0].NextPrice == 590 && bool(r.Domains[0].IsActive) && r.Domains[0].DName == "test.ru" && r.ReregBids["рф"][1] == 225
			}},
		{"get_docs_upload_uri", m("dname", "test.ru"), `{"docs_upload_sid":"123456","url":"http://www.reg.ru/user/docs/add?userdoc_secretkey=123456"}`,
			func(c C) (any, error) {
				a, _, err := c.Domain.GetDocsUploadURI(ctx, &ServiceIdent{DName: "test.ru"})
				return a, err
			},
			func(a any) bool { return a.(*DomainGetDocsUploadURIAnswer).DocsUploadSID == "123456" }},
		{"update_contacts", m("contacts", m("country", "RU", "e_mail", "test@test.ru", "person", "Svyatoslav V Ryurik"), "domains", dnames("vschizh.ru", "vschizh.su")), twoDomains,
			func(c C) (any, error) {
				a, _, err := c.Domain.UpdateContacts(ctx, &DomainUpdateContactsOptions{DomainListOptions: DomainListOptions{Domains: idents("vschizh.ru", "vschizh.su")}, Contacts: contacts})
				return a, err
			},
			func(a any) bool { return a.(*DomainServicesAnswer).Domains[1].ServiceID == 12346 }},
		{"update_private_person_flag", m("domains", dnames("vschizh.ru"), "private_person_flag", 0), `{"domains":[{"dname":"vschizh.ru","service_id":"12345","result":"success"}],"pp_flag":"is cleared"}`,
			func(c C) (any, error) {
				a, _, err := c.Domain.UpdatePrivatePersonFlag(ctx, &DomainUpdatePrivatePersonFlagOptions{DomainListOptions: DomainListOptions{Domains: idents("vschizh.ru")}})
				return a, err
			},
			func(a any) bool { return a.(*DomainUpdatePrivatePersonFlagAnswer).PPFlag == "is cleared" }},
		{"register_ns", m("dname", "test.com", "ns0", "ns0.test.com", "ns0ip", "1.2.3.4"), nsResp,
			func(c C) (any, error) {
				a, _, err := c.Domain.RegisterNS(ctx, &DomainNSOptions{ServiceIdent: ServiceIdent{DName: "test.com"}, NS0: "ns0.test.com", NS0IP: "1.2.3.4"})
				return a, err
			},
			func(a any) bool { return a.(*DomainNSAnswer).Resp["actiontype"] == "AddCns" }},
		{"delete_ns", m("dname", "test.com", "ns0", "ns0.test.com", "ns0ip", "1.2.3.4"), nsResp,
			func(c C) (any, error) {
				a, _, err := c.Domain.DeleteNS(ctx, &DomainNSOptions{ServiceIdent: ServiceIdent{DName: "test.com"}, NS0: "ns0.test.com", NS0IP: "1.2.3.4"})
				return a, err
			},
			func(a any) bool { return a.(*DomainNSAnswer).Resp["status"] == "success" }},
		{"get_nss", m("domains", dnames("test.ru", "test.su")), `{"domains":[{"dname":"test.ru","nss":[{"ns":"ns1.reg.ru"},{"ns":"ns1.test.ru","ip":["1.2.3.4","2.3.4.5"]}],"service_id":12345}]}`,
			func(c C) (any, error) {
				a, _, err := c.Domain.GetNSS(ctx, &DomainListOptions{Domains: idents("test.ru", "test.su")})
				return a, err
			},
			func(a any) bool {
				d := a.(*DomainGetNSSAnswer).Domains[0]
				return d.ServiceID == 12345 && d.NSS[0].NS == "ns1.reg.ru" && len(d.NSS[1].IP.([]any)) == 2
			}},
		{"update_nss", m("domains", dnames("test.ru", "----.ru"), "nss", m("ns0", "ns1.reg.ru", "ns1", "ns2.reg.ru")), `{"domains":[{"dname":"test.ru","result":"success","service_id":"12345"},{"dname":"----.ru","error_code":"INVALID_DOMAIN_NAME_FORMAT","result":"domain_name is invalid or unsupported zone"}]}`,
			func(c C) (any, error) {
				a, _, err := c.Domain.UpdateNSS(ctx, &DomainUpdateNSSOptions{DomainListOptions: DomainListOptions{Domains: idents("test.ru", "----.ru")}, NSS: nss})
				return a, err
			},
			func(a any) bool {
				return ErrorCode(a.(*DomainServicesAnswer).Domains[1].Err()) == "INVALID_DOMAIN_NAME_FORMAT"
			}},
		{"update_nss_undelegate", m("domain_name", "test.ru", "undelegate", 1), `{"services":[{"dname":"test.ru","result":"success","service_id":"12345"}]}`,
			func(c C) (any, error) {
				a, _, err := c.Domain.UpdateNSS(ctx, &DomainUpdateNSSOptions{DomainListOptions: DomainListOptions{ServiceIdent: ServiceIdent{DomainName: "test.ru"}}, Undelegate: true})
				return a, err
			},
			func(a any) bool { return a.(*DomainServicesAnswer).Services[0].ServiceID == 12345 }},
		{"delegate", m("domains", dnames("test.ru", "test.su")), twoDomains,
			func(c C) (any, error) {
				a, _, err := c.Domain.Delegate(ctx, &DomainListOptions{Domains: idents("test.ru", "test.su")})
				return a, err
			},
			func(a any) bool { return len(a.(*DomainServicesAnswer).Domains) == 2 }},
		{"undelegate", m("domains", dnames("test.ru", "test.su")), twoDomains,
			func(c C) (any, error) {
				a, _, err := c.Domain.Undelegate(ctx, &DomainListOptions{Domains: idents("test.ru", "test.su")})
				return a, err
			},
			func(a any) bool { return a.(*DomainServicesAnswer).Domains[0].DName == "test.ru" }},
		{"transfer_to_another_account", m("dname", "test.su", "new_user_name", "not_test"), `{"domains":[{"dname":"test.ru","result":"request_is_sent","service_id":"12345"},{"dname":"test.su","error_code":"NOT_TRANSFERED_DOMAIN_STATUS","service_id":"12346"}]}`,
			func(c C) (any, error) {
				a, _, err := c.Domain.TransferToAnotherAccount(ctx, &DomainTransferToAnotherAccountOptions{DomainListOptions: DomainListOptions{ServiceIdent: ServiceIdent{DName: "test.su"}}, NewUserName: "not_test"})
				return a, err
			},
			func(a any) bool {
				d := a.(*DomainServicesAnswer).Domains
				return d[0].Result == "request_is_sent" && ErrorCode(d[1].Err()) == "NOT_TRANSFERED_DOMAIN_STATUS"
			}},
		{"look_at_entering_list", nil, `{"messages":[{"domain_names":["test.ru"],"id":"123456"},{"domain_names":["test.su","test.com"],"id":"123457"}]}`,
			func(c C) (any, error) { a, _, err := c.Domain.LookAtEnteringList(ctx); return a, err },
			func(a any) bool {
				ms := a.(*DomainLookAtEnteringListAnswer).Messages
				return ms[1].ID == 123457 && ms[1].DomainNames[1] == "test.com"
			}},
		{"accept_or_refuse_entering_list", m("domains", []any{m("id", 123456, "action_type", "accept"), m("id", 123458, "action_type", "xxx")}), `{"domains":[{"id":123456,"action_type":"accept","result":"accepted"},{"id":123458,"action_type":"xxx","result":"action_type has incorrect format or data","error_params":{"param":"action_type"},"error_code":"PARAMETER_INCORRECT"}]}`,
			func(c C) (any, error) {
				a, _, err := c.Domain.AcceptOrRefuseEnteringList(ctx, &DomainAcceptOrRefuseEnteringListOptions{Domains: []DomainEnteringAction{{ID: 123456, ActionType: "accept"}, {ID: 123458, ActionType: "xxx"}}})
				return a, err
			},
			func(a any) bool {
				d := a.(*DomainAcceptOrRefuseEnteringListAnswer).Domains
				return d[0].Result == "accepted" && d[0].ID == 123456 && ErrorCode(d[1].Err()) == "PARAMETER_INCORRECT"
			}},
		{"request_to_transfer", m("domains", []any{m("domain_name", "reg.ru"), m("dname", "test0.ru")}), `{"domains":[{"dname":"reg.ru","result":"success"},{"dname":"test0.ru","error_code":"FREE_DOMAIN","result":"Free domain"}]}`,
			func(c C) (any, error) {
				a, _, err := c.Domain.RequestToTransfer(ctx, &DomainListOptions{Domains: []ServiceIdent{{DomainName: "reg.ru"}, {DName: "test0.ru"}}})
				return a, err
			},
			func(a any) bool { return ErrorCode(a.(*DomainServicesAnswer).Domains[1].Err()) == "FREE_DOMAIN" }},
		{"get_tld_info", m("tld", "ru"), `{"tld_info":{"email_verification":"0","idn_allowed":"0","max_dname_length":"63","max_nss_count":"13","min_dname_length":"2","min_nss_count":"0","privacy_protection_allowed":"0","regperiods":"1","reject_days_of_renew_after_create":"0","renew_grace_period_after_expdate":"31","renew_grace_period_before_expdate":"60","renewperiods":"1","restoration_grace_period":"0","tld":"ru"}}`,
			func(c C) (any, error) {
				a, _, err := c.Domain.GetTLDInfo(ctx, &DomainGetTLDInfoOptions{TLD: "ru"})
				return a, err
			},
			func(a any) bool {
				i := a.(*DomainGetTLDInfoAnswer).TLDInfo
				return i.MaxDNameLength == 63 && i.RenewGracePeriodBeforeExpDate == 60 && i.TLD == "ru"
			}},
		{"send_email_verification_letter", m("check_only", 1, "dname", "test.com"), `{"status":"ALREADY_VERIFIED"}`,
			func(c C) (any, error) {
				a, _, err := c.Domain.SendEmailVerificationLetter(ctx, &DomainSendEmailVerificationLetterOptions{DomainListOptions: DomainListOptions{ServiceIdent: ServiceIdent{DName: "test.com"}}, CheckOnly: true})
				return a, err
			},
			func(a any) bool { return a.(*DomainStatusAnswer).Status == "ALREADY_VERIFIED" }},
		{"order_trustee", m("domain_name", "test.ru"), `{"status":"OK"}`,
			func(c C) (any, error) {
				a, _, err := c.Domain.OrderTrustee(ctx, &DomainListOptions{ServiceIdent: ServiceIdent{DomainName: "test.ru"}})
				return a, err
			},
			func(a any) bool { return a.(*DomainStatusAnswer).Status == "OK" }},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			c, mux := setup(t)
			path := tt.path
			if path == "update_nss_undelegate" {
				path = "update_nss"
			}
			handle(t, mux, "/domain/"+path, tt.want, tt.answer)
			got, err := tt.call(c)
			if err != nil {
				t.Fatal(err)
			}
			if !tt.check(got) {
				t.Errorf("unexpected answer %+v", got)
			}
		})
	}
}

func TestDomainDownloadCertificate(t *testing.T) {
	c, mux := setup(t)
	mux.HandleFunc("/domain/download_certificate", func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		switch r.PostForm.Get("input_data") {
		case `{"dname":"test.ru"}`:
			fmt.Fprint(w, "%PDF-1.4 voucher")
		case `{"dname":"none.ru"}`:
			fmt.Fprint(w, `{"result":"success","answer":{"status":"SERVICE_NOT_FOUND"}}`)
		default:
			fmt.Fprint(w, `{"result":"error","error_code":"DOMAIN_NOT_FOUND","error_text":"not found"}`)
		}
	})
	ctx := context.Background()
	pdf, _, err := c.Domain.DownloadCertificate(ctx, &ServiceIdent{DName: "test.ru"})
	if err != nil || string(pdf) != "%PDF-1.4 voucher" {
		t.Fatalf("pdf = %q, err = %v", pdf, err)
	}
	if _, _, err := c.Domain.DownloadCertificate(ctx, &ServiceIdent{DName: "none.ru"}); ErrorCode(err) != "SERVICE_NOT_FOUND" {
		t.Errorf("err = %v", err)
	}
	if _, _, err := c.Domain.DownloadCertificate(ctx, &ServiceIdent{DName: "x.ru"}); ErrorCode(err) != "DOMAIN_NOT_FOUND" {
		t.Errorf("err = %v", err)
	}
}
