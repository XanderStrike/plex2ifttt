FROM golang:1.12.7-alpine as builder
WORKDIR $GOPATH/src/github.com/xanderstrike/plexlights/
RUN apk add --no-cache git
COPY . .
RUN mkdir /out
RUN mkdir /out/keystore
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/plexlights-docker

FROM scratch
LABEL maintainer="xanderstrike@gmail.com"
ADD https://github.com/golang/go/raw/master/lib/time/zoneinfo.zip /zoneinfo.zip
ENV ZONEINFO /zoneinfo.zip
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /out .
EXPOSE 8080
ENTRYPOINT ["/app/plexlights-docker"]
