package main

// sample generation
// run it from gethwrappers with go generate ./generation/generate/zk/go_generate.go

//go:generate go run wrap_zk_bytecode.go ../../../../../contracts/zksolc/shared/token/ERC677/LinkToken/LinkToken.sol/LinkToken.zbin ../../../shared/generated/link_token/zk_bytecode.go link_token LinkToken
