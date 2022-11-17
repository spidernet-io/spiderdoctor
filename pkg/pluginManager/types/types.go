// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"context"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

type ChainingPlugin interface {
	GetApiType() client.Object

	AgentEexecuteTask(logger *zap.Logger, ctx context.Context, obj runtime.Object) (failureReason string, report PluginRoundDetail, err error)

	ControllerReconcile(*zap.Logger, client.Client, context.Context, reconcile.Request) (reconcile.Result, error)
	AgentReconcile(*zap.Logger, client.Client, context.Context, reconcile.Request) (reconcile.Result, error)

	WebhookMutating(logger *zap.Logger, ctx context.Context, obj runtime.Object) error
	WebhookValidateCreate(logger *zap.Logger, ctx context.Context, obj runtime.Object) error
	WebhookValidateUpdate(logger *zap.Logger, ctx context.Context, oldObj, newObj runtime.Object) error
	WebhookValidateDelete(logger *zap.Logger, ctx context.Context, obj runtime.Object) error
}

type RoundResultStatus string

const (
	RoundResultSucceed = RoundResultStatus("succeed")
	RoundResultFail    = RoundResultStatus("fail")
)

type PluginReport struct {
	TaskName      string
	RoundNumber   int
	RoundResult   RoundResultStatus
	AgentNodeName string
	FailedReason  string
	StartTimeStam time.Time
	EndTimeStamp  time.Time
	Detail        PluginRoundDetail
}

type PluginRoundDetail map[string]interface{}

const (
	ApiMsgGetFailure      = "failed to get instance"
	ApiMsgUnknowCRD       = "unsupported crd type"
	ApiMsgUnsupportModify = "unsupported modify spec"
)
