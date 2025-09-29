package griddle

// GriddleConfig represents the complete griddle configuration file structure
type GriddleConfig struct {
	Version         int                       `yaml:"version"`
	Metadata        Metadata                  `yaml:"metadata"`
	DeployTemplates map[string]DeployTemplate `yaml:"deployTemplates"`
}

// Metadata represents the metadata section of the configuration
type Metadata struct {
	Account string `yaml:"account"`
	Project string `yaml:"project"`
	Service string `yaml:"service"`
	Owner   string `yaml:"owner"`
	Contact string `yaml:"contact"`
}

// DeployTemplate represents a deployment template with multiple instances
type DeployTemplate struct {
	Instances []Instance `yaml:"instances"`
}

// Instance represents a single deployment instance
type Instance struct {
	Name                string       `yaml:"name"`
	Chart               string       `yaml:"chart"`
	Version             string       `yaml:"version"`
	Repository          string       `yaml:"repository"`
	LocalRepositoryName string       `yaml:"localRepositoryName"`
	DependsOn           []Dependency `yaml:"dependsOn,omitempty"`
	Config              []string     `yaml:"config"`
}

// Dependency represents a dependency on another instance
type Dependency struct {
	Name string `yaml:"name"`
}
