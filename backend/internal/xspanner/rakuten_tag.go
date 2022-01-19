package xspanner

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/k-yomo/kagu-miru/backend/pkg/logging"

	"go.opentelemetry.io/otel"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

const (
	RakutenTagGroupsTableName = "rakuten_tag_groups"
	RakutenTagsTableName      = "rakuten_tags"

	TagGroupIDBrand  = 1000161
	TagGroupIDColor  = 1000111
	TagGroupIDWidth  = 1000057
	TagGroupIDDepth  = 1000058
	TagGroupIDHeight = 1000059
)

var rakutenTagsTableAllColumnsString = strings.Join(getColumnNames(RakutenTag{}), ", ")

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
	ctx, span := otel.Tracer("").Start(ctx, "xspanner.GetAllRakutenTags")
	defer span.End()

	stmt := spanner.NewStatement(fmt.Sprintf(`SELECT %s FROM rakuten_tags`, rakutenTagsTableAllColumnsString))
	iter := spannerClient.Single().Query(ctx, stmt)
	defer iter.Stop()

	var tags []*RakutenTag
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, logging.Error(ctx, fmt.Errorf("iter.Next :%w", err))
		}
		var genre RakutenTag
		if err := row.ToStruct(&genre); err != nil {
			return nil, logging.Error(ctx, fmt.Errorf("row.ToStruct :%w", err))
		}
		tags = append(tags, &genre)
	}
	return tags, nil
}
