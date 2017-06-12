FROM alpine:edge

RUN apk add --no-cache curl gcc linux-headers make musl-dev perl

RUN curl http://rakudo.org/downloads/star/rakudo-star-2017.04.tar.gz \
  | tar xzf -

RUN cd rakudo-star-2017.04                                    \
 && perl Configure.pl --backend=moar --gen-moar --prefix=/usr \
 && make install

FROM scratch

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/moar            /usr/bin/
COPY --from=0 /usr/lib/libmoar.so      /usr/lib/
COPY --from=0 /usr/share/nqp           /usr/share/nqp
COPY --from=0 /usr/share/perl6         /usr/share/perl6

ENTRYPOINT ["/usr/bin/moar"]
