FROM alpine:latest

ARG version

LABEL org.opencontainers.image.title="Toremo EventStore" \
      org.opencontainers.image.description="A simplistic EventStore designed for test purposes only." \
      org.opencontainers.image.version=${version}

RUN mkdir -p /app

WORKDIR /app

COPY bin/eventstore .

