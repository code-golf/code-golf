FROM mcr.microsoft.com/dotnet/core/sdk:3.1.401-alpine3.12 as builder

RUN mkdir /empty

ENV DOTNET_CLI_TELEMETRY_OPTOUT=1 VERSION=37.0.0

WORKDIR /source

COPY Compiler.fsproj Compiler.fs ./

RUN dotnet add package FSharp.Compiler.Service -v $VERSION \
 && dotnet publish -c Release -r linux-musl-x64 -o /compiler

# Trim additional things to improve F# compiler startup performance.
WORKDIR /trimmer

COPY Trimmer.csproj Trimmer.cs ExtraTrimmingList.txt ./

RUN dotnet run

FROM scratch

COPY --from=0 /lib/ld-musl-x86_64.so.1 \
              /lib/libcrypto.so.1.1    \
              /lib/libssl.so.1.1       /lib/
COPY --from=0 /empty                   /proc
COPY --from=0 /empty                   /tmp
COPY --from=0 /usr/lib                 /usr/lib/
COPY --from=0 /compiler                /compiler/

ENTRYPOINT ["/compiler/Compiler", "--version"]
