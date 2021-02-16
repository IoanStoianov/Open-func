package k8s

import (
	"context"
	"fmt"
	"log"

	"github.com/IoanStoianov/Open-func/pkg/types"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

// CreateService creates a clusterIP service for the deployment of a trigger
func CreateService(clientset *kubernetes.Clientset, funcSpecs types.FuncSpecs) (string, error) {
	servicesClient := clientset.CoreV1().Services(apiv1.NamespaceDefault)

	newService := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-service", funcSpecs.FuncName),
			Labels: map[string]string{
				"app": funcSpecs.FuncName,
			},
		},
		Spec: apiv1.ServiceSpec{
			Type: apiv1.ServiceTypeClusterIP,
			Selector: map[string]string{
				"app": funcSpecs.FuncName,
			},
			Ports: []apiv1.ServicePort{
				{
					Port:       80,
					TargetPort: intstr.FromInt(int(funcSpecs.FuncPort)),
				},
			},
		},
	}

	log.Println("Creating service...")
	result, err := servicesClient.Create(context.TODO(), newService, metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
		return "", err
	}

	serviceName := result.GetObjectMeta().GetName()
	log.Printf("Created service %q.\n", serviceName)

	return serviceName, nil
}

// DeleteService deletes a service in default namespace
func DeleteService(clientset *kubernetes.Clientset, name string) error {
	servicesClient := clientset.CoreV1().Services(apiv1.NamespaceDefault)

	return servicesClient.Delete(context.TODO(), name, metav1.DeleteOptions{})
}
