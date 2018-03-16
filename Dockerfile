FROM            debian:9
MAINTAINER      Luka Peschke <mail@lukapeschke.com>

RUN apt update \
    && apt install -y \
        git \
        golang \
        libudev-dev \
    && rm -rf /var/cache/apt

ADD entrypoint.sh /root/entrypoint.sh
WORKDIR /root/monitor-status
ENV GOPATH /root/go

ENTRYPOINT ["/root/entrypoint.sh"]
