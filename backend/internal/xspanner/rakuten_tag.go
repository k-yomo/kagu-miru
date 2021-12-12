package xspanner

import (
	"time"
)

const RakutenTagGroupsTableName = "rakuten_tag_groups"
const RakutenTagsTableName = "rakuten_tags"

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
