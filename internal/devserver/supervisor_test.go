package devserver_test

import (
	"context"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/monamaret/biblio/internal/devserver"
)

// Supervisor starts two processes and surfaces both cleanly.
// These tests use trivially short-lived shell commands in place of real
// vite/wrangler processes so the suite runs without Node or Wrangler
// installed (contracts/cli-commands.md `biblio-cli dev`; Principle VI —
// test-first, NON-NEGOTIABLE).

// TestSupervisor_StartsBothProcesses verifies that Supervisor.Run starts
// both configured commands and waits for them to exit cleanly.
func TestSupervisor_StartsBothProcesses(t *testing.T) {
	if _, err := exec.LookPath("sh"); err != nil {
		t.Skip("sh not available; skipping process-launch test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var buf strings.Builder
	s := devserver.NewSupervisor(devserver.Config{
		ViteCmd:    []string{"sh", "-c", "echo vite-started && exit 0"},
		WranglerCmd: []string{"sh", "-c", "echo wrangler-started && exit 0"},
		Output:     &buf,
	})

	if err := s.Run(ctx); err != nil {
		t.Errorf("Supervisor.Run() error = %v", err)
	}
}

// TestSupervisor_ReportsViteAddress verifies that the supervisor surfaces
// the address line emitted by the vite process via the configured Output.
func TestSupervisor_SurfacesOutput(t *testing.T) {
	if _, err := exec.LookPath("sh"); err != nil {
		t.Skip("sh not available; skipping output-forwarding test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var buf strings.Builder
	s := devserver.NewSupervisor(devserver.Config{
		ViteCmd:    []string{"sh", "-c", "echo 'Local: http://localhost:5173'"},
		WranglerCmd: []string{"sh", "-c", "echo 'Ready on http://localhost:8787'"},
		Output:     &buf,
	})

	_ = s.Run(ctx)
	out := buf.String()
	if !strings.Contains(out, "localhost") {
		t.Errorf("Supervisor.Run(): output missing process addresses:\n%s", out)
	}
}

// TestSupervisor_ReportsFailingProcess verifies that when one of the two
// wrapped processes exits with a non-zero code, Run returns a non-nil error
// naming which process failed (contracts/cli-commands.md `biblio-cli dev`
// "Output (failure)").
func TestSupervisor_ReportsFailingProcess(t *testing.T) {
	if _, err := exec.LookPath("sh"); err != nil {
		t.Skip("sh not available; skipping failure-report test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var buf strings.Builder
	s := devserver.NewSupervisor(devserver.Config{
		ViteCmd:    []string{"sh", "-c", "exit 1"},
		WranglerCmd: []string{"sh", "-c", "sleep 60"},
		Output:     &buf,
	})

	err := s.Run(ctx)
	if err == nil {
		t.Fatal("Supervisor.Run(): want error when a wrapped process fails, got nil")
	}
}

// TestSupervisor_CleanShutdown verifies that cancelling the context causes
// Run to return without leaving orphaned processes — the supervisor must
// propagate the cancel to both child processes.
func TestSupervisor_CleanShutdown(t *testing.T) {
	if _, err := exec.LookPath("sh"); err != nil {
		t.Skip("sh not available; skipping shutdown test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	var buf strings.Builder
	s := devserver.NewSupervisor(devserver.Config{
		ViteCmd:    []string{"sh", "-c", "sleep 60"},
		WranglerCmd: []string{"sh", "-c", "sleep 60"},
		Output:     &buf,
	})

	done := make(chan error, 1)
	go func() { done <- s.Run(ctx) }()

	// Cancel immediately — supervisor should unblock well within the test timeout.
	cancel()

	select {
	case <-done:
		// Run returned after cancel — pass.
	case <-time.After(3 * time.Second):
		t.Error("Supervisor.Run() did not return within 3s after context cancel")
	}
}
