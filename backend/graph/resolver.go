package graph

import (
	"github.com/k-yomo/kagu-miru/backend/search"
	"github.com/k-yomo/kagu-miru/backend/tracking"
)

//go:generate go run github.com/99designs/gqlgen

type Resolver struct {
	SearchClient    search.Client
	SearchIDManager *tracking.SearchIDManager
	EventLoader     tracking.EventLoader
}

func NewResolver(searchClient search.Client, searchIDManager *tracking.SearchIDManager, eventLoader tracking.EventLoader) *Resolver {
	return &Resolver{
		SearchClient:    searchClient,
		SearchIDManager: searchIDManager,
		EventLoader:     eventLoader,
	}
}
