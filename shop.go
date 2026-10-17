package regru

import (
	"context"
	"encoding/json"
)

// ShopService handles functions of the "shop" category (domain shop).
type ShopService service

// ShopNopAnswer is returned by Nop.
type ShopNopAnswer struct {
	LotID FlexInt `json:"lot_id"`
}

// Nop checks that a lot exists for opts.DName and returns its ID (shop/nop).
func (s *ShopService) Nop(ctx context.Context, opts *ServiceIdent) (*ShopNopAnswer, *Response, error) {
	return call[ShopNopAnswer](ctx, s.client, "shop/nop", opts)
}

// ShopDeleteLotOptions are the parameters of shop/delete_lot.
type ShopDeleteLotOptions struct {
	DName []string `json:"dname"`
}

// DeleteLot deletes lots by domain names (shop/delete_lot).
func (s *ShopService) DeleteLot(ctx context.Context, opts *ShopDeleteLotOptions) (*map[string]any, *Response, error) {
	return call[map[string]any](ctx, s.client, "shop/delete_lot", opts)
}

// ShopLot is a domain shop lot.
type ShopLot struct {
	DName             string    `json:"dname"`
	DNamePuny         string    `json:"dname_puny"`
	PriceType         string    `json:"price_type"` // fixed or offer
	StartPrice        FlexFloat `json:"start_price"`
	Site              Bool      `json:"site"`
	TM                Bool      `json:"tm"`
	IsOnline          Bool      `json:"is_online"`
	Description       string    `json:"description"`
	LotHits           FlexInt   `json:"lot_hits"`
	BidsCnt           FlexInt   `json:"bids_cnt"`
	FirstCreateDomain string    `json:"first_create_domain"`
	CreationDate      string    `json:"creation_date"`
	LotDate           string    `json:"lot_date"`
	LotMTime          string    `json:"lot_mtime"`
	YandexTIC         FlexInt   `json:"yandex_tic"`
	GooglePR          FlexInt   `json:"google_pr"`
	// shortcut: get_info returns JSON arrays, get_lot_list returns strings ("[ 1, 2, 3 ]", "a, b");
	// upgrade to a flexible list type once the live format is confirmed.
	Keywords    json.RawMessage `json:"keywords"`
	CategoryIDs json.RawMessage `json:"category_ids"`
}

// GetInfo returns lot information for opts.DName (shop/get_info).
func (s *ShopService) GetInfo(ctx context.Context, opts *ServiceIdent) (*ShopLot, *Response, error) {
	return call[ShopLot](ctx, s.client, "shop/get_info", opts)
}

// ShopGetLotListOptions are the parameters of shop/get_lot_list.
type ShopGetLotListOptions struct {
	Pg               int    `json:"pg,omitempty"`
	ItemsOnPage      int    `json:"itemsonpage,omitempty"` // 25 (default), 50, 100, 200, 500
	SortCol          string `json:"sortcol,omitempty"`
	SortOrder        string `json:"sortorder,omitempty"` // ASC or DESC
	Keywords         string `json:"keywords,omitempty"`
	ShowMyLots       Bool   `json:"show_my_lots,omitempty"`
	YandexTICFrom    int    `json:"yandex_tic_from,omitempty"`
	YandexTICTo      int    `json:"yandex_tic_to,omitempty"`
	GooglePRFrom     int    `json:"google_pr_from,omitempty"`
	GooglePRTo       int    `json:"google_pr_to,omitempty"`
	PriceFrom        string `json:"price_from,omitempty"`
	PriceTo          string `json:"price_to,omitempty"`
	CreationDateFrom string `json:"creation_date_from,omitempty"`
	CreationDateTo   string `json:"creation_date_to,omitempty"`
	LotDateFrom      string `json:"lot_date_from,omitempty"`
	LotDateTo        string `json:"lot_date_to,omitempty"`
	LotMTimeFrom     string `json:"lot_mtime_from,omitempty"`
	LotMTimeTo       string `json:"lot_mtime_to,omitempty"`
	DNameLengthFrom  int    `json:"dname_length_from,omitempty"`
	DNameLengthTo    int    `json:"dname_length_to,omitempty"`
}

// ShopLotList is returned by GetLotList.
type ShopLotList struct {
	LotsCnt FlexInt   `json:"lots_cnt"`
	Lots    []ShopLot `json:"lots"`
}

// GetLotList returns a page of lots (shop/get_lot_list).
func (s *ShopService) GetLotList(ctx context.Context, opts *ShopGetLotListOptions) (*ShopLotList, *Response, error) {
	return call[ShopLotList](ctx, s.client, "shop/get_lot_list", opts)
}

// ShopCategory is a lot category.
type ShopCategory struct {
	ID             FlexInt        `json:"id"`
	CategoryName   string         `json:"category_name"`
	CategoryNameEn string         `json:"category_name_en"`
	Subcategories  []ShopCategory `json:"subcategories"`
}

// ShopCategoryList is returned by GetCategoryList.
type ShopCategoryList struct {
	CategoryList []ShopCategory `json:"category_list"`
}

// GetCategoryList returns lot categories (shop/get_category_list).
func (s *ShopService) GetCategoryList(ctx context.Context, opts *ServiceIdent) (*ShopCategoryList, *Response, error) {
	return call[ShopCategoryList](ctx, s.client, "shop/get_category_list", opts)
}

// ShopGetSuggestedTagsOptions are the parameters of shop/get_suggested_tags.
type ShopGetSuggestedTagsOptions struct {
	ServiceIdent
	Limit int `json:"limit,omitempty"` // default 10, max 50
}

// ShopSuggestedTags is returned by GetSuggestedTags.
type ShopSuggestedTags struct {
	Tags []string `json:"tags"`
}

// GetSuggestedTags returns popular lot tags (shop/get_suggested_tags).
func (s *ShopService) GetSuggestedTags(ctx context.Context, opts *ShopGetSuggestedTagsOptions) (*ShopSuggestedTags, *Response, error) {
	return call[ShopSuggestedTags](ctx, s.client, "shop/get_suggested_tags", opts)
}
