# isolet
Isolet is a framework to deploy linux wargames like [Bandit](https://overthewire.org/wargames/bandit/). It uses pre-configured templates to provide isolated instance using kubernetes pods for each user.

## Configuration
You can customize the application using environment variables passed to the deployments. All the options are available in [configuration](/kubernetes/configuration)

### General

#### `WARGAME_NAME`
Name of the wargame to be deployed

#### `PUBLIC_URL`
URL of the deployed application. Required for email verification

#### `PROXY_SERVER_NAME`
Domains and subdomains to be added under server_name directive in nginx proxy.

> [!note]
> Check out the nginx documentation for format [server_name](https://nginx.org/en/docs/http/server_names.html)

#### `INSTANCE_HOSTNAME`
Domain name for accessing the spawned instances

#### `IMAGE_REGISTRY_PREFIX`
Default registry for pulling challenge images. It is prefixed to `level` to get final repo url to be 

> [!important]
> If you specify your prefix to be `docker.io/thealpha16/` the final image path for level 1 will be `docker.io/thealpha16/level1`

#### `DISCORD_FRONTEND`
Boolean to determine whether API needs `/auth` routes to be setup. If it is true, API will not authenticate the request.

> [!warning]
> This option should be used only when the API is not exposed to the public and the request is being forwarded by some other application which is properly authenticating

#### `KUBECONFIG_FILE_PATH`
Path to the kubernetes config file to access cluster from outside

> [!note]
> for more information, check out [cluster access](https://kubernetes.io/docs/tasks/access-application-cluster/access-cluster/)

#### `UI_URL`
host for frontend in case it exists. If kubernetes is being used for deployment, you can specify URL to be

```
<SERVICE_NAME_OF_UI>.<NAMESPACE>.svc.cluster.local
```

> [!info]
> for more information, head over to [dns for pods](https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/)

### Instance

#### `INSTANCE_NAMESPACE`
Namespace for deploying the user requested pods.

#### `CONCURRENT_INSTANCES` 
Number of concurrent pods that user can spawn.

#### `TERMINATION_PERIOD` 
Time in seconds to be given to the pod for graceful shutdown.

#### `INSTANCE_TIME`
Time in minutes to be added in the pod annotations after which ripper will remove the instance

#### `MAX_INSTANCE_TIME`
Time in minutes the user can extend the instance

#### `CPU_REQUEST`
Number of cores to be reserved for the pod

#### `CPU_LIMIT`
Maximum number of cores the pod can consume

#### `MEMORY_REQUEST`
Amount of memory to be reserved for the pod

#### `MEMORY_LIMIT`
Maximum amount of memory the pod can use

#### `DISK_REQUEST`
Disk space to be reserved for the pod

#### `DISK_LIMIT`
Maximum disk space the pod can utilize

> [!note]
> for more information regarding kubernetes resources, check out [resources](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/)

## API routes
| route | methods | parameters | response | sample |
|:---:|:---:|:---:|:---:|:---:|
| /api/challs | GET | NONE | challenges | [{"chall_id":0, "level":1, "name":"demo", "prompt":"solve it", "tags":["ssh", "cat"]}] |
| /api/launch | POST | chall_id, userid, level | status | {"status": "success", "message": "3b369c0b1fd5419b2f81da89cf5480d2 32747"} |
| /api/stop | POST | userid, level | status | {"status": "failure", "message": "User does not exist"} |
| /api/submit | POST | userid, level, flag | status | {"status": "failure", "message": "Flag copy detected. Incident reported!"} |
| /api/status | GET | NONE | instances | {"status": "success", "message": "[{"userid":123614343, "level":1, "password":"8f1ee93113affe32078c", "port":"32134"}]"}
