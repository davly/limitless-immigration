# limitless-immigration

**One-line:** UK Home Office Immigration Rules + Appendix FM + ETA compliance forge — corpus-pinned methodology with R166 OISC liability-footer escape.

**Category:** B2B Enterprise | RegTech / Immigration Compliance
**Target Market:** OISC-regulated immigration advisers, immigration solicitors, in-house counsel at multinationals with UK sponsor licences, education-sector visa compliance officers.
**Trinity Engines:** Reality (primary — corpus-pinned rules), Causal (secondary — eligibility-path attribution), Parallax (secondary — cross-route contradiction).

**Status:** Phase-1 scaffold (2026-05-28 marathon I51 batch). R174 5-of-5 cohort maturity FROM INCEPTION.

---

## Problem Statement

UK Home Office Immigration Rules HC 395 (as amended) is the most amended primary instrument in UK administrative law — Statements of Changes are issued ~quarterly and routinely run to hundreds of pages. Three additional moving targets compound the problem:

- **Appendix FM Minimum Income Requirement** uplifted from £18,600 to £29,000 (2024-04-11); further uplift to £38,700 was proposed but paused post-2024 General Election.
- **Electronic Travel Authorisation (ETA)** is rolling out by nationality through 2025-2026 under Nationality and Borders Act 2022.
- **Skilled Worker route shortage occupation list** and going-rate thresholds drift independently.

Every immigration compliance vendor either (a) implements ONE route and abandons the others, OR (b) implements all routes with NO corpus-version pin, leaving silent enforcement of superseded paragraphs as the failure mode.

The fundamental problem: a Home Office or appeal-tribunal audit cannot answer the question *which Immigration Rules edition was scored against, and is it still current?*

---

## R166 OISC Liability-Footer (load-bearing)

Every determination emitted by this software carries the R166 LIABILITY-FOOTER-CONST escape:

> **NOT LEGAL ADVICE.** Immigration advice is regulated in the UK by the OISC (Office of the Immigration Services Commissioner) under Immigration and Asylum Act 1999 s.83. Unregulated immigration advice is a criminal offence under IAA 1999 s.84/s.91. This software ships a Phase-1 scaffold and is NOT a substitute for OISC-regulated advice or solicitor review.

This is the OISC firewall. The library `Limitless.Immigration.Legal.LibraryRecommendsHostActs = false` declares that downstream hosts MUST surface the footer before acting on any determination.

---

## R174 5-of-5 cohort maturity (strict from inception)

- **L43 Mirror-Mark v1** — `internal/mirrormark/`. KAT-1 hex `239a7d0d3f1bbe3a98aede01e2ad818c2db60b7177c02e2f015035b2b5b7dbca` pinned + OpenSSL-reproducible.
- **R143 LOUD-ONCE-WARNING-FLAG** — `internal/honest/` with 5 LIMITLESS_IMMIGRATION_* advisories (3 Error + 2 Warn per R143.A).
- **R145.C FIREWALL-TEST-DISCIPLINE** — `internal/firewall/` with on-disk drift detection.
- **R150 PARALLEL-MAP review-metadata** — `internal/manifest/` with FreshAt + Source + SchemaVersion + Confidence + Jurisdiction + Version (Class-3 anchor).
- **R151 KAT-AS-COHORT-INVARIANT-PIN** — KAT-1 byte-identical to every cohort substrate.

---

## Phase-2 deferred backlog

See CONTEXT.md "Phase-2 deferred backlog" section. Each Phase-2 surface = separate M-slot per R145.B SIBLING-NOT-STACKED.

---

## License

Apache-2.0.
