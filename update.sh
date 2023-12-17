#!/bin/bash

set +x

echo "Select the resource to update: "
echo "1. API"
echo "2. UI"
echo "3. proxy"
echo ""
echo -n "> "
read choice

resource=""

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
    *)
        echo ""
        echo "Invalid choice"
        echo "Valid options: 1, 2, 3"
        echo "1. API   - rebuild isolet-goapi image"
        echo "2. UI    - rebuild isolet-ui image"
        echo "3. proxy - rebuild isolet-proxy image"
        echo ""
        echo "[-] exiting"
        exit
    ;;
esac

docker images | grep thealpha16/isolet-${resource}
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
docker build -t thealpha16/isolet-${resource}:${version} .

case $tag in
    "yes")
        docker rmi -f thealpha16/isolet-$resource:latest
        docker build -t thealpha16/isolet-${resource}:latest .
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

docker push thealpha16/isolet-${resource}:${version}
docker push thealpha16/isolet-${resource}:latest