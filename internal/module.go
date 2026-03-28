package internal

import (
	"context"
)

// genericModule is a stub ModuleInstance for platform module types.
// TODO: Implement actual provisioning logic per module type.
type genericModule struct {
	name       string
	moduleType string
	config     map[string]any
}

func newGenericModule(name, moduleType string, config map[string]any) *genericModule {
	return &genericModule{name: name, moduleType: moduleType, config: config}
}

func (m *genericModule) Init() error { return nil }

func (m *genericModule) Start(_ context.Context) error { return nil }

func (m *genericModule) Stop(_ context.Context) error { return nil }

// iacStateModule provides a simple in-memory IaC state store.
// TODO: Support filesystem and remote backends (GCS, S3, Azure Blob, PostgreSQL).
type iacStateModule struct {
	name   string
	config map[string]any
	state  map[string]any
}

func newIaCStateModule(name string, config map[string]any) *iacStateModule {
	return &iacStateModule{name: name, config: config, state: make(map[string]any)}
}

func (m *iacStateModule) Init() error { return nil }

func (m *iacStateModule) Start(_ context.Context) error { return nil }

func (m *iacStateModule) Stop(_ context.Context) error { return nil }

// GetState returns the current state for a resource key.
func (m *iacStateModule) GetState(key string) (map[string]any, bool) {
	v, ok := m.state[key]
	if !ok {
		return nil, false
	}
	if s, ok := v.(map[string]any); ok {
		return s, true
	}
	return nil, false
}

// SetState stores the state for a resource key.
func (m *iacStateModule) SetState(key string, value map[string]any) {
	m.state[key] = value
}
