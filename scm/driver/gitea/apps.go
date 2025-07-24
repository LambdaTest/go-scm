// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitea

import (
	"context"

	"github.com/drone/go-scm/scm"
)

type appsService struct {
	client *wrapper
}

func (s *appsService) Get(ctx context.Context, appSlug string) (*scm.App, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *appsService) ListInstallations(ctx context.Context, opts scm.ListOptions) ([]*scm.Installation, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *appsService) GetInstallation(ctx context.Context, id int64) (*scm.Installation, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *appsService) ListUserInstallations(ctx context.Context, opts scm.ListOptions) ([]*scm.Installation, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *appsService) SuspendInstallation(ctx context.Context, id int64) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *appsService) UnsuspendInstallation(ctx context.Context, id int64) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *appsService) DeleteInstallation(ctx context.Context, id int64) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *appsService) CreateInstallationToken(ctx context.Context, id int64, opts *scm.InstallationTokenOptions) (*scm.InstallationToken, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *appsService) FindOrganizationInstallation(ctx context.Context, org string) (*scm.Installation, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *appsService) FindRepositoryInstallation(ctx context.Context, owner, repo string) (*scm.Installation, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *appsService) FindRepositoryInstallationByID(ctx context.Context, id int64) (*scm.Installation, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *appsService) FindUserInstallation(ctx context.Context, user string) (*scm.Installation, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}
