FROM mcr.microsoft.com/dotnet/sdk:9.0.306-alpine3.22-amd64 AS builder

WORKDIR /f-sharp

COPY f-sharp.fs f-sharp.fsproj ./

RUN dotnet publish -c Release -o /out

# Trim additional things to improve F# compiler startup performance.
WORKDIR /trimmer

COPY trimmer.cs trimmer.csproj trimmer.txt ./

RUN dotnet run

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1         /lib/
COPY --from=0 /out/*.dll                  \
              /out/f-sharp                \
              /out/libSystem.IO.*.so      \
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

ENV DOTNET_SYSTEM_GLOBALIZATION_INVARIANT=1

ENTRYPOINT ["f-sharp"]

CMD ["--version"]
