FROM codegolf/lang-java AS builder

ENV VER=3.7.0

ADD https://github.com/Vyxal/Vyxal/releases/download/v$VER/vyxal-$VER.jar /vyxal.jar

CMD ["/opt/java/bin/java", "-jar", "vyxal.jar", "-h"]
