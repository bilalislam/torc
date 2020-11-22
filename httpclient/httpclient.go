package httpclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/valyala/fasthttp"
)

type Client struct {
	client *fasthttp.Client
}

type Request struct {
	request *fasthttp.Request
	echo    echo.Context
	client  *fasthttp.Client
}

type AppenderClient struct {
	request *fasthttp.Request
	echo    echo.Context
	client  *fasthttp.Client
}

type BuildContext struct {
	client             *fasthttp.Client
	dependencySetter   *DependencySetter
	dependencyResolver *DependencyResolver
	request            *fasthttp.Request
	echo               echo.Context
}

func NewHttpClient() *Client {
	client := &fasthttp.Client{
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	return &Client{client: client}
}

func (p *Client) NewRequest() *Request {
	var req = fasthttp.AcquireRequest()
	return &Request{
		request: req,
		client:  p.client,
	}
}

func (c *AppenderClient) AppendHeader(key string, value string) *AppenderClient {
	c.request.Header.Add(key, value)
	return c
}

func (c *AppenderClient) AppendBodyAsByte(body []byte) *AppenderClient {
	c.request.SetBody(body)
	return c
}

func (c *AppenderClient) AppendBodyAsString(body string) *AppenderClient {
	c.request.SetBodyString(body)
	return c
}

func (c *AppenderClient) AppendBody(body interface{}) *AppenderClient {
	if bodyContent, err := json.Marshal(body); err == nil {
		c.request.SetBody(bodyContent)
	} else {
		panic(err)
	}
	return c
}

func (p *Request) Get(uri string) *AppenderClient {
	p.request.Header.SetMethod("GET")
	p.request.SetRequestURI(uri)
	return &AppenderClient{
		request: p.request,
		echo:    p.echo,
		client:  p.client,
	}
}

func (p *Request) Post(uri string) *AppenderClient {
	p.request.Header.SetMethod("POST")
	p.request.SetRequestURI(uri)
	return &AppenderClient{
		request: p.request,
		client:  p.client,
	}
}

func (p *Request) Put(uri string) *AppenderClient {
	p.request.Header.SetMethod("PUT")
	p.request.SetRequestURI(uri)
	return &AppenderClient{
		request: p.request,
		client:  p.client,
	}
}

func (p *Request) Delete(uri string) *AppenderClient {
	p.request.Header.SetMethod("PUT")
	p.request.SetRequestURI(uri)
	return &AppenderClient{
		request: p.request,
		client:  p.client,
	}
}

func (p *Request) Patch(uri string) *AppenderClient {
	p.request.Header.SetMethod("PATCH")
	p.request.SetRequestURI(uri)
	return &AppenderClient{
		request: p.request,
		client:  p.client,
	}
}

func (c *AppenderClient) BuildRequest() *BuildContext {
	return &BuildContext{
		client:  c.client,
		request: c.request,
	}
}

func (c *AppenderClient) BuildRequestWithEcho(ctx echo.Context) *BuildContext {
	return &BuildContext{
		request: c.request,
		echo:    ctx,
		client:  c.client,
	}
}

type Response struct {
	StatusCode int
}

func (c *BuildContext) Call(response interface{}) (*Response, error) {
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(c.request)
	defer fasthttp.ReleaseResponse(resp)

	if c.dependencyResolver != nil {
		headerKey, resolvedContent := c.dependencyResolver.ResolveContent()
		c.request.Header.Add(headerKey, resolvedContent)
	}
	if c.echo != nil {
		c.request.Header.Add("x-correlation-id", c.echo.Request().Header.Get("x-correlation-id"))
	}
	requestError := c.client.Do(c.request, resp)
	if requestError != nil {
		return nil, requestError
	}
	if resp.StatusCode() != fasthttp.StatusOK {
		return &Response{StatusCode: resp.StatusCode()}, errors.New(fmt.Sprintf("Expected status code %d but got %d\n", fasthttp.StatusOK, resp.StatusCode()))
	}

	// Do we need to decompress the response?
	contentEncoding := resp.Header.Peek("Content-Encoding")
	var body []byte
	if bytes.EqualFold(contentEncoding, []byte("gzip")) {
		body, _ = resp.BodyGunzip()
	}
	if bytes.EqualFold(contentEncoding, []byte("brotli")) {
		body, _ = resp.BodyUnbrotli()
	} else {
		body = resp.Body()
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	if c.dependencySetter != nil {
		var realResponse interface{}
		if err := json.Unmarshal(body, &realResponse); err != nil {
			return nil, err
		}
		c.dependencySetter.SetContent(&realResponse)
	}
	return &Response{StatusCode: resp.StatusCode()}, nil
}
