// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scm

import (
	"context"
	"time"
)

type (
	// Repository represents a git repository.
	Repository struct {
		ID         string
		Namespace  string
		Name       string
		Perm       *Perm
		Branch     string
		Archived   bool
		Private    bool
		Visibility Visibility
		Clone      string
		CloneSSH   string
		Link       string
		Created    time.Time
		Updated    time.Time
	}

	// Perm represents a user's repository permissions.
	Perm struct {
		Pull  bool
		Push  bool
		Admin bool
	}

	// Hook represents a repository hook.
	Hook struct {
		ID         string
		Name       string
		Target     string
		Events     []string
		Active     bool
		SkipVerify bool
	}

	// HookInput provides the input fields required for
	// creating or updating repository webhooks.
	HookInput struct {
		Name       string
		Target     string
		Secret     string
		Events     HookEvents
		SkipVerify bool

		// NativeEvents are used to create hooks with
		// provider-specific event types that cannot be
		// abstracted or represented in HookEvents.
		NativeEvents []string
	}

	// HookEvents represents supported hook events.
	HookEvents struct {
		Branch             bool
		Deployment         bool
		Issue              bool
		IssueComment       bool
		Pipeline           bool
		PullRequest        bool
		PullRequestComment bool
		Push               bool
		ReviewComment      bool
		Tag                bool
	}

	// Status represents a commit status.
	Status struct {
		State  State
		Label  string
		Desc   string
		Target string

		// TODO(bradrydzewski) this field is only used
		// by Bitbucket which requires a user-defined
		// key (label), title and description. We need
		// to cleanup this abstraction.
		Title string
	}

	// StatusInput provides the input fields required for
	// creating or updating commit statuses.
	StatusInput struct {
		State  State
		Label  string
		Title  string
		Desc   string
		Target string
	}

	// DeployStatus represents a deployment status.
	DeployStatus struct {
		Number         int64
		State          State
		Desc           string
		Target         string
		Environment    string
		EnvironmentURL string
	}

	// RepositoryService provides access to repository resources.
	RepositoryService interface {
		// Find returns a repository by name.
		Find(context.Context, string) (*Repository, *Response, error)

		// FindHook returns a repository hook.
		FindHook(context.Context, string, string) (*Hook, *Response, error)

		// FindPerms returns repository permissions.
		FindPerms(context.Context, string) (*Perm, *Response, error)

		// List returns a list of repositories.
		List(context.Context, ListOptions) ([]*Repository, *Response, error)

		// ListV2 returns a list of repositories based on the searchTerm passed.
		ListV2(context.Context, RepoListOptions) ([]*Repository, *Response, error)

		// ListNamespace returns a list of repos in namespace
		ListNamespace(context.Context, string, ListOptions) ([]*Repository, *Response, error)

		// List2 returns a list of repositories .
		List2(context.Context, string, ListOptions) ([]*Repository, *Response, error)

		// ListRepoLanguages returns a list of repositories language with percentage.
		ListRepoLanguages(context.Context, string) (map[string]float64, *Response, error)

		// ListHooks returns a list or repository hooks.
		ListHooks(context.Context, string, ListOptions) ([]*Hook, *Response, error)

		// ListStatus returns a list of commit statuses.
		ListStatus(context.Context, string, string, ListOptions) ([]*Status, *Response, error)

		// CreateHook creates a new repository hook.
		CreateHook(context.Context, string, *HookInput) (*Hook, *Response, error)

		// CreateStatus creates a new commit status.
		CreateStatus(context.Context, string, string, *StatusInput) (*Status, *Response, error)

		// UpdateHook updates an existing repository hook.
		UpdateHook(context.Context, string, string, *HookInput) (*Hook, *Response, error)

		// DeleteHook deletes a repository hook.
		DeleteHook(context.Context, string, string) (*Response, error)
	}
)

// TODO(bradrydzewski): Add endpoint to get a repository deploy key
// TODO(bradrydzewski): Add endpoint to list repository deploy keys
// TODO(bradrydzewski): Add endpoint to create a repository deploy key
// TODO(bradrydzewski): Add endpoint to delete a repository deploy key
