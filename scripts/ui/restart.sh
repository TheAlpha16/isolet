docker rm -f isolet-ui
docker rmi -f isolet-ui

cd $(dirname "$0")/../../ui
docker build -t isolet-ui .
docker run -d --name isolet-ui \
    --network=host \
    isolet-ui
