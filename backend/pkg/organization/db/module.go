package db

import "mash/pkg/di"

func Module(c *di.Container) {
	c.Register(NewMember)
	c.Register(New)
}