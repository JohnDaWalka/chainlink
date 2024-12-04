package burn_from_mint_token_pool

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

var BurnFromMintTokenPoolMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIBurnMintERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"localTokenDecimals\",\"type\":\"uint8\"},{\"internalType\":\"address[]\",\"name\":\"allowlist\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"}],\"name\":\"AggregateValueMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"}],\"name\":\"AggregateValueRateLimitReached\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AllowListNotEnabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BucketOverfilled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"CallerIsNotARampOnRouter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotTransferToSelf\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CursedByRMN\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"DisabledNonZeroRateLimit\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"expected\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"actual\",\"type\":\"uint8\"}],\"name\":\"InvalidDecimalArgs\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"InvalidRateLimitRate\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolData\",\"type\":\"bytes\"}],\"name\":\"InvalidRemoteChainDecimals\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidRemotePoolForChain\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidSourcePoolAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"InvalidToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeProposedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"NonExistentChain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"remoteDecimals\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"localDecimals\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"remoteAmount\",\"type\":\"uint256\"}],\"name\":\"OverflowDetected\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OwnerCannotBeZero\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"PoolAlreadyAdded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RateLimitMustBeDisabled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SenderNotAllowed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenRateLimitReached\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListAdd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListRemove\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Burned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remoteToken\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"ChainAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"ChainConfigured\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Locked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Minted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"rateLimitAdmin\",\"type\":\"address\"}],\"name\":\"RateLimitAdminSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Released\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"RemotePoolAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"RemotePoolRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldRouter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"RouterUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"TokensConsumed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"addRemotePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"removes\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"adds\",\"type\":\"address[]\"}],\"name\":\"applyAllowListUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64[]\",\"name\":\"remoteChainSelectorsToRemove\",\"type\":\"uint64[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes[]\",\"name\":\"remotePoolAddresses\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"remoteTokenAddress\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"internalType\":\"structTokenPool.ChainUpdate[]\",\"name\":\"chainsToAdd\",\"type\":\"tuple[]\"}],\"name\":\"applyChainUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowList\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowListEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getCurrentInboundRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getCurrentOutboundRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRateLimitAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getRemotePools\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"\",\"type\":\"bytes[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getRemoteToken\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRmnProxy\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRouter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSupportedChains\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTokenDecimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"decimals\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"isRemotePool\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"isSupportedChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"isSupportedToken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"originalSender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"}],\"internalType\":\"structPool.LockOrBurnInV1\",\"name\":\"lockOrBurnIn\",\"type\":\"tuple\"}],\"name\":\"lockOrBurn\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"destPoolData\",\"type\":\"bytes\"}],\"internalType\":\"structPool.LockOrBurnOutV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"originalSender\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"sourcePoolData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainTokenData\",\"type\":\"bytes\"}],\"internalType\":\"structPool.ReleaseOrMintInV1\",\"name\":\"releaseOrMintIn\",\"type\":\"tuple\"}],\"name\":\"releaseOrMint\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"destinationAmount\",\"type\":\"uint256\"}],\"internalType\":\"structPool.ReleaseOrMintOutV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"removeRemotePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundConfig\",\"type\":\"tuple\"}],\"name\":\"setChainRateLimiterConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"rateLimitAdmin\",\"type\":\"address\"}],\"name\":\"setRateLimitAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"setRouter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6101006040523480156200001257600080fd5b5060405162004c6338038062004c63833981016040819052620000359162000918565b8484848484336000816200005c57604051639b15e16f60e01b815260040160405180910390fd5b600180546001600160a01b0319166001600160a01b03848116919091179091558116156200008f576200008f8162000206565b50506001600160a01b0385161580620000af57506001600160a01b038116155b80620000c257506001600160a01b038216155b15620000e1576040516342bcdf7f60e11b815260040160405180910390fd5b6001600160a01b03808616608081905290831660c0526040805163313ce56760e01b8152905163313ce567916004808201926020929091908290030181865afa92505050801562000151575060408051601f3d908101601f191682019092526200014e9181019062000a3a565b60015b1562000192578060ff168560ff161462000190576040516332ad3e0760e11b815260ff8087166004830152821660248201526044015b60405180910390fd5b505b60ff841660a052600480546001600160a01b0319166001600160a01b038316179055825115801560e052620001dc57604080516000815260208101909152620001dc908462000280565b50620001fb935050506001600160a01b038716905030600019620003dd565b505050505062000b84565b336001600160a01b038216036200023057604051636d6c4ee560e11b815260040160405180910390fd5b600080546001600160a01b0319166001600160a01b03838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60e051620002a1576040516335f4a7b360e01b815260040160405180910390fd5b60005b82518110156200032c576000838281518110620002c557620002c562000a58565b60209081029190910101519050620002df600282620004c3565b1562000322576040516001600160a01b03821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b50600101620002a4565b5060005b8151811015620003d857600082828151811062000351576200035162000a58565b6020026020010151905060006001600160a01b0316816001600160a01b0316036200037d5750620003cf565b6200038a600282620004e3565b15620003cd576040516001600160a01b03821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b60010162000330565b505050565b604051636eb1769f60e11b81523060048201526001600160a01b038381166024830152600091839186169063dd62ed3e90604401602060405180830381865afa1580156200042f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000455919062000a6e565b62000461919062000a9e565b604080516001600160a01b038616602482015260448082018490528251808303909101815260649091019091526020810180516001600160e01b0390811663095ea7b360e01b17909152919250620004bd91869190620004fa16565b50505050565b6000620004da836001600160a01b038416620005cb565b90505b92915050565b6000620004da836001600160a01b038416620006cf565b6040805180820190915260208082527f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65649082015260009062000549906001600160a01b03851690849062000721565b805190915015620003d857808060200190518101906200056a919062000ab4565b620003d85760405162461bcd60e51b815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e6044820152691bdd081cdd58d8d9595960b21b606482015260840162000187565b60008181526001830160205260408120548015620006c4576000620005f260018362000adf565b8554909150600090620006089060019062000adf565b9050808214620006745760008660000182815481106200062c576200062c62000a58565b906000526020600020015490508087600001848154811062000652576200065262000a58565b6000918252602080832090910192909255918252600188019052604090208390555b855486908062000688576200068862000af5565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050620004dd565b6000915050620004dd565b60008181526001830160205260408120546200071857508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155620004dd565b506000620004dd565b60606200073284846000856200073a565b949350505050565b6060824710156200079d5760405162461bcd60e51b815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f6044820152651c8818d85b1b60d21b606482015260840162000187565b600080866001600160a01b03168587604051620007bb919062000b31565b60006040518083038185875af1925050503d8060008114620007fa576040519150601f19603f3d011682016040523d82523d6000602084013e620007ff565b606091505b50909250905062000813878383876200081e565b979650505050505050565b60608315620008925782516000036200088a576001600160a01b0385163b6200088a5760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000604482015260640162000187565b508162000732565b620007328383815115620008a95781518083602001fd5b8060405162461bcd60e51b815260040162000187919062000b4f565b6001600160a01b0381168114620008db57600080fd5b50565b805160ff81168114620008f057600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b8051620008f081620008c5565b600080600080600060a086880312156200093157600080fd5b85516200093e81620008c5565b945060206200094f878201620008de565b60408801519095506001600160401b03808211156200096d57600080fd5b818901915089601f8301126200098257600080fd5b815181811115620009975762000997620008f5565b8060051b604051601f19603f83011681018181108582111715620009bf57620009bf620008f5565b60405291825284820192508381018501918c831115620009de57600080fd5b938501935b8285101562000a0757620009f7856200090b565b84529385019392850192620009e3565b80985050505050505062000a1e606087016200090b565b915062000a2e608087016200090b565b90509295509295909350565b60006020828403121562000a4d57600080fd5b620004da82620008de565b634e487b7160e01b600052603260045260246000fd5b60006020828403121562000a8157600080fd5b5051919050565b634e487b7160e01b600052601160045260246000fd5b80820180821115620004dd57620004dd62000a88565b60006020828403121562000ac757600080fd5b8151801515811462000ad857600080fd5b9392505050565b81810381811115620004dd57620004dd62000a88565b634e487b7160e01b600052603160045260246000fd5b60005b8381101562000b2857818101518382015260200162000b0e565b50506000910152565b6000825162000b4581846020870162000b0b565b9190910192915050565b602081526000825180602084015262000b7081604085016020870162000b0b565b601f01601f19169190910160400192915050565b60805160a05160c05160e05161402e62000c356000396000818161054f01528181611d8201526127d30152600081816105290152818161189f015261206e0152600081816102e001528181610ba901528181611a4801528181611b0201528181611b3601528181611b6901528181611bce01528181611c270152611cc90152600081816102470152818161029c01528181610708015281816121f10152818161276901526129be015261402e6000f3fe608060405234801561001057600080fd5b50600436106101cf5760003560e01c80639a4575b911610104578063c0d78655116100a2578063dc0bd97111610071578063dc0bd97114610527578063e0351e131461054d578063e8a1da1714610573578063f2fde38b1461058657600080fd5b8063c0d78655146104d9578063c4bffe2b146104ec578063c75eea9c14610501578063cf7401f31461051457600080fd5b8063acfecf91116100de578063acfecf9114610426578063af58d59f14610439578063b0f479a1146104a8578063b7946580146104c657600080fd5b80639a4575b9146103d1578063a42a7b8b146103f1578063a7cd63b71461041157600080fd5b806354c8a4f31161017157806379ba50971161014b57806379ba5097146103855780637d54534e1461038d5780638926f54f146103a05780638da5cb5b146103b357600080fd5b806354c8a4f31461033f57806362ddd3c4146103545780636d3d1a581461036757600080fd5b8063240028e8116101ad578063240028e81461028c57806324f65ee7146102d9578063390775371461030a5780634c5ef0ed1461032c57600080fd5b806301ffc9a7146101d4578063181f5a77146101fc57806321df0da714610245575b600080fd5b6101e76101e236600461317e565b610599565b60405190151581526020015b60405180910390f35b6102386040518060400160405280601b81526020017f4275726e46726f6d4d696e74546f6b656e506f6f6c20312e352e31000000000081525081565b6040516101f39190613224565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101f3565b6101e761029a366004613259565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff90811691161490565b60405160ff7f00000000000000000000000000000000000000000000000000000000000000001681526020016101f3565b61031d610318366004613276565b61067e565b604051905181526020016101f3565b6101e761033a3660046132cf565b61084d565b61035261034d36600461339e565b610897565b005b6103526103623660046132cf565b610912565b60095473ffffffffffffffffffffffffffffffffffffffff16610267565b6103526109af565b61035261039b366004613259565b610a7d565b6101e76103ae36600461340a565b610afe565b60015473ffffffffffffffffffffffffffffffffffffffff16610267565b6103e46103df366004613425565b610b15565b6040516101f39190613460565b6104046103ff36600461340a565b610bee565b6040516101f391906134b7565b610419610d59565b6040516101f39190613539565b6103526104343660046132cf565b610d6a565b61044c61044736600461340a565b610e82565b6040516101f3919081516fffffffffffffffffffffffffffffffff908116825260208084015163ffffffff1690830152604080840151151590830152606080840151821690830152608092830151169181019190915260a00190565b60045473ffffffffffffffffffffffffffffffffffffffff16610267565b6102386104d436600461340a565b610f57565b6103526104e7366004613259565b611007565b6104f46110e2565b6040516101f39190613593565b61044c61050f36600461340a565b61119a565b61035261052236600461371b565b61126c565b7f0000000000000000000000000000000000000000000000000000000000000000610267565b7f00000000000000000000000000000000000000000000000000000000000000006101e7565b61035261058136600461339e565b6112f0565b610352610594366004613259565b611802565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167faff2afbf00000000000000000000000000000000000000000000000000000000148061062c57507fffffffff0000000000000000000000000000000000000000000000000000000082167f0e64dd2900000000000000000000000000000000000000000000000000000000145b8061067857507fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a700000000000000000000000000000000000000000000000000000000145b92915050565b60408051602081019091526000815261069682611816565b60006106ef60608401356106ea6106b060c0870187613760565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611a3a92505050565b611afe565b905073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000166340c10f1961073d6060860160408701613259565b6040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b16815273ffffffffffffffffffffffffffffffffffffffff909116600482015260248101849052604401600060405180830381600087803b1580156107aa57600080fd5b505af11580156107be573d6000803e3d6000fd5b506107d3925050506060840160408501613259565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f9d228d69b5fdb8d273a2336f8fb8612d039631024ea9bf09c424a9503aa078f08360405161083191815260200190565b60405180910390a3604080516020810190915290815292915050565b600061088f83836040516108629291906137c5565b604080519182900390912067ffffffffffffffff8716600090815260076020529190912060050190611d12565b949350505050565b61089f611d2d565b61090c84848080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525050604080516020808802828101820190935287825290935087925086918291850190849080828437600092019190915250611d8092505050565b50505050565b61091a611d2d565b61092383610afe565b61096a576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff841660048201526024015b60405180910390fd5b6109aa8383838080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611f3692505050565b505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610a00576040517f02b543c600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000008082163390811790935560008054909116815560405173ffffffffffffffffffffffffffffffffffffffff909216929183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610a85611d2d565b600980547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040519081527f44676b5284b809a22248eba0da87391d79098be38bb03154be88a58bf4d091749060200160405180910390a150565b6000610678600567ffffffffffffffff8416611d12565b6040805180820190915260608082526020820152610b3282612030565b610b3f82606001356121bc565b6040516060830135815233907f696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df79060200160405180910390a26040518060400160405280610b998460200160208101906104d4919061340a565b8152602001610be66040805160ff7f000000000000000000000000000000000000000000000000000000000000000016602082015260609101604051602081830303815290604052905090565b905292915050565b67ffffffffffffffff8116600090815260076020526040812060609190610c179060050161225e565b90506000815167ffffffffffffffff811115610c3557610c356135d5565b604051908082528060200260200182016040528015610c6857816020015b6060815260200190600190039081610c535790505b50905060005b8251811015610d515760086000848381518110610c8d57610c8d6137d5565b602002602001015181526020019081526020016000208054610cae90613804565b80601f0160208091040260200160405190810160405280929190818152602001828054610cda90613804565b8015610d275780601f10610cfc57610100808354040283529160200191610d27565b820191906000526020600020905b815481529060010190602001808311610d0a57829003601f168201915b5050505050828281518110610d3e57610d3e6137d5565b6020908102919091010152600101610c6e565b509392505050565b6060610d65600261225e565b905090565b610d72611d2d565b610d7b83610afe565b610dbd576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84166004820152602401610961565b610dfd8282604051610dd09291906137c5565b604080519182900390912067ffffffffffffffff861660009081526007602052919091206005019061226b565b610e39578282826040517f74f23c7c000000000000000000000000000000000000000000000000000000008152600401610961939291906138a0565b8267ffffffffffffffff167f52d00ee4d9bd51b40168f2afc5848837288ce258784ad914278791464b3f4d768383604051610e759291906138c4565b60405180910390a2505050565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915267ffffffffffffffff8216600090815260076020908152604091829020825160a08101845260028201546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff16151594820194909452600390910154808416606083015291909104909116608082015261067890612277565b67ffffffffffffffff81166000908152600760205260409020600401805460609190610f8290613804565b80601f0160208091040260200160405190810160405280929190818152602001828054610fae90613804565b8015610ffb5780601f10610fd057610100808354040283529160200191610ffb565b820191906000526020600020905b815481529060010190602001808311610fde57829003601f168201915b50505050509050919050565b61100f611d2d565b73ffffffffffffffffffffffffffffffffffffffff811661105c576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6004805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527f02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f1684910160405180910390a15050565b606060006110f0600561225e565b90506000815167ffffffffffffffff81111561110e5761110e6135d5565b604051908082528060200260200182016040528015611137578160200160208202803683370190505b50905060005b825181101561119357828181518110611158576111586137d5565b6020026020010151828281518110611172576111726137d5565b67ffffffffffffffff9092166020928302919091019091015260010161113d565b5092915050565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915267ffffffffffffffff8216600090815260076020908152604091829020825160a08101845281546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff16151594820194909452600190910154808416606083015291909104909116608082015261067890612277565b60095473ffffffffffffffffffffffffffffffffffffffff1633148015906112ac575060015473ffffffffffffffffffffffffffffffffffffffff163314155b156112e5576040517f8e4a23d6000000000000000000000000000000000000000000000000000000008152336004820152602401610961565b6109aa838383612329565b6112f8611d2d565b60005b838110156114e5576000858583818110611317576113176137d5565b905060200201602081019061132c919061340a565b9050611343600567ffffffffffffffff831661226b565b611385576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff82166004820152602401610961565b67ffffffffffffffff811660009081526007602052604081206113aa9060050161225e565b905060005b81518110156114165761140d8282815181106113cd576113cd6137d5565b6020026020010151600760008667ffffffffffffffff1667ffffffffffffffff16815260200190815260200160002060050161226b90919063ffffffff16565b506001016113af565b5067ffffffffffffffff8216600090815260076020526040812080547fffffffffffffffffffffff0000000000000000000000000000000000000000009081168255600182018390556002820180549091169055600381018290559061147f6004830182613111565b6005820160008181611491828261314b565b505060405167ffffffffffffffff871681527f5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d859916945060200192506114d3915050565b60405180910390a150506001016112fb565b5060005b818110156117fb576000838383818110611505576115056137d5565b905060200281019061151791906138d8565b611520906139a4565b905061153181606001516000612413565b61154081608001516000612413565b80604001515160000361157f576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80516115979060059067ffffffffffffffff16612550565b6115dc5780516040517f1d5ad3c500000000000000000000000000000000000000000000000000000000815267ffffffffffffffff9091166004820152602401610961565b805167ffffffffffffffff16600090815260076020908152604091829020825160a08082018552606080870180518601516fffffffffffffffffffffffffffffffff90811680865263ffffffff42168689018190528351511515878b0181905284518a0151841686890181905294518b0151841660809889018190528954740100000000000000000000000000000000000000009283027fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff7001000000000000000000000000000000008087027fffffffffffffffffffffffff000000000000000000000000000000000000000094851690981788178216929092178d5592810290971760018c01558c519889018d52898e0180518d01518716808b528a8e019590955280515115158a8f018190528151909d01518716988a01899052518d0151909516979098018790526002890180549a90910299909316171790941695909517909255909202909117600382015590820151600482019061175f9082613b1b565b5060005b8260200151518110156117a35761179b83600001518460200151838151811061178e5761178e6137d5565b6020026020010151611f36565b600101611763565b507f8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c282600001518360400151846060015185608001516040516117e99493929190613c35565b60405180910390a150506001016114e9565b5050505050565b61180a611d2d565b6118138161255c565b50565b61182961029a60a0830160808401613259565b6118885761183d60a0820160808301613259565b6040517f961c9a4f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401610961565b73ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016632cbc26bb6118d4604084016020850161340a565b60405160e083901b7fffffffff0000000000000000000000000000000000000000000000000000000016815260809190911b77ffffffffffffffff00000000000000000000000000000000166004820152602401602060405180830381865afa158015611945573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906119699190613cce565b156119a0576040517f53ad11d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6119b86119b3604083016020840161340a565b612620565b6119d86119cb604083016020840161340a565b61033a60a0840184613760565b611a1d576119e960a0820182613760565b6040517f24eb47e50000000000000000000000000000000000000000000000000000000081526004016109619291906138c4565b611813611a30604083016020840161340a565b8260600135612746565b60008151600003611a6c57507f0000000000000000000000000000000000000000000000000000000000000000919050565b8151602014611aa957816040517f953576f70000000000000000000000000000000000000000000000000000000081526004016109619190613224565b600082806020019051810190611abf9190613ceb565b905060ff81111561067857826040517f953576f70000000000000000000000000000000000000000000000000000000081526004016109619190613224565b60007f000000000000000000000000000000000000000000000000000000000000000060ff168260ff1603611b34575081610678565b7f000000000000000000000000000000000000000000000000000000000000000060ff168260ff161115611c1f576000611b8e7f000000000000000000000000000000000000000000000000000000000000000084613d33565b9050604d8160ff161115611c02576040517fa9cb113d00000000000000000000000000000000000000000000000000000000815260ff80851660048301527f000000000000000000000000000000000000000000000000000000000000000016602482015260448101859052606401610961565b611c0d81600a613e6c565b611c179085613e7b565b915050610678565b6000611c4b837f0000000000000000000000000000000000000000000000000000000000000000613d33565b9050604d8160ff161180611c925750611c6581600a613e6c565b611c8f907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff613e7b565b84115b15611cfd576040517fa9cb113d00000000000000000000000000000000000000000000000000000000815260ff80851660048301527f000000000000000000000000000000000000000000000000000000000000000016602482015260448101859052606401610961565b611d0881600a613e6c565b61088f9085613eb6565b600081815260018301602052604081205415155b9392505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314611d7e576040517f2b5c74de00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b7f0000000000000000000000000000000000000000000000000000000000000000611dd7576040517f35f4a7b300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b8251811015611e6d576000838281518110611df757611df76137d5565b60200260200101519050611e1581600261278d90919063ffffffff16565b15611e645760405173ffffffffffffffffffffffffffffffffffffffff821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b50600101611dda565b5060005b81518110156109aa576000828281518110611e8e57611e8e6137d5565b60200260200101519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603611ed25750611f2e565b611edd6002826127af565b15611f2c5760405173ffffffffffffffffffffffffffffffffffffffff821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b600101611e71565b8051600003611f71576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805160208083019190912067ffffffffffffffff8416600090815260079092526040909120611fa39060050182612550565b611fdd5782826040517f393b8ad2000000000000000000000000000000000000000000000000000000008152600401610961929190613ecd565b6000818152600860205260409020611ff58382613b1b565b508267ffffffffffffffff167f7d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea83604051610e759190613224565b61204361029a60a0830160808401613259565b6120575761183d60a0820160808301613259565b73ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016632cbc26bb6120a3604084016020850161340a565b60405160e083901b7fffffffff0000000000000000000000000000000000000000000000000000000016815260809190911b77ffffffffffffffff00000000000000000000000000000000166004820152602401602060405180830381865afa158015612114573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906121389190613cce565b1561216f576040517f53ad11d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6121876121826060830160408401613259565b6127d1565b61219f61219a604083016020840161340a565b612850565b6118136121b2604083016020840161340a565b826060013561299e565b6040517f79cc6790000000000000000000000000000000000000000000000000000000008152306004820152602481018290527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906379cc679090604401600060405180830381600087803b15801561224a57600080fd5b505af11580156117fb573d6000803e3d6000fd5b60606000611d26836129e2565b6000611d268383612a3d565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915261230582606001516fffffffffffffffffffffffffffffffff1683600001516fffffffffffffffffffffffffffffffff16846020015163ffffffff16426122e99190613ef0565b85608001516fffffffffffffffffffffffffffffffff16612b30565b6fffffffffffffffffffffffffffffffff1682525063ffffffff4216602082015290565b61233283610afe565b612374576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84166004820152602401610961565b61237f826000612413565b67ffffffffffffffff831660009081526007602052604090206123a29083612b58565b6123ad816000612413565b67ffffffffffffffff831660009081526007602052604090206123d39060020182612b58565b7f0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b83838360405161240693929190613f03565b60405180910390a1505050565b8151156124de5781602001516fffffffffffffffffffffffffffffffff1682604001516fffffffffffffffffffffffffffffffff16101580612469575060408201516fffffffffffffffffffffffffffffffff16155b156124a257816040517f8020d1240000000000000000000000000000000000000000000000000000000081526004016109619190613f86565b80156124da576040517f433fc33d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050565b60408201516fffffffffffffffffffffffffffffffff16151580612517575060208201516fffffffffffffffffffffffffffffffff1615155b156124da57816040517fd68af9cc0000000000000000000000000000000000000000000000000000000081526004016109619190613f86565b6000611d268383612cfa565b3373ffffffffffffffffffffffffffffffffffffffff8216036125ab576040517fdad89dca00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b61262981610afe565b61266b576040517fa9902c7e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff82166004820152602401610961565b600480546040517f83826b2b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84169281019290925233602483015273ffffffffffffffffffffffffffffffffffffffff16906383826b2b90604401602060405180830381865afa1580156126ea573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061270e9190613cce565b611813576040517f728fe07b000000000000000000000000000000000000000000000000000000008152336004820152602401610961565b67ffffffffffffffff821660009081526007602052604090206124da90600201827f0000000000000000000000000000000000000000000000000000000000000000612d49565b6000611d268373ffffffffffffffffffffffffffffffffffffffff8416612a3d565b6000611d268373ffffffffffffffffffffffffffffffffffffffff8416612cfa565b7f000000000000000000000000000000000000000000000000000000000000000015611813576128026002826130cc565b611813576040517fd0d2597600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610961565b61285981610afe565b61289b576040517fa9902c7e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff82166004820152602401610961565b600480546040517fa8d87a3b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84169281019290925273ffffffffffffffffffffffffffffffffffffffff169063a8d87a3b90602401602060405180830381865afa158015612914573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906129389190613fc2565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614611813576040517f728fe07b000000000000000000000000000000000000000000000000000000008152336004820152602401610961565b67ffffffffffffffff821660009081526007602052604090206124da90827f0000000000000000000000000000000000000000000000000000000000000000612d49565b606081600001805480602002602001604051908101604052809291908181526020018280548015610ffb57602002820191906000526020600020905b815481526020019060010190808311612a1e5750505050509050919050565b60008181526001830160205260408120548015612b26576000612a61600183613ef0565b8554909150600090612a7590600190613ef0565b9050808214612ada576000866000018281548110612a9557612a956137d5565b9060005260206000200154905080876000018481548110612ab857612ab86137d5565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080612aeb57612aeb613fdf565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610678565b6000915050610678565b6000612b4f85612b408486613eb6565b612b4a908761400e565b6130fb565b95945050505050565b8154600090612b8190700100000000000000000000000000000000900463ffffffff1642613ef0565b90508015612c235760018301548354612bc9916fffffffffffffffffffffffffffffffff80821692811691859170010000000000000000000000000000000090910416612b30565b83546fffffffffffffffffffffffffffffffff919091167fffffffffffffffffffffffff0000000000000000000000000000000000000000909116177001000000000000000000000000000000004263ffffffff16021783555b60208201518354612c49916fffffffffffffffffffffffffffffffff90811691166130fb565b83548351151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00ffffffff000000000000000000000000000000009091166fffffffffffffffffffffffffffffffff92831617178455602083015160408085015183167001000000000000000000000000000000000291909216176001850155517f9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c1990612406908490613f86565b6000818152600183016020526040812054612d4157508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610678565b506000610678565b825474010000000000000000000000000000000000000000900460ff161580612d70575081155b15612d7a57505050565b825460018401546fffffffffffffffffffffffffffffffff80831692911690600090612dc090700100000000000000000000000000000000900463ffffffff1642613ef0565b90508015612e805781831115612e02576040517f9725942a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001860154612e3c9083908590849070010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff16612b30565b86547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff167001000000000000000000000000000000004263ffffffff160217875592505b84821015612f375773ffffffffffffffffffffffffffffffffffffffff8416612edf576040517ff94ebcd10000000000000000000000000000000000000000000000000000000081526004810183905260248101869052604401610961565b6040517f1a76572a000000000000000000000000000000000000000000000000000000008152600481018390526024810186905273ffffffffffffffffffffffffffffffffffffffff85166044820152606401610961565b8483101561304a5760018681015470010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff16906000908290612f7b9082613ef0565b612f85878a613ef0565b612f8f919061400e565b612f999190613e7b565b905073ffffffffffffffffffffffffffffffffffffffff8616612ff2576040517f15279c080000000000000000000000000000000000000000000000000000000081526004810182905260248101869052604401610961565b6040517fd0c8d23a000000000000000000000000000000000000000000000000000000008152600481018290526024810186905273ffffffffffffffffffffffffffffffffffffffff87166044820152606401610961565b6130548584613ef0565b86547fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff82161787556040518681529093507f1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a9060200160405180910390a1505050505050565b73ffffffffffffffffffffffffffffffffffffffff811660009081526001830160205260408120541515611d26565b600081831061310a5781611d26565b5090919050565b50805461311d90613804565b6000825580601f1061312d575050565b601f0160209004906000526020600020908101906118139190613165565b508054600082559060005260206000209081019061181391905b5b8082111561317a5760008155600101613166565b5090565b60006020828403121561319057600080fd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114611d2657600080fd5b6000815180845260005b818110156131e6576020818501810151868301820152016131ca565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081526000611d2660208301846131c0565b73ffffffffffffffffffffffffffffffffffffffff8116811461181357600080fd5b60006020828403121561326b57600080fd5b8135611d2681613237565b60006020828403121561328857600080fd5b813567ffffffffffffffff81111561329f57600080fd5b82016101008185031215611d2657600080fd5b803567ffffffffffffffff811681146132ca57600080fd5b919050565b6000806000604084860312156132e457600080fd5b6132ed846132b2565b9250602084013567ffffffffffffffff8082111561330a57600080fd5b818601915086601f83011261331e57600080fd5b81358181111561332d57600080fd5b87602082850101111561333f57600080fd5b6020830194508093505050509250925092565b60008083601f84011261336457600080fd5b50813567ffffffffffffffff81111561337c57600080fd5b6020830191508360208260051b850101111561339757600080fd5b9250929050565b600080600080604085870312156133b457600080fd5b843567ffffffffffffffff808211156133cc57600080fd5b6133d888838901613352565b909650945060208701359150808211156133f157600080fd5b506133fe87828801613352565b95989497509550505050565b60006020828403121561341c57600080fd5b611d26826132b2565b60006020828403121561343757600080fd5b813567ffffffffffffffff81111561344e57600080fd5b820160a08185031215611d2657600080fd5b60208152600082516040602084015261347c60608401826131c0565b905060208401517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0848303016040850152612b4f82826131c0565b600060208083016020845280855180835260408601915060408160051b87010192506020870160005b8281101561352c577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc088860301845261351a8583516131c0565b945092850192908501906001016134e0565b5092979650505050505050565b6020808252825182820181905260009190848201906040850190845b8181101561358757835173ffffffffffffffffffffffffffffffffffffffff1683529284019291840191600101613555565b50909695505050505050565b6020808252825182820181905260009190848201906040850190845b8181101561358757835167ffffffffffffffff16835292840192918401916001016135af565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160a0810167ffffffffffffffff81118282101715613627576136276135d5565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715613674576136746135d5565b604052919050565b801515811461181357600080fd5b80356fffffffffffffffffffffffffffffffff811681146132ca57600080fd5b6000606082840312156136bc57600080fd5b6040516060810181811067ffffffffffffffff821117156136df576136df6135d5565b60405290508082356136f08161367c565b81526136fe6020840161368a565b602082015261370f6040840161368a565b60408201525092915050565b600080600060e0848603121561373057600080fd5b613739846132b2565b925061374885602086016136aa565b915061375785608086016136aa565b90509250925092565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261379557600080fd5b83018035915067ffffffffffffffff8211156137b057600080fd5b60200191503681900382131561339757600080fd5b8183823760009101908152919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600181811c9082168061381857607f821691505b602082108103613851577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b67ffffffffffffffff84168152604060208201526000612b4f604083018486613857565b60208152600061088f602083018486613857565b600082357ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffee183360301811261390c57600080fd5b9190910192915050565b600082601f83011261392757600080fd5b813567ffffffffffffffff811115613941576139416135d5565b61397260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8401160161362d565b81815284602083860101111561398757600080fd5b816020850160208301376000918101602001919091529392505050565b600061012082360312156139b757600080fd5b6139bf613604565b6139c8836132b2565b815260208084013567ffffffffffffffff808211156139e657600080fd5b9085019036601f8301126139f957600080fd5b813581811115613a0b57613a0b6135d5565b8060051b613a1a85820161362d565b9182528381018501918581019036841115613a3457600080fd5b86860192505b83831015613a7057823585811115613a525760008081fd5b613a603689838a0101613916565b8352509186019190860190613a3a565b8087890152505050506040860135925080831115613a8d57600080fd5b5050613a9b36828601613916565b604083015250613aae36606085016136aa565b6060820152613ac03660c085016136aa565b608082015292915050565b601f8211156109aa576000816000526020600020601f850160051c81016020861015613af45750805b601f850160051c820191505b81811015613b1357828155600101613b00565b505050505050565b815167ffffffffffffffff811115613b3557613b356135d5565b613b4981613b438454613804565b84613acb565b602080601f831160018114613b9c5760008415613b665750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555613b13565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015613be957888601518255948401946001909101908401613bca565b5085821015613c2557878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b600061010067ffffffffffffffff87168352806020840152613c59818401876131c0565b8551151560408581019190915260208701516fffffffffffffffffffffffffffffffff9081166060870152908701511660808501529150613c979050565b8251151560a083015260208301516fffffffffffffffffffffffffffffffff90811660c084015260408401511660e0830152612b4f565b600060208284031215613ce057600080fd5b8151611d268161367c565b600060208284031215613cfd57600080fd5b5051919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60ff828116828216039081111561067857610678613d04565b600181815b80851115613da557817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04821115613d8b57613d8b613d04565b80851615613d9857918102915b93841c9390800290613d51565b509250929050565b600082613dbc57506001610678565b81613dc957506000610678565b8160018114613ddf5760028114613de957613e05565b6001915050610678565b60ff841115613dfa57613dfa613d04565b50506001821b610678565b5060208310610133831016604e8410600b8410161715613e28575081810a610678565b613e328383613d4c565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04821115613e6457613e64613d04565b029392505050565b6000611d2660ff841683613dad565b600082613eb1577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b808202811582820484141761067857610678613d04565b67ffffffffffffffff8316815260406020820152600061088f60408301846131c0565b8181038181111561067857610678613d04565b67ffffffffffffffff8416815260e08101613f4f60208301858051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b82511515608083015260208301516fffffffffffffffffffffffffffffffff90811660a084015260408401511660c083015261088f565b6060810161067882848051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b600060208284031215613fd457600080fd5b8151611d2681613237565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b8082018082111561067857610678613d0456fea164736f6c6343000818000a",
}

var BurnFromMintTokenPoolABI = BurnFromMintTokenPoolMetaData.ABI

var BurnFromMintTokenPoolBin = BurnFromMintTokenPoolMetaData.Bin

func DeployBurnFromMintTokenPool(auth *bind.TransactOpts, backend bind.ContractBackend, token common.Address, localTokenDecimals uint8, allowlist []common.Address, rmnProxy common.Address, router common.Address) (common.Address, *generated.Transaction, *BurnFromMintTokenPool, error) {
	parsed, err := BurnFromMintTokenPoolMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}
	if generated.IsZKSync(backend) {
		address, ethTx, contractBind, _ := generated.DeployContract(auth, parsed, common.FromHex(BurnFromMintTokenPoolZKBin), backend, token, localTokenDecimals, allowlist, rmnProxy, router)
		contractReturn := &BurnFromMintTokenPool{address: address, abi: *parsed, BurnFromMintTokenPoolCaller: BurnFromMintTokenPoolCaller{contract: contractBind}, BurnFromMintTokenPoolTransactor: BurnFromMintTokenPoolTransactor{contract: contractBind}, BurnFromMintTokenPoolFilterer: BurnFromMintTokenPoolFilterer{contract: contractBind}}
		return address, ethTx, contractReturn, err
	}
	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BurnFromMintTokenPoolBin), backend, token, localTokenDecimals, allowlist, rmnProxy, router)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, &generated.Transaction{Transaction: tx, HashZks: tx.Hash()}, &BurnFromMintTokenPool{address: address, abi: *parsed, BurnFromMintTokenPoolCaller: BurnFromMintTokenPoolCaller{contract: contract}, BurnFromMintTokenPoolTransactor: BurnFromMintTokenPoolTransactor{contract: contract}, BurnFromMintTokenPoolFilterer: BurnFromMintTokenPoolFilterer{contract: contract}}, nil
}

type BurnFromMintTokenPool struct {
	address common.Address
	abi     abi.ABI
	BurnFromMintTokenPoolCaller
	BurnFromMintTokenPoolTransactor
	BurnFromMintTokenPoolFilterer
}

type BurnFromMintTokenPoolCaller struct {
	contract *bind.BoundContract
}

type BurnFromMintTokenPoolTransactor struct {
	contract *bind.BoundContract
}

type BurnFromMintTokenPoolFilterer struct {
	contract *bind.BoundContract
}

type BurnFromMintTokenPoolSession struct {
	Contract     *BurnFromMintTokenPool
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type BurnFromMintTokenPoolCallerSession struct {
	Contract *BurnFromMintTokenPoolCaller
	CallOpts bind.CallOpts
}

type BurnFromMintTokenPoolTransactorSession struct {
	Contract     *BurnFromMintTokenPoolTransactor
	TransactOpts bind.TransactOpts
}

type BurnFromMintTokenPoolRaw struct {
	Contract *BurnFromMintTokenPool
}

type BurnFromMintTokenPoolCallerRaw struct {
	Contract *BurnFromMintTokenPoolCaller
}

type BurnFromMintTokenPoolTransactorRaw struct {
	Contract *BurnFromMintTokenPoolTransactor
}

func NewBurnFromMintTokenPool(address common.Address, backend bind.ContractBackend) (*BurnFromMintTokenPool, error) {
	abi, err := abi.JSON(strings.NewReader(BurnFromMintTokenPoolABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindBurnFromMintTokenPool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BurnFromMintTokenPool{address: address, abi: abi, BurnFromMintTokenPoolCaller: BurnFromMintTokenPoolCaller{contract: contract}, BurnFromMintTokenPoolTransactor: BurnFromMintTokenPoolTransactor{contract: contract}, BurnFromMintTokenPoolFilterer: BurnFromMintTokenPoolFilterer{contract: contract}}, nil
}

func NewBurnFromMintTokenPoolCaller(address common.Address, caller bind.ContractCaller) (*BurnFromMintTokenPoolCaller, error) {
	contract, err := bindBurnFromMintTokenPool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BurnFromMintTokenPoolCaller{contract: contract}, nil
}

func NewBurnFromMintTokenPoolTransactor(address common.Address, transactor bind.ContractTransactor) (*BurnFromMintTokenPoolTransactor, error) {
	contract, err := bindBurnFromMintTokenPool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BurnFromMintTokenPoolTransactor{contract: contract}, nil
}

func NewBurnFromMintTokenPoolFilterer(address common.Address, filterer bind.ContractFilterer) (*BurnFromMintTokenPoolFilterer, error) {
	contract, err := bindBurnFromMintTokenPool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BurnFromMintTokenPoolFilterer{contract: contract}, nil
}

func bindBurnFromMintTokenPool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BurnFromMintTokenPoolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BurnFromMintTokenPool.Contract.BurnFromMintTokenPoolCaller.contract.Call(opts, result, method, params...)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnFromMintTokenPool.Contract.BurnFromMintTokenPoolTransactor.contract.Transfer(opts)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BurnFromMintTokenPool.Contract.BurnFromMintTokenPoolTransactor.contract.Transact(opts, method, params...)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BurnFromMintTokenPool.Contract.contract.Call(opts, result, method, params...)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnFromMintTokenPool.Contract.contract.Transfer(opts)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BurnFromMintTokenPool.Contract.contract.Transact(opts, method, params...)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCaller) GetAllowList(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _BurnFromMintTokenPool.contract.Call(opts, &out, "getAllowList")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolSession) GetAllowList() ([]common.Address, error) {
	return _BurnFromMintTokenPool.Contract.GetAllowList(&_BurnFromMintTokenPool.CallOpts)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCallerSession) GetAllowList() ([]common.Address, error) {
	return _BurnFromMintTokenPool.Contract.GetAllowList(&_BurnFromMintTokenPool.CallOpts)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCaller) GetAllowListEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _BurnFromMintTokenPool.contract.Call(opts, &out, "getAllowListEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolSession) GetAllowListEnabled() (bool, error) {
	return _BurnFromMintTokenPool.Contract.GetAllowListEnabled(&_BurnFromMintTokenPool.CallOpts)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCallerSession) GetAllowListEnabled() (bool, error) {
	return _BurnFromMintTokenPool.Contract.GetAllowListEnabled(&_BurnFromMintTokenPool.CallOpts)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCaller) GetCurrentInboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _BurnFromMintTokenPool.contract.Call(opts, &out, "getCurrentInboundRateLimiterState", remoteChainSelector)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolSession) GetCurrentInboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnFromMintTokenPool.Contract.GetCurrentInboundRateLimiterState(&_BurnFromMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCallerSession) GetCurrentInboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnFromMintTokenPool.Contract.GetCurrentInboundRateLimiterState(&_BurnFromMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCaller) GetCurrentOutboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _BurnFromMintTokenPool.contract.Call(opts, &out, "getCurrentOutboundRateLimiterState", remoteChainSelector)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolSession) GetCurrentOutboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnFromMintTokenPool.Contract.GetCurrentOutboundRateLimiterState(&_BurnFromMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCallerSession) GetCurrentOutboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _BurnFromMintTokenPool.Contract.GetCurrentOutboundRateLimiterState(&_BurnFromMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCaller) GetRateLimitAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnFromMintTokenPool.contract.Call(opts, &out, "getRateLimitAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolSession) GetRateLimitAdmin() (common.Address, error) {
	return _BurnFromMintTokenPool.Contract.GetRateLimitAdmin(&_BurnFromMintTokenPool.CallOpts)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCallerSession) GetRateLimitAdmin() (common.Address, error) {
	return _BurnFromMintTokenPool.Contract.GetRateLimitAdmin(&_BurnFromMintTokenPool.CallOpts)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCaller) GetRemotePools(opts *bind.CallOpts, remoteChainSelector uint64) ([][]byte, error) {
	var out []interface{}
	err := _BurnFromMintTokenPool.contract.Call(opts, &out, "getRemotePools", remoteChainSelector)

	if err != nil {
		return *new([][]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][]byte)).(*[][]byte)

	return out0, err

}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolSession) GetRemotePools(remoteChainSelector uint64) ([][]byte, error) {
	return _BurnFromMintTokenPool.Contract.GetRemotePools(&_BurnFromMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCallerSession) GetRemotePools(remoteChainSelector uint64) ([][]byte, error) {
	return _BurnFromMintTokenPool.Contract.GetRemotePools(&_BurnFromMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCaller) GetRemoteToken(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error) {
	var out []interface{}
	err := _BurnFromMintTokenPool.contract.Call(opts, &out, "getRemoteToken", remoteChainSelector)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolSession) GetRemoteToken(remoteChainSelector uint64) ([]byte, error) {
	return _BurnFromMintTokenPool.Contract.GetRemoteToken(&_BurnFromMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCallerSession) GetRemoteToken(remoteChainSelector uint64) ([]byte, error) {
	return _BurnFromMintTokenPool.Contract.GetRemoteToken(&_BurnFromMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCaller) GetRmnProxy(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnFromMintTokenPool.contract.Call(opts, &out, "getRmnProxy")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolSession) GetRmnProxy() (common.Address, error) {
	return _BurnFromMintTokenPool.Contract.GetRmnProxy(&_BurnFromMintTokenPool.CallOpts)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCallerSession) GetRmnProxy() (common.Address, error) {
	return _BurnFromMintTokenPool.Contract.GetRmnProxy(&_BurnFromMintTokenPool.CallOpts)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCaller) GetRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnFromMintTokenPool.contract.Call(opts, &out, "getRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolSession) GetRouter() (common.Address, error) {
	return _BurnFromMintTokenPool.Contract.GetRouter(&_BurnFromMintTokenPool.CallOpts)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCallerSession) GetRouter() (common.Address, error) {
	return _BurnFromMintTokenPool.Contract.GetRouter(&_BurnFromMintTokenPool.CallOpts)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCaller) GetSupportedChains(opts *bind.CallOpts) ([]uint64, error) {
	var out []interface{}
	err := _BurnFromMintTokenPool.contract.Call(opts, &out, "getSupportedChains")

	if err != nil {
		return *new([]uint64), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint64)).(*[]uint64)

	return out0, err

}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolSession) GetSupportedChains() ([]uint64, error) {
	return _BurnFromMintTokenPool.Contract.GetSupportedChains(&_BurnFromMintTokenPool.CallOpts)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCallerSession) GetSupportedChains() ([]uint64, error) {
	return _BurnFromMintTokenPool.Contract.GetSupportedChains(&_BurnFromMintTokenPool.CallOpts)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCaller) GetToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnFromMintTokenPool.contract.Call(opts, &out, "getToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolSession) GetToken() (common.Address, error) {
	return _BurnFromMintTokenPool.Contract.GetToken(&_BurnFromMintTokenPool.CallOpts)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCallerSession) GetToken() (common.Address, error) {
	return _BurnFromMintTokenPool.Contract.GetToken(&_BurnFromMintTokenPool.CallOpts)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCaller) GetTokenDecimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _BurnFromMintTokenPool.contract.Call(opts, &out, "getTokenDecimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolSession) GetTokenDecimals() (uint8, error) {
	return _BurnFromMintTokenPool.Contract.GetTokenDecimals(&_BurnFromMintTokenPool.CallOpts)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCallerSession) GetTokenDecimals() (uint8, error) {
	return _BurnFromMintTokenPool.Contract.GetTokenDecimals(&_BurnFromMintTokenPool.CallOpts)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCaller) IsRemotePool(opts *bind.CallOpts, remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	var out []interface{}
	err := _BurnFromMintTokenPool.contract.Call(opts, &out, "isRemotePool", remoteChainSelector, remotePoolAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolSession) IsRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	return _BurnFromMintTokenPool.Contract.IsRemotePool(&_BurnFromMintTokenPool.CallOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCallerSession) IsRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	return _BurnFromMintTokenPool.Contract.IsRemotePool(&_BurnFromMintTokenPool.CallOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCaller) IsSupportedChain(opts *bind.CallOpts, remoteChainSelector uint64) (bool, error) {
	var out []interface{}
	err := _BurnFromMintTokenPool.contract.Call(opts, &out, "isSupportedChain", remoteChainSelector)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolSession) IsSupportedChain(remoteChainSelector uint64) (bool, error) {
	return _BurnFromMintTokenPool.Contract.IsSupportedChain(&_BurnFromMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCallerSession) IsSupportedChain(remoteChainSelector uint64) (bool, error) {
	return _BurnFromMintTokenPool.Contract.IsSupportedChain(&_BurnFromMintTokenPool.CallOpts, remoteChainSelector)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCaller) IsSupportedToken(opts *bind.CallOpts, token common.Address) (bool, error) {
	var out []interface{}
	err := _BurnFromMintTokenPool.contract.Call(opts, &out, "isSupportedToken", token)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolSession) IsSupportedToken(token common.Address) (bool, error) {
	return _BurnFromMintTokenPool.Contract.IsSupportedToken(&_BurnFromMintTokenPool.CallOpts, token)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCallerSession) IsSupportedToken(token common.Address) (bool, error) {
	return _BurnFromMintTokenPool.Contract.IsSupportedToken(&_BurnFromMintTokenPool.CallOpts, token)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnFromMintTokenPool.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolSession) Owner() (common.Address, error) {
	return _BurnFromMintTokenPool.Contract.Owner(&_BurnFromMintTokenPool.CallOpts)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCallerSession) Owner() (common.Address, error) {
	return _BurnFromMintTokenPool.Contract.Owner(&_BurnFromMintTokenPool.CallOpts)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _BurnFromMintTokenPool.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BurnFromMintTokenPool.Contract.SupportsInterface(&_BurnFromMintTokenPool.CallOpts, interfaceId)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BurnFromMintTokenPool.Contract.SupportsInterface(&_BurnFromMintTokenPool.CallOpts, interfaceId)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _BurnFromMintTokenPool.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolSession) TypeAndVersion() (string, error) {
	return _BurnFromMintTokenPool.Contract.TypeAndVersion(&_BurnFromMintTokenPool.CallOpts)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolCallerSession) TypeAndVersion() (string, error) {
	return _BurnFromMintTokenPool.Contract.TypeAndVersion(&_BurnFromMintTokenPool.CallOpts)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnFromMintTokenPool.contract.Transact(opts, "acceptOwnership")
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolSession) AcceptOwnership() (*types.Transaction, error) {
	return _BurnFromMintTokenPool.Contract.AcceptOwnership(&_BurnFromMintTokenPool.TransactOpts)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _BurnFromMintTokenPool.Contract.AcceptOwnership(&_BurnFromMintTokenPool.TransactOpts)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolTransactor) AddRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnFromMintTokenPool.contract.Transact(opts, "addRemotePool", remoteChainSelector, remotePoolAddress)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolSession) AddRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnFromMintTokenPool.Contract.AddRemotePool(&_BurnFromMintTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolTransactorSession) AddRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnFromMintTokenPool.Contract.AddRemotePool(&_BurnFromMintTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolTransactor) ApplyAllowListUpdates(opts *bind.TransactOpts, removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _BurnFromMintTokenPool.contract.Transact(opts, "applyAllowListUpdates", removes, adds)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _BurnFromMintTokenPool.Contract.ApplyAllowListUpdates(&_BurnFromMintTokenPool.TransactOpts, removes, adds)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolTransactorSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _BurnFromMintTokenPool.Contract.ApplyAllowListUpdates(&_BurnFromMintTokenPool.TransactOpts, removes, adds)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolTransactor) ApplyChainUpdates(opts *bind.TransactOpts, remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _BurnFromMintTokenPool.contract.Transact(opts, "applyChainUpdates", remoteChainSelectorsToRemove, chainsToAdd)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolSession) ApplyChainUpdates(remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _BurnFromMintTokenPool.Contract.ApplyChainUpdates(&_BurnFromMintTokenPool.TransactOpts, remoteChainSelectorsToRemove, chainsToAdd)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolTransactorSession) ApplyChainUpdates(remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _BurnFromMintTokenPool.Contract.ApplyChainUpdates(&_BurnFromMintTokenPool.TransactOpts, remoteChainSelectorsToRemove, chainsToAdd)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolTransactor) LockOrBurn(opts *bind.TransactOpts, lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _BurnFromMintTokenPool.contract.Transact(opts, "lockOrBurn", lockOrBurnIn)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolSession) LockOrBurn(lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _BurnFromMintTokenPool.Contract.LockOrBurn(&_BurnFromMintTokenPool.TransactOpts, lockOrBurnIn)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolTransactorSession) LockOrBurn(lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _BurnFromMintTokenPool.Contract.LockOrBurn(&_BurnFromMintTokenPool.TransactOpts, lockOrBurnIn)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolTransactor) ReleaseOrMint(opts *bind.TransactOpts, releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _BurnFromMintTokenPool.contract.Transact(opts, "releaseOrMint", releaseOrMintIn)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolSession) ReleaseOrMint(releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _BurnFromMintTokenPool.Contract.ReleaseOrMint(&_BurnFromMintTokenPool.TransactOpts, releaseOrMintIn)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolTransactorSession) ReleaseOrMint(releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _BurnFromMintTokenPool.Contract.ReleaseOrMint(&_BurnFromMintTokenPool.TransactOpts, releaseOrMintIn)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolTransactor) RemoveRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnFromMintTokenPool.contract.Transact(opts, "removeRemotePool", remoteChainSelector, remotePoolAddress)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolSession) RemoveRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnFromMintTokenPool.Contract.RemoveRemotePool(&_BurnFromMintTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolTransactorSession) RemoveRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _BurnFromMintTokenPool.Contract.RemoveRemotePool(&_BurnFromMintTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolTransactor) SetChainRateLimiterConfig(opts *bind.TransactOpts, remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _BurnFromMintTokenPool.contract.Transact(opts, "setChainRateLimiterConfig", remoteChainSelector, outboundConfig, inboundConfig)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolSession) SetChainRateLimiterConfig(remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _BurnFromMintTokenPool.Contract.SetChainRateLimiterConfig(&_BurnFromMintTokenPool.TransactOpts, remoteChainSelector, outboundConfig, inboundConfig)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolTransactorSession) SetChainRateLimiterConfig(remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _BurnFromMintTokenPool.Contract.SetChainRateLimiterConfig(&_BurnFromMintTokenPool.TransactOpts, remoteChainSelector, outboundConfig, inboundConfig)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolTransactor) SetRateLimitAdmin(opts *bind.TransactOpts, rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _BurnFromMintTokenPool.contract.Transact(opts, "setRateLimitAdmin", rateLimitAdmin)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolSession) SetRateLimitAdmin(rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _BurnFromMintTokenPool.Contract.SetRateLimitAdmin(&_BurnFromMintTokenPool.TransactOpts, rateLimitAdmin)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolTransactorSession) SetRateLimitAdmin(rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _BurnFromMintTokenPool.Contract.SetRateLimitAdmin(&_BurnFromMintTokenPool.TransactOpts, rateLimitAdmin)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolTransactor) SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error) {
	return _BurnFromMintTokenPool.contract.Transact(opts, "setRouter", newRouter)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _BurnFromMintTokenPool.Contract.SetRouter(&_BurnFromMintTokenPool.TransactOpts, newRouter)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolTransactorSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _BurnFromMintTokenPool.Contract.SetRouter(&_BurnFromMintTokenPool.TransactOpts, newRouter)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _BurnFromMintTokenPool.contract.Transact(opts, "transferOwnership", to)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _BurnFromMintTokenPool.Contract.TransferOwnership(&_BurnFromMintTokenPool.TransactOpts, to)
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _BurnFromMintTokenPool.Contract.TransferOwnership(&_BurnFromMintTokenPool.TransactOpts, to)
}

type BurnFromMintTokenPoolAllowListAddIterator struct {
	Event *BurnFromMintTokenPoolAllowListAdd

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnFromMintTokenPoolAllowListAddIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnFromMintTokenPoolAllowListAdd)
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
		it.Event = new(BurnFromMintTokenPoolAllowListAdd)
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

func (it *BurnFromMintTokenPoolAllowListAddIterator) Error() error {
	return it.fail
}

func (it *BurnFromMintTokenPoolAllowListAddIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnFromMintTokenPoolAllowListAdd struct {
	Sender common.Address
	Raw    types.Log
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) FilterAllowListAdd(opts *bind.FilterOpts) (*BurnFromMintTokenPoolAllowListAddIterator, error) {

	logs, sub, err := _BurnFromMintTokenPool.contract.FilterLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return &BurnFromMintTokenPoolAllowListAddIterator{contract: _BurnFromMintTokenPool.contract, event: "AllowListAdd", logs: logs, sub: sub}, nil
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *BurnFromMintTokenPoolAllowListAdd) (event.Subscription, error) {

	logs, sub, err := _BurnFromMintTokenPool.contract.WatchLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnFromMintTokenPoolAllowListAdd)
				if err := _BurnFromMintTokenPool.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
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

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) ParseAllowListAdd(log types.Log) (*BurnFromMintTokenPoolAllowListAdd, error) {
	event := new(BurnFromMintTokenPoolAllowListAdd)
	if err := _BurnFromMintTokenPool.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnFromMintTokenPoolAllowListRemoveIterator struct {
	Event *BurnFromMintTokenPoolAllowListRemove

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnFromMintTokenPoolAllowListRemoveIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnFromMintTokenPoolAllowListRemove)
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
		it.Event = new(BurnFromMintTokenPoolAllowListRemove)
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

func (it *BurnFromMintTokenPoolAllowListRemoveIterator) Error() error {
	return it.fail
}

func (it *BurnFromMintTokenPoolAllowListRemoveIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnFromMintTokenPoolAllowListRemove struct {
	Sender common.Address
	Raw    types.Log
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) FilterAllowListRemove(opts *bind.FilterOpts) (*BurnFromMintTokenPoolAllowListRemoveIterator, error) {

	logs, sub, err := _BurnFromMintTokenPool.contract.FilterLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return &BurnFromMintTokenPoolAllowListRemoveIterator{contract: _BurnFromMintTokenPool.contract, event: "AllowListRemove", logs: logs, sub: sub}, nil
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *BurnFromMintTokenPoolAllowListRemove) (event.Subscription, error) {

	logs, sub, err := _BurnFromMintTokenPool.contract.WatchLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnFromMintTokenPoolAllowListRemove)
				if err := _BurnFromMintTokenPool.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
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

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) ParseAllowListRemove(log types.Log) (*BurnFromMintTokenPoolAllowListRemove, error) {
	event := new(BurnFromMintTokenPoolAllowListRemove)
	if err := _BurnFromMintTokenPool.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnFromMintTokenPoolBurnedIterator struct {
	Event *BurnFromMintTokenPoolBurned

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnFromMintTokenPoolBurnedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnFromMintTokenPoolBurned)
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
		it.Event = new(BurnFromMintTokenPoolBurned)
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

func (it *BurnFromMintTokenPoolBurnedIterator) Error() error {
	return it.fail
}

func (it *BurnFromMintTokenPoolBurnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnFromMintTokenPoolBurned struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) FilterBurned(opts *bind.FilterOpts, sender []common.Address) (*BurnFromMintTokenPoolBurnedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BurnFromMintTokenPool.contract.FilterLogs(opts, "Burned", senderRule)
	if err != nil {
		return nil, err
	}
	return &BurnFromMintTokenPoolBurnedIterator{contract: _BurnFromMintTokenPool.contract, event: "Burned", logs: logs, sub: sub}, nil
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) WatchBurned(opts *bind.WatchOpts, sink chan<- *BurnFromMintTokenPoolBurned, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BurnFromMintTokenPool.contract.WatchLogs(opts, "Burned", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnFromMintTokenPoolBurned)
				if err := _BurnFromMintTokenPool.contract.UnpackLog(event, "Burned", log); err != nil {
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

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) ParseBurned(log types.Log) (*BurnFromMintTokenPoolBurned, error) {
	event := new(BurnFromMintTokenPoolBurned)
	if err := _BurnFromMintTokenPool.contract.UnpackLog(event, "Burned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnFromMintTokenPoolChainAddedIterator struct {
	Event *BurnFromMintTokenPoolChainAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnFromMintTokenPoolChainAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnFromMintTokenPoolChainAdded)
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
		it.Event = new(BurnFromMintTokenPoolChainAdded)
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

func (it *BurnFromMintTokenPoolChainAddedIterator) Error() error {
	return it.fail
}

func (it *BurnFromMintTokenPoolChainAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnFromMintTokenPoolChainAdded struct {
	RemoteChainSelector       uint64
	RemoteToken               []byte
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
	Raw                       types.Log
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) FilterChainAdded(opts *bind.FilterOpts) (*BurnFromMintTokenPoolChainAddedIterator, error) {

	logs, sub, err := _BurnFromMintTokenPool.contract.FilterLogs(opts, "ChainAdded")
	if err != nil {
		return nil, err
	}
	return &BurnFromMintTokenPoolChainAddedIterator{contract: _BurnFromMintTokenPool.contract, event: "ChainAdded", logs: logs, sub: sub}, nil
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) WatchChainAdded(opts *bind.WatchOpts, sink chan<- *BurnFromMintTokenPoolChainAdded) (event.Subscription, error) {

	logs, sub, err := _BurnFromMintTokenPool.contract.WatchLogs(opts, "ChainAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnFromMintTokenPoolChainAdded)
				if err := _BurnFromMintTokenPool.contract.UnpackLog(event, "ChainAdded", log); err != nil {
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

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) ParseChainAdded(log types.Log) (*BurnFromMintTokenPoolChainAdded, error) {
	event := new(BurnFromMintTokenPoolChainAdded)
	if err := _BurnFromMintTokenPool.contract.UnpackLog(event, "ChainAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnFromMintTokenPoolChainConfiguredIterator struct {
	Event *BurnFromMintTokenPoolChainConfigured

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnFromMintTokenPoolChainConfiguredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnFromMintTokenPoolChainConfigured)
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
		it.Event = new(BurnFromMintTokenPoolChainConfigured)
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

func (it *BurnFromMintTokenPoolChainConfiguredIterator) Error() error {
	return it.fail
}

func (it *BurnFromMintTokenPoolChainConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnFromMintTokenPoolChainConfigured struct {
	RemoteChainSelector       uint64
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
	Raw                       types.Log
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) FilterChainConfigured(opts *bind.FilterOpts) (*BurnFromMintTokenPoolChainConfiguredIterator, error) {

	logs, sub, err := _BurnFromMintTokenPool.contract.FilterLogs(opts, "ChainConfigured")
	if err != nil {
		return nil, err
	}
	return &BurnFromMintTokenPoolChainConfiguredIterator{contract: _BurnFromMintTokenPool.contract, event: "ChainConfigured", logs: logs, sub: sub}, nil
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) WatchChainConfigured(opts *bind.WatchOpts, sink chan<- *BurnFromMintTokenPoolChainConfigured) (event.Subscription, error) {

	logs, sub, err := _BurnFromMintTokenPool.contract.WatchLogs(opts, "ChainConfigured")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnFromMintTokenPoolChainConfigured)
				if err := _BurnFromMintTokenPool.contract.UnpackLog(event, "ChainConfigured", log); err != nil {
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

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) ParseChainConfigured(log types.Log) (*BurnFromMintTokenPoolChainConfigured, error) {
	event := new(BurnFromMintTokenPoolChainConfigured)
	if err := _BurnFromMintTokenPool.contract.UnpackLog(event, "ChainConfigured", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnFromMintTokenPoolChainRemovedIterator struct {
	Event *BurnFromMintTokenPoolChainRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnFromMintTokenPoolChainRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnFromMintTokenPoolChainRemoved)
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
		it.Event = new(BurnFromMintTokenPoolChainRemoved)
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

func (it *BurnFromMintTokenPoolChainRemovedIterator) Error() error {
	return it.fail
}

func (it *BurnFromMintTokenPoolChainRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnFromMintTokenPoolChainRemoved struct {
	RemoteChainSelector uint64
	Raw                 types.Log
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) FilterChainRemoved(opts *bind.FilterOpts) (*BurnFromMintTokenPoolChainRemovedIterator, error) {

	logs, sub, err := _BurnFromMintTokenPool.contract.FilterLogs(opts, "ChainRemoved")
	if err != nil {
		return nil, err
	}
	return &BurnFromMintTokenPoolChainRemovedIterator{contract: _BurnFromMintTokenPool.contract, event: "ChainRemoved", logs: logs, sub: sub}, nil
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *BurnFromMintTokenPoolChainRemoved) (event.Subscription, error) {

	logs, sub, err := _BurnFromMintTokenPool.contract.WatchLogs(opts, "ChainRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnFromMintTokenPoolChainRemoved)
				if err := _BurnFromMintTokenPool.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
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

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) ParseChainRemoved(log types.Log) (*BurnFromMintTokenPoolChainRemoved, error) {
	event := new(BurnFromMintTokenPoolChainRemoved)
	if err := _BurnFromMintTokenPool.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnFromMintTokenPoolConfigChangedIterator struct {
	Event *BurnFromMintTokenPoolConfigChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnFromMintTokenPoolConfigChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnFromMintTokenPoolConfigChanged)
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
		it.Event = new(BurnFromMintTokenPoolConfigChanged)
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

func (it *BurnFromMintTokenPoolConfigChangedIterator) Error() error {
	return it.fail
}

func (it *BurnFromMintTokenPoolConfigChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnFromMintTokenPoolConfigChanged struct {
	Config RateLimiterConfig
	Raw    types.Log
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) FilterConfigChanged(opts *bind.FilterOpts) (*BurnFromMintTokenPoolConfigChangedIterator, error) {

	logs, sub, err := _BurnFromMintTokenPool.contract.FilterLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return &BurnFromMintTokenPoolConfigChangedIterator{contract: _BurnFromMintTokenPool.contract, event: "ConfigChanged", logs: logs, sub: sub}, nil
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *BurnFromMintTokenPoolConfigChanged) (event.Subscription, error) {

	logs, sub, err := _BurnFromMintTokenPool.contract.WatchLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnFromMintTokenPoolConfigChanged)
				if err := _BurnFromMintTokenPool.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
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

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) ParseConfigChanged(log types.Log) (*BurnFromMintTokenPoolConfigChanged, error) {
	event := new(BurnFromMintTokenPoolConfigChanged)
	if err := _BurnFromMintTokenPool.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnFromMintTokenPoolLockedIterator struct {
	Event *BurnFromMintTokenPoolLocked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnFromMintTokenPoolLockedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnFromMintTokenPoolLocked)
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
		it.Event = new(BurnFromMintTokenPoolLocked)
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

func (it *BurnFromMintTokenPoolLockedIterator) Error() error {
	return it.fail
}

func (it *BurnFromMintTokenPoolLockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnFromMintTokenPoolLocked struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) FilterLocked(opts *bind.FilterOpts, sender []common.Address) (*BurnFromMintTokenPoolLockedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BurnFromMintTokenPool.contract.FilterLogs(opts, "Locked", senderRule)
	if err != nil {
		return nil, err
	}
	return &BurnFromMintTokenPoolLockedIterator{contract: _BurnFromMintTokenPool.contract, event: "Locked", logs: logs, sub: sub}, nil
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) WatchLocked(opts *bind.WatchOpts, sink chan<- *BurnFromMintTokenPoolLocked, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BurnFromMintTokenPool.contract.WatchLogs(opts, "Locked", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnFromMintTokenPoolLocked)
				if err := _BurnFromMintTokenPool.contract.UnpackLog(event, "Locked", log); err != nil {
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

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) ParseLocked(log types.Log) (*BurnFromMintTokenPoolLocked, error) {
	event := new(BurnFromMintTokenPoolLocked)
	if err := _BurnFromMintTokenPool.contract.UnpackLog(event, "Locked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnFromMintTokenPoolMintedIterator struct {
	Event *BurnFromMintTokenPoolMinted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnFromMintTokenPoolMintedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnFromMintTokenPoolMinted)
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
		it.Event = new(BurnFromMintTokenPoolMinted)
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

func (it *BurnFromMintTokenPoolMintedIterator) Error() error {
	return it.fail
}

func (it *BurnFromMintTokenPoolMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnFromMintTokenPoolMinted struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) FilterMinted(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*BurnFromMintTokenPoolMintedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BurnFromMintTokenPool.contract.FilterLogs(opts, "Minted", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &BurnFromMintTokenPoolMintedIterator{contract: _BurnFromMintTokenPool.contract, event: "Minted", logs: logs, sub: sub}, nil
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) WatchMinted(opts *bind.WatchOpts, sink chan<- *BurnFromMintTokenPoolMinted, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BurnFromMintTokenPool.contract.WatchLogs(opts, "Minted", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnFromMintTokenPoolMinted)
				if err := _BurnFromMintTokenPool.contract.UnpackLog(event, "Minted", log); err != nil {
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

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) ParseMinted(log types.Log) (*BurnFromMintTokenPoolMinted, error) {
	event := new(BurnFromMintTokenPoolMinted)
	if err := _BurnFromMintTokenPool.contract.UnpackLog(event, "Minted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnFromMintTokenPoolOwnershipTransferRequestedIterator struct {
	Event *BurnFromMintTokenPoolOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnFromMintTokenPoolOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnFromMintTokenPoolOwnershipTransferRequested)
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
		it.Event = new(BurnFromMintTokenPoolOwnershipTransferRequested)
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

func (it *BurnFromMintTokenPoolOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *BurnFromMintTokenPoolOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnFromMintTokenPoolOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnFromMintTokenPoolOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnFromMintTokenPool.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BurnFromMintTokenPoolOwnershipTransferRequestedIterator{contract: _BurnFromMintTokenPool.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *BurnFromMintTokenPoolOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnFromMintTokenPool.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnFromMintTokenPoolOwnershipTransferRequested)
				if err := _BurnFromMintTokenPool.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) ParseOwnershipTransferRequested(log types.Log) (*BurnFromMintTokenPoolOwnershipTransferRequested, error) {
	event := new(BurnFromMintTokenPoolOwnershipTransferRequested)
	if err := _BurnFromMintTokenPool.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnFromMintTokenPoolOwnershipTransferredIterator struct {
	Event *BurnFromMintTokenPoolOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnFromMintTokenPoolOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnFromMintTokenPoolOwnershipTransferred)
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
		it.Event = new(BurnFromMintTokenPoolOwnershipTransferred)
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

func (it *BurnFromMintTokenPoolOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *BurnFromMintTokenPoolOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnFromMintTokenPoolOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnFromMintTokenPoolOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnFromMintTokenPool.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BurnFromMintTokenPoolOwnershipTransferredIterator{contract: _BurnFromMintTokenPool.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BurnFromMintTokenPoolOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BurnFromMintTokenPool.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnFromMintTokenPoolOwnershipTransferred)
				if err := _BurnFromMintTokenPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) ParseOwnershipTransferred(log types.Log) (*BurnFromMintTokenPoolOwnershipTransferred, error) {
	event := new(BurnFromMintTokenPoolOwnershipTransferred)
	if err := _BurnFromMintTokenPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnFromMintTokenPoolRateLimitAdminSetIterator struct {
	Event *BurnFromMintTokenPoolRateLimitAdminSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnFromMintTokenPoolRateLimitAdminSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnFromMintTokenPoolRateLimitAdminSet)
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
		it.Event = new(BurnFromMintTokenPoolRateLimitAdminSet)
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

func (it *BurnFromMintTokenPoolRateLimitAdminSetIterator) Error() error {
	return it.fail
}

func (it *BurnFromMintTokenPoolRateLimitAdminSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnFromMintTokenPoolRateLimitAdminSet struct {
	RateLimitAdmin common.Address
	Raw            types.Log
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) FilterRateLimitAdminSet(opts *bind.FilterOpts) (*BurnFromMintTokenPoolRateLimitAdminSetIterator, error) {

	logs, sub, err := _BurnFromMintTokenPool.contract.FilterLogs(opts, "RateLimitAdminSet")
	if err != nil {
		return nil, err
	}
	return &BurnFromMintTokenPoolRateLimitAdminSetIterator{contract: _BurnFromMintTokenPool.contract, event: "RateLimitAdminSet", logs: logs, sub: sub}, nil
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) WatchRateLimitAdminSet(opts *bind.WatchOpts, sink chan<- *BurnFromMintTokenPoolRateLimitAdminSet) (event.Subscription, error) {

	logs, sub, err := _BurnFromMintTokenPool.contract.WatchLogs(opts, "RateLimitAdminSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnFromMintTokenPoolRateLimitAdminSet)
				if err := _BurnFromMintTokenPool.contract.UnpackLog(event, "RateLimitAdminSet", log); err != nil {
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

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) ParseRateLimitAdminSet(log types.Log) (*BurnFromMintTokenPoolRateLimitAdminSet, error) {
	event := new(BurnFromMintTokenPoolRateLimitAdminSet)
	if err := _BurnFromMintTokenPool.contract.UnpackLog(event, "RateLimitAdminSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnFromMintTokenPoolReleasedIterator struct {
	Event *BurnFromMintTokenPoolReleased

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnFromMintTokenPoolReleasedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnFromMintTokenPoolReleased)
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
		it.Event = new(BurnFromMintTokenPoolReleased)
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

func (it *BurnFromMintTokenPoolReleasedIterator) Error() error {
	return it.fail
}

func (it *BurnFromMintTokenPoolReleasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnFromMintTokenPoolReleased struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) FilterReleased(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*BurnFromMintTokenPoolReleasedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BurnFromMintTokenPool.contract.FilterLogs(opts, "Released", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &BurnFromMintTokenPoolReleasedIterator{contract: _BurnFromMintTokenPool.contract, event: "Released", logs: logs, sub: sub}, nil
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) WatchReleased(opts *bind.WatchOpts, sink chan<- *BurnFromMintTokenPoolReleased, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BurnFromMintTokenPool.contract.WatchLogs(opts, "Released", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnFromMintTokenPoolReleased)
				if err := _BurnFromMintTokenPool.contract.UnpackLog(event, "Released", log); err != nil {
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

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) ParseReleased(log types.Log) (*BurnFromMintTokenPoolReleased, error) {
	event := new(BurnFromMintTokenPoolReleased)
	if err := _BurnFromMintTokenPool.contract.UnpackLog(event, "Released", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnFromMintTokenPoolRemotePoolAddedIterator struct {
	Event *BurnFromMintTokenPoolRemotePoolAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnFromMintTokenPoolRemotePoolAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnFromMintTokenPoolRemotePoolAdded)
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
		it.Event = new(BurnFromMintTokenPoolRemotePoolAdded)
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

func (it *BurnFromMintTokenPoolRemotePoolAddedIterator) Error() error {
	return it.fail
}

func (it *BurnFromMintTokenPoolRemotePoolAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnFromMintTokenPoolRemotePoolAdded struct {
	RemoteChainSelector uint64
	RemotePoolAddress   []byte
	Raw                 types.Log
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) FilterRemotePoolAdded(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnFromMintTokenPoolRemotePoolAddedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnFromMintTokenPool.contract.FilterLogs(opts, "RemotePoolAdded", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &BurnFromMintTokenPoolRemotePoolAddedIterator{contract: _BurnFromMintTokenPool.contract, event: "RemotePoolAdded", logs: logs, sub: sub}, nil
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) WatchRemotePoolAdded(opts *bind.WatchOpts, sink chan<- *BurnFromMintTokenPoolRemotePoolAdded, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnFromMintTokenPool.contract.WatchLogs(opts, "RemotePoolAdded", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnFromMintTokenPoolRemotePoolAdded)
				if err := _BurnFromMintTokenPool.contract.UnpackLog(event, "RemotePoolAdded", log); err != nil {
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

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) ParseRemotePoolAdded(log types.Log) (*BurnFromMintTokenPoolRemotePoolAdded, error) {
	event := new(BurnFromMintTokenPoolRemotePoolAdded)
	if err := _BurnFromMintTokenPool.contract.UnpackLog(event, "RemotePoolAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnFromMintTokenPoolRemotePoolRemovedIterator struct {
	Event *BurnFromMintTokenPoolRemotePoolRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnFromMintTokenPoolRemotePoolRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnFromMintTokenPoolRemotePoolRemoved)
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
		it.Event = new(BurnFromMintTokenPoolRemotePoolRemoved)
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

func (it *BurnFromMintTokenPoolRemotePoolRemovedIterator) Error() error {
	return it.fail
}

func (it *BurnFromMintTokenPoolRemotePoolRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnFromMintTokenPoolRemotePoolRemoved struct {
	RemoteChainSelector uint64
	RemotePoolAddress   []byte
	Raw                 types.Log
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) FilterRemotePoolRemoved(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnFromMintTokenPoolRemotePoolRemovedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnFromMintTokenPool.contract.FilterLogs(opts, "RemotePoolRemoved", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &BurnFromMintTokenPoolRemotePoolRemovedIterator{contract: _BurnFromMintTokenPool.contract, event: "RemotePoolRemoved", logs: logs, sub: sub}, nil
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) WatchRemotePoolRemoved(opts *bind.WatchOpts, sink chan<- *BurnFromMintTokenPoolRemotePoolRemoved, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _BurnFromMintTokenPool.contract.WatchLogs(opts, "RemotePoolRemoved", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnFromMintTokenPoolRemotePoolRemoved)
				if err := _BurnFromMintTokenPool.contract.UnpackLog(event, "RemotePoolRemoved", log); err != nil {
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

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) ParseRemotePoolRemoved(log types.Log) (*BurnFromMintTokenPoolRemotePoolRemoved, error) {
	event := new(BurnFromMintTokenPoolRemotePoolRemoved)
	if err := _BurnFromMintTokenPool.contract.UnpackLog(event, "RemotePoolRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnFromMintTokenPoolRouterUpdatedIterator struct {
	Event *BurnFromMintTokenPoolRouterUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnFromMintTokenPoolRouterUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnFromMintTokenPoolRouterUpdated)
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
		it.Event = new(BurnFromMintTokenPoolRouterUpdated)
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

func (it *BurnFromMintTokenPoolRouterUpdatedIterator) Error() error {
	return it.fail
}

func (it *BurnFromMintTokenPoolRouterUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnFromMintTokenPoolRouterUpdated struct {
	OldRouter common.Address
	NewRouter common.Address
	Raw       types.Log
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) FilterRouterUpdated(opts *bind.FilterOpts) (*BurnFromMintTokenPoolRouterUpdatedIterator, error) {

	logs, sub, err := _BurnFromMintTokenPool.contract.FilterLogs(opts, "RouterUpdated")
	if err != nil {
		return nil, err
	}
	return &BurnFromMintTokenPoolRouterUpdatedIterator{contract: _BurnFromMintTokenPool.contract, event: "RouterUpdated", logs: logs, sub: sub}, nil
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *BurnFromMintTokenPoolRouterUpdated) (event.Subscription, error) {

	logs, sub, err := _BurnFromMintTokenPool.contract.WatchLogs(opts, "RouterUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnFromMintTokenPoolRouterUpdated)
				if err := _BurnFromMintTokenPool.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
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

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) ParseRouterUpdated(log types.Log) (*BurnFromMintTokenPoolRouterUpdated, error) {
	event := new(BurnFromMintTokenPoolRouterUpdated)
	if err := _BurnFromMintTokenPool.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BurnFromMintTokenPoolTokensConsumedIterator struct {
	Event *BurnFromMintTokenPoolTokensConsumed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BurnFromMintTokenPoolTokensConsumedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnFromMintTokenPoolTokensConsumed)
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
		it.Event = new(BurnFromMintTokenPoolTokensConsumed)
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

func (it *BurnFromMintTokenPoolTokensConsumedIterator) Error() error {
	return it.fail
}

func (it *BurnFromMintTokenPoolTokensConsumedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BurnFromMintTokenPoolTokensConsumed struct {
	Tokens *big.Int
	Raw    types.Log
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) FilterTokensConsumed(opts *bind.FilterOpts) (*BurnFromMintTokenPoolTokensConsumedIterator, error) {

	logs, sub, err := _BurnFromMintTokenPool.contract.FilterLogs(opts, "TokensConsumed")
	if err != nil {
		return nil, err
	}
	return &BurnFromMintTokenPoolTokensConsumedIterator{contract: _BurnFromMintTokenPool.contract, event: "TokensConsumed", logs: logs, sub: sub}, nil
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) WatchTokensConsumed(opts *bind.WatchOpts, sink chan<- *BurnFromMintTokenPoolTokensConsumed) (event.Subscription, error) {

	logs, sub, err := _BurnFromMintTokenPool.contract.WatchLogs(opts, "TokensConsumed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BurnFromMintTokenPoolTokensConsumed)
				if err := _BurnFromMintTokenPool.contract.UnpackLog(event, "TokensConsumed", log); err != nil {
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

func (_BurnFromMintTokenPool *BurnFromMintTokenPoolFilterer) ParseTokensConsumed(log types.Log) (*BurnFromMintTokenPoolTokensConsumed, error) {
	event := new(BurnFromMintTokenPoolTokensConsumed)
	if err := _BurnFromMintTokenPool.contract.UnpackLog(event, "TokensConsumed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPool) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _BurnFromMintTokenPool.abi.Events["AllowListAdd"].ID:
		return _BurnFromMintTokenPool.ParseAllowListAdd(log)
	case _BurnFromMintTokenPool.abi.Events["AllowListRemove"].ID:
		return _BurnFromMintTokenPool.ParseAllowListRemove(log)
	case _BurnFromMintTokenPool.abi.Events["Burned"].ID:
		return _BurnFromMintTokenPool.ParseBurned(log)
	case _BurnFromMintTokenPool.abi.Events["ChainAdded"].ID:
		return _BurnFromMintTokenPool.ParseChainAdded(log)
	case _BurnFromMintTokenPool.abi.Events["ChainConfigured"].ID:
		return _BurnFromMintTokenPool.ParseChainConfigured(log)
	case _BurnFromMintTokenPool.abi.Events["ChainRemoved"].ID:
		return _BurnFromMintTokenPool.ParseChainRemoved(log)
	case _BurnFromMintTokenPool.abi.Events["ConfigChanged"].ID:
		return _BurnFromMintTokenPool.ParseConfigChanged(log)
	case _BurnFromMintTokenPool.abi.Events["Locked"].ID:
		return _BurnFromMintTokenPool.ParseLocked(log)
	case _BurnFromMintTokenPool.abi.Events["Minted"].ID:
		return _BurnFromMintTokenPool.ParseMinted(log)
	case _BurnFromMintTokenPool.abi.Events["OwnershipTransferRequested"].ID:
		return _BurnFromMintTokenPool.ParseOwnershipTransferRequested(log)
	case _BurnFromMintTokenPool.abi.Events["OwnershipTransferred"].ID:
		return _BurnFromMintTokenPool.ParseOwnershipTransferred(log)
	case _BurnFromMintTokenPool.abi.Events["RateLimitAdminSet"].ID:
		return _BurnFromMintTokenPool.ParseRateLimitAdminSet(log)
	case _BurnFromMintTokenPool.abi.Events["Released"].ID:
		return _BurnFromMintTokenPool.ParseReleased(log)
	case _BurnFromMintTokenPool.abi.Events["RemotePoolAdded"].ID:
		return _BurnFromMintTokenPool.ParseRemotePoolAdded(log)
	case _BurnFromMintTokenPool.abi.Events["RemotePoolRemoved"].ID:
		return _BurnFromMintTokenPool.ParseRemotePoolRemoved(log)
	case _BurnFromMintTokenPool.abi.Events["RouterUpdated"].ID:
		return _BurnFromMintTokenPool.ParseRouterUpdated(log)
	case _BurnFromMintTokenPool.abi.Events["TokensConsumed"].ID:
		return _BurnFromMintTokenPool.ParseTokensConsumed(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (BurnFromMintTokenPoolAllowListAdd) Topic() common.Hash {
	return common.HexToHash("0x2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d8")
}

func (BurnFromMintTokenPoolAllowListRemove) Topic() common.Hash {
	return common.HexToHash("0x800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf7566")
}

func (BurnFromMintTokenPoolBurned) Topic() common.Hash {
	return common.HexToHash("0x696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df7")
}

func (BurnFromMintTokenPoolChainAdded) Topic() common.Hash {
	return common.HexToHash("0x8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c2")
}

func (BurnFromMintTokenPoolChainConfigured) Topic() common.Hash {
	return common.HexToHash("0x0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b")
}

func (BurnFromMintTokenPoolChainRemoved) Topic() common.Hash {
	return common.HexToHash("0x5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d859916")
}

func (BurnFromMintTokenPoolConfigChanged) Topic() common.Hash {
	return common.HexToHash("0x9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c19")
}

func (BurnFromMintTokenPoolLocked) Topic() common.Hash {
	return common.HexToHash("0x9f1ec8c880f76798e7b793325d625e9b60e4082a553c98f42b6cda368dd60008")
}

func (BurnFromMintTokenPoolMinted) Topic() common.Hash {
	return common.HexToHash("0x9d228d69b5fdb8d273a2336f8fb8612d039631024ea9bf09c424a9503aa078f0")
}

func (BurnFromMintTokenPoolOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (BurnFromMintTokenPoolOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (BurnFromMintTokenPoolRateLimitAdminSet) Topic() common.Hash {
	return common.HexToHash("0x44676b5284b809a22248eba0da87391d79098be38bb03154be88a58bf4d09174")
}

func (BurnFromMintTokenPoolReleased) Topic() common.Hash {
	return common.HexToHash("0x2d87480f50083e2b2759522a8fdda59802650a8055e609a7772cf70c07748f52")
}

func (BurnFromMintTokenPoolRemotePoolAdded) Topic() common.Hash {
	return common.HexToHash("0x7d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea")
}

func (BurnFromMintTokenPoolRemotePoolRemoved) Topic() common.Hash {
	return common.HexToHash("0x52d00ee4d9bd51b40168f2afc5848837288ce258784ad914278791464b3f4d76")
}

func (BurnFromMintTokenPoolRouterUpdated) Topic() common.Hash {
	return common.HexToHash("0x02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f1684")
}

func (BurnFromMintTokenPoolTokensConsumed) Topic() common.Hash {
	return common.HexToHash("0x1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a")
}

func (_BurnFromMintTokenPool *BurnFromMintTokenPool) Address() common.Address {
	return _BurnFromMintTokenPool.address
}

type BurnFromMintTokenPoolInterface interface {
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

	FilterAllowListAdd(opts *bind.FilterOpts) (*BurnFromMintTokenPoolAllowListAddIterator, error)

	WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *BurnFromMintTokenPoolAllowListAdd) (event.Subscription, error)

	ParseAllowListAdd(log types.Log) (*BurnFromMintTokenPoolAllowListAdd, error)

	FilterAllowListRemove(opts *bind.FilterOpts) (*BurnFromMintTokenPoolAllowListRemoveIterator, error)

	WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *BurnFromMintTokenPoolAllowListRemove) (event.Subscription, error)

	ParseAllowListRemove(log types.Log) (*BurnFromMintTokenPoolAllowListRemove, error)

	FilterBurned(opts *bind.FilterOpts, sender []common.Address) (*BurnFromMintTokenPoolBurnedIterator, error)

	WatchBurned(opts *bind.WatchOpts, sink chan<- *BurnFromMintTokenPoolBurned, sender []common.Address) (event.Subscription, error)

	ParseBurned(log types.Log) (*BurnFromMintTokenPoolBurned, error)

	FilterChainAdded(opts *bind.FilterOpts) (*BurnFromMintTokenPoolChainAddedIterator, error)

	WatchChainAdded(opts *bind.WatchOpts, sink chan<- *BurnFromMintTokenPoolChainAdded) (event.Subscription, error)

	ParseChainAdded(log types.Log) (*BurnFromMintTokenPoolChainAdded, error)

	FilterChainConfigured(opts *bind.FilterOpts) (*BurnFromMintTokenPoolChainConfiguredIterator, error)

	WatchChainConfigured(opts *bind.WatchOpts, sink chan<- *BurnFromMintTokenPoolChainConfigured) (event.Subscription, error)

	ParseChainConfigured(log types.Log) (*BurnFromMintTokenPoolChainConfigured, error)

	FilterChainRemoved(opts *bind.FilterOpts) (*BurnFromMintTokenPoolChainRemovedIterator, error)

	WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *BurnFromMintTokenPoolChainRemoved) (event.Subscription, error)

	ParseChainRemoved(log types.Log) (*BurnFromMintTokenPoolChainRemoved, error)

	FilterConfigChanged(opts *bind.FilterOpts) (*BurnFromMintTokenPoolConfigChangedIterator, error)

	WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *BurnFromMintTokenPoolConfigChanged) (event.Subscription, error)

	ParseConfigChanged(log types.Log) (*BurnFromMintTokenPoolConfigChanged, error)

	FilterLocked(opts *bind.FilterOpts, sender []common.Address) (*BurnFromMintTokenPoolLockedIterator, error)

	WatchLocked(opts *bind.WatchOpts, sink chan<- *BurnFromMintTokenPoolLocked, sender []common.Address) (event.Subscription, error)

	ParseLocked(log types.Log) (*BurnFromMintTokenPoolLocked, error)

	FilterMinted(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*BurnFromMintTokenPoolMintedIterator, error)

	WatchMinted(opts *bind.WatchOpts, sink chan<- *BurnFromMintTokenPoolMinted, sender []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseMinted(log types.Log) (*BurnFromMintTokenPoolMinted, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnFromMintTokenPoolOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *BurnFromMintTokenPoolOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*BurnFromMintTokenPoolOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BurnFromMintTokenPoolOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BurnFromMintTokenPoolOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*BurnFromMintTokenPoolOwnershipTransferred, error)

	FilterRateLimitAdminSet(opts *bind.FilterOpts) (*BurnFromMintTokenPoolRateLimitAdminSetIterator, error)

	WatchRateLimitAdminSet(opts *bind.WatchOpts, sink chan<- *BurnFromMintTokenPoolRateLimitAdminSet) (event.Subscription, error)

	ParseRateLimitAdminSet(log types.Log) (*BurnFromMintTokenPoolRateLimitAdminSet, error)

	FilterReleased(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*BurnFromMintTokenPoolReleasedIterator, error)

	WatchReleased(opts *bind.WatchOpts, sink chan<- *BurnFromMintTokenPoolReleased, sender []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseReleased(log types.Log) (*BurnFromMintTokenPoolReleased, error)

	FilterRemotePoolAdded(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnFromMintTokenPoolRemotePoolAddedIterator, error)

	WatchRemotePoolAdded(opts *bind.WatchOpts, sink chan<- *BurnFromMintTokenPoolRemotePoolAdded, remoteChainSelector []uint64) (event.Subscription, error)

	ParseRemotePoolAdded(log types.Log) (*BurnFromMintTokenPoolRemotePoolAdded, error)

	FilterRemotePoolRemoved(opts *bind.FilterOpts, remoteChainSelector []uint64) (*BurnFromMintTokenPoolRemotePoolRemovedIterator, error)

	WatchRemotePoolRemoved(opts *bind.WatchOpts, sink chan<- *BurnFromMintTokenPoolRemotePoolRemoved, remoteChainSelector []uint64) (event.Subscription, error)

	ParseRemotePoolRemoved(log types.Log) (*BurnFromMintTokenPoolRemotePoolRemoved, error)

	FilterRouterUpdated(opts *bind.FilterOpts) (*BurnFromMintTokenPoolRouterUpdatedIterator, error)

	WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *BurnFromMintTokenPoolRouterUpdated) (event.Subscription, error)

	ParseRouterUpdated(log types.Log) (*BurnFromMintTokenPoolRouterUpdated, error)

	FilterTokensConsumed(opts *bind.FilterOpts) (*BurnFromMintTokenPoolTokensConsumedIterator, error)

	WatchTokensConsumed(opts *bind.WatchOpts, sink chan<- *BurnFromMintTokenPoolTokensConsumed) (event.Subscription, error)

	ParseTokensConsumed(log types.Log) (*BurnFromMintTokenPoolTokensConsumed, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}

var BurnFromMintTokenPoolZKBin = ("0x0004000000000002001f00000000000200000060031002700000066b0030019d0000066b033001970003000000310355000200000001035500000001002001900000005d0000c13d0000008004000039000000400040043f000000040030008c000000850000413d000000000201043b000000e0022002700000068b0020009c000000870000a13d0000068c0020009c000000c80000a13d0000068d0020009c000000f70000213d000006930020009c000001810000213d000006960020009c000005230000613d000006970020009c000000850000c13d0000000002000416000000000002004b000000850000c13d0000000504000039000000000204041a000000800020043f000000000040043f000000000002004b00000a110000c13d000000a002000039000000400020043f0000002004000039000000000500001900000005065002100000003f07600039000006d40770019700000000072700190000066f0070009c000001230000213d000000400070043f00000000005204350000001f0560018f000000a004400039000000000006004b0000003b0000613d000000000131034f00000000036400190000000006040019000000001701043c0000000006760436000000000036004b000000370000c13d000000000005004b000000800100043d000000000001004b0000004d0000613d00000000010000190000000003020433000000000013004b0000102b0000a13d00000005031002100000000005430019000000a00330003900000000030304330000066f0330019700000000003504350000000101100039000000800300043d000000000031004b000000400000413d000000400100043d00000020030000390000000005310436000000000302043300000000003504350000004002100039000000000003004b00000a080000613d000000000500001900000000460404340000066f0660019700000000026204360000000105500039000000000035004b000000560000413d00000a080000013d0000010004000039000000400040043f0000000002000416000000000002004b000000850000c13d0000001f023000390000066c022001970000010002200039000000400020043f0000001f0530018f0000066d0630019800000100026000390000006f0000613d000000000701034f000000007807043c0000000004840436000000000024004b0000006b0000c13d000000000005004b0000007c0000613d000000000161034f0000000304500210000000000502043300000000054501cf000000000545022f000000000101043b0000010004400089000000000141022f00000000014101cf000000000151019f0000000000120435000000a00030008c000000850000413d000001000100043d0000066e0010009c000000850000213d000001200200043d001100000002001d000000ff0020008c000001100000a13d0000000001000019000019aa00010430000006a10020009c000000a20000213d000006ab0020009c000001290000a13d000006ac0020009c000001580000213d000006af0020009c000001fa0000613d000006b00020009c000000850000c13d0000000001000416000000000001004b000000850000c13d0000000001000412001900000001001d001800200000003d000080050100003900000044030000390000000004000415000000190440008a0000000504400210000006cb0200004119a819800000040f000000ff0110018f000000800010043f000006cc01000041000019a90001042e000006a20020009c0000013a0000a13d000006a30020009c000001630000213d000006a60020009c000002150000613d000006a70020009c000000850000c13d000000240030008c000000850000413d0000000002000416000000000002004b000000850000c13d0000000401100370000000000101043b0000066e0010009c000000850000213d0000000102000039000000000202041a0000066e022001970000000003000411000000000023004b000009e10000c13d0000000902000039000000000302041a0000067203300197000000000313019f000000000032041b000000800010043f00000000010004140000066b0010009c0000066b01008041000000c001100210000006e9011001c70000800d020000390000000103000039000006ea0400004100000a350000013d000006980020009c000001450000a13d000006990020009c0000016c0000213d0000069c0020009c000002980000613d0000069d0020009c000000850000c13d000000240030008c000000850000413d0000000002000416000000000002004b000000850000c13d0000000401100370000000000101043b001100000001001d0000066f0010009c000000850000213d19a815a50000040f0000001101000029000000000010043f0000000701000039000000200010043f0000004002000039000000000100001919a8196b0000040f001000000001001d000000400100043d001100000001001d19a814600000040f00000010050000290000000201500039000000000401041a000006d1004001980000000002000039000000010200c03900000011010000290000004003100039000000000023043500000080024002700000066b0220019700000020031000390000000000230435000006be02400197000000000021043500000003025000390000056e0000013d0000068e0020009c000001e10000213d000006910020009c000005380000613d000006920020009c000000850000c13d0000000001000416000000000001004b000000850000c13d0000000001000412001300000001001d001200600000003d000080050100003900000044030000390000000004000415000000130440008a0000000504400210000006cb0200004119a819800000040f000000000001004b0000000001000039000000010100c039000000800010043f000006cc01000041000019a90001042e000001400200043d0000066f0020009c000000850000213d0000001f04200039000000000034004b000000000500001900000670050080410000067004400197000000000004004b00000000060000190000067006004041000006700040009c000000000605c019000000000006004b000000850000c13d000001000420003900000000040404330000066f0040009c000009b50000a13d000006dc01000041000000000010043f0000004101000039000000040010043f000006bb01000041000019aa00010430000006b10020009c000003fb0000613d000006b20020009c000003410000613d000006b30020009c000000850000c13d0000000001000416000000000001004b000000850000c13d0000000001000412001d00000001001d001c00000000003d0000800501000039000000440300003900000000040004150000001d0440008a000005420000013d000006a80020009c0000040e0000613d000006a90020009c000003550000613d000006aa0020009c000000850000c13d0000000001000416000000000001004b000000850000c13d00000009010000390000033c0000013d0000069e0020009c000004d00000613d0000069f0020009c000003930000613d000006a00020009c000000850000c13d0000000001000416000000000001004b000000850000c13d0000000202000039000000000102041a000000800010043f000000000020043f0000002002000039000000000001004b000009e90000c13d000000a0010000390000000004020019000009f80000013d000006ad0020009c000002310000613d000006ae0020009c000000850000c13d0000000001000416000000000001004b000000850000c13d000000000103001919a8147d0000040f19a815150000040f0000028e0000013d000006a40020009c000002840000613d000006a50020009c000000850000c13d0000000001000416000000000001004b000000850000c13d00000001010000390000033c0000013d0000069a0020009c000003380000613d0000069b0020009c000000850000c13d000000240030008c000000850000413d0000000002000416000000000002004b000000850000c13d0000000401100370000000000101043b0000066f0010009c000000850000213d19a815ba0000040f0000002002000039000000400300043d001100000003001d000000000223043619a8146b0000040f00000011020000290000057b0000013d000006940020009c000005490000613d000006950020009c000000850000c13d000000e40030008c000000850000413d0000000002000416000000000002004b000000850000c13d0000000402100370000000000202043b001100000002001d0000066f0020009c000000850000213d000000e002000039000000400020043f0000002402100370000000000202043b000000000002004b0000000003000039000000010300c039000000000032004b000000850000c13d000000800020043f0000004402100370000000000202043b000006be0020009c000000850000213d000000a00020043f0000006402100370000000000202043b000006be0020009c000000850000213d000000c00020043f0000014002000039000000400020043f0000008402100370000000000202043b000000000002004b0000000003000039000000010300c039000000000032004b000000850000c13d000000e00020043f000000a402100370000000000202043b000006be0020009c000000850000213d000001000020043f000000c401100370000000000101043b000006be0010009c000000850000213d000001200010043f0000000901000039000000000101041a0000066e021001970000000001000411000000000021004b000001c20000613d0000000102000039000000000202041a0000066e02200197000000000021004b00000df00000c13d0000001101000029000000000010043f0000000601000039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f0000000100200190000000850000613d000000000101043b000000000101041a000000000001004b000003880000613d000000c00100043d000006be01100197000000800200043d000000000002004b00000ed30000c13d000000000001004b000001dd0000c13d000000a00100043d000006be0010019800000ed90000613d000000400200043d001100000002001d000006c00100004100000f460000013d0000068f0020009c000005840000613d000006900020009c000000850000c13d000000240030008c000000850000413d0000000002000416000000000002004b000000850000c13d0000000401100370000000000601043b0000066e0060009c000000850000213d0000000101000039000000000101041a0000066e011001970000000005000411000000000015004b000009e10000c13d000000000056004b00000a390000c13d000006b601000041000000800010043f000006b701000041000019aa00010430000000240030008c000000850000413d0000000002000416000000000002004b000000850000c13d0000000401100370000000000101043b001100000001001d0000066e0010009c000000850000213d0000000001000412001b00000001001d001a00000000003d0000800501000039000000440300003900000000040004150000001b0440008a0000000504400210000006cb0200004119a819800000040f0000066e01100197000000110010006b00000000010000390000000101006039000000800010043f000006cc01000041000019a90001042e0000000001000416000000000001004b000000850000c13d000000000100041a0000066e021001970000000006000411000000000026004b000009e50000c13d0000000102000039000000000302041a0000067204300197000000000464019f000000000042041b0000067201100197000000000010041b00000000010004140000066e053001970000066b0010009c0000066b01008041000000c001100210000006b4011001c70000800d020000390000000303000039000006ec0400004119a8199e0000040f0000000100200190000000850000613d00000a480000013d000000240030008c000000850000413d0000000002000416000000000002004b000000850000c13d0000000402100370000000000202043b001100000002001d0000066f0020009c000000850000213d000000110230006a000006820020009c000000850000213d000001040020008c000000850000413d000000a002000039000000400020043f0000001102000029000f00840020003d0000000f01100360000000800000043f000000000101043b001000000001001d0000066e0010009c000000850000213d000006cb01000041000000000010044300000000010004120000000400100443000000240000044300000000010004140000066b0010009c0000066b01008041000000c001100210000006dd011001c7000080050200003919a819a30000040f0000000100200190000013d70000613d0000000202000367000000000101043b0000066e01100197000000100010006b00000a6b0000c13d0000000f01000029000e0060001000920000000e01200360000000000101043b0000066f0010009c000000850000213d000000400300043d000006df0200004100000000002304350000008001100210000006e001100197000f00000003001d00000004023000390000000000120435000006cb010000410000000000100443000000000100041200000004001004430000004001000039000000240010044300000000010004140000066b0010009c0000066b01008041000000c001100210000006dd011001c7000080050200003919a819a30000040f0000000100200190000013d70000613d000000000201043b00000000010004140000066e02200197000000040020008c00000be30000c13d0000000103000031000000200030008c0000002004000039000000000403401900000c0d0000013d000000240030008c000000850000413d0000000002000416000000000002004b000000850000c13d0000000401100370000000000101043b0000066f0010009c000000850000213d19a817b80000040f000000000001004b0000000001000039000000010100c039000000400200043d00000000001204350000066b0020009c0000066b020080410000004001200210000006e8011001c7000019a90001042e000000440030008c000000850000413d0000000002000416000000000002004b000000850000c13d0000000402100370000000000202043b001100000002001d0000066f0020009c000000850000213d0000002402100370000000000202043b0000066f0020009c000000850000213d0000002304200039000000000034004b000000850000813d000f00040020003d0000000f01100360000000000101043b001000000001001d0000066f0010009c000000850000213d0000002402200039000d00000002001d000e00100020002d0000000e0030006b000000850000213d0000000101000039000000000101041a0000066e011001970000000002000411000000000012004b000009e10000c13d0000001101000029000000000010043f0000000601000039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f0000000100200190000000850000613d000000000101043b000000000101041a000000000001004b000003880000613d0000001101000029000000000010043f0000000701000039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f0000000100200190000000850000613d00000010020000290000001f02200039000006fd02200197000b00000002001d0000003f02200039000006fd02200197000000000101043b000c00000001001d000000400100043d0000000002210019000000000012004b000000000400003900000001040040390000066f0020009c000001230000213d0000000100400190000001230000c13d000000400020043f000000100200002900000000022104360000000e05000029000000000050007c000000850000213d0000001004000029000006fd03400198000e001f00400193000a00000003001d00000000033200190000000f040000290000002004400039000f00000004001d0000000204400367000002ff0000613d000000000504034f0000000006020019000000005705043c0000000006760436000000000036004b000002fb0000c13d0000000e0000006b0000030d0000613d0000000a044003600000000e050000290000000305500210000000000603043300000000065601cf000000000656022f000000000404043b0000010005500089000000000454022f00000000045401cf000000000464019f0000000000430435000000100320002900000000000304350000066b0020009c0000066b02008041000000400220021000000000010104330000066b0010009c0000066b010080410000006001100210000000000121019f00000000020004140000066b0020009c0000066b02008041000000c002200210000000000112019f000006b4011001c7000080100200003919a819a30000040f0000000100200190000000850000613d0000000c020000290000000503200039000000000201043b000000000103001919a817cc0000040f000000400700043d000000000001004b00000f500000c13d000c00000007001d000000240170003900000040020000390000000000210435000006da01000041000000000017043500000004017000390000001102000029000000000021043500000044037000390000000d01000029000000100200002919a815860000040f0000000c02000029000010450000013d0000000001000416000000000001004b000000850000c13d0000000401000039000000000101041a0000066e01100197000000800010043f000006cc01000041000019a90001042e0000000001000416000000000001004b000000850000c13d000000c001000039000000400010043f0000001b01000039000000800010043f000006f701000041000000a00010043f0000002001000039000000c00010043f0000008001000039000000e00200003919a8146b0000040f000000c00110008a0000066b0010009c0000066b010080410000006001100210000006f8011001c7000019a90001042e000000440030008c000000850000413d0000000002000416000000000002004b000000850000c13d0000000402100370000000000202043b001100000002001d0000066f0020009c000000850000213d0000002402100370000000000202043b0000066f0020009c000000850000213d0000002304200039000000000034004b000000850000813d0000000404200039000000000141034f000000000101043b001000000001001d0000066f0010009c000000850000213d0000002402200039000f00000002001d0000001001200029000000000031004b000000850000213d0000000101000039000000000101041a0000066e011001970000000002000411000000000012004b000009e10000c13d0000001101000029000000000010043f0000000601000039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f0000000100200190000000850000613d000000000101043b000000000101041a000000000001004b00000b590000c13d000000400100043d000006ba0200004100000000002104350000000402100039000000110300002900000000003204350000066b0010009c0000066b010080410000004001100210000006bb011001c7000019aa00010430000000240030008c000000850000413d0000000002000416000000000002004b000000850000c13d0000000401100370000000000101043b0000066f0010009c000000850000213d000000000010043f0000000701000039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f0000000100200190000000850000613d000000000101043b0000000501100039000000000301041a000000400200043d000f00000002001d001100000003001d0000000002320436000e00000002001d000000000010043f00000000010004140000066b0010009c0000066b01008041000000c0011002100000067a011001c7000080100200003919a819a30000040f0000000100200190000000850000613d0000001105000029000000000005004b0000000e02000029000003c70000613d000000000101043b0000000e020000290000000003000019000000000401041a000000000242043600000001011000390000000103300039000000000053004b000003c10000413d0000000f0120006a0000001f01100039000006fd011001970000000f04100029000000000014004b000000000100003900000001010040390000066f0040009c000001230000213d0000000100100190000001230000c13d000000400040043f0000000f0100002900000000020104330000066f0020009c000001230000213d00000005012002100000003f03100039000006710330019700000000034300190000066f0030009c000001230000213d000000400030043f000d00000004001d0000000005240436000000000002004b000003e90000613d00000060020000390000000003000019000000000435001900000000002404350000002003300039000000000013004b000003e40000413d000c00000005001d0000000f010000290000000001010433000000000001004b00000b620000c13d000000400100043d000000200200003900000000032104360000000d0200002900000000020204330000000000230435000000400310003900000005042002100000000005340019000000000002004b00000bc60000c13d000000000215004900000a090000013d000000240030008c000000850000413d0000000002000416000000000002004b000000850000c13d0000000401100370000000000201043b000006f900200198000000850000c13d0000000101000039000006fa0020009c000005460000613d000006fb0020009c000005460000613d000006fc0020009c000000000100c019000000800010043f000006cc01000041000019a90001042e000000440030008c000000850000413d0000000002000416000000000002004b000000850000c13d0000000402100370000000000202043b0000066f0020009c000000850000213d0000002305200039000000000035004b000000850000813d0000000405200039000000000551034f000000000905043b0000066f0090009c000000850000213d0000002407200039000000050b90021000000000087b0019000000000038004b000000850000213d0000002402100370000000000202043b0000066f0020009c000000850000213d0000002305200039000000000035004b000000850000813d0000000405200039000000000551034f000000000605043b0000066f0060009c000000850000213d0000002402200039000000050a60021000000000052a0019000000000035004b000000850000213d0000000103000039000000000303041a0000066e03300197000000000c00041100000000003c004b000009e10000c13d0000003f03b000390000067103300197000006d30030009c000001230000213d0000008003300039000f00000003001d000000400030043f000000800090043f000000000009004b000004500000613d000000000371034f000000000303043b0000066e0030009c000000850000213d000000200440003900000000003404350000002007700039000000000087004b000004450000413d000000400300043d000f00000003001d0000003f03a0003900000671033001970000000f033000290000000f0030006c000000000400003900000001040040390000066f0030009c000001230000213d0000000100400190000001230000c13d000000400030043f0000000f030000290000000003630436000e00000003001d000000000006004b0000046a0000613d0000000f03000029000000000421034f000000000404043b0000066e0040009c000000850000213d000000200330003900000000004304350000002002200039000000000052004b000004610000413d000006cb010000410000000000100443000000000100041200000004001004430000006001000039000000240010044300000000010004140000066b0010009c0000066b01008041000000c001100210000006dd011001c7000080050200003919a819a30000040f0000000100200190000013d70000613d000000000101043b000000000001004b00000ade0000613d000000800100043d000000000001004b00000fc60000c13d0000000f010000290000000001010433000000000001004b00000a480000613d00000000030000190000048a0000013d00000001033000390000000f010000290000000001010433000000000013004b00000a480000813d00000005013002100000000e0110002900000000010104330000066e04100198000004850000613d000000000040043f0000000301000039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c70000801002000039001000000003001d001100000004001d19a819a30000040f000000110400002900000010030000290000000100200190000000850000613d000000000101043b000000000101041a000000000001004b000004850000c13d0000000203000039000000000103041a0000066f0010009c000001230000213d0000000102100039000000000023041b000006790110009a000000000041041b000000000103041a000d00000001001d000000000040043f0000000301000039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f00000011030000290000000100200190000000850000613d000000000101043b0000000d02000029000000000021041b000000400100043d00000000003104350000066b0010009c0000066b01008041000000400110021000000000020004140000066b0020009c0000066b02008041000000c002200210000000000112019f0000067a011001c70000800d0200003900000001030000390000067b0400004119a8199e0000040f00000010030000290000000100200190000004850000c13d000000850000013d000000240030008c000000850000413d0000000002000416000000000002004b000000850000c13d0000000402100370000000000202043b0000066f0020009c000000850000213d0000000003230049000006820030009c000000850000213d000000a40030008c000000850000413d000000c003000039000000400030043f0000006003000039000000800030043f000000a00030043f001000840020003d0000001001100360000000000101043b001100000001001d0000066e0010009c000000850000213d000006cb01000041000000000010044300000000010004120000000400100443000000240000044300000000010004140000066b0010009c0000066b01008041000000c001100210000006dd011001c7000080050200003919a819a30000040f0000000100200190000013d70000613d0000000202000367000000000101043b0000066e01100197000000110010006b00000a4a0000c13d0000001001000029000f0060001000920000000f01200360000000000101043b0000066f0010009c000000850000213d000000400300043d000006df0200004100000000002304350000008001100210000006e001100197001000000003001d00000004023000390000000000120435000006cb010000410000000000100443000000000100041200000004001004430000004001000039000000240010044300000000010004140000066b0010009c0000066b01008041000000c001100210000006dd011001c7000080050200003919a819a30000040f0000000100200190000013d70000613d000000000201043b00000000010004140000066e02200197000000040020008c00000ae10000c13d0000000103000031000000200030008c0000002004000039000000000403401900000b0b0000013d000000240030008c000000850000413d0000000002000416000000000002004b000000850000c13d0000000401100370000000000101043b0000066e0010009c000000850000213d0000000102000039000000000202041a0000066e022001970000000003000411000000000023004b000009e10000c13d000000000001004b00000a250000c13d000006ca01000041000000800010043f000006b701000041000019aa000104300000000001000416000000000001004b000000850000c13d0000000001000412001500000001001d001400400000003d000080050100003900000044030000390000000004000415000000150440008a0000000504400210000006cb0200004119a819800000040f0000066e01100197000000800010043f000006cc01000041000019a90001042e000000240030008c000000850000413d0000000002000416000000000002004b000000850000c13d0000000401100370000000000101043b001100000001001d0000066f0010009c000000850000213d19a815a50000040f0000001101000029000000000010043f0000000701000039000000200010043f0000004002000039000000000100001919a8196b0000040f001000000001001d000000400100043d001100000001001d19a814600000040f0000001005000029000000000405041a000006d1004001980000000002000039000000010200c03900000011010000290000004003100039000000000023043500000080024002700000066b0220019700000020031000390000000000230435000006be0240019700000000002104350000000102500039000000000402041a000000800210003900000080034002700000000000320435000006be034001970000006002100039000000000032043519a8176a0000040f000000400100043d001000000001001d000000110200002919a8149a0000040f000000100200002900000000012100490000066b0010009c0000066b0100804100000060011002100000066b0020009c0000066b020080410000004002200210000000000121019f000019a90001042e000000440030008c000000850000413d0000000002000416000000000002004b000000850000c13d0000000402100370000000000202043b0000066f0020009c000000850000213d0000002304200039000000000034004b000000850000813d0000000404200039000000000441034f000000000404043b000600000004001d0000066f0040009c000000850000213d000500240020003d000000060200002900000005022002100000000502200029000000000032004b000000850000213d0000002402100370000000000202043b000200000002001d0000066f0020009c000000850000213d00000002020000290000002302200039000000000032004b000000850000813d00000002020000290000000402200039000000000121034f000000000101043b000100000001001d0000066f0010009c000000850000213d0000000201000029000300240010003d000000010100002900000005011002100000000301100029000000000031004b000000850000213d0000000101000039000000000101041a0000066e011001970000000002000411000000000012004b000009e10000c13d000000060000006b00000c530000c13d000000010000006b00000a480000613d000500000000001d0000000501000029000000050110021000000003011000290000000202000367000000000112034f000000000101043b0000000003000031000000020430006a000001430440008a00000670054001970000067006100197000000000756013f000000000056004b00000000050000190000067005004041000000000041004b00000000040000190000067004008041000006700070009c000000000504c019000000000005004b000000850000c13d001000030010002d000000100130006a000f00000001001d000006820010009c000000850000213d0000000f01000029000001200010008c000000850000413d000000400100043d000900000001001d000006bc0010009c000001230000213d0000000901000029000000a001100039000000400010043f0000001001200360000000000101043b0000066f0010009c000000850000213d00000009040000290000000001140436000800000001001d00000010010000290000002001100039000000000112034f000000000101043b0000066f0010009c000000850000213d0000001001100029001100000001001d0000001f01100039000000000031004b0000000004000019000006700400804100000670011001970000067005300197000000000751013f000000000051004b00000000010000190000067001004041000006700070009c000000000104c019000000000001004b000000850000c13d0000001101200360000000000101043b0000066f0010009c000001230000213d00000005091002100000003f049000390000067104400197000000400600043d0000000004460019000e00000006001d000000000064004b000000000700003900000001070040390000066f0040009c000001230000213d0000000100700190000001230000c13d000000400040043f0000000e040000290000000000140435000000110100002900000020081000390000000009980019000000000039004b000000850000213d000000000098004b000006670000813d0000000e0a000029000006240000013d000000200aa000390000000001b7001900000000000104350000000000ca04350000002008800039000000000098004b000006670000813d000000000182034f000000000101043b0000066f0010009c000000850000213d000000110d1000290000003f01d00039000000000031004b000000000400001900000670040080410000067001100197000000000751013f000000000051004b00000000010000190000067001004041000006700070009c000000000104c019000000000001004b000000850000c13d000000200ed000390000000001e2034f000000000b01043b0000066f00b0009c000001230000213d0000001f01b00039000006fd011001970000003f01100039000006fd01100197000000400c00043d00000000011c00190000000000c1004b000000000400003900000001040040390000066f0010009c000001230000213d0000000100400190000001230000c13d0000004004d00039000000400010043f0000000007bc043600000000014b0019000000000031004b000000850000213d0000002001e00039000000000412034f000006fd01b00198000000000e170019000006590000613d000000000f04034f000000000d07001900000000f60f043c000000000d6d04360000000000ed004b000006550000c13d0000001f0db001900000061d0000613d000000000114034f0000000304d0021000000000060e043300000000064601cf000000000646022f000000000101043b0000010004400089000000000141022f00000000014101cf000000000161019f00000000001e04350000061d0000013d00000008010000290000000e04000029000000000041043500000010010000290000004006100039000000000162034f000000000101043b0000066f0010009c000000850000213d00000010071000290000001f01700039000000000031004b000000000400001900000670040080410000067001100197000000000851013f000000000051004b00000000010000190000067001004041000006700080009c000000000104c019000000000001004b000000850000c13d000000000172034f000000000401043b0000066f0040009c000001230000213d0000001f01400039000006fd011001970000003f01100039000006fd01100197000000400500043d0000000001150019000000000051004b000000000800003900000001080040390000066f0010009c000001230000213d0000000100800190000001230000c13d0000002008700039000000400010043f00000000074504360000000001840019000000000031004b000000850000213d000000000882034f000006fd0140019800000000031700190000069f0000613d000000000908034f000000000a070019000000009b09043c000000000aba043600000000003a004b0000069b0000c13d0000001f09400190000006ac0000613d000000000118034f0000000308900210000000000903043300000000098901cf000000000989022f000000000101043b0000010008800089000000000181022f00000000018101cf000000000191019f00000000001304350000000001470019000000000001043500000009010000290000004001100039000600000001001d00000000005104350000000f01000029000000600110008a000006820010009c000000850000213d000000600010008c000000850000413d000000400300043d000006bd0030009c000001230000213d0000006001300039000000400010043f0000002001600039000000000412034f000000000404043b000000000004004b0000000005000039000000010500c039000000000054004b000000850000c13d00000000044304360000002001100039000000000512034f000000000505043b000006be0050009c000000850000213d00000000005404350000002004100039000000000142034f000000000101043b000006be0010009c000000850000213d0000004005300039000000000015043500000009010000290000006001100039000700000001001d00000000003104350000000f01000029000000c00110008a000006820010009c000000850000213d000000600010008c000000850000413d000000400100043d000006bd0010009c000001230000213d0000006003100039000000400030043f0000002003400039000000000432034f000000000404043b000000000004004b0000000005000039000000010500c039000000000054004b000000850000c13d00000000044104360000002003300039000000000532034f000000000505043b000006be0050009c000000850000213d00000000005404350000002003300039000000000232034f000000000302043b000006be0030009c000000850000213d0000004002100039000000000032043500000009030000290000008003300039000400000003001d00000000001304350000000703000029000000000303043300000040053000390000000005050433000006be065001970000000057030434000000000007004b0000070b0000613d000000000006004b0000126e0000613d0000000005050433000006be05500197000000000056004b000007100000413d0000126e0000013d000000000006004b0000125a0000c13d0000000005050433000006be005001980000125a0000c13d0000000002020433000006be022001970000000003010433000000000003004b0000071c0000613d000000000002004b000012750000613d0000000003040433000006be03300197000000000032004b000007210000413d000012750000013d000000000002004b0000125e0000c13d0000000002040433000006be002001980000125e0000c13d000000060100002900000000010104330000000001010433000000000001004b00000a7a0000613d000000090100002900000000010104330000066f01100197001100000001001d000000000010043f0000000601000039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f0000000100200190000000850000613d000000000101043b000000000101041a000000000001004b0000123f0000c13d0000000501000039000000000101041a0000066f0010009c000001230000213d00000001021000390000000503000039000000000023041b000006c20110009a0000001102000029000000000021041b000000000103041a001000000001001d000000000020043f0000000601000039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f0000000100200190000000850000613d000000000101043b0000001002000029000000000021041b000000090100002900000000010104330000066f01100197000000000010043f0000000701000039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f0000000100200190000000850000613d000000000301043b000000070100002900000000010104330000000024010434001100000004001d000000000004004b0000000004000039000000010400c039001000000004001d000000400400043d000006bc0040009c000001230000213d000f00000003001d0000000002020433000006be0220019700000040011000390000000001010433000b00000001001d000000a001400039000000400010043f000e00000002001d000d00000004001d0000000001240436000c00000001001d000006c301000041000000000010044300000000010004140000066b0010009c0000066b01008041000000c001100210000006c4011001c70000800b0200003919a819a30000040f0000000100200190000013d70000613d0000000b02000029000006be02200197000000000101043b0000000d040000290000004003400039000000100500002900000000005304350000008003400039000000000023043500000060034000390000000e0500002900000000005304350000066b011001970000000c030000290000000000130435000000110000006b0000000003000019000006c50300c0410000000f09000029000000000409041a000006c604400197000000000343019f0000008002200210000000000425019f000000000353019f0000008002100210000000000323019f000000000039041b0000000103900039000000000043041b000000400300043d000006bc0030009c000001230000213d0000000404000029000000000404043300000020054000390000000005050433000000000604043300000040044000390000000004040433000000a007300039000000400070043f0000008007300039000006be0840019700000000008704350000004007300039000000000006004b0000000006000039000000010600c039000000000067043500000020063000390000000000160435000006be015001970000006005300039000000000015043500000000001304350000000003000019000006c50300c0410000000205900039000000000605041a000006c606600197000000000363019f000000000223019f000000000212019f000000000025041b0000008002400210000000000112019f0000000302900039000000000012041b0000000601000029000000000301043300000000540304340000066f0040009c000001230000213d0000000406900039000000000106041a000000010010019000000001071002700000007f0770618f0000001f0070008c00000000020000390000000102002039000000000121013f000000010010019000000df50000c13d000000200070008c001100000006001d001000000004001d000e00000003001d000008000000413d000d00000007001d000f00000005001d000000000060043f00000000010004140000066b0010009c0000066b01008041000000c0011002100000067a011001c7000080100200003919a819a30000040f0000000100200190000000850000613d00000010040000290000001f024000390000000502200270000000200040008c0000000002004019000000000301043b0000000d010000290000001f01100039000000050110027000000000011300190000000002230019000000000012004b00000011060000290000000f05000029000008000000813d000000000002041b0000000102200039000000000012004b000007fc0000413d0000001f0040008c0000081f0000a13d000000000060043f00000000010004140000066b0010009c0000066b01008041000000c0011002100000067a011001c7000080100200003919a819a30000040f0000000100200190000000850000613d0000001007000029000006fd02700198000000000101043b0000000e080000290000082b0000613d000000010320008a0000000503300270000000000331001900000001043000390000002003000039000000110600002900000000058300190000000005050433000000000051041b00000020033000390000000101100039000000000041004b000008170000c13d0000082d0000013d000000000004004b000008230000613d0000000001050433000008240000013d00000000010000190000000302400210000006fe0220027f000006fe02200167000000000121016f0000000102400210000000000121019f000008390000013d00000020030000390000001106000029000000000072004b000008370000813d0000000302700210000000f80220018f000006fe0220027f000006fe0220016700000000038300190000000003030433000000000223016f000000000021041b000000010170021000000001011001bf000000000016041b000000080100002900000000010104330000000002010433000000000002004b0000095c0000613d0000000003000019000c00000003001d0000000502300210000000000121001900000020011000390000000001010433001000000001001d0000000031010434000000000001004b00000a7a0000613d00000009020000290000000002020433001100000002001d0000066b0010009c0000066b0100804100000060011002100000066b0030009c000f00000003001d0000066b0200004100000000020340190000004002200210000000000121019f00000000020004140000066b0020009c0000066b02008041000000c002200210000000000112019f000006b4011001c7000080100200003919a819a30000040f0000000100200190000000850000613d00000011020000290000066f02200197000000000101043b001100000001001d000b00000002001d000000000020043f0000000701000039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f0000000100200190000000850000613d000000000201043b0000001101000029000000000010043f000e00000002001d0000000601200039000d00000001001d000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f0000000100200190000000850000613d000000000101043b000000000101041a000000000001004b000010370000c13d0000000e010000290000000502100039000000000102041a0000066f0010009c000001230000213d000a00000001001d0000000101100039000000000012041b000e00000002001d000000000020043f00000000010004140000066b0010009c0000066b01008041000000c0011002100000067a011001c7000080100200003919a819a30000040f00000001002001900000001102000029000000850000613d000000000101043b0000000a01100029000000000021041b0000000e01000029000000000101041a000e00000001001d000000000020043f0000000d01000029000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f00000001002001900000001102000029000000850000613d000000000101043b0000000e03000029000000000031041b000000000020043f0000000801000039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f0000000100200190000000850000613d000000000401043b000000100100002900000000050104330000066f0050009c000001230000213d000000000104041a000000010010019000000001031002700000007f0330618f0000001f0030008c00000000020000390000000102002039000000000121013f00000001001001900000000f0700002900000df50000c13d000000200030008c001100000004001d000e00000005001d000008ec0000413d000d00000003001d000000000040043f00000000010004140000066b0010009c0000066b01008041000000c0011002100000067a011001c7000080100200003919a819a30000040f0000000100200190000000850000613d0000000e050000290000001f025000390000000502200270000000200050008c0000000002004019000000000301043b0000000d010000290000001f01100039000000050110027000000000011300190000000002230019000000000012004b0000000f070000290000001104000029000008ec0000813d000000000002041b0000000102200039000000000012004b000008e80000413d0000001f0050008c000009180000a13d000000000040043f00000000010004140000066b0010009c0000066b01008041000000c0011002100000067a011001c7000080100200003919a819a30000040f0000000100200190000000850000613d0000000e08000029000006fd02800198000000000101043b000009560000613d000000010320008a00000005033002700000000003310019000000010430003900000020030000390000000f07000029000000100600002900000000056300190000000005050433000000000051041b00000020033000390000000101100039000000000041004b000009030000c13d000000000082004b000009140000813d0000000302800210000000f80220018f000006fe0220027f000006fe0220016700000000036300190000000003030433000000000223016f000000000021041b000000010180021000000001011001bf0000001104000029000009240000013d000000000005004b0000091c0000613d00000000010704330000091d0000013d000000000100001900000010060000290000000302500210000006fe0220027f000006fe02200167000000000121016f0000000102500210000000000121019f000000000014041b000000400100043d00000020020000390000000003210436000000000206043300000000002304350000004003100039000000000002004b000009350000613d000000000400001900000000053400190000000006740019000000000606043300000000006504350000002004400039000000000024004b0000092e0000413d0000001f04200039000006fd044001970000000002320019000000000002043500000040024000390000066b0020009c0000066b0200804100000060022002100000066b0010009c0000066b010080410000004001100210000000000112019f00000000020004140000066b0020009c0000066b02008041000000c002200210000000000112019f000006b4011001c70000800d020000390000000203000039000006c8040000410000000b0500002919a8199e0000040f0000000100200190000000850000613d0000000c030000290000000103300039000000080100002900000000010104330000000002010433000000000023004b000008400000413d0000095c0000013d00000020030000390000000f070000290000001006000029000000000082004b0000090c0000413d000009140000013d00000004010000290000000002010433000000070100002900000000050104330000000601000029000000000301043300000009010000290000000004010433000000400100043d0000002006100039000001000700003900000000007604350000066f0440019700000000004104350000010007100039000000006403043400000000004704350000012003100039000000000004004b000009780000613d000000000700001900000000083700190000000009760019000000000909043300000000009804350000002007700039000000000047004b000009710000413d000000000643001900000000000604350000000076050434000000000006004b0000000006000039000000010600c039000000400810003900000000006804350000000006070433000006be066001970000006007100039000000000067043500000040055000390000000005050433000006be05500197000000800610003900000000005604350000000065020434000000000005004b0000000005000039000000010500c039000000a00710003900000000005704350000000005060433000006be05500197000000c006100039000000000056043500000040022000390000000002020433000006be02200197000000e00510003900000000002504350000001f02400039000006fd02200197000000000212004900000000023200190000066b0020009c0000066b0200804100000060022002100000066b0010009c0000066b010080410000004001100210000000000112019f00000000020004140000066b0020009c0000066b02008041000000c002200210000000000112019f000006b4011001c70000800d020000390000000103000039000006c90400004119a8199e0000040f0000000100200190000000850000613d00000005020000290000000102200039000500000002001d000000010020006c000005be0000413d00000a480000013d00000005054002100000003f065000390000067106600197000000400700043d0000000006670019000f00000007001d000000000076004b000000000700003900000001070040390000066f0060009c000001230000213d0000000100700190000001230000c13d0000010007300039000000400060043f0000000f030000290000000003430436000e00000003001d00000120022000390000000003250019000000000073004b000000850000213d000000000004004b000009d40000613d0000000e0400002900000000250204340000066e0050009c000000850000213d0000000004540436000000000032004b000009ce0000413d000001600500043d0000066e0050009c000000850000213d000001800200043d001000000002001d0000066e0020009c000000850000213d0000000003000411000000000003004b00000a4c0000c13d000000400100043d0000068a0200004100000a7c0000013d000006ed01000041000000800010043f000006b701000041000019aa00010430000006eb01000041000000800010043f000006b701000041000019aa00010430000000a005000039000006db0300004100000000040000190000000006050019000000000503041a000000000556043600000001033000390000000104400039000000000014004b000009ec0000413d000000410160008a000006fd04100197000006d30040009c000001230000213d0000008001400039000000400010043f0000000000210435000000a002400039000000800300043d0000000000320435000000c002400039000000000003004b00000a080000613d000000a004000039000000000500001900000000460404340000066e0660019700000000026204360000000105500039000000000035004b00000a020000413d00000000021200490000066b0020009c0000066b0200804100000060022002100000066b0010009c0000066b010080410000004001100210000000000112019f000019a90001042e000000a006000039000006d20400004100000000050000190000000007060019000000000604041a000000000667043600000001044000390000000105500039000000000025004b00000a140000413d000000410270008a000006fd04200197000006d30040009c000001230000213d0000008002400039000000800500043d000000400020043f0000066f0050009c000000280000a13d000001230000013d0000000402000039000000000302041a0000067204300197000000000414019f000000000042041b0000066e02300197000000800020043f000000a00010043f00000000010004140000066b0010009c0000066b01008041000000c001100210000006d5011001c70000800d020000390000000103000039000006d60400004119a8199e0000040f0000000100200190000000850000613d00000a480000013d000000000100041a0000067201100197000000000161019f000000000010041b00000000010004140000066b0010009c0000066b01008041000000c001100210000006b4011001c70000800d020000390000000303000039000006b50400004119a8199e0000040f0000000100200190000000850000613d0000000001000019000019a90001042e000000100100002900000a6c0000013d0000000106000039000000000406041a0000067204400197000000000334019f000000000036041b000000000005004b00000a7a0000613d0000066e0210019800000a7a0000613d000000100000006b00000a7a0000613d000b00000006001d000000800020043f000000c00050043f0000067301000041000000400300043d000d00000003001d00000000001304350000000001000414000000040020008c000c00000002001d00000a820000c13d00000000010004150000001f0110008a00000005011002100000000103000031000000200030008c00000020040000390000000004034019001f00000000003d00000ab00000013d0000000f01000029000000000112034f000000000101043b0000066e0010009c000000850000213d000000400200043d000006de030000410000000000320435000000040320003900000000001304350000066b0020009c0000066b020080410000004001200210000006bb011001c7000019aa00010430000000400100043d000006ca0200004100000000002104350000066b0010009c0000066b01008041000000400110021000000674011001c7000019aa000104300000000d030000290000066b0030009c0000066b0300804100000040033002100000066b0010009c0000066b01008041000000c001100210000000000131019f00000674011001c719a819a30000040f00000060031002700000066b03300197000000200030008c000000200400003900000000040340190000001f0640018f00000020074001900000000d0570002900000a9b0000613d000000000801034f0000000d09000029000000008a08043c0000000009a90436000000000059004b00000a970000c13d000000000006004b00000aa80000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f000300000001035500000000010004150000001e0110008a0000000501100210001e00000000003d000000010020019000000ac70000613d0000001f02400039000000600420018f0000000d02400029000000000042004b000000000400003900000001040040390000066f0020009c000001230000213d0000000100400190000001230000c13d000000400020043f000000200030008c000000850000413d0000000d030000290000000003030433000000ff0030008c000000850000213d0000000501100270000000000103001f0000001101000029000000ff0110018f000000000031004b00000dfd0000c13d0000001101000029000000a00010043f0000000402000039000000000102041a000006720110019700000010011001af000000000012041b0000000f010000290000000001010433000000000001004b0000000001000039000000010100c039000000e00010043f00000e0c0000613d000000400100043d000006770010009c000001230000213d0000002002100039000000400020043f0000000000010435000000e00100043d000000000001004b00000e080000c13d000000400100043d000006f00200004100000a7c0000013d00000010030000290000066b0030009c0000066b0300804100000040033002100000066b0010009c0000066b01008041000000c001100210000000000131019f000006bb011001c719a819a30000040f00000060031002700000066b03300197000000200030008c000000200400003900000000040340190000001f0640018f0000002007400190000000100570002900000afa0000613d000000000801034f0000001009000029000000008a08043c0000000009a90436000000000059004b00000af60000c13d000000000006004b00000b070000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f0003000000010355000000010020019000000c470000613d0000001f01400039000000600210018f0000001001200029000000000021004b000000000200003900000001020040390000066f0010009c000001230000213d0000000100200190000001230000c13d000000400010043f000000200030008c000000850000413d00000010020000290000000002020433000000000002004b0000000003000039000000010300c039000000000032004b000000850000c13d000000000002004b00000dfb0000c13d0000000f0100002900000020011000390000000201100367000000000101043b001000000001001d0000066e0010009c000000850000213d000006cb010000410000000000100443000000000100041200000004001004430000006001000039000000240010044300000000010004140000066b0010009c0000066b01008041000000c001100210000006dd011001c7000080050200003919a819a30000040f0000000100200190000013d70000613d000000000101043b000000000001004b0000104e0000c13d0000000f010000290000000201100367000000000101043b001000000001001d0000066f0010009c000000850000213d0000001001000029000000000010043f0000000601000039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f0000000100200190000000850000613d000000400200043d000e00000002001d0000000402200039000000000101043b000000000101041a000000000001004b0000111f0000c13d000006e7010000410000000e030000290000000000130435000000100100002900000c410000013d00000000030000310000000f01000029000000100200002919a814dd0000040f0000000002010019000000110100002919a8162c0000040f0000000001000019000019a90001042e0000000002000019001100000002001d0000000502200210001000000002001d0000000e012000290000000001010433000000000010043f0000000801000039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f0000000100200190000000850000613d000000000101043b000000000201041a000000010320019000000001052002700000007f0550618f0000001f0050008c00000000040000390000000104002039000000000043004b00000df50000c13d000000400700043d0000000004570436000000000003004b00000ba00000613d000900000004001d000a00000005001d000b00000007001d000000000010043f00000000010004140000066b0010009c0000066b01008041000000c0011002100000067a011001c7000080100200003919a819a30000040f0000000100200190000000850000613d0000000a08000029000000000008004b000000200500008a00000ba80000613d000000000201043b00000000010000190000000d060000290000000b0700002900000009090000290000000003190019000000000402041a000000000043043500000001022000390000002001100039000000000081004b00000b980000413d00000bab0000013d000006ff012001970000000000140435000000000005004b00000020010000390000000001006039000000200500008a0000000d0600002900000bab0000013d00000000010000190000000d060000290000000b070000290000003f01100039000000000251016f0000000001720019000000000021004b000000000200003900000001020040390000066f0010009c000001230000213d0000000100200190000001230000c13d000000400010043f00000000010604330000001102000029000000000021004b0000102b0000a13d00000010030000290000000c0130002900000000007104350000000001060433000000000021004b0000102b0000a13d00000001022000390000000f010000290000000001010433000000000012004b00000b630000413d000003ee0000013d00000000040000190000000d0c00002900000bd10000013d0000001f07600039000006fd077001970000000006650019000000000006043500000000057500190000000104400039000000000024004b000003f90000813d0000000006150049000000400660008a0000000003630436000000200cc0003900000000060c043300000000760604340000000005650436000000000006004b00000bc90000613d00000000080000190000000009580019000000000a870019000000000a0a04330000000000a904350000002008800039000000000068004b00000bdb0000413d00000bc90000013d0000000f030000290000066b0030009c0000066b0300804100000040033002100000066b0010009c0000066b01008041000000c001100210000000000131019f000006bb011001c719a819a30000040f00000060031002700000066b03300197000000200030008c000000200400003900000000040340190000001f0640018f00000020074001900000000f0570002900000bfc0000613d000000000801034f0000000f09000029000000008a08043c0000000009a90436000000000059004b00000bf80000c13d000000000006004b00000c090000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f0003000000010355000000010020019000000de40000613d0000001f01400039000000600210018f0000000f01200029000000000021004b000000000200003900000001020040390000066f0010009c000001230000213d0000000100200190000001230000c13d000000400010043f000000200030008c000000850000413d0000000f020000290000000002020433000000000002004b0000000003000039000000010300c039000000000032004b000000850000c13d000000000002004b00000dfb0000c13d0000000e010000290000000201100367000000000101043b000f00000001001d0000066f0010009c000000850000213d0000000f01000029000000000010043f0000000601000039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f0000000100200190000000850000613d000000400200043d000d00000002001d0000000402200039000000000101043b000000000101041a000000000001004b000010b80000c13d000006e7010000410000000d0300002900000000001304350000000f0100002900000000001204350000066b0030009c0000066b030080410000004001300210000006bb011001c7000019aa000104300000001f0530018f0000066d06300198000000400200043d000000000462001900000ec00000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b00000c4e0000c13d00000ec00000013d0000000002000019000700000002001d000000050120021000000005011000290000000201100367000000000101043b000b00000001001d0000066f0010009c000000850000213d0000000b01000029000000000010043f0000000601000039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f0000000100200190000000850000613d000000000101043b000000000301041a000000000003004b000010d40000613d0000000501000039000000000201041a000000000002004b000013e40000613d000000010130008a000000000023004b00000c8f0000613d000000000012004b0000102b0000a13d000006b80130009a000006b80220009a000000000202041a000000000021041b000000000020043f0000000601000039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c70000801002000039001100000003001d19a819a30000040f0000000100200190000000850000613d000000000101043b0000001102000029000000000021041b0000000501000039000000000301041a000000000003004b000010310000613d000000010130008a000006b80230009a000000000002041b0000000502000039000000000012041b0000000b01000029000000000010043f0000000601000039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f0000000100200190000000850000613d000000000101043b000000000001041b0000000b01000029000000000010043f0000000701000039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f0000000100200190000000850000613d000000000101043b0000000501100039000000000301041a000000400200043d000f00000002001d001100000003001d0000000002320436000a00000002001d000000000010043f00000000010004140000066b0010009c0000066b01008041000000c0011002100000067a011001c7000080100200003919a819a30000040f0000000100200190000000850000613d0000001105000029000000000005004b0000000a0200002900000cce0000613d000000000101043b0000000a020000290000000003000019000000000401041a000000000242043600000001011000390000000103300039000000000053004b00000cc80000413d0000000f0120006a0000001f01100039000006fd021001970000000f01200029000000000021004b000000000200003900000001020040390000066f0010009c000001230000213d0000000100200190000001230000c13d000000400010043f0000000f010000290000000001010433000000000001004b00000d730000613d000000000200001900000ce80000013d000000000101043b000000000001041b000000110200002900000001022000390000000f010000290000000001010433000000000012004b00000d730000813d001100000002001d0000000b01000029000000000010043f0000000701000039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f0000000100200190000000850000613d000000000301043b0000000f0100002900000000010104330000001102000029000000000021004b0000102b0000a13d00000005012002100000000a011000290000000001010433000c00000001001d000000000010043f000d00000003001d0000000601300039000e00000001001d000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f0000000100200190000000850000613d000000000101043b000000000101041a000000000001004b00000ce20000613d0000000d020000290000000503200039000000000203041a000000000002004b000013e40000613d000000000021004b001000000001001d000d00000003001d00000d540000613d000900000002001d000000000030043f00000000010004140000066b0010009c0000066b01008041000000c0011002100000067a011001c7000080100200003919a819a30000040f0000000100200190000000850000613d00000010020000290008000100200092000000000101043b0000000d04000029000000000204041a000000080020006c0000102b0000a13d0000000902000029000000010220008a0000000001120019000000000101041a000900000001001d000000000040043f00000000010004140000066b0010009c0000066b01008041000000c0011002100000067a011001c7000080100200003919a819a30000040f0000000100200190000000850000613d000000000101043b00000008011000290000000902000029000000000021041b000000000020043f0000000e01000029000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f0000000100200190000000850000613d000000000101043b0000001002000029000000000021041b0000000d03000029000000000103041a001000000001001d000000000001004b000010310000613d000000000030043f00000000010004140000066b0010009c0000066b01008041000000c0011002100000067a011001c7000080100200003919a819a30000040f0000000100200190000000850000613d0000001002000029000000010220008a000000000101043b0000000001210019000000000001041b0000000d01000029000000000021041b0000000c01000029000000000010043f0000000e01000029000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f000000010020019000000ce00000c13d000000850000013d0000000b01000029000000000010043f0000000701000039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f0000000100200190000000850000613d000000000301043b000000000003041b0000000101300039000000000001041b0000000201300039000000000001041b0000000301300039000000000001041b0000000404300039000000000104041a000000010010019000000001051002700000007f0550618f0000001f0050008c00000000020000390000000102002039000000000121013f000000010010019000000df50000c13d000000000005004b00000db50000613d0000001f0050008c00000db40000a13d000f00000005001d001100000003001d001000000004001d000000000040043f00000000010004140000066b0010009c0000066b01008041000000c0011002100000067a011001c7000080100200003919a819a30000040f0000000100200190000000850000613d000000000101043b0000000f020000290000001f02200039000000050220027000000000022100190000000103100039000000000023004b00000db00000813d000000000003041b0000000103300039000000000023004b00000dac0000413d0000001002000029000000000002041b00000000040100190000001103000029000000000004041b0000000501300039000000000201041a000000000001041b000000000002004b00000dcd0000613d001100000002001d000000000010043f00000000010004140000066b0010009c0000066b01008041000000c0011002100000067a011001c7000080100200003919a819a30000040f0000000100200190000000850000613d000000000101043b0000001102100029000000000021004b00000dcd0000813d000000000001041b0000000101100039000000000021004b00000dc90000413d000000400100043d0000000b0200002900000000002104350000066b0010009c0000066b01008041000000400110021000000000020004140000066b0020009c0000066b02008041000000c002200210000000000112019f0000067a011001c70000800d020000390000000103000039000006b90400004119a8199e0000040f0000000100200190000000850000613d00000007020000290000000102200039000000060020006c00000c540000413d000005bb0000013d0000001f0530018f0000066d06300198000000400200043d000000000462001900000ec00000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b00000deb0000c13d00000ec00000013d000006cd02000041000001400020043f000001440010043f000006ce01000041000019aa00010430000006dc01000041000000000010043f0000002201000039000000040010043f000006bb01000041000019aa00010430000006e10200004100000a7c0000013d0000002404200039000000000034043500000675030000410000000000320435000000040320003900000000001304350000066b0020009c0000066b02008041000000400120021000000676011001c7000019aa000104300000000f010000290000000001010433000000000001004b00000ef70000c13d000000400300043d0000067c010000410000000000130435000000240130003900000000020004100000000000210435001100000003001d0000000401300039000000000021043500000000010004140000000c02000029000000040020008c00000e1e0000c13d0000000103000031000000200030008c0000002004000039000000000403401900000e490000013d00000011020000290000066b0020009c0000066b0200804100000040022002100000066b0010009c0000066b01008041000000c001100210000000000121019f00000676011001c70000000c0200002919a819a30000040f00000060031002700000066b03300197000000200030008c000000200400003900000000040340190000001f0640018f0000002007400190000000110570002900000e380000613d000000000801034f0000001109000029000000008a08043c0000000009a90436000000000059004b00000e340000c13d000000000006004b00000e450000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f0003000000010355000000010020019000000eb50000613d0000001f01400039000000600210018f0000001101200029000000000021004b000000000200003900000001020040390000066f0010009c000001230000213d0000000100200190000001230000c13d000000400010043f000000200040008c000000850000413d00000011020000290000000002020433000000000002004b000013e40000c13d0000004402100039000000010400008a000000000042043500000020021000390000067d040000410000000000420435000000240410003900000000050004100000000000540435000000440400003900000000004104350000067e0010009c000001230000213d000000c004100039000000400040043f000000a0051000390000067f04000041001000000005001d000000000045043500000080051000390000002004000039000f00000005001d0000000000450435000000000401043300000000010004140000000c05000029000000040050008c00000e880000613d0000066b0020009c0000066b0200804100000040022002100000066b0040009c0000066b040080410000006003400210000000000223019f0000066b0010009c0000066b01008041000000c001100210000000000112019f0000000c0200002919a8199e0000040f000b000100200193000300000001035500000060011002700001066b0010019d0000066b03100197000000000003004b000010650000c13d001100600000003d000e00800000003d000000110100002900000000010104330000000b0000006b000010990000c13d000000000001004b000010cb0000c13d000000400100043d000006850200004100000000002104350000000402100039000000200300003900000000003204350000000f020000290000000002020433000000240310003900000000002304350000004403100039000000000002004b000000100700002900000ea80000613d000000000400001900000000053400190000000006740019000000000606043300000000006504350000002004400039000000000024004b00000ea10000413d0000001f04200039000006fd044001970000000002320019000000000002043500000044024000390000066b0020009c0000066b0200804100000060022002100000066b0010009c0000066b010080410000004001100210000000000112019f000019aa000104300000001f0530018f0000066d06300198000000400200043d000000000462001900000ec00000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b00000ebc0000c13d000000000005004b00000ecd0000613d000000000161034f0000000305500210000000000604043300000000065601cf000000000656022f000000000101043b0000010005500089000000000151022f00000000015101cf000000000161019f000000000014043500000060013002100000066b0020009c0000066b020080410000004002200210000000000112019f000019aa00010430000000000001004b00000f430000613d000000a00200043d000006be02200197000000000021004b00000f430000813d0000001101000029000000000010043f0000000701000039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f0000000100200190000000850000613d000000000101043b000000800200003919a818580000040f000001200100043d000006be01100197000000e00200043d000000000002004b00000f840000c13d000000000001004b00000ef30000c13d000001000100043d000006be0010019800000f8a0000613d000000400200043d001100000002001d000006c001000041000010950000013d000000000200001900000eff0000013d000000110200002900000001022000390000000f010000290000000001010433000000000012004b00000e0c0000813d001100000002001d00000005012002100000000e0110002900000000010104330000066e0310019800000ef90000613d000000000030043f0000000301000039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c70000801002000039001000000003001d19a819a30000040f00000010040000290000000100200190000000850000613d000000000101043b000000000101041a000000000001004b00000ef90000c13d0000000203000039000000000103041a0000066f0010009c000001230000213d0000000102100039000000000023041b000006790110009a000000000041041b000000000103041a000d00000001001d000000000040043f0000000301000039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f00000010030000290000000100200190000000850000613d000000000101043b0000000d02000029000000000021041b000000400100043d00000000003104350000066b0010009c0000066b01008041000000400110021000000000020004140000066b0020009c0000066b02008041000000c002200210000000000112019f0000067a011001c70000800d0200003900000001030000390000067b0400004119a8199e0000040f000000010020019000000ef90000c13d000000850000013d000000400200043d001100000002001d000006bf0100004100000000001204350000000402200039000000800100003919a816150000040f00000011010000290000066b0010009c0000066b01008041000000400110021000000689011001c7000019aa000104300000002001700039000000100200002900000000002104350000002001000039000000000017043500000040017000390000000a021000290000000f0300002900000002033003670000000a0000006b00000f610000613d000000000403034f0000000005010019000000004604043c0000000005650436000000000025004b00000f5d0000c13d0000000e0000006b00000f6f0000613d0000000a033003600000000e040000290000000304400210000000000502043300000000054501cf000000000545022f000000000303043b0000010004400089000000000343022f00000000034301cf000000000353019f0000000000320435000000100110002900000000000104350000066b0070009c0000066b0700804100000040017002100000000b02000029000006d70020009c000006d7020080410000006002200210000000000112019f00000000020004140000066b0020009c0000066b02008041000000c002200210000000000121019f000006d80110009a0000800d020000390000000203000039000006d904000041000000110500002900000a350000013d000000000001004b000010920000613d000001000200043d000006be02200197000000000021004b000010920000813d0000001101000029000000000010043f0000000701000039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f0000000100200190000000850000613d000000000101043b0000000201100039000000e00200003919a818580000040f000000400100043d00000011020000290000000002210436000000800300043d000000000003004b0000000003000039000000010300c0390000000000320435000000a00200043d000006be0220019700000040031000390000000000230435000000c00200043d000006be0220019700000060031000390000000000230435000000e00200043d000000000002004b0000000002000039000000010200c03900000080031000390000000000230435000001000200043d000006be02200197000000a0031000390000000000230435000001200200043d000006be02200197000000c00310003900000000002304350000066b0010009c0000066b01008041000000400110021000000000020004140000066b0020009c0000066b02008041000000c002200210000000000112019f000006cf011001c70000800d020000390000000103000039000006d00400004100000a350000013d000000000200001900000fcd0000013d00000010020000290000000102200039000000800100043d000000000012004b0000047f0000813d001000000002001d0000000501200210000000a00110003900000000010104330000066e01100197001100000001001d000000000010043f0000000301000039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f0000000100200190000000850000613d000000000101043b000000000301041a000000000003004b00000fc80000613d0000000201000039000000000201041a000000000002004b000013e40000613d000000010130008a000000000023004b000010050000613d000000000012004b0000102b0000a13d000006ee0130009a000006ee0220009a000000000202041a000000000021041b000000000020043f0000000301000039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c70000801002000039000d00000003001d19a819a30000040f0000000d030000290000000100200190000000850000613d000000000101043b000000000031041b0000000201000039000000000301041a000000000003004b000010310000613d000000010130008a000006ee0230009a000000000002041b0000000202000039000000000012041b0000001101000029000000000010043f0000000301000039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f0000000100200190000000850000613d000000000101043b000000000001041b000000400100043d000000110200002900000000002104350000066b0010009c0000066b01008041000000400110021000000000020004140000066b0020009c0000066b02008041000000c002200210000000000112019f0000067a011001c70000800d020000390000000103000039000006ef0400004119a8199e0000040f000000010020019000000fc80000c13d000000850000013d000006dc01000041000000000010043f0000003201000039000000040010043f000006bb01000041000019aa00010430000006dc01000041000000000010043f0000003101000039000000040010043f000006bb01000041000019aa00010430000000400300043d001100000003001d000000240130003900000040020000390000000000210435000006c701000041000000000013043500000004013000390000000b0200002900000000002104350000004402300039000000100100002919a8146b0000040f000000110200002900000000012100490000066b0010009c0000066b010080410000066b0020009c0000066b0200804100000060011002100000004002200210000000000121019f000019aa000104300000001001000029000000000010043f0000000301000039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f0000000100200190000000850000613d000000000101043b000000000101041a000000000001004b00000b3a0000c13d000000400100043d000006e2020000410000000000210435000000040210003900000010030000290000038d0000013d0000066f0030009c000001230000213d0000001f02300039000006fd022001970000003f02200039000006fd02200197000000400400043d0000000002240019001100000004001d000000000042004b000000000400003900000001040040390000066f0020009c000001230000213d0000000100400190000001230000c13d000000400020043f00000011020000290000000004320436000006fd023001980000001f0330018f000e00000004001d00000000012400190000000304000367000010840000613d000000000504034f0000000e06000029000000005705043c0000000006760436000000000016004b000010800000c13d000000000003004b00000e8c0000613d000000000224034f0000000303300210000000000401043300000000043401cf000000000434022f000000000202043b0000010003300089000000000232022f00000000023201cf000000000242019f000000000021043500000e8c0000013d000000400200043d001100000002001d000006bf0100004100000000001204350000000402200039000000e00100003900000f490000013d000000000001004b0000121e0000c13d000006800100004100000000001004430000000c01000029000000040010044300000000010004140000066b0010009c0000066b01008041000000c00110021000000681011001c7000080020200003919a819a30000040f0000000100200190000013d70000613d000000000101043b000000000001004b0000121a0000c13d000000400100043d00000044021000390000068803000041000000000032043500000024021000390000001d0300003900000000003204350000068502000041000000000021043500000004021000390000002003000039000000000032043500000f4b0000013d0000000401000039000000000301041a000006f1010000410000000d0400002900000000001404350000000f01000029000000000012043500000024014000390000000002000411000000000021043500000000010004140000066e02300197000000040020008c000010da0000c13d0000000103000031000000200030008c00000020040000390000000004034019000011040000013d0000000e020000290000066b0020009c0000066b0200804100000040022002100000066b0010009c0000066b010080410000006001100210000000000121019f000019aa00010430000000400100043d000006ba02000041000000000021043500000004021000390000000b030000290000038d0000013d0000000d030000290000066b0030009c0000066b0300804100000040033002100000066b0010009c0000066b01008041000000c001100210000000000131019f00000676011001c719a819a30000040f00000060031002700000066b03300197000000200030008c000000200400003900000000040340190000001f0640018f00000020074001900000000d05700029000010f30000613d000000000801034f0000000d09000029000000008a08043c0000000009a90436000000000059004b000010ef0000c13d000000000006004b000011000000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f000300000001035500000001002001900000112f0000613d0000001f01400039000000600210018f0000000d01200029000000000021004b000000000200003900000001020040390000066f0010009c000001230000213d0000000100200190000001230000c13d000000400010043f000000200030008c000000850000413d0000000d020000290000000002020433000000000002004b0000000003000039000000010300c039000000000032004b000000850000c13d000000000002004b0000127b0000c13d000006e4020000410000000000210435000000040210003900000000030004110000038d0000013d0000000401000039000000000301041a000006e3010000410000000e0400002900000000001404350000001001000029000000000012043500000000010004140000066e02300197000000040020008c0000113b0000c13d0000000103000031000000200030008c00000020040000390000000004034019000011650000013d0000001f0530018f0000066d06300198000000400200043d000000000462001900000ec00000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000011360000c13d00000ec00000013d0000000e030000290000066b0030009c0000066b0300804100000040033002100000066b0010009c0000066b01008041000000c001100210000000000131019f000006bb011001c719a819a30000040f00000060031002700000066b03300197000000200030008c000000200400003900000000040340190000001f0640018f00000020074001900000000e05700029000011540000613d000000000801034f0000000e09000029000000008a08043c0000000009a90436000000000059004b000011500000c13d000000000006004b000011610000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f00030000000103550000000100200190000012620000613d0000001f01400039000000600210018f0000000e01200029000000000021004b000000000200003900000001020040390000066f0010009c000001230000213d0000000100200190000001230000c13d000000400010043f000000200030008c000000850000413d0000000e020000290000000002020433000006c50020009c000000850000813d0000000003000411000000000023004b000012bb0000c13d00000002010003670000000f02100360000000000202043b0000066f0020009c000000850000213d0000000f030000290000004003300039000000000131034f000000000101043b001000000001001d000000000020043f0000000701000039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f0000000100200190000000850000613d000000000101043b0000001002000029000000110300002919a818c50000040f000006800100004100000000001004430000001101000029000000040010044300000000010004140000066b0010009c0000066b01008041000000c00110021000000681011001c7000080020200003919a819a30000040f0000000100200190000013d70000613d000000000101043b000000000001004b000000850000613d000000400300043d000000240130003900000010020000290000000000210435000006e501000041000000000013043500000000010004100000066e01100197000e00000003001d0000000402300039000000000012043500000000010004140000001102000029000000040020008c000011c20000613d0000000e020000290000066b0020009c0000066b0200804100000040022002100000066b0010009c0000066b01008041000000c001100210000000000121019f00000676011001c7000000110200002919a8199e0000040f00000060031002700001066b0030019d00030000000103550000000100200190000013460000613d0000000e010000290000066f0010009c000001230000213d0000000e02000029000000400020043f000000100100002900000000001204350000066b0020009c0000066b02008041000000400120021000000000020004140000066b0020009c0000066b02008041000000c002200210000000000112019f0000067a011001c70000800d020000390000000203000039000006e604000041000000000500041119a8199e0000040f0000000100200190000000850000613d0000000f010000290000000201100367000000000101043b0000066f0010009c000000850000213d19a815ba0000040f000000400200043d001100000002001d0000000002000412001700000002001d001600200000003d001000000001001d000080050100003900000044030000390000000004000415000000170440008a0000000504400210000006cb0200004119a819800000040f000000ff0310018f0000001101000029000000200210003900000000003204350000002002000039000000000021043519a814550000040f000000400100043d000e00000001001d19a814550000040f0000000e020000290000002001200039000f00000001001d00000011030000290000000000310435000000100100002900000000001204350000000003020019000000400400043d001100000004001d00000020010000390000000002140436000000000103043300000040030000390000000000320435000000600240003919a8146b0000040f000000000201001900000011040000290000000001410049000000200310008a0000000f0100002900000000010104330000004004400039000000000034043519a8146b0000040f000000110200002900000000012100490000066b0020009c0000066b0200804100000040022002100000066b0010009c0000066b010080410000006001100210000000000121019f000019a90001042e00000011010000290000000001010433000000000001004b0000122b0000613d000006820010009c000000850000213d000000200010008c000000850000413d0000000e010000290000000001010433000000000001004b0000000002000039000000010200c039000000000021004b000000850000c13d000000000001004b000012460000613d000000800100043d00000140000004430000016000100443000000a00100043d00000020030000390000018000300443000001a000100443000000c00100043d0000004002000039000001c000200443000001e0001004430000006001000039000000e00200043d000002000010044300000220002004430000010000300443000000040100003900000120001004430000068701000041000019a90001042e00000009010000290000000001010433000000400200043d000006c10300004100000000003204350000066f0110019700000a730000013d000000400100043d00000064021000390000068303000041000000000032043500000044021000390000068403000041000000000032043500000024021000390000002a030000390000000000320435000006850200004100000000002104350000000402100039000000200300003900000000003204350000066b0010009c0000066b01008041000000400110021000000686011001c7000019aa00010430000000400200043d001100000002001d000006c001000041000012710000013d000000400300043d001100000003001d000006c002000041000012780000013d0000001f0530018f0000066d06300198000000400200043d000000000462001900000ec00000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000012690000c13d00000ec00000013d000000400200043d001100000002001d000006bf0100004100000000001204350000000402200039000000000103001900000f490000013d000000400300043d001100000003001d000006bf020000410000000000230435000000040230003900000f490000013d00000002020003670000000e01200360000000000101043b0000066f0010009c000000850000213d0000000e030000290000008003300039000000000332034f000000000303043b0000000004000031000000110540006a000000230550008a00000670065001970000067007300197000000000867013f000000000067004b00000000060000190000067006004041000000000053004b00000000050000190000067005008041000006700080009c000000000605c019000000000006004b000000850000c13d0000001105000029000f00040050003d0000000f05300029000000000252034f000000000302043b0000066f0030009c000000850000213d0000000004340049000000200250003900000670054001970000067006200197000000000756013f000000000056004b00000000050000190000067005004041000000000042004b00000000040000190000067004002041000006700070009c000000000504c019000000000005004b000000850000c13d19a815150000040f000000000001004b000012bf0000c13d0000001101000029000000a4021000390000000f0100002919a814b30000040f000006f603000041000000400500043d001100000005001d0000000000350435000000040350003900000020040000390000000000430435000000240350003919a815860000040f000010440000013d000006e402000041000000000021043500000004021000390000038d0000013d00000002010003670000000e02100360000000000202043b0000066f0020009c000000850000213d0000000e03000029000d00400030003d0000000d01100360000000000101043b000e00000001001d000000000020043f0000000701000039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f0000000100200190000000850000613d000000000101043b00000002011000390000000e02000029000000100300002919a818c50000040f0000000d01000029000d00600010003d00000002030003670000000d01300360000000000101043b0000000004000031000000110240006a000000230220008a00000670052001970000067006100197000000000756013f000000000056004b00000000050000190000067005004041000000000021004b00000000020000190000067002008041000006700070009c000000000502c019000000000005004b000000850000c13d0000000f01100029000000000213034f000000000202043b0000066f0020009c000000850000213d0000000005240049000000200610003900000670015001970000067007600197000000000817013f000000000017004b00000000010000190000067001004041000000000056004b00000000050000190000067005002041000006700080009c000000000105c019000000000001004b000000850000c13d0000001f01200039000006fd011001970000003f01100039000006fd05100197000000400100043d0000000005510019000000000015004b000000000800003900000001080040390000066f0050009c000001230000213d0000000100800190000001230000c13d000000400050043f00000000052104360000000008620019000000000048004b000000850000213d000000000463034f000006fd062001980000001f0720018f0000000003650019000013200000613d000000000804034f0000000009050019000000008a08043c0000000009a90436000000000039004b0000131c0000c13d000000000007004b0000132d0000613d000000000464034f0000000306700210000000000703043300000000076701cf000000000767022f000000000404043b0000010006600089000000000464022f00000000046401cf000000000474019f0000000000430435000000000225001900000000000204350000000002010433000000200020008c000013530000613d000000000002004b000013570000c13d000006cb010000410000000000100443000000000100041200000004001004430000002001000039000000240010044300000000010004140000066b0010009c0000066b01008041000000c001100210000006dd011001c7000080050200003919a819a30000040f0000000100200190000013d70000613d000000000101043b001100000001001d000013600000013d0000066b033001970000001f0530018f0000066d06300198000000400200043d000000000462001900000ec00000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b0000134e0000c13d00000ec00000013d0000000002050433001100000002001d000000ff0020008c000013600000a13d000000400400043d001100000004001d000006f50200004100000000002404350000000402400039000000200300003900000000003204350000002402400039000010430000013d000006cb010000410000000000100443000000000100041200000004001004430000002001000039000000240010044300000000010004140000066b0010009c0000066b01008041000000c001100210000006dd011001c7000080050200003919a819a30000040f0000000100200190000013d70000613d0000001102000029000000ff0220018f000000000301043b000000ff0430018f000000000142004b000013d80000c13d0000000e01000029001100000001001d0000000d01000029000000800110008a000e00000001001d0000000201100367000000000101043b000f00000001001d0000066e0010009c000000850000213d000006800100004100000000001004430000001001000029000000040010044300000000010004140000066b0010009c0000066b01008041000000c00110021000000681011001c7000080020200003919a819a30000040f0000000100200190000013d70000613d000000000101043b000000000001004b000000850000613d000000400300043d000000240130003900000011020000290000000000210435000006f3010000410000000000130435000d00000003001d00000004013000390000000f02000029000000000021043500000000010004140000001002000029000000040020008c000013ad0000613d0000000d020000290000066b0020009c0000066b0200804100000040022002100000066b0010009c0000066b01008041000000c001100210000000000121019f00000676011001c7000000100200002919a8199e0000040f00000060031002700001066b0030019d00030000000103550000000100200190000013f00000613d0000000d010000290000066f0010009c000001230000213d0000000d01000029000000400010043f0000000e010000290000000201100367000000000601043b0000066e0060009c000000850000213d00000011010000290000000d0200002900000000001204350000066b0020009c0000066b02008041000000400120021000000000020004140000066b0020009c0000066b02008041000000c002200210000000000112019f0000067a011001c70000800d020000390000000303000039000006f404000041000000000500041119a8199e0000040f0000000100200190000000850000613d000000400100043d001000000001001d19a8144a0000040f000000110200002900000010010000290000000000210435000000400100043d00000000002104350000066b0010009c0000066b010080410000004001100210000006e8011001c7000019a90001042e000000000001042f000013e10000a13d000000ff0010008c000013e40000213d0000004d0010008c0000141f0000213d000000000001004b000013fd0000c13d0000000102000039000014080000013d0000000001240049000000ff0010008c000013ea0000a13d000006dc01000041000000000010043f0000001101000039000000040010043f000006bb01000041000019aa000104300000004e0010008c0000141f0000813d000000000001004b0000140a0000c13d00000001020000390000141c0000013d0000066b033001970000001f0530018f0000066d06300198000000400200043d000000000462001900000ec00000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000013f80000c13d00000ec00000013d0000000a030000390000000102000039000000010010019000000000043300a9000000010300603900000000022300a900000001011002720000000003040019000013ff0000c13d000000000002004b000014160000613d0000000e012000f9000013760000013d0000000a0500003900000001020000390000000004010019000000010040019000000000065500a9000000010500603900000000022500a9000000010440027200000000050600190000140d0000c13d000000000002004b0000141c0000c13d000006dc01000041000000000010043f0000001201000039000000040010043f000006bb01000041000019aa00010430000006fe022001290000000e0020006c000014310000813d000000400200043d001000000002001d000006f2010000410000000000120435000000040120003900000011020000290000000e0400002919a816230000040f000000100200002900000000012100490000066b0010009c0000066b0100804100000060011002100000066b0020009c0000066b020080410000004002200210000000000121019f000019aa00010430000000ff0210018f0000004d0020008c000013e40000213d000000000002004b000014380000c13d0000000101000039000014410000013d0000000a030000390000000101000039000000010020019000000000043300a9000000010300603900000000011300a9000000010220027200000000030400190000143a0000c13d0000000e0000006b001100000000001d000013770000613d0000000e031000b9001100000003001d0000000e023000fa000000000012004b000013770000613d000013e40000013d000007000010009c0000144f0000813d0000002001100039000000400010043f000000000001042d000006dc01000041000000000010043f0000004101000039000000040010043f000006bb01000041000019aa00010430000007010010009c0000145a0000813d0000004001100039000000400010043f000000000001042d000006dc01000041000000000010043f0000004101000039000000040010043f000006bb01000041000019aa00010430000007020010009c000014650000813d000000a001100039000000400010043f000000000001042d000006dc01000041000000000010043f0000004101000039000000040010043f000006bb01000041000019aa0001043000000000430104340000000001320436000000000003004b000014770000613d000000000200001900000000052100190000000006240019000000000606043300000000006504350000002002200039000000000032004b000014700000413d000000000231001900000000000204350000001f02300039000006fd022001970000000001210019000000000001042d000006820010009c000014980000213d000000430010008c000014980000a13d00000002020003670000000403200370000000000403043b0000066f0040009c000014980000213d0000002403200370000000000503043b0000066f0050009c000014980000213d0000002303500039000000000013004b000014980000813d0000000403500039000000000232034f000000000302043b0000066f0030009c000014980000213d00000024025000390000000005320019000000000015004b000014980000213d0000000001040019000000000001042d0000000001000019000019aa000104300000000043020434000006be03300197000000000331043600000000040404330000066b04400197000000000043043500000040032000390000000003030433000000000003004b0000000003000039000000010300c0390000004004100039000000000034043500000060032000390000000003030433000006be033001970000006004100039000000000034043500000080022000390000000002020433000006be0220019700000080031000390000000000230435000000a001100039000000000001042d0000000204000367000000000224034f000000000202043b000000000300003100000000051300490000001f0550008a00000670065001970000067007200197000000000867013f000000000067004b00000000060000190000067006002041000000000052004b00000000050000190000067005004041000006700080009c000000000605c019000000000006004b000014db0000613d0000000001120019000000000214034f000000000202043b0000066f0020009c000014db0000213d0000000003230049000000200110003900000670043001970000067005100197000000000645013f000000000045004b00000000040000190000067004004041000000000031004b00000000030000190000067003002041000006700060009c000000000403c019000000000004004b000014db0000c13d000000000001042d0000000001000019000019aa00010430000007030020009c0000150d0000813d00000000040100190000001f01200039000006fd011001970000003f01100039000006fd05100197000000400100043d0000000005510019000000000015004b000000000700003900000001070040390000066f0050009c0000150d0000213d00000001007001900000150d0000c13d000000400050043f00000000052104360000000007420019000000000037004b000015130000213d000006fd062001980000001f0720018f00000002044003670000000003650019000014fd0000613d000000000804034f0000000009050019000000008a08043c0000000009a90436000000000039004b000014f90000c13d000000000007004b0000150a0000613d000000000464034f0000000306700210000000000703043300000000076701cf000000000767022f000000000404043b0000010006600089000000000464022f00000000046401cf000000000474019f000000000043043500000000022500190000000000020435000000000001042d000006dc01000041000000000010043f0000004101000039000000040010043f000006bb01000041000019aa000104300000000001000019000019aa000104300003000000000002000300000003001d000200000002001d0000066f01100197000000000010043f0000000701000039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f0000000207000029000000030a00002900000001002001900000157e0000613d0000000003000031000000000601043b0000070300a0009c000015800000813d0000001f01a00039000006fd011001970000003f01100039000006fd02100197000000400100043d0000000002210019000000000012004b000000000500003900000001050040390000066f0020009c000015800000213d0000000100500190000015800000c13d000100000006001d000000400020043f0000000002a1043600000000057a0019000000000035004b0000157e0000213d000006fd04a001980000001f05a0018f00000002067003670000000003420019000015490000613d000000000706034f0000000008020019000000007907043c0000000008980436000000000038004b000015450000c13d000000000005004b000015560000613d000000000446034f0000000305500210000000000603043300000000065601cf000000000656022f000000000404043b0000010005500089000000000454022f00000000045401cf000000000464019f00000000004304350000000003a2001900000000000304350000066b0020009c0000066b02008041000000400220021000000000010104330000066b0010009c0000066b010080410000006001100210000000000121019f00000000020004140000066b0020009c0000066b02008041000000c002200210000000000112019f000006b4011001c7000080100200003919a819a30000040f00000001002001900000157e0000613d000000000101043b000000000010043f00000001010000290000000601100039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f00000001002001900000157e0000613d000000000101043b000000000101041a000000000001004b0000000001000039000000010100c039000000000001042d0000000001000019000019aa00010430000006dc01000041000000000010043f0000004101000039000000040010043f000006bb01000041000019aa000104300000000003230436000006fd062001980000001f0720018f00000000056300190000000201100367000015920000613d000000000801034f0000000009030019000000008a08043c0000000009a90436000000000059004b0000158e0000c13d000000000007004b0000159f0000613d000000000161034f0000000306700210000000000705043300000000076701cf000000000767022f000000000101043b0000010006600089000000000161022f00000000016101cf000000000171019f0000000000150435000000000123001900000000000104350000001f01200039000006fd011001970000000001130019000000000001042d000000400100043d000007020010009c000015b40000813d000000a002100039000000400020043f000000800210003900000000000204350000006002100039000000000002043500000040021000390000000000020435000000200210003900000000000204350000000000010435000000000001042d000006dc01000041000000000010043f0000004101000039000000040010043f000006bb01000041000019aa0001043000030000000000020000066f01100197000000000010043f0000000701000039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f0000000100200190000016070000613d000000000101043b0000000405100039000000000205041a000000010320019000000001062002700000007f0660618f0000001f0060008c00000000040000390000000104002039000000000043004b000016090000c13d000000400100043d0000000004610436000000000003004b000015f30000613d000100000004001d000200000006001d000300000001001d000000000050043f00000000010004140000066b0010009c0000066b01008041000000c0011002100000067a011001c7000080100200003919a819a30000040f0000000100200190000016070000613d0000000206000029000000000006004b000015f90000613d000000000201043b0000000005000019000000030100002900000001070000290000000003570019000000000402041a000000000043043500000001022000390000002005500039000000000065004b000015eb0000413d000015fb0000013d000006ff022001970000000000240435000000000006004b00000020050000390000000005006039000015fb0000013d000000000500001900000003010000290000003f03500039000006fd023001970000000003120019000000000023004b000000000200003900000001020040390000066f0030009c0000160f0000213d00000001002001900000160f0000c13d000000400030043f000000000001042d0000000001000019000019aa00010430000006dc01000041000000000010043f0000002201000039000000040010043f000006bb01000041000019aa00010430000006dc01000041000000000010043f0000004101000039000000040010043f000006bb01000041000019aa000104300000000043010434000000000003004b0000000003000039000000010300c03900000000033204360000000004040433000006be044001970000000000430435000000400220003900000040011000390000000001010433000006be011001970000000000120435000000000001042d00000040051000390000000000450435000000ff0330018f00000020041000390000000000340435000000ff0220018f00000000002104350000006001100039000000000001042d0007000000000002000400000001001d000600000002001d0000000021020434000000000001004b000017450000613d0000066b0010009c0000066b0100804100000060011002100000066b0020009c000500000002001d0000066b020080410000004002200210000000000121019f00000000020004140000066b0020009c0000066b02008041000000c002200210000000000112019f000006b4011001c7000080100200003919a819a30000040f00000001002001900000173d0000613d000000000101043b000700000001001d00000004010000290000066f01100197000200000001001d000000000010043f0000000701000039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f00000001002001900000173d0000613d000000000201043b0000000701000029000000000010043f000400000002001d0000000601200039000300000001001d000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f00000001002001900000173d0000613d000000000101043b000000000101041a000000000001004b0000174d0000c13d00000004010000290000000502100039000000000102041a000007030010009c0000173f0000813d000100000001001d0000000101100039000000000012041b000400000002001d000000000020043f00000000010004140000066b0010009c0000066b01008041000000c0011002100000067a011001c7000080100200003919a819a30000040f000000010020019000000007020000290000173d0000613d000000000101043b0000000101100029000000000021041b0000000401000029000000000101041a000400000001001d000000000020043f0000000301000029000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f000000010020019000000007020000290000173d0000613d000000000101043b0000000403000029000000000031041b000000000020043f0000000801000039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f00000001002001900000173d0000613d000000000801043b000000060100002900000000040104330000066f0040009c0000173f0000213d000000000108041a000000010210019000000001031002700000007f0330618f0000001f0030008c00000000010000390000000101002039000000000012004b0000000507000029000017640000c13d000000200030008c000400000008001d000700000004001d000016d00000413d000300000003001d000000000080043f00000000010004140000066b0010009c0000066b01008041000000c0011002100000067a011001c7000080100200003919a819a30000040f00000001002001900000173d0000613d00000007040000290000001f024000390000000502200270000000200040008c0000000002004019000000000301043b00000003010000290000001f01100039000000050110027000000000011300190000000002230019000000000012004b00000005070000290000000408000029000016d00000813d000000000002041b0000000102200039000000000012004b000016cc0000413d0000001f0040008c000000200a00008a000000200b000039000017000000a13d000000000080043f00000000010004140000066b0010009c0000066b01008041000000c0011002100000067a011001c7000080100200003919a819a30000040f00000001002001900000173d0000613d0000000709000029000000200a00008a0000000002a90170000000000101043b000000200b000039000017360000613d000000010320008a000000050330027000000000043100190000002003000039000000010440003900000005070000290000000606000029000000040800002900000000056300190000000005050433000000000051041b00000020033000390000000101100039000000000041004b000016ec0000c13d000000000092004b000016fd0000813d0000000302900210000000f80220018f000006fe0220027f000006fe0220016700000000036300190000000003030433000000000223016f000000000021041b000000010190021000000001011001bf0000170c0000013d000000000004004b000017040000613d0000000001070433000017050000013d000000000100001900000006060000290000000302400210000006fe0220027f000006fe02200167000000000121016f0000000102400210000000000121019f000000000018041b000000400100043d0000000003b10436000000000206043300000000002304350000004003100039000000000002004b0000171c0000613d000000000400001900000000053400190000000006740019000000000606043300000000006504350000002004400039000000000024004b000017150000413d0000001f042000390000000004a4016f0000000002230019000000000002043500000040024000390000066b0020009c0000066b0200804100000060022002100000066b0010009c0000066b010080410000004001100210000000000112019f00000000020004140000066b0020009c0000066b02008041000000c002200210000000000121019f000006b4011001c70000800d020000390000000203000039000006c804000041000000020500002919a8199e0000040f00000001002001900000173d0000613d000000000001042d00000000030b0019000000050700002900000006060000290000000408000029000000000092004b000016f50000413d000016fd0000013d0000000001000019000019aa00010430000006dc01000041000000000010043f0000004101000039000000040010043f000006bb01000041000019aa00010430000000400100043d000006ca0200004100000000002104350000066b0010009c0000066b01008041000000400110021000000674011001c7000019aa00010430000000400300043d000700000003001d000000240130003900000040020000390000000000210435000006c70100004100000000001304350000000401300039000000020200002900000000002104350000004402300039000000060100002919a8146b0000040f000000070200002900000000012100490000066b0010009c0000066b010080410000066b0020009c0000066b0200804100000060011002100000004002200210000000000121019f000019aa00010430000006dc01000041000000000010043f0000002201000039000000040010043f000006bb01000041000019aa000104300005000000000002000000400300043d000007020030009c000017b10000813d000000a002300039000000400020043f00000080023000390000000000020435000000600230003900000000000204350000004002300039000000000002043500000020023000390000000000020435000000000003043500000060021000390000000002020433000100000002001d000500000001001d0000000012010434000300000002001d000200000001001d0000000001010433000400000001001d000006c301000041000000000010044300000000010004140000066b0010009c0000066b01008041000000c001100210000006c4011001c70000800b0200003919a819a30000040f0000000100200190000017b70000613d00000004020000290000066b04200197000000000601043b000000000346004b0000000501000029000017ab0000413d00000080021000390000000002020433000006be0520019700000000023500a9000000000046004b0000179c0000613d00000000033200d9000000000053004b000017ab0000c13d0000000303000029000006be03300197000000000032001a000017ab0000413d00000000023200190000000103000029000006be03300197000006be04200197000000000023004b000000000304801900000000003104350000066b0260019700000002030000290000000000230435000000000001042d000006dc01000041000000000010043f0000001101000039000000040010043f000006bb01000041000019aa00010430000006dc01000041000000000010043f0000004101000039000000040010043f000006bb01000041000019aa00010430000000000001042f000000000010043f0000000601000039000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f0000000100200190000017ca0000613d000000000101043b000000000101041a000000000001004b0000000001000039000000010100c039000000000001042d0000000001000019000019aa000104300006000000000002000300000002001d000000000020043f000600000001001d0000000101100039000400000001001d000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f0000000100200190000018440000613d0000000603000029000000000101043b000000000101041a000000000001004b000018420000613d000000000203041a000000000002004b000018460000613d000000000021004b000500000001001d000018200000613d000200000002001d000000000030043f00000000010004140000066b0010009c0000066b01008041000000c0011002100000067a011001c7000080100200003919a819a30000040f0000000100200190000018440000613d00000005020000290001000100200092000000000101043b0000000604000029000000000204041a000000010020006c0000184c0000a13d0000000202000029000000010220008a0000000001120019000000000101041a000200000001001d000000000040043f00000000010004140000066b0010009c0000066b01008041000000c0011002100000067a011001c7000080100200003919a819a30000040f0000000100200190000018440000613d000000000101043b00000001011000290000000202000029000000000021041b000000000020043f0000000401000029000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f0000000100200190000018440000613d000000000101043b0000000502000029000000000021041b0000000603000029000000000103041a000500000001001d000000000001004b000018520000613d000000000030043f00000000010004140000066b0010009c0000066b01008041000000c0011002100000067a011001c7000080100200003919a819a30000040f0000000100200190000018440000613d0000000502000029000000010220008a000000000101043b0000000001210019000000000001041b0000000601000029000000000021041b0000000301000029000000000010043f0000000401000029000000200010043f00000000010004140000066b0010009c0000066b01008041000000c00110021000000678011001c7000080100200003919a819a30000040f0000000100200190000018440000613d000000000101043b000000000001041b0000000101000039000000000001042d0000000001000019000000000001042d0000000001000019000019aa00010430000006dc01000041000000000010043f0000001101000039000000040010043f000006bb01000041000019aa00010430000006dc01000041000000000010043f0000003201000039000000040010043f000006bb01000041000019aa00010430000006dc01000041000000000010043f0000003101000039000000040010043f000006bb01000041000019aa000104300003000000000002000100000002001d000300000001001d000000000101041a000200000001001d000006c301000041000000000010044300000000010004140000066b0010009c0000066b01008041000000c001100210000006c4011001c70000800b0200003919a819a30000040f0000000100200190000018ba0000613d000000020800002900000080028002700000066b03200197000000000201043b000000000532004b0000000307000029000018b40000413d0000000101700039000018730000c13d000000000207041a000018850000013d000000000301041a000000800630027000000000045600a900000000055400d9000000000065004b000018b40000c13d000006be05800197000000000054001a000018b40000413d000006be033001970000000004540019000000000043004b0000000003048019000006720480019700000080022002100000070402200197000000000242019f000000000232019f000000010600002900000020036000390000000004030433000006be04400197000006be05200197000000000054004b00000000050440190000070502200197000000000225019f0000000005060433000000000005004b0000000005000019000006c50500c041000000000252019f000000000027041b000000400260003900000000050204330000008005500210000000000445019f000000000041041b0000000001000039000000010100c039000000400400043d00000000011404360000000003030433000006be0330019700000000003104350000000001020433000006be01100197000000400240003900000000001204350000066b0040009c0000066b04008041000000400140021000000000020004140000066b0020009c0000066b02008041000000c002200210000000000112019f00000706011001c70000800d020000390000000103000039000007070400004119a8199e0000040f0000000100200190000018bb0000613d000000000001042d000006dc01000041000000000010043f0000001101000039000000040010043f000006bb01000041000019aa00010430000000000001042f0000000001000019000019aa000104300000066e04400197000000400510003900000000004504350000002004100039000000000034043500000000002104350000006001100039000000000001042d0006000000000002000000000401041a000006d100400198000019190000613d000000000002004b000019190000613d000600000004001d000500000002001d000200000003001d000300000001001d0000000101100039000100000001001d000000000101041a000400000001001d000006c301000041000000000010044300000000010004140000066b0010009c0000066b01008041000000c001100210000006c4011001c70000800b0200003919a819a30000040f00000001002001900000191a0000613d000000060300002900000080023002700000066b02200197000000000101043b000000000421004b000019350000413d000006be033001970000000405000029000006be02500197000018eb0000c13d00000005040000290000000305000029000018ff0000013d000000000023004b0000193d0000213d000000800650027000000000056400a900000000044500d9000000000064004b000019350000c13d000000000035001a000019350000413d0000000003350019000000800110021000000704011001970000000305000029000000000405041a0000070804400197000000000114019f000000000015041b000000000032004b00000000030240190000000504000029000000000042004b0000191b0000413d000000000143004b0000192c0000413d000006be01100197000000000205041a0000070a02200197000000000112019f000000000015041b000000400100043d00000000004104350000066b0010009c0000066b01008041000000400110021000000000020004140000066b0020009c0000066b02008041000000c002200210000000000112019f0000067a011001c70000800d0200003900000001030000390000070b0400004119a8199e0000040f00000001002001900000193b0000613d000000000001042d000000000001042f000000400100043d0000000004010019000000040110003900000002030000290000066e00300198000019450000c13d0000070f03000041000000000034043500000000002104350000002401400039000000050200002900000000002104350000066b0040009c0000066b04008041000000400140021000000676011001c7000019aa000104300000000101000029000000000101041a0000008001100272000019350000613d00000005043000690000000002140019000000010220008a000000000042004b0000194a0000813d000006dc01000041000000000010043f0000001101000039000000040010043f000006bb01000041000019aa000104300000000001000019000019aa00010430000000400100043d000007090200004100000000002104350000066b0010009c0000066b01008041000000400110021000000674011001c7000019aa000104300000070e03000041000600000004001d000000000034043500000005030000290000195e0000013d00000000021200d9000000400100043d0000000005010019000000040110003900000002040000290000066e004001980000195b0000c13d0000070d0400004100000000004504350000000000210435000000240150003900000000003104350000066b0050009c0000066b05008041000000400150021000000676011001c7000019aa000104300000070c04000041000600000005001d0000000000450435000000020400002919a818bd0000040f000000060200002900000000012100490000066b0010009c0000066b0100804100000060011002100000066b0020009c0000066b020080410000004002200210000000000121019f000019aa00010430000000000001042f0000066b0010009c0000066b0100804100000040011002100000066b0020009c0000066b020080410000006002200210000000000112019f00000000020004140000066b0020009c0000066b02008041000000c002200210000000000112019f000006b4011001c7000080100200003919a819a30000040f00000001002001900000197e0000613d000000000101043b000000000001042d0000000001000019000019aa0001043000000000050100190000000000200443000000050030008c0000198e0000413d000000040100003900000000020000190000000506200210000000000664001900000005066002700000000006060031000000000161043a0000000102200039000000000031004b000019860000413d0000066b0030009c0000066b03008041000000600130021000000000020004140000066b0020009c0000066b02008041000000c002200210000000000112019f00000710011001c7000000000205001919a819a30000040f00000001002001900000199d0000613d000000000101043b000000000001042d000000000001042f000019a1002104210000000102000039000000000001042d0000000002000019000000000001042d000019a6002104230000000102000039000000000001042d0000000002000019000000000001042d000019a800000432000019a90001042e000019aa00010430000000000000000000000000000000000000000000000000000000000000000000000000ffffffff00000000000000000000000000000000000000000000000000000001ffffffe000000000000000000000000000000000000000000000000000000000ffffffe0000000000000000000000000ffffffffffffffffffffffffffffffffffffffff000000000000000000000000000000000000000000000000ffffffffffffffff80000000000000000000000000000000000000000000000000000000000000007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0ffffffffffffffffffffffff0000000000000000000000000000000000000000313ce567000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000655a7c0e000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000044000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffffdf0200000000000000000000000000000000000040000000000000000000000000bfa87805ed57dc1f0d489ce33be4c4577d74ccde357eeeee058a32c55c44a53202000000000000000000000000000000000000200000000000000000000000002640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d8dd62ed3e00000000000000000000000000000000000000000000000000000000095ea7b300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffff3f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65641806aa1896bbf26568e884a7374b41e002500962caba6a15023a8d90e8508b8302000002000000000000000000000000000000240000000000000000000000007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6f742073756363656564000000000000000000000000000000000000000000005361666545524332303a204552433230206f7065726174696f6e20646964206e08c379a00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000840000000000000000000000000000000200000000000000000000000000000140000001000000000000000000416464726573733a2063616c6c20746f206e6f6e2d636f6e747261637400000000000000000000000000000000000000000000640000000000000000000000009b15e16f00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000009a4575b800000000000000000000000000000000000000000000000000000000c0d7865400000000000000000000000000000000000000000000000000000000dc0bd97000000000000000000000000000000000000000000000000000000000e8a1da1600000000000000000000000000000000000000000000000000000000e8a1da1700000000000000000000000000000000000000000000000000000000f2fde38b00000000000000000000000000000000000000000000000000000000dc0bd97100000000000000000000000000000000000000000000000000000000e0351e1300000000000000000000000000000000000000000000000000000000c75eea9b00000000000000000000000000000000000000000000000000000000c75eea9c00000000000000000000000000000000000000000000000000000000cf7401f300000000000000000000000000000000000000000000000000000000c0d7865500000000000000000000000000000000000000000000000000000000c4bffe2b00000000000000000000000000000000000000000000000000000000acfecf9000000000000000000000000000000000000000000000000000000000b0f479a000000000000000000000000000000000000000000000000000000000b0f479a100000000000000000000000000000000000000000000000000000000b794658000000000000000000000000000000000000000000000000000000000acfecf9100000000000000000000000000000000000000000000000000000000af58d59f000000000000000000000000000000000000000000000000000000009a4575b900000000000000000000000000000000000000000000000000000000a42a7b8b00000000000000000000000000000000000000000000000000000000a7cd63b70000000000000000000000000000000000000000000000000000000054c8a4f20000000000000000000000000000000000000000000000000000000079ba5096000000000000000000000000000000000000000000000000000000008926f54e000000000000000000000000000000000000000000000000000000008926f54f000000000000000000000000000000000000000000000000000000008da5cb5b0000000000000000000000000000000000000000000000000000000079ba5097000000000000000000000000000000000000000000000000000000007d54534e0000000000000000000000000000000000000000000000000000000054c8a4f30000000000000000000000000000000000000000000000000000000062ddd3c4000000000000000000000000000000000000000000000000000000006d3d1a5800000000000000000000000000000000000000000000000000000000240028e700000000000000000000000000000000000000000000000000000000390775360000000000000000000000000000000000000000000000000000000039077537000000000000000000000000000000000000000000000000000000004c5ef0ed00000000000000000000000000000000000000000000000000000000240028e80000000000000000000000000000000000000000000000000000000024f65ee70000000000000000000000000000000000000000000000000000000001ffc9a700000000000000000000000000000000000000000000000000000000181f5a770000000000000000000000000000000000000000000000000000000021df0da70200000000000000000000000000000000000000000000000000000000000000ed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278dad89dca000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000800000000000000000fc949c7b4a13586e39d89eead2f38644f9fb3efb5a0490b14f8fc0ceab44c2515204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d8599161e670e4b000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000024000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffff5f000000000000000000000000000000000000000000000000ffffffffffffff9f00000000000000000000000000000000ffffffffffffffffffffffffffffffff8020d12400000000000000000000000000000000000000000000000000000000d68af9cc000000000000000000000000000000000000000000000000000000001d5ad3c500000000000000000000000000000000000000000000000000000000fc949c7b4a13586e39d89eead2f38644f9fb3efb5a0490b14f8fc0ceab44c250796b89b91644bc98cd93958e4c9038275d622183e25ac5af08cc6b5d9553913202000002000000000000000000000000000000040000000000000000000000000000000000000000000000010000000000000000000000000000000000000000ffffffffffffffffffffff000000000000000000000000000000000000000000393b8ad2000000000000000000000000000000000000000000000000000000007d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c28579befe00000000000000000000000000000000000000000000000000000000310ab089e4439a4c15d089f94afb7896ff553aecb10793d0ab882de59d99a32e00000000000000000000000000000000000000200000008000000000000000008e4a23d600000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002400000140000000000000000002000000000000000000000000000000000000e00000000000000000000000000350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b0000000000000000000000ff0000000000000000000000000000000000000000036b6384b5eca791c62761152d0c79bb0604c104a5fb6f4eb0703f3154bb3db0000000000000000000000000000000000000000000000000ffffffffffffff7f00000000000000000000000000000000000000000000003fffffffffffffffe0020000000000000000000000000000000000004000000080000000000000000002dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f168400000000000000000000000000000000000000000000000000000000ffffffbffdffffffffffffffffffffffffffffffffffffc000000000000000000000000052d00ee4d9bd51b40168f2afc5848837288ce258784ad914278791464b3f4d7674f23c7c00000000000000000000000000000000000000000000000000000000405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace4e487b71000000000000000000000000000000000000000000000000000000000200000200000000000000000000000000000044000000000000000000000000961c9a4f000000000000000000000000000000000000000000000000000000002cbc26bb000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffffff0000000000000000000000000000000053ad11d800000000000000000000000000000000000000000000000000000000d0d2597600000000000000000000000000000000000000000000000000000000a8d87a3b00000000000000000000000000000000000000000000000000000000728fe07b0000000000000000000000000000000000000000000000000000000079cc679000000000000000000000000000000000000000000000000000000000696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df7a9902c7e000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000020000000000000000000000000000000000002000000080000000000000000044676b5284b809a22248eba0da87391d79098be38bb03154be88a58bf4d0917402b543c6000000000000000000000000000000000000000000000000000000008be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e02b5c74de00000000000000000000000000000000000000000000000000000000bfa87805ed57dc1f0d489ce33be4c4577d74ccde357eeeee058a32c55c44a533800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf756635f4a7b30000000000000000000000000000000000000000000000000000000083826b2b00000000000000000000000000000000000000000000000000000000a9cb113d0000000000000000000000000000000000000000000000000000000040c10f19000000000000000000000000000000000000000000000000000000009d228d69b5fdb8d273a2336f8fb8612d039631024ea9bf09c424a9503aa078f0953576f70000000000000000000000000000000000000000000000000000000024eb47e5000000000000000000000000000000000000000000000000000000004275726e46726f6d4d696e74546f6b656e506f6f6c20312e352e3100000000000000000000000000000000000000000000000000000000c0000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffff01ffc9a7000000000000000000000000000000000000000000000000000000000e64dd2900000000000000000000000000000000000000000000000000000000aff2afbf00000000000000000000000000000000000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000000000000000000000000000000000000000000000ffffffffffffffe0000000000000000000000000000000000000000000000000ffffffffffffffc0000000000000000000000000000000000000000000000000ffffffffffffff600000000000000000000000000000000000000000000000010000000000000000000000000000000000000000ffffffff00000000000000000000000000000000ffffffffffffffffffffff00ffffffff0000000000000000000000000000000002000000000000000000000000000000000000600000000000000000000000009ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c19ffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff9725942a00000000000000000000000000000000000000000000000000000000ffffffffffffffffffffffffffffffff000000000000000000000000000000001871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690ad0c8d23a0000000000000000000000000000000000000000000000000000000015279c08000000000000000000000000000000000000000000000000000000001a76572a00000000000000000000000000000000000000000000000000000000f94ebcd1000000000000000000000000000000000000000000000000000000000200000200000000000000000000000000000000000000000000000000000000")
