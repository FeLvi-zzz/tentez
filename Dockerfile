FROM gcr.io/distroless/static:nonroot

# TARGETOS and TARGETARCH are automatically set by buildkit if `--platform` is passed
# https://github.com/moby/buildkit/blob/master/frontend/dockerfile/docs/reference.md#automatic-platform-args-in-the-global-scope
ARG TARGETOS
ARG TARGETARCH

ARG REVISION=unknown

WORKDIR /

# NOTE: require the pre-built binaries in dist/
COPY dist/tentez-${TARGETOS}-${TARGETARCH} /tentez

LABEL org.opencontainers.image.title="FeLvi-zzz/tentez"

USER 65532:65532

ENTRYPOINT ["/tentez"]
