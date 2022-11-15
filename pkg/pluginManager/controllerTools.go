package pluginManager

import (
	"context"
	crd "github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1"
	"github.com/spidernet-io/spiderdoctor/pkg/types"
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

func (s *pluginControllerReconciler) GetSpiderAgentNodeNotInSucceedRecord(ctx context.Context, succeedNodeList []string) (failNodelist []string, err error) {
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

func (s *pluginControllerReconciler) UpdateStatus(logger *zap.Logger, ctx context.Context, oldStatus *crd.TaskStatus, schedulePlan *crd.SchedulePlan) (result *reconcile.Result, taskStatus *crd.TaskStatus, e error) {
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
		return nil, newStatus, nil
	}

	if *newStatus.DoneRound == *newStatus.ExpectedRound {
		// done task
		return nil, nil, nil
	}

	latestRecord := &(newStatus.History[recordLength-1])
	nowTime := time.Now()
	logger.Sugar().Debugf("current time:%v , latest history record: %+v", nowTime, latestRecord)
	logger.Sugar().Debugf("all history record: %+v", newStatus.History)

	switch {
	case nowTime.Before(latestRecord.StartTimeStamp.Time):
		// this round not start, do nothing
		logger.Sugar().Debugf("wait for starting next round ")

		// trigger when task end
		result = &reconcile.Result{
			RequeueAfter: latestRecord.DeadLineTimeStamp.Time.Sub(nowTime),
		}

	case nowTime.Before(latestRecord.DeadLineTimeStamp.Time) && nowTime.After(latestRecord.StartTimeStamp.Time):
		// still in this round , do nothing

		// trigger when task end
		result = &reconcile.Result{
			RequeueAfter: latestRecord.DeadLineTimeStamp.Time.Sub(nowTime),
		}

	case nowTime.After(latestRecord.DeadLineTimeStamp.Time):
		if int(*newStatus.DoneRound) == (recordLength - 1) {

			// update result in latestRecord
			reportNode := []string{}
			reportNode = append(reportNode, latestRecord.SucceedAgentNodeList...)
			reportNode = append(reportNode, latestRecord.FailedAgentNodeList...)
			if unknowReportNodeList, e := s.GetSpiderAgentNodeNotInSucceedRecord(ctx, reportNode); e != nil {
				return nil, nil, e
			} else {
				if len(unknowReportNodeList) > 0 || len(latestRecord.FailedAgentNodeList) > 0 {
					latestRecord.UnReportAgentNodeList = unknowReportNodeList
					n := crd.StatusHistoryRecordStatusFail
					latestRecord.Status = n
					newStatus.LastRoundStatus = &n
					logger.Sugar().Errorf("round %v failed , failedNode=%v, unknowReportNode=%v", latestRecord.RoundNumber, latestRecord.FailedAgentNodeList, unknowReportNodeList)
					if len(latestRecord.FailedAgentNodeList) > 0 {
						latestRecord.FailureReason = "some agents failed"
					} else if len(unknowReportNodeList) > 0 {
						latestRecord.FailureReason = "some agents did not report"
					}
				} else {
					n := crd.StatusHistoryRecordStatusSucceed
					latestRecord.Status = n
					newStatus.LastRoundStatus = &n
					logger.Sugar().Infof("round %v succeeded ", latestRecord.RoundNumber)
				}
			}

			// trigger when task start
			result = &reconcile.Result{
				RequeueAfter: latestRecord.StartTimeStamp.Time.Sub(nowTime),
			}

			// add next round record
			n := *(newStatus.DoneRound) + 1
			newStatus.DoneRound = &n
			newRecod := NewStatusHistoryRecord(int(n+1), schedulePlan)
			newStatus.History = append(newStatus.History, *newRecod)

			// TODO: add to workqueue to collect all report of last round, for node latestRecord.FailedAgentNodeList and latestRecord.SucceedAgentNodeList

		} else {
			// it should not get here , because once it succeeded to insert New StatusHistoryRecord, it should to to case 1
			logger.Sugar().Warnf("it has collect the report last round, just wait for the start of next round")
		}

	}

	if *newStatus.DoneRound == *newStatus.ExpectedRound {
		newStatus.Finish = true
		result = nil
	} else {
		newStatus.Finish = false
	}

	return result, newStatus, nil

}
