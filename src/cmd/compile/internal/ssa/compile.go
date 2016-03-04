// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssa

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"
)

// Compile is the main entry point for this package.
// Compile modifies f so that on return:
//   · all Values in f map to 0 or 1 assembly instructions of the target architecture
//   · the order of f.Blocks is the order to emit the Blocks
//   · the order of b.Values is the order to emit the Values in each Block
//   · f has a non-nil regAlloc field
func Compile(f *Func) {
	// TODO: debugging - set flags to control verbosity of compiler,
	// which phases to dump IR before/after, etc.
	if f.Log() {
		f.Logf("compiling %s\n", f.Name)
	}

	// hook to print function & phase if panic happens
	phaseName := "init"
	defer func() {
		if phaseName != "" {
			err := recover()
			stack := make([]byte, 16384)
			n := runtime.Stack(stack, false)
			stack = stack[:n]
			f.Fatalf("panic during %s while compiling %s:\n\n%v\n\n%s\n", phaseName, f.Name, err, stack)
		}
	}()

	// Run all the passes
	printFunc(f)
	f.Config.HTML.WriteFunc("start", f)
	checkFunc(f)
	const logMemStats = false
	for _, p := range passes {
		if !f.Config.optimize && !p.required {
			continue
		}
		f.pass = &p
		phaseName = p.name
		if f.Log() {
			f.Logf("  pass %s begin\n", p.name)
		}
		// TODO: capture logging during this pass, add it to the HTML
		var mStart runtime.MemStats
		if logMemStats || p.mem {
			runtime.ReadMemStats(&mStart)
		}

		tStart := time.Now()
		p.fn(f)
		tEnd := time.Now()

		// Need something less crude than "Log the whole intermediate result".
		if f.Log() || f.Config.HTML != nil {
			time := tEnd.Sub(tStart).Nanoseconds()
			var stats string
			if logMemStats {
				var mEnd runtime.MemStats
				runtime.ReadMemStats(&mEnd)
				nBytes := mEnd.TotalAlloc - mStart.TotalAlloc
				nAllocs := mEnd.Mallocs - mStart.Mallocs
				stats = fmt.Sprintf("[%d ns %d allocs %d bytes]", time, nAllocs, nBytes)
			} else {
				stats = fmt.Sprintf("[%d ns]", time)
			}

			f.Logf("  pass %s end %s\n", p.name, stats)
			printFunc(f)
			f.Config.HTML.WriteFunc(fmt.Sprintf("after %s <span class=\"stats\">%s</span>", phaseName, stats), f)
		}
		if p.time || p.mem {
			// Surround timing information w/ enough context to allow comparisons.
			time := tEnd.Sub(tStart).Nanoseconds()
			if p.time {
				f.logStat("TIME(ns)", time)
			}
			if p.mem {
				var mEnd runtime.MemStats
				runtime.ReadMemStats(&mEnd)
				nBytes := mEnd.TotalAlloc - mStart.TotalAlloc
				nAllocs := mEnd.Mallocs - mStart.Mallocs
				f.logStat("TIME(ns):BYTES:ALLOCS", time, nBytes, nAllocs)
			}
		}
		checkFunc(f)
	}

	// Squash error printing defer
	phaseName = ""
}

type pass struct {
	name     string
	fn       func(*Func)
	required bool
	disabled bool
	time     bool // report time to run pass
	mem      bool // report mem stats to run pass
	stats    int  // pass reports own "stats" (e.g., branches removed)
	debug    int  // pass performs some debugging. =1 should be in error-testing-friendly Warnl format.
	test     int  // pass-specific ad-hoc option, perhaps useful in development
}

// PhaseOption sets the specified flag in the specified ssa phase,
// returning empty string if this was successful or a string explaining
// the error if it was not. A version of the phase name with "_"
// replaced by " " is also checked for a match.
// See gc/lex.go for dissection of the option string. Example use:
// GO_GCFLAGS=-d=ssa/generic_cse/time,ssa/generic_cse/stats,ssa/generic_cse/debug=3 ./make.bash ...
//
func PhaseOption(phase, flag string, val int) string {
	underphase := strings.Replace(phase, "_", " ", -1)
	for i, p := range passes {
		if p.name == phase || p.name == underphase {
			switch flag {
			case "on":
				p.disabled = val == 0
			case "off":
				p.disabled = val != 0
			case "time":
				p.time = val != 0
			case "mem":
				p.mem = val != 0
			case "debug":
				p.debug = val
			case "stats":
				p.stats = val
			case "test":
				p.test = val
			default:
				return fmt.Sprintf("Did not find a flag matching %s in -d=ssa/%s debug option", flag, phase)
			}
			if p.disabled && p.required {
				return fmt.Sprintf("Cannot disable required SSA phase %s using -d=ssa/%s debug option", phase, phase)
			}
			passes[i] = p
			return ""
		}
	}
	return fmt.Sprintf("Did not find a phase matching %s in -d=ssa/... debug option", phase)
}

// list of passes for the compiler
var passes = [...]pass{
	// TODO: combine phielim and copyelim into a single pass?
	{name: "early phielim", fn: phielim},
	{name: "early copyelim", fn: copyelim},
	{name: "early deadcode", fn: deadcode}, // remove generated dead code to avoid doing pointless work during opt
	{name: "short circuit", fn: shortcircuit},
	{name: "decompose user", fn: decomposeUser, required: true},
	{name: "decompose builtin", fn: decomposeBuiltIn, required: true},
	{name: "opt", fn: opt, required: true},           // TODO: split required rules and optimizing rules
	{name: "zero arg cse", fn: zcse, required: true}, // required to merge OpSB values
	{name: "opt deadcode", fn: deadcode},             // remove any blocks orphaned during opt
	{name: "generic cse", fn: cse},
	{name: "phiopt", fn: phiopt},
	{name: "nilcheckelim", fn: nilcheckelim},
	{name: "prove", fn: prove},
	{name: "generic deadcode", fn: deadcode},
	{name: "fuse", fn: fuse},
	{name: "dse", fn: dse},
	{name: "tighten", fn: tighten}, // move values closer to their uses
	{name: "lower", fn: lower, required: true},
	{name: "lowered cse", fn: cse},
	{name: "lowered deadcode", fn: deadcode, required: true},
	{name: "checkLower", fn: checkLower, required: true},
	{name: "late phielim", fn: phielim},
	{name: "late copyelim", fn: copyelim},
	{name: "late deadcode", fn: deadcode},
	{name: "critical", fn: critical, required: true}, // remove critical edges
	{name: "likelyadjust", fn: likelyadjust},
	{name: "layout", fn: layout, required: true},       // schedule blocks
	{name: "schedule", fn: schedule, required: true},   // schedule values
	{name: "flagalloc", fn: flagalloc, required: true}, // allocate flags register
	{name: "regalloc", fn: regalloc, required: true},   // allocate int & float registers + stack slots
	{name: "trim", fn: trim},                           // remove empty blocks
}

// Double-check phase ordering constraints.
// This code is intended to document the ordering requirements
// between different phases. It does not override the passes
// list above.
type constraint struct {
	a, b string // a must come before b
}

var passOrder = [...]constraint{
	// prove reliese on common-subexpression elimination for maximum benefits.
	{"generic cse", "prove"},
	// deadcode after prove to eliminate all new dead blocks.
	{"prove", "generic deadcode"},
	// common-subexpression before dead-store elim, so that we recognize
	// when two address expressions are the same.
	{"generic cse", "dse"},
	// cse substantially improves nilcheckelim efficacy
	{"generic cse", "nilcheckelim"},
	// allow deadcode to clean up after nilcheckelim
	{"nilcheckelim", "generic deadcode"},
	// nilcheckelim generates sequences of plain basic blocks
	{"nilcheckelim", "fuse"},
	// nilcheckelim relies on opt to rewrite user nil checks
	{"opt", "nilcheckelim"},
	// tighten should happen before lowering to avoid splitting naturally paired instructions such as CMP/SET
	{"tighten", "lower"},
	// tighten will be most effective when as many values have been removed as possible
	{"generic deadcode", "tighten"},
	{"generic cse", "tighten"},
	// don't run optimization pass until we've decomposed builtin objects
	{"decompose builtin", "opt"},
	// don't layout blocks until critical edges have been removed
	{"critical", "layout"},
	// regalloc requires the removal of all critical edges
	{"critical", "regalloc"},
	// regalloc requires all the values in a block to be scheduled
	{"schedule", "regalloc"},
	// checkLower must run after lowering & subsequent dead code elim
	{"lower", "checkLower"},
	{"lowered deadcode", "checkLower"},
	// flagalloc needs instructions to be scheduled.
	{"schedule", "flagalloc"},
	// regalloc needs flags to be allocated first.
	{"flagalloc", "regalloc"},
	// trim needs regalloc to be done first.
	{"regalloc", "trim"},
}

func init() {
	for _, c := range passOrder {
		a, b := c.a, c.b
		i := -1
		j := -1
		for k, p := range passes {
			if p.name == a {
				i = k
			}
			if p.name == b {
				j = k
			}
		}
		if i < 0 {
			log.Panicf("pass %s not found", a)
		}
		if j < 0 {
			log.Panicf("pass %s not found", b)
		}
		if i >= j {
			log.Panicf("passes %s and %s out of order", a, b)
		}
	}
}
