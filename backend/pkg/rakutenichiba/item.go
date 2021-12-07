package rakutenichiba

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/k-yomo/kagu-miru/backend/pkg/httputil"
	"github.com/k-yomo/kagu-miru/backend/pkg/urlutil"
)

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

const SearchItemCountPerPage = 30
const SearchItemPageLimit = 100

type SearchItemParams struct {
	// Either one of `Keyword`, `ShopCode`, `ItemCode` or `GenreID` must be supplied
	Keyword  string
	ShopCode string
	ItemCode string
	GenreID  int
	MinPrice int
	MaxPrice int
	Page     int
	SortType SearchItemSortType
}

type SearchItemResponse struct {
	Items []struct {
		Item *Item `json:"Item"`
	} `json:"Items"`
	Hits      int `json:"hits"`
	First     int `json:"first"`
	Last      int `json:"last"`
	Count     int `json:"count"`
	Page      int `json:"page"`
	PageCount int `json:"pageCount"`
}

type SearchItemSortType string

const (
	SearchItemSortTypeStandard          SearchItemSortType = "standard"
	SearchItemSortTypeAffiliateRateAsc  SearchItemSortType = "+affiliateRate"
	SearchItemSortTypeAffiliateRateDesc SearchItemSortType = "-affiliateRate"
	SearchItemSortTypeReviewCountAsc    SearchItemSortType = "+reviewCount"
	SearchItemSortTypeReviewCountDesc   SearchItemSortType = "-reviewCount"
	SearchItemSortTypeReviewAverageAsc  SearchItemSortType = "+reviewAverage"
	SearchItemSortTypeReviewAverageDesc SearchItemSortType = "-reviewAverage"
	SearchItemSortTypeItemPriceAsc      SearchItemSortType = "+itemPrice"
	SearchItemSortTypeItemPriceDesc     SearchItemSortType = "-itemPrice"
	SearchItemSortTypeUpdateTimeAsc     SearchItemSortType = "+updateTime"
	SearchItemSortTypeUpdateTimeDesc    SearchItemSortType = "-updateTime"
)

// SearchItem searches items
// https://webservice.rakuten.co.jp/api/ichibaitemsearch/
func (c *Client) SearchItem(ctx context.Context, params *SearchItemParams) (*SearchItemResponse, error) {
	if params.Keyword == "" && params.ShopCode == "" && params.ItemCode == "" && params.GenreID == 0 {
		return nil, errors.New("either one of `Keyword`, `ShopCode`, `ItemCode` or `GenreID` must be supplied")
	}

	if params.Page == 0 {
		params.Page = 1
	}

	reqParams := map[string]string{
		"sort":        string(params.SortType),
		"affiliateId": c.affiliateID,
	}
	if params.Keyword != "" {
		reqParams["keyword"] = params.Keyword
	}
	if params.ShopCode != "" {
		reqParams["shopCode"] = params.ShopCode
	}
	if params.ItemCode != "" {
		reqParams["itemCode"] = params.ItemCode
	}
	if params.GenreID != 0 {
		reqParams["genreId"] = strconv.Itoa(params.GenreID)
	}
	if params.MinPrice > 0 {
		reqParams["minPrice"] = strconv.Itoa(params.MinPrice)
	}
	if params.MaxPrice > 0 {
		reqParams["maxPrice"] = strconv.Itoa(params.MaxPrice)
	}
	if params.Page != 0 {
		reqParams["page"] = strconv.Itoa(params.Page)
	}

	u := urlutil.CopyWithQueries(c.itemSearchAPIURL, c.buildParams(reqParams))
	var resp SearchItemResponse
	if err := httputil.GetAndUnmarshal(ctx, c.httpClient, u, &resp); err != nil {
		return nil, fmt.Errorf("httputil.GetAndUnmarshal: %w", err)
	}
	return &resp, nil
}
