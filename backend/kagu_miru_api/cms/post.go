package cms

import "time"

type Post struct {
	Slug        string      `json:"slug"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	MainImage   *Image      `json:"mainImage"`
	PublishedAt time.Time   `json:"publishedAt"`
	Categories  []*Category `json:"categories"`
}

type Category struct {
	ID    string   `json:"id"`
	Names []string `json:"names"`
}

type Image struct {
	Type  string `json:"_type"`
	Asset struct {
		Ref  string `json:"_ref"`
		Type string `json:"_type"`
	} `json:"asset"`
}
