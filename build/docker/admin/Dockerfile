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

RUN go build -ldflags="-s -w" -o admin cmd/admin/admin.go

FROM alpine

WORKDIR /build

ENV TZ Asia/Shanghai

COPY --from=builder /build/admin .
COPY configs/ /configs
RUN mkdir logs


EXPOSE 8050
ENTRYPOINT ["./admin"]
CMD ["-c", "configs/admin/admin.yaml"]
