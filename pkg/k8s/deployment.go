package k8s

import (
	"context"
	"log"
	"strconv"

	"github.com/IoanStoianov/Open-func/pkg/types"

	"k8s.io/client-go/kubernetes"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateDeployment creates a deployment based on a funcSpecs
func CreateDeployment(clientset *kubernetes.Clientset, funcSpecs types.FuncSpecs) (string, error) {

	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: funcSpecs.FuncName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(funcSpecs.Instances),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": funcSpecs.FuncName,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": funcSpecs.FuncName,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:            funcSpecs.FuncName,
							Image:           funcSpecs.ImageName,
							ImagePullPolicy: "Never",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: funcSpecs.FuncPort,
								},
							},
							EnvFrom: []apiv1.EnvFromSource{},
							Env: []apiv1.EnvVar{
								{
									Name:  "OPEN_FUNC_PORT",
									Value: strconv.Itoa(int(funcSpecs.FuncPort)),
								},
							},
						},
					},
				},
			},
		},
	}

	// Create Deployment
	log.Println("Creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}

	log.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
	return funcSpecs.FuncName, nil
}

// DeleteDeployment does what the name says
func DeleteDeployment(clientset *kubernetes.Clientset, name string) error {
	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	return deploymentsClient.Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func int32Ptr(i int32) *int32 { return &i }
