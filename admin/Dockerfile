FROM golang:1.22 AS build

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /api

FROM gcr.io/distroless/base-debian10

ENV GIN_MODE=release
ENV TZ=UTC

COPY --from=build /api /

CMD ["/api"]