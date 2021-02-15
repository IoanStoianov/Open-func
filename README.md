# Open-func
Open source runtime similar to AWS serverless lambda and Azure functions.

## Prerequisites
1. [Kubernetes CLI](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
2. [Minikube](https://minikube.sigs.k8s.io/docs/start/)
3. [Docker](https://docs.docker.com/engine/install/)

## Setup

`make build` - builds server image in cluster

`kubectl apply -f deployments/default-rbac.yaml` - adds cluster admin rights to the server

`kubectl apply -f deployments/open-func.yaml` - deploys Open-func in the cluster

## Adding the example images to the cluster

`eval $(minikube -p minikube docker-env)` - has to be run in every new terminal window before you build an image inside the cluster. An alternative would be to put it into your .profile file.

`docker build examples/<dockerfile-location> -t <image-name>`
