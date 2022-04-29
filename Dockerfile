# Build stage
FROM golang:1.18.1-alpine3.15  as builder

WORKDIR /app
COPY . .
RUN go build -o simplebank main.go

# Run stage
FROM alpine:3.15

WORKDIR /app

COPY --from=builder /app/simplebank .
COPY app.env .

EXPOSE 8000
CMD [ "/app/simplebank" ]