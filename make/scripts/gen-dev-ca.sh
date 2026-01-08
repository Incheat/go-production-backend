#!/usr/bin/env bash
set -euo pipefail

# === FIXED OUTPUT DIR (project-root/infra/security/tls/auth) ===
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
OUT_DIR="$PROJECT_ROOT/infra/security/tls/localhost"

if [[ ! -d "$OUT_DIR" ]]; then
  mkdir -p "$OUT_DIR"
fi

cd "$OUT_DIR"

echo "==> OUT_DIR=$OUT_DIR"

# Subjects (JP)
CA_SUBJ="/C=JP/ST=Tokyo/L=Tokyo/O=Dev/CN=Dev Root CA (JP)"
SRV_SUBJ="/C=JP/ST=Tokyo/L=Tokyo/O=Dev/CN=localhost"

# 1) Root CA (create if missing)
if [[ ! -f dev-root-ca.key || ! -f dev-root-ca.crt ]]; then
  echo "==> Generating dev root CA (JP)..."
  openssl genrsa -out dev-root-ca.key 4096
  openssl req -x509 -new -nodes \
    -key dev-root-ca.key \
    -sha256 -days 3650 \
    -subj "$CA_SUBJ" \
    -out dev-root-ca.crt
else
  echo "==> dev root CA already exists, skip"
fi

# 2) OpenSSL config for SAN
cat > localhost.cnf <<'CNF'
[ req ]
default_bits       = 2048
prompt             = no
default_md         = sha256
req_extensions     = req_ext
distinguished_name = dn

[ dn ]
C  = JP
ST = Tokyo
L  = Tokyo
O  = Dev
CN = localhost

[ req_ext ]
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = localhost
IP.1  = 127.0.0.1
CNF

cat > v3_server.cnf <<'CNF'
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = localhost
IP.1  = 127.0.0.1
CNF

# 3) Generate server key + cert (always re-generate)
echo "==> Generating localhost key/cert..."
rm -f localhost.key localhost.csr localhost.crt

openssl genrsa -out localhost.key 2048
openssl req -new \
  -key localhost.key \
  -out localhost.csr \
  -config localhost.cnf \
  -subj "$SRV_SUBJ"

openssl x509 -req \
  -in localhost.csr \
  -CA dev-root-ca.crt \
  -CAkey dev-root-ca.key \
  -CAcreateserial \
  -out localhost.crt \
  -days 825 -sha256 \
  -extfile v3_server.cnf

# 4) Verify
echo "==> Verify subject & SAN"
openssl x509 -in localhost.crt -noout -subject
openssl x509 -in localhost.crt -noout -text | grep -A2 "Subject Alternative Name"

echo "==> Files"
ls -la

echo ""
echo "==> Done"
echo "CA (import into Insomnia as CA): $OUT_DIR/dev-root-ca.crt"
echo "Server cert (Envoy):             $OUT_DIR/localhost.crt"
echo "Server key  (Envoy):             $OUT_DIR/localhost.key"
