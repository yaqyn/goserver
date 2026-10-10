FROM golang:1.27.1-alpine@sha256:8a5910f31396cd4d89662f56c68b3ae31d374308270a1c3bd96672ee5ed43414 AS build
WORKDIR /src
COPY go.mod ./
COPY main.go http.go ./
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X main.version=${VERSION}" -o /portlight .

FROM scratch
ARG VERSION=dev
LABEL org.opencontainers.image.title="Portlight — HTTP Server" \
      org.opencontainers.image.description="A focused Go HTTP service with clear operational signals" \
      org.opencontainers.image.source="https://github.com/yaqyn/portlight-http-server" \
      org.opencontainers.image.version="${VERSION}"
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /portlight /portlight
USER 65532:65532
ENV PORT=8991
EXPOSE 8991
STOPSIGNAL SIGTERM
ENTRYPOINT ["/portlight"]
