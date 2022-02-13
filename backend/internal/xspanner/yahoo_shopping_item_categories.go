package xspanner

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/k-yomo/kagu-miru/backend/pkg/logging"

	"go.opentelemetry.io/otel"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/spanner"
)

const YahooShoppingItemCategoriesTableName = "yahoo_shopping_item_categories"

var yahooShoppingItemCategoriesTableAllColumnsString = strings.Join(getColumnNames(YahooShoppingItemCategory{}), ", ")

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

	stmt := spanner.NewStatement(fmt.Sprintf(`SELECT %s FROM yahoo_shopping_item_categories`, yahooShoppingItemCategoriesTableAllColumnsString))
	iter := spannerClient.Single().Query(ctx, stmt)
	defer iter.Stop()

	var itemCategories []*YahooShoppingItemCategory
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, logging.Error(ctx, fmt.Errorf("iter.Next :%w", err))
		}
		var category YahooShoppingItemCategory
		if err := row.ToStruct(&category); err != nil {
			return nil, logging.Error(ctx, fmt.Errorf("row.ToStruct :%w", err))
		}
		itemCategories = append(itemCategories, &category)
	}
	return itemCategories, nil
}
