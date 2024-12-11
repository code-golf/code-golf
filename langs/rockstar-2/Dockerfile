FROM mcr.microsoft.com/dotnet/sdk:9.0.100-alpine3.20 AS builder

RUN apk add --no-cache build-base zlib-dev

RUN curl -L https://github.com/RockstarLang/rockstar/archive/refs/tags/v2.0.16.tar.gz \
  | tar xz --strip-components 1

RUN find Starship -name '*.csproj' -exec sed -Ei 's/\bnet\d+\b/net9/' {} + \
 && dotnet publish Starship/Rockstar -c Release -o rockstar \
    -p:InvariantGlobalization=true

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/sh         /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /rockstar/rockstar       /usr/bin/

COPY /rockstarwrapper /usr/bin/

ENTRYPOINT ["rockstarwrapper"]

CMD ["--version"]
