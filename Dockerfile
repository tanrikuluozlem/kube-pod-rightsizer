FROM golang:1.24-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /bin/rightsizer ./cmd/main.go

FROM gcr.io/distroless/static:nonroot

COPY --from=builder /bin/rightsizer /rightsizer

USER nonroot:nonroot

ENTRYPOINT ["/rightsizer"]
