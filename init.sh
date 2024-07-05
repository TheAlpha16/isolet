set +x

cd ./kubernetes/init
kubectl apply -f instance-namespace.yml
kubectl apply -f db-volume.yml
kubectl apply -f network-policy.yml
kubectl apply -f roles.yml

# initialize secrets and configs
cd ../configuration
kubectl apply -f app-secrets.yml
kubectl apply -f app-config.yml

# initialize applications
cd ../definition
kubectl apply -f db-main.yml
kubectl apply -f ui-main.yml
kubectl apply -f api-main.yml
kubectl apply -f proxy-main.yml
kubectl apply -f ripper-main.yml
