package pluginManager

import (
	"context"
	"fmt"
	crd "github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

func NewStatusHistoryRecord(startTime time.Time, RoundNumber int, schedulePlan *crd.SchedulePlan) *crd.StatusHistoryRecord {
	newRecod := crd.StatusHistoryRecord{
		Status:                crd.StatusHistoryRecordStatusNotstarted,
		FailureReason:         "",
		RoundNumber:           RoundNumber,
		SucceedAgentNodeList:  []string{},
		FailedAgentNodeList:   []string{},
		UnReportAgentNodeList: []string{},
	}
	newRecod.StartTimeStamp = metav1.NewTime(startTime)

	adder := time.Duration(schedulePlan.TimeoutMinute) * time.Minute
	endTime := startTime.Add(adder)
	newRecod.DeadLineTimeStamp = metav1.NewTime(endTime)

	return &newRecod
}

func GetPod(ctx context.Context, c client.Client, name, namespace string, opts ...client.ListOption) (*corev1.Pod, error) {
	d := &corev1.Pod{
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

func GetPodNodeName(ctx context.Context, c client.Client, name, namespace string) (nodeName string, err error) {
	if len(name) == 0 {
		return "", fmt.Errorf("empty pod name")
	}
	if len(namespace) == 0 {
		return "", fmt.Errorf("empty pod namespace")
	}

	if d, e := GetPod(ctx, c, name, namespace, nil); e != nil {
		return "", e
	} else {
		if d == nil {
			return "", fmt.Errorf("pod is empty")
		}
		return d.Spec.NodeName, nil
	}
}

func CheckItemInList(item string, checklist []string) (bool, error) {
	if len(item) == 0 {
		return false, fmt.Errorf("empty item")
	}
	if len(checklist) == 0 {
		return false, nil
	}
	for _, v := range checklist {
		if v == item {
			return true, nil
		}
	}
	return false, nil
}
