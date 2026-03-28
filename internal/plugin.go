// Package internal implements the workflow-plugin-platform external plugin,
// providing platform infrastructure module types and pipeline step types.
package internal

import (
	"fmt"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// platformPlugin implements sdk.PluginProvider.
type platformPlugin struct{}

// NewPlatformPlugin returns a new platformPlugin instance.
func NewPlatformPlugin() sdk.PluginProvider {
	return &platformPlugin{}
}

// Manifest returns plugin metadata.
func (p *platformPlugin) Manifest() sdk.PluginManifest {
	return sdk.PluginManifest{
		Name:        "workflow-plugin-platform",
		Version:     "0.1.0",
		Author:      "GoCodeAlone",
		Description: "Platform infrastructure modules: Kubernetes, ECS, DNS, networking, API gateway, autoscaling, multi-region, DigitalOcean, iac.state, app.container",
	}
}

// ModuleTypes returns the module type names this plugin provides.
func (p *platformPlugin) ModuleTypes() []string {
	return []string{
		"platform.kubernetes",
		"platform.ecs",
		"platform.dns",
		"platform.networking",
		"platform.apigateway",
		"platform.autoscaling",
		"platform.provider",
		"platform.resource",
		"platform.context",
		"platform.region",
		"platform.region_router",
		"platform.doks",
		"platform.do_networking",
		"platform.do_dns",
		"platform.do_app",
		"platform.do_database",
		"app.container",
		"iac.state",
		"argo.workflows",
	}
}

// CreateModule creates a module instance of the given type.
func (p *platformPlugin) CreateModule(typeName, name string, config map[string]any) (sdk.ModuleInstance, error) {
	switch typeName {
	case "platform.kubernetes":
		return newGenericModule(name, typeName, config), nil
	case "platform.ecs":
		return newGenericModule(name, typeName, config), nil
	case "platform.dns":
		return newGenericModule(name, typeName, config), nil
	case "platform.networking":
		return newGenericModule(name, typeName, config), nil
	case "platform.apigateway":
		return newGenericModule(name, typeName, config), nil
	case "platform.autoscaling":
		return newGenericModule(name, typeName, config), nil
	case "platform.provider":
		return newGenericModule(name, typeName, config), nil
	case "platform.resource":
		return newGenericModule(name, typeName, config), nil
	case "platform.context":
		return newGenericModule(name, typeName, config), nil
	case "platform.region":
		return newGenericModule(name, typeName, config), nil
	case "platform.region_router":
		return newGenericModule(name, typeName, config), nil
	case "platform.doks":
		return newGenericModule(name, typeName, config), nil
	case "platform.do_networking":
		return newGenericModule(name, typeName, config), nil
	case "platform.do_dns":
		return newGenericModule(name, typeName, config), nil
	case "platform.do_app":
		return newGenericModule(name, typeName, config), nil
	case "platform.do_database":
		return newGenericModule(name, typeName, config), nil
	case "app.container":
		return newGenericModule(name, typeName, config), nil
	case "iac.state":
		return newIaCStateModule(name, config), nil
	case "argo.workflows":
		return newGenericModule(name, typeName, config), nil
	default:
		return nil, fmt.Errorf("platform plugin: unknown module type %q", typeName)
	}
}

// StepTypes returns the step type names this plugin provides.
func (p *platformPlugin) StepTypes() []string {
	return []string{
		"step.platform_template",
		"step.k8s_plan",
		"step.k8s_apply",
		"step.k8s_status",
		"step.k8s_destroy",
		"step.ecs_plan",
		"step.ecs_apply",
		"step.ecs_status",
		"step.ecs_destroy",
		"step.iac_plan",
		"step.iac_apply",
		"step.iac_status",
		"step.iac_destroy",
		"step.iac_drift_detect",
		"step.dns_plan",
		"step.dns_apply",
		"step.dns_status",
		"step.network_plan",
		"step.network_apply",
		"step.network_status",
		"step.apigw_plan",
		"step.apigw_apply",
		"step.apigw_status",
		"step.apigw_destroy",
		"step.scaling_plan",
		"step.scaling_apply",
		"step.scaling_status",
		"step.scaling_destroy",
		"step.app_deploy",
		"step.app_status",
		"step.app_rollback",
		"step.region_deploy",
		"step.region_promote",
		"step.region_failover",
		"step.region_status",
		"step.region_weight",
		"step.region_sync",
		"step.argo_submit",
		"step.argo_status",
		"step.argo_logs",
		"step.argo_delete",
		"step.argo_list",
		"step.do_deploy",
		"step.do_status",
		"step.do_logs",
		"step.do_scale",
		"step.do_destroy",
	}
}

// CreateStep creates a step instance of the given type.
func (p *platformPlugin) CreateStep(typeName, name string, config map[string]any) (sdk.StepInstance, error) {
	switch typeName {
	case "step.platform_template",
		"step.k8s_plan", "step.k8s_apply", "step.k8s_status", "step.k8s_destroy",
		"step.ecs_plan", "step.ecs_apply", "step.ecs_status", "step.ecs_destroy",
		"step.iac_plan", "step.iac_apply", "step.iac_status", "step.iac_destroy", "step.iac_drift_detect",
		"step.dns_plan", "step.dns_apply", "step.dns_status",
		"step.network_plan", "step.network_apply", "step.network_status",
		"step.apigw_plan", "step.apigw_apply", "step.apigw_status", "step.apigw_destroy",
		"step.scaling_plan", "step.scaling_apply", "step.scaling_status", "step.scaling_destroy",
		"step.app_deploy", "step.app_status", "step.app_rollback",
		"step.region_deploy", "step.region_promote", "step.region_failover",
		"step.region_status", "step.region_weight", "step.region_sync",
		"step.argo_submit", "step.argo_status", "step.argo_logs", "step.argo_delete", "step.argo_list",
		"step.do_deploy", "step.do_status", "step.do_logs", "step.do_scale", "step.do_destroy":
		return newGenericStep(name, typeName, config), nil
	default:
		return nil, fmt.Errorf("platform plugin: unknown step type %q", typeName)
	}
}
