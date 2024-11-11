FROM debian:12.7-slim

COPY clamav-scanner /usr/bin/clamav-scanner
COPY config/* /tmp

RUN export DEBIAN_FRONTEND=noninteractive && \
    apt-get update -qqy && \
    apt-get upgrade -qqy && \
    apt-get install -qqy --no-install-recommends  \
        python3-pip \
        pipx \
        apt-transport-https \
        gnupg \
        clamav-daemon \
        clamav-freshclam \
        python3-crcmod && \
        mkdir -p /clamav/cvds && \
        chown -R clamav:clamav /clamav && \
    sed -i -e '/DatabaseMirror\|DNSDatabaseInfo\|CompressLocalDatabase\|ScriptedUpdates/d' /etc/clamav/freshclam.conf && \
    cat /tmp/freshclam.tmp >> /etc/clamav/freshclam.conf && \
    sed -i -e '/LocalSocket\|StreamMaxLength\|MaxScanSize\|MaxFileSize/d' /etc/clamav/clamd.conf && \
    cat /tmp/clamd.tmp >> /etc/clamav/clamd.conf && \
    truncate -s 0 /var/log/apt/*.log /var/log/*.log && \
    rm -rf /var/lib/apt/lists/* /var/cache/apt/archives/*

USER clamav

ENV PATH "$PATH:/var/lib/clamav/.local/bin/"

RUN pipx install cvdupdate && \
    cvdupdate config set -c /clamav/config.json -d /clamav/cvds -l /clamav/log

EXPOSE 1337

CMD [ "/usr/bin/clamav-scanner" ]
