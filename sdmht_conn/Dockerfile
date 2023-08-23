# build
FROM golang:1.20 as builder
ENV GO111MODULE='on' GOPROXY='https://goproxy.cn|https://goproxy.io|direct'
WORKDIR /go/src/pantheon/

# To use cached layer by docker when dependencies not change, copy go.mod and go.sum first.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# note default value of CGO_ENABLED is `1` in image `golang`, we should change it to `0`
RUN CGO_ENABLED=0 go build -installsuffix cgo -o .build/webinar_conn ./webinar/webinar_conn/app/webinar_conn

FROM alpine:latest
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && apk --no-cache add ca-certificates

WORKDIR /opt/pantheon/webinar/webinar_conn
COPY --from=builder /go/src/pantheon/.build .
EXPOSE 7032 4647
ENTRYPOINT ["./webinar_conn"]