/*
Copyright 2017 The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Note: the example only works with the code within the same release/branch.
package main

import (
	"encoding/json"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cs"
	gr "github.com/awesome-fc/golang-runtime"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type csDescribeClusterUserKubeconfigResp struct {
	Config string `json:"config"`
}

var k8sClientset *kubernetes.Clientset

// initializes k8s cluster client by getting config from the ContainerService
func initialize(ctx *gr.FCContext, evt map[string]string) error {
	k8sClusterID := evt["k8sClusterID"]
	fcLogger := gr.GetLogger().WithField("requestId", ctx.RequestID)
	fcLogger.Infoln("init initialize golang!")
	csClient, err := cs.NewClientWithStsToken(ctx.Region, ctx.Credentials.AccessKeyID, ctx.Credentials.AccessKeySecret, ctx.Credentials.SecurityToken)
	if err != nil {
		return err
	}

	dcuKubeConfReq := cs.CreateDescribeClusterUserKubeconfigRequest()
	dcuKubeConfReq.ClusterId = k8sClusterID
	resp, err := csClient.DescribeClusterUserKubeconfig(dcuKubeConfReq)
	if err != nil {
		return err
	}

	fcLogger.Infof("DescribeClusterUserKubeconfig resp: %s, %d", resp.RequestId, resp.GetHttpStatus())
	respBody := &csDescribeClusterUserKubeconfigResp{}
	if err := json.Unmarshal(resp.GetHttpContentBytes(), respBody); err != nil {
		panic(err)
	}

	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(respBody.Config))
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	k8sClientset = clientset
	return nil
}

func handler(ctx *gr.FCContext, event []byte) ([]byte, error) {
	evt := make(map[string]string)
	if err := json.Unmarshal(event, &evt); err != nil {
		panic(err)
	}

	if err := initialize(ctx, evt); err != nil {
		return nil, err
	}

	action := evt["action"]
	switch action {
	case "train-model":
		return startTrainingJob(ctx, evt)
	case "deploy-model":
		return deployModel(ctx, evt)
	case "get-train-status":
		return getTrainingStatus(ctx, evt)
	case "get-deployment-status":
		return getDeploymentStatus(ctx, evt)
	case "expose-service":
		return exposeService(ctx, evt)
	case "get-service-status":
		return getServiceStatus(ctx, evt)
	case "run-test-cases":
		return runTestCases(ctx, evt)
	default:
		return nil, fmt.Errorf("action %s is not supported", action)
	}

	return event, nil
}

func main() {
	gr.Start(handler, nil)
}

func int32Ptr(i int32) *int32 { return &i }

func getEnv(ossBucket, region string) []apiv1.EnvVar {
	return []apiv1.EnvVar{
		{
			Name: "ACCESS_KEY_ID",
			ValueFrom: &apiv1.EnvVarSource{
				SecretKeyRef: &apiv1.SecretKeySelector{
					Key: "ACCESS_KEY_ID",
					LocalObjectReference: apiv1.LocalObjectReference{
						Name: "fnf-secret",
					},
				},
			},
		},
		{
			Name: "ACCESS_KEY_ID_SECRET",
			ValueFrom: &apiv1.EnvVarSource{
				SecretKeyRef: &apiv1.SecretKeySelector{
					Key: "ACCESS_KEY_ID_SECRET",
					LocalObjectReference: apiv1.LocalObjectReference{
						Name: "fnf-secret",
					},
				},
			},
		},
		{
			Name: "OSS_BUCKET",
			Value: ossBucket,
		},
		{
			Name: "REGION",
			Value: region,
		},
	}
}
