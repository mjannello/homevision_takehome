package main

import (
	"github.com/cenkalti/backoff/v4"
	"homevision/web"
	"net/http"
	"time"
)

func main() {

	backoffStrategy := backoff.NewExponentialBackOff()
	backoffStrategy.MaxElapsedTime = 30 * time.Second
	backoffStrategy.MaxInterval = 5 * time.Second

	client := http.Client{}
	retryableHTTPClient := web.NewRetryableHTTPClient(&client, backoffStrategy)

	totalPages := 1
	perPage := 3
	houses, _ := web.FetchHousesInfo(totalPages, perPage, retryableHTTPClient)

	web.ProcessHouseImages(houses, retryableHTTPClient)

}
