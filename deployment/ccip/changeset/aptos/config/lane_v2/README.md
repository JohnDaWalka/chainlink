# Proposing configs for lanes

This proposes a standard and an interface for configuring lanes in a chain family agnostic way.

## Integrating a new family
### Changeset
Every chain family will have to create a changeset that supports upgrading all the contracts necessary to deploy a lane in it's family side.  

Having the chain as source requires:
- UpdateOnRampDestsConfig
- UpdateFeeQuoterPricesConfig
- UpdateFeeQuoterDestsConfig
- UpdateRouterRampsConfig  

Having the chain as destination requires:
- UpdateRouterRampsConfig
- UpdateOffRampSourcesConfig

### Interface Implementation
This changeset will take as input a struct that implements `UpdateLanesCfg` interface. And this struct needs to be added to the `NewUpdateLanesCfg` factory.

## Generic configs

The input will be a generic config `LaneConfig` which takes generic sources and destinations `ChainDefinition` configs. This structure can be extended for extra family specific configs. Aptos, for example, needs to know the version of the OnRamp contract for a lane.
