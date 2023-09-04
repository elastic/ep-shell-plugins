## Getting started

This repository implements plugin commands that can be used with `elastic-package shell` to interact with integrations.

## Development

Even though the project is "go-gettable", there is the [`Makefile`](./Makefile) present, which can be used to build,
install, format the source code among others. Some examples of the available targets are:

`make build` - build the source

`make clean` - delete binary and build folder

`make format` - format the Go code

`make check` - one-liner, used by CI to verify if source code is ready to be pushed to the repository

`make install` - build the tool source and move binary `$HOME/.elastic-package/shell_plugins`

`make gomod` - ensure go.mod and go.sum are up to date

`make licenser` - add the Elastic license header in the source code

To start developing, download and build the latest main of `plugins` binary:

```bash
git clone https://github.com/elastic/ep-shell-plugins.git
cd ep-shell-plugins
make install
```

When developing on Windows, please use the `core.autocrlf=input` or `core.autocrlf=false` option to avoid issues with CRLF line endings:
```bash
git clone --config core.autocrlf=input https://github.com/elastic/ep-shell-plugins.git
cd ep-shell-plugins
make build
```

This option can be also configured on existing clones with the following commands. Be aware that these commands
will remove uncommited changes.
```bash
git config core.autocrlf input
git rm --cached -r .
git reset --hard
```
