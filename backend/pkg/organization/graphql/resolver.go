package graphql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gosimple/slug"
	"github.com/graph-gophers/graphql-go"

	"mash/pkg/auth"
	service_auth "mash/pkg/auth/service"
	service_codebase "mash/pkg/codebase/service"
	gqlerrors "mash/pkg/graphql/errors"
	"mash/pkg/graphql/resolvers"
	"mash/pkg/organization"
	service_organization "mash/pkg/organization/service"
	service_user "mash/pkg/user/service"
)

type organizationRootResolver struct {
	service         *service_organization.Service
	authService     *service_auth.Service
	userService     *service_user.Service
	codebaseService *service_codebase.Service

	authorRootResolver    resolvers.AuthorRootResolver
	licensesRootResolver  resolvers.LicenseRootResolver
	codebasesRootResolver resolvers.CodebaseRootResolver
}

func New(
	service *service_organization.Service,
	authService *service_auth.Service,
	userService *service_user.Service,
	codebaseService *service_codebase.Service,

	authorRootResolver resolvers.AuthorRootResolver,
	licensesRootResolver resolvers.LicenseRootResolver,
	codebasesRootResolver resolvers.CodebaseRootResolver,
) resolvers.OrganizationRootResolver {
	return &organizationRootResolver{
		service:         service,
		authService:     authService,
		userService:     userService,
		codebaseService: codebaseService,

		authorRootResolver:    authorRootResolver,
		licensesRootResolver:  licensesRootResolver,
		codebasesRootResolver: codebasesRootResolver,
	}
}

func (r *organizationRootResolver) Organizations(ctx context.Context) ([]resolvers.OrganizationResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	orgs, err := r.service.ListByUserID(ctx, userID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	var res []resolvers.OrganizationResolver

	for _, org := range orgs {
		res = append(res, &organizationResolver{
			root: r,
			org:  org,
		})
	}

	return res, nil
}

func (r *organizationRootResolver) Organization(ctx context.Context, args resolvers.OrganizationArgs) (resolvers.OrganizationResolver, error) {

	var org *organization.Organization

	if args.ID != nil {
		var err error
		org, err = r.service.GetByID(ctx, string(*args.ID))
		if err != nil {
			return nil, gqlerrors.Error(err)
		}
	} else if args.ShortID != nil {
		var err error
		org, err = r.service.GetByShortID(ctx, organization.ShortOrganizationID(*args.ShortID))
		if err != nil {
			return nil, gqlerrors.Error(err)
		}
	}

	if err := r.authService.CanRead(ctx, org); err != nil {
		return nil, gqlerrors.Error(err)
	}

	return &organizationResolver{org: org, root: r}, nil
}

func (r *organizationRootResolver) CreateOrganization(ctx context.Context, args resolvers.CreateOrganizationArgs) (resolvers.OrganizationResolver, error) {
	org, err := r.service.Create(ctx, args.Input.Name)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	return &organizationResolver{root: r, org: org}, nil
}

func (r *organizationRootResolver) AddUserToOrganization(ctx context.Context, args resolvers.AddUserToOrganizationArgs) (resolvers.OrganizationResolver, error) {
	org, err := r.service.GetByID(ctx, string(args.Input.OrganizationID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.authService.CanWrite(ctx, org); err != nil {
		return nil, gqlerrors.Error(err)
	}

	user, err := r.userService.GetByEmail(ctx, args.Input.Email)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if _, err := r.service.AddMember(ctx, org.ID, user.ID); err != nil {
		return nil, gqlerrors.Error(err)
	}

	return &organizationResolver{root: r, org: org}, nil
}

type organizationResolver struct {
	root *organizationRootResolver
	org  *organization.Organization
}

func (r *organizationResolver) ID() graphql.ID {
	return graphql.ID(r.org.ID)
}

func (r *organizationResolver) ShortID() graphql.ID {
	return graphql.ID(slug.Make(r.org.Name) + "-" + string(r.org.ShortID))
}

func (r *organizationResolver) Name() string {
	return r.org.Name
}

func (r *organizationResolver) Members(ctx context.Context) ([]resolvers.AuthorResolver, error) {
	members, err := r.root.service.Members(ctx, r.org.ID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	var res []resolvers.AuthorResolver

	for _, m := range members {
		author, err := r.root.authorRootResolver.Author(ctx, graphql.ID(m.UserID))
		switch {
		case err == nil:
			res = append(res, author)
		case errors.Is(err, sql.ErrNoRows):
			// skip
		case err != nil:
			return nil, gqlerrors.Error(err)
		}
	}

	return res, nil
}

func (r *organizationResolver) Codebases(ctx context.Context) ([]resolvers.CodebaseResolver, error) {
	codebases, err := r.root.codebaseService.ListByOrganization(ctx, r.org.ID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	var res []resolvers.CodebaseResolver

	for _, cb := range codebases {
		id := graphql.ID(cb.ID)
		resolver, err := r.root.codebasesRootResolver.Codebase(ctx, resolvers.CodebaseArgs{ID: &id})
		if err != nil {
			return nil, gqlerrors.Error(err)
		}
		res = append(res, resolver)
	}

	return res, nil
}

func (r *organizationResolver) LicenseSubscriptions(ctx context.Context) ([]resolvers.LicenseResolver, error) {
	return r.root.licensesRootResolver.InternalListForOrganization(ctx, r.org.ID)
}