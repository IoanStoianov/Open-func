package k8s

import (
	"context"
	"fmt"

	"github.com/IoanStoianov/Open-func/pkg/types"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

//
func CreateService(clientset *kubernetes.Clientset, fungTrigger types.FuncTrigger) int32 {
	servicesClient := clientset.CoreV1().Services(apiv1.NamespaceDefault)

	newService := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-port-%d-service", fungTrigger.FuncName, fungTrigger.FuncPort),
			Labels: map[string]string{
				"func": fungTrigger.FuncName,
			},
		},
		Spec: apiv1.ServiceSpec{

			Type: apiv1.ServiceTypeNodePort,
			Selector: map[string]string{
				"func": fungTrigger.FuncName,
			},
			Ports: []apiv1.ServicePort{
				{
					Protocol: apiv1.ProtocolTCP,
					Port:     fungTrigger.FuncPort,
				},
			},
		},
	}

	fmt.Println("Creating service...")
	result, err := servicesClient.Create(context.TODO(), newService, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created service %q.\n", result.GetObjectMeta().GetName())

	return result.Spec.Ports[0].NodePort
}
