package db

import (
	"context"

	"github.com/k-yomo/kagu-miru/backend/internal/xspanner"
)

type Client interface {
	GetItem(ctx context.Context, itemID string) (*xspanner.Item, error)
	GetSameGroupItemsByItemID(ctx context.Context, itemID string) ([]*xspanner.Item, error)
	GetAllActiveItemCategories(ctx context.Context) ([]*xspanner.ItemCategory, error)
	GetAllItemCategoriesWithParent(ctx context.Context) ([]*xspanner.ItemCategoryWithParent, error)
	GetTopLevelItemCategories(ctx context.Context) ([]*xspanner.ItemCategory, error)
}
