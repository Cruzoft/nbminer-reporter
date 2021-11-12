FROM alpine:3.14 as alpine

RUN apk --no-cache add tzdata zip ca-certificates

WORKDIR /usr/share/zoneinfo

# -0 means no compression.  Needed because go's
# tz loader doesn't handle compressed data.
RUN zip -q -r -0 /zoneinfo.zip .

FROM golang:1.17.3-alpine3.14 as builder

WORKDIR /go/src

COPY ./go.mod ./

RUN apk add build-base 
RUN go mod download

FROM builder as binary

COPY ./ /go/src

RUN CGO_ENABLED=0 go build -ldflags="-w -s -X main.buildTime=`date +'%Y-%m-%d_%T'`" -o /go/bin/nbsim ./cmd/simulator/...

FROM binary

RUN go test ./cmd/simulator/...

FROM scratch

ENV TWITTER_LANGUAGE=es 

LABEL maintainer="ingcsmoreno@gmail.com"

# the timezone data:
ENV ZONEINFO /zoneinfo.zip
COPY --from=alpine /zoneinfo.zip /

# the tls certificates:
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=binary /go/bin/nbsim /go/bin/nbsim

ENTRYPOINT ["/go/bin/nbsim"]