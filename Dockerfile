# Stage 1: Build frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app

COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci

COPY frontend/ .
RUN npm run build

# Stage 2: Build backend
FROM golang:1.24-alpine AS backend-builder

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage 3: Production
FROM alpine:latest

WORKDIR /app

COPY --from=backend-builder /app/main .
COPY --from=frontend-builder /app/dist ./static

EXPOSE 8080

CMD ["./main"]
