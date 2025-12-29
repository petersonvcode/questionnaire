#!/bin/bash

setup_backend() {
    echo "Setting up backend..."

    # Directory for the application files
    mkdir -p $APP_DIR
    chown -R $APP_USER:$APP_GROUP $APP_DIR
    chmod -R 755 $APP_DIR
    chmod g+rwx $APP_DIR

    # Directory for the logs files
    mkdir -p $LOGS_DIR
    chown -R $APP_USER:$APP_GROUP $LOGS_DIR
    chmod -R 755 $LOGS_DIR
    chmod g+rw $LOGS_DIR

    # Directory for the database files
    mkdir -p $DATABASE_DIR
    chown -R $APP_USER:$APP_GROUP $DATABASE_DIR
    chmod -R 755 $DATABASE_DIR
    chmod g+rw $DATABASE_DIR

    # Download binary files
    gh release download latest \
        -R petersonvcode/questionnaire \
        -p webserver \
        -O "$APP_DIR/webserver"
    echo "Downloaded webserver binary to $APP_DIR/webserver"
    chmod +x "$APP_DIR/webserver"
    chown "$APP_USER:$APP_GROUP" "$APP_DIR/webserver"

    echo "Backend setup successfully."
}

setup_service() {
    LOG_FILE=${LOGS_DIR}/webserver-${ENV}.log

    cat > /etc/systemd/system/q-webserver-${ENV}.service << EOF
[Unit]
Description=Questionnaire Backend Service ${ENV}
After=network.target local-fs.target
Requires=local-fs.target

[Service]
Type=simple
User=${APP_USER}
Group=${APP_GROUP}
WorkingDirectory=${APP_DIR}
ExecStart=${APP_DIR}/webserver
Restart=always
RestartSec=5

StandardOutput=append:${LOG_FILE}
StandardError=append:${LOG_FILE}

NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=${PERSISTENCE_MOUNT_DIR} ${APP_DIR}

PrivateTmp=true
ProtectKernelTunables=true
ProtectControlGroups=true
RestrictRealtime=true
EOF

  systemctl daemon-reload
  systemctl start "q-webserver-${ENV}"

  echo "To start the service: systemctl start q-webserver-${ENV}"
  echo "To check status: systemctl status q-webserver-${ENV}"
  echo "Logs will be written to: ${LOG_FILE}"
  echo "To follow logs: tail -f ${LOG_FILE}"
  echo "To view via journald: journalctl -u q-webserver-${ENV} -f -n 100"
}