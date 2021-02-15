package k8s

import (
	"context"
	"log"

	"github.com/IoanStoianov/Open-func/pkg/types"

	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CreateJob starts a job pod for a cold start trigger
func CreateJob(clientset *kubernetes.Clientset, trigger types.ColdTriggerEvent) (string, error) {
	batchClient := clientset.BatchV1().Jobs(apiv1.NamespaceDefault)

	job := &batchv1.Job{
		ObjectMeta: v1.ObjectMeta{
			Name: trigger.FuncName,
		},
		Spec: batchv1.JobSpec{
			ActiveDeadlineSeconds: int64Ptr(10),
			Template: apiv1.PodTemplateSpec{
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:            trigger.FuncName,
							Image:           trigger.ImageName,
							ImagePullPolicy: "Never",
							Env: []apiv1.EnvVar{
								{
									Name:  "REDIS_URL",
									Value: "redis", // TODO: Get from environment
								},
							},
						},
					},
					RestartPolicy: apiv1.RestartPolicyNever,
				},
			},
		},
	}

	log.Println("Starting job...")
	result, err := batchClient.Create(context.TODO(), job, v1.CreateOptions{})
	if err != nil {
		return "", err
	}

	jobName := result.GetObjectMeta().GetName()

	log.Printf("Created deployment %q.\n", jobName)
	return jobName, nil
}

// DeleteJob does what the name says
func DeleteJob(clientset *kubernetes.Clientset, name string) error {
	batchClient := clientset.BatchV1().Jobs(apiv1.NamespaceDefault)

	return batchClient.Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func int64Ptr(i int64) *int64 { return &i }
