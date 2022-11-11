// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

// !!!!!! crd marker:
// https://github.com/kubernetes-sigs/controller-tools/blob/master/pkg/crd/markers/crd.go
// https://book.kubebuilder.io/reference/markers/crd.html
// https://github.com/kubernetes-sigs/controller-tools/blob/master/pkg/crd/markers/validation.go
// https://book.kubebuilder.io/reference/markers/crd-validation.html

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NetdoctorSpec struct {
	// +kubebuilder:validation:Optional
	Schedule *SchedulePlan `json:"schedule,omitempty"`

	EnabledIPv4 bool `json:"enabledIPv4"`

	EnabledIPv6 bool `json:"enabledIPv6"`
}

type NetdoctorStatus struct {
	// +kubebuilder:validation:Minimum=0
	ExpectedRound int64 `json:"expectedRound"`

	// +kubebuilder:validation:Minimum=0
	DoneRound *int64 `json:"doneRound"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type:=string
	// +kubebuilder:validation:Format:=date-time
	LastRoundTimeStamp *metav1.Time `json:"lastRoundTimeStamp,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type:=string
	// +kubebuilder:validation:Format:=date-time
	NextRoundTimeStamp *metav1.Time `json:"nextRoundTimeStamp,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=succeed;fail;unknown
	LastRoundStatus *string `json:"lastRoundStatus,omitempty"`
}

// scope(Namespaced or Cluster)
// +kubebuilder:resource:categories={spiderdoctor},path="netdoctors",singular="netdoctor",scope="Cluster",shortName={nd}
// +kubebuilder:printcolumn:JSONPath=".status.ExpectedRound",description="ExpectedRound",name="ExpectedRound",type=integer
// +kubebuilder:printcolumn:JSONPath=".status.DoneRound",description="DoneRound",name="DoneRound",type=integer
// +kubebuilder:printcolumn:JSONPath=".status.LastRoundStatus",description="LastRoundStatus",name="LastRoundStatus",type=integer
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +genclient
// +genclient:nonNamespaced

type Netdoctor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   NetdoctorSpec   `json:"spec,omitempty"`
	Status NetdoctorStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

type NetdoctorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Netdoctor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Netdoctor{}, &NetdoctorList{})
}
