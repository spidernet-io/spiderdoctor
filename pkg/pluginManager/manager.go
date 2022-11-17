// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package pluginManager

import (
	"fmt"
	k8sObjManager "github.com/spidernet-io/spiderdoctor/pkg/k8ObjManager"
	crd "github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1"
	"github.com/spidernet-io/spiderdoctor/pkg/lock"
	"github.com/spidernet-io/spiderdoctor/pkg/pluginManager/netdns"
	"github.com/spidernet-io/spiderdoctor/pkg/pluginManager/nethttp"
	plugintypes "github.com/spidernet-io/spiderdoctor/pkg/pluginManager/types"
	"github.com/spidernet-io/spiderdoctor/pkg/taskStatusManager"
	"github.com/spidernet-io/spiderdoctor/pkg/types"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"time"
)

var pluginLock = &lock.Mutex{}

type pluginManager struct {
	chainingPlugins map[string]plugintypes.ChainingPlugin
	logger          *zap.Logger
}
type PluginManager interface {
	RunAgentController()
	RunControllerController(healthPort int, webhookPort int, webhookTlsDir string)
}

var globalPluginManager *pluginManager

// --------------------------------------
func (s *pluginManager) RunAgentController() {
	logger := s.logger
	logger.Sugar().Infof("setup agent reconcile")

	scheme := runtime.NewScheme()
	if e := clientgoscheme.AddToScheme(scheme); e != nil {
		logger.Sugar().Fatalf("failed to add k8s scheme, reason=%v", e)
	}
	if e := crd.AddToScheme(scheme); e != nil {
		logger.Sugar().Fatalf("failed to add scheme for plugins, reason=%v", e)
	}

	n := ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     "0",
		HealthProbeBindAddress: "0",
		LeaderElection:         false,
	}
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), n)
	if err != nil {
		logger.Sugar().Fatalf("failed to NewManager, reason=%v", err)
	}

	if len(types.AgentConfig.LocalNodeName) == 0 {
		logger.Sugar().Fatalf("local node name is empty")
	}

	if e := k8sObjManager.Initk8sObjManager(mgr.GetClient()); e != nil {
		logger.Sugar().Fatalf("failed to Initk8sObjManager, error=%v", e)
	}

	for name, plugin := range s.chainingPlugins {
		logger.Sugar().Infof("run controller for plugin %v", name)
		k := &pluginAgentReconciler{
			logger:        logger.Named(name + "Reconciler"),
			plugin:        plugin,
			client:        mgr.GetClient(),
			crdKind:       name,
			taskRoundData: taskStatusManager.NewTaskStatus(),
			localNodeName: types.AgentConfig.LocalNodeName,
		}
		if e := k.SetupWithManager(mgr); e != nil {
			s.logger.Sugar().Fatalf("failed to builder reconcile for plugin %v, error=%v", name, e)
		}
	}

	// before mgr.Start, it should not use mgr.GetClient() to get api obj, because "the controller cache is not started, can not read objects"
	go func() {
		msg := "reconcile of plugin down"
		if e := mgr.Start(ctrl.SetupSignalHandler()); e != nil {
			msg += fmt.Sprintf(", error=%v", e)
		}
		s.logger.Error(msg)
		time.Sleep(5 * time.Second)
	}()

}

// --------------------------------------

func (s *pluginManager) RunControllerController(healthPort int, webhookPort int, webhookTlsDir string) {

	logger := s.logger
	scheme := runtime.NewScheme()
	if e := clientgoscheme.AddToScheme(scheme); e != nil {
		logger.Sugar().Fatalf("failed to add k8s scheme, reason=%v", e)
	}
	if e := crd.AddToScheme(scheme); e != nil {
		logger.Sugar().Fatalf("failed to add scheme for plugins, reason=%v", e)
	}

	n := ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: "0",
		// health
		HealthProbeBindAddress: "0",
		// webhook
		Port:    webhookPort,
		CertDir: webhookTlsDir,
		// lease
		LeaderElection:          true,
		LeaderElectionNamespace: types.ControllerConfig.PodNamespace,
		LeaderElectionID:        types.ControllerConfig.PodName,
	}
	if healthPort != 0 {
		n.HealthProbeBindAddress = fmt.Sprintf(":%d", healthPort)
		n.ReadinessEndpointName = "/healthy/readiness"
		n.LivenessEndpointName = "/healthy/liveness"
	}
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), n)
	if err != nil {
		logger.Sugar().Fatalf("failed to NewManager, reason=%v", err)
	}
	if healthPort != 0 {
		// could implement your checker , type Checker func(req *http.Request) error
		if err := mgr.AddHealthzCheck("/healthy/liveness", healthz.Ping); err != nil {
			logger.Sugar().Fatalf("failed to AddHealthzCheck, reason=%v", err)
		}
		if err := mgr.AddReadyzCheck("/healthy/readiness", healthz.Ping); err != nil {
			logger.Sugar().Fatalf("failed to AddReadyzCheck, reason=%v", err)
		}
		// add other route
		// mgr.GetWebhookServer().Register("/route", XXXX)
	}

	if e := k8sObjManager.Initk8sObjManager(mgr.GetClient()); e != nil {
		logger.Sugar().Fatalf("failed to Initk8sObjManager, error=%v", e)

	}

	for name, plugin := range s.chainingPlugins {
		// setup reconcile
		logger.Sugar().Infof("run controller for plugin %v", name)
		k := &pluginControllerReconciler{
			logger:  logger.Named(name + "Reconciler"),
			plugin:  plugin,
			client:  mgr.GetClient(),
			crdKind: name,
		}
		if e := k.SetupWithManager(mgr); e != nil {
			s.logger.Sugar().Fatalf("failed to builder reconcile for plugin %v, error=%v", name, e)
		}
		// setup webhook
		t := &pluginWebhookhander{
			logger:  logger.Named(name + "Webhook"),
			plugin:  plugin,
			client:  mgr.GetClient(),
			crdKind: name,
		}
		if e := t.SetupWebhook(mgr); e != nil {
			s.logger.Sugar().Fatalf("failed to builder webhook for plugin %v, error=%v", name, e)
		}
	}

	go func() {
		msg := "reconcile of plugin down"
		if e := mgr.Start(ctrl.SetupSignalHandler()); e != nil {
			msg += fmt.Sprintf(", error=%v", e)
		}
		s.logger.Error(msg)
		time.Sleep(5 * time.Second)
	}()

}

// -------------------------

func InitPluginManager(logger *zap.Logger) PluginManager {
	pluginLock.Lock()
	defer pluginLock.Unlock()

	globalPluginManager.logger = logger

	return globalPluginManager
}

const (
	// ------ add crd ------
	KindNameNethttp = "Nethttp"
	KindNameNetdns  = "Netdns"
)

func init() {
	globalPluginManager = &pluginManager{
		chainingPlugins: map[string]plugintypes.ChainingPlugin{},
	}

	// ------ add crd ------
	globalPluginManager.chainingPlugins[KindNameNethttp] = &nethttp.PluginNetHttp{}
	globalPluginManager.chainingPlugins[KindNameNetdns] = &netdns.PluginNetDns{}

}
