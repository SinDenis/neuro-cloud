FROM golang:1.18-alpine

RUN apk add --update --no-cache alpine-sdk bash ca-certificates \
      libressl \
      tar \
      git openssh openssl yajl-dev zlib-dev cyrus-sasl-dev openssl-dev build-base coreutils
WORKDIR /root
RUN git clone https://github.com/edenhill/librdkafka.git
WORKDIR /root/librdkafka
RUN /root/librdkafka/configure

WORKDIR /app

ADD . /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download
RUN go build -tags musl --ldflags "-extldflags -static" cmd/api/main.go

EXPOSE 8080

ENTRYPOINT ["./main"]