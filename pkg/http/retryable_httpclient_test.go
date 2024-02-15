package http_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	_ "homevision/pkg/http"
	http2 "homevision/pkg/http"
	"homevision/pkg/test"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRetryableHTTPClient_Do(t *testing.T) {
	type depFields struct {
		client      *test.HttpClientMock
		backoffOpts *test.BackOffOptionsMock
	}

	type output struct {
		response *http.Response
		err      error
	}
	tests := []struct {
		name   string
		in     *http.Request
		on     func(df *depFields)
		assert func(*testing.T, *output)
	}{
		{
			name: "do request successfully",
			in:   httptest.NewRequest(http.MethodGet, "/some-url", nil),
			on: func(df *depFields) {
				df.client.On("Do", mock.Anything).Return(&http.Response{}, nil)
			},
			assert: func(t *testing.T, out *output) {
				assert.NoError(t, out.err)
				assert.NotNil(t, out.response)
			},
		},
		// TODO: Should use something to mock backoff.Retry, so I could perform this test. Not it works but it waits
		// until elapsed time finishes

		//{
		//	name: "do request fails - service unavailable",
		//	in:   httptest.NewRequest(http.MethodGet, "/some-url", nil),
		//	on: func(df *depFields) {
		//		df.client.On("Do", mock.Anything).Return(&http.Response{}, errors.New("service unavailable"))
		//	},
		//	assert: func(t *testing.T, out *output) {
		//		assert.Error(t, out.err)
		//		assert.Nil(t, out.response)
		//	},
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Having
			clientMock := test.HttpClientMock{}
			backoffOptsMock := test.BackOffOptionsMock{}
			df := &depFields{client: &clientMock, backoffOpts: &backoffOptsMock}
			tt.on(df)

			// When
			retryableClient := http2.NewRetryableHTTPClient(&clientMock, &backoffOptsMock)
			response, err := retryableClient.Do(tt.in)

			// Then
			tt.assert(t, &output{response: response, err: err})
			clientMock.AssertExpectations(t)
			backoffOptsMock.AssertExpectations(t)

		})
	}

}
