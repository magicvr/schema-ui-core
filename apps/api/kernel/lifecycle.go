package kernel

import (
	"context"
	"errors"
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
			return lifecycleFailure(CodeLifecycleStartFailed, module.ID, "start", err, cleanupErr)
		}
		r.started = append(r.started, module)
	}
	return nil
}

func (r *Runtime) Ready(ctx context.Context) error {
	for _, module := range r.started {
		if module.Hooks.Ready == nil {
			continue
		}
		if err := module.Hooks.Ready(ctx); err != nil {
			cleanupErr := r.stopModules(ctx, r.started)
			r.started = nil
			return lifecycleFailure(CodeLifecycleReadyFailed, module.ID, "ready check", err, cleanupErr)
		}
	}
	return nil
}

func lifecycleFailure(fallbackCode ErrorCode, fallbackModule, phase string, err, cleanupErr error) error {
	detail := fmt.Sprintf("%s failed: %v", phase, err)
	result := &Error{Code: fallbackCode, ModuleID: fallbackModule, Detail: detail}
	var structured *Error
	if errors.As(err, &structured) {
		copy := *structured
		result = &copy
	}
	if cleanupErr != nil {
		result.Detail += fmt.Sprintf("; cleanup failed: %v", cleanupErr)
	}
	return result
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
