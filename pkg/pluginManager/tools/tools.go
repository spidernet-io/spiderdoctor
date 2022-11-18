package tools

import (
	crd "github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1"
)

func ValidataCrdSchedule(plan *crd.SchedulePlan) error {

	return nil
}

func GetDefaultSchedule() (plan *crd.SchedulePlan) {
	return &crd.SchedulePlan{
		StartAfterMinute: 0,
		RoundNumber:      1,
		IntervalMinute:   60,
		TimeoutMinute:    60,
	}
}

func GetDefaultNetSuccessCondition() (plan *crd.NetSuccessCondition) {
	n := float64(1)
	return &crd.NetSuccessCondition{
		SuccessRate: &n,
	}
}
