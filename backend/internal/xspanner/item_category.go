package xspanner

import (
	"time"

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
