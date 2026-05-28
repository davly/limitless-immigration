// Package mirrormark implements the cohort L43 Mirror-Mark v1 receipt
// algorithm — byte-identical to foundation/pkg/mirrormark and every
// cohort Go port.
//
// Mark format (byte-identical to foundation/pkg/mirrormark):
//
//	"lore@v1:" + base64url( corpusSHA[:8] || hmacSHA256(0x01 || corpusSHA || payload, key) )
//
// Resulting in a fixed 62-character string: `lore@v1:` prefix (8 chars)
// + 54-char base64url body (40 raw bytes encoded).
//
// Why limitless-immigration consumes this today:
//
//   - Limitless-immigration's output (UK Home Office Immigration Rules
//     compliance assertions + Appendix FM eligibility decisions + ETA
//     status records) crosses a regulator-grade trust boundary. A
//     decision payload that arrives with a verifiable Mirror-Mark is
//     provenance-anchored — Home Office casework / appeal tribunal /
//     immigration solicitor can cold-verify the decision was not edited
//     after upstream signed it.
//   - The corpus prefix carries the Immigration-Rules-version SHA — IS
//     the moat. A regulator can cold-verify which Immigration Rules
//     edition (HC 395 as amended) the decision was scored against.
package mirrormark

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

const (
	// MarkVersion is the 1-byte tag prefixing the HMAC input.
	MarkVersion byte = 0x01
	// MarkPrefix is the documented header-value prefix.
	MarkPrefix = "lore@v1:"
	// MarkCorpusPrefixLen is the corpus-SHA prefix length (8 bytes).
	MarkCorpusPrefixLen = 8
	// MarkBodyLen is the unencoded length of the mark body (40 bytes).
	MarkBodyLen = MarkCorpusPrefixLen + sha256.Size
)

var (
	// ErrUnknownMarkVersion — mark missing canonical prefix.
	ErrUnknownMarkVersion = errors.New("mirrormark: unknown mark version (missing 'lore@v1:' prefix)")
	// ErrMalformedMark — base64url decode failed or wrong body length.
	ErrMalformedMark = errors.New("mirrormark: malformed mark (base64url decode failed or wrong body length)")
	// ErrCorpusMismatch — corpus prefix in mark != supplied corpus SHA.
	ErrCorpusMismatch = errors.New("mirrormark: corpus prefix mismatch (mark signed by different corpus)")
	// ErrSignatureMismatch — HMAC digest mismatch (payload or key wrong).
	ErrSignatureMismatch = errors.New("mirrormark: HMAC signature mismatch (payload tampered or wrong key)")
)

// Sign returns the canonical Mirror-Mark string for the given payload.
// Byte-identical to foundation/pkg/mirrormark.Sign.
func Sign(corpusSHA [sha256.Size]byte, payload []byte, key []byte) string {
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write([]byte{MarkVersion})
	_, _ = mac.Write(corpusSHA[:])
	_, _ = mac.Write(payload)
	digest := mac.Sum(nil)

	body := make([]byte, 0, MarkBodyLen)
	body = append(body, corpusSHA[:MarkCorpusPrefixLen]...)
	body = append(body, digest...)

	return MarkPrefix + base64.RawURLEncoding.EncodeToString(body)
}

// Verify cold-checks a Mirror-Mark against (corpus, payload, key).
// Returns nil on match; one of the typed sentinel errors on any failure.
// Both byte-comparisons use hmac.Equal (constant-time) — timing-safe.
func Verify(mark string, corpusSHA [sha256.Size]byte, payload []byte, key []byte) error {
	if len(mark) < len(MarkPrefix) || mark[:len(MarkPrefix)] != MarkPrefix {
		return ErrUnknownMarkVersion
	}
	body, err := base64.RawURLEncoding.DecodeString(mark[len(MarkPrefix):])
	if err != nil {
		return ErrMalformedMark
	}
	if len(body) != MarkBodyLen {
		return ErrMalformedMark
	}
	corpusPrefix := body[:MarkCorpusPrefixLen]
	digest := body[MarkCorpusPrefixLen:]
	if !hmac.Equal(corpusPrefix, corpusSHA[:MarkCorpusPrefixLen]) {
		return ErrCorpusMismatch
	}
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write([]byte{MarkVersion})
	_, _ = mac.Write(corpusSHA[:])
	_, _ = mac.Write(payload)
	want := mac.Sum(nil)
	if !hmac.Equal(digest, want) {
		return ErrSignatureMismatch
	}
	return nil
}
