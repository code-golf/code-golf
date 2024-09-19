FROM mcr.microsoft.com/dotnet/sdk:8.0.401-alpine3.20 AS builder

RUN apk add --no-cache build-base zlib-dev

RUN curl -L https://github.com/RockstarLang/rockstar2/archive/refs/tags/v0.11.0.tar.gz \
  | tar xz --strip-components 1

# Remove 9-10
COPY /argv.patch .
RUN patch -p 0 < argv.patch

RUN dotnet publish Starship/Rockstar -c Release -o rockstar

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/sh         /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /rockstar/rockstar       /usr/bin/

COPY /rockstarwrapper /usr/bin/

# Remove 23-24
COPY --from=0 /bin /bin
RUN chmod +x /usr/bin/rockstarwrapper

ENTRYPOINT ["rockstarwrapper"]

CMD ["--version"]
