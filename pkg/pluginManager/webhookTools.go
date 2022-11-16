package pluginManager

import (
	"context"
	"k8s.io/apimachinery/pkg/runtime"
)

func (s *pluginWebhookhander) validateRequest(ctx context.Context, obj runtime.Object) error {
	// TODO: check
	return nil
}
