# SECURITY — limitless-immigration

## Phase-1 scope (load-bearing)

This is a Phase-1 scaffold. The deployment MUST NOT make automated immigration eligibility determinations without:

1. **Counsel review** — `internal/legal/ReviewedByCounsel = false` is the R166 honest-default. Phase-1 templates have NOT been reviewed by qualified immigration counsel or OISC-regulated advisers.
2. **OISC compliance** — unregulated immigration advice is a criminal offence under IAA 1999 s.91. Hosts MUST be OISC-regulated or operate under solicitor supervision.
3. **R143 advisory acknowledgement** — the 5 LIMITLESS_IMMIGRATION_* advisories MUST be visible to every operator.
4. **Corpus-SHA cold-verification** — before any live determination, the corpus pins in `internal/immigration-rules/` MUST be cold-verified against gov.uk-published canonical (`HomeOfficeImmRulesHC395SHA`, `AppendixFMFinancialRequirementSHA`, `ETARolloutSHA`).

## R166 LIABILITY-FOOTER-CONST

The constant `internal/legal/LegalLiabilityFooter` is the OISC escape phrase. Every determination payload that crosses a trust boundary MUST embed it verbatim until counsel review flips `ReviewedByCounsel` to true.

## Mirror-Mark v1 tamper-evidence

Every determination payload signed with a Mirror-Mark v1 (`lore@v1:` 62-char prefix). Drift in corpus SHA / payload / key → `ErrCorpusMismatch` / `ErrSignatureMismatch`. Constant-time comparison via `hmac.Equal`.

## Threat model — what this Phase-1 scaffold DOES NOT defend against

- Compromised signing key (no KMS integration in Phase-1).
- Compromised corpus SHA (Phase-1 placeholders sha256(CorpusID-string); Phase-2 binds to gov.uk artefacts).
- Side-channel timing attacks on base64 decode or prefix compare.
- Adversarial Statement of Changes parsing (no formal HC 395 parser in Phase-1).
