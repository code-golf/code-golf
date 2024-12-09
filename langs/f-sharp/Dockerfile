FROM mcr.microsoft.com/dotnet/sdk:9.0.100-alpine3.20 AS builder

WORKDIR /work

COPY f-sharp.fs f-sharp.fsproj ./

RUN dotnet publish -c Release -o /out

# Trim additional things to improve F# compiler startup performance.
WORKDIR /trimmer

COPY Trimmer.csproj Trimmer.cs ExtraTrimmingList.txt ./

RUN dotnet run

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1    \
              /lib/libcrypto.so.3         \
              /lib/libssl.so.3            /lib/
COPY --from=0 /out/*.dll                  \
              /out/f-sharp                \
              /out/libSystem.IO.*.so      \
              /out/libSystem.Native.so    \
              /out/libSystem.*.OpenSsl.so \
              /out/libclrjit.so           \
              /out/libcoreclr.so          \
              /out/libhostfxr.so          \
              /out/libhostpolicy.so       /usr/bin/
COPY --from=0 /usr/lib/libgcc_s.so.1      \
              /usr/lib/libstdc++.so.6     /usr/lib/

ENV DOTNET_SYSTEM_GLOBALIZATION_INVARIANT=1

ENTRYPOINT ["f-sharp"]

CMD ["--version"]
