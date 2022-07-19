FROM golang as build
WORKDIR /battery-email-api
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build

FROM debian:10-slim
RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates
COPY --from=build /battery-email-api/battery-email-api /usr/local/bin/battery-email-api
EXPOSE 80
