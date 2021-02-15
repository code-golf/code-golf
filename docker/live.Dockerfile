FROM golang:1.16rc1-alpine3.13

ENV CGO_ENABLED=0 GOPATH=

RUN apk add --no-cache build-base linux-headers tzdata

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN go build -ldflags -s \
 && gcc -Wall -Werror -Wextra -o /usr/bin/run-lang -s -static run-lang.c

FROM scratch

COPY --from=codegolf/lang-swift      ["/", "/langs/swift/rootfs/"     ] #  870 MiB
COPY --from=codegolf/lang-rust       ["/", "/langs/rust/rootfs/"      ] #  707 MiB
COPY --from=codegolf/lang-haskell    ["/", "/langs/haskell/rootfs/"   ] #  332 MiB
COPY --from=codegolf/lang-julia      ["/", "/langs/julia/rootfs/"     ] #  280 MiB
COPY --from=codegolf/lang-zig        ["/", "/langs/zig/rootfs/"       ] #  216 MiB
COPY --from=codegolf/lang-python     ["/", "/langs/python/rootfs/"    ] #  196 MiB
COPY --from=codegolf/lang-powershell ["/", "/langs/powershell/rootfs/"] #  185 MiB
COPY --from=codegolf/lang-c-sharp    ["/", "/langs/c-sharp/rootfs/"   ] #  129 MiB
COPY --from=codegolf/lang-go         ["/", "/langs/go/rootfs/"        ] #  110 MiB
COPY --from=codegolf/lang-f-sharp    ["/", "/langs/f-sharp/rootfs/"   ] #  108 MiB
COPY --from=codegolf/lang-v          ["/", "/langs/v/rootfs/"         ] # 96.6 MiB
COPY --from=codegolf/lang-fortran    ["/", "/langs/fortran/rootfs/"   ] # 85.6 MiB
COPY --from=codegolf/lang-java       ["/", "/langs/java/rootfs/"      ] # 68.8 MiB
COPY --from=codegolf/lang-hexagony   ["/", "/langs/hexagony/rootfs/"  ] # 62.6 MiB
COPY --from=codegolf/lang-raku       ["/", "/langs/raku/rootfs/"      ] # 50.6 MiB
COPY --from=codegolf/lang-lisp       ["/", "/langs/lisp/rootfs/"      ] # 33.6 MiB
COPY --from=codegolf/lang-nim        ["/", "/langs/nim/rootfs/"       ] # 21.5 MiB
COPY --from=codegolf/lang-javascript ["/", "/langs/javascript/rootfs/"] # 21.2 MiB
COPY --from=codegolf/lang-ruby       ["/", "/langs/ruby/rootfs/"      ] # 14.5 MiB
COPY --from=codegolf/lang-php        ["/", "/langs/php/rootfs/"       ] # 10.5 MiB
COPY --from=codegolf/lang-cobol      ["/", "/langs/cobol/rootfs/"     ] # 4.09 MiB
COPY --from=codegolf/lang-perl       ["/", "/langs/perl/rootfs/"      ] # 4.07 MiB
COPY --from=codegolf/lang-j          ["/", "/langs/j/rootfs/"         ] # 3.32 MiB
COPY --from=codegolf/lang-brainfuck  ["/", "/langs/brainfuck/rootfs/" ] # 1.62 MiB
COPY --from=codegolf/lang-c          ["/", "/langs/c/rootfs/"         ] # 1.61 MiB
COPY --from=codegolf/lang-bash       ["/", "/langs/bash/rootfs/"      ] # 1.15 MiB
COPY --from=codegolf/lang-sql        ["/", "/langs/sql/rootfs/"       ] #  986 KiB
COPY --from=codegolf/lang-fish       ["/", "/langs/fish/rootfs/"      ] #  486 KiB
COPY --from=codegolf/lang-lua        ["/", "/langs/lua/rootfs/"       ] #  334 KiB

COPY --from=0 /go/code-golf                      /
COPY          /*.toml /words.txt                 /
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY          /public                            /public/
COPY --from=0 /usr/bin/run-lang                  /usr/bin/
COPY --from=0 /usr/share/zoneinfo                /usr/share/zoneinfo/
COPY          /views                             /views/

CMD ["/code-golf"]
