FROM golang:1.25 AS build

RUN apt-get update && apt-get install -y \
    gcc \
    libc6-dev \
    librdkafka-dev \
    pkg-config \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o /herald

FROM gcr.io/distroless/base-debian10

ENV GIN_MODE=release
ENV TZ=UTC

COPY --from=build /herald /

CMD ["/herald"]