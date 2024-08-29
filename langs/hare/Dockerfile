FROM alpine:3.20 AS builder

RUN apk add --no-cache build-base curl qbe scdoc

RUN curl https://git.sr.ht/~sircmpwn/harec/archive/0.24.2.tar.gz | tar xz

RUN cd harec-0.24.2               \
 && cp configs/linux.mk config.mk \
 && make -j `nproc`               \
 && make install

RUN curl https://git.sr.ht/~sircmpwn/hare/archive/0.24.2.tar.gz | tar xz

RUN cd hare-0.24.2                \
 && cp configs/linux.mk config.mk \
 && make -j `nproc`               \
 && make install

FROM codegolf/lang-base

COPY --from=0 /bin /bin
COPY --from=0 /lib /lib
COPY --from=0 /usr /usr

COPY /harewrapper /usr/bin/
RUN chmod +x /usr/bin/harewrapper

ENTRYPOINT ["harewrapper"]

CMD ["-v"]
