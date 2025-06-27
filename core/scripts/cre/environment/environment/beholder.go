package environment

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/BurntSushi/toml"
	"github.com/google/go-github/v72/github"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	chipingressset "github.com/smartcontractkit/chainlink-testing-framework/framework/components/dockercompose/chip_ingress_set"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var (
	withBeholderFlag              bool
	protoConfigsFlag              []string
	redPandaKafkaURLFlag          string
	redPandaSchemaRegistryURLFlag string
	kafkaCreateTopicsFlag         []string
	kafkaRemoveTopicsFlag         bool
)

type ChipIngressConfig struct {
	ChipIngress *chipingressset.Input `toml:"chip_ingress"`
	Kafka       *KafkaConfig          `toml:"kafka"`
}

type KafkaConfig struct {
	Topics []string `toml:"topics"`
}

var startBeholderCmd = &cobra.Command{
	Use:   "start-beholder",
	Short: "Start the Beholder",
	Long:  `Start the Beholder`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if topologyFlag != TopologySimplified && topologyFlag != TopologyFull {
			return fmt.Errorf("invalid topology: %s. Valid topologies are: %s, %s", topologyFlag, TopologySimplified, TopologyFull)
		}

		// set TESTCONTAINERS_RYUK_DISABLED to true to disable Ryuk, so that Ryuk doesn't destroy the containers, when the command ends
		setErr := os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
		if setErr != nil {
			return fmt.Errorf("failed to set TESTCONTAINERS_RYUK_DISABLED environment variable: %w", setErr)
		}

		startBeholderErr := startBeholder(cmd.Context(), protoConfigsFlag)
		if startBeholderErr != nil {
			waitOnErrorTimeoutDurationFn()
			beholderRemoveErr := framework.RemoveTestStack(chipingressset.DEFAULT_STACK_NAME)
			if beholderRemoveErr != nil {
				fmt.Fprint(os.Stderr, errors.Wrap(beholderRemoveErr, manualBeholderCleanupMsg).Error())
			}
			return errors.Wrap(startBeholderErr, "failed to start Beholder")
		}

		return nil
	},
}

func startBeholder(cmdContext context.Context, protoConfigsFlag []string) (startupErr error) {
	// just in case, remove the stack if it exists
	_ = framework.RemoveTestStack(chipingressset.DEFAULT_STACK_NAME)

	defer func() {
		p := recover()

		if p != nil {
			fmt.Println("Panicked when starting Beholder")

			if err, ok := p.(error); ok {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				fmt.Fprintf(os.Stderr, "Stack trace: %s\n", string(debug.Stack()))

				startupErr = err
			} else {
				fmt.Fprintf(os.Stderr, "panic: %v\n", p)
				fmt.Fprintf(os.Stderr, "Stack trace: %s\n", string(debug.Stack()))

				startupErr = fmt.Errorf("panic: %v", p)
			}

			waitOnErrorTimeoutDurationFn()

			beholderRemoveErr := framework.RemoveTestStack(chipingressset.DEFAULT_STACK_NAME)
			if beholderRemoveErr != nil {
				fmt.Fprint(os.Stderr, errors.Wrap(beholderRemoveErr, manualBeholderCleanupMsg).Error())
			}
		}
	}()

	setErr := os.Setenv("CTF_CONFIGS", "configs/chip-ingress.toml")
	if setErr != nil {
		return fmt.Errorf("failed to set CTF_CONFIGS environment variable: %w", setErr)
	}

	// Load and validate test configuration
	in, err := framework.Load[ChipIngressConfig](nil)
	if err != nil {
		return errors.Wrap(err, "failed to load test configuration")
	}

	out, startErr := chipingressset.New(in.ChipIngress)
	if startErr != nil {
		return errors.Wrap(startErr, "failed to create Chip Ingress set")
	}

	fmt.Println()
	framework.L.Info().Msgf("Red Panda Console URL: %s", out.RedPanda.ConsoleExternalURL)

	topicsErr := chipingressset.CreateTopics(cmdContext, out.RedPanda.KafkaExternalURL, in.Kafka.Topics)
	if topicsErr != nil {
		return errors.Wrap(topicsErr, "failed to create topics")
	}

	for _, topic := range in.Kafka.Topics {
		framework.L.Info().Msgf("Topic URL: %s", fmt.Sprintf("%s/topics/%s", out.RedPanda.ConsoleExternalURL, topic))
	}
	fmt.Println()

	return parseConfigsAndRegisterProtos(cmdContext, protoConfigsFlag, out.RedPanda.SchemaRegistryExternalURL)
}

func parseConfigsAndRegisterProtos(ctx context.Context, protoConfigsFlag []string, schemaRegistryExternalURL string) error {
	var protoSchemaSets []chipingressset.ProtoSchemaSet
	for _, protoConfig := range protoConfigsFlag {
		file, fileErr := os.ReadFile(protoConfig)
		if fileErr != nil {
			return errors.Wrapf(fileErr, "failed to read proto config file: %s", protoConfig)
		}

		type wrappedProtoSchemaSets struct {
			ProtoSchemaSets []chipingressset.ProtoSchemaSet `toml:"proto_schema_sets"`
		}

		var schemaSets wrappedProtoSchemaSets
		if err := toml.Unmarshal(file, &schemaSets); err != nil {
			return errors.Wrapf(err, "failed to unmarshal proto config file: %s", protoConfig)
		}

		protoSchemaSets = append(protoSchemaSets, schemaSets.ProtoSchemaSets...)
	}

	if len(protoSchemaSets) == 0 {
		framework.L.Warn().Msg("no proto configs provided, skipping proto registration")

		return nil
	}

	for _, protoSchemaSet := range protoSchemaSets {
		framework.L.Info().Msgf("Registering and fetching proto from %s", protoSchemaSet.Repository)
		framework.L.Info().Msgf("Proto schema set config: %+v", protoSchemaSet)
	}

	var client *github.Client
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	} else {
		framework.L.Warn().Msg("GITHUB_TOKEN is not set, using unauthenticated GitHub client. This may cause rate limiting issues when downloading proto files")
		client = github.NewClient(nil)
	}

	reposErr := chipingressset.DefaultRegisterAndFetchProtos(ctx, client, protoSchemaSets, schemaRegistryExternalURL)
	if reposErr != nil {
		return errors.Wrap(reposErr, "failed to fetch and register protos")
	}
	return nil
}

var createKafkaTopicsCmd = &cobra.Command{
	Use:   "create-kafka-topics",
	Short: "Create Kafka topics",
	Long:  `Create Kafka topics (with or without removing existing topics)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if redPandaKafkaURLFlag == "" {
			return fmt.Errorf("red-panda-kafka-url cannot be empty")
		}

		if len(kafkaCreateTopicsFlag) == 0 {
			return fmt.Errorf("kafka topics list cannot be empty")
		}

		if kafkaRemoveTopicsFlag {
			topicsErr := chipingressset.DeleteAllTopics(cmd.Context(), redPandaKafkaURLFlag)
			if topicsErr != nil {
				return errors.Wrap(topicsErr, "failed to remove topics")
			}
		}

		topicsErr := chipingressset.CreateTopics(cmd.Context(), redPandaKafkaURLFlag, kafkaCreateTopicsFlag)
		if topicsErr != nil {
			return errors.Wrap(topicsErr, "failed to create topics")
		}

		return nil
	},
}

var fetchAndRegisterProtosCmd = &cobra.Command{
	Use:   "fetch-and-register-protos",
	Short: "Fetch and register protos",
	Long:  `Fetch and register protos`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if redPandaSchemaRegistryURLFlag == "" {
			return fmt.Errorf("red-panda-schema-registry-url cannot be empty")
		}

		if len(protoConfigsFlag) == 0 {
			framework.L.Warn().Msg("no proto configs provided, skipping proto registration")

			return nil
		}

		return parseConfigsAndRegisterProtos(cmd.Context(), protoConfigsFlag, redPandaSchemaRegistryURLFlag)
	},
}
