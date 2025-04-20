FROM mcr.microsoft.com/dotnet/sdk:9.0.203-alpine3.21 AS builder

WORKDIR /powershell

COPY powershell.cs powershell.csproj ./

RUN dotnet publish -c Release -o /out

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /out/*.dll               \
              /out/libSystem.Native.so \
              /out/libclrjit.so        \
              /out/libcoreclr.so       \
              /out/libhostfxr.so       \
              /out/libhostpolicy.so    \
              /out/libpsl-native.so    \
              /out/powershell          /usr/bin/
COPY --from=0 /usr/lib/libgcc_s.so.1   \
              /usr/lib/libstdc++.so.6  /usr/lib/

ENV DOTNET_SYSTEM_GLOBALIZATION_INVARIANT=1

ENTRYPOINT ["powershell"]

CMD ["--version"]
