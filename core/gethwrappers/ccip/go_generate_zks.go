// Package gethwrappers_ccip provides tools for wrapping solidity contracts with
// golang packages, using abigen.
package ccip

//go:generate go run ./generation/generate_zks/wrap.go ../../../contracts/solc/v0.8.24/ARMProxy/ARMProxy.abi ../../../contracts/solc/v0.8.24/ARMProxy/ARMProxy.bin RMNProxyContract rmn_proxy_contract ../../../contracts/zksolc/v1.5.6/ARMProxy/ARMProxy.sol/ARMProxy.zbin
//go:generate go run ./generation/generate_zks/wrap.go ../../../contracts/solc/v0.8.24/TokenAdminRegistry/TokenAdminRegistry.abi ../../../contracts/solc/v0.8.24/TokenAdminRegistry/TokenAdminRegistry.bin TokenAdminRegistry token_admin_registry ../../../contracts/zksolc/v1.5.6/TokenAdminRegistry/TokenAdminRegistry.sol/TokenAdminRegistry.zbin
//go:generate go run ./generation/generate_zks/wrap.go ../../../contracts/solc/v0.8.24/RegistryModuleOwnerCustom/RegistryModuleOwnerCustom.abi ../../../contracts/solc/v0.8.24/RegistryModuleOwnerCustom/RegistryModuleOwnerCustom.bin RegistryModuleOwnerCustom registry_module_owner_custom ../../../contracts/zksolc/v1.5.6/RegistryModuleOwnerCustom/RegistryModuleOwnerCustom.sol/RegistryModuleOwnerCustom.zbin
//go:generate go run ./generation/generate_zks/wrap.go ../../../contracts/solc/v0.8.24/OnRamp/OnRamp.abi ../../../contracts/solc/v0.8.24/OnRamp/OnRamp.bin OnRamp onramp ../../../contracts/zksolc/v1.5.6/OnRamp/OnRamp.sol/OnRamp.zbin
//go:generate go run ./generation/generate_zks/wrap.go ../../../contracts/solc/v0.8.24/OffRamp/OffRamp.abi ../../../contracts/solc/v0.8.24/OffRamp/OffRamp.bin OffRamp offramp ../../../contracts/zksolc/v1.5.6/OffRamp/OffRamp.sol/OffRamp.zbin
//go:generate go run ./generation/generate_zks/wrap.go ../../../contracts/solc/v0.8.24/RMNRemote/RMNRemote.abi ../../../contracts/solc/v0.8.24/RMNRemote/RMNRemote.bin RMNRemote rmn_remote ../../../contracts/zksolc/v1.5.6/RMNRemote/RMNRemote.sol/RMNRemote.zbin
//go:generate go run ./generation/generate_zks/wrap.go ../../../contracts/solc/v0.8.24/RMNHome/RMNHome.abi ../../../contracts/solc/v0.8.24/RMNHome/RMNHome.bin RMNHome rmn_home ../../../contracts/zksolc/v1.5.6/RMNHome/RMNHome.sol/RMNHome.zbin
//go:generate go run ./generation/generate_zks/wrap.go ../../../contracts/solc/v0.8.24/MultiAggregateRateLimiter/MultiAggregateRateLimiter.abi ../../../contracts/solc/v0.8.24/MultiAggregateRateLimiter/MultiAggregateRateLimiter.bin MultiAggregateRateLimiter multi_aggregate_rate_limiter ../../../contracts/zksolc/v1.5.6/MultiAggregateRateLimiter/MultiAggregateRateLimiter.sol/MultiAggregateRateLimiter.zbin
//go:generate go run ./generation/generate_zks/wrap.go ../../../contracts/solc/v0.8.24/Router/Router.abi ../../../contracts/solc/v0.8.24/Router/Router.bin Router router ../../../contracts/zksolc/v1.5.6/Router/Router.sol/Router.zbin
//go:generate go run ./generation/generate_zks/wrap.go ../../../contracts/solc/v0.8.24/FeeQuoter/FeeQuoter.abi ../../../contracts/solc/v0.8.24/FeeQuoter/FeeQuoter.bin FeeQuoter fee_quoter ../../../contracts/zksolc/v1.5.6/FeeQuoter/FeeQuoter.sol/FeeQuoter.zbin
//go:generate go run ./generation/generate_zks/wrap.go ../../../contracts/solc/v0.8.24/CCIPHome/CCIPHome.abi ../../../contracts/solc/v0.8.24/CCIPHome/CCIPHome.bin CCIPHome ccip_home ../../../contracts/zksolc/v1.5.6/CCIPHome/CCIPHome.sol/CCIPHome.zbin
//go:generate go run ./generation/generate_zks/wrap.go ../../../contracts/solc/v0.8.24/NonceManager/NonceManager.abi ../../../contracts/solc/v0.8.24/NonceManager/NonceManager.bin NonceManager nonce_manager ../../../contracts/zksolc/v1.5.6/NonceManager/NonceManager.sol/NonceManager.zbin

// Pools
//go:generate go run ./generation/generate_zks/wrap.go ../../../contracts/solc/v0.8.24/BurnMintTokenPool/BurnMintTokenPool.abi ../../../contracts/solc/v0.8.24/BurnMintTokenPool/BurnMintTokenPool.bin BurnMintTokenPool burn_mint_token_pool ../../../contracts/zksolc/v1.5.6/BurnMintTokenPool/BurnMintTokenPool.sol/BurnMintTokenPool.zbin
//go:generate go run ./generation/generate_zks/wrap.go ../../../contracts/solc/v0.8.24/BurnFromMintTokenPool/BurnFromMintTokenPool.abi ../../../contracts/solc/v0.8.24/BurnFromMintTokenPool/BurnFromMintTokenPool.bin BurnFromMintTokenPool burn_from_mint_token_pool ../../../contracts/zksolc/v1.5.6/BurnFromMintTokenPool/BurnFromMintTokenPool.sol/BurnFromMintTokenPool.zbin
//go:generate go run ./generation/generate_zks/wrap.go ../../../contracts/solc/v0.8.24/BurnWithFromMintTokenPool/BurnWithFromMintTokenPool.abi ../../../contracts/solc/v0.8.24/BurnWithFromMintTokenPool/BurnWithFromMintTokenPool.bin BurnWithFromMintTokenPool burn_with_from_mint_token_pool ../../../contracts/zksolc/v1.5.6/BurnWithFromMintTokenPool/BurnWithFromMintTokenPool.sol/BurnWithFromMintTokenPool.zbin
//go:generate go run ./generation/generate_zks/wrap.go ../../../contracts/solc/v0.8.24/LockReleaseTokenPool/LockReleaseTokenPool.abi ../../../contracts/solc/v0.8.24/LockReleaseTokenPool/LockReleaseTokenPool.bin LockReleaseTokenPool lock_release_token_pool ../../../contracts/zksolc/v1.5.6/LockReleaseTokenPool/LockReleaseTokenPool.sol/LockReleaseTokenPool.zbin
//go:generate go run ./generation/generate_zks/wrap.go ../../../contracts/solc/v0.8.24/TokenPool/TokenPool.abi ../../../contracts/solc/v0.8.24/TokenPool/TokenPool.bin TokenPool token_pool ../../../contracts/zksolc/v1.5.6/TokenPool/TokenPool.sol/TokenPool.zbin
//go:generate go run ./generation/generate_zks/wrap.go ../../../contracts/solc/v0.8.24/USDCTokenPool/USDCTokenPool.abi ../../../contracts/solc/v0.8.24/USDCTokenPool/USDCTokenPool.bin USDCTokenPool usdc_token_pool ../../../contracts/zksolc/v1.5.6/USDCTokenPool/USDCTokenPool.sol/USDCTokenPool.zbin

// Helpers
//go:generate go run ./generation/generate_zks/wrap.go ../../../contracts/solc/v0.8.24/MockV3Aggregator/MockV3Aggregator.abi ../../../contracts/solc/v0.8.24/MockV3Aggregator/MockV3Aggregator.bin MockV3Aggregator mock_v3_aggregator_contract ../../../contracts/zksolc/v1.5.6/MockV3Aggregator/MockV3Aggregator.sol/MockV3Aggregator.zbin
//go:generate go run ./generation/generate_zks/wrap.go ../../../contracts/solc/v0.8.24/MaybeRevertMessageReceiver/MaybeRevertMessageReceiver.abi ../../../contracts/solc/v0.8.24/MaybeRevertMessageReceiver/MaybeRevertMessageReceiver.bin MaybeRevertMessageReceiver maybe_revert_message_receiver ../../../contracts/zksolc/v1.5.6/MaybeRevertMessageReceiver/MaybeRevertMessageReceiver.sol/MaybeRevertMessageReceiver.zbin
//go:generate go run ./generation/generate_zks/wrap.go ../../../contracts/solc/v0.8.24/PingPongDemo/PingPongDemo.abi ../../../contracts/solc/v0.8.24/PingPongDemo/PingPongDemo.bin PingPongDemo ping_pong_demo ../../../contracts/zksolc/v1.5.6/PingPongDemo/PingPongDemo.sol/PingPongDemo.zbin
//go:generate go run ./generation/generate_zks/wrap.go ../../../contracts/solc/v0.8.24/MessageHasher/MessageHasher.abi ../../../contracts/solc/v0.8.24/MessageHasher/MessageHasher.bin MessageHasher message_hasher ../../../contracts/zksolc/v1.5.6/MessageHasher/MessageHasher.sol/MessageHasher.zbin
//go:generate go run ./generation/generate_zks/wrap.go ../../../contracts/solc/v0.8.24/MultiOCR3Helper/MultiOCR3Helper.abi ../../../contracts/solc/v0.8.24/MultiOCR3Helper/MultiOCR3Helper.bin MultiOCR3Helper multi_ocr3_helper ../../../contracts/zksolc/v1.5.6/MultiOCR3Helper/MultiOCR3Helper.sol/MultiOCR3Helper.zbin
//go:generate go run ./generation/generate_zks/wrap.go ../../../contracts/solc/v0.8.24/USDCReaderTester/USDCReaderTester.abi ../../../contracts/solc/v0.8.24/USDCReaderTester/USDCReaderTester.bin USDCReaderTester usdc_reader_tester ../../../contracts/zksolc/v1.5.6/USDCReaderTester/USDCReaderTester.sol/USDCReaderTester.zbin
//go:generate go run ./generation/generate_zks/wrap.go ../../../contracts/solc/v0.8.24/ReportCodec/ReportCodec.abi ../../../contracts/solc/v0.8.24/ReportCodec/ReportCodec.bin ReportCodec report_codec ../../../contracts/zksolc/v1.5.6/ReportCodec/ReportCodec.sol/ReportCodec.zbin
//go:generate go run ./generation/generate_zks/wrap.go ../../../contracts/solc/v0.8.24/EtherSenderReceiver/EtherSenderReceiver.abi ../../../contracts/solc/v0.8.24/EtherSenderReceiver/EtherSenderReceiver.bin EtherSenderReceiver ether_sender_receiver ../../../contracts/zksolc/v1.5.6/EtherSenderReceiver/EtherSenderReceiver.sol/EtherSenderReceiver.zbin
//go:generate go run ./generation/generate_zks/wrap.go ../../../contracts/solc/v0.8.24/WETH9/WETH9.abi ../../../contracts/solc/v0.8.24/WETH9/WETH9.bin WETH9 weth9 ../../../contracts/zksolc/v1.5.6/WETH9/WETH9.sol/WETH9.zbin
//go:generate go run ./generation/generate_zks/wrap.go ../../../contracts/solc/v0.8.24/MockE2EUSDCTokenMessenger/MockE2EUSDCTokenMessenger.abi ../../../contracts/solc/v0.8.24/MockE2EUSDCTokenMessenger/MockE2EUSDCTokenMessenger.bin MockE2EUSDCTokenMessenger mock_usdc_token_messenger ../../../contracts/zksolc/v1.5.6/MockE2EUSDCTokenMessenger/MockE2EUSDCTokenMessenger.sol/MockE2EUSDCTokenMessenger.zbin
//go:generate go run ./generation/generate_zks/wrap.go ../../../contracts/solc/v0.8.24/MockE2EUSDCTransmitter/MockE2EUSDCTransmitter.abi ../../../contracts/solc/v0.8.24/MockE2EUSDCTransmitter/MockE2EUSDCTransmitter.bin MockE2EUSDCTransmitter mock_usdc_token_transmitter ../../../contracts/zksolc/v1.5.6/MockE2EUSDCTransmitter/MockE2EUSDCTransmitter.sol/MockE2EUSDCTransmitter.zbin
//go:generate go run ./generation/generate_zks/wrap.go ../../../contracts/solc/v0.8.24/CCIPReaderTester/CCIPReaderTester.abi ../../../contracts/solc/v0.8.24/CCIPReaderTester/CCIPReaderTester.bin CCIPReaderTester ccip_reader_tester ../../../contracts/zksolc/v1.5.6/CCIPReaderTester/CCIPReaderTester.sol/CCIPReaderTester.zbin

// EncodingUtils
//go:generate go run ./generation/generate_zks/wrap.go ../../../contracts/solc/v0.8.24/ICCIPEncodingUtils/ICCIPEncodingUtils.abi ../../../contracts/solc/v0.8.24/ICCIPEncodingUtils/ICCIPEncodingUtils.bin EncodingUtils ccip_encoding_utils ../../../contracts/zksolc/v1.5.6/ICCIPEncodingUtils/ICCIPEncodingUtils.sol/ICCIPEncodingUtils.zbin

// To run these commands, you must either install docker, or the correct version
// of abigen. The latter can be installed with these commands, at least on linux:
//
//   git clone https://github.com/ethereum/go-ethereum
//   cd go-ethereum/cmd/abigen
//   git checkout v<version-needed>
//   go install
//
// Here, <version-needed> is the version of go-ethereum specified in chainlink's
// go.mod. This will install abigen in "$GOPATH/bin", which you should add to
// your $PATH.
//
// To reduce explicit dependencies, and in case the system does not have the
// correct version of abigen installed , the above commands spin up docker
// containers. In my hands, total running time including compilation is about
// 13s. If you're modifying solidity code and testing against go code a lot, it
// might be worthwhile to generate the the wrappers using a static container
// with abigen and solc, which will complete much faster. E.g.
//
//   abigen -sol ../../contracts/src/v0.6/VRFAll.sol -pkg vrf -out solidity_interfaces.go
//
// where VRFAll.sol simply contains `import "contract_path";` instructions for
// all the contracts you wish to target. This runs in about 0.25 seconds in my
// hands.
//
// If you're on linux, you can copy the correct version of solc out of the
// appropriate docker container. At least, the following works on ubuntu:
//
//   $ docker run --name solc ethereum/solc:0.6.2
//   $ sudo docker cp solc:/usr/bin/solc /usr/bin
//   $ docker rm solc
//
// If you need to point abigen at your solc executable, you can specify the path
// with the abigen --solc <path-to-executable> option.
