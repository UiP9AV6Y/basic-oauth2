package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/UiP9AV6Y/basic-oauth2/pkg/version"
)

var UserAgent string = "Basic-Oauth2"

func init() {
	UserAgent = UserAgent + "/" + version.Version()
}

type ClientOptions struct {
	Socket  string
	Address string
	Port    int
}

func (o *ClientOptions) Client() (*Client, error) {
	var network string
	var address string
	if o.Socket != "" {
		network = "unix"
		address = o.Socket
	} else {
		network = "tcp"
		address = net.JoinHostPort(o.Address, strconv.Itoa(o.Port))
	}

	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}
	dialContext := func(_ context.Context, _, _ string) (net.Conn, error) {
		return dialer.Dial(network, address)
	}
	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	connect := &http.Client{
		Transport: transport,
	}

	client := &Client{
		network: network,
		address: address,
		connect: connect,
	}
	return client, nil
}

type Client struct {
	network string
	address string
	connect *http.Client
}

func (c *Client) Network() string {
	return c.network
}

func (c *Client) String() string {
	return c.address
}

func (c *Client) URI(path string) string {
	return fmt.Sprint("http://", c.address, path)
}

func (c *Client) Visit(path string, consumer Consumer) error {
	uri := c.URI(path)
	req, err := c.request(uri)
	if err != nil {
		return err
	}

	res, err := c.connect.Do(req)
	if err != nil {
		return err
	}

	return c.consume(res, consumer)
}

func (c *Client) request(uri string) (*http.Request, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", UserAgent)

	return req, nil
}

func (c *Client) consume(res *http.Response, consumer Consumer) error {
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	result := consumer.UnmarshalSample()
	if err := decoder.Decode(&result); err != nil {
		return err
	}

	return consumer.Consume(result)
}
