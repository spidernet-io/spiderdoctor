package k8sObjManager

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (nm *k8sObjManager) GetService(ctx context.Context, name, namespace string) (*corev1.Service, error) {

	d := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	key := client.ObjectKeyFromObject(d)
	if e := nm.client.Get(ctx, key, d); e != nil {
		return nil, e
	}
	return d, nil
}
