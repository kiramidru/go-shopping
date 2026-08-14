package catalog

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

var (
	ErrBookNotFound = errors.New("book not found")
	ErrInvalidQuery = errors.New("invalid search query")
	ErrUpstream     = errors.New("google books api error")
)

type CatalogService interface {
	FetchBooks(ctx context.Context, req VolumeRequest) (*CatalogResponse, error)
	GetBooksByID(ctx context.Context, bookID string) (*VolumeResponse, error)
}

type Service struct {
	client  *http.Client
	baseURL string
	apiKey  string
}

func NewService(client *http.Client, baseURL, apiKey string) *Service {
	return &Service{
		client:  client,
		baseURL: baseURL,
		apiKey:  apiKey,
	}
}

func (s *Service) FetchBooks(ctx context.Context, req VolumeRequest) (*CatalogResponse, error) {
	params := url.Values{}
	if req.Query != nil {
		params.Set("q", *req.Query)
	}

	if s.apiKey != "" {
		params.Set("key", s.apiKey)
	}

	if req.MaxResults == 0 {
		req.MaxResults = 40
	}

	params.Set("maxResults", fmt.Sprintf("%d", req.MaxResults))
	if req.StartIndex > 0 {
		params.Set("startIndex", fmt.Sprintf("%d", req.StartIndex))
	}

	endpoint := s.baseURL + "/volumes?" + params.Encode()
	bookReq, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	resp, err := s.client.Do(bookReq)
	if err != nil {
		return nil, fmt.Errorf("search books: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrUpstream
	}

	var result CatalogResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode search response: %w", err)
	}
	return &result, nil
}

func (s *Service) GetBooksByID(ctx context.Context, bookID string) (*VolumeResponse, error) {
	params := url.Values{}
	if s.apiKey != "" {
		params.Set("key", s.apiKey)
	}
	endpoint := fmt.Sprintf("%s/volumes/%s?%s", s.baseURL, url.PathEscape(bookID), params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get book: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrBookNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, ErrUpstream
	}

	var volume VolumeResponse
	if err := json.NewDecoder(resp.Body).Decode(&volume); err != nil {
		return nil, fmt.Errorf("decode volume: %w", err)
	}
	return &volume, nil
}
