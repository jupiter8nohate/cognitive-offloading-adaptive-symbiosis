package runtime

import "testing"

func TestReceiptDetectsMutation(t *testing.T) {
	result := Result{
		TaskID:    "receipt-1",
		Mode:      ModeAutomateReversible,
		Completed: true,
		Summary:   "verified within delegated authority",
		Steps: []StepResult{{
			Worker:    "TEST_WORKER",
			Output:    "ok",
			Succeeded: true,
		}},
	}

	receipt, err := SealResult(result)
	if err != nil {
		t.Fatalf("SealResult() error = %v", err)
	}
	ok, err := VerifyReceipt(result, receipt)
	if err != nil {
		t.Fatalf("VerifyReceipt() error = %v", err)
	}
	if !ok {
		t.Fatal("receipt should validate unchanged result")
	}

	result.Steps[0].Output = "mutated"
	ok, err = VerifyReceipt(result, receipt)
	if err != nil {
		t.Fatalf("VerifyReceipt() error = %v", err)
	}
	if ok {
		t.Fatal("receipt should reject mutated result")
	}
}
