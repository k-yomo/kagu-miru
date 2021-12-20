package xspanner

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/spanner"
)

const YahooShoppingItemCategoriesTableName = "yahoo_shopping_item_categories"

// YahooShoppingItemCategory represents item category used in Yahoo Shopping
type YahooShoppingItemCategory struct {
	ID             int64             `spanner:"id"`
	Name           string            `spanner:"name"`
	Level          int64             `spanner:"level"`
	ParentID       spanner.NullInt64 `spanner:"parent_id"`
	ItemCategoryID string            `spanner:"item_category_id"`
	UpdatedAt      time.Time         `spanner:"updated_at"`
}

func GetAllYahooShoppingItemCategories(ctx context.Context, spannerClient *spanner.Client) ([]*YahooShoppingItemCategory, error) {
	ctx, span := otel.Tracer("").Start(ctx, "xspanner.GetAllYahooShoppingItemCategories")
	defer span.End()

	stmt := spanner.NewStatement(`SELECT * FROM yahoo_shopping_item_categories`)
	iter := spannerClient.Single().Query(ctx, stmt)
	defer iter.Stop()

	var itemCategories []*YahooShoppingItemCategory
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var genre YahooShoppingItemCategory
		if err := row.ToStruct(&genre); err != nil {
			return nil, err
		}
		itemCategories = append(itemCategories, &genre)
	}
	return itemCategories, nil
}
