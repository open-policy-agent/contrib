FROM golang:alpine
WORKDIR /go/src/app
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o opa-iptables-example .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/app/opa-iptables-example .
CMD ["./opa-iptables-example"]