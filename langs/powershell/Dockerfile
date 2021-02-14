FROM mcr.microsoft.com/dotnet/sdk:5.0.102-alpine3.12 as builder

RUN mkdir /empty

ENV DOTNET_CLI_TELEMETRY_OPTOUT=1 VERSION=7.1.2

WORKDIR /source

COPY Interpreter.csproj Interpreter.cs ./

RUN dotnet add package System.Management.Automation              -v $VERSION \
 && dotnet add package Microsoft.PowerShell.Commands.Diagnostics -v $VERSION \
 && dotnet add package Microsoft.PowerShell.Commands.Management  -v $VERSION \
 && dotnet add package Microsoft.PowerShell.Commands.Utility     -v $VERSION \
 && dotnet add package Microsoft.PowerShell.ConsoleHost          -v $VERSION \
 && dotnet add package Microsoft.WSMan.Management                -v $VERSION \
 && dotnet publish -c Release -r linux-musl-x64 -o /interpreter

RUN echo nobody:x:99:99::/: > /etc/passwd

FROM scratch

COPY --from=0 /etc/passwd              /etc/
COPY --from=0 /lib/ld-musl-x86_64.so.1 \
              /lib/libcrypto.so.1.1    \
              /lib/libssl.so.1.1       /lib/
COPY --from=0 /empty                   /proc
COPY --from=0 /empty                   /tmp
COPY --from=0 /usr/lib                 /usr/lib/
COPY --from=0 /interpreter             /interpreter/

ENTRYPOINT ["/interpreter/Interpreter", "--version"]
