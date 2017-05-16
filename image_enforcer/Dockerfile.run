FROM alpine

MAINTAINER Torin Sandall <torinsandall@gmail.com>

RUN apk --no-cache add ca-certificates && \
    update-ca-certificates

ADD clair-layer-sync_linux_amd64 /

ENTRYPOINT ["/clair-layer-sync_linux_amd64"]
