package runtime

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"
)

// Receipt is a tamper-evident digest over a normalized runtime result.
type Receipt struct {
	TaskID      string    `json:"task_id"`
	Digest      string    `json:"sha256"`
	GeneratedAt time.Time `json:"generated_at"`
}

type receiptPayload struct {
	TaskID    string       `json:"task_id"`
	Mode      Mode         `json:"mode"`
	Completed bool         `json:"completed"`
	Summary   string       `json:"summary"`
	Steps     []StepResult `json:"steps"`
	Evidence  []Evidence   `json:"evidence"`
}

// SealResult creates a deterministic SHA-256 digest over the semantically relevant
// result fields. The timestamp belongs to the receipt, not the hashed payload.
func SealResult(result Result) (Receipt, error) {
	if result.TaskID == "" {
		return Receipt{}, errors.New("result task id is required")
	}

	payload := receiptPayload{
		TaskID:    result.TaskID,
		Mode:      result.Mode,
		Completed: result.Completed,
		Summary:   result.Summary,
		Steps:     result.Steps,
		Evidence:  result.Evidence,
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		return Receipt{}, err
	}
	digest := sha256.Sum256(encoded)
	return Receipt{
		TaskID:      result.TaskID,
		Digest:      hex.EncodeToString(digest[:]),
		GeneratedAt: time.Now().UTC(),
	}, nil
}

// VerifyReceipt recomputes the digest and compares it with an existing receipt.
func VerifyReceipt(result Result, receipt Receipt) (bool, error) {
	if receipt.TaskID != result.TaskID {
		return false, nil
	}
	calculated, err := SealResult(result)
	if err != nil {
		return false, err
	}
	return calculated.Digest == receipt.Digest, nil
}
