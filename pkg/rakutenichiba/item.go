package rakutenichiba

import "strings"

type Item struct {
	ItemName        string  `json:"itemName"`
	Catchcopy       string  `json:"catchcopy"`
	ItemCaption     string  `json:"itemCaption"`
	ItemPrice       int     `json:"itemPrice"`
	PointRate       float64 `json:"pointRate"`
	ItemCode        string  `json:"itemCode"`
	ItemURL         string  `json:"itemUrl"`
	AffiliateRate   int     `json:"affiliateRate"`
	AffiliateUrl    string  `json:"affiliateUrl"`
	Availability    int     `json:"availability"`
	GenreID         string  `json:"genreId"`
	TagIDs          []int   `json:"tagIds"`
	MediumImageURLs []struct {
		ImageURL string `json:"imageUrl"`
	} `json:"mediumImageUrls"`
	SmallImageURLs []struct {
		ImageURL string `json:"imageUrl"`
	} `json:"SmallImageUrls"`
	ReviewCount   int     `json:"reviewCount"`
	ReviewAverage float64 `json:"reviewAverage"`

	// "startTime": "",
	// "endTime": "",
	// "asurakuClosingTime": "",
	// 	"pointRateStartTime": "",
	// "pointRateEndTime": "",

	ShopName          string `json:"shopName"`
	ShopCode          string `json:"shopCode"`
	ShopUrl           string `json:"shopUrl"`
	ShopAffiliateUrl  string `json:"shopAffiliateUrl"`
	ShopOfTheYearFlag int    `json:"shopOfTheYearFlag"`

	ShipOverseasFlag int `json:"shipOverseasFlag"`
	AsurakuFlag      int `json:"asurakuFlag"`
	ImageFlag        int `json:"imageFlag"`
	TaxFlag          int `json:"taxFlag"`
	PostageFlag      int `json:"postageFlag"`
	GiftFlag         int `json:"giftFlag"`
	CreditCardFlag   int `json:"creditCardFlag"`

	AsurakuArea      string `json:"asurakuArea"`
	ShipOverseasArea string `json:"shipOverseasArea"`
}

func (i *Item) ID() string {
	return strings.Split(i.ItemCode, ":")[1]
}
