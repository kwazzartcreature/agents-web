FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git make npm nodejs

RUN go install github.com/a-h/templ/cmd/templ@latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make build


FROM alpine

WORKDIR /root/

COPY --from=builder /app/build ./build

RUN chmod +x ./build/main

ENV GIN_MODE=release

EXPOSE 8000

CMD ["./build/main"]
