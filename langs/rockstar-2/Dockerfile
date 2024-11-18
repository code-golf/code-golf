FROM mcr.microsoft.com/dotnet/sdk:9.0.100-alpine3.20 AS builder

RUN apk add --no-cache build-base zlib-dev

RUN curl -L https://github.com/RockstarLang/rockstar/archive/refs/tags/v2.0.0-beta.2.0.0.tar.gz \
  | tar xz --strip-components 1

RUN dotnet publish Starship/Rockstar -c Release -o rockstar

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/sh         /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /rockstar/rockstar       /usr/bin/

COPY /rockstarwrapper /usr/bin/

ENTRYPOINT ["rockstarwrapper"]

CMD ["--version"]
