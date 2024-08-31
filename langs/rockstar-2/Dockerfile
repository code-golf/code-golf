#FROM mcr.microsoft.com/dotnet/sdk:8.0.401-alpine3.20 AS builder
FROM alpine:3.20 AS builder

ENV VER=0.9.1

RUN apk add --no-cache build-base curl libc6-compat

RUN curl -L https://github.com/RockstarLang/rockstar2/releases/download/v$VER/rockstar-v$VER-linux-x64.tar.gz \
  | tar xz --strip-components 1

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/chmod /bin/sh /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1    \
              /lib/libgcompat.so.0        \
              /lib/libucontext.so.1       /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2 /lib64/
COPY --from=0 /rockstar                   /usr/bin/
COPY --from=0 /usr/lib/libobstack.so.1    /usr/lib/

COPY /rockstarwrapper /usr/bin/
RUN chmod +x /usr/bin/rockstarwrapper

ENTRYPOINT ["rockstarwrapper"]

CMD ["--version"]
