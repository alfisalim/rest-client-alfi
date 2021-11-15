package custom_resty

import (
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"net/http"
	"time"
)

type BuilderResty struct {
	endpoint string
	headers  map[string]string
	body     interface{}
	client   *resty.Client
	//handler  rest.RestClientRepository
}

func New() *BuilderResty {
	client := resty.New()

	return &BuilderResty{
		client: client,
	}
}

// Set Endpoint
func (b *BuilderResty) SetEndpoint(endpoint string) {
	b.endpoint = endpoint
}

// Set Header
func (b *BuilderResty) SetHeader(header map[string]string) {
	for key, val := range header {
		b.headers[key] = val
	}
}

// Set Request
func (b *BuilderResty) SetBody(body interface{}) {
	b.body = body
}

// Set Request, Header, and Body at a time
func (b *BuilderResty) SetRequest(endpoint string, header map[string]string, body interface{}) {
	b.SetEndpoint(endpoint)
	b.SetHeader(header)
	b.SetBody(body)
}

// Set Timeout
func (b *BuilderResty) TimeoutSet(number int, unitTime time.Duration) {
	b.client.SetTimeout(time.Duration(number) * time.Millisecond)
}

// Post request that was built to client
func (b *BuilderResty) Post(response interface{}) (interface{}, error) {
	var data *resty.Response
	var err error

	data, err = b.client.SetPreRequestHook(b.BeforeRequest).R().SetBody(b.body).Post(b.endpoint)
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			return nil, status.Errorf(codes.DeadlineExceeded, TIMEOUT_SERVICES)
		}
		return nil, err
	}

	if data.StatusCode() != 200 {
		return nil, errors.New("response code not 200")
	}

	var body = data.Body()
	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}

	return response, nil
}

func (b *BuilderResty) BeforeRequest(r *resty.Client, h *http.Request) error {
	for k, v := range b.headers {
		h.Header[k] = []string{v}
	}
	return nil
}
