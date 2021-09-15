package index

type Platform string

const (
	PlatformRakuten = "rakuten"
)

type Status int

const (
	StatusActive Status = iota + 1
	StatusInactive
)

type Item struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Status         Status   `json:"status"`
	SellingPageURL string   `json:"sellingPageUrl"`
	Price          int      `json:"price"`
	ImageURLs      []string `json:"imageUrls"`
	AverageRating  float64  `json:"averageRating"`
	ReviewCount    int      `json:"reviewCount"`
	Platform       Platform `json:"platform"`
}
