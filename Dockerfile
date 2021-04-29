FROM golang:1.14.0-alpine as build

WORKDIR /go/src/app
COPY go.mod go.mod
COPY go.sum go.sum
COPY sample_exporter.go sample_exporter.go
#COPY loadbalancer loadbalancer
RUN export CGO_ENABLED=0 && export GOOS=linux && go get . && go build

FROM alpine
RUN apk --no-cache add ca-certificates
COPY --from=build /go/src/app/sample_exporter /go/src/app/sample_exporter

CMD ["/go/src/app/sample_exporter"]
