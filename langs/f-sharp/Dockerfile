FROM mcr.microsoft.com/dotnet/core/sdk:3.1.201-alpine3.11 as builder

ENV SDK_VERSION=3.1.201
ENV DOTNET_CLI_TELEMETRY_OPTOUT=1

WORKDIR /source/
COPY Compiler.fsproj Compiler.fs /source/
RUN /usr/bin/dotnet add package FSharp.Compiler.Service -v 35.0.0
RUN /usr/bin/dotnet publish --runtime alpine.3.11-x64 -c Release -o /compiler

# Trim additional things to improve F# compiler startup performance.
WORKDIR /trimmer/
COPY Trimmer.csproj Trimmer.cs ExtraTrimmingList.txt /trimmer/
RUN /usr/bin/dotnet run

RUN mkdir /empty

RUN /usr/bin/dotnet exec /usr/share/dotnet/sdk/$SDK_VERSION/FSharp/fsc.exe --help | head -1 > /compiler/version.txt

FROM scratch
COPY --from=0 /bin/sh                  /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1 \
              /lib/libcrypto.so.1.1    \
              /lib/libssl.so.1.1       /lib/
COPY --from=0 /usr/lib                 /usr/lib/
COPY --from=0 /compiler                /compiler/
COPY --from=0 /empty                   /tmp/
COPY --from=0 /empty                   /usr/bin/

ENTRYPOINT /compiler/Compiler --version
