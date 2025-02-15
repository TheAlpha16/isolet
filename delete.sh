#!/bin/bash

set -e

echo "[#] uninstalling isolet..."
helm uninstall isolet
echo "[+] isolet uninstalled"

echo "[#] uninstalling cert-manager..."
helm uninstall cert-manager -n cert-manager
echo "[+] cert-manager uninstalled"
