# Base
FROM golang:1-bullseye

WORKDIR /cultivator

# Copy
COPY backend /cultivator/backend
COPY backend/config /cultivator/backend/config
COPY backend/model /cultivator/backend/model
COPY backend/pull /cultivator/backend/pull
COPY backend/routes /cultivator/backend/routes
COPY backend/storage /cultivator/backend/storage
COPY harvester /cultivator/harvester
COPY node /cultivator/node
COPY node/command /cultivator/node/command
COPY go.mod /cultivator/go.mod
COPY go.sum /cultivator/go.sum
COPY Bootstrap.go /cultivator/Bootstrap.go

# Build / Setup
RUN apt-get -y update
RUN apt-get -y install unzip
RUN apt-get -y install sudo
RUN go build -o /cultivator/Cultivator

CMD ["./Cultivator"]