FROM mcr.microsoft.com/dotnet/sdk:9.0.100-alpine3.20 AS builder

RUN curl -o framework.exe https://go.microsoft.com/fwlink/?linkid=2099466

RUN chmod +x framework.exe \
 && ./framework.exe --install

ENV VER=1.2.1

RUN curl -#L https://github.com/tomtheisen/stax/archive/refs/tags/v$VER.tar.gz \
  | tar xz --strip-components 1

RUN dotnet publish StaxLang.CLI -c Release -o /usr/bin

FROM codegolf/lang-base

COPY --from=0 / /

ENTRYPOINT ["stax"]

CMD ["--version"]
