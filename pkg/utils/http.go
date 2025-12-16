package utils

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Headers [][2]string

type Params [][2]string

type Request struct {
	Method  string
	Headers Headers
	URL     string
	Params  Params
	Body    io.ReadWriter
}

const (
	MethodGet     = "GET"
	MethodHead    = "HEAD"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodPatch   = "PATCH" // RFC 5789
	MethodDelete  = "DELETE"
	MethodConnect = "CONNECT"
	MethodOptions = "OPTIONS"
	MethodTrace   = "TRACE"
)

func SendRequest(ctx context.Context, req *Request) (json.RawMessage, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport,
	}

	parsedURL, err := url.Parse(req.URL)
	if err != nil {
		return []byte{}, err
	}

	paramsValue := url.Values{}
	for _, params := range req.Params {
		paramsValue.Add(params[0], params[1])
	}

	parsedURL.RawQuery = paramsValue.Encode()

	request, err := http.NewRequestWithContext(ctx, req.Method, parsedURL.String(), req.Body)
	if err != nil {
		return []byte{}, err
	}

	for _, header := range req.Headers {
		request.Header.Set(header[0], header[1])
	}

	response, err := client.Do(request)
	if err != nil {
		return []byte{}, err
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusBadRequest {
		return body, fmt.Errorf("status: %d, %s, body: %s", response.StatusCode, response.Status, body)
	}
	return body, nil
}
