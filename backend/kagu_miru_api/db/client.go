package db

import (
	"context"

	"github.com/k-yomo/kagu-miru/backend/internal/xspanner"
)

type Client interface {
	GetAllItemCategories(ctx context.Context) ([]*xspanner.ItemCategory, error)
}
