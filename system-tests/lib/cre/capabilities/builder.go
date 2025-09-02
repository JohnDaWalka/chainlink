package capabilities

import (
	"github.com/pkg/errors"

	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

type registryConfigFns struct {
	V1 cre.CapabilityRegistryConfigFn
	V2 cre.CapabilityRegistryConfigFn
}

type Capability struct {
	flag                      cre.CapabilityFlag
	jobSpecFn                 cre.JobSpecFn
	nodeConfigFn              cre.NodeConfigFn
	gatewayJobHandlerConfigFn cre.GatewayHandlerConfigFn
	registryConfigFns         registryConfigFns
	validateFn                func(*Capability) error
}

func (c *Capability) Flag() cre.CapabilityFlag {
	return c.flag
}

func (c *Capability) JobSpecFn() cre.JobSpecFn {
	return c.jobSpecFn
}

func (c *Capability) NodeConfigFn() cre.NodeConfigFn {
	return c.nodeConfigFn
}

func (c *Capability) GatewayJobHandlerConfigFn() cre.GatewayHandlerConfigFn {
	return c.gatewayJobHandlerConfigFn
}

func (c *Capability) CapabilityRegistryV1ConfigFn() cre.CapabilityRegistryConfigFn {
	return c.registryConfigFns.V1
}

func (c *Capability) CapabilityRegistryV2ConfigFn() cre.CapabilityRegistryConfigFn {
	return c.registryConfigFns.V2
}

type Option func(*Capability)

func WithJobSpecFn(jobSpecFn cre.JobSpecFn) Option {
	return func(c *Capability) {
		c.jobSpecFn = jobSpecFn
	}
}

func WithNodeConfigFn(nodeConfigFn cre.NodeConfigFn) Option {
	return func(c *Capability) {
		c.nodeConfigFn = nodeConfigFn
	}
}

func WithGatewayJobHandlerConfigFn(gatewayJobHandlerConfigFn cre.GatewayHandlerConfigFn) Option {
	return func(c *Capability) {
		c.gatewayJobHandlerConfigFn = gatewayJobHandlerConfigFn
	}
}

func WithCapabilityRegistryV1ConfigFn(fn cre.CapabilityRegistryConfigFn) Option {
	return func(c *Capability) {
		c.registryConfigFns.V1 = fn
	}
}

func WithCapabilityRegistryV2ConfigFn(fn cre.CapabilityRegistryConfigFn) Option {
	return func(c *Capability) {
		c.registryConfigFns.V2 = fn
	}
}

func WithValidateFn(validateFn func(*Capability) error) Option {
	return func(c *Capability) {
		c.validateFn = validateFn
	}
}

func New(flag cre.CapabilityFlag, opts ...Option) (*Capability, error) {
	capability := &Capability{
		flag: flag,
		registryConfigFns: registryConfigFns{
			V1: unimplmentedConfigFn,
			V2: unimplmentedConfigFn,
		},
	}
	for _, opt := range opts {
		opt(capability)
	}

	if capability.validateFn != nil {
		if err := capability.validateFn(capability); err != nil {
			return nil, errors.Wrapf(err, "failed to validate capability %s", capability.flag)
		}
	}

	return capability, nil
}

func unimplmentedConfigFn(donFlags []cre.CapabilityFlag, nodeSetInput *cre.CapabilitiesAwareNodeSet) ([]keystone_changeset.DONCapabilityWithConfig, error) {
	return nil, errors.New("config function is not implemented")
}
