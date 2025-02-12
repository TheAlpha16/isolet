POSTGRES_HOST=0.0.0.0
POSTGRES_USER=postgres
POSTGRES_PASSWORD=6a78db13fd1f73cc781380bf6ca408
POSTGRES_DATABASE=postgres
SESSION_SECRET=5b22ff96a78db13fd1f73cc78139112e9e180bf6ca4087aac24a609cc2abd657

docker rm -f isolet-socky
docker rmi -f isolet-socky

cd $(dirname "$0")/../../socky
docker build -t isolet-socky .
docker run -d --name isolet-socky \
    --network=host \
    -e POSTGRES_USER=$POSTGRES_USER \
    -e POSTGRES_PASSWORD=$POSTGRES_PASSWORD \
    -e POSTGRES_HOST=$POSTGRES_HOST \
    -e POSTGRES_DATABASE=$POSTGRES_DATABASE \
    -e SESSION_SECRET=$SESSION_SECRET \
    isolet-socky
