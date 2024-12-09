package lock_release_token_pool

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

var LockReleaseTokenPoolMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"localTokenDecimals\",\"type\":\"uint8\"},{\"internalType\":\"address[]\",\"name\":\"allowlist\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"acceptLiquidity\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"}],\"name\":\"AggregateValueMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"}],\"name\":\"AggregateValueRateLimitReached\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AllowListNotEnabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BucketOverfilled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"CallerIsNotARampOnRouter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotTransferToSelf\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CursedByRMN\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"DisabledNonZeroRateLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientLiquidity\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"expected\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"actual\",\"type\":\"uint8\"}],\"name\":\"InvalidDecimalArgs\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"InvalidRateLimitRate\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolData\",\"type\":\"bytes\"}],\"name\":\"InvalidRemoteChainDecimals\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidRemotePoolForChain\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidSourcePoolAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"InvalidToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LiquidityNotAccepted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeProposedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"NonExistentChain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"remoteDecimals\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"localDecimals\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"remoteAmount\",\"type\":\"uint256\"}],\"name\":\"OverflowDetected\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OwnerCannotBeZero\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"PoolAlreadyAdded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RateLimitMustBeDisabled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SenderNotAllowed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenRateLimitReached\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListAdd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListRemove\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Burned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remoteToken\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"ChainAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"ChainConfigured\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"LiquidityAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"LiquidityRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"LiquidityTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Locked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Minted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"rateLimitAdmin\",\"type\":\"address\"}],\"name\":\"RateLimitAdminSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Released\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"RemotePoolAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"RemotePoolRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldRouter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"RouterUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"TokensConsumed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"addRemotePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"removes\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"adds\",\"type\":\"address[]\"}],\"name\":\"applyAllowListUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64[]\",\"name\":\"remoteChainSelectorsToRemove\",\"type\":\"uint64[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes[]\",\"name\":\"remotePoolAddresses\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"remoteTokenAddress\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"internalType\":\"structTokenPool.ChainUpdate[]\",\"name\":\"chainsToAdd\",\"type\":\"tuple[]\"}],\"name\":\"applyChainUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"canAcceptLiquidity\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowList\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowListEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getCurrentInboundRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getCurrentOutboundRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRateLimitAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRebalancer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getRemotePools\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"\",\"type\":\"bytes[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getRemoteToken\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRmnProxy\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRouter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSupportedChains\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTokenDecimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"decimals\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"isRemotePool\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"isSupportedChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"isSupportedToken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"originalSender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"}],\"internalType\":\"structPool.LockOrBurnInV1\",\"name\":\"lockOrBurnIn\",\"type\":\"tuple\"}],\"name\":\"lockOrBurn\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"destPoolData\",\"type\":\"bytes\"}],\"internalType\":\"structPool.LockOrBurnOutV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"provideLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"originalSender\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"sourcePoolData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainTokenData\",\"type\":\"bytes\"}],\"internalType\":\"structPool.ReleaseOrMintInV1\",\"name\":\"releaseOrMintIn\",\"type\":\"tuple\"}],\"name\":\"releaseOrMint\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"destinationAmount\",\"type\":\"uint256\"}],\"internalType\":\"structPool.ReleaseOrMintOutV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"removeRemotePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundConfig\",\"type\":\"tuple\"}],\"name\":\"setChainRateLimiterConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"rateLimitAdmin\",\"type\":\"address\"}],\"name\":\"setRateLimitAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"rebalancer\",\"type\":\"address\"}],\"name\":\"setRebalancer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"setRouter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101206040523480156200001257600080fd5b506040516200511f3803806200511f8339810160408190526200003591620005bb565b8585858584336000816200005c57604051639b15e16f60e01b815260040160405180910390fd5b600180546001600160a01b0319166001600160a01b03848116919091179091558116156200008f576200008f81620001f3565b50506001600160a01b0385161580620000af57506001600160a01b038116155b80620000c257506001600160a01b038216155b15620000e1576040516342bcdf7f60e11b815260040160405180910390fd5b6001600160a01b03808616608081905290831660c0526040805163313ce56760e01b8152905163313ce567916004808201926020929091908290030181865afa92505050801562000151575060408051601f3d908101601f191682019092526200014e91810190620006ee565b60015b1562000191578060ff168560ff16146200018f576040516332ad3e0760e11b815260ff80871660048301528216602482015260440160405180910390fd5b505b60ff841660a052600480546001600160a01b0319166001600160a01b038316179055825115801560e052620001db57604080516000815260208101909152620001db90846200026d565b5050505091151561010052506200075a945050505050565b336001600160a01b038216036200021d57604051636d6c4ee560e11b815260040160405180910390fd5b600080546001600160a01b0319166001600160a01b03838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60e0516200028e576040516335f4a7b360e01b815260040160405180910390fd5b60005b825181101562000319576000838281518110620002b257620002b26200070c565b60209081029190910101519050620002cc600282620003ca565b156200030f576040516001600160a01b03821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b5060010162000291565b5060005b8151811015620003c55760008282815181106200033e576200033e6200070c565b6020026020010151905060006001600160a01b0316816001600160a01b0316036200036a5750620003bc565b62000377600282620003ea565b15620003ba576040516001600160a01b03821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b6001016200031d565b505050565b6000620003e1836001600160a01b03841662000401565b90505b92915050565b6000620003e1836001600160a01b03841662000505565b60008181526001830160205260408120548015620004fa5760006200042860018362000722565b85549091506000906200043e9060019062000722565b9050808214620004aa5760008660000182815481106200046257620004626200070c565b90600052602060002001549050808760000184815481106200048857620004886200070c565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080620004be57620004be62000744565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050620003e4565b6000915050620003e4565b60008181526001830160205260408120546200054e57508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155620003e4565b506000620003e4565b6001600160a01b03811681146200056d57600080fd5b50565b805160ff811681146200058257600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b8051620005828162000557565b805180151581146200058257600080fd5b60008060008060008060c08789031215620005d557600080fd5b8651620005e28162000557565b95506020620005f388820162000570565b60408901519096506001600160401b03808211156200061157600080fd5b818a0191508a601f8301126200062657600080fd5b8151818111156200063b576200063b62000587565b8060051b604051601f19603f8301168101818110858211171562000663576200066362000587565b60405291825284820192508381018501918d8311156200068257600080fd5b938501935b82851015620006ab576200069b856200059d565b8452938501939285019262000687565b809950505050505050620006c2606088016200059d565b9250620006d260808801620005aa565b9150620006e260a088016200059d565b90509295509295509295565b6000602082840312156200070157600080fd5b620003e18262000570565b634e487b7160e01b600052603260045260246000fd5b81810381811115620003e457634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052603160045260246000fd5b60805160a05160c05160e051610100516148f46200082b600039600081816105a40152611ac601526000818161063e015281816123180152612e3101526000818161061801528181611e35015261260401526000818161036701528181610e6b01528181611fde01528181612098015281816120cc015281816120ff01528181612164015281816121bd015261225f0152600081816102ce015281816103230152818161077f015281816108510152818161094501528181611b8801528181612dc7015261301c01526148f46000f3fe608060405234801561001057600080fd5b50600436106102415760003560e01c80638da5cb5b11610145578063c0d78655116100bd578063dc0bd9711161008c578063e8a1da1711610071578063e8a1da1714610662578063eb521a4c14610675578063f2fde38b1461068857600080fd5b8063dc0bd97114610616578063e0351e131461063c57600080fd5b8063c0d78655146105c8578063c4bffe2b146105db578063c75eea9c146105f0578063cf7401f31461060357600080fd5b8063acfecf9111610114578063b0f479a1116100f9578063b0f479a114610571578063b79465801461058f578063bb98546b146105a257600080fd5b8063acfecf91146104ef578063af58d59f1461050257600080fd5b80638da5cb5b1461047c5780639a4575b91461049a578063a42a7b8b146104ba578063a7cd63b7146104da57600080fd5b80634c5ef0ed116101d85780636cfd1553116101a757806379ba50971161018c57806379ba50971461044e5780637d54534e146104565780638926f54f1461046957600080fd5b80636cfd15531461041d5780636d3d1a581461043057600080fd5b80634c5ef0ed146103d157806354c8a4f3146103e457806362ddd3c4146103f7578063663200871461040a57600080fd5b8063240028e811610214578063240028e81461031357806324f65ee7146103605780633907753714610391578063432a6ba3146103b357600080fd5b806301ffc9a7146102465780630a861f2a1461026e578063181f5a771461028357806321df0da7146102cc575b600080fd5b6102596102543660046139e3565b61069b565b60405190151581526020015b60405180910390f35b61028161027c366004613a25565b6106f7565b005b6102bf6040518060400160405280601a81526020017f4c6f636b52656c65617365546f6b656e506f6f6c20312e352e3100000000000081525081565b6040516102659190613aac565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610265565b610259610321366004613ae1565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff90811691161490565b60405160ff7f0000000000000000000000000000000000000000000000000000000000000000168152602001610265565b6103a461039f366004613afe565b6108a8565b60405190518152602001610265565b600a5473ffffffffffffffffffffffffffffffffffffffff166102ee565b6102596103df366004613b57565b6109f6565b6102816103f2366004613c26565b610a40565b610281610405366004613b57565b610abb565b610281610418366004613c92565b610b53565b61028161042b366004613ae1565b610c2f565b60095473ffffffffffffffffffffffffffffffffffffffff166102ee565b610281610c7e565b610281610464366004613ae1565b610d4c565b610259610477366004613cbe565b610dcd565b60015473ffffffffffffffffffffffffffffffffffffffff166102ee565b6104ad6104a8366004613cd9565b610de4565b6040516102659190613d14565b6104cd6104c8366004613cbe565b610eb0565b6040516102659190613d6b565b6104e261101b565b6040516102659190613ded565b6102816104fd366004613b57565b61102c565b610515610510366004613cbe565b611144565b604051610265919081516fffffffffffffffffffffffffffffffff908116825260208084015163ffffffff1690830152604080840151151590830152606080840151821690830152608092830151169181019190915260a00190565b60045473ffffffffffffffffffffffffffffffffffffffff166102ee565b6102bf61059d366004613cbe565b611219565b7f0000000000000000000000000000000000000000000000000000000000000000610259565b6102816105d6366004613ae1565b6112c9565b6105e36113a4565b6040516102659190613e47565b6105156105fe366004613cbe565b61145c565b610281610611366004613fcf565b61152e565b7f00000000000000000000000000000000000000000000000000000000000000006102ee565b7f0000000000000000000000000000000000000000000000000000000000000000610259565b610281610670366004613c26565b6115b2565b610281610683366004613a25565b611ac4565b610281610696366004613ae1565b611be0565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167fe1d405660000000000000000000000000000000000000000000000000000000014806106f157506106f182611bf4565b92915050565b600a5473ffffffffffffffffffffffffffffffffffffffff16331461074f576040517f8e4a23d60000000000000000000000000000000000000000000000000000000081523360048201526024015b60405180910390fd5b6040517f70a0823100000000000000000000000000000000000000000000000000000000815230600482015281907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa1580156107db573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107ff9190614014565b1015610837576040517fbb55fd2700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61087873ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000163383611cd8565b604051819033907fc2c3f06e49b9f15e7b4af9055e183b0d73362e033ad82a07dec9bf984017171990600090a350565b6040805160208101909152600081526108c082611dac565b600061091960608401356109146108da60c087018761402d565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250611fd092505050565b612094565b905061096c61092e6060850160408601613ae1565b73ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169083611cd8565b61097c6060840160408501613ae1565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f2d87480f50083e2b2759522a8fdda59802650a8055e609a7772cf70c07748f52836040516109da91815260200190565b60405180910390a3604080516020810190915290815292915050565b6000610a388383604051610a0b929190614092565b604080519182900390912067ffffffffffffffff87166000908152600760205291909120600501906122a8565b949350505050565b610a486122c3565b610ab58484808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152505060408051602080880282810182019093528782529093508792508691829185019084908082843760009201919091525061231692505050565b50505050565b610ac36122c3565b610acc83610dcd565b610b0e576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84166004820152602401610746565b610b4e8383838080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506124cc92505050565b505050565b610b5b6122c3565b6040517f0a861f2a0000000000000000000000000000000000000000000000000000000081526004810182905273ffffffffffffffffffffffffffffffffffffffff831690630a861f2a90602401600060405180830381600087803b158015610bc357600080fd5b505af1158015610bd7573d6000803e3d6000fd5b505050508173ffffffffffffffffffffffffffffffffffffffff167f6fa7abcf1345d1d478e5ea0da6b5f26a90eadb0546ef15ed3833944fbfd1db6282604051610c2391815260200190565b60405180910390a25050565b610c376122c3565b600a80547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60005473ffffffffffffffffffffffffffffffffffffffff163314610ccf576040517f02b543c600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000008082163390811790935560008054909116815560405173ffffffffffffffffffffffffffffffffffffffff909216929183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610d546122c3565b600980547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040519081527f44676b5284b809a22248eba0da87391d79098be38bb03154be88a58bf4d091749060200160405180910390a150565b60006106f1600567ffffffffffffffff84166122a8565b6040805180820190915260608082526020820152610e01826125c6565b6040516060830135815233907f9f1ec8c880f76798e7b793325d625e9b60e4082a553c98f42b6cda368dd600089060200160405180910390a26040518060400160405280610e5b84602001602081019061059d9190613cbe565b8152602001610ea86040805160ff7f000000000000000000000000000000000000000000000000000000000000000016602082015260609101604051602081830303815290604052905090565b905292915050565b67ffffffffffffffff8116600090815260076020526040812060609190610ed990600501612752565b90506000815167ffffffffffffffff811115610ef757610ef7613e89565b604051908082528060200260200182016040528015610f2a57816020015b6060815260200190600190039081610f155790505b50905060005b82518110156110135760086000848381518110610f4f57610f4f6140a2565b602002602001015181526020019081526020016000208054610f70906140d1565b80601f0160208091040260200160405190810160405280929190818152602001828054610f9c906140d1565b8015610fe95780601f10610fbe57610100808354040283529160200191610fe9565b820191906000526020600020905b815481529060010190602001808311610fcc57829003601f168201915b5050505050828281518110611000576110006140a2565b6020908102919091010152600101610f30565b509392505050565b60606110276002612752565b905090565b6110346122c3565b61103d83610dcd565b61107f576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84166004820152602401610746565b6110bf8282604051611092929190614092565b604080519182900390912067ffffffffffffffff861660009081526007602052919091206005019061275f565b6110fb578282826040517f74f23c7c0000000000000000000000000000000000000000000000000000000081526004016107469392919061416d565b8267ffffffffffffffff167f52d00ee4d9bd51b40168f2afc5848837288ce258784ad914278791464b3f4d768383604051611137929190614191565b60405180910390a2505050565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915267ffffffffffffffff8216600090815260076020908152604091829020825160a08101845260028201546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff1615159482019490945260039091015480841660608301529190910490911660808201526106f19061276b565b67ffffffffffffffff81166000908152600760205260409020600401805460609190611244906140d1565b80601f0160208091040260200160405190810160405280929190818152602001828054611270906140d1565b80156112bd5780601f10611292576101008083540402835291602001916112bd565b820191906000526020600020905b8154815290600101906020018083116112a057829003601f168201915b50505050509050919050565b6112d16122c3565b73ffffffffffffffffffffffffffffffffffffffff811661131e576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6004805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527f02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f1684910160405180910390a15050565b606060006113b26005612752565b90506000815167ffffffffffffffff8111156113d0576113d0613e89565b6040519080825280602002602001820160405280156113f9578160200160208202803683370190505b50905060005b82518110156114555782818151811061141a5761141a6140a2565b6020026020010151828281518110611434576114346140a2565b67ffffffffffffffff909216602092830291909101909101526001016113ff565b5092915050565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915267ffffffffffffffff8216600090815260076020908152604091829020825160a08101845281546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff1615159482019490945260019091015480841660608301529190910490911660808201526106f19061276b565b60095473ffffffffffffffffffffffffffffffffffffffff16331480159061156e575060015473ffffffffffffffffffffffffffffffffffffffff163314155b156115a7576040517f8e4a23d6000000000000000000000000000000000000000000000000000000008152336004820152602401610746565b610b4e83838361281d565b6115ba6122c3565b60005b838110156117a75760008585838181106115d9576115d96140a2565b90506020020160208101906115ee9190613cbe565b9050611605600567ffffffffffffffff831661275f565b611647576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff82166004820152602401610746565b67ffffffffffffffff8116600090815260076020526040812061166c90600501612752565b905060005b81518110156116d8576116cf82828151811061168f5761168f6140a2565b6020026020010151600760008667ffffffffffffffff1667ffffffffffffffff16815260200190815260200160002060050161275f90919063ffffffff16565b50600101611671565b5067ffffffffffffffff8216600090815260076020526040812080547fffffffffffffffffffffff000000000000000000000000000000000000000000908116825560018201839055600282018054909116905560038101829055906117416004830182613976565b600582016000818161175382826139b0565b505060405167ffffffffffffffff871681527f5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d85991694506020019250611795915050565b60405180910390a150506001016115bd565b5060005b81811015611abd5760008383838181106117c7576117c76140a2565b90506020028101906117d991906141a5565b6117e290614271565b90506117f381606001516000612907565b61180281608001516000612907565b806040015151600003611841576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80516118599060059067ffffffffffffffff16612a44565b61189e5780516040517f1d5ad3c500000000000000000000000000000000000000000000000000000000815267ffffffffffffffff9091166004820152602401610746565b805167ffffffffffffffff16600090815260076020908152604091829020825160a08082018552606080870180518601516fffffffffffffffffffffffffffffffff90811680865263ffffffff42168689018190528351511515878b0181905284518a0151841686890181905294518b0151841660809889018190528954740100000000000000000000000000000000000000009283027fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff7001000000000000000000000000000000008087027fffffffffffffffffffffffff000000000000000000000000000000000000000094851690981788178216929092178d5592810290971760018c01558c519889018d52898e0180518d01518716808b528a8e019590955280515115158a8f018190528151909d01518716988a01899052518d0151909516979098018790526002890180549a909102999093161717909416959095179092559092029091176003820155908201516004820190611a2190826143e8565b5060005b826020015151811015611a6557611a5d836000015184602001518381518110611a5057611a506140a2565b60200260200101516124cc565b600101611a25565b507f8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c28260000151836040015184606001518560800151604051611aab9493929190614502565b60405180910390a150506001016117ab565b5050505050565b7f0000000000000000000000000000000000000000000000000000000000000000611b1b576040517fe93f8fa400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600a5473ffffffffffffffffffffffffffffffffffffffff163314611b6e576040517f8e4a23d6000000000000000000000000000000000000000000000000000000008152336004820152602401610746565b611bb073ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016333084612a50565b604051819033907fc17cea59c2955cb181b03393209566960365771dbba9dc3d510180e7cb31208890600090a350565b611be86122c3565b611bf181612aae565b50565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167faff2afbf000000000000000000000000000000000000000000000000000000001480611c8757507fffffffff0000000000000000000000000000000000000000000000000000000082167f0e64dd2900000000000000000000000000000000000000000000000000000000145b806106f157507fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a7000000000000000000000000000000000000000000000000000000001492915050565b60405173ffffffffffffffffffffffffffffffffffffffff8316602482015260448101829052610b4e9084907fa9059cbb00000000000000000000000000000000000000000000000000000000906064015b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090931692909217909152612b72565b611dbf61032160a0830160808401613ae1565b611e1e57611dd360a0820160808301613ae1565b6040517f961c9a4f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091166004820152602401610746565b73ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016632cbc26bb611e6a6040840160208501613cbe565b60405160e083901b7fffffffff0000000000000000000000000000000000000000000000000000000016815260809190911b77ffffffffffffffff00000000000000000000000000000000166004820152602401602060405180830381865afa158015611edb573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611eff919061459b565b15611f36576040517f53ad11d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611f4e611f496040830160208401613cbe565b612c7e565b611f6e611f616040830160208401613cbe565b6103df60a084018461402d565b611fb357611f7f60a082018261402d565b6040517f24eb47e5000000000000000000000000000000000000000000000000000000008152600401610746929190614191565b611bf1611fc66040830160208401613cbe565b8260600135612da4565b6000815160000361200257507f0000000000000000000000000000000000000000000000000000000000000000919050565b815160201461203f57816040517f953576f70000000000000000000000000000000000000000000000000000000081526004016107469190613aac565b6000828060200190518101906120559190614014565b905060ff8111156106f157826040517f953576f70000000000000000000000000000000000000000000000000000000081526004016107469190613aac565b60007f000000000000000000000000000000000000000000000000000000000000000060ff168260ff16036120ca5750816106f1565b7f000000000000000000000000000000000000000000000000000000000000000060ff168260ff1611156121b55760006121247f0000000000000000000000000000000000000000000000000000000000000000846145e7565b9050604d8160ff161115612198576040517fa9cb113d00000000000000000000000000000000000000000000000000000000815260ff80851660048301527f000000000000000000000000000000000000000000000000000000000000000016602482015260448101859052606401610746565b6121a381600a614720565b6121ad908561472f565b9150506106f1565b60006121e1837f00000000000000000000000000000000000000000000000000000000000000006145e7565b9050604d8160ff16118061222857506121fb81600a614720565b612225907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff61472f565b84115b15612293576040517fa9cb113d00000000000000000000000000000000000000000000000000000000815260ff80851660048301527f000000000000000000000000000000000000000000000000000000000000000016602482015260448101859052606401610746565b61229e81600a614720565b610a38908561476a565b600081815260018301602052604081205415155b9392505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314612314576040517f2b5c74de00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b7f000000000000000000000000000000000000000000000000000000000000000061236d576040517f35f4a7b300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b825181101561240357600083828151811061238d5761238d6140a2565b602002602001015190506123ab816002612deb90919063ffffffff16565b156123fa5760405173ffffffffffffffffffffffffffffffffffffffff821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b50600101612370565b5060005b8151811015610b4e576000828281518110612424576124246140a2565b60200260200101519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff160361246857506124c4565b612473600282612e0d565b156124c25760405173ffffffffffffffffffffffffffffffffffffffff821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b600101612407565b8051600003612507576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805160208083019190912067ffffffffffffffff84166000908152600790925260409091206125399060050182612a44565b6125735782826040517f393b8ad2000000000000000000000000000000000000000000000000000000008152600401610746929190614781565b600081815260086020526040902061258b83826143e8565b508267ffffffffffffffff167f7d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea836040516111379190613aac565b6125d961032160a0830160808401613ae1565b6125ed57611dd360a0820160808301613ae1565b73ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016632cbc26bb6126396040840160208501613cbe565b60405160e083901b7fffffffff0000000000000000000000000000000000000000000000000000000016815260809190911b77ffffffffffffffff00000000000000000000000000000000166004820152602401602060405180830381865afa1580156126aa573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906126ce919061459b565b15612705576040517f53ad11d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61271d6127186060830160408401613ae1565b612e2f565b6127356127306040830160208401613cbe565b612eae565b611bf16127486040830160208401613cbe565b8260600135612ffc565b606060006122bc83613040565b60006122bc838361309b565b6040805160a0810182526000808252602082018190529181018290526060810182905260808101919091526127f982606001516fffffffffffffffffffffffffffffffff1683600001516fffffffffffffffffffffffffffffffff16846020015163ffffffff16426127dd91906147a4565b85608001516fffffffffffffffffffffffffffffffff1661318e565b6fffffffffffffffffffffffffffffffff1682525063ffffffff4216602082015290565b61282683610dcd565b612868576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84166004820152602401610746565b612873826000612907565b67ffffffffffffffff8316600090815260076020526040902061289690836131b6565b6128a1816000612907565b67ffffffffffffffff831660009081526007602052604090206128c790600201826131b6565b7f0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b8383836040516128fa939291906147b7565b60405180910390a1505050565b8151156129d25781602001516fffffffffffffffffffffffffffffffff1682604001516fffffffffffffffffffffffffffffffff1610158061295d575060408201516fffffffffffffffffffffffffffffffff16155b1561299657816040517f8020d124000000000000000000000000000000000000000000000000000000008152600401610746919061483a565b80156129ce576040517f433fc33d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050565b60408201516fffffffffffffffffffffffffffffffff16151580612a0b575060208201516fffffffffffffffffffffffffffffffff1615155b156129ce57816040517fd68af9cc000000000000000000000000000000000000000000000000000000008152600401610746919061483a565b60006122bc8383613358565b60405173ffffffffffffffffffffffffffffffffffffffff80851660248301528316604482015260648101829052610ab59085907f23b872dd0000000000000000000000000000000000000000000000000000000090608401611d2a565b3373ffffffffffffffffffffffffffffffffffffffff821603612afd576040517fdad89dca00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b6000612bd4826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff166133a79092919063ffffffff16565b805190915015610b4e5780806020019051810190612bf2919061459b565b610b4e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152608401610746565b612c8781610dcd565b612cc9576040517fa9902c7e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff82166004820152602401610746565b600480546040517f83826b2b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84169281019290925233602483015273ffffffffffffffffffffffffffffffffffffffff16906383826b2b90604401602060405180830381865afa158015612d48573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612d6c919061459b565b611bf1576040517f728fe07b000000000000000000000000000000000000000000000000000000008152336004820152602401610746565b67ffffffffffffffff821660009081526007602052604090206129ce90600201827f00000000000000000000000000000000000000000000000000000000000000006133b6565b60006122bc8373ffffffffffffffffffffffffffffffffffffffff841661309b565b60006122bc8373ffffffffffffffffffffffffffffffffffffffff8416613358565b7f000000000000000000000000000000000000000000000000000000000000000015611bf157612e60600282613739565b611bf1576040517fd0d2597600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610746565b612eb781610dcd565b612ef9576040517fa9902c7e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff82166004820152602401610746565b600480546040517fa8d87a3b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84169281019290925273ffffffffffffffffffffffffffffffffffffffff169063a8d87a3b90602401602060405180830381865afa158015612f72573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612f969190614876565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614611bf1576040517f728fe07b000000000000000000000000000000000000000000000000000000008152336004820152602401610746565b67ffffffffffffffff821660009081526007602052604090206129ce90827f00000000000000000000000000000000000000000000000000000000000000006133b6565b6060816000018054806020026020016040519081016040528092919081815260200182805480156112bd57602002820191906000526020600020905b81548152602001906001019080831161307c5750505050509050919050565b600081815260018301602052604081205480156131845760006130bf6001836147a4565b85549091506000906130d3906001906147a4565b90508082146131385760008660000182815481106130f3576130f36140a2565b9060005260206000200154905080876000018481548110613116576131166140a2565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061314957613149614893565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506106f1565b60009150506106f1565b60006131ad8561319e848661476a565b6131a890876148c2565b613768565b95945050505050565b81546000906131df90700100000000000000000000000000000000900463ffffffff16426147a4565b905080156132815760018301548354613227916fffffffffffffffffffffffffffffffff8082169281169185917001000000000000000000000000000000009091041661318e565b83546fffffffffffffffffffffffffffffffff919091167fffffffffffffffffffffffff0000000000000000000000000000000000000000909116177001000000000000000000000000000000004263ffffffff16021783555b602082015183546132a7916fffffffffffffffffffffffffffffffff9081169116613768565b83548351151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00ffffffff000000000000000000000000000000009091166fffffffffffffffffffffffffffffffff92831617178455602083015160408085015183167001000000000000000000000000000000000291909216176001850155517f9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c19906128fa90849061483a565b600081815260018301602052604081205461339f575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556106f1565b5060006106f1565b6060610a38848460008561377e565b825474010000000000000000000000000000000000000000900460ff1615806133dd575081155b156133e757505050565b825460018401546fffffffffffffffffffffffffffffffff8083169291169060009061342d90700100000000000000000000000000000000900463ffffffff16426147a4565b905080156134ed578183111561346f576040517f9725942a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60018601546134a99083908590849070010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff1661318e565b86547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff167001000000000000000000000000000000004263ffffffff160217875592505b848210156135a45773ffffffffffffffffffffffffffffffffffffffff841661354c576040517ff94ebcd10000000000000000000000000000000000000000000000000000000081526004810183905260248101869052604401610746565b6040517f1a76572a000000000000000000000000000000000000000000000000000000008152600481018390526024810186905273ffffffffffffffffffffffffffffffffffffffff85166044820152606401610746565b848310156136b75760018681015470010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff169060009082906135e890826147a4565b6135f2878a6147a4565b6135fc91906148c2565b613606919061472f565b905073ffffffffffffffffffffffffffffffffffffffff861661365f576040517f15279c080000000000000000000000000000000000000000000000000000000081526004810182905260248101869052604401610746565b6040517fd0c8d23a000000000000000000000000000000000000000000000000000000008152600481018290526024810186905273ffffffffffffffffffffffffffffffffffffffff87166044820152606401610746565b6136c185846147a4565b86547fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff82161787556040518681529093507f1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a9060200160405180910390a1505050505050565b73ffffffffffffffffffffffffffffffffffffffff8116600090815260018301602052604081205415156122bc565b600081831061377757816122bc565b5090919050565b606082471015613810576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f60448201527f722063616c6c00000000000000000000000000000000000000000000000000006064820152608401610746565b6000808673ffffffffffffffffffffffffffffffffffffffff16858760405161383991906148d5565b60006040518083038185875af1925050503d8060008114613876576040519150601f19603f3d011682016040523d82523d6000602084013e61387b565b606091505b509150915061388c87838387613897565b979650505050505050565b6060831561392d5782516000036139265773ffffffffffffffffffffffffffffffffffffffff85163b613926576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610746565b5081610a38565b610a3883838151156139425781518083602001fd5b806040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107469190613aac565b508054613982906140d1565b6000825580601f10613992575050565b601f016020900490600052602060002090810190611bf191906139ca565b5080546000825590600052602060002090810190611bf191905b5b808211156139df57600081556001016139cb565b5090565b6000602082840312156139f557600080fd5b81357fffffffff00000000000000000000000000000000000000000000000000000000811681146122bc57600080fd5b600060208284031215613a3757600080fd5b5035919050565b60005b83811015613a59578181015183820152602001613a41565b50506000910152565b60008151808452613a7a816020860160208601613a3e565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006122bc6020830184613a62565b73ffffffffffffffffffffffffffffffffffffffff81168114611bf157600080fd5b600060208284031215613af357600080fd5b81356122bc81613abf565b600060208284031215613b1057600080fd5b813567ffffffffffffffff811115613b2757600080fd5b820161010081850312156122bc57600080fd5b803567ffffffffffffffff81168114613b5257600080fd5b919050565b600080600060408486031215613b6c57600080fd5b613b7584613b3a565b9250602084013567ffffffffffffffff80821115613b9257600080fd5b818601915086601f830112613ba657600080fd5b813581811115613bb557600080fd5b876020828501011115613bc757600080fd5b6020830194508093505050509250925092565b60008083601f840112613bec57600080fd5b50813567ffffffffffffffff811115613c0457600080fd5b6020830191508360208260051b8501011115613c1f57600080fd5b9250929050565b60008060008060408587031215613c3c57600080fd5b843567ffffffffffffffff80821115613c5457600080fd5b613c6088838901613bda565b90965094506020870135915080821115613c7957600080fd5b50613c8687828801613bda565b95989497509550505050565b60008060408385031215613ca557600080fd5b8235613cb081613abf565b946020939093013593505050565b600060208284031215613cd057600080fd5b6122bc82613b3a565b600060208284031215613ceb57600080fd5b813567ffffffffffffffff811115613d0257600080fd5b820160a081850312156122bc57600080fd5b602081526000825160406020840152613d306060840182613a62565b905060208401517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08483030160408501526131ad8282613a62565b600060208083016020845280855180835260408601915060408160051b87010192506020870160005b82811015613de0577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc0888603018452613dce858351613a62565b94509285019290850190600101613d94565b5092979650505050505050565b6020808252825182820181905260009190848201906040850190845b81811015613e3b57835173ffffffffffffffffffffffffffffffffffffffff1683529284019291840191600101613e09565b50909695505050505050565b6020808252825182820181905260009190848201906040850190845b81811015613e3b57835167ffffffffffffffff1683529284019291840191600101613e63565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405160a0810167ffffffffffffffff81118282101715613edb57613edb613e89565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715613f2857613f28613e89565b604052919050565b8015158114611bf157600080fd5b80356fffffffffffffffffffffffffffffffff81168114613b5257600080fd5b600060608284031215613f7057600080fd5b6040516060810181811067ffffffffffffffff82111715613f9357613f93613e89565b6040529050808235613fa481613f30565b8152613fb260208401613f3e565b6020820152613fc360408401613f3e565b60408201525092915050565b600080600060e08486031215613fe457600080fd5b613fed84613b3a565b9250613ffc8560208601613f5e565b915061400b8560808601613f5e565b90509250925092565b60006020828403121561402657600080fd5b5051919050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261406257600080fd5b83018035915067ffffffffffffffff82111561407d57600080fd5b602001915036819003821315613c1f57600080fd5b8183823760009101908152919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600181811c908216806140e557607f821691505b60208210810361411e577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b67ffffffffffffffff841681526040602082015260006131ad604083018486614124565b602081526000610a38602083018486614124565b600082357ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffee18336030181126141d957600080fd5b9190910192915050565b600082601f8301126141f457600080fd5b813567ffffffffffffffff81111561420e5761420e613e89565b61423f60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601613ee1565b81815284602083860101111561425457600080fd5b816020850160208301376000918101602001919091529392505050565b6000610120823603121561428457600080fd5b61428c613eb8565b61429583613b3a565b815260208084013567ffffffffffffffff808211156142b357600080fd5b9085019036601f8301126142c657600080fd5b8135818111156142d8576142d8613e89565b8060051b6142e7858201613ee1565b918252838101850191858101903684111561430157600080fd5b86860192505b8383101561433d5782358581111561431f5760008081fd5b61432d3689838a01016141e3565b8352509186019190860190614307565b808789015250505050604086013592508083111561435a57600080fd5b5050614368368286016141e3565b60408301525061437b3660608501613f5e565b606082015261438d3660c08501613f5e565b608082015292915050565b601f821115610b4e576000816000526020600020601f850160051c810160208610156143c15750805b601f850160051c820191505b818110156143e0578281556001016143cd565b505050505050565b815167ffffffffffffffff81111561440257614402613e89565b6144168161441084546140d1565b84614398565b602080601f83116001811461446957600084156144335750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b1785556143e0565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b828110156144b657888601518255948401946001909101908401614497565b50858210156144f257878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b600061010067ffffffffffffffff8716835280602084015261452681840187613a62565b8551151560408581019190915260208701516fffffffffffffffffffffffffffffffff90811660608701529087015116608085015291506145649050565b8251151560a083015260208301516fffffffffffffffffffffffffffffffff90811660c084015260408401511660e08301526131ad565b6000602082840312156145ad57600080fd5b81516122bc81613f30565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60ff82811682821603908111156106f1576106f16145b8565b600181815b8085111561465957817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0482111561463f5761463f6145b8565b8085161561464c57918102915b93841c9390800290614605565b509250929050565b600082614670575060016106f1565b8161467d575060006106f1565b8160018114614693576002811461469d576146b9565b60019150506106f1565b60ff8411156146ae576146ae6145b8565b50506001821b6106f1565b5060208310610133831016604e8410600b84101617156146dc575081810a6106f1565b6146e68383614600565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04821115614718576147186145b8565b029392505050565b60006122bc60ff841683614661565b600082614765577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b80820281158282048414176106f1576106f16145b8565b67ffffffffffffffff83168152604060208201526000610a386040830184613a62565b818103818111156106f1576106f16145b8565b67ffffffffffffffff8416815260e0810161480360208301858051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b82511515608083015260208301516fffffffffffffffffffffffffffffffff90811660a084015260408401511660c0830152610a38565b606081016106f182848051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b60006020828403121561488857600080fd5b81516122bc81613abf565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b808201808211156106f1576106f16145b8565b600082516141d9818460208701613a3e56fea164736f6c6343000818000a",
}

var LockReleaseTokenPoolABI = LockReleaseTokenPoolMetaData.ABI

var LockReleaseTokenPoolBin = LockReleaseTokenPoolMetaData.Bin

func DeployLockReleaseTokenPool(auth *bind.TransactOpts, backend bind.ContractBackend, token common.Address, localTokenDecimals uint8, allowlist []common.Address, rmnProxy common.Address, acceptLiquidity bool, router common.Address) (common.Address, *types.Transaction, *LockReleaseTokenPool, error) {
	parsed, err := LockReleaseTokenPoolMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LockReleaseTokenPoolBin), backend, token, localTokenDecimals, allowlist, rmnProxy, acceptLiquidity, router)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LockReleaseTokenPool{address: address, abi: *parsed, LockReleaseTokenPoolCaller: LockReleaseTokenPoolCaller{contract: contract}, LockReleaseTokenPoolTransactor: LockReleaseTokenPoolTransactor{contract: contract}, LockReleaseTokenPoolFilterer: LockReleaseTokenPoolFilterer{contract: contract}}, nil
}

type LockReleaseTokenPool struct {
	address common.Address
	abi     abi.ABI
	LockReleaseTokenPoolCaller
	LockReleaseTokenPoolTransactor
	LockReleaseTokenPoolFilterer
}

type LockReleaseTokenPoolCaller struct {
	contract *bind.BoundContract
}

type LockReleaseTokenPoolTransactor struct {
	contract *bind.BoundContract
}

type LockReleaseTokenPoolFilterer struct {
	contract *bind.BoundContract
}

type LockReleaseTokenPoolSession struct {
	Contract     *LockReleaseTokenPool
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type LockReleaseTokenPoolCallerSession struct {
	Contract *LockReleaseTokenPoolCaller
	CallOpts bind.CallOpts
}

type LockReleaseTokenPoolTransactorSession struct {
	Contract     *LockReleaseTokenPoolTransactor
	TransactOpts bind.TransactOpts
}

type LockReleaseTokenPoolRaw struct {
	Contract *LockReleaseTokenPool
}

type LockReleaseTokenPoolCallerRaw struct {
	Contract *LockReleaseTokenPoolCaller
}

type LockReleaseTokenPoolTransactorRaw struct {
	Contract *LockReleaseTokenPoolTransactor
}

func NewLockReleaseTokenPool(address common.Address, backend bind.ContractBackend) (*LockReleaseTokenPool, error) {
	abi, err := abi.JSON(strings.NewReader(LockReleaseTokenPoolABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindLockReleaseTokenPool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPool{address: address, abi: abi, LockReleaseTokenPoolCaller: LockReleaseTokenPoolCaller{contract: contract}, LockReleaseTokenPoolTransactor: LockReleaseTokenPoolTransactor{contract: contract}, LockReleaseTokenPoolFilterer: LockReleaseTokenPoolFilterer{contract: contract}}, nil
}

func NewLockReleaseTokenPoolCaller(address common.Address, caller bind.ContractCaller) (*LockReleaseTokenPoolCaller, error) {
	contract, err := bindLockReleaseTokenPool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolCaller{contract: contract}, nil
}

func NewLockReleaseTokenPoolTransactor(address common.Address, transactor bind.ContractTransactor) (*LockReleaseTokenPoolTransactor, error) {
	contract, err := bindLockReleaseTokenPool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolTransactor{contract: contract}, nil
}

func NewLockReleaseTokenPoolFilterer(address common.Address, filterer bind.ContractFilterer) (*LockReleaseTokenPoolFilterer, error) {
	contract, err := bindLockReleaseTokenPool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolFilterer{contract: contract}, nil
}

func bindLockReleaseTokenPool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LockReleaseTokenPoolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LockReleaseTokenPool.Contract.LockReleaseTokenPoolCaller.contract.Call(opts, result, method, params...)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.LockReleaseTokenPoolTransactor.contract.Transfer(opts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.LockReleaseTokenPoolTransactor.contract.Transact(opts, method, params...)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LockReleaseTokenPool.Contract.contract.Call(opts, result, method, params...)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.contract.Transfer(opts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.contract.Transact(opts, method, params...)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) CanAcceptLiquidity(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "canAcceptLiquidity")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) CanAcceptLiquidity() (bool, error) {
	return _LockReleaseTokenPool.Contract.CanAcceptLiquidity(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) CanAcceptLiquidity() (bool, error) {
	return _LockReleaseTokenPool.Contract.CanAcceptLiquidity(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) GetAllowList(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "getAllowList")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) GetAllowList() ([]common.Address, error) {
	return _LockReleaseTokenPool.Contract.GetAllowList(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) GetAllowList() ([]common.Address, error) {
	return _LockReleaseTokenPool.Contract.GetAllowList(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) GetAllowListEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "getAllowListEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) GetAllowListEnabled() (bool, error) {
	return _LockReleaseTokenPool.Contract.GetAllowListEnabled(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) GetAllowListEnabled() (bool, error) {
	return _LockReleaseTokenPool.Contract.GetAllowListEnabled(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) GetCurrentInboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "getCurrentInboundRateLimiterState", remoteChainSelector)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) GetCurrentInboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _LockReleaseTokenPool.Contract.GetCurrentInboundRateLimiterState(&_LockReleaseTokenPool.CallOpts, remoteChainSelector)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) GetCurrentInboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _LockReleaseTokenPool.Contract.GetCurrentInboundRateLimiterState(&_LockReleaseTokenPool.CallOpts, remoteChainSelector)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) GetCurrentOutboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "getCurrentOutboundRateLimiterState", remoteChainSelector)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) GetCurrentOutboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _LockReleaseTokenPool.Contract.GetCurrentOutboundRateLimiterState(&_LockReleaseTokenPool.CallOpts, remoteChainSelector)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) GetCurrentOutboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _LockReleaseTokenPool.Contract.GetCurrentOutboundRateLimiterState(&_LockReleaseTokenPool.CallOpts, remoteChainSelector)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) GetRateLimitAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "getRateLimitAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) GetRateLimitAdmin() (common.Address, error) {
	return _LockReleaseTokenPool.Contract.GetRateLimitAdmin(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) GetRateLimitAdmin() (common.Address, error) {
	return _LockReleaseTokenPool.Contract.GetRateLimitAdmin(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) GetRebalancer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "getRebalancer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) GetRebalancer() (common.Address, error) {
	return _LockReleaseTokenPool.Contract.GetRebalancer(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) GetRebalancer() (common.Address, error) {
	return _LockReleaseTokenPool.Contract.GetRebalancer(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) GetRemotePools(opts *bind.CallOpts, remoteChainSelector uint64) ([][]byte, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "getRemotePools", remoteChainSelector)

	if err != nil {
		return *new([][]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][]byte)).(*[][]byte)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) GetRemotePools(remoteChainSelector uint64) ([][]byte, error) {
	return _LockReleaseTokenPool.Contract.GetRemotePools(&_LockReleaseTokenPool.CallOpts, remoteChainSelector)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) GetRemotePools(remoteChainSelector uint64) ([][]byte, error) {
	return _LockReleaseTokenPool.Contract.GetRemotePools(&_LockReleaseTokenPool.CallOpts, remoteChainSelector)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) GetRemoteToken(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "getRemoteToken", remoteChainSelector)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) GetRemoteToken(remoteChainSelector uint64) ([]byte, error) {
	return _LockReleaseTokenPool.Contract.GetRemoteToken(&_LockReleaseTokenPool.CallOpts, remoteChainSelector)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) GetRemoteToken(remoteChainSelector uint64) ([]byte, error) {
	return _LockReleaseTokenPool.Contract.GetRemoteToken(&_LockReleaseTokenPool.CallOpts, remoteChainSelector)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) GetRmnProxy(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "getRmnProxy")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) GetRmnProxy() (common.Address, error) {
	return _LockReleaseTokenPool.Contract.GetRmnProxy(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) GetRmnProxy() (common.Address, error) {
	return _LockReleaseTokenPool.Contract.GetRmnProxy(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) GetRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "getRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) GetRouter() (common.Address, error) {
	return _LockReleaseTokenPool.Contract.GetRouter(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) GetRouter() (common.Address, error) {
	return _LockReleaseTokenPool.Contract.GetRouter(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) GetSupportedChains(opts *bind.CallOpts) ([]uint64, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "getSupportedChains")

	if err != nil {
		return *new([]uint64), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint64)).(*[]uint64)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) GetSupportedChains() ([]uint64, error) {
	return _LockReleaseTokenPool.Contract.GetSupportedChains(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) GetSupportedChains() ([]uint64, error) {
	return _LockReleaseTokenPool.Contract.GetSupportedChains(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) GetToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "getToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) GetToken() (common.Address, error) {
	return _LockReleaseTokenPool.Contract.GetToken(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) GetToken() (common.Address, error) {
	return _LockReleaseTokenPool.Contract.GetToken(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) GetTokenDecimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "getTokenDecimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) GetTokenDecimals() (uint8, error) {
	return _LockReleaseTokenPool.Contract.GetTokenDecimals(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) GetTokenDecimals() (uint8, error) {
	return _LockReleaseTokenPool.Contract.GetTokenDecimals(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) IsRemotePool(opts *bind.CallOpts, remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "isRemotePool", remoteChainSelector, remotePoolAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) IsRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	return _LockReleaseTokenPool.Contract.IsRemotePool(&_LockReleaseTokenPool.CallOpts, remoteChainSelector, remotePoolAddress)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) IsRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	return _LockReleaseTokenPool.Contract.IsRemotePool(&_LockReleaseTokenPool.CallOpts, remoteChainSelector, remotePoolAddress)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) IsSupportedChain(opts *bind.CallOpts, remoteChainSelector uint64) (bool, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "isSupportedChain", remoteChainSelector)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) IsSupportedChain(remoteChainSelector uint64) (bool, error) {
	return _LockReleaseTokenPool.Contract.IsSupportedChain(&_LockReleaseTokenPool.CallOpts, remoteChainSelector)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) IsSupportedChain(remoteChainSelector uint64) (bool, error) {
	return _LockReleaseTokenPool.Contract.IsSupportedChain(&_LockReleaseTokenPool.CallOpts, remoteChainSelector)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) IsSupportedToken(opts *bind.CallOpts, token common.Address) (bool, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "isSupportedToken", token)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) IsSupportedToken(token common.Address) (bool, error) {
	return _LockReleaseTokenPool.Contract.IsSupportedToken(&_LockReleaseTokenPool.CallOpts, token)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) IsSupportedToken(token common.Address) (bool, error) {
	return _LockReleaseTokenPool.Contract.IsSupportedToken(&_LockReleaseTokenPool.CallOpts, token)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) Owner() (common.Address, error) {
	return _LockReleaseTokenPool.Contract.Owner(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) Owner() (common.Address, error) {
	return _LockReleaseTokenPool.Contract.Owner(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _LockReleaseTokenPool.Contract.SupportsInterface(&_LockReleaseTokenPool.CallOpts, interfaceId)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _LockReleaseTokenPool.Contract.SupportsInterface(&_LockReleaseTokenPool.CallOpts, interfaceId)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LockReleaseTokenPool.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) TypeAndVersion() (string, error) {
	return _LockReleaseTokenPool.Contract.TypeAndVersion(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolCallerSession) TypeAndVersion() (string, error) {
	return _LockReleaseTokenPool.Contract.TypeAndVersion(&_LockReleaseTokenPool.CallOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LockReleaseTokenPool.contract.Transact(opts, "acceptOwnership")
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) AcceptOwnership() (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.AcceptOwnership(&_LockReleaseTokenPool.TransactOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.AcceptOwnership(&_LockReleaseTokenPool.TransactOpts)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactor) AddRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _LockReleaseTokenPool.contract.Transact(opts, "addRemotePool", remoteChainSelector, remotePoolAddress)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) AddRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.AddRemotePool(&_LockReleaseTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactorSession) AddRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.AddRemotePool(&_LockReleaseTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactor) ApplyAllowListUpdates(opts *bind.TransactOpts, removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPool.contract.Transact(opts, "applyAllowListUpdates", removes, adds)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.ApplyAllowListUpdates(&_LockReleaseTokenPool.TransactOpts, removes, adds)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactorSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.ApplyAllowListUpdates(&_LockReleaseTokenPool.TransactOpts, removes, adds)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactor) ApplyChainUpdates(opts *bind.TransactOpts, remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _LockReleaseTokenPool.contract.Transact(opts, "applyChainUpdates", remoteChainSelectorsToRemove, chainsToAdd)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) ApplyChainUpdates(remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.ApplyChainUpdates(&_LockReleaseTokenPool.TransactOpts, remoteChainSelectorsToRemove, chainsToAdd)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactorSession) ApplyChainUpdates(remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.ApplyChainUpdates(&_LockReleaseTokenPool.TransactOpts, remoteChainSelectorsToRemove, chainsToAdd)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactor) LockOrBurn(opts *bind.TransactOpts, lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _LockReleaseTokenPool.contract.Transact(opts, "lockOrBurn", lockOrBurnIn)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) LockOrBurn(lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.LockOrBurn(&_LockReleaseTokenPool.TransactOpts, lockOrBurnIn)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactorSession) LockOrBurn(lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.LockOrBurn(&_LockReleaseTokenPool.TransactOpts, lockOrBurnIn)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactor) ProvideLiquidity(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _LockReleaseTokenPool.contract.Transact(opts, "provideLiquidity", amount)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) ProvideLiquidity(amount *big.Int) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.ProvideLiquidity(&_LockReleaseTokenPool.TransactOpts, amount)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactorSession) ProvideLiquidity(amount *big.Int) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.ProvideLiquidity(&_LockReleaseTokenPool.TransactOpts, amount)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactor) ReleaseOrMint(opts *bind.TransactOpts, releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _LockReleaseTokenPool.contract.Transact(opts, "releaseOrMint", releaseOrMintIn)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) ReleaseOrMint(releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.ReleaseOrMint(&_LockReleaseTokenPool.TransactOpts, releaseOrMintIn)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactorSession) ReleaseOrMint(releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.ReleaseOrMint(&_LockReleaseTokenPool.TransactOpts, releaseOrMintIn)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactor) RemoveRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _LockReleaseTokenPool.contract.Transact(opts, "removeRemotePool", remoteChainSelector, remotePoolAddress)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) RemoveRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.RemoveRemotePool(&_LockReleaseTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactorSession) RemoveRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.RemoveRemotePool(&_LockReleaseTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactor) SetChainRateLimiterConfig(opts *bind.TransactOpts, remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _LockReleaseTokenPool.contract.Transact(opts, "setChainRateLimiterConfig", remoteChainSelector, outboundConfig, inboundConfig)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) SetChainRateLimiterConfig(remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.SetChainRateLimiterConfig(&_LockReleaseTokenPool.TransactOpts, remoteChainSelector, outboundConfig, inboundConfig)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactorSession) SetChainRateLimiterConfig(remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.SetChainRateLimiterConfig(&_LockReleaseTokenPool.TransactOpts, remoteChainSelector, outboundConfig, inboundConfig)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactor) SetRateLimitAdmin(opts *bind.TransactOpts, rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPool.contract.Transact(opts, "setRateLimitAdmin", rateLimitAdmin)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) SetRateLimitAdmin(rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.SetRateLimitAdmin(&_LockReleaseTokenPool.TransactOpts, rateLimitAdmin)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactorSession) SetRateLimitAdmin(rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.SetRateLimitAdmin(&_LockReleaseTokenPool.TransactOpts, rateLimitAdmin)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactor) SetRebalancer(opts *bind.TransactOpts, rebalancer common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPool.contract.Transact(opts, "setRebalancer", rebalancer)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) SetRebalancer(rebalancer common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.SetRebalancer(&_LockReleaseTokenPool.TransactOpts, rebalancer)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactorSession) SetRebalancer(rebalancer common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.SetRebalancer(&_LockReleaseTokenPool.TransactOpts, rebalancer)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactor) SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPool.contract.Transact(opts, "setRouter", newRouter)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.SetRouter(&_LockReleaseTokenPool.TransactOpts, newRouter)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactorSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.SetRouter(&_LockReleaseTokenPool.TransactOpts, newRouter)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactor) TransferLiquidity(opts *bind.TransactOpts, from common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LockReleaseTokenPool.contract.Transact(opts, "transferLiquidity", from, amount)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) TransferLiquidity(from common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.TransferLiquidity(&_LockReleaseTokenPool.TransactOpts, from, amount)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactorSession) TransferLiquidity(from common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.TransferLiquidity(&_LockReleaseTokenPool.TransactOpts, from, amount)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPool.contract.Transact(opts, "transferOwnership", to)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.TransferOwnership(&_LockReleaseTokenPool.TransactOpts, to)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.TransferOwnership(&_LockReleaseTokenPool.TransactOpts, to)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactor) WithdrawLiquidity(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _LockReleaseTokenPool.contract.Transact(opts, "withdrawLiquidity", amount)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolSession) WithdrawLiquidity(amount *big.Int) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.WithdrawLiquidity(&_LockReleaseTokenPool.TransactOpts, amount)
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolTransactorSession) WithdrawLiquidity(amount *big.Int) (*types.Transaction, error) {
	return _LockReleaseTokenPool.Contract.WithdrawLiquidity(&_LockReleaseTokenPool.TransactOpts, amount)
}

type LockReleaseTokenPoolAllowListAddIterator struct {
	Event *LockReleaseTokenPoolAllowListAdd

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolAllowListAddIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolAllowListAdd)
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
		it.Event = new(LockReleaseTokenPoolAllowListAdd)
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

func (it *LockReleaseTokenPoolAllowListAddIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolAllowListAddIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolAllowListAdd struct {
	Sender common.Address
	Raw    types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterAllowListAdd(opts *bind.FilterOpts) (*LockReleaseTokenPoolAllowListAddIterator, error) {

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolAllowListAddIterator{contract: _LockReleaseTokenPool.contract, event: "AllowListAdd", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAllowListAdd) (event.Subscription, error) {

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolAllowListAdd)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseAllowListAdd(log types.Log) (*LockReleaseTokenPoolAllowListAdd, error) {
	event := new(LockReleaseTokenPoolAllowListAdd)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolAllowListRemoveIterator struct {
	Event *LockReleaseTokenPoolAllowListRemove

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolAllowListRemoveIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolAllowListRemove)
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
		it.Event = new(LockReleaseTokenPoolAllowListRemove)
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

func (it *LockReleaseTokenPoolAllowListRemoveIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolAllowListRemoveIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolAllowListRemove struct {
	Sender common.Address
	Raw    types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterAllowListRemove(opts *bind.FilterOpts) (*LockReleaseTokenPoolAllowListRemoveIterator, error) {

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolAllowListRemoveIterator{contract: _LockReleaseTokenPool.contract, event: "AllowListRemove", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAllowListRemove) (event.Subscription, error) {

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolAllowListRemove)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseAllowListRemove(log types.Log) (*LockReleaseTokenPoolAllowListRemove, error) {
	event := new(LockReleaseTokenPoolAllowListRemove)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolBurnedIterator struct {
	Event *LockReleaseTokenPoolBurned

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolBurnedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolBurned)
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
		it.Event = new(LockReleaseTokenPoolBurned)
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

func (it *LockReleaseTokenPoolBurnedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolBurnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolBurned struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterBurned(opts *bind.FilterOpts, sender []common.Address) (*LockReleaseTokenPoolBurnedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "Burned", senderRule)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolBurnedIterator{contract: _LockReleaseTokenPool.contract, event: "Burned", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchBurned(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolBurned, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "Burned", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolBurned)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "Burned", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseBurned(log types.Log) (*LockReleaseTokenPoolBurned, error) {
	event := new(LockReleaseTokenPoolBurned)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "Burned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolChainAddedIterator struct {
	Event *LockReleaseTokenPoolChainAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolChainAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolChainAdded)
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
		it.Event = new(LockReleaseTokenPoolChainAdded)
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

func (it *LockReleaseTokenPoolChainAddedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolChainAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolChainAdded struct {
	RemoteChainSelector       uint64
	RemoteToken               []byte
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
	Raw                       types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterChainAdded(opts *bind.FilterOpts) (*LockReleaseTokenPoolChainAddedIterator, error) {

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "ChainAdded")
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolChainAddedIterator{contract: _LockReleaseTokenPool.contract, event: "ChainAdded", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchChainAdded(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolChainAdded) (event.Subscription, error) {

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "ChainAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolChainAdded)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "ChainAdded", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseChainAdded(log types.Log) (*LockReleaseTokenPoolChainAdded, error) {
	event := new(LockReleaseTokenPoolChainAdded)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "ChainAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolChainConfiguredIterator struct {
	Event *LockReleaseTokenPoolChainConfigured

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolChainConfiguredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolChainConfigured)
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
		it.Event = new(LockReleaseTokenPoolChainConfigured)
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

func (it *LockReleaseTokenPoolChainConfiguredIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolChainConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolChainConfigured struct {
	RemoteChainSelector       uint64
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
	Raw                       types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterChainConfigured(opts *bind.FilterOpts) (*LockReleaseTokenPoolChainConfiguredIterator, error) {

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "ChainConfigured")
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolChainConfiguredIterator{contract: _LockReleaseTokenPool.contract, event: "ChainConfigured", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchChainConfigured(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolChainConfigured) (event.Subscription, error) {

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "ChainConfigured")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolChainConfigured)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "ChainConfigured", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseChainConfigured(log types.Log) (*LockReleaseTokenPoolChainConfigured, error) {
	event := new(LockReleaseTokenPoolChainConfigured)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "ChainConfigured", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolChainRemovedIterator struct {
	Event *LockReleaseTokenPoolChainRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolChainRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolChainRemoved)
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
		it.Event = new(LockReleaseTokenPoolChainRemoved)
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

func (it *LockReleaseTokenPoolChainRemovedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolChainRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolChainRemoved struct {
	RemoteChainSelector uint64
	Raw                 types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterChainRemoved(opts *bind.FilterOpts) (*LockReleaseTokenPoolChainRemovedIterator, error) {

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "ChainRemoved")
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolChainRemovedIterator{contract: _LockReleaseTokenPool.contract, event: "ChainRemoved", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolChainRemoved) (event.Subscription, error) {

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "ChainRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolChainRemoved)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseChainRemoved(log types.Log) (*LockReleaseTokenPoolChainRemoved, error) {
	event := new(LockReleaseTokenPoolChainRemoved)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolConfigChangedIterator struct {
	Event *LockReleaseTokenPoolConfigChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolConfigChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolConfigChanged)
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
		it.Event = new(LockReleaseTokenPoolConfigChanged)
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

func (it *LockReleaseTokenPoolConfigChangedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolConfigChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolConfigChanged struct {
	Config RateLimiterConfig
	Raw    types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterConfigChanged(opts *bind.FilterOpts) (*LockReleaseTokenPoolConfigChangedIterator, error) {

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolConfigChangedIterator{contract: _LockReleaseTokenPool.contract, event: "ConfigChanged", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolConfigChanged) (event.Subscription, error) {

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolConfigChanged)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseConfigChanged(log types.Log) (*LockReleaseTokenPoolConfigChanged, error) {
	event := new(LockReleaseTokenPoolConfigChanged)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolLiquidityAddedIterator struct {
	Event *LockReleaseTokenPoolLiquidityAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolLiquidityAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolLiquidityAdded)
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
		it.Event = new(LockReleaseTokenPoolLiquidityAdded)
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

func (it *LockReleaseTokenPoolLiquidityAddedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolLiquidityAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolLiquidityAdded struct {
	Provider common.Address
	Amount   *big.Int
	Raw      types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterLiquidityAdded(opts *bind.FilterOpts, provider []common.Address, amount []*big.Int) (*LockReleaseTokenPoolLiquidityAddedIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "LiquidityAdded", providerRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolLiquidityAddedIterator{contract: _LockReleaseTokenPool.contract, event: "LiquidityAdded", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchLiquidityAdded(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolLiquidityAdded, provider []common.Address, amount []*big.Int) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "LiquidityAdded", providerRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolLiquidityAdded)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "LiquidityAdded", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseLiquidityAdded(log types.Log) (*LockReleaseTokenPoolLiquidityAdded, error) {
	event := new(LockReleaseTokenPoolLiquidityAdded)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "LiquidityAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolLiquidityRemovedIterator struct {
	Event *LockReleaseTokenPoolLiquidityRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolLiquidityRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolLiquidityRemoved)
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
		it.Event = new(LockReleaseTokenPoolLiquidityRemoved)
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

func (it *LockReleaseTokenPoolLiquidityRemovedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolLiquidityRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolLiquidityRemoved struct {
	Provider common.Address
	Amount   *big.Int
	Raw      types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterLiquidityRemoved(opts *bind.FilterOpts, provider []common.Address, amount []*big.Int) (*LockReleaseTokenPoolLiquidityRemovedIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "LiquidityRemoved", providerRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolLiquidityRemovedIterator{contract: _LockReleaseTokenPool.contract, event: "LiquidityRemoved", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchLiquidityRemoved(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolLiquidityRemoved, provider []common.Address, amount []*big.Int) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "LiquidityRemoved", providerRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolLiquidityRemoved)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "LiquidityRemoved", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseLiquidityRemoved(log types.Log) (*LockReleaseTokenPoolLiquidityRemoved, error) {
	event := new(LockReleaseTokenPoolLiquidityRemoved)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "LiquidityRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolLiquidityTransferredIterator struct {
	Event *LockReleaseTokenPoolLiquidityTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolLiquidityTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolLiquidityTransferred)
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
		it.Event = new(LockReleaseTokenPoolLiquidityTransferred)
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

func (it *LockReleaseTokenPoolLiquidityTransferredIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolLiquidityTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolLiquidityTransferred struct {
	From   common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterLiquidityTransferred(opts *bind.FilterOpts, from []common.Address) (*LockReleaseTokenPoolLiquidityTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "LiquidityTransferred", fromRule)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolLiquidityTransferredIterator{contract: _LockReleaseTokenPool.contract, event: "LiquidityTransferred", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchLiquidityTransferred(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolLiquidityTransferred, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "LiquidityTransferred", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolLiquidityTransferred)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "LiquidityTransferred", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseLiquidityTransferred(log types.Log) (*LockReleaseTokenPoolLiquidityTransferred, error) {
	event := new(LockReleaseTokenPoolLiquidityTransferred)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "LiquidityTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolLockedIterator struct {
	Event *LockReleaseTokenPoolLocked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolLockedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolLocked)
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
		it.Event = new(LockReleaseTokenPoolLocked)
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

func (it *LockReleaseTokenPoolLockedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolLockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolLocked struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterLocked(opts *bind.FilterOpts, sender []common.Address) (*LockReleaseTokenPoolLockedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "Locked", senderRule)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolLockedIterator{contract: _LockReleaseTokenPool.contract, event: "Locked", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchLocked(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolLocked, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "Locked", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolLocked)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "Locked", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseLocked(log types.Log) (*LockReleaseTokenPoolLocked, error) {
	event := new(LockReleaseTokenPoolLocked)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "Locked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolMintedIterator struct {
	Event *LockReleaseTokenPoolMinted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolMintedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolMinted)
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
		it.Event = new(LockReleaseTokenPoolMinted)
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

func (it *LockReleaseTokenPoolMintedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolMinted struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterMinted(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*LockReleaseTokenPoolMintedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "Minted", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolMintedIterator{contract: _LockReleaseTokenPool.contract, event: "Minted", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchMinted(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolMinted, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "Minted", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolMinted)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "Minted", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseMinted(log types.Log) (*LockReleaseTokenPoolMinted, error) {
	event := new(LockReleaseTokenPoolMinted)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "Minted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolOwnershipTransferRequestedIterator struct {
	Event *LockReleaseTokenPoolOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolOwnershipTransferRequested)
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
		it.Event = new(LockReleaseTokenPoolOwnershipTransferRequested)
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

func (it *LockReleaseTokenPoolOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LockReleaseTokenPoolOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolOwnershipTransferRequestedIterator{contract: _LockReleaseTokenPool.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolOwnershipTransferRequested)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseOwnershipTransferRequested(log types.Log) (*LockReleaseTokenPoolOwnershipTransferRequested, error) {
	event := new(LockReleaseTokenPoolOwnershipTransferRequested)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolOwnershipTransferredIterator struct {
	Event *LockReleaseTokenPoolOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolOwnershipTransferred)
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
		it.Event = new(LockReleaseTokenPoolOwnershipTransferred)
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

func (it *LockReleaseTokenPoolOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LockReleaseTokenPoolOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolOwnershipTransferredIterator{contract: _LockReleaseTokenPool.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolOwnershipTransferred)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseOwnershipTransferred(log types.Log) (*LockReleaseTokenPoolOwnershipTransferred, error) {
	event := new(LockReleaseTokenPoolOwnershipTransferred)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolRateLimitAdminSetIterator struct {
	Event *LockReleaseTokenPoolRateLimitAdminSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolRateLimitAdminSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolRateLimitAdminSet)
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
		it.Event = new(LockReleaseTokenPoolRateLimitAdminSet)
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

func (it *LockReleaseTokenPoolRateLimitAdminSetIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolRateLimitAdminSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolRateLimitAdminSet struct {
	RateLimitAdmin common.Address
	Raw            types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterRateLimitAdminSet(opts *bind.FilterOpts) (*LockReleaseTokenPoolRateLimitAdminSetIterator, error) {

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "RateLimitAdminSet")
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolRateLimitAdminSetIterator{contract: _LockReleaseTokenPool.contract, event: "RateLimitAdminSet", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchRateLimitAdminSet(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolRateLimitAdminSet) (event.Subscription, error) {

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "RateLimitAdminSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolRateLimitAdminSet)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "RateLimitAdminSet", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseRateLimitAdminSet(log types.Log) (*LockReleaseTokenPoolRateLimitAdminSet, error) {
	event := new(LockReleaseTokenPoolRateLimitAdminSet)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "RateLimitAdminSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolReleasedIterator struct {
	Event *LockReleaseTokenPoolReleased

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolReleasedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolReleased)
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
		it.Event = new(LockReleaseTokenPoolReleased)
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

func (it *LockReleaseTokenPoolReleasedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolReleasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolReleased struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterReleased(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*LockReleaseTokenPoolReleasedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "Released", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolReleasedIterator{contract: _LockReleaseTokenPool.contract, event: "Released", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchReleased(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolReleased, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "Released", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolReleased)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "Released", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseReleased(log types.Log) (*LockReleaseTokenPoolReleased, error) {
	event := new(LockReleaseTokenPoolReleased)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "Released", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolRemotePoolAddedIterator struct {
	Event *LockReleaseTokenPoolRemotePoolAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolRemotePoolAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolRemotePoolAdded)
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
		it.Event = new(LockReleaseTokenPoolRemotePoolAdded)
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

func (it *LockReleaseTokenPoolRemotePoolAddedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolRemotePoolAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolRemotePoolAdded struct {
	RemoteChainSelector uint64
	RemotePoolAddress   []byte
	Raw                 types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterRemotePoolAdded(opts *bind.FilterOpts, remoteChainSelector []uint64) (*LockReleaseTokenPoolRemotePoolAddedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "RemotePoolAdded", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolRemotePoolAddedIterator{contract: _LockReleaseTokenPool.contract, event: "RemotePoolAdded", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchRemotePoolAdded(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolRemotePoolAdded, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "RemotePoolAdded", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolRemotePoolAdded)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "RemotePoolAdded", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseRemotePoolAdded(log types.Log) (*LockReleaseTokenPoolRemotePoolAdded, error) {
	event := new(LockReleaseTokenPoolRemotePoolAdded)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "RemotePoolAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolRemotePoolRemovedIterator struct {
	Event *LockReleaseTokenPoolRemotePoolRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolRemotePoolRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolRemotePoolRemoved)
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
		it.Event = new(LockReleaseTokenPoolRemotePoolRemoved)
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

func (it *LockReleaseTokenPoolRemotePoolRemovedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolRemotePoolRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolRemotePoolRemoved struct {
	RemoteChainSelector uint64
	RemotePoolAddress   []byte
	Raw                 types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterRemotePoolRemoved(opts *bind.FilterOpts, remoteChainSelector []uint64) (*LockReleaseTokenPoolRemotePoolRemovedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "RemotePoolRemoved", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolRemotePoolRemovedIterator{contract: _LockReleaseTokenPool.contract, event: "RemotePoolRemoved", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchRemotePoolRemoved(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolRemotePoolRemoved, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "RemotePoolRemoved", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolRemotePoolRemoved)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "RemotePoolRemoved", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseRemotePoolRemoved(log types.Log) (*LockReleaseTokenPoolRemotePoolRemoved, error) {
	event := new(LockReleaseTokenPoolRemotePoolRemoved)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "RemotePoolRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolRouterUpdatedIterator struct {
	Event *LockReleaseTokenPoolRouterUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolRouterUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolRouterUpdated)
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
		it.Event = new(LockReleaseTokenPoolRouterUpdated)
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

func (it *LockReleaseTokenPoolRouterUpdatedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolRouterUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolRouterUpdated struct {
	OldRouter common.Address
	NewRouter common.Address
	Raw       types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterRouterUpdated(opts *bind.FilterOpts) (*LockReleaseTokenPoolRouterUpdatedIterator, error) {

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "RouterUpdated")
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolRouterUpdatedIterator{contract: _LockReleaseTokenPool.contract, event: "RouterUpdated", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolRouterUpdated) (event.Subscription, error) {

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "RouterUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolRouterUpdated)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseRouterUpdated(log types.Log) (*LockReleaseTokenPoolRouterUpdated, error) {
	event := new(LockReleaseTokenPoolRouterUpdated)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LockReleaseTokenPoolTokensConsumedIterator struct {
	Event *LockReleaseTokenPoolTokensConsumed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LockReleaseTokenPoolTokensConsumedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockReleaseTokenPoolTokensConsumed)
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
		it.Event = new(LockReleaseTokenPoolTokensConsumed)
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

func (it *LockReleaseTokenPoolTokensConsumedIterator) Error() error {
	return it.fail
}

func (it *LockReleaseTokenPoolTokensConsumedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LockReleaseTokenPoolTokensConsumed struct {
	Tokens *big.Int
	Raw    types.Log
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) FilterTokensConsumed(opts *bind.FilterOpts) (*LockReleaseTokenPoolTokensConsumedIterator, error) {

	logs, sub, err := _LockReleaseTokenPool.contract.FilterLogs(opts, "TokensConsumed")
	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPoolTokensConsumedIterator{contract: _LockReleaseTokenPool.contract, event: "TokensConsumed", logs: logs, sub: sub}, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) WatchTokensConsumed(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolTokensConsumed) (event.Subscription, error) {

	logs, sub, err := _LockReleaseTokenPool.contract.WatchLogs(opts, "TokensConsumed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LockReleaseTokenPoolTokensConsumed)
				if err := _LockReleaseTokenPool.contract.UnpackLog(event, "TokensConsumed", log); err != nil {
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

func (_LockReleaseTokenPool *LockReleaseTokenPoolFilterer) ParseTokensConsumed(log types.Log) (*LockReleaseTokenPoolTokensConsumed, error) {
	event := new(LockReleaseTokenPoolTokensConsumed)
	if err := _LockReleaseTokenPool.contract.UnpackLog(event, "TokensConsumed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_LockReleaseTokenPool *LockReleaseTokenPool) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _LockReleaseTokenPool.abi.Events["AllowListAdd"].ID:
		return _LockReleaseTokenPool.ParseAllowListAdd(log)
	case _LockReleaseTokenPool.abi.Events["AllowListRemove"].ID:
		return _LockReleaseTokenPool.ParseAllowListRemove(log)
	case _LockReleaseTokenPool.abi.Events["Burned"].ID:
		return _LockReleaseTokenPool.ParseBurned(log)
	case _LockReleaseTokenPool.abi.Events["ChainAdded"].ID:
		return _LockReleaseTokenPool.ParseChainAdded(log)
	case _LockReleaseTokenPool.abi.Events["ChainConfigured"].ID:
		return _LockReleaseTokenPool.ParseChainConfigured(log)
	case _LockReleaseTokenPool.abi.Events["ChainRemoved"].ID:
		return _LockReleaseTokenPool.ParseChainRemoved(log)
	case _LockReleaseTokenPool.abi.Events["ConfigChanged"].ID:
		return _LockReleaseTokenPool.ParseConfigChanged(log)
	case _LockReleaseTokenPool.abi.Events["LiquidityAdded"].ID:
		return _LockReleaseTokenPool.ParseLiquidityAdded(log)
	case _LockReleaseTokenPool.abi.Events["LiquidityRemoved"].ID:
		return _LockReleaseTokenPool.ParseLiquidityRemoved(log)
	case _LockReleaseTokenPool.abi.Events["LiquidityTransferred"].ID:
		return _LockReleaseTokenPool.ParseLiquidityTransferred(log)
	case _LockReleaseTokenPool.abi.Events["Locked"].ID:
		return _LockReleaseTokenPool.ParseLocked(log)
	case _LockReleaseTokenPool.abi.Events["Minted"].ID:
		return _LockReleaseTokenPool.ParseMinted(log)
	case _LockReleaseTokenPool.abi.Events["OwnershipTransferRequested"].ID:
		return _LockReleaseTokenPool.ParseOwnershipTransferRequested(log)
	case _LockReleaseTokenPool.abi.Events["OwnershipTransferred"].ID:
		return _LockReleaseTokenPool.ParseOwnershipTransferred(log)
	case _LockReleaseTokenPool.abi.Events["RateLimitAdminSet"].ID:
		return _LockReleaseTokenPool.ParseRateLimitAdminSet(log)
	case _LockReleaseTokenPool.abi.Events["Released"].ID:
		return _LockReleaseTokenPool.ParseReleased(log)
	case _LockReleaseTokenPool.abi.Events["RemotePoolAdded"].ID:
		return _LockReleaseTokenPool.ParseRemotePoolAdded(log)
	case _LockReleaseTokenPool.abi.Events["RemotePoolRemoved"].ID:
		return _LockReleaseTokenPool.ParseRemotePoolRemoved(log)
	case _LockReleaseTokenPool.abi.Events["RouterUpdated"].ID:
		return _LockReleaseTokenPool.ParseRouterUpdated(log)
	case _LockReleaseTokenPool.abi.Events["TokensConsumed"].ID:
		return _LockReleaseTokenPool.ParseTokensConsumed(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (LockReleaseTokenPoolAllowListAdd) Topic() common.Hash {
	return common.HexToHash("0x2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d8")
}

func (LockReleaseTokenPoolAllowListRemove) Topic() common.Hash {
	return common.HexToHash("0x800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf7566")
}

func (LockReleaseTokenPoolBurned) Topic() common.Hash {
	return common.HexToHash("0x696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df7")
}

func (LockReleaseTokenPoolChainAdded) Topic() common.Hash {
	return common.HexToHash("0x8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c2")
}

func (LockReleaseTokenPoolChainConfigured) Topic() common.Hash {
	return common.HexToHash("0x0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b")
}

func (LockReleaseTokenPoolChainRemoved) Topic() common.Hash {
	return common.HexToHash("0x5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d859916")
}

func (LockReleaseTokenPoolConfigChanged) Topic() common.Hash {
	return common.HexToHash("0x9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c19")
}

func (LockReleaseTokenPoolLiquidityAdded) Topic() common.Hash {
	return common.HexToHash("0xc17cea59c2955cb181b03393209566960365771dbba9dc3d510180e7cb312088")
}

func (LockReleaseTokenPoolLiquidityRemoved) Topic() common.Hash {
	return common.HexToHash("0xc2c3f06e49b9f15e7b4af9055e183b0d73362e033ad82a07dec9bf9840171719")
}

func (LockReleaseTokenPoolLiquidityTransferred) Topic() common.Hash {
	return common.HexToHash("0x6fa7abcf1345d1d478e5ea0da6b5f26a90eadb0546ef15ed3833944fbfd1db62")
}

func (LockReleaseTokenPoolLocked) Topic() common.Hash {
	return common.HexToHash("0x9f1ec8c880f76798e7b793325d625e9b60e4082a553c98f42b6cda368dd60008")
}

func (LockReleaseTokenPoolMinted) Topic() common.Hash {
	return common.HexToHash("0x9d228d69b5fdb8d273a2336f8fb8612d039631024ea9bf09c424a9503aa078f0")
}

func (LockReleaseTokenPoolOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (LockReleaseTokenPoolOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (LockReleaseTokenPoolRateLimitAdminSet) Topic() common.Hash {
	return common.HexToHash("0x44676b5284b809a22248eba0da87391d79098be38bb03154be88a58bf4d09174")
}

func (LockReleaseTokenPoolReleased) Topic() common.Hash {
	return common.HexToHash("0x2d87480f50083e2b2759522a8fdda59802650a8055e609a7772cf70c07748f52")
}

func (LockReleaseTokenPoolRemotePoolAdded) Topic() common.Hash {
	return common.HexToHash("0x7d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea")
}

func (LockReleaseTokenPoolRemotePoolRemoved) Topic() common.Hash {
	return common.HexToHash("0x52d00ee4d9bd51b40168f2afc5848837288ce258784ad914278791464b3f4d76")
}

func (LockReleaseTokenPoolRouterUpdated) Topic() common.Hash {
	return common.HexToHash("0x02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f1684")
}

func (LockReleaseTokenPoolTokensConsumed) Topic() common.Hash {
	return common.HexToHash("0x1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a")
}

func (_LockReleaseTokenPool *LockReleaseTokenPool) Address() common.Address {
	return _LockReleaseTokenPool.address
}

type LockReleaseTokenPoolInterface interface {
	CanAcceptLiquidity(opts *bind.CallOpts) (bool, error)

	GetAllowList(opts *bind.CallOpts) ([]common.Address, error)

	GetAllowListEnabled(opts *bind.CallOpts) (bool, error)

	GetCurrentInboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error)

	GetCurrentOutboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error)

	GetRateLimitAdmin(opts *bind.CallOpts) (common.Address, error)

	GetRebalancer(opts *bind.CallOpts) (common.Address, error)

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

	ProvideLiquidity(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	ReleaseOrMint(opts *bind.TransactOpts, releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error)

	RemoveRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error)

	SetChainRateLimiterConfig(opts *bind.TransactOpts, remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error)

	SetRateLimitAdmin(opts *bind.TransactOpts, rateLimitAdmin common.Address) (*types.Transaction, error)

	SetRebalancer(opts *bind.TransactOpts, rebalancer common.Address) (*types.Transaction, error)

	SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error)

	TransferLiquidity(opts *bind.TransactOpts, from common.Address, amount *big.Int) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	WithdrawLiquidity(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error)

	FilterAllowListAdd(opts *bind.FilterOpts) (*LockReleaseTokenPoolAllowListAddIterator, error)

	WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAllowListAdd) (event.Subscription, error)

	ParseAllowListAdd(log types.Log) (*LockReleaseTokenPoolAllowListAdd, error)

	FilterAllowListRemove(opts *bind.FilterOpts) (*LockReleaseTokenPoolAllowListRemoveIterator, error)

	WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolAllowListRemove) (event.Subscription, error)

	ParseAllowListRemove(log types.Log) (*LockReleaseTokenPoolAllowListRemove, error)

	FilterBurned(opts *bind.FilterOpts, sender []common.Address) (*LockReleaseTokenPoolBurnedIterator, error)

	WatchBurned(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolBurned, sender []common.Address) (event.Subscription, error)

	ParseBurned(log types.Log) (*LockReleaseTokenPoolBurned, error)

	FilterChainAdded(opts *bind.FilterOpts) (*LockReleaseTokenPoolChainAddedIterator, error)

	WatchChainAdded(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolChainAdded) (event.Subscription, error)

	ParseChainAdded(log types.Log) (*LockReleaseTokenPoolChainAdded, error)

	FilterChainConfigured(opts *bind.FilterOpts) (*LockReleaseTokenPoolChainConfiguredIterator, error)

	WatchChainConfigured(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolChainConfigured) (event.Subscription, error)

	ParseChainConfigured(log types.Log) (*LockReleaseTokenPoolChainConfigured, error)

	FilterChainRemoved(opts *bind.FilterOpts) (*LockReleaseTokenPoolChainRemovedIterator, error)

	WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolChainRemoved) (event.Subscription, error)

	ParseChainRemoved(log types.Log) (*LockReleaseTokenPoolChainRemoved, error)

	FilterConfigChanged(opts *bind.FilterOpts) (*LockReleaseTokenPoolConfigChangedIterator, error)

	WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolConfigChanged) (event.Subscription, error)

	ParseConfigChanged(log types.Log) (*LockReleaseTokenPoolConfigChanged, error)

	FilterLiquidityAdded(opts *bind.FilterOpts, provider []common.Address, amount []*big.Int) (*LockReleaseTokenPoolLiquidityAddedIterator, error)

	WatchLiquidityAdded(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolLiquidityAdded, provider []common.Address, amount []*big.Int) (event.Subscription, error)

	ParseLiquidityAdded(log types.Log) (*LockReleaseTokenPoolLiquidityAdded, error)

	FilterLiquidityRemoved(opts *bind.FilterOpts, provider []common.Address, amount []*big.Int) (*LockReleaseTokenPoolLiquidityRemovedIterator, error)

	WatchLiquidityRemoved(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolLiquidityRemoved, provider []common.Address, amount []*big.Int) (event.Subscription, error)

	ParseLiquidityRemoved(log types.Log) (*LockReleaseTokenPoolLiquidityRemoved, error)

	FilterLiquidityTransferred(opts *bind.FilterOpts, from []common.Address) (*LockReleaseTokenPoolLiquidityTransferredIterator, error)

	WatchLiquidityTransferred(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolLiquidityTransferred, from []common.Address) (event.Subscription, error)

	ParseLiquidityTransferred(log types.Log) (*LockReleaseTokenPoolLiquidityTransferred, error)

	FilterLocked(opts *bind.FilterOpts, sender []common.Address) (*LockReleaseTokenPoolLockedIterator, error)

	WatchLocked(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolLocked, sender []common.Address) (event.Subscription, error)

	ParseLocked(log types.Log) (*LockReleaseTokenPoolLocked, error)

	FilterMinted(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*LockReleaseTokenPoolMintedIterator, error)

	WatchMinted(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolMinted, sender []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseMinted(log types.Log) (*LockReleaseTokenPoolMinted, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LockReleaseTokenPoolOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*LockReleaseTokenPoolOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LockReleaseTokenPoolOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*LockReleaseTokenPoolOwnershipTransferred, error)

	FilterRateLimitAdminSet(opts *bind.FilterOpts) (*LockReleaseTokenPoolRateLimitAdminSetIterator, error)

	WatchRateLimitAdminSet(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolRateLimitAdminSet) (event.Subscription, error)

	ParseRateLimitAdminSet(log types.Log) (*LockReleaseTokenPoolRateLimitAdminSet, error)

	FilterReleased(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*LockReleaseTokenPoolReleasedIterator, error)

	WatchReleased(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolReleased, sender []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseReleased(log types.Log) (*LockReleaseTokenPoolReleased, error)

	FilterRemotePoolAdded(opts *bind.FilterOpts, remoteChainSelector []uint64) (*LockReleaseTokenPoolRemotePoolAddedIterator, error)

	WatchRemotePoolAdded(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolRemotePoolAdded, remoteChainSelector []uint64) (event.Subscription, error)

	ParseRemotePoolAdded(log types.Log) (*LockReleaseTokenPoolRemotePoolAdded, error)

	FilterRemotePoolRemoved(opts *bind.FilterOpts, remoteChainSelector []uint64) (*LockReleaseTokenPoolRemotePoolRemovedIterator, error)

	WatchRemotePoolRemoved(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolRemotePoolRemoved, remoteChainSelector []uint64) (event.Subscription, error)

	ParseRemotePoolRemoved(log types.Log) (*LockReleaseTokenPoolRemotePoolRemoved, error)

	FilterRouterUpdated(opts *bind.FilterOpts) (*LockReleaseTokenPoolRouterUpdatedIterator, error)

	WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolRouterUpdated) (event.Subscription, error)

	ParseRouterUpdated(log types.Log) (*LockReleaseTokenPoolRouterUpdated, error)

	FilterTokensConsumed(opts *bind.FilterOpts) (*LockReleaseTokenPoolTokensConsumedIterator, error)

	WatchTokensConsumed(opts *bind.WatchOpts, sink chan<- *LockReleaseTokenPoolTokensConsumed) (event.Subscription, error)

	ParseTokensConsumed(log types.Log) (*LockReleaseTokenPoolTokensConsumed, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}

var LockReleaseTokenPoolZKBin = ("0x000400000000000200210000000000020000006003100270000006cc0030019d000006cc03300197000300000031035500020000000103550000000100200190000000280000c13d0000008002000039000000400020043f000000040030008c000000500000413d000000000401043b000000e004400270000006df0040009c000000520000a13d000006e00040009c000000690000213d000006ed0040009c000000d60000a13d000006ee0040009c0000031b0000a13d000006ef0040009c000005ff0000613d000006f00040009c000005010000613d000006f10040009c000000500000c13d0000000001000416000000000001004b000000500000c13d0000000001000412001700000001001d001600800000003d000080050100003900000044030000390000000004000415000000170440008a000003540000013d0000012004000039000000400040043f0000000002000416000000000002004b000000500000c13d0000001f02300039000006cd022001970000012002200039000000400020043f0000001f0530018f000006ce0630019800000120026000390000003a0000613d000000000701034f000000007807043c0000000004840436000000000024004b000000360000c13d000000000005004b000000470000613d000000000161034f0000000304500210000000000502043300000000054501cf000000000545022f000000000101043b0000010004400089000000000141022f00000000014101cf000000000151019f0000000000120435000000c00030008c000000500000413d000001200100043d000006cf0010009c000000500000213d000001400200043d001100000002001d000000ff0020008c000000bd0000a13d000000000100001900001b2e00010430000006f90040009c000000880000a13d000006fa0040009c000001750000a13d000006fb0040009c0000035d0000a13d000006fc0040009c00000a380000613d000006fd0040009c0000054d0000613d000006fe0040009c000000500000c13d000000240030008c000000500000413d0000000002000416000000000002004b000000500000c13d0000000401100370000000000101043b000006d00010009c000000500000213d1b2c193c0000040f0000046c0000013d000006e10040009c0000012f0000a13d000006e20040009c000003460000a13d000006e30040009c000006080000613d000006e40040009c0000051a0000613d000006e50040009c000000500000c13d000000240030008c000000500000413d0000000002000416000000000002004b000000500000c13d0000000401100370000000000601043b000006cf0060009c000000500000213d0000000101000039000000000101041a000006cf011001970000000005000411000000000015004b00000a510000c13d000000000056004b00000aee0000c13d0000071301000041000000800010043f000007140100004100001b2e00010430000007060040009c0000023d0000213d0000070c0040009c000003660000213d0000070f0040009c0000056b0000613d000007100040009c000000500000c13d000000240030008c000000500000413d0000000002000416000000000002004b000000500000c13d0000000401100370000000000201043b0000000a01000039000000000101041a000006cf011001970000000003000411000000000031004b00000a920000c13d001100000002001d000007150100004100000000001004430000000001000412000000040010044300000024000004430000000001000414000006cc0010009c000006cc01008041000000c00110021000000716011001c700008005020000391b2c1b270000040f0000000100200190000013b30000613d000000000301043b000000400b00043d000007600100004100000000001b04350000000401b00039000000000200041000000000002104350000000001000414001000000003001d000006cf02300197000000040020008c00000afb0000c13d0000000103000031000000200030008c0000002004000039000000000403401900000b270000013d000001600200043d000006d00020009c000000500000213d0000001f04200039000000000034004b0000000005000019000006d105008041000006d104400197000000000004004b0000000006000019000006d106004041000006d10040009c000000000605c019000000000006004b000000500000c13d00000120042000390000000004040433000006d00040009c0000037e0000a13d0000074101000041000000000010043f0000004101000039000000040010043f000007180100004100001b2e00010430000006f40040009c000002540000213d000006f70040009c000003b10000613d000006f80040009c000000500000c13d000000240030008c000000500000413d0000000002000416000000000002004b000000500000c13d0000000402100370000000000202043b000006d00020009c000000500000213d0000000003230049000007200030009c000000500000213d000000a40030008c000000500000413d000000c003000039000000400030043f0000006003000039000000800030043f000000a00030043f001000840020003d0000001001100360000000000101043b001100000001001d000006cf0010009c000000500000213d000007150100004100000000001004430000000001000412000000040010044300000024000004430000000001000414000006cc0010009c000006cc01008041000000c00110021000000716011001c700008005020000391b2c1b270000040f0000000100200190000013b30000613d0000000202000367000000000101043b000006cf01100197000000110010006b00000b3a0000c13d0000001001000029000f0060001000920000000f01200360000000000101043b000006d00010009c000000500000213d000000400300043d0000074302000041000000000023043500000080011002100000074401100197001000000003001d000000040230003900000000001204350000071501000041000000000010044300000000010004120000000400100443000000400100003900000024001004430000000001000414000006cc0010009c000006cc01008041000000c00110021000000716011001c700008005020000391b2c1b270000040f0000000100200190000013b30000613d000000000201043b0000000001000414000006cf02200197000000040020008c00000bac0000c13d0000000103000031000000200030008c0000002004000039000000000403401900000bd60000013d000006e80040009c000002650000213d000006eb0040009c000003b60000613d000006ec0040009c000000500000c13d0000000002000416000000000002004b000000500000c13d0000000504000039000000000204041a000000800020043f000000000040043f000000000002004b00000a550000c13d000000a002000039000000400020043f000000000400001900000005064002100000003f0560003900000738055001970000000005250019000006d00050009c000000d00000213d000000400050043f00000000044204360000001f0560018f000000000006004b000001530000613d000000000131034f00000000036400190000000006040019000000001701043c0000000006760436000000000036004b0000014f0000c13d000000000005004b000000800100043d000000000001004b000001650000613d00000000010000190000000003020433000000000013004b000010ed0000a13d00000005031002100000000005340019000000a0033000390000000003030433000006d00330019700000000003504350000000101100039000000800300043d000000000031004b000001580000413d000000400100043d00000020030000390000000005310436000000000302043300000000003504350000004002100039000000000003004b00000a890000613d00000000050000190000000046040434000006d00660019700000000026204360000000105500039000000000035004b0000016e0000413d00000a890000013d000007010040009c000002c50000213d000007040040009c000004660000613d000007050040009c000000500000c13d000000440030008c000000500000413d0000000004000416000000000004004b000000500000c13d0000000404100370000000000404043b000006d00040009c000000500000213d0000002305400039000000000035004b000000500000813d0000000405400039000000000551034f000000000905043b000006d00090009c000000500000213d0000002407400039000000050b90021000000000087b0019000000000038004b000000500000213d0000002404100370000000000404043b000006d00040009c000000500000213d0000002305400039000000000035004b000000500000813d0000000405400039000000000551034f000000000605043b000006d00060009c000000500000213d0000002404400039000000050a60021000000000054a0019000000000035004b000000500000213d0000000103000039000000000303041a000006cf03300197000000000c00041100000000003c004b00000a510000c13d0000003f03b00039000006d203300197000007400030009c000000d00000213d0000008003300039000f00000003001d000000400030043f000000800090043f000000000009004b000001bd0000613d000000000371034f000000000303043b000006cf0030009c000000500000213d000000200220003900000000003204350000002007700039000000000087004b000001b20000413d000000400200043d000f00000002001d0000003f02a00039000006d2022001970000000f022000290000000f0020006c00000000030000390000000103004039000006d00020009c000000d00000213d0000000100300190000000d00000c13d000000400020043f0000000f020000290000000002620436000e00000002001d000000000006004b000001d70000613d0000000f02000029000000000341034f000000000303043b000006cf0030009c000000500000213d000000200220003900000000003204350000002004400039000000000054004b000001ce0000413d0000071501000041000000000010044300000000010004120000000400100443000000600100003900000024001004430000000001000414000006cc0010009c000006cc01008041000000c00110021000000716011001c700008005020000391b2c1b270000040f0000000100200190000013b30000613d000000000101043b000000000001004b00000d510000613d000000800100043d000000000001004b000010880000c13d0000000f010000290000000001010433000000000001004b00000b950000613d0000000003000019000001f70000013d00000001033000390000000f010000290000000001010433000000000013004b00000b950000813d00000005013002100000000e011000290000000001010433000006cf04100198000001f20000613d000000000040043f0000000301000039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c70000801002000039001000000003001d001100000004001d1b2c1b270000040f000000110400002900000010030000290000000100200190000000500000613d000000000101043b000000000101041a000000000001004b000001f20000c13d0000000203000039000000000103041a000006d00010009c000000d00000213d0000000102100039000000000023041b000006da0110009a000000000041041b000000000103041a000d00000001001d000000000040043f0000000301000039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f00000011030000290000000100200190000000500000613d000000000101043b0000000d02000029000000000021041b000000400100043d0000000000310435000006cc0010009c000006cc0100804100000040011002100000000002000414000006cc0020009c000006cc02008041000000c002200210000000000112019f000006db011001c70000800d020000390000000103000039000006dc040000411b2c1b220000040f00000010030000290000000100200190000001f20000c13d000000500000013d000007070040009c000003750000213d0000070a0040009c0000057c0000613d0000070b0040009c000000500000c13d0000000001000416000000000001004b000000500000c13d0000000001000412001b00000001001d001a00200000003d0000800501000039000000440300003900000000040004150000001b0440008a000000050440021000000715020000411b2c1b040000040f000000ff0110018f000000800010043f000007300100004100001b2d0001042e000006f50040009c000003cb0000613d000006f60040009c000000500000c13d0000000001000416000000000001004b000000500000c13d0000000202000039000000000102041a000000800010043f000000000020043f0000002002000039000000000001004b00000a6a0000c13d000000a001000039000000000402001900000a790000013d000006e90040009c000004330000613d000006ea0040009c000000500000c13d000000e40030008c000000500000413d0000000002000416000000000002004b000000500000c13d0000000402100370000000000202043b001100000002001d000006d00020009c000000500000213d000000e002000039000000400020043f0000002402100370000000000202043b000000000002004b0000000003000039000000010300c039000000000032004b000000500000c13d000000800020043f0000004402100370000000000202043b000007220020009c000000500000213d000000a00020043f0000006402100370000000000202043b000007220020009c000000500000213d000000c00020043f0000014002000039000000400020043f0000008402100370000000000202043b000000000002004b0000000003000039000000010300c039000000000032004b000000500000c13d000000e00020043f000000a402100370000000000202043b000007220020009c000000500000213d000001000020043f000000c401100370000000000101043b000007220010009c000000500000213d000001200010043f0000000901000039000000000101041a000006cf021001970000000001000411000000000021004b000002a60000613d0000000102000039000000000202041a000006cf02200197000000000021004b00000f6e0000c13d0000001101000029000000000010043f0000000601000039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d000000000101043b000000000101041a000000000001004b000004e50000613d000000c00100043d0000072201100197000000800200043d000000000002004b00000f870000c13d000000000001004b000002c10000c13d000000a00100043d000007220010019800000f8d0000613d000000400200043d001100000002001d0000072501000041000010090000013d000007020040009c000004760000613d000007030040009c000000500000c13d000000440030008c000000500000413d0000000002000416000000000002004b000000500000c13d0000000402100370000000000202043b001100000002001d000006cf0020009c000000500000213d0000002401100370000000000301043b0000000101000039000000000101041a000006cf011001970000000002000411000000000012004b00000a510000c13d001000000003001d00000750010000410000000000100443000000110100002900000004001004430000000001000414000006cc0010009c000006cc01008041000000c00110021000000751011001c700008002020000391b2c1b270000040f0000000100200190000013b30000613d000000000101043b000000000001004b000000500000613d000000400500043d0000075201000041000000000015043500000004015000390000001004000029000000000041043500000000010004140000001102000029000000040020008c000003090000613d000006cc0050009c000006cc0200004100000000020540190000004002200210000006cc0010009c000006cc01008041000000c001100210000000000121019f00000718011001c70000001102000029000f00000005001d1b2c1b220000040f0000000f0500002900000010040000290000006003100270000106cc0030019d0003000000010355000000010020019000000b970000613d000006d00050009c000000d00000213d000000400050043f0000000000450435000006cc0050009c000006cc0500804100000040015002100000000002000414000006cc0020009c000006cc02008041000000c002200210000000000112019f000006db011001c70000800d0200003900000002030000390000075304000041000000110500002900000b920000013d000006f20040009c000004b20000613d000006f30040009c000000500000c13d000000240030008c000000500000413d0000000002000416000000000002004b000000500000c13d0000000401100370000000000101043b001100000001001d000006d00010009c000000500000213d1b2c15740000040f0000001101000029000000000010043f0000000701000039000000200010043f000000400200003900000000010000191b2c1aef0000040f001000000001001d000000400100043d001100000001001d1b2c142f0000040f00000010050000290000000201500039000000000401041a00000734004001980000000002000039000000010200c0390000001101000029000000400310003900000000002304350000008002400270000006cc0220019700000020031000390000000000230435000007220240019700000000002104350000000302500039000004580000013d000006e60040009c000004f00000613d000006e70040009c000000500000c13d0000000001000416000000000001004b000000500000c13d0000000001000412001300000001001d001200600000003d000080050100003900000044030000390000000004000415000000130440008a000000050440021000000715020000411b2c1b040000040f000000000001004b0000000001000039000000010100c039000000800010043f000007300100004100001b2d0001042e000006ff0040009c0000053b0000613d000007000040009c000000500000c13d0000000001000416000000000001004b000000500000c13d0000000901000039000006030000013d0000070d0040009c000005970000613d0000070e0040009c000000500000c13d0000000001000416000000000001004b000000500000c13d0000000001000412001f00000001001d001e00000000003d0000800501000039000000440300003900000000040004150000001f0440008a000004fa0000013d000007080040009c000005ab0000613d000007090040009c000000500000c13d0000000001000416000000000001004b000000500000c13d0000000a01000039000006030000013d00000005054002100000003f06500039000006d206600197000000400700043d0000000006670019000f00000007001d000000000076004b00000000070000390000000107004039000006d00060009c000000d00000213d0000000100700190000000d00000c13d0000012007300039000000400060043f0000000f030000290000000003430436000e00000003001d00000140022000390000000003250019000000000073004b000000500000213d000000000004004b0000039d0000613d0000000e040000290000000025020434000006cf0050009c000000500000213d0000000004540436000000000032004b000003970000413d000001800300043d000006cf0030009c000000500000213d000001a00400043d000000000004004b0000000002000039000000010200c039000c00000004001d000000000024004b000000500000c13d000001c00200043d001000000002001d000006cf0020009c000000500000213d0000000002000411000000000002004b00000b690000c13d000000400100043d000006de0200004100000ba60000013d0000000001000416000000000001004b000000500000c13d0000000101000039000006030000013d000000240030008c000000500000413d0000000002000416000000000002004b000000500000c13d0000000401100370000000000101043b000006cf0010009c000000500000213d0000000102000039000000000202041a000006cf022001970000000003000411000000000023004b00000a510000c13d000000000001004b00000ad40000c13d0000072f01000041000000800010043f000007140100004100001b2e00010430000000240030008c000000500000413d0000000002000416000000000002004b000000500000c13d0000000401100370000000000101043b000006d00010009c000000500000213d000000000010043f0000000701000039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d000000000101043b0000000501100039000000000301041a000000400200043d000f00000002001d001100000003001d0000000002320436000e00000002001d000000000010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006db011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d0000001105000029000000000005004b0000000e02000029000003ff0000613d000000000101043b0000000e020000290000000003000019000000000401041a000000000242043600000001011000390000000103300039000000000053004b000003f90000413d0000000f0120006a0000001f0110003900000769011001970000000f04100029000000000014004b00000000010000390000000101004039000006d00040009c000000d00000213d0000000100100190000000d00000c13d000000400040043f0000000f010000290000000002010433000006d00020009c000000d00000213d00000005012002100000003f03100039000006d2033001970000000003430019000006d00030009c000000d00000213d000000400030043f000d00000004001d0000000005240436000000000002004b000004210000613d00000060020000390000000003000019000000000435001900000000002404350000002003300039000000000013004b0000041c0000413d000c00000005001d0000000f010000290000000001010433000000000001004b00000c240000c13d000000400100043d000000200200003900000000032104360000000d0200002900000000020204330000000000230435000000400310003900000005042002100000000005340019000000000002004b00000d540000c13d000000000215004900000a8a0000013d000000240030008c000000500000413d0000000002000416000000000002004b000000500000c13d0000000401100370000000000101043b001100000001001d000006d00010009c000000500000213d1b2c15740000040f0000001101000029000000000010043f0000000701000039000000200010043f000000400200003900000000010000191b2c1aef0000040f001000000001001d000000400100043d001100000001001d1b2c142f0000040f0000001005000029000000000405041a00000734004001980000000002000039000000010200c0390000001101000029000000400310003900000000002304350000008002400270000006cc0220019700000020031000390000000000230435000007220240019700000000002104350000000102500039000000000402041a0000008002100039000000800340027000000000003204350000072203400197000000600210003900000000003204351b2c18220000040f000000400100043d001000000001001d00000011020000291b2c14690000040f0000001002000029000005110000013d0000000001000416000000000001004b000000500000c13d00000000010300191b2c144c0000040f1b2c14e40000040f000000000001004b0000000001000039000000010100c039000000400200043d0000000000120435000006cc0020009c000006cc0200804100000040012002100000074b011001c700001b2d0001042e000000440030008c000000500000413d0000000002000416000000000002004b000000500000c13d0000000402100370000000000202043b001100000002001d000006d00020009c000000500000213d0000002402100370000000000202043b000006d00020009c000000500000213d0000002304200039000000000034004b000000500000813d0000000404200039000000000141034f000000000101043b001000000001001d000006d00010009c000000500000213d0000002402200039000f00000002001d0000001001200029000000000031004b000000500000213d0000000101000039000000000101041a000006cf011001970000000002000411000000000012004b00000a510000c13d0000001101000029000000000010043f0000000601000039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d000000000101043b000000000101041a000000000001004b000004e50000613d00000000030000310000000f0100002900000010020000291b2c14ac0000040f000000000201001900000011010000291b2c16e40000040f000000000100001900001b2d0001042e000000440030008c000000500000413d0000000002000416000000000002004b000000500000c13d0000000402100370000000000202043b001100000002001d000006d00020009c000000500000213d0000002402100370000000000202043b000006d00020009c000000500000213d0000002304200039000000000034004b000000500000813d000f00040020003d0000000f01100360000000000101043b001000000001001d000006d00010009c000000500000213d0000002402200039000d00000002001d000e00100020002d0000000e0030006b000000500000213d0000000101000039000000000101041a000006cf011001970000000002000411000000000012004b00000a510000c13d0000001101000029000000000010043f0000000601000039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d000000000101043b000000000101041a000000000001004b00000c880000c13d000000400100043d0000071f020000410000000000210435000000040210003900000011030000290000000000320435000006cc0010009c000006cc01008041000000400110021000000718011001c700001b2e000104300000000001000416000000000001004b000000500000c13d0000000001000412001500000001001d001400400000003d000080050100003900000044030000390000000004000415000000150440008a000000050440021000000715020000411b2c1b040000040f000006cf01100197000000800010043f000007300100004100001b2d0001042e000000240030008c000000500000413d0000000002000416000000000002004b000000500000c13d0000000401100370000000000101043b000006d00010009c000000500000213d1b2c15890000040f0000002002000039000000400300043d001100000003001d00000000022304361b2c143a0000040f00000011020000290000000001210049000006cc0010009c000006cc010080410000006001100210000006cc0020009c000006cc020080410000004002200210000000000121019f00001b2d0001042e000000240030008c000000500000413d0000000003000416000000000003004b000000500000c13d0000000401100370000000000101043b001100000001001d000007150100004100000000001004430000000001000412000000040010044300000024002004430000000001000414000006cc0010009c000006cc01008041000000c00110021000000716011001c700008005020000391b2c1b270000040f0000000100200190000013b30000613d000000400300043d000000000101043b000000000001004b00000a9b0000c13d0000071c010000410000000000130435000006cc0030009c000006cc030080410000004001300210000006d5011001c700001b2e00010430000000240030008c000000500000413d0000000002000416000000000002004b000000500000c13d0000000401100370000000000101043b001100000001001d000006cf0010009c000000500000213d1b2c16d50000040f0000000a01000039000000000201041a000006d30220019700000011022001af000000000021041b000000000100001900001b2d0001042e000000240030008c000000500000413d0000000002000416000000000002004b000000500000c13d0000000401100370000000000101043b000006cf0010009c000000500000213d0000000102000039000000000202041a000006cf022001970000000003000411000000000023004b00000a510000c13d0000000902000039000000000302041a000006d303300197000000000313019f000000000032041b000000800010043f0000000001000414000006cc0010009c000006cc01008041000000c0011002100000074c011001c70000800d0200003900000001030000390000074d0400004100000ae40000013d000000240030008c000000500000413d0000000002000416000000000002004b000000500000c13d0000000401100370000000000201043b0000076300200198000000500000c13d0000000101000039000007640020009c00000acd0000213d000007670020009c000006050000613d000007680020009c000006050000613d00000ad10000013d000000240030008c000000500000413d0000000002000416000000000002004b000000500000c13d0000000401100370000000000101043b001100000001001d000006cf0010009c000000500000213d0000000001000412001d00000001001d001c00000000003d0000800501000039000000440300003900000000040004150000001d0440008a000000050440021000000715020000411b2c1b040000040f000006cf01100197000000110010006b00000000010000390000000101006039000000800010043f000007300100004100001b2d0001042e0000000001000416000000000001004b000000500000c13d000000c001000039000000400010043f0000001a01000039000000800010043f0000075d01000041000000a00010043f0000002001000039000000c00010043f0000008001000039000000e0020000391b2c143a0000040f000000c00110008a000006cc0010009c000006cc0100804100000060011002100000075e011001c700001b2d0001042e000000240030008c000000500000413d0000000002000416000000000002004b000000500000c13d0000000402100370000000000202043b001100000002001d000006d00020009c000000500000213d000000110230006a000007200020009c000000500000213d000001040020008c000000500000413d000000a002000039000000400020043f0000001102000029000f00840020003d0000000f01100360000000800000043f000000000101043b001000000001001d000006cf0010009c000000500000213d000007150100004100000000001004430000000001000412000000040010044300000024000004430000000001000414000006cc0010009c000006cc01008041000000c00110021000000716011001c700008005020000391b2c1b270000040f0000000100200190000013b30000613d0000000202000367000000000101043b000e00000001001d000006cf01100197000000100010006b00000b5a0000c13d0000000f01000029000d0060001000920000000d01200360000000000101043b000006d00010009c000000500000213d000000400300043d0000074302000041000000000023043500000080011002100000074401100197000f00000003001d000000040230003900000000001204350000071501000041000000000010044300000000010004120000000400100443000000400100003900000024001004430000000001000414000006cc0010009c000006cc01008041000000c00110021000000716011001c700008005020000391b2c1b270000040f0000000100200190000013b30000613d000000000201043b0000000001000414000006cf02200197000000040020008c00000d7d0000c13d0000000103000031000000200030008c0000002004000039000000000403401900000da70000013d0000000001000416000000000001004b000000500000c13d0000000401000039000000000101041a000006cf01100197000000800010043f000007300100004100001b2d0001042e000000440030008c000000500000413d0000000002000416000000000002004b000000500000c13d0000000402100370000000000202043b000006d00020009c000000500000213d0000002304200039000000000034004b000000500000813d0000000404200039000000000441034f000000000404043b000600000004001d000006d00040009c000000500000213d000500240020003d000000060200002900000005022002100000000502200029000000000032004b000000500000213d0000002402100370000000000202043b000200000002001d000006d00020009c000000500000213d00000002020000290000002302200039000000000032004b000000500000813d00000002020000290000000402200039000000000121034f000000000101043b000100000001001d000006d00010009c000000500000213d0000000201000029000300240010003d000000010100002900000005011002100000000301100029000000000031004b000000500000213d0000000101000039000000000101041a000006cf011001970000000002000411000000000012004b00000a510000c13d000000060000006b00000ddd0000c13d000000010000006b00000b950000613d000500000000001d0000000501000029000000050110021000000003011000290000000202000367000000000112034f000000000101043b0000000003000031000000020430006a000001430440008a000006d105400197000006d106100197000000000756013f000000000056004b0000000005000019000006d105004041000000000041004b0000000004000019000006d104008041000006d10070009c000000000504c019000000000005004b000000500000c13d001000030010002d000000100130006a000f00000001001d000007200010009c000000500000213d0000000f01000029000001200010008c000000500000413d000000400100043d000900000001001d0000071a0010009c000000d00000213d0000000901000029000000a001100039000000400010043f0000001001200360000000000101043b000006d00010009c000000500000213d00000009040000290000000001140436000800000001001d00000010010000290000002001100039000000000112034f000000000101043b000006d00010009c000000500000213d0000001001100029001100000001001d0000001f01100039000000000031004b0000000004000019000006d104008041000006d101100197000006d105300197000000000751013f000000000051004b0000000001000019000006d101004041000006d10070009c000000000104c019000000000001004b000000500000c13d0000001101200360000000000101043b000006d00010009c000000d00000213d00000005091002100000003f04900039000006d204400197000000400600043d0000000004460019000e00000006001d000000000064004b00000000070000390000000107004039000006d00040009c000000d00000213d0000000100700190000000d00000c13d000000400040043f0000000e040000290000000000140435000000110100002900000020081000390000000009980019000000000039004b000000500000213d000000000098004b000006eb0000813d0000000e0a000029000006a80000013d000000200aa000390000000001b7001900000000000104350000000000ca04350000002008800039000000000098004b000006eb0000813d000000000182034f000000000101043b000006d00010009c000000500000213d000000110d1000290000003f01d00039000000000031004b0000000004000019000006d104008041000006d101100197000000000751013f000000000051004b0000000001000019000006d101004041000006d10070009c000000000104c019000000000001004b000000500000c13d000000200ed000390000000001e2034f000000000b01043b000006d000b0009c000000d00000213d0000001f01b0003900000769011001970000003f011000390000076901100197000000400c00043d00000000011c00190000000000c1004b00000000040000390000000104004039000006d00010009c000000d00000213d0000000100400190000000d00000c13d0000004004d00039000000400010043f0000000007bc043600000000014b0019000000000031004b000000500000213d0000002001e00039000000000412034f0000076901b00198000000000e170019000006dd0000613d000000000f04034f000000000d07001900000000f60f043c000000000d6d04360000000000ed004b000006d90000c13d0000001f0db00190000006a10000613d000000000114034f0000000304d0021000000000060e043300000000064601cf000000000646022f000000000101043b0000010004400089000000000141022f00000000014101cf000000000161019f00000000001e0435000006a10000013d00000008010000290000000e04000029000000000041043500000010010000290000004006100039000000000162034f000000000101043b000006d00010009c000000500000213d00000010071000290000001f01700039000000000031004b0000000004000019000006d104008041000006d101100197000000000851013f000000000051004b0000000001000019000006d101004041000006d10080009c000000000104c019000000000001004b000000500000c13d000000000172034f000000000401043b000006d00040009c000000d00000213d0000001f0140003900000769011001970000003f011000390000076901100197000000400500043d0000000001150019000000000051004b00000000080000390000000108004039000006d00010009c000000d00000213d0000000100800190000000d00000c13d0000002008700039000000400010043f00000000074504360000000001840019000000000031004b000000500000213d000000000882034f00000769014001980000000003170019000007230000613d000000000908034f000000000a070019000000009b09043c000000000aba043600000000003a004b0000071f0000c13d0000001f09400190000007300000613d000000000118034f0000000308900210000000000903043300000000098901cf000000000989022f000000000101043b0000010008800089000000000181022f00000000018101cf000000000191019f00000000001304350000000001470019000000000001043500000009010000290000004001100039000600000001001d00000000005104350000000f01000029000000600110008a000007200010009c000000500000213d000000600010008c000000500000413d000000400300043d000007210030009c000000d00000213d0000006001300039000000400010043f0000002001600039000000000412034f000000000404043b000000000004004b0000000005000039000000010500c039000000000054004b000000500000c13d00000000044304360000002001100039000000000512034f000000000505043b000007220050009c000000500000213d00000000005404350000002004100039000000000142034f000000000101043b000007220010009c000000500000213d0000004005300039000000000015043500000009010000290000006001100039000700000001001d00000000003104350000000f01000029000000c00110008a000007200010009c000000500000213d000000600010008c000000500000413d000000400100043d000007210010009c000000d00000213d0000006003100039000000400030043f0000002003400039000000000432034f000000000404043b000000000004004b0000000005000039000000010500c039000000000054004b000000500000c13d00000000044104360000002003300039000000000532034f000000000505043b000007220050009c000000500000213d00000000005404350000002003300039000000000232034f000000000302043b000007220030009c000000500000213d0000004002100039000000000032043500000009030000290000008003300039000400000003001d0000000000130435000000070300002900000000030304330000004005300039000000000505043300000722065001970000000057030434000000000007004b0000078f0000613d000000000006004b000012880000613d00000000050504330000072205500197000000000056004b000007940000413d000012880000013d000000000006004b000012740000c13d00000000050504330000072200500198000012740000c13d000000000202043300000722022001970000000003010433000000000003004b000007a00000613d000000000002004b0000128f0000613d00000000030404330000072203300197000000000032004b000007a50000413d0000128f0000013d000000000002004b000012780000c13d00000000020404330000072200200198000012780000c13d000000060100002900000000010104330000000001010433000000000001004b00000ba40000613d00000009010000290000000001010433000006d001100197001100000001001d000000000010043f0000000601000039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d000000000101043b000000000101041a000000000001004b0000126d0000c13d0000000501000039000000000101041a000006d00010009c000000d00000213d00000001021000390000000503000039000000000023041b000007270110009a0000001102000029000000000021041b000000000103041a001000000001001d000000000020043f0000000601000039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d000000000101043b0000001002000029000000000021041b00000009010000290000000001010433000006d001100197000000000010043f0000000701000039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d000000000301043b000000070100002900000000010104330000000024010434001100000004001d000000000004004b0000000004000039000000010400c039001000000004001d000000400400043d0000071a0040009c000000d00000213d000f00000003001d0000000002020433000007220220019700000040011000390000000001010433000b00000001001d000000a001400039000000400010043f000e00000002001d000d00000004001d0000000001240436000c00000001001d000007280100004100000000001004430000000001000414000006cc0010009c000006cc01008041000000c00110021000000729011001c70000800b020000391b2c1b270000040f0000000100200190000013b30000613d0000000b020000290000072202200197000000000101043b0000000d040000290000004003400039000000100500002900000000005304350000008003400039000000000023043500000060034000390000000e050000290000000000530435000006cc011001970000000c030000290000000000130435000000110000006b00000000030000190000072a0300c0410000000f09000029000000000409041a0000072b04400197000000000343019f0000008002200210000000000425019f000000000353019f0000008002100210000000000323019f000000000039041b0000000103900039000000000043041b000000400300043d0000071a0030009c000000d00000213d0000000404000029000000000404043300000020054000390000000005050433000000000604043300000040044000390000000004040433000000a007300039000000400070043f0000008007300039000007220840019700000000008704350000004007300039000000000006004b0000000006000039000000010600c039000000000067043500000020063000390000000000160435000007220150019700000060053000390000000000150435000000000013043500000000030000190000072a0300c0410000000205900039000000000605041a0000072b06600197000000000363019f000000000223019f000000000212019f000000000025041b0000008002400210000000000112019f0000000302900039000000000012041b000000060100002900000000030104330000000065030434000006d00050009c000000d00000213d0000000404900039000000000104041a000000010010019000000001071002700000007f0770618f0000001f0070008c00000000020000390000000102002039000000000121013f000000010010019000000f7f0000c13d000000200070008c001100000004001d001000000005001d000e00000003001d000008840000413d000d00000007001d000f00000006001d000000000040043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006db011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d00000010050000290000001f025000390000000502200270000000200050008c0000000002004019000000000301043b0000000d010000290000001f01100039000000050110027000000000011300190000000002230019000000000012004b00000011040000290000000f06000029000008840000813d000000000002041b0000000102200039000000000012004b000008800000413d0000001f0050008c000008a20000a13d000000000040043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006db011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d00000010060000290000076902600198000000000101043b0000000e07000029000008ae0000613d000000010320008a000000050330027000000000033100190000000104300039000000200300003900000000057300190000000005050433000000000051041b00000020033000390000000101100039000000000041004b0000089a0000c13d000008af0000013d000000000005004b000008a60000613d0000000001060433000008a70000013d000000000100001900000003025002100000076a0220027f0000076a02200167000000000121016f0000000102500210000000000121019f000008bc0000013d0000002003000039000000000062004b000008b90000813d0000000302600210000000f80220018f0000076a0220027f0000076a0220016700000000037300190000000003030433000000000223016f000000000021041b000000010160021000000001011001bf0000001104000029000000000014041b000000080100002900000000010104330000000002010433000000000002004b000009df0000613d0000000003000019000c00000003001d0000000502300210000000000121001900000020011000390000000001010433001000000001001d0000000031010434000000000001004b00000ba40000613d00000009020000290000000002020433001100000002001d000006cc0010009c000006cc010080410000006001100210000006cc0030009c000f00000003001d000006cc0200004100000000020340190000004002200210000000000121019f0000000002000414000006cc0020009c000006cc02008041000000c002200210000000000112019f00000711011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d0000001102000029000006d002200197000000000101043b001100000001001d000b00000002001d000000000020043f0000000701000039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d000000000201043b0000001101000029000000000010043f000e00000002001d0000000601200039000d00000001001d000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d000000000101043b000000000101041a000000000001004b000011120000c13d0000000e010000290000000502100039000000000102041a000006d00010009c000000d00000213d000a00000001001d0000000101100039000000000012041b000e00000002001d000000000020043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006db011001c700008010020000391b2c1b270000040f00000001002001900000001102000029000000500000613d000000000101043b0000000a01100029000000000021041b0000000e01000029000000000101041a000e00000001001d000000000020043f0000000d01000029000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f00000001002001900000001102000029000000500000613d000000000101043b0000000e03000029000000000031041b000000000020043f0000000801000039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d000000000401043b00000010010000290000000005010433000006d00050009c000000d00000213d000000000104041a000000010010019000000001031002700000007f0330618f0000001f0030008c00000000020000390000000102002039000000000121013f00000001001001900000000f0700002900000f7f0000c13d000000200030008c001100000004001d000e00000005001d0000096f0000413d000d00000003001d000000000040043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006db011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d0000000e050000290000001f025000390000000502200270000000200050008c0000000002004019000000000301043b0000000d010000290000001f01100039000000050110027000000000011300190000000002230019000000000012004b0000000f0700002900000011040000290000096f0000813d000000000002041b0000000102200039000000000012004b0000096b0000413d0000001f0050008c0000099b0000a13d000000000040043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006db011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d0000000e080000290000076902800198000000000101043b000009d90000613d000000010320008a00000005033002700000000003310019000000010430003900000020030000390000000f07000029000000100600002900000000056300190000000005050433000000000051041b00000020033000390000000101100039000000000041004b000009860000c13d000000000082004b000009970000813d0000000302800210000000f80220018f0000076a0220027f0000076a0220016700000000036300190000000003030433000000000223016f000000000021041b000000010180021000000001011001bf0000001104000029000009a70000013d000000000005004b0000099f0000613d0000000001070433000009a00000013d0000000001000019000000100600002900000003025002100000076a0220027f0000076a02200167000000000121016f0000000102500210000000000121019f000000000014041b000000400100043d00000020020000390000000003210436000000000206043300000000002304350000004003100039000000000002004b000009b80000613d000000000400001900000000053400190000000006740019000000000606043300000000006504350000002004400039000000000024004b000009b10000413d0000001f042000390000076904400197000000000232001900000000000204350000004002400039000006cc0020009c000006cc020080410000006002200210000006cc0010009c000006cc010080410000004001100210000000000112019f0000000002000414000006cc0020009c000006cc02008041000000c002200210000000000112019f00000711011001c70000800d0200003900000002030000390000072d040000410000000b050000291b2c1b220000040f0000000100200190000000500000613d0000000c030000290000000103300039000000080100002900000000010104330000000002010433000000000023004b000008c30000413d000009df0000013d00000020030000390000000f070000290000001006000029000000000082004b0000098f0000413d000009970000013d00000004010000290000000002010433000000070100002900000000050104330000000601000029000000000301043300000009010000290000000004010433000000400100043d000000200610003900000100070000390000000000760435000006d00440019700000000004104350000010007100039000000006403043400000000004704350000012003100039000000000004004b000009fb0000613d000000000700001900000000083700190000000009760019000000000909043300000000009804350000002007700039000000000047004b000009f40000413d000000000643001900000000000604350000000076050434000000000006004b0000000006000039000000010600c039000000400810003900000000006804350000000006070433000007220660019700000060071000390000000000670435000000400550003900000000050504330000072205500197000000800610003900000000005604350000000065020434000000000005004b0000000005000039000000010500c039000000a007100039000000000057043500000000050604330000072205500197000000c0061000390000000000560435000000400220003900000000020204330000072202200197000000e00510003900000000002504350000001f02400039000007690220019700000000021200490000000002320019000006cc0020009c000006cc020080410000006002200210000006cc0010009c000006cc010080410000004001100210000000000112019f0000000002000414000006cc0020009c000006cc02008041000000c002200210000000000112019f00000711011001c70000800d0200003900000001030000390000072e040000411b2c1b220000040f0000000100200190000000500000613d00000005020000290000000102200039000500000002001d000000010020006c000006420000413d00000b950000013d0000000001000416000000000001004b000000500000c13d000000000100041a000006cf021001970000000006000411000000000026004b00000a970000c13d0000000102000039000000000302041a000006d304300197000000000464019f000000000042041b000006d301100197000000000010041b0000000001000414000006cf05300197000006cc0010009c000006cc01008041000000c00110021000000711011001c70000800d0200003900000003030000390000074f0400004100000b920000013d0000075401000041000000800010043f000007140100004100001b2e00010430000000a007000039000007350400004100000000060000190000000005070019000000000704041a000000000775043600000001044000390000000106600039000000000026004b00000a580000413d000007360250009a000007370020009c000000d00000413d000000410250008a00000769022001970000008002200039000000800400043d000000400020043f000006d00040009c000000d00000213d000001410000013d000000a0050000390000073f0300004100000000040000190000000006050019000000000503041a000000000556043600000001033000390000000104400039000000000014004b00000a6d0000413d000000410160008a0000076904100197000007400040009c000000d00000213d0000008001400039000000400010043f0000000000210435000000a002400039000000800300043d0000000000320435000000c002400039000000000003004b00000a890000613d000000a00400003900000000050000190000000046040434000006cf0660019700000000026204360000000105500039000000000035004b00000a830000413d0000000002120049000006cc0020009c000006cc020080410000006002200210000006cc0010009c000006cc010080410000004001100210000000000112019f00001b2d0001042e0000071701000041000000800010043f000000840030043f0000075f0100004100001b2e000104300000074e01000041000000800010043f000007140100004100001b2e000104300000000a01000039000000000101041a000006cf011001970000000004000411000000000041004b00000ae50000c13d00000020013000390000071902000041000000000021043500000064013000390000001102000029000000000021043500000044013000390000000002000410000000000021043500000024013000390000000000410435000000640100003900000000001304350000071a0030009c000000d00000213d000000a001300039000000400010043f000007150100004100000000001004430000000001000412000000040010044300000024000004430000000001000414000006cc0010009c000006cc01008041000000c00110021000000716011001c70000800502000039001000000003001d1b2c1b270000040f0000000100200190000013b30000613d000000000101043b00000010020000291b2c18700000040f0000000001000414000006cc0010009c000006cc01008041000000c00110021000000711011001c70000800d0200003900000003030000390000071b0400004100000b900000013d000007650020009c000006050000613d000007660020009c000006050000613d000000800000043f000007300100004100001b2d0001042e0000000402000039000000000302041a000006d304300197000000000414019f000000000042041b000006cf02300197000000800020043f000000a00010043f0000000001000414000006cc0010009c000006cc01008041000000c00110021000000739011001c70000800d0200003900000001030000390000073a0400004100000b920000013d0000071701000041000000000013043500000004013000390000000000410435000006cc0030009c000006cc03008041000000400130021000000718011001c700001b2e00010430000000000100041a000006d301100197000000000161019f000000000010041b0000000001000414000006cc0010009c000006cc01008041000000c00110021000000711011001c70000800d020000390000000303000039000007120400004100000b920000013d000006cc00b0009c000006cc0300004100000000030b40190000004003300210000006cc0010009c000006cc01008041000000c001100210000000000131019f00000718011001c7000f0000000b001d1b2c1b270000040f0000000f0b0000290000006003100270000006cc03300197000000200030008c000000200400003900000000040340190000001f0640018f000000200740019000000000057b001900000b160000613d000000000801034f00000000090b0019000000008a08043c0000000009a90436000000000059004b00000b120000c13d000000000006004b00000b230000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f0003000000010355000000010020019000000b3c0000613d0000001f01400039000000600210018f0000000001b20019000000000021004b00000000020000390000000102004039000006d00010009c000000d00000213d0000000100200190000000d00000c13d000000400010043f000000200030008c000000500000413d00000000020b04330000001103000029000000000032004b00000b850000813d000007620200004100000ba60000013d000000100100002900000b5b0000013d0000001f0530018f000006ce06300198000000400200043d000000000462001900000b470000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b00000b430000c13d000000000005004b00000b540000613d000000000161034f0000000305500210000000000604043300000000065601cf000000000656022f000000000101043b0000010005500089000000000151022f00000000015101cf000000000161019f00000000001404350000006001300210000006cc0020009c000006cc020080410000004002200210000000000112019f00001b2e000104300000000f01000029000000000112034f000000000101043b000006cf0010009c000000500000213d000000400200043d0000074203000041000000000032043500000004032000390000000000130435000006cc0020009c000006cc02008041000000400120021000000718011001c700001b2e000104300000000104000039000000000504041a000006d305500197000000000225019f000000000024041b000000000003004b00000ba40000613d000006cf0210019800000ba40000613d000000100000006b00000ba40000613d000000800010043f000000c00030043f000006d401000041000000400300043d000d00000003001d00000000001304350000000001000414000000040020008c00000cf50000c13d0000000001000415000000210110008a00000005011002100000000103000031000000200030008c0000002004000039000000000403401900000d220000013d000000100100002900000000020004111b2c15f20000040f0000000001000414000006cc0010009c000006cc01008041000000c00110021000000711011001c70000800d0200003900000003030000390000076104000041000000000500041100000011060000291b2c1b220000040f0000000100200190000000500000613d000000000100001900001b2d0001042e000006cc033001970000001f0530018f000006ce06300198000000400200043d000000000462001900000b470000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b00000b9f0000c13d00000b470000013d000000400100043d0000072f020000410000000000210435000006cc0010009c000006cc010080410000004001100210000006d5011001c700001b2e000104300000001003000029000006cc0030009c000006cc030080410000004003300210000006cc0010009c000006cc01008041000000c001100210000000000131019f00000718011001c71b2c1b270000040f0000006003100270000006cc03300197000000200030008c000000200400003900000000040340190000001f0640018f0000002007400190000000100570002900000bc50000613d000000000801034f0000001009000029000000008a08043c0000000009a90436000000000059004b00000bc10000c13d000000000006004b00000bd20000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f0003000000010355000000010020019000000d710000613d0000001f01400039000000600210018f0000001001200029000000000021004b00000000020000390000000102004039000006d00010009c000000d00000213d0000000100200190000000d00000c13d000000400010043f000000200030008c000000500000413d00000010020000290000000002020433000000000002004b0000000003000039000000010300c039000000000032004b000000500000c13d000000000002004b00000f850000c13d0000000f0100002900000020011000390000000201100367000000000101043b001000000001001d000006cf0010009c000000500000213d0000071501000041000000000010044300000000010004120000000400100443000000600100003900000024001004430000000001000414000006cc0010009c000006cc01008041000000c00110021000000716011001c700008005020000391b2c1b270000040f0000000100200190000013b30000613d000000000101043b000000000001004b000011290000c13d0000000f010000290000000201100367000000000101043b001000000001001d000006d00010009c000000500000213d0000001001000029000000000010043f0000000601000039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d000000400200043d000e00000002001d0000000402200039000000000101043b000000000101041a000000000001004b000011a50000c13d0000074a010000410000000e030000290000000000130435000000100100002900000ddb0000013d0000000002000019001100000002001d0000000502200210001000000002001d0000000e012000290000000001010433000000000010043f0000000801000039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d000000000101043b000000000201041a000000010320019000000001052002700000007f0550618f0000001f0050008c00000000040000390000000104002039000000000043004b00000f7f0000c13d000000400700043d0000000004570436000000000003004b00000c620000613d000900000004001d000a00000005001d000b00000007001d000000000010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006db011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d0000000a08000029000000000008004b000000200500008a00000c6a0000613d000000000201043b00000000010000190000000d060000290000000b0700002900000009090000290000000003190019000000000402041a000000000043043500000001022000390000002001100039000000000081004b00000c5a0000413d00000c6d0000013d0000076b012001970000000000140435000000000005004b00000020010000390000000001006039000000200500008a0000000d0600002900000c6d0000013d00000000010000190000000d060000290000000b070000290000003f01100039000000000251016f0000000001720019000000000021004b00000000020000390000000102004039000006d00010009c000000d00000213d0000000100200190000000d00000c13d000000400010043f00000000010604330000001102000029000000000021004b000010ed0000a13d00000010030000290000000c0130002900000000007104350000000001060433000000000021004b000010ed0000a13d00000001022000390000000f010000290000000001010433000000000012004b00000c250000413d000004260000013d0000001101000029000000000010043f0000000701000039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d00000010020000290000001f022000390000076902200197000b00000002001d0000003f022000390000076902200197000000000101043b000c00000001001d000000400100043d0000000002210019000000000012004b00000000040000390000000104004039000006d00020009c000000d00000213d0000000100400190000000d00000c13d000000400020043f000000100200002900000000022104360000000e05000029000000000050007c000000500000213d00000010040000290000076903400198000e001f00400193000a00000003001d00000000033200190000000f040000290000002004400039000f00000004001d000000020440036700000cbc0000613d000000000504034f0000000006020019000000005705043c0000000006760436000000000036004b00000cb80000c13d0000000e0000006b00000cca0000613d0000000a044003600000000e050000290000000305500210000000000603043300000000065601cf000000000656022f000000000404043b0000010005500089000000000454022f00000000045401cf000000000464019f000000000043043500000010032000290000000000030435000006cc0020009c000006cc0200804100000040022002100000000001010433000006cc0010009c000006cc010080410000006001100210000000000121019f0000000002000414000006cc0020009c000006cc02008041000000c002200210000000000112019f00000711011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d0000000c020000290000000503200039000000000201043b00000000010300191b2c19500000040f000000400700043d000000000001004b000010130000c13d000c00000007001d0000002401700039000000400200003900000000002104350000073e01000041000000000017043500000004017000390000001102000029000000000021043500000044037000390000000d0100002900000010020000291b2c15550000040f0000000c02000029000011200000013d0000000d03000029000006cc0030009c000006cc030080410000004003300210000006cc0010009c000006cc01008041000000c001100210000000000131019f000006d5011001c71b2c1b270000040f0000006003100270000006cc03300197000000200030008c000000200400003900000000040340190000001f0640018f00000020074001900000000d0570002900000d0e0000613d000000000801034f0000000d09000029000000008a08043c0000000009a90436000000000059004b00000d0a0000c13d000000000006004b00000d1b0000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f00030000000103550000000001000415000000200110008a0000000501100210000000010020019000000d390000613d0000001f02400039000000600420018f0000000d02400029000000000042004b00000000040000390000000104004039000006d00020009c000000d00000213d0000000100400190000000d00000c13d000000400020043f000000200030008c000000500000413d0000000d030000290000000003030433000000ff0030008c000000500000213d0000000501100270000000000103001f0000001101000029000000ff0110018f000000000031004b00000fab0000c13d0000001101000029000000a00010043f0000000401000039000000000201041a000006d30220019700000010022001af000000000021041b0000000f010000290000000001010433000000000001004b0000000001000039000000010100c039000000e00010043f0000000001000019000010fa0000613d000000400100043d000006d80010009c000000d00000213d0000002002100039000000400020043f0000000000010435000000e00100043d000000000001004b00000fb60000c13d000000400100043d000007570200004100000ba60000013d00000000040000190000000d0c00002900000d5f0000013d0000001f0760003900000769077001970000000006650019000000000006043500000000057500190000000104400039000000000024004b000004310000813d0000000006150049000000400660008a0000000003630436000000200cc0003900000000060c043300000000760604340000000005650436000000000006004b00000d570000613d00000000080000190000000009580019000000000a870019000000000a0a04330000000000a904350000002008800039000000000068004b00000d690000413d00000d570000013d0000001f0530018f000006ce06300198000000400200043d000000000462001900000b470000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b00000d780000c13d00000b470000013d0000000f03000029000006cc0030009c000006cc030080410000004003300210000006cc0010009c000006cc01008041000000c001100210000000000131019f00000718011001c71b2c1b270000040f0000006003100270000006cc03300197000000200030008c000000200400003900000000040340190000001f0640018f00000020074001900000000f0570002900000d960000613d000000000801034f0000000f09000029000000008a08043c0000000009a90436000000000059004b00000d920000c13d000000000006004b00000da30000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f0003000000010355000000010020019000000f730000613d0000001f01400039000000600210018f0000000f01200029000000000021004b00000000020000390000000102004039000006d00010009c000000d00000213d0000000100200190000000d00000c13d000000400010043f000000200030008c000000500000413d0000000f020000290000000002020433000000000002004b0000000003000039000000010300c039000000000032004b000000500000c13d000000000002004b00000f850000c13d0000000d010000290000000201100367000000000101043b000f00000001001d000006d00010009c000000500000213d0000000f01000029000000000010043f0000000601000039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d000000400200043d000c00000002001d0000000402200039000000000101043b000000000101041a000000000001004b000011470000c13d0000074a010000410000000c0300002900000000001304350000000f01000029000000000012043500000ae90000013d0000000002000019000700000002001d000000050120021000000005011000290000000201100367000000000101043b000b00000001001d000006d00010009c000000500000213d0000000b01000029000000000010043f0000000601000039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d000000000101043b000000000301041a000000000003004b0000115a0000613d0000000501000039000000000201041a000000000002004b000013c00000613d000000010130008a000000000023004b00000e190000613d000000000012004b000010ed0000a13d0000071d0130009a0000071d0220009a000000000202041a000000000021041b000000000020043f0000000601000039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c70000801002000039001100000003001d1b2c1b270000040f0000000100200190000000500000613d000000000101043b0000001102000029000000000021041b0000000501000039000000000301041a000000000003004b000010f30000613d000000010130008a0000071d0230009a000000000002041b0000000502000039000000000012041b0000000b01000029000000000010043f0000000601000039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d000000000101043b000000000001041b0000000b01000029000000000010043f0000000701000039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d000000000101043b0000000501100039000000000301041a000000400200043d000f00000002001d001100000003001d0000000002320436000a00000002001d000000000010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006db011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d0000001105000029000000000005004b0000000a0200002900000e580000613d000000000101043b0000000a020000290000000003000019000000000401041a000000000242043600000001011000390000000103300039000000000053004b00000e520000413d0000000f0120006a0000001f0110003900000769021001970000000f01200029000000000021004b00000000020000390000000102004039000006d00010009c000000d00000213d0000000100200190000000d00000c13d000000400010043f0000000f010000290000000001010433000000000001004b00000efd0000613d000000000200001900000e700000013d000000110200002900000001022000390000000f010000290000000001010433000000000012004b00000efd0000813d001100000002001d0000000b01000029000000000010043f0000000701000039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d000000000301043b0000000f0100002900000000010104330000001102000029000000000021004b000010ed0000a13d00000005012002100000000a011000290000000001010433000c00000001001d000000000010043f000d00000003001d0000000601300039000e00000001001d000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d000000000101043b000000000101041a000000000001004b00000e6a0000613d0000000d020000290000000503200039000000000203041a000000000002004b000013c00000613d000000000021004b001000000001001d000d00000003001d00000edc0000613d000900000002001d000000000030043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006db011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d00000010020000290008000100200092000000000101043b0000000d04000029000000000204041a000000080020006c000010ed0000a13d0000000902000029000000010220008a0000000001120019000000000101041a000900000001001d000000000040043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006db011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d000000000101043b00000008011000290000000902000029000000000021041b000000000020043f0000000e01000029000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d000000000101043b0000001002000029000000000021041b0000000d03000029000000000103041a001000000001001d000000000001004b000010f30000613d000000000030043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006db011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d0000001002000029000000010220008a000000000101043b0000000001210019000000000001041b0000000d01000029000000000021041b0000000c01000029000000000010043f0000000e01000029000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d000000000101043b000000000001041b00000e6a0000013d0000000b01000029000000000010043f0000000701000039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d000000000301043b000000000003041b0000000101300039000000000001041b0000000201300039000000000001041b0000000301300039000000000001041b0000000404300039000000000104041a000000010010019000000001051002700000007f0550618f0000001f0050008c00000000020000390000000102002039000000000121013f000000010010019000000f7f0000c13d000000000005004b00000f3f0000613d0000001f0050008c00000f3e0000a13d000f00000005001d001100000003001d001000000004001d000000000040043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006db011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d000000000101043b0000000f020000290000001f02200039000000050220027000000000022100190000000103100039000000000023004b00000f3a0000813d000000000003041b0000000103300039000000000023004b00000f360000413d0000001002000029000000000002041b00000000040100190000001103000029000000000004041b0000000501300039000000000201041a000000000001041b000000000002004b00000f570000613d001100000002001d000000000010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006db011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d000000000101043b0000001102100029000000000021004b00000f570000813d000000000001041b0000000101100039000000000021004b00000f530000413d000000400100043d0000000b020000290000000000210435000006cc0010009c000006cc0100804100000040011002100000000002000414000006cc0020009c000006cc02008041000000c002200210000000000112019f000006db011001c70000800d0200003900000001030000390000071e040000411b2c1b220000040f0000000100200190000000500000613d00000007020000290000000102200039000000060020006c00000dde0000413d0000063f0000013d0000071702000041000001400020043f000001440010043f000007310100004100001b2e000104300000001f0530018f000006ce06300198000000400200043d000000000462001900000b470000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b00000f7a0000c13d00000b470000013d0000074101000041000000000010043f0000002201000039000000040010043f000007180100004100001b2e00010430000007450200004100000ba60000013d000000000001004b000010060000613d000000a00200043d0000072202200197000000000021004b000010060000813d0000001101000029000000000010043f0000000701000039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d000000000101043b00000080020000391b2c19dc0000040f000001200100043d0000072201100197000000e00200043d000000000002004b000010460000c13d000000000001004b00000fa70000c13d000001000100043d00000722001001980000104c0000613d000000400200043d001100000002001d0000072501000041000011430000013d00000024042000390000000000340435000006d603000041000000000032043500000004032000390000000000130435000006cc0020009c000006cc020080410000004001200210000006d7011001c700001b2e000104300000000f020000290000000002020433000000000002004b000010fa0000613d000000000200001900000fc20000013d000000110200002900000001022000390000000f010000290000000001010433000000000012004b000010f90000813d001100000002001d00000005012002100000000e011000290000000001010433000006cf0310019800000fbc0000613d000000000030043f0000000301000039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c70000801002000039001000000003001d1b2c1b270000040f00000010040000290000000100200190000000500000613d000000000101043b000000000101041a000000000001004b00000fbc0000c13d0000000203000039000000000103041a000006d00010009c000000d00000213d0000000102100039000000000023041b000006da0110009a000000000041041b000000000103041a000d00000001001d000000000040043f0000000301000039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f00000010030000290000000100200190000000500000613d000000000101043b0000000d02000029000000000021041b000000400100043d0000000000310435000006cc0010009c000006cc0100804100000040011002100000000002000414000006cc0020009c000006cc02008041000000c002200210000000000112019f000006db011001c70000800d020000390000000103000039000006dc040000411b2c1b220000040f000000010020019000000fbc0000c13d000000500000013d000000400200043d001100000002001d00000723010000410000000000120435000000040220003900000080010000391b2c15e40000040f0000001101000029000006cc0010009c000006cc01008041000000400110021000000724011001c700001b2e000104300000002001700039000000100200002900000000002104350000002001000039000000000017043500000040017000390000000a021000290000000f0300002900000002033003670000000a0000006b000010240000613d000000000403034f0000000005010019000000004604043c0000000005650436000000000025004b000010200000c13d0000000e0000006b000010320000613d0000000a033003600000000e040000290000000304400210000000000502043300000000054501cf000000000545022f000000000303043b0000010004400089000000000343022f00000000034301cf000000000353019f000000000032043500000010011000290000000000010435000006cc0070009c000006cc0700804100000040017002100000000b020000290000073b0020009c0000073b020080410000006002200210000000000112019f0000000002000414000006cc0020009c000006cc02008041000000c002200210000000000121019f0000073c0110009a0000800d0200003900000002030000390000073d04000041000003190000013d000000000001004b000011400000613d000001000200043d0000072202200197000000000021004b000011400000813d0000001101000029000000000010043f0000000701000039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d000000000101043b0000000201100039000000e0020000391b2c19dc0000040f000000400100043d00000011020000290000000002210436000000800300043d000000000003004b0000000003000039000000010300c0390000000000320435000000a00200043d000007220220019700000040031000390000000000230435000000c00200043d000007220220019700000060031000390000000000230435000000e00200043d000000000002004b0000000002000039000000010200c03900000080031000390000000000230435000001000200043d0000072202200197000000a0031000390000000000230435000001200200043d0000072202200197000000c0031000390000000000230435000006cc0010009c000006cc0100804100000040011002100000000002000414000006cc0020009c000006cc02008041000000c002200210000000000112019f00000732011001c70000800d020000390000000103000039000007330400004100000ae40000013d00000000020000190000108f0000013d00000010020000290000000102200039000000800100043d000000000012004b000001ec0000813d001000000002001d0000000501200210000000a0011000390000000001010433000006cf01100197001100000001001d000000000010043f0000000301000039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d000000000101043b000000000301041a000000000003004b0000108a0000613d0000000201000039000000000201041a000000000002004b000013c00000613d000000010130008a000000000023004b000010c70000613d000000000012004b000010ed0000a13d000007550130009a000007550220009a000000000202041a000000000021041b000000000020043f0000000301000039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c70000801002000039000d00000003001d1b2c1b270000040f0000000d030000290000000100200190000000500000613d000000000101043b000000000031041b0000000201000039000000000301041a000000000003004b000010f30000613d000000010130008a000007550230009a000000000002041b0000000202000039000000000012041b0000001101000029000000000010043f0000000301000039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d000000000101043b000000000001041b000000400100043d00000011020000290000000000210435000006cc0010009c000006cc0100804100000040011002100000000002000414000006cc0020009c000006cc02008041000000c002200210000000000112019f000006db011001c70000800d02000039000000010300003900000756040000411b2c1b220000040f00000001002001900000108a0000c13d000000500000013d0000074101000041000000000010043f0000003201000039000000040010043f000007180100004100001b2e000104300000074101000041000000000010043f0000003101000039000000040010043f000007180100004100001b2e00010430000000e00100043d0000000c05000029000001000050043f000000800200043d00000140000004430000016000200443000000a00200043d00000020030000390000018000300443000001a0002004430000004002000039000000c00400043d000001c000200443000001e000400443000000600200003900000200002004430000022000100443000000800100003900000240001004430000026000500443000001000030044300000005010000390000012000100443000006dd0100004100001b2d0001042e000000400300043d001100000003001d0000002401300039000000400200003900000000002104350000072c01000041000000000013043500000004013000390000000b020000290000000000210435000000440230003900000010010000291b2c143a0000040f00000011020000290000000001210049000006cc0010009c000006cc01008041000006cc0020009c000006cc0200804100000060011002100000004002200210000000000121019f00001b2e000104300000001001000029000000000010043f0000000301000039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d000000000101043b000000000101041a000000000001004b00000c050000c13d000000400100043d0000074602000041000000000021043500000004021000390000001003000029000004ea0000013d000000400200043d001100000002001d000007230100004100000000001204350000000402200039000000e0010000390000100c0000013d0000000401000039000000000301041a00000758010000410000000c0400002900000000001404350000000f0100002900000000001204350000002401400039000000000200041100000000002104350000000001000414000006cf02300197000000040020008c000011600000c13d0000000103000031000000200030008c000000200400003900000000040340190000118a0000013d000000400100043d0000071f02000041000000000021043500000004021000390000000b03000029000004ea0000013d0000000c03000029000006cc0030009c000006cc030080410000004003300210000006cc0010009c000006cc01008041000000c001100210000000000131019f000006d7011001c71b2c1b270000040f0000006003100270000006cc03300197000000200030008c000000200400003900000000040340190000001f0640018f00000020074001900000000c05700029000011790000613d000000000801034f0000000c09000029000000008a08043c0000000009a90436000000000059004b000011750000c13d000000000006004b000011860000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f00030000000103550000000100200190000011b50000613d0000001f01400039000000600210018f0000000c01200029000000000021004b00000000020000390000000102004039000006d00010009c000000d00000213d0000000100200190000000d00000c13d000000400010043f000000200030008c000000500000413d0000000c020000290000000002020433000000000002004b0000000003000039000000010300c039000000000032004b000000500000c13d000000000002004b000012950000c13d0000074802000041000000000021043500000004021000390000000003000411000004ea0000013d0000000401000039000000000301041a00000747010000410000000e040000290000000000140435000000100100002900000000001204350000000001000414000006cf02300197000000040020008c000011c10000c13d0000000103000031000000200030008c00000020040000390000000004034019000011eb0000013d0000001f0530018f000006ce06300198000000400200043d000000000462001900000b470000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000011bc0000c13d00000b470000013d0000000e03000029000006cc0030009c000006cc030080410000004003300210000006cc0010009c000006cc01008041000000c001100210000000000131019f00000718011001c71b2c1b270000040f0000006003100270000006cc03300197000000200030008c000000200400003900000000040340190000001f0640018f00000020074001900000000e05700029000011da0000613d000000000801034f0000000e09000029000000008a08043c0000000009a90436000000000059004b000011d60000c13d000000000006004b000011e70000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f000300000001035500000001002001900000127c0000613d0000001f01400039000000600210018f0000000e01200029000000000021004b00000000020000390000000102004039000006d00010009c000000d00000213d0000000100200190000000d00000c13d000000400010043f000000200030008c000000500000413d0000000e0200002900000000020204330000072a0020009c000000500000813d0000000003000411000000000023004b000012d50000c13d00000002010003670000000f02100360000000000202043b000006d00020009c000000500000213d0000000f030000290000004003300039000000000131034f000000000101043b001000000001001d000000000020043f0000000701000039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d000000000101043b000000100200002900000011030000291b2c1a490000040f000000400100043d00000010020000290000000000210435000006cc0010009c000006cc0100804100000040011002100000000002000414000006cc0020009c000006cc02008041000000c002200210000000000112019f000006db011001c70000800d020000390000000203000039000007490400004100000000050004111b2c1b220000040f0000000100200190000000500000613d0000000f010000290000000201100367000000000101043b000006d00010009c000000500000213d1b2c15890000040f000000400200043d001100000002001d0000000002000412001900000002001d001800200000003d001000000001001d000080050100003900000044030000390000000004000415000000190440008a000000050440021000000715020000411b2c1b040000040f000000ff0310018f000000110100002900000020021000390000000000320435000000200200003900000000002104351b2c14240000040f000000400100043d000e00000001001d1b2c14240000040f0000000e020000290000002001200039000f00000001001d00000011030000290000000000310435000000100100002900000000001204350000000003020019000000400400043d001100000004001d0000002001000039000000000214043600000000010304330000004003000039000000000032043500000060024000391b2c143a0000040f000000000201001900000011040000290000000001410049000000200310008a0000000f010000290000000001010433000000400440003900000000003404351b2c143a0000040f00000011020000290000000001210049000006cc0020009c000006cc020080410000004002200210000006cc0010009c000006cc010080410000006001100210000000000121019f00001b2d0001042e00000009010000290000000001010433000000400200043d00000726030000410000000000320435000006d00110019700000b620000013d000000400200043d001100000002001d00000725010000410000128b0000013d000000400300043d001100000003001d0000072502000041000012920000013d0000001f0530018f000006ce06300198000000400200043d000000000462001900000b470000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000012830000c13d00000b470000013d000000400200043d001100000002001d00000723010000410000000000120435000000040220003900000000010300190000100c0000013d000000400300043d001100000003001d0000072302000041000000000023043500000004023000390000100c0000013d00000002020003670000000d01200360000000000101043b000006d00010009c000000500000213d0000000d030000290000008003300039000000000332034f000000000303043b0000000004000031000000110540006a000000230550008a000006d106500197000006d107300197000000000867013f000000000067004b0000000006000019000006d106004041000000000053004b0000000005000019000006d105008041000006d10080009c000000000605c019000000000006004b000000500000c13d0000001105000029000f00040050003d0000000f05300029000000000252034f000000000302043b000006d00030009c000000500000213d00000000043400490000002002500039000006d105400197000006d106200197000000000756013f000000000056004b0000000005000019000006d105004041000000000042004b0000000004000019000006d104002041000006d10070009c000000000504c019000000000005004b000000500000c13d1b2c14e40000040f000000000001004b000012d90000c13d0000001101000029000000a4021000390000000f010000291b2c14820000040f0000075c03000041000000400500043d001100000005001d000000000035043500000004035000390000002004000039000000000043043500000024035000391b2c15550000040f0000111f0000013d000007480200004100000000002104350000000402100039000004ea0000013d00000002010003670000000d02100360000000000202043b000006d00020009c000000500000213d0000000d03000029000c00400030003d0000000c01100360000000000101043b000d00000001001d000000000020043f0000000701000039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000100200190000000500000613d000000000101043b00000002011000390000000d0200002900000010030000291b2c1a490000040f0000000c01000029001000600010003d00000002030003670000001001300360000000000101043b0000000004000031000000110240006a000000230220008a000006d105200197000006d106100197000000000756013f000000000056004b0000000005000019000006d105004041000000000021004b0000000002000019000006d102008041000006d10070009c000000000502c019000000000005004b000000500000c13d0000000f01100029000000000213034f000000000202043b000006d00020009c000000500000213d00000000052400490000002006100039000006d101500197000006d107600197000000000817013f000000000017004b0000000001000019000006d101004041000000000056004b0000000005000019000006d105002041000006d10080009c000000000105c019000000000001004b000000500000c13d0000001f0120003900000769011001970000003f011000390000076905100197000000400100043d0000000005510019000000000015004b00000000080000390000000108004039000006d00050009c000000d00000213d0000000100800190000000d00000c13d000000400050043f00000000052104360000000008620019000000000048004b000000500000213d000000000463034f00000769062001980000001f0720018f00000000036500190000133a0000613d000000000804034f0000000009050019000000008a08043c0000000009a90436000000000039004b000013360000c13d000000000007004b000013470000613d000000000464034f0000000306700210000000000703043300000000076701cf000000000767022f000000000404043b0000010006600089000000000464022f00000000046401cf000000000474019f0000000000430435000000000225001900000000000204350000000002010433000000200020008c000013600000613d000000000002004b000013640000c13d0000071501000041000000000010044300000000010004120000000400100443000000200100003900000024001004430000000001000414000006cc0010009c000006cc01008041000000c00110021000000716011001c700008005020000391b2c1b270000040f0000000100200190000013b30000613d000000000101043b001100000001001d0000136d0000013d0000000002050433001100000002001d000000ff0020008c0000136d0000a13d000000400400043d001100000004001d0000075b02000041000000000024043500000004024000390000002003000039000000000032043500000024024000390000111e0000013d0000071501000041000000000010044300000000010004120000000400100443000000200100003900000024001004430000000001000414000006cc0010009c000006cc01008041000000c00110021000000716011001c700008005020000391b2c1b270000040f0000000100200190000013b30000613d0000001102000029000000ff0220018f000000000301043b000000ff0430018f000000000142004b000013b40000c13d0000000d01000029001100000001001d0000001001000029000000800110008a001000000001001d0000000201100367000000000201043b000006cf0020009c000000500000213d0000000e0100002900000011030000291b2c15f20000040f00000010010000290000000201100367000000000601043b000006cf0060009c000000500000213d000000400100043d00000011020000290000000000210435000006cc0010009c000006cc0100804100000040011002100000000002000414000006cc0020009c000006cc02008041000000c002200210000000000112019f000006db011001c70000800d0200003900000003030000390000075a0400004100000000050004111b2c1b220000040f0000000100200190000000500000613d000000400100043d001000000001001d1b2c14190000040f000000110200002900000010010000290000000000210435000000400100043d0000000000210435000006cc0010009c000006cc0100804100000040011002100000074b011001c700001b2d0001042e000000000001042f000013bd0000a13d000001000010008c000013c00000813d0000004d0010008c000013ee0000213d000000000001004b000013cc0000c13d0000000102000039000013d70000013d0000000001240049000000ff0010008c000013c60000a13d0000074101000041000000000010043f0000001101000039000000040010043f000007180100004100001b2e000104300000004e0010008c000013ee0000813d000000000001004b000013d90000c13d0000000102000039000013eb0000013d0000000a030000390000000102000039000000010010019000000000043300a9000000010300603900000000022300a900000001011002720000000003040019000013ce0000c13d000000000002004b000013e50000613d0000000d012000f9000013830000013d0000000a0500003900000001020000390000000004010019000000010040019000000000065500a9000000010500603900000000022500a900000001044002720000000005060019000013dc0000c13d000000000002004b000013eb0000c13d0000074101000041000000000010043f0000001201000039000000040010043f000007180100004100001b2e000104300000076a022001290000000d0020006c000014000000813d000000400200043d001000000002001d00000759010000410000000000120435000000040120003900000011020000290000000d040000291b2c16cc0000040f00000010020000290000000001210049000006cc0010009c000006cc010080410000006001100210000006cc0020009c000006cc020080410000004002200210000000000121019f00001b2e00010430000000ff0210018f0000004d0020008c000013c00000213d000000000002004b000014070000c13d0000000101000039000014100000013d0000000a030000390000000101000039000000010020019000000000043300a9000000010300603900000000011300a900000001022002720000000003040019000014090000c13d0000000d0000006b001100000000001d000013840000613d0000000d031000b9001100000003001d0000000d023000fa000000000012004b000013840000613d000013c00000013d0000076c0010009c0000141e0000813d0000002001100039000000400010043f000000000001042d0000074101000041000000000010043f0000004101000039000000040010043f000007180100004100001b2e000104300000076d0010009c000014290000813d0000004001100039000000400010043f000000000001042d0000074101000041000000000010043f0000004101000039000000040010043f000007180100004100001b2e000104300000076e0010009c000014340000813d000000a001100039000000400010043f000000000001042d0000074101000041000000000010043f0000004101000039000000040010043f000007180100004100001b2e0001043000000000430104340000000001320436000000000003004b000014460000613d000000000200001900000000052100190000000006240019000000000606043300000000006504350000002002200039000000000032004b0000143f0000413d000000000231001900000000000204350000001f0230003900000769022001970000000001210019000000000001042d000007200010009c000014670000213d000000430010008c000014670000a13d00000002020003670000000403200370000000000403043b000006d00040009c000014670000213d0000002403200370000000000503043b000006d00050009c000014670000213d0000002303500039000000000013004b000014670000813d0000000403500039000000000232034f000000000302043b000006d00030009c000014670000213d00000024025000390000000005320019000000000015004b000014670000213d0000000001040019000000000001042d000000000100001900001b2e000104300000000043020434000007220330019700000000033104360000000004040433000006cc04400197000000000043043500000040032000390000000003030433000000000003004b0000000003000039000000010300c039000000400410003900000000003404350000006003200039000000000303043300000722033001970000006004100039000000000034043500000080022000390000000002020433000007220220019700000080031000390000000000230435000000a001100039000000000001042d0000000204000367000000000224034f000000000202043b000000000300003100000000051300490000001f0550008a000006d106500197000006d107200197000000000867013f000000000067004b0000000006000019000006d106002041000000000052004b0000000005000019000006d105004041000006d10080009c000000000605c019000000000006004b000014aa0000613d0000000001120019000000000214034f000000000202043b000006d00020009c000014aa0000213d00000000032300490000002001100039000006d104300197000006d105100197000000000645013f000000000045004b0000000004000019000006d104004041000000000031004b0000000003000019000006d103002041000006d10060009c000000000403c019000000000004004b000014aa0000c13d000000000001042d000000000100001900001b2e000104300000076f0020009c000014dc0000813d00000000040100190000001f0120003900000769011001970000003f011000390000076905100197000000400100043d0000000005510019000000000015004b00000000070000390000000107004039000006d00050009c000014dc0000213d0000000100700190000014dc0000c13d000000400050043f00000000052104360000000007420019000000000037004b000014e20000213d00000769062001980000001f0720018f00000002044003670000000003650019000014cc0000613d000000000804034f0000000009050019000000008a08043c0000000009a90436000000000039004b000014c80000c13d000000000007004b000014d90000613d000000000464034f0000000306700210000000000703043300000000076701cf000000000767022f000000000404043b0000010006600089000000000464022f00000000046401cf000000000474019f000000000043043500000000022500190000000000020435000000000001042d0000074101000041000000000010043f0000004101000039000000040010043f000007180100004100001b2e00010430000000000100001900001b2e000104300003000000000002000300000003001d000200000002001d000006d001100197000000000010043f0000000701000039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000207000029000000030a00002900000001002001900000154d0000613d0000000003000031000000000601043b0000076f00a0009c0000154f0000813d0000001f01a0003900000769011001970000003f011000390000076902100197000000400100043d0000000002210019000000000012004b00000000050000390000000105004039000006d00020009c0000154f0000213d00000001005001900000154f0000c13d000100000006001d000000400020043f0000000002a1043600000000057a0019000000000035004b0000154d0000213d0000076904a001980000001f05a0018f00000002067003670000000003420019000015180000613d000000000706034f0000000008020019000000007907043c0000000008980436000000000038004b000015140000c13d000000000005004b000015250000613d000000000446034f0000000305500210000000000603043300000000065601cf000000000656022f000000000404043b0000010005500089000000000454022f00000000045401cf000000000464019f00000000004304350000000003a200190000000000030435000006cc0020009c000006cc0200804100000040022002100000000001010433000006cc0010009c000006cc010080410000006001100210000000000121019f0000000002000414000006cc0020009c000006cc02008041000000c002200210000000000112019f00000711011001c700008010020000391b2c1b270000040f00000001002001900000154d0000613d000000000101043b000000000010043f00000001010000290000000601100039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f00000001002001900000154d0000613d000000000101043b000000000101041a000000000001004b0000000001000039000000010100c039000000000001042d000000000100001900001b2e000104300000074101000041000000000010043f0000004101000039000000040010043f000007180100004100001b2e00010430000000000323043600000769062001980000001f0720018f00000000056300190000000201100367000015610000613d000000000801034f0000000009030019000000008a08043c0000000009a90436000000000059004b0000155d0000c13d000000000007004b0000156e0000613d000000000161034f0000000306700210000000000705043300000000076701cf000000000767022f000000000101043b0000010006600089000000000161022f00000000016101cf000000000171019f0000000000150435000000000123001900000000000104350000001f0120003900000769011001970000000001130019000000000001042d000000400100043d0000076e0010009c000015830000813d000000a002100039000000400020043f000000800210003900000000000204350000006002100039000000000002043500000040021000390000000000020435000000200210003900000000000204350000000000010435000000000001042d0000074101000041000000000010043f0000004101000039000000040010043f000007180100004100001b2e000104300003000000000002000006d001100197000000000010043f0000000701000039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000100200190000015d60000613d000000000101043b0000000405100039000000000205041a000000010320019000000001062002700000007f0660618f0000001f0060008c00000000040000390000000104002039000000000043004b000015d80000c13d000000400100043d0000000004610436000000000003004b000015c20000613d000100000004001d000200000006001d000300000001001d000000000050043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006db011001c700008010020000391b2c1b270000040f0000000100200190000015d60000613d0000000206000029000000000006004b000015c80000613d000000000201043b0000000005000019000000030100002900000001070000290000000003570019000000000402041a000000000043043500000001022000390000002005500039000000000065004b000015ba0000413d000015ca0000013d0000076b022001970000000000240435000000000006004b00000020050000390000000005006039000015ca0000013d000000000500001900000003010000290000003f0350003900000769023001970000000003120019000000000023004b00000000020000390000000102004039000006d00030009c000015de0000213d0000000100200190000015de0000c13d000000400030043f000000000001042d000000000100001900001b2e000104300000074101000041000000000010043f0000002201000039000000040010043f000007180100004100001b2e000104300000074101000041000000000010043f0000004101000039000000040010043f000007180100004100001b2e000104300000000043010434000000000003004b0000000003000039000000010300c039000000000332043600000000040404330000072204400197000000000043043500000040022000390000004001100039000000000101043300000722011001970000000000120435000000000001042d0002000000000002000000400400043d00000044054000390000000000350435000000200340003900000770050000410000000000530435000006cf022001970000002405400039000000000025043500000044020000390000000000240435000007710040009c000016800000813d0000008009400039000000400090043f000007720040009c000016800000213d000006cf0a100197000000c001400039000000400010043f00000020010000390000000000190435000000a00140003900000773020000410000000000210435000000000204043300000000010004140000000400a0008c0000163e0000c13d00000001020000390000000101000031000000000001004b000016560000613d000006d00010009c000016800000213d0000001f0410003900000769044001970000003f044000390000076904400197000000400c00043d00000000044c00190000000000c4004b00000000050000390000000105004039000006d00040009c000016800000213d0000000100500190000016800000c13d000000400040043f000000000b1c043600000769031001980000001f0410018f00000000013b00190000000305000367000016300000613d000000000605034f00000000070b0019000000006806043c0000000007870436000000000017004b0000162c0000c13d000000000004004b000016580000613d000000000335034f0000000304400210000000000501043300000000054501cf000000000545022f000000000303043b0000010004400089000000000343022f00000000034301cf000000000353019f0000000000310435000016580000013d000006cc0030009c000006cc030080410000004003300210000006cc0020009c000006cc020080410000006002200210000000000232019f000006cc0010009c000006cc01008041000000c001100210000000000112019f00000000020a0019000200000009001d00010000000a001d1b2c1b220000040f000000010a0000290000000209000029000000010220018f00030000000103550000006001100270000106cc0010019d000006cc01100197000000000001004b000016140000c13d000000600c000039000000800b00003900000000030c0433000000000002004b000016880000613d000000000003004b000016730000c13d00020000000c001d00010000000b001d000007500100004100000000001004430000000400a004430000000001000414000006cc0010009c000006cc01008041000000c00110021000000751011001c700008002020000391b2c1b270000040f0000000100200190000016ba0000613d000000000101043b000000000001004b0000000201000029000016bb0000613d0000000003010433000000000003004b000000010b0000290000167f0000613d000007200030009c000016860000213d0000001f0030008c000016860000a13d00000000010b0433000000000001004b0000000002000039000000010200c039000000000021004b000016860000c13d000000000001004b0000169e0000613d000000000001042d0000074101000041000000000010043f0000004101000039000000040010043f000007180100004100001b2e00010430000000000100001900001b2e00010430000000000003004b000016b20000c13d0000000001090019000000400400043d000200000004001d0000077602000041000000000024043500000004034000390000002002000039000000000023043500000024024000391b2c143a0000040f00000002020000290000000001210049000006cc0010009c000006cc01008041000006cc0020009c000006cc0200804100000060011002100000004002200210000000000121019f00001b2e00010430000000400100043d00000064021000390000077403000041000000000032043500000044021000390000077503000041000000000032043500000024021000390000002a03000039000000000032043500000776020000410000000000210435000000040210003900000020030000390000000000320435000006cc0010009c000006cc01008041000000400110021000000777011001c700001b2e00010430000006cc00b0009c000006cc0b0080410000004002b00210000006cc0030009c000006cc030080410000006001300210000000000121019f00001b2e00010430000000000001042f000000400100043d00000044021000390000077803000041000000000032043500000024021000390000001d03000039000000000032043500000776020000410000000000210435000000040210003900000020030000390000000000320435000006cc0010009c000006cc01008041000000400110021000000724011001c700001b2e0001043000000040051000390000000000450435000000ff0330018f00000020041000390000000000340435000000ff0220018f00000000002104350000006001100039000000000001042d0000000101000039000000000101041a000006cf011001970000000002000411000000000012004b000016dc0000c13d000000000001042d000000400100043d00000754020000410000000000210435000006cc0010009c000006cc010080410000004001100210000006d5011001c700001b2e000104300007000000000002000400000001001d000600000002001d0000000021020434000000000001004b000017fd0000613d000006cc0010009c000006cc010080410000006001100210000006cc0020009c000500000002001d000006cc020080410000004002200210000000000121019f0000000002000414000006cc0020009c000006cc02008041000000c002200210000000000112019f00000711011001c700008010020000391b2c1b270000040f0000000100200190000017f50000613d000000000101043b000700000001001d0000000401000029000006d001100197000200000001001d000000000010043f0000000701000039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000100200190000017f50000613d000000000201043b0000000701000029000000000010043f000400000002001d0000000601200039000300000001001d000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000100200190000017f50000613d000000000101043b000000000101041a000000000001004b000018050000c13d00000004010000290000000502100039000000000102041a0000076f0010009c000017f70000813d000100000001001d0000000101100039000000000012041b000400000002001d000000000020043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006db011001c700008010020000391b2c1b270000040f00000001002001900000000702000029000017f50000613d000000000101043b0000000101100029000000000021041b0000000401000029000000000101041a000400000001001d000000000020043f0000000301000029000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f00000001002001900000000702000029000017f50000613d000000000101043b0000000403000029000000000031041b000000000020043f0000000801000039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000100200190000017f50000613d000000000801043b00000006010000290000000004010433000006d00040009c000017f70000213d000000000108041a000000010210019000000001031002700000007f0330618f0000001f0030008c00000000010000390000000101002039000000000012004b00000005070000290000181c0000c13d000000200030008c000400000008001d000700000004001d000017880000413d000300000003001d000000000080043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006db011001c700008010020000391b2c1b270000040f0000000100200190000017f50000613d00000007040000290000001f024000390000000502200270000000200040008c0000000002004019000000000301043b00000003010000290000001f01100039000000050110027000000000011300190000000002230019000000000012004b00000005070000290000000408000029000017880000813d000000000002041b0000000102200039000000000012004b000017840000413d0000001f0040008c000000200a00008a000000200b000039000017b80000a13d000000000080043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006db011001c700008010020000391b2c1b270000040f0000000100200190000017f50000613d0000000709000029000000200a00008a0000000002a90170000000000101043b000000200b000039000017ee0000613d000000010320008a000000050330027000000000043100190000002003000039000000010440003900000005070000290000000606000029000000040800002900000000056300190000000005050433000000000051041b00000020033000390000000101100039000000000041004b000017a40000c13d000000000092004b000017b50000813d0000000302900210000000f80220018f0000076a0220027f0000076a0220016700000000036300190000000003030433000000000223016f000000000021041b000000010190021000000001011001bf000017c40000013d000000000004004b000017bc0000613d0000000001070433000017bd0000013d0000000001000019000000060600002900000003024002100000076a0220027f0000076a02200167000000000121016f0000000102400210000000000121019f000000000018041b000000400100043d0000000003b10436000000000206043300000000002304350000004003100039000000000002004b000017d40000613d000000000400001900000000053400190000000006740019000000000606043300000000006504350000002004400039000000000024004b000017cd0000413d0000001f042000390000000004a4016f000000000223001900000000000204350000004002400039000006cc0020009c000006cc020080410000006002200210000006cc0010009c000006cc010080410000004001100210000000000112019f0000000002000414000006cc0020009c000006cc02008041000000c002200210000000000121019f00000711011001c70000800d0200003900000002030000390000072d0400004100000002050000291b2c1b220000040f0000000100200190000017f50000613d000000000001042d00000000030b0019000000050700002900000006060000290000000408000029000000000092004b000017ad0000413d000017b50000013d000000000100001900001b2e000104300000074101000041000000000010043f0000004101000039000000040010043f000007180100004100001b2e00010430000000400100043d0000072f020000410000000000210435000006cc0010009c000006cc010080410000004001100210000006d5011001c700001b2e00010430000000400300043d000700000003001d0000002401300039000000400200003900000000002104350000072c010000410000000000130435000000040130003900000002020000290000000000210435000000440230003900000006010000291b2c143a0000040f00000007020000290000000001210049000006cc0010009c000006cc01008041000006cc0020009c000006cc0200804100000060011002100000004002200210000000000121019f00001b2e000104300000074101000041000000000010043f0000002201000039000000040010043f000007180100004100001b2e000104300005000000000002000000400300043d0000076e0030009c000018690000813d000000a002300039000000400020043f00000080023000390000000000020435000000600230003900000000000204350000004002300039000000000002043500000020023000390000000000020435000000000003043500000060021000390000000002020433000100000002001d000500000001001d0000000012010434000300000002001d000200000001001d0000000001010433000400000001001d000007280100004100000000001004430000000001000414000006cc0010009c000006cc01008041000000c00110021000000729011001c70000800b020000391b2c1b270000040f00000001002001900000186f0000613d0000000402000029000006cc04200197000000000601043b000000000346004b0000000501000029000018630000413d00000080021000390000000002020433000007220520019700000000023500a9000000000046004b000018540000613d00000000033200d9000000000053004b000018630000c13d00000003030000290000072203300197000000000032001a000018630000413d0000000002320019000000010300002900000722033001970000072204200197000000000023004b00000000030480190000000000310435000006cc0260019700000002030000290000000000230435000000000001042d0000074101000041000000000010043f0000001101000039000000040010043f000007180100004100001b2e000104300000074101000041000000000010043f0000004101000039000000040010043f000007180100004100001b2e00010430000000000001042f0002000000000002000000400900043d0000076d0090009c000018f00000813d000006cf0a1001970000004001900039000000400010043f00000020019000390000077303000041000000000031043500000020010000390000000000190435000000002302043400000000010004140000000400a0008c000018ae0000c13d00000001020000390000000101000031000000000001004b000018c60000613d000006d00010009c000018f00000213d0000001f0410003900000769044001970000003f044000390000076904400197000000400c00043d00000000044c00190000000000c4004b00000000050000390000000105004039000006d00040009c000018f00000213d0000000100500190000018f00000c13d000000400040043f000000000b1c043600000769031001980000001f0410018f00000000013b00190000000305000367000018a00000613d000000000605034f00000000070b0019000000006806043c0000000007870436000000000017004b0000189c0000c13d000000000004004b000018c80000613d000000000335034f0000000304400210000000000501043300000000054501cf000000000545022f000000000303043b0000010004400089000000000343022f00000000034301cf000000000353019f0000000000310435000018c80000013d000006cc0030009c000006cc030080410000006003300210000006cc0020009c000006cc020080410000004002200210000000000223019f000006cc0010009c000006cc01008041000000c001100210000000000112019f00000000020a0019000200000009001d00010000000a001d1b2c1b220000040f000000010a0000290000000209000029000000010220018f00030000000103550000006001100270000106cc0010019d000006cc01100197000000000001004b000018840000c13d000000600c000039000000800b00003900000000030c0433000000000002004b000018f80000613d000000000003004b000018e30000c13d00020000000c001d00010000000b001d000007500100004100000000001004430000000400a004430000000001000414000006cc0010009c000006cc01008041000000c00110021000000751011001c700008002020000391b2c1b270000040f00000001002001900000192a0000613d000000000101043b000000000001004b00000002010000290000192b0000613d0000000003010433000000000003004b000000010b000029000018ef0000613d000007200030009c000018f60000213d0000001f0030008c000018f60000a13d00000000010b0433000000000001004b0000000002000039000000010200c039000000000021004b000018f60000c13d000000000001004b0000190e0000613d000000000001042d0000074101000041000000000010043f0000004101000039000000040010043f000007180100004100001b2e00010430000000000100001900001b2e00010430000000000003004b000019220000c13d0000000001090019000000400400043d000200000004001d0000077602000041000000000024043500000004034000390000002002000039000000000023043500000024024000391b2c143a0000040f00000002020000290000000001210049000006cc0010009c000006cc01008041000006cc0020009c000006cc0200804100000060011002100000004002200210000000000121019f00001b2e00010430000000400100043d00000064021000390000077403000041000000000032043500000044021000390000077503000041000000000032043500000024021000390000002a03000039000000000032043500000776020000410000000000210435000000040210003900000020030000390000000000320435000006cc0010009c000006cc01008041000000400110021000000777011001c700001b2e00010430000006cc00b0009c000006cc0b0080410000004002b00210000006cc0030009c000006cc030080410000006001300210000000000121019f00001b2e00010430000000000001042f000000400100043d00000044021000390000077803000041000000000032043500000024021000390000001d03000039000000000032043500000776020000410000000000210435000000040210003900000020030000390000000000320435000006cc0010009c000006cc01008041000000400110021000000724011001c700001b2e00010430000000000010043f0000000601000039000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f00000001002001900000194e0000613d000000000101043b000000000101041a000000000001004b0000000001000039000000010100c039000000000001042d000000000100001900001b2e000104300006000000000002000300000002001d000000000020043f000600000001001d0000000101100039000400000001001d000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000100200190000019c80000613d0000000603000029000000000101043b000000000101041a000000000001004b000019c60000613d000000000203041a000000000002004b000019ca0000613d000000000021004b000500000001001d000019a40000613d000200000002001d000000000030043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006db011001c700008010020000391b2c1b270000040f0000000100200190000019c80000613d00000005020000290001000100200092000000000101043b0000000604000029000000000204041a000000010020006c000019d00000a13d0000000202000029000000010220008a0000000001120019000000000101041a000200000001001d000000000040043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006db011001c700008010020000391b2c1b270000040f0000000100200190000019c80000613d000000000101043b00000001011000290000000202000029000000000021041b000000000020043f0000000401000029000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000100200190000019c80000613d000000000101043b0000000502000029000000000021041b0000000603000029000000000103041a000500000001001d000000000001004b000019d60000613d000000000030043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006db011001c700008010020000391b2c1b270000040f0000000100200190000019c80000613d0000000502000029000000010220008a000000000101043b0000000001210019000000000001041b0000000601000029000000000021041b0000000301000029000000000010043f0000000401000029000000200010043f0000000001000414000006cc0010009c000006cc01008041000000c001100210000006d9011001c700008010020000391b2c1b270000040f0000000100200190000019c80000613d000000000101043b000000000001041b0000000101000039000000000001042d0000000001000019000000000001042d000000000100001900001b2e000104300000074101000041000000000010043f0000001101000039000000040010043f000007180100004100001b2e000104300000074101000041000000000010043f0000003201000039000000040010043f000007180100004100001b2e000104300000074101000041000000000010043f0000003101000039000000040010043f000007180100004100001b2e000104300003000000000002000100000002001d000300000001001d000000000101041a000200000001001d000007280100004100000000001004430000000001000414000006cc0010009c000006cc01008041000000c00110021000000729011001c70000800b020000391b2c1b270000040f000000010020019000001a3e0000613d00000002080000290000008002800270000006cc03200197000000000201043b000000000532004b000000030700002900001a380000413d0000000101700039000019f70000c13d000000000207041a00001a090000013d000000000301041a000000800630027000000000045600a900000000055400d9000000000065004b00001a380000c13d0000072205800197000000000054001a00001a380000413d00000722033001970000000004540019000000000043004b0000000003048019000006d30480019700000080022002100000077902200197000000000242019f000000000232019f00000001060000290000002003600039000000000403043300000722044001970000072205200197000000000054004b00000000050440190000077a02200197000000000225019f0000000005060433000000000005004b00000000050000190000072a0500c041000000000252019f000000000027041b000000400260003900000000050204330000008005500210000000000445019f000000000041041b0000000001000039000000010100c039000000400400043d00000000011404360000000003030433000007220330019700000000003104350000000001020433000007220110019700000040024000390000000000120435000006cc0040009c000006cc0400804100000040014002100000000002000414000006cc0020009c000006cc02008041000000c002200210000000000112019f0000077b011001c70000800d0200003900000001030000390000077c040000411b2c1b220000040f000000010020019000001a3f0000613d000000000001042d0000074101000041000000000010043f0000001101000039000000040010043f000007180100004100001b2e00010430000000000001042f000000000100001900001b2e00010430000006cf04400197000000400510003900000000004504350000002004100039000000000034043500000000002104350000006001100039000000000001042d0006000000000002000000000401041a000007340040019800001a9d0000613d000000000002004b00001a9d0000613d000600000004001d000500000002001d000200000003001d000300000001001d0000000101100039000100000001001d000000000101041a000400000001001d000007280100004100000000001004430000000001000414000006cc0010009c000006cc01008041000000c00110021000000729011001c70000800b020000391b2c1b270000040f000000010020019000001a9e0000613d00000006030000290000008002300270000006cc02200197000000000101043b000000000421004b00001ab90000413d00000722033001970000000405000029000007220250019700001a6f0000c13d0000000504000029000000030500002900001a830000013d000000000023004b00001ac10000213d000000800650027000000000056400a900000000044500d9000000000064004b00001ab90000c13d000000000035001a00001ab90000413d0000000003350019000000800110021000000779011001970000000305000029000000000405041a0000077d04400197000000000114019f000000000015041b000000000032004b00000000030240190000000504000029000000000042004b00001a9f0000413d000000000143004b00001ab00000413d0000072201100197000000000205041a0000077f02200197000000000112019f000000000015041b000000400100043d0000000000410435000006cc0010009c000006cc0100804100000040011002100000000002000414000006cc0020009c000006cc02008041000000c002200210000000000112019f000006db011001c70000800d02000039000000010300003900000780040000411b2c1b220000040f000000010020019000001abf0000613d000000000001042d000000000001042f000000400100043d000000000401001900000004011000390000000203000029000006cf0030019800001ac90000c13d000007840300004100000000003404350000000000210435000000240140003900000005020000290000000000210435000006cc0040009c000006cc040080410000004001400210000006d7011001c700001b2e000104300000000101000029000000000101041a000000800110027200001ab90000613d00000005043000690000000002140019000000010220008a000000000042004b00001ace0000813d0000074101000041000000000010043f0000001101000039000000040010043f000007180100004100001b2e00010430000000000100001900001b2e00010430000000400100043d0000077e020000410000000000210435000006cc0010009c000006cc010080410000004001100210000006d5011001c700001b2e000104300000078303000041000600000004001d0000000000340435000000050300002900001ae20000013d00000000021200d9000000400100043d000000000501001900000004011000390000000204000029000006cf0040019800001adf0000c13d00000782040000410000000000450435000000000021043500000024015000390000000000310435000006cc0050009c000006cc050080410000004001500210000006d7011001c700001b2e000104300000078104000041000600000005001d000000000045043500000002040000291b2c1a410000040f00000006020000290000000001210049000006cc0010009c000006cc010080410000006001100210000006cc0020009c000006cc020080410000004002200210000000000121019f00001b2e00010430000000000001042f000006cc0010009c000006cc010080410000004001100210000006cc0020009c000006cc020080410000006002200210000000000112019f0000000002000414000006cc0020009c000006cc02008041000000c002200210000000000112019f00000711011001c700008010020000391b2c1b270000040f000000010020019000001b020000613d000000000101043b000000000001042d000000000100001900001b2e0001043000000000050100190000000000200443000000050030008c00001b120000413d000000040100003900000000020000190000000506200210000000000664001900000005066002700000000006060031000000000161043a0000000102200039000000000031004b00001b0a0000413d000006cc0030009c000006cc0300804100000060013002100000000002000414000006cc0020009c000006cc02008041000000c002200210000000000112019f00000785011001c700000000020500191b2c1b270000040f000000010020019000001b210000613d000000000101043b000000000001042d000000000001042f00001b25002104210000000102000039000000000001042d0000000002000019000000000001042d00001b2a002104230000000102000039000000000001042d0000000002000019000000000001042d00001b2c0000043200001b2d0001042e00001b2e00010430000000000000000000000000000000000000000000000000000000000000000000000000ffffffff00000000000000000000000000000000000000000000000000000001ffffffe000000000000000000000000000000000000000000000000000000000ffffffe0000000000000000000000000ffffffffffffffffffffffffffffffffffffffff000000000000000000000000000000000000000000000000ffffffffffffffff80000000000000000000000000000000000000000000000000000000000000007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0ffffffffffffffffffffffff0000000000000000000000000000000000000000313ce567000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000655a7c0e000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000044000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffffdf0200000000000000000000000000000000000040000000000000000000000000bfa87805ed57dc1f0d489ce33be4c4577d74ccde357eeeee058a32c55c44a53202000000000000000000000000000000000000200000000000000000000000002640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d800000002000000000000000000000000000001800000010000000000000000009b15e16f00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008da5cb5a00000000000000000000000000000000000000000000000000000000c0d7865400000000000000000000000000000000000000000000000000000000dc0bd97000000000000000000000000000000000000000000000000000000000e8a1da1600000000000000000000000000000000000000000000000000000000e8a1da1700000000000000000000000000000000000000000000000000000000eb521a4c00000000000000000000000000000000000000000000000000000000f2fde38b00000000000000000000000000000000000000000000000000000000dc0bd97100000000000000000000000000000000000000000000000000000000e0351e1300000000000000000000000000000000000000000000000000000000c75eea9b00000000000000000000000000000000000000000000000000000000c75eea9c00000000000000000000000000000000000000000000000000000000cf7401f300000000000000000000000000000000000000000000000000000000c0d7865500000000000000000000000000000000000000000000000000000000c4bffe2b00000000000000000000000000000000000000000000000000000000acfecf9000000000000000000000000000000000000000000000000000000000b0f479a000000000000000000000000000000000000000000000000000000000b0f479a100000000000000000000000000000000000000000000000000000000b794658000000000000000000000000000000000000000000000000000000000bb98546b00000000000000000000000000000000000000000000000000000000acfecf9100000000000000000000000000000000000000000000000000000000af58d59f00000000000000000000000000000000000000000000000000000000a42a7b8a00000000000000000000000000000000000000000000000000000000a42a7b8b00000000000000000000000000000000000000000000000000000000a7cd63b7000000000000000000000000000000000000000000000000000000008da5cb5b000000000000000000000000000000000000000000000000000000009a4575b9000000000000000000000000000000000000000000000000000000004c5ef0ec000000000000000000000000000000000000000000000000000000006cfd15520000000000000000000000000000000000000000000000000000000079ba50960000000000000000000000000000000000000000000000000000000079ba5097000000000000000000000000000000000000000000000000000000007d54534e000000000000000000000000000000000000000000000000000000008926f54f000000000000000000000000000000000000000000000000000000006cfd1553000000000000000000000000000000000000000000000000000000006d3d1a580000000000000000000000000000000000000000000000000000000062ddd3c30000000000000000000000000000000000000000000000000000000062ddd3c40000000000000000000000000000000000000000000000000000000066320087000000000000000000000000000000000000000000000000000000004c5ef0ed0000000000000000000000000000000000000000000000000000000054c8a4f300000000000000000000000000000000000000000000000000000000240028e70000000000000000000000000000000000000000000000000000000039077536000000000000000000000000000000000000000000000000000000003907753700000000000000000000000000000000000000000000000000000000432a6ba300000000000000000000000000000000000000000000000000000000240028e80000000000000000000000000000000000000000000000000000000024f65ee700000000000000000000000000000000000000000000000000000000181f5a7600000000000000000000000000000000000000000000000000000000181f5a770000000000000000000000000000000000000000000000000000000021df0da70000000000000000000000000000000000000000000000000000000001ffc9a7000000000000000000000000000000000000000000000000000000000a861f2a0200000000000000000000000000000000000000000000000000000000000000ed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278dad89dca000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000800000000000000000310ab089e4439a4c15d089f94afb7896ff553aecb10793d0ab882de59d99a32e02000002000000000000000000000000000000440000000000000000000000008e4a23d600000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002400000000000000000000000023b872dd00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffff5fc17cea59c2955cb181b03393209566960365771dbba9dc3d510180e7cb312088e93f8fa400000000000000000000000000000000000000000000000000000000fc949c7b4a13586e39d89eead2f38644f9fb3efb5a0490b14f8fc0ceab44c2515204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d8599161e670e4b000000000000000000000000000000000000000000000000000000007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000000000000000000000000000000000000000000000ffffffffffffff9f00000000000000000000000000000000ffffffffffffffffffffffffffffffff8020d124000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000064000000000000000000000000d68af9cc000000000000000000000000000000000000000000000000000000001d5ad3c500000000000000000000000000000000000000000000000000000000fc949c7b4a13586e39d89eead2f38644f9fb3efb5a0490b14f8fc0ceab44c250796b89b91644bc98cd93958e4c9038275d622183e25ac5af08cc6b5d9553913202000002000000000000000000000000000000040000000000000000000000000000000000000000000000010000000000000000000000000000000000000000ffffffffffffffffffffff000000000000000000000000000000000000000000393b8ad2000000000000000000000000000000000000000000000000000000007d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c28579befe000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020000000800000000000000000000000000000000000000000000000000000002400000140000000000000000002000000000000000000000000000000000000e00000000000000000000000000350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b0000000000000000000000ff0000000000000000000000000000000000000000036b6384b5eca791c62761152d0c79bb0604c104a5fb6f4eb0703f3154bb3db0000000000000000000000000000000000000000000000000ffffffffffffffc1ffffffffffffffffffffffffffffffffffffffffffffffff000000000000008000000000000000000000000000000000000000000000003fffffffffffffffe0020000000000000000000000000000000000004000000080000000000000000002dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f168400000000000000000000000000000000000000000000000000000000ffffffbffdffffffffffffffffffffffffffffffffffffc000000000000000000000000052d00ee4d9bd51b40168f2afc5848837288ce258784ad914278791464b3f4d7674f23c7c00000000000000000000000000000000000000000000000000000000405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace000000000000000000000000000000000000000000000000ffffffffffffff7f4e487b7100000000000000000000000000000000000000000000000000000000961c9a4f000000000000000000000000000000000000000000000000000000002cbc26bb000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffffff0000000000000000000000000000000053ad11d800000000000000000000000000000000000000000000000000000000d0d2597600000000000000000000000000000000000000000000000000000000a8d87a3b00000000000000000000000000000000000000000000000000000000728fe07b000000000000000000000000000000000000000000000000000000009f1ec8c880f76798e7b793325d625e9b60e4082a553c98f42b6cda368dd60008a9902c7e000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000020000000000000000000000000000000000002000000080000000000000000044676b5284b809a22248eba0da87391d79098be38bb03154be88a58bf4d0917402b543c6000000000000000000000000000000000000000000000000000000008be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e01806aa1896bbf26568e884a7374b41e002500962caba6a15023a8d90e8508b8302000002000000000000000000000000000000240000000000000000000000000a861f2a000000000000000000000000000000000000000000000000000000006fa7abcf1345d1d478e5ea0da6b5f26a90eadb0546ef15ed3833944fbfd1db622b5c74de00000000000000000000000000000000000000000000000000000000bfa87805ed57dc1f0d489ce33be4c4577d74ccde357eeeee058a32c55c44a533800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf756635f4a7b30000000000000000000000000000000000000000000000000000000083826b2b00000000000000000000000000000000000000000000000000000000a9cb113d000000000000000000000000000000000000000000000000000000002d87480f50083e2b2759522a8fdda59802650a8055e609a7772cf70c07748f52953576f70000000000000000000000000000000000000000000000000000000024eb47e5000000000000000000000000000000000000000000000000000000004c6f636b52656c65617365546f6b656e506f6f6c20312e352e310000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000002400000080000000000000000070a0823100000000000000000000000000000000000000000000000000000000c2c3f06e49b9f15e7b4af9055e183b0d73362e033ad82a07dec9bf9840171719bb55fd270000000000000000000000000000000000000000000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffaff2afbeffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1d4056600000000000000000000000000000000000000000000000000000000aff2afbf0000000000000000000000000000000000000000000000000000000001ffc9a7000000000000000000000000000000000000000000000000000000000e64dd2900000000000000000000000000000000000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000000000000000000000000000000000000000000000ffffffffffffffe0000000000000000000000000000000000000000000000000ffffffffffffffc0000000000000000000000000000000000000000000000000ffffffffffffff600000000000000000000000000000000000000000000000010000000000000000a9059cbb00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffff80000000000000000000000000000000000000000000000000ffffffffffffff3f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65646f742073756363656564000000000000000000000000000000000000000000005361666545524332303a204552433230206f7065726174696f6e20646964206e08c379a0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000084000000000000000000000000416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000000000000000000000000000ffffffff00000000000000000000000000000000ffffffffffffffffffffff00ffffffff0000000000000000000000000000000002000000000000000000000000000000000000600000000000000000000000009ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c19ffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff9725942a00000000000000000000000000000000000000000000000000000000ffffffffffffffffffffffffffffffff000000000000000000000000000000001871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690ad0c8d23a0000000000000000000000000000000000000000000000000000000015279c08000000000000000000000000000000000000000000000000000000001a76572a00000000000000000000000000000000000000000000000000000000f94ebcd10000000000000000000000000000000000000000000000000000000002000002000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")

func DeployLockReleaseTokenPoolZK(auth *bind.TransactOpts, backend bind.ContractBackend, token common.Address, localTokenDecimals uint8, allowlist []common.
	Address, rmnProxy common.Address, acceptLiquidity bool, router common.Address) (common.Address, *generated.Transaction, *LockReleaseTokenPool, error) {
	parsed, err := LockReleaseTokenPoolMetaData.GetAbi()
	if err != nil {
		return common.Address{},
			nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")

	}
	address,
		ethTx, contract,
		err := generated.
		DeployContract(
			auth, parsed, common.FromHex(LockReleaseTokenPoolZKBin),
			backend, token, localTokenDecimals,
			allowlist, rmnProxy,
			acceptLiquidity,

			router)
	if err != nil {
		return common.Address{}, nil,
			nil, err
	}
	return address,
		ethTx, &LockReleaseTokenPool{address: address, abi: *parsed, LockReleaseTokenPoolCaller: LockReleaseTokenPoolCaller{contract: contract},
			LockReleaseTokenPoolTransactor: LockReleaseTokenPoolTransactor{contract: contract}, LockReleaseTokenPoolFilterer: LockReleaseTokenPoolFilterer{contract: contract}}, nil
}
