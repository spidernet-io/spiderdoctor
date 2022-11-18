// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package k8sObjManager

import (
	"context"
	"fmt"
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
		return nil, fmt.Errorf("failed to get service %v/%v, reason=%v", namespace, name, e)
	}
	return d, nil
}
