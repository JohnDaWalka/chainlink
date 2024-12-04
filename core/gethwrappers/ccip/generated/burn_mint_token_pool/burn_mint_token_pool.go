package burn_mint_token_pool

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

type PoolLockOrBurnInV1 struct {
	Receiver            []byte
	RemoteChainSelector uint64
	OriginalSender      common.Address
	Amount              *big.Int
	LocalToken          common.Address
}

type PoolLockOrBurnOutV1 struct {
	DestTokenAddress []byte
	DestPoolData     []byte
}

type PoolReleaseOrMintInV1 struct {
	OriginalSender      []byte
	RemoteChainSelector uint64
	Receiver            common.Address
	Amount              *big.Int
	LocalToken          common.Address
	SourcePoolAddress   []byte
	SourcePoolData      []byte
	OffchainTokenData   []byte
}

type PoolReleaseOrMintOutV1 struct {
	DestinationAmount *big.Int
}

type RateLimiterConfig struct {
	IsEnabled bool
	Capacity  *big.Int
	Rate      *big.Int
}

type RateLimiterTokenBucket struct {
	Tokens      *big.Int
	LastUpdated uint32
	IsEnabled   bool
	Capacity    *big.Int
	Rate        *big.Int
}

type TokenPoolChainUpdate struct {
	RemoteChainSelector       uint64
	RemotePoolAddresses       [][]byte
	RemoteTokenAddress        []byte
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
}

var BurnMintTokenPoolMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIBurnMintERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"localTokenDecimals\",\"type\":\"uint8\"},{\"internalType\":\"address[]\",\"name\":\"allowlist\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"}],\"name\":\"AggregateValueMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"}],\"name\":\"AggregateValueRateLimitReached\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AllowListNotEnabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BucketOverfilled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"CallerIsNotARampOnRouter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotTransferToSelf\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CursedByRMN\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"DisabledNonZeroRateLimit\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"expected\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"actual\",\"type\":\"uint8\"}],\"name\":\"InvalidDecimalArgs\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"InvalidRateLimitRate\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolData\",\"type\":\"bytes\"}],\"name\":\"InvalidRemoteChainDecimals\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidRemotePoolForChain\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidSourcePoolAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"InvalidToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeProposedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"NonExistentChain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"remoteDecimals\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"localDecimals\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"remoteAmount\",\"type\":\"uint256\"}],\"name\":\"OverflowDetected\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OwnerCannotBeZero\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"PoolAlreadyAdded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RateLimitMustBeDisabled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SenderNotAllowed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenRateLimitReached\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListAdd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListRemove\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Burned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remoteToken\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"ChainAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"ChainConfigured\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Locked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Minted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"rateLimitAdmin\",\"type\":\"address\"}],\"name\":\"RateLimitAdminSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Released\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"RemotePoolAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"RemotePoolRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldRouter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"RouterUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"TokensConsumed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"addRemotePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"removes\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"adds\",\"type\":\"address[]\"}],\"name\":\"applyAllowListUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64[]\",\"name\":\"remoteChainSelectorsToRemove\",\"type\":\"uint64[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes[]\",\"name\":\"remotePoolAddresses\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"remoteTokenAddress\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"internalType\":\"structTokenPool.ChainUpdate[]\",\"name\":\"chainsToAdd\",\"type\":\"tuple[]\"}],\"name\":\"applyChainUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowList\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowListEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getCurrentInboundRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getCurrentOutboundRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRateLimitAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getRemotePools\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"\",\"type\":\"bytes[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getRemoteToken\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRmnProxy\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRouter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSupportedChains\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTokenDecimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"decimals\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"isRemotePool\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"isSupportedChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"isSupportedToken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"originalSender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"}],\"internalType\":\"structPool.LockOrBurnInV1\",\"name\":\"lockOrBurnIn\",\"type\":\"tuple\"}],\"name\":\"lockOrBurn\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"destPoolData\",\"type\":\"bytes\"}],\"internalType\":\"structPool.LockOrBurnOutV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"originalSender\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"sourcePoolData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainTokenData\",\"type\":\"bytes\"}],\"internalType\":\"structPool.ReleaseOrMintInV1\",\"name\":\"releaseOrMintIn\",\"type\":\"tuple\"}],\"name\":\"releaseOrMint\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"destinationAmount\",\"type\":\"uint256\"}],\"internalType\":\"structPool.ReleaseOrMintOutV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"removeRemotePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundConfig\",\"type\":\"tuple\"}],\"name\":\"setChainRateLimiterConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"rateLimitAdmin\",\"type\":\"address\"}],\"name\":\"setRateLimitAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"setRouter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6101006040523480156200001257600080fd5b5060405162004809380380620048098339810160408190526200003591620005a2565b8484848484336000816200005c57604051639b15e16f60e01b815260040160405180910390fd5b600180546001600160a01b0319166001600160a01b03848116919091179091558116156200008f576200008f81620001eb565b50506001600160a01b0385161580620000af57506001600160a01b038116155b80620000c257506001600160a01b038216155b15620000e1576040516342bcdf7f60e11b815260040160405180910390fd5b6001600160a01b03808616608081905290831660c0526040805163313ce56760e01b8152905163313ce567916004808201926020929091908290030181865afa92505050801562000151575060408051601f3d908101601f191682019092526200014e91810190620006c4565b60015b1562000191578060ff168560ff16146200018f576040516332ad3e0760e11b815260ff80871660048301528216602482015260440160405180910390fd5b505b60ff841660a052600480546001600160a01b0319166001600160a01b038316179055825115801560e052620001db57604080516000815260208101909152620001db908462000265565b5050505050505050505062000730565b336001600160a01b038216036200021557604051636d6c4ee560e11b815260040160405180910390fd5b600080546001600160a01b0319166001600160a01b03838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60e05162000286576040516335f4a7b360e01b815260040160405180910390fd5b60005b825181101562000311576000838281518110620002aa57620002aa620006e2565b60209081029190910101519050620002c4600282620003c2565b1562000307576040516001600160a01b03821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b5060010162000289565b5060005b8151811015620003bd576000828281518110620003365762000336620006e2565b6020026020010151905060006001600160a01b0316816001600160a01b031603620003625750620003b4565b6200036f600282620003e2565b15620003b2576040516001600160a01b03821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b60010162000315565b505050565b6000620003d9836001600160a01b038416620003f9565b90505b92915050565b6000620003d9836001600160a01b038416620004fd565b60008181526001830160205260408120548015620004f257600062000420600183620006f8565b85549091506000906200043690600190620006f8565b9050808214620004a25760008660000182815481106200045a576200045a620006e2565b9060005260206000200154905080876000018481548110620004805762000480620006e2565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080620004b657620004b66200071a565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050620003dc565b6000915050620003dc565b60008181526001830160205260408120546200054657508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155620003dc565b506000620003dc565b6001600160a01b03811681146200056557600080fd5b50565b805160ff811681146200057a57600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b80516200057a816200054f565b600080600080600060a08688031215620005bb57600080fd5b8551620005c8816200054f565b94506020620005d987820162000568565b60408801519095506001600160401b0380821115620005f757600080fd5b818901915089601f8301126200060c57600080fd5b8151818111156200062157620006216200057f565b8060051b604051601f19603f830116810181811085821117156200064957620006496200057f565b60405291825284820192508381018501918c8311156200066857600080fd5b938501935b828510156200069157620006818562000595565b845293850193928501926200066d565b809850505050505050620006a86060870162000595565b9150620006b86080870162000595565b90509295509295909350565b600060208284031215620006d757600080fd5b620003d98262000568565b634e487b7160e01b600052603260045260246000fd5b81810381811115620003dc57634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052603160045260246000fd5b60805160a05160c05160e051614028620007e16000396000818161054f01528181611d8201526127cd0152600081816105290152818161189f015261206e0152600081816102e001528181610ba901528181611a4801528181611b0201528181611b3601528181611b6901528181611bce01528181611c270152611cc90152600081816102470152818161029c01528181610708015281816121eb0152818161276301526129b801526140286000f3fe608060405234801561001057600080fd5b50600436106101cf5760003560e01c80639a4575b911610104578063c0d78655116100a2578063dc0bd97111610071578063dc0bd97114610527578063e0351e131461054d578063e8a1da1714610573578063f2fde38b1461058657600080fd5b8063c0d78655146104d9578063c4bffe2b146104ec578063c75eea9c14610501578063cf7401f31461051457600080fd5b8063acfecf91116100de578063acfecf9114610426578063af58d59f14610439578063b0f479a1146104a8578063b7946580146104c657600080fd5b80639a4575b9146103d1578063a42a7b8b146103f1578063a7cd63b71461041157600080fd5b806354c8a4f31161017157806379ba50971161014b57806379ba5097146103855780637d54534e1461038d5780638926f54f146103a05780638da5cb5b146103b357600080fd5b806354c8a4f31461033f57806362ddd3c4146103545780636d3d1a581461036757600080fd5b8063240028e8116101ad578063240028e81461028c57806324f65ee7146102d9578063390775371461030a5780634c5ef0ed1461032c57600080fd5b806301ffc9a7146101d4578063181f5a77146101fc57806321df0da714610245575b600080fd5b6101e76101e2366004613178565b610599565b60405190151581526020015b60405180910390f35b6102386040518060400160405280601781526020017f4275726e4d696e74546f6b656e506f6f6c20312e352e3100000000000000000081525081565b6040516101f3919061321e565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101f3565b6101e761029a366004613253565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff90811691161490565b60405160ff7f00000000000000000000000000000000000000000000000000000000000000001681526020016101f3565b61031d610318366004613270565b61067e565b604051905181526020016101f3565b6101e761033a3660046132c9565b61084d565b61035261034d366004613398565b610897565b005b6103526103623660046132c9565b610912565b60095473ffffffffffffffffffffffffffffffffffffffff16610267565b6103526109af565b61035261039b366004613253565b610a7d565b6101e76103ae366004613404565b610afe565b60015473ffffffffffffffffffffffffffffffffffffffff16610267565b6103e46103df36600461341f565b610b15565b6040516101f3919061345a565b6104046103ff366004613404565b610bee565b6040516101f391906134b1565b610419610d59565b6040516101f39190613533565b6103526104343660046132c9565b610d6a565b61044c610447366004613404565b610e82565b6040516101f3919081516fffffffffffffffffffffffffffffffff908116825260208084015163ffffffff1690830152604080840151151590830152606080840151821690830152608092830151169181019190915260a00190565b60045473ffffffffffffffffffffffffffffffffffffffff16610267565b6102386104d4366004613404565b610f57565b6103526104e7366004613253565b611007565b6104f46110e2565b6040516101f3919061358d565b61044c61050f366004613404565b61119a565b610352610522366004613715565b61126c565b7f0000000000000000000000000000000000000000000000000000000000000000610267565b7f00000000000000000000000000000000000000000000000000000000000000006101e7565b610352610581366004613398565b6112f0565b610352610594366004613253565b611802565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167faff2afbf00000000000000000000000000000000000000000000000000000000148061062c57507fffffffff0000000000000000000000000000000000000000000000000000000082167f0e64dd2900000000000000000000000000000000000000000000000000000000145b8061067857507fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a700000000000000000000000000000000000000000000000000000000145b92915050565b60408051602081019091526000815261069682611816565b60006106ef60608401356106ea6106b060c087018761375a565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611a3a92505050565b611afe565b905073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000166340c10f1961073d6060860160408701613253565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b16815273ffffffffffffffffffffffffffffffffffffffff909116600482015260248101849052604401600060405180830381600087803b1580156107aa57600080fd5b505af11580156107be573d6000803e3d6000fd5b506107d3925050506060840160408501613253565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f9d228d69b5fdb8d273a2336f8fb8612d039631024ea9bf09c424a9503aa078f08360405161083191815260200190565b60405180910390a3604080516020810190915290815292915050565b600061088f83836040516108629291906137bf565b604080519182900390912067ffffffffffffffff8716600090815260076020529190912060050190611d12565b949350505050565b61089f611d2d565b61090c84848080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525050604080516020808802828101820190935287825290935087925086918291850190849080828437600092019190915250611d8092505050565b50505050565b61091a611d2d565b61092383610afe565b61096a576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff841660048201526024015b60405180910390fd5b6109aa8383838080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611f3692505050565b505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610a00576040517f02b543c600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000008082163390811790935560008054909116815560405173ffffffffffffffffffffffffffffffffffffffff909216929183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610a85611d2d565b600980547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040519081527f44676b5284b809a22248eba0da87391d79098be38bb03154be88a58bf4d091749060200160405180910390a150565b6000610678600567ffffffffffffffff8416611d12565b6040805180820190915260608082526020820152610b3282612030565b610b3f82606001356121bc565b6040516060830135815233907f696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df79060200160405180910390a26040518060400160405280610b998460200160208101906104d49190613404565b8152602001610be66040805160ff7f000000000000000000000000000000000000000000000000000000000000000016602082015260609101604051602081830303815290604052905090565b905292915050565b67ffffffffffffffff8116600090815260076020526040812060609190610c1790600501612258565b90506000815167ffffffffffffffff811115610c3557610c356135cf565b604051908082528060200260200182016040528015610c6857816020015b6060815260200190600190039081610c535790505b50905060005b8251811015610d515760086000848381518110610c8d57610c8d6137cf565b602002602001015181526020019081526020016000208054610cae906137fe565b80601f0160208091040260200160405190810160405280929190818152602001828054610cda906137fe565b8015610d275780601f10610cfc57610100808354040283529160200191610d27565b820191906000526020600020905b815481529060010190602001808311610d0a57829003601f168201915b5050505050828281518110610d3e57610d3e6137cf565b6020908102919091010152600101610c6e565b509392505050565b6060610d656002612258565b905090565b610d72611d2d565b610d7b83610afe565b610dbd576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84166004820152602401610961565b610dfd8282604051610dd09291906137bf565b604080519182900390912067ffffffffffffffff8616600090815260076020529190912060050190612265565b610e39578282826040517f74f23c7c0000000000000000000000000000000000000000000000000000000081526004016109619392919061389a565b8267ffffffffffffffff167f52d00ee4d9bd51b40168f2afc5848837288ce258784ad914278791464b3f4d768383604051610e759291906138be565b60405180910390a2505050565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915267ffffffffffffffff8216600090815260076020908152604091829020825160a08101845260028201546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff16151594820194909452600390910154808416606083015291909104909116608082015261067890612271565b67ffffffffffffffff81166000908152600760205260409020600401805460609190610f82906137fe565b80601f0160208091040260200160405190810160405280929190818152602001828054610fae906137fe565b8015610ffb5780601f10610fd057610100808354040283529160200191610ffb565b820191906000526020600020905b815481529060010190602001808311610fde57829003601f168201915b50505050509050919050565b61100f611d2d565b73ffffffffffffffffffffffffffffffffffffffff811661105c576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6004805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527f02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f1684910160405180910390a15050565b606060006110f06005612258565b90506000815167ffffffffffffffff81111561110e5761110e6135cf565b604051908082528060200260200182016040528015611137578160200160208202803683370190505b50905060005b825181101561119357828181518110611158576111586137cf565b6020026020010151828281518110611172576111726137cf565b67ffffffffffffffff9092166020928302919091019091015260010161113d565b5092915050565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915267ffffffffffffffff8216600090815260076020908152604091829020825160a08101845281546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff16151594820194909452600190910154808416606083015291909104909116608082015261067890612271565b60095473ffffffffffffffffffffffffffffffffffffffff1633148015906112ac575060015473ffffffffffffffffffffffffffffffffffffffff163314155b156112e5576040517f8e4a23d6000000000000000000000000000000000000000000000000000000008152336004820152602401610961565b6109aa838383612323565b6112f8611d2d565b60005b838110156114e5576000858583818110611317576113176137cf565b905060200201602081019061132c9190613404565b9050611343600567ffffffffffffffff8316612265565b611385576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff82166004820152602401610961565b67ffffffffffffffff811660009081526007602052604081206113aa90600501612258565b905060005b81518110156114165761140d8282815181106113cd576113cd6137cf565b6020026020010151600760008667ffffffffffffffff1667ffffffffffffffff16815260200190815260200160002060050161226590919063ffffffff16565b506001016113af565b5067ffffffffffffffff8216600090815260076020526040812080547fffffffffffffffffffffff0000000000000000000000000000000000000000009081168255600182018390556002820180549091169055600381018290559061147f600483018261310b565b60058201600081816114918282613145565b505060405167ffffffffffffffff871681527f5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d859916945060200192506114d3915050565b60405180910390a150506001016112fb565b5060005b818110156117fb576000838383818110611505576115056137cf565b905060200281019061151791906138d2565b6115209061399e565b90506115318160600151600061240d565b6115408160800151600061240d565b80604001515160000361157f576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80516115979060059067ffffffffffffffff1661254a565b6115dc5780516040517f1d5ad3c500000000000000000000000000000000000000000000000000000000815267ffffffffffffffff9091166004820152602401610961565b805167ffffffffffffffff16600090815260076020908152604091829020825160a08082018552606080870180518601516fffffffffffffffffffffffffffffffff90811680865263ffffffff42168689018190528351511515878b0181905284518a0151841686890181905294518b0151841660809889018190528954740100000000000000000000000000000000000000009283027fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff7001000000000000000000000000000000008087027fffffffffffffffffffffffff000000000000000000000000000000000000000094851690981788178216929092178d5592810290971760018c01558c519889018d52898e0180518d01518716808b528a8e019590955280515115158a8f018190528151909d01518716988a01899052518d0151909516979098018790526002890180549a90910299909316171790941695909517909255909202909117600382015590820151600482019061175f9082613b15565b5060005b8260200151518110156117a35761179b83600001518460200151838151811061178e5761178e6137cf565b6020026020010151611f36565b600101611763565b507f8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c282600001518360400151846060015185608001516040516117e99493929190613c2f565b60405180910390a150506001016114e9565b5050505050565b61180a611d2d565b61181381612556565b50565b61182961029a60a0830160808401613253565b6118885761183d60a0820160808301613253565b6040517f961c9a4f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401610961565b73ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016632cbc26bb6118d46040840160208501613404565b60405160e083901b7fffffffff0000000000000000000000000000000000000000000000000000000016815260809190911b77ffffffffffffffff00000000000000000000000000000000166004820152602401602060405180830381865afa158015611945573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906119699190613cc8565b156119a0576040517f53ad11d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6119b86119b36040830160208401613404565b61261a565b6119d86119cb6040830160208401613404565b61033a60a084018461375a565b611a1d576119e960a082018261375a565b6040517f24eb47e50000000000000000000000000000000000000000000000000000000081526004016109619291906138be565b611813611a306040830160208401613404565b8260600135612740565b60008151600003611a6c57507f0000000000000000000000000000000000000000000000000000000000000000919050565b8151602014611aa957816040517f953576f7000000000000000000000000000000000000000000000000000000008152600401610961919061321e565b600082806020019051810190611abf9190613ce5565b905060ff81111561067857826040517f953576f7000000000000000000000000000000000000000000000000000000008152600401610961919061321e565b60007f000000000000000000000000000000000000000000000000000000000000000060ff168260ff1603611b34575081610678565b7f000000000000000000000000000000000000000000000000000000000000000060ff168260ff161115611c1f576000611b8e7f000000000000000000000000000000000000000000000000000000000000000084613d2d565b9050604d8160ff161115611c02576040517fa9cb113d00000000000000000000000000000000000000000000000000000000815260ff80851660048301527f000000000000000000000000000000000000000000000000000000000000000016602482015260448101859052606401610961565b611c0d81600a613e66565b611c179085613e75565b915050610678565b6000611c4b837f0000000000000000000000000000000000000000000000000000000000000000613d2d565b9050604d8160ff161180611c925750611c6581600a613e66565b611c8f907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff613e75565b84115b15611cfd576040517fa9cb113d00000000000000000000000000000000000000000000000000000000815260ff80851660048301527f000000000000000000000000000000000000000000000000000000000000000016602482015260448101859052606401610961565b611d0881600a613e66565b61088f9085613eb0565b600081815260018301602052604081205415155b9392505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314611d7e576040517f2b5c74de00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b7f0000000000000000000000000000000000000000000000000000000000000000611dd7576040517f35f4a7b300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b8251811015611e6d576000838281518110611df757611df76137cf565b60200260200101519050611e1581600261278790919063ffffffff16565b15611e645760405173ffffffffffffffffffffffffffffffffffffffff821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b50600101611dda565b5060005b81518110156109aa576000828281518110611e8e57611e8e6137cf565b60200260200101519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603611ed25750611f2e565b611edd6002826127a9565b15611f2c5760405173ffffffffffffffffffffffffffffffffffffffff821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b600101611e71565b8051600003611f71576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805160208083019190912067ffffffffffffffff8416600090815260079092526040909120611fa3906005018261254a565b611fdd5782826040517f393b8ad2000000000000000000000000000000000000000000000000000000008152600401610961929190613ec7565b6000818152600860205260409020611ff58382613b15565b508267ffffffffffffffff167f7d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea83604051610e75919061321e565b61204361029a60a0830160808401613253565b6120575761183d60a0820160808301613253565b73ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016632cbc26bb6120a36040840160208501613404565b60405160e083901b7fffffffff0000000000000000000000000000000000000000000000000000000016815260809190911b77ffffffffffffffff00000000000000000000000000000000166004820152602401602060405180830381865afa158015612114573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906121389190613cc8565b1561216f576040517f53ad11d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6121876121826060830160408401613253565b6127cb565b61219f61219a6040830160208401613404565b61284a565b6118136121b26040830160208401613404565b8260600135612998565b6040517f42966c68000000000000000000000000000000000000000000000000000000008152600481018290527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906342966c6890602401600060405180830381600087803b15801561224457600080fd5b505af11580156117fb573d6000803e3d6000fd5b60606000611d26836129dc565b6000611d268383612a37565b6040805160a0810182526000808252602082018190529181018290526060810182905260808101919091526122ff82606001516fffffffffffffffffffffffffffffffff1683600001516fffffffffffffffffffffffffffffffff16846020015163ffffffff16426122e39190613eea565b85608001516fffffffffffffffffffffffffffffffff16612b2a565b6fffffffffffffffffffffffffffffffff1682525063ffffffff4216602082015290565b61232c83610afe565b61236e576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84166004820152602401610961565b61237982600061240d565b67ffffffffffffffff8316600090815260076020526040902061239c9083612b52565b6123a781600061240d565b67ffffffffffffffff831660009081526007602052604090206123cd9060020182612b52565b7f0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b83838360405161240093929190613efd565b60405180910390a1505050565b8151156124d85781602001516fffffffffffffffffffffffffffffffff1682604001516fffffffffffffffffffffffffffffffff16101580612463575060408201516fffffffffffffffffffffffffffffffff16155b1561249c57816040517f8020d1240000000000000000000000000000000000000000000000000000000081526004016109619190613f80565b80156124d4576040517f433fc33d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050565b60408201516fffffffffffffffffffffffffffffffff16151580612511575060208201516fffffffffffffffffffffffffffffffff1615155b156124d457816040517fd68af9cc0000000000000000000000000000000000000000000000000000000081526004016109619190613f80565b6000611d268383612cf4565b3373ffffffffffffffffffffffffffffffffffffffff8216036125a5576040517fdad89dca00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b61262381610afe565b612665576040517fa9902c7e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff82166004820152602401610961565b600480546040517f83826b2b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84169281019290925233602483015273ffffffffffffffffffffffffffffffffffffffff16906383826b2b90604401602060405180830381865afa1580156126e4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906127089190613cc8565b611813576040517f728fe07b000000000000000000000000000000000000000000000000000000008152336004820152602401610961565b67ffffffffffffffff821660009081526007602052604090206124d490600201827f0000000000000000000000000000000000000000000000000000000000000000612d43565b6000611d268373ffffffffffffffffffffffffffffffffffffffff8416612a37565b6000611d268373ffffffffffffffffffffffffffffffffffffffff8416612cf4565b7f000000000000000000000000000000000000000000000000000000000000000015611813576127fc6002826130c6565b611813576040517fd0d2597600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610961565b61285381610afe565b612895576040517fa9902c7e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff82166004820152602401610961565b600480546040517fa8d87a3b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84169281019290925273ffffffffffffffffffffffffffffffffffffffff169063a8d87a3b90602401602060405180830381865afa15801561290e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906129329190613fbc565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614611813576040517f728fe07b000000000000000000000000000000000000000000000000000000008152336004820152602401610961565b67ffffffffffffffff821660009081526007602052604090206124d490827f0000000000000000000000000000000000000000000000000000000000000000612d43565b606081600001805480602002602001604051908101604052809291908181526020018280548015610ffb57602002820191906000526020600020905b815481526020019060010190808311612a185750505050509050919050565b60008181526001830160205260408120548015612b20576000612a5b600183613eea565b8554909150600090612a6f90600190613eea565b9050808214612ad4576000866000018281548110612a8f57612a8f6137cf565b9060005260206000200154905080876000018481548110612ab257612ab26137cf565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080612ae557612ae5613fd9565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610678565b6000915050610678565b6000612b4985612b3a8486613eb0565b612b449087614008565b6130f5565b95945050505050565b8154600090612b7b90700100000000000000000000000000000000900463ffffffff1642613eea565b90508015612c1d5760018301548354612bc3916fffffffffffffffffffffffffffffffff80821692811691859170010000000000000000000000000000000090910416612b2a565b83546fffffffffffffffffffffffffffffffff919091167fffffffffffffffffffffffff0000000000000000000000000000000000000000909116177001000000000000000000000000000000004263ffffffff16021783555b60208201518354612c43916fffffffffffffffffffffffffffffffff90811691166130f5565b83548351151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00ffffffff000000000000000000000000000000009091166fffffffffffffffffffffffffffffffff92831617178455602083015160408085015183167001000000000000000000000000000000000291909216176001850155517f9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c1990612400908490613f80565b6000818152600183016020526040812054612d3b57508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610678565b506000610678565b825474010000000000000000000000000000000000000000900460ff161580612d6a575081155b15612d7457505050565b825460018401546fffffffffffffffffffffffffffffffff80831692911690600090612dba90700100000000000000000000000000000000900463ffffffff1642613eea565b90508015612e7a5781831115612dfc576040517f9725942a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001860154612e369083908590849070010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff16612b2a565b86547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff167001000000000000000000000000000000004263ffffffff160217875592505b84821015612f315773ffffffffffffffffffffffffffffffffffffffff8416612ed9576040517ff94ebcd10000000000000000000000000000000000000000000000000000000081526004810183905260248101869052604401610961565b6040517f1a76572a000000000000000000000000000000000000000000000000000000008152600481018390526024810186905273ffffffffffffffffffffffffffffffffffffffff85166044820152606401610961565b848310156130445760018681015470010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff16906000908290612f759082613eea565b612f7f878a613eea565b612f899190614008565b612f939190613e75565b905073ffffffffffffffffffffffffffffffffffffffff8616612fec576040517f15279c080000000000000000000000000000000000000000000000000000000081526004810182905260248101869052604401610961565b6040517fd0c8d23a000000000000000000000000000000000000000000000000000000008152600481018290526024810186905273ffffffffffffffffffffffffffffffffffffffff87166044820152606401610961565b61304e8584613eea565b86547fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff82161787556040518681529093507f1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a9060200160405180910390a1505050505050565b73ffffffffffffffffffffffffffffffffffffffff811660009081526001830160205260408120541515611d26565b60008183106131045781611d26565b5090919050565b508054613117906137fe565b6000825580601f10613127575050565b601f016020900490600052602060002090810190611813919061315f565b508054600082559060005260206000209081019061181391905b5b808211156131745760008155600101613160565b5090565b60006020828403121561318a57600080fd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114611d2657600080fd5b6000815180845260005b818110156131e0576020818501810151868301820152016131c4565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081526000611d2660208301846131ba565b73ffffffffffffffffffffffffffffffffffffffff8116811461181357600080fd5b60006020828403121561326557600080fd5b8135611d2681613231565b60006020828403121561328257600080fd5b813567ffffffffffffffff81111561329957600080fd5b82016101008185031215611d2657600080fd5b803567ffffffffffffffff811681146132c457600080fd5b919050565b6000806000604084860312156132de57600080fd5b6132e7846132ac565b9250602084013567ffffffffffffffff8082111561330457600080fd5b818601915086601f83011261331857600080fd5b81358181111561332757600080fd5b87602082850101111561333957600080fd5b6020830194508093505050509250925092565b60008083601f84011261335e57600080fd5b50813567ffffffffffffffff81111561337657600080fd5b6020830191508360208260051b850101111561339157600080fd5b9250929050565b600080600080604085870312156133ae57600080fd5b843567ffffffffffffffff808211156133c657600080fd5b6133d28883890161334c565b909650945060208701359150808211156133eb57600080fd5b506133f88782880161334c565b95989497509550505050565b60006020828403121561341657600080fd5b611d26826132ac565b60006020828403121561343157600080fd5b813567ffffffffffffffff81111561344857600080fd5b820160a08185031215611d2657600080fd5b60208152600082516040602084015261347660608401826131ba565b905060208401517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0848303016040850152612b4982826131ba565b600060208083016020845280855180835260408601915060408160051b87010192506020870160005b82811015613526577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08886030184526135148583516131ba565b945092850192908501906001016134da565b5092979650505050505050565b6020808252825182820181905260009190848201906040850190845b8181101561358157835173ffffffffffffffffffffffffffffffffffffffff168352928401929184019160010161354f565b50909695505050505050565b6020808252825182820181905260009190848201906040850190845b8181101561358157835167ffffffffffffffff16835292840192918401916001016135a9565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160a0810167ffffffffffffffff81118282101715613621576136216135cf565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561366e5761366e6135cf565b604052919050565b801515811461181357600080fd5b80356fffffffffffffffffffffffffffffffff811681146132c457600080fd5b6000606082840312156136b657600080fd5b6040516060810181811067ffffffffffffffff821117156136d9576136d96135cf565b60405290508082356136ea81613676565b81526136f860208401613684565b602082015261370960408401613684565b60408201525092915050565b600080600060e0848603121561372a57600080fd5b613733846132ac565b925061374285602086016136a4565b915061375185608086016136a4565b90509250925092565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261378f57600080fd5b83018035915067ffffffffffffffff8211156137aa57600080fd5b60200191503681900382131561339157600080fd5b8183823760009101908152919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600181811c9082168061381257607f821691505b60208210810361384b577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b67ffffffffffffffff84168152604060208201526000612b49604083018486613851565b60208152600061088f602083018486613851565b600082357ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffee183360301811261390657600080fd5b9190910192915050565b600082601f83011261392157600080fd5b813567ffffffffffffffff81111561393b5761393b6135cf565b61396c60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601613627565b81815284602083860101111561398157600080fd5b816020850160208301376000918101602001919091529392505050565b600061012082360312156139b157600080fd5b6139b96135fe565b6139c2836132ac565b815260208084013567ffffffffffffffff808211156139e057600080fd5b9085019036601f8301126139f357600080fd5b813581811115613a0557613a056135cf565b8060051b613a14858201613627565b9182528381018501918581019036841115613a2e57600080fd5b86860192505b83831015613a6a57823585811115613a4c5760008081fd5b613a5a3689838a0101613910565b8352509186019190860190613a34565b8087890152505050506040860135925080831115613a8757600080fd5b5050613a9536828601613910565b604083015250613aa836606085016136a4565b6060820152613aba3660c085016136a4565b608082015292915050565b601f8211156109aa576000816000526020600020601f850160051c81016020861015613aee5750805b601f850160051c820191505b81811015613b0d57828155600101613afa565b505050505050565b815167ffffffffffffffff811115613b2f57613b2f6135cf565b613b4381613b3d84546137fe565b84613ac5565b602080601f831160018114613b965760008415613b605750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555613b0d565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015613be357888601518255948401946001909101908401613bc4565b5085821015613c1f57878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b600061010067ffffffffffffffff87168352806020840152613c53818401876131ba565b8551151560408581019190915260208701516fffffffffffffffffffffffffffffffff9081166060870152908701511660808501529150613c919050565b8251151560a083015260208301516fffffffffffffffffffffffffffffffff90811660c084015260408401511660e0830152612b49565b600060208284031215613cda57600080fd5b8151611d2681613676565b600060208284031215613cf757600080fd5b5051919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60ff828116828216039081111561067857610678613cfe565b600181815b80851115613d9f57817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04821115613d8557613d85613cfe565b80851615613d9257918102915b93841c9390800290613d4b565b509250929050565b600082613db657506001610678565b81613dc357506000610678565b8160018114613dd95760028114613de357613dff565b6001915050610678565b60ff841115613df457613df4613cfe565b50506001821b610678565b5060208310610133831016604e8410600b8410161715613e22575081810a610678565b613e2c8383613d46565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04821115613e5e57613e5e613cfe565b029392505050565b6000611d2660ff841683613da7565b600082613eab577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b808202811582820484141761067857610678613cfe565b67ffffffffffffffff8316815260406020820152600061088f60408301846131ba565b8181038181111561067857610678613cfe565b67ffffffffffffffff8416815260e08101613f4960208301858051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b82511515608083015260208301516fffffffffffffffffffffffffffffffff90811660a084015260408401511660c083015261088f565b6060810161067882848051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b600060208284031215613fce57600080fd5b8151611d2681613231565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b8082018082111561067857610678613cfe56fea164736f6c6343000818000a",
}

var BurnMintTokenPoolABI = BurnMintTokenPoolMetaData.ABI

var BurnMintTokenPoolBin = BurnMintTokenPoolMetaData.Bin

func DeployBurnMintTokenPool(auth *bind.TransactOpts, backend bind.ContractBackend, token common.Address, localTokenDecimals uint8, allowlist []common.Address, rmnProxy common.Address, router common.Address) (common.Address, *generated.Transaction, *BurnMintTokenPool, error) {
	parsed, err := BurnMintTokenPoolMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}
	if generated.IsZKSync(backend) {
		address, ethTx, contractBind, _ := generated.DeployContract(auth, parsed, common.FromHex(BurnMintTokenPoolZKBin), backend, token, localTokenDecimals, allowlist, rmnProxy, router)
		contractReturn := &BurnMintTokenPool{address: address, abi: *parsed, BurnMintTokenPoolCaller: BurnMintTokenPoolCaller{contract: contractBind}, BurnMintTokenPoolTransactor: BurnMintTokenPoolTransactor{contract: contractBind}, BurnMintTokenPoolFilterer: BurnMintTokenPoolFilterer{contract: contractBind}}
		return address, ethTx, contractReturn, err
	}
	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BurnMintTokenPoolBin), backend, token, localTokenDecimals, allowlist, rmnProxy, router)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, &generated.Transaction{Transaction: tx, Hash_zks: tx.Hash()}, &BurnMintTokenPool{address: address, abi: *parsed, BurnMintTokenPoolCaller: BurnMintTokenPoolCaller{contract: contract}, BurnMintTokenPoolTransactor: BurnMintTokenPoolTransactor{contract: contract}, BurnMintTokenPoolFilterer: BurnMintTokenPoolFilterer{contract: contract}}, nil
}

type BurnMintTokenPool struct {
	address common.Address
	abi     abi.ABI
	BurnMintTokenPoolCaller
	BurnMintTokenPoolTransactor
	BurnMintTokenPoolFilterer
}

type BurnMintTokenPoolCaller struct {
	contract *bind.BoundContract
}

type BurnMintTokenPoolTransactor struct {
	contract *bind.BoundContract
}

type BurnMintTokenPoolFilterer struct {
	contract *bind.BoundContract
}

type BurnMintTokenPoolSession struct {
	Contract     *BurnMintTokenPool
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type BurnMintTokenPoolCallerSession struct {
	Contract *BurnMintTokenPoolCaller
	CallOpts bind.CallOpts
}

type BurnMintTokenPoolTransactorSession struct {
	Contract     *BurnMintTokenPoolTransactor
	TransactOpts bind.TransactOpts
}

type BurnMintTokenPoolRaw struct {
	Contract *BurnMintTokenPool
}

type BurnMintTokenPoolCallerRaw struct {
	Contract *BurnMintTokenPoolCaller
}

type BurnMintTokenPoolTransactorRaw struct {
	Contract *BurnMintTokenPoolTransactor
}

func NewBurnMintTokenPool(address common.Address, backend bind.ContractBackend) (*BurnMintTokenPool, error) {
	abi, err := abi.JSON(strings.NewReader(BurnMintTokenPoolABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindBurnMintTokenPool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPool{address: address, abi: abi, BurnMintTokenPoolCaller: BurnMintTokenPoolCaller{contract: contract}, BurnMintTokenPoolTransactor: BurnMintTokenPoolTransactor{contract: contract}, BurnMintTokenPoolFilterer: BurnMintTokenPoolFilterer{contract: contract}}, nil
}

func NewBurnMintTokenPoolCaller(address common.Address, caller bind.ContractCaller) (*BurnMintTokenPoolCaller, error) {
	contract, err := bindBurnMintTokenPool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolCaller{contract: contract}, nil
}

func NewBurnMintTokenPoolTransactor(address common.Address, transactor bind.ContractTransactor) (*BurnMintTokenPoolTransactor, error) {
	contract, err := bindBurnMintTokenPool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolTransactor{contract: contract}, nil
}

func NewBurnMintTokenPoolFilterer(address common.Address, filterer bind.ContractFilterer) (*BurnMintTokenPoolFilterer, error) {
	contract, err := bindBurnMintTokenPool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolFilterer{contract: contract}, nil
}

func bindBurnMintTokenPool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BurnMintTokenPoolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BurnMintTokenPool.Contract.BurnMintTokenPoolCaller.contract.Call(opts, result, method, params...)
}

func (_BurnMintTokenPool *BurnMintTokenPoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.BurnMintTokenPoolTransactor.contract.Transfer(opts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.BurnMintTokenPoolTransactor.contract.Transact(opts, method, params...)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BurnMintTokenPool.Contract.contract.Call(opts, result, method, params...)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.contract.Transfer(opts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.contract.Transact(opts, method, params...)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCaller) GetAllowList(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _BurnMintTokenPool.contract.Call(opts, &out, "getAllowList")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) GetAllowList() ([]common.Address, error) {
	return _BurnMintTokenPool.Contract.GetAllowList(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCallerSession) GetAllowList() ([]common.Address, error) {
	return _BurnMintTokenPool.Contract.GetAllowList(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCaller) GetAllowListEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _BurnMintTokenPool.contract.Call(opts, &out, "getAllowListEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) GetAllowListEnabled() (bool, error) {
	return _BurnMintTokenPool.Contract.GetAllowListEnabled(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCallerSession) GetAllowListEnabled() (bool, error) {
	return _BurnMintTokenPool.Contract.GetAllowListEnabled(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCaller) GetCurrentInboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _BurnMintTokenPool.contract.Call(opts, &out, "getCurrentInboundRateLimiterState", remoteChainSelector)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) GetCurrentInboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnMintTokenPool.Contract.GetCurrentInboundRateLimiterState(&_BurnMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCallerSession) GetCurrentInboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnMintTokenPool.Contract.GetCurrentInboundRateLimiterState(&_BurnMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCaller) GetCurrentOutboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _BurnMintTokenPool.contract.Call(opts, &out, "getCurrentOutboundRateLimiterState", remoteChainSelector)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) GetCurrentOutboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnMintTokenPool.Contract.GetCurrentOutboundRateLimiterState(&_BurnMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCallerSession) GetCurrentOutboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnMintTokenPool.Contract.GetCurrentOutboundRateLimiterState(&_BurnMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCaller) GetRateLimitAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnMintTokenPool.contract.Call(opts, &out, "getRateLimitAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) GetRateLimitAdmin() (common.Address, error) {
	return _BurnMintTokenPool.Contract.GetRateLimitAdmin(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCallerSession) GetRateLimitAdmin() (common.Address, error) {
	return _BurnMintTokenPool.Contract.GetRateLimitAdmin(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCaller) GetRemotePools(opts *bind.CallOpts, remoteChainSelector uint64) ([][]byte, error) {
	var out []interface{}
	err := _BurnMintTokenPool.contract.Call(opts, &out, "getRemotePools", remoteChainSelector)

	if err != nil {
		return *new([][]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][]byte)).(*[][]byte)

	return out0, err

}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) GetRemotePools(remoteChainSelector uint64) ([][]byte, error) {
	return _BurnMintTokenPool.Contract.GetRemotePools(&_BurnMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCallerSession) GetRemotePools(remoteChainSelector uint64) ([][]byte, error) {
	return _BurnMintTokenPool.Contract.GetRemotePools(&_BurnMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCaller) GetRemoteToken(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error) {
	var out []interface{}
	err := _BurnMintTokenPool.contract.Call(opts, &out, "getRemoteToken", remoteChainSelector)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) GetRemoteToken(remoteChainSelector uint64) ([]byte, error) {
	return _BurnMintTokenPool.Contract.GetRemoteToken(&_BurnMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCallerSession) GetRemoteToken(remoteChainSelector uint64) ([]byte, error) {
	return _BurnMintTokenPool.Contract.GetRemoteToken(&_BurnMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCaller) GetRmnProxy(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnMintTokenPool.contract.Call(opts, &out, "getRmnProxy")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) GetRmnProxy() (common.Address, error) {
	return _BurnMintTokenPool.Contract.GetRmnProxy(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCallerSession) GetRmnProxy() (common.Address, error) {
	return _BurnMintTokenPool.Contract.GetRmnProxy(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCaller) GetRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnMintTokenPool.contract.Call(opts, &out, "getRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) GetRouter() (common.Address, error) {
	return _BurnMintTokenPool.Contract.GetRouter(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCallerSession) GetRouter() (common.Address, error) {
	return _BurnMintTokenPool.Contract.GetRouter(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCaller) GetSupportedChains(opts *bind.CallOpts) ([]uint64, error) {
	var out []interface{}
	err := _BurnMintTokenPool.contract.Call(opts, &out, "getSupportedChains")

	if err != nil {
		return *new([]uint64), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint64)).(*[]uint64)

	return out0, err

}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) GetSupportedChains() ([]uint64, error) {
	return _BurnMintTokenPool.Contract.GetSupportedChains(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCallerSession) GetSupportedChains() ([]uint64, error) {
	return _BurnMintTokenPool.Contract.GetSupportedChains(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCaller) GetToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnMintTokenPool.contract.Call(opts, &out, "getToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) GetToken() (common.Address, error) {
	return _BurnMintTokenPool.Contract.GetToken(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCallerSession) GetToken() (common.Address, error) {
	return _BurnMintTokenPool.Contract.GetToken(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCaller) GetTokenDecimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _BurnMintTokenPool.contract.Call(opts, &out, "getTokenDecimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) GetTokenDecimals() (uint8, error) {
	return _BurnMintTokenPool.Contract.GetTokenDecimals(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCallerSession) GetTokenDecimals() (uint8, error) {
	return _BurnMintTokenPool.Contract.GetTokenDecimals(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCaller) IsRemotePool(opts *bind.CallOpts, remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	var out []interface{}
	err := _BurnMintTokenPool.contract.Call(opts, &out, "isRemotePool", remoteChainSelector, remotePoolAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) IsRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	return _BurnMintTokenPool.Contract.IsRemotePool(&_BurnMintTokenPool.CallOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCallerSession) IsRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	return _BurnMintTokenPool.Contract.IsRemotePool(&_BurnMintTokenPool.CallOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCaller) IsSupportedChain(opts *bind.CallOpts, remoteChainSelector uint64) (bool, error) {
	var out []interface{}
	err := _BurnMintTokenPool.contract.Call(opts, &out, "isSupportedChain", remoteChainSelector)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) IsSupportedChain(remoteChainSelector uint64) (bool, error) {
	return _BurnMintTokenPool.Contract.IsSupportedChain(&_BurnMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCallerSession) IsSupportedChain(remoteChainSelector uint64) (bool, error) {
	return _BurnMintTokenPool.Contract.IsSupportedChain(&_BurnMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCaller) IsSupportedToken(opts *bind.CallOpts, token common.Address) (bool, error) {
	var out []interface{}
	err := _BurnMintTokenPool.contract.Call(opts, &out, "isSupportedToken", token)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) IsSupportedToken(token common.Address) (bool, error) {
	return _BurnMintTokenPool.Contract.IsSupportedToken(&_BurnMintTokenPool.CallOpts, token)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCallerSession) IsSupportedToken(token common.Address) (bool, error) {
	return _BurnMintTokenPool.Contract.IsSupportedToken(&_BurnMintTokenPool.CallOpts, token)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnMintTokenPool.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) Owner() (common.Address, error) {
	return _BurnMintTokenPool.Contract.Owner(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCallerSession) Owner() (common.Address, error) {
	return _BurnMintTokenPool.Contract.Owner(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _BurnMintTokenPool.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BurnMintTokenPool.Contract.SupportsInterface(&_BurnMintTokenPool.CallOpts, interfaceId)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BurnMintTokenPool.Contract.SupportsInterface(&_BurnMintTokenPool.CallOpts, interfaceId)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _BurnMintTokenPool.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) TypeAndVersion() (string, error) {
	return _BurnMintTokenPool.Contract.TypeAndVersion(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolCallerSession) TypeAndVersion() (string, error) {
	return _BurnMintTokenPool.Contract.TypeAndVersion(&_BurnMintTokenPool.CallOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnMintTokenPool.contract.Transact(opts, "acceptOwnership")
}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) AcceptOwnership() (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.AcceptOwnership(&_BurnMintTokenPool.TransactOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.AcceptOwnership(&_BurnMintTokenPool.TransactOpts)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactor) AddRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnMintTokenPool.contract.Transact(opts, "addRemotePool", remoteChainSelector, remotePoolAddress)
}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) AddRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.AddRemotePool(&_BurnMintTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactorSession) AddRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.AddRemotePool(&_BurnMintTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactor) ApplyAllowListUpdates(opts *bind.TransactOpts, removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPool.contract.Transact(opts, "applyAllowListUpdates", removes, adds)
}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.ApplyAllowListUpdates(&_BurnMintTokenPool.TransactOpts, removes, adds)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactorSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.ApplyAllowListUpdates(&_BurnMintTokenPool.TransactOpts, removes, adds)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactor) ApplyChainUpdates(opts *bind.TransactOpts, remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _BurnMintTokenPool.contract.Transact(opts, "applyChainUpdates", remoteChainSelectorsToRemove, chainsToAdd)
}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) ApplyChainUpdates(remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.ApplyChainUpdates(&_BurnMintTokenPool.TransactOpts, remoteChainSelectorsToRemove, chainsToAdd)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactorSession) ApplyChainUpdates(remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.ApplyChainUpdates(&_BurnMintTokenPool.TransactOpts, remoteChainSelectorsToRemove, chainsToAdd)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactor) LockOrBurn(opts *bind.TransactOpts, lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _BurnMintTokenPool.contract.Transact(opts, "lockOrBurn", lockOrBurnIn)
}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) LockOrBurn(lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.LockOrBurn(&_BurnMintTokenPool.TransactOpts, lockOrBurnIn)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactorSession) LockOrBurn(lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.LockOrBurn(&_BurnMintTokenPool.TransactOpts, lockOrBurnIn)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactor) ReleaseOrMint(opts *bind.TransactOpts, releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _BurnMintTokenPool.contract.Transact(opts, "releaseOrMint", releaseOrMintIn)
}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) ReleaseOrMint(releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.ReleaseOrMint(&_BurnMintTokenPool.TransactOpts, releaseOrMintIn)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactorSession) ReleaseOrMint(releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.ReleaseOrMint(&_BurnMintTokenPool.TransactOpts, releaseOrMintIn)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactor) RemoveRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnMintTokenPool.contract.Transact(opts, "removeRemotePool", remoteChainSelector, remotePoolAddress)
}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) RemoveRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.RemoveRemotePool(&_BurnMintTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactorSession) RemoveRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.RemoveRemotePool(&_BurnMintTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactor) SetChainRateLimiterConfig(opts *bind.TransactOpts, remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _BurnMintTokenPool.contract.Transact(opts, "setChainRateLimiterConfig", remoteChainSelector, outboundConfig, inboundConfig)
}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) SetChainRateLimiterConfig(remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.SetChainRateLimiterConfig(&_BurnMintTokenPool.TransactOpts, remoteChainSelector, outboundConfig, inboundConfig)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactorSession) SetChainRateLimiterConfig(remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.SetChainRateLimiterConfig(&_BurnMintTokenPool.TransactOpts, remoteChainSelector, outboundConfig, inboundConfig)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactor) SetRateLimitAdmin(opts *bind.TransactOpts, rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPool.contract.Transact(opts, "setRateLimitAdmin", rateLimitAdmin)
}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) SetRateLimitAdmin(rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.SetRateLimitAdmin(&_BurnMintTokenPool.TransactOpts, rateLimitAdmin)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactorSession) SetRateLimitAdmin(rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.SetRateLimitAdmin(&_BurnMintTokenPool.TransactOpts, rateLimitAdmin)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactor) SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPool.contract.Transact(opts, "setRouter", newRouter)
}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.SetRouter(&_BurnMintTokenPool.TransactOpts, newRouter)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactorSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.SetRouter(&_BurnMintTokenPool.TransactOpts, newRouter)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPool.contract.Transact(opts, "transferOwnership", to)
}

func (_BurnMintTokenPool *BurnMintTokenPoolSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.TransferOwnership(&_BurnMintTokenPool.TransactOpts, to)
}

func (_BurnMintTokenPool *BurnMintTokenPoolTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _BurnMintTokenPool.Contract.TransferOwnership(&_BurnMintTokenPool.TransactOpts, to)
}

type BurnMintTokenPoolAllowListAddIterator struct {
	Event *BurnMintTokenPoolAllowListAdd

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolAllowListAddIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolAllowListAdd)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(BurnMintTokenPoolAllowListAdd)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *BurnMintTokenPoolAllowListAddIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolAllowListAddIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolAllowListAdd struct {
	Sender common.Address
	Raw    types.Log
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) FilterAllowListAdd(opts *bind.FilterOpts) (*BurnMintTokenPoolAllowListAddIterator, error) {

	logs, sub, err := _BurnMintTokenPool.contract.FilterLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolAllowListAddIterator{contract: _BurnMintTokenPool.contract, event: "AllowListAdd", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAllowListAdd) (event.Subscription, error) {

	logs, sub, err := _BurnMintTokenPool.contract.WatchLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolAllowListAdd)
				if err := _BurnMintTokenPool.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) ParseAllowListAdd(log types.Log) (*BurnMintTokenPoolAllowListAdd, error) {
	event := new(BurnMintTokenPoolAllowListAdd)
	if err := _BurnMintTokenPool.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolAllowListRemoveIterator struct {
	Event *BurnMintTokenPoolAllowListRemove

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolAllowListRemoveIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolAllowListRemove)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(BurnMintTokenPoolAllowListRemove)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *BurnMintTokenPoolAllowListRemoveIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolAllowListRemoveIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolAllowListRemove struct {
	Sender common.Address
	Raw    types.Log
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) FilterAllowListRemove(opts *bind.FilterOpts) (*BurnMintTokenPoolAllowListRemoveIterator, error) {

	logs, sub, err := _BurnMintTokenPool.contract.FilterLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolAllowListRemoveIterator{contract: _BurnMintTokenPool.contract, event: "AllowListRemove", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAllowListRemove) (event.Subscription, error) {

	logs, sub, err := _BurnMintTokenPool.contract.WatchLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolAllowListRemove)
				if err := _BurnMintTokenPool.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) ParseAllowListRemove(log types.Log) (*BurnMintTokenPoolAllowListRemove, error) {
	event := new(BurnMintTokenPoolAllowListRemove)
	if err := _BurnMintTokenPool.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolBurnedIterator struct {
	Event *BurnMintTokenPoolBurned

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolBurnedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolBurned)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(BurnMintTokenPoolBurned)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *BurnMintTokenPoolBurnedIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolBurnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolBurned struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) FilterBurned(opts *bind.FilterOpts, sender []common.Address) (*BurnMintTokenPoolBurnedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BurnMintTokenPool.contract.FilterLogs(opts, "Burned", senderRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolBurnedIterator{contract: _BurnMintTokenPool.contract, event: "Burned", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) WatchBurned(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolBurned, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BurnMintTokenPool.contract.WatchLogs(opts, "Burned", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolBurned)
				if err := _BurnMintTokenPool.contract.UnpackLog(event, "Burned", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) ParseBurned(log types.Log) (*BurnMintTokenPoolBurned, error) {
	event := new(BurnMintTokenPoolBurned)
	if err := _BurnMintTokenPool.contract.UnpackLog(event, "Burned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolChainAddedIterator struct {
	Event *BurnMintTokenPoolChainAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolChainAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolChainAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(BurnMintTokenPoolChainAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *BurnMintTokenPoolChainAddedIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolChainAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolChainAdded struct {
	RemoteChainSelector       uint64
	RemoteToken               []byte
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
	Raw                       types.Log
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) FilterChainAdded(opts *bind.FilterOpts) (*BurnMintTokenPoolChainAddedIterator, error) {

	logs, sub, err := _BurnMintTokenPool.contract.FilterLogs(opts, "ChainAdded")
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolChainAddedIterator{contract: _BurnMintTokenPool.contract, event: "ChainAdded", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) WatchChainAdded(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolChainAdded) (event.Subscription, error) {

	logs, sub, err := _BurnMintTokenPool.contract.WatchLogs(opts, "ChainAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolChainAdded)
				if err := _BurnMintTokenPool.contract.UnpackLog(event, "ChainAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) ParseChainAdded(log types.Log) (*BurnMintTokenPoolChainAdded, error) {
	event := new(BurnMintTokenPoolChainAdded)
	if err := _BurnMintTokenPool.contract.UnpackLog(event, "ChainAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolChainConfiguredIterator struct {
	Event *BurnMintTokenPoolChainConfigured

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolChainConfiguredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolChainConfigured)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(BurnMintTokenPoolChainConfigured)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *BurnMintTokenPoolChainConfiguredIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolChainConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolChainConfigured struct {
	RemoteChainSelector       uint64
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
	Raw                       types.Log
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) FilterChainConfigured(opts *bind.FilterOpts) (*BurnMintTokenPoolChainConfiguredIterator, error) {

	logs, sub, err := _BurnMintTokenPool.contract.FilterLogs(opts, "ChainConfigured")
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolChainConfiguredIterator{contract: _BurnMintTokenPool.contract, event: "ChainConfigured", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) WatchChainConfigured(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolChainConfigured) (event.Subscription, error) {

	logs, sub, err := _BurnMintTokenPool.contract.WatchLogs(opts, "ChainConfigured")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolChainConfigured)
				if err := _BurnMintTokenPool.contract.UnpackLog(event, "ChainConfigured", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) ParseChainConfigured(log types.Log) (*BurnMintTokenPoolChainConfigured, error) {
	event := new(BurnMintTokenPoolChainConfigured)
	if err := _BurnMintTokenPool.contract.UnpackLog(event, "ChainConfigured", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolChainRemovedIterator struct {
	Event *BurnMintTokenPoolChainRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolChainRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolChainRemoved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(BurnMintTokenPoolChainRemoved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *BurnMintTokenPoolChainRemovedIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolChainRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolChainRemoved struct {
	RemoteChainSelector uint64
	Raw                 types.Log
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) FilterChainRemoved(opts *bind.FilterOpts) (*BurnMintTokenPoolChainRemovedIterator, error) {

	logs, sub, err := _BurnMintTokenPool.contract.FilterLogs(opts, "ChainRemoved")
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolChainRemovedIterator{contract: _BurnMintTokenPool.contract, event: "ChainRemoved", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolChainRemoved) (event.Subscription, error) {

	logs, sub, err := _BurnMintTokenPool.contract.WatchLogs(opts, "ChainRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolChainRemoved)
				if err := _BurnMintTokenPool.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) ParseChainRemoved(log types.Log) (*BurnMintTokenPoolChainRemoved, error) {
	event := new(BurnMintTokenPoolChainRemoved)
	if err := _BurnMintTokenPool.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolConfigChangedIterator struct {
	Event *BurnMintTokenPoolConfigChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolConfigChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolConfigChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(BurnMintTokenPoolConfigChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *BurnMintTokenPoolConfigChangedIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolConfigChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolConfigChanged struct {
	Config RateLimiterConfig
	Raw    types.Log
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) FilterConfigChanged(opts *bind.FilterOpts) (*BurnMintTokenPoolConfigChangedIterator, error) {

	logs, sub, err := _BurnMintTokenPool.contract.FilterLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolConfigChangedIterator{contract: _BurnMintTokenPool.contract, event: "ConfigChanged", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolConfigChanged) (event.Subscription, error) {

	logs, sub, err := _BurnMintTokenPool.contract.WatchLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolConfigChanged)
				if err := _BurnMintTokenPool.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) ParseConfigChanged(log types.Log) (*BurnMintTokenPoolConfigChanged, error) {
	event := new(BurnMintTokenPoolConfigChanged)
	if err := _BurnMintTokenPool.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolLockedIterator struct {
	Event *BurnMintTokenPoolLocked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolLockedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolLocked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(BurnMintTokenPoolLocked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *BurnMintTokenPoolLockedIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolLockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolLocked struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) FilterLocked(opts *bind.FilterOpts, sender []common.Address) (*BurnMintTokenPoolLockedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BurnMintTokenPool.contract.FilterLogs(opts, "Locked", senderRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolLockedIterator{contract: _BurnMintTokenPool.contract, event: "Locked", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) WatchLocked(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolLocked, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BurnMintTokenPool.contract.WatchLogs(opts, "Locked", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolLocked)
				if err := _BurnMintTokenPool.contract.UnpackLog(event, "Locked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) ParseLocked(log types.Log) (*BurnMintTokenPoolLocked, error) {
	event := new(BurnMintTokenPoolLocked)
	if err := _BurnMintTokenPool.contract.UnpackLog(event, "Locked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolMintedIterator struct {
	Event *BurnMintTokenPoolMinted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolMintedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolMinted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(BurnMintTokenPoolMinted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *BurnMintTokenPoolMintedIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolMinted struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) FilterMinted(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*BurnMintTokenPoolMintedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BurnMintTokenPool.contract.FilterLogs(opts, "Minted", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolMintedIterator{contract: _BurnMintTokenPool.contract, event: "Minted", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) WatchMinted(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolMinted, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BurnMintTokenPool.contract.WatchLogs(opts, "Minted", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolMinted)
				if err := _BurnMintTokenPool.contract.UnpackLog(event, "Minted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) ParseMinted(log types.Log) (*BurnMintTokenPoolMinted, error) {
	event := new(BurnMintTokenPoolMinted)
	if err := _BurnMintTokenPool.contract.UnpackLog(event, "Minted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolOwnershipTransferRequestedIterator struct {
	Event *BurnMintTokenPoolOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolOwnershipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(BurnMintTokenPoolOwnershipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *BurnMintTokenPoolOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnMintTokenPoolOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnMintTokenPool.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolOwnershipTransferRequestedIterator{contract: _BurnMintTokenPool.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnMintTokenPool.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolOwnershipTransferRequested)
				if err := _BurnMintTokenPool.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) ParseOwnershipTransferRequested(log types.Log) (*BurnMintTokenPoolOwnershipTransferRequested, error) {
	event := new(BurnMintTokenPoolOwnershipTransferRequested)
	if err := _BurnMintTokenPool.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolOwnershipTransferredIterator struct {
	Event *BurnMintTokenPoolOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(BurnMintTokenPoolOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *BurnMintTokenPoolOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnMintTokenPoolOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnMintTokenPool.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolOwnershipTransferredIterator{contract: _BurnMintTokenPool.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnMintTokenPool.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolOwnershipTransferred)
				if err := _BurnMintTokenPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) ParseOwnershipTransferred(log types.Log) (*BurnMintTokenPoolOwnershipTransferred, error) {
	event := new(BurnMintTokenPoolOwnershipTransferred)
	if err := _BurnMintTokenPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolRateLimitAdminSetIterator struct {
	Event *BurnMintTokenPoolRateLimitAdminSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolRateLimitAdminSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolRateLimitAdminSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(BurnMintTokenPoolRateLimitAdminSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *BurnMintTokenPoolRateLimitAdminSetIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolRateLimitAdminSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolRateLimitAdminSet struct {
	RateLimitAdmin common.Address
	Raw            types.Log
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) FilterRateLimitAdminSet(opts *bind.FilterOpts) (*BurnMintTokenPoolRateLimitAdminSetIterator, error) {

	logs, sub, err := _BurnMintTokenPool.contract.FilterLogs(opts, "RateLimitAdminSet")
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolRateLimitAdminSetIterator{contract: _BurnMintTokenPool.contract, event: "RateLimitAdminSet", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) WatchRateLimitAdminSet(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolRateLimitAdminSet) (event.Subscription, error) {

	logs, sub, err := _BurnMintTokenPool.contract.WatchLogs(opts, "RateLimitAdminSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolRateLimitAdminSet)
				if err := _BurnMintTokenPool.contract.UnpackLog(event, "RateLimitAdminSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) ParseRateLimitAdminSet(log types.Log) (*BurnMintTokenPoolRateLimitAdminSet, error) {
	event := new(BurnMintTokenPoolRateLimitAdminSet)
	if err := _BurnMintTokenPool.contract.UnpackLog(event, "RateLimitAdminSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolReleasedIterator struct {
	Event *BurnMintTokenPoolReleased

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolReleasedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolReleased)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(BurnMintTokenPoolReleased)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *BurnMintTokenPoolReleasedIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolReleasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolReleased struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) FilterReleased(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*BurnMintTokenPoolReleasedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BurnMintTokenPool.contract.FilterLogs(opts, "Released", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolReleasedIterator{contract: _BurnMintTokenPool.contract, event: "Released", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) WatchReleased(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolReleased, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BurnMintTokenPool.contract.WatchLogs(opts, "Released", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolReleased)
				if err := _BurnMintTokenPool.contract.UnpackLog(event, "Released", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) ParseReleased(log types.Log) (*BurnMintTokenPoolReleased, error) {
	event := new(BurnMintTokenPoolReleased)
	if err := _BurnMintTokenPool.contract.UnpackLog(event, "Released", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolRemotePoolAddedIterator struct {
	Event *BurnMintTokenPoolRemotePoolAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolRemotePoolAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolRemotePoolAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(BurnMintTokenPoolRemotePoolAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *BurnMintTokenPoolRemotePoolAddedIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolRemotePoolAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolRemotePoolAdded struct {
	RemoteChainSelector uint64
	RemotePoolAddress   []byte
	Raw                 types.Log
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) FilterRemotePoolAdded(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnMintTokenPoolRemotePoolAddedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnMintTokenPool.contract.FilterLogs(opts, "RemotePoolAdded", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolRemotePoolAddedIterator{contract: _BurnMintTokenPool.contract, event: "RemotePoolAdded", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) WatchRemotePoolAdded(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolRemotePoolAdded, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnMintTokenPool.contract.WatchLogs(opts, "RemotePoolAdded", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolRemotePoolAdded)
				if err := _BurnMintTokenPool.contract.UnpackLog(event, "RemotePoolAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) ParseRemotePoolAdded(log types.Log) (*BurnMintTokenPoolRemotePoolAdded, error) {
	event := new(BurnMintTokenPoolRemotePoolAdded)
	if err := _BurnMintTokenPool.contract.UnpackLog(event, "RemotePoolAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolRemotePoolRemovedIterator struct {
	Event *BurnMintTokenPoolRemotePoolRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolRemotePoolRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolRemotePoolRemoved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(BurnMintTokenPoolRemotePoolRemoved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *BurnMintTokenPoolRemotePoolRemovedIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolRemotePoolRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolRemotePoolRemoved struct {
	RemoteChainSelector uint64
	RemotePoolAddress   []byte
	Raw                 types.Log
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) FilterRemotePoolRemoved(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnMintTokenPoolRemotePoolRemovedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnMintTokenPool.contract.FilterLogs(opts, "RemotePoolRemoved", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolRemotePoolRemovedIterator{contract: _BurnMintTokenPool.contract, event: "RemotePoolRemoved", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) WatchRemotePoolRemoved(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolRemotePoolRemoved, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnMintTokenPool.contract.WatchLogs(opts, "RemotePoolRemoved", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolRemotePoolRemoved)
				if err := _BurnMintTokenPool.contract.UnpackLog(event, "RemotePoolRemoved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) ParseRemotePoolRemoved(log types.Log) (*BurnMintTokenPoolRemotePoolRemoved, error) {
	event := new(BurnMintTokenPoolRemotePoolRemoved)
	if err := _BurnMintTokenPool.contract.UnpackLog(event, "RemotePoolRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolRouterUpdatedIterator struct {
	Event *BurnMintTokenPoolRouterUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolRouterUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolRouterUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(BurnMintTokenPoolRouterUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *BurnMintTokenPoolRouterUpdatedIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolRouterUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolRouterUpdated struct {
	OldRouter common.Address
	NewRouter common.Address
	Raw       types.Log
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) FilterRouterUpdated(opts *bind.FilterOpts) (*BurnMintTokenPoolRouterUpdatedIterator, error) {

	logs, sub, err := _BurnMintTokenPool.contract.FilterLogs(opts, "RouterUpdated")
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolRouterUpdatedIterator{contract: _BurnMintTokenPool.contract, event: "RouterUpdated", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolRouterUpdated) (event.Subscription, error) {

	logs, sub, err := _BurnMintTokenPool.contract.WatchLogs(opts, "RouterUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolRouterUpdated)
				if err := _BurnMintTokenPool.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) ParseRouterUpdated(log types.Log) (*BurnMintTokenPoolRouterUpdated, error) {
	event := new(BurnMintTokenPoolRouterUpdated)
	if err := _BurnMintTokenPool.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnMintTokenPoolTokensConsumedIterator struct {
	Event *BurnMintTokenPoolTokensConsumed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnMintTokenPoolTokensConsumedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnMintTokenPoolTokensConsumed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(BurnMintTokenPoolTokensConsumed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *BurnMintTokenPoolTokensConsumedIterator) Error() error {
	return it.fail
}

func (it *BurnMintTokenPoolTokensConsumedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnMintTokenPoolTokensConsumed struct {
	Tokens *big.Int
	Raw    types.Log
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) FilterTokensConsumed(opts *bind.FilterOpts) (*BurnMintTokenPoolTokensConsumedIterator, error) {

	logs, sub, err := _BurnMintTokenPool.contract.FilterLogs(opts, "TokensConsumed")
	if err != nil {
		return nil, err
	}
	return &BurnMintTokenPoolTokensConsumedIterator{contract: _BurnMintTokenPool.contract, event: "TokensConsumed", logs: logs, sub: sub}, nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) WatchTokensConsumed(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolTokensConsumed) (event.Subscription, error) {

	logs, sub, err := _BurnMintTokenPool.contract.WatchLogs(opts, "TokensConsumed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnMintTokenPoolTokensConsumed)
				if err := _BurnMintTokenPool.contract.UnpackLog(event, "TokensConsumed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_BurnMintTokenPool *BurnMintTokenPoolFilterer) ParseTokensConsumed(log types.Log) (*BurnMintTokenPoolTokensConsumed, error) {
	event := new(BurnMintTokenPoolTokensConsumed)
	if err := _BurnMintTokenPool.contract.UnpackLog(event, "TokensConsumed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_BurnMintTokenPool *BurnMintTokenPool) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _BurnMintTokenPool.abi.Events["AllowListAdd"].ID:
		return _BurnMintTokenPool.ParseAllowListAdd(log)
	case _BurnMintTokenPool.abi.Events["AllowListRemove"].ID:
		return _BurnMintTokenPool.ParseAllowListRemove(log)
	case _BurnMintTokenPool.abi.Events["Burned"].ID:
		return _BurnMintTokenPool.ParseBurned(log)
	case _BurnMintTokenPool.abi.Events["ChainAdded"].ID:
		return _BurnMintTokenPool.ParseChainAdded(log)
	case _BurnMintTokenPool.abi.Events["ChainConfigured"].ID:
		return _BurnMintTokenPool.ParseChainConfigured(log)
	case _BurnMintTokenPool.abi.Events["ChainRemoved"].ID:
		return _BurnMintTokenPool.ParseChainRemoved(log)
	case _BurnMintTokenPool.abi.Events["ConfigChanged"].ID:
		return _BurnMintTokenPool.ParseConfigChanged(log)
	case _BurnMintTokenPool.abi.Events["Locked"].ID:
		return _BurnMintTokenPool.ParseLocked(log)
	case _BurnMintTokenPool.abi.Events["Minted"].ID:
		return _BurnMintTokenPool.ParseMinted(log)
	case _BurnMintTokenPool.abi.Events["OwnershipTransferRequested"].ID:
		return _BurnMintTokenPool.ParseOwnershipTransferRequested(log)
	case _BurnMintTokenPool.abi.Events["OwnershipTransferred"].ID:
		return _BurnMintTokenPool.ParseOwnershipTransferred(log)
	case _BurnMintTokenPool.abi.Events["RateLimitAdminSet"].ID:
		return _BurnMintTokenPool.ParseRateLimitAdminSet(log)
	case _BurnMintTokenPool.abi.Events["Released"].ID:
		return _BurnMintTokenPool.ParseReleased(log)
	case _BurnMintTokenPool.abi.Events["RemotePoolAdded"].ID:
		return _BurnMintTokenPool.ParseRemotePoolAdded(log)
	case _BurnMintTokenPool.abi.Events["RemotePoolRemoved"].ID:
		return _BurnMintTokenPool.ParseRemotePoolRemoved(log)
	case _BurnMintTokenPool.abi.Events["RouterUpdated"].ID:
		return _BurnMintTokenPool.ParseRouterUpdated(log)
	case _BurnMintTokenPool.abi.Events["TokensConsumed"].ID:
		return _BurnMintTokenPool.ParseTokensConsumed(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (BurnMintTokenPoolAllowListAdd) Topic() common.Hash {
	return common.HexToHash("0x2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d8")
}

func (BurnMintTokenPoolAllowListRemove) Topic() common.Hash {
	return common.HexToHash("0x800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf7566")
}

func (BurnMintTokenPoolBurned) Topic() common.Hash {
	return common.HexToHash("0x696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df7")
}

func (BurnMintTokenPoolChainAdded) Topic() common.Hash {
	return common.HexToHash("0x8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c2")
}

func (BurnMintTokenPoolChainConfigured) Topic() common.Hash {
	return common.HexToHash("0x0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b")
}

func (BurnMintTokenPoolChainRemoved) Topic() common.Hash {
	return common.HexToHash("0x5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d859916")
}

func (BurnMintTokenPoolConfigChanged) Topic() common.Hash {
	return common.HexToHash("0x9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c19")
}

func (BurnMintTokenPoolLocked) Topic() common.Hash {
	return common.HexToHash("0x9f1ec8c880f76798e7b793325d625e9b60e4082a553c98f42b6cda368dd60008")
}

func (BurnMintTokenPoolMinted) Topic() common.Hash {
	return common.HexToHash("0x9d228d69b5fdb8d273a2336f8fb8612d039631024ea9bf09c424a9503aa078f0")
}

func (BurnMintTokenPoolOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (BurnMintTokenPoolOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (BurnMintTokenPoolRateLimitAdminSet) Topic() common.Hash {
	return common.HexToHash("0x44676b5284b809a22248eba0da87391d79098be38bb03154be88a58bf4d09174")
}

func (BurnMintTokenPoolReleased) Topic() common.Hash {
	return common.HexToHash("0x2d87480f50083e2b2759522a8fdda59802650a8055e609a7772cf70c07748f52")
}

func (BurnMintTokenPoolRemotePoolAdded) Topic() common.Hash {
	return common.HexToHash("0x7d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea")
}

func (BurnMintTokenPoolRemotePoolRemoved) Topic() common.Hash {
	return common.HexToHash("0x52d00ee4d9bd51b40168f2afc5848837288ce258784ad914278791464b3f4d76")
}

func (BurnMintTokenPoolRouterUpdated) Topic() common.Hash {
	return common.HexToHash("0x02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f1684")
}

func (BurnMintTokenPoolTokensConsumed) Topic() common.Hash {
	return common.HexToHash("0x1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a")
}

func (_BurnMintTokenPool *BurnMintTokenPool) Address() common.Address {
	return _BurnMintTokenPool.address
}

type BurnMintTokenPoolInterface interface {
	GetAllowList(opts *bind.CallOpts) ([]common.Address, error)

	GetAllowListEnabled(opts *bind.CallOpts) (bool, error)

	GetCurrentInboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error)

	GetCurrentOutboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error)

	GetRateLimitAdmin(opts *bind.CallOpts) (common.Address, error)

	GetRemotePools(opts *bind.CallOpts, remoteChainSelector uint64) ([][]byte, error)

	GetRemoteToken(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error)

	GetRmnProxy(opts *bind.CallOpts) (common.Address, error)

	GetRouter(opts *bind.CallOpts) (common.Address, error)

	GetSupportedChains(opts *bind.CallOpts) ([]uint64, error)

	GetToken(opts *bind.CallOpts) (common.Address, error)

	GetTokenDecimals(opts *bind.CallOpts) (uint8, error)

	IsRemotePool(opts *bind.CallOpts, remoteChainSelector uint64, remotePoolAddress []byte) (bool, error)

	IsSupportedChain(opts *bind.CallOpts, remoteChainSelector uint64) (bool, error)

	IsSupportedToken(opts *bind.CallOpts, token common.Address) (bool, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	AddRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error)

	ApplyAllowListUpdates(opts *bind.TransactOpts, removes []common.Address, adds []common.Address) (*types.Transaction, error)

	ApplyChainUpdates(opts *bind.TransactOpts, remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error)

	LockOrBurn(opts *bind.TransactOpts, lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error)

	ReleaseOrMint(opts *bind.TransactOpts, releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error)

	RemoveRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error)

	SetChainRateLimiterConfig(opts *bind.TransactOpts, remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error)

	SetRateLimitAdmin(opts *bind.TransactOpts, rateLimitAdmin common.Address) (*types.Transaction, error)

	SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterAllowListAdd(opts *bind.FilterOpts) (*BurnMintTokenPoolAllowListAddIterator, error)

	WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAllowListAdd) (event.Subscription, error)

	ParseAllowListAdd(log types.Log) (*BurnMintTokenPoolAllowListAdd, error)

	FilterAllowListRemove(opts *bind.FilterOpts) (*BurnMintTokenPoolAllowListRemoveIterator, error)

	WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolAllowListRemove) (event.Subscription, error)

	ParseAllowListRemove(log types.Log) (*BurnMintTokenPoolAllowListRemove, error)

	FilterBurned(opts *bind.FilterOpts, sender []common.Address) (*BurnMintTokenPoolBurnedIterator, error)

	WatchBurned(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolBurned, sender []common.Address) (event.Subscription, error)

	ParseBurned(log types.Log) (*BurnMintTokenPoolBurned, error)

	FilterChainAdded(opts *bind.FilterOpts) (*BurnMintTokenPoolChainAddedIterator, error)

	WatchChainAdded(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolChainAdded) (event.Subscription, error)

	ParseChainAdded(log types.Log) (*BurnMintTokenPoolChainAdded, error)

	FilterChainConfigured(opts *bind.FilterOpts) (*BurnMintTokenPoolChainConfiguredIterator, error)

	WatchChainConfigured(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolChainConfigured) (event.Subscription, error)

	ParseChainConfigured(log types.Log) (*BurnMintTokenPoolChainConfigured, error)

	FilterChainRemoved(opts *bind.FilterOpts) (*BurnMintTokenPoolChainRemovedIterator, error)

	WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolChainRemoved) (event.Subscription, error)

	ParseChainRemoved(log types.Log) (*BurnMintTokenPoolChainRemoved, error)

	FilterConfigChanged(opts *bind.FilterOpts) (*BurnMintTokenPoolConfigChangedIterator, error)

	WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolConfigChanged) (event.Subscription, error)

	ParseConfigChanged(log types.Log) (*BurnMintTokenPoolConfigChanged, error)

	FilterLocked(opts *bind.FilterOpts, sender []common.Address) (*BurnMintTokenPoolLockedIterator, error)

	WatchLocked(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolLocked, sender []common.Address) (event.Subscription, error)

	ParseLocked(log types.Log) (*BurnMintTokenPoolLocked, error)

	FilterMinted(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*BurnMintTokenPoolMintedIterator, error)

	WatchMinted(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolMinted, sender []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseMinted(log types.Log) (*BurnMintTokenPoolMinted, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnMintTokenPoolOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*BurnMintTokenPoolOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnMintTokenPoolOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*BurnMintTokenPoolOwnershipTransferred, error)

	FilterRateLimitAdminSet(opts *bind.FilterOpts) (*BurnMintTokenPoolRateLimitAdminSetIterator, error)

	WatchRateLimitAdminSet(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolRateLimitAdminSet) (event.Subscription, error)

	ParseRateLimitAdminSet(log types.Log) (*BurnMintTokenPoolRateLimitAdminSet, error)

	FilterReleased(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*BurnMintTokenPoolReleasedIterator, error)

	WatchReleased(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolReleased, sender []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseReleased(log types.Log) (*BurnMintTokenPoolReleased, error)

	FilterRemotePoolAdded(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnMintTokenPoolRemotePoolAddedIterator, error)

	WatchRemotePoolAdded(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolRemotePoolAdded, remoteChainSelector []uint64) (event.Subscription, error)

	ParseRemotePoolAdded(log types.Log) (*BurnMintTokenPoolRemotePoolAdded, error)

	FilterRemotePoolRemoved(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnMintTokenPoolRemotePoolRemovedIterator, error)

	WatchRemotePoolRemoved(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolRemotePoolRemoved, remoteChainSelector []uint64) (event.Subscription, error)

	ParseRemotePoolRemoved(log types.Log) (*BurnMintTokenPoolRemotePoolRemoved, error)

	FilterRouterUpdated(opts *bind.FilterOpts) (*BurnMintTokenPoolRouterUpdatedIterator, error)

	WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolRouterUpdated) (event.Subscription, error)

	ParseRouterUpdated(log types.Log) (*BurnMintTokenPoolRouterUpdated, error)

	FilterTokensConsumed(opts *bind.FilterOpts) (*BurnMintTokenPoolTokensConsumedIterator, error)

	WatchTokensConsumed(opts *bind.WatchOpts, sink chan<- *BurnMintTokenPoolTokensConsumed) (event.Subscription, error)

	ParseTokensConsumed(log types.Log) (*BurnMintTokenPoolTokensConsumed, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}

var BurnMintTokenPoolZKBin string = ("0x0003000000000002001f000000000002000000600310027000020000000103550000061c0030019d0000061c0330019700000001002001900000005c0000c13d0000008004000039000000400040043f000000040030008c000000840000413d000000000201043b000000e0022002700000062f0020009c000000860000a13d000006300020009c000000c70000a13d000006310020009c000000f60000213d000006370020009c000001800000213d0000063a0020009c000005220000613d0000063b0020009c000000840000c13d0000000002000416000000000002004b000000840000c13d0000000504000039000000000204041a000000800020043f000000000040043f000000000002004b00000a100000c13d000000a002000039000000400020043f0000002004000039000000000500001900000005065002100000003f076000390000067a077001970000000007270019000006200070009c000001220000213d000000400070043f00000000005204350000001f0560018f000000a004400039000000000006004b0000003a0000613d000000000131034f00000000036400190000000006040019000000001701043c0000000006760436000000000036004b000000360000c13d000000000005004b000000800100043d000000000001004b0000004c0000613d00000000010000190000000003020433000000000013004b00000f830000a13d00000005031002100000000005430019000000a0033000390000000003030433000006200330019700000000003504350000000101100039000000800300043d000000000031004b0000003f0000413d000000400100043d00000020030000390000000005310436000000000302043300000000003504350000004002100039000000000003004b00000a070000613d00000000050000190000000046040434000006200660019700000000026204360000000105500039000000000035004b000000550000413d00000a070000013d0000010004000039000000400040043f0000000002000416000000000002004b000000840000c13d0000001f023000390000061d022001970000010002200039000000400020043f0000001f0530018f0000061e0630019800000100026000390000006e0000613d000000000701034f000000007807043c0000000004840436000000000024004b0000006a0000c13d000000000005004b0000007b0000613d000000000161034f0000000304500210000000000502043300000000054501cf000000000545022f000000000101043b0000010004400089000000000141022f00000000014101cf000000000151019f0000000000120435000000a00030008c000000840000413d000001000100043d0000061f0010009c000000840000213d000001200200043d001100000002001d000000ff0020008c0000010f0000a13d00000000010000190000186c00010430000006450020009c000000a10000213d0000064f0020009c000001280000a13d000006500020009c000001570000213d000006530020009c000001f90000613d000006540020009c000000840000c13d0000000001000416000000000001004b000000840000c13d0000000001000412001900000001001d001800200000003d000080050100003900000044030000390000000004000415000000190440008a00000005044002100000067102000041186a18420000040f000000ff0110018f000000800010043f00000672010000410000186b0001042e000006460020009c000001390000a13d000006470020009c000001620000213d0000064a0020009c000002140000613d0000064b0020009c000000840000c13d000000240030008c000000840000413d0000000002000416000000000002004b000000840000c13d0000000401100370000000000101043b0000061f0010009c000000840000213d0000000102000039000000000202041a0000061f022001970000000003000411000000000023004b000009e00000c13d0000000902000039000000000302041a0000062303300197000000000313019f000000000032041b000000800010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000691011001c70000800d020000390000000103000039000006920400004100000a340000013d0000063c0020009c000001440000a13d0000063d0020009c0000016b0000213d000006400020009c000002970000613d000006410020009c000000840000c13d000000240030008c000000840000413d0000000002000416000000000002004b000000840000c13d0000000401100370000000000101043b001100000001001d000006200010009c000000840000213d186a14670000040f0000001101000029000000000010043f0000000701000039000000200010043f00000040020000390000000001000019186a182d0000040f001000000001001d000000400100043d001100000001001d186a13220000040f00000010050000290000000201500039000000000401041a00000677004001980000000002000039000000010200c03900000011010000290000004003100039000000000023043500000080024002700000061c02200197000000200310003900000000002304350000066302400197000000000021043500000003025000390000056d0000013d000006320020009c000001e00000213d000006350020009c000005370000613d000006360020009c000000840000c13d0000000001000416000000000001004b000000840000c13d0000000001000412001300000001001d001200600000003d000080050100003900000044030000390000000004000415000000130440008a00000005044002100000067102000041186a18420000040f000000000001004b0000000001000039000000010100c039000000800010043f00000672010000410000186b0001042e000001400200043d000006200020009c000000840000213d0000001f04200039000000000034004b000000000500001900000621050080410000062104400197000000000004004b00000000060000190000062106004041000006210040009c000000000605c019000000000006004b000000840000c13d00000100042000390000000004040433000006200040009c000009b40000a13d0000068201000041000000000010043f0000004101000039000000040010043f0000065f010000410000186c00010430000006550020009c000003fa0000613d000006560020009c000003400000613d000006570020009c000000840000c13d0000000001000416000000000001004b000000840000c13d0000000001000412001d00000001001d001c00000000003d0000800501000039000000440300003900000000040004150000001d0440008a000005410000013d0000064c0020009c0000040d0000613d0000064d0020009c000003540000613d0000064e0020009c000000840000c13d0000000001000416000000000001004b000000840000c13d00000009010000390000033b0000013d000006420020009c000004cf0000613d000006430020009c000003920000613d000006440020009c000000840000c13d0000000001000416000000000001004b000000840000c13d0000000202000039000000000102041a000000800010043f000000000020043f0000002002000039000000000001004b000009e80000c13d000000a0010000390000000004020019000009f70000013d000006510020009c000002300000613d000006520020009c000000840000c13d0000000001000416000000000001004b000000840000c13d0000000001030019186a133f0000040f186a13d70000040f0000028d0000013d000006480020009c000002830000613d000006490020009c000000840000c13d0000000001000416000000000001004b000000840000c13d00000001010000390000033b0000013d0000063e0020009c000003370000613d0000063f0020009c000000840000c13d000000240030008c000000840000413d0000000002000416000000000002004b000000840000c13d0000000401100370000000000101043b000006200010009c000000840000213d186a147c0000040f0000002002000039000000400300043d001100000003001d0000000002230436186a132d0000040f00000011020000290000057a0000013d000006380020009c000005480000613d000006390020009c000000840000c13d000000e40030008c000000840000413d0000000002000416000000000002004b000000840000c13d0000000402100370000000000202043b001100000002001d000006200020009c000000840000213d000000e002000039000000400020043f0000002402100370000000000202043b000000000002004b0000000003000039000000010300c039000000000032004b000000840000c13d000000800020043f0000004402100370000000000202043b000006630020009c000000840000213d000000a00020043f0000006402100370000000000202043b000006630020009c000000840000213d000000c00020043f0000014002000039000000400020043f0000008402100370000000000202043b000000000002004b0000000003000039000000010300c039000000000032004b000000840000c13d000000e00020043f000000a402100370000000000202043b000006630020009c000000840000213d000001000020043f000000c401100370000000000101043b000006630010009c000000840000213d000001200010043f0000000901000039000000000101041a0000061f021001970000000001000411000000000021004b000001c10000613d0000000102000039000000000202041a0000061f02200197000000000021004b00000dfb0000c13d0000001101000029000000000010043f0000000601000039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000100200190000000840000613d000000000101043b000000000101041a000000000001004b000003870000613d000000c00100043d0000066301100197000000800200043d000000000002004b00000e630000c13d000000000001004b000001dc0000c13d000000a00100043d000006630010019800000e690000613d000000400200043d001100000002001d000006660100004100000e8a0000013d000006330020009c000005830000613d000006340020009c000000840000c13d000000240030008c000000840000413d0000000002000416000000000002004b000000840000c13d0000000401100370000000000601043b0000061f0060009c000000840000213d0000000101000039000000000101041a0000061f011001970000000005000411000000000015004b000009e00000c13d000000000056004b00000a380000c13d0000065a01000041000000800010043f0000065b010000410000186c00010430000000240030008c000000840000413d0000000002000416000000000002004b000000840000c13d0000000401100370000000000101043b001100000001001d0000061f0010009c000000840000213d0000000001000412001b00000001001d001a00000000003d0000800501000039000000440300003900000000040004150000001b0440008a00000005044002100000067102000041186a18420000040f0000061f01100197000000110010006b00000000010000390000000101006039000000800010043f00000672010000410000186b0001042e0000000001000416000000000001004b000000840000c13d000000000100041a0000061f021001970000000006000411000000000026004b000009e40000c13d0000000102000039000000000302041a0000062304300197000000000464019f000000000042041b0000062301100197000000000010041b00000000010004140000061f053001970000061c0010009c0000061c01008041000000c00110021000000658011001c70000800d0200003900000003030000390000069404000041186a18600000040f0000000100200190000000840000613d00000a470000013d000000240030008c000000840000413d0000000002000416000000000002004b000000840000c13d0000000402100370000000000202043b001100000002001d000006200020009c000000840000213d000000110230006a000006600020009c000000840000213d000001040020008c000000840000413d000000a002000039000000400020043f0000001102000029000f00840020003d0000000f01100360000000800000043f000000000101043b001000000001001d0000061f0010009c000000840000213d0000067101000041000000000010044300000000010004120000000400100443000000240000044300000000010004140000061c0010009c0000061c01008041000000c00110021000000683011001c70000800502000039186a18650000040f0000000100200190000012990000613d0000000202000367000000000101043b0000061f01100197000000100010006b00000a670000c13d0000000f01000029000e0060001000920000000e01200360000000000101043b000006200010009c000000840000213d000000400300043d0000068502000041000000000023043500000080011002100000068601100197000f00000003001d0000000402300039000000000012043500000671010000410000000000100443000000000100041200000004001004430000004001000039000000240010044300000000010004140000061c0010009c0000061c01008041000000c00110021000000683011001c70000800502000039186a18650000040f0000000100200190000012990000613d000000000201043b00000000010004140000061f02200197000000040020008c00000bdd0000c13d0000000103000031000000200030008c0000002004000039000000000403401900000c060000013d000000240030008c000000840000413d0000000002000416000000000002004b000000840000c13d0000000401100370000000000101043b000006200010009c000000840000213d186a167a0000040f000000000001004b0000000001000039000000010100c039000000400200043d00000000001204350000061c0020009c0000061c02008041000000400120021000000690011001c70000186b0001042e000000440030008c000000840000413d0000000002000416000000000002004b000000840000c13d0000000402100370000000000202043b001100000002001d000006200020009c000000840000213d0000002402100370000000000202043b000006200020009c000000840000213d0000002304200039000000000034004b000000840000813d000f00040020003d0000000f01100360000000000101043b001000000001001d000006200010009c000000840000213d0000002402200039000d00000002001d000e00100020002d0000000e0030006b000000840000213d0000000101000039000000000101041a0000061f011001970000000002000411000000000012004b000009e00000c13d0000001101000029000000000010043f0000000601000039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000100200190000000840000613d000000000101043b000000000101041a000000000001004b000003870000613d0000001101000029000000000010043f0000000701000039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000100200190000000840000613d00000010020000290000001f02200039000006a502200197000b00000002001d0000003f02200039000006a502200197000000000101043b000c00000001001d000000400100043d0000000002210019000000000012004b00000000040000390000000104004039000006200020009c000001220000213d0000000100400190000001220000c13d000000400020043f000000100200002900000000022104360000000e05000029000000000050007c000000840000213d0000001004000029000006a503400198000e001f00400193000a00000003001d00000000033200190000000f040000290000002004400039000f00000004001d0000000204400367000002fe0000613d000000000504034f0000000006020019000000005705043c0000000006760436000000000036004b000002fa0000c13d0000000e0000006b0000030c0000613d0000000a044003600000000e050000290000000305500210000000000603043300000000065601cf000000000656022f000000000404043b0000010005500089000000000454022f00000000045401cf000000000464019f0000000000430435000000100320002900000000000304350000061c0020009c0000061c02008041000000400220021000000000010104330000061c0010009c0000061c010080410000006001100210000000000121019f00000000020004140000061c0020009c0000061c02008041000000c002200210000000000112019f00000658011001c70000801002000039186a18650000040f0000000100200190000000840000613d0000000c020000290000000503200039000000000201043b0000000001030019186a168e0000040f000000400700043d000000000001004b00000e940000c13d000c00000007001d0000002401700039000000400200003900000000002104350000068001000041000000000017043500000004017000390000001102000029000000000021043500000044037000390000000d010000290000001002000029186a14480000040f0000000c0200002900000f9d0000013d0000000001000416000000000001004b000000840000c13d0000000401000039000000000101041a0000061f01100197000000800010043f00000672010000410000186b0001042e0000000001000416000000000001004b000000840000c13d000000c001000039000000400010043f0000001701000039000000800010043f0000069f01000041000000a00010043f0000002001000039000000c00010043f0000008001000039000000e002000039186a132d0000040f000000c00110008a0000061c0010009c0000061c010080410000006001100210000006a0011001c70000186b0001042e000000440030008c000000840000413d0000000002000416000000000002004b000000840000c13d0000000402100370000000000202043b001100000002001d000006200020009c000000840000213d0000002402100370000000000202043b000006200020009c000000840000213d0000002304200039000000000034004b000000840000813d0000000404200039000000000141034f000000000101043b001000000001001d000006200010009c000000840000213d0000002402200039000f00000002001d0000001001200029000000000031004b000000840000213d0000000101000039000000000101041a0000061f011001970000000002000411000000000012004b000009e00000c13d0000001101000029000000000010043f0000000601000039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000100200190000000840000613d000000000101043b000000000101041a000000000001004b00000b530000c13d000000400100043d0000065e0200004100000000002104350000000402100039000000110300002900000000003204350000061c0010009c0000061c0100804100000040011002100000065f011001c70000186c00010430000000240030008c000000840000413d0000000002000416000000000002004b000000840000c13d0000000401100370000000000101043b000006200010009c000000840000213d000000000010043f0000000701000039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000100200190000000840000613d000000000101043b0000000501100039000000000301041a000000400200043d000f00000002001d001100000003001d0000000002320436000e00000002001d000000000010043f00000000010004140000061c0010009c0000061c01008041000000c0011002100000062b011001c70000801002000039186a18650000040f0000000100200190000000840000613d0000001105000029000000000005004b0000000e02000029000003c60000613d000000000101043b0000000e020000290000000003000019000000000401041a000000000242043600000001011000390000000103300039000000000053004b000003c00000413d0000000f0120006a0000001f01100039000006a5011001970000000f04100029000000000014004b00000000010000390000000101004039000006200040009c000001220000213d0000000100100190000001220000c13d000000400040043f0000000f010000290000000002010433000006200020009c000001220000213d00000005012002100000003f0310003900000622033001970000000003430019000006200030009c000001220000213d000000400030043f000d00000004001d0000000005240436000000000002004b000003e80000613d00000060020000390000000003000019000000000435001900000000002404350000002003300039000000000013004b000003e30000413d000c00000005001d0000000f010000290000000001010433000000000001004b00000b5c0000c13d000000400100043d000000200200003900000000032104360000000d0200002900000000020204330000000000230435000000400310003900000005042002100000000005340019000000000002004b00000bc00000c13d000000000215004900000a080000013d000000240030008c000000840000413d0000000002000416000000000002004b000000840000c13d0000000401100370000000000201043b000006a100200198000000840000c13d0000000101000039000006a20020009c000005450000613d000006a30020009c000005450000613d000006a40020009c000000000100c019000000800010043f00000672010000410000186b0001042e000000440030008c000000840000413d0000000002000416000000000002004b000000840000c13d0000000402100370000000000202043b000006200020009c000000840000213d0000002305200039000000000035004b000000840000813d0000000405200039000000000551034f000000000905043b000006200090009c000000840000213d0000002407200039000000050b90021000000000087b0019000000000038004b000000840000213d0000002402100370000000000202043b000006200020009c000000840000213d0000002305200039000000000035004b000000840000813d0000000405200039000000000551034f000000000605043b000006200060009c000000840000213d0000002402200039000000050a60021000000000052a0019000000000035004b000000840000213d0000000103000039000000000303041a0000061f03300197000000000c00041100000000003c004b000009e00000c13d0000003f03b000390000062203300197000006790030009c000001220000213d0000008003300039000f00000003001d000000400030043f000000800090043f000000000009004b0000044f0000613d000000000371034f000000000303043b0000061f0030009c000000840000213d000000200440003900000000003404350000002007700039000000000087004b000004440000413d000000400300043d000f00000003001d0000003f03a0003900000622033001970000000f033000290000000f0030006c00000000040000390000000104004039000006200030009c000001220000213d0000000100400190000001220000c13d000000400030043f0000000f030000290000000003630436000e00000003001d000000000006004b000004690000613d0000000f03000029000000000421034f000000000404043b0000061f0040009c000000840000213d000000200330003900000000004304350000002002200039000000000052004b000004600000413d00000671010000410000000000100443000000000100041200000004001004430000006001000039000000240010044300000000010004140000061c0010009c0000061c01008041000000c00110021000000683011001c70000800502000039186a18650000040f0000000100200190000012990000613d000000000101043b000000000001004b00000ad90000613d000000800100043d000000000001004b00000f1e0000c13d0000000f010000290000000001010433000000000001004b00000a470000613d0000000003000019000004890000013d00000001033000390000000f010000290000000001010433000000000013004b00000a470000813d00000005013002100000000e0110002900000000010104330000061f04100198000004840000613d000000000040043f0000000301000039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039001000000003001d001100000004001d186a18650000040f000000110400002900000010030000290000000100200190000000840000613d000000000101043b000000000101041a000000000001004b000004840000c13d0000000203000039000000000103041a000006200010009c000001220000213d0000000102100039000000000023041b0000062a0110009a000000000041041b000000000103041a000d00000001001d000000000040043f0000000301000039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f00000011030000290000000100200190000000840000613d000000000101043b0000000d02000029000000000021041b000000400100043d00000000003104350000061c0010009c0000061c01008041000000400110021000000000020004140000061c0020009c0000061c02008041000000c002200210000000000112019f0000062b011001c70000800d0200003900000001030000390000062c04000041186a18600000040f00000010030000290000000100200190000004840000c13d000000840000013d000000240030008c000000840000413d0000000002000416000000000002004b000000840000c13d0000000402100370000000000202043b000006200020009c000000840000213d0000000003230049000006600030009c000000840000213d000000a40030008c000000840000413d000000c003000039000000400030043f0000006003000039000000800030043f000000a00030043f001000840020003d0000001001100360000000000101043b001100000001001d0000061f0010009c000000840000213d0000067101000041000000000010044300000000010004120000000400100443000000240000044300000000010004140000061c0010009c0000061c01008041000000c00110021000000683011001c70000800502000039186a18650000040f0000000100200190000012990000613d0000000202000367000000000101043b0000061f01100197000000110010006b00000a490000c13d0000001001000029000f0060001000920000000f01200360000000000101043b000006200010009c000000840000213d000000400300043d0000068502000041000000000023043500000080011002100000068601100197001000000003001d0000000402300039000000000012043500000671010000410000000000100443000000000100041200000004001004430000004001000039000000240010044300000000010004140000061c0010009c0000061c01008041000000c00110021000000683011001c70000800502000039186a18650000040f0000000100200190000012990000613d000000000201043b00000000010004140000061f02200197000000040020008c00000adc0000c13d0000000103000031000000200030008c0000002004000039000000000403401900000b050000013d000000240030008c000000840000413d0000000002000416000000000002004b000000840000c13d0000000401100370000000000101043b0000061f0010009c000000840000213d0000000102000039000000000202041a0000061f022001970000000003000411000000000023004b000009e00000c13d000000000001004b00000a240000c13d0000067001000041000000800010043f0000065b010000410000186c000104300000000001000416000000000001004b000000840000c13d0000000001000412001500000001001d001400400000003d000080050100003900000044030000390000000004000415000000150440008a00000005044002100000067102000041186a18420000040f0000061f01100197000000800010043f00000672010000410000186b0001042e000000240030008c000000840000413d0000000002000416000000000002004b000000840000c13d0000000401100370000000000101043b001100000001001d000006200010009c000000840000213d186a14670000040f0000001101000029000000000010043f0000000701000039000000200010043f00000040020000390000000001000019186a182d0000040f001000000001001d000000400100043d001100000001001d186a13220000040f0000001005000029000000000405041a00000677004001980000000002000039000000010200c03900000011010000290000004003100039000000000023043500000080024002700000061c0220019700000020031000390000000000230435000006630240019700000000002104350000000102500039000000000402041a000000800210003900000080034002700000000000320435000006630340019700000060021000390000000000320435186a162c0000040f000000400100043d001000000001001d0000001102000029186a135c0000040f000000100200002900000000012100490000061c0010009c0000061c0100804100000060011002100000061c0020009c0000061c020080410000004002200210000000000121019f0000186b0001042e000000440030008c000000840000413d0000000002000416000000000002004b000000840000c13d0000000402100370000000000202043b000006200020009c000000840000213d0000002304200039000000000034004b000000840000813d0000000404200039000000000441034f000000000404043b000600000004001d000006200040009c000000840000213d000500240020003d000000060200002900000005022002100000000502200029000000000032004b000000840000213d0000002402100370000000000202043b000200000002001d000006200020009c000000840000213d00000002020000290000002302200039000000000032004b000000840000813d00000002020000290000000402200039000000000121034f000000000101043b000100000001001d000006200010009c000000840000213d0000000201000029000300240010003d000000010100002900000005011002100000000301100029000000000031004b000000840000213d0000000101000039000000000101041a0000061f011001970000000002000411000000000012004b000009e00000c13d000000060000006b00000c4c0000c13d000000010000006b00000a470000613d000500000000001d0000000501000029000000050110021000000003011000290000000202000367000000000112034f000000000101043b0000000003000031000000020430006a000001430440008a00000621054001970000062106100197000000000756013f000000000056004b00000000050000190000062105004041000000000041004b00000000040000190000062104008041000006210070009c000000000504c019000000000005004b000000840000c13d001000030010002d000000100130006a000f00000001001d000006600010009c000000840000213d0000000f01000029000001200010008c000000840000413d000000400100043d000900000001001d000006610010009c000001220000213d0000000901000029000000a001100039000000400010043f0000001001200360000000000101043b000006200010009c000000840000213d00000009040000290000000001140436000800000001001d00000010010000290000002001100039000000000112034f000000000101043b000006200010009c000000840000213d0000001001100029001100000001001d0000001f01100039000000000031004b0000000004000019000006210400804100000621011001970000062105300197000000000751013f000000000051004b00000000010000190000062101004041000006210070009c000000000104c019000000000001004b000000840000c13d0000001101200360000000000101043b000006200010009c000001220000213d00000005091002100000003f049000390000062204400197000000400600043d0000000004460019000e00000006001d000000000064004b00000000070000390000000107004039000006200040009c000001220000213d0000000100700190000001220000c13d000000400040043f0000000e040000290000000000140435000000110100002900000020081000390000000009980019000000000039004b000000840000213d000000000098004b000006660000813d0000000e0a000029000006230000013d000000200aa000390000000001b7001900000000000104350000000000ca04350000002008800039000000000098004b000006660000813d000000000182034f000000000101043b000006200010009c000000840000213d000000110d1000290000003f01d00039000000000031004b000000000400001900000621040080410000062101100197000000000751013f000000000051004b00000000010000190000062101004041000006210070009c000000000104c019000000000001004b000000840000c13d000000200ed000390000000001e2034f000000000b01043b0000062000b0009c000001220000213d0000001f01b00039000006a5011001970000003f01100039000006a501100197000000400c00043d00000000011c00190000000000c1004b00000000040000390000000104004039000006200010009c000001220000213d0000000100400190000001220000c13d0000004004d00039000000400010043f0000000007bc043600000000014b0019000000000031004b000000840000213d0000002001e00039000000000412034f000006a501b00198000000000e170019000006580000613d000000000f04034f000000000d07001900000000f60f043c000000000d6d04360000000000ed004b000006540000c13d0000001f0db001900000061c0000613d000000000114034f0000000304d0021000000000060e043300000000064601cf000000000646022f000000000101043b0000010004400089000000000141022f00000000014101cf000000000161019f00000000001e04350000061c0000013d00000008010000290000000e04000029000000000041043500000010010000290000004006100039000000000162034f000000000101043b000006200010009c000000840000213d00000010071000290000001f01700039000000000031004b000000000400001900000621040080410000062101100197000000000851013f000000000051004b00000000010000190000062101004041000006210080009c000000000104c019000000000001004b000000840000c13d000000000172034f000000000401043b000006200040009c000001220000213d0000001f01400039000006a5011001970000003f01100039000006a501100197000000400500043d0000000001150019000000000051004b00000000080000390000000108004039000006200010009c000001220000213d0000000100800190000001220000c13d0000002008700039000000400010043f00000000074504360000000001840019000000000031004b000000840000213d000000000882034f000006a50140019800000000031700190000069e0000613d000000000908034f000000000a070019000000009b09043c000000000aba043600000000003a004b0000069a0000c13d0000001f09400190000006ab0000613d000000000118034f0000000308900210000000000903043300000000098901cf000000000989022f000000000101043b0000010008800089000000000181022f00000000018101cf000000000191019f00000000001304350000000001470019000000000001043500000009010000290000004001100039000600000001001d00000000005104350000000f01000029000000600110008a000006600010009c000000840000213d000000600010008c000000840000413d000000400300043d000006620030009c000001220000213d0000006001300039000000400010043f0000002001600039000000000412034f000000000404043b000000000004004b0000000005000039000000010500c039000000000054004b000000840000c13d00000000044304360000002001100039000000000512034f000000000505043b000006630050009c000000840000213d00000000005404350000002004100039000000000142034f000000000101043b000006630010009c000000840000213d0000004005300039000000000015043500000009010000290000006001100039000700000001001d00000000003104350000000f01000029000000c00110008a000006600010009c000000840000213d000000600010008c000000840000413d000000400100043d000006620010009c000001220000213d0000006003100039000000400030043f0000002003400039000000000432034f000000000404043b000000000004004b0000000005000039000000010500c039000000000054004b000000840000c13d00000000044104360000002003300039000000000532034f000000000505043b000006630050009c000000840000213d00000000005404350000002003300039000000000232034f000000000302043b000006630030009c000000840000213d0000004002100039000000000032043500000009030000290000008003300039000400000003001d0000000000130435000000070300002900000000030304330000004005300039000000000505043300000663065001970000000057030434000000000007004b0000070a0000613d000000000006004b000011310000613d00000000050504330000066305500197000000000056004b0000070f0000413d000011310000013d000000000006004b0000111d0000c13d000000000505043300000663005001980000111d0000c13d000000000202043300000663022001970000000003010433000000000003004b0000071b0000613d000000000002004b000011380000613d00000000030404330000066303300197000000000032004b000007200000413d000011380000013d000000000002004b000011210000c13d00000000020404330000066300200198000011210000c13d000000060100002900000000010104330000000001010433000000000001004b00000a760000613d000000090100002900000000010104330000062001100197001100000001001d000000000010043f0000000601000039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000100200190000000840000613d000000000101043b000000000101041a000000000001004b000011160000c13d0000000501000039000000000101041a000006200010009c000001220000213d00000001021000390000000503000039000000000023041b000006680110009a0000001102000029000000000021041b000000000103041a001000000001001d000000000020043f0000000601000039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000100200190000000840000613d000000000101043b0000001002000029000000000021041b000000090100002900000000010104330000062001100197000000000010043f0000000701000039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000100200190000000840000613d000000000301043b000000070100002900000000010104330000000024010434001100000004001d000000000004004b0000000004000039000000010400c039001000000004001d000000400400043d000006610040009c000001220000213d000f00000003001d0000000002020433000006630220019700000040011000390000000001010433000b00000001001d000000a001400039000000400010043f000e00000002001d000d00000004001d0000000001240436000c00000001001d0000066901000041000000000010044300000000010004140000061c0010009c0000061c01008041000000c0011002100000066a011001c70000800b02000039186a18650000040f0000000100200190000012990000613d0000000b020000290000066302200197000000000101043b0000000d040000290000004003400039000000100500002900000000005304350000008003400039000000000023043500000060034000390000000e0500002900000000005304350000061c011001970000000c030000290000000000130435000000110000006b00000000030000190000066b0300c0410000000f09000029000000000409041a0000066c04400197000000000343019f0000008002200210000000000425019f000000000353019f0000008002100210000000000323019f000000000039041b0000000103900039000000000043041b000000400300043d000006610030009c000001220000213d0000000404000029000000000404043300000020054000390000000005050433000000000604043300000040044000390000000004040433000000a007300039000000400070043f0000008007300039000006630840019700000000008704350000004007300039000000000006004b0000000006000039000000010600c039000000000067043500000020063000390000000000160435000006630150019700000060053000390000000000150435000000000013043500000000030000190000066b0300c0410000000205900039000000000605041a0000066c06600197000000000363019f000000000223019f000000000212019f000000000025041b0000008002400210000000000112019f0000000302900039000000000012041b000000060100002900000000030104330000000054030434000006200040009c000001220000213d0000000406900039000000000106041a000000010010019000000001071002700000007f0770618f0000001f0070008c00000000020000390000000102002039000000000121013f000000010010019000000e000000c13d000000200070008c001100000006001d001000000004001d000e00000003001d000007ff0000413d000d00000007001d000f00000005001d000000000060043f00000000010004140000061c0010009c0000061c01008041000000c0011002100000062b011001c70000801002000039186a18650000040f0000000100200190000000840000613d00000010040000290000001f024000390000000502200270000000200040008c0000000002004019000000000301043b0000000d010000290000001f01100039000000050110027000000000011300190000000002230019000000000012004b00000011060000290000000f05000029000007ff0000813d000000000002041b0000000102200039000000000012004b000007fb0000413d0000001f0040008c0000081e0000a13d000000000060043f00000000010004140000061c0010009c0000061c01008041000000c0011002100000062b011001c70000801002000039186a18650000040f0000000100200190000000840000613d0000001007000029000006a502700198000000000101043b0000000e080000290000082a0000613d000000010320008a0000000503300270000000000331001900000001043000390000002003000039000000110600002900000000058300190000000005050433000000000051041b00000020033000390000000101100039000000000041004b000008160000c13d0000082c0000013d000000000004004b000008220000613d0000000001050433000008230000013d00000000010000190000000302400210000006a60220027f000006a602200167000000000121016f0000000102400210000000000121019f000008380000013d00000020030000390000001106000029000000000072004b000008360000813d0000000302700210000000f80220018f000006a60220027f000006a60220016700000000038300190000000003030433000000000223016f000000000021041b000000010170021000000001011001bf000000000016041b000000080100002900000000010104330000000002010433000000000002004b0000095b0000613d0000000003000019000c00000003001d0000000502300210000000000121001900000020011000390000000001010433001000000001001d0000000031010434000000000001004b00000a760000613d00000009020000290000000002020433001100000002001d0000061c0010009c0000061c0100804100000060011002100000061c0030009c000f00000003001d0000061c0200004100000000020340190000004002200210000000000121019f00000000020004140000061c0020009c0000061c02008041000000c002200210000000000112019f00000658011001c70000801002000039186a18650000040f0000000100200190000000840000613d00000011020000290000062002200197000000000101043b001100000001001d000b00000002001d000000000020043f0000000701000039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000100200190000000840000613d000000000201043b0000001101000029000000000010043f000e00000002001d0000000601200039000d00000001001d000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000100200190000000840000613d000000000101043b000000000101041a000000000001004b00000f8f0000c13d0000000e010000290000000502100039000000000102041a000006200010009c000001220000213d000a00000001001d0000000101100039000000000012041b000e00000002001d000000000020043f00000000010004140000061c0010009c0000061c01008041000000c0011002100000062b011001c70000801002000039186a18650000040f00000001002001900000001102000029000000840000613d000000000101043b0000000a01100029000000000021041b0000000e01000029000000000101041a000e00000001001d000000000020043f0000000d01000029000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f00000001002001900000001102000029000000840000613d000000000101043b0000000e03000029000000000031041b000000000020043f0000000801000039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000100200190000000840000613d000000000401043b00000010010000290000000005010433000006200050009c000001220000213d000000000104041a000000010010019000000001031002700000007f0330618f0000001f0030008c00000000020000390000000102002039000000000121013f00000001001001900000000f0700002900000e000000c13d000000200030008c001100000004001d000e00000005001d000008eb0000413d000d00000003001d000000000040043f00000000010004140000061c0010009c0000061c01008041000000c0011002100000062b011001c70000801002000039186a18650000040f0000000100200190000000840000613d0000000e050000290000001f025000390000000502200270000000200050008c0000000002004019000000000301043b0000000d010000290000001f01100039000000050110027000000000011300190000000002230019000000000012004b0000000f070000290000001104000029000008eb0000813d000000000002041b0000000102200039000000000012004b000008e70000413d0000001f0050008c000009170000a13d000000000040043f00000000010004140000061c0010009c0000061c01008041000000c0011002100000062b011001c70000801002000039186a18650000040f0000000100200190000000840000613d0000000e08000029000006a502800198000000000101043b000009550000613d000000010320008a00000005033002700000000003310019000000010430003900000020030000390000000f07000029000000100600002900000000056300190000000005050433000000000051041b00000020033000390000000101100039000000000041004b000009020000c13d000000000082004b000009130000813d0000000302800210000000f80220018f000006a60220027f000006a60220016700000000036300190000000003030433000000000223016f000000000021041b000000010180021000000001011001bf0000001104000029000009230000013d000000000005004b0000091b0000613d00000000010704330000091c0000013d000000000100001900000010060000290000000302500210000006a60220027f000006a602200167000000000121016f0000000102500210000000000121019f000000000014041b000000400100043d00000020020000390000000003210436000000000206043300000000002304350000004003100039000000000002004b000009340000613d000000000400001900000000053400190000000006740019000000000606043300000000006504350000002004400039000000000024004b0000092d0000413d0000001f04200039000006a5044001970000000002320019000000000002043500000040024000390000061c0020009c0000061c0200804100000060022002100000061c0010009c0000061c010080410000004001100210000000000112019f00000000020004140000061c0020009c0000061c02008041000000c002200210000000000112019f00000658011001c70000800d0200003900000002030000390000066e040000410000000b05000029186a18600000040f0000000100200190000000840000613d0000000c030000290000000103300039000000080100002900000000010104330000000002010433000000000023004b0000083f0000413d0000095b0000013d00000020030000390000000f070000290000001006000029000000000082004b0000090b0000413d000009130000013d00000004010000290000000002010433000000070100002900000000050104330000000601000029000000000301043300000009010000290000000004010433000000400100043d000000200610003900000100070000390000000000760435000006200440019700000000004104350000010007100039000000006403043400000000004704350000012003100039000000000004004b000009770000613d000000000700001900000000083700190000000009760019000000000909043300000000009804350000002007700039000000000047004b000009700000413d000000000643001900000000000604350000000076050434000000000006004b0000000006000039000000010600c039000000400810003900000000006804350000000006070433000006630660019700000060071000390000000000670435000000400550003900000000050504330000066305500197000000800610003900000000005604350000000065020434000000000005004b0000000005000039000000010500c039000000a007100039000000000057043500000000050604330000066305500197000000c0061000390000000000560435000000400220003900000000020204330000066302200197000000e00510003900000000002504350000001f02400039000006a502200197000000000212004900000000023200190000061c0020009c0000061c0200804100000060022002100000061c0010009c0000061c010080410000004001100210000000000112019f00000000020004140000061c0020009c0000061c02008041000000c002200210000000000112019f00000658011001c70000800d0200003900000001030000390000066f04000041186a18600000040f0000000100200190000000840000613d00000005020000290000000102200039000500000002001d000000010020006c000005bd0000413d00000a470000013d00000005054002100000003f065000390000062206600197000000400700043d0000000006670019000f00000007001d000000000076004b00000000070000390000000107004039000006200060009c000001220000213d0000000100700190000001220000c13d0000010007300039000000400060043f0000000f030000290000000003430436000e00000003001d00000120022000390000000003250019000000000073004b000000840000213d000000000004004b000009d30000613d0000000e0400002900000000250204340000061f0050009c000000840000213d0000000004540436000000000032004b000009cd0000413d000001600300043d0000061f0030009c000000840000213d000001800200043d001000000002001d0000061f0020009c000000840000213d0000000002000411000000000002004b00000a4b0000c13d000000400100043d0000062e0200004100000a780000013d0000069501000041000000800010043f0000065b010000410000186c000104300000069301000041000000800010043f0000065b010000410000186c00010430000000a005000039000006810300004100000000040000190000000006050019000000000503041a000000000556043600000001033000390000000104400039000000000014004b000009eb0000413d000000410160008a000006a504100197000006790040009c000001220000213d0000008001400039000000400010043f0000000000210435000000a002400039000000800300043d0000000000320435000000c002400039000000000003004b00000a070000613d000000a004000039000000000500001900000000460404340000061f0660019700000000026204360000000105500039000000000035004b00000a010000413d00000000021200490000061c0020009c0000061c0200804100000060022002100000061c0010009c0000061c010080410000004001100210000000000112019f0000186b0001042e000000a006000039000006780400004100000000050000190000000007060019000000000604041a000000000667043600000001044000390000000105500039000000000025004b00000a130000413d000000410270008a000006a504200197000006790040009c000001220000213d0000008002400039000000800500043d000000400020043f000006200050009c000000270000a13d000001220000013d0000000402000039000000000302041a0000062304300197000000000414019f000000000042041b0000061f02300197000000800020043f000000a00010043f00000000010004140000061c0010009c0000061c01008041000000c0011002100000067b011001c70000800d0200003900000001030000390000067c04000041186a18600000040f0000000100200190000000840000613d00000a470000013d000000000100041a0000062301100197000000000161019f000000000010041b00000000010004140000061c0010009c0000061c01008041000000c00110021000000658011001c70000800d0200003900000003030000390000065904000041186a18600000040f0000000100200190000000840000613d00000000010000190000186b0001042e000000100100002900000a680000013d0000000104000039000000000504041a0000062305500197000000000225019f000000000024041b000000000003004b00000a760000613d0000061f0210019800000a760000613d000000100000006b00000a760000613d000000800020043f000000c00030043f0000062401000041000000400300043d000d00000003001d00000000001304350000000001000414000000040020008c00000a7e0000c13d00000000010004150000001f0110008a00000005011002100000000103000031000000200030008c0000002004000039000000000403401900000aaa0000013d0000000f01000029000000000112034f000000000101043b0000061f0010009c000000840000213d000000400200043d00000684030000410000000000320435000000040320003900000000001304350000061c0020009c0000061c0200804100000040012002100000065f011001c70000186c00010430000000400100043d000006700200004100000000002104350000061c0010009c0000061c01008041000000400110021000000625011001c70000186c000104300000000d030000290000061c0030009c0000061c0300804100000040033002100000061c0010009c0000061c01008041000000c001100210000000000131019f00000625011001c7186a18650000040f00000060031002700000061c03300197000000200030008c000000200400003900000000040340190000001f0640018f00000020074001900000000d0570002900000a970000613d000000000801034f0000000d09000029000000008a08043c0000000009a90436000000000059004b00000a930000c13d000000000006004b00000aa40000613d000000000171034f0000000306600210000000000705043300000000076701cf000000000767022f000000000101043b0000010006600089000000000161022f00000000016101cf000000000171019f0000000000150435000100000003001f00000000010004150000001e0110008a0000000501100210000000010020019000000ac10000613d0000001f02400039000000600420018f0000000d02400029000000000042004b00000000040000390000000104004039000006200020009c000001220000213d0000000100400190000001220000c13d000000400020043f000000200030008c000000840000413d0000000d030000290000000003030433000000ff0030008c000000840000213d0000000501100270000000000103001f0000001101000029000000ff0110018f000000000031004b00000e080000c13d0000001101000029000000a00010043f0000000402000039000000000102041a000006230110019700000010011001af000000000012041b0000000f010000290000000001010433000000000001004b0000000001000039000000010100c039000000e00010043f000000000100001900000ec90000613d000000400100043d000006280010009c000001220000213d0000002002100039000000400020043f0000000000010435000000e00100043d000000000001004b00000e130000c13d000000400100043d000006980200004100000a780000013d00000010030000290000061c0030009c0000061c0300804100000040033002100000061c0010009c0000061c01008041000000c001100210000000000131019f0000065f011001c7186a18650000040f00000060031002700000061c03300197000000200030008c000000200400003900000000040340190000001f0640018f0000002007400190000000100570002900000af50000613d000000000801034f0000001009000029000000008a08043c0000000009a90436000000000059004b00000af10000c13d000000000006004b00000b020000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f000000010020019000000c400000613d0000001f01400039000000600210018f0000001001200029000000000021004b00000000020000390000000102004039000006200010009c000001220000213d0000000100200190000001220000c13d000000400010043f000000200030008c000000840000413d00000010020000290000000002020433000000000002004b0000000003000039000000010300c039000000000032004b000000840000c13d000000000002004b00000e060000c13d0000000f0100002900000020011000390000000201100367000000000101043b001000000001001d0000061f0010009c000000840000213d00000671010000410000000000100443000000000100041200000004001004430000006001000039000000240010044300000000010004140000061c0010009c0000061c01008041000000c00110021000000683011001c70000800502000039186a18650000040f0000000100200190000012990000613d000000000101043b000000000001004b00000fa60000c13d0000000f010000290000000201100367000000000101043b001000000001001d000006200010009c000000840000213d0000001001000029000000000010043f0000000601000039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000100200190000000840000613d000000400200043d000e00000002001d0000000402200039000000000101043b000000000101041a000000000001004b000010210000c13d0000068f010000410000000e030000290000000000130435000000100100002900000c3a0000013d00000000030000310000000f010000290000001002000029186a139f0000040f00000000020100190000001101000029186a14ee0000040f00000000010000190000186b0001042e0000000002000019001100000002001d0000000502200210001000000002001d0000000e012000290000000001010433000000000010043f0000000801000039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000100200190000000840000613d000000000101043b000000000201041a000000010320019000000001052002700000007f0550618f0000001f0050008c00000000040000390000000104002039000000000043004b00000e000000c13d000000400700043d0000000004570436000000000003004b00000b9a0000613d000900000004001d000a00000005001d000b00000007001d000000000010043f00000000010004140000061c0010009c0000061c01008041000000c0011002100000062b011001c70000801002000039186a18650000040f0000000100200190000000840000613d0000000a08000029000000000008004b000000200500008a00000ba20000613d000000000201043b00000000010000190000000d060000290000000b0700002900000009090000290000000003190019000000000402041a000000000043043500000001022000390000002001100039000000000081004b00000b920000413d00000ba50000013d000006a7012001970000000000140435000000000005004b00000020010000390000000001006039000000200500008a0000000d0600002900000ba50000013d00000000010000190000000d060000290000000b070000290000003f01100039000000000251016f0000000001720019000000000021004b00000000020000390000000102004039000006200010009c000001220000213d0000000100200190000001220000c13d000000400010043f00000000010604330000001102000029000000000021004b00000f830000a13d00000010030000290000000c0130002900000000007104350000000001060433000000000021004b00000f830000a13d00000001022000390000000f010000290000000001010433000000000012004b00000b5d0000413d000003ed0000013d00000000040000190000000d0c00002900000bcb0000013d0000001f07600039000006a5077001970000000006650019000000000006043500000000057500190000000104400039000000000024004b000003f80000813d0000000006150049000000400660008a0000000003630436000000200cc0003900000000060c043300000000760604340000000005650436000000000006004b00000bc30000613d00000000080000190000000009580019000000000a870019000000000a0a04330000000000a904350000002008800039000000000068004b00000bd50000413d00000bc30000013d0000000f030000290000061c0030009c0000061c0300804100000040033002100000061c0010009c0000061c01008041000000c001100210000000000131019f0000065f011001c7186a18650000040f00000060031002700000061c03300197000000200030008c000000200400003900000000040340190000001f0640018f00000020074001900000000f0570002900000bf60000613d000000000801034f0000000f09000029000000008a08043c0000000009a90436000000000059004b00000bf20000c13d000000000006004b00000c030000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f000000010020019000000ddd0000613d0000001f01400039000000600210018f0000000f01200029000000000021004b00000000020000390000000102004039000006200010009c000001220000213d0000000100200190000001220000c13d000000400010043f000000200030008c000000840000413d0000000f020000290000000002020433000000000002004b0000000003000039000000010300c039000000000032004b000000840000c13d000000000002004b00000e060000c13d0000000e010000290000000201100367000000000101043b000f00000001001d000006200010009c000000840000213d0000000f01000029000000000010043f0000000601000039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000100200190000000840000613d000000400200043d000d00000002001d0000000402200039000000000101043b000000000101041a000000000001004b00000fc40000c13d0000068f010000410000000d0300002900000000001304350000000f0100002900000000001204350000061c0030009c0000061c0300804100000040013002100000065f011001c70000186c000104300000001f0530018f0000061e06300198000000400200043d000000000462001900000de80000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b00000c470000c13d00000de80000013d0000000002000019000700000002001d000000050120021000000005011000290000000201100367000000000101043b000b00000001001d000006200010009c000000840000213d0000000b01000029000000000010043f0000000601000039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000100200190000000840000613d000000000101043b000000000301041a000000000003004b00000fd70000613d0000000501000039000000000201041a000000000002004b000012a60000613d000000010130008a000000000023004b00000c880000613d000000000012004b00000f830000a13d0000065c0130009a0000065c0220009a000000000202041a000000000021041b000000000020043f0000000601000039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039001100000003001d186a18650000040f0000000100200190000000840000613d000000000101043b0000001102000029000000000021041b0000000501000039000000000301041a000000000003004b00000f890000613d000000010130008a0000065c0230009a000000000002041b0000000502000039000000000012041b0000000b01000029000000000010043f0000000601000039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000100200190000000840000613d000000000101043b000000000001041b0000000b01000029000000000010043f0000000701000039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000100200190000000840000613d000000000101043b0000000501100039000000000301041a000000400200043d000f00000002001d001100000003001d0000000002320436000a00000002001d000000000010043f00000000010004140000061c0010009c0000061c01008041000000c0011002100000062b011001c70000801002000039186a18650000040f0000000100200190000000840000613d0000001105000029000000000005004b0000000a0200002900000cc70000613d000000000101043b0000000a020000290000000003000019000000000401041a000000000242043600000001011000390000000103300039000000000053004b00000cc10000413d0000000f0120006a0000001f01100039000006a5021001970000000f01200029000000000021004b00000000020000390000000102004039000006200010009c000001220000213d0000000100200190000001220000c13d000000400010043f0000000f010000290000000001010433000000000001004b00000d6c0000613d000000000200001900000ce10000013d000000000101043b000000000001041b000000110200002900000001022000390000000f010000290000000001010433000000000012004b00000d6c0000813d001100000002001d0000000b01000029000000000010043f0000000701000039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000100200190000000840000613d000000000301043b0000000f0100002900000000010104330000001102000029000000000021004b00000f830000a13d00000005012002100000000a011000290000000001010433000c00000001001d000000000010043f000d00000003001d0000000601300039000e00000001001d000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000100200190000000840000613d000000000101043b000000000101041a000000000001004b00000cdb0000613d0000000d020000290000000503200039000000000203041a000000000002004b000012a60000613d000000000021004b001000000001001d000d00000003001d00000d4d0000613d000900000002001d000000000030043f00000000010004140000061c0010009c0000061c01008041000000c0011002100000062b011001c70000801002000039186a18650000040f0000000100200190000000840000613d00000010020000290008000100200092000000000101043b0000000d04000029000000000204041a000000080020006c00000f830000a13d0000000902000029000000010220008a0000000001120019000000000101041a000900000001001d000000000040043f00000000010004140000061c0010009c0000061c01008041000000c0011002100000062b011001c70000801002000039186a18650000040f0000000100200190000000840000613d000000000101043b00000008011000290000000902000029000000000021041b000000000020043f0000000e01000029000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000100200190000000840000613d000000000101043b0000001002000029000000000021041b0000000d03000029000000000103041a001000000001001d000000000001004b00000f890000613d000000000030043f00000000010004140000061c0010009c0000061c01008041000000c0011002100000062b011001c70000801002000039186a18650000040f0000000100200190000000840000613d0000001002000029000000010220008a000000000101043b0000000001210019000000000001041b0000000d01000029000000000021041b0000000c01000029000000000010043f0000000e01000029000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f000000010020019000000cd90000c13d000000840000013d0000000b01000029000000000010043f0000000701000039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000100200190000000840000613d000000000301043b000000000003041b0000000101300039000000000001041b0000000201300039000000000001041b0000000301300039000000000001041b0000000404300039000000000104041a000000010010019000000001051002700000007f0550618f0000001f0050008c00000000020000390000000102002039000000000121013f000000010010019000000e000000c13d000000000005004b00000dae0000613d0000001f0050008c00000dad0000a13d000f00000005001d001100000003001d001000000004001d000000000040043f00000000010004140000061c0010009c0000061c01008041000000c0011002100000062b011001c70000801002000039186a18650000040f0000000100200190000000840000613d000000000101043b0000000f020000290000001f02200039000000050220027000000000022100190000000103100039000000000023004b00000da90000813d000000000003041b0000000103300039000000000023004b00000da50000413d0000001002000029000000000002041b00000000040100190000001103000029000000000004041b0000000501300039000000000201041a000000000001041b000000000002004b00000dc60000613d001100000002001d000000000010043f00000000010004140000061c0010009c0000061c01008041000000c0011002100000062b011001c70000801002000039186a18650000040f0000000100200190000000840000613d000000000101043b0000001102100029000000000021004b00000dc60000813d000000000001041b0000000101100039000000000021004b00000dc20000413d000000400100043d0000000b0200002900000000002104350000061c0010009c0000061c01008041000000400110021000000000020004140000061c0020009c0000061c02008041000000c002200210000000000112019f0000062b011001c70000800d0200003900000001030000390000065d04000041186a18600000040f0000000100200190000000840000613d00000007020000290000000102200039000000060020006c00000c4d0000413d000005ba0000013d0000001f0530018f0000061e06300198000000400200043d000000000462001900000de80000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b00000de40000c13d000000000005004b00000df50000613d000000000161034f0000000305500210000000000604043300000000065601cf000000000656022f000000000101043b0000010005500089000000000151022f00000000015101cf000000000161019f000000000014043500000060013002100000061c0020009c0000061c020080410000004002200210000000000112019f0000186c000104300000067302000041000001400020043f000001440010043f00000674010000410000186c000104300000068201000041000000000010043f0000002201000039000000040010043f0000065f010000410000186c00010430000006870200004100000a780000013d0000002404200039000000000034043500000626030000410000000000320435000000040320003900000000001304350000061c0020009c0000061c02008041000000400120021000000627011001c70000186c000104300000000f020000290000000002020433000000000002004b00000ec90000613d000000000200001900000e1f0000013d000000110200002900000001022000390000000f010000290000000001010433000000000012004b00000ec80000813d001100000002001d00000005012002100000000e0110002900000000010104330000061f0310019800000e190000613d000000000030043f0000000301000039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039001000000003001d186a18650000040f00000010040000290000000100200190000000840000613d000000000101043b000000000101041a000000000001004b00000e190000c13d0000000203000039000000000103041a000006200010009c000001220000213d0000000102100039000000000023041b0000062a0110009a000000000041041b000000000103041a000d00000001001d000000000040043f0000000301000039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f00000010030000290000000100200190000000840000613d000000000101043b0000000d02000029000000000021041b000000400100043d00000000003104350000061c0010009c0000061c01008041000000400110021000000000020004140000061c0020009c0000061c02008041000000c002200210000000000112019f0000062b011001c70000800d0200003900000001030000390000062c04000041186a18600000040f000000010020019000000e190000c13d000000840000013d000000000001004b00000e870000613d000000a00200043d0000066302200197000000000021004b00000e870000813d0000001101000029000000000010043f0000000701000039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000100200190000000840000613d000000000101043b0000008002000039186a171a0000040f000001200100043d0000066301100197000000e00200043d000000000002004b00000edc0000c13d000000000001004b00000e830000c13d000001000100043d000006630010019800000ee20000613d000000400200043d001100000002001d000006660100004100000fc00000013d000000400200043d001100000002001d0000066401000041000000000012043500000004022000390000008001000039186a14d70000040f00000011010000290000061c0010009c0000061c01008041000000400110021000000665011001c70000186c000104300000002001700039000000100200002900000000002104350000002001000039000000000017043500000040017000390000000a021000290000000f0300002900000002033003670000000a0000006b00000ea50000613d000000000403034f0000000005010019000000004604043c0000000005650436000000000025004b00000ea10000c13d0000000e0000006b00000eb30000613d0000000a033003600000000e040000290000000304400210000000000502043300000000054501cf000000000545022f000000000303043b0000010004400089000000000343022f00000000034301cf000000000353019f0000000000320435000000100110002900000000000104350000061c0070009c0000061c0700804100000040017002100000000b020000290000067d0020009c0000067d020080410000006002200210000000000112019f00000000020004140000061c0020009c0000061c02008041000000c002200210000000000121019f0000067e0110009a0000800d0200003900000002030000390000067f04000041000000110500002900000a340000013d000000e00100043d000000800200043d00000140000004430000016000200443000000a00200043d00000020030000390000018000300443000001a000200443000000c00200043d0000004004000039000001c000400443000001e0002004430000006002000039000002000020044300000220001004430000010000300443000000040100003900000120001004430000062d010000410000186b0001042e000000000001004b00000fbd0000613d000001000200043d0000066302200197000000000021004b00000fbd0000813d0000001101000029000000000010043f0000000701000039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000100200190000000840000613d000000000101043b0000000201100039000000e002000039186a171a0000040f000000400100043d00000011020000290000000002210436000000800300043d000000000003004b0000000003000039000000010300c0390000000000320435000000a00200043d000006630220019700000040031000390000000000230435000000c00200043d000006630220019700000060031000390000000000230435000000e00200043d000000000002004b0000000002000039000000010200c03900000080031000390000000000230435000001000200043d0000066302200197000000a0031000390000000000230435000001200200043d0000066302200197000000c00310003900000000002304350000061c0010009c0000061c01008041000000400110021000000000020004140000061c0020009c0000061c02008041000000c002200210000000000112019f00000675011001c70000800d020000390000000103000039000006760400004100000a340000013d000000000200001900000f250000013d00000010020000290000000102200039000000800100043d000000000012004b0000047e0000813d001000000002001d0000000501200210000000a00110003900000000010104330000061f01100197001100000001001d000000000010043f0000000301000039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000100200190000000840000613d000000000101043b000000000301041a000000000003004b00000f200000613d0000000201000039000000000201041a000000000002004b000012a60000613d000000010130008a000000000023004b00000f5d0000613d000000000012004b00000f830000a13d000006960130009a000006960220009a000000000202041a000000000021041b000000000020043f0000000301000039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039000d00000003001d186a18650000040f0000000d030000290000000100200190000000840000613d000000000101043b000000000031041b0000000201000039000000000301041a000000000003004b00000f890000613d000000010130008a000006960230009a000000000002041b0000000202000039000000000012041b0000001101000029000000000010043f0000000301000039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000100200190000000840000613d000000000101043b000000000001041b000000400100043d000000110200002900000000002104350000061c0010009c0000061c01008041000000400110021000000000020004140000061c0020009c0000061c02008041000000c002200210000000000112019f0000062b011001c70000800d0200003900000001030000390000069704000041186a18600000040f000000010020019000000f200000c13d000000840000013d0000068201000041000000000010043f0000003201000039000000040010043f0000065f010000410000186c000104300000068201000041000000000010043f0000003101000039000000040010043f0000065f010000410000186c00010430000000400300043d001100000003001d0000002401300039000000400200003900000000002104350000066d01000041000000000013043500000004013000390000000b02000029000000000021043500000044023000390000001001000029186a132d0000040f000000110200002900000000012100490000061c0010009c0000061c010080410000061c0020009c0000061c0200804100000060011002100000004002200210000000000121019f0000186c000104300000001001000029000000000010043f0000000301000039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000100200190000000840000613d000000000101043b000000000101041a000000000001004b00000b340000c13d000000400100043d00000688020000410000000000210435000000040210003900000010030000290000038c0000013d000000400200043d001100000002001d000006640100004100000000001204350000000402200039000000e00100003900000e8d0000013d0000000401000039000000000301041a00000699010000410000000d0400002900000000001404350000000f01000029000000000012043500000024014000390000000002000411000000000021043500000000010004140000061f02300197000000040020008c00000fdd0000c13d0000000103000031000000200030008c00000020040000390000000004034019000010060000013d000000400100043d0000065e02000041000000000021043500000004021000390000000b030000290000038c0000013d0000000d030000290000061c0030009c0000061c0300804100000040033002100000061c0010009c0000061c01008041000000c001100210000000000131019f00000627011001c7186a18650000040f00000060031002700000061c03300197000000200030008c000000200400003900000000040340190000001f0640018f00000020074001900000000d0570002900000ff60000613d000000000801034f0000000d09000029000000008a08043c0000000009a90436000000000059004b00000ff20000c13d000000000006004b000010030000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f0000000100200190000010310000613d0000001f01400039000000600210018f0000000d01200029000000000021004b00000000020000390000000102004039000006200010009c000001220000213d0000000100200190000001220000c13d000000400010043f000000200030008c000000840000413d0000000d020000290000000002020433000000000002004b0000000003000039000000010300c039000000000032004b000000840000c13d000000000002004b0000113e0000c13d0000068a020000410000000000210435000000040210003900000000030004110000038c0000013d0000000401000039000000000301041a00000689010000410000000e0400002900000000001404350000001001000029000000000012043500000000010004140000061f02300197000000040020008c0000103d0000c13d0000000103000031000000200030008c00000020040000390000000004034019000010660000013d0000001f0530018f0000061e06300198000000400200043d000000000462001900000de80000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000010380000c13d00000de80000013d0000000e030000290000061c0030009c0000061c0300804100000040033002100000061c0010009c0000061c01008041000000c001100210000000000131019f0000065f011001c7186a18650000040f00000060031002700000061c03300197000000200030008c000000200400003900000000040340190000001f0640018f00000020074001900000000e05700029000010560000613d000000000801034f0000000e09000029000000008a08043c0000000009a90436000000000059004b000010520000c13d000000000006004b000010630000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f0000000100200190000011250000613d0000001f01400039000000600210018f0000000e01200029000000000021004b00000000020000390000000102004039000006200010009c000001220000213d0000000100200190000001220000c13d000000400010043f000000200030008c000000840000413d0000000e0200002900000000020204330000066b0020009c000000840000813d0000000003000411000000000023004b0000117e0000c13d00000002010003670000000f02100360000000000202043b000006200020009c000000840000213d0000000f030000290000004003300039000000000131034f000000000101043b001000000001001d000000000020043f0000000701000039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000100200190000000840000613d000000000101043b00000010020000290000001103000029186a17870000040f0000068b0100004100000000001004430000001101000029000000040010044300000000010004140000061c0010009c0000061c01008041000000c0011002100000068c011001c70000800202000039186a18650000040f0000000100200190000012990000613d000000000101043b000000000001004b000000840000613d000000400200043d0000068d010000410000000000120435000e00000002001d00000004012000390000001002000029000000000021043500000000010004140000001102000029000000040020008c000010be0000613d0000000e020000290000061c0020009c0000061c0200804100000040022002100000061c0010009c0000061c01008041000000c001100210000000000121019f0000065f011001c70000001102000029186a18600000040f00000060031002700001061c0030019d0000000100200190000012090000613d0000000e01000029000006200010009c000001220000213d0000000e02000029000000400020043f000000100100002900000000001204350000061c0020009c0000061c02008041000000400120021000000000020004140000061c0020009c0000061c02008041000000c002200210000000000112019f0000062b011001c70000800d0200003900000002030000390000068e040000410000000005000411186a18600000040f0000000100200190000000840000613d0000000f010000290000000201100367000000000101043b000006200010009c000000840000213d186a147c0000040f000000400200043d001100000002001d0000000002000412001700000002001d001600200000003d001000000001001d000080050100003900000044030000390000000004000415000000170440008a00000005044002100000067102000041186a18420000040f000000ff0310018f00000011010000290000002002100039000000000032043500000020020000390000000000210435186a13170000040f000000400100043d000e00000001001d186a13170000040f0000000e020000290000002001200039000f00000001001d00000011030000290000000000310435000000100100002900000000001204350000000003020019000000400400043d001100000004001d000000200100003900000000021404360000000001030433000000400300003900000000003204350000006002400039186a132d0000040f000000000201001900000011040000290000000001410049000000200310008a0000000f01000029000000000101043300000040044000390000000000340435186a132d0000040f000000110200002900000000012100490000061c0020009c0000061c0200804100000040022002100000061c0010009c0000061c010080410000006001100210000000000121019f0000186b0001042e00000009010000290000000001010433000000400200043d00000667030000410000000000320435000006200110019700000a6f0000013d000000400200043d001100000002001d0000066601000041000011340000013d000000400300043d001100000003001d00000666020000410000113b0000013d0000001f0530018f0000061e06300198000000400200043d000000000462001900000de80000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b0000112c0000c13d00000de80000013d000000400200043d001100000002001d000006640100004100000000001204350000000402200039000000000103001900000e8d0000013d000000400300043d001100000003001d00000664020000410000000000230435000000040230003900000e8d0000013d00000002020003670000000e01200360000000000101043b000006200010009c000000840000213d0000000e030000290000008003300039000000000332034f000000000303043b0000000004000031000000110540006a000000230550008a00000621065001970000062107300197000000000867013f000000000067004b00000000060000190000062106004041000000000053004b00000000050000190000062105008041000006210080009c000000000605c019000000000006004b000000840000c13d0000001105000029000f00040050003d0000000f05300029000000000252034f000000000302043b000006200030009c000000840000213d0000000004340049000000200250003900000621054001970000062106200197000000000756013f000000000056004b00000000050000190000062105004041000000000042004b00000000040000190000062104002041000006210070009c000000000504c019000000000005004b000000840000c13d186a13d70000040f000000000001004b000011820000c13d0000001101000029000000a4021000390000000f01000029186a13750000040f0000069e03000041000000400500043d001100000005001d00000000003504350000000403500039000000200400003900000000004304350000002403500039186a14480000040f00000f9c0000013d0000068a02000041000000000021043500000004021000390000038c0000013d00000002010003670000000e02100360000000000202043b000006200020009c000000840000213d0000000e03000029000d00400030003d0000000d01100360000000000101043b000e00000001001d000000000020043f0000000701000039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000100200190000000840000613d000000000101043b00000002011000390000000e020000290000001003000029186a17870000040f0000000d01000029000d00600010003d00000002030003670000000d01300360000000000101043b0000000004000031000000110240006a000000230220008a00000621052001970000062106100197000000000756013f000000000056004b00000000050000190000062105004041000000000021004b00000000020000190000062102008041000006210070009c000000000502c019000000000005004b000000840000c13d0000000f01100029000000000213034f000000000202043b000006200020009c000000840000213d0000000005240049000000200610003900000621015001970000062107600197000000000817013f000000000017004b00000000010000190000062101004041000000000056004b00000000050000190000062105002041000006210080009c000000000105c019000000000001004b000000840000c13d0000001f01200039000006a5011001970000003f01100039000006a505100197000000400100043d0000000005510019000000000015004b00000000080000390000000108004039000006200050009c000001220000213d0000000100800190000001220000c13d000000400050043f00000000052104360000000008620019000000000048004b000000840000213d000000000463034f000006a5062001980000001f0720018f0000000003650019000011e30000613d000000000804034f0000000009050019000000008a08043c0000000009a90436000000000039004b000011df0000c13d000000000007004b000011f00000613d000000000464034f0000000306700210000000000703043300000000076701cf000000000767022f000000000404043b0000010006600089000000000464022f00000000046401cf000000000474019f0000000000430435000000000225001900000000000204350000000002010433000000200020008c000012160000613d000000000002004b0000121a0000c13d00000671010000410000000000100443000000000100041200000004001004430000002001000039000000240010044300000000010004140000061c0010009c0000061c01008041000000c00110021000000683011001c70000800502000039186a18650000040f0000000100200190000012990000613d000000000101043b001100000001001d000012230000013d0000061c033001970000001f0530018f0000061e06300198000000400200043d000000000462001900000de80000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000012110000c13d00000de80000013d0000000002050433001100000002001d000000ff0020008c000012230000a13d000000400400043d001100000004001d0000069d020000410000000000240435000000040240003900000020030000390000000000320435000000240240003900000f9b0000013d00000671010000410000000000100443000000000100041200000004001004430000002001000039000000240010044300000000010004140000061c0010009c0000061c01008041000000c00110021000000683011001c70000800502000039186a18650000040f0000000100200190000012990000613d0000001102000029000000ff0220018f000000000301043b000000ff0430018f000000000142004b0000129a0000c13d0000000e01000029001100000001001d0000000d01000029000000800110008a000e00000001001d0000000201100367000000000101043b000f00000001001d0000061f0010009c000000840000213d0000068b0100004100000000001004430000001001000029000000040010044300000000010004140000061c0010009c0000061c01008041000000c0011002100000068c011001c70000800202000039186a18650000040f0000000100200190000012990000613d000000000101043b000000000001004b000000840000613d000000400300043d0000002401300039000000110200002900000000002104350000069b010000410000000000130435000d00000003001d00000004013000390000000f02000029000000000021043500000000010004140000001002000029000000040020008c0000126f0000613d0000000d020000290000061c0020009c0000061c0200804100000040022002100000061c0010009c0000061c01008041000000c001100210000000000121019f00000627011001c70000001002000029186a18600000040f00000060031002700001061c0030019d0000000100200190000012b20000613d0000000d01000029000006200010009c000001220000213d0000000d01000029000000400010043f0000000e010000290000000201100367000000000601043b0000061f0060009c000000840000213d00000011010000290000000d0200002900000000001204350000061c0020009c0000061c02008041000000400120021000000000020004140000061c0020009c0000061c02008041000000c002200210000000000112019f0000062b011001c70000800d0200003900000003030000390000069c040000410000000005000411186a18600000040f0000000100200190000000840000613d000000400100043d001000000001001d186a130c0000040f000000110200002900000010010000290000000000210435000000400100043d00000000002104350000061c0010009c0000061c01008041000000400110021000000690011001c70000186b0001042e000000000001042f000012a30000a13d000001000010008c000012a60000813d0000004d0010008c000012e10000213d000000000001004b000012bf0000c13d0000000102000039000012ca0000013d0000000001240049000000ff0010008c000012ac0000a13d0000068201000041000000000010043f0000001101000039000000040010043f0000065f010000410000186c000104300000004e0010008c000012e10000813d000000000001004b000012cc0000c13d0000000102000039000012de0000013d0000061c033001970000001f0530018f0000061e06300198000000400200043d000000000462001900000de80000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000012ba0000c13d00000de80000013d0000000a030000390000000102000039000000010010019000000000043300a9000000010300603900000000022300a900000001011002720000000003040019000012c10000c13d000000000002004b000012d80000613d0000000e012000f9000012390000013d0000000a0500003900000001020000390000000004010019000000010040019000000000065500a9000000010500603900000000022500a900000001044002720000000005060019000012cf0000c13d000000000002004b000012de0000c13d0000068201000041000000000010043f0000001201000039000000040010043f0000065f010000410000186c00010430000006a6022001290000000e0020006c000012f30000813d000000400200043d001000000002001d0000069a010000410000000000120435000000040120003900000011020000290000000e04000029186a14e50000040f000000100200002900000000012100490000061c0010009c0000061c0100804100000060011002100000061c0020009c0000061c020080410000004002200210000000000121019f0000186c00010430000000ff0210018f0000004d0020008c000012a60000213d000000000002004b000012fa0000c13d0000000101000039000013030000013d0000000a030000390000000101000039000000010020019000000000043300a9000000010300603900000000011300a900000001022002720000000003040019000012fc0000c13d0000000e0000006b001100000000001d0000123a0000613d0000000e031000b9001100000003001d0000000e023000fa000000000012004b0000123a0000613d000012a60000013d000006a80010009c000013110000813d0000002001100039000000400010043f000000000001042d0000068201000041000000000010043f0000004101000039000000040010043f0000065f010000410000186c00010430000006a90010009c0000131c0000813d0000004001100039000000400010043f000000000001042d0000068201000041000000000010043f0000004101000039000000040010043f0000065f010000410000186c00010430000006aa0010009c000013270000813d000000a001100039000000400010043f000000000001042d0000068201000041000000000010043f0000004101000039000000040010043f0000065f010000410000186c0001043000000000430104340000000001320436000000000003004b000013390000613d000000000200001900000000052100190000000006240019000000000606043300000000006504350000002002200039000000000032004b000013320000413d000000000231001900000000000204350000001f02300039000006a5022001970000000001210019000000000001042d000006600010009c0000135a0000213d000000430010008c0000135a0000a13d00000002020003670000000403200370000000000403043b000006200040009c0000135a0000213d0000002403200370000000000503043b000006200050009c0000135a0000213d0000002303500039000000000013004b0000135a0000813d0000000403500039000000000232034f000000000302043b000006200030009c0000135a0000213d00000024025000390000000005320019000000000015004b0000135a0000213d0000000001040019000000000001042d00000000010000190000186c0001043000000000430204340000066303300197000000000331043600000000040404330000061c04400197000000000043043500000040032000390000000003030433000000000003004b0000000003000039000000010300c039000000400410003900000000003404350000006003200039000000000303043300000663033001970000006004100039000000000034043500000080022000390000000002020433000006630220019700000080031000390000000000230435000000a001100039000000000001042d0000000204000367000000000224034f000000000202043b000000000300003100000000051300490000001f0550008a00000621065001970000062107200197000000000867013f000000000067004b00000000060000190000062106002041000000000052004b00000000050000190000062105004041000006210080009c000000000605c019000000000006004b0000139d0000613d0000000001120019000000000214034f000000000202043b000006200020009c0000139d0000213d0000000003230049000000200110003900000621043001970000062105100197000000000645013f000000000045004b00000000040000190000062104004041000000000031004b00000000030000190000062103002041000006210060009c000000000403c019000000000004004b0000139d0000c13d000000000001042d00000000010000190000186c00010430000006ab0020009c000013cf0000813d00000000040100190000001f01200039000006a5011001970000003f01100039000006a505100197000000400100043d0000000005510019000000000015004b00000000070000390000000107004039000006200050009c000013cf0000213d0000000100700190000013cf0000c13d000000400050043f00000000052104360000000007420019000000000037004b000013d50000213d000006a5062001980000001f0720018f00000002044003670000000003650019000013bf0000613d000000000804034f0000000009050019000000008a08043c0000000009a90436000000000039004b000013bb0000c13d000000000007004b000013cc0000613d000000000464034f0000000306700210000000000703043300000000076701cf000000000767022f000000000404043b0000010006600089000000000464022f00000000046401cf000000000474019f000000000043043500000000022500190000000000020435000000000001042d0000068201000041000000000010043f0000004101000039000000040010043f0000065f010000410000186c0001043000000000010000190000186c000104300003000000000002000300000003001d000200000002001d0000062001100197000000000010043f0000000701000039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000207000029000000030a0000290000000100200190000014400000613d0000000003000031000000000601043b000006ab00a0009c000014420000813d0000001f01a00039000006a5011001970000003f01100039000006a502100197000000400100043d0000000002210019000000000012004b00000000050000390000000105004039000006200020009c000014420000213d0000000100500190000014420000c13d000100000006001d000000400020043f0000000002a1043600000000057a0019000000000035004b000014400000213d000006a504a001980000001f05a0018f000000020670036700000000034200190000140b0000613d000000000706034f0000000008020019000000007907043c0000000008980436000000000038004b000014070000c13d000000000005004b000014180000613d000000000446034f0000000305500210000000000603043300000000065601cf000000000656022f000000000404043b0000010005500089000000000454022f00000000045401cf000000000464019f00000000004304350000000003a2001900000000000304350000061c0020009c0000061c02008041000000400220021000000000010104330000061c0010009c0000061c010080410000006001100210000000000121019f00000000020004140000061c0020009c0000061c02008041000000c002200210000000000112019f00000658011001c70000801002000039186a18650000040f0000000100200190000014400000613d000000000101043b000000000010043f00000001010000290000000601100039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000100200190000014400000613d000000000101043b000000000101041a000000000001004b0000000001000039000000010100c039000000000001042d00000000010000190000186c000104300000068201000041000000000010043f0000004101000039000000040010043f0000065f010000410000186c000104300000000003230436000006a5062001980000001f0720018f00000000056300190000000201100367000014540000613d000000000801034f0000000009030019000000008a08043c0000000009a90436000000000059004b000014500000c13d000000000007004b000014610000613d000000000161034f0000000306700210000000000705043300000000076701cf000000000767022f000000000101043b0000010006600089000000000161022f00000000016101cf000000000171019f0000000000150435000000000123001900000000000104350000001f01200039000006a5011001970000000001130019000000000001042d000000400100043d000006aa0010009c000014760000813d000000a002100039000000400020043f000000800210003900000000000204350000006002100039000000000002043500000040021000390000000000020435000000200210003900000000000204350000000000010435000000000001042d0000068201000041000000000010043f0000004101000039000000040010043f0000065f010000410000186c0001043000030000000000020000062001100197000000000010043f0000000701000039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000100200190000014c90000613d000000000101043b0000000405100039000000000205041a000000010320019000000001062002700000007f0660618f0000001f0060008c00000000040000390000000104002039000000000043004b000014cb0000c13d000000400100043d0000000004610436000000000003004b000014b50000613d000100000004001d000200000006001d000300000001001d000000000050043f00000000010004140000061c0010009c0000061c01008041000000c0011002100000062b011001c70000801002000039186a18650000040f0000000100200190000014c90000613d0000000206000029000000000006004b000014bb0000613d000000000201043b0000000005000019000000030100002900000001070000290000000003570019000000000402041a000000000043043500000001022000390000002005500039000000000065004b000014ad0000413d000014bd0000013d000006a7022001970000000000240435000000000006004b00000020050000390000000005006039000014bd0000013d000000000500001900000003010000290000003f03500039000006a5023001970000000003120019000000000023004b00000000020000390000000102004039000006200030009c000014d10000213d0000000100200190000014d10000c13d000000400030043f000000000001042d00000000010000190000186c000104300000068201000041000000000010043f0000002201000039000000040010043f0000065f010000410000186c000104300000068201000041000000000010043f0000004101000039000000040010043f0000065f010000410000186c000104300000000043010434000000000003004b0000000003000039000000010300c039000000000332043600000000040404330000066304400197000000000043043500000040022000390000004001100039000000000101043300000663011001970000000000120435000000000001042d00000040051000390000000000450435000000ff0330018f00000020041000390000000000340435000000ff0220018f00000000002104350000006001100039000000000001042d0007000000000002000400000001001d000600000002001d0000000021020434000000000001004b000016070000613d0000061c0010009c0000061c0100804100000060011002100000061c0020009c000500000002001d0000061c020080410000004002200210000000000121019f00000000020004140000061c0020009c0000061c02008041000000c002200210000000000112019f00000658011001c70000801002000039186a18650000040f0000000100200190000015ff0000613d000000000101043b000700000001001d00000004010000290000062001100197000200000001001d000000000010043f0000000701000039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000100200190000015ff0000613d000000000201043b0000000701000029000000000010043f000400000002001d0000000601200039000300000001001d000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000100200190000015ff0000613d000000000101043b000000000101041a000000000001004b0000160f0000c13d00000004010000290000000502100039000000000102041a000006ab0010009c000016010000813d000100000001001d0000000101100039000000000012041b000400000002001d000000000020043f00000000010004140000061c0010009c0000061c01008041000000c0011002100000062b011001c70000801002000039186a18650000040f00000001002001900000000702000029000015ff0000613d000000000101043b0000000101100029000000000021041b0000000401000029000000000101041a000400000001001d000000000020043f0000000301000029000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f00000001002001900000000702000029000015ff0000613d000000000101043b0000000403000029000000000031041b000000000020043f0000000801000039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000100200190000015ff0000613d000000000801043b00000006010000290000000004010433000006200040009c000016010000213d000000000108041a000000010210019000000001031002700000007f0330618f0000001f0030008c00000000010000390000000101002039000000000012004b0000000507000029000016260000c13d000000200030008c000400000008001d000700000004001d000015920000413d000300000003001d000000000080043f00000000010004140000061c0010009c0000061c01008041000000c0011002100000062b011001c70000801002000039186a18650000040f0000000100200190000015ff0000613d00000007040000290000001f024000390000000502200270000000200040008c0000000002004019000000000301043b00000003010000290000001f01100039000000050110027000000000011300190000000002230019000000000012004b00000005070000290000000408000029000015920000813d000000000002041b0000000102200039000000000012004b0000158e0000413d0000001f0040008c000000200a00008a000000200b000039000015c20000a13d000000000080043f00000000010004140000061c0010009c0000061c01008041000000c0011002100000062b011001c70000801002000039186a18650000040f0000000100200190000015ff0000613d0000000709000029000000200a00008a0000000002a90170000000000101043b000000200b000039000015f80000613d000000010320008a000000050330027000000000043100190000002003000039000000010440003900000005070000290000000606000029000000040800002900000000056300190000000005050433000000000051041b00000020033000390000000101100039000000000041004b000015ae0000c13d000000000092004b000015bf0000813d0000000302900210000000f80220018f000006a60220027f000006a60220016700000000036300190000000003030433000000000223016f000000000021041b000000010190021000000001011001bf000015ce0000013d000000000004004b000015c60000613d0000000001070433000015c70000013d000000000100001900000006060000290000000302400210000006a60220027f000006a602200167000000000121016f0000000102400210000000000121019f000000000018041b000000400100043d0000000003b10436000000000206043300000000002304350000004003100039000000000002004b000015de0000613d000000000400001900000000053400190000000006740019000000000606043300000000006504350000002004400039000000000024004b000015d70000413d0000001f042000390000000004a4016f0000000002230019000000000002043500000040024000390000061c0020009c0000061c0200804100000060022002100000061c0010009c0000061c010080410000004001100210000000000112019f00000000020004140000061c0020009c0000061c02008041000000c002200210000000000121019f00000658011001c70000800d0200003900000002030000390000066e040000410000000205000029186a18600000040f0000000100200190000015ff0000613d000000000001042d00000000030b0019000000050700002900000006060000290000000408000029000000000092004b000015b70000413d000015bf0000013d00000000010000190000186c000104300000068201000041000000000010043f0000004101000039000000040010043f0000065f010000410000186c00010430000000400100043d000006700200004100000000002104350000061c0010009c0000061c01008041000000400110021000000625011001c70000186c00010430000000400300043d000700000003001d0000002401300039000000400200003900000000002104350000066d01000041000000000013043500000004013000390000000202000029000000000021043500000044023000390000000601000029186a132d0000040f000000070200002900000000012100490000061c0010009c0000061c010080410000061c0020009c0000061c0200804100000060011002100000004002200210000000000121019f0000186c000104300000068201000041000000000010043f0000002201000039000000040010043f0000065f010000410000186c000104300005000000000002000000400300043d000006aa0030009c000016730000813d000000a002300039000000400020043f00000080023000390000000000020435000000600230003900000000000204350000004002300039000000000002043500000020023000390000000000020435000000000003043500000060021000390000000002020433000100000002001d000500000001001d0000000012010434000300000002001d000200000001001d0000000001010433000400000001001d0000066901000041000000000010044300000000010004140000061c0010009c0000061c01008041000000c0011002100000066a011001c70000800b02000039186a18650000040f0000000100200190000016790000613d00000004020000290000061c04200197000000000601043b000000000346004b00000005010000290000166d0000413d00000080021000390000000002020433000006630520019700000000023500a9000000000046004b0000165e0000613d00000000033200d9000000000053004b0000166d0000c13d00000003030000290000066303300197000000000032001a0000166d0000413d0000000002320019000000010300002900000663033001970000066304200197000000000023004b000000000304801900000000003104350000061c0260019700000002030000290000000000230435000000000001042d0000068201000041000000000010043f0000001101000039000000040010043f0000065f010000410000186c000104300000068201000041000000000010043f0000004101000039000000040010043f0000065f010000410000186c00010430000000000001042f000000000010043f0000000601000039000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f00000001002001900000168c0000613d000000000101043b000000000101041a000000000001004b0000000001000039000000010100c039000000000001042d00000000010000190000186c000104300006000000000002000300000002001d000000000020043f000600000001001d0000000101100039000400000001001d000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000100200190000017060000613d0000000603000029000000000101043b000000000101041a000000000001004b000017040000613d000000000203041a000000000002004b000017080000613d000000000021004b000500000001001d000016e20000613d000200000002001d000000000030043f00000000010004140000061c0010009c0000061c01008041000000c0011002100000062b011001c70000801002000039186a18650000040f0000000100200190000017060000613d00000005020000290001000100200092000000000101043b0000000604000029000000000204041a000000010020006c0000170e0000a13d0000000202000029000000010220008a0000000001120019000000000101041a000200000001001d000000000040043f00000000010004140000061c0010009c0000061c01008041000000c0011002100000062b011001c70000801002000039186a18650000040f0000000100200190000017060000613d000000000101043b00000001011000290000000202000029000000000021041b000000000020043f0000000401000029000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000100200190000017060000613d000000000101043b0000000502000029000000000021041b0000000603000029000000000103041a000500000001001d000000000001004b000017140000613d000000000030043f00000000010004140000061c0010009c0000061c01008041000000c0011002100000062b011001c70000801002000039186a18650000040f0000000100200190000017060000613d0000000502000029000000010220008a000000000101043b0000000001210019000000000001041b0000000601000029000000000021041b0000000301000029000000000010043f0000000401000029000000200010043f00000000010004140000061c0010009c0000061c01008041000000c00110021000000629011001c70000801002000039186a18650000040f0000000100200190000017060000613d000000000101043b000000000001041b0000000101000039000000000001042d0000000001000019000000000001042d00000000010000190000186c000104300000068201000041000000000010043f0000001101000039000000040010043f0000065f010000410000186c000104300000068201000041000000000010043f0000003201000039000000040010043f0000065f010000410000186c000104300000068201000041000000000010043f0000003101000039000000040010043f0000065f010000410000186c000104300003000000000002000100000002001d000300000001001d000000000101041a000200000001001d0000066901000041000000000010044300000000010004140000061c0010009c0000061c01008041000000c0011002100000066a011001c70000800b02000039186a18650000040f00000001002001900000177c0000613d000000020800002900000080028002700000061c03200197000000000201043b000000000532004b0000000307000029000017760000413d0000000101700039000017350000c13d000000000207041a000017470000013d000000000301041a000000800630027000000000045600a900000000055400d9000000000065004b000017760000c13d0000066305800197000000000054001a000017760000413d00000663033001970000000004540019000000000043004b000000000304801900000623048001970000008002200210000006ac02200197000000000242019f000000000232019f00000001060000290000002003600039000000000403043300000663044001970000066305200197000000000054004b0000000005044019000006ad02200197000000000225019f0000000005060433000000000005004b00000000050000190000066b0500c041000000000252019f000000000027041b000000400260003900000000050204330000008005500210000000000445019f000000000041041b0000000001000039000000010100c039000000400400043d000000000114043600000000030304330000066303300197000000000031043500000000010204330000066301100197000000400240003900000000001204350000061c0040009c0000061c04008041000000400140021000000000020004140000061c0020009c0000061c02008041000000c002200210000000000112019f000006ae011001c70000800d020000390000000103000039000006af04000041186a18600000040f00000001002001900000177d0000613d000000000001042d0000068201000041000000000010043f0000001101000039000000040010043f0000065f010000410000186c00010430000000000001042f00000000010000190000186c000104300000061f04400197000000400510003900000000004504350000002004100039000000000034043500000000002104350000006001100039000000000001042d0006000000000002000000000401041a0000067700400198000017db0000613d000000000002004b000017db0000613d000600000004001d000500000002001d000200000003001d000300000001001d0000000101100039000100000001001d000000000101041a000400000001001d0000066901000041000000000010044300000000010004140000061c0010009c0000061c01008041000000c0011002100000066a011001c70000800b02000039186a18650000040f0000000100200190000017dc0000613d000000060300002900000080023002700000061c02200197000000000101043b000000000421004b000017f70000413d000006630330019700000004050000290000066302500197000017ad0000c13d00000005040000290000000305000029000017c10000013d000000000023004b000017ff0000213d000000800650027000000000056400a900000000044500d9000000000064004b000017f70000c13d000000000035001a000017f70000413d00000000033500190000008001100210000006ac011001970000000305000029000000000405041a000006b004400197000000000114019f000000000015041b000000000032004b00000000030240190000000504000029000000000042004b000017dd0000413d000000000143004b000017ee0000413d0000066301100197000000000205041a000006b202200197000000000112019f000000000015041b000000400100043d00000000004104350000061c0010009c0000061c01008041000000400110021000000000020004140000061c0020009c0000061c02008041000000c002200210000000000112019f0000062b011001c70000800d020000390000000103000039000006b304000041186a18600000040f0000000100200190000017fd0000613d000000000001042d000000000001042f000000400100043d0000000004010019000000040110003900000002030000290000061f00300198000018070000c13d000006b703000041000000000034043500000000002104350000002401400039000000050200002900000000002104350000061c0040009c0000061c04008041000000400140021000000627011001c70000186c000104300000000101000029000000000101041a0000008001100272000017f70000613d00000005043000690000000002140019000000010220008a000000000042004b0000180c0000813d0000068201000041000000000010043f0000001101000039000000040010043f0000065f010000410000186c0001043000000000010000190000186c00010430000000400100043d000006b10200004100000000002104350000061c0010009c0000061c01008041000000400110021000000625011001c70000186c00010430000006b603000041000600000004001d00000000003404350000000503000029000018200000013d00000000021200d9000000400100043d0000000005010019000000040110003900000002040000290000061f004001980000181d0000c13d000006b50400004100000000004504350000000000210435000000240150003900000000003104350000061c0050009c0000061c05008041000000400150021000000627011001c70000186c00010430000006b404000041000600000005001d00000000004504350000000204000029186a177f0000040f000000060200002900000000012100490000061c0010009c0000061c0100804100000060011002100000061c0020009c0000061c020080410000004002200210000000000121019f0000186c00010430000000000001042f0000061c0010009c0000061c0100804100000040011002100000061c0020009c0000061c020080410000006002200210000000000112019f00000000020004140000061c0020009c0000061c02008041000000c002200210000000000112019f00000658011001c70000801002000039186a18650000040f0000000100200190000018400000613d000000000101043b000000000001042d00000000010000190000186c0001043000000000050100190000000000200443000000050030008c000018500000413d000000040100003900000000020000190000000506200210000000000664001900000005066002700000000006060031000000000161043a0000000102200039000000000031004b000018480000413d0000061c0030009c0000061c03008041000000600130021000000000020004140000061c0020009c0000061c02008041000000c002200210000000000112019f000006b8011001c70000000002050019186a18650000040f00000001002001900000185f0000613d000000000101043b000000000001042d000000000001042f00001863002104210000000102000039000000000001042d0000000002000019000000000001042d00001868002104230000000102000039000000000001042d0000000002000019000000000001042d0000186a000004320000186b0001042e0000186c0001043000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ffffffff00000000000000000000000000000000000000000000000000000001ffffffe000000000000000000000000000000000000000000000000000000000ffffffe0000000000000000000000000ffffffffffffffffffffffffffffffffffffffff000000000000000000000000000000000000000000000000ffffffffffffffff80000000000000000000000000000000000000000000000000000000000000007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0ffffffffffffffffffffffff0000000000000000000000000000000000000000313ce567000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000655a7c0e000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000044000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffffdf0200000000000000000000000000000000000040000000000000000000000000bfa87805ed57dc1f0d489ce33be4c4577d74ccde357eeeee058a32c55c44a53202000000000000000000000000000000000000200000000000000000000000002640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d800000002000000000000000000000000000001400000010000000000000000009b15e16f00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000009a4575b800000000000000000000000000000000000000000000000000000000c0d7865400000000000000000000000000000000000000000000000000000000dc0bd97000000000000000000000000000000000000000000000000000000000e8a1da1600000000000000000000000000000000000000000000000000000000e8a1da1700000000000000000000000000000000000000000000000000000000f2fde38b00000000000000000000000000000000000000000000000000000000dc0bd97100000000000000000000000000000000000000000000000000000000e0351e1300000000000000000000000000000000000000000000000000000000c75eea9b00000000000000000000000000000000000000000000000000000000c75eea9c00000000000000000000000000000000000000000000000000000000cf7401f300000000000000000000000000000000000000000000000000000000c0d7865500000000000000000000000000000000000000000000000000000000c4bffe2b00000000000000000000000000000000000000000000000000000000acfecf9000000000000000000000000000000000000000000000000000000000b0f479a000000000000000000000000000000000000000000000000000000000b0f479a100000000000000000000000000000000000000000000000000000000b794658000000000000000000000000000000000000000000000000000000000acfecf9100000000000000000000000000000000000000000000000000000000af58d59f000000000000000000000000000000000000000000000000000000009a4575b900000000000000000000000000000000000000000000000000000000a42a7b8b00000000000000000000000000000000000000000000000000000000a7cd63b70000000000000000000000000000000000000000000000000000000054c8a4f20000000000000000000000000000000000000000000000000000000079ba5096000000000000000000000000000000000000000000000000000000008926f54e000000000000000000000000000000000000000000000000000000008926f54f000000000000000000000000000000000000000000000000000000008da5cb5b0000000000000000000000000000000000000000000000000000000079ba5097000000000000000000000000000000000000000000000000000000007d54534e0000000000000000000000000000000000000000000000000000000054c8a4f30000000000000000000000000000000000000000000000000000000062ddd3c4000000000000000000000000000000000000000000000000000000006d3d1a5800000000000000000000000000000000000000000000000000000000240028e700000000000000000000000000000000000000000000000000000000390775360000000000000000000000000000000000000000000000000000000039077537000000000000000000000000000000000000000000000000000000004c5ef0ed00000000000000000000000000000000000000000000000000000000240028e80000000000000000000000000000000000000000000000000000000024f65ee70000000000000000000000000000000000000000000000000000000001ffc9a700000000000000000000000000000000000000000000000000000000181f5a770000000000000000000000000000000000000000000000000000000021df0da70200000000000000000000000000000000000000000000000000000000000000ed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278dad89dca000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000800000000000000000fc949c7b4a13586e39d89eead2f38644f9fb3efb5a0490b14f8fc0ceab44c2515204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d8599161e670e4b0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000240000000000000000000000007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000000000000000000000000000000000000000000000ffffffffffffff5f000000000000000000000000000000000000000000000000ffffffffffffff9f00000000000000000000000000000000ffffffffffffffffffffffffffffffff8020d124000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000064000000000000000000000000d68af9cc000000000000000000000000000000000000000000000000000000001d5ad3c500000000000000000000000000000000000000000000000000000000fc949c7b4a13586e39d89eead2f38644f9fb3efb5a0490b14f8fc0ceab44c250796b89b91644bc98cd93958e4c9038275d622183e25ac5af08cc6b5d9553913202000002000000000000000000000000000000040000000000000000000000000000000000000000000000010000000000000000000000000000000000000000ffffffffffffffffffffff000000000000000000000000000000000000000000393b8ad2000000000000000000000000000000000000000000000000000000007d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c28579befe00000000000000000000000000000000000000000000000000000000310ab089e4439a4c15d089f94afb7896ff553aecb10793d0ab882de59d99a32e00000000000000000000000000000000000000200000008000000000000000008e4a23d600000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002400000140000000000000000002000000000000000000000000000000000000e00000000000000000000000000350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b0000000000000000000000ff0000000000000000000000000000000000000000036b6384b5eca791c62761152d0c79bb0604c104a5fb6f4eb0703f3154bb3db0000000000000000000000000000000000000000000000000ffffffffffffff7f00000000000000000000000000000000000000000000003fffffffffffffffe0020000000000000000000000000000000000004000000080000000000000000002dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f168400000000000000000000000000000000000000000000000000000000ffffffbffdffffffffffffffffffffffffffffffffffffc000000000000000000000000052d00ee4d9bd51b40168f2afc5848837288ce258784ad914278791464b3f4d7674f23c7c00000000000000000000000000000000000000000000000000000000405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace4e487b71000000000000000000000000000000000000000000000000000000000200000200000000000000000000000000000044000000000000000000000000961c9a4f000000000000000000000000000000000000000000000000000000002cbc26bb000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffffff0000000000000000000000000000000053ad11d800000000000000000000000000000000000000000000000000000000d0d2597600000000000000000000000000000000000000000000000000000000a8d87a3b00000000000000000000000000000000000000000000000000000000728fe07b000000000000000000000000000000000000000000000000000000001806aa1896bbf26568e884a7374b41e002500962caba6a15023a8d90e8508b83020000020000000000000000000000000000002400000000000000000000000042966c6800000000000000000000000000000000000000000000000000000000696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df7a9902c7e000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000020000000000000000000000000000000000002000000080000000000000000044676b5284b809a22248eba0da87391d79098be38bb03154be88a58bf4d0917402b543c6000000000000000000000000000000000000000000000000000000008be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e02b5c74de00000000000000000000000000000000000000000000000000000000bfa87805ed57dc1f0d489ce33be4c4577d74ccde357eeeee058a32c55c44a533800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf756635f4a7b30000000000000000000000000000000000000000000000000000000083826b2b00000000000000000000000000000000000000000000000000000000a9cb113d0000000000000000000000000000000000000000000000000000000040c10f19000000000000000000000000000000000000000000000000000000009d228d69b5fdb8d273a2336f8fb8612d039631024ea9bf09c424a9503aa078f0953576f70000000000000000000000000000000000000000000000000000000024eb47e5000000000000000000000000000000000000000000000000000000004275726e4d696e74546f6b656e506f6f6c20312e352e310000000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffff01ffc9a7000000000000000000000000000000000000000000000000000000000e64dd2900000000000000000000000000000000000000000000000000000000aff2afbf00000000000000000000000000000000000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000000000000000000000000000000000000000000000ffffffffffffffe0000000000000000000000000000000000000000000000000ffffffffffffffc0000000000000000000000000000000000000000000000000ffffffffffffff600000000000000000000000000000000000000000000000010000000000000000000000000000000000000000ffffffff00000000000000000000000000000000ffffffffffffffffffffff00ffffffff0000000000000000000000000000000002000000000000000000000000000000000000600000000000000000000000009ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c19ffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff9725942a00000000000000000000000000000000000000000000000000000000ffffffffffffffffffffffffffffffff000000000000000000000000000000001871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690ad0c8d23a0000000000000000000000000000000000000000000000000000000015279c08000000000000000000000000000000000000000000000000000000001a76572a00000000000000000000000000000000000000000000000000000000f94ebcd1000000000000000000000000000000000000000000000000000000000200000200000000000000000000000000000000000000000000000000000000")
