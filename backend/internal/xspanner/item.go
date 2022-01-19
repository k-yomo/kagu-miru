package xspanner

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/k-yomo/kagu-miru/backend/pkg/logging"

	"cloud.google.com/go/spanner"
	"github.com/k-yomo/kagu-miru/backend/internal/xerror"
	"github.com/k-yomo/kagu-miru/backend/internal/xitem"
	"go.opentelemetry.io/otel"
	"google.golang.org/api/iterator"
)

const ItemsTableName = "items"

var itemsTableAllColumnsString = strings.Join(getColumnNames(Item{}), ", ")

type Item struct {
	ID            string             `spanner:"id"`
	GroupID       spanner.NullString `spanner:"group_id"` // temporally nullable, will fix to NOT NULL after migration
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
	WidthRange    []int64            `spanner:"width_range"`  // [gte, lte]
	DepthRange    []int64            `spanner:"depth_range"`  // [gte, lte]
	HeightRange   []int64            `spanner:"height_range"` // [gte, lte]
	JANCode       spanner.NullString `spanner:"jan_code"`
	Platform      xitem.Platform     `spanner:"platform"`
	UpdatedAt     time.Time          `spanner:"updated_at"`
}

func GetItem(ctx context.Context, spannerClient *spanner.Client, itemID string) (*Item, error) {
	ctx, span := otel.Tracer("").Start(ctx, "xspanner.GetItem")
	defer span.End()

	stmt := spanner.Statement{
		SQL:    fmt.Sprintf(`SELECT %s FROM items WHERE id = @item_id`, itemsTableAllColumnsString),
		Params: map[string]interface{}{"item_id": itemID},
	}
	iter := spannerClient.Single().Query(ctx, stmt)
	defer iter.Stop()

	row, err := iter.Next()
	if err != nil {
		if err == iterator.Done {
			return nil, xerror.NewNotFound(fmt.Errorf("item '%s' is not found", itemID))
		}
		return nil, err
	}
	var item Item
	if err := row.ToStruct(&item); err != nil {
		return nil, logging.Error(ctx, fmt.Errorf("unmarshal spanner result to item struct :%w", err))
	}

	return &item, nil
}

func GetSameGroupItemsByItemID(ctx context.Context, spannerClient *spanner.Client, itemID string) ([]*Item, error) {
	ctx, span := otel.Tracer("").Start(ctx, "xspanner.GetSameGroupItemsByItemID")
	defer span.End()

	stmt := spanner.Statement{
		SQL: fmt.Sprintf(`
SELECT %s
FROM items
WHERE 
	group_id = (
		SELECT group_id 
		FROM items 
		WHERE id = @item_id
	)
`, itemsTableAllColumnsString),
		Params: map[string]interface{}{"item_id": itemID},
	}
	iter := spannerClient.Single().Query(ctx, stmt)
	defer iter.Stop()

	var items []*Item
	for {
		row, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			return nil, logging.Error(ctx, fmt.Errorf("iter.Next :%w", err))
		}
		var item Item
		if err := row.ToStruct(&item); err != nil {
			return nil, logging.Error(ctx, fmt.Errorf("row.ToStruct :%w", err))
		}
		items = append(items, &item)
	}

	return items, nil
}
