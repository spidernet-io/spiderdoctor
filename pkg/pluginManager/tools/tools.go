package tools

import (
	"fmt"
	crd "github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1"
)

func ValidataCrdSchedule(plan *crd.SchedulePlan) error {

	if plan == nil {
		return fmt.Errorf("Schedule is empty ")
	}

	if plan.StartAfterMinute < 0 {
		return fmt.Errorf("Schedule.StartAfterMinute %v must not be smaller than 0 ", plan.StartAfterMinute)
	}

	if plan.TimeoutMinute < 1 {
		return fmt.Errorf("Schedule.TimeoutMinute %v must not be smaller than 1 ", plan.TimeoutMinute)
	}
	if plan.IntervalMinute < 1 {
		return fmt.Errorf("Schedule.IntervalMinute %v must not be smaller than 1 ", plan.IntervalMinute)
	}
	if plan.TimeoutMinute > plan.IntervalMinute {
		return fmt.Errorf("Schedule.TimeoutMinute %v must not be bigger than Schedule.IntervalMinute %v ", plan.TimeoutMinute, plan.IntervalMinute)
	}

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
