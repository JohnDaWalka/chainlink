package worker

import (
	"errors"
	"testing"
	"time"
)

// Example types to demonstrate type safety
type WorkflowConfig struct {
	ID       string
	Settings map[string]string
}

type FundingInfo struct {
	Amount   int64
	TxHashes []string
}

func TestTypedBackgroundRunner_TypedResults(t *testing.T) {
	runner := NewTypedBackgroundRunner()

	// Add tasks with different return types - all fully typed
	stringResult := Add(runner, "String Task", func() (string, error) {
		time.Sleep(50 * time.Millisecond)
		return "hello world", nil
	})

	intResult := Add(runner, "Int Task", func() (int, error) {
		time.Sleep(30 * time.Millisecond)
		return 42, nil
	})

	structResult := Add(runner, "Struct Task", func() (WorkflowConfig, error) {
		return WorkflowConfig{
			ID:       "workflow-123",
			Settings: map[string]string{"key": "value"},
		}, nil
	})

	sliceResult := Add(runner, "Slice Task", func() ([]string, error) {
		return []string{"item1", "item2", "item3"}, nil
	})

	// Wait for all tasks to complete
	err := runner.Wait()
	if err != nil {
		t.Fatalf("Expected no error from Wait(), got: %v", err)
	}

	// Get typed results - no type assertions needed!
	str, err := stringResult.Get()
	if err != nil {
		t.Fatalf("String task failed: %v", err)
	}
	if str != "hello world" {
		t.Errorf("Expected 'hello world', got: %s", str)
	}

	num, err := intResult.Get()
	if err != nil {
		t.Fatalf("Int task failed: %v", err)
	}
	if num != 42 {
		t.Errorf("Expected 42, got: %d", num)
	}

	config, err := structResult.Get()
	if err != nil {
		t.Fatalf("Struct task failed: %v", err)
	}
	if config.ID != "workflow-123" || config.Settings["key"] != "value" {
		t.Errorf("Unexpected struct result: %+v", config)
	}

	slice, err := sliceResult.Get()
	if err != nil {
		t.Fatalf("Slice task failed: %v", err)
	}
	if len(slice) != 3 || slice[0] != "item1" {
		t.Errorf("Unexpected slice result: %v", slice)
	}
}

func TestTypedBackgroundRunner_ErrorHandling(t *testing.T) {
	runner := NewTypedBackgroundRunner()

	successResult := Add(runner, "Success Task", func() (string, error) {
		return "success", nil
	})

	errorResult := Add(runner, "Error Task", func() (int, error) {
		return 0, errors.New("task failed")
	})

	err := runner.Wait()
	if err == nil {
		t.Fatal("Wait() should return error for task errors")
	}

	// Check individual results
	successValue, err := successResult.Get()
	if err != nil {
		t.Errorf("Success task should not fail, got: %v", err)
	}
	if successValue != "success" {
		t.Errorf("Expected 'success', got: %s", successValue)
	}

	errorValue, err := errorResult.Get()
	if err == nil {
		t.Error("Error task should have failed")
	}
	if errorValue != 0 {
		t.Errorf("Error task should return zero value, got: %d", errorValue)
	}
	if err.Error() != "task failed" {
		t.Errorf("Expected 'task failed', got: %v", err)
	}
}

func TestTypedBackgroundRunner_VoidTasks(t *testing.T) {
	runner := NewTypedBackgroundRunner()

	var sideEffect1, sideEffect2 bool

	voidResult1 := AddVoid(runner, "Void Task 1", func() error {
		sideEffect1 = true
		return nil
	})

	voidResult2 := AddVoid(runner, "Void Task 2", func() error {
		sideEffect2 = true
		return errors.New("void task error")
	})

	err := runner.Wait()
	if err == nil {
		t.Fatal("Wait() did not fail")
	}

	// Check void results
	_, err = voidResult1.Get()
	if err != nil {
		t.Errorf("Void task 1 should succeed, got: %v", err)
	}

	_, err = voidResult2.Get()
	if err == nil {
		t.Error("Void task 2 should have failed")
	}

	if !sideEffect1 || !sideEffect2 {
		t.Error("Side effects should have occurred")
	}
}

func TestTypedBackgroundRunner_Panic(t *testing.T) {
	runner := NewTypedBackgroundRunner()

	_ = Add(runner, "Panicking Task", func() (string, error) {
		panic("test panic")
	})

	// Should panic when Wait() is called
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Expected a panic, but didn't get one")
		}
	}()

	_ = runner.Wait()
}

func TestTypedBackgroundRunner_IsReady(t *testing.T) {
	runner := NewTypedBackgroundRunner()

	slowResult := Add(runner, "Slow Task", func() (string, error) {
		time.Sleep(100 * time.Millisecond)
		return "done", nil
	})

	fastResult := Add(runner, "Fast Task", func() (string, error) {
		return "quick", nil
	})

	// Check readiness before completion
	time.Sleep(20 * time.Millisecond)

	if slowResult.IsReady() {
		t.Error("Slow task should not be ready yet")
	}

	// Fast task might be ready by now, but let's wait for both
	err := runner.Wait()
	if err != nil {
		t.Fatalf("Wait() failed: %v", err)
	}

	// Now both should be ready
	if !slowResult.IsReady() {
		t.Error("Slow task should be ready after Wait()")
	}
	if !fastResult.IsReady() {
		t.Error("Fast task should be ready after Wait()")
	}

	// Get results
	slowValue, err := slowResult.Get()
	if err != nil || slowValue != "done" {
		t.Errorf("Unexpected slow result: %s, %v", slowValue, err)
	}

	fastValue, err := fastResult.Get()
	if err != nil || fastValue != "quick" {
		t.Errorf("Unexpected fast result: %s, %v", fastValue, err)
	}
}

// Example of real-world usage pattern
func TestTypedBackgroundRunner_RealWorldExample(t *testing.T) {
	runner := NewTypedBackgroundRunner()

	// Simulate complex setup tasks with different return types
	workflowResult := Add(runner, "Configure Workflow Registry", func() (*WorkflowConfig, error) {
		// Simulate complex contract configuration
		time.Sleep(50 * time.Millisecond)
		return &WorkflowConfig{
			ID:       "registry-456",
			Settings: map[string]string{"timeout": "30s"},
		}, nil
	})

	fundingResult := Add(runner, "Fund Chainlink nodes", func() (FundingInfo, error) {
		// Simulate funding operations
		time.Sleep(30 * time.Millisecond)
		return FundingInfo{
			Amount:   1000000000000000000, // 1 ETH in wei
			TxHashes: []string{"0xabc123", "0xdef456"},
		}, nil
	})

	configResult := AddVoid(runner, "Load Node Configs", func() error {
		// Simulate configuration loading
		time.Sleep(20 * time.Millisecond)
		return nil
	})

	// Wait for all background tasks
	err := runner.Wait()
	if err != nil {
		t.Fatalf("Background setup failed: %v", err)
	}

	// Get typed results
	workflow, err := workflowResult.Get()
	if err != nil {
		t.Fatalf("Workflow configuration failed: %v", err)
	}

	funding, err := fundingResult.Get()
	if err != nil {
		t.Fatalf("Node funding failed: %v", err)
	}

	_, err = configResult.Get()
	if err != nil {
		t.Fatalf("Config loading failed: %v", err)
	}

	// Use the typed results
	if workflow.ID != "registry-456" {
		t.Errorf("Unexpected workflow ID: %s", workflow.ID)
	}
	if funding.Amount != 1000000000000000000 {
		t.Errorf("Unexpected funding amount: %d", funding.Amount)
	}
	if len(funding.TxHashes) != 2 {
		t.Errorf("Expected 2 tx hashes, got: %d", len(funding.TxHashes))
	}

	// Verify all expected tasks were tracked
	taskNames := runner.GetTaskNames()
	expectedNames := []string{"Configure Workflow Registry", "Fund Chainlink nodes", "Load Node Configs"}
	if len(taskNames) != len(expectedNames) {
		t.Errorf("Expected %d tasks, got %d", len(expectedNames), len(taskNames))
	}
}
