#!/bin/bash

# Creates the service user/group used by the backend service
setup_service_user() {
    useradd --system --shell /bin/bash $APP_USER || echo "User $APP_USER already exists"
    groupadd $APP_GROUP || echo "Group $APP_GROUP already exists"
    usermod -a -G $APP_GROUP ec2-user
    usermod -a -G $APP_GROUP $APP_USER

    # SSH Setup
    PUB_KEY=`aws ssm get-parameter \
        --name q-backend-ssh-pub-key-${ENV} \
        --query "Parameter.Value" \
        --output text 2>/dev/null || echo error`
    if [ "$PUB_KEY" = "error" ] || [ -z "$PUB_KEY" ]; then
        echo "ERROR: Failed to get SSH public key from Parameter Store"
        exit 1
    fi
    echo "SSH public key: $PUB_KEY"
    mkdir -p /home/$APP_USER/.ssh
    chown $APP_USER:$APP_GROUP /home/$APP_USER/.ssh
    chmod 700 /home/$APP_USER/.ssh
    touch /home/$APP_USER/.ssh/authorized_keys
    chown $APP_USER:$APP_GROUP /home/$APP_USER/.ssh/authorized_keys
    chmod 600 /home/$APP_USER/.ssh/authorized_keys
    echo "$PUB_KEY" >> /home/$APP_USER/.ssh/authorized_keys
    # Allowing app user to stop/start service
    echo "%$APP_USER ALL= NOPASSWD: /bin/systemctl stop q-webserver-dev" >> /etc/sudoers.d/q-user
    echo "%$APP_USER ALL= NOPASSWD: /bin/systemctl start q-webserver-dev" >> /etc/sudoers.d/q-user

    echo "Service user $APP_USER created successfully."
    echo "Service group $APP_GROUP created successfully."
}

# This is suposed to be run on the backend instance startup
# It's objective is to mount a persistent AWS EBS volume to the instance
# In this volume database and other persistent data will be stored
setup_persistence_volume() {
  # The persistence volume is not available immediately after the instance is created
  # so we need to wait for it.
  echo "Waiting for EBS volume $PERSISTENCE_VOLUME_DEVICE_NAME to be available..."
  TIMEOUT=180
  ELAPSED=0
  while [ ! -b "$PERSISTENCE_VOLUME_DEVICE_NAME" ] && [ $ELAPSED -lt $TIMEOUT ]; do
      echo "Device $PERSISTENCE_VOLUME_DEVICE_NAME not found, waiting... (${ELAPSED}s)"
      sleep 3
      ELAPSED=$((ELAPSED + 3))
  done

  if [ ! -b "$PERSISTENCE_VOLUME_DEVICE_NAME" ]; then
      echo "ERROR: EBS volume $PERSISTENCE_VOLUME_DEVICE_NAME not available after ${TIMEOUT} seconds"
      exit 1
  fi

  echo "EBS volume $PERSISTENCE_VOLUME_DEVICE_NAME is available, proceeding with mount..."

  mkdir -p "$PERSISTENCE_MOUNT_DIR"
  # Might fail on first run. If so, run this command manually once: 'sudo mkfs -t ext4 /dev/xvdf'
  mount "$PERSISTENCE_VOLUME_DEVICE_NAME" "$PERSISTENCE_MOUNT_DIR"
  echo "EBS volume $PERSISTENCE_VOLUME_DEVICE_NAME mounted to $PERSISTENCE_MOUNT_DIR"
}

# Will get the security group id where temporary rules will be added
# so that certbot can verify the domain ownership
get_security_group_id() {
    SG_ID=$(aws ec2 describe-security-groups \
        --region "$REGION" \
        --filters "Name=tag:Id,Values=q-backend-main-sg" \
        --query "SecurityGroups[0].GroupId" \
        --output text 2>/dev/null) || {
        echo "ERROR: Failed to describe security groups"
        exit 1
    }
    if [ -z "$SG_ID" ] || [ "$SG_ID" = "None" ]; then
        echo "ERROR: Could not find security group with tag Id=q-backend-main-sg"
        exit 1
    fi
    export SG_ID
    echo "Using security group: $SG_ID"
}

# Sets up HTTPS certificate using certbot and http verification
setup_certificate() {
    echo "Setting up certificate..."

    get_certificate_from_s3 || {
        echo "Failed to get certificate from S3. Creating new certificate..."
        new_certbot_certificate
        save_certificate_to_s3
        echo "New certificate created and saved to S3."
    }

    echo "Certificate setup successfully."
}

new_certbot_certificate() {
    echo "Creating new certbot certificate..."

    get_security_group_id

    TMP_RULES=/tmp/certbot_temp_rules
    : > "$TMP_RULES"

    # Will remove the temporary rules required by certbot
    cleanup_certbot_sg() {
        if [ -s "$TMP_RULES" ]; then
            IDS=$(tr '\n' ' ' < "$TMP_RULES")
            aws ec2 revoke-security-group-ingress --region "$REGION" --security-group-rule-ids $IDS --group-id "$SG_ID" || true
            rm -f "$TMP_RULES"
        fi
    }
    trap cleanup_certbot_sg EXIT

    # Allows HTTP/HTTPS traffic from the internet to the instance
    # This is needed to allow certbot to verify the domain ownership
    open_ingress_ipv4_and_ipv6() {
        local port=$1
        # IPv4
        aws ec2 authorize-security-group-ingress \
            --region "$REGION" \
            --group-id "$SG_ID" \
            --ip-permissions "IpProtocol=tcp,FromPort=$port,ToPort=$port,IpRanges=[{CidrIp=0.0.0.0/0,Description=TEMP_CERTBOT_OPEN}]" \
            --output json | jq -r '.SecurityGroupRules[]? | .SecurityGroupRuleId // empty' >> "$TMP_RULES"

        # IPv6
        aws ec2 authorize-security-group-ingress \
            --region "$REGION" \
            --group-id "$SG_ID" \
            --ip-permissions "IpProtocol=tcp,FromPort=$port,ToPort=$port,Ipv6Ranges=[{CidrIpv6=::/0,Description=TEMP_CERTBOT_OPEN}]" \
            --output json | jq -r '.SecurityGroupRules[]? | .SecurityGroupRuleId // empty' >> "$TMP_RULES"
    }

    # Certbot uses port 80
    open_ingress_ipv4_and_ipv6 80

    dnf install -y python3 python3-devel augeas-devel gcc
    python3 -m venv /opt/certbot/
    /opt/certbot/bin/pip install --upgrade pip
    /opt/certbot/bin/pip install certbot
    ln -sf /opt/certbot/bin/certbot /usr/bin/certbot

    export CERTBOT_EMAIL="petersonvgama@gmail.com"
    if [ -z "$CERTBOT_EMAIL" ]; then
        echo "ERROR: CERTBOT_EMAIL not set; export CERTBOT_EMAIL for unattended issuance"; exit 1
    fi

    # Obtain cert via standalone HTTP-01 (binds :80) 
    certbot certonly \
        --standalone \
        --preferred-challenges http \
        --http-01-port 80 \
        -d "$CERTBOT_DOMAIN" \
        --non-interactive \
        --agree-tos \
        -m "$CERTBOT_EMAIL" \
        --keep-until-expiring

    chown -R $APP_USER:$APP_GROUP /etc/letsencrypt
    chmod -R 755 /etc/letsencrypt
    chmod g+rwx /etc/letsencrypt

    echo "New certbot certificate created successfully."
}

get_certificate_from_s3() {
    aws s3 cp s3://$SERVER_ASSETS_BUCKET/certs/$CERTBOT_DOMAIN/fullchain.pem /etc/letsencrypt/live/$CERTBOT_DOMAIN/fullchain.pem
}

save_certificate_to_s3() {
    aws s3 cp /etc/letsencrypt/live/$CERTBOT_DOMAIN/fullchain.pem s3://$SERVER_ASSETS_BUCKET/certs/$CERTBOT_DOMAIN/fullchain.pem
}

install_gh_cli() {
    type -p yum-config-manager >/dev/null || yum install yum-utils -y
    yum-config-manager --add-repo https://cli.github.com/packages/rpm/gh-cli.repo
    yum install gh -y
    echo "$GITHUB_TOKEN" > /tmp/github_token.txt
    gh auth login --with-token < /tmp/github_token.txt
    rm -f /tmp/github_token.txt
}