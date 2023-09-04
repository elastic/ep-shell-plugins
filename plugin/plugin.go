// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package plugin

import (
	"context"
	"sync"

	"github.com/elastic/elastic-package/pkg/shell"
)

type CtxKey string

type Plugin struct {
	m    sync.Mutex
	cmds map[string]shell.Command
	ctx  context.Context
}

func NewPlugin() *Plugin {
	return &Plugin{
		cmds: map[string]shell.Command{},
		ctx:  context.Background(),
	}
}

func (p *Plugin) Commands() map[string]shell.Command {
	return p.cmds
}

func (p *Plugin) AddValueToCtx(k, v any) {
	p.m.Lock()
	p.ctx = context.WithValue(p.ctx, k, v)
	p.m.Unlock()
}

func (p *Plugin) GetValueFromCtx(k any) any {
	p.m.Lock()
	defer p.m.Unlock()
	return p.ctx.Value(k)
}

func (p *Plugin) RegisterCommand(cmd shell.Command) {
	p.m.Lock()
	p.cmds[cmd.Name()] = cmd
	p.m.Unlock()
}
