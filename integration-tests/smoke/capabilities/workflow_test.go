package capabilities_test

import (
	"bufio"
	"bytes"
	"cmp"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-yaml/yaml"
	"github.com/google/go-github/v41/github"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"

	pkgworkflows "github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/feeds_consumer"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/workflow/generated/workflow_registry_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	workflow_registry_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/workflowregistry"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
)

// Copying this to avoid dependency on the core repo
func GetChainType(chainType string) (uint8, error) {
	switch chainType {
	case "evm":
		return 1, nil
	// case Solana:
	// 	return 2, nil
	// case Cosmos:
	// 	return 3, nil
	// case StarkNet:
	// 	return 4, nil
	// case Aptos:
	// 	return 5, nil
	default:
		return 0, fmt.Errorf("unexpected chaintype.ChainType: %#v", chainType)
	}
}

// Copying this to avoid dependency on the core repo
func MarshalMultichainPublicKey(ost map[string]types.OnchainPublicKey) (types.OnchainPublicKey, error) {
	pubKeys := make([][]byte, 0, len(ost))
	for k, pubKey := range ost {
		typ, err := GetChainType(k)
		if err != nil {
			// skipping unknown key type
			continue
		}
		buf := new(bytes.Buffer)
		if err = binary.Write(buf, binary.LittleEndian, typ); err != nil {
			return nil, err
		}
		length := len(pubKey)
		if length < 0 || length > math.MaxUint16 {
			return nil, errors.New("pubKey doesn't fit into uint16")
		}
		if err = binary.Write(buf, binary.LittleEndian, uint16(length)); err != nil {
			return nil, err
		}
		_, _ = buf.Write(pubKey)
		pubKeys = append(pubKeys, buf.Bytes())
	}
	// sort keys based on encoded type to make encoding deterministic
	slices.SortFunc(pubKeys, func(a, b []byte) int { return cmp.Compare(a[0], b[0]) })
	return bytes.Join(pubKeys, nil), nil
}

type WorkflowConfig struct {
	UseChainlinkCLI bool                    `toml:"use_chainlink_cli"`
	ChainlinkCLI    *ChainlinkCLIConfig     `toml:"chainlink_cli"`
	UseExising      bool                    `toml:"use_existing"`
	Existing        *ExistingWorkflowConfig `toml:"existing"`
}

type ExistingWorkflowConfig struct {
	BinaryURL string `toml:"binary_url"`
	ConfigURL string `toml:"config_url"`
}

type ChainlinkCLIConfig struct {
	FolderLocation *string `toml:"folder_location"`
}

type WorkflowTestConfig struct {
	BlockchainA    *blockchain.Input `toml:"blockchain_a" validate:"required"`
	NodeSet        *ns.Input         `toml:"nodeset" validate:"required"`
	WorkflowConfig *WorkflowConfig   `toml:"workflow_config" validate:"required"`
	JD             *jd.Input         `toml:"jd" validate:"required"`
}

type OCR3Config struct {
	Signers               [][]byte
	Transmitters          []common.Address
	F                     uint8
	OnchainConfig         []byte
	OffchainConfigVersion uint64
	OffchainConfig        []byte
}

type NodeInfo struct {
	OcrKeyBundleID            string
	TransmitterAddress        string
	PeerID                    string
	Signer                    common.Address
	OffchainPublicKey         [32]byte
	OnchainPublicKey          types.OnchainPublicKey
	ConfigEncryptionPublicKey [32]byte
}

func extractKey(value string) string {
	parts := strings.Split(value, "_")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return value
}

func downloadGHAssetFromLatestRelease(owner, repository, releaseType, assetName, ghToken string) ([]byte, error) {
	var content []byte
	if ghToken == "" {
		return content, errors.New("no github token provided")
	}

	if (releaseType == test_env.AUTOMATIC_LATEST_TAG) || (releaseType == test_env.AUTOMATIC_STABLE_LATEST_TAG) {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: ghToken},
		)
		tc := oauth2.NewClient(ctx, ts)

		ghClient := github.NewClient(tc)

		latestTags, _, err := ghClient.Repositories.ListReleases(context.Background(), owner, repository, &github.ListOptions{PerPage: 20})
		if err != nil {
			return content, errors.Wrapf(err, "failed to list releases for %s", repository)
		}

		var latestRelease *github.RepositoryRelease
		for _, tag := range latestTags {
			if releaseType == test_env.AUTOMATIC_STABLE_LATEST_TAG {
				if tag.Prerelease != nil && *tag.Prerelease {
					continue
				}
				if tag.Draft != nil && *tag.Draft {
					continue
				}
			}
			if tag.TagName != nil {
				latestRelease = tag
				break
			}
		}

		if latestRelease == nil {
			return content, errors.New("failed to find latest release with automatic tag: " + releaseType)
		}

		var assetID int64
		for _, asset := range latestRelease.Assets {
			if strings.Contains(asset.GetName(), assetName) {
				assetID = asset.GetID()
				break
			}
		}

		if assetID == 0 {
			return content, fmt.Errorf("failed to find asset %s for %s", assetName, *latestRelease.TagName)
		}

		asset, _, err := ghClient.Repositories.DownloadReleaseAsset(context.Background(), owner, repository, assetID, tc)
		if err != nil {
			return content, errors.Wrapf(err, "failed to download asset %s for %s", assetName, *latestRelease.TagName)
		}

		content, err = io.ReadAll(asset)
		if err != nil {
			return content, err
		}

		return content, nil
	}

	return content, errors.New("no automatic tag provided")
}

func getNodesInfo(
	t *testing.T,
	nodes []*clclient.ChainlinkClient,
) (nodesInfo []NodeInfo) {
	nodesInfo = make([]NodeInfo, len(nodes))

	for i, node := range nodes {
		// OCR Keys
		ocr2Keys, err := node.MustReadOCR2Keys()
		require.NoError(t, err)
		nodesInfo[i].OcrKeyBundleID = ocr2Keys.Data[0].ID

		firstOCR2Key := ocr2Keys.Data[0].Attributes
		nodesInfo[i].Signer = common.HexToAddress(extractKey(firstOCR2Key.OnChainPublicKey))

		pubKeys := make(map[string]types.OnchainPublicKey)
		ethOnchainPubKey, err := hex.DecodeString(extractKey(firstOCR2Key.OnChainPublicKey))
		require.NoError(t, err)
		pubKeys["evm"] = ethOnchainPubKey

		multichainPubKey, err := MarshalMultichainPublicKey(pubKeys)
		require.NoError(t, err)
		nodesInfo[i].OnchainPublicKey = multichainPubKey

		offchainPublicKeyBytes, err := hex.DecodeString(extractKey(firstOCR2Key.OffChainPublicKey))
		require.NoError(t, err)
		var offchainPublicKey [32]byte
		copy(offchainPublicKey[:], offchainPublicKeyBytes)
		nodesInfo[i].OffchainPublicKey = offchainPublicKey

		sharedSecretEncryptionPublicKeyBytes, err := hex.DecodeString(extractKey(firstOCR2Key.ConfigPublicKey))
		require.NoError(t, err)
		var sharedSecretEncryptionPublicKey [32]byte
		copy(sharedSecretEncryptionPublicKey[:], sharedSecretEncryptionPublicKeyBytes)
		nodesInfo[i].ConfigEncryptionPublicKey = sharedSecretEncryptionPublicKey

		// ETH Keys
		ethKeys, err := node.MustReadETHKeys()
		require.NoError(t, err)
		nodesInfo[i].TransmitterAddress = ethKeys.Data[0].Attributes.Address

		// P2P Keys
		p2pKeys, err := node.MustReadP2PKeys()
		require.NoError(t, err)
		nodesInfo[i].PeerID = p2pKeys.Data[0].Attributes.PeerID
	}

	return nodesInfo
}

func GenerateWorkflowIDFromStrings(owner string, name string, workflow []byte, config []byte, secretsURL string) (string, error) {
	ownerWithoutPrefix := owner
	if strings.HasPrefix(owner, "0x") {
		ownerWithoutPrefix = owner[2:]
	}

	ownerb, err := hex.DecodeString(ownerWithoutPrefix)
	if err != nil {
		return "", err
	}

	wid, err := pkgworkflows.GenerateWorkflowID(ownerb, name, workflow, config, secretsURL)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(wid[:]), nil
}

func isInstalled(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func download(url string) ([]byte, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancelFn()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return data, nil
}

func downloadAndDecode(url string) ([]byte, error) {
	data, err := download(url)
	if err != nil {
		return nil, err
	}

	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 content: %w", err)
	}

	return decoded, nil
}

type ChainlinkCliSettings struct {
	DevPlatform  DevPlatform  `yaml:"dev-platform"`
	UserWorkflow UserWorkflow `yaml:"user-workflow"`
	Logging      Logging      `yaml:"logging"`
	McmsConfig   McmsConfig   `yaml:"mcms-config"`
	Contracts    Contracts    `yaml:"contracts"`
	Rpcs         []RPC        `yaml:"rpcs"`
}

type DevPlatform struct {
	CapabilitiesRegistryAddress string `yaml:"capabilities-registry-contract-address"`
	DonID                       uint32 `yaml:"don-id"`
	WorkflowRegistryAddress     string `yaml:"workflow-registry-contract-address"`
}

type UserWorkflow struct {
	WorkflowOwnerAddress string `yaml:"workflow-owner-address"`
}

type Logging struct {
	SethConfigPath string `yaml:"seth-config-path"`
}

type McmsConfig struct {
	ProposalsDirectory string `yaml:"proposals-directory"`
}

type Contracts struct {
	ContractRegistry []ContractRegistry `yaml:"registries"`
}

type ContractRegistry struct {
	Name          string `yaml:"name"`
	Address       string `yaml:"address"`
	ChainSelector uint64 `yaml:"chain-selector"`
}

type RPC struct {
	ChainSelector uint64 `yaml:"chain-selector"`
	URL           string `yaml:"url"`
}

type PoRWorkflowConfig struct {
	FeedID          string `json:"feed_id"`
	URL             string `json:"url"`
	ConsumerAddress string `json:"consumer_address"`
}

const (
	chainlinkCliAssetFile       = "cre_v1.0.2_linux_amd64.tar.gz"
	cronCapabilityAssetFile     = "amd64_cron"
	e2eJobDistributorEnvVarName = "E2E_JD_IMAGE"
)

func downloadAndInstallChainlinkCLI(ghToken string) error {
	content, err := downloadGHAssetFromLatestRelease("smartcontractkit", "dev-platform", test_env.AUTOMATIC_LATEST_TAG, chainlinkCliAssetFile, ghToken)
	if err != nil {
		return err
	}

	tmpfile, err := os.CreateTemp("", chainlinkCliAssetFile)
	if err != nil {
		return err
	}
	defer tmpfile.Close()

	if _, err := tmpfile.Write(content); err != nil {
		return err
	}

	cmd := exec.Command("tar", "-xvf", tmpfile.Name()) // #nosec G204
	err = cmd.Run()

	if err != nil {
		return err
	}

	cmd = exec.Command("chmod", "+x", "chainlink-cli")
	err = cmd.Run()

	if err != nil {
		return err
	}

	if isInstalled := isInstalled("chainlink-cli"); !isInstalled {
		return errors.New("failed to install chainlink-cli or it is not available in the PATH")
	}

	return nil
}

func downloadCronCapability(ghToken string) (string, error) {
	content, err := downloadGHAssetFromLatestRelease("smartcontractkit", "capabilities", test_env.AUTOMATIC_LATEST_TAG, cronCapabilityAssetFile, ghToken)
	if err != nil {
		return "", err
	}

	fileName := cronCapabilityAssetFile
	file, err := os.Create(cronCapabilityAssetFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := file.Write(content); err != nil {
		return "", err
	}

	return fileName, nil
}

func validateInputsAndEnvVars(t *testing.T, testConfig *WorkflowTestConfig) {
	require.NotEmpty(t, os.Getenv("PRIVATE_KEY"), "PRIVATE_KEY env var must be set")
	if !testConfig.WorkflowConfig.UseChainlinkCLI {
		require.True(t, testConfig.WorkflowConfig.UseExising, "if you are not using chainlink-cli you must use an existing workflow")
	}

	ghToken := os.Getenv("GITHUB_API_TOKEN")
	_, err := downloadCronCapability(ghToken)
	require.NoError(t, err, "failed to download cron capability. Make sure token has content:read permissions to the capabilities repo")

	// TODO this part should ideally happen outside of the test, but due to how our reusable e2e test workflow is structured now
	// we cannot execute this part in workflow steps (it doesn't support any pre-execution hooks)
	if os.Getenv("IS_CI") == "true" {
		require.NotEmpty(t, os.Getenv(ctfconfig.E2E_TEST_CHAINLINK_IMAGE_ENV), "missing env var: "+ctfconfig.E2E_TEST_CHAINLINK_IMAGE_ENV)
		require.NotEmpty(t, os.Getenv(ctfconfig.E2E_TEST_CHAINLINK_VERSION_ENV), "missing env var: "+ctfconfig.E2E_TEST_CHAINLINK_VERSION_ENV)
		require.NotEmpty(t, os.Getenv(e2eJobDistributorEnvVarName), "missing env var: "+e2eJobDistributorEnvVarName)

		if testConfig.WorkflowConfig.UseChainlinkCLI {
			err = downloadAndInstallChainlinkCLI(ghToken)
			require.NoError(t, err, "failed to download and install chainlink-cli. Make sure token has content:read permissions to the dev-platform repo")
		}

		require.True(t, testConfig.WorkflowConfig.UseExising, "only existing workflow can be used in CI as of now due to issues with generating a gist read:write token")
	}

	if testConfig.WorkflowConfig.UseChainlinkCLI {
		require.True(t, isInstalled("chainlink-cli"), "chainlink-cli is required for this test. Please install it, add to path and run again")

		if !testConfig.WorkflowConfig.UseExising {
			require.NotEmpty(t, os.Getenv("GITHUB_API_TOKEN"), "GITHUB_API_TOKEN must be set to use chainlink-cli. It requires gist:read and gist:write permissions")
		} else {
			require.NotEmpty(t, testConfig.WorkflowConfig.ChainlinkCLI.FolderLocation, "folder_location must be set in the chainlink_cli config")
		}
	}
}

// copied from Bala's unmerged PR: https://github.com/smartcontractkit/chainlink/pull/15751
func getNodeInfo(nodeOut *ns.Output, bootstrapNodeCount int) ([]devenv.NodeInfo, error) {
	var nodeInfo []devenv.NodeInfo
	for i := 1; i <= len(nodeOut.CLNodes); i++ {
		p2pURL, err := url.Parse(nodeOut.CLNodes[i-1].Node.DockerP2PUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to parse p2p url: %w", err)
		}
		if i <= bootstrapNodeCount {
			nodeInfo = append(nodeInfo, devenv.NodeInfo{
				IsBootstrap: true,
				Name:        fmt.Sprintf("bootstrap-%d", i),
				P2PPort:     p2pURL.Port(),
				CLConfig: nodeclient.ChainlinkConfig{
					URL:        nodeOut.CLNodes[i-1].Node.HostURL,
					Email:      nodeOut.CLNodes[i-1].Node.APIAuthUser,
					Password:   nodeOut.CLNodes[i-1].Node.APIAuthPassword,
					InternalIP: nodeOut.CLNodes[i-1].Node.InternalIP,
				},
			})
		} else {
			nodeInfo = append(nodeInfo, devenv.NodeInfo{
				IsBootstrap: false,
				Name:        fmt.Sprintf("node-%d", i),
				P2PPort:     p2pURL.Port(),
				CLConfig: nodeclient.ChainlinkConfig{
					URL:        nodeOut.CLNodes[i-1].Node.HostURL,
					Email:      nodeOut.CLNodes[i-1].Node.APIAuthUser,
					Password:   nodeOut.CLNodes[i-1].Node.APIAuthPassword,
					InternalIP: nodeOut.CLNodes[i-1].Node.InternalIP,
				},
			})
		}
	}
	return nodeInfo, nil
}

func buildChainlinkDeploymentEnv(t *testing.T, jdOutput *jd.Output, nodeOutput *ns.Output, bs *blockchain.Output, sc *seth.Client) (*deployment.Environment, *devenv.DON, uint64) {
	lgr := logger.TestLogger(t)

	chainSelector, err := chainselectors.SelectorFromChainId(sc.Cfg.Network.ChainID)
	require.NoError(t, err, "failed to get chain selector for chain id %d", sc.Cfg.Network.ChainID)

	nodeInfo, err := getNodeInfo(nodeOutput, 1)
	require.NoError(t, err, "failed to get node info")

	jdConfig := devenv.JDConfig{
		GRPC:     jdOutput.HostGRPCUrl,
		WSRPC:    jdOutput.DockerWSRPCUrl,
		Creds:    insecure.NewCredentials(),
		NodeInfo: nodeInfo,
	}

	require.Len(t, bs.Nodes, 1, "expected only one node in the blockchain output")

	devenvConfig := devenv.EnvironmentConfig{
		JDConfig: jdConfig,
		Chains: []devenv.ChainConfig{
			{
				ChainID:   sc.Cfg.Network.ChainID,
				ChainName: sc.Cfg.Network.Name,
				ChainType: strings.ToUpper(bs.Family),
				WSRPCs: []devenv.CribRPCs{{
					External: bs.Nodes[0].HostWSUrl,
					Internal: bs.Nodes[0].DockerInternalWSUrl,
				}},
				HTTPRPCs: []devenv.CribRPCs{{
					External: bs.Nodes[0].HostHTTPUrl,
					Internal: bs.Nodes[0].DockerInternalHTTPUrl,
				}},
				DeployerKey: sc.NewTXOpts(seth.WithNonce(nil)), // set nonce to nil, so that it will be fetched from the chain
			},
		},
	}

	env, don, err := devenv.NewEnvironment(context.Background, lgr, devenvConfig)
	require.NoError(t, err, "failed to create environment")

	return env, don, chainSelector
}

type keystoneContracts struct {
	forwarderAddress           common.Address
	ocr3CapabilityAddress      common.Address
	capabilityRegistryAddrress common.Address
}

func deployKeystoneContracts(t *testing.T, testLogger zerolog.Logger, ctfEnv *deployment.Environment, chainSelector uint64) keystoneContracts {
	// Deploy keystone forwarder contract
	forwarderAddress := deployKeystoneForwarder(t, testLogger, ctfEnv, chainSelector)

	// Deploy OCR3 contract
	ocr3CapabilityAddress := deployOCR3(t, testLogger, ctfEnv, chainSelector)

	// Deploy capabilities registry contract
	capRegAddr := deployCapabilitiesRegistry(t, testLogger, ctfEnv, chainSelector)

	return keystoneContracts{
		forwarderAddress:           forwarderAddress,
		ocr3CapabilityAddress:      ocr3CapabilityAddress,
		capabilityRegistryAddrress: capRegAddr,
	}
}

func deployOCR3(t *testing.T, testLogger zerolog.Logger, ctfEnv *deployment.Environment, chainSelector uint64) common.Address {
	output, err := keystone_changeset.DeployOCR3(*ctfEnv, chainSelector)
	require.NoError(t, err, "failed to deploy OCR3 Capability contract")

	err = ctfEnv.ExistingAddresses.Merge(output.AddressBook)
	require.NoError(t, err, "failed to merge address book")

	addresses, err := ctfEnv.ExistingAddresses.AddressesForChain(chainSelector)
	require.NoError(t, err, "failed to get addresses for chain %d from the address book", chainSelector)

	var forwarderAddress common.Address
	for addrStr, tv := range addresses {
		if strings.Contains(tv.String(), "OCR3Capability") {
			forwarderAddress = common.HexToAddress(addrStr)
			testLogger.Info().Msgf("Deployed OCR3Capability contract at %s", forwarderAddress.Hex())
			break
		}
	}

	return forwarderAddress
}

func deployCapabilitiesRegistry(t *testing.T, testLogger zerolog.Logger, ctfEnv *deployment.Environment, chainSelector uint64) common.Address {
	output, err := keystone_changeset.DeployCapabilityRegistry(*ctfEnv, chainSelector)
	require.NoError(t, err, "failed to deploy Capabilities Registry contract")

	err = ctfEnv.ExistingAddresses.Merge(output.AddressBook)
	require.NoError(t, err, "failed to merge address book")

	addresses, err := ctfEnv.ExistingAddresses.AddressesForChain(chainSelector)
	require.NoError(t, err, "failed to get addresses for chain %d from the address book", chainSelector)

	var forwarderAddress common.Address
	for addrStr, tv := range addresses {
		if strings.Contains(tv.String(), "CapabilitiesRegistry") {
			forwarderAddress = common.HexToAddress(addrStr)
			testLogger.Info().Msgf("Deployed Capabilities Registry contract at %s", forwarderAddress.Hex())
			break
		}
	}

	return forwarderAddress
}

func deployKeystoneForwarder(t *testing.T, testLogger zerolog.Logger, ctfEnv *deployment.Environment, chainSelector uint64) common.Address {
	output, err := keystone_changeset.DeployForwarder(*ctfEnv, keystone_changeset.DeployForwarderRequest{
		ChainSelectors: []uint64{chainSelector},
	})
	require.NoError(t, err, "failed to deploy forwarder contract")

	err = ctfEnv.ExistingAddresses.Merge(output.AddressBook)
	require.NoError(t, err, "failed to merge address book")

	addresses, err := ctfEnv.ExistingAddresses.AddressesForChain(chainSelector)
	require.NoError(t, err, "failed to get addresses for chain %d from the address book", chainSelector)

	var forwarderAddress common.Address
	for addrStr, tv := range addresses {
		if strings.Contains(tv.String(), "KeystoneForwarder") {
			forwarderAddress = common.HexToAddress(addrStr)
			testLogger.Info().Msgf("Deployed KeystoneForwarder contract at %s", forwarderAddress.Hex())
			break
		}
	}

	return forwarderAddress
}

func prepareWorkflowRegistry(t *testing.T, testLogger zerolog.Logger, ctfEnv *deployment.Environment, chainSelector uint64, sc *seth.Client, donID uint32) common.Address {
	output, err := workflow_registry_changeset.Deploy(*ctfEnv, chainSelector)
	require.NoError(t, err, "failed to deploy workflow registry contract")

	err = ctfEnv.ExistingAddresses.Merge(output.AddressBook)
	require.NoError(t, err, "failed to merge address book")

	addresses, err := ctfEnv.ExistingAddresses.AddressesForChain(chainSelector)
	require.NoError(t, err, "failed to get addresses for chain %d from the address book", chainSelector)

	var workflowRegistryAddr common.Address
	for addrStr, tv := range addresses {
		if strings.Contains(tv.String(), "WorkflowRegistry") {
			workflowRegistryAddr = common.HexToAddress(addrStr)
			testLogger.Info().Msgf("Deployed WorkflowRegistry contract at %s", workflowRegistryAddr.Hex())
		}
	}

	// Configure Workflow Registry contract
	_, err = workflow_registry_changeset.UpdateAllowedDons(*ctfEnv, &workflow_registry_changeset.UpdateAllowedDonsRequest{
		RegistryChainSel: chainSelector,
		DonIDs:           []uint32{donID},
		Allowed:          true,
	})
	require.NoError(t, err, "failed to update allowed Dons")

	_, err = workflow_registry_changeset.UpdateAuthorizedAddresses(*ctfEnv, &workflow_registry_changeset.UpdateAuthorizedAddressesRequest{
		RegistryChainSel: chainSelector,
		Addresses:        []string{sc.MustGetRootKeyAddress().Hex()},
		Allowed:          true,
	})
	require.NoError(t, err, "failed to update authorized addresses")

	return workflowRegistryAddr
}

func prepareFeedsConsumer(t *testing.T, testLogger zerolog.Logger, ctfEnv *deployment.Environment, chainSelector uint64, sc *seth.Client, forwarderAddress common.Address, workflowName string) common.Address {
	output, err := keystone_changeset.DeployFeedsConsumer(*ctfEnv, &keystone_changeset.DeployFeedsConsumerRequest{
		ChainSelector: chainSelector,
	})
	require.NoError(t, err, "failed to deploy feeds_consumer contract")

	err = ctfEnv.ExistingAddresses.Merge(output.AddressBook)
	require.NoError(t, err, "failed to merge address book")

	addresses, err := ctfEnv.ExistingAddresses.AddressesForChain(chainSelector)
	require.NoError(t, err, "failed to get addresses for chain %d from the address book", chainSelector)

	var feedsConsumerAddress common.Address
	for addrStr, tv := range addresses {
		if strings.Contains(tv.String(), "FeedConsumer") {
			testLogger.Info().Msgf("Deployed FeedConsumer contract at %s", feedsConsumerAddress.Hex())
			feedsConsumerAddress = common.HexToAddress(addrStr)
			break
		}
	}

	require.NotEmpty(t, feedsConsumerAddress, "failed to find FeedConsumer address in the address book")

	// configure Keystone Feeds Consumer contract, so it can accept reports from the forwarder contract,
	// that come from our workflow that is owned by the root private key
	feedsConsumerInstance, err := feeds_consumer.NewKeystoneFeedsConsumer(feedsConsumerAddress, sc.Client)
	require.NoError(t, err, "failed to create feeds consumer instance")

	// Prepare hex-encoded and truncated workflow name
	var workflowNameBytes [10]byte
	var HashTruncateName = func(name string) string {
		// Compute SHA-256 hash of the input string
		hash := sha256.Sum256([]byte(name))

		// Encode as hex to ensure UTF8
		var hashBytes []byte = hash[:]
		resultHex := hex.EncodeToString(hashBytes)

		// Truncate to 10 bytes
		truncated := []byte(resultHex)[:10]
		return string(truncated)
	}

	truncated := HashTruncateName(workflowName)
	copy(workflowNameBytes[:], []byte(truncated))

	_, decodeErr := sc.Decode(feedsConsumerInstance.SetConfig(
		sc.NewTXOpts(),
		[]common.Address{forwarderAddress},           // allowed senders
		[]common.Address{sc.MustGetRootKeyAddress()}, // allowed workflow owners
		// here we need to use hex-encoded workflow name converted to []byte
		[][10]byte{workflowNameBytes}, // allowed workflow names
	))
	require.NoError(t, decodeErr, "failed to set config for feeds consumer")

	return feedsConsumerAddress
}

func registerWorkflowDirectly(t *testing.T, in *WorkflowTestConfig, sc *seth.Client, workflowRegistryAddr common.Address, donID uint32, workflowName string) {
	require.NotEmpty(t, in.WorkflowConfig.Existing.BinaryURL)
	workFlowData, err := downloadAndDecode(in.WorkflowConfig.Existing.BinaryURL)
	require.NoError(t, err, "failed to download and decode workflow binary")

	var configData []byte
	if in.WorkflowConfig.Existing.ConfigURL != "" {
		configData, err = download(in.WorkflowConfig.Existing.ConfigURL)
		require.NoError(t, err, "failed to download workflow config")
	}

	// use non-encoded workflow name
	workflowID, idErr := GenerateWorkflowIDFromStrings(sc.MustGetRootKeyAddress().Hex(), workflowName, workFlowData, configData, "")
	require.NoError(t, idErr, "failed to generate workflow ID")

	workflowRegistryInstance, err := workflow_registry_wrapper.NewWorkflowRegistry(workflowRegistryAddr, sc.Client)
	require.NoError(t, err, "failed to create workflow registry instance")

	// use non-encoded workflow name
	_, decodeErr := sc.Decode(workflowRegistryInstance.RegisterWorkflow(sc.NewTXOpts(), workflowName, [32]byte(common.Hex2Bytes(workflowID)), donID, uint8(0), in.WorkflowConfig.Existing.BinaryURL, in.WorkflowConfig.Existing.ConfigURL, ""))
	require.NoError(t, decodeErr, "failed to register workflow")
}

//revive:disable // ignore confusing-results
func compileWorkflowWithChainlinkCli(t *testing.T, in *WorkflowTestConfig, feedsConsumerAddress common.Address, settingsFile *os.File) (string, string) {
	feedID := "0x018BFE88407000400000000000000000"

	configFile, err := os.CreateTemp("", "config.json")
	require.NoError(t, err, "failed to create workflow config file")

	workflowConfig := PoRWorkflowConfig{
		FeedID:          feedID,
		URL:             "https://api.real-time-reserves.verinumus.io/v1/chainlink/proof-of-reserves/TrueUSD",
		ConsumerAddress: feedsConsumerAddress.Hex(),
	}

	configMarshalled, err := json.Marshal(workflowConfig)
	require.NoError(t, err, "failed to marshal workflow config")

	_, err = configFile.Write(configMarshalled)
	require.NoError(t, err, "failed to write workflow config file")

	var outputBuffer bytes.Buffer

	compileCmd := exec.Command("chainlink-cli", "workflow", "compile", "-S", settingsFile.Name(), "-c", configFile.Name(), "main.go") // #nosec G204
	compileCmd.Stdout = &outputBuffer
	compileCmd.Stderr = &outputBuffer
	compileCmd.Dir = *in.WorkflowConfig.ChainlinkCLI.FolderLocation
	err = compileCmd.Start()
	require.NoError(t, err, "failed to start compile command")

	err = compileCmd.Wait()
	fmt.Println("Compile output:\n", outputBuffer.String())

	require.NoError(t, err, "failed to wait for compile command")

	re := regexp.MustCompile(`Gist URL=([^\s]+)`)
	matches := re.FindAllStringSubmatch(outputBuffer.String(), -1)
	require.Len(t, matches, 2, "failed to find 2 gist URLs in compile output")

	ansiEscapePattern := `\x1b\[[0-9;]*m`
	re = regexp.MustCompile(ansiEscapePattern)

	workflowGistURL := re.ReplaceAllString(matches[0][1], "")
	workflowConfigURL := re.ReplaceAllString(matches[1][1], "")

	require.NotEmpty(t, workflowGistURL, "failed to find workflow gist URL")
	require.NotEmpty(t, workflowConfigURL, "failed to find workflow config gist URL")

	return workflowGistURL, workflowConfigURL
}

func preapreChainlinkCliSettingsFile(t *testing.T, sc *seth.Client, capRegAddr, workflowRegistryAddr common.Address, donID uint32, chainSelector uint64, rpcHTTPURL string) *os.File {
	// create chainlink-cli settings file
	settingsFile, err := os.CreateTemp("", ".chainlink-cli-settings.yaml")
	require.NoError(t, err, "failed to create chainlink-cli settings file")

	settings := ChainlinkCliSettings{
		DevPlatform: DevPlatform{
			CapabilitiesRegistryAddress: capRegAddr.Hex(),
			DonID:                       donID,
			WorkflowRegistryAddress:     workflowRegistryAddr.Hex(),
		},
		UserWorkflow: UserWorkflow{
			WorkflowOwnerAddress: sc.MustGetRootKeyAddress().Hex(),
		},
		Logging: Logging{},
		McmsConfig: McmsConfig{
			ProposalsDirectory: "./",
		},
		Contracts: Contracts{
			ContractRegistry: []ContractRegistry{
				{
					Name:          "CapabilitiesRegistry",
					Address:       capRegAddr.Hex(),
					ChainSelector: chainSelector,
				},
				{
					Name:          "WorkflowRegistry",
					Address:       workflowRegistryAddr.Hex(),
					ChainSelector: chainSelector,
				},
			},
		},
		Rpcs: []RPC{
			{
				ChainSelector: chainSelector,
				URL:           rpcHTTPURL,
			},
		},
	}

	settingsMarshalled, err := yaml.Marshal(settings)
	require.NoError(t, err, "failed to marshal chainlink-cli settings")

	_, err = settingsFile.Write(settingsMarshalled)
	require.NoError(t, err, "failed to write chainlink-cli settings file")

	return settingsFile
}

func registerWorkflow(t *testing.T, in *WorkflowTestConfig, sc *seth.Client, capRegAddr, workflowRegistryAddr, feedsConsumerAddress common.Address, donID uint32, chainSelector uint64, workflowName, pkey, rpcHTTPURL string) {
	// Register workflow directly using the provided binary and config URLs
	// This is a legacy solution, probably we can remove it soon
	if in.WorkflowConfig.UseExising && !in.WorkflowConfig.UseChainlinkCLI {
		registerWorkflowDirectly(t, in, sc, workflowRegistryAddr, donID, workflowName)

		return
	}

	// These two env vars are required by the chainlink-cli
	err := os.Setenv("WORKFLOW_OWNER_ADDRESS", sc.MustGetRootKeyAddress().Hex())
	require.NoError(t, err, "failed to set WORKFLOW_OWNER_ADDRESS env var")

	err = os.Setenv("ETH_PRIVATE_KEY", pkey)
	require.NoError(t, err, "failed to set ETH_PRIVATE_KEY env var")

	// create chainlink-cli settings file
	settingsFile := preapreChainlinkCliSettingsFile(t, sc, capRegAddr, workflowRegistryAddr, donID, chainSelector, rpcHTTPURL)

	var workflowGistURL string
	var workflowConfigURL string

	// compile and upload the workflow, if we are not using an existing one
	if !in.WorkflowConfig.UseExising {
		err := os.Setenv("GITHUB_API_TOKEN", os.Getenv("GITHUB_GIST_API_TOKEN"))
		require.NoError(t, err, "failed to set GITHUB_API_TOKEN env var")
		workflowGistURL, workflowConfigURL = compileWorkflowWithChainlinkCli(t, in, feedsConsumerAddress, settingsFile)
	} else {
		workflowGistURL = in.WorkflowConfig.Existing.BinaryURL
		workflowConfigURL = in.WorkflowConfig.Existing.ConfigURL
	}

	// register the workflow
	registerCmd := exec.Command("chainlink-cli", "workflow", "register", workflowName, "-b", workflowGistURL, "-c", workflowConfigURL, "-S", settingsFile.Name(), "-v")
	registerCmd.Stdout = os.Stdout
	registerCmd.Stderr = os.Stderr
	err = registerCmd.Run()
	require.NoError(t, err, "failed to register workflow using chainlink-cli")
}

func startAndFundNodes(t *testing.T, in *WorkflowTestConfig, bc *blockchain.Output, sc *seth.Client) (*ns.Output, []NodeInfo) {
	// Hack for CI that allows us to dynamically set the chainlink image and version
	// CTFv2 currently doesn't support dynamic image and version setting
	if os.Getenv("IS_CI") == "true" {
		// Due to how we pass custom env vars to reusable workflow we need to use placeholders, so first we need to resolve what's the name of the target environment variable
		// that stores chainlink version and then we can use it to resolve the image name
		image := fmt.Sprintf("%s:%s", os.Getenv(ctfconfig.E2E_TEST_CHAINLINK_IMAGE_ENV), ctfconfig.MustReadEnvVar_String(ctfconfig.E2E_TEST_CHAINLINK_VERSION_ENV))
		for _, nodeSpec := range in.NodeSet.NodeSpecs {
			nodeSpec.Node.Image = image
		}
	}

	nodeset, err := ns.NewSharedDBNodeSet(in.NodeSet, bc)
	require.NoError(t, err, "failed to deploy node set")

	nodeClients, err := clclient.New(nodeset.CLNodes)
	require.NoError(t, err, "failed to create chainlink clients")

	nodesInfo := getNodesInfo(t, nodeClients)

	// Fund all nodes
	for _, nodeInfo := range nodesInfo {
		_, err := actions.SendFunds(zerolog.Logger{}, sc, actions.FundsToSendPayload{
			ToAddress:  common.HexToAddress(nodeInfo.TransmitterAddress),
			Amount:     big.NewInt(5000000000000000000),
			PrivateKey: sc.MustGetRootPrivateKey(),
		})
		require.NoError(t, err)
	}

	return nodeset, nodesInfo
}

func configureNodes(t *testing.T, nodesInfo []NodeInfo, in *WorkflowTestConfig, bc *blockchain.Output, capRegAddr common.Address, workflowRegistryAddr common.Address, forwarderAddress common.Address) (*ns.Output, []*clclient.ChainlinkClient) {
	bootstrapNodeInfo := nodesInfo[0]
	workflowNodesetInfo := nodesInfo[1:]

	// configure the bootstrap node
	in.NodeSet.NodeSpecs[0].Node.TestConfigOverrides = fmt.Sprintf(`
				[Feature]
				LogPoller = true

				[OCR2]
				Enabled = true
				DatabaseTimeout = '1s'
				ContractPollInterval = '1s'

				[P2P.V2]
				Enabled = true
				ListenAddresses = ['0.0.0.0:5001']
				DefaultBootstrappers = ['%s@localhost:5001']

				[Capabilities.Peering.V2]
				Enabled = true
				ListenAddresses = ['0.0.0.0:6690']
				DefaultBootstrappers = ['%s@localhost:6690']

				# This is needed for the target capability to be initialized
				[[EVM]]
				ChainID = '%s'

				[[EVM.Nodes]]
				Name = 'anvil'
				WSURL = '%s'
				HTTPURL = '%s'
			`,
		bootstrapNodeInfo.PeerID,
		bootstrapNodeInfo.PeerID,
		bc.ChainID,
		bc.Nodes[0].DockerInternalWSUrl,
		bc.Nodes[0].DockerInternalHTTPUrl,
	)

	// configure worker nodes with p2p, peering capabilitity (for DON-2-DON communication),
	// capability (external) registry, workflow registry and gateway connector (required for reading from workflow registry and for external communication)
	for i := range workflowNodesetInfo {
		in.NodeSet.NodeSpecs[i+1].Node.TestConfigOverrides = fmt.Sprintf(`
				[Feature]
				LogPoller = true

				[OCR2]
				Enabled = true
				DatabaseTimeout = '1s'
				ContractPollInterval = '1s'

				[P2P.V2]
				Enabled = true
				ListenAddresses = ['0.0.0.0:5001']
				# assuming that node0 is the bootstrap node
				DefaultBootstrappers = ['%s@node0:5001']

				[Capabilities.Peering.V2]
				Enabled = true
				ListenAddresses = ['0.0.0.0:6690']
				# assuming that node0 is the bootstrap node
				DefaultBootstrappers = ['%s@node0:6690']

				# This is needed for the target capability to be initialized
				[[EVM]]
				ChainID = '%s'

				[[EVM.Nodes]]
				Name = 'anvil'
				WSURL = '%s'
				HTTPURL = '%s'

				[EVM.Workflow]
				FromAddress = '%s'
				ForwarderAddress = '%s'
				GasLimitDefault = 400_000

				[Capabilities.ExternalRegistry]
				Address = '%s'
				NetworkID = 'evm'
				ChainID = '%s'

				[Capabilities.WorkflowRegistry]
				Address = "%s"
				NetworkID = "evm"
				ChainID = "%s"

				[Capabilities.GatewayConnector]
				DonID = "1"
				ChainIDForNodeKey = "%s"
				NodeAddress = '%s'

				[[Capabilities.GatewayConnector.Gateways]]
				Id = "por_gateway"
				URL = "%s"
			`,
			bootstrapNodeInfo.PeerID,
			bootstrapNodeInfo.PeerID,
			bc.ChainID,
			bc.Nodes[0].DockerInternalWSUrl,
			bc.Nodes[0].DockerInternalHTTPUrl,
			workflowNodesetInfo[i].TransmitterAddress,
			forwarderAddress.Hex(),
			capRegAddr,
			bc.ChainID,
			workflowRegistryAddr.Hex(),
			bc.ChainID,
			bc.ChainID,
			workflowNodesetInfo[i].TransmitterAddress,
			"ws://node0:5003/node", // bootstrap node exposes gateway port on 5003
		)
	}

	// we need to restart all nodes for configuration changes to take effect
	nodeset, err := ns.UpgradeNodeSet(t, in.NodeSet, bc, 5*time.Second)
	require.NoError(t, err, "failed to upgrade node set")

	// we need to recreate chainlink clients after the nodes are restarted
	nodeClients, err := clclient.New(nodeset.CLNodes)
	require.NoError(t, err, "failed to create chainlink clients")

	return nodeset, nodeClients
}

func createNodeJobs(t *testing.T, nodeClients []*clclient.ChainlinkClient, nodesInfo []NodeInfo, bc *blockchain.Output, ocr3CapabilityAddress common.Address) {
	bootstrapNodeInfo := nodesInfo[0]
	workflowNodesetInfo := nodesInfo[1:]

	// Create gateway and bootstrap (ocr3) jobs for the bootstrap node
	bootstrapNode := nodeClients[0]
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		bootstrapJobSpec := fmt.Sprintf(`
				type = "bootstrap"
				schemaVersion = 1
				name = "Botostrap"
				contractID = "%s"
				contractConfigTrackerPollInterval = "1s"
				contractConfigConfirmations = 1
				relay = "evm"

				[relayConfig]
				chainID = %s
				providerType = "ocr3-capability"
			`, ocr3CapabilityAddress, bc.ChainID)
		r, _, bootErr := bootstrapNode.CreateJobRaw(bootstrapJobSpec)
		assert.NoError(t, bootErr, "failed to create bootstrap job for the bootstrap node")
		assert.Empty(t, r.Errors, "failed to create bootstrap job for the bootstrap node")

		gatewayJobSpec := fmt.Sprintf(`
				type = "gateway"
				schemaVersion = 1
				name = "PoR Gateway"
				forwardingAllowed = false

				[gatewayConfig.ConnectionManagerConfig]
				AuthChallengeLen = 10
				AuthGatewayId = "por_gateway"
				AuthTimestampToleranceSec = 5
				HeartbeatIntervalSec = 20

				[[gatewayConfig.Dons]]
				DonId = "1"
				F = 1
				HandlerName = "web-api-capabilities"
					[gatewayConfig.Dons.HandlerConfig]
					MaxAllowedMessageAgeSec = 1_000

						[gatewayConfig.Dons.HandlerConfig.NodeRateLimiter]
						GlobalBurst = 10
						GlobalRPS = 50
						PerSenderBurst = 10
						PerSenderRPS = 10

					[[gatewayConfig.Dons.Members]]
					Address = "%s"
					Name = "Workflow Node 1"
					[[gatewayConfig.Dons.Members]]
					Address = "%s"
					Name = "Workflow Node 2"
					[[gatewayConfig.Dons.Members]]
					Address = "%s"
					Name = "Workflow Node 3"
					[[gatewayConfig.Dons.Members]]
					Address = "%s"
					Name = "Workflow Node 4"

				[gatewayConfig.NodeServerConfig]
				HandshakeTimeoutMillis = 1_000
				MaxRequestBytes = 100_000
				Path = "/node"
				Port = 5_003 #this is the port the other nodes will use to connect to the gateway
				ReadTimeoutMillis = 1_000
				RequestTimeoutMillis = 10_000
				WriteTimeoutMillis = 1_000

				[gatewayConfig.UserServerConfig]
				ContentTypeHeader = "application/jsonrpc"
				MaxRequestBytes = 100_000
				Path = "/"
				Port = 5_002
				ReadTimeoutMillis = 1_000
				RequestTimeoutMillis = 10_000
				WriteTimeoutMillis = 1_000

				[gatewayConfig.HTTPClientConfig]
				MaxResponseBytes = 100_000_000
			`,
			// ETH keys of the workflow nodes
			workflowNodesetInfo[0].TransmitterAddress,
			workflowNodesetInfo[1].TransmitterAddress,
			workflowNodesetInfo[2].TransmitterAddress,
			workflowNodesetInfo[3].TransmitterAddress,
		)

		r, _, gatewayErr := bootstrapNode.CreateJobRaw(gatewayJobSpec)
		assert.NoError(t, gatewayErr, "failed to create gateway job for the bootstrap node")
		assert.Empty(t, r.Errors, "failed to create gateway job for the bootstrap node")
	}()

	// for each capability that's required by the workflow, create a job for workflow each node
	for i, nodeClient := range nodeClients {
		// First node is a bootstrap node, so we skip it
		if i == 0 {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			// since we are using a capability that is not bundled-in, we need to copy it to the Docker container
			// and point the job to the copied binary
			cronJobSpec := fmt.Sprintf(`
					type = "standardcapabilities"
					schemaVersion = 1
					name = "cron-capabilities"
					forwardingAllowed = false
					command = "/home/capabilities/%s"
					config = ""
				`, cronCapabilityAssetFile)

			response, _, errCron := nodeClient.CreateJobRaw(cronJobSpec)
			assert.NoError(t, errCron, "failed to create cron job")
			assert.Empty(t, response.Errors, "failed to create cron job")

			computeJobSpec := `
					type = "standardcapabilities"
					schemaVersion = 1
					name = "compute-capabilities"
					forwardingAllowed = false
					command = "__builtin_custom-compute-action"
					config = """
					NumWorkers = 3
						[rateLimiter]
						globalRPS = 20.0
						globalBurst = 30
						perSenderRPS = 1.0
						perSenderBurst = 5
					"""
				`

			response, _, errCompute := nodeClient.CreateJobRaw(computeJobSpec)
			assert.NoError(t, errCompute, "failed to create compute job")
			assert.Empty(t, response.Errors, "failed to create compute job")

			consensusJobSpec := fmt.Sprintf(`
					type = "offchainreporting2"
					schemaVersion = 1
					name = "Keystone OCR3 Consensus Capability"
					contractID = "%s"
					ocrKeyBundleID = "%s"
					p2pv2Bootstrappers = [
						"%s@%s",
					]
					relay = "evm"
					pluginType = "plugin"
					transmitterID = "%s"

					[relayConfig]
					chainID = "%s"

					[pluginConfig]
					command = "/usr/local/bin/chainlink-ocr3-capability"
					ocrVersion = 3
					pluginName = "ocr-capability"
					providerType = "ocr3-capability"
					telemetryType = "plugin"

					[onchainSigningStrategy]
					strategyName = 'multi-chain'
					[onchainSigningStrategy.config]
					evm = "%s"
					`,
				ocr3CapabilityAddress,
				nodesInfo[i].OcrKeyBundleID,
				bootstrapNodeInfo.PeerID,
				"node0:5001",
				nodesInfo[i].TransmitterAddress,
				bc.ChainID,
				nodesInfo[i].OcrKeyBundleID,
			)
			fmt.Println("consensusJobSpec", consensusJobSpec)
			response, _, errCons := nodeClient.CreateJobRaw(consensusJobSpec)
			assert.NoError(t, errCons, "failed to create consensus job")
			assert.Empty(t, response.Errors, "failed to create consensus job")
		}()
	}
	wg.Wait()
}

func configureWorkflowDON(t *testing.T, ctfEnv *deployment.Environment, don *devenv.DON, chainSelector uint64) {
	kcrAllCaps := []kcr.CapabilitiesRegistryCapability{
		{
			LabelledName:   "offchain_reporting",
			Version:        "1.0.0",
			CapabilityType: 2, // CONSENSUS
			ResponseType:   0, // REPORT
		},
		{
			LabelledName:   "write_geth-testnet",
			Version:        "1.0.0",
			CapabilityType: 3, // TARGET
			ResponseType:   1, // OBSERVATION_IDENTICAL
		},
		{
			LabelledName:   "cron-trigger",
			Version:        "1.0.0",
			CapabilityType: uint8(0), // trigger
		},
		{
			LabelledName:   "custom-compute",
			Version:        "1.0.0",
			CapabilityType: uint8(1), // action
		},
	}

	peerIds := make([]string, len(don.Nodes)-1)
	for i, node := range don.Nodes {
		if i == 0 {
			continue
		}
		peerIds[i-1] = node.PeerId
	}

	nop := keystone_changeset.NOP{
		Name:  "NOP 1",
		Nodes: peerIds,
	}

	donName := "keystone-don"
	donCap := keystone_changeset.DonCapabilities{
		Name: donName,
		F:    1,
		Nops: []keystone_changeset.NOP{nop},
		Capabilities: keystone_changeset.DONCapabilityWithConfig{
			Capability: kcrAllCaps,
		},
	}

	oracleConfig := keystone_changeset.OracleConfig{
		DeltaProgressMillis:               5000,
		DeltaResendMillis:                 5000,
		DeltaInitialMillis:                5000,
		DeltaRoundMillis:                  2000,
		DeltaGraceMillis:                  500,
		DeltaCertifiedCommitRequestMillis: 1000,
		DeltaStageMillis:                  30000,
		MaxRoundsPerEpoch:                 10,
		TransmissionSchedule:              []int{1, 2, 3, 4},
		MaxDurationQueryMillis:            1000,
		MaxDurationObservationMillis:      1000,
		MaxDurationAcceptMillis:           1000,
		MaxDurationTransmitMillis:         1000,
		MaxFaultyOracles:                  1,
	}

	cfg := keystone_changeset.InitialContractsCfg{
		RegistryChainSel: chainSelector,
		Dons:             []keystone_changeset.DonCapabilities{donCap},
		OCR3Config:       &oracleConfig,
	}

	_, err := keystone_changeset.ConfigureInitialContractsChangeset(*ctfEnv, cfg)
	require.NoError(t, err, "failed to configure initial contracts")
}

func startJobDistributor(t *testing.T, in *WorkflowTestConfig) *jd.Output {
	if os.Getenv("IS_CI") == "true" {
		jdImage := ctfconfig.MustReadEnvVar_String(e2eJobDistributorEnvVarName)
		in.JD.Image = jdImage
	}
	jdOutput, err := jd.NewJD(in.JD)
	require.NoError(t, err, "failed to create new job distributor")

	return jdOutput
}

// This function is used to go through Chainlink Node logs and look for entries related to report transmissions.
// Once such a log entry is found, it looks for transaction hash and then it tries to decode the transaction and print the result.
func debugReportTransmission(t *testing.T, l zerolog.Logger, ns *ns.Output, wsRPCURL string) {
	var logFiles []*os.File

	// when tests run in parallel, we need to make sure that we only process logs that belong to nodes created by the current test
	// that is required, because some tests might have custom log messages that are allowed, but only for that test (e.g. because they restart the CL node)
	var belongsToCurrentEnv = func(filePath string) bool {
		for _, clNode := range ns.CLNodes {
			if clNode == nil {
				continue
			}
			if strings.EqualFold(filePath, clNode.Node.ContainerName+".log") {
				return true
			}
		}
		return false
	}

	logsDir := "logs/docker-" + t.Name()

	fileWalkErr := filepath.Walk(logsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && belongsToCurrentEnv(info.Name()) {
			file, fileErr := os.Open(path)
			if fileErr != nil {
				return fmt.Errorf("failed to open file %s: %w", path, fileErr)
			}
			logFiles = append(logFiles, file)
		}
		return nil
	})

	if len(logFiles) != len(ns.CLNodes) {
		l.Warn().Int("Expected", len(ns.CLNodes)).Int("Got", len(logFiles)).Msg("Number of log files does not match number of nodes. Some logs might be missing.")
	}

	if fileWalkErr != nil {
		l.Error().Err(fileWalkErr).Msg("Error walking through log files. Will not look for report transmission transaction hashes")
		return
	}

	/*
	 Example log entry:
	 2025-01-28T14:44:48.080Z [DEBUG] Node sent transaction                              multinode@v0.0.0-20250121205514-f73e2f86c23b/transaction_sender.go:180 chainID=1337 logger=EVM.1337.TransactionSender tx={"type":"0x0","chainId":"0x539","nonce":"0x0","to":"0xcf7ed3acca5a467e9e704c703e8d87f634fb0fc9","gas":"0x61a80","gasPrice":"0x3b9aca00","maxPriorityFeePerGas":null,"maxFeePerGas":null,"value":"0x0","input":"0x11289565000000000000000000000000a513e6e4b8f2a923d98304ec87f64353c4d5c853000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000001c00000000000000000000000000000000000000000000000000000000000000240000000000000000000000000000000000000000000000000000000000000010d010f715db03509d388f706e16137722000e26aa650a64ac826ae8e5679cdf57fd96798ed50000000010000000100000a9c593aaed2f5371a5bc0779d1b8ea6f9c7d37bfcbb876a0a9444dbd36f64306466323239353031f39fd6e51aad88f6f4ce6ab8827279cfffb92266000100000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000001018bfe88407000400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000bb5c162c8000000000000000000000000000000000000000000000000000000006798ed37000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000060000e700d4c57250eac9dc925c951154c90c1b6017944322fb2075055d8bdbe19000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000041561c171b7465e8efef35572ef82adedb49ea71b8344a34a54ce5e853f80ca1ad7d644ebe710728f21ebfc3e2407bd90173244f744faa011c3a57213c8c585de90000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004165e6f3623acc43f163a58761655841bfebf3f6b4ea5f8d34c64188036b0ac23037ebbd3854b204ca26d828675395c4b9079ca068d9798326eb8c93f26570a1080100000000000000000000000000000000000000000000000000000000000000","v":"0xa96","r":"0x168547e96e7088c212f85a4e8dddce044bbb2abfd5ccc8a5451fdfcb812c94e5","s":"0x2a735a3df046632c2aaa7e583fe161113f3345002e6c9137bbfa6800a63f28a4","hash":"0x3fc5508310f8deef09a46ad594dcc5dc9ba415319ef1dfa3136335eb9e87ff4d"} version=2.19.0@05c05a9
	*/
	reportTransmissionTxHashPattern := regexp.MustCompile(`"hash":"(0x[0-9a-fA-F]+)"`)

	wg := &sync.WaitGroup{}

	// let's be prudent and assume that in extreme scenario when feed price isn't updated, but
	// transmission is still sent, we might have multiple transmissions per node, and if we want
	// to avoid blocking on the channel, we need to have a higher buffer
	resultsCh := make(chan string, len(logFiles)*4)

	// iterate overall log files looking for log entries containing "Node sent transaction" text
	// extract transaction hash from the log entry
	for _, f := range logFiles {
		wg.Add(1)
		file := f

		go func() {
			defer file.Close()
			defer wg.Done()

			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)

			for scanner.Scan() {
				jsonLogLine := scanner.Text()

				if !strings.Contains(jsonLogLine, "Node sent transaction") {
					continue
				}

				match := reportTransmissionTxHashPattern.MatchString(jsonLogLine)
				if match {
					resultsCh <- reportTransmissionTxHashPattern.FindStringSubmatch(jsonLogLine)[1]
				}
			}
		}()
	}

	wg.Wait()
	close(resultsCh)

	// required as Seth prints transaction traces to stdout with debug level
	_ = os.Setenv(seth.LogLevelEnvVar, "debug")

	sc, err := seth.NewClientBuilder().
		WithRpcUrl(wsRPCURL).
		WithReadOnlyMode().
		WithGethWrappersFolders([]string{"../../../core/gethwrappers/keystone/generated"}). // point Seth to the folder with keystone geth wrappers, so that it can load contract ABIs
		Build()

	if err != nil {
		l.Error().Err(err).Msg("Failed to create seth client")
		return
	}

	transmissionsFound := false
	for txHash := range resultsCh {
		transmissionsFound = true

		// set tracing level to all to trace also successful transactions
		sc.Cfg.TracingLevel = seth.TracingLevel_All
		tx, _, err := sc.Client.TransactionByHash(context.Background(), common.HexToHash(txHash))
		if err != nil {
			l.Warn().Err(err).Msgf("Failed to get transaction by hash %s", txHash)
			continue
		}
		_, decodedErr := sc.DecodeTx(tx)

		if decodedErr != nil {
			l.Error().Err(decodedErr).Msgf("Transmission transaction %s failed due to %s", txHash, decodedErr.Error())
			continue
		}
	}

	if !transmissionsFound {
		l.Error().Msg("No report transmissions found in Chainlink Node logs. This might be due to a bug in the node or contracts or node/job misconfiguration. Or issues with the test itself.")
	}
}

func logTestInfo(l zerolog.Logger, feedId, workflowName, feedConsumerAddr, forwarderAddr string) {
	l.Info().Msg("Test configuration:")
	l.Info().Msgf("Feed ID: %s", feedId)
	l.Info().Msgf("Workflow name: %s", workflowName)
	l.Info().Msgf("FeedConsumer address: %s", feedConsumerAddr)
	l.Info().Msgf("KeystoneForwarder address: %s", forwarderAddr)
}

/*
!!! ATTENTION !!!

Do not use this test as a template for your tests. It's hacky, since we were working under time pressure. We will soon refactor it follow best practices
and a golden example. Apart from its structure what is currently missing is:
- using Job Distribution to create jobs for the nodes
- using only `chainlink-cli` to register the workflow (it's there, but doesn't work in CI due to insufficient Github token permissions)
- using a mock service to provide the feed data
*/
func TestKeystoneWithOCR3Workflow(t *testing.T) {
	testLogger := framework.L

	// Define test configuration
	donID := uint32(1)
	workflowName := "abcdefgasd"
	feedID := "018bfe8840700040000000000000000000000000000000000000000000000000" // without 0x prefix!
	feedBytes := common.HexToHash(feedID)

	// we need to use double-pointers, so that what's captured in the cleanup function is a pointer, not the actual object,
	// which is only set later in the test, after the cleanup function is defined
	var nodes **ns.Output
	var wsRPCURL *string

	// clean up is LIFO, so we need to make sure we execute the debug report transmission after logs are written down
	// by function added to clean up by framework.Load() method.
	t.Cleanup(func() {
		if t.Failed() {
			if nodes == nil {
				testLogger.Warn().Msg("nodeset output is nil, skipping debug report transmission")
				return
			}
			// if the test fails, let's debug transactions of the report transmissions
			debugReportTransmission(t, testLogger, *nodes, *wsRPCURL)
		}
	})

	in, err := framework.Load[WorkflowTestConfig](t)
	require.NoError(t, err, "couldn't load test config")
	validateInputsAndEnvVars(t, in)

	pkey := os.Getenv("PRIVATE_KEY")

	// Create a new blockchain network
	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)

	sc, err := seth.NewClientBuilder().
		WithRpcUrl(bc.Nodes[0].HostWSUrl).
		WithPrivateKeys([]string{pkey}).
		Build()
	require.NoError(t, err, "failed to create seth client")

	// Start job distributor
	jdOutput := startJobDistributor(t, in)

	// Deploy and fund the DON
	nodeOutput, nodesInfo := startAndFundNodes(t, in, bc, sc)

	// Prepare the chainlink/deployment environment
	ctfEnv, don, chainSelector := buildChainlinkDeploymentEnv(t, jdOutput, nodeOutput, bc, sc)

	// Deploy keystone contracts
	keystoneContracts := deployKeystoneContracts(t, testLogger, ctfEnv, chainSelector)

	// Deploy and pre-configure workflow registry contract
	workflowRegistryAddr := prepareWorkflowRegistry(t, testLogger, ctfEnv, chainSelector, sc, donID)

	// Deploy and configure Keystone Feeds Consumer contract
	feedsConsumerAddress := prepareFeedsConsumer(t, testLogger, ctfEnv, chainSelector, sc, keystoneContracts.forwarderAddress, workflowName)

	// Register the workflow (either via chainlink-cli or by calling the workflow registry directly)
	registerWorkflow(t, in, sc, keystoneContracts.capabilityRegistryAddrress, workflowRegistryAddr, feedsConsumerAddress, donID, chainSelector, workflowName, pkey, bc.Nodes[0].HostHTTPUrl)

	// Create OCR3 and capability jobs for each node without JD
	ns, nodeClients := configureNodes(t, nodesInfo, in, bc, keystoneContracts.capabilityRegistryAddrress, workflowRegistryAddr, keystoneContracts.forwarderAddress)
	createNodeJobs(t, nodeClients, nodesInfo, bc, keystoneContracts.ocr3CapabilityAddress)

	// Log extra information that might help debugging
	t.Cleanup(func() {
		logTestInfo(testLogger, feedID, workflowName, feedsConsumerAddress.Hex(), keystoneContracts.forwarderAddress.Hex())
	})

	// set variables that are needed for the cleanup function, which debugs report transmissions
	nodes = &ns
	wsRPCURL = &bc.Nodes[0].HostWSUrl

	// Wait for OCR listeners to be ready before setting the configuration.
	// If the ConfigSet event is missed, OCR protocol will not start.
	// TODO make it fluent!
	testLogger.Info().Msg("Waiting 30s for OCR listeners to be ready...")
	time.Sleep(30 * time.Second)
	testLogger.Info().Msg("Proceeding to set OCR3 configuration.")

	// Configure the workflow DON
	configureWorkflowDON(t, ctfEnv, don, chainSelector)

	// It can take a while before the first report is produced, particularly on CI.
	timeout := 10 * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	feedsConsumerInstance, err := feeds_consumer.NewKeystoneFeedsConsumer(feedsConsumerAddress, sc.Client)
	require.NoError(t, err, "failed to create feeds consumer instance")

	testLogger.Info().Msg("Waiting for feed to update...")
	startTime := time.Now()

	for {
		select {
		case <-ctx.Done():
			t.Fatalf("feed did not update, timeout after %s", timeout)
		case <-time.After(10 * time.Second):
			elapsed := time.Since(startTime).Round(time.Second)
			price, _, err := feedsConsumerInstance.GetPrice(
				sc.NewCallOpts(),
				feedBytes,
			)
			require.NoError(t, err, "failed to get price from Keystone Consumer contract")

			if price.String() != "0" {
				testLogger.Info().Msgf("Feed updated after %s - price set, price=%s", elapsed, price)
				return
			}
			testLogger.Info().Msgf("Feed not updated yet, waiting for %s", elapsed)
		}
	}
}
