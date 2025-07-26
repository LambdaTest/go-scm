// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/h2non/gock"
)

func TestAppsCreateInstallationToken(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Post("/app/installations/123/access_tokens").
		Reply(201).
		Type("application/json").
		File("testdata/installation_token.json")

	client := NewDefault()
	opts := &scm.InstallationTokenOptions{
		RepositoryIDs: []int64{123456},
		Repositories:  []string{"owner/repo"},
		Permissions: &scm.InstallationPermissions{
			Contents: "read",
			Issues:   "write",
		},
	}

	token, res, err := client.Apps.CreateInstallationToken(context.Background(), 123, opts)
	if err != nil {
		t.Error(err)
		return
	}

	if got, want := res.Status, 201; got != want {
		t.Errorf("Want response status %d, got %d", want, got)
	}

	if got, want := token.Token, "ghs_16C7e42F292c6912E7710c838347Ae178B4a"; got != want {
		t.Errorf("Want token %s, got %s", want, got)
	}

	if got, want := len(token.Repositories), 1; got != want {
		t.Errorf("Want %d repositories, got %d", want, got)
	}

	if got, want := token.Repositories[0].Name, "Hello-World"; got != want {
		t.Errorf("Want repository name %s, got %s", want, got)
	}
}

func TestAppsGet(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/app").
		Reply(200).
		Type("application/json").
		File("testdata/app.json")

	client := NewDefault()
	app, res, err := client.Apps.Get(context.Background(), "")
	if err != nil {
		t.Error(err)
		return
	}

	if got, want := res.Status, 200; got != want {
		t.Errorf("Want response status %d, got %d", want, got)
	}

	if got, want := app.ID, int64(1); got != want {
		t.Errorf("Want app ID %d, got %d", want, got)
	}

	if got, want := app.Name, "octocat-app"; got != want {
		t.Errorf("Want app name %s, got %s", want, got)
	}
}

func TestAppsGetInstallation(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/app/installations/123").
		Reply(200).
		Type("application/json").
		File("testdata/installation.json")

	client := NewDefault()
	installation, res, err := client.Apps.GetInstallation(context.Background(), 123)
	if err != nil {
		t.Error(err)
		return
	}

	if got, want := res.Status, 200; got != want {
		t.Errorf("Want response status %d, got %d", want, got)
	}

	if got, want := installation.ID, int64(123); got != want {
		t.Errorf("Want installation ID %d, got %d", want, got)
	}

	if got, want := installation.AppID, int64(1); got != want {
		t.Errorf("Want app ID %d, got %d", want, got)
	}
}

func TestAppsListInstallations(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/app/installations").
		Reply(200).
		Type("application/json").
		File("testdata/installations.json")

	client := NewDefault()
	installations, res, err := client.Apps.ListInstallations(context.Background(), scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	if got, want := res.Status, 200; got != want {
		t.Errorf("Want response status %d, got %d", want, got)
	}

	if got, want := len(installations), 1; got != want {
		t.Errorf("Want %d installations, got %d", want, got)
	}

	if got, want := installations[0].ID, int64(123); got != want {
		t.Errorf("Want installation ID %d, got %d", want, got)
	}
}

func TestAppsFindRepositoryInstallation(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/Hello-World/installation").
		Reply(200).
		Type("application/json").
		File("testdata/installation.json")

	client := NewDefault()
	installation, res, err := client.Apps.FindRepositoryInstallation(context.Background(), "octocat", "Hello-World")
	if err != nil {
		t.Error(err)
		return
	}

	if got, want := res.Status, 200; got != want {
		t.Errorf("Want response status %d, got %d", want, got)
	}

	if got, want := installation.ID, int64(123); got != want {
		t.Errorf("Want installation ID %d, got %d", want, got)
	}
}
