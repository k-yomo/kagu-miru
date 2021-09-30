package rakutenichiba

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"github.com/k-yomo/kagu-miru/pkg/httputil"

	"github.com/k-yomo/kagu-miru/pkg/urlutil"
)

const (
	itemSearchAPIURL  = "https://app.rakuten.co.jp/services/api/IchibaItem/Search/20170706"
	genreSearchAPIURL = "https://app.rakuten.co.jp/services/api/IchibaGenre/Search/20140222"
)

type Client struct {
	mu sync.Mutex

	applicationIDs []string
	appIDIndex     int

	affiliateID string

	itemSearchAPIURL  *url.URL
	genreSearchAPIURL *url.URL
	httpClient        *http.Client
}

func NewClient(appIDs []string, affiliateID string) *Client {
	itemAPIURL, _ := url.Parse(itemSearchAPIURL)
	genreAPIURL, _ := url.Parse(genreSearchAPIURL)
	return &Client{
		applicationIDs:    appIDs,
		affiliateID:       affiliateID,
		itemSearchAPIURL:  itemAPIURL,
		genreSearchAPIURL: genreAPIURL,
		httpClient:        http.DefaultClient,
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
func (c *Client) SearchGenre(ctx context.Context, genreID string) (*SearchGenreResponse, error) {
	u := urlutil.CopyWithQueries(c.genreSearchAPIURL, c.buildParams(map[string]string{"genreId": genreID}))
	var resp SearchGenreResponse
	if err := httputil.GetAndUnmarshal(ctx, c.httpClient, u, &resp); err != nil {
		return nil, fmt.Errorf("httputil.GetAndUnmarshal: %w", err)
	}
	return &resp, nil
}

// ApplicationIDNum returns number of application ids available
// It's be useful for rate limiting
func (c *Client) ApplicationIDNum() int {
	return len(c.applicationIDs)
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

func (c *Client) buildParams(params map[string]string) map[string]string {
	p := map[string]string{
		"format":        "json",
		"applicationId": c.getApplicationID(),
	}
	for k, v := range params {
		p[k] = v
	}

	return p
}

func (c *Client) getApplicationID() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	idx := c.appIDIndex
	if idx == len(c.applicationIDs)-1 {
		c.appIDIndex = 0
	} else {
		c.appIDIndex++
	}
	return c.applicationIDs[idx]
}
