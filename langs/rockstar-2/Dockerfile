FROM mcr.microsoft.com/dotnet/sdk:8.0.401-alpine3.20 AS builder

RUN apk add --no-cache build-base zlib-dev

WORKDIR /work

COPY /Rockstar        /Rockstar
COPY /Rockstar.Engine /Rockstar.Engine

RUN dotnet publish /Rockstar -c Release -o /out

FROM codegolf/lang-base

COPY --from=0 / /

COPY --from=0 /out /usr/bin

COPY /rockstarwrapper /usr/bin/
RUN chmod +x /usr/bin/rockstarwrapper

ENTRYPOINT ["rockstarwrapper"]

CMD ["--version"]
