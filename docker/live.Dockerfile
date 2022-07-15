FROM golang:1.19rc1-alpine3.16

ENV CGO_ENABLED=0 GOPATH=

RUN apk add --no-cache brotli build-base linux-headers npm tzdata zopfli

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN go build -ldflags -s -trimpath \
 && gcc -Wall -Werror -Wextra -o /usr/bin/run-lang -s -static run-lang.c

RUN ./esbuild \
 && find dist \( -name '*.js' -or -name '*.map' \) \
  | xargs -i -n1 -P`nproc` sh -c 'brotli {} && zopfli {}'

FROM scratch

COPY --from=codegolf/lang-swift      ["/", "/langs/swift/rootfs/"     ] #  561 MiB
COPY --from=codegolf/lang-rust       ["/", "/langs/rust/rootfs/"      ] #  513 MiB
COPY --from=codegolf/lang-haskell    ["/", "/langs/haskell/rootfs/"   ] #  309 MiB
COPY --from=codegolf/lang-julia      ["/", "/langs/julia/rootfs/"     ] #  279 MiB
COPY --from=codegolf/lang-basic      ["/", "/langs/basic/rootfs/"     ] #  268 MiB
COPY --from=codegolf/lang-zig        ["/", "/langs/zig/rootfs/"       ] #  262 MiB
COPY --from=codegolf/lang-d          ["/", "/langs/d/rootfs/"         ] #  245 MiB
COPY --from=codegolf/lang-crystal    ["/", "/langs/crystal/rootfs/"   ] #  222 MiB
COPY --from=codegolf/lang-powershell ["/", "/langs/powershell/rootfs/"] #  177 MiB
COPY --from=codegolf/lang-elixir     ["/", "/langs/elixir/rootfs/"    ] #  166 MiB
COPY --from=codegolf/lang-c-sharp    ["/", "/langs/c-sharp/rootfs/"   ] #  145 MiB
COPY --from=codegolf/lang-go         ["/", "/langs/go/rootfs/"        ] #  145 MiB
COPY --from=codegolf/lang-f-sharp    ["/", "/langs/f-sharp/rootfs/"   ] #  140 MiB
COPY --from=codegolf/lang-cpp        ["/", "/langs/cpp/rootfs/"       ] #  117 MiB
COPY --from=codegolf/lang-assembly   ["/", "/langs/assembly/rootfs/"  ] # 79.9 MiB
COPY --from=codegolf/lang-fortran    ["/", "/langs/fortran/rootfs/"   ] # 77.1 MiB
COPY --from=codegolf/lang-hexagony   ["/", "/langs/hexagony/rootfs/"  ] # 63.2 MiB
COPY --from=codegolf/lang-python     ["/", "/langs/python/rootfs/"    ] # 55.4 MiB
COPY --from=codegolf/lang-raku       ["/", "/langs/raku/rootfs/"      ] # 54.1 MiB
COPY --from=codegolf/lang-java       ["/", "/langs/java/rootfs/"      ] # 51.1 MiB
COPY --from=codegolf/lang-v          ["/", "/langs/v/rootfs/"         ] # 38.5 MiB
COPY --from=codegolf/lang-prolog     ["/", "/langs/prolog/rootfs/"    ] # 31.7 MiB
COPY --from=codegolf/lang-pascal     ["/", "/langs/pascal/rootfs/"    ] # 31.3 MiB
COPY --from=codegolf/lang-lisp       ["/", "/langs/lisp/rootfs/"      ] # 30.9 MiB
COPY --from=codegolf/lang-ruby       ["/", "/langs/ruby/rootfs/"      ] # 24.2 MiB
COPY --from=codegolf/lang-viml       ["/", "/langs/viml/rootfs/"      ] # 22.3 MiB
COPY --from=codegolf/lang-javascript ["/", "/langs/javascript/rootfs/"] # 21.5 MiB
COPY --from=codegolf/lang-nim        ["/", "/langs/nim/rootfs/"       ] # 13.5 MiB
COPY --from=codegolf/lang-php        ["/", "/langs/php/rootfs/"       ] # 10.5 MiB
COPY --from=codegolf/lang-perl       ["/", "/langs/perl/rootfs/"      ] # 5.32 MiB
COPY --from=codegolf/lang-fish       ["/", "/langs/fish/rootfs/"      ] # 4.99 MiB
COPY --from=codegolf/lang-j          ["/", "/langs/j/rootfs/"         ] # 4.84 MiB
COPY --from=codegolf/lang-brainfuck  ["/", "/langs/brainfuck/rootfs/" ] # 4.57 MiB
COPY --from=codegolf/lang-cobol      ["/", "/langs/cobol/rootfs/"     ] # 4.09 MiB
COPY --from=codegolf/lang-c          ["/", "/langs/c/rootfs/"         ] # 1.61 MiB
COPY --from=codegolf/lang-bash       ["/", "/langs/bash/rootfs/"      ] # 1.19 MiB
COPY --from=codegolf/lang-sql        ["/", "/langs/sql/rootfs/"       ] # 1.11 MiB
COPY --from=codegolf/lang-lua        ["/", "/langs/lua/rootfs/"       ] #  342 KiB
COPY --from=codegolf/lang-k          ["/", "/langs/k/rootfs/"         ] #  258 KiB
COPY --from=codegolf/lang-sed        ["/", "/langs/sed/rootfs/"       ] #  228 KiB

COPY --from=0 /go/code-golf /go/esbuild.json     /
COPY          /css                               /css/
COPY --from=0 /go/dist                           /dist/
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY          /public                            /public/
COPY          /svg                               /svg/
COPY --from=0 /usr/bin/run-lang                  /usr/bin/
COPY --from=0 /usr/share/zoneinfo                /usr/share/zoneinfo/
COPY          /views                             /views/

CMD ["/code-golf"]
