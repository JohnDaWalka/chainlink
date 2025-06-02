package metering

import (
	"context"
	"strconv"
	"sync"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	billing "github.com/smartcontractkit/chainlink-protos/billing/go"
)

const (
	testAccountID           = "accountId"
	testWorkflowID          = "workflowId"
	testWorkflowExecutionID = "workflowExecutionId"
)

type mockBillingClient struct{}

func (m *mockBillingClient) SubmitWorkflowReceipt(context.Context, *billing.SubmitWorkflowReceiptRequest) (*billing.SubmitWorkflowReceiptResponse, error) {
	return &billing.SubmitWorkflowReceiptResponse{Success: true}, nil
}
func (m *mockBillingClient) ReserveCredits(context.Context, *billing.ReserveCreditsRequest) (*billing.ReserveCreditsResponse, error) {
	return &billing.ReserveCreditsResponse{Success: true}, nil
}
func newMockBillingClient() BillingClient {
	return &mockBillingClient{}
}

func TestReport(t *testing.T) {
	t.Parallel()

	testA := "a"
	testUnitA := SpendUnit(testA)
	testB := "b"
	testUnitB := SpendUnit(testB)

	t.Run("MedianSpend returns median for multiple spend units", func(t *testing.T) {
		t.Parallel()

		billingClient := newMockBillingClient()
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t))
		report.client = billingClient
		err := report.Initialize(t.Context())
		require.NoError(t, err)
		err = report.balance.Add(100)
		require.NoError(t, err)

		steps := []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "abc", SpendUnit: testA, SpendValue: "1"},
			{Peer2PeerID: "xyz", SpendUnit: testA, SpendValue: "2"},
			{Peer2PeerID: "abc", SpendUnit: testA, SpendValue: "3"},
			{Peer2PeerID: "abc", SpendUnit: testB, SpendValue: "0.1"},
			{Peer2PeerID: "xyz", SpendUnit: testB, SpendValue: "0.2"},
			{Peer2PeerID: "abc", SpendUnit: testB, SpendValue: "0.3"},
		}

		for idx := range steps {
			_, err := report.ReserveByLimits(strconv.Itoa(idx), capabilities.CapabilityInfo{}, []SpendTuple{{Value: 1, Unit: "SomeUnit"}})
			require.NoError(t, err)
			require.NoError(t, report.SetStep(strconv.Itoa(idx), steps))
		}

		expected := map[SpendUnit]SpendValue{
			testUnitA: testUnitB.IntToSpendValue(2),
			testUnitB: testUnitB.DecimalToSpendValue(decimal.NewFromFloat(0.2)),
		}

		median := report.MedianSpend()

		require.Len(t, median, 2)
		require.Contains(t, maps.Keys(median), testUnitA)
		require.Contains(t, maps.Keys(median), testUnitB)

		assert.Equal(t, expected[testUnitA].String(), median[testUnitA].String())
		assert.Equal(t, expected[testUnitB].String(), median[testUnitB].String())
	})

	t.Run("MedianSpend returns median single spend value", func(t *testing.T) {
		t.Parallel()

		billingClient := newMockBillingClient()
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t))
		report.client = billingClient
		err := report.Initialize(t.Context())
		require.NoError(t, err)
		err = report.balance.Add(100)
		require.NoError(t, err)

		steps := []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "abc", SpendUnit: "a", SpendValue: "1"},
		}

		for idx := range steps {
			_, err := report.ReserveByLimits(strconv.Itoa(idx), capabilities.CapabilityInfo{}, []SpendTuple{{Value: 1, Unit: "SomeUnit"}})
			require.NoError(t, err)
			require.NoError(t, report.SetStep(strconv.Itoa(idx), steps))
		}

		expected := map[SpendUnit]SpendValue{
			testUnitA: testUnitA.IntToSpendValue(1),
		}

		median := report.MedianSpend()

		require.Len(t, median, 1)
		require.Contains(t, maps.Keys(median), testUnitA)

		assert.Equal(t, expected[testUnitA].String(), median[testUnitA].String())
	})

	t.Run("MedianSpend returns median odd number of spend values", func(t *testing.T) {
		t.Parallel()

		billingClient := newMockBillingClient()
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t))
		report.client = billingClient
		err := report.Initialize(t.Context())
		require.NoError(t, err)
		err = report.balance.Add(100)
		require.NoError(t, err)

		steps := []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "abc", SpendUnit: testA, SpendValue: "1"},
			{Peer2PeerID: "abc", SpendUnit: testA, SpendValue: "3"},
			{Peer2PeerID: "xyz", SpendUnit: testA, SpendValue: "2"},
		}

		for idx := range steps {
			_, err := report.ReserveByLimits(strconv.Itoa(idx), capabilities.CapabilityInfo{}, []SpendTuple{{Value: 1, Unit: "SomeUnit"}})
			require.NoError(t, err)
			require.NoError(t, report.SetStep(strconv.Itoa(idx), steps))
		}

		expected := map[SpendUnit]SpendValue{
			testUnitA: testUnitA.IntToSpendValue(2),
		}

		median := report.MedianSpend()

		require.Len(t, median, 1)
		require.Contains(t, maps.Keys(median), testUnitA)

		assert.Equal(t, expected[testUnitA].String(), median[testUnitA].String())
	})

	t.Run("MedianSpend returns median as average for even number of spend values", func(t *testing.T) {
		t.Parallel()

		billingClient := newMockBillingClient()
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t))
		report.client = billingClient
		err := report.Initialize(t.Context())
		require.NoError(t, err)
		err = report.balance.Add(100)
		require.NoError(t, err)

		steps := []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "xyz", SpendUnit: testA, SpendValue: "42"},
			{Peer2PeerID: "abc", SpendUnit: testA, SpendValue: "1"},
			{Peer2PeerID: "abc", SpendUnit: testA, SpendValue: "3"},
			{Peer2PeerID: "xyz", SpendUnit: testA, SpendValue: "2"},
		}

		for idx := range steps {
			_, err := report.ReserveByLimits(strconv.Itoa(idx), capabilities.CapabilityInfo{}, []SpendTuple{{Value: 1, Unit: "SomeUnit"}})
			require.NoError(t, err)
			require.NoError(t, report.SetStep(strconv.Itoa(idx), steps))
		}

		expected := map[SpendUnit]SpendValue{
			testUnitA: testUnitA.DecimalToSpendValue(decimal.NewFromFloat(2.5)),
		}

		median := report.MedianSpend()

		require.Len(t, median, 1)
		require.Contains(t, maps.Keys(median), testUnitA)

		assert.Equal(t, expected[testUnitA].String(), median[testUnitA].String())
	})

	t.Run("ReserveStep returns error if step already exists", func(t *testing.T) {
		t.Parallel()
		billingClient := newMockBillingClient()
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t))
		report.client = billingClient
		err := report.Initialize(t.Context())
		require.NoError(t, err)
		err = report.balance.Add(100)
		require.NoError(t, err)
		_, err = report.ReserveByLimits("ref1", capabilities.CapabilityInfo{}, []SpendTuple{{Value: 1, Unit: "SomeUnit"}})
		require.NoError(t, err)
		_, err = report.ReserveByLimits("ref1", capabilities.CapabilityInfo{}, []SpendTuple{{Value: 1, Unit: "SomeUnit"}})
		require.Error(t, err)
	})

	t.Run("SetStep returns error if reserve is not called first", func(t *testing.T) {
		t.Parallel()
		billingClient := newMockBillingClient()
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t))
		report.client = billingClient
		err := report.Initialize(t.Context())
		require.NoError(t, err)
		err = report.balance.Add(100)
		require.NoError(t, err)

		steps := []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "xyz", SpendUnit: testA, SpendValue: "42"},
			{Peer2PeerID: "abc", SpendUnit: testA, SpendValue: "1"},
		}

		require.Error(t, report.SetStep("ref1", steps))
	})

	t.Run("SetStep returns error if step already exists", func(t *testing.T) {
		t.Parallel()
		billingClient := newMockBillingClient()
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t))
		report.client = billingClient
		err := report.Initialize(t.Context())
		require.NoError(t, err)
		err = report.balance.Add(100)
		require.NoError(t, err)

		steps := []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "xyz", SpendUnit: testA, SpendValue: "42"},
			{Peer2PeerID: "abc", SpendUnit: testA, SpendValue: "1"},
		}

		_, err = report.ReserveByLimits("ref1", capabilities.CapabilityInfo{}, []SpendTuple{{Value: 1, Unit: "SomeUnit"}})
		require.NoError(t, err)
		require.NoError(t, report.SetStep("ref1", steps))
		require.Error(t, report.SetStep("ref1", steps))
	})
}

// Test_MeterReports tests the Add, Get, Delete, and Len methods of a MeterReports.
// It also tests concurrent safe access.
func Test_MeterReports(t *testing.T) {
	mr := NewReports(nil)
	assert.Equal(t, 0, mr.Len())
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer wg.Done()
		report := NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t))
		mr.Add("exec1", report)
		r, ok := mr.Get("exec1")
		assert.True(t, ok)
		//nolint:errcheck // depending on the concurrent timing, this may or may not err
		report.ReserveByLimits("ref1", capabilities.CapabilityInfo{}, []SpendTuple{{Value: 1, Unit: "SomeUnit"}})
		//nolint:errcheck // depending on the concurrent timing, this may or may not err
		r.SetStep("ref1", []capabilities.MeteringNodeDetail{})
		mr.Delete("exec1")
	}()
	go func() {
		defer wg.Done()
		mr.Add("exec2", NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t)))
		r, ok := mr.Get("exec2")
		assert.True(t, ok)
		_, err := r.ReserveByLimits("ref1", capabilities.CapabilityInfo{}, []SpendTuple{{Value: 1, Unit: "SomeUnit"}})
		assert.NoError(t, err)
		err = r.SetStep("ref1", []capabilities.MeteringNodeDetail{})
		assert.NoError(t, err)
		mr.Delete("exec2")
	}()
	go func() {
		defer wg.Done()
		mr.Add("exec1", NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t)))
		r, ok := mr.Get("exec1")
		assert.True(t, ok)
		//nolint:errcheck // depending on the concurrent timing, this may or may not err
		r.ReserveByLimits("ref1", capabilities.CapabilityInfo{}, []SpendTuple{{Value: 1, Unit: "SomeUnit"}})
		//nolint:errcheck // depending on the concurrent timing, this may or may not err
		r.SetStep("ref1", []capabilities.MeteringNodeDetail{})
		mr.Delete("exec1")
	}()

	wg.Wait()
	assert.Equal(t, 0, mr.Len())
}

func Test_MeterReportsLength(t *testing.T) {
	mr := NewReports(nil)

	mr.Add("exec1", NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t)))
	mr.Add("exec2", NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t)))
	mr.Add("exec3", NewReport(testAccountID, testWorkflowID, testWorkflowExecutionID, logger.TestSugared(t)))
	assert.Equal(t, 3, mr.Len())

	mr.Delete("exec2")
	assert.Equal(t, 2, mr.Len())
}
