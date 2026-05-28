package legal

import (
	"strings"
	"testing"
)

func TestLegalLiabilityFooter_NonEmpty(t *testing.T) {
	if LegalLiabilityFooter == "" {
		t.Fatal("empty")
	}
}

func TestLegalLiabilityFooter_ContainsOISC(t *testing.T) {
	if !strings.Contains(LegalLiabilityFooter, "OISC") {
		t.Fatal("missing OISC reference")
	}
}

func TestLegalLiabilityFooter_ContainsIAA1999(t *testing.T) {
	if !strings.Contains(LegalLiabilityFooter, "IAA 1999") {
		t.Fatal("missing IAA 1999 reference")
	}
}

func TestLegalLiabilityFooter_ContainsNotLegalAdvice(t *testing.T) {
	if !strings.Contains(LegalLiabilityFooter, "NOT LEGAL ADVICE") {
		t.Fatal("missing escape phrase")
	}
}

func TestReviewedByCounsel_HonestDefaultFalse(t *testing.T) {
	if ReviewedByCounsel {
		t.Fatal("R166 honest-default: must be false in Phase-1 scaffold")
	}
}

func TestLibraryRecommendsHostActs_HonestDefaultFalse(t *testing.T) {
	if LibraryRecommendsHostActs {
		t.Fatal("R166 honest-default: must be false")
	}
}
