package xspanner

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"

	"cloud.google.com/go/spanner"
	"github.com/k-yomo/kagu-miru/backend/internal/xitem"
)

const ItemsTableName = "items"

type Item struct {
	ID            string             `spanner:"id"`
	Name          string             `spanner:"name"`
	Description   string             `spanner:"description"`
	Status        int64              `spanner:"status"`
	URL           string             `spanner:"url"`
	AffiliateURL  string             `spanner:"affiliate_url"`
	Price         int64              `spanner:"price"`
	ImageURLs     []string           `spanner:"image_urls"`
	AverageRating float64            `spanner:"average_rating"`
	ReviewCount   int64              `spanner:"review_count"`
	CategoryID    string             `spanner:"category_id"`
	BrandName     spanner.NullString `spanner:"brand_name"`
	Colors        []string           `spanner:"colors"`
	WidthRange    []int              `spanner:"width_range"`  // [gte, lte]
	DepthRange    []int              `spanner:"depth_range"`  // [gte, lte]
	HeightRange   []int              `spanner:"height_range"` // [gte, lte]
	// TagIDs        []int    `spanner:"tag_ids"`
	JANCode   spanner.NullString `spanner:"jan_code"`
	Platform  xitem.Platform     `spanner:"platform"`
	UpdatedAt time.Time          `spanner:"updated_at"`
}

func GetItem(ctx context.Context, spannerClient *spanner.Client, itemID string) (*Item, error) {
	ctx, span := otel.Tracer("").Start(ctx, "xspanner.GetItem")
	defer span.End()

	stmt := spanner.Statement{
		SQL:    `SELECT * FROM items WHERE id = @item_id`,
		Params: map[string]interface{}{"item_id": itemID},
	}
	iter := spannerClient.Single().Query(ctx, stmt)
	defer iter.Stop()

	row, err := iter.Next()
	if err != nil {
		return nil, err
	}
	var item Item
	if err := row.ToStruct(&item); err != nil {
		return nil, err
	}

	return &item, nil
}
