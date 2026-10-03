FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /janus ./cmd/janus

FROM gcr.io/distroless/static-debian12
COPY --from=builder /janus /janus
COPY configs/janus.example.yaml /etc/janus/janus.yaml
USER nonroot:nonroot
ENTRYPOINT ["/janus", "serve", "--config", "/etc/janus/janus.yaml"]
