// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scm

import (
	"context"
	"time"
)

type (
	// App represents an SCM App.
	App struct {
		ID          int64
		Slug        string
		NodeID      string
		Owner       *User
		Name        string
		Description string
		ExternalURL string
		HTMLURL     string
		CreatedAt   time.Time
		UpdatedAt   time.Time
		Permissions *InstallationPermissions
		Events      []string
	}

	// Installation represents an Apps installation.
	Installation struct {
		ID                     int64
		NodeID                 string
		AppID                  int64
		AppSlug                string
		TargetID               int64
		Account                *User
		AccessTokensURL        string
		RepositoriesURL        string
		HTMLURL                string
		TargetType             string
		SingleFileName         string
		RepositorySelection    string
		Events                 []string
		SingleFilePaths        []string
		Permissions            *InstallationPermissions
		CreatedAt              time.Time
		UpdatedAt              time.Time
		HasMultipleSingleFiles bool
		SuspendedBy            *User
		SuspendedAt            *time.Time
	}

	// InstallationToken represents an installation token.
	InstallationToken struct {
		Token        string
		ExpiresAt    time.Time
		Permissions  *InstallationPermissions
		Repositories []*Repository
	}

	// InstallationTokenOptions allow restricting a token's access to specific repositories.
	InstallationTokenOptions struct {
		// The IDs of the repositories that the installation token can access.
		// Providing repository IDs restricts the access of an installation token to specific repositories.
		RepositoryIDs []int64 `json:"repository_ids,omitempty"`

		// The names of the repositories that the installation token can access.
		// Providing repository names restricts the access of an installation token to specific repositories.
		Repositories []string `json:"repositories,omitempty"`

		// The permissions granted to the access token.
		// The permissions object includes the permission names and their access type.
		Permissions *InstallationPermissions `json:"permissions,omitempty"`
	}

	// InstallationPermissions lists the repository and organization permissions for an installation.
	InstallationPermissions struct {
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

	// AppsService provides access to Apps-related functions.
	AppsService interface {
		// Get returns a single App. Passing the empty string will get
		// the authenticated App.
		Get(ctx context.Context, appSlug string) (*App, *Response, error)

		// ListInstallations lists the installations that the current App has.
		ListInstallations(ctx context.Context, opts ListOptions) ([]*Installation, *Response, error)

		// GetInstallation returns the specified installation.
		GetInstallation(ctx context.Context, id int64) (*Installation, *Response, error)

		// ListUserInstallations lists installations that are accessible to the authenticated user.
		ListUserInstallations(ctx context.Context, opts ListOptions) ([]*Installation, *Response, error)

		// SuspendInstallation suspends the specified installation.
		SuspendInstallation(ctx context.Context, id int64) (*Response, error)

		// UnsuspendInstallation unsuspends the specified installation.
		UnsuspendInstallation(ctx context.Context, id int64) (*Response, error)

		// DeleteInstallation deletes the specified installation.
		DeleteInstallation(ctx context.Context, id int64) (*Response, error)

		// CreateInstallationToken creates a new installation token.
		CreateInstallationToken(ctx context.Context, id int64, opts *InstallationTokenOptions) (*InstallationToken, *Response, error)

		// FindOrganizationInstallation finds the organization's installation information.
		FindOrganizationInstallation(ctx context.Context, org string) (*Installation, *Response, error)

		// FindRepositoryInstallation finds the repository's installation information.
		FindRepositoryInstallation(ctx context.Context, owner, repo string) (*Installation, *Response, error)

		// FindRepositoryInstallationByID finds the repository's installation information by repository ID.
		FindRepositoryInstallationByID(ctx context.Context, id int64) (*Installation, *Response, error)

		// FindUserInstallation finds the user's installation information.
		FindUserInstallation(ctx context.Context, user string) (*Installation, *Response, error)
	}
)
