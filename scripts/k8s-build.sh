#!/usr/bin/env bash

mkdir build -p

if ! [[ -d build/cert ]]; then
    echo "Directory  build/cert/ not found."
    echo "Creating files in build/cert."
    mkdir -p build/cert
    ./scripts/webhook-create-signed-cert.sh build/cert
fi

# deployment
cp manifests/marknamespace.deployment.yml build/marknamespace.deployment.yml

# service
cp manifests/marknamespace.service.yml build/marknamespace.service.yml

# configmap
cat manifests/marknamespace.configmap.yml \
    <(cat manifests/marknamespace-conf.yml | sed 's/^/    /') > \
    build/marknamespace.configmap.yml

# secret
sed "s/<CRT-PEM>/$(cat build/cert/webhook-server-tls.crt | base64 -w0)/" \
    manifests/marknamespace.secret.tpl.yml \
    | sed "s/<KEY-PEM>/$(cat build/cert/webhook-server-tls.key | base64 -w0)/" \
    > build/marknamespace.secret.tpl.yml

# webhook
sed "s/<CA-BUNDLE>/$(cat build/cert/ca.crt | base64 -w0)/" \
    manifests/marknamespace.webhook.tpl.yml \
    > build/marknamespace.webhook.yml
