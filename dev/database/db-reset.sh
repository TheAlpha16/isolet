POSTGRES_USER=postgres
POSTGRES_PASSWORD=1c25d2068cdd9c3452013dc106e3be8a

docker rm -f test_database
docker rmi -f isolet-database

cd $(dirname "$0")/../../database
docker build -t isolet-database .
docker run -d --name test_database -p 5432:5432 -e POSTGRES_USER=$POSTGRES_USER -e POSTGRES_PASSWORD=$POSTGRES_PASSWORD isolet-database

sleep 3
cd $(dirname "$0")/../scripts/database
PGPASSWORD="$POSTGRES_PASSWORD" psql -h localhost -p 5432 -U $POSTGRES_USER -d postgres -f "./fake_data.sql"