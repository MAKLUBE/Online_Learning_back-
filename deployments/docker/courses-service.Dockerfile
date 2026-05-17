FROM golang:1.22-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o courses-service ./cmd/courses-service

FROM alpine:3.20
WORKDIR /app
COPY --from=build /app/courses-service .
EXPOSE 50053
ENTRYPOINT ["./courses-service"]
