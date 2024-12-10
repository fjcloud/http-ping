FROM registry.redhat.io/ubi9/go-toolset:latest as builder

WORKDIR /app
COPY . .

USER root
RUN chown -R 1001:1001 /app
USER 1001

RUN go mod init websocket-tester && \
    go mod tidy && \
    CGO_ENABLED=0 go build -buildvcs=false -ldflags="-s -w" -o main .

FROM registry.redhat.io/ubi9/ubi-minimal:latest

WORKDIR /app
COPY --from=builder /app/main .
COPY static/ static/

EXPOSE 8080

CMD ["./main"]
