// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NethttpSpec struct {
	// +kubebuilder:validation:Optional
	Schedule *SchedulePlan `json:"schedule,omitempty"`

	// +kubebuilder:validation:Optional
	Target *NethttpTarget `json:"target,omitempty"`

	// +kubebuilder:validation:Optional
	Request *NethttpRequest `json:"request,omitempty"`

	// +kubebuilder:validation:Optional
	SuccessCondition *NetSuccessCondition `json:"success,omitempty"`
}

type NethttpRequest struct {
	// +kubebuilder:default=true
	// +kubebuilder:validation:Optional
	TestIPv4 *bool `json:"testIPv4,omitempty"`

	// +kubebuilder:default=false
	// +kubebuilder:validation:Optional
	TestIPv6 *bool `json:"testIPv6,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=1
	DurationInSecond *uint64 `json:"durationInSecond,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=1
	QPS *uint64 `json:"qps,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=1
	PerRequestTimeoutInSecond *uint64 `json:"perRequestTimeoutInSecond,omitempty"`
}

type NethttpTarget struct {

	// +kubebuilder:default=true
	TestEndpoint bool `json:"testEndpoint,omitempty"`

	// +kubebuilder:default=true
	TestNodePort bool `json:"testNodePort,omitempty"`

	// +kubebuilder:default=false
	TestIngress bool `json:"testIngress,omitempty"`
}

// scope(Namespaced or Cluster)
// +kubebuilder:resource:categories={spiderdoctor},path="nethttps",singular="nethttp",scope="Cluster"
// +kubebuilder:printcolumn:JSONPath=".status.Finish",description="Finish",name="Finish",type=boolean
// +kubebuilder:printcolumn:JSONPath=".status.ExpectedRound",description="ExpectedRound",name="ExpectedRound",type=integer
// +kubebuilder:printcolumn:JSONPath=".status.DoneRound",description="DoneRound",name="DoneRound",type=integer
// +kubebuilder:printcolumn:JSONPath=".status.LastRoundStatus",description="LastRoundStatus",name="LastRoundStatus",type=integer
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +genclient
// +genclient:nonNamespaced

type Nethttp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   NethttpSpec `json:"spec,omitempty"`
	Status TaskStatus  `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

type NethttpList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Nethttp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Nethttp{}, &NethttpList{})
}
