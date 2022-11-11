// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package v1

type SchedulePlan struct {
	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=0
	RoundNumber int64 `json:"roundNumber"`

	// +kubebuilder:validation:Optional
	Interval int64 `json:"interval,omitempty"`

	// +kubebuilder:default=60
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Optional
	TimeoutMinute int64 `json:"timeoutMinute,omitempty"`
}
