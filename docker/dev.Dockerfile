FROM golang:1.23.0-alpine3.20

ENV CGO_ENABLED=0 GOPATH= TZ=Europe/London

# curl is used for the e2e healthcheck.
RUN apk add --no-cache build-base curl git linux-headers tzdata \
 && GOBIN=/bin go install github.com/cespare/reflex@latest

COPY --from=codegolf/lang-swift        ["/", "/langs/swift/rootfs/"       ] #  462 MiB
COPY --from=codegolf/lang-haskell      ["/", "/langs/haskell/rootfs/"     ] #  388 MiB
COPY --from=codegolf/lang-go           ["/", "/langs/go/rootfs/"          ] #  382 MiB
COPY --from=codegolf/lang-rust         ["/", "/langs/rust/rootfs/"        ] #  331 MiB
COPY --from=codegolf/lang-julia        ["/", "/langs/julia/rootfs/"       ] #  310 MiB
COPY --from=codegolf/lang-d            ["/", "/langs/d/rootfs/"           ] #  309 MiB
COPY --from=codegolf/lang-zig          ["/", "/langs/zig/rootfs/"         ] #  300 MiB
COPY --from=codegolf/lang-crystal      ["/", "/langs/crystal/rootfs/"     ] #  270 MiB
COPY --from=codegolf/lang-dart         ["/", "/langs/dart/rootfs/"        ] #  235 MiB
COPY --from=codegolf/lang-basic        ["/", "/langs/basic/rootfs/"       ] #  205 MiB
COPY --from=codegolf/lang-powershell   ["/", "/langs/powershell/rootfs/"  ] #  178 MiB
COPY --from=codegolf/lang-factor       ["/", "/langs/factor/rootfs/"      ] #  172 MiB
COPY --from=codegolf/lang-rockstar     ["/", "/langs/rockstar/rootfs/"    ] #  162 MiB
COPY --from=codegolf/lang-elixir       ["/", "/langs/elixir/rootfs/"      ] #  158 MiB
COPY --from=codegolf/lang-c-sharp      ["/", "/langs/c-sharp/rootfs/"     ] #  151 MiB
COPY --from=codegolf/lang-f-sharp      ["/", "/langs/f-sharp/rootfs/"     ] #  148 MiB
COPY --from=codegolf/lang-cpp          ["/", "/langs/cpp/rootfs/"         ] #  142 MiB
COPY --from=codegolf/lang-coconut      ["/", "/langs/coconut/rootfs/"     ] #  127 MiB
COPY --from=codegolf/lang-assembly     ["/", "/langs/assembly/rootfs/"    ] #  109 MiB
COPY --from=codegolf/lang-fortran      ["/", "/langs/fortran/rootfs/"     ] # 94.1 MiB
COPY --from=codegolf/lang-ocaml        ["/", "/langs/ocaml/rootfs/"       ] # 91.4 MiB
COPY --from=codegolf/lang-r            ["/", "/langs/r/rootfs/"           ] # 81.9 MiB
COPY --from=codegolf/lang-raku         ["/", "/langs/raku/rootfs/"        ] # 73.6 MiB
COPY --from=codegolf/lang-python       ["/", "/langs/python/rootfs/"      ] # 70.8 MiB
COPY --from=codegolf/lang-coffeescript ["/", "/langs/coffeescript/rootfs/"] # 70.7 MiB
COPY --from=codegolf/lang-civet        ["/", "/langs/civet/rootfs/"       ] # 68.6 MiB
COPY --from=codegolf/lang-clojure      ["/", "/langs/clojure/rootfs/"     ] # 65.8 MiB
COPY --from=codegolf/lang-v            ["/", "/langs/v/rootfs/"           ] # 60.9 MiB
COPY --from=codegolf/lang-java         ["/", "/langs/java/rootfs/"        ] # 58.1 MiB
COPY --from=codegolf/lang-prolog       ["/", "/langs/prolog/rootfs/"      ] # 49.2 MiB
COPY --from=codegolf/lang-javascript   ["/", "/langs/javascript/rootfs/"  ] # 41.1 MiB
COPY --from=codegolf/lang-uiua         ["/", "/langs/uiua/rootfs/"        ] # 31.5 MiB
COPY --from=codegolf/lang-lisp         ["/", "/langs/lisp/rootfs/"        ] # 31.1 MiB
COPY --from=codegolf/lang-pascal       ["/", "/langs/pascal/rootfs/"      ] # 31.1 MiB
COPY --from=codegolf/lang-golfscript   ["/", "/langs/golfscript/rootfs/"  ] # 28.6 MiB
COPY --from=codegolf/lang-ruby         ["/", "/langs/ruby/rootfs/"        ] # 28.6 MiB
COPY --from=codegolf/lang-viml         ["/", "/langs/viml/rootfs/"        ] # 24.3 MiB
COPY --from=codegolf/lang-nim          ["/", "/langs/nim/rootfs/"         ] # 15.0 MiB
COPY --from=codegolf/lang-j            ["/", "/langs/j/rootfs/"           ] # 11.2 MiB
COPY --from=codegolf/lang-tex          ["/", "/langs/tex/rootfs/"         ] # 9.67 MiB
COPY --from=codegolf/lang-hexagony     ["/", "/langs/hexagony/rootfs/"    ] # 8.81 MiB
COPY --from=codegolf/lang-php          ["/", "/langs/php/rootfs/"         ] # 8.40 MiB
COPY --from=codegolf/lang-arturo       ["/", "/langs/arturo/rootfs/"      ] # 6.21 MiB
COPY --from=codegolf/lang-perl         ["/", "/langs/perl/rootfs/"        ] # 5.51 MiB
COPY --from=codegolf/lang-tcl          ["/", "/langs/tcl/rootfs/"         ] # 5.25 MiB
COPY --from=codegolf/lang-fish         ["/", "/langs/fish/rootfs/"        ] # 4.66 MiB
COPY --from=codegolf/lang-rexx         ["/", "/langs/rexx/rootfs/"        ] # 4.63 MiB
COPY --from=codegolf/lang-cobol        ["/", "/langs/cobol/rootfs/"       ] # 4.56 MiB
COPY --from=codegolf/lang-jq           ["/", "/langs/jq/rootfs/"          ] # 3.41 MiB
COPY --from=codegolf/lang-hush         ["/", "/langs/hush/rootfs/"        ] # 3.27 MiB
COPY --from=codegolf/lang-forth        ["/", "/langs/forth/rootfs/"       ] # 2.86 MiB
COPY --from=codegolf/lang-awk          ["/", "/langs/awk/rootfs/"         ] # 1.77 MiB
COPY --from=codegolf/lang-c            ["/", "/langs/c/rootfs/"           ] # 1.70 MiB
COPY --from=codegolf/lang-rebol        ["/", "/langs/rebol/rootfs/"       ] # 1.22 MiB
COPY --from=codegolf/lang-bash         ["/", "/langs/bash/rootfs/"        ] # 1.20 MiB
COPY --from=codegolf/lang-sql          ["/", "/langs/sql/rootfs/"         ] # 1.20 MiB
COPY --from=codegolf/lang-berry        ["/", "/langs/berry/rootfs/"       ] #  973 KiB
COPY --from=codegolf/lang-janet        ["/", "/langs/janet/rootfs/"       ] #  869 KiB
COPY --from=codegolf/lang-fennel       ["/", "/langs/fennel/rootfs/"      ] #  636 KiB
COPY --from=codegolf/lang-k            ["/", "/langs/k/rootfs/"           ] #  621 KiB
COPY --from=codegolf/lang-wren         ["/", "/langs/wren/rootfs/"        ] #  496 KiB
COPY --from=codegolf/lang-lua          ["/", "/langs/lua/rootfs/"         ] #  354 KiB
COPY --from=codegolf/lang-sed          ["/", "/langs/sed/rootfs/"         ] #  236 KiB
COPY --from=codegolf/lang-brainfuck    ["/", "/langs/brainfuck/rootfs/"   ] # 51.2 KiB

COPY run-lang.c ./

RUN gcc -Wall -Werror -Wextra -o /usr/bin/run-lang -s -static run-lang.c

# reflex reruns a command when files change.
CMD reflex -sd none -r '\.(css|go|html|json|pem|svg|toml|txt)$' -R '_test\.go$' -- go run .
