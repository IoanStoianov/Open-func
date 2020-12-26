# Open-func
Open source runtime similar to AWS serverless lambda and Azure functions.

### Prerequisite
1. [Kubernetes CLI](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
2. [Minikube](https://minikube.sigs.k8s.io/docs/start/)
3. [Docker](https://docs.docker.com/engine/install/)

### Setup
`cd open-func`

`docker build . -t open-func`

`kubectl create -f open-func.yml`

`eval $(minikube -p minikube docker-env)` has to be run in every new terminal window before you build an image. An alternative would be to put it into your .profile file.
