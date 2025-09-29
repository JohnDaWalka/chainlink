package infra

type Type = string

const (
	CRIB          Type = "crib"
	Docker        Type = "docker"
	GriddleDevenv Type = "griddle-devenv"
)

type CribProvider = string

const (
	AWS  CribProvider = "aws"
	Kind CribProvider = "kind"
)

type Input struct {
	Type               string              `toml:"type" validate:"oneof=crib docker griddle-devenv"`
	CRIB               *CRIBInput          `toml:"crib"`
	GriddleDevenvInput *GriddleDevenvInput `toml:"devenv"`
}

// CRIBInput Deprecated, use GriddleDevenvInput instead
type CRIBInput struct {
	Namespace string `toml:"namespace" validate:"required"`
	// absolute path to the folder with CRIB CRE
	FolderLocation string `toml:"folder_location" validate:"required"`
	Provider       string `toml:"provider" validate:"oneof=aws kind"`
	// required for cost attribution in AWS
	TeamInput *Team `toml:"team_input" validate:"required_if=Provider aws"`
}

type GriddleDevenvInput struct {
	Namespace string `toml:"namespace" validate:"required"`
	// AWS account for deployment
	GriddleMetadata *GriddleMetadata `toml:"metadata" validate:"required"`
}

type GriddleMetadata struct {
	Account string `toml:"account" validate:"required"`
	Project string `toml:"project" validate:"required"`
	Service string `toml:"service" validate:"required"`
	Owner   string `toml:"owner" validate:"required"`
	Contact string `toml:"contact" validate:"required"`
}

// k8s cost attribution
type Team struct {
	Team       string `toml:"team" validate:"required"`
	Product    string `toml:"product" validate:"required"`
	CostCenter string `toml:"cost_center" validate:"required"`
	Component  string `toml:"component" validate:"required"`
}
