FROM golang:1.16beta1-alpine3.12

ENV CGO_ENABLED=0 GOPATH= TZ=Europe/London

RUN apk add --no-cache build-base git linux-headers tzdata \
 && GOBIN=/bin go get github.com/cespare/reflex

COPY --from=codegolf/lang-swift      ["/", "/langs/swift/rootfs/"     ] #  869 MiB
COPY --from=codegolf/lang-rust       ["/", "/langs/rust/rootfs/"      ] #  619 MiB
COPY --from=codegolf/lang-haskell    ["/", "/langs/haskell/rootfs/"   ] #  332 MiB
COPY --from=codegolf/lang-julia      ["/", "/langs/julia/rootfs/"     ] #  279 MiB
COPY --from=codegolf/lang-zig        ["/", "/langs/zig/rootfs/"       ] #  216 MiB
COPY --from=codegolf/lang-python     ["/", "/langs/python/rootfs/"    ] #  205 MiB
COPY --from=codegolf/lang-powershell ["/", "/langs/powershell/rootfs/"] #  185 MiB
COPY --from=codegolf/lang-c-sharp    ["/", "/langs/c-sharp/rootfs/"   ] #  128 MiB
COPY --from=codegolf/lang-go         ["/", "/langs/go/rootfs/"        ] #  110 MiB
COPY --from=codegolf/lang-f-sharp    ["/", "/langs/f-sharp/rootfs/"   ] #  108 MiB
COPY --from=codegolf/lang-v          ["/", "/langs/v/rootfs/"         ] # 96.6 MiB
COPY --from=codegolf/lang-fortran    ["/", "/langs/fortran/rootfs/"   ] # 85.7 MiB
COPY --from=codegolf/lang-java       ["/", "/langs/java/rootfs/"      ] # 68.7 MiB
COPY --from=codegolf/lang-raku       ["/", "/langs/raku/rootfs/"      ] # 50.6 MiB
COPY --from=codegolf/lang-lisp       ["/", "/langs/lisp/rootfs/"      ] # 35.4 MiB
COPY --from=codegolf/lang-nim        ["/", "/langs/nim/rootfs/"       ] # 21.8 MiB
COPY --from=codegolf/lang-javascript ["/", "/langs/javascript/rootfs/"] # 21.1 MiB
COPY --from=codegolf/lang-ruby       ["/", "/langs/ruby/rootfs/"      ] # 14.5 MiB
COPY --from=codegolf/lang-php        ["/", "/langs/php/rootfs/"       ] # 10.5 MiB
COPY --from=codegolf/lang-cobol      ["/", "/langs/cobol/rootfs/"     ] # 6.12 MiB
COPY --from=codegolf/lang-perl       ["/", "/langs/perl/rootfs/"      ] # 4.04 MiB
COPY --from=codegolf/lang-j          ["/", "/langs/j/rootfs/"         ] #  3.3 MiB
COPY --from=codegolf/lang-brainfuck  ["/", "/langs/brainfuck/rootfs/" ] # 1.59 MiB
COPY --from=codegolf/lang-c          ["/", "/langs/c/rootfs/"         ] # 1.58 MiB
COPY --from=codegolf/lang-bash       ["/", "/langs/bash/rootfs/"      ] # 1.19 MiB
COPY --from=codegolf/lang-sql        ["/", "/langs/sql/rootfs/"       ] # 1.03 MiB
COPY --from=codegolf/lang-fish       ["/", "/langs/fish/rootfs/"      ] #  585 KiB
COPY --from=codegolf/lang-lua        ["/", "/langs/lua/rootfs/"       ] #  318 KiB

COPY run-lang.c ./

RUN gcc -Wall -Werror -Wextra -o /usr/bin/run-lang -s -static run-lang.c

# reflex reruns a command when files change.
CMD reflex -sd none -r '\.(css|go|html|js|pem|svg|toml)$' -R '_test\.go$' -- go run -tags dev .
