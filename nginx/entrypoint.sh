#!/bin/sh
chmod +x entrypoint.sh

# Define file paths
CERT_FILE=/etc/nginx/certs/selfsigned.crt
KEY_FILE=/etc/nginx/certs/selfsigned.key

# Install OpenSSL
apk add --no-cache openssl

# Ensure the certificates directory exists
mkdir -p /etc/nginx/certs

# Check if certificate and key files exist
if [ ! -f "$CERT_FILE" ] || [ ! -f "$KEY_FILE" ]; then
    echo "Certificate or key does not exist. Generating new certificate and key..."
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout $KEY_FILE -out $CERT_FILE -subj "/C=US/ST=State/L=City/O=Organization/OU=OrgUnit/CN=localhost"
else
    echo "Using existing certificate and key."
fi

# Log NGINX version
nginx -v

# Start NGINX
nginx -g "daemon off;"