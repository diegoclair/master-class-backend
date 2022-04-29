# Build stage
FROM golang:1.18.1-alpine3.15  as builder

WORKDIR /app
COPY . .
RUN go build -o simplebank main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.15

WORKDIR /app

COPY --from=builder /app/simplebank .
COPY --from=builder /app/migrate ./migrate
COPY db/migrations ./migrations
COPY app.env .
COPY start.sh .

## Add the wait script to the image
ENV WAIT_VERSION 2.9.0
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/$WAIT_VERSION/wait /wait
RUN chmod +x /wait

EXPOSE 8000
CMD [ /wait && /app/simplebank ]
#ENTRYPOINT [ /app/start.sh ]