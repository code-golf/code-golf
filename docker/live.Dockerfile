FROM golang:1.25.4-alpine3.22

ENV CGO_ENABLED=0 GOAMD64=v4 GOEXPERIMENT=jsonv2 GOPATH=

RUN apk add --no-cache brotli build-base linux-headers npm tzdata zopfli

# Fetch modules.
COPY go.mod go.sum ./

RUN go mod download

# Build assets.
COPY node_modules ./node_modules
COPY fonts        ./fonts
COPY css          ./css
COPY js           ./js
COPY svg          ./svg

COPY esbuild package-lock.json package.json ./

RUN ./esbuild \
 && find dist \( -name '*.css' -or  -name '*.js' -or -name '*.map' -or -name '*.svg' \) \
  | xargs -i -n1 -P`nproc` sh -c 'brotli {} && zopfli {}'

# Build lang hasher.
COPY cmd/hash-langs ./cmd/hash-langs

RUN go build -ldflags -s -trimpath ./cmd/hash-langs

# Build website.
COPY . ./

RUN go build -ldflags -s -trimpath \
 && gcc -Wall -Werror -Wextra -o /usr/bin/run-lang -s -static run-lang.c

FROM scratch

COPY --from=codegolf/lang-swift        / /langs/swift/rootfs/
COPY --from=codegolf/lang-rust         / /langs/rust/rootfs/
COPY --from=codegolf/lang-julia        / /langs/julia/rootfs/
COPY --from=codegolf/lang-haskell      / /langs/haskell/rootfs/
COPY --from=codegolf/lang-go           / /langs/go/rootfs/
COPY --from=codegolf/lang-odin         / /langs/odin/rootfs/
COPY --from=codegolf/lang-crystal      / /langs/crystal/rootfs/
COPY --from=codegolf/lang-zig          / /langs/zig/rootfs/
COPY --from=codegolf/lang-factor       / /langs/factor/rootfs/
COPY --from=codegolf/lang-kotlin       / /langs/kotlin/rootfs/
COPY --from=codegolf/lang-powershell   / /langs/powershell/rootfs/
COPY --from=codegolf/lang-assembly     / /langs/assembly/rootfs/
COPY --from=codegolf/lang-dart         / /langs/dart/rootfs/
COPY --from=codegolf/lang-f-sharp      / /langs/f-sharp/rootfs/
COPY --from=codegolf/lang-cpp          / /langs/cpp/rootfs/
COPY --from=codegolf/lang-c-sharp      / /langs/c-sharp/rootfs/
COPY --from=codegolf/lang-vala         / /langs/vala/rootfs/
COPY --from=codegolf/lang-d            / /langs/d/rootfs/
COPY --from=codegolf/lang-ocaml        / /langs/ocaml/rootfs/
COPY --from=codegolf/lang-scala        / /langs/scala/rootfs/
COPY --from=codegolf/lang-coconut      / /langs/coconut/rootfs/
COPY --from=codegolf/lang-civet        / /langs/civet/rootfs/
COPY --from=codegolf/lang-coffeescript / /langs/coffeescript/rootfs/
COPY --from=codegolf/lang-stax         / /langs/stax/rootfs/
COPY --from=codegolf/lang-vyxal        / /langs/vyxal/rootfs/
COPY --from=codegolf/lang-raku         / /langs/raku/rootfs/
COPY --from=codegolf/lang-gleam        / /langs/gleam/rootfs/
COPY --from=codegolf/lang-basic        / /langs/basic/rootfs/
COPY --from=codegolf/lang-fortran      / /langs/fortran/rootfs/
COPY --from=codegolf/lang-groovy       / /langs/groovy/rootfs/
COPY --from=codegolf/lang-clojure      / /langs/clojure/rootfs/
COPY --from=codegolf/lang-cjam         / /langs/cjam/rootfs/
COPY --from=codegolf/lang-elixir       / /langs/elixir/rootfs/
COPY --from=codegolf/lang-java         / /langs/java/rootfs/
COPY --from=codegolf/lang-erlang       / /langs/erlang/rootfs/
COPY --from=codegolf/lang-prolog       / /langs/prolog/rootfs/
COPY --from=codegolf/lang-javascript   / /langs/javascript/rootfs/
COPY --from=codegolf/lang-haxe         / /langs/haxe/rootfs/
COPY --from=codegolf/lang-hy           / /langs/hy/rootfs/
COPY --from=codegolf/lang-iogii        / /langs/iogii/rootfs/
COPY --from=codegolf/lang-v            / /langs/v/rootfs/
COPY --from=codegolf/lang-golfscript   / /langs/golfscript/rootfs/
COPY --from=codegolf/lang-ruby         / /langs/ruby/rootfs/
COPY --from=codegolf/lang-common-lisp  / /langs/common-lisp/rootfs/
COPY --from=codegolf/lang-r            / /langs/r/rootfs/
COPY --from=codegolf/lang-racket       / /langs/racket/rootfs/
COPY --from=codegolf/lang-05ab1e       / /langs/05ab1e/rootfs/
COPY --from=codegolf/lang-pascal       / /langs/pascal/rootfs/
COPY --from=codegolf/lang-python       / /langs/python/rootfs/
COPY --from=codegolf/lang-viml         / /langs/viml/rootfs/
COPY --from=codegolf/lang-uiua         / /langs/uiua/rootfs/
COPY --from=codegolf/lang-qore         / /langs/qore/rootfs/
COPY --from=codegolf/lang-nim          / /langs/nim/rootfs/
COPY --from=codegolf/lang-apl          / /langs/apl/rootfs/
COPY --from=codegolf/lang-harbour      / /langs/harbour/rootfs/
COPY --from=codegolf/lang-picat        / /langs/picat/rootfs/
COPY --from=codegolf/lang-j            / /langs/j/rootfs/
COPY --from=codegolf/lang-hare         / /langs/hare/rootfs/
COPY --from=codegolf/lang-php          / /langs/php/rootfs/
COPY --from=codegolf/lang-tex          / /langs/tex/rootfs/
COPY --from=codegolf/lang-hexagony     / /langs/hexagony/rootfs/
COPY --from=codegolf/lang-scheme       / /langs/scheme/rootfs/
COPY --from=codegolf/lang-egel         / /langs/egel/rootfs/
COPY --from=codegolf/lang-tcl          / /langs/tcl/rootfs/
COPY --from=codegolf/lang-perl         / /langs/perl/rootfs/
COPY --from=codegolf/lang-arturo       / /langs/arturo/rootfs/
COPY --from=codegolf/lang-luau         / /langs/luau/rootfs/
COPY --from=codegolf/lang-cobol        / /langs/cobol/rootfs/
COPY --from=codegolf/lang-squirrel     / /langs/squirrel/rootfs/
COPY --from=codegolf/lang-rockstar     / /langs/rockstar/rootfs/
COPY --from=codegolf/lang-hush         / /langs/hush/rootfs/
COPY --from=codegolf/lang-awk          / /langs/awk/rootfs/
COPY --from=codegolf/lang-c            / /langs/c/rootfs/
COPY --from=codegolf/lang-umka         / /langs/umka/rootfs/
COPY --from=codegolf/lang-bash         / /langs/bash/rootfs/
COPY --from=codegolf/lang-bqn          / /langs/bqn/rootfs/
COPY --from=codegolf/lang-algol-68     / /langs/algol-68/rootfs/
COPY --from=codegolf/lang-rebol        / /langs/rebol/rootfs/
COPY --from=codegolf/lang-sql          / /langs/sql/rootfs/
COPY --from=codegolf/lang-rexx         / /langs/rexx/rootfs/
COPY --from=codegolf/lang-berry        / /langs/berry/rootfs/
COPY --from=codegolf/lang-jq           / /langs/jq/rootfs/
COPY --from=codegolf/lang-janet        / /langs/janet/rootfs/
COPY --from=codegolf/lang-knight       / /langs/knight/rootfs/
COPY --from=codegolf/lang-fennel       / /langs/fennel/rootfs/
COPY --from=codegolf/lang-k            / /langs/k/rootfs/
COPY --from=codegolf/lang-forth        / /langs/forth/rootfs/
COPY --from=codegolf/lang-fish         / /langs/fish/rootfs/
COPY --from=codegolf/lang-wren         / /langs/wren/rootfs/
COPY --from=codegolf/lang-lua          / /langs/lua/rootfs/
COPY --from=codegolf/lang-befunge      / /langs/befunge/rootfs/
COPY --from=codegolf/lang-sed          / /langs/sed/rootfs/
COPY --from=codegolf/lang-brainfuck    / /langs/brainfuck/rootfs/

COPY --from=0 /go/hash-langs /

RUN ["/hash-langs"]

COPY --from=0 /go/code-golf                      /
COPY --from=0 /go/dist                           /dist/
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY          /public                            /public/
COPY --from=0 /usr/bin/run-lang                  /usr/bin/
COPY --from=0 /usr/share/zoneinfo                /usr/share/zoneinfo/

CMD ["/code-golf"]
