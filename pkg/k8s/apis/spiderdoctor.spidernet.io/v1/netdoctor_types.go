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

	// +kubebuilder:validation:Optional
	Target *NetTarget `json:"target,omitempty"`

	// +kubebuilder:default=true
	// +kubebuilder:validation:Optional
	TestIPv4 *bool `json:"testIPv4,omitempty"`

	// +kubebuilder:default=false
	// +kubebuilder:validation:Optional
	TestIPv6 *bool `json:"testIPv6,omitempty"`

	// +kubebuilder:validation:Optional
	EachTimeInSecond *uint64 `json:"eachTimeInSecond,omitempty"`

	// +kubebuilder:validation:Optional
	EachQPS *uint64 `json:"eachQPS,omitempty"`

	// +kubebuilder:validation:Optional
	FailureCondition *NetFailureCondition `json:"failureCondition,omitempty"`
}

type NetFailureCondition struct {
	// +kubebuilder:default=1
	// +kubebuilder:validation:Optional
	MinAccessFailure *uint64 `json:"minAccessFailure,omitempty"`

	// +kubebuilder:default=5000
	// +kubebuilder:validation:Optional
	MinAccessDelayMs *uint64 `json:"minAccessDelayMs,omitempty"`
}

type NetdoctorStatus struct {
	// +kubebuilder:validation:Minimum=0
	ExpectedRound int64 `json:"expectedRound"`

	// +kubebuilder:validation:Minimum=0
	DoneRound *int64 `json:"doneRound"`

	Finish bool `json:"finish"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type:=string
	// +kubebuilder:validation:Format:=date-time
	LastRoundFinishTimeStamp *metav1.Time `json:"lastRoundFinishTimeStamp,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=succeed;fail;unknown
	LastRoundStatus *string `json:"lastRoundStatus,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type:=string
	// +kubebuilder:validation:Format:=date-time
	NextRoundStartTimeStamp *metav1.Time `json:"nextRoundStartTimeStamp,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type:=string
	// +kubebuilder:validation:Format:=date-time
	NextRoundDeadLineTimeStamp *metav1.Time `json:"nextRoundDeadLineTimeStamp,omitempty"`
}

// scope(Namespaced or Cluster)
// +kubebuilder:resource:categories={spiderdoctor},path="netdoctors",singular="netdoctor",scope="Cluster",shortName={nd}
// +kubebuilder:printcolumn:JSONPath=".status.Finish",description="Finish",name="Finish",type=bool
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
