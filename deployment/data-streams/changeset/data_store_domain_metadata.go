package changeset

import "fmt"

type OffchainConfig struct {
	DeltaGrace   string `json:"deltaGrace"`
	DeltaInitial string `json:"deltaInitial"`
}
type DonMetadata struct {
	ID                  string `json:"id"`
	ConfiguratorAddress string `json:"configuratorAddress"`
	OffchainConfig      OffchainConfig
	Streams             []int `json:"streams"`
}

// DataStreamsMetadata is a struct that can be used as a default metadata type.
type DataStreamsMetadata struct {
	DONs []DonMetadata
}

// DefaultMetadata implements the Cloneable interface
func (d DataStreamsMetadata) Clone() DataStreamsMetadata { return d }

func (d DataStreamsMetadata) GetDonById(id string) (DonMetadata, error) {
	for _, don := range d.DONs {
		if don.ID == id {
			return don, nil
		}
	}
	return DonMetadata{}, fmt.Errorf("don with id %s not found", id)
}
