package xspanner

import (
	"context"
	"time"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

const (
	RakutenTagGroupsTableName = "rakuten_tag_groups"
	RakutenTagsTableName      = "rakuten_tags"

	BrandTagGroupID = 1000161
	ColorTagGroupID = 1000111
)

// RakutenTagGroup represents Tag group in Rakuten
type RakutenTagGroup struct {
	ID        int64     `spanner:"id"`
	Name      string    `spanner:"name"`
	UpdatedAt time.Time `spanner:"updated_at"`
}

// RakutenTag represents Tag in Rakuten
type RakutenTag struct {
	ID         int64     `spanner:"id"`
	Name       string    `spanner:"name"`
	TagGroupID int64     `spanner:"tag_group_id"`
	UpdatedAt  time.Time `spanner:"updated_at"`
}

func GetAllRakutenTags(ctx context.Context, spannerClient *spanner.Client) ([]*RakutenTag, error) {
	stmt := spanner.NewStatement(`SELECT * FROM rakuten_tags`)
	iter := spannerClient.Single().Query(ctx, stmt)
	defer iter.Stop()

	var tags []*RakutenTag
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var genre RakutenTag
		if err := row.ToStruct(&genre); err != nil {
			return nil, err
		}
		tags = append(tags, &genre)
	}
	return tags, nil
}
