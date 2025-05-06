FROM golang:1.24.2 AS builder

ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG REVISION=unknown

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd/ cmd/
COPY internal/ internal/
COPY *.go .

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -ldflags="-X github.com/FeLvi-zzz/tentez.Revision=${REVISION}" -o tentez ./cmd/tentez/main.go

FROM gcr.io/distroless/static:nonroot

LABEL org.opencontainers.image.title="FeLvi-zzz/tentez"

WORKDIR /
COPY --from=builder /app/tentez /tentez
USER 65532:65532

ENTRYPOINT ["/tentez"]
