FROM node:20-alpine AS frontend
WORKDIR /app/web
COPY web/package*.json ./
RUN npm ci
COPY web/ .
RUN npm run build

FROM golang:1.23-alpine AS backend
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /app/web/dist ./web/dist
RUN go build -o server ./cmd/server

FROM alpine:3.19
WORKDIR /app
COPY --from=backend /app/server .
COPY --from=backend /app/web/dist ./web/dist
EXPOSE 8080
CMD ["./server"]
