#Build stage
FROM golang:alpine3.16 AS builder
WORKDIR /app
COPY . .
RUN apk add curl
RUN go build -o main main.go

#Run stage
FROM alpine:3.16 
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./db/migration
EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]