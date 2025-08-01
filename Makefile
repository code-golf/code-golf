DATE    := $(shell date +%Y-%m-%d)
GOFILES := $(shell find . -name '*.go' ! -path './.go*')
SHELL   := /bin/bash

export COMPOSE_PATH_SEPARATOR=:
export COMPOSE_FILE=docker/core.yml:docker/dev.yml

bench:
	@go test -bench B -benchmem ./...

bump:
	@go get -u
	@go mod tidy -compat=1.24
	@npm upgrade

cert:
	@mkcert -install localhost
	@chmod +r localhost-key.pem

.PHONY: db
db:
	@ssh -t root@code.golf sudo -iu postgres psql code-golf

db-dev:
	@docker compose exec db psql -U postgres code-golf

db-diff:
	@diff --color --label live --label dev --strip-trailing-cr -su   \
	    <(ssh root@code.golf sudo -iu postgres pg_dump -Os code-golf \
	    | sed -E 's/ \(Debian .+//')                                 \
	    <(docker compose exec -T db pg_dump -OsU postgres code-golf)

db-dump:
	@rm -f sql/*.gz

	@ssh root@code.golf sudo -iu postgres pg_dump -aZ9 code-golf \
	    > sql/code-golf-$(DATE).sql.gz

	@zcat sql/*.gz | zstd -fqo ~/Dropbox/code-golf/code-golf-$(DATE).sql.zst

dev:
	@touch docker/.env
	@docker compose rm -f
	@docker compose up --build

# e2e-iterate is useful when you have made a small change to test code only
# and want to re-run. Note that logs are not automatically shown when tests
# fail, because they make it harder to see test results and this target isn't
# used by CI.
e2e-iterate: export COMPOSE_FILE=docker/core.yml:docker/e2e.yml
e2e-iterate: export COMPOSE_PROJECT_NAME=code-golf-e2e
e2e-iterate:
	@docker compose run --rm e2e

e2e: export COMPOSE_FILE=docker/core.yml:docker/e2e.yml
e2e: export COMPOSE_PROJECT_NAME=code-golf-e2e
e2e:
# TODO Pass arguments to run specific tests.
	@./esbuild
	@touch docker/.env
	@docker compose rm -fsv &>/dev/null
	@docker compose build --pull -q
	@docker compose run --rm e2e || (docker compose logs; false)
	@docker compose rm -fsv &>/dev/null

fmt:
	@gofmt -s  -w $(GOFILES)
	@goimports -w $(GOFILES)

font:
	@docker build -t code-golf-font -f docker/font.Dockerfile docker
	@id=`docker create code-golf-font`;                                                \
	    docker cp "$$id:twemoji-colr/build/Twemoji Mozilla.woff2" fonts/twemoji.woff2; \
	    docker rm $$id

mathjax-fonts:
	@rm -rf public/mathjax-fonts
	@cp -r node_modules/mathjax/es5/output/chtml/fonts/woff-v2 public/mathjax-fonts

lint:
	@docker run --rm -v $(CURDIR):/app -w /app \
	    golangci/golangci-lint:v2.3.0 golangci-lint run

	@node_modules/.bin/eslint js

# Calls "make logs" at the end to make sure we didn't break the site.
live:
	@docker buildx build --pull --push \
	    --file docker/live.Dockerfile --tag codegolf/code-golf .

	@ssh root@code.golf "                                        \
	    docker pull codegolf/code-golf &&                        \
	    docker stop code-golf;                                   \
	    docker rm code-golf;                                     \
	    docker run                                               \
	        --detach                                             \
	        --env-file   /etc/code-golf.env                      \
	        --init                                               \
	        --name       code-golf                               \
	        --network    caddy                                   \
	        --pids-limit 1024                                    \
	        --privileged                                         \
	        --read-only                                          \
	        --restart    always                                  \
	        --volume     /var/run/postgresql:/var/run/postgresql \
	    codegolf/code-golf &&                                    \
	    docker system prune -f"

	@make logs

logs:
	@ssh root@code.golf docker logs --tail 5 -f code-golf

svgo:
	@node_modules/.bin/svgo -f svg

test:
	@go test ./...

.PHONY: xt
xt:
	@prove6 xt
