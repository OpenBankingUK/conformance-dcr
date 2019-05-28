FROM golang:1.12-alpine as gobuilder
RUN apk update && apk add git make bash

# disable crosscompiling
#
# A normal compiled app is dynamically linked to the libraries it needs to run (i.e., all the C libraries it binds to).
# Unfortunately, scratch is empty, so there are no libraries and no loadpath for it to look in. What we have to do is modify our build script to statically compile our app with all libraries built in.
#
# https://github.com/AlessioCoser/minimal-docker-container-for-golang
ENV CGO_ENABLED=0
# compile linux only
ENV GOOS=linux
ENV GOARCH=amd64

# For caching technique, see: https://medium.com/@petomalina/using-go-mod-download-to-speed-up-golang-docker-builds-707591336888
# All these steps will be cached
WORKDIR /app

# COPY the source code.
COPY . .

# Final image to run the binary
FROM scratch
LABEL MAINTAINER Open Banking Ltd

WORKDIR /app

EXPOSE 8443

ENTRYPOINT ["/app/"]
