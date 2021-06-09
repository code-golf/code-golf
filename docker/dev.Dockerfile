FROM golang:1.16.5-alpine3.13

ENV CGO_ENABLED=0 GOPATH= TZ=Europe/London

RUN apk add --no-cache build-base git linux-headers tzdata \
 && GOBIN=/bin go get github.com/cespare/reflex

COPY --from=codegolf/lang-swift      ["/", "/langs/swift/rootfs/"     ] #  691 MiB
COPY --from=codegolf/lang-rust       ["/", "/langs/rust/rootfs/"      ] #  669 MiB
COPY --from=codegolf/lang-haskell    ["/", "/langs/haskell/rootfs/"   ] #  332 MiB
COPY --from=codegolf/lang-julia      ["/", "/langs/julia/rootfs/"     ] #  287 MiB
COPY --from=codegolf/lang-crystal    ["/", "/langs/crystal/rootfs/"   ] #  284 MiB
COPY --from=codegolf/lang-zig        ["/", "/langs/zig/rootfs/"       ] #  241 MiB
COPY --from=codegolf/lang-powershell ["/", "/langs/powershell/rootfs/"] #  185 MiB
COPY --from=codegolf/lang-c-sharp    ["/", "/langs/c-sharp/rootfs/"   ] #  129 MiB
COPY --from=codegolf/lang-go         ["/", "/langs/go/rootfs/"        ] #  124 MiB
COPY --from=codegolf/lang-f-sharp    ["/", "/langs/f-sharp/rootfs/"   ] #  119 MiB
COPY --from=codegolf/lang-v          ["/", "/langs/v/rootfs/"         ] # 98.1 MiB
COPY --from=codegolf/lang-fortran    ["/", "/langs/fortran/rootfs/"   ] # 85.8 MiB
COPY --from=codegolf/lang-java       ["/", "/langs/java/rootfs/"      ] # 69.2 MiB
COPY --from=codegolf/lang-hexagony   ["/", "/langs/hexagony/rootfs/"  ] # 62.6 MiB
COPY --from=codegolf/lang-python     ["/", "/langs/python/rootfs/"    ] #   57 MiB
COPY --from=codegolf/lang-raku       ["/", "/langs/raku/rootfs/"      ] # 53.9 MiB
COPY --from=codegolf/lang-assembly   ["/", "/langs/assembly/rootfs/"  ] # 48.7 MiB
COPY --from=codegolf/lang-lisp       ["/", "/langs/lisp/rootfs/"      ] # 33.6 MiB
COPY --from=codegolf/lang-javascript ["/", "/langs/javascript/rootfs/"] # 21.6 MiB
COPY --from=codegolf/lang-nim        ["/", "/langs/nim/rootfs/"       ] # 21.6 MiB
COPY --from=codegolf/lang-ruby       ["/", "/langs/ruby/rootfs/"      ] # 14.5 MiB
COPY --from=codegolf/lang-php        ["/", "/langs/php/rootfs/"       ] # 10.5 MiB
COPY --from=codegolf/lang-brainfuck  ["/", "/langs/brainfuck/rootfs/" ] # 4.55 MiB
COPY --from=codegolf/lang-perl       ["/", "/langs/perl/rootfs/"      ] # 4.32 MiB
COPY --from=codegolf/lang-cobol      ["/", "/langs/cobol/rootfs/"     ] # 4.09 MiB
COPY --from=codegolf/lang-j          ["/", "/langs/j/rootfs/"         ] # 3.32 MiB
COPY --from=codegolf/lang-c          ["/", "/langs/c/rootfs/"         ] # 1.61 MiB
COPY --from=codegolf/lang-bash       ["/", "/langs/bash/rootfs/"      ] # 1.15 MiB
COPY --from=codegolf/lang-sql        ["/", "/langs/sql/rootfs/"       ] # 1.02 MiB
COPY --from=codegolf/lang-fish       ["/", "/langs/fish/rootfs/"      ] #  477 KiB
COPY --from=codegolf/lang-lua        ["/", "/langs/lua/rootfs/"       ] #  338 KiB

COPY run-lang.c ./

RUN gcc -Wall -Werror -Wextra -o /usr/bin/run-lang -s -static run-lang.c

# reflex reruns a command when files change.
CMD reflex -sd none -r '\.(css|go|html|json|pem|svg|toml)$' -R '_test\.go$' -- go run .
