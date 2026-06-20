# --- Etapa de compilación ---
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Instalar git (necesario para algunas dependencias)
RUN apk add --no-cache git

# Descargar dependencias primero (mejor caché de capas)
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código y compilar un binario estático
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./main.go

# --- Etapa final (imagen mínima) ---
FROM alpine:latest

# Certificados raíz para conexiones HTTPS (Supabase, Google Drive)
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/server .

# Render inyecta el puerto vía variable PORT
EXPOSE 8080

CMD ["./server"]
