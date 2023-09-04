// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package command

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"

	"github.com/elastic/elastic-package/pkg/shell"

	"github.com/elastic/ep-shell-plugins/plugin"
)

var _ shell.Command = &writefileCmd{}

type writefileCmd struct {
	pkgsContextKey    plugin.CtxKey
	p                 *plugin.Plugin
	flags             *pflag.FlagSet
	name, usage, desc string
}

func AddWritefileCmd(p *plugin.Plugin, pkgsContextKey plugin.CtxKey) {
	flags := pflag.NewFlagSet("", pflag.ContinueOnError)
	flags.String("path", "", "Path to the file (relative to the package root).")
	flags.String("contents", "", "Contents of the file")
	cmd := &writefileCmd{
		pkgsContextKey: pkgsContextKey,
		p:              p,
		flags:          flags,
		name:           "write-file",
		usage:          "write-file --path path --contents contents",
		desc:           "Writes a file in each of the packages in context 'Shell.Packages'.",
	}
	p.RegisterCommand(cmd)
}

func (c *writefileCmd) Name() string  { return c.name }
func (c *writefileCmd) Usage() string { return c.usage }
func (c *writefileCmd) Desc() string  { return c.desc }

func (c *writefileCmd) Exec(wd string, args []string, _, _ io.Writer) error {
	packages, ok := c.p.GetValueFromCtx(c.pkgsContextKey).([]string)
	if !ok {
		return errors.New("no packages found in the context")
	}

	if err := c.flags.Parse(args); err != nil {
		return err
	}

	for _, pkg := range packages {
		packageRoot := filepath.Join(wd, pkg)
		// check if we are in packages folder
		if _, err := os.Stat(packageRoot); err != nil {
			// check if we are in integrations root folder
			packageRoot = filepath.Join(wd, "packages", pkg)
			if _, err := os.Stat(packageRoot); err != nil {
				return errors.New("you need to be in integrations root folder or in the packages folder")
			}
		}

		path, _ := c.flags.GetString("path")
		path = filepath.Join(packageRoot, path)

		contents, _ := c.flags.GetString("contents")

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		f, err := os.Create(path)
		if err != nil {
			return err
		}

		if _, err := f.WriteString(strings.ReplaceAll(contents, `\n`, "\n")); err != nil {
			f.Close()
			return err
		}

		f.Close()
	}
	return nil
}
