#!/bin/bash
echo "- Generating self-signed certificate, please wait"
cd cert && go run /usr/local/go/src/crypto/tls/generate_cert.go -rsa-bits=2048 --host=localhost
echo "- certificate generated"