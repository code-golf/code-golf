FROM mcr.microsoft.com/dotnet/sdk:10.0.403-alpine3.24-amd64 AS builder

RUN apk add --no-cache gcc musl-dev zlib-dev

ENV VER=2.0.38

WORKDIR /rockstar

ADD https://github.com/RockstarLang/rockstar.git#v$VER .

RUN find Starship -name '*.csproj' -exec sed -Ei 's/\bnet\d+\b/net10/' {} + \
 && dotnet publish Starship/Rockstar -c Release -o /out                     \
    -p:InvariantGlobalization=true -p:PublishAot=true -p:StaticExecutable=true

FROM codegolf/lang-base

COPY --from=0 /out/rockstar /usr/bin/

ENTRYPOINT ["rockstar"]

CMD ["--version"]
