//go:build !cloud
// +build !cloud

package module

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/downloads/graphql"
)

func Module(c *di.Container) {
	c.Import(graphql.Module)
}
