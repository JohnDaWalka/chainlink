package common

import (
	"crypto/ecdsa"
	"encoding/binary"
	"slices"

	"github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/hex"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func Uint32ToBytes(val uint32) []byte {
	result := make([]byte, 4)
	binary.BigEndian.PutUint32(result, val)
	return result
}

func BytesToUint32(data []byte) uint32 {
	return binary.BigEndian.Uint32(data)
}

// input string can't have any 0x0 characters
func StringToAlignedBytes(input string, size int) []byte {
	aligned := make([]byte, size)
	copy(aligned, input)
	return aligned
}

func AlignedBytesToString(data []byte) string {
	idx := slices.IndexFunc(data, func(b byte) bool { return b == 0 })
	if idx == -1 {
		return string(data)
	}
	return string(data[:idx])
}

func Flatten(data ...[]byte) []byte {
	return flatten(data...)
}
func flatten(data ...[]byte) []byte {
	var result []byte
	for _, d := range data {
		result = append(result, d...)
	}
	return result
}

func SignData(privateKey *ecdsa.PrivateKey, data ...[]byte) ([]byte, error) {
	return utils.GenerateEthSignature(privateKey, flatten(data...))
}

func ValidateMessageAndSetSigner(m *gateway.Message) (err error) {
	err = m.Validate()
	if err != nil {
		return err
	}
	rawData := gateway.GetRawMessageBody(&m.Body)
	signatureBytes, err := hex.DecodeString(m.Signature)
	if err != nil {
		return err
	}
	signer, err := ExtractSigner(signatureBytes, rawData...)
	if err != nil {
		return err
	}
	m.Body.Sender = utils.StringToHex(string(signer))
	return nil
}

func ExtractSigner(signature []byte, data ...[]byte) (signerAddress []byte, err error) {
	addr, err := utils.GetSignersEthAddress(flatten(data...), signature)
	if err != nil {
		return nil, err
	}
	return addr.Bytes(), nil
}
