package main

import (
	"fmt"

	gr "github.com/awesome-fc/golang-runtime"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Creates a job to do the training with given training image
func startTrainingJob(ctx *gr.FCContext, evt map[string]string) ([]byte, error) {
	trainJobNamePrefix := evt["trainJobNamePrefix"]
	demoVersion := evt["demoVersion"]
	trainJobName := fmt.Sprintf("%s-%s", trainJobNamePrefix, demoVersion)

	trainImage := evt["trainImage"]
	ossBucket := evt["ossBucket"]

	fcLogger := gr.GetLogger().WithField("requestId", ctx.RequestID)
	backoffLimit := int32(4)
	jobsClient := k8sClientset.BatchV1().Jobs(apiv1.NamespaceDefault)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: trainJobName,
		},

		Spec: batchv1.JobSpec{
			BackoffLimit:            &backoffLimit,
			Template:                apiv1.PodTemplateSpec{
				Spec:                apiv1.PodSpec{
					Containers: []apiv1.Container{{
						Name: trainJobName,
						Image: trainImage,
						Command: []string{"python",  "train.py"},
						Env: getEnv(ossBucket, ctx.Region),
					}},
					RestartPolicy: apiv1.RestartPolicyNever,
				},
			},
		},
	}

	// Create job
	fcLogger.Infof("Creating job...")
	jobResult, err := jobsClient.Create(job)
	if err != nil {
		panic(err)
	}
	fcLogger.Infof("Created job %q.\n", jobResult.GetObjectMeta().GetName())
	return []byte(fmt.Sprintf(`{"result": "%s", "trainStatus": "%s"}`, jobResult.Name, jobResult.Status.String())), nil
}

func getTrainingStatus(ctx *gr.FCContext, evt map[string]string) ([]byte, error) {
	trainJobNamePrefix := evt["trainJobNamePrefix"]
	demoVersion := evt["demoVersion"]
	trainJobName := fmt.Sprintf("%s-%s", trainJobNamePrefix, demoVersion)

	jobsClient := k8sClientset.BatchV1().Jobs(apiv1.NamespaceDefault)

	job, err := jobsClient.Get(trainJobName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	status := "running"
	if job.Status.Succeeded > 0 {
		status = "succeeded"
	}
	if job.Status.Failed > 0{
		status = "failed"
	}

	fcLogger := gr.GetLogger().WithField("requestId", ctx.RequestID)
	fcLogger.Infof("Job status %+v", job.Status)

	return []byte(fmt.Sprintf(`{"trainName": "%s", "status": "%s"}`, trainJobName, status)), nil
}
