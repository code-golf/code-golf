FROM crystallang/crystal:1.13.2-alpine AS builder

ENV VER=1.13.2

RUN apk add --no-cache build-base curl \
                       gc gc-dev \
                       libevent libevent-dev libevent-static \
                       libstdc++ libstdc++-dev \
                       llvm18 llvm18-dev llvm18-libs llvm18-linker-tools llvm18-static \
                       ncurses ncurses-dev ncurses-libs ncurses-static \
                       pcre2 pcre2-dev pcre2-tools

RUN curl -L https://github.com/crystal-lang/crystal/archive/refs/tags/$VER.tar.gz \
  | tar xz

RUN cd crystal-$VER \
 && LDFLAGS="-lstdc++ -lgc -lpthread -ldl -levent -lffi -lz -ltinfo -lm" \
 make -j`nproc` crystal interpreter=1 release=1 static=1

# DO IT!!!

FROM codegolf/lang-base

COPY --from=0 /usr/bin/crystal /usr/bin/

ENTRYPOINT ["crystal"]

CMD ["-v"]
