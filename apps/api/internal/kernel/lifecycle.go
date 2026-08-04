package kernel

import (
	"context"
	"fmt"
)

type Runtime struct {
	plan    Plan
	started []Module
}

func NewRuntime(plan Plan) *Runtime {
	return &Runtime{plan: plan}
}

func (r *Runtime) Start(ctx context.Context) error {
	for _, module := range r.plan.Modules {
		if module.Hooks.Start == nil {
			r.started = append(r.started, module)
			continue
		}
		if err := module.Hooks.Start(ctx); err != nil {
			cleanupErr := r.stopModules(ctx, r.started)
			r.started = nil
			if cleanupErr != nil {
				return kernelError(CodeLifecycleStartFailed, module.ID, "start failed: %v; cleanup failed: %v", err, cleanupErr)
			}
			return kernelError(CodeLifecycleStartFailed, module.ID, "start failed: %v", err)
		}
		r.started = append(r.started, module)
	}
	return nil
}

func (r *Runtime) Ready(ctx context.Context) error {
	for _, module := range r.plan.Modules {
		if module.Hooks.Ready == nil {
			continue
		}
		if err := module.Hooks.Ready(ctx); err != nil {
			return kernelError(CodeLifecycleReadyFailed, module.ID, "ready check failed: %v", err)
		}
	}
	return nil
}

func (r *Runtime) Stop(ctx context.Context) error {
	err := r.stopModules(ctx, r.started)
	r.started = nil
	return err
}

func (r *Runtime) stopModules(ctx context.Context, modules []Module) error {
	var firstErr error
	for i := len(modules) - 1; i >= 0; i-- {
		module := modules[i]
		if module.Hooks.Stop == nil {
			continue
		}
		if err := module.Hooks.Stop(ctx); err != nil && firstErr == nil {
			firstErr = kernelError(CodeLifecycleStopFailed, module.ID, "stop failed: %v", err)
		}
	}
	return firstErr
}

func (r *Runtime) String() string {
	return fmt.Sprintf("kernel runtime modules=%v", r.plan.IDs())
}
