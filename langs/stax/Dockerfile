FROM mcr.microsoft.com/dotnet/sdk:9.0.306-alpine3.22-amd64 AS builder

RUN apk add --no-cache gcc musl-dev

ENV VER=1.2.1

RUN curl -#L https://github.com/tomtheisen/stax/archive/refs/tags/v$VER.tar.gz \
  | tar xz --strip-components 1

WORKDIR /StaxLang.CLI

COPY stax.csproj StaxLang.CLI.csproj

RUN dotnet publish -c Release -o /out

COPY stax.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/staxwrapper -s -static /stax.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1      /lib/
COPY --from=0 /out/*.dll               \
              /out/libSystem.Native.so \
              /out/libclrjit.so        \
              /out/libcoreclr.so       \
              /out/libhostfxr.so       \
              /out/libhostpolicy.so    \
              /out/stax                \
              /usr/bin/staxwrapper     /usr/bin/
COPY --from=0 /usr/lib/libgcc_s.so.1   \
              /usr/lib/libstdc++.so.6  /usr/lib/

ENTRYPOINT ["staxwrapper"]

CMD ["--version"]
