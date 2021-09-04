module github.com/zouchunxu/deployment

go 1.16

require (
	github.com/gin-gonic/gin v1.7.4
	github.com/zouchunxu/gof v0.0.0-20210902165205-8883e9090ef9
	google.golang.org/genproto v0.0.0-20210602131652-f16073e35f0c
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.26.0
	k8s.io/api v0.22.1
	k8s.io/apimachinery v0.22.1
	k8s.io/client-go v0.22.1
)

replace github.com/zouchunxu/gof => /Users/zouchunxu/web/docker/wwwroot/go/gof
