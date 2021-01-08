ARG GOPKG=github.com/UiP9AV6Y/basic-oauth2

FROM golang:1.15-buster AS builder

ARG GOPKG
WORKDIR ${GOPATH}/src/${GOPKG}

COPY go.mod go.sum ./
RUN GO111MODULE=on go mod download

COPY . .
RUN set -xe; \
  make build BUILD_DIR=/build

FROM scratch

ARG BUILD_DATE="1970-01-01T00:00:00Z"
ARG VERSION="0.0.0"
ARG VCS_URL="http://localhost/"
ARG VCS_REF="main"
LABEL org.opencontainers.image.title="basic-oauth2" \
  org.opencontainers.image.description="OIDC compatible webserver utilizing Basic Authentication" \
  org.opencontainers.image.created="$BUILD_DATE" \
  org.opencontainers.image.authors="Gordon Bleux" \
  org.opencontainers.image.source="$VCS_URL" \
  org.opencontainers.image.version="$VERSION" \
  org.opencontainers.image.revision="$VCS_REF" \
  org.opencontainers.image.vendor="Gordon Bleux" \
  org.opencontainers.image.licenses="Apache-2.0"

ENV XDG_CONFIG_HOME=/etc/basic-oauth2
COPY config/ ${XDG_CONFIG_HOME}/

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/ /bin/

HEALTHCHECK --timeout=5s --start-period=5s \
  CMD ["/bin/basic-oauth2", "server", "health"]

ENTRYPOINT ["/bin/basic-oauth2"]
CMD ["server", "run"]
