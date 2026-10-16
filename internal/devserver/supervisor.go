// Copyright 2026 Mona Maret
// SPDX-License-Identifier: Apache-2.0

// Package devserver implements the supervised process pair that powers
// `biblio-cli dev`: it starts vite dev and wrangler dev together,
// forwards their output, and shuts both down cleanly on interrupt
// (FR-005; contracts/cli-commands.md `biblio-cli dev`).
package devserver

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
)

// Config holds the two command specifications and the writer that receives
// their combined output (default os.Stdout when nil).
type Config struct {
	// ViteCmd is the argv slice for the vite dev process (e.g. ["vite", "dev"]).
	ViteCmd []string
	// WranglerCmd is the argv slice for the wrangler dev process.
	WranglerCmd []string
	// Output is where both processes' stdout/stderr are forwarded.
	// Defaults to os.Stdout when nil.
	Output io.Writer
}

// Supervisor supervises the vite dev + wrangler dev process pair.
type Supervisor struct {
	cfg Config
}

// NewSupervisor constructs a Supervisor from the given Config.
func NewSupervisor(cfg Config) *Supervisor {
	if cfg.Output == nil {
		cfg.Output = os.Stdout
	}
	return &Supervisor{cfg: cfg}
}

// Run starts both wrapped processes and blocks until one of three things
// happens:
//
//  1. ctx is cancelled (clean shutdown — both processes are killed).
//  2. Either process exits with a non-zero code — Run returns a descriptive
//     error naming the failing process.
//  3. Both processes exit with code 0 — Run returns nil.
//
// Output from both processes is forwarded to cfg.Output as it arrives.
func (s *Supervisor) Run(ctx context.Context) error {
	viteErr := make(chan error, 1)
	wranglerErr := make(chan error, 1)

	viteCmd, err := s.startProcess(ctx, "vite", s.cfg.ViteCmd)
	if err != nil {
		return err
	}
	wranglerCmd, err := s.startProcess(ctx, "wrangler", s.cfg.WranglerCmd)
	if err != nil {
		killProcessGroup(viteCmd)
		_ = viteCmd.Wait()
		return err
	}

	go func() { viteErr <- viteCmd.Wait() }()
	go func() { wranglerErr <- wranglerCmd.Wait() }()

	var firstErr error
	var once sync.Once

	for i := 0; i < 2; i++ {
		select {
		case err := <-viteErr:
			once.Do(func() {
				if err != nil {
					firstErr = fmt.Errorf("devserver: vite exited unexpectedly: %w", err)
				}
			})
			killProcessGroup(wranglerCmd)
		case err := <-wranglerErr:
			once.Do(func() {
				if err != nil {
					firstErr = fmt.Errorf("devserver: wrangler exited unexpectedly: %w", err)
				}
			})
			killProcessGroup(viteCmd)
		case <-ctx.Done():
			killProcessGroup(viteCmd)
			killProcessGroup(wranglerCmd)
			<-viteErr
			<-wranglerErr
			return nil
		}
	}

	return firstErr
}

// startProcess builds and starts a cmd from argv with ctx attached, wiring
// stdout/stderr to s.cfg.Output.
func (s *Supervisor) startProcess(ctx context.Context, name string, argv []string) (*exec.Cmd, error) {
	if len(argv) == 0 {
		return nil, fmt.Errorf("devserver: %s command is empty", name)
	}
	cmd := exec.CommandContext(ctx, argv[0], argv[1:]...)
	cmd.Stdout = s.cfg.Output
	cmd.Stderr = s.cfg.Output
	setProcessGroup(cmd)
	// exec.CommandContext's default Cancel hook kills only cmd's own
	// process when ctx is done — overriding it to group-kill closes the
	// same race for context cancellation that the explicit kills in Run()
	// close for the "other process exited" paths (issue #192).
	cmd.Cancel = func() error {
		killProcessGroup(cmd)
		return nil
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("devserver: starting %s: %w", name, err)
	}
	return cmd, nil
}
