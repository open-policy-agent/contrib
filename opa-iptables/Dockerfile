FROM alpine

RUN apk --no-cache add iptables ca-certificates && \
    update-ca-certificates

ADD opa-iptables /

ENTRYPOINT ["/opa-iptables"]