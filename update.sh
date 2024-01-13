#!/bin/bash

set +x

echo "Select the resource to update: "
echo "1. API"
echo "2. UI"
echo "3. proxy"
echo "4. ripper"
echo ""
echo -n "> "
read choice

resource=""
registry="thealpha16"

case $choice in
    "1")
        resource="goapi"
    ;;
    "2")
        resource="ui"
    ;;
    "3")
        resource="proxy"
    ;;
    "4")
        resource="ripper"
    ;;
    *)
        echo ""
        echo "Invalid choice"
        echo "Valid options: 1, 2, 3"
        echo "1. API    - rebuild isolet-goapi image"
        echo "2. UI     - rebuild isolet-ui image"
        echo "3. proxy  - rebuild isolet-proxy image"
        echo "4. ripper - rebuild isolet-ripper image"
        echo ""
        echo "[-] exiting"
        exit
    ;;
esac

docker images | grep ${registry}/isolet-${resource}
echo ""

echo "Choose version to tag"
echo -n "> "
read version
echo ""

echo "tag this to latest version? (yes/no)"
echo -n "> "
read tag
echo ""

echo "version: $version tag: $tag"
echo "Confirm modification (yes/no)?"
echo -n "> "
read confirm
echo ""

case $confirm in
    "yes")
        echo "[+] proceeding"
    ;;
    "no")
        echo "[-] exiting"
        exit
    ;;
    *)
        echo "[-] Choose 'yes' or 'no'"
        echo ""
        echo "[-] exiting"
        exit
    ;;
esac

cd ./${resource}
docker buildx build --tag ${registry}/isolet-${resource}:${version} --platform linux/arm64/v8,linux/amd64 --builder bob --push .

case $tag in
    "yes")
        docker rmi -f ${registry}/isolet-$resource:latest
        docker buildx build --tag ${registry}/isolet-${resource}:latest --platform linux/arm64/v8,linux/amd64 --builder bob --push .
    ;;
    "no")
        echo "[*] Not removing latest tag"
    ;;
    *)
        echo "[-] Choose 'yes' or 'no' for tagging"
        echo ""
        echo "[-] exiting"
        exit
    ;;
esac