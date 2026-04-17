FROM mcr.microsoft.com/dotnet/sdk:10.0.202-alpine3.23-amd64 AS builder

RUN apk add --no-cache gcc musl-dev zlib-dev

ENV VER=770406a

WORKDIR /hexagony

RUN curl -#L https://github.com/SirBogman/Hexagony/tarball/$VER \
  | tar xz --strip-components 1

COPY hexagony.csproj Hexagony.csproj

RUN dotnet publish -c Release -r linux-musl-x64

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1     /lib/
COPY --from=0 /out/Hexagony           /usr/bin/hexagony
COPY --from=0 /usr/lib/libgcc_s.so.1  \
              /usr/lib/libstdc++.so.6 /usr/lib/

ENTRYPOINT ["hexagony"]
