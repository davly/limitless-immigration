// Package legal implements the R166 LIABILITY-FOOTER-CONST + REVIEWED-
// BY-COUNSEL-FALSE sentinel pattern for limitless-immigration.
//
// Every immigration eligibility determination emitted by this software
// crosses a regulator-grade trust boundary (OISC IAA 1999 s.84). The
// liability footer escape is load-bearing per R166.
package legal

// LegalLiabilityFooter is the R166 escape that EVERY determination MUST
// carry until counsel review flips ReviewedByCounsel to true.
const LegalLiabilityFooter = "NOT LEGAL ADVICE. Immigration advice is regulated in the UK by the OISC (Office of the Immigration Services Commissioner) under Immigration and Asylum Act 1999 s.83. Unregulated immigration advice is a criminal offence under IAA 1999 s.84/s.91. This software ships a Phase-1 scaffold and is NOT a substitute for OISC-regulated advice or solicitor review. Consult an OISC-regulated adviser or solicitor before relying on any determination."

// ReviewedByCounsel is the R166 honest-default. Phase-1 ships false;
// counsel review + flip lives on its own R145.B sibling branch.
const ReviewedByCounsel = false

// LibraryRecommendsHostActs — R166 declaration that this library
// recommends host applications act on its determinations only after the
// host has surfaced the LegalLiabilityFooter to the operator.
const LibraryRecommendsHostActs = false
