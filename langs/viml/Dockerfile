FROM alpine:3.16 as builder

RUN apk add --no-cache build-base git ncurses-dev

RUN git clone --branch v9.0.0000 --depth 1 https://github.com/vim/vim.git

RUN cd vim/src \
 && ./configure --enable-multibyte --prefix=/usr --with-features=normal \
 && make install

# Remove the docs.
RUN rm -r /usr/share/vim/vim90/doc

# Quiet warnings about missing ftdetect.
RUN mkdir /usr/share/vim/vim90/ftdetect \
 && touch /usr/share/vim/vim90/ftdetect/vim.vim

FROM codegolf/lang-base

COPY --from=0 /bin                      /bin
COPY --from=0 /lib/ld-musl-x86_64.so.1  /lib/
COPY --from=0 /usr/lib/libncurses.so    \
              /usr/lib/libncursesw.so.6 /usr/lib/
COPY --from=0 /usr/share/vim/vim90      /usr/share/vim/vim90
COPY --from=0 /usr/bin/vim              /usr/bin/

COPY viml /usr/bin/

ENTRYPOINT ["viml"]

CMD ["--version"]
