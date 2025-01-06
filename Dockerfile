FROM golang:1.23 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /config-server .

#Second stage of build
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app/

COPY --from=build /config-server .

EXPOSE 8080

CMD ["./config-server", "api"]
