ENV_FILE=$(dirname "$0")/../../api/.env
KUBECONFIG_FILE=$HOME/.kube/config

echo "[#] Please be warned!"
echo "[#] API will not be able to access the Kubernetes cluster if you are using gke-gcloud-auth-plugin."

if [ ! -f "$ENV_FILE" ]; then
    echo "Error: $ENV_FILE file not found!"
    exit 1
fi

if [ ! -f "$KUBECONFIG_FILE" ]; then
    echo "Error: $KUBECONFIG_FILE file not found!"
    exit 1
fi

export $(grep -v '^#' "$ENV_FILE" | xargs)

docker rm -f isolet-api
docker rmi -f isolet-api

cd $(dirname "$0")/../../api
docker build -t isolet-api .
docker run -d --name isolet-api \
    --network=host \
    --env-file $ENV_FILE \
    -v $KUBECONFIG_FILE:/root/.kube/config \
    isolet-api
