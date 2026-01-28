FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git make curl

WORKDIR /build
RUN git clone https://github.com/argoeu/argo-web-api.git .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o argo-web-api .

FROM alpine:3.18

RUN apk add --no-cache ca-certificates openssl

RUN mkdir -p /etc/pki/tls/certs /etc/pki/tls/private /app

COPY --from=builder /build/argo-web-api /app/argo-web-api

RUN chmod +x /app/argo-web-api

RUN openssl req -x509 -newkey rsa:4096 -keyout /etc/pki/tls/private/localhost.key -out /etc/pki/tls/certs/localhost.crt -days 365 -nodes -subj "/CN=localhost"

WORKDIR /app

EXPOSE 443

CMD ["/app/argo-web-api"]
