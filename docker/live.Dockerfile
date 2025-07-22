FROM golang:1.24.5-alpine3.22

ENV CGO_ENABLED=0 GOAMD64=v4 GOPATH=

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

COPY --from=codegolf/lang-swift        ["/", "/langs/swift/rootfs/"       ] #  543 MiB
COPY --from=codegolf/lang-rust         ["/", "/langs/rust/rootfs/"        ] #  392 MiB
COPY --from=codegolf/lang-haskell      ["/", "/langs/haskell/rootfs/"     ] #  339 MiB
COPY --from=codegolf/lang-go           ["/", "/langs/go/rootfs/"          ] #  324 MiB
COPY --from=codegolf/lang-julia        ["/", "/langs/julia/rootfs/"       ] #  300 MiB
COPY --from=codegolf/lang-odin         ["/", "/langs/odin/rootfs/"        ] #  299 MiB
COPY --from=codegolf/lang-crystal      ["/", "/langs/crystal/rootfs/"     ] #  271 MiB
COPY --from=codegolf/lang-basic        ["/", "/langs/basic/rootfs/"       ] #  200 MiB
COPY --from=codegolf/lang-zig          ["/", "/langs/zig/rootfs/"         ] #  181 MiB
COPY --from=codegolf/lang-factor       ["/", "/langs/factor/rootfs/"      ] #  173 MiB
COPY --from=codegolf/lang-dart         ["/", "/langs/dart/rootfs/"        ] #  168 MiB
COPY --from=codegolf/lang-cpp          ["/", "/langs/cpp/rootfs/"         ] #  151 MiB
COPY --from=codegolf/lang-kotlin       ["/", "/langs/kotlin/rootfs/"      ] #  140 MiB
COPY --from=codegolf/lang-powershell   ["/", "/langs/powershell/rootfs/"  ] #  136 MiB
COPY --from=codegolf/lang-assembly     ["/", "/langs/assembly/rootfs/"    ] #  118 MiB
COPY --from=codegolf/lang-f-sharp      ["/", "/langs/f-sharp/rootfs/"     ] #  110 MiB
COPY --from=codegolf/lang-c-sharp      ["/", "/langs/c-sharp/rootfs/"     ] #  108 MiB
COPY --from=codegolf/lang-vala         ["/", "/langs/vala/rootfs/"        ] # 99.6 MiB
COPY --from=codegolf/lang-d            ["/", "/langs/d/rootfs/"           ] # 95.1 MiB
COPY --from=codegolf/lang-scala        ["/", "/langs/scala/rootfs/"       ] # 89.0 MiB
COPY --from=codegolf/lang-coconut      ["/", "/langs/coconut/rootfs/"     ] # 84.5 MiB
COPY --from=codegolf/lang-ocaml        ["/", "/langs/ocaml/rootfs/"       ] # 83.3 MiB
COPY --from=codegolf/lang-vyxal        ["/", "/langs/vyxal/rootfs/"       ] # 83.0 MiB
COPY --from=codegolf/lang-civet        ["/", "/langs/civet/rootfs/"       ] # 81.3 MiB
COPY --from=codegolf/lang-coffeescript ["/", "/langs/coffeescript/rootfs/"] # 80.7 MiB
COPY --from=codegolf/lang-groovy       ["/", "/langs/groovy/rootfs/"      ] # 76.5 MiB
COPY --from=codegolf/lang-stax         ["/", "/langs/stax/rootfs/"        ] # 76.5 MiB
COPY --from=codegolf/lang-raku         ["/", "/langs/raku/rootfs/"        ] # 75.6 MiB
COPY --from=codegolf/lang-gleam        ["/", "/langs/gleam/rootfs/"       ] # 73.5 MiB
COPY --from=codegolf/lang-fortran      ["/", "/langs/fortran/rootfs/"     ] # 72.8 MiB
COPY --from=codegolf/lang-clojure      ["/", "/langs/clojure/rootfs/"     ] # 67.3 MiB
COPY --from=codegolf/lang-cjam         ["/", "/langs/cjam/rootfs/"        ] # 62.7 MiB
COPY --from=codegolf/lang-java         ["/", "/langs/java/rootfs/"        ] # 57.3 MiB
COPY --from=codegolf/lang-elixir       ["/", "/langs/elixir/rootfs/"      ] # 56.0 MiB
COPY --from=codegolf/lang-prolog       ["/", "/langs/prolog/rootfs/"      ] # 49.3 MiB
COPY --from=codegolf/lang-erlang       ["/", "/langs/erlang/rootfs/"      ] # 49.0 MiB
COPY --from=codegolf/lang-javascript   ["/", "/langs/javascript/rootfs/"  ] # 43.9 MiB
COPY --from=codegolf/lang-hy           ["/", "/langs/hy/rootfs/"          ] # 41.7 MiB
COPY --from=codegolf/lang-haxe         ["/", "/langs/haxe/rootfs/"        ] # 41.1 MiB
COPY --from=codegolf/lang-v            ["/", "/langs/v/rootfs/"           ] # 38.3 MiB
COPY --from=codegolf/lang-iogii        ["/", "/langs/iogii/rootfs/"       ] # 33.4 MiB
COPY --from=codegolf/lang-golfscript   ["/", "/langs/golfscript/rootfs/"  ] # 33.3 MiB
COPY --from=codegolf/lang-ruby         ["/", "/langs/ruby/rootfs/"        ] # 33.2 MiB
COPY --from=codegolf/lang-common-lisp  ["/", "/langs/common-lisp/rootfs/" ] # 31.0 MiB
COPY --from=codegolf/lang-r            ["/", "/langs/r/rootfs/"           ] # 30.5 MiB
COPY --from=codegolf/lang-racket       ["/", "/langs/racket/rootfs/"      ] # 29.4 MiB
COPY --from=codegolf/lang-python       ["/", "/langs/python/rootfs/"      ] # 29.1 MiB
COPY --from=codegolf/lang-05ab1e       ["/", "/langs/05ab1e/rootfs/"      ] # 28.2 MiB
COPY --from=codegolf/lang-pascal       ["/", "/langs/pascal/rootfs/"      ] # 26.8 MiB
COPY --from=codegolf/lang-viml         ["/", "/langs/viml/rootfs/"        ] # 24.1 MiB
COPY --from=codegolf/lang-uiua         ["/", "/langs/uiua/rootfs/"        ] # 19.5 MiB
COPY --from=codegolf/lang-qore         ["/", "/langs/qore/rootfs/"        ] # 17.3 MiB
COPY --from=codegolf/lang-apl          ["/", "/langs/apl/rootfs/"         ] # 15.0 MiB
COPY --from=codegolf/lang-nim          ["/", "/langs/nim/rootfs/"         ] # 14.9 MiB
COPY --from=codegolf/lang-harbour      ["/", "/langs/harbour/rootfs/"     ] # 12.2 MiB
COPY --from=codegolf/lang-egel         ["/", "/langs/egel/rootfs/"        ] # 11.6 MiB
COPY --from=codegolf/lang-hare         ["/", "/langs/hare/rootfs/"        ] # 11.5 MiB
COPY --from=codegolf/lang-j            ["/", "/langs/j/rootfs/"           ] # 11.3 MiB
COPY --from=codegolf/lang-picat        ["/", "/langs/picat/rootfs/"       ] # 11.1 MiB
COPY --from=codegolf/lang-tex          ["/", "/langs/tex/rootfs/"         ] # 9.67 MiB
COPY --from=codegolf/lang-hexagony     ["/", "/langs/hexagony/rootfs/"    ] # 8.90 MiB
COPY --from=codegolf/lang-php          ["/", "/langs/php/rootfs/"         ] # 8.39 MiB
COPY --from=codegolf/lang-scheme       ["/", "/langs/scheme/rootfs/"      ] # 7.40 MiB
COPY --from=codegolf/lang-tcl          ["/", "/langs/tcl/rootfs/"         ] # 5.68 MiB
COPY --from=codegolf/lang-perl         ["/", "/langs/perl/rootfs/"        ] # 5.66 MiB
COPY --from=codegolf/lang-arturo       ["/", "/langs/arturo/rootfs/"      ] # 5.23 MiB
COPY --from=codegolf/lang-algol-68     ["/", "/langs/algol-68/rootfs/"    ] # 4.67 MiB
COPY --from=codegolf/lang-rockstar     ["/", "/langs/rockstar/rootfs/"    ] # 4.26 MiB
COPY --from=codegolf/lang-cobol        ["/", "/langs/cobol/rootfs/"       ] # 4.21 MiB
COPY --from=codegolf/lang-squirrel     ["/", "/langs/squirrel/rootfs/"    ] # 4.06 MiB
COPY --from=codegolf/lang-rexx         ["/", "/langs/rexx/rootfs/"        ] # 2.98 MiB
COPY --from=codegolf/lang-forth        ["/", "/langs/forth/rootfs/"       ] # 2.85 MiB
COPY --from=codegolf/lang-befunge      ["/", "/langs/befunge/rootfs/"     ] # 2.48 MiB
COPY --from=codegolf/lang-hush         ["/", "/langs/hush/rootfs/"        ] # 2.30 MiB
COPY --from=codegolf/lang-awk          ["/", "/langs/awk/rootfs/"         ] # 1.91 MiB
COPY --from=codegolf/lang-c            ["/", "/langs/c/rootfs/"           ] # 1.73 MiB
COPY --from=codegolf/lang-bash         ["/", "/langs/bash/rootfs/"        ] # 1.27 MiB
COPY --from=codegolf/lang-bqn          ["/", "/langs/bqn/rootfs/"         ] # 1.26 MiB
COPY --from=codegolf/lang-rebol        ["/", "/langs/rebol/rootfs/"       ] # 1.23 MiB
COPY --from=codegolf/lang-sql          ["/", "/langs/sql/rootfs/"         ] # 1.22 MiB
COPY --from=codegolf/lang-jq           ["/", "/langs/jq/rootfs/"          ] #  989 KiB
COPY --from=codegolf/lang-berry        ["/", "/langs/berry/rootfs/"       ] #  986 KiB
COPY --from=codegolf/lang-janet        ["/", "/langs/janet/rootfs/"       ] #  898 KiB
COPY --from=codegolf/lang-fennel       ["/", "/langs/fennel/rootfs/"      ] #  681 KiB
COPY --from=codegolf/lang-k            ["/", "/langs/k/rootfs/"           ] #  605 KiB
COPY --from=codegolf/lang-fish         ["/", "/langs/fish/rootfs/"        ] #  549 KiB
COPY --from=codegolf/lang-wren         ["/", "/langs/wren/rootfs/"        ] #  505 KiB
COPY --from=codegolf/lang-lua          ["/", "/langs/lua/rootfs/"         ] #  362 KiB
COPY --from=codegolf/lang-sed          ["/", "/langs/sed/rootfs/"         ] #  244 KiB
COPY --from=codegolf/lang-brainfuck    ["/", "/langs/brainfuck/rootfs/"   ] # 51.1 KiB

COPY --from=0 /go/hash-langs /

RUN ["/hash-langs"]

COPY --from=0 /go/code-golf                      /
COPY --from=0 /go/dist                           /dist/
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY          /public                            /public/
COPY --from=0 /usr/bin/run-lang                  /usr/bin/
COPY --from=0 /usr/share/zoneinfo                /usr/share/zoneinfo/

CMD ["/code-golf"]
