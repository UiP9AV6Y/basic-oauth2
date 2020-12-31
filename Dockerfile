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

ENV XDG_CONFIG_HOME=/etc/basic-oauth2
COPY config/ ${XDG_CONFIG_HOME}/

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/ /bin/

HEALTHCHECK --timeout=5s --start-period=5s \
  CMD ["/bin/basic-oauth", "server", "health"]

ENTRYPOINT ["/bin/basic-oauth2"]
CMD ["server", "run"]
