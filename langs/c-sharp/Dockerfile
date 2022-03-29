FROM mcr.microsoft.com/dotnet/sdk:6.0.201-alpine3.15 as builder

RUN mkdir /empty

ENV DOTNET_CLI_TELEMETRY_OPTOUT=1 VERSION=4.2.0-1.final

WORKDIR /source

COPY Compiler.csproj Compiler.cs ./

RUN dotnet add package Microsoft.CodeAnalysis.CSharp -v $VERSION \
 && dotnet publish -c Release -r linux-musl-x64 --self-contained -o /compiler

FROM scratch

COPY --from=0 /empty                   /dev
COPY --from=0 /lib/ld-musl-x86_64.so.1 \
              /lib/libcrypto.so.1.1    \
              /lib/libssl.so.1.1       /lib/
COPY --from=0 /empty                   /proc
COPY --from=0 /empty                   /tmp
COPY --from=0 /usr/lib                 /usr/lib/
COPY --from=0 /compiler                /compiler/

ENTRYPOINT ["/compiler/Compiler"]

CMD ["--version"]
