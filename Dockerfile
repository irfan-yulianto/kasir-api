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
FROM alpine:3.21

RUN addgroup -g 1000 app && adduser -D -u 1000 -G app app
RUN apk add --no-cache wget

WORKDIR /app

COPY --from=backend-builder /app/main .
COPY --from=frontend-builder /app/dist ./static

RUN chown -R app:app /app

USER app

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./main"]
