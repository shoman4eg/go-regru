package regru

import (
	"context"
	"testing"
)

func TestUserService(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		path, answer string
		want         any
		call         func(*Client) (any, error)
		ok           func(any) bool
	}{
		{"/user/nop", `{}`, nil, func(c *Client) (any, error) { _, _, err := c.User.Nop(ctx); return nil, err }, nil},
		{"/user/create", `{"user_id":"777"}`,
			map[string]any{"set_me_as_referrer": 1, "user_country_code": "RU", "user_email": "test@example.com", "user_login": "othertest", "user_mailnotify": 0, "user_password": "xxxxx"},
			func(c *Client) (any, error) {
				off := Bool(false)
				r, _, err := c.User.Create(ctx, &UserCreateOptions{SetMeAsReferrer: true, CountryCode: "RU", Email: "test@example.com", Login: "othertest", MailNotify: &off, Password: "xxxxx"})
				return r, err
			},
			func(v any) bool { return v.(*UserCreateAnswer).UserID == 777 }},
		{"/user/get_statistics", `{"active_domains_cnt":"5","active_domains_get_ctrl_cnt":"3","domain_folders_cnt":"2","renew_domains_cnt":"4","renew_domains_get_ctrl_cnt":"1","undelegated_domains_cnt":"6"}`, map[string]any{},
			func(c *Client) (any, error) {
				r, _, err := c.User.GetStatistics(ctx, &UserGetStatisticsOptions{})
				return r, err
			},
			func(v any) bool {
				s := v.(*UserStatistics)
				return s.ActiveDomainsCnt == 5 && s.UndelegatedDomainsCnt == 6
			}},
		{"/user/get_balance", `{"credit":"10","currency":"RUR","prepay":"1000"}`, map[string]any{"currency": "RUR"},
			func(c *Client) (any, error) {
				r, _, err := c.User.GetBalance(ctx, &UserGetBalanceOptions{Currency: "RUR"})
				return r, err
			},
			func(v any) bool {
				b := v.(*UserBalance)
				return b.Prepay == 1000 && b.Credit == 10 && b.Currency == "RUR"
			}},
		{"/user/set_reseller_url", `{}`, map[string]any{"servtype": "srv_seowizard", "url": "http://test2.com", "url_type": "refill"},
			func(c *Client) (any, error) {
				_, _, err := c.User.SetResellerURL(ctx, &UserResellerURLOptions{ServType: "srv_seowizard", URL: "http://test2.com", URLType: "refill"})
				return nil, err
			}, nil},
		{"/user/get_reseller_url", `{"url":"http://test2.com"}`, map[string]any{"servtype": "srv_seowizard", "url_type": "refill"},
			func(c *Client) (any, error) {
				r, _, err := c.User.GetResellerURL(ctx, &UserResellerURLOptions{ServType: "srv_seowizard", URLType: "refill"})
				return r, err
			},
			func(v any) bool { return v.(*UserResellerURL).URL == "http://test2.com" }},
		{"/user/get_persons", `{"persons":[{"esia_identified":0,"is_entrepreneur":0,"person_contacts":{"country":"RU","first_name":"Иван"},"person_id":1,"type":"pp","user_id":1}]}`, nil,
			func(c *Client) (any, error) { r, _, err := c.User.GetPersons(ctx); return r, err },
			func(v any) bool {
				p := v.(*UserPersons).Persons
				return len(p) == 1 && p[0].PersonID == 1 && p[0].Type == "pp" && !bool(p[0].ESIAIdentified) && len(p[0].PersonContacts) > 0
			}},
	}
	c, mux := setup(t)
	for _, tt := range tests {
		handle(t, mux, tt.path, tt.want, tt.answer)
		got, err := tt.call(c)
		if err != nil {
			t.Errorf("%s: %v", tt.path, err)
		} else if tt.ok != nil && !tt.ok(got) {
			t.Errorf("%s: got %+v", tt.path, got)
		}
	}
}
