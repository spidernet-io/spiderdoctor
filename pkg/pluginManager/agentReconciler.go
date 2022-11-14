package pluginManager

import (
	"context"
	"fmt"
	plugintypes "github.com/spidernet-io/spiderdoctor/pkg/pluginManager/types"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

type pluginAgentReconciler struct {
	client client.Client
	p      plugintypes.ChainingPlugin
	logger *zap.Logger
}

func (s *pluginAgentReconciler) Reconcile(ctx context.Context, r reconcile.Request) (reconcile.Result, error) {
	return s.p.AgentReconcile(s.logger, s.client, ctx, r)
}

var _ reconcile.Reconciler = &pluginAgentReconciler{}

func (s *pluginManager) runAgentReconcile() {
	logger := s.logger

	scheme := runtime.NewScheme()
	if e:=clientgoscheme.AddToScheme(scheme);e != nil {
		logger.Sugar().Fatalf("failed to add k8s scheme, reason=%v", e)
	}
	for name, plugin := range s.chainingPlugins {
		if e := plugin.AddToScheme(scheme); e != nil {
			logger.Sugar().Fatalf("failed to add scheme for plugin, reason=%v", name, e)
		}
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

	for name, plugin := range s.chainingPlugins {
		logger.Sugar().Infof("run controller for plugin %v", name)
		k := &pluginAgentReconciler{
			logger: logger.Named(name + "Reconciler"),
			p:      plugin,
		}
		b, e := ctrl.NewControllerManagedBy(mgr).For(plugin.GetApiType()).Complete(k)
		if e != nil {
			s.logger.Sugar().Fatalf("failed to builder reconcile for plugin %v, error=%v", name, e)
		}
	}

	go func() {
		msg := fmt.Sprintf("reconcile of plugin down")
		if e := mgr.Start(context.Background()); e != nil {
			msg += fmt.Sprintf(", error=%v", e)
		}
		s.logger.Error(msg)
		time.Sleep(5 * time.Second)
	}()
}
