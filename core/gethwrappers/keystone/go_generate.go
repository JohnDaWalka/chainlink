// Package gethwrappers provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package gethwrappers

// Keystone

//go:generate go run ../generation/wrap.go keystone BalanceReader balance_reader

//go:generate go run ../generation/wrap.go keystone CapabilitiesRegistry capabilities_registry
//go:generate go run ../generation/generate/zk/wrap_zk_bytecode.go ../../../contracts/zksolc/keystone/CapabilitiesRegistry/CapabilitiesRegistry.sol/CapabilitiesRegistry.zbin ../keystone/generated/capabilities_registry/zk_bytecode.go capabilities_registry
//go:generate cp generated/capabilities_registry/zk_bytecode.go generated/capabilities_registry_1_1_0

//go:generate go run ../generation/wrap.go keystone KeystoneFeedsConsumer feeds_consumer
//go:generate go run ../generation/wrap.go keystone KeystoneForwarder forwarder
//go:generate go run ../generation/wrap.go keystone OCR3Capability ocr3_capability
