package yahoo_shopping

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/k-yomo/kagu-miru/backend/pkg/httputil"
	"github.com/k-yomo/kagu-miru/backend/pkg/urlutil"
)

const (
	itemSearchAPIURL     = "https://shopping.yahooapis.jp/ShoppingWebService/V3/itemSearch"
	categorySearchAPIURL = "https://shopping.yahooapis.jp/ShoppingWebService/V1/categorySearch"
)

const (
	maxResultCount      = 50
	maxResultTotalCount = 1000
)

const CategoryFurnitureID = 2506

type Client struct {
	mu sync.Mutex

	applicationIDs []string
	appIDIndex     int

	itemSearchAPIURL     *url.URL
	categorySearchAPIURL *url.URL
	httpClient           *http.Client
}

func NewClient(appIDs []string) *Client {
	itemAPIURL, _ := url.Parse(itemSearchAPIURL)
	categoryAPIURL, _ := url.Parse(categorySearchAPIURL)
	return &Client{
		applicationIDs:       appIDs,
		itemSearchAPIURL:     itemAPIURL,
		categorySearchAPIURL: categoryAPIURL,
		httpClient:           http.DefaultClient,
	}
}

func (c *Client) ApplicationIDNum() int {
	return len(c.applicationIDs)
}

type Category struct {
	ID       int         `xml:"Id"`
	ParentID string      `xml:"ParentId"`
	URL      string      `xml:"Url"`
	Title    string      `xml:"Title>Short"`
	IsAdult  int         `xml:"IsAdult"`
	Children []*Category `xml:"Children>Child"`
}

type GetCategoryResponse struct {
	Result struct {
		Categories struct {
			Current  *Category   `xml:"Current"`
			Children []*Category `xml:"Children>Child"`
		} `xml:"Categories"`
	} `xml:"Result"`
}

// https://developer.yahoo.co.jp/webapi/shopping/shopping/v1/categorysearch.html
func (c *Client) SearchCategory(ctx context.Context, categoryID int) (*GetCategoryResponse, error) {
	u := urlutil.CopyWithQueries(c.categorySearchAPIURL, c.buildParams(map[string]string{"category_id": strconv.Itoa(categoryID)}))
	var resp GetCategoryResponse
	if err := httputil.GetAndXMLUnmarshal(ctx, c.httpClient, u, &resp); err != nil {
		return nil, fmt.Errorf("httputil.GetAndUnmarshal: %w", err)
	}
	return &resp, nil
}

// GetCategoryWithAllChildren gets genre with all lower hierarchy genre
func (c *Client) GetCategoryWithAllChildren(ctx context.Context, categoryID int) (*Category, error) {
	res, err := c.SearchCategory(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	category := res.Result.Categories.Current
	if err := c.setChildCategories(ctx, category); err != nil {
		return nil, fmt.Errorf("c.setChildCategories: %w", err)
	}
	return category, nil
}

func (c *Client) setChildCategories(ctx context.Context, category *Category) error {
	time.Sleep((time.Duration(1000.0 / float64(c.ApplicationIDNum()))) * time.Millisecond)

	res, err := c.SearchCategory(ctx, category.ID)
	if err != nil {
		return fmt.Errorf("rakutenIchibaAPIClient.SearchGenre: %w", err)
	}
	categories := make([]*Category, 0, len(category.Children))
	for _, child := range res.Result.Categories.Children {
		if err := c.setChildCategories(ctx, child); err != nil {
			return err
		}
		categories = append(categories, child)
	}
	category.Children = categories
	return nil
}

type SearchItemSortType string

const (
	SearchItemSortTypeScoreDesc       SearchItemSortType = "-score"
	SearchItemSortTypePriceAsc        SearchItemSortType = "+price"
	SearchItemSortTypePriceDesc       SearchItemSortType = "-price"
	SearchItemSortTypeReviewCountDesc SearchItemSortType = "-review_count"
)

type SearchItemParams struct {
	Query      string
	CategoryID int
	PriceFrom  int
	PriceTo    int
	Page       int
	SortType   SearchItemSortType
}

type SearchItemResponse struct {
	TotalResultsAvailable int     `json:"totalResultsAvailable"`
	TotalResultsReturned  int     `json:"totalResultsReturned"`
	FirstResultPosition   int     `json:"firstResultPosition"`
	Hits                  []*Item `json:"hits"`
}

// https://developer.yahoo.co.jp/webapi/shopping/shopping/v3/itemsearch.html
func (c *Client) SearchItem(ctx context.Context, params *SearchItemParams) (*SearchItemResponse, error) {
	if params.Page == 0 {
		params.Page = 1
	}

	start := (params.Page-1)*maxResultCount + 1

	reqParams := map[string]string{
		"sort":    string(params.SortType),
		"results": "50",
		"start":   strconv.Itoa(start),
		"page":    strconv.Itoa(start),
	}
	if params.Query != "" {
		reqParams["query"] = params.Query
	}
	if params.CategoryID != 0 {
		reqParams["genre_category_id"] = strconv.Itoa(params.CategoryID)
	}
	if params.PriceFrom > 0 {
		reqParams["price_from"] = strconv.Itoa(params.PriceFrom)
	}
	if params.PriceTo > 0 {
		reqParams["price_to"] = strconv.Itoa(params.PriceTo)
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
		"appid": c.getApplicationID(),
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
