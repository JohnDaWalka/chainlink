package ccip

//go:generate go run ../generation/wrap.go ccip Router router latest
//go:generate go run ../generation/generate/zk/wrap_zk_bytecode.go ../../../contracts/zksolc/ccip/Router/Router.sol/Router.zbin ../ccip/generated/latest/router/zk_bytecode.go router

//go:generate go run ../generation/wrap.go ccip CCIPHome ccip_home latest
//go:generate go run ../generation/generate/zk/wrap_zk_bytecode.go ../../../contracts/zksolc/ccip/capability/CCIPHome/CCIPHome.sol/CCIPHome.zbin ../ccip/generated/latest/ccip_home/zk_bytecode.go ccip_home
//go:generate cp generated/latest/ccip_home/zk_bytecode.go generated/v1_6_0/ccip_home

//go:generate go run ../generation/wrap.go ccip OnRamp onramp latest
//go:generate go run ../generation/generate/zk/wrap_zk_bytecode.go ../../../contracts/zksolc/ccip/onRamp/OnRamp/OnRamp.sol/OnRamp.zbin ../ccip/generated/latest/onramp/zk_bytecode.go onramp

//go:generate go run ../generation/wrap.go ccip OffRamp offramp latest
//go:generate go run ../generation/generate/zk/wrap_zk_bytecode.go ../../../contracts/zksolc/ccip/offRamp/OffRamp/OffRamp.sol/OffRamp.zbin ../ccip/generated/latest/offramp/zk_bytecode.go offramp

//go:generate go run ../generation/wrap.go ccip OnRampWithMessageTransformer onramp_with_message_transformer latest
//go:generate go run ../generation/generate/zk/wrap_zk_bytecode.go ../../../contracts/zksolc/ccip/onRamp/OnRampWithMessageTransformer/OnRampWithMessageTransformer.sol/OnRampWithMessageTransformer.zbin ../ccip/generated/latest/onramp_with_message_transformer/zk_bytecode.go onramp_with_message_transformer

//go:generate go run ../generation/wrap.go ccip OffRampWithMessageTransformer offramp_with_message_transformer latest
//go:generate go run ../generation/generate/zk/wrap_zk_bytecode.go ../../../contracts/zksolc/ccip/offRamp/OffRampWithMessageTransformer/OffRampWithMessageTransformer.sol/OffRampWithMessageTransformer.zbin ../ccip/generated/latest/offramp_with_message_transformer/zk_bytecode.go offramp_with_message_transformer

//go:generate go run ../generation/wrap.go ccip FeeQuoter fee_quoter latest
//go:generate go run ../generation/generate/zk/wrap_zk_bytecode.go ../../../contracts/zksolc/ccip/FeeQuoter/FeeQuoter.sol/FeeQuoter.zbin ../ccip/generated/latest/fee_quoter/zk_bytecode.go fee_quoter

//go:generate go run ../generation/wrap.go ccip NonceManager nonce_manager latest
//go:generate go run ../generation/generate/zk/wrap_zk_bytecode.go ../../../contracts/zksolc/ccip/NonceManager/NonceManager.sol/NonceManager.zbin ../ccip/generated/latest/nonce_manager/zk_bytecode.go nonce_manager

//go:generate go run ../generation/wrap.go ccip MultiAggregateRateLimiter multi_aggregate_rate_limiter latest
//go:generate go run ../generation/generate/zk/wrap_zk_bytecode.go ../../../contracts/zksolc/ccip/MultiAggregateRateLimiter/MultiAggregateRateLimiter.sol/MultiAggregateRateLimiter.zbin ../ccip/generated/latest/multi_aggregate_rate_limiter/zk_bytecode.go multi_aggregate_rate_limiter

//go:generate go run ../generation/wrap.go ccip TokenAdminRegistry token_admin_registry latest
//go:generate go run ../generation/generate/zk/wrap_zk_bytecode.go ../../../contracts/zksolc/ccip/tokenAdminRegistry/TokenAdminRegistry/TokenAdminRegistry.sol/TokenAdminRegistry.zbin ../ccip/generated/latest/token_admin_registry/zk_bytecode.go token_admin_registry

//go:generate go run ../generation/wrap.go ccip RegistryModuleOwnerCustom registry_module_owner_custom latest
//go:generate go run ../generation/generate/zk/wrap_zk_bytecode.go ../../../contracts/zksolc/ccip/tokenAdminRegistry/RegistryModuleOwnerCustom/RegistryModuleOwnerCustom.sol/RegistryModuleOwnerCustom.zbin ../ccip/generated/latest/registry_module_owner_custom/zk_bytecode.go registry_module_owner_custom

//go:generate go run ../generation/wrap.go ccip RMNProxy rmn_proxy_contract latest
//go:generate go run ../generation/generate/zk/wrap_zk_bytecode.go ../../../contracts/zksolc/ccip/rmn/RMNProxy/RMNProxy.sol/RMNProxy.zbin ../ccip/generated/latest/rmn_proxy_contract/zk_bytecode.go rmn_proxy_contract

//go:generate go run ../generation/wrap.go ccip RMNRemote rmn_remote latest
//go:generate go run ../generation/generate/zk/wrap_zk_bytecode.go ../../../contracts/zksolc/ccip/rmn/RMNRemote/RMNRemote.sol/RMNRemote.zbin ../ccip/generated/latest/rmn_remote/zk_bytecode.go rmn_remote

//go:generate go run ../generation/wrap.go ccip RMNHome rmn_home latest
//go:generate go run ../generation/generate/zk/wrap_zk_bytecode.go ../../../contracts/zksolc/ccip/rmn/RMNHome/RMNHome.sol/RMNHome.zbin ../ccip/generated/latest/rmn_home/zk_bytecode.go rmn_home
//go:generate cp generated/latest/rmn_home/zk_bytecode.go generated/v1_6_0/rmn_home

// Pools
//go:generate go run ../generation/wrap.go ccip BurnMintTokenPool burn_mint_token_pool latest
//go:generate go run ../generation/generate/zk/wrap_zk_bytecode.go ../../../contracts/zksolc/ccip/pools/BurnMintTokenPool/BurnMintTokenPool.sol/BurnMintTokenPool.zbin ../ccip/generated/latest/burn_mint_token_pool/zk_bytecode.go burn_mint_token_pool

//go:generate go run ../generation/wrap.go ccip BurnFromMintTokenPool burn_from_mint_token_pool latest
//go:generate go run ../generation/generate/zk/wrap_zk_bytecode.go ../../../contracts/zksolc/ccip/pools/BurnFromMintTokenPool/BurnFromMintTokenPool.sol/BurnFromMintTokenPool.zbin ../ccip/generated/latest/burn_from_mint_token_pool/zk_bytecode.go burn_from_mint_token_pool

//go:generate go run ../generation/wrap.go ccip BurnWithFromMintTokenPool burn_with_from_mint_token_pool latest
//go:generate go run ../generation/generate/zk/wrap_zk_bytecode.go ../../../contracts/zksolc/ccip/pools/BurnWithFromMintTokenPool/BurnWithFromMintTokenPool.sol/BurnWithFromMintTokenPool.zbin ../ccip/generated/latest/burn_with_from_mint_token_pool/zk_bytecode.go burn_with_from_mint_token_pool

//go:generate go run ../generation/wrap.go ccip LockReleaseTokenPool lock_release_token_pool latest
//go:generate go run ../generation/generate/zk/wrap_zk_bytecode.go ../../../contracts/zksolc/ccip/pools/LockReleaseTokenPool/LockReleaseTokenPool.sol/LockReleaseTokenPool.zbin ../ccip/generated/latest/lock_release_token_pool/zk_bytecode.go lock_release_token_pool

//go:generate go run ../generation/wrap.go ccip TokenPool token_pool latest

//go:generate go run ../generation/wrap.go ccip USDCTokenPool usdc_token_pool latest
//go:generate go run ../generation/generate/zk/wrap_zk_bytecode.go ../../../contracts/zksolc/ccip/pools/USDC/USDCTokenPool/USDCTokenPool.sol/USDCTokenPool.zbin ../ccip/generated/latest/usdc_token_pool/zk_bytecode.go usdc_token_pool

//go:generate go run ../generation/wrap.go ccip SiloedLockReleaseTokenPool siloed_lock_release_token_pool latest
//go:generate go run ../generation/generate/zk/wrap_zk_bytecode.go ../../../contracts/zksolc/ccip/pools/SiloedLockReleaseTokenPool/SiloedLockReleaseTokenPool.sol/SiloedLockReleaseTokenPool.zbin ../ccip/generated/latest/siloed_lock_release_token_pool/zk_bytecode.go siloed_lock_release_token_pool

//go:generate go run ../generation/wrap.go ccip BurnToAddressMintTokenPool burn_to_address_mint_token_pool latest
//go:generate go run ../generation/generate/zk/wrap_zk_bytecode.go ../../../contracts/zksolc/ccip/pools/BurnToAddressMintTokenPool/BurnToAddressMintTokenPool.sol/BurnToAddressMintTokenPool.zbin ../ccip/generated/latest/burn_to_address_mint_token_pool/zk_bytecode.go burn_to_address_mint_token_pool

// Helpers
//go:generate go run ../generation/wrap.go ccip MaybeRevertMessageReceiver maybe_revert_message_receiver latest
//go:generate go run ../generation/generate/zk/wrap_zk_bytecode.go ../../../contracts/zksolc/ccip/test/helpers/receivers/MaybeRevertMessageReceiver/MaybeRevertMessageReceiver.sol/MaybeRevertMessageReceiver.zbin ../ccip/generated/latest/maybe_revert_message_receiver/zk_bytecode.go maybe_revert_message_receiver

//go:generate go run ../generation/wrap.go ccip LogMessageDataReceiver log_message_data_receiver latest
//go:generate go run ../generation/generate/zk/wrap_zk_bytecode.go ../../../contracts/zksolc/ccip/test/helpers/receivers/LogMessageDataReceiver/LogMessageDataReceiver.sol/LogMessageDataReceiver.zbin ../ccip/generated/latest/log_message_data_receiver/zk_bytecode.go log_message_data_receiver

//go:generate go run ../generation/wrap.go ccip PingPongDemo ping_pong_demo latest
//go:generate go run ../generation/generate/zk/wrap_zk_bytecode.go ../../../contracts/zksolc/ccip/applications/PingPongDemo/PingPongDemo.sol/PingPongDemo.zbin ../ccip/generated/latest/ping_pong_demo/zk_bytecode.go ping_pong_demo

//go:generate go run ../generation/wrap.go ccip MessageHasher message_hasher latest
//go:generate go run ../generation/generate/zk/wrap_zk_bytecode.go ../../../contracts/zksolc/ccip/test/helpers/MessageHasher/MessageHasher.sol/MessageHasher.zbin ../ccip/generated/latest/message_hasher/zk_bytecode.go message_hasher

//go:generate go run ../generation/wrap.go ccip MultiOCR3Helper multi_ocr3_helper latest
//go:generate go run ../generation/generate/zk/wrap_zk_bytecode.go ../../../contracts/zksolc/ccip/test/helpers/MultiOCR3Helper/MultiOCR3Helper.sol/MultiOCR3Helper.zbin ../ccip/generated/latest/multi_ocr3_helper/zk_bytecode.go multi_ocr3_helper

//go:generate go run ../generation/wrap.go ccip USDCReaderTester usdc_reader_tester latest
//go:generate go run ../generation/generate/zk/wrap_zk_bytecode.go ../../../contracts/zksolc/ccip/test/helpers/USDCReaderTester/USDCReaderTester.sol/USDCReaderTester.zbin ../ccip/generated/latest/usdc_reader_tester/zk_bytecode.go usdc_reader_tester

//go:generate go run ../generation/wrap.go ccip ReportCodec report_codec latest
//go:generate go run ../generation/generate/zk/wrap_zk_bytecode.go ../../../contracts/zksolc/ccip/test/helpers/ReportCodec/ReportCodec.sol/ReportCodec.zbin ../ccip/generated/latest/report_codec/zk_bytecode.go report_codec

//go:generate go run ../generation/wrap.go ccip EtherSenderReceiver ether_sender_receiver latest
//go:generate go run ../generation/generate/zk/wrap_zk_bytecode.go ../../../contracts/zksolc/ccip/applications/EtherSenderReceiver/EtherSenderReceiver.sol/EtherSenderReceiver.zbin ../ccip/generated/latest/ether_sender_receiver/zk_bytecode.go ether_sender_receiver

//go:generate go run ../generation/wrap.go ccip MockE2EUSDCTokenMessenger mock_usdc_token_messenger latest
//go:generate go run ../generation/generate/zk/wrap_zk_bytecode.go ../../../contracts/zksolc/ccip/test/mocks/MockE2EUSDCTokenMessenger/MockE2EUSDCTokenMessenger.sol/MockE2EUSDCTokenMessenger.zbin ../ccip/generated/latest/mock_usdc_token_messenger/zk_bytecode.go mock_usdc_token_messenger

//go:generate go run ../generation/wrap.go ccip MockE2EUSDCTransmitter mock_usdc_token_transmitter latest
//go:generate go run ../generation/generate/zk/wrap_zk_bytecode.go ../../../contracts/zksolc/ccip/test/mocks/MockE2EUSDCTransmitter/MockE2EUSDCTransmitter.sol/MockE2EUSDCTransmitter.zbin ../ccip/generated/latest/mock_usdc_token_transmitter/zk_bytecode.go mock_usdc_token_transmitter

//go:generate go run ../generation/wrap.go ccip CCIPReaderTester ccip_reader_tester latest
//go:generate go run ../generation/generate/zk/wrap_zk_bytecode.go ../../../contracts/zksolc/ccip/test/helpers/CCIPReaderTester/CCIPReaderTester.sol/CCIPReaderTester.zbin ../ccip/generated/latest/ccip_reader_tester/zk_bytecode.go ccip_reader_tester

// EncodingUtils
//go:generate go run ../generation/wrap.go ccip EncodingUtils ccip_encoding_utils latest
