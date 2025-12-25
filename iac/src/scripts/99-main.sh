#!/bin/bash
# This script calls the other scripts in the correct order
# in order to setup the backend webserver instance

set -e

# Initial configuration (mainly env vars)
set_env_from_ec2_tag
get_ssm_configuration

# Setup backend dependencies
setup_persistence_volume &
setup_service_user &
wait

setup_certificate

# Download ana setup backend webserver
setup_backend

# Setup systemd service
setup_service

echo all done !!