FROM alpine:latest

LABEL org.opencontainers.image.title="Toremo EventStore" \
      org.opencontainers.image.description="A simplistic EventStore designed for test purposes only."

RUN mkdir -p /app

WORKDIR /app

COPY bin/eventstore .

