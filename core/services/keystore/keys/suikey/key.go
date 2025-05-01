package suikey

import (
	"crypto"
	"crypto/ed25519"
	cryptorand "crypto/rand"
	"fmt"
	"io"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/internal"
)

// Key represents a Sui account
type Key struct {
	raw    internal.Raw
	signFn func(io.Reader, []byte, crypto.SignerOpts) ([]byte, error)
	pubKey ed25519.PublicKey
}

// New creates new Key
func New() (Key, error) {
	return newFrom(cryptorand.Reader)
}

// MustNewInsecure returns an Account if no error
func MustNewInsecure(reader io.Reader) Key {
	key, err := newFrom(reader)
	if err != nil {
		panic(err)
	}
	return key
}

// newFrom creates a new Account from a provided random reader
func newFrom(reader io.Reader) (Key, error) {
	pub, priv, err := ed25519.GenerateKey(reader)
	if err != nil {
		return Key{}, err
	}
	return Key{
		raw:    internal.NewRaw(priv.Seed()),
		signFn: priv.Sign,
		pubKey: pub,
	}, nil
}

// ID gets Account ID
func (s Key) ID() string {
	return s.PublicKeyStr()
}

// Address returns the Sui address
func (s Key) Address() string {
	return fmt.Sprintf("%064x", s.pubKey)
}

// GetPublic gets Account's public key
func (s Key) GetPublic() ed25519.PublicKey {
	return s.pubKey
}

// PublicKeyStr returns hex encoded public key
func (s Key) PublicKeyStr() string {
	return fmt.Sprintf("%064x", s.pubKey)
}

// Raw returns the seed from private key
func (s Key) Raw() internal.Raw { return s.raw }

// Sign is used to sign a message
func (s Key) Sign(msg []byte) ([]byte, error) {
	return s.signFn(cryptorand.Reader, msg, crypto.Hash(0)) // no specific hash function used
}

// KeyFor creates an Account from a raw key
func KeyFor(raw internal.Raw) Key {
	privKey := ed25519.NewKeyFromSeed(internal.Bytes(raw))
	pubKey := privKey.Public().(ed25519.PublicKey)
	return Key{
		raw:    raw,
		signFn: privKey.Sign,
		pubKey: pubKey,
	}
}

// String returns the public key as a hex string
func (s Key) String() string {
	return s.ID()
}

// GoString returns the public key as a hex string
func (s Key) GoString() string {
	return s.String()
}

// ToV2 returns the key as a Raw
func (s Key) ToV2() internal.Raw {
	return s.Raw()
}
