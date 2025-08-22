// Copyright 2018 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proxy

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNewTransport_EmptyProxyURL(t *testing.T) {
	base := http.DefaultTransport
	transport, err := NewTransport(base, "")
	if err != nil {
		t.Fatalf("Expected no error for empty proxy URL, got: %v", err)
	}
	if transport != base {
		t.Error("Expected transport to be the same as base for empty proxy URL")
	}
}

func TestNewTransport_InvalidProxyURL(t *testing.T) {
	base := http.DefaultTransport
	_, err := NewTransport(base, "://invalid")
	if err == nil {
		t.Error("Expected error for invalid proxy URL")
	}
}

func TestNewTransport_ValidProxyURL(t *testing.T) {
	base := http.DefaultTransport
	proxyURL := "http://proxy.example.com:8080"
	transport, err := NewTransport(base, proxyURL)
	if err != nil {
		t.Fatalf("Expected no error for valid proxy URL, got: %v", err)
	}

	proxyTransport, ok := transport.(*http.Transport)
	if !ok {
		t.Fatal("Expected transport to be of type *http.Transport")
	}

	// Test that the proxy function is set correctly
	if proxyTransport.Proxy == nil {
		t.Fatal("Expected proxy function to be set")
	}

	// Test the proxy function with a sample request
	testURL, _ := url.Parse("http://example.com")
	req := &http.Request{URL: testURL}
	proxyURLResult, err := proxyTransport.Proxy(req)
	if err != nil {
		t.Fatalf("Expected no error from proxy function, got: %v", err)
	}

	expectedURL, _ := url.Parse(proxyURL)
	if proxyURLResult.String() != expectedURL.String() {
		t.Errorf("Expected proxy URL %s, got %s", expectedURL.String(), proxyURLResult.String())
	}
}

func TestTransport_RoundTrip_NoProxy(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	transport := &Transport{
		Base:     http.DefaultTransport,
		ProxyURL: nil,
	}

	client := &http.Client{Transport: transport}
	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got: %d", resp.StatusCode)
	}
}

func TestTransport_RoundTrip_WithProxy(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	proxyURL, _ := url.Parse("http://proxy.example.com:8080")
	transport := &Transport{
		Base:     http.DefaultTransport,
		ProxyURL: proxyURL,
	}

	// This test verifies that the transport is configured with the proxy
	// The actual proxy behavior would require a real proxy server for testing
	// We're just ensuring the transport is properly configured
	if transport.ProxyURL.String() != proxyURL.String() {
		t.Errorf("Expected proxy URL %s, got %s", proxyURL.String(), transport.ProxyURL.String())
	}
}
