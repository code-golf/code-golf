FROM mcr.microsoft.com/dotnet/sdk:9.0.301-alpine3.22-amd64 AS builder

WORKDIR /c-sharp

COPY c-sharp.cs c-sharp.csproj ./

RUN dotnet publish -c Release -o /out

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1    /lib/
COPY --from=0 /out/*.dll                  \
              /out/c-sharp                \
              /out/libSystem.Native.so    \
              /out/libSystem.*.OpenSsl.so \
              /out/libclrjit.so           \
              /out/libcoreclr.so          \
              /out/libhostfxr.so          \
              /out/libhostpolicy.so       /usr/bin/
COPY --from=0 /usr/lib/libcrypto.so.3     \
              /usr/lib/libgcc_s.so.1      \
              /usr/lib/libssl.so.3        \
              /usr/lib/libstdc++.so.6     /usr/lib/

ENTRYPOINT ["c-sharp"]

CMD ["--version"]
