// Package manifest implements R150 PARALLEL-MAP review-metadata envelope
// for limitless-immigration, extended with the R150 Class-3 jurisdiction-
// version anchor for UK Home Office Immigration Rules corpus pinning.
package manifest

import (
	"sort"
	"time"
)

// SchemaVersion pins the R150 envelope version.
const SchemaVersion = 1

// FreshAtUnknown is the sentinel value for FreshAt fields (honest-TODO).
var FreshAtUnknown = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

const (
	SourceHomeOfficeImmRules     = "UK Home Office Immigration Rules HC 395 (as amended) — gov.uk/government/collections/immigration-rules"
	SourceAppendixFMSpouse       = "UK Home Office Immigration Rules Appendix FM (Family Members) + Appendix FM-SE (Financial Evidence)"
	SourceETARollout             = "UK Home Office Electronic Travel Authorisation rollout schedule — gov.uk/guidance/apply-for-an-electronic-travel-authorisation-eta"
	SourceOISC                   = "Office of the Immigration Services Commissioner (OISC) — IAA 1999 s.83 + s.84"
	SourceMethodologyCorpusPkg   = "limitless-immigration internal/immigration-rules package"
	SourceContextDoc             = "limitless-immigration CONTEXT.md"
	SourceR85ParityMarker        = "limitless-immigration R85 CLEAN-PARITY"
)

type Confidence int

const (
	ConfidenceHigh   Confidence = 3
	ConfidenceMedium Confidence = 2
	ConfidenceLow    Confidence = 1
)

type Jurisdiction string

const (
	JurisdictionUK   Jurisdiction = "UK"
	JurisdictionNone Jurisdiction = ""
)

type Entry struct {
	Key           string
	Description   string
	FreshAt       time.Time
	Source        string
	SchemaVersion int
	Confidence    Confidence
	Jurisdiction  Jurisdiction
	Version       string
}

func (e Entry) IsStale(now time.Time, maxAge time.Duration) bool {
	if e.FreshAt.Equal(FreshAtUnknown) {
		return true
	}
	return now.Sub(e.FreshAt) > maxAge
}

type Manifest []Entry

func (m Manifest) SortedKeys() []string {
	keys := make([]string, 0, len(m))
	for _, e := range m {
		keys = append(keys, e.Key)
	}
	sort.Strings(keys)
	return keys
}

func (m Manifest) StaleEntries(now time.Time, maxAge time.Duration) []Entry {
	var out []Entry
	for _, e := range m {
		if e.IsStale(now, maxAge) {
			out = append(out, e)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	return out
}

func AllSources() []string {
	return []string{
		SourceAppendixFMSpouse,
		SourceContextDoc,
		SourceETARollout,
		SourceHomeOfficeImmRules,
		SourceMethodologyCorpusPkg,
		SourceOISC,
		SourceR85ParityMarker,
	}
}

// Seed returns the canonical R150 manifest for limitless-immigration.
func Seed() Manifest {
	scaffold := time.Date(2026, 5, 28, 0, 0, 0, 0, time.UTC)
	tbdPhase2 := FreshAtUnknown
	return Manifest{
		{
			Key:           "corpus.uk.home_office_immigration_rules_hc395",
			Description:   "UK Home Office Immigration Rules HC 395 (as amended). Statements of Changes are issued ~quarterly. Local SHA pin awaiting Phase-2 cold-verify against gov.uk-published canonical.",
			FreshAt:       tbdPhase2,
			Source:        SourceHomeOfficeImmRules,
			SchemaVersion: SchemaVersion,
			Confidence:    ConfidenceLow,
			Jurisdiction:  JurisdictionUK,
			Version:       "as-amended-2024",
		},
		{
			Key:           "corpus.uk.appendix_fm_financial_requirement",
			Description:   "Appendix FM minimum income requirement (MIR) for partner / spouse routes — currently £29,000 (2024-04-11 uplift). Further uplift to £38,700 was proposed but paused.",
			FreshAt:       tbdPhase2,
			Source:        SourceAppendixFMSpouse,
			SchemaVersion: SchemaVersion,
			Confidence:    ConfidenceLow,
			Jurisdiction:  JurisdictionUK,
			Version:       "2024-04-11",
		},
		{
			Key:           "corpus.uk.eta_rollout_schedule",
			Description:   "UK Electronic Travel Authorisation (ETA) rollout schedule by nationality. Rolling phase-in 2025-2026 under Nationality and Borders Act 2022.",
			FreshAt:       tbdPhase2,
			Source:        SourceETARollout,
			SchemaVersion: SchemaVersion,
			Confidence:    ConfidenceLow,
			Jurisdiction:  JurisdictionUK,
			Version:       "2025-rollout",
		},
		{
			Key:           "regulation.uk.iaa_1999_section_84",
			Description:   "Immigration and Asylum Act 1999 s.84 — unregulated immigration advice is a criminal offence. OISC regulates immigration advisers under s.83.",
			FreshAt:       scaffold,
			Source:        SourceOISC,
			SchemaVersion: SchemaVersion,
			Confidence:    ConfidenceHigh,
			Jurisdiction:  JurisdictionUK,
			Version:       "1999-c.33",
		},
		{
			Key:           "regulation.uk.nationality_borders_act_2022",
			Description:   "Nationality and Borders Act 2022. Inserts UKBA 2007 s.11A enabling ETA scheme.",
			FreshAt:       scaffold,
			Source:        SourceHomeOfficeImmRules,
			SchemaVersion: SchemaVersion,
			Confidence:    ConfidenceHigh,
			Jurisdiction:  JurisdictionUK,
			Version:       "2022-c.36",
		},
		{
			Key:           "cohort.l43.mirrormark_v1",
			Description:   "L43 Mirror-Mark v1 receipt algorithm byte-identical to foundation/pkg/mirrormark.",
			FreshAt:       scaffold,
			Source:        SourceMethodologyCorpusPkg,
			SchemaVersion: SchemaVersion,
			Confidence:    ConfidenceHigh,
			Jurisdiction:  JurisdictionNone,
			Version:       "v1",
		},
		{
			Key:           "cohort.r151.kat1_canonical_hex",
			Description:   "R151 KAT-1 cross-substrate hex anchor: 239a7d0d3f1bbe3a98aede01e2ad818c2db60b7177c02e2f015035b2b5b7dbca.",
			FreshAt:       scaffold,
			Source:        SourceMethodologyCorpusPkg,
			SchemaVersion: SchemaVersion,
			Confidence:    ConfidenceHigh,
			Jurisdiction:  JurisdictionNone,
			Version:       "v1",
		},
		{
			Key:           "placeholder.counsel_review_status",
			Description:   "R166 LIABILITY-FOOTER-CONST honest-default: ReviewedByCounsel = false. Phase-1 templates have NOT been reviewed by qualified immigration counsel or OISC-regulated advisers.",
			FreshAt:       scaffold,
			Source:        SourceContextDoc,
			SchemaVersion: SchemaVersion,
			Confidence:    ConfidenceLow,
			Jurisdiction:  JurisdictionNone,
			Version:       "phase-1",
		},
		{
			Key:           "r85.parity.code_vs_context",
			Description:   "R85 CLEAN-PARITY anchor — CONTEXT.md status row vs runtime ground truth.",
			FreshAt:       scaffold,
			Source:        SourceR85ParityMarker,
			SchemaVersion: SchemaVersion,
			Confidence:    ConfidenceHigh,
			Jurisdiction:  JurisdictionNone,
			Version:       "v1",
		},
	}
}
