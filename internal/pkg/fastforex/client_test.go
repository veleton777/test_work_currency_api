package fastforex_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/veleton777/test_work_blum/internal/pkg/fastforex"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ClientSuite struct {
	suite.Suite
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}

func (s *ClientSuite) TestConvertMethod() {
	testCases := []struct {
		name    string
		handler func(res http.ResponseWriter, req *http.Request)
		expRes  float64
		expErr  error
	}{
		{
			name: "success",
			handler: func(res http.ResponseWriter, req *http.Request) {
				res.WriteHeader(http.StatusOK)
				res.Write([]byte(`{"result": {"BTC": 0.00001444655}}`))
			},
			expRes: 0.00001444655,
			expErr: nil,
		},
		{
			name: "invalid_response_status",
			handler: func(res http.ResponseWriter, req *http.Request) {
				res.WriteHeader(http.StatusBadGateway)
				res.Write([]byte(`{}`))
			},
			expRes: 0,
			expErr: errors.New("response status not ok"),
		},
		{
			name: "invalid_response_json",
			handler: func(res http.ResponseWriter, req *http.Request) {
				res.WriteHeader(http.StatusOK)
				res.Write([]byte(`another text`))
			},
			expRes: 0,
			expErr: errors.New("unmarshal json body to struct"),
		},
		{
			name: "invalid_response_text",
			handler: func(res http.ResponseWriter, req *http.Request) {
				res.WriteHeader(http.StatusOK)
				res.Write([]byte(`{"result": {"ETH": 123}}`))
			},
			expRes: 0,
			expErr: errors.New("invalid response"),
		},
	}

	for _, c := range testCases {
		s.T().Run(c.name, func(t *testing.T) {
			testSrv := httptest.NewServer(http.HandlerFunc(c.handler))
			defer testSrv.Close()

			ctx := context.Background()

			cl := fastforex.NewClient(testSrv.URL, "apiKey", &http.Client{})
			res, err := cl.Convert(ctx, "USD", "BTC", 1)

			require.Equal(t, res, c.expRes)

			if c.expErr != nil {
				require.ErrorContains(t, err, c.expErr.Error())

				return
			}

			require.NoError(t, err)
		})
	}
}
