FROM golang:1.27.1-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY main.go http.go ./
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X main.version=${VERSION}" -o /portlight .

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /portlight /portlight
USER 65532:65532
ENV PORT=8991
EXPOSE 8991
STOPSIGNAL SIGTERM
ENTRYPOINT ["/portlight"]
