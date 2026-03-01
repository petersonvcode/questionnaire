#!/bin/bash

set -e
IP=72.62.106.70
KEY="~/.ssh/hostinger_vps"
USER=root

ssh -i "$KEY" "$USER@$IP"