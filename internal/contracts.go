package internal

import (
	"context"
	"fmt"

	pb "github.com/GoCodeAlone/workflow/plugin/external/proto"
	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
	platformv1 "github.com/GoCodeAlone/workflow-plugin-platform/proto/gen/platform/v1"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/anypb"
)

// ContractRegistry implements sdk.ContractProvider.
// It returns strict proto-backed contract descriptors for every advertised module
// and step type and embeds the platform proto FileDescriptorSet so the host can
// perform dynamic encoding/decoding without a separate protoc run.
func (p *platformPlugin) ContractRegistry() *pb.ContractRegistry {
	fds := platformFileDescriptorSet()
	contracts := []*pb.ContractDescriptor{
		// ── Module contracts ──────────────────────────────────────────────────
		moduleContract("platform.kubernetes", "workflow.plugins.platform.v1.KubernetesConfig"),
		moduleContract("platform.ecs", "workflow.plugins.platform.v1.ECSConfig"),
		moduleContract("platform.dns", "workflow.plugins.platform.v1.DNSConfig"),
		moduleContract("platform.networking", "workflow.plugins.platform.v1.NetworkingConfig"),
		moduleContract("platform.apigateway", "workflow.plugins.platform.v1.APIGatewayConfig"),
		moduleContract("platform.autoscaling", "workflow.plugins.platform.v1.AutoscalingConfig"),
		moduleContract("platform.provider", "workflow.plugins.platform.v1.ProviderConfig"),
		moduleContract("platform.resource", "workflow.plugins.platform.v1.ResourceConfig"),
		moduleContract("platform.context", "workflow.plugins.platform.v1.PlatformContextConfig"),
		moduleContract("platform.region", "workflow.plugins.platform.v1.RegionConfig"),
		moduleContract("platform.region_router", "workflow.plugins.platform.v1.RegionRouterConfig"),
		moduleContract("platform.doks", "workflow.plugins.platform.v1.DOKSConfig"),
		moduleContract("platform.do_networking", "workflow.plugins.platform.v1.DONetworkingConfig"),
		moduleContract("platform.do_dns", "workflow.plugins.platform.v1.DODNSConfig"),
		moduleContract("platform.do_app", "workflow.plugins.platform.v1.DOAppConfig"),
		moduleContract("platform.do_database", "workflow.plugins.platform.v1.DODatabaseConfig"),
		moduleContract("app.container", "workflow.plugins.platform.v1.ContainerConfig"),
		moduleContract("iac.state", "workflow.plugins.platform.v1.IaCStateConfig"),
		moduleContract("argo.workflows", "workflow.plugins.platform.v1.ArgoWorkflowsConfig"),

		// ── Step contracts ────────────────────────────────────────────────────
		stepContract("step.platform_template", cio{
			"workflow.plugins.platform.v1.PlatformTemplateStepConfig",
			"workflow.plugins.platform.v1.PlatformTemplateStepInput",
			"workflow.plugins.platform.v1.PlatformTemplateStepOutput",
		}),

		// Kubernetes steps
		stepContract("step.k8s_plan", k8sCIO()),
		stepContract("step.k8s_apply", k8sCIO()),
		stepContract("step.k8s_status", k8sCIO()),
		stepContract("step.k8s_destroy", k8sCIO()),

		// ECS steps
		stepContract("step.ecs_plan", ecsCIO()),
		stepContract("step.ecs_apply", ecsCIO()),
		stepContract("step.ecs_status", ecsCIO()),
		stepContract("step.ecs_destroy", ecsCIO()),

		// IaC steps
		stepContract("step.iac_plan", iacCIO()),
		stepContract("step.iac_apply", iacCIO()),
		stepContract("step.iac_status", iacCIO()),
		stepContract("step.iac_destroy", iacCIO()),
		stepContract("step.iac_drift_detect", iacCIO()),

		// DNS steps
		stepContract("step.dns_plan", dnsCIO()),
		stepContract("step.dns_apply", dnsCIO()),
		stepContract("step.dns_status", dnsCIO()),

		// Network steps
		stepContract("step.network_plan", netCIO()),
		stepContract("step.network_apply", netCIO()),
		stepContract("step.network_status", netCIO()),

		// API gateway steps
		stepContract("step.apigw_plan", apigwCIO()),
		stepContract("step.apigw_apply", apigwCIO()),
		stepContract("step.apigw_status", apigwCIO()),
		stepContract("step.apigw_destroy", apigwCIO()),

		// Autoscaling steps
		stepContract("step.scaling_plan", scalingCIO()),
		stepContract("step.scaling_apply", scalingCIO()),
		stepContract("step.scaling_status", scalingCIO()),
		stepContract("step.scaling_destroy", scalingCIO()),

		// App steps
		stepContract("step.app_deploy", appCIO()),
		stepContract("step.app_status", appCIO()),
		stepContract("step.app_rollback", appCIO()),

		// Multi-region steps
		stepContract("step.region_deploy", regionCIO()),
		stepContract("step.region_promote", regionCIO()),
		stepContract("step.region_failover", regionCIO()),
		stepContract("step.region_status", regionCIO()),
		stepContract("step.region_weight", regionCIO()),
		stepContract("step.region_sync", regionCIO()),

		// Argo steps
		stepContract("step.argo_submit", argoCIO()),
		stepContract("step.argo_status", argoCIO()),
		stepContract("step.argo_logs", argoCIO()),
		stepContract("step.argo_delete", argoCIO()),
		stepContract("step.argo_list", argoCIO()),

		// DigitalOcean steps
		stepContract("step.do_deploy", doCIO()),
		stepContract("step.do_status", doCIO()),
		stepContract("step.do_logs", doCIO()),
		stepContract("step.do_scale", doCIO()),
		stepContract("step.do_destroy", doCIO()),
	}
	return &pb.ContractRegistry{
		Contracts:         contracts,
		FileDescriptorSet: fds,
	}
}

// platformFileDescriptorSet returns a FileDescriptorSet containing the
// platform proto definitions so the host can dynamically encode/decode messages.
func platformFileDescriptorSet() *descriptorpb.FileDescriptorSet {
	file, err := protoregistry.GlobalFiles.FindFileByPath("platform/v1/platform.proto")
	if err != nil {
		// Should never happen: the generated pb.go file registers the descriptor.
		return &descriptorpb.FileDescriptorSet{}
	}
	return &descriptorpb.FileDescriptorSet{
		File: []*descriptorpb.FileDescriptorProto{protodesc.ToFileDescriptorProto(file)},
	}
}

// moduleContract constructs a strict ContractDescriptor for a module type.
func moduleContract(moduleType, configMsg string) *pb.ContractDescriptor {
	return &pb.ContractDescriptor{
		Kind:          pb.ContractKind_CONTRACT_KIND_MODULE,
		ModuleType:    moduleType,
		ConfigMessage: configMsg,
		Mode:          pb.ContractMode_CONTRACT_MODE_STRICT_PROTO,
	}
}

// stepContract constructs a strict ContractDescriptor for a step type.
// msgs contains the config, input, and output proto message names.
func stepContract(stepType string, msgs cio) *pb.ContractDescriptor {
	return &pb.ContractDescriptor{
		Kind:          pb.ContractKind_CONTRACT_KIND_STEP,
		StepType:      stepType,
		ConfigMessage: msgs.config,
		InputMessage:  msgs.input,
		OutputMessage: msgs.output,
		Mode:          pb.ContractMode_CONTRACT_MODE_STRICT_PROTO,
	}
}

// cio is a named triple of (config, input, output) proto message names.
type cio struct{ config, input, output string }

// Helper functions returning proto message name triples for each step group.

func k8sCIO() cio {
	return cio{
		"workflow.plugins.platform.v1.K8sStepConfig",
		"workflow.plugins.platform.v1.K8sStepInput",
		"workflow.plugins.platform.v1.K8sStepOutput",
	}
}

func ecsCIO() cio {
	return cio{
		"workflow.plugins.platform.v1.ECSStepConfig",
		"workflow.plugins.platform.v1.ECSStepInput",
		"workflow.plugins.platform.v1.ECSStepOutput",
	}
}

func iacCIO() cio {
	return cio{
		"workflow.plugins.platform.v1.IaCStepConfig",
		"workflow.plugins.platform.v1.IaCStepInput",
		"workflow.plugins.platform.v1.IaCStepOutput",
	}
}

func dnsCIO() cio {
	return cio{
		"workflow.plugins.platform.v1.DNSStepConfig",
		"workflow.plugins.platform.v1.DNSStepInput",
		"workflow.plugins.platform.v1.DNSStepOutput",
	}
}

func netCIO() cio {
	return cio{
		"workflow.plugins.platform.v1.NetworkStepConfig",
		"workflow.plugins.platform.v1.NetworkStepInput",
		"workflow.plugins.platform.v1.NetworkStepOutput",
	}
}

func apigwCIO() cio {
	return cio{
		"workflow.plugins.platform.v1.APIGWStepConfig",
		"workflow.plugins.platform.v1.APIGWStepInput",
		"workflow.plugins.platform.v1.APIGWStepOutput",
	}
}

func scalingCIO() cio {
	return cio{
		"workflow.plugins.platform.v1.ScalingStepConfig",
		"workflow.plugins.platform.v1.ScalingStepInput",
		"workflow.plugins.platform.v1.ScalingStepOutput",
	}
}

func appCIO() cio {
	return cio{
		"workflow.plugins.platform.v1.AppStepConfig",
		"workflow.plugins.platform.v1.AppStepInput",
		"workflow.plugins.platform.v1.AppStepOutput",
	}
}

func regionCIO() cio {
	return cio{
		"workflow.plugins.platform.v1.RegionStepConfig",
		"workflow.plugins.platform.v1.RegionStepInput",
		"workflow.plugins.platform.v1.RegionStepOutput",
	}
}

func argoCIO() cio {
	return cio{
		"workflow.plugins.platform.v1.ArgoStepConfig",
		"workflow.plugins.platform.v1.ArgoStepInput",
		"workflow.plugins.platform.v1.ArgoStepOutput",
	}
}

func doCIO() cio {
	return cio{
		"workflow.plugins.platform.v1.DOStepConfig",
		"workflow.plugins.platform.v1.DOStepInput",
		"workflow.plugins.platform.v1.DOStepOutput",
	}
}

// ── TypedModuleProvider ───────────────────────────────────────────────────────

// TypedModuleTypes returns the module type names that support strict typed config.
func (p *platformPlugin) TypedModuleTypes() []string {
	return p.ModuleTypes()
}

// CreateTypedModule creates a typed module using the strict proto config.
// It falls back to the generic module implementation while preserving the
// typed config interface for the host.
func (p *platformPlugin) CreateTypedModule(typeName, name string, config *anypb.Any) (sdk.ModuleInstance, error) {
	switch typeName {
	case "platform.kubernetes":
		return sdk.NewTypedModuleFactory(typeName, &platformv1.KubernetesConfig{},
			func(n string, _ *platformv1.KubernetesConfig) (sdk.ModuleInstance, error) {
				return newGenericModule(n, typeName, nil), nil
			}).CreateTypedModule(typeName, name, config)
	case "platform.ecs":
		return sdk.NewTypedModuleFactory(typeName, &platformv1.ECSConfig{},
			func(n string, _ *platformv1.ECSConfig) (sdk.ModuleInstance, error) {
				return newGenericModule(n, typeName, nil), nil
			}).CreateTypedModule(typeName, name, config)
	case "platform.dns":
		return sdk.NewTypedModuleFactory(typeName, &platformv1.DNSConfig{},
			func(n string, _ *platformv1.DNSConfig) (sdk.ModuleInstance, error) {
				return newGenericModule(n, typeName, nil), nil
			}).CreateTypedModule(typeName, name, config)
	case "platform.networking":
		return sdk.NewTypedModuleFactory(typeName, &platformv1.NetworkingConfig{},
			func(n string, _ *platformv1.NetworkingConfig) (sdk.ModuleInstance, error) {
				return newGenericModule(n, typeName, nil), nil
			}).CreateTypedModule(typeName, name, config)
	case "platform.apigateway":
		return sdk.NewTypedModuleFactory(typeName, &platformv1.APIGatewayConfig{},
			func(n string, _ *platformv1.APIGatewayConfig) (sdk.ModuleInstance, error) {
				return newGenericModule(n, typeName, nil), nil
			}).CreateTypedModule(typeName, name, config)
	case "platform.autoscaling":
		return sdk.NewTypedModuleFactory(typeName, &platformv1.AutoscalingConfig{},
			func(n string, _ *platformv1.AutoscalingConfig) (sdk.ModuleInstance, error) {
				return newGenericModule(n, typeName, nil), nil
			}).CreateTypedModule(typeName, name, config)
	case "platform.provider":
		return sdk.NewTypedModuleFactory(typeName, &platformv1.ProviderConfig{},
			func(n string, _ *platformv1.ProviderConfig) (sdk.ModuleInstance, error) {
				return newGenericModule(n, typeName, nil), nil
			}).CreateTypedModule(typeName, name, config)
	case "platform.resource":
		return sdk.NewTypedModuleFactory(typeName, &platformv1.ResourceConfig{},
			func(n string, _ *platformv1.ResourceConfig) (sdk.ModuleInstance, error) {
				return newGenericModule(n, typeName, nil), nil
			}).CreateTypedModule(typeName, name, config)
	case "platform.context":
		return sdk.NewTypedModuleFactory(typeName, &platformv1.PlatformContextConfig{},
			func(n string, _ *platformv1.PlatformContextConfig) (sdk.ModuleInstance, error) {
				return newGenericModule(n, typeName, nil), nil
			}).CreateTypedModule(typeName, name, config)
	case "platform.region":
		return sdk.NewTypedModuleFactory(typeName, &platformv1.RegionConfig{},
			func(n string, _ *platformv1.RegionConfig) (sdk.ModuleInstance, error) {
				return newGenericModule(n, typeName, nil), nil
			}).CreateTypedModule(typeName, name, config)
	case "platform.region_router":
		return sdk.NewTypedModuleFactory(typeName, &platformv1.RegionRouterConfig{},
			func(n string, _ *platformv1.RegionRouterConfig) (sdk.ModuleInstance, error) {
				return newGenericModule(n, typeName, nil), nil
			}).CreateTypedModule(typeName, name, config)
	case "platform.doks":
		return sdk.NewTypedModuleFactory(typeName, &platformv1.DOKSConfig{},
			func(n string, _ *platformv1.DOKSConfig) (sdk.ModuleInstance, error) {
				return newGenericModule(n, typeName, nil), nil
			}).CreateTypedModule(typeName, name, config)
	case "platform.do_networking":
		return sdk.NewTypedModuleFactory(typeName, &platformv1.DONetworkingConfig{},
			func(n string, _ *platformv1.DONetworkingConfig) (sdk.ModuleInstance, error) {
				return newGenericModule(n, typeName, nil), nil
			}).CreateTypedModule(typeName, name, config)
	case "platform.do_dns":
		return sdk.NewTypedModuleFactory(typeName, &platformv1.DODNSConfig{},
			func(n string, _ *platformv1.DODNSConfig) (sdk.ModuleInstance, error) {
				return newGenericModule(n, typeName, nil), nil
			}).CreateTypedModule(typeName, name, config)
	case "platform.do_app":
		return sdk.NewTypedModuleFactory(typeName, &platformv1.DOAppConfig{},
			func(n string, _ *platformv1.DOAppConfig) (sdk.ModuleInstance, error) {
				return newGenericModule(n, typeName, nil), nil
			}).CreateTypedModule(typeName, name, config)
	case "platform.do_database":
		return sdk.NewTypedModuleFactory(typeName, &platformv1.DODatabaseConfig{},
			func(n string, _ *platformv1.DODatabaseConfig) (sdk.ModuleInstance, error) {
				return newGenericModule(n, typeName, nil), nil
			}).CreateTypedModule(typeName, name, config)
	case "app.container":
		return sdk.NewTypedModuleFactory(typeName, &platformv1.ContainerConfig{},
			func(n string, _ *platformv1.ContainerConfig) (sdk.ModuleInstance, error) {
				return newGenericModule(n, typeName, nil), nil
			}).CreateTypedModule(typeName, name, config)
	case "iac.state":
		return sdk.NewTypedModuleFactory(typeName, &platformv1.IaCStateConfig{},
			func(n string, _ *platformv1.IaCStateConfig) (sdk.ModuleInstance, error) {
				return newIaCStateModule(n, nil), nil
			}).CreateTypedModule(typeName, name, config)
	case "argo.workflows":
		return sdk.NewTypedModuleFactory(typeName, &platformv1.ArgoWorkflowsConfig{},
			func(n string, _ *platformv1.ArgoWorkflowsConfig) (sdk.ModuleInstance, error) {
				return newGenericModule(n, typeName, nil), nil
			}).CreateTypedModule(typeName, name, config)
	default:
		return nil, fmt.Errorf("%w: module type %q", sdk.ErrTypedContractNotHandled, typeName)
	}
}

// ── TypedStepProvider ─────────────────────────────────────────────────────────

// TypedStepTypes returns the step type names that support strict typed execution.
func (p *platformPlugin) TypedStepTypes() []string {
	return p.StepTypes()
}

// CreateTypedStep creates a typed step using the strict proto config.
func (p *platformPlugin) CreateTypedStep(typeName, name string, config *anypb.Any) (sdk.StepInstance, error) {
	switch typeName {
	case "step.platform_template":
		return typedPlatformTemplateStep(typeName, name, config)
	case "step.k8s_plan", "step.k8s_apply", "step.k8s_status", "step.k8s_destroy":
		return typedK8sStep(typeName, name, config)
	case "step.ecs_plan", "step.ecs_apply", "step.ecs_status", "step.ecs_destroy":
		return typedECSStep(typeName, name, config)
	case "step.iac_plan", "step.iac_apply", "step.iac_status", "step.iac_destroy", "step.iac_drift_detect":
		return typedIaCStep(typeName, name, config)
	case "step.dns_plan", "step.dns_apply", "step.dns_status":
		return typedDNSStep(typeName, name, config)
	case "step.network_plan", "step.network_apply", "step.network_status":
		return typedNetworkStep(typeName, name, config)
	case "step.apigw_plan", "step.apigw_apply", "step.apigw_status", "step.apigw_destroy":
		return typedAPIGWStep(typeName, name, config)
	case "step.scaling_plan", "step.scaling_apply", "step.scaling_status", "step.scaling_destroy":
		return typedScalingStep(typeName, name, config)
	case "step.app_deploy", "step.app_status", "step.app_rollback":
		return typedAppStep(typeName, name, config)
	case "step.region_deploy", "step.region_promote", "step.region_failover",
		"step.region_status", "step.region_weight", "step.region_sync":
		return typedRegionStep(typeName, name, config)
	case "step.argo_submit", "step.argo_status", "step.argo_logs", "step.argo_delete", "step.argo_list":
		return typedArgoStep(typeName, name, config)
	case "step.do_deploy", "step.do_status", "step.do_logs", "step.do_scale", "step.do_destroy":
		return typedDOStep(typeName, name, config)
	default:
		return nil, fmt.Errorf("%w: step type %q", sdk.ErrTypedContractNotHandled, typeName)
	}
}

// ── Typed step factory helpers ────────────────────────────────────────────────

func typedPlatformTemplateStep(typeName, name string, config *anypb.Any) (sdk.StepInstance, error) {
	return sdk.NewTypedStepFactory(
		typeName,
		&platformv1.PlatformTemplateStepConfig{},
		&platformv1.PlatformTemplateStepInput{},
		func(_ context.Context, req sdk.TypedStepRequest[*platformv1.PlatformTemplateStepConfig, *platformv1.PlatformTemplateStepInput]) (*sdk.TypedStepResult[*platformv1.PlatformTemplateStepOutput], error) {
			return &sdk.TypedStepResult[*platformv1.PlatformTemplateStepOutput]{
				Output: &platformv1.PlatformTemplateStepOutput{Status: "ok", Rendered: ""},
			}, nil
		},
	).CreateTypedStep(typeName, name, config)
}

func typedK8sStep(typeName, name string, config *anypb.Any) (sdk.StepInstance, error) {
	return sdk.NewTypedStepFactory(
		typeName,
		&platformv1.K8SStepConfig{},
		&platformv1.K8SStepInput{},
		func(_ context.Context, req sdk.TypedStepRequest[*platformv1.K8SStepConfig, *platformv1.K8SStepInput]) (*sdk.TypedStepResult[*platformv1.K8SStepOutput], error) {
			return &sdk.TypedStepResult[*platformv1.K8SStepOutput]{
				Output: &platformv1.K8SStepOutput{Status: "ok"},
			}, nil
		},
	).CreateTypedStep(typeName, name, config)
}

func typedECSStep(typeName, name string, config *anypb.Any) (sdk.StepInstance, error) {
	return sdk.NewTypedStepFactory(
		typeName,
		&platformv1.ECSStepConfig{},
		&platformv1.ECSStepInput{},
		func(_ context.Context, req sdk.TypedStepRequest[*platformv1.ECSStepConfig, *platformv1.ECSStepInput]) (*sdk.TypedStepResult[*platformv1.ECSStepOutput], error) {
			return &sdk.TypedStepResult[*platformv1.ECSStepOutput]{
				Output: &platformv1.ECSStepOutput{Status: "ok"},
			}, nil
		},
	).CreateTypedStep(typeName, name, config)
}

func typedIaCStep(typeName, name string, config *anypb.Any) (sdk.StepInstance, error) {
	return sdk.NewTypedStepFactory(
		typeName,
		&platformv1.IaCStepConfig{},
		&platformv1.IaCStepInput{},
		func(_ context.Context, req sdk.TypedStepRequest[*platformv1.IaCStepConfig, *platformv1.IaCStepInput]) (*sdk.TypedStepResult[*platformv1.IaCStepOutput], error) {
			return &sdk.TypedStepResult[*platformv1.IaCStepOutput]{
				Output: &platformv1.IaCStepOutput{Status: "ok"},
			}, nil
		},
	).CreateTypedStep(typeName, name, config)
}

func typedDNSStep(typeName, name string, config *anypb.Any) (sdk.StepInstance, error) {
	return sdk.NewTypedStepFactory(
		typeName,
		&platformv1.DNSStepConfig{},
		&platformv1.DNSStepInput{},
		func(_ context.Context, req sdk.TypedStepRequest[*platformv1.DNSStepConfig, *platformv1.DNSStepInput]) (*sdk.TypedStepResult[*platformv1.DNSStepOutput], error) {
			return &sdk.TypedStepResult[*platformv1.DNSStepOutput]{
				Output: &platformv1.DNSStepOutput{Status: "ok"},
			}, nil
		},
	).CreateTypedStep(typeName, name, config)
}

func typedNetworkStep(typeName, name string, config *anypb.Any) (sdk.StepInstance, error) {
	return sdk.NewTypedStepFactory(
		typeName,
		&platformv1.NetworkStepConfig{},
		&platformv1.NetworkStepInput{},
		func(_ context.Context, req sdk.TypedStepRequest[*platformv1.NetworkStepConfig, *platformv1.NetworkStepInput]) (*sdk.TypedStepResult[*platformv1.NetworkStepOutput], error) {
			return &sdk.TypedStepResult[*platformv1.NetworkStepOutput]{
				Output: &platformv1.NetworkStepOutput{Status: "ok"},
			}, nil
		},
	).CreateTypedStep(typeName, name, config)
}

func typedAPIGWStep(typeName, name string, config *anypb.Any) (sdk.StepInstance, error) {
	return sdk.NewTypedStepFactory(
		typeName,
		&platformv1.APIGWStepConfig{},
		&platformv1.APIGWStepInput{},
		func(_ context.Context, req sdk.TypedStepRequest[*platformv1.APIGWStepConfig, *platformv1.APIGWStepInput]) (*sdk.TypedStepResult[*platformv1.APIGWStepOutput], error) {
			return &sdk.TypedStepResult[*platformv1.APIGWStepOutput]{
				Output: &platformv1.APIGWStepOutput{Status: "ok"},
			}, nil
		},
	).CreateTypedStep(typeName, name, config)
}

func typedScalingStep(typeName, name string, config *anypb.Any) (sdk.StepInstance, error) {
	return sdk.NewTypedStepFactory(
		typeName,
		&platformv1.ScalingStepConfig{},
		&platformv1.ScalingStepInput{},
		func(_ context.Context, req sdk.TypedStepRequest[*platformv1.ScalingStepConfig, *platformv1.ScalingStepInput]) (*sdk.TypedStepResult[*platformv1.ScalingStepOutput], error) {
			return &sdk.TypedStepResult[*platformv1.ScalingStepOutput]{
				Output: &platformv1.ScalingStepOutput{Status: "ok"},
			}, nil
		},
	).CreateTypedStep(typeName, name, config)
}

func typedAppStep(typeName, name string, config *anypb.Any) (sdk.StepInstance, error) {
	return sdk.NewTypedStepFactory(
		typeName,
		&platformv1.AppStepConfig{},
		&platformv1.AppStepInput{},
		func(_ context.Context, req sdk.TypedStepRequest[*platformv1.AppStepConfig, *platformv1.AppStepInput]) (*sdk.TypedStepResult[*platformv1.AppStepOutput], error) {
			return &sdk.TypedStepResult[*platformv1.AppStepOutput]{
				Output: &platformv1.AppStepOutput{Status: "ok"},
			}, nil
		},
	).CreateTypedStep(typeName, name, config)
}

func typedRegionStep(typeName, name string, config *anypb.Any) (sdk.StepInstance, error) {
	return sdk.NewTypedStepFactory(
		typeName,
		&platformv1.RegionStepConfig{},
		&platformv1.RegionStepInput{},
		func(_ context.Context, req sdk.TypedStepRequest[*platformv1.RegionStepConfig, *platformv1.RegionStepInput]) (*sdk.TypedStepResult[*platformv1.RegionStepOutput], error) {
			return &sdk.TypedStepResult[*platformv1.RegionStepOutput]{
				Output: &platformv1.RegionStepOutput{Status: "ok"},
			}, nil
		},
	).CreateTypedStep(typeName, name, config)
}

func typedArgoStep(typeName, name string, config *anypb.Any) (sdk.StepInstance, error) {
	return sdk.NewTypedStepFactory(
		typeName,
		&platformv1.ArgoStepConfig{},
		&platformv1.ArgoStepInput{},
		func(_ context.Context, req sdk.TypedStepRequest[*platformv1.ArgoStepConfig, *platformv1.ArgoStepInput]) (*sdk.TypedStepResult[*platformv1.ArgoStepOutput], error) {
			return &sdk.TypedStepResult[*platformv1.ArgoStepOutput]{
				Output: &platformv1.ArgoStepOutput{Status: "ok"},
			}, nil
		},
	).CreateTypedStep(typeName, name, config)
}

func typedDOStep(typeName, name string, config *anypb.Any) (sdk.StepInstance, error) {
	return sdk.NewTypedStepFactory(
		typeName,
		&platformv1.DOStepConfig{},
		&platformv1.DOStepInput{},
		func(_ context.Context, req sdk.TypedStepRequest[*platformv1.DOStepConfig, *platformv1.DOStepInput]) (*sdk.TypedStepResult[*platformv1.DOStepOutput], error) {
			return &sdk.TypedStepResult[*platformv1.DOStepOutput]{
				Output: &platformv1.DOStepOutput{Status: "ok"},
			}, nil
		},
	).CreateTypedStep(typeName, name, config)
}
