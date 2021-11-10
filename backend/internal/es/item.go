package es

import (
	"fmt"
	"strings"
	"time"
)

type Platform string

const (
	PlatformRakuten       Platform = "rakuten"
	PlatformYahooShopping Platform = "yahoo_shopping"
)

type Status int

const (
	StatusActive Status = iota + 1
	StatusInactive
)

type Item struct {
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
	TagIDs        []int    `json:"tag_ids"`
	JANCode       string   `json:"jan_code,omitempty"`
	Platform      Platform `json:"platform"`
	IndexedAt     int64    `json:"indexed_at"` // unix millis
}

func (i *Item) IndexedTime() time.Time {
	return time.UnixMilli(i.IndexedAt)
}

// ItemID returns original item's ID on the platform
func (i *Item) ItemID() string {
	return strings.Split(":", i.ID)[1]
}

func ItemUniqueID(platform Platform, itemID string) string {
	return fmt.Sprintf("%s:%s", platform, itemID)
}

const (
	ItemFieldID            = "id"
	ItemFieldName          = "name"
	ItemFieldDescription   = "description"
	ItemFieldStatus        = "status"
	ItemFieldURL           = "url"
	ItemFieldAffiliateURL  = "affiliate_url"
	ItemFieldPrice         = "price"
	ItemFieldImageURLs     = "image_urls"
	ItemFieldAverageRating = "average_rating"
	ItemFieldReviewCount   = "review_count"
	ItemFieldCategoryID    = "category_id"
	ItemFieldCategoryIDs   = "category_ids"
	ItemFieldCategoryNames = "category_names"
	ItemFieldTagIDs        = "tag_ids"
	ItemFieldJANCode       = "jan_code"
	ItemFieldPlatform      = "platform"
	ItemFieldIndexedAt     = "indexed_at"
)

var AllItemFields = []string{
	ItemFieldID,
	ItemFieldName,
	ItemFieldDescription,
	ItemFieldStatus,
	ItemFieldURL,
	ItemFieldAffiliateURL,
	ItemFieldPrice,
	ItemFieldImageURLs,
	ItemFieldAverageRating,
	ItemFieldReviewCount,
	ItemFieldCategoryID,
	ItemFieldCategoryIDs,
	ItemFieldCategoryNames,
	ItemFieldTagIDs,
	ItemFieldJANCode,
	ItemFieldPlatform,
	ItemFieldIndexedAt,
}
