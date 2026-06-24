package eta

import (
	"testing"
	"time"

	"github.com/davly/limitless-immigration/internal/legal"
	"github.com/davly/limitless-immigration/internal/manifest"
)

func d(s string) time.Time { t, _ := time.Parse("2006-01-02", s); return t }

func table() []ETARule {
	return []ETARule{
		{NationalityISO3: "USA", EffectiveFrom: d("2025-01-08"), Required: true},
		{NationalityISO3: "GBR", Required: false},
	}
}

func TestEffectiveFromBoundary(t *testing.T) {
	// USA wave is in force from 2025-01-08 INCLUSIVE.
	cases := []struct {
		date string
		want ETAStatus
	}{
		{"2025-01-07", StatusNotYetInForce}, // day before
		{"2025-01-08", StatusRequired},      // day of (inclusive)
		{"2025-01-09", StatusRequired},      // day after
	}
	for _, c := range cases {
		sig, ok := Classify("USA", d(c.date), table())
		if !ok || sig.Status != c.want {
			t.Errorf("USA @ %s = (%s, ok=%v), want %s", c.date, sig.Status, ok, c.want)
		}
	}
}

func TestEffectiveFromBoundaryIsCalendarDateNotInstant(t *testing.T) {
	// EffectiveFrom models a gov.uk-published CALENDAR wave date (UTC-midnight via
	// time.Parse("2006-01-02")). The in-force test must compare CALENDAR DATES, not
	// raw instants — otherwise a caller east of UTC (or any non-midnight clock) on
	// the very in-force calendar day gets flipped to NOT_YET_IN_FORCE, contradicting
	// the eta.go "effective-from is inclusive" contract.
	//
	// USA wave: in force from the 2025-01-08 CALENDAR day, inclusive.
	tokyo, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatalf("load Asia/Tokyo: %v", err)
	}

	// Same calendar day as effective-from (2025-01-08), but the underlying instant
	// is 2025-01-07 21:00 UTC — which is .Before() the UTC-midnight EffectiveFrom.
	// On the buggy code this returns NOT_YET_IN_FORCE; correct answer is REQUIRED.
	travel := time.Date(2025, 1, 8, 6, 0, 0, 0, tokyo)
	sig, ok := Classify("USA", travel, table())
	if !ok || sig.Status != StatusRequired {
		t.Errorf("USA @ 2025-01-08 06:00 Asia/Tokyo (same calendar day as effective-from) = (%s, ok=%v), want REQUIRED (effective-from is inclusive on the calendar date)", sig.Status, ok)
	}

	// A non-midnight clock on the in-force day itself must also be REQUIRED.
	travelUTCAfternoon := time.Date(2025, 1, 8, 15, 30, 0, 0, time.UTC)
	if sig, ok := Classify("USA", travelUTCAfternoon, table()); !ok || sig.Status != StatusRequired {
		t.Errorf("USA @ 2025-01-08 15:30 UTC = (%s, ok=%v), want REQUIRED", sig.Status, ok)
	}

	// Boundary preserved on the other side: the genuine day-before calendar date
	// must still read NOT_YET_IN_FORCE, even when zoned west of UTC such that the
	// instant could appear "on or after" midnight.
	la, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		t.Fatalf("load America/Los_Angeles: %v", err)
	}
	// 2025-01-07 20:00 Los Angeles == 2025-01-08 04:00 UTC. The CALENDAR travel
	// date is 2025-01-07 (day before in-force) so the signal must be NOT_YET_IN_FORCE
	// even though the raw instant is after UTC-midnight effective-from.
	travelDayBeforeWest := time.Date(2025, 1, 7, 20, 0, 0, 0, la)
	if sig, ok := Classify("USA", travelDayBeforeWest, table()); !ok || sig.Status != StatusNotYetInForce {
		t.Errorf("USA @ 2025-01-07 20:00 America/Los_Angeles (calendar day before) = (%s, ok=%v), want NOT_YET_IN_FORCE", sig.Status, ok)
	}
}

func TestExemptNationality(t *testing.T) {
	// British citizens are exempt regardless of date.
	for _, date := range []string{"2020-01-01", "2030-01-01"} {
		sig, ok := Classify("GBR", d(date), table())
		if !ok || sig.Status != StatusNotRequired {
			t.Errorf("GBR @ %s = (%s, ok=%v), want NOT_REQUIRED", date, sig.Status, ok)
		}
	}
}

func TestUnknownNationalityNeverGuesses(t *testing.T) {
	sig, ok := Classify("XXX", d("2025-06-01"), table())
	if ok || sig.Status != StatusUnknown {
		t.Errorf("unknown nationality = (%s, ok=%v), want (UNKNOWN, false)", sig.Status, ok)
	}
	// even UNKNOWN must carry the liability footer + caveat.
	if sig.Footer != legal.LegalLiabilityFooter || sig.Caveat == "" {
		t.Error("UNKNOWN signal must still carry footer + caveat")
	}
}

func TestSignalAlwaysCarriesHonestEnvelope(t *testing.T) {
	sig, _ := Classify("usa", d("2025-06-01"), table()) // lowercase ISO accepted
	if sig.NationalityISO3 != "USA" {
		t.Errorf("ISO not normalised: %q", sig.NationalityISO3)
	}
	if sig.Footer != legal.LegalLiabilityFooter {
		t.Error("R166 liability footer must always be present")
	}
	if sig.Confidence != manifest.ConfidenceLow {
		t.Errorf("seed-table signal must be Confidence=Low, got %v", sig.Confidence)
	}
	if sig.Jurisdiction != manifest.JurisdictionUK {
		t.Errorf("jurisdiction must be UK, got %q", sig.Jurisdiction)
	}
	if sig.CorpusPinPrefix == "" {
		t.Error("signal must carry the CorpusETARollout pin prefix (provenance) — the dead pin is now load-bearing")
	}
}

func TestSeedTableClassifies(t *testing.T) {
	// the shipped illustrative table resolves its own entries.
	if sig, ok := Classify("FRA", d("2025-04-02"), SeedTable()); !ok || sig.Status != StatusRequired {
		t.Errorf("FRA @ effective-from = (%s, %v), want REQUIRED", sig.Status, ok)
	}
	if sig, ok := Classify("IRL", d("2025-04-02"), SeedTable()); !ok || sig.Status != StatusNotRequired {
		t.Errorf("IRL = (%s, %v), want NOT_REQUIRED (CTA exempt)", sig.Status, ok)
	}
}
