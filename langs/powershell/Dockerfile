FROM mcr.microsoft.com/dotnet/sdk:6.0.300-alpine3.15 as builder

WORKDIR /source

COPY powershell.cs powershell.csproj ./

RUN dotnet publish -c Release -o /mybin

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 \
              /lib/libcrypto.so.1.1    \
              /lib/libssl.so.1.1       /lib/
COPY --from=0 /mybin                   /usr/bin
COPY --from=0 /usr/lib                 /usr/lib

ENTRYPOINT ["powershell"]

CMD ["--version"]
