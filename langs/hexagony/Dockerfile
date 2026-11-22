FROM mcr.microsoft.com/dotnet/sdk:11.0.100-alpine3.24-amd64 AS builder

RUN apk add --no-cache gcc musl-dev zlib-dev

ENV VER=770406a

WORKDIR /hexagony

ADD https://github.com/SirBogman/Hexagony.git?commit=$VER .

RUN sed -Ei 's/\bnet\d+\b/net11/' Hexagony.csproj \
 && dotnet publish -c Release -o /out             \
    -p:InvariantGlobalization=true -p:PublishAot=true -p:StaticExecutable=true

FROM codegolf/lang-base

COPY --from=0 /out/Hexagony /usr/bin/hexagony

ENTRYPOINT ["hexagony"]
