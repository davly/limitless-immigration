package honest

import (
	"bytes"
	"strings"
	"testing"
)

func TestLoudOncePrefix(t *testing.T) {
	if LoudOncePrefix != "[LOUD-ONCE-WARNING]" {
		t.Fatalf("LoudOncePrefix drift: %q", LoudOncePrefix)
	}
}

func TestLoudOnce_EmitsOnFirstCall(t *testing.T) {
	Reset()
	var buf bytes.Buffer
	LoudOnce(Advisory{Code: "TEST_A", Severity: SeverityInfo, Message: "msg", DocLink: "d"}, &buf)
	if !strings.Contains(buf.String(), "TEST_A") {
		t.Errorf("missing Code: %q", buf.String())
	}
}

func TestLoudOnce_SilentOnSubsequent(t *testing.T) {
	Reset()
	var buf bytes.Buffer
	adv := Advisory{Code: "TEST_ONCE", Severity: SeverityInfo, Message: "m", DocLink: "d"}
	LoudOnce(adv, &buf)
	buf.Reset()
	LoudOnce(adv, &buf)
	LoudOnce(adv, &buf)
	if buf.Len() != 0 {
		t.Fatalf("leaked: %q", buf.String())
	}
}

func TestCanonicalAdvisories_Count5(t *testing.T) {
	if len(CanonicalAdvisories()) != 5 {
		t.Fatalf("got %d, want 5", len(CanonicalAdvisories()))
	}
}

func TestCanonicalAdvisories_AllFieldsNonEmpty(t *testing.T) {
	for i, a := range CanonicalAdvisories() {
		if a.Code == "" || a.Severity == "" || a.Message == "" || a.DocLink == "" {
			t.Errorf("advisory %d: empty field", i)
		}
	}
}

func TestCanonicalAdvisories_UniqueCodes(t *testing.T) {
	seen := map[string]int{}
	for i, a := range CanonicalAdvisories() {
		if prev, ok := seen[a.Code]; ok {
			t.Errorf("dup %q at %d and %d", a.Code, prev, i)
		}
		seen[a.Code] = i
	}
}

func TestCanonicalAdvisories_AllStartWithLIMITLESS_IMMIGRATION(t *testing.T) {
	for _, a := range CanonicalAdvisories() {
		if !strings.HasPrefix(a.Code, "LIMITLESS_IMMIGRATION_") {
			t.Errorf("Code %q missing prefix", a.Code)
		}
	}
}

func TestCanonicalAdvisories_R143A_Ladder_3Error_2Warn(t *testing.T) {
	var nE, nW int
	for _, a := range CanonicalAdvisories() {
		switch a.Severity {
		case SeverityError:
			nE++
		case SeverityWarn:
			nW++
		}
	}
	if nE != 3 || nW != 2 {
		t.Fatalf("severity ladder drift: %dE %dW (want 3E 2W)", nE, nW)
	}
}

func TestFindAdvisory_ByCanonicalCode(t *testing.T) {
	for _, expected := range CanonicalAdvisories() {
		got, ok := FindAdvisory(expected.Code)
		if !ok || got.Code != expected.Code {
			t.Errorf("FindAdvisory(%q) drift", expected.Code)
		}
	}
}

func TestFindAdvisory_UnknownCode(t *testing.T) {
	_, ok := FindAdvisory("DOES_NOT_EXIST")
	if ok {
		t.Fatal("got ok=true for unknown")
	}
}

func TestReset_ClearsRegistry(t *testing.T) {
	Reset()
	var buf bytes.Buffer
	adv := Advisory{Code: "RST", Severity: SeverityInfo, Message: "m", DocLink: "d"}
	LoudOnce(adv, &buf)
	first := buf.String()
	buf.Reset()
	LoudOnce(adv, &buf)
	if buf.Len() != 0 {
		t.Fatal("expected silent on second call")
	}
	Reset()
	LoudOnce(adv, &buf)
	if buf.String() != first {
		t.Error("post-Reset emission drift")
	}
}
