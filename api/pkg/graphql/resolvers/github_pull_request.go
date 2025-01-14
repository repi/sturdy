package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/github"
)

type GitHubPullRequestRootResolver interface {
	// Internal
	InternalPullRequest(*github.PullRequest) (GitHubPullRequestResolver, error)
	InternalByCodebaseIDAndHeadSHA(context.Context, codebases.ID, string) (GitHubPullRequestResolver, error)
	InternalGitHubPullRequestByWorkspaceID(ctx context.Context, args GitHubPullRequestArgs) (GitHubPullRequestResolver, error)

	// Mutations
	CreateOrUpdateGitHubPullRequest(ctx context.Context, args CreateOrUpdateGitHubPullRequestArgs) (GitHubPullRequestResolver, error)
	MergeGitHubPullRequest(ctx context.Context, args MergeGitHubPullRequestArgs) (GitHubPullRequestResolver, error)

	// Subscriptions
	UpdatedGitHubPullRequest(ctx context.Context, args UpdatedGitHubPullRequestArgs) (<-chan GitHubPullRequestResolver, error)
}

type CreateOrUpdateGitHubPullRequestArgs struct {
	Input CreateOrUpdateGitHubPullRequestInput
}

type CreateOrUpdateGitHubPullRequestInput struct {
	WorkspaceID graphql.ID

	// PatchIDs is deprecated and not used
	PatchIDs *[]string
}

type GitHubPullRequestArgs struct {
	WorkspaceID *graphql.ID
}

type UpdatedGitHubPullRequestArgs struct {
	WorkspaceID graphql.ID
}

type MergeGitHubPullRequestArgs struct {
	Input MergeGitHubPullRequestInput
}

type MergeGitHubPullRequestInput struct {
	WorkspaceID graphql.ID
}

type GitHubPullRequestResolver interface {
	ID() graphql.ID
	PullRequestNumber() int32
	Open() bool
	Merged() bool
	MergedAt() *int32
	Base() string
	Workspace(context.Context) (WorkspaceResolver, error)
	Statuses(context.Context) ([]GitHubPullRequestStatusResolver, error)
	State() (GitHubPullRequestState, error)
}

type GitHubPullRequestState string

const (
	GitHubPullRequestStateOpen    GitHubPullRequestState = "Open"
	GitHubPullRequestStateClosed  GitHubPullRequestState = "Closed"
	GitHubPullRequestStateMerging GitHubPullRequestState = "Merging"
	GitHubPullRequestStateMerged  GitHubPullRequestState = "Merged"
)
