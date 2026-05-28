// Package honest implements the cohort R143 LOUD-ONCE-WARNING-FLAG
// discipline for limitless-immigration.
//
// 5 honest-defaults surfaces:
//
//  1. LIMITLESS_IMMIGRATION_HC_395_VERSION_PIN_REQUIRED — UK Immigration
//     Rules HC 395 is amended ~quarterly via Statements of Changes; the
//     local corpus SHA must be cold-verified against gov.uk before any
//     live decision.
//  2. LIMITLESS_IMMIGRATION_APPENDIX_FM_FINANCIAL_REQUIREMENT_FROZEN —
//     Appendix FM minimum income threshold (£29,000 from 2024-04-11, was
//     £18,600 pre-Sunak reform) is a moving target. Hard-coded values
//     are wrong-confidence-risk per knowledge-bedrock guidance.
//  3. LIMITLESS_IMMIGRATION_ETA_ROLLOUT_INCOMPLETE — Electronic Travel
//     Authorisation rolling out by nationality through 2025-2026; status
//     by nationality drifts continuously.
//  4. LIMITLESS_IMMIGRATION_DECISION_NOT_LEGAL_ADVICE — every emitted
//     determination MUST carry R166 LIABILITY-FOOTER-CONST escape; this
//     scaffold is NOT a substitute for OISC-regulated immigration advice.
//  5. LIMITLESS_IMMIGRATION_REVIEWED_BY_COUNSEL_FALSE — R166 honest-
//     default. Phase-1 templates have NOT been reviewed by qualified
//     immigration counsel or OISC-regulated advisers.
package honest

import (
	"fmt"
	"io"
	"sync"
)

const LoudOncePrefix = "[LOUD-ONCE-WARNING]"

// Severity follows R143.A SEVERITY-LADDER-CONVENTION.
type Severity string

const (
	SeverityInfo  Severity = "INFO"
	SeverityWarn  Severity = "WARN"
	SeverityError Severity = "ERROR"
)

// Advisory is one R143 advisory entry.
type Advisory struct {
	Code     string
	Severity Severity
	Message  string
	DocLink  string
}

var canonicalAdvisories = []Advisory{
	{
		Code:     "LIMITLESS_IMMIGRATION_HC_395_VERSION_PIN_REQUIRED",
		Severity: SeverityError,
		Message:  "UK Immigration Rules HC 395 (as amended) ships under a Statement of Changes cycle (~quarterly). The local corpus SHA pinned in internal/immigration-rules/ MUST be cold-verified against gov.uk-published canonical before any live decision. Stale pin = wrong-confidence-risk per knowledge-bedrock guidance: silent enforcement of superseded paragraphs is the worst failure mode in immigration adjudication.",
		DocLink:  "SECURITY.md",
	},
	{
		Code:     "LIMITLESS_IMMIGRATION_APPENDIX_FM_FINANCIAL_REQUIREMENT_FROZEN",
		Severity: SeverityError,
		Message:  "Appendix FM Minimum Income Requirement (MIR) for partner / spouse routes is a moving target. From 2024-04-11 the threshold is £29,000 (up from £18,600 pre-Sunak reform); further uplift to £38,700 was proposed but paused post-2024 General Election. Hard-coded thresholds in code are wrong-confidence-risk: any quoted figure MUST be sourced from internal/immigration-rules/ corpus + cold-verified against current gov.uk policy paper.",
		DocLink:  "SECURITY.md",
	},
	{
		Code:     "LIMITLESS_IMMIGRATION_ETA_ROLLOUT_INCOMPLETE",
		Severity: SeverityWarn,
		Message:  "UK Electronic Travel Authorisation (ETA) (under Nationality and Borders Act 2022 + UKBA 2007 s.11A as amended) is rolling out by nationality through 2025-2026. ETA-required nationality lists drift; a determination of 'ETA required: yes/no' MUST be timestamped + sourced from internal/immigration-rules/ corpus. Wave-by-wave rollout is documented at gov.uk/guidance/apply-for-an-electronic-travel-authorisation-eta.",
		DocLink:  "SECURITY.md",
	},
	{
		Code:     "LIMITLESS_IMMIGRATION_DECISION_NOT_LEGAL_ADVICE",
		Severity: SeverityError,
		Message:  "Every immigration eligibility determination emitted by this software MUST carry the R166 LIABILITY-FOOTER-CONST escape: 'NOT LEGAL ADVICE. Immigration advice is regulated in the UK by the OISC (Office of the Immigration Services Commissioner) under the Immigration and Asylum Act 1999 s.84. Unregulated immigration advice is a criminal offence under IAA 1999 s.91. Consult an OISC-regulated adviser or solicitor before relying on any determination.'",
		DocLink:  "SECURITY.md",
	},
	{
		Code:     "LIMITLESS_IMMIGRATION_REVIEWED_BY_COUNSEL_FALSE",
		Severity: SeverityWarn,
		Message:  "R166 LIABILITY-FOOTER-CONST honest-default. Phase-1 scaffold ships ReviewedByCounsel = false. Placeholder narrative templates + eligibility-rule scaffolds have NOT been reviewed by qualified immigration counsel or OISC-regulated advisers. Operator MUST commission counsel review + flip ReviewedByCounsel to true on its own R145.B sibling branch before any live deployment.",
		DocLink:  "SECURITY.md",
	},
}

var (
	registryMu sync.RWMutex
	registry   = map[string]*sync.Once{}
)

// LoudOnce emits the advisory exactly once per package process-lifetime.
func LoudOnce(adv Advisory, w io.Writer) {
	registryMu.RLock()
	once, ok := registry[adv.Code]
	registryMu.RUnlock()
	if !ok {
		registryMu.Lock()
		once, ok = registry[adv.Code]
		if !ok {
			once = &sync.Once{}
			registry[adv.Code] = once
		}
		registryMu.Unlock()
	}
	once.Do(func() {
		_, _ = fmt.Fprintf(w, "%s %s %s: %s (see %s)\n",
			LoudOncePrefix, adv.Severity, adv.Code, adv.Message, adv.DocLink)
	})
}

// Reset clears the once-gate registry. Test-only.
func Reset() {
	registryMu.Lock()
	registry = map[string]*sync.Once{}
	registryMu.Unlock()
}

// CanonicalAdvisories returns a defensive copy.
func CanonicalAdvisories() []Advisory {
	out := make([]Advisory, len(canonicalAdvisories))
	copy(out, canonicalAdvisories)
	return out
}

// FindAdvisory looks up a canonical advisory by Code.
func FindAdvisory(code string) (Advisory, bool) {
	for _, a := range canonicalAdvisories {
		if a.Code == code {
			return a, true
		}
	}
	return Advisory{}, false
}
