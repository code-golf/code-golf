FROM golang:1.17.0-alpine3.14

ENV CGO_ENABLED=0 GOPATH= TZ=Europe/London

# curl is used for the e2e healthcheck.
RUN apk add --no-cache build-base curl git linux-headers tzdata \
 && GOBIN=/bin go install github.com/cespare/reflex@latest

COPY --from=codegolf/lang-rust       ["/", "/langs/rust/rootfs/"      ] #  718 MiB
COPY --from=codegolf/lang-swift      ["/", "/langs/swift/rootfs/"     ] #  693 MiB
COPY --from=codegolf/lang-haskell    ["/", "/langs/haskell/rootfs/"   ] #  332 MiB
COPY --from=codegolf/lang-julia      ["/", "/langs/julia/rootfs/"     ] #  288 MiB
COPY --from=codegolf/lang-crystal    ["/", "/langs/crystal/rootfs/"   ] #  263 MiB
COPY --from=codegolf/lang-zig        ["/", "/langs/zig/rootfs/"       ] #  241 MiB
COPY --from=codegolf/lang-powershell ["/", "/langs/powershell/rootfs/"] #  185 MiB
COPY --from=codegolf/lang-go         ["/", "/langs/go/rootfs/"        ] #  136 MiB
COPY --from=codegolf/lang-c-sharp    ["/", "/langs/c-sharp/rootfs/"   ] #  130 MiB
COPY --from=codegolf/lang-f-sharp    ["/", "/langs/f-sharp/rootfs/"   ] #  119 MiB
COPY --from=codegolf/lang-fortran    ["/", "/langs/fortran/rootfs/"   ] # 80.6 MiB
COPY --from=codegolf/lang-java       ["/", "/langs/java/rootfs/"      ] # 69.2 MiB
COPY --from=codegolf/lang-hexagony   ["/", "/langs/hexagony/rootfs/"  ] # 62.6 MiB
COPY --from=codegolf/lang-python     ["/", "/langs/python/rootfs/"    ] # 57.2 MiB
COPY --from=codegolf/lang-raku       ["/", "/langs/raku/rootfs/"      ] # 54.5 MiB
COPY --from=codegolf/lang-assembly   ["/", "/langs/assembly/rootfs/"  ] # 48.7 MiB
COPY --from=codegolf/lang-v          ["/", "/langs/v/rootfs/"         ] # 35.4 MiB
COPY --from=codegolf/lang-lisp       ["/", "/langs/lisp/rootfs/"      ] # 33.6 MiB
COPY --from=codegolf/lang-pascal     ["/", "/langs/pascal/rootfs/"    ] # 31.3 MiB
COPY --from=codegolf/lang-viml       ["/", "/langs/viml/rootfs/"      ] # 21.3 MiB
COPY --from=codegolf/lang-javascript ["/", "/langs/javascript/rootfs/"] # 20.7 MiB
COPY --from=codegolf/lang-ruby       ["/", "/langs/ruby/rootfs/"      ] # 14.4 MiB
COPY --from=codegolf/lang-nim        ["/", "/langs/nim/rootfs/"       ] # 12.1 MiB
COPY --from=codegolf/lang-php        ["/", "/langs/php/rootfs/"       ] # 10.5 MiB
COPY --from=codegolf/lang-j          ["/", "/langs/j/rootfs/"         ] # 4.77 MiB
COPY --from=codegolf/lang-brainfuck  ["/", "/langs/brainfuck/rootfs/" ] # 4.55 MiB
COPY --from=codegolf/lang-perl       ["/", "/langs/perl/rootfs/"      ] # 4.32 MiB
COPY --from=codegolf/lang-cobol      ["/", "/langs/cobol/rootfs/"     ] #  4.1 MiB
COPY --from=codegolf/lang-c          ["/", "/langs/c/rootfs/"         ] # 1.61 MiB
COPY --from=codegolf/lang-bash       ["/", "/langs/bash/rootfs/"      ] # 1.14 MiB
COPY --from=codegolf/lang-sql        ["/", "/langs/sql/rootfs/"       ] # 1.03 MiB
COPY --from=codegolf/lang-fish       ["/", "/langs/fish/rootfs/"      ] #  477 KiB
COPY --from=codegolf/lang-lua        ["/", "/langs/lua/rootfs/"       ] #  338 KiB

COPY run-lang.c ./

RUN gcc -Wall -Werror -Wextra -o /usr/bin/run-lang -s -static run-lang.c

# reflex reruns a command when files change.
CMD reflex -sd none -r '\.(css|go|html|json|pem|svg|toml)$' -R '_test\.go$' -- go run .
