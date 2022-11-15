package pluginManager

import (
	"context"
	crd "github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1"
	"github.com/spidernet-io/spiderdoctor/pkg/types"
	"go.uber.org/zap"
	"time"
)

func (s *pluginControllerReconciler) GetSpiderAgentNodeNotInSucceedRecod(ctx context.Context, succeedNodeList []string) (failNodelist []string, err error) {
	allNodeList, e := GetDaemonsetPodNodeNameList(ctx, s.client, types.ControllerConfig.SpiderDoctorAgentDaemonsetName, types.ControllerConfig.PodNamespace)
	if e != nil {
		return nil, e
	}

	failNodelist = []string{}
OUTER:
	for _, v := range allNodeList {
		for _, m := range succeedNodeList {
			if m == v {
				continue OUTER
			}
		}
		failNodelist = append(failNodelist, v)
	}
	return failNodelist, nil
}

func (s *pluginControllerReconciler) UpdateStatus(logger *zap.Logger, ctx context.Context, oldStatus *crd.TaskStatus, schedulePlan *crd.SchedulePlan) (*crd.TaskStatus, error) {
	newStatus := oldStatus.DeepCopy()
	recordLength := len(newStatus.History)

	// init new instance first
	if newStatus.ExpectedRound == nil || recordLength == 0 {
		n := schedulePlan.RoundNumber
		newStatus.ExpectedRound = &n
		m := int64(0)
		newStatus.DoneRound = &m
		newRecod := NewStatusHistoryRecord(1, schedulePlan)
		newStatus.History = append(newStatus.History, *newRecod)
		logger.Debug("initialize the status of new instance")
		// updating status firstly , it will trigger to handle it next round
		return newStatus, nil
	}

	if *newStatus.DoneRound == *newStatus.ExpectedRound {
		// done task
		return nil, nil
	}

	latestRecord := &(newStatus.History[recordLength-1])
	nowTime := time.Now()
	logger.Sugar().Debugf("current time:%v , latest history record: %+v", nowTime, latestRecord)
	logger.Sugar().Debugf("all history record: %+v", newStatus.History)

	switch {
	case nowTime.Before(latestRecord.StartTimeStamp.Time):
		// this round not start, do nothing
		logger.Sugar().Debugf("wait for starting next round ")

	case nowTime.Before(latestRecord.DeadLineTimeStamp.Time) && nowTime.After(latestRecord.StartTimeStamp.Time):
		// still in this round , do nothing

	case nowTime.After(latestRecord.DeadLineTimeStamp.Time):
		if int(*newStatus.DoneRound) == (recordLength - 1) {

			// update result in latestRecord
			if failedNodeList, e := s.GetSpiderAgentNodeNotInSucceedRecod(ctx, latestRecord.SucceedAgentNodeList); e != nil {
				return nil, e
			} else {
				if len(failedNodeList) > 0 {
					latestRecord.FailedAgentNodeList = failedNodeList
					n := crd.StatusHistoryRecordStatusSucceed
					latestRecord.Status = n
					newStatus.LastRoundStatus = &n
					logger.Sugar().Errorf("round %v failed , failedNode=%v", latestRecord.RoundNumber, failedNodeList)
				} else {
					n := crd.StatusHistoryRecordStatusSucceed
					latestRecord.Status = n
					newStatus.LastRoundStatus = &n
					logger.Sugar().Infof("round %v succeeded ", latestRecord.RoundNumber)
				}
			}

			// add next round record
			n := *(newStatus.DoneRound) + 1
			newStatus.DoneRound = &n
			newRecod := NewStatusHistoryRecord(int(*(newStatus.DoneRound)), schedulePlan)
			newStatus.History = append(newStatus.History, *newRecod)

			// TODO: add to workqueue to collect all report of last round

		} else {
			// it should not get here , because once it succeeded to insert New StatusHistoryRecord, it should to to case 1
			logger.Sugar().Warnf("it has collect the report last round, just wait for the start of next round")
		}

	}

	if *newStatus.DoneRound == *newStatus.ExpectedRound {
		newStatus.Finish = true
	} else {
		newStatus.Finish = false
	}

	return newStatus, nil

}
