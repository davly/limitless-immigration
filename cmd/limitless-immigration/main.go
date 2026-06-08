// Command limitless-immigration — UK Home Office Immigration Rules
// + Appendix FM + ETA compliance forge CLI.
//
// Phase-1 scaffold. Ships:
//
//   - `corpus list`     — list pinned immigration-rules corpus SHAs
//   - `advisories list` — list R143 LIMITLESS_IMMIGRATION_* advisories
//   - `manifest list`   — list R150 schematised-knowledge entries
//   - `eta status`      — per-nationality ETA eligibility SIGNAL (not a determination)
//   - `version`         — print version
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/davly/limitless-immigration/internal/eta"
	"github.com/davly/limitless-immigration/internal/honest"
	immigrationrules "github.com/davly/limitless-immigration/internal/immigration-rules"
	"github.com/davly/limitless-immigration/internal/manifest"
)

const version = "0.1.0-i51-scaffold"

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	switch os.Args[1] {
	case "corpus":
		corpusCmd(os.Args[2:])
	case "advisories":
		advisoriesCmd(os.Args[2:])
	case "manifest":
		manifestCmd(os.Args[2:])
	case "eta":
		etaCmd(os.Args[2:])
	case "version":
		fmt.Println("limitless-immigration", version)
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `Usage: limitless-immigration <command>

Commands:
  corpus list     -- list pinned immigration-rules corpus SHAs
  advisories list -- list R143 advisories
  manifest list   -- list R150 schematised-knowledge entries
  eta status <ISO3> <YYYY-MM-DD> -- per-nationality ETA eligibility SIGNAL (not a determination)
  version         -- print version`)
}

func etaCmd(args []string) {
	if len(args) < 3 || args[0] != "status" {
		fmt.Fprintln(os.Stderr, "Usage: eta status <ISO3> <YYYY-MM-DD>")
		os.Exit(2)
	}
	date, err := time.Parse("2006-01-02", args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid date %q (want YYYY-MM-DD): %v\n", args[2], err)
		os.Exit(2)
	}
	sig, ok := eta.Classify(args[1], date, eta.SeedTable())
	fmt.Printf("ETA signal: %s  (nationality %s, travel %s)\n", sig.Status, sig.NationalityISO3, args[2])
	if !ok {
		fmt.Println("  note: nationality not in the seed table -> UNKNOWN. The engine never guesses.")
	}
	if !sig.EffectiveFrom.IsZero() {
		fmt.Printf("  effective-from: %s\n", sig.EffectiveFrom.Format("2006-01-02"))
	}
	fmt.Printf("  confidence: Low (illustrative seed table) | jurisdiction: %s | corpus pin: %s\n",
		sig.Jurisdiction, sig.CorpusPinPrefix)
	fmt.Printf("  %s\n", sig.Caveat)
	fmt.Printf("  %s\n", sig.Footer)
}

func corpusCmd(args []string) {
	if len(args) < 1 || args[0] != "list" {
		fmt.Fprintln(os.Stderr, "Usage: corpus list")
		os.Exit(2)
	}
	for _, p := range immigrationrules.AllPins() {
		fmt.Printf("%s\n  sha256: %s\n  prefix: %s\n", p.ID, p.HexSHA(), p.PrefixHex())
	}
}

func advisoriesCmd(args []string) {
	if len(args) < 1 || args[0] != "list" {
		fmt.Fprintln(os.Stderr, "Usage: advisories list")
		os.Exit(2)
	}
	for _, a := range honest.CanonicalAdvisories() {
		fmt.Printf("[%s] %s\n  %s\n", a.Severity, a.Code, a.Message)
	}
}

func manifestCmd(args []string) {
	if len(args) < 1 || args[0] != "list" {
		fmt.Fprintln(os.Stderr, "Usage: manifest list")
		os.Exit(2)
	}
	for _, e := range manifest.Seed() {
		fmt.Printf("%s\n  desc: %s\n  source: %s\n  jurisdiction: %s\n  version: %s\n",
			e.Key, e.Description, e.Source, e.Jurisdiction, e.Version)
	}
}
