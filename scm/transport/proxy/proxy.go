// Copyright 2018 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proxy

import (
	"net/http"
	"net/url"
)

// Transport is an http.RoundTripper that makes HTTP
// requests through a proxy, wrapping a base RoundTripper
type Transport struct {
	Base     http.RoundTripper
	ProxyURL *url.URL
}

// RoundTrip makes the request through the configured proxy.
func (t *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	// If no proxy is configured, use the base transport
	if t.ProxyURL == nil {
		return t.base().RoundTrip(r)
	}

	// Create a new transport with the proxy configuration
	proxyTransport := &http.Transport{
		Proxy: func(_ *http.Request) (*url.URL, error) {
			return t.ProxyURL, nil
		},
	}

	// If we have a base transport, copy its configuration
	if t.Base != nil {
		if baseTransport, ok := t.Base.(*http.Transport); ok {
			proxyTransport.TLSClientConfig = baseTransport.TLSClientConfig
			proxyTransport.DialContext = baseTransport.DialContext
			proxyTransport.MaxIdleConns = baseTransport.MaxIdleConns
			proxyTransport.MaxIdleConnsPerHost = baseTransport.MaxIdleConnsPerHost
			proxyTransport.IdleConnTimeout = baseTransport.IdleConnTimeout
			proxyTransport.TLSHandshakeTimeout = baseTransport.TLSHandshakeTimeout
			proxyTransport.ExpectContinueTimeout = baseTransport.ExpectContinueTimeout
		}
	}

	return proxyTransport.RoundTrip(r)
}

// base returns the base transport. If no base transport
// is configured, the default transport is returned.
func (t *Transport) base() http.RoundTripper {
	if t.Base != nil {
		return t.Base
	}
	return http.DefaultTransport
}

// NewTransport creates a new proxy transport with the given proxy URL.
// If proxyURL is empty or nil, it returns the base transport unchanged.
func NewTransport(base http.RoundTripper, proxyURL string) (http.RoundTripper, error) {
	if proxyURL == "" {
		return base, nil
	}

	parsedURL, err := url.Parse(proxyURL)
	if err != nil {
		return nil, err
	}

	return &Transport{
		Base:     base,
		ProxyURL: parsedURL,
	}, nil
}
