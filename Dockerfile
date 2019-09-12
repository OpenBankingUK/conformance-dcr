FROM golang:1.12-alpine as gobuilder
RUN apk update && apk add git make bash

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app

COPY . .
RUN make build

# Print version information
RUN /app/dcr -version

# Final image to run the binary
FROM scratch
LABEL MAINTAINER Open Banking Ltd

WORKDIR /app

COPY --from=gobuilder /app/dcr /app/

CMD ["/app/dcr"]
