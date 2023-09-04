// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package command

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Masterminds/semver/v3"
	"github.com/spf13/pflag"

	"github.com/elastic/elastic-package/pkg/packages"
	"github.com/elastic/elastic-package/pkg/packages/changelog"
	"github.com/elastic/elastic-package/pkg/shell"

	"github.com/elastic/ep-shell-plugins/plugin"
)

var _ shell.Command = &changelogCmd{}

type changelogCmd struct {
	pkgsContextKey    plugin.CtxKey
	p                 *plugin.Plugin
	name, usage, desc string
	flags             *pflag.FlagSet
}

func AddChangelogCmd(p *plugin.Plugin, pkgsContextKey plugin.CtxKey) {
	flags := pflag.NewFlagSet("", pflag.ContinueOnError)
	flags.String("next", "", "changelog entry is added in the next `major`, `minor` or `patch` version")
	flags.String("description", "", "description for the changelog entry")
	flags.String("type", "", "type of change (bugfix, enhancement or breaking-change) for the changelog entry")
	flags.String("link", "", "link to the pull request or issue with more information about the changelog entry")
	cmd := &changelogCmd{
		pkgsContextKey: pkgsContextKey,
		p:              p,
		name:           "changelog",
		usage:          "changelog --next {major|minor|patch} --description desc --type {bugfix|enhancement|breaking-change} --link link",
		desc:           "Add an entry to the changelog file in each of the packages in context 'Shell.Packages'.",
		flags:          flags,
	}
	p.RegisterCommand(cmd)
}

func (c *changelogCmd) Name() string  { return c.name }
func (c *changelogCmd) Usage() string { return c.usage }
func (c *changelogCmd) Desc() string  { return c.desc }

func (c *changelogCmd) Exec(wd string, args []string, _, _ io.Writer) error {
	packages, ok := c.p.GetValueFromCtx(c.pkgsContextKey).([]string)
	if !ok {
		return errors.New("no packages found in the context")
	}

	if err := c.flags.Parse(args); err != nil {
		return err
	}

	for _, pkg := range packages {
		packageRoot := pkg
		// check if we are in packages folder
		if _, err := os.Stat(filepath.Join(wd, pkg)); err != nil {
			// check if we are in integrations root folder
			packageRoot = filepath.Join(wd, "packages", pkg)
			if _, err := os.Stat(packageRoot); err != nil {
				return errors.New("you need to be in integrations root folder or in the packages folder")
			}
		}
		if err := changelogAddCmdForRoot(packageRoot, c.flags, args); err != nil {
			return err
		}
	}
	return nil
}

func changelogAddCmdForRoot(packageRoot string, flags *pflag.FlagSet, args []string) error {
	nextMode, _ := flags.GetString("next")
	v, err := changelogCmdVersion(nextMode, packageRoot)
	if err != nil {
		return err
	}
	version := v.String()

	description, _ := flags.GetString("description")
	changeType, _ := flags.GetString("type")
	link, _ := flags.GetString("link")

	entry := changelog.Revision{
		Version: version,
		Changes: []changelog.Entry{
			{
				Description: description,
				Type:        changeType,
				Link:        link,
			},
		},
	}

	if err := patchChangelogFile(packageRoot, entry); err != nil {
		return err
	}

	if err := setManifestVersion(packageRoot, version); err != nil {
		return err
	}

	return nil
}

func changelogCmdVersion(nextMode, packageRoot string) (*semver.Version, error) {
	revisions, err := changelog.ReadChangelogFromPackageRoot(packageRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to read current changelog: %w", err)
	}
	if len(revisions) == 0 {
		return semver.MustParse("0.0.0"), nil
	}

	version, err := semver.NewVersion(revisions[0].Version)
	if err != nil {
		return nil, fmt.Errorf("invalid version in changelog %q: %w", revisions[0].Version, err)
	}

	switch nextMode {
	case "":
		break
	case "major":
		v := version.IncMajor()
		version = &v
	case "minor":
		v := version.IncMinor()
		version = &v
	case "patch":
		v := version.IncPatch()
		version = &v
	default:
		return nil, fmt.Errorf("invalid value for \"next\": %s", nextMode)
	}

	return version, nil
}

// patchChangelogFile looks for the proper place to add the new revision in the changelog,
// trying to conserve original format and comments.
func patchChangelogFile(packageRoot string, patch changelog.Revision) error {
	changelogPath := filepath.Join(packageRoot, changelog.PackageChangelogFile)
	d, err := os.ReadFile(changelogPath)
	if err != nil {
		return err
	}

	d, err = changelog.PatchYAML(d, patch)
	if err != nil {
		return err
	}

	return os.WriteFile(changelogPath, d, 0644)
}

func setManifestVersion(packageRoot string, version string) error {
	manifestPath := filepath.Join(packageRoot, packages.PackageManifestFile)
	d, err := os.ReadFile(manifestPath)
	if err != nil {
		return err
	}

	d, err = changelog.SetManifestVersion(d, version)
	if err != nil {
		return err
	}

	return os.WriteFile(manifestPath, d, 0644)
}
