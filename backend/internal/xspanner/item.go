package xspanner

import (
	"time"

	"github.com/k-yomo/kagu-miru/backend/internal/xitem"

	"cloud.google.com/go/spanner"
)

const ItemsTableName = "items"

type Item struct {
	ID            string   `spanner:"id"`
	Name          string   `spanner:"name"`
	Description   string   `spanner:"description"`
	Status        int64    `spanner:"status"`
	URL           string   `spanner:"url"`
	AffiliateURL  string   `spanner:"affiliate_url"`
	Price         int64    `spanner:"price"`
	ImageURLs     []string `spanner:"image_urls"`
	AverageRating float64  `spanner:"average_rating"`
	ReviewCount   int64    `spanner:"review_count"`
	CategoryID    string   `spanner:"category_id"`
	// TagIDs        []int    `spanner:"tag_ids"`
	JANCode   spanner.NullString `spanner:"jan_code"`
	Platform  xitem.Platform     `spanner:"platform"`
	UpdatedAt time.Time          `spanner:"updated_at"`
}
