#!/bin/bash

set -e
source $(dirname $0)/helpers.sh

print_usage_and_exit() {
    echo "Usage: $(basename $0) <environment>"
    echo "Environments: dev, prod"
    exit 1
}

ENVIRONMENT=$1
IP=72.62.106.70
KEY="~/.ssh/hostinger_vps_ed25519"
SSH_USER=root
SERVICE_USER=q-${ENVIRONMENT}-service

validate_environment $ENVIRONMENT || print_usage_and_exit

echo "Creating service user: $SERVICE_USER"
ssh -i "$KEY" "$SSH_USER@$IP" 'bash -s' <<EOF
id $SERVICE_USER >/dev/null 2>&1 && {
    echo "Service user $SERVICE_USER already exists"
    exit 0
}
useradd $SERVICE_USER --system --no-create-home -s /sbin/nologin --user-group
echo "Service user $SERVICE_USER created successfully"
EOF
