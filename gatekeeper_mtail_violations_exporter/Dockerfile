# build mtail binary
FROM golang:alpine AS builder

WORKDIR /go/src/github.com/google/mtail

# mtail repo includes a tag for each "release"
# build the mtail binary using the rc36 release
RUN apk add --update --no-cache --virtual build-dependencies git make \
 && git clone https://github.com/google/mtail /go/src/github.com/google/mtail \
 && git checkout v3.0.0-rc36 \
 && make depclean && make install_deps \
 && PREFIX=/go make STATIC=y -B install

# package mtail binary in a small scratch image (total size ~12MB)
FROM scratch
COPY --from=builder /go/bin/mtail /usr/bin/mtail
ENTRYPOINT ["/usr/bin/mtail"]
EXPOSE 3903
WORKDIR /tmp