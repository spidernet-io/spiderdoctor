package k8sObjManager

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apitypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (nm *k8sObjManager) GetNode(ctx context.Context, nodeName string) (*corev1.Node, error) {
	var node corev1.Node
	if err := nm.client.Get(ctx, apitypes.NamespacedName{Name: nodeName}, &node); err != nil {
		return nil, err
	}

	return &node, nil
}

func (nm *k8sObjManager) ListNodes(ctx context.Context, opts ...client.ListOption) (*corev1.NodeList, error) {
	var nodeList corev1.NodeList
	if err := nm.client.List(ctx, &nodeList, opts...); err != nil {
		return nil, err
	}

	return &nodeList, nil
}

func (nm *k8sObjManager) ListSelectedNodes(ctx context.Context, labelSelector *metav1.LabelSelector) ([]string, error) {
	selector, err := metav1.LabelSelectorAsSelector(labelSelector)
	if err != nil {
		return nil, err
	}

	nodeList, err := nm.ListNodes(
		ctx,
		client.MatchingLabelsSelector{Selector: selector},
	)
	if err != nil {
		return nil, err
	}

	if len(nodeList.Items) == 0 {
		return nil, nil
	}

	v := []string{}
	for _, t := range nodeList.Items {
		v = append(v, t.Name)
	}
	return v, nil
}

func (nm *k8sObjManager) MatchNodeSelected(ctx context.Context, nodeName string, labelSelector *metav1.LabelSelector) (bool, error) {
	selector, err := metav1.LabelSelectorAsSelector(labelSelector)
	if err != nil {
		return false, err
	}

	nodeList, err := nm.ListNodes(
		ctx,
		client.MatchingLabelsSelector{Selector: selector},
		client.MatchingFields{metav1.ObjectNameField: nodeName},
	)
	if err != nil {
		return false, err
	}

	if len(nodeList.Items) == 0 {
		return false, nil
	}

	return true, nil
}
