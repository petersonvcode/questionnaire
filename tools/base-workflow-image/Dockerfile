FROM ubuntu:24.04

ARG ARCH=amd64
ARG NVM_VERSION=0.40.3
ARG GO_VERSION
ARG NODE_VERSION
ARG PNPM_VERSION

RUN if [ -z "$ARCH" ]; then \
    echo "ARCH not set. Aborting."; \
    exit 1; \
elif [ -z "$GO_VERSION" ]; then \
    echo "GO_VERSION not set. Aborting."; \
    exit 1; \
elif [ -z "$NODE_VERSION" ]; then \
    echo "NODE_VERSION not set. Aborting."; \
    exit 1; \
elif [ -z "$PNPM_VERSION" ]; then \
    echo "PNPM_VERSION not set. Aborting."; \
    exit 1; \
elif [ -z "$NVM_VERSION" ]; then \
    echo "NVM_VERSION not set. Aborting."; \
    exit 1; \
fi

RUN apt-get update && apt-get upgrade -y
RUN apt-get install -y curl build-essential

# Installing golang
RUN curl -fsSL https://go.dev/dl/go${GO_VERSION}.linux-${ARCH}.tar.gz | tar -xz -C /usr/local
ENV PATH="/usr/local/go/bin:${PATH}"

# Installing nvm
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Create a script file sourced by both interactive and non-interactive bash shells
ENV BASH_ENV="/root/.bash_env"
RUN touch "${BASH_ENV}"
RUN echo '. "${BASH_ENV}"' >> ~/.bashrc
RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v${NVM_VERSION}/install.sh | PROFILE="${BASH_ENV}" bash
RUN nvm install ${NODE_VERSION}
RUN nvm use ${NODE_VERSION}
RUN nvm alias default ${NODE_VERSION}
RUN npm install -g corepack@latest
RUN corepack enable
RUN corepack install --global pnpm@${PNPM_VERSION}