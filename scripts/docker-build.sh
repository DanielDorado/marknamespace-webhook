#!/usr/bin/env bash

docker build . -f docker/Dockerfile -t default-route-openshift-image-registry.apps-crc.testing/openshift/marknamespace:latest
docker login default-route-openshift-image-registry.apps-crc.testing --username kubeadmin --password $(oc whoami -t)
docker push default-route-openshift-image-registry.apps-crc.testing/openshift/marknamespace:latest
