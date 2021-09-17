package rakuten

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/k-yomo/kagu-miru/pkg/urlutil"
)

const (
	ichibaItemAPIBaseURL  = "https://app.rakuten.co.jp/services/api/IchibaItem/Search/20170706"
	ichibaGenreAPIBaseURL = "https://app.rakuten.co.jp/services/api/IchibaGenre/Search/20140222"
)

type IchibaClient struct {
	applicationIDs []string
	appIDIndex     int

	affiliateID string

	ichibaItemAPIBaseURL  *url.URL
	ichibaGenreAPIBaseURL *url.URL
	httpClient            *http.Client
}

func NewIchibaClient(appIDs []string, affiliateID string) *IchibaClient {
	ichibaItemAPIURL, _ := url.Parse(ichibaItemAPIBaseURL)
	ichibaGenreAPIURL, _ := url.Parse(ichibaGenreAPIBaseURL)
	return &IchibaClient{
		applicationIDs:        appIDs,
		affiliateID:           affiliateID,
		ichibaItemAPIBaseURL:  ichibaItemAPIURL,
		ichibaGenreAPIBaseURL: ichibaGenreAPIURL,
		httpClient:            http.DefaultClient,
	}
}

const GenreFurnitureID = "100804"

type Genre struct {
	ID    int    `json:"genreId"`
	Name  string `json:"genreName"`
	Level int    `json:"genreLevel"`
}

type SearchGenreResponse struct {
	Parents  []*Genre `json:"parents"`
	Current  *Genre   `json:"current"`
	Children []struct {
		Child *Genre `json:"child"`
	} `json:"children"`
}

// SearchGenre searches parent, current and children genre of given ID
// https://webservice.rakuten.co.jp/api/ichibagenresearch/
func (i *IchibaClient) SearchGenre(ctx context.Context, genreID string) (*SearchGenreResponse, error) {
	u := urlutil.CopyWithQueries(i.ichibaGenreAPIBaseURL, i.buildParams(map[string]string{"genreId": genreID}))
	var resp SearchGenreResponse
	if err := getAndUnmarshal(ctx, i.httpClient, u, &resp); err != nil {
		return nil, fmt.Errorf("getAndUnmarshal: %w", err)
	}
	return &resp, nil
}

// ApplicationIDNum returns number of application ids available
// It's be useful for rate limiting
func (i *IchibaClient) ApplicationIDNum() int {
	return len(i.applicationIDs)
}

const SearchItemCountPerPage = 30
const SearchItemPageLimit = 100

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
	return fmt.Sprintf("rakuten_%s", i.ItemCode)
}

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
func (i *IchibaClient) SearchItem(ctx context.Context, params *SearchItemParams) (*SearchItemResponse, error) {
	if params.Keyword == "" && params.ShopCode == "" && params.ItemCode == "" && params.GenreID == 0 {
		return nil, errors.New("either one of `Keyword`, `ShopCode`, `ItemCode` or `GenreID` must be supplied")
	}

	if params.Page == 0 {
		params.Page = 1
	}

	reqParams := map[string]string{
		"sort":        string(params.SortType),
		"affiliateId": i.affiliateID,
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

	u := urlutil.CopyWithQueries(i.ichibaItemAPIBaseURL, i.buildParams(reqParams))
	var resp SearchItemResponse
	if err := getAndUnmarshal(ctx, i.httpClient, u, &resp); err != nil {
		return nil, fmt.Errorf("getAndUnmarshal: %w", err)
	}
	return &resp, nil
}

func (i *IchibaClient) buildParams(params map[string]string) map[string]string {
	p := map[string]string{
		"format":        "json",
		"applicationId": i.getApplicationID(),
	}
	for k, v := range params {
		p[k] = v
	}

	return p
}

func (i *IchibaClient) getApplicationID() string {
	idx := i.appIDIndex
	if idx == len(i.applicationIDs)-1 {
		i.appIDIndex = 0
	} else {
		i.appIDIndex++
	}
	return i.applicationIDs[idx]
}
