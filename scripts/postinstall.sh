#!/bin/sh
set -e

# Create bloggulus group (if it doesn't exist)
if ! getent group bloggulus >/dev/null; then
    groupadd --system bloggulus
fi

# Create bloggulus user (if it doesn't exist)
if ! getent passwd bloggulus >/dev/null; then
    useradd                                \
        --system                           \
        --gid bloggulus                    \
        --shell /usr/sbin/nologin          \
        --comment "bloggulus feed reader"  \
        bloggulus
fi

# Update config file permissions (idempotent)
chown root:bloggulus /etc/bloggulus.conf
chmod 0640 /etc/bloggulus.conf

# Reload systemd to pickup bloggulus.service
systemctl daemon-reload
