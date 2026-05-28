package mirrormark

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"strings"
	"testing"
)

// Cohort-canonical KAT-1 mark literal. Byte-identical to every cohort
// Go port. limitless-immigration joins as the UK Home Office Immigration
// Rules + Appendix FM + ETA substrate.
const kat1Mark = "lore@v1:AAAAAAAAAAAjmn0NPxu-Opiu3gHirYGMLbYLcXfALi8BUDWytbfbyg"

// KAT-1 HMAC-SHA256 digest hex. R151 KAT-AS-COHORT-INVARIANT-CROSS-
// SUBSTRATE-PIN canonical anchor.
const kat1DigestHex = "239a7d0d3f1bbe3a98aede01e2ad818c2db60b7177c02e2f015035b2b5b7dbca"

// TestVerify_KAT1Mark — cohort substrate-parity oracle.
func TestVerify_KAT1Mark(t *testing.T) {
	var zeroCorpus [sha256.Size]byte
	if err := Verify(kat1Mark, zeroCorpus, []byte{}, []byte{}); err != nil {
		t.Fatalf("KAT-1 cohort literal failed Verify: %v\n\nThe limitless-immigration mirrormark algorithm has drifted from the cohort.", err)
	}
}

// TestSign_ProducesKAT1Mark — Sign reproduces published literal.
func TestSign_ProducesKAT1Mark(t *testing.T) {
	var zeroCorpus [sha256.Size]byte
	got := Sign(zeroCorpus, []byte{}, []byte{})
	if got != kat1Mark {
		t.Fatalf("Sign for KAT-1 input drift:\n  got:  %q\n  want: %q", got, kat1Mark)
	}
}

// TestKAT1Digest_EmbeddedInKAT1Mark — connects mark literal -> OpenSSL.
func TestKAT1Digest_EmbeddedInKAT1Mark(t *testing.T) {
	encoded := strings.TrimPrefix(kat1Mark, MarkPrefix)
	body, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("KAT-1 mark body not valid base64.RawURLEncoding: %v", err)
	}
	if len(body) != MarkBodyLen {
		t.Fatalf("KAT-1 body length: got %d want %d", len(body), MarkBodyLen)
	}
	gotDigestHex := hex.EncodeToString(body[MarkCorpusPrefixLen:])
	if gotDigestHex != kat1DigestHex {
		t.Fatalf("KAT-1 embedded-digest drift:\n  got:      %s\n  expected: %s", gotDigestHex, kat1DigestHex)
	}
}

// TestSign_RoundtripVerify — happy path.
func TestSign_RoundtripVerify(t *testing.T) {
	var corpus [sha256.Size]byte
	for i := range corpus {
		corpus[i] = byte(i)
	}
	key := []byte("immigration_test_key")
	payload := []byte(`{"applicant":"AB123456","route":"Skilled Worker"}`)
	mark := Sign(corpus, payload, key)
	if err := Verify(mark, corpus, payload, key); err != nil {
		t.Fatalf("Verify rejected fresh mark: %v", err)
	}
}

// TestVerify_RejectsMissingPrefix.
func TestVerify_RejectsMissingPrefix(t *testing.T) {
	var corpus [sha256.Size]byte
	err := Verify("not-a-mark", corpus, []byte{}, []byte("k"))
	if err != ErrUnknownMarkVersion {
		t.Fatalf("got %v, want ErrUnknownMarkVersion", err)
	}
}

// TestVerify_RejectsMalformedBase64.
func TestVerify_RejectsMalformedBase64(t *testing.T) {
	var corpus [sha256.Size]byte
	err := Verify("lore@v1:!!!not-base64!!!", corpus, []byte{}, []byte("k"))
	if err != ErrMalformedMark {
		t.Fatalf("got %v, want ErrMalformedMark", err)
	}
}

// TestVerify_RejectsWrongCorpus.
func TestVerify_RejectsWrongCorpus(t *testing.T) {
	var corpusA, corpusB [sha256.Size]byte
	for i := range corpusA {
		corpusA[i] = 0x11
		corpusB[i] = 0x22
	}
	mark := Sign(corpusA, []byte("p"), []byte("k"))
	err := Verify(mark, corpusB, []byte("p"), []byte("k"))
	if err != ErrCorpusMismatch {
		t.Fatalf("got %v, want ErrCorpusMismatch", err)
	}
}

// TestVerify_RejectsTamperedPayload.
func TestVerify_RejectsTamperedPayload(t *testing.T) {
	var corpus [sha256.Size]byte
	for i := range corpus {
		corpus[i] = 0x44
	}
	mark := Sign(corpus, []byte("original"), []byte("k"))
	err := Verify(mark, corpus, []byte("tampered"), []byte("k"))
	if err != ErrSignatureMismatch {
		t.Fatalf("got %v, want ErrSignatureMismatch", err)
	}
}

// TestVerify_RejectsTamperedKey.
func TestVerify_RejectsTamperedKey(t *testing.T) {
	var corpus [sha256.Size]byte
	for i := range corpus {
		corpus[i] = 0x55
	}
	mark := Sign(corpus, []byte("p"), []byte("alice"))
	err := Verify(mark, corpus, []byte("p"), []byte("bob"))
	if err != ErrSignatureMismatch {
		t.Fatalf("got %v, want ErrSignatureMismatch", err)
	}
}

// TestMarkLength_FixedAt62.
func TestMarkLength_FixedAt62(t *testing.T) {
	var corpus [sha256.Size]byte
	for i := range corpus {
		corpus[i] = byte(i * 3)
	}
	mark := Sign(corpus, []byte("anything"), []byte("k"))
	if len(mark) != 62 {
		t.Fatalf("Mark length: got %d, want 62", len(mark))
	}
}

// TestMarkPrefix_Pinned.
func TestMarkPrefix_Pinned(t *testing.T) {
	if MarkPrefix != "lore@v1:" {
		t.Fatalf("MarkPrefix drift: %q", MarkPrefix)
	}
}
