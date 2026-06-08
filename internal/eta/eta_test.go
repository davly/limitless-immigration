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
