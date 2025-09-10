package crypto

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type P2PKey struct {
	EncryptedJSON []byte
	PeerID        p2pkey.PeerID
	Password      string
}

func NewP2PKey(password string) (*P2PKey, error) {
	key, err := p2pkey.NewV2()
	if err != nil {
		return nil, err
	}
	d, err := key.ToEncryptedJSON(password, utils.DefaultScryptParams)
	if err != nil {
		return nil, err
	}

	return &P2PKey{
		EncryptedJSON: d,
		PeerID:        key.PeerID(),
		Password:      password,
	}, nil
}

type P2PKeys struct {
	Keys     []P2PKey
	Password string
}

func GenerateP2PKeys(password string, n int) (*P2PKeys, error) {
	result := &P2PKeys{
		Password: password,
		Keys:     make([]P2PKey, 0, n),
	}
	for i := 0; i < n; i++ {
		key, err := NewP2PKey(password)
		if err != nil {
			return nil, err
		}

		result.Keys = append(result.Keys, P2PKey{
			EncryptedJSON: key.EncryptedJSON,
			PeerID:        key.PeerID,
			Password:      password,
		})
	}
	return result, nil
}
