FROM golang:1.20-alpine as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY cmd /build/cmd
COPY internal /build/internal

RUN CGO_ENABLED=0 GOOS=linux go build -o /server cmd/server/main.go

FROM alpine

WORKDIR /

COPY --from=builder /server /server

ENTRYPOINT ["/server"]
