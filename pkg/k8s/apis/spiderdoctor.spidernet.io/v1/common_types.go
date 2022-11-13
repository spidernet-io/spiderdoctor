// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package v1

type SchedulePlan struct {
	// +kubebuilder:default=0
	// +kubebuilder:validation:Minimum=0
	StartAfterMinute int64 `json:"startAfterMinute"`

	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=0
	RoundNumber int64 `json:"roundNumber"`

	// +kubebuilder:default=360
	// +kubebuilder:validation:Optional
	IntervalMinute int64 `json:"intervalMinute,omitempty"`

	// +kubebuilder:default=60
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Optional
	TimeoutMinute int64 `json:"timeoutMinute,omitempty"`
}

type NetTarget struct {
	// +kubebuilder:validation:Optional
	Service *TargetService `json:"service,omitempty"`

	// +kubebuilder:validation:Optional
	HostAddress *string `json:"hostAddress,omitempty"`
}

type TargetService struct {
	ServiceName string `json:"serviceName,omitempty"`

	// +kubebuilder:validation:Optional
	TestEndpoint *bool `json:"testEndpoint,omitempty"`

	// +kubebuilder:validation:Optional
	TestNodePort *bool `json:"testNodePort,omitempty"`

	// +kubebuilder:validation:Optional
	TestIngress *bool `json:"testIngress,omitempty"`
}