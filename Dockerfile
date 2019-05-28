FROM golang:1.12-alpine as gobuilder
RUN apk update && apk add git make bash

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN mkdir /app
WORKDIR /app

COPY . .
RUN make build

# Final image to run the binary
FROM scratch
LABEL MAINTAINER Open Banking Ltd

WORKDIR /app

COPY manifests/* /manifests
COPY --from=gobuilder /app/dcr /app/

CMD ["/app/dcr"]
