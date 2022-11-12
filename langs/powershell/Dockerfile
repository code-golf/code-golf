FROM mcr.microsoft.com/dotnet/sdk:7.0.100-alpine3.16 as builder

WORKDIR /work

COPY powershell.cs powershell.csproj ./

RUN dotnet publish -c Release -o /out

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1  /lib/
COPY --from=0 /out                      /usr/bin
COPY --from=0 /usr/lib/libgcc_s.so.1    \
              /usr/lib/libicudata.so.71 \
              /usr/lib/libicui18n.so.71 \
              /usr/lib/libicuio.so.71   \
              /usr/lib/libicuuc.so.71   \
              /usr/lib/libstdc++.so.6   /usr/lib/
COPY --from=0 /usr/share/icu            /usr/share/icu

ENTRYPOINT ["powershell"]

CMD ["--version"]
