FROM mcr.microsoft.com/dotnet/sdk:7.0.201-alpine3.17 as builder

WORKDIR /work

COPY c-sharp.cs c-sharp.csproj ./

RUN dotnet publish -c Release -o /out

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1  \
              /lib/libcrypto.so.3       \
              /lib/libssl.so.3          /lib/
COPY --from=0 /out                      /usr/bin
COPY --from=0 /usr/lib/libgcc_s.so.1    \
              /usr/lib/libicudata.so.72 \
              /usr/lib/libicui18n.so.72 \
              /usr/lib/libicuio.so.72   \
              /usr/lib/libicuuc.so.72   \
              /usr/lib/libstdc++.so.6   /usr/lib/
COPY --from=0 /usr/share/icu            /usr/share/icu

ENTRYPOINT ["c-sharp"]

CMD ["--version"]
