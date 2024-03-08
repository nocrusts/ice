// Package wine implements wine program command routines for
// interacting with a wineprefix [Prefix]
package wine

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	ErrWineRootAbs  = errors.New("wineroot is not absolute")
	ErrWineNotFound = errors.New("wine64 not found in system or wineroot")
)

// Prefix is a representation of a wineprefix, which is where
// WINE stores its data and is equivalent to a C:\ drive.
type Prefix struct {
	// Path to a wine installation.
	Root string

	// Stdout and Stderr specify the descendant Prefix wine call's
	// standard output and error. This is mostly reserved for logging purposes.
	// By default, they will be set to their os counterparts.
	Stderr io.Writer
	Stdout io.Writer

	wine string
	dir  string
}

func (p Prefix) String() string {
	return p.dir
}

// Wine64 returns a path to the system or wineroot's 'wine64'.
// Wine64 will attempt to resolve for a [ULWGL launcher] if
// it is present and set necessary environment variables.
//
// [ULWGL launcher]: https://github.com/Open-Wine-Components/ULWGL-launcher
func Wine64(root string) (string, error) {
	wineLook := "wine64"

	if root != "" {
		if !filepath.IsAbs(root) {
			return "", ErrWineRootAbs
		}

		if strings.Contains(strings.ToLower(root), "ulwgl") {
			slog.Info("Detected ULWGL Wineroot!")

			wineLook = filepath.Join(root, "ulwgl-run")
			os.Setenv("STORE", "none")
		} else {
			wineLook = filepath.Join(root, "bin", wineLook)
		}
	}

	wine, err := exec.LookPath(wineLook)
	if err != nil {
		return "", err
	}

	return wine, nil
}

// New returns a new Prefix.
//
// [Wine64] will be used to verify the named root or if
// wine is installed; it will be looked in $PATH only once -
// if the wine executable changes it will not be re-looked.
//
// dir must be an absolute path and has correct permissions
// to modify.
func New(dir string, root string) (*Prefix, error) {
	w, err := Wine64(root)
	if err != nil {
		return nil, fmt.Errorf("bad wine: %w", err)
	}

	// Always ensure its created, wine will complain if the root
	// directory doesnt exist
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create prefix: %s", err)
	}

	return &Prefix{
		Root:   root,
		Stderr: os.Stderr,
		Stdout: os.Stdout,
		wine:   w,
		dir:    dir,
	}, nil
}

// Dir returns the directory of the Prefix.
func (p *Prefix) Dir() string {
	return p.dir
}

// Wine returns a new Cmd with the prefix's Wine as the named program.
func (p *Prefix) Wine(exe string, arg ...string) *Cmd {
	arg = append([]string{exe}, arg...)
	cmd := p.Command(p.wine, arg...)

	if strings.Contains(strings.ToLower(p.wine), "ulwgl") {
		cmd.Env = append(cmd.Environ(), "PROTON_VERB=runinprefix")
	}

	return cmd
}

// Kill kills the Prefix's processes.
func (p *Prefix) Kill() error {
	return p.Wine("wineboot", "-k").Run()
}

// Init performs initialization for first Wine instance.
func (p *Prefix) Init() error {
	return p.Wine("wineboot", "-i").Headless().Run()
}

// Update updates the wineprefix directory.
func (p *Prefix) Update() error {
	return p.Wine("wineboot", "-u").Headless().Run()
}

// Version returns the wineprefix's Wine version.
func (p *Prefix) Version() string {
	cmd := p.Wine("--version")
	cmd.Stdout = nil // required for Output()
	cmd.Stderr = nil

	ver, _ := cmd.Output()
	if len(ver) < 0 {
		return "unknown"
	}

	// remove newline
	return string(ver[:len(ver)-1])
}
