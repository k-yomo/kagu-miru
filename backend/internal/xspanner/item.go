package xspanner

import (
	"time"

	"cloud.google.com/go/spanner"
)

const ItemsTableName = "items"

type Item struct {
	ID            string   `spanner:"id"`
	Name          string   `spanner:"name"`
	Description   string   `spanner:"description"`
	Status        int      `spanner:"status"`
	URL           string   `spanner:"url"`
	AffiliateURL  string   `spanner:"affiliate_url"`
	Price         int      `spanner:"price"`
	ImageURLs     []string `spanner:"image_urls"`
	AverageRating float64  `spanner:"average_rating"`
	ReviewCount   int      `spanner:"review_count"`
	CategoryID    string   `spanner:"category_id"`
	// TagIDs        []int    `spanner:"tag_ids"`
	JANCode   spanner.NullString `spanner:"jan_code"`
	Platform  string             `spanner:"platform"`
	UpdatedAt time.Time          `spanner:"updated_at"`
}
