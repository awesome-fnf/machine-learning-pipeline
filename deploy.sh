set -e

# Build and deploy the FC golang function
echo "Building Golang function binary"
cd functions/main

# Build golang binary with the public Docker image that has necessary K8s dependencies
docker run -it -v `pwd`:/go/src/fnf-fc-golang-k8s-builder registry.cn-hangzhou.aliyuncs.com/function-flow-demo/fnf-fc-golang-k8s-builder:v1 bash -c "cd /go/src/fnf-fc-golang-k8s-builder; go build -o bin/bootstrap *.go"
chmod 777 bin/bootstrap
zip -j code.zip bin/bootstrap
cd ../..

echo "Deploy Golang function"
fun deploy