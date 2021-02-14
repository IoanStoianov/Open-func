package k8s

import (
	"context"
	"fmt"
	"strconv"

	"github.com/IoanStoianov/Open-func/pkg/types"

	"k8s.io/client-go/kubernetes"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//
func CreateDeployment(clientset *kubernetes.Clientset, funcTrigger types.FuncTrigger) {

	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: funcTrigger.FuncName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": funcTrigger.FuncName,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": funcTrigger.FuncName,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:            funcTrigger.FuncName,
							Image:           funcTrigger.ImageName,
							ImagePullPolicy: "Never",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: funcTrigger.FuncPort,
								},
							},
							EnvFrom: []apiv1.EnvFromSource{},
							Env: []apiv1.EnvVar{
								{
									Name:  "OPEN_FUNC_PORT",
									Value: strconv.Itoa(int(funcTrigger.FuncPort)),
								},
							},
						},
					},
				},
			},
		},
	}

	// Create Deployment
	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
}

func int32Ptr(i int32) *int32 { return &i }
