// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitlab

import (
	"context"
	"fmt"
	"strconv"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/internal/null"
)

type organizationService struct {
	client *wrapper
}

func (s *organizationService) Find(ctx context.Context, name string) (*scm.Organization, *scm.Response, error) {
	path := fmt.Sprintf("api/v4/groups/%s", name)
	out := new(organization)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertOrganization(out), res, err
}

func (s *organizationService) FindMembership(ctx context.Context, name, userID string) (*scm.Membership, *scm.Response, error) {
	path := fmt.Sprintf("api/v4/groups/%s/billable_members/%s/memberships", name, userID)
	out := new(membership)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertMembership(out), res, err
}

func (s *organizationService) ListMemberships(ctx context.Context, orgNameList []string, username string, opts scm.ListOptions) ([]*scm.Membership, *scm.Response, error) {
	userInfo := new(user)
	out := []*scm.Membership{}
	res, err := s.client.do(ctx, "GET", "/api/v4/user", nil, userInfo)
	if err != nil {
		return nil, res, err
	}

	for _, v := range orgNameList {
		var member *scm.Membership
		if username == v {
			member = &scm.Membership{
				Role:   scm.RoleAdmin,
				Active: true,
				Organization: scm.Organization{
					Name:   username,
					Avatar: "",
				},
			}
		} else {
			member, _, err = s.FindMembership(ctx, v, strconv.Itoa(userInfo.ID))
			if err != nil {
				return nil, res, err
			}
		}

		out = append(out, member)
	}

	return out, nil, nil
}

func (s *organizationService) List(ctx context.Context, opts scm.ListOptions) ([]*scm.Organization, *scm.Response, error) {
	path := fmt.Sprintf("api/v4/groups?%s", encodeListOptions(opts))
	out := []*organization{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertOrganizationList(out), res, err
}

type organization struct {
	Name   string      `json:"name"`
	Path   string      `json:"path"`
	Avatar null.String `json:"avatar_url"`
}

func convertOrganizationList(from []*organization) []*scm.Organization {
	to := []*scm.Organization{}
	for _, v := range from {
		to = append(to, convertOrganization(v))
	}
	return to
}

func convertOrganization(from *organization) *scm.Organization {
	return &scm.Organization{
		Name:   from.Path,
		Avatar: from.Avatar.String,
	}
}

func convertMembership(from *membership) *scm.Membership {
	to := new(scm.Membership)
	to.Active = true
	to.Organization.Name = from.OrgName

	switch from.Access.Role {
	case "Owner":
		to.Role = scm.RoleAdmin
	case "Maintainer", "Developer":
		to.Role = scm.RoleMember
	default:
		to.Role = scm.RoleViewer
	}

	return to
}

type membership struct {
	OrgName string `json:"source_full_name"`
	Access  struct {
		Role string `json:"string_value"`
	} `json:"access_level"`
}
