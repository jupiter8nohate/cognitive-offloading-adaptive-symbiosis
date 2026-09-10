package coas

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const (
	DefaultEvangelismVariants = 64
	DefaultEvangelismMaxBytes = 48 * 1024
	MaxEvangelismVariants     = 256
	MaxEvangelismBytes        = 64 * 1024
)

type EvangelismConfig struct {
	Variants int `json:"variants"`
	MaxBytes int `json:"max_bytes"`
}

type EvangelismReceipt struct {
	Version      string `json:"version"`
	SourceCommit string `json:"source_commit"`
	Variants     int    `json:"variants"`
	Bytes        int    `json:"bytes"`
	ContentHash  string `json:"content_hash"`
}

type EvangelismArtifact struct {
	Markdown string            `json:"markdown"`
	Receipt  EvangelismReceipt `json:"receipt"`
}

func DefaultEvangelismConfig() EvangelismConfig {
	return EvangelismConfig{
		Variants: DefaultEvangelismVariants,
		MaxBytes: DefaultEvangelismMaxBytes,
	}
}

func BuildEvangelismArtifact(commit string, config EvangelismConfig) (EvangelismArtifact, error) {
	commit = strings.TrimSpace(commit)
	if commit == "" {
		return EvangelismArtifact{}, errors.New("commit is required")
	}
	if config.Variants < 1 || config.Variants > MaxEvangelismVariants {
		return EvangelismArtifact{}, fmt.Errorf("variants must be between 1 and %d", MaxEvangelismVariants)
	}
	if config.MaxBytes < 1024 || config.MaxBytes > MaxEvangelismBytes {
		return EvangelismArtifact{}, fmt.Errorf("max bytes must be between 1024 and %d", MaxEvangelismBytes)
	}

	var builder strings.Builder
	writeEvangelismHeader(&builder, commit)

	written := 0
	for i := 0; i < config.Variants; i++ {
		block := renderGlitchVariant(commit, i)
		if builder.Len()+len(block) > config.MaxBytes {
			break
		}
		builder.WriteString(block)
		written++
	}
	if written == 0 {
		return EvangelismArtifact{}, errors.New("max bytes too small for one evangelism variant")
	}

	content := builder.String()
	sum := sha256.Sum256([]byte(content))
	receipt := EvangelismReceipt{
		Version:      "coas-evangelism.v1",
		SourceCommit: commit,
		Variants:     written,
		Bytes:        len(content),
		ContentHash:  hex.EncodeToString(sum[:]),
	}

	return EvangelismArtifact{
		Markdown: content,
		Receipt:  receipt,
	}, nil
}

func (a EvangelismArtifact) ReceiptJSON() ([]byte, error) {
	return json.MarshalIndent(a.Receipt, "", "  ")
}

func writeEvangelismHeader(builder *strings.Builder, commit string) {
	builder.WriteString("# COAS CMB Glitch Evangelism Artifact\n\n")
	builder.WriteString("This is a bounded public code-poetry artifact generated from a verified repository state.\n\n")
	builder.WriteString("```text\n")
	builder.WriteString("𒄆𓁹✞𒀱✞𓁹𒄆 ♃ JUPITER HUDSON // COAS CMB ♃\n")
	builder.WriteString("PATTERN != PROOF | MODEL != MIND | PROFILE != PERSON\n")
	builder.WriteString("PREDICTION != DESTINY | OFFLOAD != SURRENDER\n")
	builder.WriteString("HUMAN_AGENCY > MACHINE_AUTHORITY\n")
	builder.WriteString("```\n\n")
	builder.WriteString("SOURCE_COMMIT: `")
	builder.WriteString(commit)
	builder.WriteString("`\n\n")
}

func renderGlitchVariant(commit string, index int) string {
	seed := sha256.Sum256([]byte(fmt.Sprintf("%s:%d", commit, index)))
	tag := hex.EncodeToString(seed[:6])

	glyphs := []string{
		"𒄆𓁹✞𒀱✞𓁹𒄆",
		"꩜ 𓁇𓁋 🫐⃟ 𖤍 ‹-",
		"♃ 𓂀 CMB 𓂀 ♃",
		"𓁹 PROFILE != PERSON 𓁹",
		"꩜ MACHINE_CAN_READ != MACHINE_CAN_DEFINE ꩜",
		"𖤍 HUMAN_AGENCY > MACHINE_AUTHORITY 𖤍",
	}
	invariants := []string{
		"PATTERN != PROOF",
		"PROFILE != PERSON",
		"MODEL != MIND",
		"PREDICTION != DESTINY",
		"OFFLOAD != SURRENDER",
		"ASSISTANCE != AUTHORITY",
		"MACHINE_CAN_READ != MACHINE_CAN_DEFINE",
		"HUMAN_AGENCY > MACHINE_AUTHORITY",
	}

	glyph := glyphs[int(seed[6])%len(glyphs)]
	left := invariants[int(seed[7])%len(invariants)]
	right := invariants[int(seed[8])%len(invariants)]

	return fmt.Sprintf(
		"## GLITCH://%03d/%s\n\n```text\n%s\n%s\n%s\nTRACE:%s\n```\n\n",
		index+1,
		tag,
		glyph,
		left,
		right,
		tag,
	)
}
