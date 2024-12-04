package usdc_token_pool

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

type USDCTokenPoolDomain struct {
	AllowedCaller    [32]byte
	DomainIdentifier uint32
	Enabled          bool
}

type USDCTokenPoolDomainUpdate struct {
	AllowedCaller     [32]byte
	DomainIdentifier  uint32
	DestChainSelector uint64
	Enabled           bool
}

var USDCTokenPoolMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractITokenMessenger\",\"name\":\"tokenMessenger\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"allowlist\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"}],\"name\":\"AggregateValueMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"}],\"name\":\"AggregateValueRateLimitReached\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AllowListNotEnabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BucketOverfilled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"CallerIsNotARampOnRouter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotTransferToSelf\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CursedByRMN\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"DisabledNonZeroRateLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"expected\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"actual\",\"type\":\"uint8\"}],\"name\":\"InvalidDecimalArgs\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"expected\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"got\",\"type\":\"uint32\"}],\"name\":\"InvalidDestinationDomain\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"allowedCaller\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"domainIdentifier\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"internalType\":\"structUSDCTokenPool.DomainUpdate\",\"name\":\"domain\",\"type\":\"tuple\"}],\"name\":\"InvalidDomain\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"}],\"name\":\"InvalidMessageVersion\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"expected\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"got\",\"type\":\"uint64\"}],\"name\":\"InvalidNonce\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"rateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"InvalidRateLimitRate\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"}],\"name\":\"InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolData\",\"type\":\"bytes\"}],\"name\":\"InvalidRemoteChainDecimals\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidRemotePoolForChain\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"expected\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"got\",\"type\":\"uint32\"}],\"name\":\"InvalidSourceDomain\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidSourcePoolAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"InvalidToken\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"}],\"name\":\"InvalidTokenMessengerVersion\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustBeProposedOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"NonExistentChain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"remoteDecimals\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"localDecimals\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"remoteAmount\",\"type\":\"uint256\"}],\"name\":\"OverflowDetected\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OwnerCannotBeZero\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"PoolAlreadyAdded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RateLimitMustBeDisabled\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"SenderNotAllowed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"capacity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenMaxCapacityExceeded\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minWaitInSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"TokenRateLimitReached\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"domain\",\"type\":\"uint64\"}],\"name\":\"UnknownDomain\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnlockingUSDCFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListAdd\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"AllowListRemove\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Burned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remoteToken\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"ChainAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"name\":\"ChainConfigured\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"ChainRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"indexed\":false,\"internalType\":\"structRateLimiter.Config\",\"name\":\"config\",\"type\":\"tuple\"}],\"name\":\"ConfigChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenMessenger\",\"type\":\"address\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"allowedCaller\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"domainIdentifier\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"indexed\":false,\"internalType\":\"structUSDCTokenPool.DomainUpdate[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"name\":\"DomainsSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Locked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Minted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"rateLimitAdmin\",\"type\":\"address\"}],\"name\":\"RateLimitAdminSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Released\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"RemotePoolAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"RemotePoolRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldRouter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"RouterUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"TokensConsumed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"SUPPORTED_USDC_VERSION\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"addRemotePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"removes\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"adds\",\"type\":\"address[]\"}],\"name\":\"applyAllowListUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64[]\",\"name\":\"remoteChainSelectorsToRemove\",\"type\":\"uint64[]\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes[]\",\"name\":\"remotePoolAddresses\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"remoteTokenAddress\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundRateLimiterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundRateLimiterConfig\",\"type\":\"tuple\"}],\"internalType\":\"structTokenPool.ChainUpdate[]\",\"name\":\"chainsToAdd\",\"type\":\"tuple[]\"}],\"name\":\"applyChainUpdates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowList\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowListEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getCurrentInboundRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getCurrentOutboundRateLimiterState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"tokens\",\"type\":\"uint128\"},{\"internalType\":\"uint32\",\"name\":\"lastUpdated\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.TokenBucket\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainSelector\",\"type\":\"uint64\"}],\"name\":\"getDomain\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"allowedCaller\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"domainIdentifier\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"internalType\":\"structUSDCTokenPool.Domain\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRateLimitAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getRemotePools\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"\",\"type\":\"bytes[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"getRemoteToken\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRmnProxy\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"rmnProxy\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRouter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSupportedChains\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTokenDecimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"decimals\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_localDomainIdentifier\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_messageTransmitter\",\"outputs\":[{\"internalType\":\"contractIMessageTransmitter\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"i_tokenMessenger\",\"outputs\":[{\"internalType\":\"contractITokenMessenger\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"isRemotePool\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"}],\"name\":\"isSupportedChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"isSupportedToken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"originalSender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"}],\"internalType\":\"structPool.LockOrBurnInV1\",\"name\":\"lockOrBurnIn\",\"type\":\"tuple\"}],\"name\":\"lockOrBurn\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"destTokenAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"destPoolData\",\"type\":\"bytes\"}],\"internalType\":\"structPool.LockOrBurnOutV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"originalSender\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"sourcePoolAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"sourcePoolData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"offchainTokenData\",\"type\":\"bytes\"}],\"internalType\":\"structPool.ReleaseOrMintInV1\",\"name\":\"releaseOrMintIn\",\"type\":\"tuple\"}],\"name\":\"releaseOrMint\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"destinationAmount\",\"type\":\"uint256\"}],\"internalType\":\"structPool.ReleaseOrMintOutV1\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"remotePoolAddress\",\"type\":\"bytes\"}],\"name\":\"removeRemotePool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"remoteChainSelector\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"outboundConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isEnabled\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"capacity\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"rate\",\"type\":\"uint128\"}],\"internalType\":\"structRateLimiter.Config\",\"name\":\"inboundConfig\",\"type\":\"tuple\"}],\"name\":\"setChainRateLimiterConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"allowedCaller\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"domainIdentifier\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"destChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"internalType\":\"structUSDCTokenPool.DomainUpdate[]\",\"name\":\"domains\",\"type\":\"tuple[]\"}],\"name\":\"setDomains\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"rateLimitAdmin\",\"type\":\"address\"}],\"name\":\"setRateLimitAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newRouter\",\"type\":\"address\"}],\"name\":\"setRouter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6101606040523480156200001257600080fd5b50604051620054b7380380620054b7833981016040819052620000359162000b93565b836006848484336000816200005d57604051639b15e16f60e01b815260040160405180910390fd5b600180546001600160a01b0319166001600160a01b03848116919091179091558116156200009057620000908162000493565b50506001600160a01b0385161580620000b057506001600160a01b038116155b80620000c357506001600160a01b038216155b15620000e2576040516342bcdf7f60e11b815260040160405180910390fd5b6001600160a01b03808616608081905290831660c0526040805163313ce56760e01b8152905163313ce567916004808201926020929091908290030181865afa92505050801562000152575060408051601f3d908101601f191682019092526200014f9181019062000cb9565b60015b1562000193578060ff168560ff161462000191576040516332ad3e0760e11b815260ff8087166004830152821660248201526044015b60405180910390fd5b505b60ff841660a052600480546001600160a01b0319166001600160a01b038316179055825115801560e052620001dd57604080516000815260208101909152620001dd90846200050d565b5050506001600160a01b03871691506200020c9050576040516306b7c75960e31b815260040160405180910390fd5b6000856001600160a01b0316632c1219216040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200024d573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000273919062000ce5565b90506000816001600160a01b03166354fd4d506040518163ffffffff1660e01b8152600401602060405180830381865afa158015620002b6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620002dc919062000d05565b905063ffffffff8116156200030d576040516334697c6b60e11b815263ffffffff8216600482015260240162000188565b6000876001600160a01b0316639cdbb1816040518163ffffffff1660e01b8152600401602060405180830381865afa1580156200034e573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000374919062000d05565b905063ffffffff811615620003a5576040516316ba39c560e31b815263ffffffff8216600482015260240162000188565b6001600160a01b038089166101005283166101208190526040805163234d8e3d60e21b81529051638d3638f4916004808201926020929091908290030181865afa158015620003f8573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200041e919062000d05565b63ffffffff16610140526101005160805162000449916001600160a01b03909116906000196200066a565b6040516001600160a01b03891681527f2e902d38f15b233cbb63711add0fca4545334d3a169d60c0a616494d7eea95449060200160405180910390a1505050505050505062000e52565b336001600160a01b03821603620004bd57604051636d6c4ee560e11b815260040160405180910390fd5b600080546001600160a01b0319166001600160a01b03838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60e0516200052e576040516335f4a7b360e01b815260040160405180910390fd5b60005b8251811015620005b957600083828151811062000552576200055262000d2d565b602090810291909101015190506200056c60028262000750565b15620005af576040516001600160a01b03821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b5060010162000531565b5060005b815181101562000665576000828281518110620005de57620005de62000d2d565b6020026020010151905060006001600160a01b0316816001600160a01b0316036200060a57506200065c565b6200061760028262000770565b156200065a576040516001600160a01b03821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b600101620005bd565b505050565b604051636eb1769f60e11b81523060048201526001600160a01b038381166024830152600091839186169063dd62ed3e90604401602060405180830381865afa158015620006bc573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620006e2919062000d43565b620006ee919062000d73565b604080516001600160a01b038616602482015260448082018490528251808303909101815260649091019091526020810180516001600160e01b0390811663095ea7b360e01b179091529192506200074a918691906200078716565b50505050565b600062000767836001600160a01b03841662000858565b90505b92915050565b600062000767836001600160a01b0384166200095c565b6040805180820190915260208082527f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c656490820152600090620007d6906001600160a01b038516908490620009ae565b805190915015620006655780806020019051810190620007f7919062000d89565b620006655760405162461bcd60e51b815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e6044820152691bdd081cdd58d8d9595960b21b606482015260840162000188565b60008181526001830160205260408120548015620009515760006200087f60018362000dad565b8554909150600090620008959060019062000dad565b905080821462000901576000866000018281548110620008b957620008b962000d2d565b9060005260206000200154905080876000018481548110620008df57620008df62000d2d565b6000918252602080832090910192909255918252600188019052604090208390555b855486908062000915576200091562000dc3565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506200076a565b60009150506200076a565b6000818152600183016020526040812054620009a5575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556200076a565b5060006200076a565b6060620009bf8484600085620009c7565b949350505050565b60608247101562000a2a5760405162461bcd60e51b815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f6044820152651c8818d85b1b60d21b606482015260840162000188565b600080866001600160a01b0316858760405162000a48919062000dff565b60006040518083038185875af1925050503d806000811462000a87576040519150601f19603f3d011682016040523d82523d6000602084013e62000a8c565b606091505b50909250905062000aa08783838762000aab565b979650505050505050565b6060831562000b1f57825160000362000b17576001600160a01b0385163b62000b175760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000604482015260640162000188565b5081620009bf565b620009bf838381511562000b365781518083602001fd5b8060405162461bcd60e51b815260040162000188919062000e1d565b6001600160a01b038116811462000b6857600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b805162000b8e8162000b52565b919050565b600080600080600060a0868803121562000bac57600080fd5b855162000bb98162000b52565b8095505060208087015162000bce8162000b52565b60408801519095506001600160401b038082111562000bec57600080fd5b818901915089601f83011262000c0157600080fd5b81518181111562000c165762000c1662000b6b565b8060051b604051601f19603f8301168101818110858211171562000c3e5762000c3e62000b6b565b60405291825284820192508381018501918c83111562000c5d57600080fd5b938501935b8285101562000c865762000c768562000b81565b8452938501939285019262000c62565b80985050505050505062000c9d6060870162000b81565b915062000cad6080870162000b81565b90509295509295909350565b60006020828403121562000ccc57600080fd5b815160ff8116811462000cde57600080fd5b9392505050565b60006020828403121562000cf857600080fd5b815162000cde8162000b52565b60006020828403121562000d1857600080fd5b815163ffffffff8116811462000cde57600080fd5b634e487b7160e01b600052603260045260246000fd5b60006020828403121562000d5657600080fd5b5051919050565b634e487b7160e01b600052601160045260246000fd5b808201808211156200076a576200076a62000d5d565b60006020828403121562000d9c57600080fd5b8151801515811462000cde57600080fd5b818103818111156200076a576200076a62000d5d565b634e487b7160e01b600052603160045260246000fd5b60005b8381101562000df657818101518382015260200162000ddc565b50506000910152565b6000825162000e1381846020870162000dd9565b9190910192915050565b602081526000825180602084015262000e3e81604085016020870162000dd9565b601f01601f19169190910160400192915050565b60805160a05160c05160e0516101005161012051610140516145af62000f08600039600081816104170152818161113c01528181612107015261216501526000818161072b0152610a750152600081816103dd01526110520152600081816106dc0152818161221d0152612bcc01526000818161061801528181611eb40152612509015260006103660152600081816102cd015281816103220152818161101c01528181612b620152612db701526145af6000f3fe608060405234801561001057600080fd5b50600436106102405760003560e01c80639a4575b911610145578063c4bffe2b116100bd578063dfadfa351161008c578063e8a1da1711610071578063e8a1da1714610700578063f2fde38b14610713578063fbf84dd71461072657600080fd5b8063dfadfa351461063c578063e0351e13146106da57600080fd5b8063c4bffe2b146105db578063c75eea9c146105f0578063cf7401f314610603578063dc0bd9711461061657600080fd5b8063acfecf9111610114578063b0f479a1116100f9578063b0f479a114610597578063b7946580146105b5578063c0d78655146105c857600080fd5b8063acfecf9114610515578063af58d59f1461052857600080fd5b80639a4575b9146104b85780639fdf13ff146104d8578063a42a7b8b146104e0578063a7cd63b71461050057600080fd5b806354c8a4f3116101d85780636d3d1a58116101a75780637d54534e1161018c5780637d54534e146104745780638926f54f146104875780638da5cb5b1461049a57600080fd5b80636d3d1a581461044e57806379ba50971461046c57600080fd5b806354c8a4f3146103c55780636155cda0146103d857806362ddd3c4146103ff5780636b716b0d1461041257600080fd5b8063240028e811610214578063240028e81461031257806324f65ee71461035f57806339077537146103905780634c5ef0ed146103b257600080fd5b806241d3c11461024557806301ffc9a71461025a578063181f5a771461028257806321df0da7146102cb575b600080fd5b610258610253366004613577565b61074d565b005b61026d6102683660046135ec565b6108ea565b60405190151581526020015b60405180910390f35b6102be6040518060400160405280601381526020017f55534443546f6b656e506f6f6c20312e352e310000000000000000000000000081525081565b6040516102799190613692565b7f00000000000000000000000000000000000000000000000000000000000000005b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610279565b61026d6103203660046136c7565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff90811691161490565b60405160ff7f0000000000000000000000000000000000000000000000000000000000000000168152602001610279565b6103a361039e3660046136e4565b6109cf565b60405190518152602001610279565b61026d6103c0366004613736565b610bb4565b6102586103d3366004613807565b610bfe565b6102ed7f000000000000000000000000000000000000000000000000000000000000000081565b61025861040d366004613736565b610c79565b6104397f000000000000000000000000000000000000000000000000000000000000000081565b60405163ffffffff9091168152602001610279565b60095473ffffffffffffffffffffffffffffffffffffffff166102ed565b610258610d11565b6102586104823660046136c7565b610ddf565b61026d610495366004613873565b610e60565b60015473ffffffffffffffffffffffffffffffffffffffff166102ed565b6104cb6104c6366004613890565b610e77565b60405161027991906138cb565b610439600081565b6104f36104ee366004613873565b6111b7565b6040516102799190613922565b610508611322565b60405161027991906139a4565b610258610523366004613736565b611333565b61053b610536366004613873565b61144b565b604051610279919081516fffffffffffffffffffffffffffffffff908116825260208084015163ffffffff1690830152604080840151151590830152606080840151821690830152608092830151169181019190915260a00190565b60045473ffffffffffffffffffffffffffffffffffffffff166102ed565b6102be6105c3366004613873565b611520565b6102586105d63660046136c7565b6115d0565b6105e36116a4565b60405161027991906139fe565b61053b6105fe366004613873565b61175c565b610258610611366004613b8b565b61182e565b7f00000000000000000000000000000000000000000000000000000000000000006102ed565b6106b061064a366004613873565b60408051606080820183526000808352602080840182905292840181905267ffffffffffffffff949094168452600a82529282902082519384018352805484526001015463ffffffff811691840191909152640100000000900460ff1615159082015290565b604080518251815260208084015163ffffffff169082015291810151151590820152606001610279565b7f000000000000000000000000000000000000000000000000000000000000000061026d565b61025861070e366004613807565b6118b2565b6102586107213660046136c7565b611dc4565b6102ed7f000000000000000000000000000000000000000000000000000000000000000081565b610755611dd8565b60005b818110156108ac57600083838381811061077457610774613bd2565b90506080020180360381019061078a9190613c15565b805190915015806107a75750604081015167ffffffffffffffff16155b1561081657604080517fa087bd2900000000000000000000000000000000000000000000000000000000815282516004820152602083015163ffffffff1660248201529082015167ffffffffffffffff1660448201526060820151151560648201526084015b60405180910390fd5b60408051606080820183528351825260208085015163ffffffff9081168285019081529286015115158486019081529585015167ffffffffffffffff166000908152600a90925293902091518255516001918201805494511515640100000000027fffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000000909516919093161792909217905501610758565b507f1889010d2535a0ab1643678d1da87fbbe8b87b2f585b47ddb72ec622aef9ee5682826040516108de929190613c8f565b60405180910390a15050565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167faff2afbf00000000000000000000000000000000000000000000000000000000148061097d57507fffffffff0000000000000000000000000000000000000000000000000000000082167f0e64dd2900000000000000000000000000000000000000000000000000000000145b806109c957507fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a700000000000000000000000000000000000000000000000000000000145b92915050565b6040805160208101909152600081526109e782611e2b565b60006109f660c0840184613d16565b810190610a039190613d7b565b90506000610a1460e0850185613d16565b810190610a219190613e48565b9050610a3181600001518361204f565b805160208201516040517f57ecfd2800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016926357ecfd2892610aa892600401613ed9565b6020604051808303816000875af1158015610ac7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610aeb9190613efe565b610b21576040517fbf969f2200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610b3160608501604086016136c7565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f9d228d69b5fdb8d273a2336f8fb8612d039631024ea9bf09c424a9503aa078f08660600135604051610b9391815260200190565b60405180910390a35050604080516020810190915260609092013582525090565b6000610bf68383604051610bc9929190613f1b565b604080519182900390912067ffffffffffffffff8716600090815260076020529190912060050190612200565b949350505050565b610c06611dd8565b610c738484808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152505060408051602080880282810182019093528782529093508792508691829185019084908082843760009201919091525061221b92505050565b50505050565b610c81611dd8565b610c8a83610e60565b610ccc576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8416600482015260240161080d565b610d0c8383838080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506123d192505050565b505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610d62576040517f02b543c600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000008082163390811790935560008054909116815560405173ffffffffffffffffffffffffffffffffffffffff909216929183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b610de7611dd8565b600980547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040519081527f44676b5284b809a22248eba0da87391d79098be38bb03154be88a58bf4d091749060200160405180910390a150565b60006109c9600567ffffffffffffffff8416612200565b6040805180820190915260608082526020820152610e94826124cb565b6000600a81610ea96040860160208701613873565b67ffffffffffffffff168152602080820192909252604090810160002081516060810183528154815260019091015463ffffffff81169382019390935264010000000090920460ff161515908201819052909150610f5057610f116040840160208501613873565b6040517fd201c48a00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff909116600482015260240161080d565b610f5a8380613d16565b9050602014610fa157610f6d8380613d16565b6040517fa3c8cf0900000000000000000000000000000000000000000000000000000000815260040161080d929190613f74565b6000610fad8480613d16565b810190610fba9190613f88565b602083015183516040517ff856ddb60000000000000000000000000000000000000000000000000000000081526060880135600482015263ffffffff90921660248301526044820183905273ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000008116606484015260848301919091529192506000917f0000000000000000000000000000000000000000000000000000000000000000169063f856ddb69060a4016020604051808303816000875af115801561109b573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906110bf9190613fa1565b6040516060870135815290915033907f696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df79060200160405180910390a2604051806040016040528061111c8760200160208101906105c39190613873565b815260408051808201825267ffffffffffffffff851680825263ffffffff7f00000000000000000000000000000000000000000000000000000000000000008116602093840190815284518085019390935251169281019290925290910190606001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152919052905295945050505050565b67ffffffffffffffff81166000908152600760205260408120606091906111e090600501612657565b90506000815167ffffffffffffffff8111156111fe576111fe613a40565b60405190808252806020026020018201604052801561123157816020015b606081526020019060019003908161121c5790505b50905060005b825181101561131a576008600084838151811061125657611256613bd2565b60200260200101518152602001908152602001600020805461127790613fbe565b80601f01602080910402602001604051908101604052809291908181526020018280546112a390613fbe565b80156112f05780601f106112c5576101008083540402835291602001916112f0565b820191906000526020600020905b8154815290600101906020018083116112d357829003601f168201915b505050505082828151811061130757611307613bd2565b6020908102919091010152600101611237565b509392505050565b606061132e6002612657565b905090565b61133b611dd8565b61134483610e60565b611386576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8416600482015260240161080d565b6113c68282604051611399929190613f1b565b604080519182900390912067ffffffffffffffff8616600090815260076020529190912060050190612664565b611402578282826040517f74f23c7c00000000000000000000000000000000000000000000000000000000815260040161080d93929190614011565b8267ffffffffffffffff167f52d00ee4d9bd51b40168f2afc5848837288ce258784ad914278791464b3f4d76838360405161143e929190613f74565b60405180910390a2505050565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915267ffffffffffffffff8216600090815260076020908152604091829020825160a08101845260028201546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff1615159482019490945260039091015480841660608301529190910490911660808201526109c990612670565b67ffffffffffffffff8116600090815260076020526040902060040180546060919061154b90613fbe565b80601f016020809104026020016040519081016040528092919081815260200182805461157790613fbe565b80156115c45780601f10611599576101008083540402835291602001916115c4565b820191906000526020600020905b8154815290600101906020018083116115a757829003601f168201915b50505050509050919050565b6115d8611dd8565b73ffffffffffffffffffffffffffffffffffffffff8116611625576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6004805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff000000000000000000000000000000000000000083168117909355604080519190921680825260208201939093527f02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f168491016108de565b606060006116b26005612657565b90506000815167ffffffffffffffff8111156116d0576116d0613a40565b6040519080825280602002602001820160405280156116f9578160200160208202803683370190505b50905060005b82518110156117555782818151811061171a5761171a613bd2565b602002602001015182828151811061173457611734613bd2565b67ffffffffffffffff909216602092830291909101909101526001016116ff565b5092915050565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915267ffffffffffffffff8216600090815260076020908152604091829020825160a08101845281546fffffffffffffffffffffffffffffffff808216835270010000000000000000000000000000000080830463ffffffff16958401959095527401000000000000000000000000000000000000000090910460ff1615159482019490945260019091015480841660608301529190910490911660808201526109c990612670565b60095473ffffffffffffffffffffffffffffffffffffffff16331480159061186e575060015473ffffffffffffffffffffffffffffffffffffffff163314155b156118a7576040517f8e4a23d600000000000000000000000000000000000000000000000000000000815233600482015260240161080d565b610d0c838383612722565b6118ba611dd8565b60005b83811015611aa75760008585838181106118d9576118d9613bd2565b90506020020160208101906118ee9190613873565b9050611905600567ffffffffffffffff8316612664565b611947576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8216600482015260240161080d565b67ffffffffffffffff8116600090815260076020526040812061196c90600501612657565b905060005b81518110156119d8576119cf82828151811061198f5761198f613bd2565b6020026020010151600760008667ffffffffffffffff1667ffffffffffffffff16815260200190815260200160002060050161266490919063ffffffff16565b50600101611971565b5067ffffffffffffffff8216600090815260076020526040812080547fffffffffffffffffffffff00000000000000000000000000000000000000000090811682556001820183905560028201805490911690556003810182905590611a41600483018261350a565b6005820160008181611a538282613544565b505060405167ffffffffffffffff871681527f5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d85991694506020019250611a95915050565b60405180910390a150506001016118bd565b5060005b81811015611dbd576000838383818110611ac757611ac7613bd2565b9050602002810190611ad99190614035565b611ae290614112565b9050611af38160600151600061280c565b611b028160800151600061280c565b806040015151600003611b41576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8051611b599060059067ffffffffffffffff16612949565b611b9e5780516040517f1d5ad3c500000000000000000000000000000000000000000000000000000000815267ffffffffffffffff909116600482015260240161080d565b805167ffffffffffffffff16600090815260076020908152604091829020825160a08082018552606080870180518601516fffffffffffffffffffffffffffffffff90811680865263ffffffff42168689018190528351511515878b0181905284518a0151841686890181905294518b0151841660809889018190528954740100000000000000000000000000000000000000009283027fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff7001000000000000000000000000000000008087027fffffffffffffffffffffffff000000000000000000000000000000000000000094851690981788178216929092178d5592810290971760018c01558c519889018d52898e0180518d01518716808b528a8e019590955280515115158a8f018190528151909d01518716988a01899052518d0151909516979098018790526002890180549a909102999093161717909416959095179092559092029091176003820155908201516004820190611d21908261421a565b5060005b826020015151811015611d6557611d5d836000015184602001518381518110611d5057611d50613bd2565b60200260200101516123d1565b600101611d25565b507f8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c28260000151836040015184606001518560800151604051611dab9493929190614334565b60405180910390a15050600101611aab565b5050505050565b611dcc611dd8565b611dd581612955565b50565b60015473ffffffffffffffffffffffffffffffffffffffff163314611e29576040517f2b5c74de00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b611e3e61032060a08301608084016136c7565b611e9d57611e5260a08201608083016136c7565b6040517f961c9a4f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff909116600482015260240161080d565b73ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016632cbc26bb611ee96040840160208501613873565b60405160e083901b7fffffffff0000000000000000000000000000000000000000000000000000000016815260809190911b77ffffffffffffffff00000000000000000000000000000000166004820152602401602060405180830381865afa158015611f5a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611f7e9190613efe565b15611fb5576040517f53ad11d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611fcd611fc86040830160208401613873565b612a19565b611fed611fe06040830160208401613873565b6103c060a0840184613d16565b61203257611ffe60a0820182613d16565b6040517f24eb47e500000000000000000000000000000000000000000000000000000000815260040161080d929190613f74565b611dd56120456040830160208401613873565b8260600135612b3f565b600482015163ffffffff81161561209a576040517f68d2f8d600000000000000000000000000000000000000000000000000000000815263ffffffff8216600482015260240161080d565b6008830151600c8401516014850151602085015163ffffffff8085169116146121055760208501516040517fe366a11700000000000000000000000000000000000000000000000000000000815263ffffffff9182166004820152908416602482015260440161080d565b7f000000000000000000000000000000000000000000000000000000000000000063ffffffff168263ffffffff161461219a576040517f77e4802600000000000000000000000000000000000000000000000000000000815263ffffffff7f0000000000000000000000000000000000000000000000000000000000000000811660048301528316602482015260440161080d565b845167ffffffffffffffff8281169116146121f85784516040517ff917ffea00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff9182166004820152908216602482015260440161080d565b505050505050565b600081815260018301602052604081205415155b9392505050565b7f0000000000000000000000000000000000000000000000000000000000000000612272576040517f35f4a7b300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005b825181101561230857600083828151811061229257612292613bd2565b602002602001015190506122b0816002612b8690919063ffffffff16565b156122ff5760405173ffffffffffffffffffffffffffffffffffffffff821681527f800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf75669060200160405180910390a15b50600101612275565b5060005b8151811015610d0c57600082828151811061232957612329613bd2565b60200260200101519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff160361236d57506123c9565b612378600282612ba8565b156123c75760405173ffffffffffffffffffffffffffffffffffffffff821681527f2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d89060200160405180910390a15b505b60010161230c565b805160000361240c576040517f8579befe00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805160208083019190912067ffffffffffffffff841660009081526007909252604090912061243e9060050182612949565b6124785782826040517f393b8ad200000000000000000000000000000000000000000000000000000000815260040161080d9291906143cd565b6000818152600860205260409020612490838261421a565b508267ffffffffffffffff167f7d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea8360405161143e9190613692565b6124de61032060a08301608084016136c7565b6124f257611e5260a08201608083016136c7565b73ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016632cbc26bb61253e6040840160208501613873565b60405160e083901b7fffffffff0000000000000000000000000000000000000000000000000000000016815260809190911b77ffffffffffffffff00000000000000000000000000000000166004820152602401602060405180830381865afa1580156125af573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906125d39190613efe565b1561260a576040517f53ad11d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61262261261d60608301604084016136c7565b612bca565b61263a6126356040830160208401613873565b612c49565b611dd561264d6040830160208401613873565b8260600135612d97565b6060600061221483612ddb565b60006122148383612e36565b6040805160a0810182526000808252602082018190529181018290526060810182905260808101919091526126fe82606001516fffffffffffffffffffffffffffffffff1683600001516fffffffffffffffffffffffffffffffff16846020015163ffffffff16426126e2919061441f565b85608001516fffffffffffffffffffffffffffffffff16612f29565b6fffffffffffffffffffffffffffffffff1682525063ffffffff4216602082015290565b61272b83610e60565b61276d576040517f1e670e4b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8416600482015260240161080d565b61277882600061280c565b67ffffffffffffffff8316600090815260076020526040902061279b9083612f51565b6127a681600061280c565b67ffffffffffffffff831660009081526007602052604090206127cc9060020182612f51565b7f0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b8383836040516127ff93929190614432565b60405180910390a1505050565b8151156128d75781602001516fffffffffffffffffffffffffffffffff1682604001516fffffffffffffffffffffffffffffffff16101580612862575060408201516fffffffffffffffffffffffffffffffff16155b1561289b57816040517f8020d12400000000000000000000000000000000000000000000000000000000815260040161080d91906144b5565b80156128d3576040517f433fc33d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050565b60408201516fffffffffffffffffffffffffffffffff16151580612910575060208201516fffffffffffffffffffffffffffffffff1615155b156128d357816040517fd68af9cc00000000000000000000000000000000000000000000000000000000815260040161080d91906144b5565b600061221483836130f3565b3373ffffffffffffffffffffffffffffffffffffffff8216036129a4576040517fdad89dca00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff838116918217835560015460405192939116917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b612a2281610e60565b612a64576040517fa9902c7e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8216600482015260240161080d565b600480546040517f83826b2b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84169281019290925233602483015273ffffffffffffffffffffffffffffffffffffffff16906383826b2b90604401602060405180830381865afa158015612ae3573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612b079190613efe565b611dd5576040517f728fe07b00000000000000000000000000000000000000000000000000000000815233600482015260240161080d565b67ffffffffffffffff821660009081526007602052604090206128d390600201827f0000000000000000000000000000000000000000000000000000000000000000613142565b60006122148373ffffffffffffffffffffffffffffffffffffffff8416612e36565b60006122148373ffffffffffffffffffffffffffffffffffffffff84166130f3565b7f000000000000000000000000000000000000000000000000000000000000000015611dd557612bfb6002826134c5565b611dd5576040517fd0d2597600000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8216600482015260240161080d565b612c5281610e60565b612c94576040517fa9902c7e00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff8216600482015260240161080d565b600480546040517fa8d87a3b00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff84169281019290925273ffffffffffffffffffffffffffffffffffffffff169063a8d87a3b90602401602060405180830381865afa158015612d0d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612d3191906144f1565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614611dd5576040517f728fe07b00000000000000000000000000000000000000000000000000000000815233600482015260240161080d565b67ffffffffffffffff821660009081526007602052604090206128d390827f0000000000000000000000000000000000000000000000000000000000000000613142565b6060816000018054806020026020016040519081016040528092919081815260200182805480156115c457602002820191906000526020600020905b815481526020019060010190808311612e175750505050509050919050565b60008181526001830160205260408120548015612f1f576000612e5a60018361441f565b8554909150600090612e6e9060019061441f565b9050808214612ed3576000866000018281548110612e8e57612e8e613bd2565b9060005260206000200154905080876000018481548110612eb157612eb1613bd2565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080612ee457612ee461450e565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506109c9565b60009150506109c9565b6000612f4885612f39848661453d565b612f439087614554565b6134f4565b95945050505050565b8154600090612f7a90700100000000000000000000000000000000900463ffffffff164261441f565b9050801561301c5760018301548354612fc2916fffffffffffffffffffffffffffffffff80821692811691859170010000000000000000000000000000000090910416612f29565b83546fffffffffffffffffffffffffffffffff919091167fffffffffffffffffffffffff0000000000000000000000000000000000000000909116177001000000000000000000000000000000004263ffffffff16021783555b60208201518354613042916fffffffffffffffffffffffffffffffff90811691166134f4565b83548351151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00ffffffff000000000000000000000000000000009091166fffffffffffffffffffffffffffffffff92831617178455602083015160408085015183167001000000000000000000000000000000000291909216176001850155517f9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c19906127ff9084906144b5565b600081815260018301602052604081205461313a575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556109c9565b5060006109c9565b825474010000000000000000000000000000000000000000900460ff161580613169575081155b1561317357505050565b825460018401546fffffffffffffffffffffffffffffffff808316929116906000906131b990700100000000000000000000000000000000900463ffffffff164261441f565b9050801561327957818311156131fb576040517f9725942a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60018601546132359083908590849070010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff16612f29565b86547fffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffff167001000000000000000000000000000000004263ffffffff160217875592505b848210156133305773ffffffffffffffffffffffffffffffffffffffff84166132d8576040517ff94ebcd1000000000000000000000000000000000000000000000000000000008152600481018390526024810186905260440161080d565b6040517f1a76572a000000000000000000000000000000000000000000000000000000008152600481018390526024810186905273ffffffffffffffffffffffffffffffffffffffff8516604482015260640161080d565b848310156134435760018681015470010000000000000000000000000000000090046fffffffffffffffffffffffffffffffff16906000908290613374908261441f565b61337e878a61441f565b6133889190614554565b6133929190614567565b905073ffffffffffffffffffffffffffffffffffffffff86166133eb576040517f15279c08000000000000000000000000000000000000000000000000000000008152600481018290526024810186905260440161080d565b6040517fd0c8d23a000000000000000000000000000000000000000000000000000000008152600481018290526024810186905273ffffffffffffffffffffffffffffffffffffffff8716604482015260640161080d565b61344d858461441f565b86547fffffffffffffffffffffffffffffffff00000000000000000000000000000000166fffffffffffffffffffffffffffffffff82161787556040518681529093507f1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a9060200160405180910390a1505050505050565b73ffffffffffffffffffffffffffffffffffffffff811660009081526001830160205260408120541515612214565b60008183106135035781612214565b5090919050565b50805461351690613fbe565b6000825580601f10613526575050565b601f016020900490600052602060002090810190611dd5919061355e565b5080546000825590600052602060002090810190611dd591905b5b80821115613573576000815560010161355f565b5090565b6000806020838503121561358a57600080fd5b823567ffffffffffffffff808211156135a257600080fd5b818501915085601f8301126135b657600080fd5b8135818111156135c557600080fd5b8660208260071b85010111156135da57600080fd5b60209290920196919550909350505050565b6000602082840312156135fe57600080fd5b81357fffffffff000000000000000000000000000000000000000000000000000000008116811461221457600080fd5b6000815180845260005b8181101561365457602081850181015186830182015201613638565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081526000612214602083018461362e565b73ffffffffffffffffffffffffffffffffffffffff81168114611dd557600080fd5b6000602082840312156136d957600080fd5b8135612214816136a5565b6000602082840312156136f657600080fd5b813567ffffffffffffffff81111561370d57600080fd5b8201610100818503121561221457600080fd5b67ffffffffffffffff81168114611dd557600080fd5b60008060006040848603121561374b57600080fd5b833561375681613720565b9250602084013567ffffffffffffffff8082111561377357600080fd5b818601915086601f83011261378757600080fd5b81358181111561379657600080fd5b8760208285010111156137a857600080fd5b6020830194508093505050509250925092565b60008083601f8401126137cd57600080fd5b50813567ffffffffffffffff8111156137e557600080fd5b6020830191508360208260051b850101111561380057600080fd5b9250929050565b6000806000806040858703121561381d57600080fd5b843567ffffffffffffffff8082111561383557600080fd5b613841888389016137bb565b9096509450602087013591508082111561385a57600080fd5b50613867878288016137bb565b95989497509550505050565b60006020828403121561388557600080fd5b813561221481613720565b6000602082840312156138a257600080fd5b813567ffffffffffffffff8111156138b957600080fd5b820160a0818503121561221457600080fd5b6020815260008251604060208401526138e7606084018261362e565b905060208401517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0848303016040850152612f48828261362e565b600060208083016020845280855180835260408601915060408160051b87010192506020870160005b82811015613997577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc088860301845261398585835161362e565b9450928501929085019060010161394b565b5092979650505050505050565b6020808252825182820181905260009190848201906040850190845b818110156139f257835173ffffffffffffffffffffffffffffffffffffffff16835292840192918401916001016139c0565b50909695505050505050565b6020808252825182820181905260009190848201906040850190845b818110156139f257835167ffffffffffffffff1683529284019291840191600101613a1a565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff81118282101715613a9257613a92613a40565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715613adf57613adf613a40565b604052919050565b8015158114611dd557600080fd5b80356fffffffffffffffffffffffffffffffff81168114613b1557600080fd5b919050565b600060608284031215613b2c57600080fd5b6040516060810181811067ffffffffffffffff82111715613b4f57613b4f613a40565b6040529050808235613b6081613ae7565b8152613b6e60208401613af5565b6020820152613b7f60408401613af5565b60408201525092915050565b600080600060e08486031215613ba057600080fd5b8335613bab81613720565b9250613bba8560208601613b1a565b9150613bc98560808601613b1a565b90509250925092565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b803563ffffffff81168114613b1557600080fd5b600060808284031215613c2757600080fd5b6040516080810181811067ffffffffffffffff82111715613c4a57613c4a613a40565b60405282358152613c5d60208401613c01565b60208201526040830135613c7081613720565b60408201526060830135613c8381613ae7565b60608201529392505050565b6020808252818101839052600090604080840186845b87811015613d09578135835263ffffffff613cc1868401613c01565b168584015283820135613cd381613720565b67ffffffffffffffff1683850152606082810135613cf081613ae7565b1515908401526080928301929190910190600101613ca5565b5090979650505050505050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112613d4b57600080fd5b83018035915067ffffffffffffffff821115613d6657600080fd5b60200191503681900382131561380057600080fd5b600060408284031215613d8d57600080fd5b613d95613a6f565b8235613da081613720565b8152613dae60208401613c01565b60208201529392505050565b600082601f830112613dcb57600080fd5b813567ffffffffffffffff811115613de557613de5613a40565b613e1660207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601613a98565b818152846020838601011115613e2b57600080fd5b816020850160208301376000918101602001919091529392505050565b600060208284031215613e5a57600080fd5b813567ffffffffffffffff80821115613e7257600080fd5b9083019060408286031215613e8657600080fd5b613e8e613a6f565b823582811115613e9d57600080fd5b613ea987828601613dba565b825250602083013582811115613ebe57600080fd5b613eca87828601613dba565b60208301525095945050505050565b604081526000613eec604083018561362e565b8281036020840152612f48818561362e565b600060208284031215613f1057600080fd5b815161221481613ae7565b8183823760009101908152919050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b602081526000610bf6602083018486613f2b565b600060208284031215613f9a57600080fd5b5035919050565b600060208284031215613fb357600080fd5b815161221481613720565b600181811c90821680613fd257607f821691505b60208210810361400b577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b67ffffffffffffffff84168152604060208201526000612f48604083018486613f2b565b600082357ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffee183360301811261406957600080fd5b9190910192915050565b600082601f83011261408457600080fd5b8135602067ffffffffffffffff808311156140a1576140a1613a40565b8260051b6140b0838201613a98565b93845285810183019383810190888611156140ca57600080fd5b84880192505b85831015614106578235848111156140e85760008081fd5b6140f68a87838c0101613dba565b83525091840191908401906140d0565b98975050505050505050565b6000610120823603121561412557600080fd5b60405160a0810167ffffffffffffffff828210818311171561414957614149613a40565b816040528435915061415a82613720565b9082526020840135908082111561417057600080fd5b61417c36838701614073565b6020840152604085013591508082111561419557600080fd5b506141a236828601613dba565b6040830152506141b53660608501613b1a565b60608201526141c73660c08501613b1a565b608082015292915050565b601f821115610d0c576000816000526020600020601f850160051c810160208610156141fb5750805b601f850160051c820191505b818110156121f857828155600101614207565b815167ffffffffffffffff81111561423457614234613a40565b614248816142428454613fbe565b846141d2565b602080601f83116001811461429b57600084156142655750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b1785556121f8565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b828110156142e8578886015182559484019460019091019084016142c9565b508582101561432457878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b600061010067ffffffffffffffff871683528060208401526143588184018761362e565b8551151560408581019190915260208701516fffffffffffffffffffffffffffffffff90811660608701529087015116608085015291506143969050565b8251151560a083015260208301516fffffffffffffffffffffffffffffffff90811660c084015260408401511660e0830152612f48565b67ffffffffffffffff83168152604060208201526000610bf6604083018461362e565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b818103818111156109c9576109c96143f0565b67ffffffffffffffff8416815260e0810161447e60208301858051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b82511515608083015260208301516fffffffffffffffffffffffffffffffff90811660a084015260408401511660c0830152610bf6565b606081016109c982848051151582526020808201516fffffffffffffffffffffffffffffffff9081169184019190915260409182015116910152565b60006020828403121561450357600080fd5b8151612214816136a5565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b80820281158282048414176109c9576109c96143f0565b808201808211156109c9576109c96143f0565b60008261459d577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b50049056fea164736f6c6343000818000a",
}

var USDCTokenPoolABI = USDCTokenPoolMetaData.ABI

var USDCTokenPoolBin = USDCTokenPoolMetaData.Bin

func DeployUSDCTokenPool(auth *bind.TransactOpts, backend bind.ContractBackend, tokenMessenger common.Address, token common.Address, allowlist []common.Address, rmnProxy common.Address, router common.Address) (common.Address, *generated.Transaction, *USDCTokenPool, error) {
	parsed, err := USDCTokenPoolMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}
	if generated.IsZKSync(backend) {
		address, ethTx, contractBind, _ := generated.DeployContract(auth, parsed, common.FromHex(USDCTokenPoolZKBin), backend, tokenMessenger, token, allowlist, rmnProxy, router)
		contractReturn := &USDCTokenPool{address: address, abi: *parsed, USDCTokenPoolCaller: USDCTokenPoolCaller{contract: contractBind}, USDCTokenPoolTransactor: USDCTokenPoolTransactor{contract: contractBind}, USDCTokenPoolFilterer: USDCTokenPoolFilterer{contract: contractBind}}
		return address, ethTx, contractReturn, err
	}
	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(USDCTokenPoolBin), backend, tokenMessenger, token, allowlist, rmnProxy, router)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, &generated.Transaction{Transaction: tx, HashZks: tx.Hash()}, &USDCTokenPool{address: address, abi: *parsed, USDCTokenPoolCaller: USDCTokenPoolCaller{contract: contract}, USDCTokenPoolTransactor: USDCTokenPoolTransactor{contract: contract}, USDCTokenPoolFilterer: USDCTokenPoolFilterer{contract: contract}}, nil
}

type USDCTokenPool struct {
	address common.Address
	abi     abi.ABI
	USDCTokenPoolCaller
	USDCTokenPoolTransactor
	USDCTokenPoolFilterer
}

type USDCTokenPoolCaller struct {
	contract *bind.BoundContract
}

type USDCTokenPoolTransactor struct {
	contract *bind.BoundContract
}

type USDCTokenPoolFilterer struct {
	contract *bind.BoundContract
}

type USDCTokenPoolSession struct {
	Contract     *USDCTokenPool
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type USDCTokenPoolCallerSession struct {
	Contract *USDCTokenPoolCaller
	CallOpts bind.CallOpts
}

type USDCTokenPoolTransactorSession struct {
	Contract     *USDCTokenPoolTransactor
	TransactOpts bind.TransactOpts
}

type USDCTokenPoolRaw struct {
	Contract *USDCTokenPool
}

type USDCTokenPoolCallerRaw struct {
	Contract *USDCTokenPoolCaller
}

type USDCTokenPoolTransactorRaw struct {
	Contract *USDCTokenPoolTransactor
}

func NewUSDCTokenPool(address common.Address, backend bind.ContractBackend) (*USDCTokenPool, error) {
	abi, err := abi.JSON(strings.NewReader(USDCTokenPoolABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindUSDCTokenPool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &USDCTokenPool{address: address, abi: abi, USDCTokenPoolCaller: USDCTokenPoolCaller{contract: contract}, USDCTokenPoolTransactor: USDCTokenPoolTransactor{contract: contract}, USDCTokenPoolFilterer: USDCTokenPoolFilterer{contract: contract}}, nil
}

func NewUSDCTokenPoolCaller(address common.Address, caller bind.ContractCaller) (*USDCTokenPoolCaller, error) {
	contract, err := bindUSDCTokenPool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolCaller{contract: contract}, nil
}

func NewUSDCTokenPoolTransactor(address common.Address, transactor bind.ContractTransactor) (*USDCTokenPoolTransactor, error) {
	contract, err := bindUSDCTokenPool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolTransactor{contract: contract}, nil
}

func NewUSDCTokenPoolFilterer(address common.Address, filterer bind.ContractFilterer) (*USDCTokenPoolFilterer, error) {
	contract, err := bindUSDCTokenPool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolFilterer{contract: contract}, nil
}

func bindUSDCTokenPool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := USDCTokenPoolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_USDCTokenPool *USDCTokenPoolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _USDCTokenPool.Contract.USDCTokenPoolCaller.contract.Call(opts, result, method, params...)
}

func (_USDCTokenPool *USDCTokenPoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.USDCTokenPoolTransactor.contract.Transfer(opts)
}

func (_USDCTokenPool *USDCTokenPoolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.USDCTokenPoolTransactor.contract.Transact(opts, method, params...)
}

func (_USDCTokenPool *USDCTokenPoolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _USDCTokenPool.Contract.contract.Call(opts, result, method, params...)
}

func (_USDCTokenPool *USDCTokenPoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.contract.Transfer(opts)
}

func (_USDCTokenPool *USDCTokenPoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.contract.Transact(opts, method, params...)
}

func (_USDCTokenPool *USDCTokenPoolCaller) SUPPORTEDUSDCVERSION(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "SUPPORTED_USDC_VERSION")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) SUPPORTEDUSDCVERSION() (uint32, error) {
	return _USDCTokenPool.Contract.SUPPORTEDUSDCVERSION(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) SUPPORTEDUSDCVERSION() (uint32, error) {
	return _USDCTokenPool.Contract.SUPPORTEDUSDCVERSION(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCaller) GetAllowList(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "getAllowList")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) GetAllowList() ([]common.Address, error) {
	return _USDCTokenPool.Contract.GetAllowList(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) GetAllowList() ([]common.Address, error) {
	return _USDCTokenPool.Contract.GetAllowList(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCaller) GetAllowListEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "getAllowListEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) GetAllowListEnabled() (bool, error) {
	return _USDCTokenPool.Contract.GetAllowListEnabled(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) GetAllowListEnabled() (bool, error) {
	return _USDCTokenPool.Contract.GetAllowListEnabled(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCaller) GetCurrentInboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "getCurrentInboundRateLimiterState", remoteChainSelector)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) GetCurrentInboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _USDCTokenPool.Contract.GetCurrentInboundRateLimiterState(&_USDCTokenPool.CallOpts, remoteChainSelector)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) GetCurrentInboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _USDCTokenPool.Contract.GetCurrentInboundRateLimiterState(&_USDCTokenPool.CallOpts, remoteChainSelector)
}

func (_USDCTokenPool *USDCTokenPoolCaller) GetCurrentOutboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "getCurrentOutboundRateLimiterState", remoteChainSelector)

	if err != nil {
		return *new(RateLimiterTokenBucket), err
	}

	out0 := *abi.ConvertType(out[0], new(RateLimiterTokenBucket)).(*RateLimiterTokenBucket)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) GetCurrentOutboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _USDCTokenPool.Contract.GetCurrentOutboundRateLimiterState(&_USDCTokenPool.CallOpts, remoteChainSelector)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) GetCurrentOutboundRateLimiterState(remoteChainSelector uint64) (RateLimiterTokenBucket, error) {
	return _USDCTokenPool.Contract.GetCurrentOutboundRateLimiterState(&_USDCTokenPool.CallOpts, remoteChainSelector)
}

func (_USDCTokenPool *USDCTokenPoolCaller) GetDomain(opts *bind.CallOpts, chainSelector uint64) (USDCTokenPoolDomain, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "getDomain", chainSelector)

	if err != nil {
		return *new(USDCTokenPoolDomain), err
	}

	out0 := *abi.ConvertType(out[0], new(USDCTokenPoolDomain)).(*USDCTokenPoolDomain)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) GetDomain(chainSelector uint64) (USDCTokenPoolDomain, error) {
	return _USDCTokenPool.Contract.GetDomain(&_USDCTokenPool.CallOpts, chainSelector)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) GetDomain(chainSelector uint64) (USDCTokenPoolDomain, error) {
	return _USDCTokenPool.Contract.GetDomain(&_USDCTokenPool.CallOpts, chainSelector)
}

func (_USDCTokenPool *USDCTokenPoolCaller) GetRateLimitAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "getRateLimitAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) GetRateLimitAdmin() (common.Address, error) {
	return _USDCTokenPool.Contract.GetRateLimitAdmin(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) GetRateLimitAdmin() (common.Address, error) {
	return _USDCTokenPool.Contract.GetRateLimitAdmin(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCaller) GetRemotePools(opts *bind.CallOpts, remoteChainSelector uint64) ([][]byte, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "getRemotePools", remoteChainSelector)

	if err != nil {
		return *new([][]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][]byte)).(*[][]byte)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) GetRemotePools(remoteChainSelector uint64) ([][]byte, error) {
	return _USDCTokenPool.Contract.GetRemotePools(&_USDCTokenPool.CallOpts, remoteChainSelector)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) GetRemotePools(remoteChainSelector uint64) ([][]byte, error) {
	return _USDCTokenPool.Contract.GetRemotePools(&_USDCTokenPool.CallOpts, remoteChainSelector)
}

func (_USDCTokenPool *USDCTokenPoolCaller) GetRemoteToken(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "getRemoteToken", remoteChainSelector)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) GetRemoteToken(remoteChainSelector uint64) ([]byte, error) {
	return _USDCTokenPool.Contract.GetRemoteToken(&_USDCTokenPool.CallOpts, remoteChainSelector)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) GetRemoteToken(remoteChainSelector uint64) ([]byte, error) {
	return _USDCTokenPool.Contract.GetRemoteToken(&_USDCTokenPool.CallOpts, remoteChainSelector)
}

func (_USDCTokenPool *USDCTokenPoolCaller) GetRmnProxy(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "getRmnProxy")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) GetRmnProxy() (common.Address, error) {
	return _USDCTokenPool.Contract.GetRmnProxy(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) GetRmnProxy() (common.Address, error) {
	return _USDCTokenPool.Contract.GetRmnProxy(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCaller) GetRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "getRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) GetRouter() (common.Address, error) {
	return _USDCTokenPool.Contract.GetRouter(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) GetRouter() (common.Address, error) {
	return _USDCTokenPool.Contract.GetRouter(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCaller) GetSupportedChains(opts *bind.CallOpts) ([]uint64, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "getSupportedChains")

	if err != nil {
		return *new([]uint64), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint64)).(*[]uint64)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) GetSupportedChains() ([]uint64, error) {
	return _USDCTokenPool.Contract.GetSupportedChains(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) GetSupportedChains() ([]uint64, error) {
	return _USDCTokenPool.Contract.GetSupportedChains(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCaller) GetToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "getToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) GetToken() (common.Address, error) {
	return _USDCTokenPool.Contract.GetToken(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) GetToken() (common.Address, error) {
	return _USDCTokenPool.Contract.GetToken(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCaller) GetTokenDecimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "getTokenDecimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) GetTokenDecimals() (uint8, error) {
	return _USDCTokenPool.Contract.GetTokenDecimals(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) GetTokenDecimals() (uint8, error) {
	return _USDCTokenPool.Contract.GetTokenDecimals(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCaller) ILocalDomainIdentifier(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "i_localDomainIdentifier")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) ILocalDomainIdentifier() (uint32, error) {
	return _USDCTokenPool.Contract.ILocalDomainIdentifier(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) ILocalDomainIdentifier() (uint32, error) {
	return _USDCTokenPool.Contract.ILocalDomainIdentifier(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCaller) IMessageTransmitter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "i_messageTransmitter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) IMessageTransmitter() (common.Address, error) {
	return _USDCTokenPool.Contract.IMessageTransmitter(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) IMessageTransmitter() (common.Address, error) {
	return _USDCTokenPool.Contract.IMessageTransmitter(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCaller) ITokenMessenger(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "i_tokenMessenger")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) ITokenMessenger() (common.Address, error) {
	return _USDCTokenPool.Contract.ITokenMessenger(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) ITokenMessenger() (common.Address, error) {
	return _USDCTokenPool.Contract.ITokenMessenger(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCaller) IsRemotePool(opts *bind.CallOpts, remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "isRemotePool", remoteChainSelector, remotePoolAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) IsRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	return _USDCTokenPool.Contract.IsRemotePool(&_USDCTokenPool.CallOpts, remoteChainSelector, remotePoolAddress)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) IsRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	return _USDCTokenPool.Contract.IsRemotePool(&_USDCTokenPool.CallOpts, remoteChainSelector, remotePoolAddress)
}

func (_USDCTokenPool *USDCTokenPoolCaller) IsSupportedChain(opts *bind.CallOpts, remoteChainSelector uint64) (bool, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "isSupportedChain", remoteChainSelector)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) IsSupportedChain(remoteChainSelector uint64) (bool, error) {
	return _USDCTokenPool.Contract.IsSupportedChain(&_USDCTokenPool.CallOpts, remoteChainSelector)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) IsSupportedChain(remoteChainSelector uint64) (bool, error) {
	return _USDCTokenPool.Contract.IsSupportedChain(&_USDCTokenPool.CallOpts, remoteChainSelector)
}

func (_USDCTokenPool *USDCTokenPoolCaller) IsSupportedToken(opts *bind.CallOpts, token common.Address) (bool, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "isSupportedToken", token)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) IsSupportedToken(token common.Address) (bool, error) {
	return _USDCTokenPool.Contract.IsSupportedToken(&_USDCTokenPool.CallOpts, token)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) IsSupportedToken(token common.Address) (bool, error) {
	return _USDCTokenPool.Contract.IsSupportedToken(&_USDCTokenPool.CallOpts, token)
}

func (_USDCTokenPool *USDCTokenPoolCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) Owner() (common.Address, error) {
	return _USDCTokenPool.Contract.Owner(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) Owner() (common.Address, error) {
	return _USDCTokenPool.Contract.Owner(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _USDCTokenPool.Contract.SupportsInterface(&_USDCTokenPool.CallOpts, interfaceId)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _USDCTokenPool.Contract.SupportsInterface(&_USDCTokenPool.CallOpts, interfaceId)
}

func (_USDCTokenPool *USDCTokenPoolCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _USDCTokenPool.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_USDCTokenPool *USDCTokenPoolSession) TypeAndVersion() (string, error) {
	return _USDCTokenPool.Contract.TypeAndVersion(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolCallerSession) TypeAndVersion() (string, error) {
	return _USDCTokenPool.Contract.TypeAndVersion(&_USDCTokenPool.CallOpts)
}

func (_USDCTokenPool *USDCTokenPoolTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _USDCTokenPool.contract.Transact(opts, "acceptOwnership")
}

func (_USDCTokenPool *USDCTokenPoolSession) AcceptOwnership() (*types.Transaction, error) {
	return _USDCTokenPool.Contract.AcceptOwnership(&_USDCTokenPool.TransactOpts)
}

func (_USDCTokenPool *USDCTokenPoolTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _USDCTokenPool.Contract.AcceptOwnership(&_USDCTokenPool.TransactOpts)
}

func (_USDCTokenPool *USDCTokenPoolTransactor) AddRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _USDCTokenPool.contract.Transact(opts, "addRemotePool", remoteChainSelector, remotePoolAddress)
}

func (_USDCTokenPool *USDCTokenPoolSession) AddRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.AddRemotePool(&_USDCTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_USDCTokenPool *USDCTokenPoolTransactorSession) AddRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.AddRemotePool(&_USDCTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_USDCTokenPool *USDCTokenPoolTransactor) ApplyAllowListUpdates(opts *bind.TransactOpts, removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _USDCTokenPool.contract.Transact(opts, "applyAllowListUpdates", removes, adds)
}

func (_USDCTokenPool *USDCTokenPoolSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.ApplyAllowListUpdates(&_USDCTokenPool.TransactOpts, removes, adds)
}

func (_USDCTokenPool *USDCTokenPoolTransactorSession) ApplyAllowListUpdates(removes []common.Address, adds []common.Address) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.ApplyAllowListUpdates(&_USDCTokenPool.TransactOpts, removes, adds)
}

func (_USDCTokenPool *USDCTokenPoolTransactor) ApplyChainUpdates(opts *bind.TransactOpts, remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _USDCTokenPool.contract.Transact(opts, "applyChainUpdates", remoteChainSelectorsToRemove, chainsToAdd)
}

func (_USDCTokenPool *USDCTokenPoolSession) ApplyChainUpdates(remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.ApplyChainUpdates(&_USDCTokenPool.TransactOpts, remoteChainSelectorsToRemove, chainsToAdd)
}

func (_USDCTokenPool *USDCTokenPoolTransactorSession) ApplyChainUpdates(remoteChainSelectorsToRemove []uint64, chainsToAdd []TokenPoolChainUpdate) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.ApplyChainUpdates(&_USDCTokenPool.TransactOpts, remoteChainSelectorsToRemove, chainsToAdd)
}

func (_USDCTokenPool *USDCTokenPoolTransactor) LockOrBurn(opts *bind.TransactOpts, lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _USDCTokenPool.contract.Transact(opts, "lockOrBurn", lockOrBurnIn)
}

func (_USDCTokenPool *USDCTokenPoolSession) LockOrBurn(lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.LockOrBurn(&_USDCTokenPool.TransactOpts, lockOrBurnIn)
}

func (_USDCTokenPool *USDCTokenPoolTransactorSession) LockOrBurn(lockOrBurnIn PoolLockOrBurnInV1) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.LockOrBurn(&_USDCTokenPool.TransactOpts, lockOrBurnIn)
}

func (_USDCTokenPool *USDCTokenPoolTransactor) ReleaseOrMint(opts *bind.TransactOpts, releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _USDCTokenPool.contract.Transact(opts, "releaseOrMint", releaseOrMintIn)
}

func (_USDCTokenPool *USDCTokenPoolSession) ReleaseOrMint(releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.ReleaseOrMint(&_USDCTokenPool.TransactOpts, releaseOrMintIn)
}

func (_USDCTokenPool *USDCTokenPoolTransactorSession) ReleaseOrMint(releaseOrMintIn PoolReleaseOrMintInV1) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.ReleaseOrMint(&_USDCTokenPool.TransactOpts, releaseOrMintIn)
}

func (_USDCTokenPool *USDCTokenPoolTransactor) RemoveRemotePool(opts *bind.TransactOpts, remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _USDCTokenPool.contract.Transact(opts, "removeRemotePool", remoteChainSelector, remotePoolAddress)
}

func (_USDCTokenPool *USDCTokenPoolSession) RemoveRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.RemoveRemotePool(&_USDCTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_USDCTokenPool *USDCTokenPoolTransactorSession) RemoveRemotePool(remoteChainSelector uint64, remotePoolAddress []byte) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.RemoveRemotePool(&_USDCTokenPool.TransactOpts, remoteChainSelector, remotePoolAddress)
}

func (_USDCTokenPool *USDCTokenPoolTransactor) SetChainRateLimiterConfig(opts *bind.TransactOpts, remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _USDCTokenPool.contract.Transact(opts, "setChainRateLimiterConfig", remoteChainSelector, outboundConfig, inboundConfig)
}

func (_USDCTokenPool *USDCTokenPoolSession) SetChainRateLimiterConfig(remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.SetChainRateLimiterConfig(&_USDCTokenPool.TransactOpts, remoteChainSelector, outboundConfig, inboundConfig)
}

func (_USDCTokenPool *USDCTokenPoolTransactorSession) SetChainRateLimiterConfig(remoteChainSelector uint64, outboundConfig RateLimiterConfig, inboundConfig RateLimiterConfig) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.SetChainRateLimiterConfig(&_USDCTokenPool.TransactOpts, remoteChainSelector, outboundConfig, inboundConfig)
}

func (_USDCTokenPool *USDCTokenPoolTransactor) SetDomains(opts *bind.TransactOpts, domains []USDCTokenPoolDomainUpdate) (*types.Transaction, error) {
	return _USDCTokenPool.contract.Transact(opts, "setDomains", domains)
}

func (_USDCTokenPool *USDCTokenPoolSession) SetDomains(domains []USDCTokenPoolDomainUpdate) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.SetDomains(&_USDCTokenPool.TransactOpts, domains)
}

func (_USDCTokenPool *USDCTokenPoolTransactorSession) SetDomains(domains []USDCTokenPoolDomainUpdate) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.SetDomains(&_USDCTokenPool.TransactOpts, domains)
}

func (_USDCTokenPool *USDCTokenPoolTransactor) SetRateLimitAdmin(opts *bind.TransactOpts, rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _USDCTokenPool.contract.Transact(opts, "setRateLimitAdmin", rateLimitAdmin)
}

func (_USDCTokenPool *USDCTokenPoolSession) SetRateLimitAdmin(rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.SetRateLimitAdmin(&_USDCTokenPool.TransactOpts, rateLimitAdmin)
}

func (_USDCTokenPool *USDCTokenPoolTransactorSession) SetRateLimitAdmin(rateLimitAdmin common.Address) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.SetRateLimitAdmin(&_USDCTokenPool.TransactOpts, rateLimitAdmin)
}

func (_USDCTokenPool *USDCTokenPoolTransactor) SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error) {
	return _USDCTokenPool.contract.Transact(opts, "setRouter", newRouter)
}

func (_USDCTokenPool *USDCTokenPoolSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.SetRouter(&_USDCTokenPool.TransactOpts, newRouter)
}

func (_USDCTokenPool *USDCTokenPoolTransactorSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.SetRouter(&_USDCTokenPool.TransactOpts, newRouter)
}

func (_USDCTokenPool *USDCTokenPoolTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _USDCTokenPool.contract.Transact(opts, "transferOwnership", to)
}

func (_USDCTokenPool *USDCTokenPoolSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.TransferOwnership(&_USDCTokenPool.TransactOpts, to)
}

func (_USDCTokenPool *USDCTokenPoolTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _USDCTokenPool.Contract.TransferOwnership(&_USDCTokenPool.TransactOpts, to)
}

type USDCTokenPoolAllowListAddIterator struct {
	Event *USDCTokenPoolAllowListAdd

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCTokenPoolAllowListAddIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCTokenPoolAllowListAdd)
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
		it.Event = new(USDCTokenPoolAllowListAdd)
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

func (it *USDCTokenPoolAllowListAddIterator) Error() error {
	return it.fail
}

func (it *USDCTokenPoolAllowListAddIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCTokenPoolAllowListAdd struct {
	Sender common.Address
	Raw    types.Log
}

func (_USDCTokenPool *USDCTokenPoolFilterer) FilterAllowListAdd(opts *bind.FilterOpts) (*USDCTokenPoolAllowListAddIterator, error) {

	logs, sub, err := _USDCTokenPool.contract.FilterLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolAllowListAddIterator{contract: _USDCTokenPool.contract, event: "AllowListAdd", logs: logs, sub: sub}, nil
}

func (_USDCTokenPool *USDCTokenPoolFilterer) WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolAllowListAdd) (event.Subscription, error) {

	logs, sub, err := _USDCTokenPool.contract.WatchLogs(opts, "AllowListAdd")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCTokenPoolAllowListAdd)
				if err := _USDCTokenPool.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
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

func (_USDCTokenPool *USDCTokenPoolFilterer) ParseAllowListAdd(log types.Log) (*USDCTokenPoolAllowListAdd, error) {
	event := new(USDCTokenPoolAllowListAdd)
	if err := _USDCTokenPool.contract.UnpackLog(event, "AllowListAdd", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type USDCTokenPoolAllowListRemoveIterator struct {
	Event *USDCTokenPoolAllowListRemove

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCTokenPoolAllowListRemoveIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCTokenPoolAllowListRemove)
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
		it.Event = new(USDCTokenPoolAllowListRemove)
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

func (it *USDCTokenPoolAllowListRemoveIterator) Error() error {
	return it.fail
}

func (it *USDCTokenPoolAllowListRemoveIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCTokenPoolAllowListRemove struct {
	Sender common.Address
	Raw    types.Log
}

func (_USDCTokenPool *USDCTokenPoolFilterer) FilterAllowListRemove(opts *bind.FilterOpts) (*USDCTokenPoolAllowListRemoveIterator, error) {

	logs, sub, err := _USDCTokenPool.contract.FilterLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolAllowListRemoveIterator{contract: _USDCTokenPool.contract, event: "AllowListRemove", logs: logs, sub: sub}, nil
}

func (_USDCTokenPool *USDCTokenPoolFilterer) WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolAllowListRemove) (event.Subscription, error) {

	logs, sub, err := _USDCTokenPool.contract.WatchLogs(opts, "AllowListRemove")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCTokenPoolAllowListRemove)
				if err := _USDCTokenPool.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
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

func (_USDCTokenPool *USDCTokenPoolFilterer) ParseAllowListRemove(log types.Log) (*USDCTokenPoolAllowListRemove, error) {
	event := new(USDCTokenPoolAllowListRemove)
	if err := _USDCTokenPool.contract.UnpackLog(event, "AllowListRemove", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type USDCTokenPoolBurnedIterator struct {
	Event *USDCTokenPoolBurned

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCTokenPoolBurnedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCTokenPoolBurned)
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
		it.Event = new(USDCTokenPoolBurned)
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

func (it *USDCTokenPoolBurnedIterator) Error() error {
	return it.fail
}

func (it *USDCTokenPoolBurnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCTokenPoolBurned struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_USDCTokenPool *USDCTokenPoolFilterer) FilterBurned(opts *bind.FilterOpts, sender []common.Address) (*USDCTokenPoolBurnedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _USDCTokenPool.contract.FilterLogs(opts, "Burned", senderRule)
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolBurnedIterator{contract: _USDCTokenPool.contract, event: "Burned", logs: logs, sub: sub}, nil
}

func (_USDCTokenPool *USDCTokenPoolFilterer) WatchBurned(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolBurned, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _USDCTokenPool.contract.WatchLogs(opts, "Burned", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCTokenPoolBurned)
				if err := _USDCTokenPool.contract.UnpackLog(event, "Burned", log); err != nil {
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

func (_USDCTokenPool *USDCTokenPoolFilterer) ParseBurned(log types.Log) (*USDCTokenPoolBurned, error) {
	event := new(USDCTokenPoolBurned)
	if err := _USDCTokenPool.contract.UnpackLog(event, "Burned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type USDCTokenPoolChainAddedIterator struct {
	Event *USDCTokenPoolChainAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCTokenPoolChainAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCTokenPoolChainAdded)
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
		it.Event = new(USDCTokenPoolChainAdded)
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

func (it *USDCTokenPoolChainAddedIterator) Error() error {
	return it.fail
}

func (it *USDCTokenPoolChainAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCTokenPoolChainAdded struct {
	RemoteChainSelector       uint64
	RemoteToken               []byte
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
	Raw                       types.Log
}

func (_USDCTokenPool *USDCTokenPoolFilterer) FilterChainAdded(opts *bind.FilterOpts) (*USDCTokenPoolChainAddedIterator, error) {

	logs, sub, err := _USDCTokenPool.contract.FilterLogs(opts, "ChainAdded")
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolChainAddedIterator{contract: _USDCTokenPool.contract, event: "ChainAdded", logs: logs, sub: sub}, nil
}

func (_USDCTokenPool *USDCTokenPoolFilterer) WatchChainAdded(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolChainAdded) (event.Subscription, error) {

	logs, sub, err := _USDCTokenPool.contract.WatchLogs(opts, "ChainAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCTokenPoolChainAdded)
				if err := _USDCTokenPool.contract.UnpackLog(event, "ChainAdded", log); err != nil {
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

func (_USDCTokenPool *USDCTokenPoolFilterer) ParseChainAdded(log types.Log) (*USDCTokenPoolChainAdded, error) {
	event := new(USDCTokenPoolChainAdded)
	if err := _USDCTokenPool.contract.UnpackLog(event, "ChainAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type USDCTokenPoolChainConfiguredIterator struct {
	Event *USDCTokenPoolChainConfigured

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCTokenPoolChainConfiguredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCTokenPoolChainConfigured)
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
		it.Event = new(USDCTokenPoolChainConfigured)
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

func (it *USDCTokenPoolChainConfiguredIterator) Error() error {
	return it.fail
}

func (it *USDCTokenPoolChainConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCTokenPoolChainConfigured struct {
	RemoteChainSelector       uint64
	OutboundRateLimiterConfig RateLimiterConfig
	InboundRateLimiterConfig  RateLimiterConfig
	Raw                       types.Log
}

func (_USDCTokenPool *USDCTokenPoolFilterer) FilterChainConfigured(opts *bind.FilterOpts) (*USDCTokenPoolChainConfiguredIterator, error) {

	logs, sub, err := _USDCTokenPool.contract.FilterLogs(opts, "ChainConfigured")
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolChainConfiguredIterator{contract: _USDCTokenPool.contract, event: "ChainConfigured", logs: logs, sub: sub}, nil
}

func (_USDCTokenPool *USDCTokenPoolFilterer) WatchChainConfigured(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolChainConfigured) (event.Subscription, error) {

	logs, sub, err := _USDCTokenPool.contract.WatchLogs(opts, "ChainConfigured")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCTokenPoolChainConfigured)
				if err := _USDCTokenPool.contract.UnpackLog(event, "ChainConfigured", log); err != nil {
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

func (_USDCTokenPool *USDCTokenPoolFilterer) ParseChainConfigured(log types.Log) (*USDCTokenPoolChainConfigured, error) {
	event := new(USDCTokenPoolChainConfigured)
	if err := _USDCTokenPool.contract.UnpackLog(event, "ChainConfigured", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type USDCTokenPoolChainRemovedIterator struct {
	Event *USDCTokenPoolChainRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCTokenPoolChainRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCTokenPoolChainRemoved)
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
		it.Event = new(USDCTokenPoolChainRemoved)
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

func (it *USDCTokenPoolChainRemovedIterator) Error() error {
	return it.fail
}

func (it *USDCTokenPoolChainRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCTokenPoolChainRemoved struct {
	RemoteChainSelector uint64
	Raw                 types.Log
}

func (_USDCTokenPool *USDCTokenPoolFilterer) FilterChainRemoved(opts *bind.FilterOpts) (*USDCTokenPoolChainRemovedIterator, error) {

	logs, sub, err := _USDCTokenPool.contract.FilterLogs(opts, "ChainRemoved")
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolChainRemovedIterator{contract: _USDCTokenPool.contract, event: "ChainRemoved", logs: logs, sub: sub}, nil
}

func (_USDCTokenPool *USDCTokenPoolFilterer) WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolChainRemoved) (event.Subscription, error) {

	logs, sub, err := _USDCTokenPool.contract.WatchLogs(opts, "ChainRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCTokenPoolChainRemoved)
				if err := _USDCTokenPool.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
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

func (_USDCTokenPool *USDCTokenPoolFilterer) ParseChainRemoved(log types.Log) (*USDCTokenPoolChainRemoved, error) {
	event := new(USDCTokenPoolChainRemoved)
	if err := _USDCTokenPool.contract.UnpackLog(event, "ChainRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type USDCTokenPoolConfigChangedIterator struct {
	Event *USDCTokenPoolConfigChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCTokenPoolConfigChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCTokenPoolConfigChanged)
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
		it.Event = new(USDCTokenPoolConfigChanged)
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

func (it *USDCTokenPoolConfigChangedIterator) Error() error {
	return it.fail
}

func (it *USDCTokenPoolConfigChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCTokenPoolConfigChanged struct {
	Config RateLimiterConfig
	Raw    types.Log
}

func (_USDCTokenPool *USDCTokenPoolFilterer) FilterConfigChanged(opts *bind.FilterOpts) (*USDCTokenPoolConfigChangedIterator, error) {

	logs, sub, err := _USDCTokenPool.contract.FilterLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolConfigChangedIterator{contract: _USDCTokenPool.contract, event: "ConfigChanged", logs: logs, sub: sub}, nil
}

func (_USDCTokenPool *USDCTokenPoolFilterer) WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolConfigChanged) (event.Subscription, error) {

	logs, sub, err := _USDCTokenPool.contract.WatchLogs(opts, "ConfigChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCTokenPoolConfigChanged)
				if err := _USDCTokenPool.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
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

func (_USDCTokenPool *USDCTokenPoolFilterer) ParseConfigChanged(log types.Log) (*USDCTokenPoolConfigChanged, error) {
	event := new(USDCTokenPoolConfigChanged)
	if err := _USDCTokenPool.contract.UnpackLog(event, "ConfigChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type USDCTokenPoolConfigSetIterator struct {
	Event *USDCTokenPoolConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCTokenPoolConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCTokenPoolConfigSet)
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
		it.Event = new(USDCTokenPoolConfigSet)
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

func (it *USDCTokenPoolConfigSetIterator) Error() error {
	return it.fail
}

func (it *USDCTokenPoolConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCTokenPoolConfigSet struct {
	TokenMessenger common.Address
	Raw            types.Log
}

func (_USDCTokenPool *USDCTokenPoolFilterer) FilterConfigSet(opts *bind.FilterOpts) (*USDCTokenPoolConfigSetIterator, error) {

	logs, sub, err := _USDCTokenPool.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolConfigSetIterator{contract: _USDCTokenPool.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_USDCTokenPool *USDCTokenPoolFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolConfigSet) (event.Subscription, error) {

	logs, sub, err := _USDCTokenPool.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCTokenPoolConfigSet)
				if err := _USDCTokenPool.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_USDCTokenPool *USDCTokenPoolFilterer) ParseConfigSet(log types.Log) (*USDCTokenPoolConfigSet, error) {
	event := new(USDCTokenPoolConfigSet)
	if err := _USDCTokenPool.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type USDCTokenPoolDomainsSetIterator struct {
	Event *USDCTokenPoolDomainsSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCTokenPoolDomainsSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCTokenPoolDomainsSet)
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
		it.Event = new(USDCTokenPoolDomainsSet)
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

func (it *USDCTokenPoolDomainsSetIterator) Error() error {
	return it.fail
}

func (it *USDCTokenPoolDomainsSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCTokenPoolDomainsSet struct {
	Arg0 []USDCTokenPoolDomainUpdate
	Raw  types.Log
}

func (_USDCTokenPool *USDCTokenPoolFilterer) FilterDomainsSet(opts *bind.FilterOpts) (*USDCTokenPoolDomainsSetIterator, error) {

	logs, sub, err := _USDCTokenPool.contract.FilterLogs(opts, "DomainsSet")
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolDomainsSetIterator{contract: _USDCTokenPool.contract, event: "DomainsSet", logs: logs, sub: sub}, nil
}

func (_USDCTokenPool *USDCTokenPoolFilterer) WatchDomainsSet(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolDomainsSet) (event.Subscription, error) {

	logs, sub, err := _USDCTokenPool.contract.WatchLogs(opts, "DomainsSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCTokenPoolDomainsSet)
				if err := _USDCTokenPool.contract.UnpackLog(event, "DomainsSet", log); err != nil {
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

func (_USDCTokenPool *USDCTokenPoolFilterer) ParseDomainsSet(log types.Log) (*USDCTokenPoolDomainsSet, error) {
	event := new(USDCTokenPoolDomainsSet)
	if err := _USDCTokenPool.contract.UnpackLog(event, "DomainsSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type USDCTokenPoolLockedIterator struct {
	Event *USDCTokenPoolLocked

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCTokenPoolLockedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCTokenPoolLocked)
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
		it.Event = new(USDCTokenPoolLocked)
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

func (it *USDCTokenPoolLockedIterator) Error() error {
	return it.fail
}

func (it *USDCTokenPoolLockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCTokenPoolLocked struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_USDCTokenPool *USDCTokenPoolFilterer) FilterLocked(opts *bind.FilterOpts, sender []common.Address) (*USDCTokenPoolLockedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _USDCTokenPool.contract.FilterLogs(opts, "Locked", senderRule)
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolLockedIterator{contract: _USDCTokenPool.contract, event: "Locked", logs: logs, sub: sub}, nil
}

func (_USDCTokenPool *USDCTokenPoolFilterer) WatchLocked(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolLocked, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _USDCTokenPool.contract.WatchLogs(opts, "Locked", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCTokenPoolLocked)
				if err := _USDCTokenPool.contract.UnpackLog(event, "Locked", log); err != nil {
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

func (_USDCTokenPool *USDCTokenPoolFilterer) ParseLocked(log types.Log) (*USDCTokenPoolLocked, error) {
	event := new(USDCTokenPoolLocked)
	if err := _USDCTokenPool.contract.UnpackLog(event, "Locked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type USDCTokenPoolMintedIterator struct {
	Event *USDCTokenPoolMinted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCTokenPoolMintedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCTokenPoolMinted)
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
		it.Event = new(USDCTokenPoolMinted)
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

func (it *USDCTokenPoolMintedIterator) Error() error {
	return it.fail
}

func (it *USDCTokenPoolMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCTokenPoolMinted struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log
}

func (_USDCTokenPool *USDCTokenPoolFilterer) FilterMinted(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*USDCTokenPoolMintedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _USDCTokenPool.contract.FilterLogs(opts, "Minted", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolMintedIterator{contract: _USDCTokenPool.contract, event: "Minted", logs: logs, sub: sub}, nil
}

func (_USDCTokenPool *USDCTokenPoolFilterer) WatchMinted(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolMinted, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _USDCTokenPool.contract.WatchLogs(opts, "Minted", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCTokenPoolMinted)
				if err := _USDCTokenPool.contract.UnpackLog(event, "Minted", log); err != nil {
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

func (_USDCTokenPool *USDCTokenPoolFilterer) ParseMinted(log types.Log) (*USDCTokenPoolMinted, error) {
	event := new(USDCTokenPoolMinted)
	if err := _USDCTokenPool.contract.UnpackLog(event, "Minted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type USDCTokenPoolOwnershipTransferRequestedIterator struct {
	Event *USDCTokenPoolOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCTokenPoolOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCTokenPoolOwnershipTransferRequested)
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
		it.Event = new(USDCTokenPoolOwnershipTransferRequested)
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

func (it *USDCTokenPoolOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *USDCTokenPoolOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCTokenPoolOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_USDCTokenPool *USDCTokenPoolFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*USDCTokenPoolOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _USDCTokenPool.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolOwnershipTransferRequestedIterator{contract: _USDCTokenPool.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_USDCTokenPool *USDCTokenPoolFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _USDCTokenPool.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCTokenPoolOwnershipTransferRequested)
				if err := _USDCTokenPool.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_USDCTokenPool *USDCTokenPoolFilterer) ParseOwnershipTransferRequested(log types.Log) (*USDCTokenPoolOwnershipTransferRequested, error) {
	event := new(USDCTokenPoolOwnershipTransferRequested)
	if err := _USDCTokenPool.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type USDCTokenPoolOwnershipTransferredIterator struct {
	Event *USDCTokenPoolOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCTokenPoolOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCTokenPoolOwnershipTransferred)
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
		it.Event = new(USDCTokenPoolOwnershipTransferred)
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

func (it *USDCTokenPoolOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *USDCTokenPoolOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCTokenPoolOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_USDCTokenPool *USDCTokenPoolFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*USDCTokenPoolOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _USDCTokenPool.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolOwnershipTransferredIterator{contract: _USDCTokenPool.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_USDCTokenPool *USDCTokenPoolFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _USDCTokenPool.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCTokenPoolOwnershipTransferred)
				if err := _USDCTokenPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_USDCTokenPool *USDCTokenPoolFilterer) ParseOwnershipTransferred(log types.Log) (*USDCTokenPoolOwnershipTransferred, error) {
	event := new(USDCTokenPoolOwnershipTransferred)
	if err := _USDCTokenPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type USDCTokenPoolRateLimitAdminSetIterator struct {
	Event *USDCTokenPoolRateLimitAdminSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCTokenPoolRateLimitAdminSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCTokenPoolRateLimitAdminSet)
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
		it.Event = new(USDCTokenPoolRateLimitAdminSet)
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

func (it *USDCTokenPoolRateLimitAdminSetIterator) Error() error {
	return it.fail
}

func (it *USDCTokenPoolRateLimitAdminSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCTokenPoolRateLimitAdminSet struct {
	RateLimitAdmin common.Address
	Raw            types.Log
}

func (_USDCTokenPool *USDCTokenPoolFilterer) FilterRateLimitAdminSet(opts *bind.FilterOpts) (*USDCTokenPoolRateLimitAdminSetIterator, error) {

	logs, sub, err := _USDCTokenPool.contract.FilterLogs(opts, "RateLimitAdminSet")
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolRateLimitAdminSetIterator{contract: _USDCTokenPool.contract, event: "RateLimitAdminSet", logs: logs, sub: sub}, nil
}

func (_USDCTokenPool *USDCTokenPoolFilterer) WatchRateLimitAdminSet(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolRateLimitAdminSet) (event.Subscription, error) {

	logs, sub, err := _USDCTokenPool.contract.WatchLogs(opts, "RateLimitAdminSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCTokenPoolRateLimitAdminSet)
				if err := _USDCTokenPool.contract.UnpackLog(event, "RateLimitAdminSet", log); err != nil {
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

func (_USDCTokenPool *USDCTokenPoolFilterer) ParseRateLimitAdminSet(log types.Log) (*USDCTokenPoolRateLimitAdminSet, error) {
	event := new(USDCTokenPoolRateLimitAdminSet)
	if err := _USDCTokenPool.contract.UnpackLog(event, "RateLimitAdminSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type USDCTokenPoolReleasedIterator struct {
	Event *USDCTokenPoolReleased

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCTokenPoolReleasedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCTokenPoolReleased)
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
		it.Event = new(USDCTokenPoolReleased)
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

func (it *USDCTokenPoolReleasedIterator) Error() error {
	return it.fail
}

func (it *USDCTokenPoolReleasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCTokenPoolReleased struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log
}

func (_USDCTokenPool *USDCTokenPoolFilterer) FilterReleased(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*USDCTokenPoolReleasedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _USDCTokenPool.contract.FilterLogs(opts, "Released", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolReleasedIterator{contract: _USDCTokenPool.contract, event: "Released", logs: logs, sub: sub}, nil
}

func (_USDCTokenPool *USDCTokenPoolFilterer) WatchReleased(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolReleased, sender []common.Address, recipient []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _USDCTokenPool.contract.WatchLogs(opts, "Released", senderRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCTokenPoolReleased)
				if err := _USDCTokenPool.contract.UnpackLog(event, "Released", log); err != nil {
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

func (_USDCTokenPool *USDCTokenPoolFilterer) ParseReleased(log types.Log) (*USDCTokenPoolReleased, error) {
	event := new(USDCTokenPoolReleased)
	if err := _USDCTokenPool.contract.UnpackLog(event, "Released", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type USDCTokenPoolRemotePoolAddedIterator struct {
	Event *USDCTokenPoolRemotePoolAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCTokenPoolRemotePoolAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCTokenPoolRemotePoolAdded)
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
		it.Event = new(USDCTokenPoolRemotePoolAdded)
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

func (it *USDCTokenPoolRemotePoolAddedIterator) Error() error {
	return it.fail
}

func (it *USDCTokenPoolRemotePoolAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCTokenPoolRemotePoolAdded struct {
	RemoteChainSelector uint64
	RemotePoolAddress   []byte
	Raw                 types.Log
}

func (_USDCTokenPool *USDCTokenPoolFilterer) FilterRemotePoolAdded(opts *bind.FilterOpts, remoteChainSelector []uint64) (*USDCTokenPoolRemotePoolAddedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _USDCTokenPool.contract.FilterLogs(opts, "RemotePoolAdded", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolRemotePoolAddedIterator{contract: _USDCTokenPool.contract, event: "RemotePoolAdded", logs: logs, sub: sub}, nil
}

func (_USDCTokenPool *USDCTokenPoolFilterer) WatchRemotePoolAdded(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolRemotePoolAdded, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _USDCTokenPool.contract.WatchLogs(opts, "RemotePoolAdded", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCTokenPoolRemotePoolAdded)
				if err := _USDCTokenPool.contract.UnpackLog(event, "RemotePoolAdded", log); err != nil {
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

func (_USDCTokenPool *USDCTokenPoolFilterer) ParseRemotePoolAdded(log types.Log) (*USDCTokenPoolRemotePoolAdded, error) {
	event := new(USDCTokenPoolRemotePoolAdded)
	if err := _USDCTokenPool.contract.UnpackLog(event, "RemotePoolAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type USDCTokenPoolRemotePoolRemovedIterator struct {
	Event *USDCTokenPoolRemotePoolRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCTokenPoolRemotePoolRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCTokenPoolRemotePoolRemoved)
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
		it.Event = new(USDCTokenPoolRemotePoolRemoved)
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

func (it *USDCTokenPoolRemotePoolRemovedIterator) Error() error {
	return it.fail
}

func (it *USDCTokenPoolRemotePoolRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCTokenPoolRemotePoolRemoved struct {
	RemoteChainSelector uint64
	RemotePoolAddress   []byte
	Raw                 types.Log
}

func (_USDCTokenPool *USDCTokenPoolFilterer) FilterRemotePoolRemoved(opts *bind.FilterOpts, remoteChainSelector []uint64) (*USDCTokenPoolRemotePoolRemovedIterator, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _USDCTokenPool.contract.FilterLogs(opts, "RemotePoolRemoved", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolRemotePoolRemovedIterator{contract: _USDCTokenPool.contract, event: "RemotePoolRemoved", logs: logs, sub: sub}, nil
}

func (_USDCTokenPool *USDCTokenPoolFilterer) WatchRemotePoolRemoved(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolRemotePoolRemoved, remoteChainSelector []uint64) (event.Subscription, error) {

	var remoteChainSelectorRule []interface{}
	for _, remoteChainSelectorItem := range remoteChainSelector {
		remoteChainSelectorRule = append(remoteChainSelectorRule, remoteChainSelectorItem)
	}

	logs, sub, err := _USDCTokenPool.contract.WatchLogs(opts, "RemotePoolRemoved", remoteChainSelectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCTokenPoolRemotePoolRemoved)
				if err := _USDCTokenPool.contract.UnpackLog(event, "RemotePoolRemoved", log); err != nil {
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

func (_USDCTokenPool *USDCTokenPoolFilterer) ParseRemotePoolRemoved(log types.Log) (*USDCTokenPoolRemotePoolRemoved, error) {
	event := new(USDCTokenPoolRemotePoolRemoved)
	if err := _USDCTokenPool.contract.UnpackLog(event, "RemotePoolRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type USDCTokenPoolRouterUpdatedIterator struct {
	Event *USDCTokenPoolRouterUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCTokenPoolRouterUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCTokenPoolRouterUpdated)
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
		it.Event = new(USDCTokenPoolRouterUpdated)
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

func (it *USDCTokenPoolRouterUpdatedIterator) Error() error {
	return it.fail
}

func (it *USDCTokenPoolRouterUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCTokenPoolRouterUpdated struct {
	OldRouter common.Address
	NewRouter common.Address
	Raw       types.Log
}

func (_USDCTokenPool *USDCTokenPoolFilterer) FilterRouterUpdated(opts *bind.FilterOpts) (*USDCTokenPoolRouterUpdatedIterator, error) {

	logs, sub, err := _USDCTokenPool.contract.FilterLogs(opts, "RouterUpdated")
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolRouterUpdatedIterator{contract: _USDCTokenPool.contract, event: "RouterUpdated", logs: logs, sub: sub}, nil
}

func (_USDCTokenPool *USDCTokenPoolFilterer) WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolRouterUpdated) (event.Subscription, error) {

	logs, sub, err := _USDCTokenPool.contract.WatchLogs(opts, "RouterUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCTokenPoolRouterUpdated)
				if err := _USDCTokenPool.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
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

func (_USDCTokenPool *USDCTokenPoolFilterer) ParseRouterUpdated(log types.Log) (*USDCTokenPoolRouterUpdated, error) {
	event := new(USDCTokenPoolRouterUpdated)
	if err := _USDCTokenPool.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type USDCTokenPoolTokensConsumedIterator struct {
	Event *USDCTokenPoolTokensConsumed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCTokenPoolTokensConsumedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCTokenPoolTokensConsumed)
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
		it.Event = new(USDCTokenPoolTokensConsumed)
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

func (it *USDCTokenPoolTokensConsumedIterator) Error() error {
	return it.fail
}

func (it *USDCTokenPoolTokensConsumedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCTokenPoolTokensConsumed struct {
	Tokens *big.Int
	Raw    types.Log
}

func (_USDCTokenPool *USDCTokenPoolFilterer) FilterTokensConsumed(opts *bind.FilterOpts) (*USDCTokenPoolTokensConsumedIterator, error) {

	logs, sub, err := _USDCTokenPool.contract.FilterLogs(opts, "TokensConsumed")
	if err != nil {
		return nil, err
	}
	return &USDCTokenPoolTokensConsumedIterator{contract: _USDCTokenPool.contract, event: "TokensConsumed", logs: logs, sub: sub}, nil
}

func (_USDCTokenPool *USDCTokenPoolFilterer) WatchTokensConsumed(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolTokensConsumed) (event.Subscription, error) {

	logs, sub, err := _USDCTokenPool.contract.WatchLogs(opts, "TokensConsumed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCTokenPoolTokensConsumed)
				if err := _USDCTokenPool.contract.UnpackLog(event, "TokensConsumed", log); err != nil {
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

func (_USDCTokenPool *USDCTokenPoolFilterer) ParseTokensConsumed(log types.Log) (*USDCTokenPoolTokensConsumed, error) {
	event := new(USDCTokenPoolTokensConsumed)
	if err := _USDCTokenPool.contract.UnpackLog(event, "TokensConsumed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_USDCTokenPool *USDCTokenPool) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _USDCTokenPool.abi.Events["AllowListAdd"].ID:
		return _USDCTokenPool.ParseAllowListAdd(log)
	case _USDCTokenPool.abi.Events["AllowListRemove"].ID:
		return _USDCTokenPool.ParseAllowListRemove(log)
	case _USDCTokenPool.abi.Events["Burned"].ID:
		return _USDCTokenPool.ParseBurned(log)
	case _USDCTokenPool.abi.Events["ChainAdded"].ID:
		return _USDCTokenPool.ParseChainAdded(log)
	case _USDCTokenPool.abi.Events["ChainConfigured"].ID:
		return _USDCTokenPool.ParseChainConfigured(log)
	case _USDCTokenPool.abi.Events["ChainRemoved"].ID:
		return _USDCTokenPool.ParseChainRemoved(log)
	case _USDCTokenPool.abi.Events["ConfigChanged"].ID:
		return _USDCTokenPool.ParseConfigChanged(log)
	case _USDCTokenPool.abi.Events["ConfigSet"].ID:
		return _USDCTokenPool.ParseConfigSet(log)
	case _USDCTokenPool.abi.Events["DomainsSet"].ID:
		return _USDCTokenPool.ParseDomainsSet(log)
	case _USDCTokenPool.abi.Events["Locked"].ID:
		return _USDCTokenPool.ParseLocked(log)
	case _USDCTokenPool.abi.Events["Minted"].ID:
		return _USDCTokenPool.ParseMinted(log)
	case _USDCTokenPool.abi.Events["OwnershipTransferRequested"].ID:
		return _USDCTokenPool.ParseOwnershipTransferRequested(log)
	case _USDCTokenPool.abi.Events["OwnershipTransferred"].ID:
		return _USDCTokenPool.ParseOwnershipTransferred(log)
	case _USDCTokenPool.abi.Events["RateLimitAdminSet"].ID:
		return _USDCTokenPool.ParseRateLimitAdminSet(log)
	case _USDCTokenPool.abi.Events["Released"].ID:
		return _USDCTokenPool.ParseReleased(log)
	case _USDCTokenPool.abi.Events["RemotePoolAdded"].ID:
		return _USDCTokenPool.ParseRemotePoolAdded(log)
	case _USDCTokenPool.abi.Events["RemotePoolRemoved"].ID:
		return _USDCTokenPool.ParseRemotePoolRemoved(log)
	case _USDCTokenPool.abi.Events["RouterUpdated"].ID:
		return _USDCTokenPool.ParseRouterUpdated(log)
	case _USDCTokenPool.abi.Events["TokensConsumed"].ID:
		return _USDCTokenPool.ParseTokensConsumed(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (USDCTokenPoolAllowListAdd) Topic() common.Hash {
	return common.HexToHash("0x2640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d8")
}

func (USDCTokenPoolAllowListRemove) Topic() common.Hash {
	return common.HexToHash("0x800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf7566")
}

func (USDCTokenPoolBurned) Topic() common.Hash {
	return common.HexToHash("0x696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df7")
}

func (USDCTokenPoolChainAdded) Topic() common.Hash {
	return common.HexToHash("0x8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c2")
}

func (USDCTokenPoolChainConfigured) Topic() common.Hash {
	return common.HexToHash("0x0350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b")
}

func (USDCTokenPoolChainRemoved) Topic() common.Hash {
	return common.HexToHash("0x5204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d859916")
}

func (USDCTokenPoolConfigChanged) Topic() common.Hash {
	return common.HexToHash("0x9ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c19")
}

func (USDCTokenPoolConfigSet) Topic() common.Hash {
	return common.HexToHash("0x2e902d38f15b233cbb63711add0fca4545334d3a169d60c0a616494d7eea9544")
}

func (USDCTokenPoolDomainsSet) Topic() common.Hash {
	return common.HexToHash("0x1889010d2535a0ab1643678d1da87fbbe8b87b2f585b47ddb72ec622aef9ee56")
}

func (USDCTokenPoolLocked) Topic() common.Hash {
	return common.HexToHash("0x9f1ec8c880f76798e7b793325d625e9b60e4082a553c98f42b6cda368dd60008")
}

func (USDCTokenPoolMinted) Topic() common.Hash {
	return common.HexToHash("0x9d228d69b5fdb8d273a2336f8fb8612d039631024ea9bf09c424a9503aa078f0")
}

func (USDCTokenPoolOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (USDCTokenPoolOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (USDCTokenPoolRateLimitAdminSet) Topic() common.Hash {
	return common.HexToHash("0x44676b5284b809a22248eba0da87391d79098be38bb03154be88a58bf4d09174")
}

func (USDCTokenPoolReleased) Topic() common.Hash {
	return common.HexToHash("0x2d87480f50083e2b2759522a8fdda59802650a8055e609a7772cf70c07748f52")
}

func (USDCTokenPoolRemotePoolAdded) Topic() common.Hash {
	return common.HexToHash("0x7d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea")
}

func (USDCTokenPoolRemotePoolRemoved) Topic() common.Hash {
	return common.HexToHash("0x52d00ee4d9bd51b40168f2afc5848837288ce258784ad914278791464b3f4d76")
}

func (USDCTokenPoolRouterUpdated) Topic() common.Hash {
	return common.HexToHash("0x02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f1684")
}

func (USDCTokenPoolTokensConsumed) Topic() common.Hash {
	return common.HexToHash("0x1871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690a")
}

func (_USDCTokenPool *USDCTokenPool) Address() common.Address {
	return _USDCTokenPool.address
}

type USDCTokenPoolInterface interface {
	SUPPORTEDUSDCVERSION(opts *bind.CallOpts) (uint32, error)

	GetAllowList(opts *bind.CallOpts) ([]common.Address, error)

	GetAllowListEnabled(opts *bind.CallOpts) (bool, error)

	GetCurrentInboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error)

	GetCurrentOutboundRateLimiterState(opts *bind.CallOpts, remoteChainSelector uint64) (RateLimiterTokenBucket, error)

	GetDomain(opts *bind.CallOpts, chainSelector uint64) (USDCTokenPoolDomain, error)

	GetRateLimitAdmin(opts *bind.CallOpts) (common.Address, error)

	GetRemotePools(opts *bind.CallOpts, remoteChainSelector uint64) ([][]byte, error)

	GetRemoteToken(opts *bind.CallOpts, remoteChainSelector uint64) ([]byte, error)

	GetRmnProxy(opts *bind.CallOpts) (common.Address, error)

	GetRouter(opts *bind.CallOpts) (common.Address, error)

	GetSupportedChains(opts *bind.CallOpts) ([]uint64, error)

	GetToken(opts *bind.CallOpts) (common.Address, error)

	GetTokenDecimals(opts *bind.CallOpts) (uint8, error)

	ILocalDomainIdentifier(opts *bind.CallOpts) (uint32, error)

	IMessageTransmitter(opts *bind.CallOpts) (common.Address, error)

	ITokenMessenger(opts *bind.CallOpts) (common.Address, error)

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

	SetDomains(opts *bind.TransactOpts, domains []USDCTokenPoolDomainUpdate) (*types.Transaction, error)

	SetRateLimitAdmin(opts *bind.TransactOpts, rateLimitAdmin common.Address) (*types.Transaction, error)

	SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterAllowListAdd(opts *bind.FilterOpts) (*USDCTokenPoolAllowListAddIterator, error)

	WatchAllowListAdd(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolAllowListAdd) (event.Subscription, error)

	ParseAllowListAdd(log types.Log) (*USDCTokenPoolAllowListAdd, error)

	FilterAllowListRemove(opts *bind.FilterOpts) (*USDCTokenPoolAllowListRemoveIterator, error)

	WatchAllowListRemove(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolAllowListRemove) (event.Subscription, error)

	ParseAllowListRemove(log types.Log) (*USDCTokenPoolAllowListRemove, error)

	FilterBurned(opts *bind.FilterOpts, sender []common.Address) (*USDCTokenPoolBurnedIterator, error)

	WatchBurned(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolBurned, sender []common.Address) (event.Subscription, error)

	ParseBurned(log types.Log) (*USDCTokenPoolBurned, error)

	FilterChainAdded(opts *bind.FilterOpts) (*USDCTokenPoolChainAddedIterator, error)

	WatchChainAdded(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolChainAdded) (event.Subscription, error)

	ParseChainAdded(log types.Log) (*USDCTokenPoolChainAdded, error)

	FilterChainConfigured(opts *bind.FilterOpts) (*USDCTokenPoolChainConfiguredIterator, error)

	WatchChainConfigured(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolChainConfigured) (event.Subscription, error)

	ParseChainConfigured(log types.Log) (*USDCTokenPoolChainConfigured, error)

	FilterChainRemoved(opts *bind.FilterOpts) (*USDCTokenPoolChainRemovedIterator, error)

	WatchChainRemoved(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolChainRemoved) (event.Subscription, error)

	ParseChainRemoved(log types.Log) (*USDCTokenPoolChainRemoved, error)

	FilterConfigChanged(opts *bind.FilterOpts) (*USDCTokenPoolConfigChangedIterator, error)

	WatchConfigChanged(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolConfigChanged) (event.Subscription, error)

	ParseConfigChanged(log types.Log) (*USDCTokenPoolConfigChanged, error)

	FilterConfigSet(opts *bind.FilterOpts) (*USDCTokenPoolConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*USDCTokenPoolConfigSet, error)

	FilterDomainsSet(opts *bind.FilterOpts) (*USDCTokenPoolDomainsSetIterator, error)

	WatchDomainsSet(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolDomainsSet) (event.Subscription, error)

	ParseDomainsSet(log types.Log) (*USDCTokenPoolDomainsSet, error)

	FilterLocked(opts *bind.FilterOpts, sender []common.Address) (*USDCTokenPoolLockedIterator, error)

	WatchLocked(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolLocked, sender []common.Address) (event.Subscription, error)

	ParseLocked(log types.Log) (*USDCTokenPoolLocked, error)

	FilterMinted(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*USDCTokenPoolMintedIterator, error)

	WatchMinted(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolMinted, sender []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseMinted(log types.Log) (*USDCTokenPoolMinted, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*USDCTokenPoolOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*USDCTokenPoolOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*USDCTokenPoolOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*USDCTokenPoolOwnershipTransferred, error)

	FilterRateLimitAdminSet(opts *bind.FilterOpts) (*USDCTokenPoolRateLimitAdminSetIterator, error)

	WatchRateLimitAdminSet(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolRateLimitAdminSet) (event.Subscription, error)

	ParseRateLimitAdminSet(log types.Log) (*USDCTokenPoolRateLimitAdminSet, error)

	FilterReleased(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address) (*USDCTokenPoolReleasedIterator, error)

	WatchReleased(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolReleased, sender []common.Address, recipient []common.Address) (event.Subscription, error)

	ParseReleased(log types.Log) (*USDCTokenPoolReleased, error)

	FilterRemotePoolAdded(opts *bind.FilterOpts, remoteChainSelector []uint64) (*USDCTokenPoolRemotePoolAddedIterator, error)

	WatchRemotePoolAdded(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolRemotePoolAdded, remoteChainSelector []uint64) (event.Subscription, error)

	ParseRemotePoolAdded(log types.Log) (*USDCTokenPoolRemotePoolAdded, error)

	FilterRemotePoolRemoved(opts *bind.FilterOpts, remoteChainSelector []uint64) (*USDCTokenPoolRemotePoolRemovedIterator, error)

	WatchRemotePoolRemoved(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolRemotePoolRemoved, remoteChainSelector []uint64) (event.Subscription, error)

	ParseRemotePoolRemoved(log types.Log) (*USDCTokenPoolRemotePoolRemoved, error)

	FilterRouterUpdated(opts *bind.FilterOpts) (*USDCTokenPoolRouterUpdatedIterator, error)

	WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolRouterUpdated) (event.Subscription, error)

	ParseRouterUpdated(log types.Log) (*USDCTokenPoolRouterUpdated, error)

	FilterTokensConsumed(opts *bind.FilterOpts) (*USDCTokenPoolTokensConsumedIterator, error)

	WatchTokensConsumed(opts *bind.WatchOpts, sink chan<- *USDCTokenPoolTokensConsumed) (event.Subscription, error)

	ParseTokensConsumed(log types.Log) (*USDCTokenPoolTokensConsumed, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}

var USDCTokenPoolZKBin = ("0x000400000000000200250000000000020000006003100270000007c10030019d000007c103300197000300000031035500020000000103550000000100200190000000320000c13d000000800b0000390000004000b0043f000000040030008c00000a830000413d000000000201043b000000e002200270000007ea0020009c000000730000a13d000007eb0020009c000000840000213d000007f80020009c000000b40000a13d000007f90020009c000001490000a13d000007fa0020009c000005cb0000613d000007fb0020009c000004cc0000613d000007fc0020009c00000a830000c13d000000240030008c00000a830000413d0000000002000416000000000002004b00000a830000c13d0000000401100370000000000101043b000007c40010009c00000a830000213d0000000102000039000000000202041a000007c4022001970000000003000411000000000023004b00000a190000c13d000000000001004b00000a6a0000c13d0000083301000041000000800010043f000008210100004100001f03000104300000016004000039000000400040043f0000000002000416000000000002004b00000a830000c13d0000001f02300039000007c2022001970000016002200039000000400020043f0000001f0530018f000007c3063001980000016002600039000000440000613d000000000701034f000000007807043c0000000004840436000000000024004b000000400000c13d000000000005004b000000510000613d000000000161034f0000000304500210000000000502043300000000054501cf000000000545022f000000000101043b0000010004400089000000000141022f00000000014101cf000000000151019f0000000000120435000000a00030008c00000a830000413d000001600100043d000e00000001001d000007c40010009c00000a830000213d000001800100043d000007c40010009c00000a830000213d000001a00200043d000007c50020009c00000a830000213d0000001f04200039000000000034004b0000000005000019000007c605008041000007c604400197000000000004004b0000000006000019000007c606004041000007c60040009c000000000605c019000000000006004b00000a830000c13d00000160042000390000000004040433000007c50040009c000001d50000a13d0000085801000041000000000010043f0000004101000039000000040010043f000007d50100004100001f0300010430000008040020009c000000990000a13d000008050020009c000000ec0000a13d000008060020009c000001980000a13d000008070020009c000009fb0000613d000008080020009c000004ff0000613d000008090020009c00000a830000c13d0000000001000416000000000001004b00000a830000c13d0000000101000039000005cf0000013d000007ec0020009c000000c00000a13d000007ed0020009c000001810000a13d000007ee0020009c000005d40000613d000007ef0020009c000004e50000613d000007f00020009c00000a830000c13d0000000001000416000000000001004b00000a830000c13d0000000001000412001500000001001d001400a00000003d000080050100003900000044030000390000000004000415000000150440008a000001c30000013d000008110020009c000000fd0000213d000008170020009c000001b50000213d0000081a0020009c000005130000613d0000081b0020009c00000a830000c13d000000240030008c00000a830000413d0000000002000416000000000002004b00000a830000c13d0000000401100370000000000201043b000008710020019800000a830000c13d0000000101000039000008720020009c000005d10000613d000008730020009c000005d10000613d000008740020009c000000000100c019000000800010043f0000081d0100004100001f020001042e000007ff0020009c000001140000213d000008020020009c000002010000613d000008030020009c00000a830000c13d0000000001000416000000000001004b00000a830000c13d000000800000043f0000081d0100004100001f020001042e000007f30020009c000001250000213d000007f60020009c000002550000613d000007f70020009c00000a830000c13d000000240030008c00000a830000413d0000000002000416000000000002004b00000a830000c13d0000000401100370000000000101043b001300000001001d000007c50010009c00000a830000213d1f011c370000040f0000001301000029000000000010043f0000000701000039000000200010043f000000400200003900000000010000191f011ec40000040f001200000001001d000000400100043d001300000001001d1f011ae40000040f0000001205000029000000000405041a0000083a004001980000000002000039000000010200c0390000001301000029000000400310003900000000002304350000008002400270000007c10220019700000020031000390000000000230435000008290240019700000000002104350000000102500039000001730000013d0000080c0020009c000001340000213d0000080f0020009c0000035d0000613d000008100020009c00000a830000c13d0000000001000416000000000001004b00000a830000c13d0000000001000412001d00000001001d001c00800000003d0000800501000039000000440300003900000000040004150000001d0440008a000001c30000013d000008120020009c000001ca0000213d000008150020009c0000054a0000613d000008160020009c00000a830000c13d0000000001000416000000000001004b00000a830000c13d0000000001000412001f00000001001d001e00200000003d0000800501000039000000440300003900000000040004150000001f0440008a00000005044002100000081c020000411f011ed90000040f000000ff0110018f000000800010043f0000081d0100004100001f020001042e000008000020009c000002990000613d000008010020009c00000a830000c13d0000000001000416000000000001004b00000a830000c13d0000000202000039000000000102041a000000800010043f000000000020043f0000002002000039000000000001004b00000a310000c13d000000a001000039000000000402001900000a400000013d000007f40020009c000003010000613d000007f50020009c00000a830000c13d0000000001000416000000000001004b00000a830000c13d0000000001000412001900000001001d001800400000003d000080050100003900000044030000390000000004000415000000190440008a000001c30000013d0000080d0020009c000004210000613d0000080e0020009c00000a830000c13d0000000001000416000000000001004b00000a830000c13d0000000001000412001b00000001001d001a00c00000003d0000800501000039000000440300003900000000040004150000001b0440008a00000005044002100000081c020000411f011ed90000040f000007c101100197000000800010043f0000081d0100004100001f020001042e000007fd0020009c0000045d0000613d000007fe0020009c00000a830000c13d000000240030008c00000a830000413d0000000002000416000000000002004b00000a830000c13d0000000401100370000000000101043b001300000001001d000007c50010009c00000a830000213d1f011c370000040f0000001301000029000000000010043f0000000701000039000000200010043f000000400200003900000000010000191f011ec40000040f001200000001001d000000400100043d001300000001001d1f011ae40000040f00000012050000290000000201500039000000000401041a0000083a004001980000000002000039000000010200c0390000001301000029000000400310003900000000002304350000008002400270000007c10220019700000020031000390000000000230435000008290240019700000000002104350000000302500039000000000402041a0000008002100039000000800340027000000000003204350000082903400197000000600210003900000000003204351f011de50000040f000000400100043d001200000001001d00000013020000291f011b290000040f0000001202000029000004dc0000013d000007f10020009c0000049a0000613d000007f20020009c00000a830000c13d0000000001000416000000000001004b00000a830000c13d0000000001000412001700000001001d001600600000003d000080050100003900000044030000390000000004000415000000170440008a00000005044002100000081c020000411f011ed90000040f000000000001004b0000000001000039000000010100c039000000800010043f0000081d0100004100001f020001042e0000080a0020009c000004fa0000613d0000080b0020009c00000a830000c13d0000000001000416000000000001004b00000a830000c13d000000000100041a000007c4021001970000000006000411000000000026004b00000a590000c13d0000000102000039000000000302041a000007c804300197000000000464019f000000000042041b000007c801100197000000000010041b0000000001000414000007c405300197000007c10010009c000007c101008041000000c0011002100000081e011001c70000800d020000390000000303000039000008620400004100000a7a0000013d000008180020009c000005650000613d000008190020009c00000a830000c13d0000000001000416000000000001004b00000a830000c13d0000000001000412002300000001001d002200000000003d000080050100003900000044030000390000000004000415000000230440008a00000005044002100000081c020000411f011ed90000040f000007c401100197000000800010043f0000081d0100004100001f020001042e000008130020009c000005790000613d000008140020009c00000a830000c13d0000000001000416000000000001004b00000a830000c13d00000000010300191f011b0c0000040f1f011ba40000040f000005090000013d00000005054002100000003f06500039000007c706600197000000400700043d0000000006670019001100000007001d000000000076004b00000000070000390000000107004039000007c50060009c0000006d0000213d00000001007001900000006d0000c13d0000016007300039000000400060043f00000011030000290000000003430436001000000003001d00000180022000390000000003250019000000000073004b00000a830000213d000000000004004b000001f40000613d00000010040000290000000025020434000007c40050009c00000a830000213d0000000004540436000000000032004b000001ee0000413d000001c00300043d000007c40030009c00000a830000213d000001e00200043d001300000002001d000007c40020009c00000a830000213d0000000002000411000000000002004b00000a850000c13d000000400100043d000007e90200004100000b370000013d000000240030008c00000a830000413d0000000002000416000000000002004b00000a830000c13d0000000402100370000000000202043b001300000002001d000007c50020009c00000a830000213d000000130230006a000007df0020009c00000a830000213d000000a40020008c00000a830000413d000000c002000039000000400020043f0000006002000039000000800020043f000000a00020043f0000001302000029001200840020003d0000001201100360000000000101043b001100000001001d000007c40010009c00000a830000213d0000081c0100004100000000001004430000000001000412000000040010044300000024000004430000000001000414000007c10010009c000007c101008041000000c00110021000000847011001c700008005020000391f011efc0000040f000000010020019000001ae30000613d000000000101043b000000110210014f0000000201000367000007c40020019800000a7f0000c13d000000120200002900110060002000920000001101100360000000000101043b000007c50010009c00000a830000213d000000400300043d000008490200004100000000002304350000008001100210001200000003001d000000040230003900000000001204350000081c01000041000000000010044300000000010004120000000400100443000000400100003900000024001004430000000001000414000007c10010009c000007c101008041000000c00110021000000847011001c700008005020000391f011efc0000040f000000010020019000001ae30000613d000000000201043b0000000001000414000007c402200197000000040020008c00000cb90000c13d0000000103000031000000200030008c0000002004000039000000000403401900000ce30000013d0000000002000416000000000002004b00000a830000c13d0000000504000039000000000204041a000000800020043f000000000040043f000000000002004b00000a1d0000c13d000000a002000039000000400020043f0000002004000039000000000500001900000005065002100000003f076000390000083d077001970000000007270019000007c50070009c0000006d0000213d000000400070043f00000000005204350000001f0560018f000000a004400039000000000006004b000002750000613d000000000131034f00000000036400190000000006040019000000001701043c0000000006760436000000000036004b000002710000c13d000000000005004b000000800100043d000000000001004b000002870000613d00000000010000190000000003020433000000000013004b0000127b0000a13d00000005031002100000000005430019000000a0033000390000000003030433000007c50330019700000000003504350000000101100039000000800300043d000000000031004b0000027a0000413d000000400100043d00000020030000390000000003310436000000000402043300000000004304350000004003100039000000000004004b000002970000613d000000000500001900000020022000390000000006020433000007c50660019700000000036304360000000105500039000000000045004b000002900000413d000000000213004900000a510000013d000000240030008c00000a830000413d0000000002000416000000000002004b00000a830000c13d0000000401100370000000000101043b000007c50010009c00000a830000213d000000000010043f0000000701000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000101043b0000000501100039000000000301041a000000400200043d001100000002001d001300000003001d0000000002320436001000000002001d000000000010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007d0011001c700008010020000391f011efc0000040f000000010020019000000a830000613d0000001305000029000000000005004b0000001002000029000002cd0000613d000000000101043b00000010020000290000000003000019000000000401041a000000000242043600000001011000390000000103300039000000000053004b000002c70000413d000000110120006a0000001f011000390000087a011001970000001104100029000000000014004b00000000010000390000000101004039000007c50040009c0000006d0000213d00000001001001900000006d0000c13d000000400040043f00000011010000290000000002010433000007c50020009c0000006d0000213d00000005012002100000003f03100039000007c7033001970000000003430019000007c50030009c0000006d0000213d000000400030043f000f00000004001d0000000005240436000000000002004b000002ef0000613d00000060020000390000000003000019000000000435001900000000002404350000002003300039000000000013004b000002ea0000413d000e00000005001d00000011010000290000000001010433000000000001004b00000b9a0000c13d000000400100043d000000200200003900000000032104360000000f0200002900000000020204330000000000230435000000400310003900000005042002100000000005340019000000000002004b00000c9c0000c13d000000000215004900000a510000013d000000e40030008c00000a830000413d0000000002000416000000000002004b00000a830000c13d0000000402100370000000000202043b001300000002001d000007c50020009c00000a830000213d000000e002000039000000400020043f0000002402100370000000000202043b000000000002004b0000000003000039000000010300c039000000000032004b00000a830000c13d000000800020043f0000004402100370000000000202043b000008290020009c00000a830000213d000000a00020043f0000006402100370000000000202043b000008290020009c00000a830000213d000000c00020043f0000014002000039000000400020043f0000008402100370000000000202043b000000000002004b0000000003000039000000010300c039000000000032004b00000a830000c13d000000e00020043f000000a402100370000000000202043b000008290020009c00000a830000213d000001000020043f000000c401100370000000000101043b000008290010009c00000a830000213d000001200010043f0000000901000039000000000101041a000007c4021001970000000001000411000000000021004b0000033e0000613d0000000102000039000000000202041a000007c402200197000000000021004b00000f450000c13d0000001301000029000000000010043f0000000601000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000101043b000000000101041a000000000001004b0000048f0000613d000000c00100043d0000082901100197000000800200043d000000000002004b000010010000c13d000000000001004b000003590000c13d000000a00100043d0000082900100198000010070000613d000000400200043d001300000002001d0000082b01000041000010940000013d000000440030008c00000a830000413d0000000002000416000000000002004b00000a830000c13d0000000402100370000000000202043b000007c50020009c00000a830000213d0000002304200039000000000034004b00000a830000813d0000000404200039000000000441034f000000000804043b000007c50080009c00000a830000213d0000002406200039000000050a80021000000000076a0019000000000037004b00000a830000213d0000002402100370000000000202043b000007c50020009c00000a830000213d0000002304200039000000000034004b00000a830000813d0000000404200039000000000441034f000000000504043b000007c50050009c00000a830000213d000000240220003900000005095002100000000004290019000000000034004b00000a830000213d000000000c0b00190000000103000039000000000303041a000007c403300197000000000b00041100000000003b004b00000a190000c13d0000003f03a00039000007c7033001970000083c0030009c0000006d0000213d0000008003300039001100000003001d000000400030043f000000800080043f000000000008004b000003a00000613d000000000361034f000000000303043b000007c40030009c00000a830000213d000000200cc0003900000000003c04350000002006600039000000000076004b000003950000413d000000400300043d001100000003001d0000003f03900039000007c7033001970000001103300029000000110030006c00000000060000390000000106004039000007c50030009c0000006d0000213d00000001006001900000006d0000c13d000000400030043f00000011030000290000000003530436001000000003001d000000000005004b000003ba0000613d0000001103000029000000000521034f000000000505043b000007c40050009c00000a830000213d000000200330003900000000005304350000002002200039000000000042004b000003b10000413d0000081c01000041000000000010044300000000010004120000000400100443000000600100003900000024001004430000000001000414000007c10010009c000007c101008041000000c00110021000000847011001c700008005020000391f011efc0000040f000000010020019000001ae30000613d000000000101043b000000000001004b00000b970000613d000000800100043d000000000001004b000010e40000c13d00000011010000290000000001010433000000000001004b00000a7d0000613d0000000003000019000003da0000013d000000010330003900000011010000290000000001010433000000000013004b00000a7d0000813d000000050130021000000010011000290000000001010433001307c40010019c000003d50000613d0000001301000029000000000010043f0000000301000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c70000801002000039001200000003001d1f011efc0000040f0000001203000029000000010020019000000a830000613d000000000101043b000000000101041a000000000001004b000003d50000c13d0000000201000039000000000101041a000007c50010009c0000006d0000213d00000001021000390000000203000039000000000023041b000007cf0110009a0000001302000029000000000021041b000000000103041a000f00000001001d000000000020043f0000000301000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000101043b0000000f02000029000000000021041b000000400100043d00000013020000290000000000210435000007c10010009c000007c10100804100000040011002100000000002000414000007c10020009c000007c102008041000000c002200210000000000112019f000007d0011001c70000800d020000390000000103000039000007d1040000411f011ef70000040f00000012030000290000000100200190000003d50000c13d00000a830000013d000000440030008c00000a830000413d0000000002000416000000000002004b00000a830000c13d0000000402100370000000000202043b001300000002001d000007c50020009c00000a830000213d0000002402100370000000000202043b000007c50020009c00000a830000213d0000002304200039000000000034004b00000a830000813d0000000404200039000000000141034f000000000101043b001200000001001d000007c50010009c00000a830000213d0000002402200039001100000002001d0000001201200029000000000031004b00000a830000213d0000000101000039000000000101041a000007c4011001970000000002000411000000000012004b00000a190000c13d0000001301000029000000000010043f0000000601000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000101043b000000000101041a000000000001004b0000048f0000613d0000000003000031000000110100002900000012020000291f011b6c0000040f000000000201001900000013010000291f011ca70000040f000000000100001900001f020001042e000000440030008c00000a830000413d0000000002000416000000000002004b00000a830000c13d0000000402100370000000000202043b001300000002001d000007c50020009c00000a830000213d0000002402100370000000000202043b000007c50020009c00000a830000213d0000002304200039000000000034004b00000a830000813d001100040020003d0000001101100360000000000101043b001200000001001d000007c50010009c00000a830000213d0000001201200029001000240010003d000000100030006b00000a830000213d0000000101000039000000000101041a000007c4011001970000000002000411000000000012004b00000a190000c13d0000001301000029000000000010043f0000000601000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000101043b000000000101041a000000000001004b00000bff0000c13d000000400100043d00000824020000410000000000210435000000040210003900000013030000290000000000320435000007c10010009c000007c1010080410000004001100210000007d5011001c700001f0300010430000000240030008c00000a830000413d0000000002000416000000000002004b00000a830000c13d0000000401100370000000000101043b000007c50010009c00000a830000213d000000e002000039000000400020043f000000800000043f000000a00000043f000000c00000043f000000000010043f0000000a01000039000000200010043f000000400200003900000000010000191f011ec40000040f001300000001001d000000e0010000391f011aef0000040f0000001302000029000000000102041a000000e00010043f0000000102200039000000000202041a000007c103200197000001000030043f00000834002001980000000002000039000000010200c039000001200020043f000000400200043d0000000001120436000001000300043d000007c1033001970000000000310435000001200100043d000000000001004b0000000001000039000000010100c03900000040032000390000000000130435000007c10020009c000007c102008041000000400120021000000835011001c700001f020001042e000000240030008c00000a830000413d0000000002000416000000000002004b00000a830000c13d0000000401100370000000000101043b000007c50010009c00000a830000213d1f011c4c0000040f0000002002000039000000400300043d001300000003001d00000000022304361f011afa0000040f00000013020000290000000001210049000007c10010009c000007c1010080410000006001100210000007c10020009c000007c1020080410000004002200210000000000121019f00001f020001042e000000240030008c00000a830000413d0000000002000416000000000002004b00000a830000c13d0000000401100370000000000601043b000007c40060009c00000a830000213d0000000101000039000000000101041a000007c4011001970000000005000411000000000015004b00000a190000c13d000000000056004b00000a5d0000c13d0000082001000041000000800010043f000008210100004100001f03000104300000000001000416000000000001004b00000a830000c13d0000000901000039000005cf0000013d000000240030008c00000a830000413d0000000002000416000000000002004b00000a830000c13d0000000401100370000000000101043b000007c50010009c00000a830000213d1f011e410000040f000000000001004b0000000001000039000000010100c039000000400200043d0000000000120435000007c10020009c000007c10200804100000040012002100000085e011001c700001f020001042e000000240030008c00000a830000413d0000000002000416000000000002004b00000a830000c13d0000000402100370000000000202043b000007c50020009c00000a830000213d0000002304200039000000000034004b00000a830000813d0000000404200039000000000141034f000000000101043b000f00000001001d000007c50010009c00000a830000213d00000024092000390000000f0100002900000007011002100000000001910019000000000031004b00000a830000213d0000000101000039000000000101041a000007c4011001970000000002000411000000000012004b00000a190000c13d0000000f0000006b00000aad0000c13d0000002001000039000000800010043f000000c0020000390000008001000039000000a00000043f0000000002120049000007c10020009c000007c1020080410000006002200210000007c10010009c000007c1010080410000004001100210000000000112019f0000000002000414000007c10020009c000007c102008041000000c002200210000000000121019f0000081e011001c70000800d020000390000000103000039000008790400004100000a7a0000013d000000240030008c00000a830000413d0000000002000416000000000002004b00000a830000c13d0000000401100370000000000101043b001300000001001d000007c40010009c00000a830000213d0000000001000412002100000001001d002000000000003d000080050100003900000044030000390000000004000415000000210440008a00000005044002100000081c020000411f011ed90000040f000007c401100197000000130010006b00000000010000390000000101006039000000800010043f0000081d0100004100001f020001042e0000000001000416000000000001004b00000a830000c13d000000c001000039000000400010043f0000001301000039000000800010043f0000086f01000041000000a00010043f0000002001000039000000c00010043f0000008001000039000000e0020000391f011afa0000040f000000c00110008a000007c10010009c000007c101008041000000600110021000000870011001c700001f020001042e000000240030008c00000a830000413d0000000002000416000000000002004b00000a830000c13d0000000402100370000000000202043b001300000002001d000007c50020009c00000a830000213d000000130230006a000007df0020009c00000a830000213d000001040020008c00000a830000413d000000a002000039000000400020043f0000001302000029001200840020003d0000001201100360000000800000043f000000000101043b001100000001001d000007c40010009c00000a830000213d0000081c0100004100000000001004430000000001000412000000040010044300000024000004430000000001000414000007c10010009c000007c101008041000000c00110021000000847011001c700008005020000391f011efc0000040f000000010020019000001ae30000613d000000000101043b000000110210014f0000000201000367000007c40020019800000a7f0000c13d000000120200002900110060002000920000001101100360000000000101043b000007c50010009c00000a830000213d000000400300043d000008490200004100000000002304350000008001100210001200000003001d000000040230003900000000001204350000081c01000041000000000010044300000000010004120000000400100443000000400100003900000024001004430000000001000414000007c10010009c000007c101008041000000c00110021000000847011001c700008005020000391f011efc0000040f000000010020019000001ae30000613d000000000201043b0000000001000414000007c402200197000000040020008c00000d3c0000c13d0000000103000031000000200030008c0000002004000039000000000403401900000d660000013d0000000001000416000000000001004b00000a830000c13d0000000401000039000000000101041a000007c401100197000000800010043f0000081d0100004100001f020001042e000000440030008c00000a830000413d0000000002000416000000000002004b00000a830000c13d0000000402100370000000000202043b000007c50020009c00000a830000213d0000002304200039000000000034004b00000a830000813d0000000404200039000000000441034f000000000404043b000800000004001d000007c50040009c00000a830000213d000700240020003d000000080200002900000005022002100000000702200029000000000032004b00000a830000213d0000002402100370000000000202043b000200000002001d000007c50020009c00000a830000213d00000002020000290000002302200039000000000032004b00000a830000813d00000002020000290000000402200039000000000121034f000000000101043b000100000001001d000007c50010009c00000a830000213d0000000201000029000500240010003d000000010100002900000005011002100000000501100029000000000031004b00000a830000213d0000000101000039000000000101041a000007c4011001970000000002000411000000000012004b00000a190000c13d000000080000006b00000d9c0000c13d000008250100004100000000001004430000000001000414000007c10010009c000007c101008041000000c00110021000000826011001c70000800b020000391f011efc0000040f000000010020019000001ae30000613d000000000101043b000000010000006b00000a7d0000613d000007c101100197000400000001001d0003008000100218000700000000001d0000000701000029000000050110021000000005021000290000000201000367000000000221034f000000000202043b0000000003000031000000020430006a000001430440008a000007c605400197000007c606200197000000000756013f000000000056004b0000000005000019000007c605004041000000000042004b0000000004000019000007c604008041000007c60070009c000000000504c019000000000005004b00000a830000c13d001200050020002d000000120230006a001100000002001d000007df0020009c00000a830000213d0000001102000029000001200020008c00000a830000413d000000400200043d000b00000002001d000008270020009c0000006d0000213d0000000b02000029000000a002200039000000400020043f0000001202100360000000000202043b000007c50020009c00000a830000213d0000000b040000290000000002240436000a00000002001d00000012020000290000002002200039000000000221034f000000000202043b000007c50020009c00000a830000213d0000001202200029001300000002001d0000001f02200039000000000032004b0000000004000019000007c604008041000007c602200197000007c605300197000000000752013f000000000052004b0000000002000019000007c602004041000007c60070009c000000000204c019000000000002004b00000a830000c13d0000001302100360000000000202043b000007c50020009c0000006d0000213d00000005092002100000003f04900039000007c704400197000000400600043d0000000004460019001000000006001d000000000064004b00000000070000390000000107004039000007c50040009c0000006d0000213d00000001007001900000006d0000c13d000000400040043f00000010040000290000000000240435000000130200002900000020082000390000000009890019000000000039004b00000a830000213d000000000098004b000006c60000813d000000100a000029000006830000013d000000200aa000390000000002b7001900000000000204350000000000ca04350000002008800039000000000098004b000006c60000813d000000000281034f000000000202043b000007c50020009c00000a830000213d000000130d2000290000003f02d00039000000000032004b0000000004000019000007c604008041000007c602200197000000000752013f000000000052004b0000000002000019000007c602004041000007c60070009c000000000204c019000000000002004b00000a830000c13d000000200ed000390000000002e1034f000000000b02043b000007c500b0009c0000006d0000213d0000001f02b000390000087a022001970000003f022000390000087a02200197000000400c00043d00000000022c00190000000000c2004b00000000040000390000000104004039000007c50020009c0000006d0000213d00000001004001900000006d0000c13d0000004004d00039000000400020043f0000000007bc043600000000024b0019000000000032004b00000a830000213d0000002002e00039000000000421034f0000087a02b00198000000000e270019000006b80000613d000000000f04034f000000000d07001900000000f60f043c000000000d6d04360000000000ed004b000006b40000c13d0000001f0db001900000067c0000613d000000000224034f0000000304d0021000000000060e043300000000064601cf000000000646022f000000000202043b0000010004400089000000000242022f00000000024201cf000000000262019f00000000002e04350000067c0000013d0000000a020000290000001004000029000000000042043500000012020000290000004006200039000000000261034f000000000202043b000007c50020009c00000a830000213d00000012072000290000001f02700039000000000032004b0000000004000019000007c604008041000007c602200197000000000852013f000000000052004b0000000002000019000007c602004041000007c60080009c000000000204c019000000000002004b00000a830000c13d000000000271034f000000000402043b000007c50040009c0000006d0000213d0000001f024000390000087a022001970000003f022000390000087a02200197000000400500043d0000000002250019000000000052004b00000000080000390000000108004039000007c50020009c0000006d0000213d00000001008001900000006d0000c13d0000002008700039000000400020043f00000000074504360000000002840019000000000032004b00000a830000213d000000000881034f0000087a024001980000000003270019000006fe0000613d000000000908034f000000000a070019000000009b09043c000000000aba043600000000003a004b000006fa0000c13d0000001f094001900000070b0000613d000000000228034f0000000308900210000000000903043300000000098901cf000000000989022f000000000202043b0000010008800089000000000282022f00000000028201cf000000000292019f0000000000230435000000000247001900000000000204350000000b020000290000004002200039000900000002001d00000000005204350000001102000029000000600220008a000007df0020009c00000a830000213d000000600020008c00000a830000413d000000400300043d000008280030009c0000006d0000213d0000006002300039000000400020043f0000002002600039000000000421034f000000000404043b000000000004004b0000000005000039000000010500c039000000000054004b00000a830000c13d00000000044304360000002002200039000000000521034f000000000505043b000008290050009c00000a830000213d00000000005404350000002004200039000000000241034f000000000202043b000008290020009c00000a830000213d000000400530003900000000002504350000000b020000290000006002200039000800000002001d00000000003204350000001102000029000000c00220008a000007df0020009c00000a830000213d000000600020008c00000a830000413d000000400200043d000008280020009c0000006d0000213d0000006003200039000000400030043f0000002003400039000000000431034f000000000404043b000000000004004b0000000005000039000000010500c039000000000054004b00000a830000c13d00000000044204360000002003300039000000000531034f000000000505043b000008290050009c00000a830000213d00000000005404350000002003300039000000000131034f000000000301043b000008290030009c00000a830000213d000000400120003900000000003104350000000b030000290000008003300039000600000003001d0000000000230435000000080300002900000000030304330000004005300039000000000505043300000829065001970000000057030434000000000007004b0000076a0000613d000000000006004b0000147b0000613d00000000050504330000082905500197000000000056004b0000076f0000413d0000147b0000013d000000000006004b000014670000c13d00000000050504330000082900500198000014670000c13d000000000101043300000829011001970000000003020433000000000003004b0000077b0000613d000000000001004b000014820000613d00000000030404330000082903300197000000000031004b000007800000413d000014820000013d000000000001004b0000146b0000c13d000000000104043300000829001001980000146b0000c13d000000090100002900000000010104330000000001010433000000000001004b00000b350000613d0000000b010000290000000001010433000007c501100197001300000001001d000000000010043f0000000601000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000101043b000000000101041a000000000001004b000013be0000c13d0000000501000039000000000101041a000007c50010009c0000006d0000213d00000001021000390000000503000039000000000023041b0000082d0110009a0000001302000029000000000021041b000000000103041a001200000001001d000000000020043f0000000601000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000101043b0000001202000029000000000021041b0000000b010000290000000001010433000007c501100197000000000010043f0000000701000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000101043b000000400200043d000008270020009c0000006d0000213d0000000803000029000000000303043300000020043000390000000004040433000000000503043300000040033000390000000003030433000000a006200039000000400060043f0000008006200039000008290730019700000000007604350000004006200039000000000005004b0000000005000039000000010500c0390000000000560435000000200520003900000004060000290000000000650435000008290440019700000060052000390000000000450435000000000042043500000000020000190000082e0200c041000000000501041a0000082f05500197000000000252019f00000003022001af000000000242019f000000000021041b0000008002300210000000000242019f0000000103100039000000000023041b000000400200043d000008270020009c0000006d0000213d0000000603000029000000000303043300000020043000390000000004040433000000000503043300000040033000390000000003030433000000a006200039000000400060043f0000008006200039000008290730019700000000007604350000004006200039000000000005004b0000000005000039000000010500c0390000000000560435000000200520003900000004060000290000000000650435000008290440019700000060052000390000000000450435000000000042043500000000020000190000082e0200c0410000000205100039000000000605041a0000082f06600197000000000262019f00000003022001af000000000242019f000000000025041b0000008002300210000000000242019f0000000303100039000000000023041b00000009020000290000000002020433001100000002001d0000000032020434001000000003001d001300000002001d000007c50020009c0000006d0000213d0000000403100039000000000103041a000000010010019000000001051002700000007f0550618f0000001f0050008c00000000020000390000000102002039000000000121013f0000000100100190000019940000c13d000000200050008c0000001304000029001200000003001d000008470000413d000f00000005001d000000000030043f0000000001000414000007c10010009c000007c101008041000000c001100210000007d0011001c700008010020000391f011efc0000040f000000010020019000000a830000613d00000013040000290000001f024000390000000502200270000000200040008c0000000002004019000000000301043b0000000f010000290000001f01100039000000050110027000000000011300190000000002230019000000000012004b0000001203000029000008470000813d000000000002041b0000000102200039000000000012004b000008430000413d000000200040008c000008650000413d000000000030043f0000000001000414000007c10010009c000007c101008041000000c001100210000007d0011001c700008010020000391f011efc0000040f000000010020019000000a830000613d00000013060000290000087a02600198000000000101043b0000001107000029000008720000613d000000010320008a000000050330027000000000033100190000000104300039000000200300003900000000057300190000000005050433000000000051041b00000020033000390000000101100039000000000041004b0000085d0000c13d000008730000013d000000000004004b000008700000613d00000003014002100000087b0110027f0000087b0110016700000010020000290000000002020433000000000112016f0000000102400210000000000121019f000008800000013d0000000001000019000008800000013d0000002003000039000000000062004b0000087d0000813d0000000302600210000000f80220018f0000087b0220027f0000087b0220016700000000037300190000000003030433000000000223016f000000000021041b000000010160021000000001011001bf0000001203000029000000000013041b0000000a0100002900000000010104330000000002010433000000000002004b000009a20000613d0000000003000019000f00000003001d0000000502300210000000000121001900000020011000390000000001010433001200000001001d0000000031010434000000000001004b00000b350000613d0000000b020000290000000002020433001300000002001d000007c10010009c000007c1010080410000006001100210000007c10030009c000c00000003001d000007c10200004100000000020340190000004002200210000000000121019f0000000002000414000007c10020009c000007c102008041000000c002200210000000000112019f0000081e011001c700008010020000391f011efc0000040f000000010020019000000a830000613d0000001302000029000007c502200197000000000101043b001300000001001d000e00000002001d000000000020043f0000000701000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000201043b0000001301000029000000000010043f001100000002001d0000000601200039001000000001001d000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000101043b000000000101041a000000000001004b000011f80000c13d00000011010000290000000502100039000000000102041a000007c50010009c0000006d0000213d000d00000001001d0000000101100039000000000012041b001100000002001d000000000020043f0000000001000414000007c10010009c000007c101008041000000c001100210000007d0011001c700008010020000391f011efc0000040f0000000100200190000000130200002900000a830000613d000000000101043b0000000d01100029000000000021041b0000001101000029000000000101041a001100000001001d000000000020043f0000001001000029000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f0000000100200190000000130200002900000a830000613d000000000101043b0000001103000029000000000031041b000000000020043f0000000801000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000301043b00000012010000290000000004010433000007c50040009c0000006d0000213d000000000103041a000000010010019000000001051002700000007f0550618f0000001f0050008c00000000020000390000000102002039000000000121013f0000000100100190000019940000c13d000000200050008c001300000003001d001100000004001d000009310000413d001000000005001d000000000030043f0000000001000414000007c10010009c000007c101008041000000c001100210000007d0011001c700008010020000391f011efc0000040f000000010020019000000a830000613d00000011040000290000001f024000390000000502200270000000200040008c0000000002004019000000000301043b00000010010000290000001f01100039000000050110027000000000011300190000000002230019000000000012004b0000001303000029000009310000813d000000000002041b0000000102200039000000000012004b0000092d0000413d000000200040008c0000095d0000413d000000000030043f0000000001000414000007c10010009c000007c101008041000000c001100210000007d0011001c700008010020000391f011efc0000040f000000010020019000000a830000613d00000011080000290000087a02800198000000000101043b0000099c0000613d000000010320008a00000005033002700000000003310019000000010430003900000020030000390000000c07000029000000120600002900000000056300190000000005050433000000000051041b00000020033000390000000101100039000000000041004b000009480000c13d000000000082004b000009590000813d0000000302800210000000f80220018f0000087b0220027f0000087b0220016700000000036300190000000003030433000000000223016f000000000021041b000000010180021000000001011001bf00000013030000290000096a0000013d000000000004004b0000000c07000029000009680000613d00000003014002100000087b0110027f0000087b011001670000000002070433000000000112016f0000000102400210000000000121019f000009690000013d00000000010000190000001206000029000000000013041b000000400100043d00000020020000390000000003210436000000000206043300000000002304350000004003100039000000000002004b0000097b0000613d000000000400001900000000053400190000000006740019000000000606043300000000006504350000002004400039000000000024004b000009740000413d0000001f042000390000087a04400197000000000232001900000000000204350000004002400039000007c10020009c000007c1020080410000006002200210000007c10010009c000007c1010080410000004001100210000000000112019f0000000002000414000007c10020009c000007c102008041000000c002200210000000000112019f0000081e011001c70000800d02000039000000020300003900000831040000410000000e050000291f011ef70000040f000000010020019000000a830000613d0000000f0300002900000001033000390000000a0100002900000000010104330000000002010433000000000023004b000008870000413d000009a20000013d00000020030000390000000c070000290000001206000029000000000082004b000009510000413d000009590000013d0000000601000029000000000201043300000008010000290000000005010433000000090100002900000000030104330000000b010000290000000004010433000000400100043d000000200610003900000100070000390000000000760435000007c50440019700000000004104350000010007100039000000006403043400000000004704350000012003100039000000000004004b000009be0000613d000000000700001900000000083700190000000009760019000000000909043300000000009804350000002007700039000000000047004b000009b70000413d000000000643001900000000000604350000000076050434000000000006004b0000000006000039000000010600c039000000400810003900000000006804350000000006070433000008290660019700000060071000390000000000670435000000400550003900000000050504330000082905500197000000800610003900000000005604350000000065020434000000000005004b0000000005000039000000010500c039000000a007100039000000000057043500000000050604330000082905500197000000c0061000390000000000560435000000400220003900000000020204330000082902200197000000e00510003900000000002504350000001f024000390000087a0220019700000000021200490000000002320019000007c10020009c000007c1020080410000006002200210000007c10010009c000007c1010080410000004001100210000000000112019f0000000002000414000007c10020009c000007c102008041000000c002200210000000000112019f0000081e011001c70000800d02000039000000010300003900000832040000411f011ef70000040f000000010020019000000a830000613d00000007020000290000000102200039000700000002001d000000010020006c0000061d0000413d00000a7d0000013d000000240030008c00000a830000413d0000000002000416000000000002004b00000a830000c13d0000000401100370000000000101043b000007c40010009c00000a830000213d0000000102000039000000000202041a000007c4022001970000000003000411000000000023004b00000a190000c13d0000000902000039000000000302041a000007c803300197000000000313019f000000000032041b000000800010043f0000000001000414000007c10010009c000007c101008041000000c0011002100000085f011001c70000800d020000390000000103000039000008600400004100000a7a0000013d0000087501000041000000800010043f000008210100004100001f03000104300000083b04000041000000a00500003900000000060000190000000007050019000000000504041a000000000557043600000001044000390000000106600039000000000026004b00000a200000413d000000410270008a0000087a042001970000083c0040009c0000006d0000213d0000008002400039000000800500043d000000400020043f000007c50050009c0000006d0000213d000002620000013d000000a005000039000008460300004100000000040000190000000006050019000000000503041a000000000556043600000001033000390000000104400039000000000014004b00000a340000413d000000410160008a0000087a041001970000083c0040009c0000006d0000213d0000008001400039000000400010043f0000000000210435000000a002400039000000800300043d0000000000320435000000c002400039000000000003004b00000a500000613d000000a00400003900000000050000190000000046040434000007c40660019700000000026204360000000105500039000000000035004b00000a4a0000413d0000000002120049000007c10020009c000007c1020080410000006002200210000007c10010009c000007c1010080410000004001100210000000000112019f00001f020001042e0000086101000041000000800010043f000008210100004100001f0300010430000000000100041a000007c801100197000000000161019f000000000010041b0000000001000414000007c10010009c000007c101008041000000c0011002100000081e011001c70000800d0200003900000003030000390000081f0400004100000a7a0000013d0000000402000039000000000302041a000007c804300197000000000414019f000000000042041b000007c402300197000000800020043f000000a00010043f0000000001000414000007c10010009c000007c101008041000000c0011002100000083e011001c70000800d0200003900000001030000390000083f040000411f011ef70000040f000000010020019000000a830000613d000000000100001900001f020001042e0000001201100360000000000101043b000007c40010009c00000aa30000a13d000000000100001900001f03000104300000000105000039000000000405041a000007c804400197000000000224019f000000000025041b000000000003004b00000b350000613d000007c40210019800000b350000613d000000130000006b00000b350000613d000d00000005001d000000800010043f000000c00030043f000007c901000041000000400300043d001200000003001d00000000001304350000000001000414000000040020008c00000b3d0000c13d0000000001000415000000250110008a00000005011002100000000103000031000000200030008c00000020040000390000000004034019002500000000003d00000b6b0000013d000000400200043d0000084803000041000000000032043500000004032000390000000000130435000007c10020009c000007c1020080410000004001200210000007d5011001c700001f0300010430000000000a000019000e00000009001d0000000701a0021000000000039100190000000001300079000007df0010009c00000a830000213d000000800010008c00000a830000413d000000400100043d0000083c0010009c0000006d0000213d0000008002100039000000400020043f0000000206000367000000000236034f000000000202043b00000000042104360000002005300039000000000356034f000000000303043b000007c10030009c00000a830000213d00000000003404350000002007500039000000000576034f000000000805043b000007c50080009c00000a830000213d000000400510003900000000008504350000002007700039000000000676034f000000000606043b000000000006004b0000000007000039000000010700c039000000000076004b00000a830000c13d00000060071000390000000000670435000000400b00043d000000000002004b00000f2d0000613d000000000008004b00000f2d0000613d0000082800b0009c0000006d0000213d0000006001b00039000000400010043f0000004001b00039001200000001001d000000000061043500000000012b0436001100000001001d00000000003104350000000001050433000007c501100197000000000010043f0000000a01000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c7000080100200003900130000000a001d00100000000b001d1f011efc0000040f000000130a0000290000000e09000029000000010020019000000a830000613d00000010020000290000000002020433000000000101043b000000000021041b00000011020000290000000002020433000007c1022001970000000101100039000000000301041a0000087603300197000000000223019f00000012030000290000000003030433000000000003004b00000877030000410000000003006019000000000232019f000000000021041b000000010aa000390000000f00a0006c00000aaf0000413d000000400100043d00000020021000390000000f03000029000000000032043500000020020000390000000000210435000000400210003900000002030003670000000004000019000000000593034f000000000505043b00000000065204360000002005900039000000000753034f000000000707043b000007c10070009c00000a830000213d00000000007604350000002005500039000000000653034f000000000606043b000007c50060009c00000a830000213d000000400720003900000000006704350000002005500039000000000553034f000000000505043b000000000005004b0000000006000039000000010600c039000000000065004b00000a830000c13d000000600620003900000000005604350000008009900039000000800220003900000001044000390000000f0040006c00000b150000413d000005380000013d000000400100043d00000833020000410000000000210435000007c10010009c000007c1010080410000004001100210000007ca011001c700001f03000104300000001203000029000007c10030009c000007c1030080410000004003300210000007c10010009c000007c101008041000000c001100210000000000131019f000007ca011001c71f011efc0000040f0000006003100270000007c103300197000000200030008c000000200400003900000000040340190000001f0640018f0000002007400190000000120570002900000b560000613d000000000801034f0000001209000029000000008a08043c0000000009a90436000000000059004b00000b520000c13d000000000006004b00000b630000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f00030000000103550000000001000415000000240110008a0000000501100210002400000000003d000000010020019000000b800000613d0000001f02400039000000600420018f0000001202400029000000000042004b00000000040000390000000104004039000007c50020009c0000006d0000213d00000001004001900000006d0000c13d000000400020043f000000200030008c00000a830000413d00000012030000290000000003030433000000ff0030008c00000a830000213d0000000501100270000000000103001f000000060030008c00000f620000c13d0000000601000039000000a00010043f0000000401000039000000000201041a000007c80220019700000013022001af000000000021041b00000011010000290000000001010433000000000001004b0000000001000039000000010100c039000000e00010043f00000f720000613d000000400100043d000007cd0010009c0000006d0000213d0000002002100039000000400020043f0000000000010435000000e00100043d000000000001004b00000f6e0000c13d000000400100043d000008650200004100000b370000013d0000000002000019001300000002001d0000000502200210001200000002001d00000010012000290000000001010433000000000010043f0000000801000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000101043b000000000201041a000000010320019000000001052002700000007f0550618f0000001f0050008c00000000040000390000000104002039000000000442013f0000000100400190000019940000c13d000000400700043d0000000004570436000000000003004b00000bd90000613d000b00000004001d000c00000005001d000d00000007001d000000000010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007d0011001c700008010020000391f011efc0000040f000000010020019000000a830000613d0000000c08000029000000000008004b000000200500008a00000be10000613d000000000201043b00000000010000190000000f060000290000000d070000290000000b090000290000000003190019000000000402041a000000000043043500000001022000390000002001100039000000000081004b00000bd10000413d00000be40000013d0000087c012001970000000000140435000000000005004b00000020010000390000000001006039000000200500008a0000000f0600002900000be40000013d00000000010000190000000f060000290000000d070000290000003f01100039000000000251016f0000000001720019000000000021004b00000000020000390000000102004039000007c50010009c0000006d0000213d00000001002001900000006d0000c13d000000400010043f00000000010604330000001302000029000000000021004b0000127b0000a13d00000012030000290000000e0130002900000000007104350000000001060433000000000021004b0000127b0000a13d000000010220003900000011010000290000000001010433000000000012004b00000b9b0000413d000002f40000013d0000001301000029000000000010043f0000000701000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d00000012020000290000001f022000390000087a02200197000e00000002001d0000003f022000390000087a02200197000000000101043b000f00000001001d000000400100043d0000000002210019000000000012004b00000000040000390000000104004039000007c50020009c0000006d0000213d00000001004001900000006d0000c13d000000400020043f000000120200002900000000022104360000001005000029000000000050007c00000a830000213d00000012040000290000087a054001980010001f00400193000d00000005001d000000000352001900000011040000290000002004400039001100000004001d000000020440036700000c330000613d000000000504034f0000000006020019000000005705043c0000000006760436000000000036004b00000c2f0000c13d000000100000006b00000c410000613d0000000d0440036000000010050000290000000305500210000000000603043300000000065601cf000000000656022f000000000404043b0000010005500089000000000454022f00000000045401cf000000000464019f000000000043043500000012032000290000000000030435000007c10020009c000007c10200804100000040022002100000000001010433000007c10010009c000007c1010080410000006001100210000000000121019f0000000002000414000007c10020009c000007c102008041000000c002200210000000000112019f0000081e011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000101043b000b00000001001d000000000010043f0000000f010000290000000601100039000c00000001001d000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000101043b000000000101041a000a00000001001d000000000001004b000011490000c13d000000400100043d0000004402100039000000120300002900000000003204350000002402100039000000400300003900000000003204350000084302000041000000000021043500000004021000390000001303000029000000000032043500000064021000390000000d03200029000000110400002900000002044003670000000d0000006b00000c820000613d000000000504034f0000000006020019000000005705043c0000000006760436000000000036004b00000c7e0000c13d000000100000006b00000c900000613d0000000d0440036000000010050000290000000305500210000000000603043300000000065601cf000000000656022f000000000404043b0000010005500089000000000454022f00000000045401cf000000000464019f000000000043043500000012022000290000000000020435000007c10010009c000007c10100804100000040011002100000000e02000029000008440020009c00000844020080410000006002200210000000000112019f000008450110009a00001f030001043000000000040000190000000f0c00002900000ca70000013d0000001f076000390000087a077001970000000006650019000000000006043500000000057500190000000104400039000000000024004b000002ff0000813d0000000006150049000000400660008a0000000003630436000000200cc0003900000000060c043300000000760604340000000005650436000000000006004b00000c9f0000613d00000000080000190000000009580019000000000a870019000000000a0a04330000000000a904350000002008800039000000000068004b00000cb10000413d00000c9f0000013d0000001203000029000007c10030009c000007c1030080410000004003300210000007c10010009c000007c101008041000000c001100210000000000131019f000007d5011001c71f011efc0000040f0000006003100270000007c103300197000000200030008c000000200400003900000000040340190000001f0640018f0000002007400190000000120570002900000cd20000613d000000000801034f0000001209000029000000008a08043c0000000009a90436000000000059004b00000cce0000c13d000000000006004b00000cdf0000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f0003000000010355000000010020019000000f4a0000613d0000001f01400039000000600210018f0000001201200029000000000021004b00000000020000390000000102004039000007c50010009c0000006d0000213d00000001002001900000006d0000c13d000000400010043f000000200030008c00000a830000413d00000012020000290000000002020433000000000002004b0000000003000039000000010300c039000000000032004b00000a830000c13d000000000002004b000010250000c13d000000110100002900000020011000390000000201100367000000000101043b001200000001001d000007c40010009c00000a830000213d0000081c01000041000000000010044300000000010004120000000400100443000000600100003900000024001004430000000001000414000007c10010009c000007c101008041000000c00110021000000847011001c700008005020000391f011efc0000040f000000010020019000001ae30000613d000000000101043b000000000001004b0000120f0000c13d00000011010000290000000201100367000000000101043b001200000001001d000007c50010009c00000a830000213d0000001201000029000000000010043f0000000601000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000400200043d001000000002001d0000000402200039000000000101043b000000000101041a000000000001004b00000d960000613d0000000401000039000000000301041a0000084c0100004100000010040000290000000000140435000000120100002900000000001204350000000001000414000007c402300197000000040020008c000013220000c13d0000000103000031000000200030008c000000200400003900000000040340190000134c0000013d0000001203000029000007c10030009c000007c1030080410000004003300210000007c10010009c000007c101008041000000c001100210000000000131019f000007d5011001c71f011efc0000040f0000006003100270000007c103300197000000200030008c000000200400003900000000040340190000001f0640018f0000002007400190000000120570002900000d550000613d000000000801034f0000001209000029000000008a08043c0000000009a90436000000000059004b00000d510000c13d000000000006004b00000d620000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f0003000000010355000000010020019000000f560000613d0000001f01400039000000600210018f0000001201200029000000000021004b00000000020000390000000102004039000007c50010009c0000006d0000213d00000001002001900000006d0000c13d000000400010043f000000200030008c00000a830000413d00000012020000290000000002020433000000000002004b0000000003000039000000010300c039000000000032004b00000a830000c13d000000000002004b000010250000c13d00000011010000290000000201100367000000000101043b001200000001001d000007c50010009c00000a830000213d0000001201000029000000000010043f0000000601000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000400200043d001000000002001d0000000402200039000000000101043b000000000101041a000000000001004b000012260000c13d0000085d010000410000001003000029000000000013043500000012010000290000000000120435000012ce0000013d0000000002000019000900000002001d000000050120021000000007011000290000000201100367000000000101043b000d00000001001d000007c50010009c00000a830000213d0000000d01000029000000000010043f0000000601000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000101043b000000000301041a000000000003004b000012810000613d0000000501000039000000000201041a000000000002004b000011500000613d000000010130008a000000000023004b00000dd80000613d000000000012004b0000127b0000a13d000008220130009a000008220220009a000000000202041a000000000021041b000000000020043f0000000601000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c70000801002000039001300000003001d1f011efc0000040f000000010020019000000a830000613d000000000101043b0000001302000029000000000021041b0000000501000039000000000301041a000000000003004b000012750000613d000000010130008a000008220230009a000000000002041b0000000502000039000000000012041b0000000d01000029000000000010043f0000000601000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000101043b000000000001041b0000000d01000029000000000010043f0000000701000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000101043b0000000501100039000000000301041a000000400200043d001100000002001d001300000003001d0000000002320436000c00000002001d000000000010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007d0011001c700008010020000391f011efc0000040f000000010020019000000a830000613d0000001305000029000000000005004b0000000c0200002900000e170000613d000000000101043b0000000c020000290000000003000019000000000401041a000000000242043600000001011000390000000103300039000000000053004b00000e110000413d000000110120006a0000001f011000390000087a021001970000001101200029000000000021004b00000000020000390000000102004039000007c50010009c0000006d0000213d00000001002001900000006d0000c13d000000400010043f00000011010000290000000001010433000000000001004b00000ebc0000613d000000000200001900000e2f0000013d0000001302000029000000010220003900000011010000290000000001010433000000000012004b00000ebc0000813d001300000002001d0000000d01000029000000000010043f0000000701000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000301043b000000110100002900000000010104330000001302000029000000000021004b0000127b0000a13d00000005012002100000000c011000290000000001010433000e00000001001d000000000010043f000f00000003001d0000000601300039001000000001001d000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000101043b000000000101041a000000000001004b00000e290000613d0000000f020000290000000503200039000000000203041a000000000002004b000011500000613d000000000021004b001200000001001d000f00000003001d00000e9b0000613d000b00000002001d000000000030043f0000000001000414000007c10010009c000007c101008041000000c001100210000007d0011001c700008010020000391f011efc0000040f000000010020019000000a830000613d0000001202000029000a000100200092000000000101043b0000000f04000029000000000204041a0000000a0020006c0000127b0000a13d0000000b02000029000000010220008a0000000001120019000000000101041a000b00000001001d000000000040043f0000000001000414000007c10010009c000007c101008041000000c001100210000007d0011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000101043b0000000a011000290000000b02000029000000000021041b000000000020043f0000001001000029000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000101043b0000001202000029000000000021041b0000000f03000029000000000103041a001200000001001d000000000001004b000012750000613d000000000030043f0000000001000414000007c10010009c000007c101008041000000c001100210000007d0011001c700008010020000391f011efc0000040f000000010020019000000a830000613d0000001202000029000000010220008a000000000101043b0000000001210019000000000001041b0000000f01000029000000000021041b0000000e01000029000000000010043f0000001001000029000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000101043b000000000001041b00000e290000013d0000000d01000029000000000010043f0000000701000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000301043b000000000003041b0000000101300039000000000001041b0000000201300039000000000001041b0000000301300039000000000001041b0000000404300039000000000104041a000000010010019000000001051002700000007f0550618f0000001f0050008c00000000020000390000000102002039000000000121013f0000000100100190000019940000c13d000000000005004b00000efe0000613d0000001f0050008c00000efd0000a13d001100000005001d001300000003001d001200000004001d000000000040043f0000000001000414000007c10010009c000007c101008041000000c001100210000007d0011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000101043b00000011020000290000001f02200039000000050220027000000000022100190000000103100039000000000023004b00000ef90000813d000000000003041b0000000103300039000000000023004b00000ef50000413d0000001202000029000000000002041b00000000040100190000001303000029000000000004041b0000000501300039000000000201041a000000000001041b000000000002004b00000f160000613d001300000002001d000000000010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007d0011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000101043b0000001302100029000000000021004b00000f160000813d000000000001041b0000000101100039000000000021004b00000f120000413d000000400100043d0000000d020000290000000000210435000007c10010009c000007c10100804100000040011002100000000002000414000007c10020009c000007c102008041000000c002200210000000000112019f000007d0011001c70000800d02000039000000010300003900000823040000411f011ef70000040f000000010020019000000a830000613d00000009020000290000000102200039000000080020006c00000d9d0000413d0000060b0000013d000008780200004100000000002b043500000000010104330000000402b0003900000000001204350000000001040433000007c1011001970000002402b0003900000000001204350000000001050433000007c5011001970000004402b0003900000000001204350000000001070433000000000001004b0000000001000039000000010100c0390000006402b000390000000000120435000007c100b0009c000007c10b0080410000004001b00210000007e3011001c700001f03000104300000083602000041000001400020043f000001440010043f000008370100004100001f03000104300000001f0530018f000007c306300198000000400200043d0000000004620019000010320000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b00000f510000c13d000010320000013d0000001f0530018f000007c306300198000000400200043d0000000004620019000010320000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b00000f5d0000c13d000010320000013d00000024012000390000000000310435000007cb010000410000000000120435000000040120003900000006030000390000000000310435000007c10020009c000007c1020080410000004001200210000007cc011001c700001f030001043000000011010000290000000001010433000000000001004b000010450000c13d0000000e01000029001207c40010019c000000400100043d001300000001001d00000f7f0000c13d000007e80100004100000013020000290000000000120435000007c10020009c000007c1020080410000004001200210000007ca011001c700001f0300010430000007d2010000410000001302000029000000000012043500000000010004140000001202000029000000040020008c00000f8b0000c13d0000000103000031000000200030008c0000002004000039000000000403401900000fb60000013d0000001302000029000007c10020009c000007c1020080410000004002200210000007c10010009c000007c101008041000000c001100210000000000121019f000007ca011001c700000012020000291f011efc0000040f0000006003100270000007c103300197000000200030008c000000200400003900000000040340190000001f0640018f0000002007400190000000130570002900000fa50000613d000000000801034f0000001309000029000000008a08043c0000000009a90436000000000059004b00000fa10000c13d000000000006004b00000fb20000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f00030000000103550000000100200190000010270000613d0000001f01400039000000600110018f0000001305100029000000000015004b00000000020000390000000102004039001100000005001d000007c50050009c0000006d0000213d00000001002001900000006d0000c13d0000001102000029000000400020043f000000200040008c00000a830000413d00000013020000290000000002020433001300000002001d000007c40020009c00000a830000213d000007d3020000410000001104000029000000000024043500000000020004140000001304000029000000040040008c000011560000c13d0000001102100029001000000002001d000007c50020009c0000006d0000213d0000001002000029000000400020043f00000011020000290000000002020433000007c10020009c00000a830000213d000000000002004b000012c90000c13d000007d6020000410000001004000029000000000024043500000000020004140000001204000029000000040040008c000012df0000c13d0000001002100029001100000002001d000007c50020009c0000006d0000213d0000001102000029000000400020043f00000010020000290000000002020433000007c10020009c00000a830000213d000000000002004b000014880000c13d0000000e02000029000001000020043f0000001304000029000001200040043f000007d802000041000000110500002900000000002504350000000002000414000000040040008c0000148b0000c13d0000001102100029001300000002001d000007c50020009c0000006d0000213d0000001302000029000000400020043f000014c10000013d000000000001004b000010910000613d000000a00200043d0000082902200197000000000021004b000010910000813d0000001301000029000000000010043f0000000701000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000101043b00000080020000391f011e550000040f000001200100043d0000082901100197000000e00200043d000000000002004b000010a20000c13d000000000001004b000010210000c13d000001000100043d0000082900100198000010a80000613d000000400200043d001300000002001d0000082b010000410000118f0000013d0000084a0200004100000b370000013d0000001f0530018f000007c306300198000000400200043d0000000004620019000010320000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b0000102e0000c13d000000000005004b0000103f0000613d000000000161034f0000000305500210000000000604043300000000065601cf000000000656022f000000000101043b0000010005500089000000000151022f00000000015101cf000000000161019f00000000001404350000006001300210000007c10020009c000007c1020080410000004002200210000000000112019f00001f030001043000000000020000190000104d0000013d0000001302000029000000010220003900000011010000290000000001010433000000000012004b00000f720000813d001300000002001d000000050120021000000010011000290000000001010433000007c403100198000010470000613d000000000030043f0000000301000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c70000801002000039001200000003001d1f011efc0000040f0000001204000029000000010020019000000a830000613d000000000101043b000000000101041a000000000001004b000010470000c13d0000000203000039000000000103041a000007c50010009c0000006d0000213d0000000102100039000000000023041b000007cf0110009a000000000041041b000000000103041a000f00000001001d000000000040043f0000000301000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f0000001203000029000000010020019000000a830000613d000000000101043b0000000f02000029000000000021041b000000400100043d0000000000310435000007c10010009c000007c10100804100000040011002100000000002000414000007c10020009c000007c102008041000000c002200210000000000112019f000007d0011001c70000800d020000390000000103000039000007d1040000411f011ef70000040f0000000100200190000010470000c13d00000a830000013d000000400200043d001300000002001d0000082a010000410000000000120435000000040120003900000080020000391f011e320000040f00000013020000290000000001210049000007c10010009c000007c1010080410000006001100210000007c10020009c000007c1020080410000004002200210000000000121019f00001f0300010430000000000001004b0000118c0000613d000001000200043d0000082902200197000000000021004b0000118c0000813d0000001301000029000000000010043f0000000701000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000101043b0000000201100039000000e0020000391f011e550000040f000000400100043d00000013020000290000000002210436000000800300043d000000000003004b0000000003000039000000010300c0390000000000320435000000a00200043d000008290220019700000040031000390000000000230435000000c00200043d000008290220019700000060031000390000000000230435000000e00200043d000000000002004b0000000002000039000000010200c03900000080031000390000000000230435000001000200043d0000082902200197000000a0031000390000000000230435000001200200043d0000082902200197000000c0031000390000000000230435000007c10010009c000007c10100804100000040011002100000000002000414000007c10020009c000007c102008041000000c002200210000000000112019f00000838011001c70000800d020000390000000103000039000008390400004100000a7a0000013d0000000002000019000010eb0000013d00000012020000290000000102200039000000800100043d000000000012004b000003cf0000813d001200000002001d0000000501200210000000a0011000390000000001010433000007c401100197001300000001001d000000000010043f0000000301000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000101043b000000000301041a000000000003004b000010e60000613d0000000201000039000000000201041a000000000002004b000011500000613d000000010130008a000000000023004b000011230000613d000000000012004b0000127b0000a13d000008630130009a000008630220009a000000000202041a000000000021041b000000000020043f0000000301000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c70000801002000039000f00000003001d1f011efc0000040f0000000f03000029000000010020019000000a830000613d000000000101043b000000000031041b0000000201000039000000000301041a000000000003004b000012750000613d000000010130008a000008630230009a000000000002041b0000000202000039000000000012041b0000001301000029000000000010043f0000000301000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000101043b000000000001041b000000400100043d00000013020000290000000000210435000007c10010009c000007c10100804100000040011002100000000002000414000007c10020009c000007c102008041000000c002200210000000000112019f000007d0011001c70000800d02000039000000010300003900000864040000411f011ef70000040f0000000100200190000010e60000c13d00000a830000013d0000000f010000290000000501100039000f00000001001d000000000101041a000900000001001d000000000001004b000011930000c13d0000085801000041000000000010043f0000001101000039000000040010043f000007d50100004100001f03000104300000001101000029000007c10010009c000007c1010080410000004001100210000007c10020009c000007c102008041000000c002200210000000000112019f000007ca011001c700000013020000291f011efc0000040f0000006003100270000007c103300197000000200030008c000000200400003900000000040340190000001f0640018f00000020074001900000001105700029000011700000613d000000000801034f0000001109000029000000008a08043c0000000009a90436000000000059004b0000116c0000c13d000000000006004b0000117d0000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f00030000000103550000000100200190000011ec0000613d0000001f01400039000000600110018f0000001102100029001000000002001d000007c50020009c0000006d0000213d0000001002000029000000400020043f000000200030008c00000fd70000813d00000a830000013d000000400200043d001300000002001d0000082a0100004100000000001204350000000401200039000000e002000039000010970000013d00000009020000290000000a0020006b0000123b0000c13d0000000f01000029000000000010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007d0011001c700008010020000391f011efc0000040f000000010020019000000a830000613d0000000a02000029000000010220008a000000000101043b0000000001210019000000000001041b0000000f01000029000000000021041b0000000b01000029000000000010043f0000000c01000029000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000101043b000000000001041b000000400100043d0000002002100039000000120300002900000000003204350000002002000039000000000021043500000040021000390000000d03200029000000110400002900000002044003670000000d0000006b000011c90000613d000000000504034f0000000006020019000000005705043c0000000006760436000000000036004b000011c50000c13d000000100000006b000011d70000613d0000000d0440036000000010050000290000000305500210000000000603043300000000065601cf000000000656022f000000000404043b0000010005500089000000000454022f00000000045401cf000000000464019f000000000043043500000012022000290000000000020435000007c10010009c000007c10100804100000040011002100000000e02000029000008400020009c00000840020080410000006002200210000000000112019f0000000002000414000007c10020009c000007c102008041000000c002200210000000000121019f000008410110009a0000800d0200003900000002030000390000084204000041000000130500002900000a7a0000013d0000001f0530018f000007c306300198000000400200043d0000000004620019000010320000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000011f30000c13d000010320000013d000000400300043d001300000003001d0000002401300039000000400200003900000000002104350000083001000041000000000013043500000004013000390000000e020000290000000000210435000000440230003900000012010000291f011afa0000040f00000013020000290000000001210049000007c10010009c000007c101008041000007c10020009c000007c10200804100000060011002100000004002200210000000000121019f00001f03000104300000001201000029000000000010043f0000000301000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000101043b000000000101041a000000000001004b00000d120000c13d000000400100043d0000084b02000041000000000021043500000004021000390000001203000029000004940000013d0000000401000039000000000301041a0000086601000041000000100400002900000000001404350000001201000029000000000012043500000024014000390000000002000411000007c402200197001200000002001d00000000002104350000000001000414000007c402300197000000040020008c000012870000c13d0000000103000031000000200030008c00000020040000390000000004034019000012b10000013d0000000f01000029000000000010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007d0011001c700008010020000391f011efc0000040f000000010020019000000a830000613d0000000a020000290008000100200092000000000101043b0000000f02000029000000000202041a000000080020006c0000127b0000a13d0000000902000029000000010220008a0000000001120019000000000101041a000900000001001d0000000f01000029000000000010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007d0011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000101043b00000008011000290000000902000029000000000021041b000000000020043f0000000c01000029000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000101043b0000000a02000029000000000021041b0000000f01000029000000000101041a000a00000001001d000000000001004b000011960000c13d0000085801000041000000000010043f0000003101000039000000040010043f000007d50100004100001f03000104300000085801000041000000000010043f0000003201000039000000040010043f000007d50100004100001f0300010430000000400100043d0000082402000041000000000021043500000004021000390000000d03000029000004940000013d0000001003000029000007c10030009c000007c1030080410000004003300210000007c10010009c000007c101008041000000c001100210000000000131019f000007cc011001c71f011efc0000040f0000006003100270000007c103300197000000200030008c000000200400003900000000040340190000001f0640018f00000020074001900000001005700029000012a00000613d000000000801034f0000001009000029000000008a08043c0000000009a90436000000000059004b0000129c0000c13d000000000006004b000012ad0000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f00030000000103550000000100200190000012d30000613d0000001f01400039000000600210018f0000001001200029000000000021004b00000000020000390000000102004039000007c50010009c0000006d0000213d00000001002001900000006d0000c13d000000400010043f000000200030008c00000a830000413d00000010020000290000000002020433000000000002004b0000000003000039000000010300c039000000000032004b00000a830000c13d000000000002004b000013c50000c13d0000084d02000041000012220000013d000007d4010000410000001003000029000000000013043500000004013000390000000000210435000007c10030009c000007c1030080410000004001300210000007d5011001c700001f03000104300000001f0530018f000007c306300198000000400200043d0000000004620019000010320000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000012da0000c13d000010320000013d000007c10020009c000007c102008041000000c0012002100000001002000029001000000002001d000007c10020009c000007c1020080410000004002200210000000000112019f000007ca011001c700000012020000291f011efc0000040f0000006003100270000007c103300197000000200030008c000000200400003900000000040340190000001f0640018f00000020074001900000001005700029000012fa0000613d000000000801034f0000001009000029000000008a08043c0000000009a90436000000000059004b000012f60000c13d000000000006004b000013070000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f00030000000103550000000100200190000013160000613d0000001f01400039000000600110018f0000001002100029001100000002001d000007c50020009c0000006d0000213d0000001102000029000000400020043f000000200030008c00000fea0000813d00000a830000013d0000001f0530018f000007c306300198000000400200043d0000000004620019000010320000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b0000131d0000c13d000010320000013d0000001003000029000007c10030009c000007c1030080410000004003300210000007c10010009c000007c101008041000000c001100210000000000131019f000007d5011001c71f011efc0000040f0000006003100270000007c103300197000000200030008c000000200400003900000000040340190000001f0640018f000000200740019000000010057000290000133b0000613d000000000801034f0000001009000029000000008a08043c0000000009a90436000000000059004b000013370000c13d000000000006004b000013480000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f000300000001035500000001002001900000146f0000613d0000001f01400039000000600210018f0000001001200029000000000021004b00000000020000390000000102004039000007c50010009c0000006d0000213d00000001002001900000006d0000c13d000000400010043f000000200030008c00000a830000413d00000010020000290000000002020433000007c40020009c00000a830000213d0000000003000411000000000023004b000014dc0000c13d00000002010003670000001102100360000000000202043b000007c50020009c00000a830000213d0000001103000029001000400030003d0000001001100360000000000101043b001200000001001d000000000020043f0000000701000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000101043b000f00000001001d0000081c0100004100000000001004430000000001000412000000040010044300000024000004430000000001000414000007c10010009c000007c101008041000000c00110021000000847011001c700008005020000391f011efc0000040f000000010020019000001ae30000613d000000000101043b000d00000001001d0000000f01000029000000000101041a000e00000001001d0000083a00100198000018190000613d000000120000006b000018190000613d0000000f010000290000000101100039000b00000001001d000000000101041a000c00000001001d000008250100004100000000001004430000000001000414000007c10010009c000007c101008041000000c00110021000000826011001c70000800b020000391f011efc0000040f000000010020019000001ae30000613d0000000e020000290000008002200270000007c102200197000000000101043b000000000421004b000011500000413d0000000e0200002900000829032001970000000c020000290000082902200197000015dd0000c13d0000000d01000029000007c404100197000000120020006c000017cd0000813d000000400100043d00000000050100190000000401100039000000000004004b000017fb0000c13d0000085c0300004100000000003504350000000000210435000000240150003900000012020000290000000000210435000007c10050009c000007c1050080410000004001500210000007cc011001c700001f03000104300000000b010000290000000001010433000000400200043d0000082c030000410000000000320435000007c50110019700000aa60000013d00000002020003670000001101200360000000000101043b000007c50010009c00000a830000213d00000011030000290000008003300039000000000332034f000000000403043b0000000003000031000000130530006a000000230550008a000007c606500197000007c607400197000000000867013f000000000067004b0000000006000019000007c606004041000000000054004b0000000005000019000007c605008041000007c60080009c000000000605c019000000000006004b00000a830000c13d0000001305000029001000040050003d0000001004400029000000000242034f000000000202043b001200000002001d000007c50020009c00000a830000213d000000120230006a0000002006400039000007c603200197000007c604600197000000000534013f000000000034004b0000000003000019000007c603004041000f00000006001d000000000026004b0000000002000019000007c602002041000007c60050009c000000000302c019000000000003004b00000a830000c13d000000000010043f0000000701000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d00000012020000290000001f022000390000087a022001970000003f022000390000087a02200197000000000101043b000e00000001001d000000400100043d0000000002210019000000000012004b00000000030000390000000103004039000007c50020009c0000006d0000213d00000001003001900000006d0000c13d000000400020043f000000120400002900000000024104360000000f04400029000000000040007c00000a830000213d00000012050000290000087a045001980000001f0550018f0000000f0300002900000002063003670000000003420019000014250000613d000000000706034f0000000008020019000000007907043c0000000008980436000000000038004b000014210000c13d000000000005004b000014320000613d000000000446034f0000000305500210000000000603043300000000065601cf000000000656022f000000000404043b0000010005500089000000000454022f00000000045401cf000000000464019f000000000043043500000012032000290000000000030435000007c10020009c000007c10200804100000040022002100000000001010433000007c10010009c000007c1010080410000006001100210000000000121019f0000000002000414000007c10020009c000007c102008041000000c002200210000000000112019f0000081e011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000101043b000000000010043f0000000e010000290000000601100039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000101043b000000000101041a000000000001004b000015eb0000c13d0000001301000029000000a40210003900000010010000291f011b420000040f0000086e03000041000000400500043d001300000005001d0000000000350435000000000301001900000000040200190000000401500039000000000203001900000000030400191f011c150000040f000010980000013d000000400200043d001300000002001d0000082b010000410000147e0000013d000000400300043d001300000003001d0000082b01000041000014850000013d0000001f0530018f000007c306300198000000400200043d0000000004620019000010320000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000014760000c13d000010320000013d000000400200043d001300000002001d0000082a01000041000000000012043500000004012000390000000002030019000010970000013d000000400300043d001300000003001d0000082a0100004100000000001304350000000401300039000010970000013d000007d7010000410000001103000029000012cb0000013d000007c10020009c000007c102008041000000c0012002100000001102000029001100000002001d000007c10020009c000007c1020080410000004002200210000000000112019f000007ca011001c700000013020000291f011efc0000040f0000006003100270000007c103300197000000200030008c000000200400003900000000040340190000001f0640018f00000020074001900000001105700029000014a60000613d000000000801034f0000001109000029000000008a08043c0000000009a90436000000000059004b000014a20000c13d000000000006004b000014b30000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f00030000000103550000000100200190000014e30000613d0000001f01400039000000600110018f0000001102100029001300000002001d000007c50020009c0000006d0000213d0000001302000029000000400020043f000000200030008c00000a830000413d00000011020000290000000002020433000007c10020009c00000a830000213d000001400020043f000000800400043d000001000200043d000007d905000041000000130700002900000000005704350000000405700039000000000600041000000000006504350000002405700039000007c402200197001000000002001d00000000002504350000000002000414000007c404400197001100000004001d000000040040008c000014ef0000c13d0000001301100029000007c50010009c0000006d0000213d000000400010043f000015230000013d0000084d0200004100000000002104350000000002000411000007c40220019700000004031000390000000000230435000004950000013d0000001f0530018f000007c306300198000000400200043d0000000004620019000010320000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000014ea0000c13d000010320000013d0000001301000029001300000001001d000007c10010009c000007c1010080410000004001100210000007c10020009c000007c102008041000000c002200210000000000112019f000007cc011001c700000011020000291f011efc0000040f0000006003100270000007c103300197000000200030008c000000200400003900000000040340190000001f0640018f000000200740019000000013057000290000150a0000613d000000000801034f0000001309000029000000008a08043c0000000009a90436000000000059004b000015060000c13d000000000006004b000015170000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f00030000000103550000000100200190000015810000613d0000001f01400039000000600110018f0000001301100029000007c50010009c0000006d0000213d000000400010043f000000200030008c00000a830000413d00000013020000290000000002020433000000000002004b000011500000c13d0000004402100039000000010400008a00000000004204350000002002100039000007da04000041000000000042043500000024041000390000001005000029000000000054043500000044040000390000000000410435000007db0010009c0000006d0000213d000000c004100039000000400040043f000000a005100039000007dc04000041001300000005001d000000000045043500000080051000390000002004000039000f00000005001d0000000000450435000000000401043300000000010004140000001105000029000000040050008c000015550000613d000007c10020009c000007c1020080410000004002200210000007c10040009c000007c1040080410000006003400210000000000223019f000007c10010009c000007c101008041000000c001100210000000000112019f00000011020000291f011ef70000040f000d00010020019300030000000103550000006001100270000107c10010019d000007c103100197000000000003004b0000158d0000c13d001000600000003d000e00800000003d000000100100002900000000010104330000000d0000006b000015ba0000c13d000000000001004b000015e20000c13d000000400100043d000007e20200004100000000002104350000000402100039000000200300003900000000003204350000000f020000290000000002020433000000240310003900000000002304350000004403100039000000000002004b000015740000613d000000000400001900000000053400190000001306400029000000000606043300000000006504350000002004400039000000000024004b0000156d0000413d0000001f042000390000087a04400197000000000232001900000000000204350000004402400039000007c10020009c000007c1020080410000006002200210000007c10010009c000007c1010080410000004001100210000000000112019f00001f03000104300000001f0530018f000007c306300198000000400200043d0000000004620019000010320000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000015880000c13d000010320000013d000007c50030009c0000006d0000213d0000001f023000390000087a022001970000003f022000390000087a02200197000000400400043d0000000002240019001000000004001d000000000042004b00000000040000390000000104004039000007c50020009c0000006d0000213d00000001004001900000006d0000c13d000000400020043f000000100200002900000000043204360000087a023001980000001f0330018f000e00000004001d00000000012400190000000304000367000015ac0000613d000000000504034f0000000e06000029000000005705043c0000000006760436000000000016004b000015a80000c13d000000000003004b000015590000613d000000000224034f0000000303300210000000000401043300000000043401cf000000000434022f000000000202043b0000010003300089000000000232022f00000000023201cf000000000242019f0000000000210435000015590000013d000000000001004b000018b90000c13d000007dd010000410000000000100443000000110100002900000004001004430000000001000414000007c10010009c000007c101008041000000c001100210000007de011001c700008002020000391f011efc0000040f000000010020019000001ae30000613d000000000101043b000000000001004b000018b50000c13d000000400100043d0000004402100039000007e603000041000000000032043500000024021000390000001d030000390000000000320435000007e2020000410000000000210435000000040210003900000020030000390000000000320435000007c10010009c000007c1010080410000004001100210000007e7011001c700001f0300010430000000000023004b000017e80000a13d000000400100043d000008670200004100000b370000013d0000000e02000029000007c10020009c000007c1020080410000004002200210000007c10010009c000007c1010080410000006001100210000000000121019f00001f030001043000000002010003670000001102100360000000000202043b000007c50020009c00000a830000213d0000001103000029001100400030003d0000001101100360000000000101043b001200000001001d000000000020043f0000000701000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000101043b000f00000001001d0000081c0100004100000000001004430000000001000412000000040010044300000024000004430000000001000414000007c10010009c000007c101008041000000c00110021000000847011001c700008005020000391f011efc0000040f000000010020019000001ae30000613d000000000101043b000d00000001001d0000000f010000290000000201100039000c00000001001d000000000101041a000e00000001001d0000083a00100198000016550000613d000000120000006b000016550000613d0000000f010000290000000301100039000b00000001001d000000000101041a000f00000001001d000008250100004100000000001004430000000001000414000007c10010009c000007c101008041000000c00110021000000826011001c70000800b020000391f011efc0000040f000000010020019000001ae30000613d0000000e020000290000008002200270000007c102200197000000000101043b000000000421004b000011500000413d0000000e0200002900000829032001970000000f020000290000082902200197000019140000c13d0000000d01000029000007c404100197000000120020006c000013ae0000413d000000120130006c000017cf0000413d00000829011001970000000c03000029000000000203041a0000085002200197000000000112019f000000000013041b000000400100043d00000012020000290000000000210435000007c10010009c000007c10100804100000040011002100000000002000414000007c10020009c000007c102008041000000c002200210000000000112019f000007d0011001c70000800d02000039000000010300003900000851040000411f011ef70000040f000000010020019000000a830000613d000000110100002900000060051000390000000201000367000000000251034f000000000202043b0000000003000031000000130430006a000000230440008a000007c606400197000007c607200197000000000867013f000000000067004b0000000007000019000007c607004041000000000042004b0000000009000019000007c609008041000007c60080009c000000000709c019000000000007004b00000a830000c13d0000001002200029000000000721034f000000000707043b000007c50070009c00000a830000213d000000400070008c00000a830000413d00000000087300490000002007200039000000000087004b0000000002000019000007c602002041000007c608800197000007c609700197000000000a89013f000000000089004b0000000008000019000007c608004041000007c600a0009c000000000802c019000000000008004b00000a830000c13d000000400200043d001200000002001d000008570020009c0000006d0000213d00000012020000290000004002200039000000400020043f000000000271034f000000000202043b000007c50020009c00000a830000213d000000120800002900000000022804360000002007700039000000000771034f000000000707043b000007c10070009c00000a830000213d00000000007204350000002005500039000000000551034f000000000505043b000007c607500197000000000867013f000000000067004b0000000006000019000007c606004041000000000045004b0000000004000019000007c604008041000007c60080009c000000000604c019000000000006004b00000a830000c13d0000001005500029000000000451034f000000000404043b000007c50040009c00000a830000213d00000000034300490000002005500039000007c606300197000007c607500197000000000867013f000000000067004b0000000006000019000007c606004041000000000035004b0000000003000019000007c603002041000007c60080009c000000000603c019000000000006004b00000a830000c13d000000200040008c00000a830000413d000000000351034f000000000603043b000007c50060009c00000a830000213d000000000354001900000000045600190000000005430049000007df0050009c00000a830000213d000000400050008c00000a830000413d000000400500043d001100000005001d000008570050009c0000006d0000213d000000000541034f00000011060000290000004006600039000000400060043f000000000505043b000007c50050009c00000a830000213d00000000084500190000001f05800039000000000035004b0000000007000019000007c607008041000007c609500197000007c605300197000000000a59013f000000000059004b0000000009000019000007c609004041000007c600a0009c000000000907c019000000000009004b00000a830000c13d000000000781034f000000000707043b000007c50070009c0000006d0000213d0000001f097000390000087a099001970000003f099000390000087a099001970000000009690019000007c50090009c0000006d0000213d0000002008800039000000400090043f00000000007604350000000009870019000000000039004b00000a830000213d000000000a81034f0000087a0b7001980000001f0c70018f000000110800002900000060088000390000000009b80019000016fb0000613d000000000d0a034f000000000e08001900000000df0d043c000000000efe043600000000009e004b000016f70000c13d00000000000c004b000017080000613d000000000aba034f000000030bc00210000000000c090433000000000cbc01cf000000000cbc022f000000000a0a043b000001000bb00089000000000aba022f000000000aba01cf000000000aca019f0000000000a904350000000007870019000000000007043500000011070000290000000006670436000f00000006001d0000002006400039000000000661034f000000000606043b000007c50060009c00000a830000213d00000000064600190000001f04600039000000000034004b0000000007000019000007c607008041000007c604400197000000000854013f000000000054004b0000000004000019000007c604004041000007c60080009c000000000407c019000000000004004b00000a830000c13d000000000461034f000000000404043b000007c50040009c0000006d0000213d0000001f054000390000087a055001970000003f055000390000087a07500197000000400500043d0000000007750019000000000057004b00000000080000390000000108004039000007c50070009c0000006d0000213d00000001008001900000006d0000c13d0000002008600039000000400070043f00000000064504360000000007840019000000000037004b00000a830000213d000000000381034f0000087a074001980000001f0840018f0000000001760019000017420000613d000000000903034f000000000a060019000000009b09043c000000000aba043600000000001a004b0000173e0000c13d000000000008004b0000174f0000613d000000000373034f0000000307800210000000000801043300000000087801cf000000000878022f000000000303043b0000010007700089000000000373022f00000000037301cf000000000383019f0000000000310435000000000146001900000000000104350000000f0100002900000000005104350000001101000029000000000101043300000004031000390000000003030433000007c103300198000019a60000c13d0000000002020433000007c10220019700000008031000390000000003030433000007c103300197000000000023004b000019ab0000c13d00000014021000390000000002020433000d00000002001d0000000c011000390000000001010433000e00000001001d0000081c01000041000000000010044300000000010004120000000400100443000000c00100003900000024001004430000000001000414000007c10010009c000007c101008041000000c00110021000000847011001c700008005020000391f011efc0000040f000000010020019000001ae30000613d000000400200043d001000000002001d0000000402200039000000000101043b0000000e0310014f000007c100300198000019b70000c13d0000000d01000029000007c50110019700000012030000290000000003030433000007c503300197000000000031004b000019c50000c13d0000000f010000290000000001010433000000110300002900000000030304330000086b0400004100000010060000290000000000460435000000400400003900000000004204350000004404600039000000005303043400000000003404350000006404600039000000000003004b0000179a0000613d000000000600001900000000074600190000000008650019000000000808043300000000008704350000002006600039000000000036004b000017930000413d000000000534001900000000000504350000001f033000390000087a033001970000000003340019000000000223004900000010040000290000002404400039000000000024043500000000140104340000000002430436001200000002001d001100000004001d000000000004004b000017b10000613d000000000200001900000012032000290000000004210019000000000404043300000000004304350000002002200039000000110020006c000017aa0000413d0000001102000029000000120120002900000000000104350000081c01000041000000000010044300000000010004120000000400100443000000a00100003900000024001004430000000001000414000007c10010009c000007c101008041000000c00110021000000847011001c700008005020000391f011efc0000040f000000010020019000001ae30000613d000000000201043b0000000001000414000007c402200197000000040020008c000019cc0000c13d0000000103000031000000200030008c00000020040000390000000004034019000019ff0000013d000000120130006c000018010000813d0000000b01000029000000000101041a0000008001100272000011500000613d00000012053000690000000002150019000000010220008a000000000052004b000011500000413d00000000021200d9000000400100043d00000000060100190000000401100039000000000004004b0000190c0000c13d0000085a040000410000000000460435000000000021043500000024016000390000000000310435000007c10060009c000007c1060080410000004001600210000007cc011001c700001f03000104300000000c05000029000000800650027000000000056400a900000000044500d9000000000064004b000011500000c13d000000000035001a000011500000413d000000000335001900000080011002100000084e011001970000000f05000029000000000405041a0000084f04400197000000000114019f000000000015041b000000000032004b0000000003024019000013aa0000013d0000085b03000041001300000005001d000000000035043500000012030000291f011ebb0000040f000010980000013d00000829011001970000000f03000029000000000203041a0000085002200197000000000112019f000000000013041b000000400100043d00000012020000290000000000210435000007c10010009c000007c10100804100000040011002100000000002000414000007c10020009c000007c102008041000000c002200210000000000112019f000007d0011001c70000800d02000039000000010300003900000851040000411f011ef70000040f000000010020019000000a830000613d00000011010000290000000201100367000000000101043b000007c50010009c00000a830000213d000000000010043f0000000a01000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000400200043d000008280020009c0000006d0000213d000000000101043b0000006003200039000000400030043f000000000301041a001200000003001d00000000033204360000000101100039000000000101041a000000400220003900000834001001980000000004000039000000010400c0390000000000420435000007c101100197000f00000001001d00000000001304350000000202000367000018ac0000613d0000001101000029000000200110008a000000000112034f000000000401043b0000000003000031000000130130006a000000230110008a000007c605100197000007c606400197000000000756013f000000000056004b0000000005000019000007c605004041000000000014004b0000000001000019000007c601008041000007c60070009c000000000501c019000000000005004b00000a830000c13d000000130100002900000004011000390000000005140019000000000452034f000000000404043b000007c50040009c00000a830000213d00000000064300490000002003500039000007c605600197000007c607300197000000000857013f000000000057004b0000000005000019000007c605004041000000000063004b0000000006000019000007c606002041000007c60080009c000000000506c019000000000005004b00000a830000c13d000000200040008c000019100000c13d000000000132034f0000001002200360000000000202043b001000000002001d000000000101043b000e00000001001d000000400200043d0000085401000041001300000002001d00000000001204350000081c0100004100000000001004430000000001000412000000040010044300000024000004430000000001000414000007c10010009c000007c101008041000000c00110021000000847011001c700008005020000391f011efc0000040f000000010020019000001ae30000613d000000000101043b000000130400002900000084024000390000001203000029000000000032043500000044024000390000000e03000029000000000032043500000024024000390000000f0300002900000000003204350000000402400039000000100300002900000000003204350000006402400039000007c40110019700000000001204350000081c01000041000000000010044300000000010004120000000400100443000000800100003900000024001004430000000001000414000007c10010009c000007c101008041000000c00110021000000847011001c700008005020000391f011efc0000040f000000010020019000001ae30000613d000000000201043b0000000001000414000007c402200197000000040020008c000019290000c13d0000000104000031000000200040008c0000002004008039000019530000013d00000013010000290000002401100039000000000112034f000000000101043b000007c50010009c00000a830000213d000000400200043d000008520300004100000aa50000013d00000010010000290000000001010433000000000001004b000018c60000613d000007df0010009c00000a830000213d000000200010008c00000a830000413d0000000e010000290000000001010433000000000001004b0000000002000039000000010200c039000000000021004b00000a830000c13d000000000001004b000018f80000613d000000400100043d00000012020000290000000000210435000007c10010009c000007c10100804100000040011002100000000002000414000007c10020009c000007c102008041000000c002200210000000000121019f000007d0011001c70000800d020000390000000103000039000007e4040000411f011ef70000040f000000010020019000000a830000613d000000800100043d00000140000004430000016000100443000000a00100043d00000020030000390000018000300443000001a000100443000000c00100043d0000004002000039000001c000200443000001e0001004430000006001000039000000e00200043d000002000010044300000220002004430000008001000039000001000200043d00000240001004430000026000200443000000a001000039000001200200043d0000028000100443000002a000200443000000c001000039000001400200043d000002c000100443000002e000200443000001000030044300000007010000390000012000100443000007e50100004100001f020001042e000000400100043d0000006402100039000007e00300004100000000003204350000004402100039000007e103000041000000000032043500000024021000390000002a030000390000000000320435000007e2020000410000000000210435000000040210003900000020030000390000000000320435000007c10010009c000007c1010080410000004001100210000007e3011001c700001f03000104300000085905000041001300000006001d0000000000560435000017ff0000013d00000000020100191f011b420000040f00000853030000410000145d0000013d000000000023004b000015df0000213d0000000f05000029000000800650027000000000056400a900000000044500d9000000000064004b000011500000c13d000000000035001a000011500000413d000000000335001900000080011002100000084e011001970000000c05000029000000000405041a0000084f04400197000000000114019f000000000015041b000000000032004b0000000003024019000016370000013d0000001303000029000007c10030009c000007c1030080410000004003300210000007c10010009c000007c101008041000000c001100210000000000131019f00000855011001c71f011ef70000040f0000006003100270000007c103300197000000200030008c000000200400003900000000040340190000001f0640018f00000020074001900000001305700029000019420000613d000000000801034f0000001309000029000000008a08043c0000000009a90436000000000059004b0000193e0000c13d000000000006004b0000194f0000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f000300000001035500000001002001900000199a0000613d0000001f01400039000000600210018f0000001301200029000000000021004b00000000020000390000000102004039000007c50010009c0000006d0000213d00000001002001900000006d0000c13d000000400010043f000000200040008c00000a830000413d00000013020000290000000002020433000f00000002001d000007c50020009c00000a830000213d00000010020000290000000000210435000007c10010009c000007c10100804100000040011002100000000002000414000007c10020009c000007c102008041000000c002200210000000000112019f000007d0011001c70000800d020000390000000203000039000008560400004100000000050004111f011ef70000040f000000010020019000000a830000613d00000011010000290000000201100367000000000101043b000007c50010009c00000a830000213d000000000010043f0000000701000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000000101043b0000000401100039000000000201041a000000010320019000000001042002700000007f0440618f001300000004001d0000001f0040008c00000000040000390000000104002039000000000043004b00001a170000613d0000085801000041000000000010043f0000002201000039000000040010043f000007d50100004100001f03000104300000001f0530018f000007c306300198000000400200043d0000000004620019000010320000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b000019a10000c13d000010320000013d000000400100043d000007d40200004100000000002104350000000402100039000004940000013d000000400100043d000000240410003900000000003404350000086803000041000000000031043500000004031000390000000000230435000007c10010009c000007c1010080410000004001100210000007cc011001c700001f0300010430000008690300004100000010040000290000000000340435000007c10110019700000000001204350000000e01000029000007c10110019700000024024000390000000000120435000007c10040009c000007c1040080410000004001400210000007cc011001c700001f03000104300000086a0400004100000010050000290000000000450435000000000032043500000024025000390000000000120435000013b90000013d00000011030000290000001f033000390000087a03300197000000100500002900000000035300490000001203300029000007c10030009c000007c1030080410000006003300210000007c10050009c000007c10400004100000000040540190000004004400210000000000343019f000007c10010009c000007c101008041000000c001100210000000000131019f1f011ef70000040f0000006003100270000007c103300197000000200030008c000000200400003900000000040340190000001f0640018f00000020074001900000001005700029000019ee0000613d000000000801034f0000001009000029000000008a08043c0000000009a90436000000000059004b000019ea0000c13d000000000006004b000019fb0000613d000000000771034f0000000306600210000000000805043300000000086801cf000000000868022f000000000707043b0000010006600089000000000767022f00000000066701cf000000000686019f0000000000650435000100000003001f0003000000010355000000010020019000001a350000613d0000001f01400039000000600210018f0000001001200029000000000021004b00000000020000390000000102004039000007c50010009c0000006d0000213d00000001002001900000006d0000c13d000000400010043f000000200030008c00000a830000413d00000010020000290000000002020433000000000002004b0000000003000039000000010300c039000000000032004b00000a830000c13d000000000002004b00001a410000c13d0000086d0200004100000b370000013d000000400400043d001100000004001d00000013050000290000000004540436001200000004001d000000000003004b00001a6b0000613d000000000010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007d0011001c700008010020000391f011efc0000040f000000010020019000000a830000613d000000130000006b000000000200001900001a710000613d000000000101043b00000000020000190000001203200029000000000401041a000000000043043500000001011000390000002002200039000000130020006c00001a2d0000413d00001a710000013d0000001f0530018f000007c306300198000000400200043d0000000004620019000010320000613d000000000701034f0000000008020019000000007907043c0000000008980436000000000048004b00001a3c0000c13d000010320000013d000000130200002900000044022000390000000203000367000000000423034f000000000604043b000007c40060009c00000a830000213d0000002002200039000000000223034f000000000202043b001300000002001d0000000000210435000007c10010009c000007c10100804100000040011002100000000002000414000007c10020009c000007c102008041000000c002200210000000000112019f000007d0011001c70000800d0200003900000003030000390000086c0400004100000000050004111f011ef70000040f000000010020019000000a830000613d000000400100043d000007cd0010009c0000006d0000213d0000002002100039000000400020043f00000013020000290000000000210435000000400100043d0000000000210435000007c10010009c000007c10100804100000040011002100000085e011001c700001f020001042e0000087c0120019700000012020000290000000000120435000000130000006b000000200200003900000000020060390000003f012000390000087a011001970000001102100029000000000012004b00000000010000390000000101004039001300000002001d000007c50020009c0000006d0000213d00000001001001900000006d0000c13d0000001301000029000000400010043f000008570010009c0000006d0000213d00000013020000290000004001200039000000400010043f0000000f01000029000007c5011001970000000001120436001200000001001d0000081c01000041000000000010044300000000010004120000000400100443000000c00100003900000024001004430000000001000414000007c10010009c000007c101008041000000c00110021000000847011001c700008005020000391f011efc0000040f000000010020019000001ae30000613d000000000101043b000007c1011001970000001204000029000000000014043500000013010000290000000001010433000007c502100197000000400100043d000000200310003900000000002304350000000002040433000007c1022001970000004003100039000000000023043500000040020000390000000000210435000008280010009c0000006d0000213d0000006003100039000000400030043f000008570030009c0000006d0000213d000000a002100039000000400020043f00000011020000290000000000230435000000800210003900000000001204350000002004000039000000400100043d00000000044104360000000003030433000000400500003900000000005404350000006004100039000000005303043400000000003404350000008004100039000000000003004b00001ac60000613d000000000600001900000000074600190000000008650019000000000808043300000000008704350000002006600039000000000036004b00001abf0000413d000000000534001900000000000504350000001f033000390000087a0330019700000000053400190000000003150049000000200330008a00000000020204330000004004100039000000000034043500000000430204340000000002350436000000000003004b00001adc0000613d000000000500001900000000062500190000000007540019000000000707043300000000007604350000002005500039000000000035004b00001ad50000413d0000001f043000390000087a04400197000000000332001900000000000304350000000003140049000000000223001900000a510000013d000000000001042f0000087d0010009c00001ae90000813d000000a001100039000000400010043f000000000001042d0000085801000041000000000010043f0000004101000039000000040010043f000007d50100004100001f03000104300000087e0010009c00001af40000813d0000006001100039000000400010043f000000000001042d0000085801000041000000000010043f0000004101000039000000040010043f000007d50100004100001f030001043000000000430104340000000001320436000000000003004b00001b060000613d000000000200001900000000052100190000000006240019000000000606043300000000006504350000002002200039000000000032004b00001aff0000413d000000000231001900000000000204350000001f023000390000087a022001970000000001210019000000000001042d000007df0010009c00001b270000213d000000430010008c00001b270000a13d00000002020003670000000403200370000000000403043b000007c50040009c00001b270000213d0000002403200370000000000503043b000007c50050009c00001b270000213d0000002303500039000000000013004b00001b270000813d0000000403500039000000000232034f000000000302043b000007c50030009c00001b270000213d00000024025000390000000005320019000000000015004b00001b270000213d0000000001040019000000000001042d000000000100001900001f03000104300000000043020434000008290330019700000000033104360000000004040433000007c104400197000000000043043500000040032000390000000003030433000000000003004b0000000003000039000000010300c039000000400410003900000000003404350000006003200039000000000303043300000829033001970000006004100039000000000034043500000080022000390000000002020433000008290220019700000080031000390000000000230435000000a001100039000000000001042d0000000204000367000000000224034f000000000202043b000000000300003100000000051300490000001f0550008a000007c606500197000007c607200197000000000867013f000000000067004b0000000006000019000007c606002041000000000052004b0000000005000019000007c605004041000007c60080009c000000000605c019000000000006004b00001b6a0000613d0000000001120019000000000214034f000000000202043b000007c50020009c00001b6a0000213d00000000032300490000002001100039000007c604300197000007c605100197000000000645013f000000000045004b0000000004000019000007c604004041000000000031004b0000000003000019000007c603002041000007c60060009c000000000403c019000000000004004b00001b6a0000c13d000000000001042d000000000100001900001f03000104300000087f0020009c00001b9c0000813d00000000040100190000001f012000390000087a011001970000003f011000390000087a05100197000000400100043d0000000005510019000000000015004b00000000070000390000000107004039000007c50050009c00001b9c0000213d000000010070019000001b9c0000c13d000000400050043f00000000052104360000000007420019000000000037004b00001ba20000213d0000087a062001980000001f0720018f0000000204400367000000000365001900001b8c0000613d000000000804034f0000000009050019000000008a08043c0000000009a90436000000000039004b00001b880000c13d000000000007004b00001b990000613d000000000464034f0000000306700210000000000703043300000000076701cf000000000767022f000000000404043b0000010006600089000000000464022f00000000046401cf000000000474019f000000000043043500000000022500190000000000020435000000000001042d0000085801000041000000000010043f0000004101000039000000040010043f000007d50100004100001f0300010430000000000100001900001f03000104300003000000000002000300000003001d000200000002001d000007c501100197000000000010043f0000000701000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f0000000207000029000000030a000029000000010020019000001c0d0000613d0000000003000031000000000601043b0000087f00a0009c00001c0f0000813d0000001f01a000390000087a011001970000003f011000390000087a02100197000000400100043d0000000002210019000000000012004b00000000050000390000000105004039000007c50020009c00001c0f0000213d000000010050019000001c0f0000c13d000100000006001d000000400020043f0000000002a1043600000000057a0019000000000035004b00001c0d0000213d0000087a04a001980000001f05a0018f0000000206700367000000000342001900001bd80000613d000000000706034f0000000008020019000000007907043c0000000008980436000000000038004b00001bd40000c13d000000000005004b00001be50000613d000000000446034f0000000305500210000000000603043300000000065601cf000000000656022f000000000404043b0000010005500089000000000454022f00000000045401cf000000000464019f00000000004304350000000003a200190000000000030435000007c10020009c000007c10200804100000040022002100000000001010433000007c10010009c000007c1010080410000006001100210000000000121019f0000000002000414000007c10020009c000007c102008041000000c002200210000000000112019f0000081e011001c700008010020000391f011efc0000040f000000010020019000001c0d0000613d000000000101043b000000000010043f00000001010000290000000601100039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000001c0d0000613d000000000101043b000000000101041a000000000001004b0000000001000039000000010100c039000000000001042d000000000100001900001f03000104300000085801000041000000000010043f0000004101000039000000040010043f000007d50100004100001f03000104300000002004000039000000000441043600000000003404350000087a063001980000001f0730018f00000040011000390000000005610019000000020220036700001c240000613d000000000802034f0000000009010019000000008a08043c0000000009a90436000000000059004b00001c200000c13d000000000007004b00001c310000613d000000000262034f0000000306700210000000000705043300000000076701cf000000000767022f000000000202043b0000010006600089000000000262022f00000000026201cf000000000272019f0000000000250435000000000231001900000000000204350000001f023000390000087a022001970000000001120019000000000001042d000000400100043d0000087d0010009c00001c460000813d000000a002100039000000400020043f000000800210003900000000000204350000006002100039000000000002043500000040021000390000000000020435000000200210003900000000000204350000000000010435000000000001042d0000085801000041000000000010043f0000004101000039000000040010043f000007d50100004100001f03000104300003000000000002000007c501100197000000000010043f0000000701000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000001c990000613d000000000101043b0000000405100039000000000205041a000000010320019000000001062002700000007f0660618f0000001f0060008c00000000040000390000000104002039000000000043004b00001c9b0000c13d000000400100043d0000000004610436000000000003004b00001c850000613d000100000004001d000200000006001d000300000001001d000000000050043f0000000001000414000007c10010009c000007c101008041000000c001100210000007d0011001c700008010020000391f011efc0000040f000000010020019000001c990000613d0000000206000029000000000006004b00001c8b0000613d000000000201043b0000000005000019000000030100002900000001070000290000000003570019000000000402041a000000000043043500000001022000390000002005500039000000000065004b00001c7d0000413d00001c8d0000013d0000087c022001970000000000240435000000000006004b0000002005000039000000000500603900001c8d0000013d000000000500001900000003010000290000003f035000390000087a023001970000000003120019000000000023004b00000000020000390000000102004039000007c50030009c00001ca10000213d000000010020019000001ca10000c13d000000400030043f000000000001042d000000000100001900001f03000104300000085801000041000000000010043f0000002201000039000000040010043f000007d50100004100001f03000104300000085801000041000000000010043f0000004101000039000000040010043f000007d50100004100001f03000104300007000000000002000400000001001d000600000002001d0000000021020434000000000001004b00001dc00000613d000007c10010009c000007c1010080410000006001100210000007c10020009c000500000002001d000007c1020080410000004002200210000000000121019f0000000002000414000007c10020009c000007c102008041000000c002200210000000000112019f0000081e011001c700008010020000391f011efc0000040f000000010020019000001db80000613d000000000101043b000700000001001d0000000401000029000007c501100197000200000001001d000000000010043f0000000701000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000001db80000613d000000000201043b0000000701000029000000000010043f000400000002001d0000000601200039000300000001001d000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000001db80000613d000000000101043b000000000101041a000000000001004b00001dc80000c13d00000004010000290000000502100039000000000102041a0000087f0010009c00001dba0000813d000100000001001d0000000101100039000000000012041b000400000002001d000000000020043f0000000001000414000007c10010009c000007c101008041000000c001100210000007d0011001c700008010020000391f011efc0000040f0000000100200190000000070200002900001db80000613d000000000101043b0000000101100029000000000021041b0000000401000029000000000101041a000400000001001d000000000020043f0000000301000029000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f0000000100200190000000070200002900001db80000613d000000000101043b0000000403000029000000000031041b000000000020043f0000000801000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000001db80000613d000000000801043b00000006010000290000000004010433000007c50040009c00001dba0000213d000000000108041a000000010210019000000001031002700000007f0330618f0000001f0030008c00000000010000390000000101002039000000000012004b000000050700002900001ddf0000c13d000000200030008c000400000008001d000700000004001d00001d4b0000413d000300000003001d000000000080043f0000000001000414000007c10010009c000007c101008041000000c001100210000007d0011001c700008010020000391f011efc0000040f000000010020019000001db80000613d00000007040000290000001f024000390000000502200270000000200040008c0000000002004019000000000301043b00000003010000290000001f01100039000000050110027000000000011300190000000002230019000000000012004b0000000507000029000000040800002900001d4b0000813d000000000002041b0000000102200039000000000012004b00001d470000413d000000200040008c000000200a00008a000000200b00003900001d7b0000413d000000000080043f0000000001000414000007c10010009c000007c101008041000000c001100210000007d0011001c700008010020000391f011efc0000040f000000010020019000001db80000613d0000000709000029000000200a00008a0000000002a90170000000000101043b000000200b00003900001db10000613d000000010320008a000000050330027000000000043100190000002003000039000000010440003900000005070000290000000606000029000000040800002900000000056300190000000005050433000000000051041b00000020033000390000000101100039000000000041004b00001d670000c13d000000000092004b00001d780000813d0000000302900210000000f80220018f0000087b0220027f0000087b0220016700000000036300190000000003030433000000000223016f000000000021041b000000010190021000000001011001bf00001d870000013d000000000004004b00001d850000613d00000003014002100000087b0110027f0000087b011001670000000002070433000000000112016f0000000102400210000000000121019f00001d860000013d00000000010000190000000606000029000000000018041b000000400100043d0000000003b10436000000000206043300000000002304350000004003100039000000000002004b00001d970000613d000000000400001900000000053400190000000006740019000000000606043300000000006504350000002004400039000000000024004b00001d900000413d0000001f042000390000000004a4016f000000000223001900000000000204350000004002400039000007c10020009c000007c1020080410000006002200210000007c10010009c000007c1010080410000004001100210000000000112019f0000000002000414000007c10020009c000007c102008041000000c002200210000000000121019f0000081e011001c70000800d020000390000000203000039000008310400004100000002050000291f011ef70000040f000000010020019000001db80000613d000000000001042d00000000030b0019000000050700002900000006060000290000000408000029000000000092004b00001d700000413d00001d780000013d000000000100001900001f03000104300000085801000041000000000010043f0000004101000039000000040010043f000007d50100004100001f0300010430000000400100043d00000833020000410000000000210435000007c10010009c000007c1010080410000004001100210000007ca011001c700001f0300010430000000400300043d000700000003001d00000024013000390000004002000039000000000021043500000830010000410000000000130435000000040130003900000002020000290000000000210435000000440230003900000006010000291f011afa0000040f00000007020000290000000001210049000007c10010009c000007c101008041000007c10020009c000007c10200804100000060011002100000004002200210000000000121019f00001f03000104300000085801000041000000000010043f0000002201000039000000040010043f000007d50100004100001f03000104300005000000000002000000400300043d0000087d0030009c00001e2b0000813d000000a002300039000000400020043f00000080023000390000000000020435000000600230003900000000000204350000004002300039000000000002043500000020023000390000000000020435000000000003043500000060021000390000000002020433000100000002001d000500000001001d0000000012010434000300000002001d000200000001001d0000000001010433000400000001001d000008250100004100000000001004430000000001000414000007c10010009c000007c101008041000000c00110021000000826011001c70000800b020000391f011efc0000040f000000010020019000001e310000613d0000000402000029000007c104200197000000000601043b000000000346004b000000050100002900001e250000413d00000080021000390000000002020433000008290520019700000000023500a9000000000046004b00001e170000613d00000000033200d9000000000053004b00001e250000c13d00000003030000290000082903300197000000000032001a00001e250000413d000000000232001900000001030000290000082903300197000000000023004b00000000030280190000000000310435000007c10260019700000002030000290000000000230435000000000001042d0000085801000041000000000010043f0000001101000039000000040010043f000007d50100004100001f03000104300000085801000041000000000010043f0000004101000039000000040010043f000007d50100004100001f0300010430000000000001042f0000000043020434000000000003004b0000000003000039000000010300c0390000000003310436000000000404043300000829044001970000000000430435000000400220003900000000020204330000082902200197000000400310003900000000002304350000006001100039000000000001042d000000000010043f0000000601000039000000200010043f0000000001000414000007c10010009c000007c101008041000000c001100210000007ce011001c700008010020000391f011efc0000040f000000010020019000001e530000613d000000000101043b000000000101041a000000000001004b0000000001000039000000010100c039000000000001042d000000000100001900001f03000104300003000000000002000100000002001d000300000001001d000000000101041a000200000001001d000008250100004100000000001004430000000001000414000007c10010009c000007c101008041000000c00110021000000826011001c70000800b020000391f011efc0000040f000000010020019000001eb80000613d00000002080000290000008002800270000007c103200197000000000201043b000000000532004b000000030700002900001eb20000413d000000010170003900001e700000c13d000000000207041a00001e830000013d000000000301041a000000800630027000000000045600a900000000055400d9000000000065004b00001eb20000c13d0000082905800197000000000054001a00001eb20000413d00000000045400190000082903300197000000000043004b000000000304801900000080022002100000084e02200197000000000223019f000000000307041a000007c803300197000000000232019f00000001060000290000002003600039000000000403043300000829044001970000082905200197000000000054004b00000000050440190000088002200197000000000225019f0000000005060433000000000005004b00000000050000190000082e0500c041000000000252019f000000000027041b000000400260003900000000050204330000008005500210000000000445019f000000000041041b0000000001000039000000010100c039000000400400043d00000000011404360000000003030433000008290330019700000000003104350000000001020433000008290110019700000040024000390000000000120435000007c10040009c000007c10400804100000040014002100000000002000414000007c10020009c000007c102008041000000c002200210000000000112019f00000881011001c70000800d02000039000000010300003900000882040000411f011ef70000040f000000010020019000001eb90000613d000000000001042d0000085801000041000000000010043f0000001101000039000000040010043f000007d50100004100001f0300010430000000000001042f000000000100001900001f0300010430000007c404400197000000400510003900000000004504350000002004100039000000000034043500000000002104350000006001100039000000000001042d000000000001042f000007c10010009c000007c1010080410000004001100210000007c10020009c000007c1020080410000006002200210000000000112019f0000000002000414000007c10020009c000007c102008041000000c002200210000000000112019f0000081e011001c700008010020000391f011efc0000040f000000010020019000001ed70000613d000000000101043b000000000001042d000000000100001900001f030001043000000000050100190000000000200443000000050030008c00001ee70000413d000000040100003900000000020000190000000506200210000000000664001900000005066002700000000006060031000000000161043a0000000102200039000000000031004b00001edf0000413d000007c10030009c000007c10300804100000060013002100000000002000414000007c10020009c000007c102008041000000c002200210000000000112019f00000883011001c700000000020500191f011efc0000040f000000010020019000001ef60000613d000000000101043b000000000001042d000000000001042f00001efa002104210000000102000039000000000001042d0000000002000019000000000001042d00001eff002104230000000102000039000000000001042d0000000002000019000000000001042d00001f010000043200001f020001042e00001f030001043000000000000000000000000000000000000000000000000000000000ffffffff00000000000000000000000000000000000000000000000000000001ffffffe000000000000000000000000000000000000000000000000000000000ffffffe0000000000000000000000000ffffffffffffffffffffffffffffffffffffffff000000000000000000000000000000000000000000000000ffffffffffffffff80000000000000000000000000000000000000000000000000000000000000007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0ffffffffffffffffffffffff0000000000000000000000000000000000000000313ce567000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000655a7c0e000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000044000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffffdf0200000000000000000000000000000000000040000000000000000000000000bfa87805ed57dc1f0d489ce33be4c4577d74ccde357eeeee058a32c55c44a53202000000000000000000000000000000000000200000000000000000000000002640d4d76caf8bf478aabfa982fa4e1c4eb71a37f93cd15e80dbc657911546d82c1219210000000000000000000000000000000000000000000000000000000054fd4d500000000000000000000000000000000000000000000000000000000068d2f8d60000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000240000000000000000000000009cdbb18100000000000000000000000000000000000000000000000000000000b5d1ce28000000000000000000000000000000000000000000000000000000008d3638f400000000000000000000000000000000000000000000000000000000dd62ed3e00000000000000000000000000000000000000000000000000000000095ea7b300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffff3f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65641806aa1896bbf26568e884a7374b41e002500962caba6a15023a8d90e8508b8302000002000000000000000000000000000000240000000000000000000000007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6f742073756363656564000000000000000000000000000000000000000000005361666545524332303a204552433230206f7065726174696f6e20646964206e08c379a00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000840000000000000000000000002e902d38f15b233cbb63711add0fca4545334d3a169d60c0a616494d7eea95440000000200000000000000000000000000000200000001000000000000000000416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000000000000000000000000000000000000000006400000000000000000000000035be3ac8000000000000000000000000000000000000000000000000000000009b15e16f00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000009a4575b800000000000000000000000000000000000000000000000000000000c4bffe2a00000000000000000000000000000000000000000000000000000000dfadfa3400000000000000000000000000000000000000000000000000000000e8a1da1600000000000000000000000000000000000000000000000000000000e8a1da1700000000000000000000000000000000000000000000000000000000f2fde38b00000000000000000000000000000000000000000000000000000000fbf84dd700000000000000000000000000000000000000000000000000000000dfadfa3500000000000000000000000000000000000000000000000000000000e0351e1300000000000000000000000000000000000000000000000000000000cf7401f200000000000000000000000000000000000000000000000000000000cf7401f300000000000000000000000000000000000000000000000000000000dc0bd97100000000000000000000000000000000000000000000000000000000c4bffe2b00000000000000000000000000000000000000000000000000000000c75eea9c00000000000000000000000000000000000000000000000000000000acfecf9000000000000000000000000000000000000000000000000000000000b0f479a000000000000000000000000000000000000000000000000000000000b0f479a100000000000000000000000000000000000000000000000000000000b794658000000000000000000000000000000000000000000000000000000000c0d7865500000000000000000000000000000000000000000000000000000000acfecf9100000000000000000000000000000000000000000000000000000000af58d59f00000000000000000000000000000000000000000000000000000000a42a7b8a00000000000000000000000000000000000000000000000000000000a42a7b8b00000000000000000000000000000000000000000000000000000000a7cd63b7000000000000000000000000000000000000000000000000000000009a4575b9000000000000000000000000000000000000000000000000000000009fdf13ff0000000000000000000000000000000000000000000000000000000054c8a4f2000000000000000000000000000000000000000000000000000000006d3d1a57000000000000000000000000000000000000000000000000000000007d54534d000000000000000000000000000000000000000000000000000000007d54534e000000000000000000000000000000000000000000000000000000008926f54f000000000000000000000000000000000000000000000000000000008da5cb5b000000000000000000000000000000000000000000000000000000006d3d1a580000000000000000000000000000000000000000000000000000000079ba50970000000000000000000000000000000000000000000000000000000062ddd3c30000000000000000000000000000000000000000000000000000000062ddd3c4000000000000000000000000000000000000000000000000000000006b716b0d0000000000000000000000000000000000000000000000000000000054c8a4f3000000000000000000000000000000000000000000000000000000006155cda000000000000000000000000000000000000000000000000000000000240028e700000000000000000000000000000000000000000000000000000000390775360000000000000000000000000000000000000000000000000000000039077537000000000000000000000000000000000000000000000000000000004c5ef0ed00000000000000000000000000000000000000000000000000000000240028e80000000000000000000000000000000000000000000000000000000024f65ee700000000000000000000000000000000000000000000000000000000181f5a7600000000000000000000000000000000000000000000000000000000181f5a770000000000000000000000000000000000000000000000000000000021df0da7000000000000000000000000000000000000000000000000000000000041d3c10000000000000000000000000000000000000000000000000000000001ffc9a7310ab089e4439a4c15d089f94afb7896ff553aecb10793d0ab882de59d99a32e00000000000000000000000000000000000000200000008000000000000000000200000000000000000000000000000000000000000000000000000000000000ed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278dad89dca000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000800000000000000000fc949c7b4a13586e39d89eead2f38644f9fb3efb5a0490b14f8fc0ceab44c2515204aec90a3c794d8e90fded8b46ae9c7c552803e7e832e0c1d358396d8599161e670e4b00000000000000000000000000000000000000000000000000000000796b89b91644bc98cd93958e4c9038275d622183e25ac5af08cc6b5d955391320200000200000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffff5f000000000000000000000000000000000000000000000000ffffffffffffff9f00000000000000000000000000000000ffffffffffffffffffffffffffffffff8020d12400000000000000000000000000000000000000000000000000000000d68af9cc000000000000000000000000000000000000000000000000000000001d5ad3c500000000000000000000000000000000000000000000000000000000fc949c7b4a13586e39d89eead2f38644f9fb3efb5a0490b14f8fc0ceab44c2500000000000000000000000010000000000000000000000000000000000000000ffffffffffffffffffffff000000000000000000000000000000000000000000393b8ad2000000000000000000000000000000000000000000000000000000007d628c9a1796743d365ab521a8b2a4686e419b3269919dc9145ea2ce853b54ea8d340f17e19058004c20453540862a9c62778504476f6756755cb33bcd6c38c28579befe00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ff0000000000000000000000000000000000000000000000600000000000000000000000008e4a23d600000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002400000140000000000000000002000000000000000000000000000000000000e00000000000000000000000000350d63aa5f270e01729d00d627eeb8f3429772b1818c016c66a588a864f912b0000000000000000000000ff0000000000000000000000000000000000000000036b6384b5eca791c62761152d0c79bb0604c104a5fb6f4eb0703f3154bb3db0000000000000000000000000000000000000000000000000ffffffffffffff7f00000000000000000000000000000000000000000000003fffffffffffffffe0020000000000000000000000000000000000004000000080000000000000000002dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f168400000000000000000000000000000000000000000000000000000000ffffffbffdffffffffffffffffffffffffffffffffffffc000000000000000000000000052d00ee4d9bd51b40168f2afc5848837288ce258784ad914278791464b3f4d7674f23c7c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ffffff9bffffffffffffffffffffffffffffffffffffff9c000000000000000000000000405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace0200000200000000000000000000000000000044000000000000000000000000961c9a4f000000000000000000000000000000000000000000000000000000002cbc26bb0000000000000000000000000000000000000000000000000000000053ad11d800000000000000000000000000000000000000000000000000000000d0d2597600000000000000000000000000000000000000000000000000000000a8d87a3b00000000000000000000000000000000000000000000000000000000728fe07b00000000000000000000000000000000000000000000000000000000000000000000000000000000ffffffff00000000000000000000000000000000ffffffffffffffffffffffff00000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000000000000000000000000000001871cdf8010e63f2eb8384381a68dfa7416dc571a5517e66e88b2d2d0c0a690ad201c48a00000000000000000000000000000000000000000000000000000000a3c8cf0900000000000000000000000000000000000000000000000000000000f856ddb60000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a4000000000000000000000000696de425f79f4a40bc6d2122ca50507f0efbeabbff86a84871b7196ab8ea8df7000000000000000000000000000000000000000000000000ffffffffffffffbf4e487b7100000000000000000000000000000000000000000000000000000000d0c8d23a0000000000000000000000000000000000000000000000000000000015279c08000000000000000000000000000000000000000000000000000000001a76572a00000000000000000000000000000000000000000000000000000000f94ebcd100000000000000000000000000000000000000000000000000000000a9902c7e000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000020000000000000000000000000000000000002000000080000000000000000044676b5284b809a22248eba0da87391d79098be38bb03154be88a58bf4d0917402b543c6000000000000000000000000000000000000000000000000000000008be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0bfa87805ed57dc1f0d489ce33be4c4577d74ccde357eeeee058a32c55c44a533800671136ab6cfee9fbe5ed1fb7ca417811aca3cf864800d127b927adedf756635f4a7b30000000000000000000000000000000000000000000000000000000083826b2b000000000000000000000000000000000000000000000000000000009725942a00000000000000000000000000000000000000000000000000000000e366a1170000000000000000000000000000000000000000000000000000000077e4802600000000000000000000000000000000000000000000000000000000f917ffea0000000000000000000000000000000000000000000000000000000057ecfd28000000000000000000000000000000000000000000000000000000009d228d69b5fdb8d273a2336f8fb8612d039631024ea9bf09c424a9503aa078f0bf969f220000000000000000000000000000000000000000000000000000000024eb47e50000000000000000000000000000000000000000000000000000000055534443546f6b656e506f6f6c20312e352e31000000000000000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffff01ffc9a7000000000000000000000000000000000000000000000000000000000e64dd2900000000000000000000000000000000000000000000000000000000aff2afbf000000000000000000000000000000000000000000000000000000002b5c74de00000000000000000000000000000000000000000000000000000000ffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000000000000000000000000000000000000000000000000000000000000100000000a087bd29000000000000000000000000000000000000000000000000000000001889010d2535a0ab1643678d1da87fbbe8b87b2f585b47ddb72ec622aef9ee56ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000000000000000000000000000000000000000000000ffffffffffffff60000000000000000000000000000000000000000000000000ffffffffffffffa00000000000000000000000000000000000000000000000010000000000000000ffffffffffffffffffffff00ffffffff0000000000000000000000000000000002000000000000000000000000000000000000600000000000000000000000009ea3374b67bf275e6bb9c8ae68f9cae023e1c528b4b27e092f0bb209d3531c1902000002000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
