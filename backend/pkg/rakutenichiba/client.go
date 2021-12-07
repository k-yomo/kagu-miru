package rakutenichiba

import (
	"net/http"
	"net/url"
	"sync"
)

const (
	itemSearchAPIURL  = "https://app.rakuten.co.jp/services/api/IchibaItem/Search/20170706"
	genreSearchAPIURL = "https://app.rakuten.co.jp/services/api/IchibaGenre/Search/20140222"
	tagSearchAPIURL   = "https://app.rakuten.co.jp/services/api/IchibaTag/Search/20140222"
)

type Client struct {
	mu sync.Mutex

	applicationIDs []string
	appIDIndex     int

	affiliateID string

	itemSearchAPIURL  *url.URL
	genreSearchAPIURL *url.URL
	tagSearchAPIURL   *url.URL
	httpClient        *http.Client
}

func NewClient(appIDs []string, affiliateID string) *Client {
	itemAPIURL, _ := url.Parse(itemSearchAPIURL)
	genreAPIURL, _ := url.Parse(genreSearchAPIURL)
	tagAPIURL, _ := url.Parse(tagSearchAPIURL)
	return &Client{
		applicationIDs:    appIDs,
		affiliateID:       affiliateID,
		itemSearchAPIURL:  itemAPIURL,
		genreSearchAPIURL: genreAPIURL,
		tagSearchAPIURL:   tagAPIURL,
		httpClient:        http.DefaultClient,
	}
}

// ApplicationIDNum returns number of application ids available
// It's be useful for rate limiting
func (c *Client) ApplicationIDNum() int {
	return len(c.applicationIDs)
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
