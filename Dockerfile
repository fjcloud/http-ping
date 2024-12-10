FROM registry.redhat.io/ubi9/go-toolset:latest as builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 go build -buildvcs=false -ldflags="-s -w" -o main .

FROM registry.redhat.io/ubi9/ubi-minimal:latest

WORKDIR /app
COPY --from=builder /app/main .
COPY static/ static/

EXPOSE 8080

CMD ["./main"]
