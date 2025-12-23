#!/bin/bash

set -e

VALID_ENVIRONMENTS=(dev prod)

validate_environment() {
    environment=$1
    if ! echo "${VALID_ENVIRONMENTS[@]}" | grep -q "$environment"; then
        echo >&2 "Invalid environment: $environment. Valid environments: ${VALID_ENVIRONMENTS[@]}";
        exit 1;
    fi
}

check_aws_profile() {
    environment=$1
    test -n "$environment" || {
        echo >&2 "Environment not specified. Aborting.";
        return 1;
    }

    profile=q-${environment}
    profiles=$(aws configure list-profiles)
    if ! echo "$profiles" | grep -q "$profile"; then
        echo >&2 "AWS credentials for $profile not found. Aborting.";
        echo >&2 "Run 'aws configure --profile $profile' to configure the credentials.";
        return 2;
    fi
}

set_aws_profile() {
  environment=$1
  test -n "$environment" || {
    echo >&2 "Environment not specified. Aborting.";
    return 1;
  }
  export AWS_PROFILE="q-${environment}"
  echo "Using AWS_PROFILE: $AWS_PROFILE"
}

get_parameter_store_value() {
  local key=$1
  local value
  value=$(aws ssm get-parameter \
    --name "$key" \
    --query "Parameter.Value" \
    --output text 2>/dev/null) || {
    echo "ERROR: Failed to read SSM parameter '$key'";
    exit 2;
  }
  if [[ -v "$value" || "$value" == "None" ]]; then
    echo "ERROR: No value found for SSM parameter '$key'"
    exit 3
  fi
  echo "$value"
}