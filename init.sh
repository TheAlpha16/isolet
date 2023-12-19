set +x

# initialize secrets and configs
cd ./kubernetes/configuration
kubectl apply -f app-secrets.yml.local
kubectl apply -f app-config.yml
kubectl apply -f roles.yml

# initialize applications
cd ../definition
kubectl apply -f db-main.yml
kubectl apply -f ui-main.yml
kubectl apply -f goapi-main.yml
kubectl apply -f proxy-main.yml