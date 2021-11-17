// Package robtex a simple Robtex API client.
// https://www.robtex.com/api/
package robtex

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"
)

// API endpoints.
const (
	FreeAPIBaseURL = "https://freeapi.robtex.com"
	ProAPIBaseURL  = "https://proapi.robtex.com"
)

// Client a Robtex API client.
type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	apiKey     string
}

// New created a new Client.
func New(apiKey string) *Client {
	baseURL, _ := url.Parse(FreeAPIBaseURL)
	if apiKey != "" {
		baseURL, _ = url.Parse(ProAPIBaseURL)
	}

	return &Client{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		baseURL:    baseURL,
		apiKey:     apiKey,
	}
}

// IPQuery This API returns the current forward and reverse of an IP number, together with GEO-location data and network data.
// The format returned is JSON.
// Most keys are self-explanatory.
// ex: https://freeapi.robtex.com/ipquery/199.19.54.1
func (c Client) IPQuery(ctx context.Context, ip string) (*IPQueryResponse, error) {
	req, err := c.newRequest(ctx, "ipquery", ip)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%d: %s", resp.StatusCode, string(data))
	}

	var r IPQueryResponse
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

// ASQuery Returns an array of networks related to a specific AS number.
// Currently, only returns networks actually in global bgp table, but plans are to extend it.
// ex: https://freeapi.robtex.com/asquery/1234
func (c Client) ASQuery(ctx context.Context, number string) (*ASQueryResponse, error) {
	req, err := c.newRequest(ctx, "asquery", number)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%d: %s", resp.StatusCode, string(data))
	}

	var r ASQueryResponse
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

// PassiveDNSForward This API returns ldjson format, that is one JSON object per line.
// The format used is Passive DNS - Common Output Format (https://tools.ietf.org/html/draft-dulaunoy-dnsop-passive-dns-cof-03).
// ex: https://freeapi.robtex.com/pdns/forward/a.iana-servers.net
func (c Client) PassiveDNSForward(ctx context.Context, domain string) ([]PassiveDNS, error) {
	req, err := c.newRequest(ctx, "pdns", "forward", domain)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%d: %s", resp.StatusCode, string(data))
	}

	decoder := json.NewDecoder(resp.Body)

	var results []PassiveDNS
	for decoder.More() {
		var r PassiveDNS
		err = decoder.Decode(&r)
		if err != nil {
			return nil, err
		}
		results = append(results, r)
	}

	return results, nil
}

// PassiveDNSReverse This API returns ldjson format, that is one JSON object per line.
// The format used is Passive DNS - Common Output Format (https://tools.ietf.org/html/draft-dulaunoy-dnsop-passive-dns-cof-03).
// ex: https://freeapi.robtex.com/pdns/reverse/199.43.132.53
func (c Client) PassiveDNSReverse(ctx context.Context, ip string) ([]PassiveDNS, error) {
	req, err := c.newRequest(ctx, "pdns", "reverse", ip)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%d: %s", resp.StatusCode, string(data))
	}

	decoder := json.NewDecoder(resp.Body)

	var results []PassiveDNS
	for decoder.More() {
		var r PassiveDNS
		err = decoder.Decode(&r)
		if err != nil {
			return nil, err
		}
		results = append(results, r)
	}

	return results, nil
}

func (c Client) newRequest(ctx context.Context, k ...string) (*http.Request, error) {
	endpoint, err := c.baseURL.Parse(path.Join(c.baseURL.Path, path.Join(k...)))
	if err != nil {
		return nil, err
	}

	if c.apiKey != "" {
		query := endpoint.Query()
		query.Set("key", c.apiKey)
		endpoint.RawQuery = query.Encode()
	}

	return http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), http.NoBody)
}
