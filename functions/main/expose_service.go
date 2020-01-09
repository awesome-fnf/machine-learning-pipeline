package main

import (
	"fmt"

	gr "github.com/awesome-fc/golang-runtime"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Creates a service for the given model serving deployment with LoadBalancer type such that the predictions/inferences
// can be accessed via an internet IP
func exposeService(ctx *gr.FCContext, evt map[string]string) ([]byte, error) {
	deploymentNamePrefix := evt["deploymentNamePrefix"]
	serviceNamePrefix := evt["serviceNamePrefix"]
	demoVersion := evt["demoVersion"]
	serviceName := fmt.Sprintf("%s-%s", serviceNamePrefix, demoVersion)
	deploymentName := fmt.Sprintf("%s-%s", deploymentNamePrefix, demoVersion)

	fcLogger := gr.GetLogger().WithField("requestId", ctx.RequestID)
	servicesClient := k8sClientset.CoreV1().Services(apiv1.NamespaceDefault)

	selector := map[string]string{}
	selector["app"] = deploymentName
	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceName,
		},
		Spec: apiv1.ServiceSpec{
			Ports:                    []apiv1.ServicePort{
				{
					Port: 8501,
					TargetPort: intstr.FromInt(8501),
					Protocol: apiv1.ProtocolTCP,
				},
			},
			Selector:                 selector,
			Type:                     "LoadBalancer",
		},
	}

	// Create Deployment
	fcLogger.Infof("Creating service...")
	result, err := servicesClient.Create(service)
	if err != nil {
		panic(err)
	}
	fcLogger.Infof("Created service %q.\n", result.GetObjectMeta().GetName())

	return []byte(fmt.Sprintf(`{"name": "%s", "result": "%s"}`, result.Name, result.String())), nil
}

// Get the input service and returns the corresponding serving endpoint: http://{lb-ip}:8501
func getServiceStatus(ctx *gr.FCContext, evt map[string]string) ([]byte, error) {
	serviceNamePrefix := evt["serviceNamePrefix"]
	demoVersion := evt["demoVersion"]
	servicetName := fmt.Sprintf("%s-%s", serviceNamePrefix, demoVersion)

	servicesClient := k8sClientset.CoreV1().Services(apiv1.NamespaceDefault)

	service, err := servicesClient.Get(servicetName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if len(service.Status.LoadBalancer.Ingress) > 0 {
		lbIP := service.Status.LoadBalancer.Ingress[0].IP
		return []byte(fmt.Sprintf(`{"servingEndpoint": "http://%s:8501", "serviceStatus": "%s"}`, lbIP, "success")), nil
	}

	return []byte(fmt.Sprintf(`{"servingEndpoint": "", "serviceStatus": "%s"}`, "not-ready")), nil
}