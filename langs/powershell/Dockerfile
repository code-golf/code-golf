FROM mcr.microsoft.com/dotnet/core/sdk:3.1.300-alpine3.11 as builder

ENV DOTNET_CLI_TELEMETRY_OPTOUT=1

WORKDIR /source
COPY Interpreter.csproj Interpreter.cs /source/

ENV VERSION=7.0.1

RUN /usr/bin/dotnet add package System.Management.Automation -v $VERSION
RUN /usr/bin/dotnet add package Microsoft.PowerShell.Commands.Diagnostics -v $VERSION
RUN /usr/bin/dotnet add package Microsoft.PowerShell.Commands.Management -v $VERSION
RUN /usr/bin/dotnet add package Microsoft.PowerShell.Commands.Utility -v $VERSION
RUN /usr/bin/dotnet add package Microsoft.PowerShell.ConsoleHost -v $VERSION
RUN /usr/bin/dotnet add package Microsoft.WSMan.Management -v $VERSION

RUN /usr/bin/dotnet publish --runtime alpine.3.11-x64 -c Release -o /interpreter

RUN mkdir /empty

RUN echo nobody:x:99:99::/: > /etc/passwd

FROM scratch

COPY --from=0 /bin/sh                  /bin/
COPY --from=0 /etc/passwd              /etc/
COPY --from=0 /lib/ld-musl-x86_64.so.1 \
              /lib/libcrypto.so.1.1    \
              /lib/libssl.so.1.1       /lib/
COPY --from=0 /usr/lib                 /usr/lib/
COPY --from=0 /interpreter             /interpreter/
COPY --from=0 /empty                   /tmp/

ENTRYPOINT /interpreter/Interpreter --version
