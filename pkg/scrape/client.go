package scrape

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/ptsypyshev/gb-golang-level3-new/pkg/htmlmeta"
)

var client = http.DefaultClient

var ErrStatusCodeInvalid = errors.New("status code invalid")

func Parse(ctx context.Context, url string) (*htmlmeta.Meta, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("http NewRequestWithContext: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http client Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrStatusCodeInvalid
	}

	meta, err := htmlmeta.Parse(ctx, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("htmlmeta Parse: %w", err)
	}

	return meta, nil
}
