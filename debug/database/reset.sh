docker rm -f test_database
docker rmi -f isolet-database

cd ../../database
docker build -t isolet-database .
docker run -d --name test_database -p 5432:5432 isolet-database
