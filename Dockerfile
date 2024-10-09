# ARG DEBIAN_IMAGE=debian:stable-slim
# ARG BASE=gcr.io/distroless/static-debian11:nonroot
# FROM --platform=$BUILDPLATFORM ${DEBIAN_IMAGE} AS build
# SHELL [ "/bin/sh", "-ec" ]

# RUN export DEBCONF_NONINTERACTIVE_SEEN=true \
#            DEBIAN_FRONTEND=noninteractive \
#            DEBIAN_PRIORITY=critical \
#            TERM=linux ; \
#     apt-get -qq update ; \
#     apt-get -yyqq upgrade ; \
#     apt-get -yyqq install ca-certificates libcap2-bin; \
#     apt-get clean


# FROM golang:1.19.0-alpine
FROM golang:1.21.3-alpine

# get dig
RUN apk add bind-tools

WORKDIR /coredns/

COPY go.mod .
COPY go.sum .

RUN go mod download && go mod verify

# and then add the code
ADD . /coredns

RUN go generate

# RUN CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -a -o  
RUN CGO_ENABLED=0 GOOS=linux go build -a -o coredns coredns.go

# COPY coredns /coredns
# RUN setcap cap_net_bind_service=+ep /coredns

# FROM --platform=$TARGETPLATFORM ${BASE}
# COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# COPY --from=build /coredns /coredns
# USER nonroot:nonroot

EXPOSE 53 53/udp
# ENTRYPOINT ["/coredns"]
