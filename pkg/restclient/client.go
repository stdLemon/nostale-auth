package restclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
)

const jsonContentType = "application/json"

type Client struct {
	httpClient     *http.Client
	defaultHeaders http.Header
}

type Request struct {
	httpClient  *http.Client
	headers     http.Header
	body        any
	result      any
	errorResult any
}

type Response struct {
	raw         *http.Response
	result      any
	errorResult any
}

func New(defaultHeaders http.Header) Client {
	if defaultHeaders == nil {
		defaultHeaders = make(http.Header)
	}

	return Client{
		httpClient:     &http.Client{},
		defaultHeaders: defaultHeaders.Clone(),
	}
}

func (c Client) NewRequest() *Request {
	return &Request{
		httpClient: c.httpClient,
		headers:    c.defaultHeaders.Clone(),
	}
}

func (r *Request) SetHeader(key, value string) *Request {
	if r.headers == nil {
		r.headers = make(http.Header)
	}

	r.headers.Set(key, value)
	return r
}

func (r *Request) SetHeaders(headers http.Header) *Request {
	if headers == nil {
		return r
	}

	if r.headers == nil {
		r.headers = make(http.Header)
	}

	maps.Copy(r.headers, headers)
	return r
}

func (r *Request) SetBody(body any) *Request {
	r.body = body
	return r
}

func (r *Request) SetResult(result any) *Request {
	r.result = result
	return r
}

func (r *Request) SetError(errResult any) *Request {
	r.errorResult = errResult
	return r
}

func (r *Request) Get(endpoint string) (*Response, error) {
	return r.do(http.MethodGet, endpoint)
}

func (r *Request) Post(endpoint string) (*Response, error) {
	return r.do(http.MethodPost, endpoint)
}

func (r *Request) Put(endpoint string) (*Response, error) {
	return r.do(http.MethodPut, endpoint)
}

func (r *Request) Delete(endpoint string) (*Response, error) {
	return r.do(http.MethodDelete, endpoint)
}

func (r *Response) StatusCode() int {
	if r == nil || r.raw == nil {
		return 0
	}

	return r.raw.StatusCode
}

func (r *Response) Status() string {
	if r == nil || r.raw == nil {
		return ""
	}

	return r.raw.Status
}

func (r *Response) isError() bool {
	return r.StatusCode() >= http.StatusMultipleChoices
}

func (r *Response) Result() any {
	if r == nil {
		return nil
	}

	return r.result
}

func (r *Response) Error() any {
	if r == nil {
		return nil
	}

	return r.errorResult
}

func (r *Response) Header() http.Header {
	if r == nil || r.raw == nil {
		return nil
	}

	return r.raw.Header.Clone()
}

func (r *Request) do(method, endpoint string) (*Response, error) {
	requestBody, err := r.buildBodyReader()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, endpoint, requestBody)
	if err != nil {
		return nil, err
	}

	request.Header = r.headers.Clone()
	if request.Header == nil {
		request.Header = make(http.Header)
	}

	request.Header.Set("Accept", jsonContentType)

	if r.body != nil {
		request.Header.Set("Content-Type", jsonContentType)
	}

	httpResp, err := r.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	resp := &Response{raw: httpResp}

	if resp.isError() {
		statusErr := fmt.Errorf("unexpected response status: %s", resp.Status())
		if r.errorResult == nil {
			return resp, statusErr
		}
		if err := decodeJSONBody(httpResp.Body, r.errorResult); err != nil {
			return resp, fmt.Errorf("%w, failed to decode error response body: %w", statusErr, err)
		}
		resp.errorResult = r.errorResult
		return resp, statusErr
	}

	if r.result == nil {
		return resp, nil
	}
	if err := decodeJSONBody(httpResp.Body, r.result); err != nil {
		return resp, fmt.Errorf("%w, failed to decode response body", err)
	}
	resp.result = r.result

	return resp, nil
}

func (r *Request) buildBodyReader() (io.Reader, error) {
	if r.body == nil {
		return nil, nil
	}

	payload, err := json.Marshal(r.body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(payload), nil
}

func decodeJSONBody(body io.Reader, out any) error {
	decoder := json.NewDecoder(body)
	return decoder.Decode(out)
}
