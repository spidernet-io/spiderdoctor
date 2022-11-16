package k8sObjManager

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (nm *k8sObjManager) GetDaemonset(ctx context.Context, name, namespace string) (*appsv1.DaemonSet, error) {
	d := &appsv1.DaemonSet{
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

func (nm *k8sObjManager) ListDaemonsetPodNodes(ctx context.Context, daemonsetName, daemonsetNameSpace string) ([]string, error) {

	dae, e := nm.GetDaemonset(ctx, daemonsetName, daemonsetNameSpace)
	if e != nil {
		return nil, fmt.Errorf("failed to get daemonset, error=%v", e)
	}

	podLable := dae.Spec.Template.Labels
	opts := []client.ListOption{
		client.MatchingLabelsSelector{
			Selector: labels.SelectorFromSet(podLable),
		},
	}
	podlist, e := nm.GetPodList(ctx, opts...)
	if e != nil {
		return nil, fmt.Errorf("failed to get pod list, error=%v", e)
	}

	nodelist := []string{}
	for _, v := range podlist {
		nodelist = append(nodelist, v.Spec.NodeName)
	}
	return nodelist, nil
}
