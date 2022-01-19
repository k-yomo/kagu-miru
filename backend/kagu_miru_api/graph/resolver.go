package graph

import (
	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/cms"
	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/db"
	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/queryclassifier"
	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/search"
	"github.com/k-yomo/kagu-miru/backend/kagu_miru_api/tracking"
)

//go:generate go run github.com/99designs/gqlgen

type Resolver struct {
	DBClient              db.Client
	SearchClient          search.Client
	QueryClassifierClient queryclassifier.QueryClassifier
	CMSClient             cms.Client
	SearchIDManager       *tracking.SearchIDManager
	EventLoader           tracking.EventLoader
}

func NewResolver(
	dbClient db.Client,
	searchClient search.Client,
	queryClassifierClient queryclassifier.QueryClassifier,
	searchIDManager *tracking.SearchIDManager,
	cmsClient cms.Client,
	eventLoader tracking.EventLoader,
) *Resolver {
	return &Resolver{
		DBClient:              dbClient,
		SearchClient:          searchClient,
		QueryClassifierClient: queryClassifierClient,
		SearchIDManager:       searchIDManager,
		CMSClient:             cmsClient,
		EventLoader:           eventLoader,
	}
}
