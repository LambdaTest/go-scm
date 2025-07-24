// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"time"

	"github.com/drone/go-scm/scm"
)

type appsService struct {
	client *wrapper
}

type app struct {
	ID          int64                    `json:"id,omitempty"`
	Slug        string                   `json:"slug,omitempty"`
	NodeID      string                   `json:"node_id,omitempty"`
	Owner       *user                    `json:"owner,omitempty"`
	Name        string                   `json:"name,omitempty"`
	Description string                   `json:"description,omitempty"`
	ExternalURL string                   `json:"external_url,omitempty"`
	HTMLURL     string                   `json:"html_url,omitempty"`
	CreatedAt   time.Time                `json:"created_at,omitempty"`
	UpdatedAt   time.Time                `json:"updated_at,omitempty"`
	Permissions *installationPermissions `json:"permissions,omitempty"`
	Events      []string                 `json:"events,omitempty"`
}

type installation struct {
	ID                     int64                    `json:"id,omitempty"`
	NodeID                 string                   `json:"node_id,omitempty"`
	AppID                  int64                    `json:"app_id,omitempty"`
	AppSlug                string                   `json:"app_slug,omitempty"`
	TargetID               int64                    `json:"target_id,omitempty"`
	Account                *user                    `json:"account,omitempty"`
	AccessTokensURL        string                   `json:"access_tokens_url,omitempty"`
	RepositoriesURL        string                   `json:"repositories_url,omitempty"`
	HTMLURL                string                   `json:"html_url,omitempty"`
	TargetType             string                   `json:"target_type,omitempty"`
	SingleFileName         string                   `json:"single_file_name,omitempty"`
	RepositorySelection    string                   `json:"repository_selection,omitempty"`
	Events                 []string                 `json:"events,omitempty"`
	SingleFilePaths        []string                 `json:"single_file_paths,omitempty"`
	Permissions            *installationPermissions `json:"permissions,omitempty"`
	CreatedAt              time.Time                `json:"created_at,omitempty"`
	UpdatedAt              time.Time                `json:"updated_at,omitempty"`
	HasMultipleSingleFiles *bool                    `json:"has_multiple_single_files,omitempty"`
	SuspendedBy            *user                    `json:"suspended_by,omitempty"`
	SuspendedAt            *time.Time               `json:"suspended_at,omitempty"`
}

type installationToken struct {
	Token        string                   `json:"token,omitempty"`
	ExpiresAt    time.Time                `json:"expires_at,omitempty"`
	Permissions  *installationPermissions `json:"permissions,omitempty"`
	Repositories []*repository            `json:"repositories,omitempty"`
}

type installationPermissions struct {
	Actions                       string `json:"actions,omitempty"`
	Administration                string `json:"administration,omitempty"`
	Blocking                      string `json:"blocking,omitempty"`
	Checks                        string `json:"checks,omitempty"`
	Contents                      string `json:"contents,omitempty"`
	ContentReferences             string `json:"content_references,omitempty"`
	Deployments                   string `json:"deployments,omitempty"`
	Emails                        string `json:"emails,omitempty"`
	Environments                  string `json:"environments,omitempty"`
	Followers                     string `json:"followers,omitempty"`
	Issues                        string `json:"issues,omitempty"`
	Metadata                      string `json:"metadata,omitempty"`
	Members                       string `json:"members,omitempty"`
	OrganizationAdministration    string `json:"organization_administration,omitempty"`
	OrganizationCustomRoles       string `json:"organization_custom_roles,omitempty"`
	OrganizationHooks             string `json:"organization_hooks,omitempty"`
	OrganizationPackages          string `json:"organization_packages,omitempty"`
	OrganizationPlan              string `json:"organization_plan,omitempty"`
	OrganizationPreReceiveHooks   string `json:"organization_pre_receive_hooks,omitempty"`
	OrganizationProjects          string `json:"organization_projects,omitempty"`
	OrganizationSecrets           string `json:"organization_secrets,omitempty"`
	OrganizationSelfHostedRunners string `json:"organization_self_hosted_runners,omitempty"`
	OrganizationUserBlocking      string `json:"organization_user_blocking,omitempty"`
	Packages                      string `json:"packages,omitempty"`
	Pages                         string `json:"pages,omitempty"`
	PullRequests                  string `json:"pull_requests,omitempty"`
	RepositoryHooks               string `json:"repository_hooks,omitempty"`
	RepositoryProjects            string `json:"repository_projects,omitempty"`
	RepositoryPreReceiveHooks     string `json:"repository_pre_receive_hooks,omitempty"`
	Secrets                       string `json:"secrets,omitempty"`
	SecretScanningAlerts          string `json:"secret_scanning_alerts,omitempty"`
	SecurityEvents                string `json:"security_events,omitempty"`
	SingleFile                    string `json:"single_file,omitempty"`
	Statuses                      string `json:"statuses,omitempty"`
	TeamDiscussions               string `json:"team_discussions,omitempty"`
	VulnerabilityAlerts           string `json:"vulnerability_alerts,omitempty"`
	Workflows                     string `json:"workflows,omitempty"`
}

func (s *appsService) Get(ctx context.Context, appSlug string) (*scm.App, *scm.Response, error) {
	var path string
	if appSlug != "" {
		path = fmt.Sprintf("apps/%s", appSlug)
	} else {
		path = "app"
	}

	out := new(app)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertApp(out), res, err
}

func (s *appsService) ListInstallations(ctx context.Context, opts scm.ListOptions) ([]*scm.Installation, *scm.Response, error) {
	path := fmt.Sprintf("app/installations?%s", encodeListOptions(opts))
	out := []*installation{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertInstallationList(out), res, err
}

func (s *appsService) GetInstallation(ctx context.Context, id int64) (*scm.Installation, *scm.Response, error) {
	path := fmt.Sprintf("app/installations/%d", id)
	out := new(installation)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertInstallation(out), res, err
}

func (s *appsService) ListUserInstallations(ctx context.Context, opts scm.ListOptions) ([]*scm.Installation, *scm.Response, error) {
	path := fmt.Sprintf("user/installations?%s", encodeListOptions(opts))
	var result struct {
		Installations []*installation `json:"installations"`
	}
	res, err := s.client.do(ctx, "GET", path, nil, &result)
	return convertInstallationList(result.Installations), res, err
}

func (s *appsService) SuspendInstallation(ctx context.Context, id int64) (*scm.Response, error) {
	path := fmt.Sprintf("app/installations/%d/suspended", id)
	res, err := s.client.do(ctx, "PUT", path, nil, nil)
	return res, err
}

func (s *appsService) UnsuspendInstallation(ctx context.Context, id int64) (*scm.Response, error) {
	path := fmt.Sprintf("app/installations/%d/suspended", id)
	res, err := s.client.do(ctx, "DELETE", path, nil, nil)
	return res, err
}

func (s *appsService) DeleteInstallation(ctx context.Context, id int64) (*scm.Response, error) {
	path := fmt.Sprintf("app/installations/%d", id)
	res, err := s.client.do(ctx, "DELETE", path, nil, nil)
	return res, err
}

func (s *appsService) CreateInstallationToken(ctx context.Context, id int64, opts *scm.InstallationTokenOptions) (*scm.InstallationToken, *scm.Response, error) {
	path := fmt.Sprintf("app/installations/%d/access_tokens", id)
	out := new(installationToken)
	res, err := s.client.do(ctx, "POST", path, opts, out)
	return convertInstallationToken(out), res, err
}

func (s *appsService) FindOrganizationInstallation(ctx context.Context, org string) (*scm.Installation, *scm.Response, error) {
	path := fmt.Sprintf("orgs/%s/installation", org)
	out := new(installation)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertInstallation(out), res, err
}

func (s *appsService) FindRepositoryInstallation(ctx context.Context, owner, repo string) (*scm.Installation, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/%s/installation", owner, repo)
	out := new(installation)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertInstallation(out), res, err
}

func (s *appsService) FindRepositoryInstallationByID(ctx context.Context, id int64) (*scm.Installation, *scm.Response, error) {
	path := fmt.Sprintf("repositories/%d/installation", id)
	out := new(installation)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertInstallation(out), res, err
}

func (s *appsService) FindUserInstallation(ctx context.Context, user string) (*scm.Installation, *scm.Response, error) {
	path := fmt.Sprintf("users/%s/installation", user)
	out := new(installation)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertInstallation(out), res, err
}

func convertApp(from *app) *scm.App {
	if from == nil {
		return nil
	}
	return &scm.App{
		ID:          from.ID,
		Slug:        from.Slug,
		NodeID:      from.NodeID,
		Owner:       convertUser(from.Owner),
		Name:        from.Name,
		Description: from.Description,
		ExternalURL: from.ExternalURL,
		HTMLURL:     from.HTMLURL,
		CreatedAt:   from.CreatedAt,
		UpdatedAt:   from.UpdatedAt,
		Permissions: convertInstallationPermissions(from.Permissions),
		Events:      from.Events,
	}
}

func convertInstallation(from *installation) *scm.Installation {
	if from == nil {
		return nil
	}
	hasMultipleSingleFiles := false
	if from.HasMultipleSingleFiles != nil {
		hasMultipleSingleFiles = *from.HasMultipleSingleFiles
	}
	return &scm.Installation{
		ID:                     from.ID,
		NodeID:                 from.NodeID,
		AppID:                  from.AppID,
		AppSlug:                from.AppSlug,
		TargetID:               from.TargetID,
		Account:                convertUser(from.Account),
		AccessTokensURL:        from.AccessTokensURL,
		RepositoriesURL:        from.RepositoriesURL,
		HTMLURL:                from.HTMLURL,
		TargetType:             from.TargetType,
		SingleFileName:         from.SingleFileName,
		RepositorySelection:    from.RepositorySelection,
		Events:                 from.Events,
		SingleFilePaths:        from.SingleFilePaths,
		Permissions:            convertInstallationPermissions(from.Permissions),
		CreatedAt:              from.CreatedAt,
		UpdatedAt:              from.UpdatedAt,
		HasMultipleSingleFiles: hasMultipleSingleFiles,
		SuspendedBy:            convertUser(from.SuspendedBy),
		SuspendedAt:            from.SuspendedAt,
	}
}

func convertInstallationList(from []*installation) []*scm.Installation {
	to := []*scm.Installation{}
	for _, v := range from {
		to = append(to, convertInstallation(v))
	}
	return to
}

func convertInstallationToken(from *installationToken) *scm.InstallationToken {
	if from == nil {
		return nil
	}
	return &scm.InstallationToken{
		Token:        from.Token,
		ExpiresAt:    from.ExpiresAt,
		Permissions:  convertInstallationPermissions(from.Permissions),
		Repositories: convertRepositoryList(from.Repositories),
	}
}

func convertInstallationPermissions(from *installationPermissions) *scm.InstallationPermissions {
	if from == nil {
		return nil
	}
	return &scm.InstallationPermissions{
		Actions:                       from.Actions,
		Administration:                from.Administration,
		Blocking:                      from.Blocking,
		Checks:                        from.Checks,
		Contents:                      from.Contents,
		ContentReferences:             from.ContentReferences,
		Deployments:                   from.Deployments,
		Emails:                        from.Emails,
		Environments:                  from.Environments,
		Followers:                     from.Followers,
		Issues:                        from.Issues,
		Metadata:                      from.Metadata,
		Members:                       from.Members,
		OrganizationAdministration:    from.OrganizationAdministration,
		OrganizationCustomRoles:       from.OrganizationCustomRoles,
		OrganizationHooks:             from.OrganizationHooks,
		OrganizationPackages:          from.OrganizationPackages,
		OrganizationPlan:              from.OrganizationPlan,
		OrganizationPreReceiveHooks:   from.OrganizationPreReceiveHooks,
		OrganizationProjects:          from.OrganizationProjects,
		OrganizationSecrets:           from.OrganizationSecrets,
		OrganizationSelfHostedRunners: from.OrganizationSelfHostedRunners,
		OrganizationUserBlocking:      from.OrganizationUserBlocking,
		Packages:                      from.Packages,
		Pages:                         from.Pages,
		PullRequests:                  from.PullRequests,
		RepositoryHooks:               from.RepositoryHooks,
		RepositoryProjects:            from.RepositoryProjects,
		RepositoryPreReceiveHooks:     from.RepositoryPreReceiveHooks,
		Secrets:                       from.Secrets,
		SecretScanningAlerts:          from.SecretScanningAlerts,
		SecurityEvents:                from.SecurityEvents,
		SingleFile:                    from.SingleFile,
		Statuses:                      from.Statuses,
		TeamDiscussions:               from.TeamDiscussions,
		VulnerabilityAlerts:           from.VulnerabilityAlerts,
		Workflows:                     from.Workflows,
	}
}
