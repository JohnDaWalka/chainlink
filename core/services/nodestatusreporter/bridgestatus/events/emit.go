package events

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
)

// EmitBridgeStatusEvent emits a Bridge Status event through the provided custmsg.MessageEmitter
func EmitBridgeStatusEvent(ctx context.Context, emitter custmsg.MessageEmitter, event *BridgeStatusEvent) error {
	if event.Timestamp == "" {
		event.Timestamp = time.Now().Format(time.RFC3339Nano)
	}

	eventBytes, err := proto.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal BridgeStatusEvent: %w", err)
	}

	enrichedEmitter := emitter.With(
		"beholder_data_schema", SchemaBridgeStatus,
		"beholder_domain", "platform",
		"beholder_entity", fmt.Sprintf("%s.%s", ProtoPkg, BridgeStatusEventEntity),
	)

	return enrichedEmitter.Emit(ctx, string(eventBytes))
}
