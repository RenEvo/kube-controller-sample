FROM golang:1.11.2-alpine as builder

WORKDIR /go-modules
COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o kube-controller main.go

FROM alpine

RUN adduser -S -D -H -h /home/controller controller
COPY --from=builder /go-modules/kube-controller /usr/local/bin/
USER controller

WORKDIR /home/controller

ENTRYPOINT ["/usr/local/bin/kube-controller"]
CMD []