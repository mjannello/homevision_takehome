package test

// TODO: This file could be splitted into different files for every mock, inside each pkg

import (
	"github.com/stretchr/testify/mock"
	http2 "homevision/pkg/http"
	"homevision/pkg/web"
	"net/http"
	"time"
)

type HttpClientMock struct {
	mock.Mock
}

func (h *HttpClientMock) Do(req *http.Request) (*http.Response, error) {
	args := h.Called(req)
	if args.Get(1) == nil {
		return args.Get(0).(*http.Response), nil
	}
	return args.Get(0).(*http.Response), args.Get(1).(error)
}

type BackOffOptionsMock struct {
	mock.Mock
}

func (boo *BackOffOptionsMock) NextBackOff() time.Duration {
	return 0
}

func (boo *BackOffOptionsMock) Reset() {}

type RetryableHTTPClientMock struct {
	mock.Mock
}

func (m *RetryableHTTPClientMock) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

type ProcessorMock struct {
	mock.Mock
}

func (m *ProcessorMock) FetchHousesInfo(totalPages, perPage int, httpClient http2.RetryableHTTPClient) ([]web.House, error) {
	args := m.Called(totalPages, perPage, httpClient)
	houses, _ := args.Get(0).([]web.House)
	err, _ := args.Get(1).(error)
	return houses, err
}

func (m *ProcessorMock) ProcessHouseImages(houses []web.House, client http2.RetryableHTTPClient) {
	m.Called(houses, client)
}
