# Base
FROM golang:1.17.0-alpine3.14

# Setup
RUN mkdir /cultivator
ADD . /cultivator
WORKDIR /cultivator
RUN go mod download
RUN go build -o cultivator .

# Run
CMD ["/cultivator/cultivator"]
