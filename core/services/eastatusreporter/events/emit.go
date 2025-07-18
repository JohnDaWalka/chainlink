package events

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
)

// EmitEAStatusEvent emits an EA Status event through the provided custmsg.MessageEmitter
func EmitEAStatusEvent(ctx context.Context, emitter custmsg.MessageEmitter, event *EAStatusEvent) error {
	if event.Timestamp == "" {
		event.Timestamp = time.Now().Format(time.RFC3339Nano)
	}

	eventBytes, err := proto.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal EAStatusEvent: %w", err)
	}

	enrichedEmitter := emitter.With(
		"beholder_data_schema", SchemaEAStatus,
		"beholder_domain", "platform",
		"beholder_entity", fmt.Sprintf("%s.%s", ProtoPkg, EAStatusEventEntity),
	)

	return enrichedEmitter.Emit(ctx, string(eventBytes))
}
