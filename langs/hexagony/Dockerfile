FROM mcr.microsoft.com/dotnet/sdk:10.0.400-alpine3.24-amd64 AS builder

RUN apk add --no-cache gcc musl-dev zlib-dev

ENV VER=770406a

WORKDIR /hexagony

ADD https://github.com/SirBogman/Hexagony.git?commit=$VER .

COPY hexagony.csproj Hexagony.csproj

RUN dotnet publish -c Release -r linux-musl-x64

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1     /lib/
COPY --from=0 /out/Hexagony           /usr/bin/hexagony
COPY --from=0 /usr/lib/libgcc_s.so.1  \
              /usr/lib/libstdc++.so.6 /usr/lib/

ENTRYPOINT ["hexagony"]
