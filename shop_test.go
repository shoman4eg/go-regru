package regru

import (
	"context"
	"testing"
)

func TestShopService(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		path, answer string
		want         any
		call         func(*Client) (any, error)
		ok           func(any) bool
	}{
		{"/shop/nop", `{"lot_id":"123"}`, map[string]any{"dname": "qqq.ru"},
			func(c *Client) (any, error) {
				r, _, err := c.Shop.Nop(ctx, &ServiceIdent{DName: "qqq.ru"})
				return r, err
			},
			func(v any) bool { return v.(*ShopNopAnswer).LotID == 123 }},
		{"/shop/delete_lot", `{}`, map[string]any{"dname": []string{"domain1.ru", "domain2.ru", "domain3.ru"}},
			func(c *Client) (any, error) {
				r, _, err := c.Shop.DeleteLot(ctx, &ShopDeleteLotOptions{DName: []string{"domain1.ru", "domain2.ru", "domain3.ru"}})
				return r, err
			}, nil},
		{"/shop/get_info", `{"category_ids":["1","2","3"],"description":"desc","dname":"test-shop-api-1.ru","dname_puny":"test-shop-api-1.ru","keywords":["k1","k2","k3"],"price_type":"fixed","site":"1","start_price":"200.00","tm":"1"}`,
			map[string]any{"dname": "test-shop-api-1.ru"},
			func(c *Client) (any, error) {
				r, _, err := c.Shop.GetInfo(ctx, &ServiceIdent{DName: "test-shop-api-1.ru"})
				return r, err
			},
			func(v any) bool {
				l := v.(*ShopLot)
				return l.StartPrice == 200 && bool(l.Site) && string(l.Keywords) == `["k1","k2","k3"]`
			}},
		{"/shop/get_lot_list", `{"lots":[{"bids_cnt":"0","category_ids":"[ 1, 2, 3 ]","creation_date":"2010-11-12","dname":"test-shop-api-1.ru","dname_puny":"test-shop-api-1.ru","first_create_domain":"2010-11-12","google_pr":"0","is_online":"1","keywords":"test, shop, api","lot_date":"2011-12-21 17:40:36","lot_hits":"null","lot_mtime":"2016-04-08 13:55:02","price_type":"fixed","site":"0","start_price":"3000.00","tm":"0","yandex_tic":"0"}],"lots_cnt":"1"}`,
			map[string]any{"itemsonpage": 100, "show_my_lots": 1},
			func(c *Client) (any, error) {
				r, _, err := c.Shop.GetLotList(ctx, &ShopGetLotListOptions{ItemsOnPage: 100, ShowMyLots: true})
				return r, err
			},
			func(v any) bool {
				l := v.(*ShopLotList)
				return l.LotsCnt == 1 && l.Lots[0].StartPrice == 3000 && bool(l.Lots[0].IsOnline)
			}},
		{"/shop/get_category_list", `{"category_list":[{"category_name":"Бизнес","category_name_en":"Business","subcategories":[{"category_name":"Бухгалтерское дело","category_name_en":"Accounting profession","id":"32"}]}]}`,
			map[string]any{"dname": "qqq.ru"},
			func(c *Client) (any, error) {
				r, _, err := c.Shop.GetCategoryList(ctx, &ServiceIdent{DName: "qqq.ru"})
				return r, err
			},
			func(v any) bool {
				l := v.(*ShopCategoryList).CategoryList
				return l[0].CategoryNameEn == "Business" && l[0].Subcategories[0].ID == 32
			}},
		{"/shop/get_suggested_tags", `{"tags":["бизнес","интернет","авто","бренд","туризм"]}`, map[string]any{"limit": 25},
			func(c *Client) (any, error) {
				r, _, err := c.Shop.GetSuggestedTags(ctx, &ShopGetSuggestedTagsOptions{Limit: 25})
				return r, err
			},
			func(v any) bool { return len(v.(*ShopSuggestedTags).Tags) == 5 }},
	}
	c, mux := setup(t)
	for _, tt := range tests {
		handle(t, mux, tt.path, tt.want, tt.answer)
		if got, err := tt.call(c); err != nil {
			t.Errorf("%s: %v", tt.path, err)
		} else if tt.ok != nil && !tt.ok(got) {
			t.Errorf("%s: got %+v", tt.path, got)
		}
	}
}
