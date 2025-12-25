#!/bin/bash

set -e
source $(dirname $0)/../../bin/helpers.sh
ENVIRONMENT=$1

print_usage_and_exit() {
    echo "Usage: $(basename $0) <environment>"
    echo "Environments: dev, prod"
    exit 1
}

run() {
    validate_environment $ENVIRONMENT || print_usage_and_exit

    GO_VERSION=$(grep go ../../backend/go.mod | head -1 | cut -d' ' -f2)
    echo "Go version: $GO_VERSION"
    test -n "$GO_VERSION" || {
        echo "Go version not found. Aborting.";
        exit 1;
    }

    NODE_VERSION=$(cat ../../frontend/.nvmrc | tr -d 'vV')
    echo "Node version: $NODE_VERSION"
    test -n "$NODE_VERSION" || {
        echo "Node version not found. Aborting.";
        exit 1;
    }

    PNPM_VERSION=$(cat ../../frontend/package.json | jq -r '.packageManager' | cut -d'@' -f2)
    echo "Pnpm version: $PNPM_VERSION"
    test -n "$PNPM_VERSION" || {
        echo "Pnpm version not found. Aborting.";
        exit 1;
    }

    NVM_VERSION=$(cat ../../frontend/.nvmrc | tr -d 'vV')
    echo "Nvm version: $NVM_VERSION"
    test -n "$NVM_VERSION" || {
        echo "Nvm version not found. Aborting.";
        exit 1;
    }

    docker build \
        -t q-workflows-${ENVIRONMENT} \
        --build-arg GO_VERSION="$GO_VERSION" \
        --build-arg NODE_VERSION="$NODE_VERSION" \
        --build-arg PNPM_VERSION="$PNPM_VERSION" \
        --build-arg NVM_VERSION="$NVM_VERSION" \
        .

    echo "Docker image built successfully."
}

run