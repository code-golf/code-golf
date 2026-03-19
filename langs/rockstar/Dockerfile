FROM mcr.microsoft.com/dotnet/sdk:10.0.201-alpine3.23-amd64 AS builder

RUN apk add --no-cache gcc musl-dev zlib-dev

ENV VER=2.0.31

WORKDIR /rockstar

ADD https://github.com/RockstarLang/rockstar.git#v$VER .

WORKDIR /rockstar/Starship

COPY rockstar.csproj Rockstar/Rockstar.csproj

COPY rockstar.engine.csproj Rockstar.Engine/

RUN dotnet publish Rockstar -c Release -o /out

COPY rockstar.c .

RUN gcc -Wall -Werror -Wextra -o /usr/bin/rockstar -s rockstar.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1 /lib/
COPY --from=0 /usr/bin/rockstar   /usr/bin/
COPY --from=0 /out/rockstar       /usr/local/bin/

ENV PATH=/usr/bin

ENTRYPOINT ["rockstar"]

CMD ["--version"]
