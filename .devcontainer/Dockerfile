FROM mcr.microsoft.com/devcontainers/go:dev-1.25-trixie

SHELL ["/bin/bash", "-c"]

RUN apt-get update && apt-get -y install --no-install-recommends \
    bash \
    make \
    git \
    curl \
    ca-certificates \
    findutils \
    coreutils

RUN curl -sSfL https://taskfile.dev/install.sh | sh -s -- -b /usr/local/bin
