FROM golang:1.20.4-alpine3.18

ENV CGO_ENABLED=0 GOAMD64=v4 GOPATH=

RUN apk add --no-cache brotli build-base linux-headers npm tzdata zopfli

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN go build -ldflags -s -trimpath \
 && gcc -Wall -Werror -Wextra -o /usr/bin/run-lang -s -static run-lang.c

RUN ./esbuild \
 && find dist \( -name '*.css' -or  -name '*.js' -or -name '*.map' \) \
  | xargs -i -n1 -P`nproc` sh -c 'brotli {} && zopfli {}'

FROM scratch

COPY --from=codegolf/lang-swift      ["/", "/langs/swift/rootfs/"     ] #  484 MiB
COPY --from=codegolf/lang-rust       ["/", "/langs/rust/rootfs/"      ] #  426 MiB
COPY --from=codegolf/lang-haskell    ["/", "/langs/haskell/rootfs/"   ] #  406 MiB
COPY --from=codegolf/lang-julia      ["/", "/langs/julia/rootfs/"     ] #  301 MiB
COPY --from=codegolf/lang-d          ["/", "/langs/d/rootfs/"         ] #  296 MiB
COPY --from=codegolf/lang-zig        ["/", "/langs/zig/rootfs/"       ] #  262 MiB
COPY --from=codegolf/lang-crystal    ["/", "/langs/crystal/rootfs/"   ] #  242 MiB
COPY --from=codegolf/lang-dart       ["/", "/langs/dart/rootfs/"      ] #  240 MiB
COPY --from=codegolf/lang-basic      ["/", "/langs/basic/rootfs/"     ] #  216 MiB
COPY --from=codegolf/lang-powershell ["/", "/langs/powershell/rootfs/"] #  176 MiB
COPY --from=codegolf/lang-elixir     ["/", "/langs/elixir/rootfs/"    ] #  168 MiB
COPY --from=codegolf/lang-f-sharp    ["/", "/langs/f-sharp/rootfs/"   ] #  150 MiB
COPY --from=codegolf/lang-go         ["/", "/langs/go/rootfs/"        ] #  150 MiB
COPY --from=codegolf/lang-c-sharp    ["/", "/langs/c-sharp/rootfs/"   ] #  149 MiB
COPY --from=codegolf/lang-cpp        ["/", "/langs/cpp/rootfs/"       ] #  118 MiB
COPY --from=codegolf/lang-ocaml      ["/", "/langs/ocaml/rootfs/"     ] # 99.2 MiB
COPY --from=codegolf/lang-assembly   ["/", "/langs/assembly/rootfs/"  ] # 90.9 MiB
COPY --from=codegolf/lang-fortran    ["/", "/langs/fortran/rootfs/"   ] # 85.6 MiB
COPY --from=codegolf/lang-r          ["/", "/langs/r/rootfs/"         ] # 76.8 MiB
COPY --from=codegolf/lang-python     ["/", "/langs/python/rootfs/"    ] #   74 MiB
COPY --from=codegolf/lang-raku       ["/", "/langs/raku/rootfs/"      ] # 68.9 MiB
COPY --from=codegolf/lang-java       ["/", "/langs/java/rootfs/"      ] # 51.1 MiB
COPY --from=codegolf/lang-prolog     ["/", "/langs/prolog/rootfs/"    ] # 50.9 MiB
COPY --from=codegolf/lang-v          ["/", "/langs/v/rootfs/"         ] # 48.5 MiB
COPY --from=codegolf/lang-lisp       ["/", "/langs/lisp/rootfs/"      ] # 30.9 MiB
COPY --from=codegolf/lang-pascal     ["/", "/langs/pascal/rootfs/"    ] # 30.7 MiB
COPY --from=codegolf/lang-javascript ["/", "/langs/javascript/rootfs/"] # 24.3 MiB
COPY --from=codegolf/lang-golfscript ["/", "/langs/golfscript/rootfs/"] # 24.2 MiB
COPY --from=codegolf/lang-ruby       ["/", "/langs/ruby/rootfs/"      ] # 24.1 MiB
COPY --from=codegolf/lang-viml       ["/", "/langs/viml/rootfs/"      ] # 23.1 MiB
COPY --from=codegolf/lang-nim        ["/", "/langs/nim/rootfs/"       ] # 13.6 MiB
COPY --from=codegolf/lang-j          ["/", "/langs/j/rootfs/"         ] #   11 MiB
COPY --from=codegolf/lang-php        ["/", "/langs/php/rootfs/"       ] # 10.5 MiB
COPY --from=codegolf/lang-tex        ["/", "/langs/tex/rootfs/"       ] # 9.58 MiB
COPY --from=codegolf/lang-hexagony   ["/", "/langs/hexagony/rootfs/"  ] # 7.98 MiB
COPY --from=codegolf/lang-perl       ["/", "/langs/perl/rootfs/"      ] # 5.34 MiB
COPY --from=codegolf/lang-tcl        ["/", "/langs/tcl/rootfs/"       ] # 5.23 MiB
COPY --from=codegolf/lang-fish       ["/", "/langs/fish/rootfs/"      ] # 4.92 MiB
COPY --from=codegolf/lang-brainfuck  ["/", "/langs/brainfuck/rootfs/" ] #  4.5 MiB
COPY --from=codegolf/lang-cobol      ["/", "/langs/cobol/rootfs/"     ] # 4.12 MiB
COPY --from=codegolf/lang-awk        ["/", "/langs/awk/rootfs/"       ] # 1.72 MiB
COPY --from=codegolf/lang-c          ["/", "/langs/c/rootfs/"         ] # 1.63 MiB
COPY --from=codegolf/lang-janet      ["/", "/langs/janet/rootfs/"     ] # 1.37 MiB
COPY --from=codegolf/lang-bash       ["/", "/langs/bash/rootfs/"      ] # 1.19 MiB
COPY --from=codegolf/lang-sql        ["/", "/langs/sql/rootfs/"       ] # 1.13 MiB
COPY --from=codegolf/lang-k          ["/", "/langs/k/rootfs/"         ] #  526 KiB
COPY --from=codegolf/lang-wren       ["/", "/langs/wren/rootfs/"      ] #  484 KiB
COPY --from=codegolf/lang-lua        ["/", "/langs/lua/rootfs/"       ] #  342 KiB
COPY --from=codegolf/lang-sed        ["/", "/langs/sed/rootfs/"       ] #  232 KiB

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
