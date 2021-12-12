package rakutenichiba

type TagGroup struct {
	ID   int    `json:"tagGroupId"`
	Name string `json:"tagGroupName"`
	Tags []*struct {
		Tag *Tag `json:"tag"`
	} `json:"tags"`
}

type Tag struct {
	ID   int    `json:"tagId"`
	Name string `json:"tagName"`
	// Almost all tags have 0 for parent id
	// So not sure if it's used in Rakuten
	ParentID int `json:"parentTagId"`
}
