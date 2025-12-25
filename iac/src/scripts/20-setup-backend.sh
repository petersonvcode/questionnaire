#!/bin/bash

setup_backend() {
    echo "Setting up backend..."
    USER=q-user
    GROUP=apps

    # Directory for the application files
    mkdir -p $APP_DIR
    chown -R $USER:$GROUP $APP_DIR
    chmod -R 755 $APP_DIR
    chmod g+rwx $APP_DIR

    # Directory for the database files
    mkdir -p $DATABASE_DIR
    chown -R $USER:$GROUP $DATABASE_DIR
    chmod -R 755 $DATABASE_DIR
    chmod g+rw $DATABASE_DIR

    # Download binary files

    echo "Backend setup successfully."
}

setup_service() {
    echo oi
}