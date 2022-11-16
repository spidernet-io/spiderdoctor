package pluginManager

import (
	"context"
	"k8s.io/apimachinery/pkg/runtime"
)

func (s *pluginWebhookhander) CheckRequest(ctx context.Context, obj runtime.Object) error {
	return nil
}
