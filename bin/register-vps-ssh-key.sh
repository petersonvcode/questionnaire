#!/bin/bash

set -e

# Replace these variables as needed
IP=72.62.106.70
PUB_KEY="~/.ssh/hostinger_vps_ed25519.pub"
USER=root
PW=$(cat /home/pet/Work/repos/pessoal/vps-root-pw)

EXPANDED_PUB_KEY="${PUB_KEY/#\~/$HOME}" # Handle ~ in path
if [ ! -f "$EXPANDED_PUB_KEY" ]; then
    echo "Public key file not found: $PUB_KEY"
    echo "Please generate the public key file and try again."
    echo "ssh-keygen -t ed25519 -f $PUB_KEY -C 'your@email.com'"
    exit 1
fi

if [ -z "$PW" ]; then
    echo "Password not found!"
    exit 1
fi

echo use password: "$PW"
ssh-copy-id -i "$EXPANDED_PUB_KEY" "$USER@$IP"

