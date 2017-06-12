FROM alpine:edge

RUN apk add --no-cache curl gcc make musl-dev

RUN curl http://www.cpan.org/src/5.0/perl-5.26.0.tar.xz | tar xJf -

RUN cd perl-5.26.0                              \
 && ./Configure                                 \
    -Accflags='-DNO_MATHOMS -DPERL_DISABLE_PMC' \
    -Dprefix=/usr                               \
    -des                                        \
 && make                                        \
 && make install-strip

RUN find /usr/lib/perl5 -regex '.*\.\(h\|pod\)' -delete

FROM scratch

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/perl            /usr/bin/
COPY --from=0 /usr/lib/perl5           /usr/lib/perl5

ENTRYPOINT ["/usr/bin/perl"]
