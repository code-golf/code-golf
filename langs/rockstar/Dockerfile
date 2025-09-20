FROM mcr.microsoft.com/dotnet/sdk:9.0.305-alpine3.22-amd64 AS builder

RUN apk add --no-cache gcc musl-dev zlib-dev

ENV VER=2.0.31

RUN curl -#L https://github.com/RockstarLang/rockstar/archive/refs/tags/v$VER.tar.gz \
  | tar xz --strip-components 1

RUN find Starship -name '*.csproj' -exec sed -Ei 's/\bnet\d+\b/net9/' {} + \
 && dotnet publish Starship/Rockstar -c Release -o /usr/bin \
    -p:InvariantGlobalization=true

COPY rockstar.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/rockstarwrapper -s -static rockstar.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/rockstar        \
              /usr/bin/rockstarwrapper /usr/bin/

ENTRYPOINT ["rockstarwrapper"]

CMD ["--version"]
