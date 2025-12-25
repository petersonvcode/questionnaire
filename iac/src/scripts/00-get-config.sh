#!/bin/bash

# Gets the environment from the EC2 tag 'environment'
# Gets the region from the EC2 instance identity document
# Will set those values as the environment variables ENV and REGION
set_env_from_ec2_tag() {
    EC2_TOKEN=`curl -s "http://169.254.169.254/latest/api/token" \
        -X PUT \
        -H "X-aws-ec2-metadata-token-ttl-seconds: 21600"`

    ENV=`curl -s http://169.254.169.254/latest/meta-data/tags/instance/environment \
        -H "X-aws-ec2-metadata-token: $EC2_TOKEN"`
    test -n "$ENV" || {
        echo "Environment not found. Aborting."
        exit 1
    }
    echo "Environment: $ENV"
    export ENV

    REGION=`curl -s http://169.254.169.254/latest/dynamic/instance-identity/document \
        -H "X-aws-ec2-metadata-token: $EC2_TOKEN" \
        | jq -r '.region'`
    test -n "$REGION" || {
        echo "Region not found. Aborting."
        exit 1
    }
    echo "Region: $REGION"
    export REGION

    APP_USER=q-user
    export APP_USER
    APP_GROUP=apps
    export APP_GROUP
}

# Gets the configuration from Parameter Store
# Will set the following environment variables:
# - BACKEND_DOMAIN
# - PERSISTENCE_VOLUME_DEVICE_NAME
# - PERSISTENCE_MOUNT_DIR
# - DATABASE_DIR
# Requires the environment variable ENV to be set
get_ssm_configuration() {
    config=`aws ssm get-parameter \
        --name "q-backend-conf-${ENV}" \
        --query "Parameter.Value" \
        --output text 2>/dev/null`
    test -n "$config" || {
        echo "Configuration not found. Aborting."
        exit 1
    }
    echo "Configuration: $config"

    BACKEND_DOMAIN=$(echo "$config" | jq -r '.backend_domain')
    PERSISTENCE_VOLUME_DEVICE_NAME=$(echo "$config" | jq -r '.persistence_volume_device_name')
    PERSISTENCE_MOUNT_DIR=$(echo "$config" | jq -r '.persistence_mount_dir')
    DATABASE_DIR=$(echo "$config" | jq -r '.database_dir')
    APP_DIR=$(echo "$config" | jq -r '.app_dir')

    # Ensuring all variables are set
    test -n "$BACKEND_DOMAIN" || {
        echo "BACKEND_DOMAIN not found. Aborting."
        exit 1
    }
    export BACKEND_DOMAIN
    test -n "$PERSISTENCE_VOLUME_DEVICE_NAME" || {
        echo "PERSISTENCE_VOLUME_DEVICE_NAME not found. Aborting."
        exit 1
    }
    export PERSISTENCE_VOLUME_DEVICE_NAME
    test -n "$PERSISTENCE_MOUNT_DIR" || {
        echo "PERSISTENCE_MOUNT_DIR not found. Aborting."
        exit 1
    }
    export PERSISTENCE_MOUNT_DIR
    test -n "$DATABASE_DIR" || {
        echo "DATABASE_DIR not found. Aborting."
        exit 1
    }
    export DATABASE_DIR
    test -n "$APP_DIR" || {
        echo "APP_DIR not found. Aborting."
        exit 1
    }
    export APP_DIR
    
    echo Set configuration as environment variables
}
