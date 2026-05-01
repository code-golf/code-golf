FROM mcr.microsoft.com/dotnet/sdk:10.0.203-alpine3.23-amd64 AS builder

RUN apk add --no-cache gcc musl-dev zlib-dev

ENV VER=2.0.31

WORKDIR /rockstar

ADD https://github.com/RockstarLang/rockstar.git#v$VER .

RUN find Starship -name '*.csproj' -exec sed -Ei 's/\bnet\d+\b/net10/' {} + \
 && dotnet publish Starship/Rockstar -c Release -o /usr/local/bin \
    -p:InvariantGlobalization=true

COPY rockstar.c .

RUN gcc -Wall -Werror -Wextra -o /usr/bin/rockstar -s -static rockstar.c

FROM codegolf/lang-base

COPY --from=0 --parents /lib/ld-musl-*.so.1     \
                        /usr/bin/rockstar       \
                        /usr/local/bin/rockstar /

ENV PATH=/usr/bin

ENTRYPOINT ["rockstar"]

CMD ["--version"]
