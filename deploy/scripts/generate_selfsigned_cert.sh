#!/bin/bash

CERT_DIR="./deploy/certs"

mkdir -p = "$CERT_DIR"

echo "Генерация самоподписанного SSL-сертификата"

openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout "$CERT_DIR/server.key" \
  -out "$CERT_DIR/server.crt" \
  -subj "/C=RU/ST=UFA/L=UFA/O=DevOpsTeam/OU=IT/CN=localhost"

echo "Сертификаты успешно созданы в папке $CERT_DIR"
