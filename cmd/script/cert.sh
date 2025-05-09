#!/bin/bash

CERT_DIR="certs"
KEY_FILE="$CERT_DIR/server.key"
CERT_FILE="$CERT_DIR/server.crt"

mkdir -p "$CERT_DIR"

openssl req -x509 -newkey rsa:4096 -sha256 -days 365 -nodes \
  -keyout "$KEY_FILE" \
  -out "$CERT_FILE" \
  -subj "/C=US/ST=State/L=City/O=Organization/OU=Dev/CN=localhost"

echo "✅ TLS certificate and key generated:"
echo "  Key:  $KEY_FILE"
echo "  Cert: $CERT_FILE"
