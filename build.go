//go:build mage

// This file is used by mage to build the project.
package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	mgtarget "github.com/magefile/mage/target"
)

const (
	GO           = "GO"
	CC           = "CC"
	CXX          = "CXX"
	CLANG_FORMAT = "CLANG_FORMAT"
	EH_FRAME_BIN = "EH_FRAME_BIN"
)

var (
	Default = All

	// toolBinaries is a map of tool ENV vars to their default binary.
	// If the ENV var is not set, the default binary will be used.
	// If the default is NOT a path, it will be looked up in the PATH.
	toolBinaries = map[string]string{
		GO:           mg.GoCmd(),
		CC:           "zig cc",
		CXX:          "zig c++",
		CLANG_FORMAT: "clang-format",
	}
	// goToolBinaries is a map of go tool ENV vars to their default binary.
	// It will be run with `go run`.
	// e.g go run github.com/parca-dev/parca-agent/tree/main/cmd/eh-frame@latest
	goToolBinaries = map[string]string{
		EH_FRAME_BIN: "github.com/parca-dev/parca-agent/cmd/eh-frame@latest",
	}

	host = targetFromGoArch(runtime.GOARCH)
)

type Target struct {
	goarch  string
	goos    string
	cTarget string
}

var targets = []Target{
	{
		goarch:  "amd64",
		goos:    "linux",
		cTarget: "x86_64-linux-gnu",
	},
	{
		goarch:  "arm64",
		goos:    "linux",
		cTarget: "aarch64-linux-gnu",
	},
}

func targetFromGoArch(goArch string) Target {
	for _, t := range targets {
		if t.goarch == goArch {
			return t
		}
	}
	panic(fmt.Errorf("unknown go arch %s", goArch))
}

func (t Target) String() string {
	return fmt.Sprintf("%s (%s) on (%s)", t.goarch, t.cTarget, host)
}

func (t Target) CrossCCCmd() string {
	return fmt.Sprintf("%s -target %s", bin(CC), t.cTarget)
}

func (t Target) CrossCXXCmd() string {
	return fmt.Sprintf("%s -target %s", bin(CXX), t.cTarget)
}

const (
	out      = "out"
	vendored = "vendored"

	tables        = "tables"
	normalTables  = "normal"
	finalTables   = "final"
	compactTables = "compact"
)

// fullOutPath returns the full path to the output binary.
func fullOutPath(arch, name string) string {
	return path.Join(out, arch, name)
}

// fullTablesPath returns the full path to the output tables.
func fullTablesPath(arch, name string) string {
	return path.Join(tables, arch, name)
}

var Directories = []string{
	out,
	tables,
}

// All runs all the targets.
func All() error {
	var err error
	var format Format
	err = errors.Join(err, format.All())
	var build Build
	err = errors.Join(err, build.All())
	var generate Generate
	err = errors.Join(err, generate.All())
	return err
}

const (
	goGlob  = "./src/*.go"
	cppGlob = "./src/*.cpp"
)

// Format runs all the formatters.
type Format mg.Namespace

// Go runs gofmt on all the Go code.
func (Format) Go() error {
	matches, err := filepath.Glob(goGlob)
	if err != nil {
		return err
	}
	return sh.Run(bin(GO), "fmt", strings.Join(matches, " "))
}

// Cpp runs clang-format on all the C code.
func (Format) Cpp() error {
	matches, err := filepath.Glob(cppGlob)
	if err != nil {
		return err
	}
	args := []string{"-i"}
	args = append(args, matches...)
	return sh.Run(bin(CLANG_FORMAT), args...)
}

// All runs all the linters.
// clang-format for C code.
// gofmt for Go code.
func (Format) All() error {
	mg.Deps(Format.Go, Format.Cpp)
	return nil
}

// Build runs all the builders.
type Build mg.Namespace

// Build builds all the binaries.
func (Build) All() error {
	mg.Deps(Build.Go, Build.Cpp)
	return nil
}

// Variant is a build variant.
type Variant struct {
	name  string
	src   string
	flags []string
}

var goBuildVariants = []Variant{
	{
		name:  "basic-go",
		src:   "./src/main.go",
		flags: nil,
	},
}

// Go builds all the Go binaries.
func (Build) Go() error {
	mg.Deps(dirs)

	var errs error
	for _, vr := range goBuildVariants {
		if err := buildGo(host, vr.src, vr.name, vr.flags...); err != nil {
			errs = errors.Join(errs, err)
		}
	}

	return errs
}

var cppBuildVariants = []Variant{
	{
		name:  "basic-cpp",
		src:   "./src/basic-cpp.cpp",
		flags: nil,
	},
	{
		name:  "basic-cpp-with-debuginfo",
		src:   "./src/basic-cpp.cpp",
		flags: []string{"-g"},
	},
	{
		name:  "basic-cpp-no-fp",
		src:   "./src/basic-cpp.cpp",
		flags: []string{"-fno-omit-frame-pointer"},
	},
	{
		name:  "basic-cpp-no-fp-with-debuginfo",
		src:   "./src/basic-cpp.cpp",
		flags: []string{"-fno-omit-frame-pointer", "-g"},
	},
	{
		name:  "basic-cpp-no-pie-with-compressed-debuginfo",
		src:   "./src/basic-cpp.cpp",
		flags: []string{"-fno-pie", "-fno-PIE", "-g", "-gz=zlib"},
	},
	{
		name:  "basic-cpp-no-pie-and-no-fp",
		src:   "./src/basic-cpp.cpp",
		flags: []string{"-fno-pie", "-fno-PIE", "-fno-omit-frame-pointer"},
	},
	{
		name:  "basic-cpp-no-pie-and-no-fp-with-compressed-debuginfo",
		src:   "./src/basic-cpp.cpp",
		flags: []string{"-fno-pie", "-fno-PIE", "-fno-omit-frame-pointer", "-g", "-gz=zlib"},
	},
	{
		name:  "basic-cpp-plt",
		src:   "./src/basic-cpp-plt.cpp",
		flags: nil,
	},
	{
		name:  "basic-cpp-plt-with-debuginfo",
		src:   "./src/basic-cpp-plt.cpp",
		flags: []string{"-g"},
	},
	{
		name:  "basic-cpp-plt-pie",
		src:   "./src/basic-cpp-plt.cpp",
		flags: []string{"-fPIE", "-pie"},
	},
	{
		name:  "basic-cpp-plt-hardened",
		src:   "./src/basic-cpp-plt.cpp",
		flags: []string{"-fPIE", "-pie", "-fstack-protector-all", "-D_FORTIFY_SOURCE=2", "-Wl,-z,now", "-Wl,-z,relro", "-O2"},
	},
	{
		name:  "basic-cpp-plt-no-pie",
		src:   "./src/basic-cpp-plt.cpp",
		flags: []string{"-fno-pie", "-fno-PIE"},
	},
	{
		name:  "basic-cpp-multithreaded-no-fp",
		src:   "./src/basic-cpp-multithreaded.cpp",
		flags: []string{"-fno-omit-frame-pointer"},
	},
	// -c -g -gz=zlib
	{
		name:  "basic-cpp-jit",
		src:   "./src/basic-cpp-jit.cpp",
		flags: []string{"-g"},
	},
	{
		name:  "basic-cpp-jit-no-fp",
		src:   "./src/basic-cpp-jit.cpp",
		flags: []string{"-fno-omit-frame-pointer"},
	},
}

// Cpp builds all the C binaries.
func (Build) Cpp() error {
	mg.Deps(dirs)

	var errs error
	for _, vr := range cppBuildVariants {
		if err := buildCpp(host, vr.src, vr.name, vr.flags...); err != nil {
			errs = errors.Join(errs, err)
		}
	}

	return errs
}

// Cross builds all the binaries for all the targets.
func (Build) Cross() error {
	mg.Deps(dirs)

	var errs error
	for _, target := range targets {
		log.Printf("Building for %s\n", target)
		for _, vr := range goBuildVariants {
			if err := buildGo(target, vr.src, vr.name, vr.flags...); err != nil {
				errs = errors.Join(errs, err)
			}
		}
		for _, vr := range cppBuildVariants {
			if err := buildCpp(target, vr.src, vr.name, vr.flags...); err != nil {
				errs = errors.Join(errs, err)
			}
		}
	}
	return errs
}

// buildCpp builds a C binary for the given target if it has changed.
func buildCpp(target Target, src, out string, flags ...string) error {
	out = fullOutPath(target.goarch, out)
	changed, err := mgtarget.Path(out, src)
	if err != nil {
		return err
	}
	if !changed {
		return nil
	}
	return run(target.CrossCXXCmd(), append([]string{src, "-o", out}, flags...)...)
}

// buildGo builds a Go binary for the given target if it has changed.
func buildGo(target Target, src, out string, flags ...string) error {
	out = fullOutPath(target.goarch, out)
	changed, err := mgtarget.Path(out, src)
	if err != nil {
		return err
	}
	if !changed {
		return nil
	}
	return sh.RunWith(map[string]string{
		"GOARCH": target.goarch,
		"GOOS":   target.goos,
	}, bin(GO), append([]string{"build", "-o", out, src}, flags...)...)
}

// Generate runs all the generators.
type Generate mg.Namespace

// All generates all the tables.
func (Generate) All() error {
	mg.SerialDeps(Generate.Normal, Generate.Final, Generate.Compact)
	return nil
}

// generateTables generates all the tables for the given flags.
func generateTables(src, dst string, flags ...string) error {
	for _, target := range targets {
		err := filepath.Walk(path.Join(src, target.goarch), func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			outputFilePath := fullTablesPath(target.goarch, path.Join(dst, info.Name()+".txt"))

			chaged, err := mgtarget.Path(outputFilePath, p)
			if err != nil {
				return err
			}
			if !chaged {
				return nil
			}

			output, err := runOutGoTool(EH_FRAME_BIN, append([]string{"--executable", p}, flags...)...)
			if err != nil {
				return err
			}

			if err := os.MkdirAll(path.Dir(outputFilePath), 0755); err != nil {
				return err
			}

			return os.WriteFile(outputFilePath, []byte(output), 0644)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Normal generates all the normal tables.
func (Generate) Normal() error {
	var errs error
	errs = errors.Join(errs, generateTables(out, normalTables))
	errs = errors.Join(errs, generateTables(vendored, normalTables))
	return errs
}

// Final generates all the final tables.
func (Generate) Final() error {
	var errs error
	errs = errors.Join(errs, generateTables(out, finalTables, "--final"))
	errs = errors.Join(errs, generateTables(vendored, finalTables, "--final"))
	return errs
}

// Compact generates all the compact tables.
func (Generate) Compact() error {
	var errs error
	errs = errors.Join(errs, generateTables(out, compactTables, "--compact"))
	errs = errors.Join(errs, generateTables(vendored, compactTables, "--compact"))
	return errs
}

// Clean cleans all the generated files and binaries.
func Clean() error {
	for _, dir := range Directories {
		if err := os.RemoveAll(dir); err != nil {
			return err
		}
	}
	return nil
}

// Env prints the environment variables.
func Env() error {
	for _, vr := range os.Environ() {
		fmt.Println(vr)
	}
	return nil
}

// Helpers:

// dirs creates all the directories.
func dirs() error {
	for _, target := range targets {
		for _, dir := range Directories {
			if err := os.MkdirAll(path.Join(dir, target.goarch), 0755); err != nil {
				return err
			}
		}
	}
	return nil
}

// getOrDefault returns the ENV var or the default binary.
func getOrDefault(key string) string {
	if cmd := os.Getenv(key); cmd != "" {
		return cmd
	}
	if val, ok := toolBinaries[key]; ok {
		return val
	}
	panic(fmt.Errorf("no default binary for %s", key))
}

// bin returns the full path to the binary.
func bin(cmd string) string {
	exe := getOrDefault(cmd)
	if strings.HasPrefix(exe, "./") || strings.HasPrefix(exe, "../") || strings.HasPrefix(exe, "/") {
		if _, err := os.Stat(exe); os.IsNotExist(err) {
			panic(fmt.Sprintf("binary %s does not exist: %s", cmd, exe))
		}
	}

	parts := strings.Split(exe, " ")
	if len(parts) > 1 {
		exe = parts[0]
	}
	if _, err := exec.LookPath(exe); err != nil {
		panic(fmt.Sprintf("binary %s does not exist in PATH: %s", cmd, exe))
	}

	return strings.Join(parts, " ")
}

// run runs the command with the given args.
func run(cmd string, args ...string) error {
	parts := strings.Split(cmd, " ")
	if len(parts) > 1 {
		return sh.Run(parts[0], append(parts[1:], args...)...)
	}
	return sh.Run(cmd, args...)
}

var (
	goRun    = sh.RunCmd(bin(GO), "run")
	goRunOut = sh.OutCmd(bin(GO), "run")
)

// runGoTool runs the go tool with the given args.
func runGoTool(cmd string, args ...string) error {
	return goRun(append([]string{goToolBinaries[cmd]}, args...)...)
}

// runOutGoTool runs the go tool with the given args and returns the output.
func runOutGoTool(cmd string, args ...string) (string, error) {
	return goRunOut(append([]string{goToolBinaries[cmd]}, args...)...)
}
