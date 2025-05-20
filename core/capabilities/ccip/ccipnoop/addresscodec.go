package ccipnoop

type NoopAddressCodec struct{}

func (n NoopAddressCodec) AddressBytesToString(addr []byte) (string, error) {
	return string(addr), nil
}

func (n NoopAddressCodec) AddressStringToBytes(addr string) ([]byte, error) {
	return []byte(addr), nil
}
