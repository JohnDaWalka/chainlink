package crib

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
)

type OutputReader struct {
	cribEnvStateDirPath string
}

// NewOutputReader creates new instance
func NewOutputReader(cribEnvStateDirPath string) *OutputReader {
	return &OutputReader{cribEnvStateDirPath: cribEnvStateDirPath}
}

func (r *OutputReader) ReadNodesDetails() (NodesDetails, error) {
	var result NodesDetails
	byteValue, err := r.readCRIBDataFile(NodesDetailsFileName)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		return result, errors.Wrap(err, "error unmarshalling result")
	}

	return result, nil
}

func (r *OutputReader) ReadRMNNodeConfigs() ([]RMNNodeConfig, error) {
	var result []RMNNodeConfig
	byteValue, err := r.readCRIBDataFile(RMNNodeIdentitiesFileName)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling result")
	}

	return result, nil
}

func (r *OutputReader) ReadChainConfigs() ([]devenv.ChainConfig, error) {
	var result []devenv.ChainConfig
	byteValue, err := r.readCRIBDataFile(ChainsConfigsFileName)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling result")
	}

	return result, nil
}

func (r *OutputReader) ReadAddressBook() (*deployment.AddressBookMap, error) {
	var result map[uint64]map[string]deployment.TypeAndVersion
	byteValue, err := r.readCRIBDataFile(AddressBookFileName)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling result")
	}

	return deployment.NewMemoryAddressBookFromMap(result), nil
}

func (r *OutputReader) ReadJDOutput() (*jd.Output, error) {
	var result jd.Output
	byteValue, err := r.readCRIBDataFile(JDOutputFileName)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling JSON")
	}

	return &result, nil
}

func (r *OutputReader) ReadBlockchainOutputs() ([]blockchain.Output, error) {
	var result []blockchain.Output
	byteValue, err := r.readCRIBDataFile(BlockChainsOutputFileName)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling JSON")
	}

	return result, nil
}

func (r *OutputReader) ReadNodeSetOutput() (*simple_node_set.Output, error) {
	var nodeOuts []clnode.NodeOut
	byteValue, err := r.readCRIBDataFile(NodeSetOutputFileName)
	if err != nil {
		return nil, errors.Wrap(err, "error reading node set output file")
	}

	err = json.Unmarshal(byteValue, &nodeOuts)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling JSON")
	}

	nodes := make([]*clnode.Output, len(nodeOuts))
	for i, nodeOut := range nodeOuts {
		nodes[i] = &clnode.Output{Node: &nodeOut}
	}

	return &simple_node_set.Output{
		UseCache: false,
		CLNodes:  nodes,
	}, nil
}

func (r *OutputReader) readCRIBDataFile(fileName string) ([]byte, error) {
	dataDirPath := path.Join(r.cribEnvStateDirPath, "data")
	file, err := os.Open(fmt.Sprintf("%s/%s", dataDirPath, fileName))
	if err != nil {
		return nil, errors.Wrap(err, "error opening file")
	}
	defer file.Close()

	// Read the file's content into a byte slice
	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.Wrap(err, "error reading file")
	}
	return byteValue, nil
}
