#!/bin/bash

set -e
source $(dirname $0)/../bin/helpers.sh
ENVIRONMENT=$1

print_usage_and_exit() {
    echo "Usage: $(basename $0) <environment>"
    echo "Environments: dev, prod"
    exit 1
}

prepare_environment() {
    local website=$(get_parameter_store_value "q-${ENVIRONMENT}-website")
    echo "Website: $website"
    test -n "$website" || {
        echo "Website details not found in Parameter Store. Aborting.";
        exit 1;
    }
    DISTRIBUTION_ID=$(echo "$website" | jq -r '.distribution_id')
    BUCKET_NAME=$(echo "$website" | jq -r '.bucket_name')
    ADDRESSES=$(echo "$website" | jq -r '.addresses')
    echo "Distribution ID: $DISTRIBUTION_ID"
    echo "Bucket Name: $BUCKET_NAME"
}

build_frontend() {
    echo "Building frontend..."
    pnpm build
}

deploy_frontend() {
    echo "Deploying frontend..."
    echo "Deleting old files..."
    aws s3 rm s3://${BUCKET_NAME} --recursive
    echo "Uploading new files..."
    aws s3 cp dist s3://${BUCKET_NAME} --recursive
    aws cloudfront create-invalidation --distribution-id ${DISTRIBUTION_ID} --paths "/*" --no-cli-pager
}

success_message() {
    echo "Frontend deployed successfully."
    echo "Addresses: $ADDRESSES"
}

run() {
    validate_environment $ENVIRONMENT || print_usage_and_exit
    check_aws_profile $ENVIRONMENT || print_usage_and_exit
    set_aws_profile $ENVIRONMENT
    prepare_environment
    build_frontend
    deploy_frontend
    success_message
}

run