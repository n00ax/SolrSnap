GOCMD=go

build-all:
	go build
build-docker:
	docker run -v $(pwd):/solr-snapshot-service -w /solr-snapshot-service 7d07942c16d5 make build-all