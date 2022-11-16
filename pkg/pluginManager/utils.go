package pluginManager

import (
	"context"
	"fmt"
	crd "github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

func NewStatusHistoryRecord(RoundNumber int, schedulePlan *crd.SchedulePlan) *crd.StatusHistoryRecord {
	newRecod := crd.StatusHistoryRecord{
		Status:                crd.StatusHistoryRecordStatusNotstarted,
		FailureReason:         "",
		RoundNumber:           RoundNumber,
		SucceedAgentNodeList:  []string{},
		FailedAgentNodeList:   []string{},
		UnReportAgentNodeList: []string{},
	}
	startTime := time.Now().Add(time.Duration(schedulePlan.StartAfterMinute) * time.Minute)
	newRecod.StartTimeStamp = metav1.NewTime(startTime)

	adder := time.Duration(schedulePlan.TimeoutMinute) * time.Minute
	endTime := startTime.Add(adder)
	newRecod.DeadLineTimeStamp = metav1.NewTime(endTime)

	return &newRecod
}

func GetPodList(ctx context.Context, c client.Client, opts ...client.ListOption) ([]corev1.Pod, error) {
	var podlist corev1.PodList
	if e := c.List(ctx, &podlist, opts...); e != nil {
		return nil, e
	}
	return podlist.Items, nil
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

func GetDaemonsetPodNodeNameList(ctx context.Context, c client.Client, daemonsetName, namespace string) ([]string, error) {

	dae, e := GetDaemonset(ctx, c, daemonsetName, namespace)
	if e != nil {
		return nil, fmt.Errorf("failed to get daemonset, error=%v", e)
	}

	podLable := dae.Spec.Template.Labels
	opts := []client.ListOption{
		client.MatchingLabelsSelector{
			Selector: labels.SelectorFromSet(podLable),
		},
	}
	podlist, e := GetPodList(ctx, c, opts...)
	if e != nil {
		return nil, fmt.Errorf("failed to get pod list, error=%v", e)
	}

	nodelist := []string{}
	for _, v := range podlist {
		nodelist = append(nodelist, v.Spec.NodeName)
	}
	return nodelist, nil
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
