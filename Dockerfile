# Etapa 1: compila el binario Go dentro de una imagen con toolchain de Go.
FROM golang:1.25-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# CGO_ENABLED=0 genera un binario estático, útil para una imagen final pequeña.
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/piscina-api ./cmd/piscina-api

# Etapa 2: imagen final mínima, sin compilador Go.
FROM alpine:3.22

WORKDIR /app

COPY --from=build /out/piscina-api /app/piscina-api
COPY web /app/web

# Usuario no-root por seguridad.
RUN adduser -D -u 10001 appuser && chown -R appuser:appuser /app
USER appuser

EXPOSE 8080

CMD ["/app/piscina-api"]
