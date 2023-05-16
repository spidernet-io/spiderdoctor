// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PluginReport
// +k8s:openapi-gen=true
type PluginReport struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec PluginReportSpec `json:"spec,omitempty"`
}

// PluginReportList
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PluginReportList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []PluginReport `json:"items"`
}

// PluginReportSpec defines the desired state of PluginReport
type PluginReportSpec struct {
	TaskName string `json:"taskName"`
	//TaskSpec       interface{} `json:"taskSpec"`
	RoundNumber    int         `json:"roundNumber"`
	RoundResult    string      `json:"roundResult"`
	NodeName       string      `json:"nodeName"`
	PodName        string      `json:"podName"`
	FailedReason   string      `json:"failedReason"`
	StartTimeStamp metav1.Time `json:"startTimeStamp"`
	EndTimeStamp   metav1.Time `json:"endTimeStamp"`
	RoundDuraiton  string      `json:"roundDuraiton"`
	ReportType     string      `json:"reportType"`
	//Detail         interface{} `json:"detail"`
}

func init() {
	SchemeBuilder.Register(&PluginReport{}, &PluginReportList{})
}
