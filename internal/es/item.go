package es

import "time"

type Platform string

const (
	PlatformRakuten Platform = "rakuten"
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
	GenreID       int      `json:"genre_id"`
	TagIDs        []int    `json:"tag_ids"`
	JANCode       string   `json:"jan_code,omitempty"`
	Platform      Platform `json:"platform"`
	IndexedAt     int64    `json:"indexed_at"` // unix millis
}

func (i *Item) IndexedTime() time.Time {
	return time.UnixMilli(i.IndexedAt)
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
	ItemFieldGenreID       = "genre_id"
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
	ItemFieldGenreID,
	ItemFieldTagIDs,
	ItemFieldJANCode,
	ItemFieldPlatform,
	ItemFieldIndexedAt,
}
