package search

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/k-yomo/kagu-miru/backend/internal/es"
	"github.com/k-yomo/kagu-miru/backend/internal/xitem"
	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/graph/gqlmodel"
	"github.com/k-yomo/kagu-miru/backend/pkg/logging"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
)

func mapGraphqlPlatformToPlatform(platform gqlmodel.ItemSellingPlatform) (xitem.Platform, error) {
	switch platform {
	case gqlmodel.ItemSellingPlatformRakuten:
		return xitem.PlatformRakuten, nil
	case gqlmodel.ItemSellingPlatformYahooShopping:
		return xitem.PlatformYahooShopping, nil
	case gqlmodel.ItemSellingPlatformPaypayMall:
		return xitem.PlatformPayPayMall, nil
	default:
		return "", fmt.Errorf("unknown platform %s", platform.String())
	}
}

func mapGraphqlItemColorToSearchItemColor(color gqlmodel.ItemColor) string {
	switch color {
	case gqlmodel.ItemColorWhite:
		return "ホワイト"
	case gqlmodel.ItemColorYellow:
		return "イエロー"
	case gqlmodel.ItemColorOrange:
		return "オレンジ"
	case gqlmodel.ItemColorPink:
		return "ピンク"
	case gqlmodel.ItemColorRed:
		return "レッド"
	case gqlmodel.ItemColorBeige:
		return "ベージュ"
	case gqlmodel.ItemColorSilver:
		return "シルバー"
	case gqlmodel.ItemColorGold:
		return "ゴールド"
	case gqlmodel.ItemColorGray:
		return "グレー"
	case gqlmodel.ItemColorPurple:
		return "パープル"
	case gqlmodel.ItemColorBrown:
		return "ブラウン"
	case gqlmodel.ItemColorGreen:
		return "グリーン"
	case gqlmodel.ItemColorBlue:
		return "ブルー"
	case gqlmodel.ItemColorBlack:
		return "ブラック"
	case gqlmodel.ItemColorNavy:
		return "ネイビー"
	case gqlmodel.ItemColorKhaki:
		return "カーキ"
	case gqlmodel.ItemColorWineRed:
		return "ワインレッド"
	case gqlmodel.ItemColorTransparent:
		return "透明"
	default:
		return ""
	}
}

func mapElasticsearchHitsToItems(ctx context.Context, hits []*elastic.SearchHit) []*es.Item {
	items := make([]*es.Item, 0, len(hits))
	for _, hit := range hits {
		var item es.Item
		if err := json.Unmarshal(hit.Source, &item); err != nil {
			logging.Logger(ctx).Error("Failed to unmarshal hit.Source into es.Item", zap.String("source", string(hit.Source)))
			continue
		}

		items = append(items, &item)
	}

	return items
}
