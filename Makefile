GOFILES := $(shell find -name '*.go' ! -path './.go*')
SHELL   := /bin/bash

bump:
	@go get -u
	@go mod tidy

dev:
	@docker-compose rm -f
	@docker-compose up --build

dump-db:
	@ssh kino pg_dump -Os code_golf > db/0.schema.sql
	@perl -pi -e 'chomp if eof' db/0.schema.sql

fmt:
	@gofmt -s  -w $(GOFILES)
	@goimports -w $(GOFILES)
