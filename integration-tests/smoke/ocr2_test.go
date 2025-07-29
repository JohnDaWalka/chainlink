package smoke

import (
	"bufio"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/onsi/gomega"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/parrot"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	ctf_docker "github.com/smartcontractkit/chainlink-testing-framework/lib/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	"github.com/smartcontractkit/chainlink/v2/core/config/env"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

type ocr2test struct {
	name                string
	env                 map[string]string
	chainReaderAndCodec bool
}

func defaultTestData() ocr2test {
	return ocr2test{
		name: "n/a",
		env: map[string]string{
			string(env.EVMPlugin.Cmd): "", // not yet supported
		},
		chainReaderAndCodec: false,
	}
}

// Tests a basic OCRv2 median feed
func TestOCRv2Basic(t *testing.T) {
	t.Parallel()
	noPlugins := map[string]string{
		string(env.EVMPlugin.Cmd):    "",
		string(env.MedianPlugin.Cmd): "",
	}
	plugins := map[string]string{
		string(env.EVMPlugin.Cmd):    "", // not yet supported
		string(env.MedianPlugin.Cmd): "chainlink-feeds",
	}
	for _, test := range []ocr2test{
		{"legacy", noPlugins, false},
		{"legacy-chain-reader", noPlugins, true},
		{"plugins", plugins, false},
		{"plugins-chain-reader", plugins, true},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			l := logging.GetTestLogger(t)

			_, aggregatorContracts, sethClient, parrotClient := prepareORCv2SmokeTestEnv(t, test, l, 5)

			route := &parrot.Route{
				Method:             parrot.MethodAny,
				Path:               "/ocr2",
				ResponseBody:       10,
				ResponseStatusCode: http.StatusOK,
			}
			err := parrotClient.SetAdapterRoute(route)
			require.NoError(t, err, "Failed to set route in mock adapter")
			err = actions.WatchNewOCRRound(l, sethClient, 2, contracts.V2OffChainAgrregatorToOffChainAggregatorWithRounds(aggregatorContracts), time.Minute*5)
			require.NoError(t, err)

			roundData, err := aggregatorContracts[0].GetRound(testcontext.Get(t), big.NewInt(2))
			require.NoError(t, err, "Error getting latest OCR answer")
			require.Equal(t, int64(10), roundData.Answer.Int64(),
				"Expected latest answer from OCR contract to be 10 but got %d",
				roundData.Answer.Int64(),
			)
		})
	}
}

// Tests that just calling requestNewRound() will properly induce more rounds
func TestOCRv2Request(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)

	_, aggregatorContracts, sethClient, _ := prepareORCv2SmokeTestEnv(t, defaultTestData(), l, 5)

	// Keep the mockserver value the same and continually request new rounds
	for round := 2; round <= 4; round++ {
		err := actions.StartNewRound(contracts.V2OffChainAgrregatorToOffChainAggregatorWithRounds(aggregatorContracts))
		require.NoError(t, err, "Error starting new OCR2 round")
		err = actions.WatchNewOCRRound(l, sethClient, int64(round), contracts.V2OffChainAgrregatorToOffChainAggregatorWithRounds(aggregatorContracts), time.Minute*5)
		require.NoError(t, err, "Error watching for new OCR2 round")
		roundData, err := aggregatorContracts[0].GetRound(testcontext.Get(t), big.NewInt(int64(round)))
		require.NoError(t, err, "Getting latest answer from OCR contract shouldn't fail")
		require.Equal(t, int64(5), roundData.Answer.Int64(),
			"Expected round %d answer from OCR contract to be 5 but got %d",
			round,
			roundData.Answer.Int64(),
		)
	}
}

func TestOCRv2JobReplacement(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)

	nodeSetOutput, aggregatorContracts, sethClient, parrotClient := prepareORCv2SmokeTestEnv(t, defaultTestData(), l, 5)
	var allNodeClients []*nodeclient.ChainlinkClient
	for _, node := range nodeSetOutput.CLNodes {
		clClient, err := nodeclient.NewChainlinkClient(&nodeclient.ChainlinkConfig{
			URL:      node.Node.ExternalURL,
			Email:    node.Node.APIAuthUser,
			Password: node.Node.APIAuthPassword,
		}, l)
		require.NoError(t, err, "Error creating chainlink client")
		allNodeClients = append(allNodeClients, clClient)
	}

	workerNodeClients := allNodeClients[1:]
	bootstrapNodeClient := allNodeClients[0]

	route := &parrot.Route{
		Method:             parrot.MethodAny,
		Path:               "/ocr2",
		ResponseBody:       10,
		ResponseStatusCode: http.StatusOK,
	}
	err := parrotClient.SetAdapterRoute(route)
	require.NoError(t, err, "Failed to set route in mock adapter")
	err = actions.WatchNewOCRRound(l, sethClient, 2, contracts.V2OffChainAgrregatorToOffChainAggregatorWithRounds(aggregatorContracts), time.Minute*5)
	require.NoError(t, err, "Error watching for new OCR2 round")

	roundData, err := aggregatorContracts[0].GetRound(testcontext.Get(t), big.NewInt(2))
	require.NoError(t, err, "Error getting latest OCR answer")
	require.Equal(t, int64(10), roundData.Answer.Int64(),
		"Expected latest answer from OCR contract to be 10 but got %d",
		roundData.Answer.Int64(),
	)

	err = actions.DeleteJobs(allNodeClients)
	require.NoError(t, err)

	err = actions.DeleteBridges(allNodeClients)
	require.NoError(t, err)

	route.ResponseBody = 15
	require.NoError(t, err, "Failed to set route in mock adapter")
	require.GreaterOrEqual(t, sethClient.ChainID, int64(0), "Chain ID should be greater than or equal to 0")
	err = actions.CreateOCRv2JobsLocal(
		aggregatorContracts,
		bootstrapNodeClient,
		workerNodeClients,
		parrotClient,
		route,
		uint64(sethClient.ChainID), //nolint:gosec // Conversion from int64 to uint64 is safe
		false,
		false,
	)
	require.NoError(t, err, "Error creating OCRv2 jobs")

	err = actions.WatchNewOCRRound(l, sethClient, 3, contracts.V2OffChainAgrregatorToOffChainAggregatorWithRounds(aggregatorContracts), time.Minute*3)
	require.NoError(t, err, "Error watching for new OCR2 round")

	roundData, err = aggregatorContracts[0].GetRound(testcontext.Get(t), big.NewInt(3))
	require.NoError(t, err, "Getting latest answer from OCR contract shouldn't fail")
	require.Equal(t, int64(15), roundData.Answer.Int64(),
		"Expected latest answer from OCR contract to be 15 but got %d",
		roundData.Answer.Int64(),
	)
}

type ocr2Config struct {
	Blockchain *blockchain.Input `toml:"blockchain"`
	NodeSets   []*ns.Input       `toml:"nodesets" validate:"required"`
}

type ocr2contractConfigMock struct {
	ocr2Config
}

func (o ocr2contractConfigMock) UseExistingOffChainAggregatorsContracts() bool {
	return false
}

func (o ocr2contractConfigMock) ConfigureExistingOffChainAggregatorsContracts() bool {
	return false
}

func (o ocr2contractConfigMock) NumberOfContractsToDeploy() int {
	return 1
}

func (o ocr2contractConfigMock) OffChainAggregatorsContractsAddresses() []common.Address {
	return []common.Address{}
}

func prepareORCv2SmokeTestEnv(t *testing.T, testData ocr2test, l zerolog.Logger, firstRoundResult int) (*ns.Output, []contracts.OffchainAggregatorV2, *seth.Client, *test_env.Parrot) {
	setErr := os.Setenv("CTF_CONFIGS", "ocr2.toml")
	require.NoError(t, setErr, "Error setting CTF_CONFIGS")

	conf, confErr := framework.Load[ocr2Config](t)
	require.NoError(t, confErr, "Error loading config")

	for idx, nodeSpec := range conf.NodeSets[0].NodeSpecs {
		nodeSpec.Node.EnvVars = testData.env
		conf.NodeSets[0].NodeSpecs[idx] = nodeSpec
	}

	bcOut, bcOutErr := blockchain.NewBlockchainNetwork(conf.Blockchain)
	require.NoError(t, bcOutErr, "Error creating blockchain network")

	sethClient, sethErr := seth.NewClientBuilder().
		WithRpcUrl(bcOut.Nodes[0].ExternalWSUrl).
		WithPrivateKeys([]string{"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}).
		WithProtections(false, false, seth.MustMakeDuration(time.Second)).
		Build()
	require.NoError(t, sethErr, "Error creating seth client")

	nodeSetOutput, nodesetErr := ns.NewSharedDBNodeSet(conf.NodeSets[0], bcOut)
	require.NoError(t, nodesetErr, "Error creating node set")

	var allNodeClients []*nodeclient.ChainlinkClient
	for _, node := range nodeSetOutput.CLNodes {
		clClient, err := nodeclient.NewChainlinkClient(&nodeclient.ChainlinkConfig{
			URL:      node.Node.ExternalURL,
			Email:    node.Node.APIAuthUser,
			Password: node.Node.APIAuthPassword,
		}, l)
		require.NoError(t, err, "Error creating chainlink client")
		allNodeClients = append(allNodeClients, clClient)
	}

	workerNodeClients := allNodeClients[1:]
	bootstrapNodeClient := allNodeClients[0]

	linkContract, linkContractErr := contracts.DeployLinkTokenContract(l, sethClient)
	require.NoError(t, linkContractErr, "Error deploying link token contract")

	fundErr := actions.FundChainlinkNodesFromRootAddress(l, sethClient, contracts.ChainlinkClientToChainlinkNodeWithKeysAndAddress(workerNodeClients), big.NewFloat(0.5))
	require.NoError(t, fundErr, "Error funding Chainlink nodes")

	t.Cleanup(func() {
		// ignore error, we will see failures in the logs anyway
		_ = actions.ReturnFundsFromNodes(l, sethClient, contracts.ChainlinkClientToChainlinkNodeWithKeysAndAddress(allNodeClients))
	})

	// Gather transmitters
	var transmitters []string
	for _, node := range workerNodeClients {
		addr, err := node.PrimaryEthAddress()
		if err != nil {
			require.NoError(t, fmt.Errorf("error getting node's primary ETH address: %w", err))
		}
		transmitters = append(transmitters, addr)
	}

	ocrOffChainOptions := contracts.DefaultOffChainAggregatorOptions()
	aggregatorContracts, err := actions.SetupOCRv2Contracts(l, sethClient, ocr2contractConfigMock{}, common.HexToAddress(linkContract.Address()), transmitters, ocrOffChainOptions)
	require.NoError(t, err, "Error deploying OCRv2 aggregator contracts")

	if sethClient.ChainID < 0 {
		t.Errorf("negative chain ID: %d", sethClient.ChainID)
	}
	ocrRoute := &parrot.Route{
		Method:             parrot.MethodAny,
		Path:               "/ocr2",
		ResponseBody:       firstRoundResult,
		ResponseStatusCode: http.StatusOK,
	}
	require.GreaterOrEqual(t, sethClient.ChainID, int64(0), "Chain ID should be greater than or equal to 0")

	// TODO: move to CTFv2
	parrot := test_env.NewParrot([]string{framework.DefaultNetworkName})
	pErr := parrot.StartContainer()
	require.NoError(t, pErr, "Error starting parrot")

	err = actions.CreateOCRv2JobsLocal(
		aggregatorContracts,
		bootstrapNodeClient,
		workerNodeClients,
		parrot,
		ocrRoute,
		uint64(sethClient.ChainID), //nolint:gosec // Conversion from int64 to uint64 is safe
		false,
		testData.chainReaderAndCodec,
	)
	require.NoError(t, err, "Error creating OCRv2 jobs")

	// if !config.OCR2.UseExistingOffChainAggregatorsContracts() || (config.OCR2.UseExistingOffChainAggregatorsContracts() && config.OCR2.ConfigureExistingOffChainAggregatorsContracts()) {
	ocrV2Config, err := actions.BuildMedianOCR2ConfigLocal(workerNodeClients, ocrOffChainOptions)
	require.NoError(t, err, "Error building OCRv2 config")

	err = actions.ConfigureOCRv2AggregatorContracts(ocrV2Config, aggregatorContracts)
	require.NoError(t, err, "Error configuring OCRv2 aggregator contracts")
	// }

	assertCorrectNodeConfiguration(t, l, testData, nodeSetOutput)

	err = actions.WatchNewOCRRound(l, sethClient, 1, contracts.V2OffChainAgrregatorToOffChainAggregatorWithRounds(aggregatorContracts), time.Minute*5)
	require.NoError(t, err, "Error watching for new OCR2 round")
	roundData, err := aggregatorContracts[0].GetRound(testcontext.Get(t), big.NewInt(1))
	require.NoError(t, err, "Getting latest answer from OCR contract shouldn't fail")
	require.Equal(t, int64(firstRoundResult), roundData.Answer.Int64(),
		"Expected latest answer from OCR contract to be 5 but got %d",
		roundData.Answer.Int64(),
	)

	return nodeSetOutput, aggregatorContracts, sethClient, parrot
}

func assertCorrectNodeConfiguration(t *testing.T, l zerolog.Logger, testData ocr2test, nodeSetOutput *ns.Output) {
	l.Info().Msg("Checking if all nodes have correct plugin configuration applied")

	// we have to use gomega here, because sometimes there's a delay in the logs being written (especially in the CI)
	// and this check fails on the first execution, and we don't want to add any hardcoded sleeps

	gom := gomega.NewGomegaWithT(t)
	gom.Eventually(func(g gomega.Gomega) {
		allNodesHaveCorrectConfig := false

		var expectedPatterns []string
		expectedNodeCount := len(nodeSetOutput.CLNodes) - 1

		if testData.env[string(env.MedianPlugin.Cmd)] != "" {
			expectedPatterns = append(expectedPatterns, `Registered loopp.*OCR2.*Median.*`)
		}

		if testData.chainReaderAndCodec {
			expectedPatterns = append(expectedPatterns, `relayConfig.chainReader`)
		} else {
			expectedPatterns = append(expectedPatterns, "ChainReader missing from RelayConfig; falling back to internal MedianContract")
		}

		logFilePaths := make(map[string]string)
		tempLogsDir := os.TempDir()

		var nodesToInclude []string
		for i := 1; i < len(nodeSetOutput.CLNodes); i++ {
			nodesToInclude = append(nodesToInclude, nodeSetOutput.CLNodes[i].Node.ContainerName+".log")
		}

		// save all log files in temp dir
		loggingErr := ctf_docker.WriteAllContainersLogs(l, tempLogsDir)
		if loggingErr != nil {
			l.Debug().Err(loggingErr).Msg("Error writing all containers logs. Trying again...")

			// try again
			return
		}

		var fileNameIncludeFilter = func(name string) bool {
			for _, n := range nodesToInclude {
				if strings.EqualFold(name, n) {
					return true
				}
			}
			return false
		}

		// find log files for CL nodes
		fileWalkErr := filepath.Walk(tempLogsDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				if os.IsPermission(err) {
					return nil
				}
				return err
			}
			if !info.IsDir() && fileNameIncludeFilter(info.Name()) {
				absPath, err := filepath.Abs(path)
				if err != nil {
					return err
				}
				logFilePaths[strings.TrimSuffix(info.Name(), ".log")] = absPath
			}
			return nil
		})

		if fileWalkErr != nil {
			l.Debug().Err(fileWalkErr).Msg("Error walking through log files. Trying again...")

			return
		}

		if len(logFilePaths) != expectedNodeCount {
			l.Debug().Msgf("Expected number of log files to match number of nodes (excluding bootstrap node). Expected: %d, Found: %d. Trying again...", expectedNodeCount, len(logFilePaths))

			return
		}

		// search for expected pattern in log file
		var searchForLineInFile = func(filePath string, pattern string) bool {
			file, fileErr := os.Open(filePath)
			if fileErr != nil {
				return false
			}

			defer func(file *os.File) {
				_ = file.Close()
			}(file)

			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)
			pc := regexp.MustCompile(pattern)

			for scanner.Scan() {
				jsonLogLine := scanner.Text()
				if pc.MatchString(jsonLogLine) {
					return true
				}
			}
			return false
		}

		wg := sync.WaitGroup{}
		resultsCh := make(chan map[string][]string, len(logFilePaths))

		// process all logs in parallel
		for nodeName, logFilePath := range logFilePaths {
			wg.Add(1)
			filePath := logFilePath
			go func() {
				defer wg.Done()
				var patternsFound []string
				for _, pattern := range expectedPatterns {
					found := searchForLineInFile(filePath, pattern)
					if found {
						patternsFound = append(patternsFound, pattern)
					}
				}
				resultsCh <- map[string][]string{nodeName: patternsFound}
			}()
		}

		wg.Wait()
		close(resultsCh)

		var correctlyConfiguredNodes []string
		var incorrectlyConfiguredNodes []string

		// check results
		for result := range resultsCh {
			for nodeName, patternsFound := range result {
				if len(patternsFound) == len(expectedPatterns) {
					correctlyConfiguredNodes = append(correctlyConfiguredNodes, nodeName)
				} else {
					incorrectlyConfiguredNodes = append(incorrectlyConfiguredNodes, nodeName)
				}
			}
		}

		allNodesHaveCorrectConfig = len(correctlyConfiguredNodes) == expectedNodeCount

		g.Expect(allNodesHaveCorrectConfig).To(gomega.BeTrue(), "%d nodes' logs were missing expected plugin configuration entries. Correctly configured nodes: %s. Nodes with missing configuration: %s. Expected log patterns: %s", expectedNodeCount-len(correctlyConfiguredNodes), strings.Join(correctlyConfiguredNodes, ", "), strings.Join(incorrectlyConfiguredNodes, ", "), strings.Join(expectedPatterns, ", "))
	}, "1m", "10s").Should(gomega.Succeed())

	l.Info().Msg("All nodes have correct plugin configuration applied")
}
