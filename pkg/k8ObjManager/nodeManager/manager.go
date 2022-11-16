// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package nodeManager

import (
	"context"
	"errors"
	"fmt"
	"github.com/spidernet-io/spiderdoctor/pkg/lock"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	apitypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type NodeManager interface {
	GetNode(ctx context.Context, nodeName string) (*corev1.Node, error)
	ListNodes(ctx context.Context, opts ...client.ListOption) (*corev1.NodeList, error)
	MatchNodeSelected(ctx context.Context, nodeName string, labelSelector *metav1.LabelSelector) (bool, error)
	ListSelectedNodes(ctx context.Context, labelSelector *metav1.LabelSelector) ([]string, error)
	ListDaemonsetPodNodes(ctx context.Context, daemonsetName, daemonsetNameSpace string) ([]string, error)
}

type nodeManager struct {
	client client.Client
}

var l lock.Mutex
var globalNodeManager *nodeManager

func InitNodeManager(client client.Client) (NodeManager, error) {
	if client == nil {
		return nil, errors.New("k8s client must be specified")
	}
	l.Lock()
	defer l.Unlock()

	if globalNodeManager == nil {
		globalNodeManager = &nodeManager{
			client: client,
		}
	}
	return globalNodeManager, nil
}

func (nm *nodeManager) GetNode(ctx context.Context, nodeName string) (*corev1.Node, error) {
	var node corev1.Node
	if err := nm.client.Get(ctx, apitypes.NamespacedName{Name: nodeName}, &node); err != nil {
		return nil, err
	}

	return &node, nil
}

func (nm *nodeManager) ListNodes(ctx context.Context, opts ...client.ListOption) (*corev1.NodeList, error) {
	var nodeList corev1.NodeList
	if err := nm.client.List(ctx, &nodeList, opts...); err != nil {
		return nil, err
	}

	return &nodeList, nil
}

func (nm *nodeManager) ListSelectedNodes(ctx context.Context, labelSelector *metav1.LabelSelector) ([]string, error) {
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

func (nm *nodeManager) MatchNodeSelected(ctx context.Context, nodeName string, labelSelector *metav1.LabelSelector) (bool, error) {
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

// ------------------

func GetPodList(ctx context.Context, c client.Client, opts ...client.ListOption) ([]corev1.Pod, error) {
	var podlist corev1.PodList
	if e := c.List(ctx, &podlist, opts...); e != nil {
		return nil, e
	}
	return podlist.Items, nil
}

func GetDaemonset(ctx context.Context, c client.Client, name, namespace string) (*appsv1.DaemonSet, error) {
	d := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	key := client.ObjectKeyFromObject(d)
	if e := c.Get(ctx, key, d); e != nil {
		return nil, e
	}
	return d, nil
}

func (nm *nodeManager) ListDaemonsetPodNodes(ctx context.Context, daemonsetName, daemonsetNameSpace string) ([]string, error) {

	dae, e := GetDaemonset(ctx, nm.client, daemonsetName, daemonsetNameSpace)
	if e != nil {
		return nil, fmt.Errorf("failed to get daemonset, error=%v", e)
	}

	podLable := dae.Spec.Template.Labels
	opts := []client.ListOption{
		client.MatchingLabelsSelector{
			Selector: labels.SelectorFromSet(podLable),
		},
	}
	podlist, e := GetPodList(ctx, nm.client, opts...)
	if e != nil {
		return nil, fmt.Errorf("failed to get pod list, error=%v", e)
	}

	nodelist := []string{}
	for _, v := range podlist {
		nodelist = append(nodelist, v.Spec.NodeName)
	}
	return nodelist, nil
}
