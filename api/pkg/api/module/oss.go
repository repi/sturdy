//go:build !enterprise && !cloud
// +build !enterprise,!cloud

package module

import (
	"getsturdy.com/api/pkg/api"
	"getsturdy.com/api/pkg/di"
	module_queue "getsturdy.com/api/pkg/queue/module"
	module_remote "getsturdy.com/api/pkg/remote/module"
)

func Module(c *di.Container) {
	common(c)
	c.Import(module_queue.Module)
	c.Import(module_remote.Module)

	c.Register(api.ProvideAPI, new(api.Starter))
}
