// Copyright 2026 Mona Maret
// SPDX-License-Identifier: Apache-2.0

//go:build !windows

package devserver_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/monamaret/biblio/internal/devserver"
)

// TestSupervisor_KillsGrandchildProcesses reproduces issue #192: `npx vite
// dev` / `npx wrangler dev` don't exec into the real dev server, they spawn
// it as a *child* process. Killing only the wrapper (what the pre-fix
// supervisor did) orphans that child, leaving a hanging dev server behind.
// The wrapping shell here mirrors that exact shape — it backgrounds a
// grandchild (which execs into `sleep` so it's trivially identifiable by
// PID) and waits on it, just as npx waits on vite/wrangler. Process groups
// are POSIX-only, so this test (and the kill path it guards) is Unix-only.
func TestSupervisor_KillsGrandchildProcesses(t *testing.T) {
	if _, err := exec.LookPath("sh"); err != nil {
		t.Skip("sh not available; skipping grandchild-cleanup test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	pidFile := filepath.Join(t.TempDir(), "grandchild.pid")
	wrapper := []string{"sh", "-c", `sh -c "echo \$\$ > ` + pidFile + `; exec sleep 60" & wait`}

	var buf strings.Builder
	s := devserver.NewSupervisor(devserver.Config{
		ViteCmd:     wrapper,
		WranglerCmd: []string{"sh", "-c", "sleep 60"},
		Output:      &buf,
	})

	done := make(chan error, 1)
	go func() { done <- s.Run(ctx) }()

	// Wait for the grandchild to report its PID before triggering shutdown.
	var grandchildPID int
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		data, err := os.ReadFile(pidFile)
		if err == nil && strings.TrimSpace(string(data)) != "" {
			grandchildPID, err = strconv.Atoi(strings.TrimSpace(string(data)))
			if err != nil {
				t.Fatalf("parsing grandchild pid file: %v", err)
			}
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if grandchildPID == 0 {
		t.Fatal("grandchild process never reported its PID")
	}

	cancel()

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("Supervisor.Run() did not return within 3s after context cancel")
	}

	// Give the SIGKILL a moment to land, then confirm the grandchild is
	// actually gone — not just that Run() returned.
	time.Sleep(200 * time.Millisecond)
	if err := syscall.Kill(grandchildPID, 0); err == nil {
		t.Errorf("grandchild process %d is still alive after Supervisor.Run() returned — orphaned (issue #192)", grandchildPID)
	}
}
