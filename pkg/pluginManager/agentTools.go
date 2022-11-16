package pluginManager

import (
	"context"
	"fmt"
	k8sObjManager "github.com/spidernet-io/spiderdoctor/pkg/k8ObjManager"
	crd "github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1"
	plugintypes "github.com/spidernet-io/spiderdoctor/pkg/pluginManager/types"
	"github.com/spidernet-io/spiderdoctor/pkg/taskStatusManager"
	"github.com/spidernet-io/spiderdoctor/pkg/types"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

func (s *pluginAgentReconciler) CallPluginImplementRoundTask(logger *zap.Logger, obj runtime.Object, schedulePlan *crd.SchedulePlan, taskName string, roundNumber int) {
	taskRoundName := fmt.Sprintf("%s.round%d", taskName, roundNumber)

	roundDuration := time.Duration(schedulePlan.TimeoutMinute) * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), roundDuration)
	defer cancel()
	taskSucceed := make(chan bool)
	logger.Sugar().Infof("call plugin to implement with timeout %v minute", schedulePlan.TimeoutMinute)

	go func() {
		msg := plugintypes.PluginReport{
			TaskName:      taskName,
			RoundNumber:   roundNumber,
			AgentNodeName: s.localNodeName,
			StartTimeStam: time.Now(),
		}
		failureReason, report, e := s.plugin.AgentEexecuteTask(logger, ctx, obj)
		if e != nil {
			logger.Sugar().Errorf("plugin failed to implement the round task, error=%v", e)
			taskSucceed <- false
		} else {
			if len(failureReason) == 0 {
				taskSucceed <- true
				msg.RoundResult = plugintypes.RoundResultSucceed
				msg.FailedReason = failureReason
			} else {
				taskSucceed <- false
				msg.RoundResult = plugintypes.RoundResultFail
				msg.FailedReason = ""
			}
		}
		msg.EndTimeStamp = time.Now()
		if report != nil {
			msg.Detail = report
		}

		// output to staout
		fmt.Printf("%+v\n", msg)

		// TODO: write report to disk for controler to collect
	}()

	select {
	case <-ctx.Done():
		logger.Sugar().Errorf("timeout for getting result from plugin, the round task failed")
		s.taskRoundData.SetTask(taskRoundName, taskStatusManager.RoundStatusFail)
	case r := <-taskSucceed:
		logger.Sugar().Infof("succed to call plugin to implement round task, succeed=%v", r)
		if r {
			s.taskRoundData.SetTask(taskRoundName, taskStatusManager.RoundStatusSucceeded)
		} else {
			s.taskRoundData.SetTask(taskRoundName, taskStatusManager.RoundStatusFail)
		}
	}

	// delete data
	go func() {
		time.Sleep(roundDuration)
		s.taskRoundData.DeleteTask(taskRoundName)
	}()
}

func (s *pluginAgentReconciler) HandleAgentTaskRound(logger *zap.Logger, ctx context.Context, oldStatus *crd.TaskStatus, schedulePlan *crd.SchedulePlan, obj runtime.Object, taskName string) (result *reconcile.Result, taskStatus *crd.TaskStatus, e error) {
	newStatus := oldStatus.DeepCopy()
	recordLength := len(newStatus.History)
	nowTime := time.Now()

	// check node selector whether need to implement it
	if schedulePlan.SourceAgentNodeSelector != nil {
		if ok, e := k8sObjManager.GetK8sObjManager().MatchNodeSelected(ctx, types.AgentConfig.LocalNodeName, schedulePlan.SourceAgentNodeSelector); e != nil {
			msg := fmt.Sprintf("failed to MatchNodeSelected, error=%v", e)
			logger.Error(msg)
			return nil, nil, fmt.Errorf(msg)
		} else {
			if !ok {
				logger.Sugar().Infof("local node is not selected by the task, node selector=%v , ignore", schedulePlan.SourceAgentNodeSelector)
				return nil, nil, nil
			}
		}
	}

	if newStatus.ExpectedRound == nil || recordLength == 0 || *newStatus.DoneRound == *newStatus.ExpectedRound {
		// not start or all finish
		return nil, nil, nil
	}

	latestRecord := &(newStatus.History[recordLength-1])
	logger.Sugar().Debugf("current time:%v , latest history record: %+v", nowTime, latestRecord)
	// logger.Sugar().Debugf("all history record: %+v", newStatus.History)

	if latestRecord.Status != crd.StatusHistoryRecordStatusOngoing {
		logger.Sugar().Debugf("ignore task %v , no opportunity to implement ", taskName)
		return nil, nil, nil
	}

	// check whether we have reported the round result
	if len(latestRecord.SucceedAgentNodeList) != 0 || len(latestRecord.FailedAgentNodeList) != 0 {
		v := []string{}
		v = append(v, latestRecord.SucceedAgentNodeList...)
		v = append(v, latestRecord.FailedAgentNodeList...)
		logger.Sugar().Debugf("check whether localNode %v has report ", s.localNodeName)

		if ok, e := CheckItemInList(s.localNodeName, v); e != nil {
			logger.Sugar().Errorf("failed to check local node in task record")
			// no need to requeue
			return nil, nil, nil
		} else {
			if ok {
				logger.Sugar().Debugf("ignore task %v , it has reported for the round result", taskName)
				return nil, nil, nil
			}
		}
	}

	taskRoundName := fmt.Sprintf("%s.round%d", taskName, latestRecord.RoundNumber)
	nextInterval := time.Duration(types.AgentConfig.Configmap.TaskPollIntervalInSecond) * time.Second

	if status, existed := s.taskRoundData.CheckTask(taskRoundName); !existed {
		// mark to started it
		s.taskRoundData.SetTask(taskRoundName, taskStatusManager.RoundStatusOngoing)

		// we still have not reported the result for an ongoing round. do it
		go s.CallPluginImplementRoundTask(logger.Named(taskRoundName), obj, schedulePlan, taskName, latestRecord.RoundNumber)
		logger.Sugar().Infof("task %v , trigger to implement task round, and try to poll report after %v second", taskRoundName, types.AgentConfig.Configmap.TaskPollIntervalInSecond)

		// trigger to poll result after interval
		result = &reconcile.Result{
			RequeueAfter: nextInterval,
		}

	} else {
		if status == taskStatusManager.RoundStatusOngoing {
			// the task is on going
			logger.Sugar().Infof("task %v is going, try to poll report after %v second", taskRoundName, types.AgentConfig.Configmap.TaskPollIntervalInSecond)

			// trigger to poll result after interval
			result = &reconcile.Result{
				RequeueAfter: nextInterval,
			}

		} else {
			// the task finish, report
			if status == taskStatusManager.RoundStatusSucceeded {
				logger.Sugar().Infof("task %v , report to succeed", taskRoundName)
				latestRecord.SucceedAgentNodeList = append(latestRecord.SucceedAgentNodeList, s.localNodeName)
			} else {
				logger.Sugar().Infof("task %v , report to fail", taskRoundName)
				latestRecord.FailedAgentNodeList = append(latestRecord.FailedAgentNodeList, s.localNodeName)
			}
			// requeue immediately to make sure the update succeed , not conflicted
			result = &reconcile.Result{
				Requeue: true,
			}

		}
	}

	return result, newStatus, nil
}
