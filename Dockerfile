FROM golang:1.22 AS build
LABEL maintainer "Karim Boucher <me@karimboucher.com>"

WORKDIR /go/src/app
COPY . .

WORKDIR /go/src/app/cmd/carthage
RUN CGO_ENABLED=0 go build

FROM scratch
COPY --from=build /go/src/app/cmd/carthage/carthage /bin/carthage
ENTRYPOINT ["/bin/carthage"]
