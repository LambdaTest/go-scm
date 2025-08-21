// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package github implements a GitHub client.
package github

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/transport/proxy"
)

// New returns a new GitHub API client.
// This function maintains backward compatibility and creates a client without proxy.
func New(uri string) (*scm.Client, error) {
	return NewWithProxy(uri, "")
}

// NewWithProxy returns a new GitHub API client with optional proxy support.
// If proxyURL is empty or nil, no proxy will be used.
// If proxyURL is provided, all HTTP requests will be routed through the specified proxy.
func NewWithProxy(uri, proxyURL string) (*scm.Client, error) {
	base, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(base.Path, "/") {
		base.Path = base.Path + "/"
	}
	// home := base.String()
	// if home == "https://api.github.com/" {
	// 	home = "https://github.com/"
	// }

	client := &wrapper{new(scm.Client)}
	client.BaseURL = base
	// initialize services
	client.Driver = scm.DriverGithub
	client.Linker = &linker{websiteAddress(base)}
	client.Apps = &appsService{client}
	client.Contents = &contentService{client}
	client.Git = &gitService{client}
	client.Issues = &issueService{client}
	client.Milestones = &milestoneService{client}
	client.Organizations = &organizationService{client}
	client.PullRequests = &pullService{&issueService{client}}
	client.Repositories = &RepositoryService{client}
	client.Releases = &releaseService{client}
	client.Reviews = &reviewService{client}
	client.Users = &userService{client}
	client.Webhooks = &webhookService{client}
<<<<<<< Updated upstream
	
	client.Client.Client = &http.Client{Transport: http.DefaultTransport}
=======
>>>>>>> Stashed changes

	if proxyURL != "" {
		transport, err := proxy.NewTransport(http.DefaultTransport, proxyURL)
		if err != nil {
			return nil, err
		}
		client.Client.Client.Transport = transport
	}

	return client.Client, nil
}

// NewDefault returns a new GitHub API client using the
// default api.github.com address.
func NewDefault() *scm.Client {
	client, _ := New("https://api.github.com")
	return client
}

// NewDefaultWithProxy returns a new GitHub API client using the
// default api.github.com address with optional proxy support.
// If proxyURL is empty or nil, no proxy will be used.
// If proxyURL is provided, all HTTP requests will be routed through the specified proxy.
func NewDefaultWithProxy(proxyURL string) *scm.Client {
	client, _ := NewWithProxy("https://api.github.com", proxyURL)
	return client
}

// wraper wraps the Client to provide high level helper functions
// for making http requests and unmarshaling the response.
type wrapper struct {
	*scm.Client
}

// do wraps the Client.Do function by creating the Request and
// unmarshalling the response.
func (c *wrapper) do(ctx context.Context, method, path string, in, out interface{}) (*scm.Response, error) {
	req := &scm.Request{
		Method: method,
		Path:   path,
	}
	// if we are posting or putting data, we need to
	// write it to the body of the request.
	if in != nil {
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(in)
		req.Header = map[string][]string{
			"Content-Type": {"application/json"},
		}
		req.Body = buf
	}

	// execute the http request
	res, err := c.Client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// parse the github request id.
	res.ID = res.Header.Get("X-GitHub-Request-Id")

	// parse the github rate limit details.
	res.Rate.Limit, _ = strconv.Atoi(
		res.Header.Get("X-RateLimit-Limit"),
	)
	res.Rate.Remaining, _ = strconv.Atoi(
		res.Header.Get("X-RateLimit-Remaining"),
	)
	res.Rate.Reset, _ = strconv.ParseInt(
		res.Header.Get("X-RateLimit-Reset"), 10, 64,
	)

	// snapshot the request rate limit
	c.Client.SetRate(res.Rate)

	// if an error is encountered, unmarshal and return the
	// error response.
	if res.Status > 300 {
		err := new(Error)
		json.NewDecoder(res.Body).Decode(err)
		return res, err
	}

	if out == nil {
		return res, nil
	}

	// if a json response is expected, parse and return
	// the json response.
	return res, json.NewDecoder(res.Body).Decode(out)
}

// Error represents a Github error.
type Error struct {
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

// helper function converts the github API url to
// the website url.
func websiteAddress(u *url.URL) string {
	host, proto := u.Host, u.Scheme
	switch host {
	case "api.github.com":
		return "https://github.com/"
	}
	return proto + "://" + host + "/"
}
