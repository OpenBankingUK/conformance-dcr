FROM golang:1.13.6-alpine as gobuilder
RUN apk update && apk add git make bash ca-certificates

ENV TERM xterm-color
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make build

RUN /app/dcr -version

# Final image to run the binary
FROM scratch
LABEL MAINTAINER="Open Banking Ltd"

COPY --from=gobuilder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=gobuilder /app/dcr /usr/bin/dcr

EXPOSE 8080

ENTRYPOINT ["dcr"]
CMD ["--help"]
