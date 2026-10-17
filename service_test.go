package regru

import (
	"context"
	"testing"
)

func TestServiceService(t *testing.T) {
	type M = map[string]any
	ctx := context.Background()
	tests := []struct {
		path  string
		want  any
		reply string
		run   func(s *ServiceService) (any, error)
		check func(got any) bool
	}{
		{"nop", M{"domain_name": "test.ru"},
			`{"services":[{"dname":"test.ru","result":"success","service_id":"12345","servtype":"domain"}]}`,
			func(s *ServiceService) (any, error) {
				a, _, err := s.Nop(ctx, &ServiceListOptions{ServiceIdent: ServiceIdent{DomainName: "test.ru"}})
				return a, err
			},
			func(g any) bool {
				a := g.(*ServiceListAnswer)
				return a.Services[0].ServiceID == 12345 && a.Services[0].Err() == nil
			}},
		{"get_prices", nil,
			`{"currency":"RUR","pricegroup":"0","show_renew_data":"0","prices":[{"servortld":"srv_addip","unit":"month","subtype":"","price":"30"}]}`,
			func(s *ServiceService) (any, error) { a, _, err := s.GetPrices(ctx, nil); return a, err },
			func(g any) bool { a := g.(*ServicePricesAnswer); return a.Prices[0].Price == 30 && a.Currency == "RUR" }},
		{"get_servtype_details", M{"servtype": "srv_hosting_ispmgr"},
			`[{"commonname":"Web-forwarding","extparams":{},"is_renewable":"1","periods_new":"1-10","periods_renew":"1-10","price_new":"120.00","price_renew":"120.00","servtype":"srv_webfwd","subtype":"","unit":"year"}]`,
			func(s *ServiceService) (any, error) {
				a, _, err := s.GetServtypeDetails(ctx, &ServiceGetServtypeDetailsOptions{ServType: "srv_hosting_ispmgr"})
				return a, err
			},
			func(g any) bool {
				a := *g.(*ServiceServtypeDetailsAnswer)
				return a[0].PriceNew == 120 && bool(a[0].IsRenewable)
			}},
		{"create", M{"dname": "test.ru", "email": "test@test.ru", "period": 12, "servtype": "srv_websitebuilder", "subtype": "start"},
			`{"bill_id":"37601633","descr":"service srv_websitebuilder is ordered for domain test.ru","pay_notes":"No charge required","payment":"0.00","service_id":"27249497","uplink_service_id":"0"}`,
			func(s *ServiceService) (any, error) {
				a, _, err := s.Create(ctx, ServiceCreateOptions{"dname": "test.ru", "email": "test@test.ru", "period": 12, "servtype": "srv_websitebuilder", "subtype": "start"})
				return a, err
			},
			func(g any) bool {
				a := g.(*ServiceCreateAnswer)
				return a.ServiceID == 27249497 && a.BillID == 37601633
			}},
		{"check_create", M{"domain_name": "qqq.ru", "period": 1, "servtype": "srv_hosting_plesk", "subtype": "Host-1-0420"},
			`{"descr":"everything is ok for creating","param":{"dname":"qqq.ru","servtype":"srv_hosting_plesk","subtype":"Host-1-0420","period":"1"}}`,
			func(s *ServiceService) (any, error) {
				a, _, err := s.CheckCreate(ctx, &ServiceCheckCreateOptions{DomainName: "qqq.ru", Period: 1, ServType: "srv_hosting_plesk", SubType: "Host-1-0420"})
				return a, err
			},
			func(g any) bool { return g.(*ServiceCheckCreateAnswer).Param["dname"] == "qqq.ru" }},
		{"delete", M{"domain_name": "test.ru", "servtype": "srv_hosting_ispmgr"}, `null`,
			func(s *ServiceService) (any, error) {
				_, err := s.Delete(ctx, &ServiceIdent{DomainName: "test.ru", ServType: "srv_hosting_ispmgr"})
				return nil, err
			}, nil},
		{"get_info", M{"services": []M{{"domain_name": "test.ru"}, {"service_id": 111111}}},
			`{"services":[{"creation_date":"2001-01-01","expiration_date":"2101-01-01","result":"success","service_id":"12345","servtype":"domain","state":"A","subtype":"test","user2_login":"test@test.com"}]}`,
			func(s *ServiceService) (any, error) {
				a, _, err := s.GetInfo(ctx, &ServiceGetInfoOptions{ServiceListOptions: ServiceListOptions{Services: []ServiceIdent{{DomainName: "test.ru"}, {ServiceID: 111111}}}})
				return a, err
			},
			func(g any) bool {
				a := g.(*ServiceInfoAnswer).Services[0]
				return a.State == "A" && a.User2Login == "test@test.com"
			}},
		{"get_list", M{"servtype": "domain"},
			`{"services":[{"creation_date":"2009-04-18","dname":"test.ru","expiration_date":"2011-04-18","service_id":"111","servtype":"domain","state":"A","subtype":"test","uplink_service_id":"0"}]}`,
			func(s *ServiceService) (any, error) {
				a, _, err := s.GetList(ctx, &ServiceGetListOptions{ServType: "domain"})
				return a, err
			},
			func(g any) bool {
				a := g.(*ServiceInfoAnswer).Services[0]
				return a.DName == "test.ru" && a.ExpirationDate == "2011-04-18"
			}},
		{"get_folders", M{"domain_name": "test.ru"},
			`{"folders":[{"folder_id":"12345","folder_name":"test_folder"}]}`,
			func(s *ServiceService) (any, error) {
				a, _, err := s.GetFolders(ctx, &ServiceIdent{DomainName: "test.ru"})
				return a, err
			},
			func(g any) bool { return g.(*ServiceFoldersAnswer).Folders[0].FolderID == 12345 }},
		{"get_details", M{"domain_name": "vschizh.ru", "show_contacts_only": 1},
			`{"services":[{"contacts":{"country":"RU","e_mail":"test@test.ru"},"dname":"vschizh.ru","service_id":"12345","servtype":"domain","result":"success"}]}`,
			func(s *ServiceService) (any, error) {
				a, _, err := s.GetDetails(ctx, &ServiceGetDetailsOptions{ServiceListOptions: ServiceListOptions{ServiceIdent: ServiceIdent{DomainName: "vschizh.ru"}}, ShowContactsOnly: true})
				return a, err
			},
			func(g any) bool { return g.(*ServiceDetailsAnswer).Services[0].Contacts["e_mail"] == "test@test.ru" }},
		{"get_dedicated_server_list", nil,
			`{"server_list":[{"server_id":"171","cpu_content":"Intel Pentium 4 2.8Ghz","cpu_count":"1","cpu_core":"2","hdd_content":"IDE","hdd_count":"2","hdd_size":"120Gb","ram_content":"RAM","ram_count":"2","ram_size":"512Mb","month_traf":"1000","price_retail":"3080"}]}`,
			func(s *ServiceService) (any, error) { a, _, err := s.GetDedicatedServerList(ctx, nil); return a, err },
			func(g any) bool {
				a := g.(*ServiceDedicatedServerListAnswer).ServerList[0]
				return a.ServerID == 171 && a.PriceRetail == 3080
			}},
		{"update", M{"dname": "qqq.ru", "fwd_action": "addfwd", "fwdfrom": "/", "fwdto": "http://www.reg.ru", "servtype": "srv_webfwd", "webfwd_type": "redirect"},
			`{"descr":"service srv_webfwd is updated for domain qqq.ru"}`,
			func(s *ServiceService) (any, error) {
				a, _, err := s.Update(ctx, ServiceUpdateOptions{"dname": "qqq.ru", "fwd_action": "addfwd", "fwdfrom": "/", "fwdto": "http://www.reg.ru", "servtype": "srv_webfwd", "webfwd_type": "redirect"})
				return a, err
			},
			func(g any) bool { return g.(*ServiceUpdateAnswer).Descr != "" }},
		{"renew", M{"period": 2, "service_id": 123456},
			`{"period":"2","payment":"100","currency":"RUR","status":"renew_success","dname":"test12345.ru","bill_id":"123456","servtype":"domain"}`,
			func(s *ServiceService) (any, error) {
				a, _, err := s.Renew(ctx, &ServiceRenewOptions{ServiceListOptions: ServiceListOptions{ServiceIdent: ServiceIdent{ServiceID: 123456}}, Period: 2})
				return a, err
			},
			func(g any) bool { a := g.(*ServiceRenewAnswer); return a.Status == "renew_success" && a.Payment == 100 }},
		{"get_bills", M{"dname": "qqq.ru"},
			`{"services":[{"bills":["123456","234567"],"dname":"qqq.ru","service_id":"12345","servtype":"domain"}]}`,
			func(s *ServiceService) (any, error) {
				a, _, err := s.GetBills(ctx, &ServiceListOptions{ServiceIdent: ServiceIdent{DName: "qqq.ru"}})
				return a, err
			},
			func(g any) bool { return g.(*ServiceBillsAnswer).Services[0].Bills[1] == 234567 }},
		{"check_readonly", M{"service_id": 654321},
			`{"services":[{"error_code":"DATA_IS_READ_ONLY","error_detail":"DATA_IS_READ_ONLY","is_readonly":1,"service_id":"654321"},{"error_code":null,"error_detail":null,"is_readonly":0,"service_id":"123456"}]}`,
			func(s *ServiceService) (any, error) {
				a, _, err := s.CheckReadonly(ctx, &ServiceListOptions{ServiceIdent: ServiceIdent{ServiceID: 654321}})
				return a, err
			},
			func(g any) bool {
				a := g.(*ServiceReadonlyAnswer).Services
				return bool(a[0].IsReadonly) && a[0].Err() != nil && !bool(a[1].IsReadonly) && a[1].Err() == nil
			}},
		{"set_autorenew_flag", M{"dname": "qqq.ru", "flag_value": 0}, `null`,
			func(s *ServiceService) (any, error) {
				_, err := s.SetAutorenewFlag(ctx, &ServiceSetAutorenewFlagOptions{ServiceIdent: ServiceIdent{DName: "qqq.ru"}})
				return nil, err
			}, nil},
		{"suspend", M{"dname": "qqq.ru"}, `null`,
			func(s *ServiceService) (any, error) {
				_, err := s.Suspend(ctx, &ServiceIdent{DName: "qqq.ru"})
				return nil, err
			}, nil},
		{"resume", M{"dname": "qqq.ru"}, `null`,
			func(s *ServiceService) (any, error) {
				_, err := s.Resume(ctx, &ServiceIdent{DName: "qqq.ru"})
				return nil, err
			}, nil},
		{"get_depreciated_period", M{"dname": "qqq.ru", "servtype": "srv_hosting_ispmgr"}, `{"depreciated_period":"1"}`,
			func(s *ServiceService) (any, error) {
				a, _, err := s.GetDepreciatedPeriod(ctx, &ServiceIdent{DName: "qqq.ru", ServType: "srv_hosting_ispmgr"})
				return a, err
			},
			func(g any) bool { return g.(*ServiceDepreciatedPeriodAnswer).DepreciatedPeriod == 1 }},
		{"upgrade", M{"dname": "qqq.ru", "period": 2, "servtype": "srv_hosting_ispmgr", "subtype": "Host-2"},
			`{"bill_id":"123","new_service_id":"-1","returned_amount":"100","withdrawed_amount":"100"}`,
			func(s *ServiceService) (any, error) {
				a, _, err := s.Upgrade(ctx, &ServiceUpgradeOptions{ServiceIdent: ServiceIdent{DName: "qqq.ru", ServType: "srv_hosting_ispmgr", SubType: "Host-2"}, Period: 2})
				return a, err
			},
			func(g any) bool {
				a := g.(*ServiceUpgradeAnswer)
				return a.NewServiceID == -1 && a.WithdrawedAmount == 100
			}},
		{"partcontrol_grant", M{"newlogin": "test_user", "service_id": 1}, `{"newlogin":"test_user","service_id":"1"}`,
			func(s *ServiceService) (any, error) {
				a, _, err := s.PartcontrolGrant(ctx, &ServicePartcontrolGrantOptions{ServiceIdent: ServiceIdent{ServiceID: 1}, NewLogin: "test_user"})
				return a, err
			},
			func(g any) bool {
				a := g.(*ServicePartcontrolAnswer)
				return a.NewLogin == "test_user" && a.ServiceID == 1
			}},
		{"partcontrol_revoke", M{"service_id": 1}, `{"service_id":"1"}`,
			func(s *ServiceService) (any, error) {
				a, _, err := s.PartcontrolRevoke(ctx, &ServiceIdent{ServiceID: 1})
				return a, err
			},
			func(g any) bool { return g.(*ServicePartcontrolAnswer).ServiceID == 1 }},
		{"resend_mail", M{"servtype": "srv_ssl_certificate"}, `{"dname":"test.ru","result":"success","service_id":"111"}`,
			func(s *ServiceService) (any, error) {
				a, _, err := s.ResendMail(ctx, &ServiceResendMailOptions{ServiceIdent: ServiceIdent{ServType: "srv_ssl_certificate"}})
				return a, err
			},
			func(g any) bool { a := g.(*ServiceItem); return a.ServiceID == 111 && a.Result == "success" }},
		{"refill", M{"amount": 10, "currency": "UAH", "service_id": 123456},
			`{"bill_id":"1000000","currency":"RUR","pay_notes":"Amount successfully charged","pay_type":"prepay","payment":"36.88","service_id":"123456"}`,
			func(s *ServiceService) (any, error) {
				a, _, err := s.Refill(ctx, &ServiceRefillOptions{ServiceIdent: ServiceIdent{ServiceID: 123456}, Amount: 10, Currency: "UAH"})
				return a, err
			},
			func(g any) bool { a := g.(*ServiceRefillAnswer); return a.Payment == 36.88 && a.BillID == 1000000 }},
		{"get_balance", M{"service_id": 1}, `{"balance":"10","currency":"USD"}`,
			func(s *ServiceService) (any, error) {
				a, _, err := s.GetBalance(ctx, &ServiceIdent{ServiceID: 1})
				return a, err
			},
			func(g any) bool { return g.(*ServiceBalanceAnswer).Balance == 10 }},
		{"seowizard_manage_link", M{"service_id": 1}, `{"link":"http://reg.ru/seowizard.html"}`,
			func(s *ServiceService) (any, error) {
				a, _, err := s.SeowizardManageLink(ctx, &ServiceIdent{ServiceID: 1})
				return a, err
			},
			func(g any) bool { return g.(*ServiceLinkAnswer).Link == "http://reg.ru/seowizard.html" }},
		{"get_websitebuilder_link", M{"service_id": 1}, `{"link":"http://www.reg.ru/ru/?login_hash=123"}`,
			func(s *ServiceService) (any, error) {
				a, _, err := s.GetWebsitebuilderLink(ctx, &ServiceIdent{ServiceID: 1})
				return a, err
			},
			func(g any) bool { return g.(*ServiceLinkAnswer).Link != "" }},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			c, mux := setup(t)
			handle(t, mux, "/service/"+tt.path, tt.want, tt.reply)
			got, err := tt.run(c.Service)
			if err != nil {
				t.Fatal(err)
			}
			if tt.check != nil && !tt.check(got) {
				t.Errorf("unexpected answer %+v", got)
			}
		})
	}
}
