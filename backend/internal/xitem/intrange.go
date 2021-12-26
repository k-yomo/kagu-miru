package xitem

// LteUnlimited is magic number used to represent unlimited number
const LteUnlimited = 999999

type IntRange struct {
	Gte int `json:"gte"`
	Lte int `json:"lte"`
}

func NewIntRange(gte int, lte *int) *IntRange {
	if gte > 0 && lte == nil {
		return &IntRange{
			Gte: gte,
			Lte: LteUnlimited,
		}
	} else {
		return &IntRange{
			Gte: gte,
			Lte: *lte,
		}
	}
}
