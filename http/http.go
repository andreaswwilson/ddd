package http

import (
	log "ddd/logger"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	defaultBaseUrl = "https://httpbin.org"
)

type Client struct {
	client  *http.Client
	baseURL *url.URL
	token   string
}

func NewClient(token string, options ...ClientOptionFunc) (*Client, error) {
	c := &Client{
		token: token,
	}

	c.client = &http.Client{}

	err := c.setBaseURL(defaultBaseUrl)
	if err != nil {
		return nil, err
	}

	for _, fn := range options {
		if fn == nil {
			continue
		}
		if err := fn(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c *Client) setBaseURL(urlStr string) error {
	// Make sure the given URL end with a slash
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}

	baseURL, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	// Update the base URL of the client.
	c.baseURL = baseURL

	return nil
}

func (c *Client) NewRequest(method, path string) (*http.Request, error) {
	log.Info("heia")
	log.Error("heia2")
	u := *c.baseURL
	unescaped, err := url.PathUnescape(path)
	if err != nil {
		return nil, err
	}
	parsedURL, err := url.Parse(unescaped)
	if err != nil {
		return nil, err
	}
	u.Path = c.baseURL.Path + parsedURL.Path
	u.RawQuery = parsedURL.RawQuery

	reqHeaders := make(http.Header)
	reqHeaders.Set("Accept", "application/json")
	if method == http.MethodGet {
		reqHeaders.Set("Content-Type", "application/json")
	}
	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	defer io.Copy(io.Discard, resp.Body)

	err = checkResponse(resp)
	if err != nil {
		return nil, err
	}

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}
	return resp, err
}

type ErrorResponse struct {
	Body     []byte
	Response *http.Response
	Message  string
}

func (e *ErrorResponse) Error() string {
	path, _ := url.QueryUnescape(e.Response.Request.URL.Path)
	url := fmt.Sprintf("%s://%s%s", e.Response.Request.URL.Scheme, e.Response.Request.URL.Host, path)

	if e.Message == "" {
		return fmt.Sprintf("%s %s: %s", e.Response.Request.Method, url, e.Response.Status)
	} else {
		return fmt.Sprintf("%s %s: %s %s", e.Response.Request.Method, url, e.Response.Status, e.Message)
	}
}

func checkResponse(r *http.Response) error {
	switch r.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusNoContent, http.StatusNotModified:
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := io.ReadAll(r.Body)
	if err == nil && strings.TrimSpace(string(data)) != "" {
		errorResponse.Body = data
		errorResponse.Message = string(data)
	}

	return errorResponse
}
