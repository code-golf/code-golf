FROM mcr.microsoft.com/dotnet/sdk:5.0.301-alpine3.13 as builder

ENV DOTNET_CLI_TELEMETRY_OPTOUT=1 VERSION=ac545d2

RUN mkdir /empty

RUN curl -L https://github.com/SirBogman/Hexagony/tarball/$VERSION | tar xz \
 && cd /SirBogman-Hexagony-$VERSION                                         \
 && VERSION= dotnet publish -c Release -r linux-musl-x64 -o /hexagony       \
    -p:InvariantGlobalization=true -p:PublishReadyToRun=true -p:PublishTrimmed=true

FROM scratch

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /empty                   /proc
COPY --from=0 /empty                   /tmp
COPY --from=0 /usr/lib/libgcc_s.so.1   \
              /usr/lib/libstdc++.so.6  /usr/lib/
COPY --from=0 /hexagony                /hexagony/

ENTRYPOINT ["/hexagony/Hexagony"]
