# Open-func
Open source runtime similar to AWS serverless lambda and Azure functions.

### Prerequisite
1. [Kubernetes CLI](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
2. [Minikube](https://minikube.sigs.k8s.io/docs/start/)
3. [Docker](https://docs.docker.com/engine/install/)

### Setup
`cd open-func`

`kubectl apply -f fabric8-rbac.yaml` this add permission so that open-func can deploy containers

`eval $(minikube -p minikube docker-env)` has to be run in every new terminal window before you build an image. An alternative would be to put it into your .profile file.

`docker build . -t open-func`

`kubectl create -f open-func.yml`

`kubectl expose pod open-func-x242s --type="NodePort" --port 8090`




