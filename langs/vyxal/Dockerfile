FROM codegolf/lang-java AS builder

ENV VER=3.10.3

ADD --chmod=755 https://github.com/Vyxal/Vyxal/releases/download/v$VER/vyxal-$VER.jar vyxal.jar

ENTRYPOINT ["/opt/java/bin/java", "-jar", "vyxal.jar"]

CMD ["--version"]
