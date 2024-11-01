package txm

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type DummyKeystore struct {
	privateKeyMap map[common.Address]*ecdsa.PrivateKey
}

func NewKeystore() *DummyKeystore {
	return &DummyKeystore{privateKeyMap: make(map[common.Address]*ecdsa.PrivateKey)}
}

func (k *DummyKeystore) Add(privateKeyString string, address common.Address) error {
	privateKey, err := crypto.HexToECDSA(privateKeyString)
	if err != nil {
		return err
	}
	k.privateKeyMap[address] = privateKey
	return nil
}

func (k *DummyKeystore) SignTx(_ context.Context, fromAddress common.Address, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	if key, exists := k.privateKeyMap[fromAddress]; exists {
		return types.SignTx(tx, types.LatestSignerForChainID(chainID), key)
	}
	return nil, fmt.Errorf("private key for address: %v not found", fromAddress)
}
