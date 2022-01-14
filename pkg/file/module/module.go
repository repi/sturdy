package module

import (
	"mash/pkg/di"
	"mash/pkg/file/graphql"
)

func Module(c *di.Container) {
	c.Import(graphql.Module)
}
