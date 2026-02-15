FROM golang:1.25.3 AS builder

WORKDIR /app

COPY . .

RUN cd cmd && CGO_ENABLED=0 GOOS=linux go build -o app main.go


FROM alpine AS app

WORKDIR /

COPY --from=builder /app/cmd/app ./

COPY --from=builder /app/db/migrations ./migrations

ENTRYPOINT ["./app"]
