FROM mcr.microsoft.com/dotnet/sdk:10.0.203-alpine3.23-amd64 AS builder

RUN apk add --no-cache gcc musl-dev

ENV VER=1.2.1

WORKDIR /stax

ADD https://github.com/tomtheisen/stax.git#v$VER .

COPY stax.csproj StaxLang.CLI/StaxLang.CLI.csproj

RUN dotnet publish StaxLang.CLI -c Release -o /out

COPY stax.c .

RUN gcc -Wall -Werror -Wextra -o /usr/bin/stax -s -static stax.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1      /lib/
COPY --from=0 /usr/bin/stax            /usr/bin/
COPY --from=0 /usr/lib/libgcc_s.so.1   \
              /usr/lib/libstdc++.so.6  /usr/lib/
COPY --from=0 /out/*.dll               \
              /out/libSystem.Native.so \
              /out/libclrjit.so        \
              /out/libcoreclr.so       \
              /out/libhostfxr.so       \
              /out/libhostpolicy.so    \
              /out/stax                /usr/local/bin/

ENV DOTNET_SYSTEM_GLOBALIZATION_INVARIANT=1 PATH=/usr/bin

ENTRYPOINT ["stax"]

CMD ["--version"]
