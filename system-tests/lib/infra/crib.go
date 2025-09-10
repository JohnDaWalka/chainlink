package infra

import (
	"fmt"
)

type Type = string

const (
	CRIB   Type = "crib"
	Docker Type = "docker"
)

type CribProvider = string

const (
	AWS  CribProvider = "aws"
	Kind CribProvider = "kind"
)

type Input struct {
	Type string     `toml:"type" validate:"oneof=crib docker"`
	CRIB *CRIBInput `toml:"crib"`
}

func (i *Input) IsCRIB() bool {
	return i.Type == CRIB
}

func (i *Input) IsDocker() bool {
	return i.Type == Docker
}

// Unfortunately, we need to construct some of these URLs before any environment is created, because they are used
// in CL node configs. This introduces a coupling between Helm charts used by CRIB and Docker container names used by CTFv2.
func (i *Input) InternalHost(nodeIndex int, isBootstrap bool, donName string) string {
	if i.IsCRIB() {
		if isBootstrap {
			return fmt.Sprintf("%s-bt-%d", donName, nodeIndex)
		}
		return fmt.Sprintf("%s-%d", donName, nodeIndex)
	}

	return fmt.Sprintf("%s-node%d", donName, nodeIndex)
}

func (i *Input) InternalGatewayHost(nodeIndex int, isBootstrap bool, donName string) string {
	if i.IsCRIB() {
		host := fmt.Sprintf("%s-%d", donName, nodeIndex)
		if isBootstrap {
			host = fmt.Sprintf("%s-bt-%d", donName, nodeIndex)
		}
		host += "-gtwnode"

		return host
	}

	return fmt.Sprintf("%s-node%d", donName, nodeIndex)
}

func (i *Input) ExternalGatewayHost() string {
	if i.IsCRIB() {
		return i.CRIB.Namespace + "-gateway.main.stage.cldev.sh"
	}

	return "localhost"
}

func (i *Input) ExternalGatewayPort(dockerPort int) int {
	if i.IsCRIB() {
		return 80
	}

	return dockerPort
}

type CRIBInput struct {
	Namespace string `toml:"namespace" validate:"required"`
	// absolute path to the folder with CRIB CRE
	FolderLocation string `toml:"folder_location" validate:"required"`
	Provider       string `toml:"provider" validate:"oneof=aws kind"`
	// required for cost attribution in AWS
	TeamInput *Team `toml:"team_input" validate:"required_if=Provider aws"`
}

// k8s cost attribution
type Team struct {
	Team       string `toml:"team" validate:"required"`
	Product    string `toml:"product" validate:"required"`
	CostCenter string `toml:"cost_center" validate:"required"`
	Component  string `toml:"component" validate:"required"`
}
