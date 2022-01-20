package item_fetcher

import "github.com/k-yomo/kagu-miru/backend/internal/xitem"

func ItemOrderingKey(item *xitem.Item) string {
	orderingKey := item.JANCode
	if orderingKey == "" {
		orderingKey = item.Name
	}
	return orderingKey
}
