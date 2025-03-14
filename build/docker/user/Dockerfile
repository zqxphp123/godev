FROM golang:1.20.4-alpine3.18 AS builder

WORKDIR /build

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOPROXY https://goproxy.cn,direct
ENV GOARCH amd64

COPY go.mod .
#COPY go.sum .
RUN go mod download
COPY . .

RUN go build -ldflags="-s -w" -o user cmd/user/user.go

FROM alpine

WORKDIR /build

ENV TZ Asia/Shanghai

COPY --from=builder /build/user .
COPY configs/ /configs
RUN mkdir logs


EXPOSE 8021
ENTRYPOINT ["./user"]
CMD ["-c", "configs/user/srv.yaml"]
