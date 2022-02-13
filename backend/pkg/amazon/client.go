package amazon

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/time/rate"

	"go.uber.org/multierr"

	"github.com/utekaravinash/gopaapi5"
	"github.com/utekaravinash/gopaapi5/api"
)

const (
	// maxResultCountPerPage = 10
	maxPage          = 10
	maxAPICallPerDay = 8640
)

var TransactionPerDayExhausted = errors.New("TPD is exhausted")

type Client struct {
	apiClient  *gopaapi5.Client
	apiCallNum int

	rateLimiter *rate.Limiter
}

func NewClient(partnerTag, accessKey, secretKey string) (*Client, error) {
	apiClient, err := gopaapi5.NewClient(accessKey, secretKey, partnerTag, api.Japan)
	if err != nil {
		return nil, err
	}
	return &Client{
		apiClient:   apiClient,
		rateLimiter: rate.NewLimiter(1, 1),
	}, nil
}

// GetBrowseNodes gets browse nodes (categories) by id list
func (c *Client) GetBrowseNodes(ctx context.Context, browseNodeIDs []string) ([]api.BrowseNode, error) {
	if c.tpdExhausted() {
		return nil, TransactionPerDayExhausted
	}
	if c.tpdExhausted() {
		return nil, TransactionPerDayExhausted
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("failed to wait :%w", err)
	}
	c.apiCallNum++
	resp, err := c.apiClient.GetBrowseNodes(ctx, &api.GetBrowseNodesParams{
		BrowseNodeIds: browseNodeIDs,
		Resources: []api.Resource{
			api.BrowseNodesAncestor,
			api.BrowseNodesChildren,
		},
		LanguagesOfPreference: []api.Language{api.JapaneseJapan},
	})
	if err != nil {
		return nil, err
	}
	var errs []error
	if len(resp.Errors) > 0 {
		for _, e := range resp.Errors {
			errs = append(errs, fmt.Errorf("type: %v, code: %v, message: %v", e.Type, e.Code, e.Message))
		}
		return nil, multierr.Combine(errs...)
	}

	return resp.BrowseNodesResult.BrowseNodes, err
}

func (c *Client) SearchItems(ctx context.Context, params *api.SearchItemsParams) (*api.SearchResult, error) {
	if c.tpdExhausted() {
		return nil, TransactionPerDayExhausted
	}
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("failed to wait: %w", err)
	}

	c.apiCallNum++
	resp, err := c.apiClient.SearchItems(ctx, params)
	if err != nil {
		return nil, err
	}
	var errs []error
	if len(resp.Errors) > 0 {
		for _, e := range resp.Errors {
			errs = append(errs, fmt.Errorf("type: %v, code: %v, message: %v", e.Type, e.Code, e.Message))
		}
		return nil, multierr.Combine(errs...)
	}

	return &resp.SearchResult, nil
}

func (c *Client) tpdExhausted() bool {
	return c.apiCallNum >= maxAPICallPerDay
}
