package main

import (
	"fmt"

	gr "github.com/awesome-fc/golang-runtime"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Creates a deployment for model serving with the given serving image
func deployModel(ctx *gr.FCContext, evt map[string]string) ([]byte, error) {
	deploymentNamePrefix := evt["deploymentNamePrefix"]
	demoVersion := evt["demoVersion"]
	deploymentName := fmt.Sprintf("%s-%s", deploymentNamePrefix, demoVersion)

	servingImage := evt["servingImage"]
	ossBucket := evt["ossBucket"]

	fcLogger := gr.GetLogger().WithField("requestId", ctx.RequestID)
	deploymentsClient := k8sClientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": deploymentName,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": deploymentName,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  deploymentName,
							Image: servingImage,
							Env:   getEnv(ossBucket, ctx.Region),
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 8501,
								},
							},
						},
					},
				},
			},
		},
	}

	// Create Deployment
	fcLogger.Infof("Creating deployment...")
	result, err := deploymentsClient.Create(deployment)
	if err != nil {
		panic(err)
	}
	fcLogger.Infof("Created deployment %q.\n", result.GetObjectMeta().GetName())

	return []byte(fmt.Sprintf(`{"name": "%s", "result": "%s"}`, result.Name, result.String())), nil
}

func getDeploymentStatus(ctx *gr.FCContext, evt map[string]string) ([]byte, error) {
	deploymentNamePrefix := evt["deploymentNamePrefix"]
	demoVersion := evt["demoVersion"]
	deploymentName := fmt.Sprintf("%s-%s", deploymentNamePrefix, demoVersion)

	deploymentsClient := k8sClientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	deployment, err := deploymentsClient.Get(deploymentName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	status := "running"
	if deployment.Status.AvailableReplicas > 0 {
		status = "succeeded"
	}

	podsGetter := k8sClientset.CoreV1().Pods(apiv1.NamespaceDefault)

	pods, err := podsGetter.List(metav1.ListOptions{
		TypeMeta:            metav1.TypeMeta{},
		LabelSelector:       fmt.Sprintf("app=%s", deploymentName),
	})
	if err != nil {
		return nil, err
	}

	podName := ""

	if len(pods.Items) > 0 {
		podName = pods.Items[0].Name
	}

	return []byte(fmt.Sprintf(`{"deploymentName": "%s", "deploymentStatus": "%s", "podName": "%s"}`, deploymentName, status, podName)), nil
}
