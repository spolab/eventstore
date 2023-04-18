FROM alpine:latest

ARG version

LABEL org.opencontainers.image.source="https://github.com/spolab/eventstore"
LABEL org.opencontainers.image.description="The super-simplistic event store"
LABEL org.opencontainers.image.licenses="MIT"

RUN apk add --no-cache gcompat && \
    mkdir -p /app

WORKDIR /app

COPY bin/eventstore .

ENTRYPOINT [ "bin/eventstore" ]