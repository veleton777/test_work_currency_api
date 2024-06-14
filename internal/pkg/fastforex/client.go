package fastforex

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

var (
	errInvalidResponse     = errors.New("invalid response")
	errResponseStatusNotOK = errors.New("response status not ok")
)

const convertURI = "convert"

type httpClient interface {
	Do(r *http.Request) (*http.Response, error)
}

type Client struct {
	host       string
	apiKey     string
	httpClient httpClient
}

func NewClient(host, apiKey string, httpClient *http.Client) *Client {
	return &Client{
		host:       host,
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

func (c *Client) Convert(ctx context.Context, from, to string, amount float64) (float64, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.host+"/"+convertURI, nil)
	if err != nil {
		return 0, errors.Wrap(err, "create http request")
	}

	values := make(url.Values)
	values.Add("api_key", c.apiKey)
	values.Add("from", from)
	values.Add("to", to)
	values.Add("amount", strconv.FormatFloat(amount, 'f', -1, 64))

	req.URL.RawQuery = values.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, errors.Wrap(err, "request to fast forex")
	}

	if resp.StatusCode != http.StatusOK {
		return 0, errResponseStatusNotOK
	}

	defer resp.Body.Close()

	var result ConvertResp
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, errors.Wrap(err, "unmarshal json body to struct")
	}

	v, ok := result.Result[to]
	if !ok {
		return 0, errInvalidResponse
	}

	return v, nil
}
