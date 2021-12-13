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
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Status        Status   `json:"status"`
	URL           string   `json:"url"`
	AffiliateURL  string   `json:"affiliate_url"`
	Price         int      `json:"price"`
	ImageURLs     []string `json:"image_urls"`
	AverageRating float64  `json:"average_rating"`
	ReviewCount   int      `json:"review_count"`
	CategoryID    string   `json:"category_id"`
	CategoryIDs   []string `json:"category_ids"`
	CategoryNames []string `json:"category_names"`
	BrandName     string   `json:"brand_name"`
	Colors        []string `json:"colors"`
	TagIDs        []int    `json:"tag_ids"`
	JANCode       string   `json:"jan_code,omitempty"`
	Platform      Platform `json:"platform"`
}

func ItemUniqueID(platform Platform, itemID string) string {
	return fmt.Sprintf("%s:%s", platform, itemID)
}
