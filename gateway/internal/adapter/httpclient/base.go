package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Impervguin/ds-lab2/gateway/internal/usecase"
)

const UserNameHeader = "X-User-Name"

const requestTimeout = 5 * time.Second

type baseClient struct {
	httpClient *http.Client
	baseURL    string
}

func newBaseClient(baseURL string) baseClient {
	return baseClient{
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{Timeout: requestTimeout},
	}
}

type call struct {
	method   string
	path     string
	query    url.Values
	username string
	body     any
	out      any
	statuses map[int]error
}

func (c *baseClient) do(ctx context.Context, call call) error {
	request, err := c.newRequest(ctx, call)
	if err != nil {
		return err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("call %s %s: %w", call.method, call.path, err)
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode >= http.StatusBadRequest {
		if mapped, ok := call.statuses[response.StatusCode]; ok {
			return mapped
		}
		return fmt.Errorf("call %s %s: %w", call.method, call.path, unexpectedStatus(response))
	}

	if call.out == nil {
		return nil
	}
	if err := json.NewDecoder(response.Body).Decode(call.out); err != nil {
		return fmt.Errorf("decode %s %s response: %w", call.method, call.path, err)
	}
	return nil
}

func (c *baseClient) newRequest(ctx context.Context, call call) (*http.Request, error) {
	var body io.Reader
	if call.body != nil {
		encoded, err := json.Marshal(call.body)
		if err != nil {
			return nil, fmt.Errorf("encode %s %s request: %w", call.method, call.path, err)
		}
		body = bytes.NewReader(encoded)
	}

	target := c.baseURL + call.path
	if len(call.query) > 0 {
		target += "?" + call.query.Encode()
	}

	request, err := http.NewRequestWithContext(ctx, call.method, target, body)
	if err != nil {
		return nil, fmt.Errorf("build %s %s request: %w", call.method, call.path, err)
	}
	if call.body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	if call.username != "" {
		request.Header.Set(UserNameHeader, call.username)
	}
	return request, nil
}

func unexpectedStatus(response *http.Response) error {
	var payload struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil || payload.Message == "" {
		return fmt.Errorf("unexpected status %d", response.StatusCode)
	}
	return fmt.Errorf("unexpected status %d: %s", response.StatusCode, payload.Message)
}

type pageResponse[T any] struct {
	Page          int `json:"page"`
	PageSize      int `json:"pageSize"`
	TotalElements int `json:"totalElements"`
	Items         []T `json:"items"`
}

func toPage[T, D any](response pageResponse[T], convert func(T) D) usecase.Page[D] {
	items := make([]D, 0, len(response.Items))
	for _, item := range response.Items {
		items = append(items, convert(item))
	}

	return usecase.Page[D]{
		Page:          response.Page,
		PageSize:      response.PageSize,
		TotalElements: response.TotalElements,
		Items:         items,
	}
}

func paging(page, size int) url.Values {
	return url.Values{
		"page": {fmt.Sprint(page)},
		"size": {fmt.Sprint(size)},
	}
}
