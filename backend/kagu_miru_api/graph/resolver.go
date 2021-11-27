package graph

import (
	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/db"
	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/search"
	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/tracking"
)

//go:generate go run github.com/99designs/gqlgen

type Resolver struct {
	DBClient        db.Client
	SearchClient    search.Client
	SearchIDManager *tracking.SearchIDManager
	EventLoader     tracking.EventLoader
}

func NewResolver(dbClient db.Client, searchClient search.Client, searchIDManager *tracking.SearchIDManager, eventLoader tracking.EventLoader) *Resolver {
	return &Resolver{
		DBClient:        dbClient,
		SearchClient:    searchClient,
		SearchIDManager: searchIDManager,
		EventLoader:     eventLoader,
	}
}
