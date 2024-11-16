FROM golang:1.23.3

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /ripper
ENV GIN_MODE=release
ENV TZ=UTC

CMD ["/ripper"]