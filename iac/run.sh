#!/bin/bash

set -e
source $(dirname $0)/../bin/helpers.sh
COMMAND=$1
ENVIRONMENT=$2
VALID_COMMANDS=(init plan apply destroy)
CONFIG_FILE="config/$ENVIRONMENT.conf"
VARS_FILE="config/$ENVIRONMENT.tfvars"
BASE_FOLDER="src"

print_usage_and_exit() {
  echo "Usage: $(basename $0) <command> <environment>"
  echo "Commands: init, plan, apply, destroy"
  echo "Environments: dev, prod"
  exit 1
}

validate_command() {
  test -n "$COMMAND" || {
    echo >&2 "Command not specified. Aborting.";
    return 1;
  }

  if ! echo "${VALID_COMMANDS[@]}" | grep -q "$COMMAND"; then
    echo >&2 "Invalid command: $COMMAND. Valid commands: ${VALID_COMMANDS[@]}";
    return 2;
  fi
}

validate_dependencies() {
  command -v terraform >/dev/null 2>&1 || { 
    echo >&2 "Terraform is required but it's not installed. Aborting.";
    exit 1;
  }

  command -v aws >/dev/null 2>&1 || { 
    echo >&2 "AWS CLI is required but it's not installed. Aborting.";
    exit 2;
  }
}

ensure_state_bucket_exists() {
  local bucket_name
  local region
  bucket_name=$(awk -F '"' '/bucket/ {print $2}' "$CONFIG_FILE")
  region=$(awk -F '"' '/region/ {print $2}' "$CONFIG_FILE")
  aws s3api head-bucket --bucket "$bucket_name" --region "$region" >/dev/null 2>&1 && {
    echo >&2 "State bucket '$bucket_name' found."
  } || {
    echo >&2 "State bucket '$bucket_name' not found. Creating it..."
    local location_constraint
    if [ "$region" == "us-east-1" ]; then
      location_constraint=""
    else
      location_constraint="--create-bucket-configuration LocationConstraint=$region"
    fi
    aws s3api create-bucket --bucket "$bucket_name" --region "$region" --no-cli-pager $location_constraint >/dev/null || {
      echo >&2 "Failed to create state bucket '$bucket_name'. Aborting.";
      exit 8;
    }
    echo >&2 "State bucket '$bucket_name' created."
  }
}

run_command() {
  echo "Running command '$COMMAND' for environment '$ENVIRONMENT' ..."

  if [ "$COMMAND" == "init" ]; then
    rm -rf "$BASE_FOLDER/.terraform"
    rm -rf "$BASE_FOLDER/.terraform.lock.hcl"
    terraform -chdir="$BASE_FOLDER" init -backend-config="../$CONFIG_FILE"
    echo "Terraform initialized."

  elif [ "$COMMAND" == "plan" ]; then
    test -d "$BASE_FOLDER/.terraform" || {
      echo >&2 "Terraform directory '$BASE_FOLDER/.terraform' not found. Aborting.";
      echo >&2 "Run '$0 init $ENVIRONMENT' first.";
      exit 9;
    }
    terraform -chdir="$BASE_FOLDER" plan -var-file="../$VARS_FILE" -out=planfile
    echo "Terraform plan created."
    echo "Run '$(basename $0) apply $ENVIRONMENT' to apply the plan."

  elif [ "$COMMAND" == "apply" ]; then
    test -f "$BASE_FOLDER/planfile" || {
      echo >&2 "Plan file 'planfile' not found. Aborting.";
      echo >&2 "Run '$0 plan $ENVIRONMENT' first.";
      exit 10;
    }
    terraform -chdir="$BASE_FOLDER" apply --auto-approve planfile
    echo "Terraform plan applied."

  elif [ "$COMMAND" == "destroy" ]; then
    terraform -chdir="$BASE_FOLDER" destroy -var-file="../$VARS_FILE"
    echo "Terraform destroyed."
  else
    echo "Invalid command: $COMMAND. Valid commands: ${VALID_COMMANDS[@]}";
    exit 11;
  fi

  echo "Command '$COMMAND' for environment '$ENVIRONMENT' completed successfully."
}

run() {
  validate_command || print_usage_and_exit
  validate_environment $ENVIRONMENT || print_usage_and_exit
  check_aws_profile $ENVIRONMENT || print_usage_and_exit
  set_aws_profile $ENVIRONMENT
  validate_dependencies
  ensure_state_bucket_exists
  run_command
}

run
