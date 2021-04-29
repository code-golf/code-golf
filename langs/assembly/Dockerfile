FROM alpine:3.13 as builder

ENV VERSION=57b1109

RUN apk add --no-cache build-base curl

RUN curl -L https://github.com/NewDefectus/defAsm/tarball/$VERSION | tar xz \
 && mv NewDefectus-defAsm-$VERSION defAsm

FROM gcr.io/distroless/nodejs

COPY --from=0 /defAsm /defAsm
