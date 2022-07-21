FROM mcr.microsoft.com/dotnet/sdk:6.0.302-alpine3.15 as builder

ENV DOTNET_CLI_TELEMETRY_OPTOUT=1 VERSION=41.0.5

WORKDIR /source

COPY Compiler.fsproj Compiler.fs ./

RUN dotnet add package FSharp.Compiler.Service -v $VERSION \
 && dotnet publish -c Release -r linux-musl-x64 --self-contained -o /compiler

# Trim additional things to improve F# compiler startup performance.
WORKDIR /trimmer

COPY Trimmer.csproj Trimmer.cs ExtraTrimmingList.txt ./

RUN dotnet run

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 \
              /lib/libcrypto.so.1.1    \
              /lib/libssl.so.1.1       /lib/
COPY --from=0 /usr/lib                 /usr/lib/
COPY --from=0 /compiler                /compiler/

ENTRYPOINT ["/compiler/Compiler"]

CMD ["--version"]
