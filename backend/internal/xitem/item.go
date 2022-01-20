package xitem

import (
	"fmt"
)

type Platform string

const (
	PlatformRakuten       Platform = "rakuten"
	PlatformYahooShopping Platform = "yahoo_shopping"
	PlatformPayPayMall    Platform = "paypay_mall"
)

type Status int

const (
	StatusActive Status = iota + 1
	StatusInactive
)

type Item struct {
	// must set an ID generated from ItemUniqueID
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Status        Status    `json:"status"`
	URL           string    `json:"url"`
	AffiliateURL  string    `json:"affiliate_url"`
	Price         int       `json:"price"`
	ImageURLs     []string  `json:"image_urls"`
	AverageRating float64   `json:"average_rating"`
	ReviewCount   int       `json:"review_count"`
	CategoryID    string    `json:"category_id"`
	CategoryIDs   []string  `json:"category_ids"`
	CategoryNames []string  `json:"category_names"`
	BrandName     string    `json:"brand_name"`
	Colors        []string  `json:"colors"`
	WidthRange    *IntRange `json:"widthRange,omitempty"`
	DepthRange    *IntRange `json:"depthRange,omitempty"`
	HeightRange   *IntRange `json:"heightRange,omitempty"`
	JANCode       string    `json:"jan_code,omitempty"`
	Platform      Platform  `json:"platform"`
}

func (i *Item) IsIndexable() bool {
	var excludeCategoryIDs = map[string]struct{}{
		"406451": {}, // ベッド用部品・メンテナンス用品
		"566177": {}, // 収納家具用部品
		"215720": {}, // 蛍光灯
		"566178": {}, // 電球
		"568590": {}, // 誘導灯
		"215716": {}, // 照明器具部品
		"101860": {}, // その他
		"566188": {}, // デスク用部品
		"566193": {}, // カーテン・ブラインド用アクセサリー
		"207738": {}, // 温度計・湿度計
		"500349": {}, // 火鉢
		"101859": {}, // その他
	}

	for _, categoryID := range i.CategoryIDs {
		if _, ok := excludeCategoryIDs[categoryID]; ok {
			return false
		}
	}

	return true
}

func ItemUniqueID(platform Platform, itemID string) string {
	return fmt.Sprintf("%s:%s", platform, itemID)
}
