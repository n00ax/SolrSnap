GOCMD=go
GIT_COMMIT=$(git rev-list -1 HEAD)

build-all:
	$(GOCMD) build -ldflags "-X config.variables.GitCommit=\"New\"" -o solrsnap .