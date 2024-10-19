image="docker.io/thealpha16/nc-test:latest"
container="nc-test"

docker rm -f $container
docker rmi -f $image

docker build -t $image .
docker run -d --name $container -e FLAG=dssfds -e USERNAME=prime -p 53581:53581 $image