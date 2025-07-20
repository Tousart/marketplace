FROM golang:1.23.5-alpine AS build

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

FROM alpine:latest AS runner

WORKDIR /marketplace

COPY --from=build /build/main ./main

COPY --from=build /build/config/config.yaml ./config.yaml

COPY --from=build /build/resources/id.txt ./resources/id.txt

CMD ["./main", "--config=./config.yaml"]