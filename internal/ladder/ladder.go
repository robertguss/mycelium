// Package ladder is the pure perspective-ladder grammar (no filesystem, no CLI).
package ladder

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// Rung is second-opinion or council.
type Rung string

const (
	RungSecondOpinion Rung = "second-opinion"
	RungCouncil       Rung = "council"
)

// Adapter is cursor or manual.
type Adapter string

const (
	AdapterCursor Adapter = "cursor"
	AdapterManual Adapter = "manual"
)

// CostClass is cheap | quick | standard | high-stakes.
type CostClass string

const (
	CostCheap      CostClass = "cheap"
	CostQuick      CostClass = "quick"
	CostStandard   CostClass = "standard"
	CostHighStakes CostClass = "high-stakes"
)

const SeedDissentToken = "SEED-DISSENT"

var (
	ErrEmpty     = errors.New("ladder: empty value")
	ErrUnknown   = errors.New("ladder: unknown value")
	ErrCard      = errors.New("ladder: cardinality")
	ErrSeedDissent = errors.New("ladder: SEED-DISSENT missing in retained dissent")
)

// ParseRung accepts exact tokens only (no case fold, no trim).
func ParseRung(s string) (Rung, error) {
	switch s {
	case string(RungSecondOpinion):
		return RungSecondOpinion, nil
	case string(RungCouncil):
		return RungCouncil, nil
	case "":
		return "", ErrEmpty
	default:
		return "", fmt.Errorf("%w: rung %q", ErrUnknown, s)
	}
}

// ParseAdapter accepts exact tokens only.
func ParseAdapter(s string) (Adapter, error) {
	switch s {
	case string(AdapterCursor):
		return AdapterCursor, nil
	case string(AdapterManual):
		return AdapterManual, nil
	case "":
		return "", ErrEmpty
	default:
		return "", fmt.Errorf("%w: adapter %q", ErrUnknown, s)
	}
}

// ParseCostClass accepts exact tokens only.
func ParseCostClass(s string) (CostClass, error) {
	switch s {
	case string(CostCheap):
		return CostCheap, nil
	case string(CostQuick):
		return CostQuick, nil
	case string(CostStandard):
		return CostStandard, nil
	case string(CostHighStakes):
		return CostHighStakes, nil
	case "":
		return "", ErrEmpty
	default:
		return "", fmt.Errorf("%w: cost_class %q", ErrUnknown, s)
	}
}

// CostClassOK is the §6.1 IFF table.
func CostClassOK(rung Rung, class CostClass) bool {
	switch rung {
	case RungSecondOpinion:
		return class == CostCheap
	case RungCouncil:
		return class == CostQuick || class == CostStandard || class == CostHighStakes
	default:
		return false
	}
}

// OptInOK accepts TOML boolean true only. false, "true", missing → false.
func OptInOK(v any) bool {
	b, ok := v.(bool)
	return ok && b
}

// PromptSHA256 is hex(sha256(TrimSpace(SectionBody(body, "Prompt")))).
func PromptSHA256(body string) string {
	sec := strings.TrimSpace(SectionBody(body, "Prompt"))
	sum := sha256.Sum256([]byte(sec))
	return hex.EncodeToString(sum[:])
}

// RequiredCMPH2 returns commissioning required H2 names.
func RequiredCMPH2() []string {
	return []string{"Prompt", "Attachments", "Cost"}
}

// RequiredRPTH2 returns model-report required H2 names.
func RequiredRPTH2() []string {
	return []string{"Position", "Findings", "Dissent"}
}

// RequiredRCLH2 returns reconciliation required H2 names (v1 list).
func RequiredRCLH2() []string {
	return []string{
		"Convergence",
		"Material disagreement",
		"Evidence unique to one report",
		"Contradictory evidence",
		"Different assumptions",
		"Different scope interpretations",
		"Recommendations independently supported",
		"Questions requiring another spike",
		"Final reconciled recommendation",
		"Retained dissent",
	}
}

// Cardinality enforces §6.5. nRPT=0,nRCL=0 is WIP PASS.
func Cardinality(rung Rung, nRPT, nRCL int) error {
	switch rung {
	case RungSecondOpinion:
		if nRPT == 0 && nRCL == 0 {
			return nil
		}
		if nRPT == 1 && nRCL == 0 {
			return nil
		}
		if nRCL > 0 {
			return fmt.Errorf("%w: second-opinion forbids RCL (got %d)", ErrCard, nRCL)
		}
		return fmt.Errorf("%w: second-opinion requires exactly one RPT (got %d)", ErrCard, nRPT)
	case RungCouncil:
		if nRPT == 0 && nRCL == 0 {
			return nil
		}
		if nRPT >= 2 && nRCL == 1 {
			return nil
		}
		return fmt.Errorf("%w: council requires >=2 RPT and exactly one RCL (got %d RPT, %d RCL)", ErrCard, nRPT, nRCL)
	default:
		return fmt.Errorf("%w: unknown rung %q", ErrCard, rung)
	}
}

// SeedDissentOK fails when any RPT dissent contains SEED-DISSENT but rclRetained does not.
func SeedDissentOK(rptDissentBodies []string, rclRetained string) error {
	seeded := false
	for _, body := range rptDissentBodies {
		if strings.Contains(body, SeedDissentToken) {
			seeded = true
			break
		}
	}
	if !seeded {
		return nil
	}
	if strings.Contains(rclRetained, SeedDissentToken) {
		return nil
	}
	return ErrSeedDissent
}

// SectionBody returns bytes after the first exact H2 until the next H2 or EOF.
func SectionBody(body, h2 string) string {
	want := "## " + h2
	lines := strings.Split(body, "\n")
	start := -1
	for i, line := range lines {
		trim := strings.TrimRight(line, "\r")
		if trim == want {
			start = i + 1
			break
		}
	}
	if start < 0 {
		return ""
	}
	var b strings.Builder
	for i := start; i < len(lines); i++ {
		trim := strings.TrimRight(lines[i], "\r")
		if strings.HasPrefix(trim, "## ") {
			break
		}
		b.WriteString(lines[i])
		if i+1 < len(lines) {
			b.WriteByte('\n')
		}
	}
	return b.String()
}
