FROM golang:alpine AS build-env
WORKDIR /src
ADD . /src
RUN apk add --no-cache git
RUN go get -d -v
RUN  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o senv .

FROM alpine
ENV AWS_REGION=us-east-1
RUN apk add --no-cache ca-certificates
COPY --from=build-env /src/senv /usr/bin/senv
ENTRYPOINT ["senv"] 