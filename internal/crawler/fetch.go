package crawler

import (
	"context"
	"io"
	"net/http"
)

func fetchURL(ctx context.Context, client *http.Client, url string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return nil, 0, err
	}

	response, err := client.Do(req)

	if err != nil {
		return nil, 0, err
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, 0, err
	}

	return body, response.StatusCode, nil
}
