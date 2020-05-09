FROM mcr.microsoft.com/dotnet/core/sdk:3.1.201-alpine3.11 as builder

ENV DOTNET_CLI_TELEMETRY_OPTOUT=1

WORKDIR /source
COPY Compiler.csproj Compiler.cs /source/

# Note: the ReadyToRun optimization process failed for earlier versions
# of the CSharp package.
RUN /usr/bin/dotnet add package Microsoft.CodeAnalysis.CSharp -v 3.6.0-3.final

RUN /usr/bin/dotnet publish --runtime alpine.3.11-x64 -c Release -o /compiler

RUN mkdir /empty

FROM scratch

COPY --from=0 /bin/ln                  \
              /bin/sh                  /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1 \
              /lib/libcrypto.so.1.1    \
              /lib/libssl.so.1.1       /lib/
COPY --from=0 /usr/lib                 /usr/lib/
COPY --from=0 /compiler                /compiler/
COPY --from=0 /empty                   /tmp/
COPY --from=0 /empty                   /usr/bin/

RUN ln -s /compiler/Compiler /usr/bin/c-sharp

ENTRYPOINT /usr/bin/c-sharp --version
