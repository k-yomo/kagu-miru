package xspanner

type IntRange struct {
	Gte int `spanner:"gte"`
	Lte int `spanner:"lte"`
}
