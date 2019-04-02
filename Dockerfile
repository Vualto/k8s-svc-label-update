FROM golang:1.11.2 AS builder

COPY glide.yaml /go/src/app/glide.yaml
WORKDIR /go/src/app

RUN apt-get update && \
    apt-get install -y golang-glide git && \
    glide install

COPY . /go/src/app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o app .

FROM scratch
LABEL maintainer="roger.pales@vualto.com"

ADD build/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/app/app .

CMD ["./app"]
