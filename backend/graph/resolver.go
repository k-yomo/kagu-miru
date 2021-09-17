package graph

import "github.com/k-yomo/kagu-miru/backend/search"

//go:generate go run github.com/99designs/gqlgen

type Resolver struct {
	SearchClient search.Client
}

func NewResolver(searchClient search.Client) *Resolver {
	return &Resolver{
		SearchClient: searchClient,
	}
}
