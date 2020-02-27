GOFILES := $(shell find . -name '*.go' ! -path './.go*')
SHELL   := /bin/bash

define STUB
package routes

import "net/http"

const (
	commonCssPath = ""
	holeCssPath   = ""
	holeJsPath    = ""
)

func Asset(w http.ResponseWriter, r *http.Request) {}
endef

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

lint:
# FIXME Stub out assets if it doesn't yet exist.
ifeq ($(wildcard routes/assets.go),)
	$(file > routes/assets.go, $(STUB))
endif

	@docker run --rm -v $(CURDIR):/app -w /app golangci/golangci-lint:v1.23.7 golangci-lint run

test:
# FIXME Stub out assets if it doesn't yet exist.
ifeq ($(wildcard routes/assets.go),)
	$(file > routes/assets.go, $(STUB))
endif

	@go test ./...
