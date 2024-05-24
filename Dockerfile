FROM debian:12.5-slim

ENV PATH "$PATH:/opt/google-cloud-sdk/bin/"

COPY clamav-scanner /usr/bin/clamav-scanner
COPY config/* /tmp

RUN export DEBIAN_FRONTEND=noninteractive && \
    apt-get update -qqy && \
    apt-get upgrade -qqy && \
    apt-get install -qqy --no-install-recommends  \
        curl \
        python3-pip \
        pipx \
        apt-transport-https \
        lsb-release \
        openssh-client \
        gnupg \
        clamav-daemon \
        clamav-freshclam \
        python3-crcmod && \
        mkdir -p /clamav/cvds && \
        chown -R clamav:clamav /clamav && \
    CLOUD_SDK_REPO="cloud-sdk-$(lsb_release -c -s)" && \
    echo -n "Adding Cloud SDK apt repository: " && \
    echo "deb https://packages.cloud.google.com/apt $CLOUD_SDK_REPO main" \
        | tee /etc/apt/sources.list.d/google-cloud-sdk.list && \
    curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg \
        | gpg --dearmor -o /etc/apt/trusted.gpg.d/packages-cloud-google-apt.gpg && \
    apt-get update -qqy && \
    apt-get install -qqy --no-install-recommends google-cloud-sdk && \
    gcloud config set core/disable_usage_reporting true && \
    gcloud config set component_manager/disable_update_check true && \
    gcloud config set metrics/environment github_docker_image && \
    gcloud --version && \
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
