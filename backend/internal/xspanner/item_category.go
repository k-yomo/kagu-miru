package xspanner

import (
	"context"
	"sort"
	"time"

	"go.opentelemetry.io/otel"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/spanner"
)

const ItemCategoriesTableName = "item_categories"

type ItemCategory struct {
	ID        string             `spanner:"id"`
	Name      string             `spanner:"name"`
	Level     int64              `spanner:"level"`
	ParentID  spanner.NullString `spanner:"parent_id"`
	UpdatedAt time.Time          `spanner:"updated_at"`
}

type ItemCategoryWithParent struct {
	*ItemCategory
	Parent *ItemCategoryWithParent
}

// CategoryIDs returns hierarchized category ids from top level category id to current category id
func (i *ItemCategoryWithParent) CategoryIDs() []string {
	categoryIDs := []string{i.ID}
	parentCategory := i.Parent
	for parentCategory != nil {
		categoryIDs = append([]string{parentCategory.ID}, categoryIDs...)
		parentCategory = parentCategory.Parent
	}
	return categoryIDs
}

// CategoryNames returns hierarchized category names from top level category name to current category name
func (i *ItemCategoryWithParent) CategoryNames() []string {
	categoryNames := []string{i.Name}
	parentCategory := i.Parent
	for parentCategory != nil {
		categoryNames = append([]string{parentCategory.Name}, categoryNames...)
		parentCategory = parentCategory.Parent
	}
	return categoryNames
}

func GetAllItemCategoriesWithParent(ctx context.Context, spannerClient *spanner.Client) ([]*ItemCategoryWithParent, error) {
	ctx, span := otel.Tracer("").Start(ctx, "xspanner.GetAllItemCategoriesWithParent")
	defer span.End()

	allItemCategories, err := GetAllItemCategories(ctx, spannerClient)
	if err != nil {
		return nil, err
	}

	sort.Slice(allItemCategories, func(i, j int) bool {
		return allItemCategories[i].Level < allItemCategories[j].Level
	})

	itemCategoryWithParentMap := make(map[string]*ItemCategoryWithParent)
	for _, itemCategory := range allItemCategories {
		itemCategoryWithParent := &ItemCategoryWithParent{
			ItemCategory: itemCategory,
		}
		if itemCategory.ParentID.Valid {
			itemCategoryWithParent.Parent = itemCategoryWithParentMap[itemCategory.ParentID.StringVal]
		}

		itemCategoryWithParentMap[itemCategory.ID] = itemCategoryWithParent
	}

	itemCategoriesWithParent := make([]*ItemCategoryWithParent, 0, len(itemCategoryWithParentMap))
	for _, itemCategoryWithParent := range itemCategoryWithParentMap {
		itemCategoriesWithParent = append(itemCategoriesWithParent, itemCategoryWithParent)
	}

	return itemCategoriesWithParent, nil
}

func GetAllItemCategories(ctx context.Context, spannerClient *spanner.Client) ([]*ItemCategory, error) {
	ctx, span := otel.Tracer("").Start(ctx, "xspanner.GetAllItemCategories")
	defer span.End()

	stmt := spanner.NewStatement(`SELECT * FROM item_categories`)
	iter := spannerClient.Single().Query(ctx, stmt)
	defer iter.Stop()

	var itemCategories []*ItemCategory
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var itemCategory ItemCategory
		if err := row.ToStruct(&itemCategory); err != nil {
			return nil, err
		}
		itemCategories = append(itemCategories, &itemCategory)
	}
	return itemCategories, nil
}
