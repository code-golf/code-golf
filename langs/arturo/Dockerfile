FROM nimlang/nim:2.2.2-alpine AS builder

RUN apk add --no-cache pcre

ENV VER=3454c16

RUN curl -#L https://github.com/arturo-lang/arturo/tarball/$VER \
  | tar xz --strip-components 1

# Drop an "l" to make Nim int -> C long int. Both are 64 bit, not sure why
# this is needed on some systems and not others, seems to affect alpine more.
# See https://github.com/arturo-lang/arturo/issues/1896
# And https://gcc.gnu.org/onlinedocs/gcc/Integer-Overflow-Builtins.html
# $ echo 'echo sizeof(int)' | docker run --rm -i nimlang/nim:2.2.2-alpine nim c -r -
# 8
RUN sed -i -E 's/(__builtin_s)(add|sub|mul)l(l_overflow)/\1\2\3/g' \
    src/helpers/intrinsics.nim                                     \
 && ./build.nims --install --mode mini --release

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /root/.arturo/bin/arturo /usr/bin/
COPY --from=0 /usr/lib/libpcre.so.1    /usr/lib/

ENTRYPOINT ["arturo"]

CMD ["--version"]
