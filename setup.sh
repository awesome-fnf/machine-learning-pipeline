set -e

# Build and push the training image
cd training
docker build -t registry.$REGION.aliyuncs.com/$NAMESPACE/fnf-fc-training-tensorflow-demo:$VERSION .
docker push registry.$REGION.aliyuncs.com/$NAMESPACE/fnf-fc-training-tensorflow-demo:$VERSION
cd ..

# Build and push the serving image
cd serving
docker build -t registry.$REGION.aliyuncs.com/$NAMESPACE/fnf-fc-serving-tensorflow-demo:$VERSION .
docker push registry.$REGION.aliyuncs.com/$NAMESPACE/fnf-fc-serving-tensorflow-demo:$VERSION
cd ..

bash deploy.sh




