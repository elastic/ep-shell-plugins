VERSION_IMPORT_PATH = github.com/elastic/ep-shell-plugins/version
VERSION_COMMIT_HASH = `git describe --always --long --dirty`
VERSION_BUILD_TIME = `date +%s`
DEFAULT_VERSION_TAG ?=
VERSION_TAG = `(git describe --exact-match --tags 2>/dev/null || echo '$(DEFAULT_VERSION_TAG)') | tr -d '\n'`
VERSION_LDFLAGS = -X $(VERSION_IMPORT_PATH).CommitHash=$(VERSION_COMMIT_HASH) -X $(VERSION_IMPORT_PATH).BuildTime=$(VERSION_BUILD_TIME) -X $(VERSION_IMPORT_PATH).Tag=$(VERSION_TAG)

.PHONY: build

build:
	go build -ldflags "$(VERSION_LDFLAGS)" -o plugins

clean:
	rm -rf build
	rm -f plugins

format:
	go run golang.org/x/tools/cmd/goimports -local github.com/elastic/ep-shell-plugins/ -w .

install: build
	mkdir -p ${HOME}/.elastic-package/shell_plugins
	mv plugins ${HOME}/.elastic-package/shell_plugins

lint:
	go run honnef.co/go/tools/cmd/staticcheck ./...

licenser:
	go run github.com/elastic/go-licenser -license Elastic

gomod:
	go mod tidy

test-go:
	# -count=1 is included to invalidate the test cache. This way, if you run "make test-go" multiple times
	# you will get fresh test results each time. For instance, changing the source of mocked packages
	# does not invalidate the cache so having the -count=1 to invalidate the test cache is useful.
	go test -count 1 ./...

test: test-go

check-git-clean:
	git update-index --really-refresh
	git diff-index --quiet HEAD

check: check-static test check-git-clean

check-static: build format lint licenser gomod check-git-clean
