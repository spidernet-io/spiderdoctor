package pluginManager

import (
	"context"
	"fmt"
	k8sObjManager "github.com/spidernet-io/spiderdoctor/pkg/k8ObjManager"
	crd "github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1"
	"github.com/spidernet-io/spiderdoctor/pkg/types"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

func (s *pluginControllerReconciler) GetSpiderAgentNodeNotInRecord(ctx context.Context, succeedNodeList []string, agentNodeSelector *metav1.LabelSelector) (failNodelist []string, err error) {
	var allNodeList []string
	var e error
	if agentNodeSelector == nil {
		allNodeList, e = k8sObjManager.GetK8sObjManager().ListDaemonsetPodNodes(ctx, types.ControllerConfig.SpiderDoctorAgentDaemonsetName, types.ControllerConfig.PodNamespace)
		if e != nil {
			return nil, e
		}
		s.logger.Sugar().Debugf("all agent nodes: %+v", allNodeList)
		if len(allNodeList) == 0 {
			return nil, fmt.Errorf("failed to find agent node ")
		}
	} else {
		allNodeList, e = k8sObjManager.GetK8sObjManager().ListSelectedNodes(ctx, agentNodeSelector)
		if e != nil {
			return nil, e
		}
		s.logger.Sugar().Debugf("selected agent nodes: %+v", allNodeList)
		if len(allNodeList) == 0 {
			return nil, fmt.Errorf("failed to find agent node ")
		}
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

func (s *pluginControllerReconciler) UpdateRoundFinalStatus(logger *zap.Logger, ctx context.Context, newStatus *crd.TaskStatus, agentNodeSelector *metav1.LabelSelector, deadline bool) (roundDone bool, err error) {

	recordLength := len(newStatus.History)
	latestRecord := &(newStatus.History[recordLength-1])
	roundNumber := latestRecord.RoundNumber

	if latestRecord.Status == crd.StatusHistoryRecordStatusFail || latestRecord.Status == crd.StatusHistoryRecordStatusSucceed || latestRecord.Status == crd.StatusHistoryRecordStatusNotstarted {
		return true, nil
	}

	// when not reach deadline, ignore when nothing report
	if !deadline && len(latestRecord.SucceedAgentNodeList) == 0 && len(latestRecord.FailedAgentNodeList) == 0 {
		logger.Sugar().Debugf("round %v not report anthing", roundNumber)
		return false, nil
	}

	// update result in latestRecord
	reportNode := []string{}
	reportNode = append(reportNode, latestRecord.SucceedAgentNodeList...)
	reportNode = append(reportNode, latestRecord.FailedAgentNodeList...)
	if unknowReportNodeList, e := s.GetSpiderAgentNodeNotInRecord(ctx, reportNode, agentNodeSelector); e != nil {
		logger.Sugar().Errorf("round %v failed to GetSpiderAgentNodeNotInSucceedRecord, error=%v", roundNumber, e)
		return false, e
	} else {
		if len(unknowReportNodeList) > 0 && !deadline {
			// when not reach the deadline, ignore
			logger.Sugar().Debugf("round %v , partial agents did not reported, wait for daedline", roundNumber)
			return false, nil
		}

		// it's ok to collect round status
		if len(unknowReportNodeList) > 0 || len(latestRecord.FailedAgentNodeList) > 0 {
			latestRecord.UnReportAgentNodeList = unknowReportNodeList
			n := crd.StatusHistoryRecordStatusFail
			latestRecord.Status = n
			newStatus.LastRoundStatus = &n
			logger.Sugar().Errorf("round %v failed , failedNode=%v, unknowReportNode=%v", roundNumber, latestRecord.FailedAgentNodeList, unknowReportNodeList)

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
		return true, nil
	}

}

func (s *pluginControllerReconciler) UpdateStatus(logger *zap.Logger, ctx context.Context, oldStatus *crd.TaskStatus, schedulePlan *crd.SchedulePlan, taskName string) (result *reconcile.Result, taskStatus *crd.TaskStatus, e error) {
	newStatus := oldStatus.DeepCopy()
	recordLength := len(newStatus.History)
	nextInterval := time.Duration(types.ControllerConfig.Configmap.TaskPollIntervalInSecond) * time.Second
	nowTime := time.Now()

	// init new instance first
	if newStatus.ExpectedRound == nil || recordLength == 0 {
		n := schedulePlan.RoundNumber
		newStatus.ExpectedRound = &n
		m := int64(0)
		newStatus.DoneRound = &m

		startTime := time.Now().Add(time.Duration(schedulePlan.StartAfterMinute) * time.Minute)
		newRecod := NewStatusHistoryRecord(startTime, 1, schedulePlan)
		newStatus.History = append(newStatus.History, *newRecod)
		logger.Sugar().Debugf("initialize the first round of task : %v ", taskName, *newRecod)
		// trigger
		result = &reconcile.Result{
			Requeue: true,
		}
		// updating status firstly , it will trigger to handle it next round
		return result, newStatus, nil
	}

	if *newStatus.DoneRound == *newStatus.ExpectedRound {
		// done task
		return nil, nil, nil
	}

	latestRecord := &(newStatus.History[recordLength-1])
	roundNumber := latestRecord.RoundNumber
	logger.Sugar().Debugf("current time:%v , latest history record: %+v", nowTime, latestRecord)
	logger.Sugar().Debugf("all history record: %+v", newStatus.History)

	switch {
	case nowTime.After(latestRecord.StartTimeStamp.Time) && nowTime.Before(latestRecord.DeadLineTimeStamp.Time):

		if latestRecord.Status == crd.StatusHistoryRecordStatusNotstarted {
			latestRecord.Status = crd.StatusHistoryRecordStatusOngoing
			// requeue immediately to make sure the update succeed , not conflicted
			result = &reconcile.Result{
				Requeue: true,
			}

		} else if latestRecord.Status == crd.StatusHistoryRecordStatusOngoing {
			logger.Debug("try to poll the status of task " + taskName)
			if roundDone, e := s.UpdateRoundFinalStatus(logger, ctx, newStatus, schedulePlan.SourceAgentNodeSelector, false); e != nil {
				return nil, nil, e
			} else {
				if roundDone {
					logger.Sugar().Infof("round %v get reports from all agents ", roundNumber)

					// add new round record
					if *(newStatus.DoneRound) < *(newStatus.ExpectedRound) {
						n := *(newStatus.DoneRound) + 1
						newStatus.DoneRound = &n
						if n < *(newStatus.ExpectedRound) {
							startTime := latestRecord.StartTimeStamp.Time.Add(time.Duration(schedulePlan.IntervalMinute) * time.Minute)
							newRecod := NewStatusHistoryRecord(startTime, int(n+1), schedulePlan)
							newStatus.History = append(newStatus.History, *newRecod)
							logger.Sugar().Infof("insert new record for next round : %+v", *newRecod)
						} else {
							newStatus.Finish = true
						}
					}

					// TODO: add to workqueue to collect all report of last round, for node latestRecord.FailedAgentNodeList and latestRecord.SucceedAgentNodeList

					// requeue immediately to make sure the update succeed , not conflicted
					result = &reconcile.Result{
						Requeue: true,
					}

				} else {
					// trigger after interval
					result = &reconcile.Result{
						RequeueAfter: nextInterval,
					}
				}
			}
		} else {
			logger.Debug("ignore poll the finished round of task " + taskName)

			// trigger when deadline
			result = &reconcile.Result{
				RequeueAfter: latestRecord.DeadLineTimeStamp.Time.Sub(time.Now()),
			}
		}

	case nowTime.Before(latestRecord.StartTimeStamp.Time):
		fallthrough
	case nowTime.After(latestRecord.DeadLineTimeStamp.Time):
		if *newStatus.DoneRound == *newStatus.ExpectedRound {
			logger.Sugar().Debugf("task %s finish, ignore ", taskName)
			newStatus.Finish = true
			result = nil

		} else {

			// when task not finsih , once we update the status succeed , we will not get here , it should go to case nowTime.Before(latestRecord.StartTimeStamp.Time)
			if latestRecord.Status == crd.StatusHistoryRecordStatusOngoing {
				// here, we should update last round status

				if _, e := s.UpdateRoundFinalStatus(logger, ctx, newStatus, schedulePlan.SourceAgentNodeSelector, true); e != nil {
					return nil, nil, e
				} else {
					logger.Sugar().Infof("round %v get reports from all agents ", roundNumber)

					// add new round record
					if *(newStatus.DoneRound) < *(newStatus.ExpectedRound) {
						n := *(newStatus.DoneRound) + 1
						newStatus.DoneRound = &n

						if n < *(newStatus.ExpectedRound) {
							startTime := latestRecord.StartTimeStamp.Time.Add(time.Duration(schedulePlan.IntervalMinute) * time.Minute)
							newRecod := NewStatusHistoryRecord(startTime, int(n+1), schedulePlan)
							newStatus.History = append(newStatus.History, *newRecod)
							logger.Sugar().Infof("insert new record for next round : %+v", *newRecod)
						} else {
							newStatus.Finish = true
						}
					}

					// TODO: add to workqueue to collect all report of last round, for node latestRecord.FailedAgentNodeList and latestRecord.SucceedAgentNodeList

					// requeue immediately to make sure the update succeed , not conflicted
					result = &reconcile.Result{
						Requeue: true,
					}
				}
			} else {
				// round finish
				// trigger when next round start
				newRoundNumber := len(newStatus.History)
				currentRecord := &(newStatus.History[newRoundNumber-1])
				logger.Sugar().Infof("task %v wait for next round %v at %v", taskName, newRoundNumber, currentRecord.StartTimeStamp)
				result = &reconcile.Result{
					RequeueAfter: currentRecord.StartTimeStamp.Time.Sub(time.Now()),
				}
			}
		}
	}

	return result, newStatus, nil

}
