FROM ghcr.io/trevatk/cfgo:1.22.5-dev-cf AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod tidy && go mod verify

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /usr/src/bin/server ./cmd/server

FROM scratch

WORKDIR /usr/local/bin

USER httpbin

COPY --from=builder /usr/src/bin/server .

ENTRYPOINT ["server"]