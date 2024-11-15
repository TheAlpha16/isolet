#!/bin/bash

set +x
registry="docker.io/thealpha16"

case $1 in
    "-y")
        shift
        tag="yes"
    ;;
    *)
        tag="-"
    ;;
esac

version=$2

case $1 in
    "api")
        resource="api"
    ;;
    "ui")
        resource="ui"
    ;;
    "ripper")
        resource="ripper"
    ;;
    *)
        echo ""
        echo "[-] invalid resource"
        echo "[#] valid options: api, ui, ripper"
        echo ""
        echo "[#] usage: build.sh [-y] <resource> <version>"
        echo "[-] exiting"
        exit
    ;;
esac

case $version in
    "")
        echo ""
        echo "[-] version not provided"
        echo "[#] usage: build.sh [-y] <resource> <version>"
        echo ""
        echo "[-] exiting"
        exit
    ;;
esac

docker images | grep ${registry}/isolet-${resource}

case $tag in
    "-")
        echo -n "[?] tag this to latest version? (Y/n) "
        read tag

        case $tag in
            "N"|"n"|"no"|"No"|"NO")
                tag="no"
            ;;
            *)
                tag="yes"
            ;;
        esac
    ;;
    "yes")
        echo ""
        echo "[#] tagging to latest version"
    ;;
esac

cd $(dirname "$0")/../../${resource}
docker buildx build --tag ${registry}/isolet-${resource}:${version} --platform linux/amd64 --push .

case $tag in
    "yes")
        docker buildx build --tag ${registry}/isolet-${resource}:latest --platform linux/amd64 --push .
    ;;
    "no")
        echo "[#] not removing latest tag"
    ;;
esac