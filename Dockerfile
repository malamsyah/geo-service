FROM golang:alpine AS base

RUN apk --update add --no-cache ca-certificates openssl git tzdata && \
update-ca-certificates

WORKDIR $GOPATH/src/github.com/malamsyah/geo-service/

COPY . .

RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /server ./cmd/server

FROM scratch

COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=base /server .

CMD ["./server"]
EXPOSE 8080