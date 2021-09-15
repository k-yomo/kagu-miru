package index


type indexParams struct {
	Index *index `json:"index"`
}

type index struct {
	Index string `json:"_index"`
	ID    string `json:"_id"`
}
