FROM golang:1.26-alpine AS builder

WORKDIR /src
RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download
RUN go install -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/app ./cmd/app

FROM alpine:3.22

RUN adduser -D -u 1000 appuser && apk add --no-cache ca-certificates wget
WORKDIR /app

COPY --from=builder /out/app /app/app
COPY --from=builder /go/bin/migrate /app/migrate
COPY internal/migrations /app/internal/migrations

USER appuser
EXPOSE 8080
HEALTHCHECK CMD wget --spider -q http://localhost:8080/health || exit 1
CMD ["/app/app"]
