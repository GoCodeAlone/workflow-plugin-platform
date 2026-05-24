// Command workflow-plugin-platform is a workflow engine external plugin that
// provides platform infrastructure modules: Kubernetes, ECS, DNS, networking,
// API gateway, autoscaling, multi-region, DigitalOcean platform types, and more.
// It runs as a subprocess and communicates with the host workflow engine via
// the go-plugin gRPC protocol.
package main

import (
	"github.com/GoCodeAlone/workflow-plugin-platform/internal"
	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

func main() {
	sdk.Serve(internal.NewPlatformPlugin(), sdk.WithBuildVersion(sdk.ResolveBuildVersion(internal.Version)))
}
