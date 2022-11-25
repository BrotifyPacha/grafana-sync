package miniGrafanaClient

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type Client struct {
	baseUrl string
}

func NewClient(host string) (*Client, error) {
	err := validateHost(host)
	if err != nil {
		return nil, err
	}
	host = strings.TrimRight(host, "/") + "/"
	return &Client{baseUrl: host}, nil
}

func (c *Client) Get(apiRoute string) (bytes []byte, err error) {
	apiRoute = strings.TrimLeft(apiRoute, "/")
	resp, err := http.Get(c.baseUrl + apiRoute)
	if err != nil {
		return
	}
	bytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("Bad status code - %d, body:\n%s", resp.StatusCode, string(bytes))) 
		return
	}
	return
}

func validateHost(host string) error {
	protocolRegex := regexp.MustCompile("^https?://.*")
	if !protocolRegex.MatchString(host) {
		return errors.New("Host must contain protocol http or https")
	}
	return nil
}
