FROM golang:1.17.3-alpine3.14

ENV CGO_ENABLED=0 GOPATH=

RUN apk add --no-cache brotli build-base linux-headers npm tzdata zopfli

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN go build -ldflags -s -trimpath \
 && gcc -Wall -Werror -Wextra -o /usr/bin/run-lang -s -static run-lang.c

RUN ./esbuild \
 && find dist \( -name '*.js' -or -name '*.map' \) -exec brotli {} + \
 && find dist \( -name '*.js' -or -name '*.map' \) -exec zopfli {} +

FROM scratch

COPY --from=codegolf/lang-rust       ["/", "/langs/rust/rootfs/"      ] #  648 MiB
COPY --from=codegolf/lang-swift      ["/", "/langs/swift/rootfs/"     ] #  613 MiB
COPY --from=codegolf/lang-haskell    ["/", "/langs/haskell/rootfs/"   ] #  332 MiB
COPY --from=codegolf/lang-julia      ["/", "/langs/julia/rootfs/"     ] #  259 MiB
COPY --from=codegolf/lang-zig        ["/", "/langs/zig/rootfs/"       ] #  242 MiB
COPY --from=codegolf/lang-crystal    ["/", "/langs/crystal/rootfs/"   ] #  241 MiB
COPY --from=codegolf/lang-powershell ["/", "/langs/powershell/rootfs/"] #  177 MiB
COPY --from=codegolf/lang-c-sharp    ["/", "/langs/c-sharp/rootfs/"   ] #  144 MiB
COPY --from=codegolf/lang-f-sharp    ["/", "/langs/f-sharp/rootfs/"   ] #  138 MiB
COPY --from=codegolf/lang-go         ["/", "/langs/go/rootfs/"        ] #  136 MiB
COPY --from=codegolf/lang-fortran    ["/", "/langs/fortran/rootfs/"   ] # 80.6 MiB
COPY --from=codegolf/lang-assembly   ["/", "/langs/assembly/rootfs/"  ] #   75 MiB
COPY --from=codegolf/lang-java       ["/", "/langs/java/rootfs/"      ] # 69.2 MiB
COPY --from=codegolf/lang-hexagony   ["/", "/langs/hexagony/rootfs/"  ] # 63.1 MiB
COPY --from=codegolf/lang-python     ["/", "/langs/python/rootfs/"    ] # 58.4 MiB
COPY --from=codegolf/lang-raku       ["/", "/langs/raku/rootfs/"      ] # 46.6 MiB
COPY --from=codegolf/lang-v          ["/", "/langs/v/rootfs/"         ] # 35.4 MiB
COPY --from=codegolf/lang-lisp       ["/", "/langs/lisp/rootfs/"      ] # 33.6 MiB
COPY --from=codegolf/lang-pascal     ["/", "/langs/pascal/rootfs/"    ] # 31.3 MiB
COPY --from=codegolf/lang-prolog     ["/", "/langs/prolog/rootfs/"    ] # 30.7 MiB
COPY --from=codegolf/lang-viml       ["/", "/langs/viml/rootfs/"      ] # 21.3 MiB
COPY --from=codegolf/lang-javascript ["/", "/langs/javascript/rootfs/"] # 20.7 MiB
COPY --from=codegolf/lang-ruby       ["/", "/langs/ruby/rootfs/"      ] # 14.4 MiB
COPY --from=codegolf/lang-nim        ["/", "/langs/nim/rootfs/"       ] # 13.2 MiB
COPY --from=codegolf/lang-php        ["/", "/langs/php/rootfs/"       ] # 10.5 MiB
COPY --from=codegolf/lang-fish       ["/", "/langs/fish/rootfs/"      ] # 4.98 MiB
COPY --from=codegolf/lang-j          ["/", "/langs/j/rootfs/"         ] # 4.77 MiB
COPY --from=codegolf/lang-brainfuck  ["/", "/langs/brainfuck/rootfs/" ] # 4.56 MiB
COPY --from=codegolf/lang-perl       ["/", "/langs/perl/rootfs/"      ] # 4.32 MiB
COPY --from=codegolf/lang-cobol      ["/", "/langs/cobol/rootfs/"     ] #  4.1 MiB
COPY --from=codegolf/lang-c          ["/", "/langs/c/rootfs/"         ] # 1.61 MiB
COPY --from=codegolf/lang-bash       ["/", "/langs/bash/rootfs/"      ] # 1.14 MiB
COPY --from=codegolf/lang-sql        ["/", "/langs/sql/rootfs/"       ] # 1.03 MiB
COPY --from=codegolf/lang-lua        ["/", "/langs/lua/rootfs/"       ] #  338 KiB

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
