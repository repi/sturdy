package service

import (
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/workspace/service"
)

func Module(c *di.Container) {
	c.Register(New, new(service.Service))
}