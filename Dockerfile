FROM ghcr.io/trevatk/cfgo:1.22.5-dev-cf AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod tidy && go mod verify

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /usr/src/bin/httpbin ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /usr/local/bin

COPY --from=builder /usr/src/bin/httpbin ./

EXPOSE 8080

ENTRYPOINT [ "httpbin" ]