FROM mcr.microsoft.com/dotnet/sdk:10.0.403-alpine3.24-amd64 AS builder

RUN apk add --no-cache gcc musl-dev patch zlib-dev

ENV VER=2.0.31

WORKDIR /rockstar

ADD https://github.com/RockstarLang/rockstar.git#v$VER .

COPY stdin.patch .

RUN patch -p0 < stdin.patch

RUN find Starship -name '*.csproj' -exec sed -Ei 's/\bnet\d+\b/net10/' {} + \
 && dotnet publish Starship/Rockstar -c Release -o /usr/bin \
    -p:InvariantGlobalization=true -p:StaticExecutable=true

FROM codegolf/lang-base

COPY --from=0 /usr/bin/rockstar /usr/bin/

ENTRYPOINT ["rockstar"]

CMD ["--version"]
