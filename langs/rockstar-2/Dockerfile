FROM mcr.microsoft.com/dotnet/sdk:8.0.401-alpine3.20 AS builder

RUN apk add --no-cache g++ zlib-dev

RUN curl -L https://github.com/RockstarLang/rockstar2/archive/refs/tags/v0.9.1.tar.gz | tar xz

RUN dotnet publish /rockstar2*/Starship/Rockstar -c Release -o out

FROM codegolf/lang-base

COPY --from=0 /bin/chmod /bin/sh       /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /out                     /usr/bin

COPY /rockstarwrapper /usr/bin/
RUN chmod +x /usr/bin/rockstarwrapper

ENTRYPOINT ["rockstarwrapper"]

CMD ["--version"]
