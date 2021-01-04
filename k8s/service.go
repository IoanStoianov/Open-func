package k8s

import (
	"context"
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// apiVersion: v1
// kind: Service
// metadata:
//   name: my-service
// spec:
//   selector:
//     app: MyApp
//   ports:
//     - protocol: TCP
//       port: 80
//       targetPort: 9376

//
func CreateService(clientset *kubernetes.Clientset) {
	servicesClient := clientset.CoreV1().Services(apiv1.NamespaceDefault)

	newService := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nodejs-port-service",
			Labels: map[string]string{
				"app": "node-docker",
			},
		},
		Spec: apiv1.ServiceSpec{

			Type: apiv1.ServiceTypeNodePort,
			Selector: map[string]string{
				"app": "node-docker",
			},
			Ports: []apiv1.ServicePort{
				{
					Protocol: apiv1.ProtocolTCP,
					Port:     8000,
					NodePort: 32041,
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
}
