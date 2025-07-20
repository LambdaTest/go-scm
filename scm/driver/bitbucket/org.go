// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitbucket

import (
	"context"
	"fmt"

	"github.com/drone/go-scm/scm"
)

type organizationService struct {
	client *wrapper
}

func (s *organizationService) Find(ctx context.Context, name string) (*scm.Organization, *scm.Response, error) {
	path := fmt.Sprintf("2.0/workspaces/%s", name)
	out := new(organization)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertOrganization(out), res, err
}

func (s *organizationService) FindMembership(ctx context.Context, name, username string) (*scm.Membership, *scm.Response, error) {
	return nil, nil, scm.ErrNotSupported
}

func (s *organizationService) ListMemberships(ctx context.Context, orgNameList []string, username string, opts scm.ListOptions) ([]*scm.Membership, *scm.Response, error) {
	path := fmt.Sprintf("2.0/user/permissions/workspaces?%s", encodeListRoleOptions(opts))
	out := new(membershipList)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	copyPagination(out.pagination, res)
	return convertMembershipList(out), res, err
}

func (s *organizationService) List(ctx context.Context, opts scm.ListOptions) ([]*scm.Organization, *scm.Response, error) {
	path := fmt.Sprintf("2.0/workspaces?%s", encodeListRoleOptions(opts))
	out := new(organizationList)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	copyPagination(out.pagination, res)
	return convertOrganizationList(out), res, err
}

func convertOrganizationList(from *organizationList) []*scm.Organization {
	to := []*scm.Organization{}
	for _, v := range from.Values {
		to = append(to, convertOrganization(v))
	}
	return to
}

func convertMembershipList(from *membershipList) []*scm.Membership {
	to := []*scm.Membership{}
	for _, v := range from.Values {
		to = append(to, convertMembership(v))
	}
	return to
}

type organizationList struct {
	pagination
	Values []*organization `json:"values"`
}

type membershipList struct {
	pagination
	Values []*membership `json:"values"`
}

type membership struct {
	Permission string `json:"permission"`
	Workspace  struct {
		Slug string `json:"slug"`
	} `json:"workspace"`
}

type organization struct {
	Login string `json:"slug"`
}

func convertOrganization(from *organization) *scm.Organization {
	return &scm.Organization{
		Name:   from.Login,
		Avatar: fmt.Sprintf("https://bitbucket.org/account/%s/avatar/32/", from.Login),
	}
}

func convertMembership(from *membership) *scm.Membership {
	to := new(scm.Membership)
	to.Active = true

	switch from.Permission {
	case "owner":
		to.Role = scm.RoleAdmin
	case "collaborator":
		to.Role = scm.RoleMember
	default:
		to.Role = scm.RoleViewer
	}

	to.Organization.Name = from.Workspace.Slug
	return to
}
