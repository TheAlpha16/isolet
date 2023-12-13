set +x

# initialize secrets and configs
cd ./kubernetes/configuration
kubectl apply -f db-secrets.yml
kubectl apply -f api-secrets.yml
kubectl apply -f api-config.yml
kubectl apply -f roles.yml

# initialize database
cd ../db-def
kubectl apply -f db-volume.yml
kubectl apply -f db-main.yml

# initialize api
cd ../api-def
kubectl apply -f go-main.yml