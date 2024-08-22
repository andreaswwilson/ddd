package http

type ClientOptionFunc func(*Client) error

func WithBaseURL(urlStr string) ClientOptionFunc {
	return func(c *Client) error {
		return c.setBaseURL(urlStr)
	}
}
