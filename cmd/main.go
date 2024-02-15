package main

import (
	"github.com/cenkalti/backoff/v4"
	http2 "homevision/pkg/http"
	"homevision/pkg/web"
	"net/http"
	"time"
)

// TODO: Create a webProcessor interface that has the  FetchHousesInfo and ProcessHouseImages methods to make the main
// testable
func main() {

	backoffStrategy := backoff.NewExponentialBackOff()
	backoffStrategy.MaxElapsedTime = 30 * time.Second
	backoffStrategy.MaxInterval = 5 * time.Second

	client := http.Client{}
	retryableHTTPClient := http2.NewRetryableHTTPClient(&client, backoffStrategy)

	totalPages := 10
	perPage := 10
	houses, _ := web.FetchHousesInfo(totalPages, perPage, retryableHTTPClient)

	web.ProcessHouseImages(houses, retryableHTTPClient)
	// TODO: add a finish message
}
