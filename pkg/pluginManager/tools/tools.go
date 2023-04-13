// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package tools

import (
	"fmt"
	crd "github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1beta1"
)

func ValidataCrdSchedule(plan *crd.SchedulePlan) error {

	if plan == nil {
		return fmt.Errorf("Schedule is empty ")
	}

	if plan.Simple.StartAfterMinute < 0 {
		return fmt.Errorf("Schedule.StartAfterMinute %v must not be smaller than 0 ", plan.Simple.StartAfterMinute)
	}

	if plan.TimeoutMinute < 1 {
		return fmt.Errorf("Schedule.TimeoutMinute %v must not be smaller than 1 ", plan.TimeoutMinute)
	}
	if plan.Simple.IntervalMinute < 1 {
		return fmt.Errorf("Schedule.IntervalMinute %v must not be smaller than 1 ", plan.Simple.IntervalMinute)
	}
	if plan.TimeoutMinute > plan.Simple.IntervalMinute {
		return fmt.Errorf("Schedule.TimeoutMinute %v must not be bigger than Schedule.IntervalMinute %v ", plan.TimeoutMinute, plan.Simple.IntervalMinute)
	}

	return nil
}

func GetDefaultSchedule() (plan *crd.SchedulePlan) {
	return &crd.SchedulePlan{
		TimeoutMinute: 60,
		Simple: &crd.Simple{
			StartAfterMinute: 0,
			RoundNumber:      1,
			IntervalMinute:   60,
		}}
}

func GetDefaultNetSuccessCondition() (plan *crd.NetSuccessCondition) {
	n := float64(1)
	return &crd.NetSuccessCondition{
		SuccessRate: &n,
	}
}
