package search

import "github.com/k-yomo/kagu-miru/internal/es"

type Result struct {
	Hits *Hits `json:"hits"`
}

type Hits struct {
	Total    *Total      `json:"total"`
	MaxScore interface{} `json:"max_score"`
	Hits     []*Hit      `json:"hits"`
}

type Total struct {
	Value int `json:"value"`
}

type Hit struct {
	Index  string   `json:"_index"`
	Type   string   `json:"_type"`
	ID     string   `json:"_id"`
	Score  float64  `json:"_score"`
	Source *es.Item `json:"_source"`
}
