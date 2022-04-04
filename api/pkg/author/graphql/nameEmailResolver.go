package graphql

import (
	"getsturdy.com/api/pkg/graphql/resolvers"
	"github.com/graph-gophers/graphql-go"
)

// authorNameEmailResolver is used for "imported" git users that doesn't have a Sturdy account
type authorNameEmailResolver struct {
	name, email string
}

func (r *authorNameEmailResolver) Status() (resolvers.UserStatus, error) {
	return resolvers.UserStatusShadow, nil
}

func (r *authorNameEmailResolver) ID() graphql.ID {
	return graphql.ID(r.email + "//" + r.name)
}

func (r *authorNameEmailResolver) Name() string {
	return r.name
}

func (r *authorNameEmailResolver) AvatarUrl() *string {
	return nil
}

func (r *authorNameEmailResolver) Email() string {
	return r.email
}
