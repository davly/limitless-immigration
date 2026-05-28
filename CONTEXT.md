# CONTEXT — limitless-immigration

**Repo:** `C:\limitless\flagships\limitless-immigration`
**GitHub:** https://github.com/davly/limitless-immigration (public, Apache-2.0)
**Substrate:** Go 1.22
**Status:** Phase 1 scaffold (I51 marathon 2026-05-28 batch of 6 deferred NEW flagships).
**Cohort posture:** R174 5-of-5 cohort maturity FROM INCEPTION.

---

## What this flagship is

A UK immigration compliance scaffold that pins:

1. **UK Home Office Immigration Rules HC 395** (as amended via Statements of Changes ~quarterly).
2. **Appendix FM Financial Requirement** (£29,000 minimum income from 2024-04-11).
3. **Electronic Travel Authorisation (ETA) rollout schedule** (Nationality and Borders Act 2022).

PLUS the cohort 5-of-5 invariants — Mirror-Mark v1 / R143 LoudOnce / R145.C firewall / R150 manifest / R151 KAT-1 anchor.

---

## R174 5-of-5 cohort maturity (strict from inception)

| Package | Discipline |
|---|---|
| `internal/firewall/` | R145.C FIREWALL-TEST-DISCIPLINE |
| `internal/honest/` | R143 LOUD-ONCE-WARNING-FLAG (5 advisories) |
| `internal/mirrormark/` | L43 Mirror-Mark v1 + R151 KAT-1 anchor |
| `internal/manifest/` | R150 PARALLEL-MAP review-metadata envelope |
| `internal/legal/` | R166 LIABILITY-FOOTER-CONST + REVIEWED-BY-COUNSEL-FALSE |

Plus 1 domain-gate package:

| Package | Domain |
|---|---|
| `internal/immigration-rules/` | UK Home Office Immigration Rules + Appendix FM + ETA corpus SHA pins (THE MOAT) |

---

## R143 Advisories surfaced

5 LIMITLESS_IMMIGRATION_* advisories shipped at module-load via `internal/honest/`:

1. `LIMITLESS_IMMIGRATION_HC_395_VERSION_PIN_REQUIRED` (Error)
2. `LIMITLESS_IMMIGRATION_APPENDIX_FM_FINANCIAL_REQUIREMENT_FROZEN` (Error)
3. `LIMITLESS_IMMIGRATION_DECISION_NOT_LEGAL_ADVICE` (Error)
4. `LIMITLESS_IMMIGRATION_ETA_ROLLOUT_INCOMPLETE` (Warn)
5. `LIMITLESS_IMMIGRATION_REVIEWED_BY_COUNSEL_FALSE` (Warn)

R143.A SEVERITY-LADDER: 3 Error + 2 Warn (cohort-canonical shape).

---

## Phase-2 deferred backlog (honest disclosure)

The following Phase-2+ surfaces are NOT shipped — each is a separate M-slot:

1. **HC 395 paragraph index ingestion** — full Statement of Changes parse + per-paragraph SHA. Phase 2.
2. **Appendix FM-SE financial evidence calculator** — full evidence-bundle eligibility scoring (savings, employment, self-employment combinations). Phase 2.
3. **Skilled Worker route SOC code + going-rate eligibility engine** — Appendix Skilled Worker + Appendix Skilled Occupations. Phase 2.
4. **ETA per-nationality status engine** — nationality-to-ETA-requirement lookup with effective-from dates. Phase 2.
5. **Appendix EU + EU Settlement Scheme status tracking** — pre-Brexit status preservation. Phase 3.
6. **Counsel review of placeholder narrative templates** — Phase 1 ships R166 LIABILITY-FOOTER-CONST with `ReviewedByCounsel = false`.

The Phase-1 scaffold ships ONLY: corpus-SHA pinning + R143 advisories + Mirror-Mark + R150 manifest + R145.C firewall + R166 footer + KAT-1 R151 anchor.

---

## R85 CLEAN-PARITY anchor

This CONTEXT.md is the canonical doc-comment anchor. Divergence between this status row and runtime ground truth = R85 violation.
