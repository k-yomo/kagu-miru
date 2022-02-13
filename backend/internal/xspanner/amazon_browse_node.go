package xspanner

import (
	"context"
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/k-yomo/kagu-miru/backend/pkg/logging"
	"go.opentelemetry.io/otel"
	"google.golang.org/api/iterator"
)

const AmazonBrowseNodesTableName = "amazon_browse_nodes"

var amazonBrowseNodesTableAllColumnsString = strings.Join(getColumnNames(AmazonBrowseNode{}), ", ")

// AmazonBrowseNode represents item category used in Amazon
type AmazonBrowseNode struct {
	ID             string             `spanner:"id"`
	Name           string             `spanner:"name"`
	Level          int64              `spanner:"level"`
	ParentID       spanner.NullString `spanner:"parent_id"`
	ItemCategoryID string             `spanner:"item_category_id"`
	UpdatedAt      time.Time          `spanner:"updated_at"`
}

func GetAllAmazonBrowseNodes(ctx context.Context, spannerClient *spanner.Client) ([]*AmazonBrowseNode, error) {
	ctx, span := otel.Tracer("").Start(ctx, "xspanner.GetAllAmazonBrowseNodes")
	defer span.End()

	stmt := spanner.NewStatement(fmt.Sprintf(`SELECT %s FROM amazon_browse_nodes`, amazonBrowseNodesTableAllColumnsString))
	iter := spannerClient.Single().Query(ctx, stmt)
	defer iter.Stop()

	var browseNodes []*AmazonBrowseNode
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, logging.Error(ctx, fmt.Errorf("iter.Next :%w", err))
		}
		var browseNode AmazonBrowseNode
		if err := row.ToStruct(&browseNode); err != nil {
			return nil, logging.Error(ctx, fmt.Errorf("row.ToStruct :%w", err))
		}
		browseNodes = append(browseNodes, &browseNode)
	}
	return browseNodes, nil
}
