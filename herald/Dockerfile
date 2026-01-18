FROM golang:1.25 AS build

# Install librdkafka and build tools for confluent-kafka-go
RUN apt-get update && apt-get install -y \
    gcc \
    libc6-dev \
    librdkafka-dev \
    pkg-config \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /usr/src/app

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o /herald ./herald

FROM gcr.io/distroless/base-debian10

ENV GIN_MODE=release
ENV TZ=UTC

COPY --from=build /herald /

CMD ["/herald"]