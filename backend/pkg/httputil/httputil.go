package httputil

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func GetAndUnmarshal(ctx context.Context, httpClient *http.Client, u *url.URL, to interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return fmt.Errorf("http.NewRequestWithContext, url: %s: %w", u.String(), err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("httpClient.Do, url: %s: %w", u.String(), err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll: %w", err)
	}
	if err := json.Unmarshal(bodyBytes, to); err != nil {
		return fmt.Errorf("json.Unmarshal, body: %s: %w", bodyBytes, err)
	}
	return nil
}

func GetAndXMLUnmarshal(ctx context.Context, httpClient *http.Client, u *url.URL, to interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return fmt.Errorf("http.NewRequestWithContext, url: %s: %w", u.String(), err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("httpClient.Do, url: %s: %w", u.String(), err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll: %w", err)
	}
	if err := xml.Unmarshal(bodyBytes, to); err != nil {
		return fmt.Errorf("xml.Unmarshal, body: %s: %w", bodyBytes, err)
	}
	return nil
}
