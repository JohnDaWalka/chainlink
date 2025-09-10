package crypto

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
)

type EVMKeys struct {
	EncryptedJSONs  [][]byte
	PublicAddresses []common.Address
	Password        string
	ChainID         int
}

type EVMKey struct {
	EncryptedJSON []byte
	PublicAddress common.Address
	Password      string
	ChainID       int
}

func NewEVMKey(password string, chainID int) (*EVMKey, error) {
	key, addr, err := clclient.NewETHKey(password)
	if err != nil {
		return nil, fmt.Errorf("failed to create new EVM key: %w", err)
	}
	return &EVMKey{
		EncryptedJSON: key,
		PublicAddress: addr,
		Password:      password,
		ChainID:       chainID,
	}, nil
}

func GenerateEVMKeys(password string, n int) (*EVMKeys, error) {
	result := &EVMKeys{
		Password: password,
	}
	for range n {
		key, addr, err := clclient.NewETHKey(password)
		if err != nil {
			return result, fmt.Errorf("failed to create new EVM key: %w", err)
		}
		result.EncryptedJSONs = append(result.EncryptedJSONs, key)
		result.PublicAddresses = append(result.PublicAddresses, addr)
	}
	return result, nil
}

/*
Generates new private and public key pair

Returns a new public address and a private key
*/
func GenerateNewKeyPair() (common.Address, *ecdsa.PrivateKey, error) {
	privateKey, pkErr := crypto.GenerateKey()
	if pkErr != nil {
		return common.Address{}, nil, errors.Wrap(pkErr, "failed to generate a new private key (EOA)")
	}

	publicKeyAddr := crypto.PubkeyToAddress(privateKey.PublicKey)
	return publicKeyAddr, privateKey, nil
}

func PrivateKeyToAddress(privateKey *ecdsa.PrivateKey) (common.Address, error) {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, errors.New("error casting public key to ECDSA")
	}
	return crypto.PubkeyToAddress(*publicKeyECDSA), nil
}
