FROM alpine:3.23 AS builder

RUN apk add --no-cache dash gcc make musl-dev ncurses-dev ncurses-static

ENV LDFLAGS='-static' VER=9.2.0010

ADD https://github.com/vim/vim.git#v$VER /vim

WORKDIR /vim/src

RUN ./configure            \
    --enable-multibyte     \
    --prefix=/usr          \
    --with-features=normal \
 && make -j`nproc` install

# Quiet warnings about missing ftdetect.
RUN mkdir /usr/share/vim/vim92/ftdetect \
 && touch /usr/share/vim/vim92/ftdetect/vim.vim

COPY viml.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/viml -s -static /viml.c

FROM codegolf/lang-base

COPY --from=0           /usr/bin/dash        /bin/sh
COPY --from=0 --parents /lib/ld-musl-*.so.1  \
                        /usr/bin/vim         \
                        /usr/bin/viml        \
                        /usr/share/vim/vim92 /

ENTRYPOINT ["viml"]

CMD ["--version"]
