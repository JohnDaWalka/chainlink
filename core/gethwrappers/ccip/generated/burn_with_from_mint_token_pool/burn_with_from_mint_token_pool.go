package burn_with_from_mint_token_pool

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

var BurnWithFromMintTokenPoolMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIBurnMintERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"localTokenDecimals\",\"type\":\"uint8\"},{\"internalType\":\"address[]\",\"name\":\"allowlist\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"}],\"name\":\"AggregateValueMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"}],\"name\":\"AggregateValueRateLimitReached\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AllowListNotEnabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BucketOverfilled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"CallerIsNotARampOnRouter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotTransferToSelf\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CursedByRMN\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"DisabledNonZeroRateLimit\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"expected\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"actual\",\"type\":\"uint8\"}],\"name\":\"InvalidDecimalArgs\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"InvalidRateLimitRate\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolData\",\"type\":\"bytes\"}],\"name\":\"InvalidRemoteChainDecimals\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidRemotePoolForChain\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidSourcePoolAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"InvalidToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeProposedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"NonExistentChain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"remoteDecimals\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"localDecimals\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"remoteAmount\",\"type\":\"uint256\"}],\"name\":\"OverflowDetected\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OwnerCannotBeZero\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"PoolAlreadyAdded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RateLimitMustBeDisabled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SenderNotAllowed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenRateLimitReached\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListAdd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListRemove\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Burned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remoteToken\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"ChainAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"ChainConfigured\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Locked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Minted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"rateLimitAdmin\",\"type\":\"address\"}],\"name\":\"RateLimitAdminSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Released\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"RemotePoolAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"RemotePoolRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldRouter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"RouterUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"TokensConsumed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"addRemotePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"removes\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"adds\",\"type\":\"address[]\"}],\"name\":\"applyAllowListUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64[]\",\"name\":\"remoteChainSelectorsToRemove\",\"type\":\"uint64[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes[]\",\"name\":\"remotePoolAddresses\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"remoteTokenAddress\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"internalType\":\"structTokenPool.ChainUpdate[]\",\"name\":\"chainsToAdd\",\"type\":\"tuple[]\"}],\"name\":\"applyChainUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowList\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowListEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getCurrentInboundRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getCurrentOutboundRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRateLimitAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getRemotePools\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"\",\"type\":\"bytes[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getRemoteToken\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRmnProxy\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRouter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSupportedChains\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTokenDecimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"decimals\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"isRemotePool\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"isSupportedChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"isSupportedToken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"originalSender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"}],\"internalType\":\"structPool.LockOrBurnInV1\",\"name\":\"lockOrBurnIn\",\"type\":\"tuple\"}],\"name\":\"lockOrBurn\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"destPoolData\",\"type\":\"bytes\"}],\"internalType\":\"structPool.LockOrBurnOutV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"originalSender\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"sourcePoolData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainTokenData\",\"type\":\"bytes\"}],\"internalType\":\"structPool.ReleaseOrMintInV1\",\"name\":\"releaseOrMintIn\",\"type\":\"tuple\"}],\"name\":\"releaseOrMint\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"destinationAmount\",\"type\":\"uint256\"}],\"internalType\":\"structPool.ReleaseOrMintOutV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"removeRemotePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundConfig\",\"type\":\"tuple\"}],\"name\":\"setChainRateLimiterConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"rateLimitAdmin\",\"type\":\"address\"}],\"name\":\"setRateLimitAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"setRouter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x6101006040523480156200001257600080fd5b5060405162004c5c38038062004c5c833981016040819052620000359162000918565b8484848484336000816200005c57604051639b15e16f60e01b815260040160405180910390fd5b600180546001600160a01b0319166001600160a01b03848116919091179091558116156200008f576200008f8162000206565b50506001600160a01b0385161580620000af57506001600160a01b038116155b80620000c257506001600160a01b038216155b15620000e1576040516342bcdf7f60e11b815260040160405180910390fd5b6001600160a01b03808616608081905290831660c0526040805163313ce56760e01b8152905163313ce567916004808201926020929091908290030181865afa92505050801562000151575060408051601f3d908101601f191682019092526200014e9181019062000a3a565b60015b1562000192578060ff168560ff161462000190576040516332ad3e0760e11b815260ff8087166004830152821660248201526044015b60405180910390fd5b505b60ff841660a052600480546001600160a01b0319166001600160a01b038316179055825115801560e052620001dc57604080516000815260208101909152620001dc908462000280565b50620001fb935050506001600160a01b038716905030600019620003dd565b505050505062000b84565b336001600160a01b038216036200023057604051636d6c4ee560e11b815260040160405180910390fd5b600080546001600160a01b0319166001600160a01b03838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60e051620002a1576040516335f4a7b360e01b815260040160405180910390fd5b60005b82518110156200032c576000838281518110620002c557620002c562000a58565b60209081029190910101519050620002df600282620004c3565b1562000322576040516001600160a01b03821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b50600101620002a4565b5060005b8151811015620003d857600082828151811062000351576200035162000a58565b6020026020010151905060006001600160a01b0316816001600160a01b0316036200037d5750620003cf565b6200038a600282620004e3565b15620003cd576040516001600160a01b03821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b60010162000330565b505050565b604051636eb1769f60e11b81523060048201526001600160a01b038381166024830152600091839186169063dd62ed3e90604401602060405180830381865afa1580156200042f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000455919062000a6e565b62000461919062000a9e565b604080516001600160a01b038616602482015260448082018490528251808303909101815260649091019091526020810180516001600160e01b0390811663095ea7b360e01b17909152919250620004bd91869190620004fa16565b50505050565b6000620004da836001600160a01b038416620005cb565b90505b92915050565b6000620004da836001600160a01b038416620006cf565b6040805180820190915260208082527f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65649082015260009062000549906001600160a01b03851690849062000721565b805190915015620003d857808060200190518101906200056a919062000ab4565b620003d85760405162461bcd60e51b815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e6044820152691bdd081cdd58d8d9595960b21b606482015260840162000187565b60008181526001830160205260408120548015620006c4576000620005f260018362000adf565b8554909150600090620006089060019062000adf565b9050808214620006745760008660000182815481106200062c576200062c62000a58565b906000526020600020015490508087600001848154811062000652576200065262000a58565b6000918252602080832090910192909255918252600188019052604090208390555b855486908062000688576200068862000af5565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050620004dd565b6000915050620004dd565b60008181526001830160205260408120546200071857508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155620004dd565b506000620004dd565b60606200073284846000856200073a565b949350505050565b6060824710156200079d5760405162461bcd60e51b815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f6044820152651c8818d85b1b60d21b606482015260840162000187565b600080866001600160a01b03168587604051620007bb919062000b31565b60006040518083038185875af1925050503d8060008114620007fa576040519150601f19603f3d011682016040523d82523d6000602084013e620007ff565b606091505b50909250905062000813878383876200081e565b979650505050505050565b60608315620008925782516000036200088a576001600160a01b0385163b6200088a5760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000604482015260640162000187565b508162000732565b620007328383815115620008a95781518083602001fd5b8060405162461bcd60e51b815260040162000187919062000b4f565b6001600160a01b0381168114620008db57600080fd5b50565b805160ff81168114620008f057600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b8051620008f081620008c5565b600080600080600060a086880312156200093157600080fd5b85516200093e81620008c5565b945060206200094f878201620008de565b60408801519095506001600160401b03808211156200096d57600080fd5b818901915089601f8301126200098257600080fd5b815181811115620009975762000997620008f5565b8060051b604051601f19603f83011681018181108582111715620009bf57620009bf620008f5565b60405291825284820192508381018501918c831115620009de57600080fd5b938501935b8285101562000a0757620009f7856200090b565b84529385019392850192620009e3565b80985050505050505062000a1e606087016200090b565b915062000a2e608087016200090b565b90509295509295909350565b60006020828403121562000a4d57600080fd5b620004da82620008de565b634e487b7160e01b600052603260045260246000fd5b60006020828403121562000a8157600080fd5b5051919050565b634e487b7160e01b600052601160045260246000fd5b80820180821115620004dd57620004dd62000a88565b60006020828403121562000ac757600080fd5b8151801515811462000ad857600080fd5b9392505050565b81810381811115620004dd57620004dd62000a88565b634e487b7160e01b600052603160045260246000fd5b60005b8381101562000b2857818101518382015260200162000b0e565b50506000910152565b6000825162000b4581846020870162000b0b565b9190910192915050565b602081526000825180602084015262000b7081604085016020870162000b0b565b601f01601f19169190910160400192915050565b60805160a05160c05160e05161402762000c356000396000818161054801528181611d7b01526127cc0152600081816105220152818161189801526120670152600081816102d901528181610ba201528181611a4101528181611afb01528181611b2f01528181611b6201528181611bc701528181611c200152611cc20152600081816102400152818161029501528181610701015281816121ea0152818161276201526129b701526140276000f3fe608060405234801561001057600080fd5b50600436106101cf5760003560e01c80639a4575b911610104578063c0d78655116100a2578063dc0bd97111610071578063dc0bd97114610520578063e0351e1314610546578063e8a1da171461056c578063f2fde38b1461057f57600080fd5b8063c0d78655146104d2578063c4bffe2b146104e5578063c75eea9c146104fa578063cf7401f31461050d57600080fd5b8063acfecf91116100de578063acfecf911461041f578063af58d59f14610432578063b0f479a1146104a1578063b7946580146104bf57600080fd5b80639a4575b9146103ca578063a42a7b8b146103ea578063a7cd63b71461040a57600080fd5b806354c8a4f31161017157806379ba50971161014b57806379ba50971461037e5780637d54534e146103865780638926f54f146103995780638da5cb5b146103ac57600080fd5b806354c8a4f31461033857806362ddd3c41461034d5780636d3d1a581461036057600080fd5b8063240028e8116101ad578063240028e81461028557806324f65ee7146102d257806339077537146103035780634c5ef0ed1461032557600080fd5b806301ffc9a7146101d4578063181f5a77146101fc57806321df0da71461023e575b600080fd5b6101e76101e2366004613177565b610592565b60405190151581526020015b60405180910390f35b60408051808201909152601f81527f4275726e5769746846726f6d4d696e74546f6b656e506f6f6c20312e352e310060208201525b6040516101f3919061321d565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101f3565b6101e7610293366004613252565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff90811691161490565b60405160ff7f00000000000000000000000000000000000000000000000000000000000000001681526020016101f3565b61031661031136600461326f565b610677565b604051905181526020016101f3565b6101e76103333660046132c8565b610846565b61034b610346366004613397565b610890565b005b61034b61035b3660046132c8565b61090b565b60095473ffffffffffffffffffffffffffffffffffffffff16610260565b61034b6109a8565b61034b610394366004613252565b610a76565b6101e76103a7366004613403565b610af7565b60015473ffffffffffffffffffffffffffffffffffffffff16610260565b6103dd6103d836600461341e565b610b0e565b6040516101f39190613459565b6103fd6103f8366004613403565b610be7565b6040516101f391906134b0565b610412610d52565b6040516101f39190613532565b61034b61042d3660046132c8565b610d63565b610445610440366004613403565b610e7b565b6040516101f3919081516fffffffffffffffffffffffffffffffff908116825260208084015163ffffffff1690830152604080840151151590830152606080840151821690830152608092830151169181019190915260a00190565b60045473ffffffffffffffffffffffffffffffffffffffff16610260565b6102316104cd366004613403565b610f50565b61034b6104e0366004613252565b611000565b6104ed6110db565b6040516101f3919061358c565b610445610508366004613403565b611193565b61034b61051b366004613714565b611265565b7f0000000000000000000000000000000000000000000000000000000000000000610260565b7f00000000000000000000000000000000000000000000000000000000000000006101e7565b61034b61057a366004613397565b6112e9565b61034b61058d366004613252565b6117fb565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167faff2afbf00000000000000000000000000000000000000000000000000000000148061062557507fffffffff0000000000000000000000000000000000000000000000000000000082167f0e64dd2900000000000000000000000000000000000000000000000000000000145b8061067157507fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a700000000000000000000000000000000000000000000000000000000145b92915050565b60408051602081019091526000815261068f8261180f565b60006106e860608401356106e36106a960c0870187613759565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611a3392505050565b611af7565b905073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000166340c10f196107366060860160408701613252565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b16815273ffffffffffffffffffffffffffffffffffffffff909116600482015260248101849052604401600060405180830381600087803b1580156107a357600080fd5b505af11580156107b7573d6000803e3d6000fd5b506107cc925050506060840160408501613252565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f9d228d69b5fdb8d273a2336f8fb8612d039631024ea9bf09c424a9503aa078f08360405161082a91815260200190565b60405180910390a3604080516020810190915290815292915050565b6000610888838360405161085b9291906137be565b604080519182900390912067ffffffffffffffff8716600090815260076020529190912060050190611d0b565b949350505050565b610898611d26565b61090584848080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525050604080516020808802828101820190935287825290935087925086918291850190849080828437600092019190915250611d7992505050565b50505050565b610913611d26565b61091c83610af7565b610963576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff841660048201526024015b60405180910390fd5b6109a38383838080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611f2f92505050565b505050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146109f9576040517f02b543c600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000008082163390811790935560008054909116815560405173ffffffffffffffffffffffffffffffffffffffff909216929183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610a7e611d26565b600980547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040519081527f44676b5284b809a22248eba0da87391d79098be38bb03154be88a58bf4d091749060200160405180910390a150565b6000610671600567ffffffffffffffff8416611d0b565b6040805180820190915260608082526020820152610b2b82612029565b610b3882606001356121b5565b6040516060830135815233907f696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df79060200160405180910390a26040518060400160405280610b928460200160208101906104cd9190613403565b8152602001610bdf6040805160ff7f000000000000000000000000000000000000000000000000000000000000000016602082015260609101604051602081830303815290604052905090565b905292915050565b67ffffffffffffffff8116600090815260076020526040812060609190610c1090600501612257565b90506000815167ffffffffffffffff811115610c2e57610c2e6135ce565b604051908082528060200260200182016040528015610c6157816020015b6060815260200190600190039081610c4c5790505b50905060005b8251811015610d4a5760086000848381518110610c8657610c866137ce565b602002602001015181526020019081526020016000208054610ca7906137fd565b80601f0160208091040260200160405190810160405280929190818152602001828054610cd3906137fd565b8015610d205780601f10610cf557610100808354040283529160200191610d20565b820191906000526020600020905b815481529060010190602001808311610d0357829003601f168201915b5050505050828281518110610d3757610d376137ce565b6020908102919091010152600101610c67565b509392505050565b6060610d5e6002612257565b905090565b610d6b611d26565b610d7483610af7565b610db6576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8416600482015260240161095a565b610df68282604051610dc99291906137be565b604080519182900390912067ffffffffffffffff8616600090815260076020529190912060050190612264565b610e32578282826040517f74f23c7c00000000000000000000000000000000000000000000000000000000815260040161095a93929190613899565b8267ffffffffffffffff167f52d00ee4d9bd51b40168f2afc5848837288ce258784ad914278791464b3f4d768383604051610e6e9291906138bd565b60405180910390a2505050565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915267ffffffffffffffff8216600090815260076020908152604091829020825160a08101845260028201546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff16151594820194909452600390910154808416606083015291909104909116608082015261067190612270565b67ffffffffffffffff81166000908152600760205260409020600401805460609190610f7b906137fd565b80601f0160208091040260200160405190810160405280929190818152602001828054610fa7906137fd565b8015610ff45780601f10610fc957610100808354040283529160200191610ff4565b820191906000526020600020905b815481529060010190602001808311610fd757829003601f168201915b50505050509050919050565b611008611d26565b73ffffffffffffffffffffffffffffffffffffffff8116611055576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6004805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527f02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f1684910160405180910390a15050565b606060006110e96005612257565b90506000815167ffffffffffffffff811115611107576111076135ce565b604051908082528060200260200182016040528015611130578160200160208202803683370190505b50905060005b825181101561118c57828181518110611151576111516137ce565b602002602001015182828151811061116b5761116b6137ce565b67ffffffffffffffff90921660209283029190910190910152600101611136565b5092915050565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915267ffffffffffffffff8216600090815260076020908152604091829020825160a08101845281546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff16151594820194909452600190910154808416606083015291909104909116608082015261067190612270565b60095473ffffffffffffffffffffffffffffffffffffffff1633148015906112a5575060015473ffffffffffffffffffffffffffffffffffffffff163314155b156112de576040517f8e4a23d600000000000000000000000000000000000000000000000000000000815233600482015260240161095a565b6109a3838383612322565b6112f1611d26565b60005b838110156114de576000858583818110611310576113106137ce565b90506020020160208101906113259190613403565b905061133c600567ffffffffffffffff8316612264565b61137e576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8216600482015260240161095a565b67ffffffffffffffff811660009081526007602052604081206113a390600501612257565b905060005b815181101561140f576114068282815181106113c6576113c66137ce565b6020026020010151600760008667ffffffffffffffff1667ffffffffffffffff16815260200190815260200160002060050161226490919063ffffffff16565b506001016113a8565b5067ffffffffffffffff8216600090815260076020526040812080547fffffffffffffffffffffff00000000000000000000000000000000000000000090811682556001820183905560028201805490911690556003810182905590611478600483018261310a565b600582016000818161148a8282613144565b505060405167ffffffffffffffff871681527f5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d859916945060200192506114cc915050565b60405180910390a150506001016112f4565b5060005b818110156117f45760008383838181106114fe576114fe6137ce565b905060200281019061151091906138d1565b6115199061399d565b905061152a8160600151600061240c565b6115398160800151600061240c565b806040015151600003611578576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80516115909060059067ffffffffffffffff16612549565b6115d55780516040517f1d5ad3c500000000000000000000000000000000000000000000000000000000815267ffffffffffffffff909116600482015260240161095a565b805167ffffffffffffffff16600090815260076020908152604091829020825160a08082018552606080870180518601516fffffffffffffffffffffffffffffffff90811680865263ffffffff42168689018190528351511515878b0181905284518a0151841686890181905294518b0151841660809889018190528954740100000000000000000000000000000000000000009283027fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff7001000000000000000000000000000000008087027fffffffffffffffffffffffff000000000000000000000000000000000000000094851690981788178216929092178d5592810290971760018c01558c519889018d52898e0180518d01518716808b528a8e019590955280515115158a8f018190528151909d01518716988a01899052518d0151909516979098018790526002890180549a9091029990931617179094169590951790925590920290911760038201559082015160048201906117589082613b14565b5060005b82602001515181101561179c57611794836000015184602001518381518110611787576117876137ce565b6020026020010151611f2f565b60010161175c565b507f8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c282600001518360400151846060015185608001516040516117e29493929190613c2e565b60405180910390a150506001016114e2565b5050505050565b611803611d26565b61180c81612555565b50565b61182261029360a0830160808401613252565b6118815761183660a0820160808301613252565b6040517f961c9a4f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff909116600482015260240161095a565b73ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016632cbc26bb6118cd6040840160208501613403565b60405160e083901b7fffffffff0000000000000000000000000000000000000000000000000000000016815260809190911b77ffffffffffffffff00000000000000000000000000000000166004820152602401602060405180830381865afa15801561193e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906119629190613cc7565b15611999576040517f53ad11d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6119b16119ac6040830160208401613403565b612619565b6119d16119c46040830160208401613403565b61033360a0840184613759565b611a16576119e260a0820182613759565b6040517f24eb47e500000000000000000000000000000000000000000000000000000000815260040161095a9291906138bd565b61180c611a296040830160208401613403565b826060013561273f565b60008151600003611a6557507f0000000000000000000000000000000000000000000000000000000000000000919050565b8151602014611aa257816040517f953576f700000000000000000000000000000000000000000000000000000000815260040161095a919061321d565b600082806020019051810190611ab89190613ce4565b905060ff81111561067157826040517f953576f700000000000000000000000000000000000000000000000000000000815260040161095a919061321d565b60007f000000000000000000000000000000000000000000000000000000000000000060ff168260ff1603611b2d575081610671565b7f000000000000000000000000000000000000000000000000000000000000000060ff168260ff161115611c18576000611b877f000000000000000000000000000000000000000000000000000000000000000084613d2c565b9050604d8160ff161115611bfb576040517fa9cb113d00000000000000000000000000000000000000000000000000000000815260ff80851660048301527f00000000000000000000000000000000000000000000000000000000000000001660248201526044810185905260640161095a565b611c0681600a613e65565b611c109085613e74565b915050610671565b6000611c44837f0000000000000000000000000000000000000000000000000000000000000000613d2c565b9050604d8160ff161180611c8b5750611c5e81600a613e65565b611c88907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff613e74565b84115b15611cf6576040517fa9cb113d00000000000000000000000000000000000000000000000000000000815260ff80851660048301527f00000000000000000000000000000000000000000000000000000000000000001660248201526044810185905260640161095a565b611d0181600a613e65565b6108889085613eaf565b600081815260018301602052604081205415155b9392505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314611d77576040517f2b5c74de00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b7f0000000000000000000000000000000000000000000000000000000000000000611dd0576040517f35f4a7b300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b8251811015611e66576000838281518110611df057611df06137ce565b60200260200101519050611e0e81600261278690919063ffffffff16565b15611e5d5760405173ffffffffffffffffffffffffffffffffffffffff821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b50600101611dd3565b5060005b81518110156109a3576000828281518110611e8757611e876137ce565b60200260200101519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603611ecb5750611f27565b611ed66002826127a8565b15611f255760405173ffffffffffffffffffffffffffffffffffffffff821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b600101611e6a565b8051600003611f6a576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805160208083019190912067ffffffffffffffff8416600090815260079092526040909120611f9c9060050182612549565b611fd65782826040517f393b8ad200000000000000000000000000000000000000000000000000000000815260040161095a929190613ec6565b6000818152600860205260409020611fee8382613b14565b508267ffffffffffffffff167f7d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea83604051610e6e919061321d565b61203c61029360a0830160808401613252565b6120505761183660a0820160808301613252565b73ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016632cbc26bb61209c6040840160208501613403565b60405160e083901b7fffffffff0000000000000000000000000000000000000000000000000000000016815260809190911b77ffffffffffffffff00000000000000000000000000000000166004820152602401602060405180830381865afa15801561210d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906121319190613cc7565b15612168576040517f53ad11d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61218061217b6060830160408401613252565b6127ca565b6121986121936040830160208401613403565b612849565b61180c6121ab6040830160208401613403565b8260600135612997565b6040517f9dc29fac000000000000000000000000000000000000000000000000000000008152306004820152602481018290527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690639dc29fac90604401600060405180830381600087803b15801561224357600080fd5b505af11580156117f4573d6000803e3d6000fd5b60606000611d1f836129db565b6000611d1f8383612a36565b6040805160a0810182526000808252602082018190529181018290526060810182905260808101919091526122fe82606001516fffffffffffffffffffffffffffffffff1683600001516fffffffffffffffffffffffffffffffff16846020015163ffffffff16426122e29190613ee9565b85608001516fffffffffffffffffffffffffffffffff16612b29565b6fffffffffffffffffffffffffffffffff1682525063ffffffff4216602082015290565b61232b83610af7565b61236d576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8416600482015260240161095a565b61237882600061240c565b67ffffffffffffffff8316600090815260076020526040902061239b9083612b51565b6123a681600061240c565b67ffffffffffffffff831660009081526007602052604090206123cc9060020182612b51565b7f0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b8383836040516123ff93929190613efc565b60405180910390a1505050565b8151156124d75781602001516fffffffffffffffffffffffffffffffff1682604001516fffffffffffffffffffffffffffffffff16101580612462575060408201516fffffffffffffffffffffffffffffffff16155b1561249b57816040517f8020d12400000000000000000000000000000000000000000000000000000000815260040161095a9190613f7f565b80156124d3576040517f433fc33d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050565b60408201516fffffffffffffffffffffffffffffffff16151580612510575060208201516fffffffffffffffffffffffffffffffff1615155b156124d357816040517fd68af9cc00000000000000000000000000000000000000000000000000000000815260040161095a9190613f7f565b6000611d1f8383612cf3565b3373ffffffffffffffffffffffffffffffffffffffff8216036125a4576040517fdad89dca00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b61262281610af7565b612664576040517fa9902c7e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8216600482015260240161095a565b600480546040517f83826b2b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84169281019290925233602483015273ffffffffffffffffffffffffffffffffffffffff16906383826b2b90604401602060405180830381865afa1580156126e3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906127079190613cc7565b61180c576040517f728fe07b00000000000000000000000000000000000000000000000000000000815233600482015260240161095a565b67ffffffffffffffff821660009081526007602052604090206124d390600201827f0000000000000000000000000000000000000000000000000000000000000000612d42565b6000611d1f8373ffffffffffffffffffffffffffffffffffffffff8416612a36565b6000611d1f8373ffffffffffffffffffffffffffffffffffffffff8416612cf3565b7f00000000000000000000000000000000000000000000000000000000000000001561180c576127fb6002826130c5565b61180c576040517fd0d2597600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8216600482015260240161095a565b61285281610af7565b612894576040517fa9902c7e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8216600482015260240161095a565b600480546040517fa8d87a3b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84169281019290925273ffffffffffffffffffffffffffffffffffffffff169063a8d87a3b90602401602060405180830381865afa15801561290d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906129319190613fbb565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461180c576040517f728fe07b00000000000000000000000000000000000000000000000000000000815233600482015260240161095a565b67ffffffffffffffff821660009081526007602052604090206124d390827f0000000000000000000000000000000000000000000000000000000000000000612d42565b606081600001805480602002602001604051908101604052809291908181526020018280548015610ff457602002820191906000526020600020905b815481526020019060010190808311612a175750505050509050919050565b60008181526001830160205260408120548015612b1f576000612a5a600183613ee9565b8554909150600090612a6e90600190613ee9565b9050808214612ad3576000866000018281548110612a8e57612a8e6137ce565b9060005260206000200154905080876000018481548110612ab157612ab16137ce565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080612ae457612ae4613fd8565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610671565b6000915050610671565b6000612b4885612b398486613eaf565b612b439087614007565b6130f4565b95945050505050565b8154600090612b7a90700100000000000000000000000000000000900463ffffffff1642613ee9565b90508015612c1c5760018301548354612bc2916fffffffffffffffffffffffffffffffff80821692811691859170010000000000000000000000000000000090910416612b29565b83546fffffffffffffffffffffffffffffffff919091167fffffffffffffffffffffffff0000000000000000000000000000000000000000909116177001000000000000000000000000000000004263ffffffff16021783555b60208201518354612c42916fffffffffffffffffffffffffffffffff90811691166130f4565b83548351151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00ffffffff000000000000000000000000000000009091166fffffffffffffffffffffffffffffffff92831617178455602083015160408085015183167001000000000000000000000000000000000291909216176001850155517f9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c19906123ff908490613f7f565b6000818152600183016020526040812054612d3a57508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610671565b506000610671565b825474010000000000000000000000000000000000000000900460ff161580612d69575081155b15612d7357505050565b825460018401546fffffffffffffffffffffffffffffffff80831692911690600090612db990700100000000000000000000000000000000900463ffffffff1642613ee9565b90508015612e795781831115612dfb576040517f9725942a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001860154612e359083908590849070010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff16612b29565b86547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff167001000000000000000000000000000000004263ffffffff160217875592505b84821015612f305773ffffffffffffffffffffffffffffffffffffffff8416612ed8576040517ff94ebcd1000000000000000000000000000000000000000000000000000000008152600481018390526024810186905260440161095a565b6040517f1a76572a000000000000000000000000000000000000000000000000000000008152600481018390526024810186905273ffffffffffffffffffffffffffffffffffffffff8516604482015260640161095a565b848310156130435760018681015470010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff16906000908290612f749082613ee9565b612f7e878a613ee9565b612f889190614007565b612f929190613e74565b905073ffffffffffffffffffffffffffffffffffffffff8616612feb576040517f15279c08000000000000000000000000000000000000000000000000000000008152600481018290526024810186905260440161095a565b6040517fd0c8d23a000000000000000000000000000000000000000000000000000000008152600481018290526024810186905273ffffffffffffffffffffffffffffffffffffffff8716604482015260640161095a565b61304d8584613ee9565b86547fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff82161787556040518681529093507f1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a9060200160405180910390a1505050505050565b73ffffffffffffffffffffffffffffffffffffffff811660009081526001830160205260408120541515611d1f565b60008183106131035781611d1f565b5090919050565b508054613116906137fd565b6000825580601f10613126575050565b601f01602090049060005260206000209081019061180c919061315e565b508054600082559060005260206000209081019061180c91905b5b80821115613173576000815560010161315f565b5090565b60006020828403121561318957600080fd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114611d1f57600080fd5b6000815180845260005b818110156131df576020818501810151868301820152016131c3565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081526000611d1f60208301846131b9565b73ffffffffffffffffffffffffffffffffffffffff8116811461180c57600080fd5b60006020828403121561326457600080fd5b8135611d1f81613230565b60006020828403121561328157600080fd5b813567ffffffffffffffff81111561329857600080fd5b82016101008185031215611d1f57600080fd5b803567ffffffffffffffff811681146132c357600080fd5b919050565b6000806000604084860312156132dd57600080fd5b6132e6846132ab565b9250602084013567ffffffffffffffff8082111561330357600080fd5b818601915086601f83011261331757600080fd5b81358181111561332657600080fd5b87602082850101111561333857600080fd5b6020830194508093505050509250925092565b60008083601f84011261335d57600080fd5b50813567ffffffffffffffff81111561337557600080fd5b6020830191508360208260051b850101111561339057600080fd5b9250929050565b600080600080604085870312156133ad57600080fd5b843567ffffffffffffffff808211156133c557600080fd5b6133d18883890161334b565b909650945060208701359150808211156133ea57600080fd5b506133f78782880161334b565b95989497509550505050565b60006020828403121561341557600080fd5b611d1f826132ab565b60006020828403121561343057600080fd5b813567ffffffffffffffff81111561344757600080fd5b820160a08185031215611d1f57600080fd5b60208152600082516040602084015261347560608401826131b9565b905060208401517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0848303016040850152612b4882826131b9565b600060208083016020845280855180835260408601915060408160051b87010192506020870160005b82811015613525577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08886030184526135138583516131b9565b945092850192908501906001016134d9565b5092979650505050505050565b6020808252825182820181905260009190848201906040850190845b8181101561358057835173ffffffffffffffffffffffffffffffffffffffff168352928401929184019160010161354e565b50909695505050505050565b6020808252825182820181905260009190848201906040850190845b8181101561358057835167ffffffffffffffff16835292840192918401916001016135a8565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160a0810167ffffffffffffffff81118282101715613620576136206135ce565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561366d5761366d6135ce565b604052919050565b801515811461180c57600080fd5b80356fffffffffffffffffffffffffffffffff811681146132c357600080fd5b6000606082840312156136b557600080fd5b6040516060810181811067ffffffffffffffff821117156136d8576136d86135ce565b60405290508082356136e981613675565b81526136f760208401613683565b602082015261370860408401613683565b60408201525092915050565b600080600060e0848603121561372957600080fd5b613732846132ab565b925061374185602086016136a3565b915061375085608086016136a3565b90509250925092565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261378e57600080fd5b83018035915067ffffffffffffffff8211156137a957600080fd5b60200191503681900382131561339057600080fd5b8183823760009101908152919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600181811c9082168061381157607f821691505b60208210810361384a577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b67ffffffffffffffff84168152604060208201526000612b48604083018486613850565b602081526000610888602083018486613850565b600082357ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffee183360301811261390557600080fd5b9190910192915050565b600082601f83011261392057600080fd5b813567ffffffffffffffff81111561393a5761393a6135ce565b61396b60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601613626565b81815284602083860101111561398057600080fd5b816020850160208301376000918101602001919091529392505050565b600061012082360312156139b057600080fd5b6139b86135fd565b6139c1836132ab565b815260208084013567ffffffffffffffff808211156139df57600080fd5b9085019036601f8301126139f257600080fd5b813581811115613a0457613a046135ce565b8060051b613a13858201613626565b9182528381018501918581019036841115613a2d57600080fd5b86860192505b83831015613a6957823585811115613a4b5760008081fd5b613a593689838a010161390f565b8352509186019190860190613a33565b8087890152505050506040860135925080831115613a8657600080fd5b5050613a943682860161390f565b604083015250613aa736606085016136a3565b6060820152613ab93660c085016136a3565b608082015292915050565b601f8211156109a3576000816000526020600020601f850160051c81016020861015613aed5750805b601f850160051c820191505b81811015613b0c57828155600101613af9565b505050505050565b815167ffffffffffffffff811115613b2e57613b2e6135ce565b613b4281613b3c84546137fd565b84613ac4565b602080601f831160018114613b955760008415613b5f5750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555613b0c565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015613be257888601518255948401946001909101908401613bc3565b5085821015613c1e57878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b600061010067ffffffffffffffff87168352806020840152613c52818401876131b9565b8551151560408581019190915260208701516fffffffffffffffffffffffffffffffff9081166060870152908701511660808501529150613c909050565b8251151560a083015260208301516fffffffffffffffffffffffffffffffff90811660c084015260408401511660e0830152612b48565b600060208284031215613cd957600080fd5b8151611d1f81613675565b600060208284031215613cf657600080fd5b5051919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60ff828116828216039081111561067157610671613cfd565b600181815b80851115613d9e57817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04821115613d8457613d84613cfd565b80851615613d9157918102915b93841c9390800290613d4a565b509250929050565b600082613db557506001610671565b81613dc257506000610671565b8160018114613dd85760028114613de257613dfe565b6001915050610671565b60ff841115613df357613df3613cfd565b50506001821b610671565b5060208310610133831016604e8410600b8410161715613e21575081810a610671565b613e2b8383613d45565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04821115613e5d57613e5d613cfd565b029392505050565b6000611d1f60ff841683613da6565b600082613eaa577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b808202811582820484141761067157610671613cfd565b67ffffffffffffffff8316815260406020820152600061088860408301846131b9565b8181038181111561067157610671613cfd565b67ffffffffffffffff8416815260e08101613f4860208301858051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b82511515608083015260208301516fffffffffffffffffffffffffffffffff90811660a084015260408401511660c0830152610888565b6060810161067182848051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b600060208284031215613fcd57600080fd5b8151611d1f81613230565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b8082018082111561067157610671613cfd56fea164736f6c6343000818000a",
}

var BurnWithFromMintTokenPoolABI = BurnWithFromMintTokenPoolMetaData.ABI

var BurnWithFromMintTokenPoolBin = BurnWithFromMintTokenPoolMetaData.Bin

func DeployBurnWithFromMintTokenPool(auth *bind.TransactOpts, backend bind.ContractBackend, token common.Address, localTokenDecimals uint8, allowlist []common.Address, rmnProxy common.Address, router common.Address) (common.Address, *generated.Transaction, *BurnWithFromMintTokenPool, error) {
	parsed, err := BurnWithFromMintTokenPoolMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}
	if generated.IsZKSync(backend) {
		address, ethTx, contractBind, _ := generated.DeployContract(auth, parsed, common.FromHex(BurnWithFromMintTokenPoolZKBin), backend, token, localTokenDecimals, allowlist, rmnProxy, router)
		contractReturn := &BurnWithFromMintTokenPool{address: address, abi: *parsed, BurnWithFromMintTokenPoolCaller: BurnWithFromMintTokenPoolCaller{contract: contractBind}, BurnWithFromMintTokenPoolTransactor: BurnWithFromMintTokenPoolTransactor{contract: contractBind}, BurnWithFromMintTokenPoolFilterer: BurnWithFromMintTokenPoolFilterer{contract: contractBind}}
		return address, ethTx, contractReturn, err
	}
	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BurnWithFromMintTokenPoolBin), backend, token, localTokenDecimals, allowlist, rmnProxy, router)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, &generated.Transaction{Transaction: tx, HashZks: tx.Hash()}, &BurnWithFromMintTokenPool{address: address, abi: *parsed, BurnWithFromMintTokenPoolCaller: BurnWithFromMintTokenPoolCaller{contract: contract}, BurnWithFromMintTokenPoolTransactor: BurnWithFromMintTokenPoolTransactor{contract: contract}, BurnWithFromMintTokenPoolFilterer: BurnWithFromMintTokenPoolFilterer{contract: contract}}, nil
}

type BurnWithFromMintTokenPool struct {
	address common.Address
	abi     abi.ABI
	BurnWithFromMintTokenPoolCaller
	BurnWithFromMintTokenPoolTransactor
	BurnWithFromMintTokenPoolFilterer
}

type BurnWithFromMintTokenPoolCaller struct {
	contract *bind.BoundContract
}

type BurnWithFromMintTokenPoolTransactor struct {
	contract *bind.BoundContract
}

type BurnWithFromMintTokenPoolFilterer struct {
	contract *bind.BoundContract
}

type BurnWithFromMintTokenPoolSession struct {
	Contract     *BurnWithFromMintTokenPool
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type BurnWithFromMintTokenPoolCallerSession struct {
	Contract *BurnWithFromMintTokenPoolCaller
	CallOpts bind.CallOpts
}

type BurnWithFromMintTokenPoolTransactorSession struct {
	Contract     *BurnWithFromMintTokenPoolTransactor
	TransactOpts bind.TransactOpts
}

type BurnWithFromMintTokenPoolRaw struct {
	Contract *BurnWithFromMintTokenPool
}

type BurnWithFromMintTokenPoolCallerRaw struct {
	Contract *BurnWithFromMintTokenPoolCaller
}

type BurnWithFromMintTokenPoolTransactorRaw struct {
	Contract *BurnWithFromMintTokenPoolTransactor
}

func NewBurnWithFromMintTokenPool(address common.Address, backend bind.ContractBackend) (*BurnWithFromMintTokenPool, error) {
	abi, err := abi.JSON(strings.NewReader(BurnWithFromMintTokenPoolABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindBurnWithFromMintTokenPool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPool{address: address, abi: abi, BurnWithFromMintTokenPoolCaller: BurnWithFromMintTokenPoolCaller{contract: contract}, BurnWithFromMintTokenPoolTransactor: BurnWithFromMintTokenPoolTransactor{contract: contract}, BurnWithFromMintTokenPoolFilterer: BurnWithFromMintTokenPoolFilterer{contract: contract}}, nil
}

func NewBurnWithFromMintTokenPoolCaller(address common.Address, caller bind.ContractCaller) (*BurnWithFromMintTokenPoolCaller, error) {
	contract, err := bindBurnWithFromMintTokenPool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolCaller{contract: contract}, nil
}

func NewBurnWithFromMintTokenPoolTransactor(address common.Address, transactor bind.ContractTransactor) (*BurnWithFromMintTokenPoolTransactor, error) {
	contract, err := bindBurnWithFromMintTokenPool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolTransactor{contract: contract}, nil
}

func NewBurnWithFromMintTokenPoolFilterer(address common.Address, filterer bind.ContractFilterer) (*BurnWithFromMintTokenPoolFilterer, error) {
	contract, err := bindBurnWithFromMintTokenPool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolFilterer{contract: contract}, nil
}

func bindBurnWithFromMintTokenPool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BurnWithFromMintTokenPoolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BurnWithFromMintTokenPool.Contract.BurnWithFromMintTokenPoolCaller.contract.Call(opts, result, method, params...)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.BurnWithFromMintTokenPoolTransactor.contract.Transfer(opts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.BurnWithFromMintTokenPoolTransactor.contract.Transact(opts, method, params...)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BurnWithFromMintTokenPool.Contract.contract.Call(opts, result, method, params...)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.contract.Transfer(opts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.contract.Transact(opts, method, params...)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCaller) GetAllowList(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPool.contract.Call(opts, &out, "getAllowList")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) GetAllowList() ([]common.Address, error) {
	return _BurnWithFromMintTokenPool.Contract.GetAllowList(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerSession) GetAllowList() ([]common.Address, error) {
	return _BurnWithFromMintTokenPool.Contract.GetAllowList(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCaller) GetAllowListEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPool.contract.Call(opts, &out, "getAllowListEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) GetAllowListEnabled() (bool, error) {
	return _BurnWithFromMintTokenPool.Contract.GetAllowListEnabled(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerSession) GetAllowListEnabled() (bool, error) {
	return _BurnWithFromMintTokenPool.Contract.GetAllowListEnabled(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCaller) GetCurrentInboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPool.contract.Call(opts, &out, "getCurrentInboundRateLimiterState", remoteChainSelector)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) GetCurrentInboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnWithFromMintTokenPool.Contract.GetCurrentInboundRateLimiterState(&_BurnWithFromMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerSession) GetCurrentInboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnWithFromMintTokenPool.Contract.GetCurrentInboundRateLimiterState(&_BurnWithFromMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCaller) GetCurrentOutboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPool.contract.Call(opts, &out, "getCurrentOutboundRateLimiterState", remoteChainSelector)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) GetCurrentOutboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnWithFromMintTokenPool.Contract.GetCurrentOutboundRateLimiterState(&_BurnWithFromMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerSession) GetCurrentOutboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnWithFromMintTokenPool.Contract.GetCurrentOutboundRateLimiterState(&_BurnWithFromMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCaller) GetRateLimitAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPool.contract.Call(opts, &out, "getRateLimitAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) GetRateLimitAdmin() (common.Address, error) {
	return _BurnWithFromMintTokenPool.Contract.GetRateLimitAdmin(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerSession) GetRateLimitAdmin() (common.Address, error) {
	return _BurnWithFromMintTokenPool.Contract.GetRateLimitAdmin(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCaller) GetRemotePools(opts *bind.CallOpts, remoteChainSelector uint64) ([][]byte, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPool.contract.Call(opts, &out, "getRemotePools", remoteChainSelector)

	if err != nil {
		return *new([][]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][]byte)).(*[][]byte)

	return out0, err

}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) GetRemotePools(remoteChainSelector uint64) ([][]byte, error) {
	return _BurnWithFromMintTokenPool.Contract.GetRemotePools(&_BurnWithFromMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerSession) GetRemotePools(remoteChainSelector uint64) ([][]byte, error) {
	return _BurnWithFromMintTokenPool.Contract.GetRemotePools(&_BurnWithFromMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCaller) GetRemoteToken(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPool.contract.Call(opts, &out, "getRemoteToken", remoteChainSelector)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) GetRemoteToken(remoteChainSelector uint64) ([]byte, error) {
	return _BurnWithFromMintTokenPool.Contract.GetRemoteToken(&_BurnWithFromMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerSession) GetRemoteToken(remoteChainSelector uint64) ([]byte, error) {
	return _BurnWithFromMintTokenPool.Contract.GetRemoteToken(&_BurnWithFromMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCaller) GetRmnProxy(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPool.contract.Call(opts, &out, "getRmnProxy")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) GetRmnProxy() (common.Address, error) {
	return _BurnWithFromMintTokenPool.Contract.GetRmnProxy(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerSession) GetRmnProxy() (common.Address, error) {
	return _BurnWithFromMintTokenPool.Contract.GetRmnProxy(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCaller) GetRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPool.contract.Call(opts, &out, "getRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) GetRouter() (common.Address, error) {
	return _BurnWithFromMintTokenPool.Contract.GetRouter(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerSession) GetRouter() (common.Address, error) {
	return _BurnWithFromMintTokenPool.Contract.GetRouter(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCaller) GetSupportedChains(opts *bind.CallOpts) ([]uint64, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPool.contract.Call(opts, &out, "getSupportedChains")

	if err != nil {
		return *new([]uint64), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint64)).(*[]uint64)

	return out0, err

}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) GetSupportedChains() ([]uint64, error) {
	return _BurnWithFromMintTokenPool.Contract.GetSupportedChains(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerSession) GetSupportedChains() ([]uint64, error) {
	return _BurnWithFromMintTokenPool.Contract.GetSupportedChains(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCaller) GetToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPool.contract.Call(opts, &out, "getToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) GetToken() (common.Address, error) {
	return _BurnWithFromMintTokenPool.Contract.GetToken(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerSession) GetToken() (common.Address, error) {
	return _BurnWithFromMintTokenPool.Contract.GetToken(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCaller) GetTokenDecimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPool.contract.Call(opts, &out, "getTokenDecimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) GetTokenDecimals() (uint8, error) {
	return _BurnWithFromMintTokenPool.Contract.GetTokenDecimals(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerSession) GetTokenDecimals() (uint8, error) {
	return _BurnWithFromMintTokenPool.Contract.GetTokenDecimals(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCaller) IsRemotePool(opts *bind.CallOpts, remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPool.contract.Call(opts, &out, "isRemotePool", remoteChainSelector, remotePoolAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) IsRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	return _BurnWithFromMintTokenPool.Contract.IsRemotePool(&_BurnWithFromMintTokenPool.CallOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerSession) IsRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	return _BurnWithFromMintTokenPool.Contract.IsRemotePool(&_BurnWithFromMintTokenPool.CallOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCaller) IsSupportedChain(opts *bind.CallOpts, remoteChainSelector uint64) (bool, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPool.contract.Call(opts, &out, "isSupportedChain", remoteChainSelector)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) IsSupportedChain(remoteChainSelector uint64) (bool, error) {
	return _BurnWithFromMintTokenPool.Contract.IsSupportedChain(&_BurnWithFromMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerSession) IsSupportedChain(remoteChainSelector uint64) (bool, error) {
	return _BurnWithFromMintTokenPool.Contract.IsSupportedChain(&_BurnWithFromMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCaller) IsSupportedToken(opts *bind.CallOpts, token common.Address) (bool, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPool.contract.Call(opts, &out, "isSupportedToken", token)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) IsSupportedToken(token common.Address) (bool, error) {
	return _BurnWithFromMintTokenPool.Contract.IsSupportedToken(&_BurnWithFromMintTokenPool.CallOpts, token)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerSession) IsSupportedToken(token common.Address) (bool, error) {
	return _BurnWithFromMintTokenPool.Contract.IsSupportedToken(&_BurnWithFromMintTokenPool.CallOpts, token)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPool.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) Owner() (common.Address, error) {
	return _BurnWithFromMintTokenPool.Contract.Owner(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerSession) Owner() (common.Address, error) {
	return _BurnWithFromMintTokenPool.Contract.Owner(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPool.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BurnWithFromMintTokenPool.Contract.SupportsInterface(&_BurnWithFromMintTokenPool.CallOpts, interfaceId)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BurnWithFromMintTokenPool.Contract.SupportsInterface(&_BurnWithFromMintTokenPool.CallOpts, interfaceId)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _BurnWithFromMintTokenPool.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) TypeAndVersion() (string, error) {
	return _BurnWithFromMintTokenPool.Contract.TypeAndVersion(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolCallerSession) TypeAndVersion() (string, error) {
	return _BurnWithFromMintTokenPool.Contract.TypeAndVersion(&_BurnWithFromMintTokenPool.CallOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.contract.Transact(opts, "acceptOwnership")
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) AcceptOwnership() (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.AcceptOwnership(&_BurnWithFromMintTokenPool.TransactOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.AcceptOwnership(&_BurnWithFromMintTokenPool.TransactOpts)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactor) AddRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.contract.Transact(opts, "addRemotePool", remoteChainSelector, remotePoolAddress)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) AddRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.AddRemotePool(&_BurnWithFromMintTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactorSession) AddRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.AddRemotePool(&_BurnWithFromMintTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactor) ApplyAllowListUpdates(opts *bind.TransactOpts, removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.contract.Transact(opts, "applyAllowListUpdates", removes, adds)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.ApplyAllowListUpdates(&_BurnWithFromMintTokenPool.TransactOpts, removes, adds)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactorSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.ApplyAllowListUpdates(&_BurnWithFromMintTokenPool.TransactOpts, removes, adds)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactor) ApplyChainUpdates(opts *bind.TransactOpts, remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.contract.Transact(opts, "applyChainUpdates", remoteChainSelectorsToRemove, chainsToAdd)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) ApplyChainUpdates(remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.ApplyChainUpdates(&_BurnWithFromMintTokenPool.TransactOpts, remoteChainSelectorsToRemove, chainsToAdd)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactorSession) ApplyChainUpdates(remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.ApplyChainUpdates(&_BurnWithFromMintTokenPool.TransactOpts, remoteChainSelectorsToRemove, chainsToAdd)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactor) LockOrBurn(opts *bind.TransactOpts, lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.contract.Transact(opts, "lockOrBurn", lockOrBurnIn)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) LockOrBurn(lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.LockOrBurn(&_BurnWithFromMintTokenPool.TransactOpts, lockOrBurnIn)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactorSession) LockOrBurn(lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.LockOrBurn(&_BurnWithFromMintTokenPool.TransactOpts, lockOrBurnIn)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactor) ReleaseOrMint(opts *bind.TransactOpts, releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.contract.Transact(opts, "releaseOrMint", releaseOrMintIn)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) ReleaseOrMint(releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.ReleaseOrMint(&_BurnWithFromMintTokenPool.TransactOpts, releaseOrMintIn)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactorSession) ReleaseOrMint(releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.ReleaseOrMint(&_BurnWithFromMintTokenPool.TransactOpts, releaseOrMintIn)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactor) RemoveRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.contract.Transact(opts, "removeRemotePool", remoteChainSelector, remotePoolAddress)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) RemoveRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.RemoveRemotePool(&_BurnWithFromMintTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactorSession) RemoveRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.RemoveRemotePool(&_BurnWithFromMintTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactor) SetChainRateLimiterConfig(opts *bind.TransactOpts, remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.contract.Transact(opts, "setChainRateLimiterConfig", remoteChainSelector, outboundConfig, inboundConfig)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) SetChainRateLimiterConfig(remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.SetChainRateLimiterConfig(&_BurnWithFromMintTokenPool.TransactOpts, remoteChainSelector, outboundConfig, inboundConfig)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactorSession) SetChainRateLimiterConfig(remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.SetChainRateLimiterConfig(&_BurnWithFromMintTokenPool.TransactOpts, remoteChainSelector, outboundConfig, inboundConfig)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactor) SetRateLimitAdmin(opts *bind.TransactOpts, rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.contract.Transact(opts, "setRateLimitAdmin", rateLimitAdmin)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) SetRateLimitAdmin(rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.SetRateLimitAdmin(&_BurnWithFromMintTokenPool.TransactOpts, rateLimitAdmin)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactorSession) SetRateLimitAdmin(rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.SetRateLimitAdmin(&_BurnWithFromMintTokenPool.TransactOpts, rateLimitAdmin)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactor) SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.contract.Transact(opts, "setRouter", newRouter)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.SetRouter(&_BurnWithFromMintTokenPool.TransactOpts, newRouter)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactorSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.SetRouter(&_BurnWithFromMintTokenPool.TransactOpts, newRouter)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.contract.Transact(opts, "transferOwnership", to)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.TransferOwnership(&_BurnWithFromMintTokenPool.TransactOpts, to)
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _BurnWithFromMintTokenPool.Contract.TransferOwnership(&_BurnWithFromMintTokenPool.TransactOpts, to)
}

type BurnWithFromMintTokenPoolAllowListAddIterator struct {
	Event *BurnWithFromMintTokenPoolAllowListAdd

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolAllowListAddIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolAllowListAdd)
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
		it.Event = new(BurnWithFromMintTokenPoolAllowListAdd)
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

func (it *BurnWithFromMintTokenPoolAllowListAddIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolAllowListAddIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolAllowListAdd struct {
	Sender common.Address
	Raw    types.Log
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) FilterAllowListAdd(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolAllowListAddIterator, error) {

	logs, sub, err := _BurnWithFromMintTokenPool.contract.FilterLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolAllowListAddIterator{contract: _BurnWithFromMintTokenPool.contract, event: "AllowListAdd", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAllowListAdd) (event.Subscription, error) {

	logs, sub, err := _BurnWithFromMintTokenPool.contract.WatchLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolAllowListAdd)
				if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
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

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) ParseAllowListAdd(log types.Log) (*BurnWithFromMintTokenPoolAllowListAdd, error) {
	event := new(BurnWithFromMintTokenPoolAllowListAdd)
	if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolAllowListRemoveIterator struct {
	Event *BurnWithFromMintTokenPoolAllowListRemove

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolAllowListRemoveIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolAllowListRemove)
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
		it.Event = new(BurnWithFromMintTokenPoolAllowListRemove)
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

func (it *BurnWithFromMintTokenPoolAllowListRemoveIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolAllowListRemoveIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolAllowListRemove struct {
	Sender common.Address
	Raw    types.Log
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) FilterAllowListRemove(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolAllowListRemoveIterator, error) {

	logs, sub, err := _BurnWithFromMintTokenPool.contract.FilterLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolAllowListRemoveIterator{contract: _BurnWithFromMintTokenPool.contract, event: "AllowListRemove", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAllowListRemove) (event.Subscription, error) {

	logs, sub, err := _BurnWithFromMintTokenPool.contract.WatchLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolAllowListRemove)
				if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
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

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) ParseAllowListRemove(log types.Log) (*BurnWithFromMintTokenPoolAllowListRemove, error) {
	event := new(BurnWithFromMintTokenPoolAllowListRemove)
	if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolBurnedIterator struct {
	Event *BurnWithFromMintTokenPoolBurned

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolBurnedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolBurned)
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
		it.Event = new(BurnWithFromMintTokenPoolBurned)
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

func (it *BurnWithFromMintTokenPoolBurnedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolBurnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolBurned struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) FilterBurned(opts *bind.FilterOpts, sender []common.Address) (*BurnWithFromMintTokenPoolBurnedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPool.contract.FilterLogs(opts, "Burned", senderRule)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolBurnedIterator{contract: _BurnWithFromMintTokenPool.contract, event: "Burned", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) WatchBurned(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolBurned, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPool.contract.WatchLogs(opts, "Burned", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolBurned)
				if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "Burned", log); err != nil {
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

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) ParseBurned(log types.Log) (*BurnWithFromMintTokenPoolBurned, error) {
	event := new(BurnWithFromMintTokenPoolBurned)
	if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "Burned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolChainAddedIterator struct {
	Event *BurnWithFromMintTokenPoolChainAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolChainAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolChainAdded)
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
		it.Event = new(BurnWithFromMintTokenPoolChainAdded)
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

func (it *BurnWithFromMintTokenPoolChainAddedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolChainAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolChainAdded struct {
	RemoteChainSelector       uint64
	RemoteToken               []byte
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
	Raw                       types.Log
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) FilterChainAdded(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolChainAddedIterator, error) {

	logs, sub, err := _BurnWithFromMintTokenPool.contract.FilterLogs(opts, "ChainAdded")
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolChainAddedIterator{contract: _BurnWithFromMintTokenPool.contract, event: "ChainAdded", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) WatchChainAdded(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolChainAdded) (event.Subscription, error) {

	logs, sub, err := _BurnWithFromMintTokenPool.contract.WatchLogs(opts, "ChainAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolChainAdded)
				if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "ChainAdded", log); err != nil {
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

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) ParseChainAdded(log types.Log) (*BurnWithFromMintTokenPoolChainAdded, error) {
	event := new(BurnWithFromMintTokenPoolChainAdded)
	if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "ChainAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolChainConfiguredIterator struct {
	Event *BurnWithFromMintTokenPoolChainConfigured

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolChainConfiguredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolChainConfigured)
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
		it.Event = new(BurnWithFromMintTokenPoolChainConfigured)
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

func (it *BurnWithFromMintTokenPoolChainConfiguredIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolChainConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolChainConfigured struct {
	RemoteChainSelector       uint64
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
	Raw                       types.Log
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) FilterChainConfigured(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolChainConfiguredIterator, error) {

	logs, sub, err := _BurnWithFromMintTokenPool.contract.FilterLogs(opts, "ChainConfigured")
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolChainConfiguredIterator{contract: _BurnWithFromMintTokenPool.contract, event: "ChainConfigured", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) WatchChainConfigured(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolChainConfigured) (event.Subscription, error) {

	logs, sub, err := _BurnWithFromMintTokenPool.contract.WatchLogs(opts, "ChainConfigured")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolChainConfigured)
				if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "ChainConfigured", log); err != nil {
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

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) ParseChainConfigured(log types.Log) (*BurnWithFromMintTokenPoolChainConfigured, error) {
	event := new(BurnWithFromMintTokenPoolChainConfigured)
	if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "ChainConfigured", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolChainRemovedIterator struct {
	Event *BurnWithFromMintTokenPoolChainRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolChainRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolChainRemoved)
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
		it.Event = new(BurnWithFromMintTokenPoolChainRemoved)
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

func (it *BurnWithFromMintTokenPoolChainRemovedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolChainRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolChainRemoved struct {
	RemoteChainSelector uint64
	Raw                 types.Log
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) FilterChainRemoved(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolChainRemovedIterator, error) {

	logs, sub, err := _BurnWithFromMintTokenPool.contract.FilterLogs(opts, "ChainRemoved")
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolChainRemovedIterator{contract: _BurnWithFromMintTokenPool.contract, event: "ChainRemoved", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolChainRemoved) (event.Subscription, error) {

	logs, sub, err := _BurnWithFromMintTokenPool.contract.WatchLogs(opts, "ChainRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolChainRemoved)
				if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
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

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) ParseChainRemoved(log types.Log) (*BurnWithFromMintTokenPoolChainRemoved, error) {
	event := new(BurnWithFromMintTokenPoolChainRemoved)
	if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolConfigChangedIterator struct {
	Event *BurnWithFromMintTokenPoolConfigChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolConfigChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolConfigChanged)
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
		it.Event = new(BurnWithFromMintTokenPoolConfigChanged)
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

func (it *BurnWithFromMintTokenPoolConfigChangedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolConfigChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolConfigChanged struct {
	Config RateLimiterConfig
	Raw    types.Log
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) FilterConfigChanged(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolConfigChangedIterator, error) {

	logs, sub, err := _BurnWithFromMintTokenPool.contract.FilterLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolConfigChangedIterator{contract: _BurnWithFromMintTokenPool.contract, event: "ConfigChanged", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolConfigChanged) (event.Subscription, error) {

	logs, sub, err := _BurnWithFromMintTokenPool.contract.WatchLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolConfigChanged)
				if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
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

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) ParseConfigChanged(log types.Log) (*BurnWithFromMintTokenPoolConfigChanged, error) {
	event := new(BurnWithFromMintTokenPoolConfigChanged)
	if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolLockedIterator struct {
	Event *BurnWithFromMintTokenPoolLocked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolLockedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolLocked)
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
		it.Event = new(BurnWithFromMintTokenPoolLocked)
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

func (it *BurnWithFromMintTokenPoolLockedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolLockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolLocked struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) FilterLocked(opts *bind.FilterOpts, sender []common.Address) (*BurnWithFromMintTokenPoolLockedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPool.contract.FilterLogs(opts, "Locked", senderRule)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolLockedIterator{contract: _BurnWithFromMintTokenPool.contract, event: "Locked", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) WatchLocked(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolLocked, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPool.contract.WatchLogs(opts, "Locked", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolLocked)
				if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "Locked", log); err != nil {
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

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) ParseLocked(log types.Log) (*BurnWithFromMintTokenPoolLocked, error) {
	event := new(BurnWithFromMintTokenPoolLocked)
	if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "Locked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolMintedIterator struct {
	Event *BurnWithFromMintTokenPoolMinted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolMintedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolMinted)
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
		it.Event = new(BurnWithFromMintTokenPoolMinted)
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

func (it *BurnWithFromMintTokenPoolMintedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolMinted struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) FilterMinted(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*BurnWithFromMintTokenPoolMintedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPool.contract.FilterLogs(opts, "Minted", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolMintedIterator{contract: _BurnWithFromMintTokenPool.contract, event: "Minted", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) WatchMinted(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolMinted, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPool.contract.WatchLogs(opts, "Minted", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolMinted)
				if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "Minted", log); err != nil {
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

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) ParseMinted(log types.Log) (*BurnWithFromMintTokenPoolMinted, error) {
	event := new(BurnWithFromMintTokenPoolMinted)
	if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "Minted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolOwnershipTransferRequestedIterator struct {
	Event *BurnWithFromMintTokenPoolOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolOwnershipTransferRequested)
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
		it.Event = new(BurnWithFromMintTokenPoolOwnershipTransferRequested)
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

func (it *BurnWithFromMintTokenPoolOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnWithFromMintTokenPoolOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPool.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolOwnershipTransferRequestedIterator{contract: _BurnWithFromMintTokenPool.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPool.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolOwnershipTransferRequested)
				if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) ParseOwnershipTransferRequested(log types.Log) (*BurnWithFromMintTokenPoolOwnershipTransferRequested, error) {
	event := new(BurnWithFromMintTokenPoolOwnershipTransferRequested)
	if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolOwnershipTransferredIterator struct {
	Event *BurnWithFromMintTokenPoolOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolOwnershipTransferred)
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
		it.Event = new(BurnWithFromMintTokenPoolOwnershipTransferred)
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

func (it *BurnWithFromMintTokenPoolOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnWithFromMintTokenPoolOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPool.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolOwnershipTransferredIterator{contract: _BurnWithFromMintTokenPool.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPool.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolOwnershipTransferred)
				if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) ParseOwnershipTransferred(log types.Log) (*BurnWithFromMintTokenPoolOwnershipTransferred, error) {
	event := new(BurnWithFromMintTokenPoolOwnershipTransferred)
	if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolRateLimitAdminSetIterator struct {
	Event *BurnWithFromMintTokenPoolRateLimitAdminSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolRateLimitAdminSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolRateLimitAdminSet)
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
		it.Event = new(BurnWithFromMintTokenPoolRateLimitAdminSet)
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

func (it *BurnWithFromMintTokenPoolRateLimitAdminSetIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolRateLimitAdminSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolRateLimitAdminSet struct {
	RateLimitAdmin common.Address
	Raw            types.Log
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) FilterRateLimitAdminSet(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolRateLimitAdminSetIterator, error) {

	logs, sub, err := _BurnWithFromMintTokenPool.contract.FilterLogs(opts, "RateLimitAdminSet")
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolRateLimitAdminSetIterator{contract: _BurnWithFromMintTokenPool.contract, event: "RateLimitAdminSet", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) WatchRateLimitAdminSet(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolRateLimitAdminSet) (event.Subscription, error) {

	logs, sub, err := _BurnWithFromMintTokenPool.contract.WatchLogs(opts, "RateLimitAdminSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolRateLimitAdminSet)
				if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "RateLimitAdminSet", log); err != nil {
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

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) ParseRateLimitAdminSet(log types.Log) (*BurnWithFromMintTokenPoolRateLimitAdminSet, error) {
	event := new(BurnWithFromMintTokenPoolRateLimitAdminSet)
	if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "RateLimitAdminSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolReleasedIterator struct {
	Event *BurnWithFromMintTokenPoolReleased

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolReleasedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolReleased)
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
		it.Event = new(BurnWithFromMintTokenPoolReleased)
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

func (it *BurnWithFromMintTokenPoolReleasedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolReleasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolReleased struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) FilterReleased(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*BurnWithFromMintTokenPoolReleasedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPool.contract.FilterLogs(opts, "Released", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolReleasedIterator{contract: _BurnWithFromMintTokenPool.contract, event: "Released", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) WatchReleased(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolReleased, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPool.contract.WatchLogs(opts, "Released", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolReleased)
				if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "Released", log); err != nil {
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

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) ParseReleased(log types.Log) (*BurnWithFromMintTokenPoolReleased, error) {
	event := new(BurnWithFromMintTokenPoolReleased)
	if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "Released", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolRemotePoolAddedIterator struct {
	Event *BurnWithFromMintTokenPoolRemotePoolAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolRemotePoolAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolRemotePoolAdded)
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
		it.Event = new(BurnWithFromMintTokenPoolRemotePoolAdded)
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

func (it *BurnWithFromMintTokenPoolRemotePoolAddedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolRemotePoolAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolRemotePoolAdded struct {
	RemoteChainSelector uint64
	RemotePoolAddress   []byte
	Raw                 types.Log
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) FilterRemotePoolAdded(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnWithFromMintTokenPoolRemotePoolAddedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPool.contract.FilterLogs(opts, "RemotePoolAdded", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolRemotePoolAddedIterator{contract: _BurnWithFromMintTokenPool.contract, event: "RemotePoolAdded", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) WatchRemotePoolAdded(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolRemotePoolAdded, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPool.contract.WatchLogs(opts, "RemotePoolAdded", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolRemotePoolAdded)
				if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "RemotePoolAdded", log); err != nil {
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

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) ParseRemotePoolAdded(log types.Log) (*BurnWithFromMintTokenPoolRemotePoolAdded, error) {
	event := new(BurnWithFromMintTokenPoolRemotePoolAdded)
	if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "RemotePoolAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolRemotePoolRemovedIterator struct {
	Event *BurnWithFromMintTokenPoolRemotePoolRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolRemotePoolRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolRemotePoolRemoved)
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
		it.Event = new(BurnWithFromMintTokenPoolRemotePoolRemoved)
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

func (it *BurnWithFromMintTokenPoolRemotePoolRemovedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolRemotePoolRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolRemotePoolRemoved struct {
	RemoteChainSelector uint64
	RemotePoolAddress   []byte
	Raw                 types.Log
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) FilterRemotePoolRemoved(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnWithFromMintTokenPoolRemotePoolRemovedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPool.contract.FilterLogs(opts, "RemotePoolRemoved", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolRemotePoolRemovedIterator{contract: _BurnWithFromMintTokenPool.contract, event: "RemotePoolRemoved", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) WatchRemotePoolRemoved(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolRemotePoolRemoved, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnWithFromMintTokenPool.contract.WatchLogs(opts, "RemotePoolRemoved", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolRemotePoolRemoved)
				if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "RemotePoolRemoved", log); err != nil {
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

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) ParseRemotePoolRemoved(log types.Log) (*BurnWithFromMintTokenPoolRemotePoolRemoved, error) {
	event := new(BurnWithFromMintTokenPoolRemotePoolRemoved)
	if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "RemotePoolRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolRouterUpdatedIterator struct {
	Event *BurnWithFromMintTokenPoolRouterUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolRouterUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolRouterUpdated)
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
		it.Event = new(BurnWithFromMintTokenPoolRouterUpdated)
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

func (it *BurnWithFromMintTokenPoolRouterUpdatedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolRouterUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolRouterUpdated struct {
	OldRouter common.Address
	NewRouter common.Address
	Raw       types.Log
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) FilterRouterUpdated(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolRouterUpdatedIterator, error) {

	logs, sub, err := _BurnWithFromMintTokenPool.contract.FilterLogs(opts, "RouterUpdated")
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolRouterUpdatedIterator{contract: _BurnWithFromMintTokenPool.contract, event: "RouterUpdated", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolRouterUpdated) (event.Subscription, error) {

	logs, sub, err := _BurnWithFromMintTokenPool.contract.WatchLogs(opts, "RouterUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolRouterUpdated)
				if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
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

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) ParseRouterUpdated(log types.Log) (*BurnWithFromMintTokenPoolRouterUpdated, error) {
	event := new(BurnWithFromMintTokenPoolRouterUpdated)
	if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnWithFromMintTokenPoolTokensConsumedIterator struct {
	Event *BurnWithFromMintTokenPoolTokensConsumed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnWithFromMintTokenPoolTokensConsumedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnWithFromMintTokenPoolTokensConsumed)
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
		it.Event = new(BurnWithFromMintTokenPoolTokensConsumed)
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

func (it *BurnWithFromMintTokenPoolTokensConsumedIterator) Error() error {
	return it.fail
}

func (it *BurnWithFromMintTokenPoolTokensConsumedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnWithFromMintTokenPoolTokensConsumed struct {
	Tokens *big.Int
	Raw    types.Log
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) FilterTokensConsumed(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolTokensConsumedIterator, error) {

	logs, sub, err := _BurnWithFromMintTokenPool.contract.FilterLogs(opts, "TokensConsumed")
	if err != nil {
		return nil, err
	}
	return &BurnWithFromMintTokenPoolTokensConsumedIterator{contract: _BurnWithFromMintTokenPool.contract, event: "TokensConsumed", logs: logs, sub: sub}, nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) WatchTokensConsumed(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolTokensConsumed) (event.Subscription, error) {

	logs, sub, err := _BurnWithFromMintTokenPool.contract.WatchLogs(opts, "TokensConsumed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnWithFromMintTokenPoolTokensConsumed)
				if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "TokensConsumed", log); err != nil {
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

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPoolFilterer) ParseTokensConsumed(log types.Log) (*BurnWithFromMintTokenPoolTokensConsumed, error) {
	event := new(BurnWithFromMintTokenPoolTokensConsumed)
	if err := _BurnWithFromMintTokenPool.contract.UnpackLog(event, "TokensConsumed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPool) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _BurnWithFromMintTokenPool.abi.Events["AllowListAdd"].ID:
		return _BurnWithFromMintTokenPool.ParseAllowListAdd(log)
	case _BurnWithFromMintTokenPool.abi.Events["AllowListRemove"].ID:
		return _BurnWithFromMintTokenPool.ParseAllowListRemove(log)
	case _BurnWithFromMintTokenPool.abi.Events["Burned"].ID:
		return _BurnWithFromMintTokenPool.ParseBurned(log)
	case _BurnWithFromMintTokenPool.abi.Events["ChainAdded"].ID:
		return _BurnWithFromMintTokenPool.ParseChainAdded(log)
	case _BurnWithFromMintTokenPool.abi.Events["ChainConfigured"].ID:
		return _BurnWithFromMintTokenPool.ParseChainConfigured(log)
	case _BurnWithFromMintTokenPool.abi.Events["ChainRemoved"].ID:
		return _BurnWithFromMintTokenPool.ParseChainRemoved(log)
	case _BurnWithFromMintTokenPool.abi.Events["ConfigChanged"].ID:
		return _BurnWithFromMintTokenPool.ParseConfigChanged(log)
	case _BurnWithFromMintTokenPool.abi.Events["Locked"].ID:
		return _BurnWithFromMintTokenPool.ParseLocked(log)
	case _BurnWithFromMintTokenPool.abi.Events["Minted"].ID:
		return _BurnWithFromMintTokenPool.ParseMinted(log)
	case _BurnWithFromMintTokenPool.abi.Events["OwnershipTransferRequested"].ID:
		return _BurnWithFromMintTokenPool.ParseOwnershipTransferRequested(log)
	case _BurnWithFromMintTokenPool.abi.Events["OwnershipTransferred"].ID:
		return _BurnWithFromMintTokenPool.ParseOwnershipTransferred(log)
	case _BurnWithFromMintTokenPool.abi.Events["RateLimitAdminSet"].ID:
		return _BurnWithFromMintTokenPool.ParseRateLimitAdminSet(log)
	case _BurnWithFromMintTokenPool.abi.Events["Released"].ID:
		return _BurnWithFromMintTokenPool.ParseReleased(log)
	case _BurnWithFromMintTokenPool.abi.Events["RemotePoolAdded"].ID:
		return _BurnWithFromMintTokenPool.ParseRemotePoolAdded(log)
	case _BurnWithFromMintTokenPool.abi.Events["RemotePoolRemoved"].ID:
		return _BurnWithFromMintTokenPool.ParseRemotePoolRemoved(log)
	case _BurnWithFromMintTokenPool.abi.Events["RouterUpdated"].ID:
		return _BurnWithFromMintTokenPool.ParseRouterUpdated(log)
	case _BurnWithFromMintTokenPool.abi.Events["TokensConsumed"].ID:
		return _BurnWithFromMintTokenPool.ParseTokensConsumed(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (BurnWithFromMintTokenPoolAllowListAdd) Topic() common.Hash {
	return common.HexToHash("0x2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d8")
}

func (BurnWithFromMintTokenPoolAllowListRemove) Topic() common.Hash {
	return common.HexToHash("0x800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf7566")
}

func (BurnWithFromMintTokenPoolBurned) Topic() common.Hash {
	return common.HexToHash("0x696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df7")
}

func (BurnWithFromMintTokenPoolChainAdded) Topic() common.Hash {
	return common.HexToHash("0x8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c2")
}

func (BurnWithFromMintTokenPoolChainConfigured) Topic() common.Hash {
	return common.HexToHash("0x0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b")
}

func (BurnWithFromMintTokenPoolChainRemoved) Topic() common.Hash {
	return common.HexToHash("0x5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d859916")
}

func (BurnWithFromMintTokenPoolConfigChanged) Topic() common.Hash {
	return common.HexToHash("0x9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c19")
}

func (BurnWithFromMintTokenPoolLocked) Topic() common.Hash {
	return common.HexToHash("0x9f1ec8c880f76798e7b793325d625e9b60e4082a553c98f42b6cda368dd60008")
}

func (BurnWithFromMintTokenPoolMinted) Topic() common.Hash {
	return common.HexToHash("0x9d228d69b5fdb8d273a2336f8fb8612d039631024ea9bf09c424a9503aa078f0")
}

func (BurnWithFromMintTokenPoolOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (BurnWithFromMintTokenPoolOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (BurnWithFromMintTokenPoolRateLimitAdminSet) Topic() common.Hash {
	return common.HexToHash("0x44676b5284b809a22248eba0da87391d79098be38bb03154be88a58bf4d09174")
}

func (BurnWithFromMintTokenPoolReleased) Topic() common.Hash {
	return common.HexToHash("0x2d87480f50083e2b2759522a8fdda59802650a8055e609a7772cf70c07748f52")
}

func (BurnWithFromMintTokenPoolRemotePoolAdded) Topic() common.Hash {
	return common.HexToHash("0x7d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea")
}

func (BurnWithFromMintTokenPoolRemotePoolRemoved) Topic() common.Hash {
	return common.HexToHash("0x52d00ee4d9bd51b40168f2afc5848837288ce258784ad914278791464b3f4d76")
}

func (BurnWithFromMintTokenPoolRouterUpdated) Topic() common.Hash {
	return common.HexToHash("0x02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f1684")
}

func (BurnWithFromMintTokenPoolTokensConsumed) Topic() common.Hash {
	return common.HexToHash("0x1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a")
}

func (_BurnWithFromMintTokenPool *BurnWithFromMintTokenPool) Address() common.Address {
	return _BurnWithFromMintTokenPool.address
}

type BurnWithFromMintTokenPoolInterface interface {
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

	FilterAllowListAdd(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolAllowListAddIterator, error)

	WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAllowListAdd) (event.Subscription, error)

	ParseAllowListAdd(log types.Log) (*BurnWithFromMintTokenPoolAllowListAdd, error)

	FilterAllowListRemove(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolAllowListRemoveIterator, error)

	WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolAllowListRemove) (event.Subscription, error)

	ParseAllowListRemove(log types.Log) (*BurnWithFromMintTokenPoolAllowListRemove, error)

	FilterBurned(opts *bind.FilterOpts, sender []common.Address) (*BurnWithFromMintTokenPoolBurnedIterator, error)

	WatchBurned(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolBurned, sender []common.Address) (event.Subscription, error)

	ParseBurned(log types.Log) (*BurnWithFromMintTokenPoolBurned, error)

	FilterChainAdded(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolChainAddedIterator, error)

	WatchChainAdded(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolChainAdded) (event.Subscription, error)

	ParseChainAdded(log types.Log) (*BurnWithFromMintTokenPoolChainAdded, error)

	FilterChainConfigured(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolChainConfiguredIterator, error)

	WatchChainConfigured(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolChainConfigured) (event.Subscription, error)

	ParseChainConfigured(log types.Log) (*BurnWithFromMintTokenPoolChainConfigured, error)

	FilterChainRemoved(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolChainRemovedIterator, error)

	WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolChainRemoved) (event.Subscription, error)

	ParseChainRemoved(log types.Log) (*BurnWithFromMintTokenPoolChainRemoved, error)

	FilterConfigChanged(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolConfigChangedIterator, error)

	WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolConfigChanged) (event.Subscription, error)

	ParseConfigChanged(log types.Log) (*BurnWithFromMintTokenPoolConfigChanged, error)

	FilterLocked(opts *bind.FilterOpts, sender []common.Address) (*BurnWithFromMintTokenPoolLockedIterator, error)

	WatchLocked(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolLocked, sender []common.Address) (event.Subscription, error)

	ParseLocked(log types.Log) (*BurnWithFromMintTokenPoolLocked, error)

	FilterMinted(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*BurnWithFromMintTokenPoolMintedIterator, error)

	WatchMinted(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolMinted, sender []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseMinted(log types.Log) (*BurnWithFromMintTokenPoolMinted, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnWithFromMintTokenPoolOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*BurnWithFromMintTokenPoolOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnWithFromMintTokenPoolOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*BurnWithFromMintTokenPoolOwnershipTransferred, error)

	FilterRateLimitAdminSet(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolRateLimitAdminSetIterator, error)

	WatchRateLimitAdminSet(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolRateLimitAdminSet) (event.Subscription, error)

	ParseRateLimitAdminSet(log types.Log) (*BurnWithFromMintTokenPoolRateLimitAdminSet, error)

	FilterReleased(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*BurnWithFromMintTokenPoolReleasedIterator, error)

	WatchReleased(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolReleased, sender []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseReleased(log types.Log) (*BurnWithFromMintTokenPoolReleased, error)

	FilterRemotePoolAdded(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnWithFromMintTokenPoolRemotePoolAddedIterator, error)

	WatchRemotePoolAdded(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolRemotePoolAdded, remoteChainSelector []uint64) (event.Subscription, error)

	ParseRemotePoolAdded(log types.Log) (*BurnWithFromMintTokenPoolRemotePoolAdded, error)

	FilterRemotePoolRemoved(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnWithFromMintTokenPoolRemotePoolRemovedIterator, error)

	WatchRemotePoolRemoved(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolRemotePoolRemoved, remoteChainSelector []uint64) (event.Subscription, error)

	ParseRemotePoolRemoved(log types.Log) (*BurnWithFromMintTokenPoolRemotePoolRemoved, error)

	FilterRouterUpdated(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolRouterUpdatedIterator, error)

	WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolRouterUpdated) (event.Subscription, error)

	ParseRouterUpdated(log types.Log) (*BurnWithFromMintTokenPoolRouterUpdated, error)

	FilterTokensConsumed(opts *bind.FilterOpts) (*BurnWithFromMintTokenPoolTokensConsumedIterator, error)

	WatchTokensConsumed(opts *bind.WatchOpts, sink chan<- *BurnWithFromMintTokenPoolTokensConsumed) (event.Subscription, error)

	ParseTokensConsumed(log types.Log) (*BurnWithFromMintTokenPoolTokensConsumed, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}

var BurnWithFromMintTokenPoolZKBin = ("0x0x0004000000000002001f000000000002000000000301001900000060033002700000066e0030019d0000066e033001970003000000310355000200000001035500000001002001900000005e0000c13d0000008004000039000000400040043f000000040030008c000000860000413d000000000201043b000000e0022002700000068e0020009c000000880000a13d0000068f0020009c000000c90000a13d000006900020009c000000f80000213d000006960020009c000001820000213d000006990020009c000005240000613d0000069a0020009c000000860000c13d0000000002000416000000000002004b000000860000c13d0000000504000039000000000204041a000000800020043f000000000040043f000000000002004b00000a120000c13d000000a002000039000000400020043f0000002004000039000000000500001900000005065002100000003f07600039000006d7077001970000000007270019000006720070009c000001240000213d000000400070043f00000000005204350000001f0560018f000000a004400039000000000006004b0000003c0000613d000000000131034f00000000036400190000000006040019000000001701043c0000000006760436000000000036004b000000380000c13d000000000005004b000000800100043d000000000001004b0000004e0000613d00000000010000190000000003020433000000000013004b000010320000a13d00000005031002100000000005430019000000a0033000390000000003030433000006720330019700000000003504350000000101100039000000800300043d000000000031004b000000410000413d000000400100043d00000020030000390000000005310436000000000302043300000000003504350000004002100039000000000003004b00000a090000613d00000000050000190000000046040434000006720660019700000000026204360000000105500039000000000035004b000000570000413d00000a090000013d0000010004000039000000400040043f0000000002000416000000000002004b000000860000c13d0000001f023000390000066f022001970000010002200039000000400020043f0000001f0530018f00000670063001980000010002600039000000700000613d000000000701034f000000007807043c0000000004840436000000000024004b0000006c0000c13d000000000005004b0000007d0000613d000000000161034f0000000304500210000000000502043300000000054501cf000000000545022f000000000101043b0000010004400089000000000141022f00000000014101cf000000000151019f0000000000120435000000a00030008c000000860000413d000001000100043d000006710010009c000000860000213d000001200200043d001100000002001d000000ff0020008c000001110000a13d0000000001000019000019b600010430000006a40020009c000000a30000213d000006ae0020009c0000012a0000a13d000006af0020009c000001590000213d000006b20020009c000001fb0000613d000006b30020009c000000860000c13d0000000001000416000000000001004b000000860000c13d0000000001000412001900000001001d001800200000003d000080050100003900000044030000390000000004000415000000190440008a0000000504400210000006ce0200004119b4198c0000040f000000ff0110018f000000800010043f000006cf01000041000019b50001042e000006a50020009c0000013b0000a13d000006a60020009c000001640000213d000006a90020009c000002160000613d000006aa0020009c000000860000c13d000000240030008c000000860000413d0000000002000416000000000002004b000000860000c13d0000000401100370000000000101043b000006710010009c000000860000213d0000000102000039000000000202041a00000671022001970000000003000411000000000023004b000009e20000c13d0000000902000039000000000302041a0000067503300197000000000313019f000000000032041b000000800010043f00000000010004140000066e0010009c0000066e01008041000000c001100210000006ec011001c70000800d020000390000000103000039000006ed0400004100000a360000013d0000069b0020009c000001460000a13d0000069c0020009c0000016d0000213d0000069f0020009c000002990000613d000006a00020009c000000860000c13d000000240030008c000000860000413d0000000002000416000000000002004b000000860000c13d0000000401100370000000000101043b001100000001001d000006720010009c000000860000213d19b415b00000040f0000001101000029000000000010043f0000000701000039000000200010043f0000004002000039000000000100001919b419770000040f001000000001001d000000400100043d001100000001001d19b414b30000040f00000010050000290000000201500039000000000401041a000006d4004001980000000002000039000000010200c03900000011010000290000004003100039000000000023043500000080024002700000066e0220019700000020031000390000000000230435000006c102400197000000000021043500000003025000390000056f0000013d000006910020009c000001e20000213d000006940020009c000005390000613d000006950020009c000000860000c13d0000000001000416000000000001004b000000860000c13d0000000001000412001300000001001d001200600000003d000080050100003900000044030000390000000004000415000000130440008a0000000504400210000006ce0200004119b4198c0000040f000000000001004b0000000001000039000000010100c039000000800010043f000006cf01000041000019b50001042e000001400200043d000006720020009c000000860000213d0000001f04200039000000000034004b000000000500001900000673050080410000067304400197000000000004004b00000000060000190000067306004041000006730040009c000000000605c019000000000006004b000000860000c13d00000100042000390000000004040433000006720040009c000009b60000a13d000006df01000041000000000010043f0000004101000039000000040010043f000006be01000041000019b600010430000006b40020009c000003fc0000613d000006b50020009c000003420000613d000006b60020009c000000860000c13d0000000001000416000000000001004b000000860000c13d0000000001000412001d00000001001d001c00000000003d0000800501000039000000440300003900000000040004150000001d0440008a000005430000013d000006ab0020009c0000040f0000613d000006ac0020009c000003560000613d000006ad0020009c000000860000c13d0000000001000416000000000001004b000000860000c13d00000009010000390000033d0000013d000006a10020009c000004d10000613d000006a20020009c000003940000613d000006a30020009c000000860000c13d0000000001000416000000000001004b000000860000c13d0000000202000039000000000102041a000000800010043f000000000020043f0000002002000039000000000001004b000009ea0000c13d000000a0010000390000000004020019000009f90000013d000006b00020009c000002320000613d000006b10020009c000000860000c13d0000000001000416000000000001004b000000860000c13d000000000103001919b414670000040f19b415200000040f0000028f0000013d000006a70020009c000002850000613d000006a80020009c000000860000c13d0000000001000416000000000001004b000000860000c13d00000001010000390000033d0000013d0000069d0020009c000003390000613d0000069e0020009c000000860000c13d000000240030008c000000860000413d0000000002000416000000000002004b000000860000c13d0000000401100370000000000101043b000006720010009c000000860000213d19b415c50000040f0000002002000039000000400300043d001100000003001d000000000223043619b414550000040f00000011020000290000057c0000013d000006970020009c0000054a0000613d000006980020009c000000860000c13d000000e40030008c000000860000413d0000000002000416000000000002004b000000860000c13d0000000402100370000000000202043b001100000002001d000006720020009c000000860000213d000000e002000039000000400020043f0000002402100370000000000202043b000000000002004b0000000003000039000000010300c039000000000032004b000000860000c13d000000800020043f0000004402100370000000000202043b000006c10020009c000000860000213d000000a00020043f0000006402100370000000000202043b000006c10020009c000000860000213d000000c00020043f0000014002000039000000400020043f0000008402100370000000000202043b000000000002004b0000000003000039000000010300c039000000000032004b000000860000c13d000000e00020043f000000a402100370000000000202043b000006c10020009c000000860000213d000001000020043f000000c401100370000000000101043b000006c10010009c000000860000213d000001200010043f0000000901000039000000000101041a00000671021001970000000001000411000000000021004b000001c30000613d0000000102000039000000000202041a0000067102200197000000000021004b00000df50000c13d0000001101000029000000000010043f0000000601000039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000000860000613d000000000101043b000000000101041a000000000001004b000003890000613d000000c00100043d000006c101100197000000800200043d000000000002004b00000ed90000c13d000000000001004b000001de0000c13d000000a00100043d000006c10010019800000edf0000613d000000400200043d001100000002001d000006c30100004100000f4c0000013d000006920020009c000005850000613d000006930020009c000000860000c13d000000240030008c000000860000413d0000000002000416000000000002004b000000860000c13d0000000401100370000000000601043b000006710060009c000000860000213d0000000101000039000000000101041a00000671011001970000000005000411000000000015004b000009e20000c13d000000000056004b00000a3a0000c13d000006b901000041000000800010043f000006ba01000041000019b600010430000000240030008c000000860000413d0000000002000416000000000002004b000000860000c13d0000000401100370000000000101043b001100000001001d000006710010009c000000860000213d0000000001000412001b00000001001d001a00000000003d0000800501000039000000440300003900000000040004150000001b0440008a0000000504400210000006ce0200004119b4198c0000040f0000067101100197000000110010006b00000000010000390000000101006039000000800010043f000006cf01000041000019b50001042e0000000001000416000000000001004b000000860000c13d000000000100041a00000671021001970000000006000411000000000026004b000009e60000c13d0000000102000039000000000302041a0000067504300197000000000464019f000000000042041b0000067501100197000000000010041b000000000100041400000671053001970000066e0010009c0000066e01008041000000c001100210000006b7011001c70000800d020000390000000303000039000006ef0400004119b419aa0000040f0000000100200190000000860000613d00000a490000013d000000240030008c000000860000413d0000000002000416000000000002004b000000860000c13d0000000402100370000000000202043b001100000002001d000006720020009c000000860000213d000000110230006a000006850020009c000000860000213d000001040020008c000000860000413d000000a002000039000000400020043f0000001102000029000f00840020003d0000000f01100360000000800000043f000000000101043b001000000001001d000006710010009c000000860000213d000006ce01000041000000000010044300000000010004120000000400100443000000240000044300000000010004140000066e0010009c0000066e01008041000000c001100210000006e0011001c7000080050200003919b419af0000040f0000000100200190000013e20000613d0000000202000367000000000101043b0000067101100197000000100010006b00000a6c0000c13d0000000f01000029000e0060001000920000000e01200360000000000101043b000006720010009c000000860000213d000000400300043d000006e20200004100000000002304350000008001100210000006e301100197000f00000003001d00000004023000390000000000120435000006ce010000410000000000100443000000000100041200000004001004430000004001000039000000240010044300000000010004140000066e0010009c0000066e01008041000000c001100210000006e0011001c7000080050200003919b419af0000040f0000000100200190000013e20000613d000000000201043b00000000010004140000067102200197000000040020008c00000be60000c13d0000000103000031000000200030008c0000002004000039000000000403401900000c110000013d000000240030008c000000860000413d0000000002000416000000000002004b000000860000c13d0000000401100370000000000101043b000006720010009c000000860000213d19b417c30000040f000000000001004b0000000001000039000000010100c039000000400200043d00000000001204350000066e0020009c0000066e020080410000004001200210000006eb011001c7000019b50001042e000000440030008c000000860000413d0000000002000416000000000002004b000000860000c13d0000000402100370000000000202043b001100000002001d000006720020009c000000860000213d0000002402100370000000000202043b000006720020009c000000860000213d0000002304200039000000000034004b000000860000813d000f00040020003d0000000f01100360000000000101043b001000000001001d000006720010009c000000860000213d0000002402200039000d00000002001d000e00100020002d0000000e0030006b000000860000213d0000000101000039000000000101041a00000671011001970000000002000411000000000012004b000009e20000c13d0000001101000029000000000010043f0000000601000039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000000860000613d000000000101043b000000000101041a000000000001004b000003890000613d0000001101000029000000000010043f0000000701000039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000000860000613d00000010020000290000001f022000390000070002200197000b00000002001d0000003f022000390000070002200197000000000101043b000c00000001001d000000400100043d0000000002210019000000000012004b00000000040000390000000104004039000006720020009c000001240000213d0000000100400190000001240000c13d000000400020043f000000100200002900000000022104360000000e05000029000000000050007c000000860000213d00000010040000290000070003400198000e001f00400193000a00000003001d00000000033200190000000f040000290000002004400039000f00000004001d0000000204400367000003000000613d000000000504034f0000000006020019000000005705043c0000000006760436000000000036004b000002fc0000c13d0000000e0000006b0000030e0000613d0000000a044003600000000e050000290000000305500210000000000603043300000000065601cf000000000656022f000000000404043b0000010005500089000000000454022f00000000045401cf000000000464019f0000000000430435000000100320002900000000000304350000066e0020009c0000066e02008041000000400220021000000000010104330000066e0010009c0000066e010080410000006001100210000000000121019f00000000020004140000066e0020009c0000066e02008041000000c002200210000000000112019f000006b7011001c7000080100200003919b419af0000040f0000000100200190000000860000613d0000000c020000290000000503200039000000000201043b000000000103001919b417d70000040f000000400700043d000000000001004b00000f560000c13d000c00000007001d000000240170003900000040020000390000000000210435000006dd01000041000000000017043500000004017000390000001102000029000000000021043500000044037000390000000d01000029000000100200002919b415910000040f0000000c020000290000104c0000013d0000000001000416000000000001004b000000860000c13d0000000401000039000000000101041a0000067101100197000000800010043f000006cf01000041000019b50001042e0000000001000416000000000001004b000000860000c13d000000c001000039000000400010043f0000001f01000039000000800010043f000006fa01000041000000a00010043f0000002001000039000000c00010043f0000008001000039000000e00200003919b414550000040f000000c00110008a0000066e0010009c0000066e010080410000006001100210000006fb011001c7000019b50001042e000000440030008c000000860000413d0000000002000416000000000002004b000000860000c13d0000000402100370000000000202043b001100000002001d000006720020009c000000860000213d0000002402100370000000000202043b000006720020009c000000860000213d0000002304200039000000000034004b000000860000813d0000000404200039000000000141034f000000000101043b001000000001001d000006720010009c000000860000213d0000002402200039000f00000002001d0000001001200029000000000031004b000000860000213d0000000101000039000000000101041a00000671011001970000000002000411000000000012004b000009e20000c13d0000001101000029000000000010043f0000000601000039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000000860000613d000000000101043b000000000101041a000000000001004b00000b5c0000c13d000000400100043d000006bd0200004100000000002104350000000402100039000000110300002900000000003204350000066e0010009c0000066e010080410000004001100210000006be011001c7000019b600010430000000240030008c000000860000413d0000000002000416000000000002004b000000860000c13d0000000401100370000000000101043b000006720010009c000000860000213d000000000010043f0000000701000039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000000860000613d000000000101043b0000000501100039000000000301041a000000400200043d000f00000002001d001100000003001d0000000002320436000e00000002001d000000000010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067d011001c7000080100200003919b419af0000040f0000000100200190000000860000613d0000001105000029000000000005004b0000000e02000029000003c80000613d000000000101043b0000000e020000290000000003000019000000000401041a000000000242043600000001011000390000000103300039000000000053004b000003c20000413d0000000f0120006a0000001f0110003900000700011001970000000f04100029000000000014004b00000000010000390000000101004039000006720040009c000001240000213d0000000100100190000001240000c13d000000400040043f0000000f010000290000000002010433000006720020009c000001240000213d00000005012002100000003f0310003900000674033001970000000003430019000006720030009c000001240000213d000000400030043f000d00000004001d0000000005240436000000000002004b000003ea0000613d00000060020000390000000003000019000000000435001900000000002404350000002003300039000000000013004b000003e50000413d000c00000005001d0000000f010000290000000001010433000000000001004b00000b650000c13d000000400100043d000000200200003900000000032104360000000d0200002900000000020204330000000000230435000000400310003900000005042002100000000005340019000000000002004b00000bc90000c13d000000000215004900000a0a0000013d000000240030008c000000860000413d0000000002000416000000000002004b000000860000c13d0000000401100370000000000201043b000006fc00200198000000860000c13d0000000101000039000006fd0020009c000005470000613d000006fe0020009c000005470000613d000006ff0020009c000000000100c019000000800010043f000006cf01000041000019b50001042e000000440030008c000000860000413d0000000002000416000000000002004b000000860000c13d0000000402100370000000000202043b000006720020009c000000860000213d0000002305200039000000000035004b000000860000813d0000000405200039000000000551034f000000000905043b000006720090009c000000860000213d0000002407200039000000050b90021000000000087b0019000000000038004b000000860000213d0000002402100370000000000202043b000006720020009c000000860000213d0000002305200039000000000035004b000000860000813d0000000405200039000000000551034f000000000605043b000006720060009c000000860000213d0000002402200039000000050a60021000000000052a0019000000000035004b000000860000213d0000000103000039000000000303041a0000067103300197000000000c00041100000000003c004b000009e20000c13d0000003f03b000390000067403300197000006d60030009c000001240000213d0000008003300039000f00000003001d000000400030043f000000800090043f000000000009004b000004510000613d000000000371034f000000000303043b000006710030009c000000860000213d000000200440003900000000003404350000002007700039000000000087004b000004460000413d000000400300043d000f00000003001d0000003f03a0003900000674033001970000000f033000290000000f0030006c00000000040000390000000104004039000006720030009c000001240000213d0000000100400190000001240000c13d000000400030043f0000000f030000290000000003630436000e00000003001d000000000006004b0000046b0000613d0000000f03000029000000000421034f000000000404043b000006710040009c000000860000213d000000200330003900000000004304350000002002200039000000000052004b000004620000413d000006ce010000410000000000100443000000000100041200000004001004430000006001000039000000240010044300000000010004140000066e0010009c0000066e01008041000000c001100210000006e0011001c7000080050200003919b419af0000040f0000000100200190000013e20000613d000000000101043b000000000001004b00000ae00000613d000000800100043d000000000001004b00000fcd0000c13d0000000f010000290000000001010433000000000001004b00000a490000613d00000000030000190000048b0000013d00000001033000390000000f010000290000000001010433000000000013004b00000a490000813d00000005013002100000000e0110002900000000010104330000067104100198000004860000613d000000000040043f0000000301000039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c70000801002000039001000000003001d001100000004001d19b419af0000040f000000110400002900000010030000290000000100200190000000860000613d000000000101043b000000000101041a000000000001004b000004860000c13d0000000203000039000000000103041a000006720010009c000001240000213d0000000102100039000000000023041b0000067c0110009a000000000041041b000000000103041a000d00000001001d000000000040043f0000000301000039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f00000011030000290000000100200190000000860000613d000000000101043b0000000d02000029000000000021041b000000400100043d00000000003104350000066e0010009c0000066e01008041000000400110021000000000020004140000066e0020009c0000066e02008041000000c002200210000000000112019f0000067d011001c70000800d0200003900000001030000390000067e0400004119b419aa0000040f00000010030000290000000100200190000004860000c13d000000860000013d000000240030008c000000860000413d0000000002000416000000000002004b000000860000c13d0000000402100370000000000202043b000006720020009c000000860000213d0000000003230049000006850030009c000000860000213d000000a40030008c000000860000413d000000c003000039000000400030043f0000006003000039000000800030043f000000a00030043f001000840020003d0000001001100360000000000101043b001100000001001d000006710010009c000000860000213d000006ce01000041000000000010044300000000010004120000000400100443000000240000044300000000010004140000066e0010009c0000066e01008041000000c001100210000006e0011001c7000080050200003919b419af0000040f0000000100200190000013e20000613d0000000202000367000000000101043b0000067101100197000000110010006b00000a4b0000c13d0000001001000029000f0060001000920000000f01200360000000000101043b000006720010009c000000860000213d000000400300043d000006e20200004100000000002304350000008001100210000006e301100197001000000003001d00000004023000390000000000120435000006ce010000410000000000100443000000000100041200000004001004430000004001000039000000240010044300000000010004140000066e0010009c0000066e01008041000000c001100210000006e0011001c7000080050200003919b419af0000040f0000000100200190000013e20000613d000000000201043b00000000010004140000067102200197000000040020008c00000ae30000c13d0000000103000031000000200030008c0000002004000039000000000403401900000b0e0000013d000000240030008c000000860000413d0000000002000416000000000002004b000000860000c13d0000000401100370000000000101043b000006710010009c000000860000213d0000000102000039000000000202041a00000671022001970000000003000411000000000023004b000009e20000c13d000000000001004b00000a260000c13d000006cd01000041000000800010043f000006ba01000041000019b6000104300000000001000416000000000001004b000000860000c13d0000000001000412001500000001001d001400400000003d000080050100003900000044030000390000000004000415000000150440008a0000000504400210000006ce0200004119b4198c0000040f0000067101100197000000800010043f000006cf01000041000019b50001042e000000240030008c000000860000413d0000000002000416000000000002004b000000860000c13d0000000401100370000000000101043b001100000001001d000006720010009c000000860000213d19b415b00000040f0000001101000029000000000010043f0000000701000039000000200010043f0000004002000039000000000100001919b419770000040f001000000001001d000000400100043d001100000001001d19b414b30000040f0000001005000029000000000405041a000006d4004001980000000002000039000000010200c03900000011010000290000004003100039000000000023043500000080024002700000066e0220019700000020031000390000000000230435000006c10240019700000000002104350000000102500039000000000402041a000000800210003900000080034002700000000000320435000006c1034001970000006002100039000000000032043519b417750000040f000000400100043d001000000001001d000000110200002919b414840000040f000000100200002900000000012100490000066e0010009c0000066e0100804100000060011002100000066e0020009c0000066e020080410000004002200210000000000121019f000019b50001042e000000440030008c000000860000413d0000000002000416000000000002004b000000860000c13d0000000402100370000000000202043b000006720020009c000000860000213d0000002304200039000000000034004b000000860000813d0000000404200039000000000441034f000000000404043b000600000004001d000006720040009c000000860000213d000500240020003d000000060200002900000005022002100000000502200029000000000032004b000000860000213d0000002402100370000000000202043b000200000002001d000006720020009c000000860000213d00000002020000290000002302200039000000000032004b000000860000813d00000002020000290000000402200039000000000121034f000000000101043b000100000001001d000006720010009c000000860000213d0000000201000029000300240010003d000000010100002900000005011002100000000301100029000000000031004b000000860000213d0000000101000039000000000101041a00000671011001970000000002000411000000000012004b000009e20000c13d000000060000006b00000c570000c13d000000010000006b00000a490000613d000500000000001d0000000501000029000000050110021000000003011000290000000202000367000000000112034f000000000101043b0000000003000031000000020430006a000001430440008a00000673054001970000067306100197000000000756013f000000000056004b00000000050000190000067305004041000000000041004b00000000040000190000067304008041000006730070009c000000000504c019000000000005004b000000860000c13d001000030010002d000000100130006a000f00000001001d000006850010009c000000860000213d0000000f01000029000001200010008c000000860000413d000000400100043d000900000001001d000006bf0010009c000001240000213d0000000901000029000000a001100039000000400010043f0000001001200360000000000101043b000006720010009c000000860000213d00000009040000290000000001140436000800000001001d00000010010000290000002001100039000000000112034f000000000101043b000006720010009c000000860000213d0000001001100029001100000001001d0000001f01100039000000000031004b0000000004000019000006730400804100000673011001970000067305300197000000000751013f000000000051004b00000000010000190000067301004041000006730070009c000000000104c019000000000001004b000000860000c13d0000001101200360000000000101043b000006720010009c000001240000213d00000005091002100000003f049000390000067404400197000000400600043d0000000004460019000e00000006001d000000000064004b00000000070000390000000107004039000006720040009c000001240000213d0000000100700190000001240000c13d000000400040043f0000000e040000290000000000140435000000110100002900000020081000390000000009980019000000000039004b000000860000213d000000000098004b000006680000813d0000000e0a000029000006250000013d000000200aa000390000000001b7001900000000000104350000000000ca04350000002008800039000000000098004b000006680000813d000000000182034f000000000101043b000006720010009c000000860000213d000000110d1000290000003f01d00039000000000031004b000000000400001900000673040080410000067301100197000000000751013f000000000051004b00000000010000190000067301004041000006730070009c000000000104c019000000000001004b000000860000c13d000000200ed000390000000001e2034f000000000b01043b0000067200b0009c000001240000213d0000001f01b0003900000700011001970000003f011000390000070001100197000000400c00043d00000000011c00190000000000c1004b00000000040000390000000104004039000006720010009c000001240000213d0000000100400190000001240000c13d0000004004d00039000000400010043f0000000007bc043600000000014b0019000000000031004b000000860000213d0000002001e00039000000000412034f0000070001b00198000000000e1700190000065a0000613d000000000f04034f000000000d07001900000000f60f043c000000000d6d04360000000000ed004b000006560000c13d0000001f0db001900000061e0000613d000000000114034f0000000304d0021000000000060e043300000000064601cf000000000646022f000000000101043b0000010004400089000000000141022f00000000014101cf000000000161019f00000000001e04350000061e0000013d00000008010000290000000e04000029000000000041043500000010010000290000004006100039000000000162034f000000000101043b000006720010009c000000860000213d00000010071000290000001f01700039000000000031004b000000000400001900000673040080410000067301100197000000000851013f000000000051004b00000000010000190000067301004041000006730080009c000000000104c019000000000001004b000000860000c13d000000000172034f000000000401043b000006720040009c000001240000213d0000001f0140003900000700011001970000003f011000390000070001100197000000400500043d0000000001150019000000000051004b00000000080000390000000108004039000006720010009c000001240000213d0000000100800190000001240000c13d0000002008700039000000400010043f00000000074504360000000001840019000000000031004b000000860000213d000000000882034f00000700014001980000000003170019000006a00000613d000000000908034f000000000a070019000000009b09043c000000000aba043600000000003a004b0000069c0000c13d0000001f09400190000006ad0000613d000000000118034f0000000308900210000000000903043300000000098901cf000000000989022f000000000101043b0000010008800089000000000181022f00000000018101cf000000000191019f00000000001304350000000001470019000000000001043500000009010000290000004001100039000600000001001d00000000005104350000000f01000029000000600110008a000006850010009c000000860000213d000000600010008c000000860000413d000000400300043d000006c00030009c000001240000213d0000006001300039000000400010043f0000002001600039000000000412034f000000000404043b000000000004004b0000000005000039000000010500c039000000000054004b000000860000c13d00000000044304360000002001100039000000000512034f000000000505043b000006c10050009c000000860000213d00000000005404350000002004100039000000000142034f000000000101043b000006c10010009c000000860000213d0000004005300039000000000015043500000009010000290000006001100039000700000001001d00000000003104350000000f01000029000000c00110008a000006850010009c000000860000213d000000600010008c000000860000413d000000400100043d000006c00010009c000001240000213d0000006003100039000000400030043f0000002003400039000000000432034f000000000404043b000000000004004b0000000005000039000000010500c039000000000054004b000000860000c13d00000000044104360000002003300039000000000532034f000000000505043b000006c10050009c000000860000213d00000000005404350000002003300039000000000232034f000000000302043b000006c10030009c000000860000213d0000004002100039000000000032043500000009030000290000008003300039000400000003001d00000000001304350000000703000029000000000303043300000040053000390000000005050433000006c1065001970000000057030434000000000007004b0000070c0000613d000000000006004b000012780000613d0000000005050433000006c105500197000000000056004b000007110000413d000012780000013d000000000006004b000012640000c13d0000000005050433000006c100500198000012640000c13d0000000002020433000006c1022001970000000003010433000000000003004b0000071d0000613d000000000002004b0000127f0000613d0000000003040433000006c103300197000000000032004b000007220000413d0000127f0000013d000000000002004b000012680000c13d0000000002040433000006c100200198000012680000c13d000000060100002900000000010104330000000001010433000000000001004b00000a7b0000613d000000090100002900000000010104330000067201100197001100000001001d000000000010043f0000000601000039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000000860000613d000000000101043b000000000101041a000000000001004b000012490000c13d0000000501000039000000000101041a000006720010009c000001240000213d00000001021000390000000503000039000000000023041b000006c50110009a0000001102000029000000000021041b000000000103041a001000000001001d000000000020043f0000000601000039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000000860000613d000000000101043b0000001002000029000000000021041b000000090100002900000000010104330000067201100197000000000010043f0000000701000039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000000860000613d000000000301043b000000070100002900000000010104330000000024010434001100000004001d000000000004004b0000000004000039000000010400c039001000000004001d000000400400043d000006bf0040009c000001240000213d000f00000003001d0000000002020433000006c10220019700000040011000390000000001010433000b00000001001d000000a001400039000000400010043f000e00000002001d000d00000004001d0000000001240436000c00000001001d000006c601000041000000000010044300000000010004140000066e0010009c0000066e01008041000000c001100210000006c7011001c70000800b0200003919b419af0000040f0000000100200190000013e20000613d0000000b02000029000006c102200197000000000101043b0000000d040000290000004003400039000000100500002900000000005304350000008003400039000000000023043500000060034000390000000e0500002900000000005304350000066e011001970000000c030000290000000000130435000000110000006b0000000003000019000006c80300c0410000000f09000029000000000409041a000006c904400197000000000343019f0000008002200210000000000425019f000000000353019f0000008002100210000000000323019f000000000039041b0000000103900039000000000043041b000000400300043d000006bf0030009c000001240000213d0000000404000029000000000404043300000020054000390000000005050433000000000604043300000040044000390000000004040433000000a007300039000000400070043f0000008007300039000006c10840019700000000008704350000004007300039000000000006004b0000000006000039000000010600c039000000000067043500000020063000390000000000160435000006c1015001970000006005300039000000000015043500000000001304350000000003000019000006c80300c0410000000205900039000000000605041a000006c906600197000000000363019f000000000223019f000000000212019f000000000025041b0000008002400210000000000112019f0000000302900039000000000012041b000000060100002900000000030104330000000054030434000006720040009c000001240000213d0000000406900039000000000106041a000000010010019000000001071002700000007f0770618f0000001f0070008c00000000020000390000000102002039000000000121013f000000010010019000000dfa0000c13d000000200070008c001100000006001d001000000004001d000e00000003001d000008010000413d000d00000007001d000f00000005001d000000000060043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067d011001c7000080100200003919b419af0000040f0000000100200190000000860000613d00000010040000290000001f024000390000000502200270000000200040008c0000000002004019000000000301043b0000000d010000290000001f01100039000000050110027000000000011300190000000002230019000000000012004b00000011060000290000000f05000029000008010000813d000000000002041b0000000102200039000000000012004b000007fd0000413d0000001f0040008c000008200000a13d000000000060043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067d011001c7000080100200003919b419af0000040f0000000100200190000000860000613d00000010070000290000070002700198000000000101043b0000000e080000290000082c0000613d000000010320008a0000000503300270000000000331001900000001043000390000002003000039000000110600002900000000058300190000000005050433000000000051041b00000020033000390000000101100039000000000041004b000008180000c13d0000082e0000013d000000000004004b000008240000613d0000000001050433000008250000013d00000000010000190000000302400210000007010220027f0000070102200167000000000121016f0000000102400210000000000121019f0000083a0000013d00000020030000390000001106000029000000000072004b000008380000813d0000000302700210000000f80220018f000007010220027f000007010220016700000000038300190000000003030433000000000223016f000000000021041b000000010170021000000001011001bf000000000016041b000000080100002900000000010104330000000002010433000000000002004b0000095d0000613d0000000003000019000c00000003001d0000000502300210000000000121001900000020011000390000000001010433001000000001001d0000000031010434000000000001004b00000a7b0000613d00000009020000290000000002020433001100000002001d0000066e0010009c0000066e0100804100000060011002100000066e0030009c000f00000003001d0000066e0200004100000000020340190000004002200210000000000121019f00000000020004140000066e0020009c0000066e02008041000000c002200210000000000112019f000006b7011001c7000080100200003919b419af0000040f0000000100200190000000860000613d00000011020000290000067202200197000000000101043b001100000001001d000b00000002001d000000000020043f0000000701000039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000000860000613d000000000201043b0000001101000029000000000010043f000e00000002001d0000000601200039000d00000001001d000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000000860000613d000000000101043b000000000101041a000000000001004b0000103e0000c13d0000000e010000290000000502100039000000000102041a000006720010009c000001240000213d000a00000001001d0000000101100039000000000012041b000e00000002001d000000000020043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067d011001c7000080100200003919b419af0000040f00000001002001900000001102000029000000860000613d000000000101043b0000000a01100029000000000021041b0000000e01000029000000000101041a000e00000001001d000000000020043f0000000d01000029000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f00000001002001900000001102000029000000860000613d000000000101043b0000000e03000029000000000031041b000000000020043f0000000801000039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000000860000613d000000000401043b00000010010000290000000005010433000006720050009c000001240000213d000000000104041a000000010010019000000001031002700000007f0330618f0000001f0030008c00000000020000390000000102002039000000000121013f00000001001001900000000f0700002900000dfa0000c13d000000200030008c001100000004001d000e00000005001d000008ed0000413d000d00000003001d000000000040043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067d011001c7000080100200003919b419af0000040f0000000100200190000000860000613d0000000e050000290000001f025000390000000502200270000000200050008c0000000002004019000000000301043b0000000d010000290000001f01100039000000050110027000000000011300190000000002230019000000000012004b0000000f070000290000001104000029000008ed0000813d000000000002041b0000000102200039000000000012004b000008e90000413d0000001f0050008c000009190000a13d000000000040043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067d011001c7000080100200003919b419af0000040f0000000100200190000000860000613d0000000e080000290000070002800198000000000101043b000009570000613d000000010320008a00000005033002700000000003310019000000010430003900000020030000390000000f07000029000000100600002900000000056300190000000005050433000000000051041b00000020033000390000000101100039000000000041004b000009040000c13d000000000082004b000009150000813d0000000302800210000000f80220018f000007010220027f000007010220016700000000036300190000000003030433000000000223016f000000000021041b000000010180021000000001011001bf0000001104000029000009250000013d000000000005004b0000091d0000613d00000000010704330000091e0000013d000000000100001900000010060000290000000302500210000007010220027f0000070102200167000000000121016f0000000102500210000000000121019f000000000014041b000000400100043d00000020020000390000000003210436000000000206043300000000002304350000004003100039000000000002004b000009360000613d000000000400001900000000053400190000000006740019000000000606043300000000006504350000002004400039000000000024004b0000092f0000413d0000001f0420003900000700044001970000000002320019000000000002043500000040024000390000066e0020009c0000066e0200804100000060022002100000066e0010009c0000066e010080410000004001100210000000000112019f00000000020004140000066e0020009c0000066e02008041000000c002200210000000000112019f000006b7011001c70000800d020000390000000203000039000006cb040000410000000b0500002919b419aa0000040f0000000100200190000000860000613d0000000c030000290000000103300039000000080100002900000000010104330000000002010433000000000023004b000008410000413d0000095d0000013d00000020030000390000000f070000290000001006000029000000000082004b0000090d0000413d000009150000013d00000004010000290000000002010433000000070100002900000000050104330000000601000029000000000301043300000009010000290000000004010433000000400100043d000000200610003900000100070000390000000000760435000006720440019700000000004104350000010007100039000000006403043400000000004704350000012003100039000000000004004b000009790000613d000000000700001900000000083700190000000009760019000000000909043300000000009804350000002007700039000000000047004b000009720000413d000000000643001900000000000604350000000076050434000000000006004b0000000006000039000000010600c039000000400810003900000000006804350000000006070433000006c1066001970000006007100039000000000067043500000040055000390000000005050433000006c105500197000000800610003900000000005604350000000065020434000000000005004b0000000005000039000000010500c039000000a00710003900000000005704350000000005060433000006c105500197000000c006100039000000000056043500000040022000390000000002020433000006c102200197000000e00510003900000000002504350000001f024000390000070002200197000000000212004900000000023200190000066e0020009c0000066e0200804100000060022002100000066e0010009c0000066e010080410000004001100210000000000112019f00000000020004140000066e0020009c0000066e02008041000000c002200210000000000112019f000006b7011001c70000800d020000390000000103000039000006cc0400004119b419aa0000040f0000000100200190000000860000613d00000005020000290000000102200039000500000002001d000000010020006c000005bf0000413d00000a490000013d00000005054002100000003f065000390000067406600197000000400700043d0000000006670019000f00000007001d000000000076004b00000000070000390000000107004039000006720060009c000001240000213d0000000100700190000001240000c13d0000010007300039000000400060043f0000000f030000290000000003430436000e00000003001d00000120022000390000000003250019000000000073004b000000860000213d000000000004004b000009d50000613d0000000e040000290000000025020434000006710050009c000000860000213d0000000004540436000000000032004b000009cf0000413d000001600500043d000006710050009c000000860000213d000001800200043d001000000002001d000006710020009c000000860000213d0000000003000411000000000003004b00000a4d0000c13d000000400100043d0000068d0200004100000a7d0000013d000006f001000041000000800010043f000006ba01000041000019b600010430000006ee01000041000000800010043f000006ba01000041000019b600010430000000a005000039000006de0300004100000000040000190000000006050019000000000503041a000000000556043600000001033000390000000104400039000000000014004b000009ed0000413d000000410160008a0000070004100197000006d60040009c000001240000213d0000008001400039000000400010043f0000000000210435000000a002400039000000800300043d0000000000320435000000c002400039000000000003004b00000a090000613d000000a00400003900000000050000190000000046040434000006710660019700000000026204360000000105500039000000000035004b00000a030000413d00000000021200490000066e0020009c0000066e0200804100000060022002100000066e0010009c0000066e010080410000004001100210000000000112019f000019b50001042e000000a006000039000006d50400004100000000050000190000000007060019000000000604041a000000000667043600000001044000390000000105500039000000000025004b00000a150000413d000000410270008a0000070004200197000006d60040009c000001240000213d0000008002400039000000800500043d000000400020043f000006720050009c000000290000a13d000001240000013d0000000402000039000000000302041a0000067504300197000000000414019f000000000042041b0000067102300197000000800020043f000000a00010043f00000000010004140000066e0010009c0000066e01008041000000c001100210000006d8011001c70000800d020000390000000103000039000006d90400004119b419aa0000040f0000000100200190000000860000613d00000a490000013d000000000100041a0000067501100197000000000161019f000000000010041b00000000010004140000066e0010009c0000066e01008041000000c001100210000006b7011001c70000800d020000390000000303000039000006b80400004119b419aa0000040f0000000100200190000000860000613d0000000001000019000019b50001042e000000100100002900000a6d0000013d0000000106000039000000000406041a0000067504400197000000000334019f000000000036041b000000000005004b00000a7b0000613d000006710210019800000a7b0000613d000000100000006b00000a7b0000613d000b00000006001d000000800020043f000000c00050043f0000067601000041000000400300043d000d00000003001d00000000001304350000000001000414000000040020008c000c00000002001d00000a830000c13d00000000010004150000001f0110008a00000005011002100000000103000031000000200030008c00000020040000390000000004034019001f00000000003d00000ab20000013d0000000f01000029000000000112034f000000000101043b000006710010009c000000860000213d000000400200043d000006e1030000410000000000320435000000040320003900000000001304350000066e0020009c0000066e020080410000004001200210000006be011001c7000019b600010430000000400100043d000006cd0200004100000000002104350000066e0010009c0000066e01008041000000400110021000000677011001c7000019b6000104300000000d030000290000066e0030009c0000066e0300804100000040033002100000066e0010009c0000066e01008041000000c001100210000000000131019f00000677011001c719b419af0000040f000000000301001900000060033002700000066e03300197000000200030008c000000200400003900000000040340190000001f0640018f00000020074001900000000d0570002900000a9d0000613d000000000801034f0000000d09000029000000008a08043c0000000009a90436000000000059004b00000a990000c13d000000000006004b00000aaa0000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f000300000001035500000000010004150000001e0110008a0000000501100210001e00000000003d000000010020019000000ac90000613d0000001f02400039000000600420018f0000000d02400029000000000042004b00000000040000390000000104004039000006720020009c000001240000213d0000000100400190000001240000c13d000000400020043f000000200030008c000000860000413d0000000d030000290000000003030433000000ff0030008c000000860000213d0000000501100270000000000103001f0000001101000029000000ff0110018f000000000031004b00000e020000c13d0000001101000029000000a00010043f0000000402000039000000000102041a000006750110019700000010011001af000000000012041b0000000f010000290000000001010433000000000001004b0000000001000039000000010100c039000000e00010043f00000e110000613d000000400100043d0000067a0010009c000001240000213d0000002002100039000000400020043f0000000000010435000000e00100043d000000000001004b00000e0d0000c13d000000400100043d000006f30200004100000a7d0000013d00000010030000290000066e0030009c0000066e0300804100000040033002100000066e0010009c0000066e01008041000000c001100210000000000131019f000006be011001c719b419af0000040f000000000301001900000060033002700000066e03300197000000200030008c000000200400003900000000040340190000001f0640018f0000002007400190000000100570002900000afd0000613d000000000801034f0000001009000029000000008a08043c0000000009a90436000000000059004b00000af90000c13d000000000006004b00000b0a0000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f0003000000010355000000010020019000000c4b0000613d0000001f01400039000000600210018f0000001001200029000000000021004b00000000020000390000000102004039000006720010009c000001240000213d0000000100200190000001240000c13d000000400010043f000000200030008c000000860000413d00000010020000290000000002020433000000000002004b0000000003000039000000010300c039000000000032004b000000860000c13d000000000002004b00000e000000c13d0000000f0100002900000020011000390000000201100367000000000101043b001000000001001d000006710010009c000000860000213d000006ce010000410000000000100443000000000100041200000004001004430000006001000039000000240010044300000000010004140000066e0010009c0000066e01008041000000c001100210000006e0011001c7000080050200003919b419af0000040f0000000100200190000013e20000613d000000000101043b000000000001004b000010550000c13d0000000f010000290000000201100367000000000101043b001000000001001d000006720010009c000000860000213d0000001001000029000000000010043f0000000601000039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000000860000613d000000400200043d000e00000002001d0000000402200039000000000101043b000000000101041a000000000001004b000011270000c13d000006ea010000410000000e030000290000000000130435000000100100002900000c450000013d00000000030000310000000f01000029000000100200002919b414e80000040f0000000002010019000000110100002919b416370000040f0000000001000019000019b50001042e0000000002000019001100000002001d0000000502200210001000000002001d0000000e012000290000000001010433000000000010043f0000000801000039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000000860000613d000000000101043b000000000201041a000000010320019000000001052002700000007f0550618f0000001f0050008c00000000040000390000000104002039000000000043004b00000dfa0000c13d000000400700043d0000000004570436000000000003004b00000ba30000613d000900000004001d000a00000005001d000b00000007001d000000000010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067d011001c7000080100200003919b419af0000040f0000000100200190000000860000613d0000000a08000029000000000008004b000000200500008a00000bab0000613d000000000201043b00000000010000190000000d060000290000000b0700002900000009090000290000000003190019000000000402041a000000000043043500000001022000390000002001100039000000000081004b00000b9b0000413d00000bae0000013d00000702012001970000000000140435000000000005004b00000020010000390000000001006039000000200500008a0000000d0600002900000bae0000013d00000000010000190000000d060000290000000b070000290000003f01100039000000000251016f0000000001720019000000000021004b00000000020000390000000102004039000006720010009c000001240000213d0000000100200190000001240000c13d000000400010043f00000000010604330000001102000029000000000021004b000010320000a13d00000010030000290000000c0130002900000000007104350000000001060433000000000021004b000010320000a13d00000001022000390000000f010000290000000001010433000000000012004b00000b660000413d000003ef0000013d00000000040000190000000d0c00002900000bd40000013d0000001f0760003900000700077001970000000006650019000000000006043500000000057500190000000104400039000000000024004b000003fa0000813d0000000006150049000000400660008a0000000003630436000000200cc0003900000000060c043300000000760604340000000005650436000000000006004b00000bcc0000613d00000000080000190000000009580019000000000a870019000000000a0a04330000000000a904350000002008800039000000000068004b00000bde0000413d00000bcc0000013d0000000f030000290000066e0030009c0000066e0300804100000040033002100000066e0010009c0000066e01008041000000c001100210000000000131019f000006be011001c719b419af0000040f000000000301001900000060033002700000066e03300197000000200030008c000000200400003900000000040340190000001f0640018f00000020074001900000000f0570002900000c000000613d000000000801034f0000000f09000029000000008a08043c0000000009a90436000000000059004b00000bfc0000c13d000000000006004b00000c0d0000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f0003000000010355000000010020019000000de90000613d0000001f01400039000000600210018f0000000f01200029000000000021004b00000000020000390000000102004039000006720010009c000001240000213d0000000100200190000001240000c13d000000400010043f000000200030008c000000860000413d0000000f020000290000000002020433000000000002004b0000000003000039000000010300c039000000000032004b000000860000c13d000000000002004b00000e000000c13d0000000e010000290000000201100367000000000101043b000f00000001001d000006720010009c000000860000213d0000000f01000029000000000010043f0000000601000039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000000860000613d000000400200043d000d00000002001d0000000402200039000000000101043b000000000101041a000000000001004b000010bf0000c13d000006ea010000410000000d0300002900000000001304350000000f0100002900000000001204350000066e0030009c0000066e030080410000004001300210000006be011001c7000019b6000104300000001f0530018f0000067006300198000000400200043d000000000462001900000ec60000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b00000c520000c13d00000ec60000013d0000000002000019000700000002001d000000050120021000000005011000290000000201100367000000000101043b000b00000001001d000006720010009c000000860000213d0000000b01000029000000000010043f0000000601000039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000000860000613d000000000101043b000000000301041a000000000003004b000010db0000613d0000000501000039000000000201041a000000000002004b000013ef0000613d000000010130008a000000000023004b00000c930000613d000000000012004b000010320000a13d000006bb0130009a000006bb0220009a000000000202041a000000000021041b000000000020043f0000000601000039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c70000801002000039001100000003001d19b419af0000040f0000000100200190000000860000613d000000000101043b0000001102000029000000000021041b0000000501000039000000000301041a000000000003004b000010380000613d000000010130008a000006bb0230009a000000000002041b0000000502000039000000000012041b0000000b01000029000000000010043f0000000601000039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000000860000613d000000000101043b000000000001041b0000000b01000029000000000010043f0000000701000039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000000860000613d000000000101043b0000000501100039000000000301041a000000400200043d000f00000002001d001100000003001d0000000002320436000a00000002001d000000000010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067d011001c7000080100200003919b419af0000040f0000000100200190000000860000613d0000001105000029000000000005004b0000000a0200002900000cd20000613d000000000101043b0000000a020000290000000003000019000000000401041a000000000242043600000001011000390000000103300039000000000053004b00000ccc0000413d0000000f0120006a0000001f0110003900000700021001970000000f01200029000000000021004b00000000020000390000000102004039000006720010009c000001240000213d0000000100200190000001240000c13d000000400010043f0000000f010000290000000001010433000000000001004b00000d780000613d000000000200001900000cec0000013d000000000101043b000000000001041b000000110200002900000001022000390000000f010000290000000001010433000000000012004b00000d780000813d001100000002001d0000000b01000029000000000010043f0000000701000039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000000860000613d000000000301043b0000000f0100002900000000010104330000001102000029000000000021004b000010320000a13d00000005012002100000000a011000290000000001010433000c00000001001d000000000010043f000d00000003001d0000000601300039000e00000001001d000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000000860000613d000000000101043b000000000101041a000000000001004b00000ce60000613d0000000d020000290000000503200039000000000203041a000000000002004b000013ef0000613d000000000021004b001000000001001d000d00000003001d00000d590000613d000900000002001d000000000030043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067d011001c7000080100200003919b419af0000040f0000000100200190000000860000613d00000010020000290008000100200092000000000101043b0000000d04000029000000000204041a000000080020006c0000000003040019000010320000a13d0000000902000029000000010220008a0000000001120019000000000101041a000900000001001d000000000030043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067d011001c7000080100200003919b419af0000040f0000000100200190000000860000613d000000000101043b00000008011000290000000902000029000000000021041b000000000020043f0000000e01000029000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000000860000613d000000000101043b0000001002000029000000000021041b0000000d03000029000000000103041a001000000001001d000000000001004b000010380000613d000000000030043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067d011001c7000080100200003919b419af0000040f0000000100200190000000860000613d0000001002000029000000010220008a000000000101043b0000000001210019000000000001041b0000000d01000029000000000021041b0000000c01000029000000000010043f0000000e01000029000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f000000010020019000000ce40000c13d000000860000013d0000000b01000029000000000010043f0000000701000039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000000860000613d000000000301043b000000000003041b0000000101300039000000000001041b0000000201300039000000000001041b0000000301300039000000000001041b0000000404300039000000000104041a000000010010019000000001051002700000007f0550618f0000001f0050008c00000000020000390000000102002039000000000121013f000000010010019000000dfa0000c13d000000000005004b00000dba0000613d0000001f0050008c00000db90000a13d000f00000005001d001100000003001d001000000004001d000000000040043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067d011001c7000080100200003919b419af0000040f0000000100200190000000860000613d000000000101043b0000000f020000290000001f02200039000000050220027000000000022100190000000103100039000000000023004b00000db50000813d000000000003041b0000000103300039000000000023004b00000db10000413d0000001002000029000000000002041b00000000040100190000001103000029000000000004041b0000000501300039000000000201041a000000000001041b000000000002004b00000dd20000613d001100000002001d000000000010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067d011001c7000080100200003919b419af0000040f0000000100200190000000860000613d000000000101043b0000001102100029000000000021004b00000dd20000813d000000000001041b0000000101100039000000000021004b00000dce0000413d000000400100043d0000000b0200002900000000002104350000066e0010009c0000066e01008041000000400110021000000000020004140000066e0020009c0000066e02008041000000c002200210000000000112019f0000067d011001c70000800d020000390000000103000039000006bc0400004119b419aa0000040f0000000100200190000000860000613d00000007020000290000000102200039000000060020006c00000c580000413d000005bc0000013d0000001f0530018f0000067006300198000000400200043d000000000462001900000ec60000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b00000df00000c13d00000ec60000013d000006d002000041000001400020043f000001440010043f000006d101000041000019b600010430000006df01000041000000000010043f0000002201000039000000040010043f000006be01000041000019b600010430000006e40200004100000a7d0000013d0000002404200039000000000034043500000678030000410000000000320435000000040320003900000000001304350000066e0020009c0000066e02008041000000400120021000000679011001c7000019b6000104300000000f010000290000000001010433000000000001004b00000efd0000c13d000000400300043d0000067f010000410000000000130435000000240130003900000000020004100000000000210435001100000003001d0000000401300039000000000021043500000000010004140000000c02000029000000040020008c00000e230000c13d0000000103000031000000200030008c0000002004000039000000000403401900000e4f0000013d00000011020000290000066e0020009c0000066e0200804100000040022002100000066e0010009c0000066e01008041000000c001100210000000000121019f00000679011001c70000000c0200002919b419af0000040f000000000301001900000060033002700000066e03300197000000200030008c000000200400003900000000040340190000001f0640018f0000002007400190000000110570002900000e3e0000613d000000000801034f0000001109000029000000008a08043c0000000009a90436000000000059004b00000e3a0000c13d000000000006004b00000e4b0000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f0003000000010355000000010020019000000ebb0000613d0000001f01400039000000600210018f0000001101200029000000000021004b00000000020000390000000102004039000006720010009c000001240000213d0000000100200190000001240000c13d000000400010043f000000200040008c000000860000413d00000011020000290000000002020433000000000002004b000013ef0000c13d0000004402100039000000010400008a000000000042043500000020021000390000068004000041000000000042043500000024041000390000000005000410000000000054043500000044040000390000000000410435000006810010009c000001240000213d000000c004100039000000400040043f000000a0051000390000068204000041001000000005001d000000000045043500000080051000390000002004000039000f00000005001d0000000000450435000000000401043300000000010004140000000c05000029000000040050008c00000e8e0000613d0000066e0020009c0000066e0200804100000040022002100000066e0040009c0000066e040080410000006003400210000000000223019f0000066e0010009c0000066e01008041000000c001100210000000000112019f0000000c0200002919b419aa0000040f000b000100200193000300000001035500000060011002700001066e0010019d0000066e03100197000000000003004b0000106c0000c13d001100600000003d000e00800000003d000000110100002900000000010104330000000b0000006b000010a00000c13d000000000001004b000010d20000c13d000000400100043d000006880200004100000000002104350000000402100039000000200300003900000000003204350000000f020000290000000002020433000000240310003900000000002304350000004403100039000000000002004b000000100700002900000eae0000613d000000000400001900000000053400190000000006740019000000000606043300000000006504350000002004400039000000000024004b00000ea70000413d0000001f0420003900000700044001970000000002320019000000000002043500000044024000390000066e0020009c0000066e0200804100000060022002100000066e0010009c0000066e010080410000004001100210000000000112019f000019b6000104300000001f0530018f0000067006300198000000400200043d000000000462001900000ec60000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b00000ec20000c13d000000000005004b00000ed30000613d000000000161034f0000000305500210000000000604043300000000065601cf000000000656022f000000000101043b0000010005500089000000000151022f00000000015101cf000000000161019f000000000014043500000060013002100000066e0020009c0000066e020080410000004002200210000000000112019f000019b600010430000000000001004b00000f490000613d000000a00200043d000006c102200197000000000021004b00000f490000813d0000001101000029000000000010043f0000000701000039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000000860000613d000000000101043b000000800200003919b418640000040f000001200100043d000006c101100197000000e00200043d000000000002004b00000f8b0000c13d000000000001004b00000ef90000c13d000001000100043d000006c10010019800000f910000613d000000400200043d001100000002001d000006c3010000410000109c0000013d000000000200001900000f050000013d000000110200002900000001022000390000000f010000290000000001010433000000000012004b00000e110000813d001100000002001d00000005012002100000000e011000290000000001010433000006710310019800000eff0000613d000000000030043f0000000301000039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c70000801002000039001000000003001d19b419af0000040f00000010040000290000000100200190000000860000613d000000000101043b000000000101041a000000000001004b00000eff0000c13d0000000203000039000000000103041a000006720010009c000001240000213d0000000102100039000000000023041b0000067c0110009a000000000041041b000000000103041a000d00000001001d000000000040043f0000000301000039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f00000010030000290000000100200190000000860000613d000000000101043b0000000d02000029000000000021041b000000400100043d00000000003104350000066e0010009c0000066e01008041000000400110021000000000020004140000066e0020009c0000066e02008041000000c002200210000000000112019f0000067d011001c70000800d0200003900000001030000390000067e0400004119b419aa0000040f000000010020019000000eff0000c13d000000860000013d000000400200043d001100000002001d000006c20100004100000000001204350000000402200039000000800100003919b416200000040f00000011010000290000066e0010009c0000066e0100804100000040011002100000068c011001c7000019b60001043000000000030700190000002001300039000000100200002900000000002104350000002001000039000000000013043500000040013000390000000a021000290000000f0300002900000002033003670000000a0000006b00000f680000613d000000000403034f0000000005010019000000004604043c0000000005650436000000000025004b00000f640000c13d0000000e0000006b00000f760000613d0000000a033003600000000e040000290000000304400210000000000502043300000000054501cf000000000545022f000000000303043b0000010004400089000000000343022f00000000034301cf000000000353019f0000000000320435000000100110002900000000000104350000066e0070009c0000066e0700804100000040017002100000000b02000029000006da0020009c000006da020080410000006002200210000000000112019f00000000020004140000066e0020009c0000066e02008041000000c002200210000000000121019f000006db0110009a0000800d020000390000000203000039000006dc04000041000000110500002900000a360000013d000000000001004b000010990000613d000001000200043d000006c102200197000000000021004b000010990000813d0000001101000029000000000010043f0000000701000039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000000860000613d000000000101043b0000000201100039000000e00200003919b418640000040f000000400100043d00000011020000290000000002210436000000800300043d000000000003004b0000000003000039000000010300c0390000000000320435000000a00200043d000006c10220019700000040031000390000000000230435000000c00200043d000006c10220019700000060031000390000000000230435000000e00200043d000000000002004b0000000002000039000000010200c03900000080031000390000000000230435000001000200043d000006c102200197000000a0031000390000000000230435000001200200043d000006c102200197000000c00310003900000000002304350000066e0010009c0000066e01008041000000400110021000000000020004140000066e0020009c0000066e02008041000000c002200210000000000112019f000006d2011001c70000800d020000390000000103000039000006d30400004100000a360000013d000000000200001900000fd40000013d00000010020000290000000102200039000000800100043d000000000012004b000004800000813d001000000002001d0000000501200210000000a00110003900000000010104330000067101100197001100000001001d000000000010043f0000000301000039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000000860000613d000000000101043b000000000301041a000000000003004b00000fcf0000613d0000000201000039000000000201041a000000000002004b000013ef0000613d000000010130008a000000000023004b0000100c0000613d000000000012004b000010320000a13d000006f10130009a000006f10220009a000000000202041a000000000021041b000000000020043f0000000301000039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c70000801002000039000d00000003001d19b419af0000040f0000000d030000290000000100200190000000860000613d000000000101043b000000000031041b0000000201000039000000000301041a000000000003004b000010380000613d000000010130008a000006f10230009a000000000002041b0000000202000039000000000012041b0000001101000029000000000010043f0000000301000039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000000860000613d000000000101043b000000000001041b000000400100043d000000110200002900000000002104350000066e0010009c0000066e01008041000000400110021000000000020004140000066e0020009c0000066e02008041000000c002200210000000000112019f0000067d011001c70000800d020000390000000103000039000006f20400004119b419aa0000040f000000010020019000000fcf0000c13d000000860000013d000006df01000041000000000010043f0000003201000039000000040010043f000006be01000041000019b600010430000006df01000041000000000010043f0000003101000039000000040010043f000006be01000041000019b600010430000000400300043d001100000003001d000000240130003900000040020000390000000000210435000006ca01000041000000000013043500000004013000390000000b0200002900000000002104350000004402300039000000100100002919b414550000040f000000110200002900000000012100490000066e0010009c0000066e010080410000066e0020009c0000066e0200804100000060011002100000004002200210000000000121019f000019b6000104300000001001000029000000000010043f0000000301000039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000000860000613d000000000101043b000000000101041a000000000001004b00000b3d0000c13d000000400100043d000006e5020000410000000000210435000000040210003900000010030000290000038e0000013d000006720030009c000001240000213d0000001f0230003900000700022001970000003f022000390000070002200197000000400400043d0000000002240019001100000004001d000000000042004b00000000040000390000000104004039000006720020009c000001240000213d0000000100400190000001240000c13d000000400020043f0000001102000029000000000432043600000700023001980000001f0330018f000e00000004001d000000000124001900000003040003670000108b0000613d000000000504034f0000000e06000029000000005705043c0000000006760436000000000016004b000010870000c13d000000000003004b00000e920000613d000000000224034f0000000303300210000000000401043300000000043401cf000000000434022f000000000202043b0000010003300089000000000232022f00000000023201cf000000000242019f000000000021043500000e920000013d000000400200043d001100000002001d000006c20100004100000000001204350000000402200039000000e00100003900000f4f0000013d000000000001004b000012280000c13d000006830100004100000000001004430000000c01000029000000040010044300000000010004140000066e0010009c0000066e01008041000000c00110021000000684011001c7000080020200003919b419af0000040f0000000100200190000013e20000613d000000000101043b000000000001004b000012240000c13d000000400100043d00000044021000390000068b03000041000000000032043500000024021000390000001d0300003900000000003204350000068802000041000000000021043500000004021000390000002003000039000000000032043500000f510000013d0000000401000039000000000301041a000006f4010000410000000d0400002900000000001404350000000f01000029000000000012043500000024014000390000000002000411000000000021043500000000010004140000067102300197000000040020008c000010e10000c13d0000000103000031000000200030008c000000200400003900000000040340190000110c0000013d0000000e020000290000066e0020009c0000066e0200804100000040022002100000066e0010009c0000066e010080410000006001100210000000000121019f000019b600010430000000400100043d000006bd02000041000000000021043500000004021000390000000b030000290000038e0000013d0000000d030000290000066e0030009c0000066e0300804100000040033002100000066e0010009c0000066e01008041000000c001100210000000000131019f00000679011001c719b419af0000040f000000000301001900000060033002700000066e03300197000000200030008c000000200400003900000000040340190000001f0640018f00000020074001900000000d05700029000010fb0000613d000000000801034f0000000d09000029000000008a08043c0000000009a90436000000000059004b000010f70000c13d000000000006004b000011080000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f00030000000103550000000100200190000011370000613d0000001f01400039000000600210018f0000000d01200029000000000021004b00000000020000390000000102004039000006720010009c000001240000213d0000000100200190000001240000c13d000000400010043f000000200030008c000000860000413d0000000d020000290000000002020433000000000002004b0000000003000039000000010300c039000000000032004b000000860000c13d000000000002004b000012850000c13d000006e7020000410000000000210435000000040210003900000000030004110000038e0000013d0000000401000039000000000301041a000006e6010000410000000e0400002900000000001404350000001001000029000000000012043500000000010004140000067102300197000000040020008c000011430000c13d0000000103000031000000200030008c000000200400003900000000040340190000116e0000013d0000001f0530018f0000067006300198000000400200043d000000000462001900000ec60000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b0000113e0000c13d00000ec60000013d0000000e030000290000066e0030009c0000066e0300804100000040033002100000066e0010009c0000066e01008041000000c001100210000000000131019f000006be011001c719b419af0000040f000000000301001900000060033002700000066e03300197000000200030008c000000200400003900000000040340190000001f0640018f00000020074001900000000e057000290000115d0000613d000000000801034f0000000e09000029000000008a08043c0000000009a90436000000000059004b000011590000c13d000000000006004b0000116a0000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f000300000001035500000001002001900000126c0000613d0000001f01400039000000600210018f0000000e01200029000000000021004b00000000020000390000000102004039000006720010009c000001240000213d0000000100200190000001240000c13d000000400010043f000000200030008c000000860000413d0000000e020000290000000002020433000006c80020009c000000860000813d0000000003000411000000000023004b000012c50000c13d00000002010003670000000f02100360000000000202043b000006720020009c000000860000213d0000000f030000290000004003300039000000000131034f000000000101043b001000000001001d000000000020043f0000000701000039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000000860000613d000000000101043b0000001002000029000000110300002919b418d10000040f000006830100004100000000001004430000001101000029000000040010044300000000010004140000066e0010009c0000066e01008041000000c00110021000000684011001c7000080020200003919b419af0000040f0000000100200190000013e20000613d000000000101043b000000000001004b000000860000613d000000400300043d000000240130003900000010020000290000000000210435000006e801000041000000000013043500000000010004100000067101100197000e00000003001d0000000402300039000000000012043500000000010004140000001102000029000000040020008c000011cc0000613d0000000e020000290000066e0020009c0000066e0200804100000040022002100000066e0010009c0000066e01008041000000c001100210000000000121019f00000679011001c7000000110200002919b419aa0000040f000000000301001900000060033002700001066e0030019d00030000000103550000000100200190000013500000613d0000000e01000029000006720010009c000001240000213d0000000e02000029000000400020043f000000100100002900000000001204350000066e0020009c0000066e02008041000000400120021000000000020004140000066e0020009c0000066e02008041000000c002200210000000000112019f0000067d011001c70000800d020000390000000203000039000006e904000041000000000500041119b419aa0000040f0000000100200190000000860000613d0000000f010000290000000201100367000000000101043b000006720010009c000000860000213d19b415c50000040f000000400200043d001100000002001d0000000002000412001700000002001d001600200000003d001000000001001d000080050100003900000044030000390000000004000415000000170440008a0000000504400210000006ce0200004119b4198c0000040f000000ff0310018f0000001101000029000000200210003900000000003204350000002002000039000000000021043519b414a80000040f000000400100043d000e00000001001d19b414a80000040f0000000e020000290000002001200039000f00000001001d00000011030000290000000000310435000000100100002900000000001204350000000003020019000000400400043d001100000004001d00000020010000390000000002140436000000000103043300000040030000390000000000320435000000600240003919b414550000040f000000000201001900000011040000290000000001420049000000200310008a0000000f0100002900000000010104330000004004400039000000000034043519b414550000040f000000110200002900000000012100490000066e0020009c0000066e0200804100000040022002100000066e0010009c0000066e010080410000006001100210000000000121019f000019b50001042e00000011010000290000000001010433000000000001004b000012350000613d000006850010009c000000860000213d000000200010008c000000860000413d0000000e010000290000000001010433000000000001004b0000000002000039000000010200c039000000000021004b000000860000c13d000000000001004b000012500000613d000000800100043d00000140000004430000016000100443000000a00100043d00000020030000390000018000300443000001a000100443000000c00100043d0000004002000039000001c000200443000001e0001004430000006001000039000000e00200043d000002000010044300000220002004430000010000300443000000040100003900000120001004430000068a01000041000019b50001042e00000009010000290000000001010433000000400200043d000006c4030000410000000000320435000006720110019700000a740000013d000000400100043d00000064021000390000068603000041000000000032043500000044021000390000068703000041000000000032043500000024021000390000002a030000390000000000320435000006880200004100000000002104350000000402100039000000200300003900000000003204350000066e0010009c0000066e01008041000000400110021000000689011001c7000019b600010430000000400200043d001100000002001d000006c3010000410000127b0000013d000000400300043d001100000003001d000006c302000041000012820000013d0000001f0530018f0000067006300198000000400200043d000000000462001900000ec60000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000012730000c13d00000ec60000013d000000400200043d001100000002001d000006c20100004100000000001204350000000402200039000000000103001900000f4f0000013d000000400300043d001100000003001d000006c2020000410000000000230435000000040230003900000f4f0000013d00000002020003670000000e01200360000000000101043b000006720010009c000000860000213d0000000e030000290000008003300039000000000332034f000000000303043b0000000004000031000000110540006a000000230550008a00000673065001970000067307300197000000000867013f000000000067004b00000000060000190000067306004041000000000053004b00000000050000190000067305008041000006730080009c000000000605c019000000000006004b000000860000c13d0000001105000029000f00040050003d0000000f05300029000000000252034f000000000302043b000006720030009c000000860000213d0000000004340049000000200250003900000673054001970000067306200197000000000756013f000000000056004b00000000050000190000067305004041000000000042004b00000000040000190000067304002041000006730070009c000000000504c019000000000005004b000000860000c13d19b415200000040f000000000001004b000012c90000c13d0000001101000029000000a4021000390000000f0100002919b414be0000040f000006f903000041000000400500043d001100000005001d0000000000350435000000040350003900000020040000390000000000430435000000240350003919b415910000040f0000104b0000013d000006e702000041000000000021043500000004021000390000038e0000013d00000002010003670000000e02100360000000000202043b000006720020009c000000860000213d0000000e03000029000d00400030003d0000000d01100360000000000101043b000e00000001001d000000000020043f0000000701000039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000000860000613d000000000101043b00000002011000390000000e02000029000000100300002919b418d10000040f0000000d01000029000d00600010003d00000002030003670000000d01300360000000000101043b0000000004000031000000110240006a000000230220008a00000673052001970000067306100197000000000756013f000000000056004b00000000050000190000067305004041000000000021004b00000000020000190000067302008041000006730070009c000000000502c019000000000005004b000000860000c13d0000000f01100029000000000213034f000000000202043b000006720020009c000000860000213d0000000005240049000000200610003900000673015001970000067307600197000000000817013f000000000017004b00000000010000190000067301004041000000000056004b00000000050000190000067305002041000006730080009c000000000105c019000000000001004b000000860000c13d0000001f0120003900000700011001970000003f011000390000070005100197000000400100043d0000000005510019000000000015004b00000000080000390000000108004039000006720050009c000001240000213d0000000100800190000001240000c13d000000400050043f00000000052104360000000008620019000000000048004b000000860000213d000000000463034f00000700062001980000001f0720018f00000000036500190000132a0000613d000000000804034f0000000009050019000000008a08043c0000000009a90436000000000039004b000013260000c13d000000000007004b000013370000613d000000000464034f0000000306700210000000000703043300000000076701cf000000000767022f000000000404043b0000010006600089000000000464022f00000000046401cf000000000474019f0000000000430435000000000225001900000000000204350000000002010433000000200020008c0000135d0000613d000000000002004b000013610000c13d000006ce010000410000000000100443000000000100041200000004001004430000002001000039000000240010044300000000010004140000066e0010009c0000066e01008041000000c001100210000006e0011001c7000080050200003919b419af0000040f0000000100200190000013e20000613d000000000101043b001100000001001d0000136a0000013d0000066e033001970000001f0530018f0000067006300198000000400200043d000000000462001900000ec60000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000013580000c13d00000ec60000013d0000000002050433001100000002001d000000ff0020008c0000136a0000a13d000000400400043d001100000004001d000006f802000041000000000024043500000004024000390000002003000039000000000032043500000024024000390000104a0000013d000006ce010000410000000000100443000000000100041200000004001004430000002001000039000000240010044300000000010004140000066e0010009c0000066e01008041000000c001100210000006e0011001c7000080050200003919b419af0000040f0000000100200190000013e20000613d0000001102000029000000ff0220018f000000000301043b000000ff0430018f000000000142004b000013e30000c13d0000000e01000029001100000001001d0000000d01000029000000800110008a000e00000001001d0000000201100367000000000101043b000f00000001001d000006710010009c000000860000213d000006830100004100000000001004430000001001000029000000040010044300000000010004140000066e0010009c0000066e01008041000000c00110021000000684011001c7000080020200003919b419af0000040f0000000100200190000013e20000613d000000000101043b000000000001004b000000860000613d000000400300043d000000240130003900000011020000290000000000210435000006f6010000410000000000130435000d00000003001d00000004013000390000000f02000029000000000021043500000000010004140000001002000029000000040020008c000013b80000613d0000000d020000290000066e0020009c0000066e0200804100000040022002100000066e0010009c0000066e01008041000000c001100210000000000121019f00000679011001c7000000100200002919b419aa0000040f000000000301001900000060033002700001066e0030019d00030000000103550000000100200190000013fb0000613d0000000d01000029000006720010009c000001240000213d0000000d01000029000000400010043f0000000e010000290000000201100367000000000601043b000006710060009c000000860000213d00000011010000290000000d0200002900000000001204350000066e0020009c0000066e02008041000000400120021000000000020004140000066e0020009c0000066e02008041000000c002200210000000000112019f0000067d011001c70000800d020000390000000303000039000006f704000041000000000500041119b419aa0000040f0000000100200190000000860000613d000000400100043d001000000001001d19b4149d0000040f000000110200002900000010010000290000000000210435000000400100043d00000000002104350000066e0010009c0000066e010080410000004001100210000006eb011001c7000019b50001042e000000000001042f000013ec0000a13d000000ff0010008c000013ef0000213d0000004d0010008c0000142a0000213d000000000001004b000014080000c13d0000000102000039000014130000013d0000000001240049000000ff0010008c000013f50000a13d000006df01000041000000000010043f0000001101000039000000040010043f000006be01000041000019b6000104300000004e0010008c0000142a0000813d000000000001004b000014150000c13d0000000102000039000014270000013d0000066e033001970000001f0530018f0000067006300198000000400200043d000000000462001900000ec60000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000014030000c13d00000ec60000013d0000000a030000390000000102000039000000010010019000000000043300a9000000010300603900000000022300a9000000010110027200000000030400190000140a0000c13d000000000002004b000014210000613d0000000e012000f9000013800000013d0000000a0500003900000001020000390000000004010019000000010040019000000000065500a9000000010500603900000000022500a900000001044002720000000005060019000014180000c13d000000000002004b000014270000c13d000006df01000041000000000010043f0000001201000039000000040010043f000006be01000041000019b60001043000000701022001290000000e0020006c0000143c0000813d000000400200043d001000000002001d000006f5010000410000000000120435000000040120003900000011020000290000000e0400002919b4162e0000040f000000100200002900000000012100490000066e0010009c0000066e0100804100000060011002100000066e0020009c0000066e020080410000004002200210000000000121019f000019b600010430000000ff0210018f0000004d0020008c000013ef0000213d000000000002004b000014430000c13d00000001010000390000144c0000013d0000000a030000390000000101000039000000010020019000000000043300a9000000010300603900000000011300a900000001022002720000000003040019000014450000c13d0000000e0000006b001100000000001d000013810000613d0000000e031000b9001100000003001d0000000e023000fa000000000012004b000013810000613d000013ef0000013d00000000430104340000000001320436000000000003004b000014610000613d000000000200001900000000052100190000000006240019000000000606043300000000006504350000002002200039000000000032004b0000145a0000413d000000000231001900000000000204350000001f0230003900000700022001970000000001210019000000000001042d000006850010009c000014820000213d000000430010008c000014820000a13d00000002020003670000000403200370000000000403043b000006720040009c000014820000213d0000002403200370000000000503043b000006720050009c000014820000213d0000002303500039000000000013004b000014820000813d0000000403500039000000000232034f000000000302043b000006720030009c000014820000213d00000024025000390000000005320019000000000015004b000014820000213d0000000001040019000000000001042d0000000001000019000019b6000104300000000043020434000006c103300197000000000331043600000000040404330000066e04400197000000000043043500000040032000390000000003030433000000000003004b0000000003000039000000010300c0390000004004100039000000000034043500000060032000390000000003030433000006c1033001970000006004100039000000000034043500000080022000390000000002020433000006c10220019700000080031000390000000000230435000000a001100039000000000001042d000007030010009c000014a20000813d0000002001100039000000400010043f000000000001042d000006df01000041000000000010043f0000004101000039000000040010043f000006be01000041000019b600010430000007040010009c000014ad0000813d0000004001100039000000400010043f000000000001042d000006df01000041000000000010043f0000004101000039000000040010043f000006be01000041000019b600010430000007050010009c000014b80000813d000000a001100039000000400010043f000000000001042d000006df01000041000000000010043f0000004101000039000000040010043f000006be01000041000019b6000104300000000204000367000000000224034f000000000202043b000000000300003100000000051300490000001f0550008a00000673065001970000067307200197000000000867013f000000000067004b00000000060000190000067306002041000000000052004b00000000050000190000067305004041000006730080009c000000000605c019000000000006004b000014e60000613d0000000001120019000000000214034f000000000202043b000006720020009c000014e60000213d0000000003230049000000200110003900000673043001970000067305100197000000000645013f000000000045004b00000000040000190000067304004041000000000031004b00000000030000190000067303002041000006730060009c000000000403c019000000000004004b000014e60000c13d000000000001042d0000000001000019000019b6000104300000000004010019000007060020009c000015180000813d0000001f0120003900000700011001970000003f011000390000070005100197000000400100043d0000000005510019000000000015004b00000000070000390000000107004039000006720050009c000015180000213d0000000100700190000015180000c13d000000400050043f00000000052104360000000007420019000000000037004b0000151e0000213d00000700062001980000001f0720018f00000002044003670000000003650019000015080000613d000000000804034f0000000009050019000000008a08043c0000000009a90436000000000039004b000015040000c13d000000000007004b000015150000613d000000000464034f0000000306700210000000000703043300000000076701cf000000000767022f000000000404043b0000010006600089000000000464022f00000000046401cf000000000474019f000000000043043500000000022500190000000000020435000000000001042d000006df01000041000000000010043f0000004101000039000000040010043f000006be01000041000019b6000104300000000001000019000019b6000104300003000000000002000300000003001d000200000002001d0000067201100197000000000010043f0000000701000039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000207000029000000030a0000290000000100200190000015890000613d0000000003000031000000000601043b0000070600a0009c0000158b0000813d0000001f01a0003900000700011001970000003f011000390000070002100197000000400100043d0000000002210019000000000012004b00000000050000390000000105004039000006720020009c0000158b0000213d00000001005001900000158b0000c13d000100000006001d000000400020043f0000000002a1043600000000057a0019000000000035004b000015890000213d0000070004a001980000001f05a0018f00000002067003670000000003420019000015540000613d000000000706034f0000000008020019000000007907043c0000000008980436000000000038004b000015500000c13d000000000005004b000015610000613d000000000446034f0000000305500210000000000603043300000000065601cf000000000656022f000000000404043b0000010005500089000000000454022f00000000045401cf000000000464019f00000000004304350000000003a2001900000000000304350000066e0020009c0000066e02008041000000400220021000000000010104330000066e0010009c0000066e010080410000006001100210000000000121019f00000000020004140000066e0020009c0000066e02008041000000c002200210000000000112019f000006b7011001c7000080100200003919b419af0000040f0000000100200190000015890000613d000000000101043b000000000010043f00000001010000290000000601100039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000015890000613d000000000101043b000000000101041a000000000001004b0000000001000039000000010100c039000000000001042d0000000001000019000019b600010430000006df01000041000000000010043f0000004101000039000000040010043f000006be01000041000019b600010430000000000323043600000700062001980000001f0720018f000000000563001900000002011003670000159d0000613d000000000801034f0000000009030019000000008a08043c0000000009a90436000000000059004b000015990000c13d000000000007004b000015aa0000613d000000000161034f0000000306700210000000000705043300000000076701cf000000000767022f000000000101043b0000010006600089000000000161022f00000000016101cf000000000171019f0000000000150435000000000123001900000000000104350000001f0120003900000700011001970000000001130019000000000001042d000000400100043d000007050010009c000015bf0000813d000000a002100039000000400020043f000000800210003900000000000204350000006002100039000000000002043500000040021000390000000000020435000000200210003900000000000204350000000000010435000000000001042d000006df01000041000000000010043f0000004101000039000000040010043f000006be01000041000019b60001043000030000000000020000067201100197000000000010043f0000000701000039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000016120000613d000000000101043b0000000405100039000000000205041a000000010320019000000001062002700000007f0660618f0000001f0060008c00000000040000390000000104002039000000000043004b000016140000c13d000000400100043d0000000004610436000000000003004b000015fe0000613d000100000004001d000200000006001d000300000001001d000000000050043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067d011001c7000080100200003919b419af0000040f0000000100200190000016120000613d0000000206000029000000000006004b000016040000613d000000000201043b0000000005000019000000030100002900000001070000290000000003570019000000000402041a000000000043043500000001022000390000002005500039000000000065004b000015f60000413d000016060000013d00000702022001970000000000240435000000000006004b00000020050000390000000005006039000016060000013d000000000500001900000003010000290000003f0350003900000700023001970000000003120019000000000023004b00000000020000390000000102004039000006720030009c0000161a0000213d00000001002001900000161a0000c13d000000400030043f000000000001042d0000000001000019000019b600010430000006df01000041000000000010043f0000002201000039000000040010043f000006be01000041000019b600010430000006df01000041000000000010043f0000004101000039000000040010043f000006be01000041000019b6000104300000000043010434000000000003004b0000000003000039000000010300c03900000000033204360000000004040433000006c1044001970000000000430435000000400220003900000040011000390000000001010433000006c1011001970000000000120435000000000001042d00000040051000390000000000450435000000ff0330018f00000020041000390000000000340435000000ff0220018f00000000002104350000006001100039000000000001042d0007000000000002000400000001001d000600000002001d0000000021020434000000000001004b000017500000613d0000066e0010009c0000066e0100804100000060011002100000066e0020009c000500000002001d0000066e020080410000004002200210000000000121019f00000000020004140000066e0020009c0000066e02008041000000c002200210000000000112019f000006b7011001c7000080100200003919b419af0000040f0000000100200190000017480000613d000000000101043b000700000001001d00000004010000290000067201100197000200000001001d000000000010043f0000000701000039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000017480000613d000000000201043b0000000701000029000000000010043f000400000002001d0000000601200039000300000001001d000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000017480000613d000000000101043b000000000101041a000000000001004b000017580000c13d00000004010000290000000502100039000000000102041a000007060010009c0000174a0000813d000100000001001d0000000101100039000000000012041b000400000002001d000000000020043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067d011001c7000080100200003919b419af0000040f00000001002001900000000702000029000017480000613d000000000101043b0000000101100029000000000021041b0000000401000029000000000101041a000400000001001d000000000020043f0000000301000029000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f00000001002001900000000702000029000017480000613d000000000101043b0000000403000029000000000031041b000000000020043f0000000801000039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000017480000613d000000000801043b00000006010000290000000004010433000006720040009c0000174a0000213d000000000108041a000000010210019000000001031002700000007f0330618f0000001f0030008c00000000010000390000000101002039000000000012004b00000005070000290000176f0000c13d000000200030008c000400000008001d000700000004001d000016db0000413d000300000003001d000000000080043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067d011001c7000080100200003919b419af0000040f0000000100200190000017480000613d00000007040000290000001f024000390000000502200270000000200040008c0000000002004019000000000301043b00000003010000290000001f01100039000000050110027000000000011300190000000002230019000000000012004b00000005070000290000000408000029000016db0000813d000000000002041b0000000102200039000000000012004b000016d70000413d0000001f0040008c000000200a00008a000000200b0000390000170b0000a13d000000000080043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067d011001c7000080100200003919b419af0000040f0000000100200190000017480000613d0000000709000029000000200a00008a0000000002a90170000000000101043b000000200b000039000017410000613d000000010320008a000000050330027000000000043100190000002003000039000000010440003900000005070000290000000606000029000000040800002900000000056300190000000005050433000000000051041b00000020033000390000000101100039000000000041004b000016f70000c13d000000000092004b000017080000813d0000000302900210000000f80220018f000007010220027f000007010220016700000000036300190000000003030433000000000223016f000000000021041b000000010190021000000001011001bf000017170000013d000000000004004b0000170f0000613d0000000001070433000017100000013d000000000100001900000006060000290000000302400210000007010220027f0000070102200167000000000121016f0000000102400210000000000121019f000000000018041b000000400100043d0000000003b10436000000000206043300000000002304350000004003100039000000000002004b000017270000613d000000000400001900000000053400190000000006740019000000000606043300000000006504350000002004400039000000000024004b000017200000413d0000001f042000390000000004a4016f0000000002230019000000000002043500000040024000390000066e0020009c0000066e0200804100000060022002100000066e0010009c0000066e010080410000004001100210000000000112019f00000000020004140000066e0020009c0000066e02008041000000c002200210000000000121019f000006b7011001c70000800d020000390000000203000039000006cb04000041000000020500002919b419aa0000040f0000000100200190000017480000613d000000000001042d00000000030b0019000000050700002900000006060000290000000408000029000000000092004b000017000000413d000017080000013d0000000001000019000019b600010430000006df01000041000000000010043f0000004101000039000000040010043f000006be01000041000019b600010430000000400100043d000006cd0200004100000000002104350000066e0010009c0000066e01008041000000400110021000000677011001c7000019b600010430000000400300043d000700000003001d000000240130003900000040020000390000000000210435000006ca0100004100000000001304350000000401300039000000020200002900000000002104350000004402300039000000060100002919b414550000040f000000070200002900000000012100490000066e0010009c0000066e010080410000066e0020009c0000066e0200804100000060011002100000004002200210000000000121019f000019b600010430000006df01000041000000000010043f0000002201000039000000040010043f000006be01000041000019b6000104300005000000000002000000400300043d000007050030009c000017bc0000813d000000a002300039000000400020043f00000080023000390000000000020435000000600230003900000000000204350000004002300039000000000002043500000020023000390000000000020435000000000003043500000060021000390000000002020433000100000002001d000500000001001d0000000012010434000300000002001d000200000001001d0000000001010433000400000001001d000006c601000041000000000010044300000000010004140000066e0010009c0000066e01008041000000c001100210000006c7011001c70000800b0200003919b419af0000040f0000000100200190000017c20000613d00000004020000290000066e04200197000000000601043b000000000346004b0000000501000029000017b60000413d00000080021000390000000002020433000006c10520019700000000023500a9000000000046004b000017a70000613d00000000033200d9000000000053004b000017b60000c13d0000000303000029000006c103300197000000000032001a000017b60000413d00000000023200190000000103000029000006c103300197000006c104200197000000000023004b000000000304801900000000003104350000066e0260019700000002030000290000000000230435000000000001042d000006df01000041000000000010043f0000001101000039000000040010043f000006be01000041000019b600010430000006df01000041000000000010043f0000004101000039000000040010043f000006be01000041000019b600010430000000000001042f000000000010043f0000000601000039000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000017d50000613d000000000101043b000000000101041a000000000001004b0000000001000039000000010100c039000000000001042d0000000001000019000019b6000104300006000000000002000300000002001d000000000020043f000600000001001d0000000101100039000400000001001d000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000018500000613d0000000603000029000000000101043b000000000101041a000000000001004b0000184e0000613d000000000203041a000000000002004b000018520000613d000000000021004b000500000001001d0000182c0000613d000200000002001d000000000030043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067d011001c7000080100200003919b419af0000040f0000000100200190000018500000613d00000005020000290001000100200092000000000101043b0000000604000029000000000204041a000000010020006c0000000003040019000018580000a13d0000000202000029000000010220008a0000000001120019000000000101041a000200000001001d000000000030043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067d011001c7000080100200003919b419af0000040f0000000100200190000018500000613d000000000101043b00000001011000290000000202000029000000000021041b000000000020043f0000000401000029000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000018500000613d000000000101043b0000000502000029000000000021041b0000000603000029000000000103041a000500000001001d000000000001004b0000185e0000613d000000000030043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067d011001c7000080100200003919b419af0000040f0000000100200190000018500000613d0000000502000029000000010220008a000000000101043b0000000001210019000000000001041b0000000601000029000000000021041b0000000301000029000000000010043f0000000401000029000000200010043f00000000010004140000066e0010009c0000066e01008041000000c0011002100000067b011001c7000080100200003919b419af0000040f0000000100200190000018500000613d000000000101043b000000000001041b0000000101000039000000000001042d0000000001000019000000000001042d0000000001000019000019b600010430000006df01000041000000000010043f0000001101000039000000040010043f000006be01000041000019b600010430000006df01000041000000000010043f0000003201000039000000040010043f000006be01000041000019b600010430000006df01000041000000000010043f0000003101000039000000040010043f000006be01000041000019b6000104300003000000000002000100000002001d000300000001001d000000000101041a000200000001001d000006c601000041000000000010044300000000010004140000066e0010009c0000066e01008041000000c001100210000006c7011001c70000800b0200003919b419af0000040f0000000100200190000018c60000613d000000020800002900000080028002700000066e03200197000000000201043b000000000532004b0000000307000029000018c00000413d00000001017000390000187f0000c13d000000000207041a000018910000013d000000000301041a000000800630027000000000045600a900000000055400d9000000000065004b000018c00000c13d000006c105800197000000000054001a000018c00000413d000006c1033001970000000004540019000000000043004b0000000003048019000006750480019700000080022002100000070702200197000000000242019f000000000232019f000000010600002900000020036000390000000004030433000006c104400197000006c105200197000000000054004b00000000050440190000070802200197000000000225019f0000000005060433000000000005004b0000000005000019000006c80500c041000000000252019f000000000027041b000000400260003900000000050204330000008005500210000000000445019f000000000041041b0000000001000039000000010100c039000000400400043d00000000011404360000000003030433000006c10330019700000000003104350000000001020433000006c101100197000000400240003900000000001204350000066e0040009c0000066e04008041000000400140021000000000020004140000066e0020009c0000066e02008041000000c002200210000000000112019f00000709011001c70000800d0200003900000001030000390000070a0400004119b419aa0000040f0000000100200190000018c70000613d000000000001042d000006df01000041000000000010043f0000001101000039000000040010043f000006be01000041000019b600010430000000000001042f0000000001000019000019b6000104300000067104400197000000400510003900000000004504350000002004100039000000000034043500000000002104350000006001100039000000000001042d0006000000000002000000000401041a000006d400400198000019250000613d000000000002004b000019250000613d000600000004001d000500000002001d000200000003001d000300000001001d0000000101100039000100000001001d000000000101041a000400000001001d000006c601000041000000000010044300000000010004140000066e0010009c0000066e01008041000000c001100210000006c7011001c70000800b0200003919b419af0000040f0000000100200190000019260000613d000000060300002900000080023002700000066e02200197000000000101043b000000000421004b000019410000413d000006c1033001970000000405000029000006c102500197000018f70000c13d000000050400002900000003050000290000190b0000013d000000000023004b000019490000213d000000800650027000000000056400a900000000044500d9000000000064004b000019410000c13d000000000035001a000019410000413d0000000003350019000000800110021000000707011001970000000305000029000000000405041a0000070b04400197000000000114019f000000000015041b000000000032004b00000000030240190000000504000029000000000042004b000019270000413d000000000143004b000019380000413d000006c101100197000000000205041a0000070d02200197000000000112019f000000000015041b000000400100043d00000000004104350000066e0010009c0000066e01008041000000400110021000000000020004140000066e0020009c0000066e02008041000000c002200210000000000112019f0000067d011001c70000800d0200003900000001030000390000070e0400004119b419aa0000040f0000000100200190000019470000613d000000000001042d000000000001042f000000400100043d0000000004010019000000040110003900000002030000290000067100300198000019510000c13d0000071203000041000000000034043500000000002104350000002401400039000000050200002900000000002104350000066e0040009c0000066e04008041000000400140021000000679011001c7000019b6000104300000000101000029000000000101041a0000008001100272000019410000613d00000005043000690000000002140019000000010220008a000000000042004b000019560000813d000006df01000041000000000010043f0000001101000039000000040010043f000006be01000041000019b6000104300000000001000019000019b600010430000000400100043d0000070c0200004100000000002104350000066e0010009c0000066e01008041000000400110021000000677011001c7000019b6000104300000071103000041000600000004001d000000000034043500000005030000290000196a0000013d00000000021200d9000000400100043d0000000005010019000000040110003900000002040000290000067100400198000019670000c13d000007100400004100000000004504350000000000210435000000240150003900000000003104350000066e0050009c0000066e05008041000000400150021000000679011001c7000019b6000104300000070f04000041000600000005001d0000000000450435000000020400002919b418c90000040f000000060200002900000000012100490000066e0010009c0000066e0100804100000060011002100000066e0020009c0000066e020080410000004002200210000000000121019f000019b600010430000000000001042f0000066e0010009c0000066e0100804100000040011002100000066e0020009c0000066e020080410000006002200210000000000112019f00000000020004140000066e0020009c0000066e02008041000000c002200210000000000112019f000006b7011001c7000080100200003919b419af0000040f00000001002001900000198a0000613d000000000101043b000000000001042d0000000001000019000019b60001043000000000050100190000000000200443000000050030008c0000199a0000413d000000040100003900000000020000190000000506200210000000000664001900000005066002700000000006060031000000000161043a0000000102200039000000000031004b000019920000413d0000066e0030009c0000066e03008041000000600130021000000000020004140000066e0020009c0000066e02008041000000c002200210000000000112019f00000713011001c7000000000205001919b419af0000040f0000000100200190000019a90000613d000000000101043b000000000001042d000000000001042f000019ad002104210000000102000039000000000001042d0000000002000019000000000001042d000019b2002104230000000102000039000000000001042d0000000002000019000000000001042d000019b400000432000019b50001042e000019b600010430000000000000000000000000000000000000000000000000000000000000000000000000ffffffff00000000000000000000000000000000000000000000000000000001ffffffe000000000000000000000000000000000000000000000000000000000ffffffe0000000000000000000000000ffffffffffffffffffffffffffffffffffffffff000000000000000000000000000000000000000000000000ffffffffffffffff80000000000000000000000000000000000000000000000000000000000000007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0ffffffffffffffffffffffff0000000000000000000000000000000000000000313ce567000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000655a7c0e000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000044000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffffdf0200000000000000000000000000000000000040000000000000000000000000bfa87805ed57dc1f0d489ce33be4c4577d74ccde357eeeee058a32c55c44a53202000000000000000000000000000000000000200000000000000000000000002640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d8dd62ed3e00000000000000000000000000000000000000000000000000000000095ea7b300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffff3f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65641806aa1896bbf26568e884a7374b41e002500962caba6a15023a8d90e8508b8302000002000000000000000000000000000000240000000000000000000000007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6f742073756363656564000000000000000000000000000000000000000000005361666545524332303a204552433230206f7065726174696f6e20646964206e08c379a00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000840000000000000000000000000000000200000000000000000000000000000140000001000000000000000000416464726573733a2063616c6c20746f206e6f6e2d636f6e747261637400000000000000000000000000000000000000000000640000000000000000000000009b15e16f00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000009a4575b800000000000000000000000000000000000000000000000000000000c0d7865400000000000000000000000000000000000000000000000000000000dc0bd97000000000000000000000000000000000000000000000000000000000e8a1da1600000000000000000000000000000000000000000000000000000000e8a1da1700000000000000000000000000000000000000000000000000000000f2fde38b00000000000000000000000000000000000000000000000000000000dc0bd97100000000000000000000000000000000000000000000000000000000e0351e1300000000000000000000000000000000000000000000000000000000c75eea9b00000000000000000000000000000000000000000000000000000000c75eea9c00000000000000000000000000000000000000000000000000000000cf7401f300000000000000000000000000000000000000000000000000000000c0d7865500000000000000000000000000000000000000000000000000000000c4bffe2b00000000000000000000000000000000000000000000000000000000acfecf9000000000000000000000000000000000000000000000000000000000b0f479a000000000000000000000000000000000000000000000000000000000b0f479a100000000000000000000000000000000000000000000000000000000b794658000000000000000000000000000000000000000000000000000000000acfecf9100000000000000000000000000000000000000000000000000000000af58d59f000000000000000000000000000000000000000000000000000000009a4575b900000000000000000000000000000000000000000000000000000000a42a7b8b00000000000000000000000000000000000000000000000000000000a7cd63b70000000000000000000000000000000000000000000000000000000054c8a4f20000000000000000000000000000000000000000000000000000000079ba5096000000000000000000000000000000000000000000000000000000008926f54e000000000000000000000000000000000000000000000000000000008926f54f000000000000000000000000000000000000000000000000000000008da5cb5b0000000000000000000000000000000000000000000000000000000079ba5097000000000000000000000000000000000000000000000000000000007d54534e0000000000000000000000000000000000000000000000000000000054c8a4f30000000000000000000000000000000000000000000000000000000062ddd3c4000000000000000000000000000000000000000000000000000000006d3d1a5800000000000000000000000000000000000000000000000000000000240028e700000000000000000000000000000000000000000000000000000000390775360000000000000000000000000000000000000000000000000000000039077537000000000000000000000000000000000000000000000000000000004c5ef0ed00000000000000000000000000000000000000000000000000000000240028e80000000000000000000000000000000000000000000000000000000024f65ee70000000000000000000000000000000000000000000000000000000001ffc9a700000000000000000000000000000000000000000000000000000000181f5a770000000000000000000000000000000000000000000000000000000021df0da70200000000000000000000000000000000000000000000000000000000000000ed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278dad89dca000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000800000000000000000fc949c7b4a13586e39d89eead2f38644f9fb3efb5a0490b14f8fc0ceab44c2515204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d8599161e670e4b000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000024000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffff5f000000000000000000000000000000000000000000000000ffffffffffffff9f00000000000000000000000000000000ffffffffffffffffffffffffffffffff8020d12400000000000000000000000000000000000000000000000000000000d68af9cc000000000000000000000000000000000000000000000000000000001d5ad3c500000000000000000000000000000000000000000000000000000000fc949c7b4a13586e39d89eead2f38644f9fb3efb5a0490b14f8fc0ceab44c250796b89b91644bc98cd93958e4c9038275d622183e25ac5af08cc6b5d9553913202000002000000000000000000000000000000040000000000000000000000000000000000000000000000010000000000000000000000000000000000000000ffffffffffffffffffffff000000000000000000000000000000000000000000393b8ad2000000000000000000000000000000000000000000000000000000007d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c28579befe00000000000000000000000000000000000000000000000000000000310ab089e4439a4c15d089f94afb7896ff553aecb10793d0ab882de59d99a32e00000000000000000000000000000000000000200000008000000000000000008e4a23d600000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002400000140000000000000000002000000000000000000000000000000000000e00000000000000000000000000350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b0000000000000000000000ff0000000000000000000000000000000000000000036b6384b5eca791c62761152d0c79bb0604c104a5fb6f4eb0703f3154bb3db0000000000000000000000000000000000000000000000000ffffffffffffff7f00000000000000000000000000000000000000000000003fffffffffffffffe0020000000000000000000000000000000000004000000080000000000000000002dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f168400000000000000000000000000000000000000000000000000000000ffffffbffdffffffffffffffffffffffffffffffffffffc000000000000000000000000052d00ee4d9bd51b40168f2afc5848837288ce258784ad914278791464b3f4d7674f23c7c00000000000000000000000000000000000000000000000000000000405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace4e487b71000000000000000000000000000000000000000000000000000000000200000200000000000000000000000000000044000000000000000000000000961c9a4f000000000000000000000000000000000000000000000000000000002cbc26bb000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffffff0000000000000000000000000000000053ad11d800000000000000000000000000000000000000000000000000000000d0d2597600000000000000000000000000000000000000000000000000000000a8d87a3b00000000000000000000000000000000000000000000000000000000728fe07b000000000000000000000000000000000000000000000000000000009dc29fac00000000000000000000000000000000000000000000000000000000696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df7a9902c7e000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000020000000000000000000000000000000000002000000080000000000000000044676b5284b809a22248eba0da87391d79098be38bb03154be88a58bf4d0917402b543c6000000000000000000000000000000000000000000000000000000008be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e02b5c74de00000000000000000000000000000000000000000000000000000000bfa87805ed57dc1f0d489ce33be4c4577d74ccde357eeeee058a32c55c44a533800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf756635f4a7b30000000000000000000000000000000000000000000000000000000083826b2b00000000000000000000000000000000000000000000000000000000a9cb113d0000000000000000000000000000000000000000000000000000000040c10f19000000000000000000000000000000000000000000000000000000009d228d69b5fdb8d273a2336f8fb8612d039631024ea9bf09c424a9503aa078f0953576f70000000000000000000000000000000000000000000000000000000024eb47e5000000000000000000000000000000000000000000000000000000004275726e5769746846726f6d4d696e74546f6b656e506f6f6c20312e352e31000000000000000000000000000000000000000000000000c0000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffff01ffc9a7000000000000000000000000000000000000000000000000000000000e64dd2900000000000000000000000000000000000000000000000000000000aff2afbf00000000000000000000000000000000000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000000000000000000000000000000000000000000000ffffffffffffffe0000000000000000000000000000000000000000000000000ffffffffffffffc0000000000000000000000000000000000000000000000000ffffffffffffff600000000000000000000000000000000000000000000000010000000000000000000000000000000000000000ffffffff00000000000000000000000000000000ffffffffffffffffffffff00ffffffff0000000000000000000000000000000002000000000000000000000000000000000000600000000000000000000000009ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c19ffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff9725942a00000000000000000000000000000000000000000000000000000000ffffffffffffffffffffffffffffffff000000000000000000000000000000001871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690ad0c8d23a0000000000000000000000000000000000000000000000000000000015279c08000000000000000000000000000000000000000000000000000000001a76572a00000000000000000000000000000000000000000000000000000000f94ebcd10000000000000000000000000000000000000000000000000000000002000002000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
