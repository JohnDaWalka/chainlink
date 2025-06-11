package deployment

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"
	"github.com/zksync-sdk/zksync2-go/types"
)

type ZkSyncECDSASigner struct {
	kmsClient   *EVMKMSClient
	address     common.Address
	pubKeyBytes []byte
	chainID     *big.Int
}

func (s *ZkSyncECDSASigner) Address() common.Address {
	return s.address
}

func (s *ZkSyncECDSASigner) ChainID() *big.Int {
	return s.chainID
}

func (s *ZkSyncECDSASigner) PrivateKey() *ecdsa.PrivateKey {
	return nil
}

func (s *ZkSyncECDSASigner) SignMessage(_ context.Context, msg []byte) ([]byte, error) {
	return nil, errors.New("hashTypedData not implemented for ZkSyncECDSASigner")
}

func (s *ZkSyncECDSASigner) SignTransaction(ctx context.Context, tx *types.Transaction) ([]byte, error) {
	return nil, errors.New("hashTypedData not implemented for ZkSyncECDSASigner")
}

func (s *ZkSyncECDSASigner) SignTypedData(_ context.Context, typedData *apitypes.TypedData) ([]byte, error) {
	hash, err := accounts.HashTypedData(typedData)
	if err != nil {
		return nil, fmt.Errorf("failed to get hash of typed data: %w", err)
	}

	sig, err := s.kmsClient.SignHash(hash, s.pubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to sign hash of typed data: %w", err)
	}
	if sig[64] < 27 {
		sig[64] += 27
	}
	return sig, nil
}

func CreateZkSyncWalletWithKMS(kmsClient *EVMKMSClient, address common.Address, chainID *big.Int, pubKeyBytes []byte, client *clients.Client) (*accounts.Wallet, error) {
	signer := &ZkSyncECDSASigner{
		kmsClient:   kmsClient,
		pubKeyBytes: pubKeyBytes,
		address:     address,
		chainID:     chainID,
	}
	return accounts.NewWalletFromSigner(signer, client, nil)
}
