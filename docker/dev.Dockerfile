FROM golang:1.23.5-alpine3.21

ENV CGO_ENABLED=0 GOPATH= TZ=Europe/London

# curl is used for the e2e healthcheck.
RUN apk add --no-cache build-base curl git linux-headers tzdata \
 && GOBIN=/bin go install github.com/cespare/reflex@latest

COPY --from=codegolf/lang-swift        ["/", "/langs/swift/rootfs/"       ] #  516 MiB
COPY --from=codegolf/lang-go           ["/", "/langs/go/rootfs/"          ] #  383 MiB
COPY --from=codegolf/lang-rust         ["/", "/langs/rust/rootfs/"        ] #  356 MiB
COPY --from=codegolf/lang-haskell      ["/", "/langs/haskell/rootfs/"     ] #  339 MiB
COPY --from=codegolf/lang-julia        ["/", "/langs/julia/rootfs/"       ] #  300 MiB
COPY --from=codegolf/lang-zig          ["/", "/langs/zig/rootfs/"         ] #  298 MiB
COPY --from=codegolf/lang-odin         ["/", "/langs/odin/rootfs/"        ] #  280 MiB
COPY --from=codegolf/lang-crystal      ["/", "/langs/crystal/rootfs/"     ] #  270 MiB
COPY --from=codegolf/lang-dart         ["/", "/langs/dart/rootfs/"        ] #  239 MiB
COPY --from=codegolf/lang-kotlin       ["/", "/langs/kotlin/rootfs/"      ] #  224 MiB
COPY --from=codegolf/lang-basic        ["/", "/langs/basic/rootfs/"       ] #  205 MiB
COPY --from=codegolf/lang-factor       ["/", "/langs/factor/rootfs/"      ] #  173 MiB
COPY --from=codegolf/lang-scala        ["/", "/langs/scala/rootfs/"       ] #  164 MiB
COPY --from=codegolf/lang-vyxal        ["/", "/langs/vyxal/rootfs/"       ] #  152 MiB
COPY --from=codegolf/lang-groovy       ["/", "/langs/groovy/rootfs/"      ] #  146 MiB
COPY --from=codegolf/lang-cpp          ["/", "/langs/cpp/rootfs/"         ] #  142 MiB
COPY --from=codegolf/lang-powershell   ["/", "/langs/powershell/rootfs/"  ] #  135 MiB
COPY --from=codegolf/lang-cjam         ["/", "/langs/cjam/rootfs/"        ] #  132 MiB
COPY --from=codegolf/lang-rockstar     ["/", "/langs/rockstar/rootfs/"    ] #  127 MiB
COPY --from=codegolf/lang-assembly     ["/", "/langs/assembly/rootfs/"    ] #  114 MiB
COPY --from=codegolf/lang-f-sharp      ["/", "/langs/f-sharp/rootfs/"     ] #  109 MiB
COPY --from=codegolf/lang-c-sharp      ["/", "/langs/c-sharp/rootfs/"     ] #  107 MiB
COPY --from=codegolf/lang-fortran      ["/", "/langs/fortran/rootfs/"     ] # 95.1 MiB
COPY --from=codegolf/lang-d            ["/", "/langs/d/rootfs/"           ] # 93.3 MiB
COPY --from=codegolf/lang-gleam        ["/", "/langs/gleam/rootfs/"       ] # 91.2 MiB
COPY --from=codegolf/lang-coconut      ["/", "/langs/coconut/rootfs/"     ] # 87.3 MiB
COPY --from=codegolf/lang-ocaml        ["/", "/langs/ocaml/rootfs/"       ] # 84.5 MiB
COPY --from=codegolf/lang-stax         ["/", "/langs/stax/rootfs/"        ] # 76.5 MiB
COPY --from=codegolf/lang-raku         ["/", "/langs/raku/rootfs/"        ] # 76.0 MiB
COPY --from=codegolf/lang-elixir       ["/", "/langs/elixir/rootfs/"      ] # 74.3 MiB
COPY --from=codegolf/lang-civet        ["/", "/langs/civet/rootfs/"       ] # 71.8 MiB
COPY --from=codegolf/lang-coffeescript ["/", "/langs/coffeescript/rootfs/"] # 71.3 MiB
COPY --from=codegolf/lang-erlang       ["/", "/langs/erlang/rootfs/"      ] # 68.5 MiB
COPY --from=codegolf/lang-clojure      ["/", "/langs/clojure/rootfs/"     ] # 66.4 MiB
COPY --from=codegolf/lang-java         ["/", "/langs/java/rootfs/"        ] # 58.4 MiB
COPY --from=codegolf/lang-v            ["/", "/langs/v/rootfs/"           ] # 50.2 MiB
COPY --from=codegolf/lang-egel         ["/", "/langs/egel/rootfs/"        ] # 49.7 MiB
COPY --from=codegolf/lang-prolog       ["/", "/langs/prolog/rootfs/"      ] # 49.3 MiB
COPY --from=codegolf/lang-javascript   ["/", "/langs/javascript/rootfs/"  ] # 41.9 MiB
COPY --from=codegolf/lang-hy           ["/", "/langs/hy/rootfs/"          ] # 41.8 MiB
COPY --from=codegolf/lang-haxe         ["/", "/langs/haxe/rootfs/"        ] # 39.8 MiB
COPY --from=codegolf/lang-iogii        ["/", "/langs/iogii/rootfs/"       ] # 34.1 MiB
COPY --from=codegolf/lang-ruby         ["/", "/langs/ruby/rootfs/"        ] # 33.2 MiB
COPY --from=codegolf/lang-common-lisp  ["/", "/langs/common-lisp/rootfs/" ] # 31.1 MiB
COPY --from=codegolf/lang-r            ["/", "/langs/r/rootfs/"           ] # 30.1 MiB
COPY --from=codegolf/lang-uiua         ["/", "/langs/uiua/rootfs/"        ] # 30.0 MiB
COPY --from=codegolf/lang-python       ["/", "/langs/python/rootfs/"      ] # 29.2 MiB
COPY --from=codegolf/lang-racket       ["/", "/langs/racket/rootfs/"      ] # 29.0 MiB
COPY --from=codegolf/lang-golfscript   ["/", "/langs/golfscript/rootfs/"  ] # 28.6 MiB
COPY --from=codegolf/lang-05ab1e       ["/", "/langs/05ab1e/rootfs/"      ] # 28.2 MiB
COPY --from=codegolf/lang-pascal       ["/", "/langs/pascal/rootfs/"      ] # 26.6 MiB
COPY --from=codegolf/lang-viml         ["/", "/langs/viml/rootfs/"        ] # 24.1 MiB
COPY --from=codegolf/lang-nim          ["/", "/langs/nim/rootfs/"         ] # 15.6 MiB
COPY --from=codegolf/lang-apl          ["/", "/langs/apl/rootfs/"         ] # 15.0 MiB
COPY --from=codegolf/lang-harbour      ["/", "/langs/harbour/rootfs/"     ] # 12.2 MiB
COPY --from=codegolf/lang-hare         ["/", "/langs/hare/rootfs/"        ] # 11.4 MiB
COPY --from=codegolf/lang-j            ["/", "/langs/j/rootfs/"           ] # 11.2 MiB
COPY --from=codegolf/lang-picat        ["/", "/langs/picat/rootfs/"       ] # 11.1 MiB
COPY --from=codegolf/lang-tex          ["/", "/langs/tex/rootfs/"         ] # 9.67 MiB
COPY --from=codegolf/lang-hexagony     ["/", "/langs/hexagony/rootfs/"    ] # 8.90 MiB
COPY --from=codegolf/lang-php          ["/", "/langs/php/rootfs/"         ] # 8.40 MiB
COPY --from=codegolf/lang-scheme       ["/", "/langs/scheme/rootfs/"      ] # 7.38 MiB
COPY --from=codegolf/lang-tcl          ["/", "/langs/tcl/rootfs/"         ] # 5.68 MiB
COPY --from=codegolf/lang-perl         ["/", "/langs/perl/rootfs/"        ] # 5.54 MiB
COPY --from=codegolf/lang-arturo       ["/", "/langs/arturo/rootfs/"      ] # 5.29 MiB
COPY --from=codegolf/lang-fish         ["/", "/langs/fish/rootfs/"        ] # 4.85 MiB
COPY --from=codegolf/lang-algol-68     ["/", "/langs/algol-68/rootfs/"    ] # 4.66 MiB
COPY --from=codegolf/lang-cobol        ["/", "/langs/cobol/rootfs/"       ] # 4.61 MiB
COPY --from=codegolf/lang-rockstar-2   ["/", "/langs/rockstar-2/rootfs/"  ] # 4.26 MiB
COPY --from=codegolf/lang-squirrel     ["/", "/langs/squirrel/rootfs/"    ] # 4.08 MiB
COPY --from=codegolf/lang-befunge      ["/", "/langs/befunge/rootfs/"     ] # 3.95 MiB
COPY --from=codegolf/lang-hush         ["/", "/langs/hush/rootfs/"        ] # 3.32 MiB
COPY --from=codegolf/lang-rexx         ["/", "/langs/rexx/rootfs/"        ] # 3.06 MiB
COPY --from=codegolf/lang-forth        ["/", "/langs/forth/rootfs/"       ] # 2.86 MiB
COPY --from=codegolf/lang-awk          ["/", "/langs/awk/rootfs/"         ] # 1.80 MiB
COPY --from=codegolf/lang-c            ["/", "/langs/c/rootfs/"           ] # 1.73 MiB
COPY --from=codegolf/lang-bqn          ["/", "/langs/bqn/rootfs/"         ] # 1.26 MiB
COPY --from=codegolf/lang-rebol        ["/", "/langs/rebol/rootfs/"       ] # 1.23 MiB
COPY --from=codegolf/lang-sql          ["/", "/langs/sql/rootfs/"         ] # 1.21 MiB
COPY --from=codegolf/lang-bash         ["/", "/langs/bash/rootfs/"        ] # 1.20 MiB
COPY --from=codegolf/lang-berry        ["/", "/langs/berry/rootfs/"       ] #  986 KiB
COPY --from=codegolf/lang-jq           ["/", "/langs/jq/rootfs/"          ] #  973 KiB
COPY --from=codegolf/lang-janet        ["/", "/langs/janet/rootfs/"       ] #  902 KiB
COPY --from=codegolf/lang-fennel       ["/", "/langs/fennel/rootfs/"      ] #  673 KiB
COPY --from=codegolf/lang-k            ["/", "/langs/k/rootfs/"           ] #  631 KiB
COPY --from=codegolf/lang-wren         ["/", "/langs/wren/rootfs/"        ] #  505 KiB
COPY --from=codegolf/lang-lua          ["/", "/langs/lua/rootfs/"         ] #  362 KiB
COPY --from=codegolf/lang-sed          ["/", "/langs/sed/rootfs/"         ] #  244 KiB
COPY --from=codegolf/lang-brainfuck    ["/", "/langs/brainfuck/rootfs/"   ] # 51.1 KiB

COPY run-lang.c ./

RUN gcc -Wall -Werror -Wextra -o /usr/bin/run-lang -s -static run-lang.c

# reflex reruns a command when files change.
CMD reflex -sd none -r '\.(css|go|html|json|pem|svg|toml|txt)$' -R '_test\.go$' -- go run .
