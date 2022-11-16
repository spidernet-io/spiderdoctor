package k8sObjManager

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (nm *k8sObjManager) GetPodList(ctx context.Context, opts ...client.ListOption) ([]corev1.Pod, error) {
	var podlist corev1.PodList
	if e := nm.client.List(ctx, &podlist, opts...); e != nil {
		return nil, e
	}
	return podlist.Items, nil
}
