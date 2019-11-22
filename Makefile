GOFILES := $(shell find -name '*.go' ! -path './.go*')
SHELL   := /bin/bash

bump:
	@go get -u
	@go mod tidy

deps:
	@yay -S mkcert python-brotli python-fonttools

dev:
	@docker-compose rm -f
	@docker-compose up --build

dump-db:
	@ssh kino pg_dump -Os code_golf > db/0.schema.sql
	@perl -pi -e 'chomp if eof' db/0.schema.sql

fmt:
	@gofmt -s  -w $(GOFILES)
	@goimports -w $(GOFILES)

font:
	@pyftsubset ~/Downloads/fontawesome-pro-5.11.2-web/webfonts/fa-light-300.ttf \
		--flavor=woff2                                                           \
		--no-hinting                                                             \
		--output-file=assets/font.woff2                                          \
		--unicodes-file=font-subset.txt
