FROM golang:1.22 AS builder

WORKDIR /src

ADD . .
ARG VERSION
RUN \
  --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  CGO_ENABLED=0 \
  GOOS=linux \
  go build \
  -ldflags="-X main.Version=$VERSION" \
  -o /protoc-gen-go-o5-messaging ./cmd/protoc-gen-go-o5-messaging

FROM scratch
LABEL org.opencontainers.image.source = "https://github.com/pentops/o5-messaging"
COPY --from=builder /protoc-gen-go-o5-messaging /
ENTRYPOINT ["/protoc-gen-go-o5-messaging"]

