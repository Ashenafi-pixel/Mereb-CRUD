FROM golang:1.23.0 AS builder

WORKDIR /build
COPY go.* .
RUN go mod download
COPY . .
RUN go build -o ./mereb-crud-server cmd/merebapi/main.go

FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /build/mereb-crud-server ./mereb-crud-server

ENV DB_USER=postgres
ENV DB_PORT=5434
ENV DB_PASSWORD=password
ENV DB_HOST=host.docker.internal
ENV DB_NAME=person
ENV JWT_KEY=ultra-super-secret

CMD ["/app/mereb-crud-server"]
