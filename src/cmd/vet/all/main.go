// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// The vet/all command runs go vet on the standard library and commands.
// It compares the output against a set of whitelists
// maintained in the whitelist directory.
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/build"
	"internal/testenv"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	flagPlatforms = flag.String("p", "", "platform(s) to use e.g. linux/amd64,darwin/386")
	flagAll       = flag.Bool("all", false, "run all platforms")
	flagNoLines   = flag.Bool("n", false, "don't print line numbers")
)

var cmdGoPath string

func main() {
	log.SetPrefix("vet/all: ")
	log.SetFlags(0)

	var err error
	cmdGoPath, err = testenv.GoTool()
	if err != nil {
		log.Print("could not find cmd/go; skipping")
		// We're on a platform that can't run cmd/go.
		// We want this script to be able to run as part of all.bash,
		// so return cleanly rather than with exit code 1.
		return
	}

	flag.Parse()
	switch {
	case *flagAll && *flagPlatforms != "":
		log.Print("-all and -p flags are incompatible")
		flag.Usage()
		os.Exit(2)
	case *flagPlatforms != "":
		vetPlatforms(parseFlagPlatforms())
	case *flagAll:
		vetPlatforms(allPlatforms())
	default:
		host := platform{os: build.Default.GOOS, arch: build.Default.GOARCH}
		host.vet()
	}
}

func allPlatforms() []platform {
	var pp []platform
	cmd := exec.Command(cmdGoPath, "tool", "dist", "list")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	lines := bytes.Split(out, []byte{'\n'})
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		pp = append(pp, parsePlatform(string(line)))
	}
	return pp
}

func parseFlagPlatforms() []platform {
	var pp []platform
	components := strings.Split(*flagPlatforms, ",")
	for _, c := range components {
		pp = append(pp, parsePlatform(c))
	}
	return pp
}

func parsePlatform(s string) platform {
	vv := strings.Split(s, "/")
	if len(vv) != 2 {
		log.Fatalf("could not parse platform %s, must be of form goos/goarch", s)
	}
	return platform{os: vv[0], arch: vv[1]}
}

type whitelist map[string]int

// load adds entries from the whitelist file, if present, for os/arch to w.
func (w whitelist) load(goos string, goarch string) {
	// Look up whether goarch is a 32-bit or 64-bit architecture.
	archbits, ok := nbits[goarch]
	if !ok {
		log.Fatal("unknown bitwidth for arch %q", goarch)
	}

	// Look up whether goarch has a shared arch suffix,
	// such as mips64x for mips64 and mips64le.
	archsuff := goarch
	if x, ok := archAsmX[goarch]; ok {
		archsuff = x
	}

	// Load whitelists.
	filenames := []string{
		"all.txt",
		goos + ".txt",
		goarch + ".txt",
		goos + "_" + goarch + ".txt",
		fmt.Sprintf("%dbit.txt", archbits),
	}
	if goarch != archsuff {
		filenames = append(filenames,
			archsuff+".txt",
			goos+"_"+archsuff+".txt",
		)
	}

	// We allow error message templates using GOOS and GOARCH.
	if goos == "android" {
		goos = "linux" // so many special cases :(
	}

	// Read whitelists and do template substitution.
	replace := strings.NewReplacer("GOOS", goos, "GOARCH", goarch, "ARCHSUFF", archsuff)

	for _, filename := range filenames {
		path := filepath.Join("whitelist", filename)
		f, err := os.Open(path)
		if err != nil {
			// Allow not-exist errors; not all combinations have whitelists.
			if os.IsNotExist(err) {
				continue
			}
			log.Fatal(err)
		}
		scan := bufio.NewScanner(f)
		for scan.Scan() {
			line := scan.Text()
			if len(line) == 0 || strings.HasPrefix(line, "//") {
				continue
			}
			w[replace.Replace(line)]++
		}
		if err := scan.Err(); err != nil {
			log.Fatal(err)
		}
	}
}

type platform struct {
	os   string
	arch string
}

func (p platform) String() string {
	return p.os + "/" + p.arch
}

// ignorePathPrefixes are file path prefixes that should be ignored wholesale.
var ignorePathPrefixes = [...]string{
	// These testdata dirs have lots of intentionally broken/bad code for tests.
	"cmd/go/testdata/",
	"cmd/vet/testdata/",
	"go/printer/testdata/",
	// cmd/compile/internal/big is a vendored copy of math/big.
	// Ignore it so that we only have to deal with math/big issues once.
	"cmd/compile/internal/big/",
}

func vetPlatforms(pp []platform) {
	for _, p := range pp {
		p.vet()
	}
}

func (p platform) vet() {
	if p.arch == "s390x" {
		// TODO: reinstate when s390x gets vet support (issue 15454)
		return
	}
	fmt.Printf("go run main.go -p %s\n", p)

	// Load whitelist(s).
	w := make(whitelist)
	w.load(p.os, p.arch)

	env := append(os.Environ(), "GOOS="+p.os, "GOARCH="+p.arch)

	// Do 'go install std' before running vet.
	// It is cheap when already installed.
	// Not installing leads to non-obvious failures due to inability to typecheck.
	// TODO: If go/loader ever makes it to the standard library, have vet use it,
	// at which point vet can work off source rather than compiled packages.
	cmd := exec.Command(cmdGoPath, "install", "std")
	cmd.Env = env
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("failed to run GOOS=%s GOARCH=%s 'go install std': %v\n%s", p.os, p.arch, err, out)
	}

	// 'go tool vet .' is considerably faster than 'go vet ./...'
	// TODO: The unsafeptr checks are disabled for now,
	// because there are so many false positives,
	// and no clear way to improve vet to eliminate large chunks of them.
	// And having them in the whitelists will just cause annoyance
	// and churn when working on the runtime.
	cmd = exec.Command(cmdGoPath, "tool", "vet", "-unsafeptr=false", ".")
	cmd.Dir = filepath.Join(runtime.GOROOT(), "src")
	cmd.Env = env
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	// Process vet output.
	scan := bufio.NewScanner(stderr)
NextLine:
	for scan.Scan() {
		line := scan.Text()
		if strings.HasPrefix(line, "vet: ") {
			// Typecheck failure: Malformed syntax or multiple packages or the like.
			// This will yield nicer error messages elsewhere, so ignore them here.
			continue
		}

		fields := strings.SplitN(line, ":", 3)
		var file, lineno, msg string
		switch len(fields) {
		case 2:
			// vet message with no line number
			file, msg = fields[0], fields[1]
		case 3:
			file, lineno, msg = fields[0], fields[1], fields[2]
		default:
			log.Fatalf("could not parse vet output line:\n%s", line)
		}
		msg = strings.TrimSpace(msg)

		for _, ignore := range ignorePathPrefixes {
			if strings.HasPrefix(file, filepath.FromSlash(ignore)) {
				continue NextLine
			}
		}

		// Temporarily ignore unrecognized printf verbs from cmd.
		// The compiler now has several fancy verbs (CL 28339)
		// used with types implementing fmt.Formatters,
		// and I believe gri has plans to add many more.
		// TODO: remove when issue 17057 is fixed.
		if strings.HasPrefix(file, "cmd/") && strings.HasPrefix(msg, "unrecognized printf verb") {
			continue
		}

		key := file + ": " + msg
		if w[key] == 0 {
			// Vet error with no match in the whitelist. Print it.
			if *flagNoLines {
				fmt.Printf("%s: %s\n", file, msg)
			} else {
				fmt.Printf("%s:%s: %s\n", file, lineno, msg)
			}
			continue
		}
		w[key]--
	}
	if scan.Err() != nil {
		log.Fatalf("failed to scan vet output: %v", scan.Err())
	}
	err = cmd.Wait()
	// We expect vet to fail.
	// Make sure it has failed appropriately, though (for example, not a PathError).
	if _, ok := err.(*exec.ExitError); !ok {
		log.Fatalf("unexpected go vet execution failure: %v", err)
	}
	printedHeader := false
	if len(w) > 0 {
		for k, v := range w {
			if v != 0 {
				if !printedHeader {
					fmt.Println("unmatched whitelist entries:")
					printedHeader = true
				}
				for i := 0; i < v; i++ {
					fmt.Println(k)
				}
			}
		}
	}
}

// nbits maps from architecture names to the number of bits in a pointer.
// TODO: figure out a clean way to avoid get this info rather than listing it here yet again.
var nbits = map[string]int{
	"386":      32,
	"amd64":    64,
	"amd64p32": 32,
	"arm":      32,
	"arm64":    64,
	"mips64":   64,
	"mips64le": 64,
	"ppc64":    64,
	"ppc64le":  64,
}

// archAsmX maps architectures to the suffix usually used for their assembly files,
// if different than the arch name itself.
var archAsmX = map[string]string{
	"android":  "linux",
	"mips64":   "mips64x",
	"mips64le": "mips64x",
	"ppc64":    "ppc64x",
	"ppc64le":  "ppc64x",
}
