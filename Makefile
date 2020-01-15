GOFILES := $(shell find . -name '*.go' ! -path './.go*')
SHELL   := /bin/bash

bump:
	@go get -u
	@go mod tidy

deps:
	@yay -S mkcert python-brotli python-fonttools

dev:
	@docker-compose rm -f
	@docker-compose up --build

diff-db:
	@diff --color --label live --label dev --strip-trailing-cr -su \
		<(ssh -p 1988 code-golf.io pg_dump -Os code_golf)          \
		<(docker-compose exec db pg_dump -OsU postgres code_golf)

fmt:
	@gofmt -s  -w $(GOFILES)
	@goimports -w $(GOFILES)

font:
	@pyftsubset ~/Downloads/fontawesome-pro-5.12.0-web/webfonts/fa-light-300.ttf \
		--flavor=woff2                                                           \
		--no-hinting                                                             \
		--output-file=assets/font.woff2                                          \
		--unicodes-file=font-subset.txt
