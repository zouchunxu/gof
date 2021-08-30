build:
	GOOS=linux go build -o ./app .
	eval $(minikube docker-env)
	docker build -t in-cluster .

run:
	kubectl run --rm -i demo --image=in-cluster --image-pull-policy=IfNotPresent

rbac:
	kubectl create clusterrolebinding default-view --clusterrole=view --serviceaccount=default:default
