package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Transport struct {
	OverrideURL      *url.URL
	RequestCallback  func(*http.Request)
	ResponseCallback func(*http.Response)
	*http.Transport
}

func (t *Transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	if t.RequestCallback != nil {
		t.RequestCallback(req)
	}

	if t.OverrideURL != nil {
		port := "80"
		if s := strings.Split(req.URL.Host, ":"); len(s) > 1 {
			port = s[1]
		}

		req = cloneRequest(req)
		req.Header.Set("X-Original-Scheme", req.URL.Scheme)
		req.Header.Set("X-Original-Port", port)
		req.URL.Scheme = t.OverrideURL.Scheme
		req.URL.Host = t.OverrideURL.Host
	}

	resp, err = t.Transport.RoundTrip(req)
	if err == nil && t.ResponseCallback != nil {
		t.ResponseCallback(resp)
	}

	return
}

type Client struct {
	rootUrl    *url.URL
	httpClient *http.Client
}

func NewClient(rootUrl *url.URL, configure func(*Transport)) *Client {
	tr := &Transport{
		Transport: &http.Transport{
			Proxy: proxyFromEnvironment,
		},
	}
	if configure != nil {
		configure(tr)
	}

	maxRedirects := 10
	checkRedirect := func(req *http.Request, via []*http.Request) error {
		if req.Host == "" {
			req.Host = req.URL.Host
		}

		for key, val := range via[0].Header {
			if req.Host == via[0].Host || !strings.EqualFold(key, "Authorization") {
				req.Header[key] = val
			}
		}

		if len(via) >= maxRedirects {
			return fmt.Errorf("stopped after %d redirects", maxRedirects)
		} else {
			return nil
		}
	}

	return &Client{
		rootUrl: rootUrl,
		httpClient: &http.Client{
			Transport:     tr,
			CheckRedirect: checkRedirect,
		},
	}
}

func (c *Client) PerformRequest(method, path string, body io.Reader, configure func(*http.Request)) (*http.Response, error) {
	url, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	url = c.rootUrl.ResolveReference(url)
	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return nil, err
	}
	if configure != nil {
		configure(req)
	}

	return c.httpClient.Do(req)
}

func cloneRequest(req *http.Request) *http.Request {
	dup := new(http.Request)
	*dup = *req
	dup.URL, _ = url.Parse(req.URL.String())
	dup.Header = make(http.Header)
	for k, s := range req.Header {
		dup.Header[k] = s
	}
	return dup
}

// An implementation of http.ProxyFromEnvironment that isn't broken
func proxyFromEnvironment(req *http.Request) (*url.URL, error) {
	var proxy string

	switch req.URL.Scheme {
	case "http":
		proxy = os.Getenv("http_proxy")
	case "https":
		if proxy = os.Getenv("https_proxy"); proxy == "" {
			proxy = os.Getenv("HTTPS_PROXY")
		}
	}

	if proxy == "" {
		if proxy = os.Getenv("all_proxy"); proxy == "" {
			proxy = os.Getenv("ALL_PROXY")
		}
	}

	if proxy == "" {
		return nil, nil
	} else {
		if !strings.Contains(proxy, "://") {
			proxy = "http://" + proxy
		}

		proxyURL, err := url.Parse(proxy)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy address %q: %v", proxy, err)
		}

		return proxyURL, nil
	}
}
