FROM mcr.microsoft.com/dotnet/core/sdk:3.1.300-alpine3.11 as builder

ENV NUGET_VERSION=36.0.1
ENV DOTNET_CLI_TELEMETRY_OPTOUT=1

WORKDIR /source/
COPY Compiler.fsproj Compiler.fs /source/

RUN /usr/bin/dotnet add package FSharp.Compiler.Service -v $NUGET_VERSION
RUN /usr/bin/dotnet publish --runtime alpine.3.11-x64 -c Release -o /compiler

# Trim additional things to improve F# compiler startup performance.
WORKDIR /trimmer/
COPY Trimmer.csproj Trimmer.cs ExtraTrimmingList.txt /trimmer/
RUN /usr/bin/dotnet run

RUN mkdir /empty

FROM scratch
COPY --from=0 /bin/sh                  /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1 \
              /lib/libcrypto.so.1.1    \
              /lib/libssl.so.1.1       /lib/
COPY --from=0 /usr/lib                 /usr/lib/
COPY --from=0 /compiler                /compiler/
COPY --from=0 /empty                   /tmp/

ENTRYPOINT /compiler/Compiler --version
