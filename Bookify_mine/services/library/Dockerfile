# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS builder
WORKDIR /src

# Copiamos todo el repositorio (incluye go.mod, go.sum, proto, pkg, services, etc.)
COPY . .

# Descargar dependencias
RUN go mod download

# Compilar solo el microservicio Library
RUN CGO_ENABLED=0 go build -o /opt/service ./services/library/cmd/library

# Etapa final
FROM alpine:3.19
WORKDIR /app

# Certificados para conexiones HTTPS/TLS
RUN apk add --no-cache ca-certificates

# Copiar el binario compilado
COPY --from=builder /opt/service /usr/local/bin/service

# Exponer puertos
EXPOSE 8081
EXPOSE 50051

ENTRYPOINT ["/usr/local/bin/service"]
