// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package main

import (
	"github.com/elastic/elastic-package/pkg/shell"
	hplugin "github.com/hashicorp/go-plugin"

	"github.com/elastic/ep-shell-plugins/command"
	"github.com/elastic/ep-shell-plugins/plugin"
)

func main() {
	const (
		ctxKeyPackages plugin.CtxKey = "Shell.Packages"
		ctxKeyDB       plugin.CtxKey = "Shell.DB"
	)

	p := plugin.NewPlugin()

	command.AddChangelogCmd(p, ctxKeyPackages)
	command.AddInitdbCmd(p, ctxKeyDB)
	command.AddWhereCmd(p, ctxKeyPackages, ctxKeyDB)
	command.AddWritefileCmd(p, ctxKeyPackages)
	command.AddRunscriptCmd(p, ctxKeyPackages)
	command.AddYqCmd(p, ctxKeyPackages)

	// pluginMap is the map of plugins we can dispense.
	pluginMap := map[string]hplugin.Plugin{
		"shell_plugin": &shell.ShellPlugin{Impl: p},
	}

	hplugin.Serve(&hplugin.ServeConfig{
		HandshakeConfig: shell.Handshake(),
		Plugins:         pluginMap,
	})
}
