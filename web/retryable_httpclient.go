package web

import (
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"net/http"
	"time"
)

type BackoffOptions interface {
	NextBackOff() time.Duration
	Reset()
}

type RetryableHTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type retryableHTTPClient struct {
	client      *http.Client
	backoffOpts BackoffOptions
}

func NewRetryableHTTPClient(client *http.Client, backoffOpts BackoffOptions) RetryableHTTPClient {
	return &retryableHTTPClient{
		client:      client,
		backoffOpts: backoffOpts,
	}
}

func (r *retryableHTTPClient) Do(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	operation := func() error {
		resp, err = r.client.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode == http.StatusServiceUnavailable {
			return fmt.Errorf("service unavailable")
		}
		return nil
	}

	err = backoff.Retry(operation, r.backoffOpts)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	return resp, nil
}
