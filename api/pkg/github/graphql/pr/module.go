package pr

import "getsturdy.com/api/pkg/di"

func Module(c *di.Container) {
	c.Register(NewResolver)
}