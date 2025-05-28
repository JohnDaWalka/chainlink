package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	billing "github.com/smartcontractkit/chainlink-protos/billing/go"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type BillingService struct {
	billing.UnimplementedWorkflowServiceServer

	lggr logger.Logger
}

func (s *BillingService) ReserveCredits(
	_ context.Context,
	request *billing.ReserveCreditsRequest,
) (*billing.ReserveCreditsResponse, error) {
	s.lggr.Infof("ReserveCredits: %v", request)

	return &billing.ReserveCreditsResponse{Success: true}, nil
}

func (s *BillingService) WorkflowReceipt(
	_ context.Context,
	request *billing.SubmitWorkflowReceiptRequest,
) (*billing.SubmitWorkflowReceiptResponse, error) {
	s.lggr.Info("WorkflowReceipt")

	return &billing.SubmitWorkflowReceiptResponse{Success: true}, nil
}

func RunBillingListener(ctx context.Context, lggr logger.Logger) {
	lis, err := net.Listen("tcp", "localhost:4319")
	if err != nil {
		log.Fatalf("billing failed to listen: %v", err)
	}

	s := grpc.NewServer()

	billing.RegisterWorkflowServiceServer(s, &BillingService{lggr: lggr})
	context.AfterFunc(ctx, s.Stop)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("billing failed to serve: %v", err)
		}
	}()
}

func setupBeholder(lggr logger.Logger) error {
	writer := &lggrWriter{lggr: lggr}

	client, err := beholder.NewWriterClient(writer)
	if err != nil {
		return err
	}

	beholder.SetClient(client)

	return nil
}

type lggrWriter struct {
	lggr logger.Logger
}

func (w lggrWriter) Write(bts []byte) (int, error) {
	w.lggr.Info(string(bts))

	return len(bts), nil
}
