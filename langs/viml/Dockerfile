FROM alpine:3.14 as builder

RUN mkdir /empty

RUN apk add --no-cache build-base git ncurses-dev

RUN git clone --branch v8.2.3131 --depth 1 https://github.com/vim/vim.git

RUN cd vim/src \
 && ./configure --enable-multibyte --prefix=/usr --with-features=normal \
 && make install

# Remove the docs.
RUN rm -r /usr/share/vim/vim82/doc

# Quiet warnings about missing ftdetect.
RUN mkdir /usr/share/vim/vim82/ftdetect \
 && touch /usr/share/vim/vim82/ftdetect/vim.vim

FROM scratch

COPY --from=0 /bin                      /bin
COPY --from=0 /lib/ld-musl-x86_64.so.1  /lib/
COPY --from=0 /usr/lib/libncurses.so    \
              /usr/lib/libncursesw.so.6 /usr/lib/
COPY --from=0 /usr/share/vim/vim82      /usr/share/vim/vim82
COPY --from=0 /empty                    /proc
COPY --from=0 /empty                    /tmp
COPY --from=0 /usr/bin/vim              /usr/bin/

COPY viml /usr/bin/

ENTRYPOINT ["/usr/bin/viml"]

CMD ["--version"]
