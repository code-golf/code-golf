FROM alpine:3.22 AS builder

RUN apk add --no-cache build-base dash git ncurses-dev

ENV VER=9.1.0447

RUN git clone --branch v$VER --depth 1 https://github.com/vim/vim.git

WORKDIR /vim/src

RUN ./configure            \
    --enable-multibyte     \
    --prefix=/usr          \
    --with-features=normal \
 && make -j`nproc` install

# Remove the docs.
RUN rm -r /usr/share/vim/vim91/doc

# Quiet warnings about missing ftdetect.
RUN mkdir /usr/share/vim/vim91/ftdetect \
 && touch /usr/share/vim/vim91/ftdetect/vim.vim

COPY viml.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/viml -s -static /viml.c

FROM codegolf/lang-base

COPY --from=0 /usr/bin/dash             /bin/sh
COPY --from=0 /lib/ld-musl-x86_64.so.1  /lib/
COPY --from=0 /usr/bin/vim              \
              /usr/bin/viml             /usr/bin/
COPY --from=0 /usr/lib/libncurses.so    \
              /usr/lib/libncursesw.so.6 /usr/lib/
COPY --from=0 /usr/share/vim/vim91      /usr/share/vim/vim91

ENTRYPOINT ["viml"]

CMD ["--version"]
