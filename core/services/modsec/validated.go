package modsec

import (
	"fmt"

	"github.com/pelletier/go-toml"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

func ValidatedModsecSpec(tomlString string) (jb job.Job, err error) {
	var spec job.ModsecSpec
	tree, err := toml.Load(tomlString)
	if err != nil {
		return job.Job{}, fmt.Errorf("toml error on load: %w", err)
	}

	err = tree.Unmarshal(&spec)
	if err != nil {
		return job.Job{}, fmt.Errorf("toml unmarshal error on spec: %w", err)
	}

	err = tree.Unmarshal(&jb)
	if err != nil {
		return job.Job{}, fmt.Errorf("toml unmarshal error on job: %w", err)
	}

	jb.ModsecSpec = &spec

	if jb.Type != job.Modsec {
		return job.Job{}, fmt.Errorf("the only supported type is currently 'modsec', got %s", jb.Type)
	}

	if err = validate(jb.ModsecSpec); err != nil {
		return job.Job{}, err
	}

	return jb, nil
}
