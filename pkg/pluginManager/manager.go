package pluginManager

import (
	"github.com/spidernet-io/spiderdoctor/pkg/lock"
	plugintypes "github.com/spidernet-io/spiderdoctor/pkg/pluginManager/types"
	"github.com/spidernet-io/spiderdoctor/pkg/plugins/nethttp"
	"go.uber.org/zap"
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
	s.logger.Sugar().Infof("setup agent controller")
	s.runAgentReconcile()
}

func (s *pluginManager) RunControllerController(webhookPort int, webhookTlsDir string) {
	s.logger.Sugar().Infof("setup controller webhook")
	s.runWebhook(webhookPort, webhookTlsDir)
	s.logger.Sugar().Infof("setup controller controller")
	s.runControllerReconcile()
}

func init() {
	globalPluginManager = &pluginManager{
		chainingPlugins: map[string]plugintypes.ChainingPlugin{},
	}
	globalPluginManager.chainingPlugins["nethttp"] = &nethttp.PluginNetHttp{}

}
