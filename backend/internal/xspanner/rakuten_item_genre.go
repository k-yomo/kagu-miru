package xspanner

import (
	"context"
	"time"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/spanner"
)

const RakutenItemGenresTableName = "rakuten_item_genres"

// RakutenItemGenre represents Genre (equivalent of item category in Kagumiru) used in Rakuten
type RakutenItemGenre struct {
	ID             int64             `spanner:"id"`
	Name           string            `spanner:"name"`
	Level          int64             `spanner:"level"`
	ParentID       spanner.NullInt64 `spanner:"parent_id"`
	ItemCategoryID string            `spanner:"item_category_id"`
	UpdatedAt      time.Time         `spanner:"updated_at"`
}

func GetAllRakutenItemGenres(ctx context.Context, spannerClient *spanner.Client) ([]*RakutenItemGenre, error) {
	stmt := spanner.NewStatement(`SELECT * FROM rakuten_item_genres`)
	iter := spannerClient.Single().Query(ctx, stmt)
	defer iter.Stop()

	var itemGenres []*RakutenItemGenre
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var genre RakutenItemGenre
		if err := row.ToStruct(&genre); err != nil {
			return nil, err
		}
		itemGenres = append(itemGenres, &genre)
	}
	return itemGenres, nil
}
