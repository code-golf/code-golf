FROM codegolf/lang-ruby

FROM codegolf/lang-haskell

FROM alpine:3.23 AS builder

RUN apk add --no-cache gcc musl-dev

ENV VER=1.1.1

ADD --chmod=755 https://golfscript.com/iogii/iogii-$VER /usr/local/bin/iogii

# Rewrite shebang to bypass /usr/bin/env and go directly to /usr/bin/ruby.
RUN sed -i '1c#!/usr/bin/ruby' /usr/local/bin/iogii

COPY iogii.c .

RUN gcc -Wall -Werror -Wextra -o /usr/bin/iogii -s iogii.c

FROM codegolf/lang-base

ENV GHC='x86_64-linux-ghc-9.14.1-d531'

COPY --from=0 /lib                                                              /lib
COPY --from=0 /usr/bin/ruby                                                     /usr/bin/
COPY --from=2 /usr/bin/gcc /usr/bin/iogii                                       /usr/bin/
COPY --from=1 /usr/lib                                                          /usr/lib
COPY --from=0 /usr/lib/ruby/*/bundled_gems.rb                                   \
              /usr/lib/ruby/*/did_you_mean.rb                                   \
              /usr/lib/ruby/*/error_highlight.rb                                \
              /usr/lib/ruby/*/monitor.rb                                        \
              /usr/lib/ruby/*/rubygems.rb                                       /usr/lib/ruby/4.0.0/
COPY --from=0 /usr/lib/ruby/*/did_you_mean                                      /usr/lib/ruby/4.0.0/did_you_mean
COPY --from=0 /usr/lib/ruby/*/error_highlight                                   /usr/lib/ruby/4.0.0/error_highlight
COPY --from=0 /usr/lib/ruby/*/rubygems/basic_specification.rb                   \
              /usr/lib/ruby/*/rubygems/defaults.rb                              \
              /usr/lib/ruby/*/rubygems/dependency.rb                            \
              /usr/lib/ruby/*/rubygems/deprecate.rb                             \
              /usr/lib/ruby/*/rubygems/errors.rb                                \
              /usr/lib/ruby/*/rubygems/exceptions.rb                            \
              /usr/lib/ruby/*/rubygems/path_support.rb                          \
              /usr/lib/ruby/*/rubygems/platform.rb                              \
              /usr/lib/ruby/*/rubygems/requirement.rb                           \
              /usr/lib/ruby/*/rubygems/specification.rb                         \
              /usr/lib/ruby/*/rubygems/specification_record.rb                  \
              /usr/lib/ruby/*/rubygems/stub_specification.rb                    \
              /usr/lib/ruby/*/rubygems/target_rbconfig.rb                       \
              /usr/lib/ruby/*/rubygems/unknown_command_spell_checker.rb         \
              /usr/lib/ruby/*/rubygems/util.rb                                  \
              /usr/lib/ruby/*/rubygems/version.rb                               \
              /usr/lib/ruby/*/rubygems/win_platform.rb                          /usr/lib/ruby/4.0.0/rubygems/
COPY --from=0 /usr/lib/ruby/*/rubygems/core_ext                                 /usr/lib/ruby/4.0.0/rubygems/core_ext
COPY --from=0 /usr/lib/ruby/*/rubygems/util                                     /usr/lib/ruby/4.0.0/rubygems/util
COPY --from=0 /usr/lib/ruby/*/syntax_suggest/core_ext.rb                        /usr/lib/ruby/4.0.0/syntax_suggest/
COPY --from=0 /usr/lib/ruby/*/x86_64-linux-musl/monitor.so                      \
              /usr/lib/ruby/*/x86_64-linux-musl/rbconfig.rb                     /usr/lib/ruby/4.0.0/x86_64-linux-musl/
COPY --from=0 /usr/lib/ruby/*/x86_64-linux-musl/enc/encdb.so                    /usr/lib/ruby/4.0.0/x86_64-linux-musl/enc/
COPY --from=0 /usr/lib/ruby/*/x86_64-linux-musl/enc/trans/transdb.so            /usr/lib/ruby/4.0.0/x86_64-linux-musl/enc/trans/
COPY --from=0 /usr/lib/ruby/gems/4.0.0/specifications/default                   /usr/lib/ruby/gems/4.0.0/specifications/default
COPY --from=1 /usr/local/bin                                                    /usr/local/bin
COPY --from=2 /usr/local/bin/iogii                                              /usr/local/bin/
COPY --from=1 /usr/local/lib/settings                                           /usr/local/lib/
COPY --from=1 /usr/local/lib/package.conf.d                                     /usr/local/lib/package.conf.d
COPY --from=1 /usr/local/lib/*/libHSarray-*.so                                  \
              /usr/local/lib/*/libHSbase-*.so                                   \
              /usr/local/lib/*/libHSbinary-*.so                                 \
              /usr/local/lib/*/libHSbytestring-*.so                             \
              /usr/local/lib/*/libHScontainers-*.so                             \
              /usr/local/lib/*/libHSdeepseq-*.so                                \
              /usr/local/lib/*/libHSdirectory-*.so                              \
              /usr/local/lib/*/libHSexceptions-*.so                             \
              /usr/local/lib/*/libHSfile-io-*.so                                \
              /usr/local/lib/*/libHSfilepath-*.so                               \
              /usr/local/lib/*/libHSghc-*.so                                    \
              /usr/local/lib/*/libHSghci-*.so                                   \
              /usr/local/lib/*/libHShaskeline-*.so                              \
              /usr/local/lib/*/libHShpc-*.so                                    \
              /usr/local/lib/*/libHSmtl-*.so                                    \
              /usr/local/lib/*/libHSos-string-*.so                              \
              /usr/local/lib/*/libHSpretty-*.so                                 \
              /usr/local/lib/*/libHSprocess-*.so                                \
              /usr/local/lib/*/libHSrts-*_thr-*.so                              \
              /usr/local/lib/*/libHSsemaphore-compat-*.so                       \
              /usr/local/lib/*/libHSstm-*.so                                    \
              /usr/local/lib/*/libHStemplate-haskell-*.so                       \
              /usr/local/lib/*/libHSterminfo-*.so                               \
              /usr/local/lib/*/libHStime-*.so                                   \
              /usr/local/lib/*/libHStransformers-*.so                           \
              /usr/local/lib/*/libHSunix-*.so                                   /usr/local/lib/$GHC/
COPY --from=1 /usr/local/lib/*/base-*/Prelude.dyn_hi                            /usr/local/lib/$GHC/base-4.22.0.0-d466/
COPY --from=1 /usr/local/lib/*/base-*/Control/Monad.dyn_hi                      /usr/local/lib/$GHC/base-4.22.0.0-d466/Control/
COPY --from=1 /usr/local/lib/*/base-*/Data/Char.dyn_hi                          \
              /usr/local/lib/*/base-*/Data/List.dyn_hi                          \
              /usr/local/lib/*/base-*/Data/Maybe.dyn_hi                         \
              /usr/local/lib/*/base-*/Data/String.dyn_hi                        /usr/local/lib/$GHC/base-4.22.0.0-d466/Data/
COPY --from=1 /usr/local/lib/*/base-*/System/Environment.dyn_hi                 \
              /usr/local/lib/*/base-*/System/IO.dyn_hi                          \
              /usr/local/lib/*/base-*/System/Info.dyn_hi                        /usr/local/lib/$GHC/base-4.22.0.0-d466/System/
COPY --from=1 /usr/local/lib/*/containers-*/Data/Map.dyn_hi                     \
              /usr/local/lib/*/containers-*/Data/Set.dyn_hi                     /usr/local/lib/$GHC/containers-0.8-3e81/Data/
COPY --from=1 /usr/local/lib/*/containers-*/Data/Map                            /usr/local/lib/$GHC/containers-0.8-3e81/Data/Map
COPY --from=1 /usr/local/lib/*/containers-*/Data/Set                            /usr/local/lib/$GHC/containers-0.8-3e81/Data/Set
COPY --from=1 /usr/local/lib/*/ghc-internal-*/GHC/Internal/Base.dyn_hi          \
              /usr/local/lib/*/ghc-internal-*/GHC/Internal/CString.dyn_hi       \
              /usr/local/lib/*/ghc-internal-*/GHC/Internal/Char.dyn_hi          \
              /usr/local/lib/*/ghc-internal-*/GHC/Internal/Classes.dyn_hi       \
              /usr/local/lib/*/ghc-internal-*/GHC/Internal/Enum.dyn_hi          \
              /usr/local/lib/*/ghc-internal-*/GHC/Internal/Err.dyn_hi           \
              /usr/local/lib/*/ghc-internal-*/GHC/Internal/Float.dyn_hi         \
              /usr/local/lib/*/ghc-internal-*/GHC/Internal/GHCi.dyn_hi          \
              /usr/local/lib/*/ghc-internal-*/GHC/Internal/List.dyn_hi          \
              /usr/local/lib/*/ghc-internal-*/GHC/Internal/Maybe.dyn_hi         \
              /usr/local/lib/*/ghc-internal-*/GHC/Internal/Num.dyn_hi           \
              /usr/local/lib/*/ghc-internal-*/GHC/Internal/Prim.dyn_hi          \
              /usr/local/lib/*/ghc-internal-*/GHC/Internal/Read.dyn_hi          \
              /usr/local/lib/*/ghc-internal-*/GHC/Internal/Real.dyn_hi          \
              /usr/local/lib/*/ghc-internal-*/GHC/Internal/Show.dyn_hi          \
              /usr/local/lib/*/ghc-internal-*/GHC/Internal/TopHandler.dyn_hi    \
              /usr/local/lib/*/ghc-internal-*/GHC/Internal/Types.dyn_hi         \
              /usr/local/lib/*/ghc-internal-*/GHC/Internal/Unicode.dyn_hi       /usr/local/lib/$GHC/ghc-internal-9.1401.0-ec29/GHC/Internal/
COPY --from=1 /usr/local/lib/*/ghc-internal-*/GHC/Internal/Bignum               /usr/local/lib/$GHC/ghc-internal-9.1401.0-ec29/GHC/Internal/Bignum
COPY --from=1 /usr/local/lib/*/ghc-internal-*/GHC/Internal/Data/Foldable.dyn_hi \
              /usr/local/lib/*/ghc-internal-*/GHC/Internal/Data/Functor.dyn_hi  \
              /usr/local/lib/*/ghc-internal-*/GHC/Internal/Data/Maybe.dyn_hi    \
              /usr/local/lib/*/ghc-internal-*/GHC/Internal/Data/NonEmpty.dyn_hi \
              /usr/local/lib/*/ghc-internal-*/GHC/Internal/Data/OldList.dyn_hi  \
              /usr/local/lib/*/ghc-internal-*/GHC/Internal/Data/Tuple.dyn_hi    \
              /usr/local/lib/*/ghc-internal-*/GHC/Internal/Data/Version.dyn_hi  /usr/local/lib/$GHC/ghc-internal-9.1401.0-ec29/GHC/Internal/Data/
COPY --from=1 /usr/local/lib/*/ghc-internal-*/GHC/Internal/Exception            /usr/local/lib/$GHC/ghc-internal-9.1401.0-ec29/GHC/Internal/Exception
COPY --from=1 /usr/local/lib/*/ghc-internal-*/GHC/Internal/GHCi                 /usr/local/lib/$GHC/ghc-internal-9.1401.0-ec29/GHC/Internal/GHCi
COPY --from=1 /usr/local/lib/*/ghc-internal-*/GHC/Internal/Prim                 /usr/local/lib/$GHC/ghc-internal-9.1401.0-ec29/GHC/Internal/Prim
COPY --from=1 /usr/local/lib/*/ghc-internal-*/GHC/Internal/Stack                /usr/local/lib/$GHC/ghc-internal-9.1401.0-ec29/GHC/Internal/Stack
COPY --from=1 /usr/local/lib/*/ghc-internal-*/GHC/Internal/System               /usr/local/lib/$GHC/ghc-internal-9.1401.0-ec29/GHC/Internal/System
COPY --from=1 /usr/local/lib/*/ghc-internal-*/GHC/Internal/Text                 /usr/local/lib/$GHC/ghc-internal-9.1401.0-ec29/GHC/Internal/Text
COPY --from=1 /usr/local/lib/*/ghc-internal-*/GHC/Internal/Unsafe               /usr/local/lib/$GHC/ghc-internal-9.1401.0-ec29/GHC/Internal/Unsafe
COPY --from=1 /usr/local/lib/*/rts-*/include/DerivedConstants.h                 /usr/local/lib/$GHC/rts-1.0.3/include/

ENTRYPOINT ["iogii"]

CMD ["--version"]
