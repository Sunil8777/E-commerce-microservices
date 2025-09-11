# Build stage
FROM golang:1.24-alpine AS build
RUN apk add --no-cache gcc g++ make ca-certificates

WORKDIR /go/src/github.com/sunil8777/E-commerce-microservices

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0
RUN go build -trimpath -o /go/bin/app ./account/cmd/account

# Final stage
FROM alpine:3.19
WORKDIR /usr/bin

COPY --from=build /go/bin/app .

EXPOSE 8080
CMD ["app"]
