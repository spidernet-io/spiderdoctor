package pluginManager

import (
	"context"
	"github.com/spidernet-io/spiderdoctor/pkg/lock"
	"github.com/spidernet-io/spiderdoctor/pkg/types"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var pluginLock = &lock.Mutex{}

type PluginManager interface {
	RunAgentController()
	RunControllerController(webhookPort int, webhookTlsDir string)
}

type pluginManager struct {
	chainingPlugins map[string]ChainingPlugin
	logger          *zap.Logger
}

var globalPluginManager *pluginManager

type ChainingPlugin interface {
	GetApiType() client.Object
	CheckObjType(obj runtime.Object) bool

	ControllerReconcile(client.Client, context.Context, reconcile.Request) (reconcile.Result, error)
	AgentReconcile(client.Client, context.Context, reconcile.Request) (reconcile.Result, error)

	WebhookMutating(logger *zap.Logger, ctx context.Context, obj runtime.Object) error
	WebhookValidateCreate(logger *zap.Logger, ctx context.Context, obj runtime.Object) error
	WebhookValidateUpdate(logger *zap.Logger, ctx context.Context, oldObj, newObj runtime.Object) error
	WebhookValidateDelete(logger *zap.Logger, ctx context.Context, obj runtime.Object) error
}

func RegisterPlugin(pluginName string, p ChainingPlugin) {
	pluginLock.Lock()
	defer pluginLock.Unlock()

	if globalPluginManager == nil {
		globalPluginManager = &pluginManager{
			chainingPlugins: map[string]ChainingPlugin{},
		}
	}

	globalPluginManager.chainingPlugins[pluginName] = p
}

func NewPluginManager(logger *zap.Logger) PluginManager {
	pluginLock.Lock()
	defer pluginLock.Unlock()

	if globalPluginManager == nil {
		globalPluginManager = &pluginManager{
			logger:          logger,
			chainingPlugins: map[string]ChainingPlugin{},
		}
	} else {
		globalPluginManager.logger = logger
	}
	return globalPluginManager
}

// --------------------

type pluginControllerReconciler struct {
	client client.Client
	p      ChainingPlugin
}

func (s *pluginControllerReconciler) Reconcile(ctx context.Context, r reconcile.Request) (reconcile.Result, error) {
	return s.p.ControllerReconcile(s.client, ctx, r)
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
		go func(name string, t ChainingPlugin) {
			logger.Sugar().Infof("run controller for plugin %v", name)
			builder.For(t.GetApiType()).Owns(t.GetApiType()).Build(&pluginControllerReconciler{p: t})
		}(name, plugin)
	}
}

// --------------------

type pluginAgentReconciler struct {
	client client.Client
	p      ChainingPlugin
}

func (s *pluginAgentReconciler) Reconcile(ctx context.Context, r reconcile.Request) (reconcile.Result, error) {
	return s.p.AgentReconcile(s.client, ctx, r)
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
		go func(name string, t ChainingPlugin) {
			logger.Sugar().Infof("run controller for plugin %v", name)
			builder.For(t.GetApiType()).Owns(t.GetApiType()).Build(&pluginAgentReconciler{p: t})
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
