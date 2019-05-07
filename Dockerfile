FROM golang:1.12.2-alpine@sha256:d481168873b7516b9f34d322615d589fafb166ff5fd57d93e96f64787a58887c AS builder

WORKDIR /app
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -a -o kramp .


FROM scratch
LABEL maintainer="meneer@koenbollen.nl"

ENV LIMIT_ALBUMS=5
ENV LIMIT_BOOKS=5
ENV ENV=development
EXPOSE 8080
ENTRYPOINT ["/kramp"]

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/kramp /kramp
