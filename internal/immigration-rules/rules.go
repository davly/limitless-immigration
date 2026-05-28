// Package immigrationrules holds the THE-MOAT corpus pin for UK Home
// Office Immigration Rules (HC 395 as amended), Appendix FM, and the
// Electronic Travel Authorisation rollout schedule.
//
// PHASE-1 SCAFFOLD NOTE — the SHA values below are PLACEHOLDER derived-
// by-pinning literals. Phase-2 binds them to canonical gov.uk-published
// artefacts.
package immigrationrules

import (
	"crypto/sha256"
	"encoding/hex"
)

// CorpusID identifies which UK immigration corpus a determination was
// scored against.
type CorpusID string

const (
	// CorpusHomeOfficeImmRulesHC395 identifies UK Home Office Immigration
	// Rules HC 395 (as amended). Statement-of-Changes cycle ~quarterly.
	CorpusHomeOfficeImmRulesHC395 CorpusID = "uk_home_office_imm_rules_hc395_as_amended_2024"

	// CorpusAppendixFMFinancialRequirement identifies Appendix FM
	// (Family Members) + FM-SE (Financial Evidence) financial
	// requirement scaffold.
	CorpusAppendixFMFinancialRequirement CorpusID = "uk_appendix_fm_financial_requirement_2024_04_11"

	// CorpusETARollout identifies the Electronic Travel Authorisation
	// rollout schedule by nationality (Nationality and Borders Act 2022).
	CorpusETARollout CorpusID = "uk_eta_rollout_2025_phased"
)

// HomeOfficeImmRulesHC395SHA — PLACEHOLDER 32-byte corpus SHA pin.
var HomeOfficeImmRulesHC395SHA = sha256.Sum256([]byte(string(CorpusHomeOfficeImmRulesHC395)))

// AppendixFMFinancialRequirementSHA — PLACEHOLDER 32-byte corpus SHA pin.
var AppendixFMFinancialRequirementSHA = sha256.Sum256([]byte(string(CorpusAppendixFMFinancialRequirement)))

// ETARolloutSHA — PLACEHOLDER 32-byte corpus SHA pin.
var ETARolloutSHA = sha256.Sum256([]byte(string(CorpusETARollout)))

type CorpusPin struct {
	ID  CorpusID
	SHA [sha256.Size]byte
}

func (p CorpusPin) HexSHA() string {
	return hex.EncodeToString(p.SHA[:])
}

func (p CorpusPin) PrefixHex() string {
	return hex.EncodeToString(p.SHA[:8])
}

func AllPins() []CorpusPin {
	return []CorpusPin{
		{ID: CorpusAppendixFMFinancialRequirement, SHA: AppendixFMFinancialRequirementSHA},
		{ID: CorpusETARollout, SHA: ETARolloutSHA},
		{ID: CorpusHomeOfficeImmRulesHC395, SHA: HomeOfficeImmRulesHC395SHA},
	}
}

func PinByID(id CorpusID) (CorpusPin, bool) {
	for _, p := range AllPins() {
		if p.ID == id {
			return p, true
		}
	}
	return CorpusPin{}, false
}
