FROM alpine:3.16 as builder

RUN apk add --no-cache build-base curl

RUN curl https://ftp.gnu.org/gnu/gawk/gawk-5.2.0.tar.xz | tar xJ

RUN mv gawk-* gawk        \
 && cd gawk               \
 && ./configure           \
   --disable-extension    \
   --disable-lint         \
   LDFLAGS="-static"      \
 && make                  \
 && strip gawk

FROM codegolf/lang-base

COPY --from=0 /gawk/gawk /usr/bin/

ENTRYPOINT ["gawk"]

CMD ["--version"]
