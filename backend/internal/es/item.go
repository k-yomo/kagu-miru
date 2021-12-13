package es

import "github.com/k-yomo/kagu-miru/backend/internal/xitem"

type Item struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	Status        xitem.Status   `json:"status"`
	URL           string         `json:"url"`
	AffiliateURL  string         `json:"affiliate_url"`
	Price         int            `json:"price"`
	ImageURLs     []string       `json:"image_urls"`
	AverageRating float64        `json:"average_rating"`
	ReviewCount   int            `json:"review_count"`
	CategoryID    string         `json:"category_id"`
	CategoryIDs   []string       `json:"category_ids"`
	CategoryNames []string       `json:"category_names"`
	BrandName     string         `json:"brand_name,omitempty"`
	Colors        []string       `json:"colors"`
	TagIDs        []int          `json:"tag_ids"`
	JANCode       string         `json:"jan_code,omitempty"`
	Platform      xitem.Platform `json:"platform"`
	IndexedAt     int64          `json:"indexed_at"` // unix millis
}

func (i *Item) IsActive() bool {
	return i.Status == xitem.StatusActive
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
	ItemFieldBrandName     = "brand_name"
	ItemFieldColors        = "colors"
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
	ItemFieldBrandName,
	ItemFieldColors,
	ItemFieldTagIDs,
	ItemFieldJANCode,
	ItemFieldPlatform,
	ItemFieldIndexedAt,
}
