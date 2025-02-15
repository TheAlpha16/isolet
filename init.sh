#!/bin/bash

set -e

echo "[#] installing cert-manager..."
helm repo add jetstack https://charts.jetstack.io --force-update

helm install \
  cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --version v1.17.0 \
  --set crds.enabled=true

echo "[+] cert-manager installed"

echo "[#] installing isolet..."
helm install isolet $(dirname $0)/charts
echo "[+] isolet installed"
