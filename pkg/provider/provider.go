package provider

import (
	"fmt"
)

const (
	MsgMissingProvider         = "provider '%s' required"
	MsgUnsupportedProviderType = "unsupported provider type '%d'"
)

type (
	ProviderType int

	InstanceProviderFunc func() interface{}

	ProviderConfig struct {
		ID   string
		Type ProviderType
		Impl InstanceProviderFunc
	}

	Provider struct {
		providers map[ProviderType]ProviderConfig
	}

	GenericProvider interface {
		Close() error
	}
)

// New creates a new provider instance and configures it with provider implementations
func New(opts ...ProviderConfig) (*Provider, error) {
	p := Provider{
		providers: make(map[ProviderType]ProviderConfig),
	}

	if err := p.RegisterProviders(false, opts...); err != nil {
		return nil, err
	}

	return &p, nil
}

// Close iterates over all registered ProviderConfigs and asks them to 'close'
func (p *Provider) Close() bool {
	hasError := false
	for _, pc := range p.providers {
		err := pc.Impl().(GenericProvider).Close()
		if err != nil {
			hasError = true
		}
	}
	return hasError
}

// RegisterProviders registers one or more providers. An existing provider will be overwritten
// if ignoreExists is true, otherwise the function returns an error.
func (p *Provider) RegisterProviders(ignoreExists bool, opts ...ProviderConfig) error {
	for _, opt := range opts {
		if _, ok := p.providers[opt.Type]; ok {
			if !ignoreExists {
				return fmt.Errorf("provider of type '%s' already registered", opt.ID)
			}
		}
		p.providers[opt.Type] = opt
	}
	return nil
}

// Find returns the registered provider instance if defined.
// The bool flag is set to true if there is a provider and false otherwise.
func (p *Provider) Find(providerType ProviderType) (interface{}, bool) {
	pc, ok := p.providers[providerType]
	if !ok {
		return nil, false
	}
	return pc.Impl(), true
}

// WithProvider returns a populated ProviderConfig struct.
func WithProvider(ID string, providerType ProviderType, impl InstanceProviderFunc) ProviderConfig {
	return ProviderConfig{
		ID:   ID,
		Type: providerType,
		Impl: impl,
	}
}
