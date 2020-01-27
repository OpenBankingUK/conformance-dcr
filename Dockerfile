FROM golang:1.13.6-alpine as gobuilder
RUN apk update && apk add git make bash ca-certificates

ENV TERM xterm-color
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app

RUN printf "%b" "\033[92m" "gobuilder: copying only go.mod and go.sum for pre-cache step..." "\033[0m" "\n"
COPY go.mod go.sum ./
RUN go mod download
RUN printf "%b" "\033[92m" "gobuilder: copying all files ..." "\033[0m" "\n"
COPY . .
RUN make build

RUN printf "%b" "\033[92m" "gobuilder: printing DCR version information ..." "\033[0m" "\n" && \
    /app/dcr -version

# Final image to run the binary
FROM scratch
LABEL MAINTAINER="Open Banking Ltd"

COPY --from=gobuilder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=gobuilder /app/dcr /usr/bin/dcr

EXPOSE 8080

ENTRYPOINT ["dcr"]
CMD ["--help"]
