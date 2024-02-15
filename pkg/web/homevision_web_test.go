package web_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"homevision/pkg/test"
	"homevision/pkg/web"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestFetchHousesInfo(t *testing.T) {
	type depFields struct {
		httpClient *test.HttpClientMock
	}

	type output struct {
		houses []web.House
		err    error
	}

	tests := []struct {
		name   string
		on     func(df *depFields)
		assert func(*testing.T, *output)
	}{
		{
			name: "fetch houses successfully",
			on: func(df *depFields) {
				validHouseJSON := `{"houses": [{"id": 1, "address": "123 Main St", "homeowner": "John Doe", "price": 100000, "photoURL": "http://example.com/photo.jpg"}], "ok":true}`
				df.httpClient.On("Do", mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(validHouseJSON)),
				}, nil)
			},
			assert: func(t *testing.T, out *output) {
				assert.NoError(t, out.err)
				assert.NotNil(t, out.houses)
			},
		},
		{
			name: "error on HTTP request",
			on: func(df *depFields) {
				df.httpClient.On("Do", mock.Anything).Return(&http.Response{}, fmt.Errorf("error on request"))
			},
			assert: func(t *testing.T, out *output) {
				assert.NoError(t, out.err)
				assert.Nil(t, out.houses)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Having
			clientMock := test.HttpClientMock{}
			df := &depFields{httpClient: &clientMock}
			tt.on(df)

			// When
			houses, err := web.FetchHousesInfo(1, 10, &clientMock)

			// Then
			tt.assert(t, &output{houses: houses, err: err})
			clientMock.AssertExpectations(t)
		})
	}
}

func TestDownloadImageContent(t *testing.T) {
	type depFields struct {
		httpClient *test.HttpClientMock
	}

	type output struct {
		content string
		err     error
	}

	tests := []struct {
		name   string
		on     func(df *depFields)
		assert func(*testing.T, *output)
	}{
		{
			name: "download image content successfully",
			on: func(df *depFields) {
				df.httpClient.On("Do", mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("image content")),
				}, nil)
			},
			assert: func(t *testing.T, out *output) {
				assert.NoError(t, out.err)
				assert.Equal(t, "image content", out.content)
			},
		},
		{
			name: "error on HTTP request",
			on: func(df *depFields) {
				df.httpClient.On("Do", mock.Anything).Return(&http.Response{}, fmt.Errorf("error on request"))
			},
			assert: func(t *testing.T, out *output) {
				assert.Error(t, out.err)
				assert.Empty(t, out.content)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Having
			clientMock := test.HttpClientMock{}
			df := &depFields{httpClient: &clientMock}
			tt.on(df)

			// When
			content, err := web.DownloadImageContent("http://example.com/image.jpg", &clientMock)

			// Then
			tt.assert(t, &output{content: content, err: err})
			clientMock.AssertExpectations(t)
		})
	}
}

// TODO: use dependency to not to download every image when testing the image processor
//func TestProcessHouseImages(t *testing.T) {
//	type depFields struct {
//		httpClient *test.HttpClientMock
//	}
//
//	type output struct {
//		processedImages []web.Image
//	}
//
//	tests := []struct {
//		name   string
//		on     func(df *depFields)
//		assert func(*testing.T, *output)
//	}{
//		{
//			name: "process house images successfully",
//			on: func(df *depFields) {
//				df.httpClient.On("Do", mock.Anything).Return(&http.Response{
//					StatusCode: http.StatusOK,
//					Body:       io.NopCloser(strings.NewReader(`{"houses": [{"id": 1, "address": "123 Main St", "homeowner": "John Doe", "price": 100000, "photoURL": "http://example.com/photo1.jpg"},{"id": 2, "address": "456 Oak St", "homeowner": "Jane Doe", "price": 150000, "photoURL": "http://example.com/photo2.jpg"}],"ok": true}`)),
//				}, nil)
//			},
//			assert: func(t *testing.T, out *output) {
//				assert.Len(t, out.processedImages, 2)
//
//				assert.Equal(t, 1, out.processedImages[0].ID)
//				assert.Equal(t, "123_Main_St", out.processedImages[0].Address)
//				assert.Equal(t, "photo1.jpg", out.processedImages[0].Name)
//
//				assert.Equal(t, 2, out.processedImages[1].ID)
//				assert.Equal(t, "456_Oak_St", out.processedImages[1].Address)
//				assert.Equal(t, "photo2.jpg", out.processedImages[1].Name)
//			},
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			// Having
//			clientMock := test.HttpClientMock{}
//			df := &depFields{httpClient: &clientMock}
//			tt.on(df)
//
//			// When
//			houses := []web.House{
//				{
//					ID:        1,
//					Address:   "123 Main St",
//					Homeowner: "John Doe",
//					Price:     100000,
//					PhotoURL:  "http://example.com/photo1.jpg",
//				},
//				{
//					ID:        2,
//					Address:   "456 Oak St",
//					Homeowner: "Jane Doe",
//					Price:     150000,
//					PhotoURL:  "http://example.com/photo2.jpg",
//				},
//			}
//			web.ProcessHouseImages(houses, &clientMock)
//
//			// Then
//			tt.assert(t, &output{})
//			clientMock.AssertExpectations(t)
//		})
//	}
//}
