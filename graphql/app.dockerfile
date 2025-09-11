FROM golang:1.24-alpine AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0
RUN go build -trimpath -o app ./graphql

FROM alpine:3.19
WORKDIR /usr/bin
COPY --from=build /app/app .
EXPOSE 8080
CMD ["./app"]
