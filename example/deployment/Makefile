proto:
	protoc --proto_path=. --proto_path="./third_party/" --proto_path=./api --gin_out="paths=source_relative:." --go_out=paths=source_relative:. --go-grpc_out=require_unimplemented_servers=false:. --go-grpc_opt=paths=source_relative  ./api/*.proto

build:
	sh build.sh

run:
	kubectl apply -f deployment.yaml

rbac:
	kubectl create clusterrolebinding default-view --clusterrole=view --serviceaccount=default:default
