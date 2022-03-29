FROM mcr.microsoft.com/dotnet/sdk:6.0.201-alpine3.15 as builder

# /etc/passwd is needed.
RUN mkdir /empty /myetc \
 && echo nobody:x:99:99::/: > /myetc/passwd

WORKDIR /source

COPY powershell.cs powershell.csproj ./

RUN dotnet publish -c Release -o /mybin

FROM scratch

COPY --from=0 /empty                   /dev
COPY --from=0 /myetc                   /etc
COPY --from=0 /lib/ld-musl-x86_64.so.1 \
              /lib/libcrypto.so.1.1    \
              /lib/libssl.so.1.1       /lib/
COPY --from=0 /empty                   /proc
COPY --from=0 /empty                   /tmp
COPY --from=0 /mybin                   /usr/bin
COPY --from=0 /usr/lib                 /usr/lib

ENTRYPOINT ["powershell"]

CMD ["--version"]
