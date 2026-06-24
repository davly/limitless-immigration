// Package eta is a per-nationality UK Electronic Travel Authorisation (ETA)
// eligibility-SIGNAL engine. It is deliberately NOT a determination.
//
// ETA-required nationality lists roll out wave-by-wave through 2025-2026 and
// drift (Nationality and Borders Act 2022; honest advisory #3 in this repo), so
// every signal ships Confidence=Low, the CorpusETARollout corpus pin for
// provenance, the R166 liability footer, and a cold-verify caveat. A caller must
// re-source the answer against gov.uk before relying on it.
//
// This activates the CorpusETARollout pin (defined in immigration-rules but read
// by no logic until now) and builds CONTEXT.md phase-2 backlog item #4
// ("ETA per-nationality status engine — nationality-to-ETA-requirement lookup
// with effective-from dates").
package eta

import (
	"strings"
	"time"

	immigrationrules "github.com/davly/limitless-immigration/internal/immigration-rules"
	"github.com/davly/limitless-immigration/internal/legal"
	"github.com/davly/limitless-immigration/internal/manifest"
)

// ETAStatus is the closed set of ETA signals. A SIGNAL, never a determination.
type ETAStatus string

const (
	StatusRequired      ETAStatus = "REQUIRED"
	StatusNotRequired   ETAStatus = "NOT_REQUIRED"
	StatusNotYetInForce ETAStatus = "NOT_YET_IN_FORCE"
	StatusUnknown       ETAStatus = "UNKNOWN"
)

// ETARule is one published rollout wave. Required=false marks a nationality
// exempt from ETA (e.g. British/Irish citizens), for which the date is moot.
// Seed rules are ILLUSTRATIVE and NON-AUTHORITATIVE — cold-verify against gov.uk.
type ETARule struct {
	NationalityISO3 string
	EffectiveFrom   time.Time
	Required        bool
}

// ETASignal is the structured, non-authoritative output. The Footer and Caveat
// are ALWAYS populated — even for an UNKNOWN nationality.
type ETASignal struct {
	NationalityISO3 string
	TravelDate      time.Time
	Status          ETAStatus
	EffectiveFrom   time.Time // zero when exempt or unknown
	Confidence      manifest.Confidence
	Jurisdiction    manifest.Jurisdiction
	CorpusPinPrefix string
	Footer          string
	Caveat          string
}

const signalCaveat = "SIGNAL ONLY, not a determination. UK ETA nationality waves drift; cold-verify against gov.uk/guidance/apply-for-an-electronic-travel-authorisation-eta for the travel date before reliance."

func corpusPinPrefix() string {
	if p, ok := immigrationrules.PinByID(immigrationrules.CorpusETARollout); ok {
		return p.PrefixHex()
	}
	return ""
}

// Classify returns the ETA signal for a nationality (ISO 3166-1 alpha-3) at a
// travel date, and ok=false when the nationality is absent from the table
// (UNKNOWN — the engine never guesses). The returned signal ALWAYS carries the
// liability footer + corpus pin, including the unknown case.
func Classify(nationalityISO3 string, travelDate time.Time, table []ETARule) (ETASignal, bool) {
	iso := strings.ToUpper(strings.TrimSpace(nationalityISO3))
	sig := ETASignal{
		NationalityISO3: iso,
		TravelDate:      travelDate,
		Status:          StatusUnknown,
		Confidence:      manifest.ConfidenceLow,
		Jurisdiction:    manifest.JurisdictionUK,
		CorpusPinPrefix: corpusPinPrefix(),
		Footer:          legal.LegalLiabilityFooter,
		Caveat:          signalCaveat,
	}
	for _, r := range table {
		if strings.ToUpper(strings.TrimSpace(r.NationalityISO3)) != iso {
			continue
		}
		if !r.Required {
			sig.Status = StatusNotRequired // exempt nationality; date irrelevant
			return sig, true
		}
		sig.EffectiveFrom = r.EffectiveFrom
		// EffectiveFrom models a published gov.uk CALENDAR wave date (e.g. via
		// time.Parse("2006-01-02")), so the in-force test is a CALENDAR-DATE compare,
		// not a raw-instant compare. Reducing both operands to calendar-date
		// granularity keeps "effective-from is inclusive" honest for any caller —
		// including one whose travelDate carries a non-midnight clock or a non-UTC
		// location. Without this, a traveller whose calendar travel date is exactly
		// the in-force date but whose underlying instant differs from the UTC-midnight
		// EffectiveFrom would be wrongly flipped to NOT_YET_IN_FORCE (or vice versa).
		if beforeCalendarDate(travelDate, r.EffectiveFrom) {
			sig.Status = StatusNotYetInForce // wave not yet in force for this nationality
		} else {
			sig.Status = StatusRequired // effective-from is inclusive (calendar date)
		}
		return sig, true
	}
	return sig, false
}

// beforeCalendarDate reports whether the calendar date of a is strictly earlier
// than the calendar date of b. Each operand is reduced to its own year/month/day
// (via time.Time.Date(), in the time's own location), so the comparison ignores
// the time-of-day component entirely — matching the calendar-date contract of a
// gov.uk-published effective-from wave date and of a traveller's calendar travel
// date. Equal calendar dates are NOT before (the inclusive boundary). This is a
// pure value comparison, allocation-free and race-safe.
func beforeCalendarDate(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	if ay != by {
		return ay < by
	}
	if am != bm {
		return am < bm
	}
	return ad < bd
}

// SeedTable is an ILLUSTRATIVE, NON-AUTHORITATIVE sample of rollout waves for
// demos/tests only. Hosts MUST replace it with a cold-verified, gov.uk-sourced
// table; the dates here are placeholders, not legal facts.
func SeedTable() []ETARule {
	d := func(s string) time.Time { t, _ := time.Parse("2006-01-02", s); return t }
	return []ETARule{
		{NationalityISO3: "QAT", EffectiveFrom: d("2024-11-15"), Required: true},
		{NationalityISO3: "USA", EffectiveFrom: d("2025-01-08"), Required: true},
		{NationalityISO3: "FRA", EffectiveFrom: d("2025-04-02"), Required: true},
		{NationalityISO3: "GBR", Required: false}, // British citizens: exempt
		{NationalityISO3: "IRL", Required: false}, // Irish citizens: exempt (Common Travel Area)
	}
}
