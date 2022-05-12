FROM golang:1.18.2-alpine3.15

ENV CGO_ENABLED=0 GOPATH= TZ=Europe/London

# curl is used for the e2e healthcheck.
RUN apk add --no-cache build-base curl git linux-headers tzdata \
 && GOBIN=/bin go install github.com/cespare/reflex@latest

COPY --from=codegolf/lang-swift      ["/", "/langs/swift/rootfs/"     ] #  561 MiB
COPY --from=codegolf/lang-rust       ["/", "/langs/rust/rootfs/"      ] #  511 MiB
COPY --from=codegolf/lang-haskell    ["/", "/langs/haskell/rootfs/"   ] #  309 MiB
COPY --from=codegolf/lang-julia      ["/", "/langs/julia/rootfs/"     ] #  278 MiB
COPY --from=codegolf/lang-elixir     ["/", "/langs/elixir/rootfs/"    ] #  277 MiB
COPY --from=codegolf/lang-zig        ["/", "/langs/zig/rootfs/"       ] #  262 MiB
COPY --from=codegolf/lang-d          ["/", "/langs/d/rootfs/"         ] #  253 MiB
COPY --from=codegolf/lang-crystal    ["/", "/langs/crystal/rootfs/"   ] #  221 MiB
COPY --from=codegolf/lang-powershell ["/", "/langs/powershell/rootfs/"] #  177 MiB
COPY --from=codegolf/lang-c-sharp    ["/", "/langs/c-sharp/rootfs/"   ] #  145 MiB
COPY --from=codegolf/lang-go         ["/", "/langs/go/rootfs/"        ] #  145 MiB
COPY --from=codegolf/lang-f-sharp    ["/", "/langs/f-sharp/rootfs/"   ] #  140 MiB
COPY --from=codegolf/lang-cpp        ["/", "/langs/cpp/rootfs/"       ] #  115 MiB
COPY --from=codegolf/lang-fortran    ["/", "/langs/fortran/rootfs/"   ] # 80.2 MiB
COPY --from=codegolf/lang-assembly   ["/", "/langs/assembly/rootfs/"  ] # 79.9 MiB
COPY --from=codegolf/lang-hexagony   ["/", "/langs/hexagony/rootfs/"  ] # 63.2 MiB
COPY --from=codegolf/lang-python     ["/", "/langs/python/rootfs/"    ] # 57.8 MiB
COPY --from=codegolf/lang-raku       ["/", "/langs/raku/rootfs/"      ] # 57.4 MiB
COPY --from=codegolf/lang-java       ["/", "/langs/java/rootfs/"      ] # 51.1 MiB
COPY --from=codegolf/lang-v          ["/", "/langs/v/rootfs/"         ] # 38.1 MiB
COPY --from=codegolf/lang-lisp       ["/", "/langs/lisp/rootfs/"      ] # 33.6 MiB
COPY --from=codegolf/lang-prolog     ["/", "/langs/prolog/rootfs/"    ] # 31.7 MiB
COPY --from=codegolf/lang-pascal     ["/", "/langs/pascal/rootfs/"    ] # 31.2 MiB
COPY --from=codegolf/lang-ruby       ["/", "/langs/ruby/rootfs/"      ] # 24.3 MiB
COPY --from=codegolf/lang-javascript ["/", "/langs/javascript/rootfs/"] # 21.5 MiB
COPY --from=codegolf/lang-viml       ["/", "/langs/viml/rootfs/"      ] # 21.3 MiB
COPY --from=codegolf/lang-nim        ["/", "/langs/nim/rootfs/"       ] # 13.5 MiB
COPY --from=codegolf/lang-php        ["/", "/langs/php/rootfs/"       ] # 10.5 MiB
COPY --from=codegolf/lang-fish       ["/", "/langs/fish/rootfs/"      ] # 4.98 MiB
COPY --from=codegolf/lang-j          ["/", "/langs/j/rootfs/"         ] # 4.84 MiB
COPY --from=codegolf/lang-brainfuck  ["/", "/langs/brainfuck/rootfs/" ] # 4.56 MiB
COPY --from=codegolf/lang-perl       ["/", "/langs/perl/rootfs/"      ] # 4.32 MiB
COPY --from=codegolf/lang-cobol      ["/", "/langs/cobol/rootfs/"     ] #  4.1 MiB
COPY --from=codegolf/lang-c          ["/", "/langs/c/rootfs/"         ] # 1.61 MiB
COPY --from=codegolf/lang-bash       ["/", "/langs/bash/rootfs/"      ] # 1.19 MiB
COPY --from=codegolf/lang-sql        ["/", "/langs/sql/rootfs/"       ] # 1.09 MiB
COPY --from=codegolf/lang-lua        ["/", "/langs/lua/rootfs/"       ] #  342 KiB
COPY --from=codegolf/lang-k          ["/", "/langs/k/rootfs/"         ] #  262 KiB
COPY --from=codegolf/lang-sed        ["/", "/langs/sed/rootfs/"       ] #  228 KiB

COPY run-lang.c ./

RUN gcc -Wall -Werror -Wextra -o /usr/bin/run-lang -s -static run-lang.c

# reflex reruns a command when files change.
CMD reflex -sd none -r '\.(css|go|html|json|pem|svg|toml)$' -R '_test\.go$' -- go run .
