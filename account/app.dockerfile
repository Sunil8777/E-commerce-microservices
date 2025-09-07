# Build stage
FROM golang:1.13-alpine3.11 AS build
RUN apk add --no-cache gcc g++ make ca-certificates

WORKDIR /go/src/github.com/sunil8777/E-commerce-microservices

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /go/bin/app ./account/cmd/account

# Final stage
FROM alpine:3.11
WORKDIR /usr/bin

COPY --from=build /go/bin/app .

EXPOSE 8080
CMD ["/usr/bin/app"]
