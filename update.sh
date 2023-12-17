docker rmi -f thealpha16/isolet-goapi

cd ./api
docker build -t thealpha16/isolet-goapi .
docker push thealpha16/isolet-goapi