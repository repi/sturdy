package graphql

import (
	"context"
	"fmt"

	"github.com/graph-gophers/graphql-go"

	service_auth "getsturdy.com/api/pkg/auth/service"
	service_codebase "getsturdy.com/api/pkg/codebase/service"
	gqlerror "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/remote/enterprise/service"
	service_workspace "getsturdy.com/api/pkg/workspaces/service"
)

type remoteRootResolver struct {
	service          *service.Service
	workspaceService service_workspace.Service
	authService      *service_auth.Service
	codebaseService  *service_codebase.Service

	workspaceRootResolver resolvers.WorkspaceRootResolver
	codebaseRootResolver  resolvers.CodebaseRootResolver
}

func New(
	service *service.Service,
	workspaceService service_workspace.Service,
	authService *service_auth.Service,
	codebaseService *service_codebase.Service,

	workspaceRootResolver resolvers.WorkspaceRootResolver,
	codebaseRootResolver resolvers.CodebaseRootResolver,
) resolvers.RemoteRootResolver {
	return &remoteRootResolver{
		service:          service,
		workspaceService: workspaceService,
		authService:      authService,
		codebaseService:  codebaseService,

		workspaceRootResolver: workspaceRootResolver,
		codebaseRootResolver:  codebaseRootResolver,
	}
}

func (r *remoteRootResolver) CreateCodebaseRemote(ctx context.Context, args resolvers.CreateCodebaseRemoteArgs) (resolvers.CodebaseResolver, error) {
	cb, err := r.codebaseService.GetByID(ctx, args.Input.CodebaseID)
	if err != nil {
		return nil, gqlerror.Error(err)
	}

	if err := r.authService.CanWrite(ctx, cb); err != nil {
		return nil, gqlerror.Error(err)
	}

	err = r.service.SetRemote(
		ctx,
		args.Input.CodebaseID,
		args.Input.Name,
		args.Input.Url,
		args.Input.BasicAuthUsername,
		args.Input.BasicAuthPassword,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add remote: %w", err)
	}

	id := graphql.ID(args.Input.CodebaseID)
	return r.codebaseRootResolver.Codebase(ctx, resolvers.CodebaseArgs{ID: &id})
}

func (r *remoteRootResolver) PushWorkspaceToRemote(ctx context.Context, args resolvers.PushWorkspaceToRemoteArgs) (resolvers.WorkspaceResolver, error) {
	ws, err := r.workspaceService.GetByID(ctx, string(args.Input.WorkspaceID))
	if err != nil {
		return nil, gqlerror.Error(fmt.Errorf("could not get workspace: %w", err))
	}

	if err := r.authService.CanWrite(ctx, ws); err != nil {
		return nil, gqlerror.Error(err)
	}

	if err := r.service.Push(ctx, ws.CodebaseID, ws.ID); err != nil {
		return nil, gqlerror.Error(fmt.Errorf("failed to push workspace: %w", err))
	}

	return r.workspaceRootResolver.Workspace(ctx, resolvers.WorkspaceArgs{ID: args.Input.WorkspaceID})
}

func (r *remoteRootResolver) RemoteFetchToTrunk(ctx context.Context, args resolvers.PushWorkspaceToRemoteArgs) (resolvers.WorkspaceResolver, error) {
	ws, err := r.workspaceService.GetByID(ctx, string(args.Input.WorkspaceID))
	if err != nil {
		return nil, gqlerror.Error(fmt.Errorf("could not get workspace: %w", err))
	}

	if err := r.authService.CanWrite(ctx, ws); err != nil {
		return nil, gqlerror.Error(err)
	}

	if err := r.service.Push(ctx, ws.CodebaseID, ws.ID); err != nil {
		return nil, gqlerror.Error(fmt.Errorf("failed to push workspace: %w", err))
	}

	return r.workspaceRootResolver.Workspace(ctx, resolvers.WorkspaceArgs{ID: args.Input.WorkspaceID})
}

func (r *remoteRootResolver) FetchRemoteToTrunk(ctx context.Context, args resolvers.FetchRemoteToTrunkArgs) (resolvers.CodebaseResolver, error) {
	cb, err := r.codebaseService.GetByID(ctx, string(args.Input.CodebaseID))
	if err != nil {
		return nil, gqlerror.Error(fmt.Errorf("could not get codebase: %w", err))
	}

	if err := r.authService.CanWrite(ctx, cb); err != nil {
		return nil, gqlerror.Error(err)
	}

	if err := r.service.Pull(ctx, cb.ID); err != nil {
		return nil, gqlerror.Error(fmt.Errorf("failed to fetch codebase: %w", err))
	}

	id := graphql.ID(args.Input.CodebaseID)
	return r.codebaseRootResolver.Codebase(ctx, resolvers.CodebaseArgs{ID: &id})
}
