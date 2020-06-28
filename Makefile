GOFILES  := $(shell find . -name '*.go' ! -path './.go*')
POSTGRES := postgres:11.8-alpine
SHELL    := /bin/bash

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

cert:
	@mkcert -install localhost
	@chmod +r localhost-key.pem

.PHONY: db
db:
	@ssh -t rancher@code.golf docker run -it --rm \
	    --env-file /etc/code-golf.env $(POSTGRES) psql

db-admin:
	@ssh -t rancher@code.golf docker run -it --rm \
	    --env-file /etc/code-golf.env $(POSTGRES) psql -WU doadmin

db-diff:
	@diff --color --label live --label dev --strip-trailing-cr -su \
	    <(ssh rancher@code.golf "docker run --rm                   \
	    --env-file /etc/code-golf.env $(POSTGRES) pg_dump -Os")    \
	    <(docker-compose exec db pg_dump -OsU postgres code-golf)

db-dump:
	@rm db/*.gz

	@ssh rancher@code.golf "docker run --env-file /etc/code-golf.env \
	    --rm $(POSTGRES) sh -c 'pg_dump -a | gzip -9'"               \
	    > db/code-golf-`date +%Y-%m-%d`.sql.gz

	@cp db/*.gz ~/Dropbox/code-golf/

deps:
	@yay -S mkcert python-brotli python-fonttools

dev:
	@docker-compose rm -f
	@docker-compose up --build

fmt:
	@gofmt -s  -w $(GOFILES)
	@goimports -w $(GOFILES)

font:
	@docker build -t code-golf-font -f Dockerfile.font .
	@id=`docker create code-golf-font`;                                                 \
	    docker cp "$$id:twemoji-colr/build/Twemoji Mozilla.woff2" assets/twemoji.woff2; \
	    docker rm $$id

legacy-font:
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

	@docker run --rm -v $(CURDIR):/app -w /app golangci/golangci-lint:v1.27.0 golangci-lint run

live:
	@./build-assets

	@docker build --pull -t codegolf/code-golf .

	@docker push codegolf/code-golf

	@ssh rancher@code.golf "              \
	    docker pull codegolf/code-golf && \
	    docker stop code-golf;            \
	    docker rm code-golf;              \
	    docker run                        \
	    --cap-add      CAP_KILL           \
	    --cap-add      CAP_SETGID         \
	    --cap-add      CAP_SETUID         \
	    --cap-add      CAP_SYS_ADMIN      \
	    --cap-drop     ALL                \
	    --detach                          \
	    --env-file     /etc/code-golf.env \
	    --init                            \
	    --name         code-golf          \
	    --publish       80:1080           \
	    --publish      443:1443           \
	    --read-only                       \
	    --restart      always             \
	    --security-opt seccomp:unconfined \
	    --volume       certs:/certs       \
	    codegolf/code-golf"

logs:
	@ssh rancher@code.golf docker logs -f code-golf

test:
# FIXME Stub out assets if it doesn't yet exist.
ifeq ($(wildcard routes/assets.go),)
	$(file > routes/assets.go, $(STUB))
endif

	@go test ./...
