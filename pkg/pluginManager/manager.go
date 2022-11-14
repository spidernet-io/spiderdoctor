package pluginManager

import (
	"context"
	"github.com/spidernet-io/spiderdoctor/pkg/lock"
	plugintypes "github.com/spidernet-io/spiderdoctor/pkg/pluginManager/types"
	"github.com/spidernet-io/spiderdoctor/pkg/plugins/nethttp"
	"github.com/spidernet-io/spiderdoctor/pkg/types"
	"go.uber.org/zap"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var pluginLock = &lock.Mutex{}

type pluginManager struct {
	chainingPlugins map[string]plugintypes.ChainingPlugin
	logger          *zap.Logger
}
type PluginManager interface {
	RunAgentController()
	RunControllerController(webhookPort int, webhookTlsDir string)
}

var globalPluginManager *pluginManager

func InitPluginManager(logger *zap.Logger) PluginManager {
	pluginLock.Lock()
	defer pluginLock.Unlock()

	globalPluginManager.logger = logger

	return globalPluginManager
}

func init() {
	globalPluginManager = &pluginManager{
		chainingPlugins: map[string]plugintypes.ChainingPlugin{},
	}
	globalPluginManager.chainingPlugins["nethttp"] = &nethttp.PluginNetHttp{}

}

// --------------------

type pluginControllerReconciler struct {
	client client.Client
	p      plugintypes.ChainingPlugin
	logger *zap.Logger
}

func (s *pluginControllerReconciler) Reconcile(ctx context.Context, r reconcile.Request) (reconcile.Result, error) {
	return s.p.ControllerReconcile(s.logger, s.client, ctx, r)
}

var _ reconcile.Reconciler = &pluginControllerReconciler{}

func (s *pluginManager) runControllerController() {
	logger := s.logger

	n := ctrl.Options{
		MetricsBindAddress:      "0",
		HealthProbeBindAddress:  "0",
		LeaderElection:          true,
		LeaderElectionNamespace: types.ControllerConfig.PodNamespace,
		LeaderElectionID:        types.ControllerConfig.PodName,
	}
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), n)
	if err != nil {
		logger.Sugar().Fatalf("failed to NewManager, reason=%v", err)
	}
	builder := ctrl.NewControllerManagedBy(mgr)
	for name, plugin := range s.chainingPlugins {
		go func(name string, t plugintypes.ChainingPlugin) {
			logger.Sugar().Infof("run controller for plugin %v", name)
			builder.For(t.GetApiType()).Owns(t.GetApiType()).Build(&pluginControllerReconciler{logger: logger.Named(name + "Reconciler"), p: t})
		}(name, plugin)
	}
}

// --------------------

type pluginAgentReconciler struct {
	client client.Client
	p      plugintypes.ChainingPlugin
	logger *zap.Logger
}

func (s *pluginAgentReconciler) Reconcile(ctx context.Context, r reconcile.Request) (reconcile.Result, error) {
	return s.p.AgentReconcile(s.logger, s.client, ctx, r)
}

var _ reconcile.Reconciler = &pluginAgentReconciler{}

func (s *pluginManager) runAgentController() {
	logger := s.logger

	n := ctrl.Options{
		MetricsBindAddress:     "0",
		HealthProbeBindAddress: "0",
		LeaderElection:         false,
	}
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), n)
	if err != nil {
		logger.Sugar().Fatalf("failed to NewManager, reason=%v", err)
	}
	builder := ctrl.NewControllerManagedBy(mgr)
	for name, plugin := range s.chainingPlugins {
		go func(name string, t plugintypes.ChainingPlugin) {
			logger.Sugar().Infof("run controller for plugin %v", name)
			builder.For(t.GetApiType()).Owns(t.GetApiType()).Build(&pluginAgentReconciler{logger: logger.Named(name + "Reconciler"), p: t})
		}(name, plugin)
	}
}

func (s *pluginManager) RunAgentController() {
	s.logger.Sugar().Infof("setup agent controller")
	s.runAgentController()
}

func (s *pluginManager) RunControllerController(webhookPort int, webhookTlsDir string) {
	s.logger.Sugar().Infof("setup controller webhook")
	s.runWebhook(webhookPort, webhookTlsDir)
	s.logger.Sugar().Infof("setup controller controller")
	s.runControllerController()
}
