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

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=2
	// +kubebuilder:validation:Minimum=1
	DurationInSecond int `json:"durationInSecond,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=5
	// +kubebuilder:validation:Maximum=100
	// +kubebuilder:validation:Minimum=1
	QPS int `json:"qps,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=5
	// +kubebuilder:validation:Minimum=1
	PerRequestTimeoutInSecond int `json:"perRequestTimeoutInSecond,omitempty"`
}

type TargetAgentSepc struct {
	// +kubebuilder:default=true
	// +kubebuilder:validation:Optional
	TestIPv4 *bool `json:"testIPv4,omitempty"`

	// +kubebuilder:default=false
	// +kubebuilder:validation:Optional
	TestIPv6 *bool `json:"testIPv6,omitempty"`

	// +kubebuilder:default=true
	TestEndpoint bool `json:"testEndpoint,omitempty"`

	// +kubebuilder:default=true
	TestMultusInterface bool `json:"testMultusInterface,omitempty"`

	// +kubebuilder:default=true
	TestClusterIp bool `json:"testClusterIp,omitempty"`

	// +kubebuilder:default=true
	TestNodePort bool `json:"testNodePort,omitempty"`

	// +kubebuilder:default=false
	TestIngress bool `json:"testIngress,omitempty"`
}

type NethttpTarget struct {

	// +kubebuilder:validation:Optional
	TargetUrl *string `json:"targetUrl,omitempty"`

	// +kubebuilder:validation:Optional
	TargetAgent *TargetAgentSepc `json:"targetAgent,omitempty"`
}

// scope(Namespaced or Cluster)
// +kubebuilder:resource:categories={spiderdoctor},path="nethttps",singular="nethttp",scope="Cluster"
// +kubebuilder:printcolumn:JSONPath=".status.finish",description="finish",name="finish",type=boolean
// +kubebuilder:printcolumn:JSONPath=".status.expectedRound",description="expectedRound",name="expectedRound",type=integer
// +kubebuilder:printcolumn:JSONPath=".status.doneRound",description="doneRound",name="doneRound",type=integer
// +kubebuilder:printcolumn:JSONPath=".status.lastRoundStatus",description="lastRoundStatus",name="lastRoundStatus",type=string
// +kubebuilder:printcolumn:JSONPath=".spec.schedule.intervalMinute",description="roundIntervalMinute",name="intervalMinute",type=integer
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
