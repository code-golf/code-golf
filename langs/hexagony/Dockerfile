FROM mcr.microsoft.com/dotnet/sdk:5.0.102-alpine3.12 as builder

ENV COMMIT=b531c38

RUN mkdir /empty

ENV DOTNET_CLI_TELEMETRY_OPTOUT=1

RUN curl -L https://github.com/SirBogman/Hexagony/tarball/$COMMIT \
  | tar xz

WORKDIR /SirBogman-Hexagony-$COMMIT

RUN dotnet publish -c Release -r linux-musl-x64 -o /hexagony \
    -p:PublishReadyToRun=true -p:PublishTrimmed=true

FROM scratch

COPY --from=0 /lib/ld-musl-x86_64.so.1 \
              /lib/libcrypto.so.1.1    \
              /lib/libssl.so.1.1       /lib/
COPY --from=0 /empty                   /proc
COPY --from=0 /empty                   /tmp
COPY --from=0 /usr/lib                 /usr/lib/
COPY --from=0 /hexagony                /hexagony/

ENTRYPOINT ["/hexagony/Hexagony"]
