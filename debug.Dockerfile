# Dockerfile.debug
FROM ubuntu:22.04

ARG GO_VERSION="1.23.4"
ARG DEBIAN_FRONTEND="noninteractive"

# Create the 'runner' user/group to match GitHub's default
RUN useradd -m runner

RUN apt-get update && apt-get install -y \
    ca-certificates wget git \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*

# Install Go exactly as you do on GitHub
RUN wget https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz -O /tmp/go.tgz \
 && tar -C /usr/local -xzf /tmp/go.tgz \
 && rm /tmp/go.tgz

# Attempt to mirror GitHub's environment
ENV HOME="/home/runner"
ENV GOPATH="/home/runner/go"
ENV GOCACHE="/home/runner/.cache/go-build"
ENV GOROOT="/usr/local/go"
ENV PATH="/usr/local/go/bin:${PATH}"

# Create the standard GitHub Actions work directory
RUN mkdir -p /home/runner/work/chainlink/chainlink
WORKDIR /home/runner/work/chainlink/chainlink

# Switch to runner user (optional, but more faithful to GH environment)
USER runner

CMD ["bash"]
