FROM golang:1.25-alpine AS builder
WORKDIR /plugin

COPY go.mod go.sum ./
RUN go mod download

COPY plugin.go .
COPY examples/ examples/
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o opa ./examples/

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /plugin/opa /opa
USER nonroot:nonroot
ENTRYPOINT ["/opa"]
CMD ["run"]
