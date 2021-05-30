GOFILES  := $(shell find . -name '*.go' ! -path './.go*')
POSTGRES := postgres:13.3-alpine
SHELL    := /bin/bash

export COMPOSE_PATH_SEPARATOR=:
export COMPOSE_FILE=docker/core.yml:docker/ports.yml

define STUB
package routes

import "net/http"

const holeJsPath       = ""
const twemojiWoff2Path = ""

func Asset(w http.ResponseWriter, r *http.Request) {}
endef

bench:
# FIXME Stub out assets if it doesn't yet exist.
ifeq ($(wildcard routes/assets.go),)
	$(file > routes/assets.go, $(STUB))
endif

	@go test -bench B -benchmem ./...

bump:
	@go get -u
	@go mod tidy
	@npm upgrade
	@make editor

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

db-dev:
	@docker-compose exec db psql -U postgres code-golf

db-diff:
	@diff --color --label live --label dev --strip-trailing-cr -su \
	    <(ssh rancher@code.golf "docker run --rm                   \
	    --env-file /etc/code-golf.env $(POSTGRES) pg_dump -Os")    \
	    <(docker-compose exec db pg_dump -OsU postgres code-golf)

db-dump:
	@rm -f db/*.gz

	@ssh rancher@code.golf "docker run --env-file /etc/code-golf.env \
	    --rm $(POSTGRES) sh -c 'pg_dump -a | gzip -9'"               \
	    > db/code-golf-`date +%Y-%m-%d`.sql.gz

	@cp db/*.gz ~/Dropbox/code-golf/

deps:
	@yay -S mkcert python-brotli python-fonttools

dev:
	@touch docker/.env
	@docker-compose rm -f
	@docker-compose up --build

e2e: export COMPOSE_FILE         = docker/core.yml:docker/e2e.yml
e2e: export COMPOSE_PROJECT_NAME = code-golf-e2e
e2e:
# TODO Pass arguments to run specific tests.
	@touch docker/.env
	@docker-compose rm -fs
	@docker-compose pull -q
	@docker-compose build -q
	@docker-compose run e2e || (docker-compose logs; false)
	@docker-compose rm -fs

editor:
	@node_modules/.bin/esbuild --bundle --format=esm --outdir=public editor.js

fmt:
	@gofmt -s  -w $(GOFILES)
	@goimports -w $(GOFILES)

font:
	@docker build -t code-golf-font -f docker/font.Dockerfile docker
	@id=`docker create code-golf-font`;                                                 \
	    docker cp "$$id:twemoji-colr/build/Twemoji Mozilla.woff2" assets/twemoji.woff2; \
	    docker rm $$id

lint:
# FIXME Stub out assets if it doesn't yet exist.
ifeq ($(wildcard routes/assets.go),)
	$(file > routes/assets.go, $(STUB))
endif

	@docker run --rm -v $(CURDIR):/app -w /app \
	    golangci/golangci-lint:v1.39.0 golangci-lint run

live:
	@./build-assets

	@docker build --pull -f docker/live.Dockerfile -t codegolf/code-golf .

	@docker push codegolf/code-golf

	@ssh rancher@code.golf "              \
	    docker pull codegolf/code-golf && \
	    docker stop code-golf;            \
	    docker rm code-golf;              \
	    docker run                        \
	    --detach                          \
	    --env-file     /etc/code-golf.env \
	    --init                            \
	    --name         code-golf          \
	    --pids-limit   1024               \
	    --privileged                      \
	    --publish       80:1080           \
	    --publish      443:1443           \
	    --read-only                       \
	    --restart      always             \
	    --ulimit       core=-1            \
	    --volume       certs:/certs       \
	    codegolf/code-golf &&             \
	    docker system prune -f"

logs:
	@ssh rancher@code.golf docker logs --tail 5 -f code-golf

test:
# FIXME Stub out assets if it doesn't yet exist.
ifeq ($(wildcard routes/assets.go),)
	$(file > routes/assets.go, $(STUB))
endif

	@go test ./...
