FROM golang:1.26.3-alpine3.23 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o requests_api .

FROM alpine:3.23 as production
WORKDIR /app
COPY --from=builder /app/requests_api .
RUN adduser -D -u 1000 helpdesk
USER helpdesk
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget -qO- http://localhost:8080/health || exit 1
CMD ["./requests_api"]
