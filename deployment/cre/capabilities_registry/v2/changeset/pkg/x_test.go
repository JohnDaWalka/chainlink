package pkg

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
)

func TestCapabilityConfig_UnmarshalComplexJSON(t *testing.T) {
	jsonConfig := `{"method_configs":{"BalanceAt":{"remote_executable_config":{"request_timeout":{"seconds":30},"server_max_parallel_requests":10}},"CallContract":{"remote_executable_config":{"request_timeout":{"seconds":30},"server_max_parallel_requests":10}},"EstimateGas":{"remote_executable_config":{"request_timeout":{"seconds":30},"server_max_parallel_requests":10}},"FilterLogs":{"remote_executable_config":{"request_timeout":{"seconds":30},"server_max_parallel_requests":10}},"GetTransactionByHash":{"remote_executable_config":{"request_timeout":{"seconds":30},"server_max_parallel_requests":10}},"GetTransactionReceipt":{"remote_executable_config":{"request_timeout":{"seconds":30},"server_max_parallel_requests":10}},"HeaderByNumber":{"remote_executable_config":{"request_timeout":{"seconds":30},"server_max_parallel_requests":10}},"LogTrigger":{"remote_trigger_config":{"batchCollectionPeriod":{"nanos":200000000},"maxBatchSize":25,"messageExpiry":{"seconds":120},"minResponsesToAggregate":2,"registrationExpiry":{"seconds":60},"registrationRefresh":{"seconds":20}}},"WriteReport":{"remote_executable_config":{"delta_stage_nanos":38400000000,"request_hasher_type":"WriteReportExcludeSignatures","request_timeout":{"nanos":800000000,"seconds":268},"server_max_parallel_requests":10,"transmission_schedule":"OneAtATime"}}}}`

	t.Run("unmarshal_json_to_protobuf", func(t *testing.T) {
		// Convert JSON string to bytes
		jsonEncodedCfg := []byte(jsonConfig)

		// Create protobuf config and unmarshal options
		pbCfg := &pb.CapabilityConfig{}
		ops := protojson.UnmarshalOptions{DiscardUnknown: true}

		// Unmarshal JSON into protobuf
		err := ops.Unmarshal(jsonEncodedCfg, pbCfg)
		require.NoError(t, err, "failed to protojson unmarshal json encoded config")

		// Verify the protobuf was populated
		assert.NotNil(t, pbCfg)

		// Convert back to JSON to verify round-trip
		marshaledJSON, err := protojson.Marshal(pbCfg)
		require.NoError(t, err)

		// Parse the original and marshaled JSON to compare structures
		var original, roundTrip map[string]interface{}
		err = json.Unmarshal([]byte(jsonConfig), &original)
		require.NoError(t, err)

		err = json.Unmarshal(marshaledJSON, &roundTrip)
		require.NoError(t, err)

		// Verify key structures are preserved
		assert.Contains(t, roundTrip, "method_configs")
		methodConfigs := roundTrip["method_configs"].(map[string]interface{})

		// Check a few key methods exist
		assert.Contains(t, methodConfigs, "BalanceAt")
		assert.Contains(t, methodConfigs, "LogTrigger")
		assert.Contains(t, methodConfigs, "WriteReport")

		// Verify LogTrigger has remote_trigger_config
		logTrigger := methodConfigs["LogTrigger"].(map[string]interface{})
		assert.Contains(t, logTrigger, "remote_trigger_config")

		// Verify WriteReport has remote_executable_config
		writeReport := methodConfigs["WriteReport"].(map[string]interface{})
		assert.Contains(t, writeReport, "remote_executable_config")
	})

	t.Run("use_capability_config_unmarshal_proto", func(t *testing.T) {
		// Test using the CapabilityConfig type directly
		var config CapabilityConfig

		// First marshal to proto bytes using the existing MarshalProto method
		tempConfig := make(CapabilityConfig)
		err := json.Unmarshal([]byte(jsonConfig), &tempConfig)
		require.NoError(t, err)

		protoBytes, err := tempConfig.MarshalProto()
		require.NoError(t, err)

		// Then unmarshal back using UnmarshalProto
		err = config.UnmarshalProto(protoBytes)
		require.NoError(t, err)

		// Verify the config contains expected structure
		assert.Contains(t, config, "method_configs")
		methodConfigs := config["method_configs"].(map[string]interface{})
		assert.Contains(t, methodConfigs, "BalanceAt")
		assert.Contains(t, methodConfigs, "LogTrigger")
		assert.Contains(t, methodConfigs, "WriteReport")
	})
}
