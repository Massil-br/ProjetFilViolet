# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copie des configurations workspace et modules
COPY go.work ./
COPY backend/api/go.mod backend/api/go.sum ./backend/api/
COPY backend/game/go.mod ./backend/game/

RUN cd backend/game && go mod download

# Copie des codes source
COPY backend/api/ ./backend/api/
COPY backend/game/ ./backend/game/

# Construction du binaire
RUN cd backend/game && CGO_ENABLED=0 GOOS=linux go build -o game .

# Stage final
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/backend/game/game .
# Partager le fichier .env de l'api si besoin
COPY --from=builder /app/backend/api/.env ./ 

EXPOSE 8082

CMD ["./game"]
