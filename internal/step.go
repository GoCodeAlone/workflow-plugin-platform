package internal

import (
	"context"
	"fmt"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// genericStep is a stub StepInstance for platform step types.
// TODO: Implement actual platform operations per step type.
type genericStep struct {
	name     string
	stepType string
	config   map[string]any
}

func newGenericStep(name, stepType string, config map[string]any) *genericStep {
	return &genericStep{name: name, stepType: stepType, config: config}
}

// Execute runs the platform step, returning a stub result.
// TODO: Implement real provisioning, status, and destroy operations.
func (s *genericStep) Execute(
	_ context.Context,
	_ map[string]any,
	_ map[string]map[string]any,
	_ map[string]any,
	_ map[string]any,
	_ map[string]any,
) (*sdk.StepResult, error) {
	module, _ := s.config["module"].(string)
	resource, _ := s.config["resource"].(string)

	return &sdk.StepResult{
		Output: map[string]any{
			"status":    "ok",
			"step_type": s.stepType,
			"module":    module,
			"resource":  resource,
			"message":   fmt.Sprintf("TODO: %s not yet implemented in external plugin", s.stepType),
		},
	}, nil
}
