package pluginManager

import (
	"fmt"
	crd "github.com/spidernet-io/spiderdoctor/pkg/k8s/apis/spiderdoctor.spidernet.io/v1"
	"github.com/spidernet-io/spiderdoctor/pkg/lock"
	"github.com/spidernet-io/spiderdoctor/pkg/pluginManager/nethttp"
	plugintypes "github.com/spidernet-io/spiderdoctor/pkg/pluginManager/types"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"time"
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

	for name, plugin := range s.chainingPlugins {
		logger.Sugar().Infof("run controller for plugin %v", name)
		k := &pluginAgentReconciler{
			logger: logger.Named(name + "Reconciler"),
			p:      plugin,
			client: mgr.GetClient(),
		}
		if e := k.RunAgentReconcile(mgr); e != nil {
			s.logger.Sugar().Fatalf("failed to builder reconcile for plugin %v, error=%v", name, e)
		}
	}

	go func() {
		msg := fmt.Sprintf("reconcile of plugin down")
		if e := mgr.Start(ctrl.SetupSignalHandler()); e != nil {
			msg += fmt.Sprintf(", error=%v", e)
		}
		s.logger.Error(msg)
		time.Sleep(5 * time.Second)
	}()

}

func (s *pluginManager) RunControllerController(webhookPort int, webhookTlsDir string) {
	s.logger.Sugar().Infof("setup controller webhook")
	s.runWebhook(webhookPort, webhookTlsDir)
	s.logger.Sugar().Infof("setup controller reconcile")
	s.runControllerReconcile()
}

func init() {
	globalPluginManager = &pluginManager{
		chainingPlugins: map[string]plugintypes.ChainingPlugin{},
	}
	globalPluginManager.chainingPlugins["nethttp"] = &nethttp.PluginNetHttp{}

}
