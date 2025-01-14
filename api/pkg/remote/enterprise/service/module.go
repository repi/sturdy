package service

import (
	"getsturdy.com/api/pkg/di"
	remote_service "getsturdy.com/api/pkg/remote/service"
)

func Module(c *di.Container) {
	c.Register(New)
	c.Register(New, new(remote_service.Service))
}
