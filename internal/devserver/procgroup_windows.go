// Copyright 2026 Mona Maret
// SPDX-License-Identifier: Apache-2.0

//go:build windows

package devserver

import "os/exec"

// setProcessGroup is a no-op on Windows: process groups are POSIX-only.
// killProcessGroup below falls back to the pre-#192 single-process kill,
// so Windows behavior is unchanged rather than regressed; CI only targets
// ubuntu-latest today, so the orphaned-grandchild fix only needs to land
// for Unix.
func setProcessGroup(cmd *exec.Cmd) {}

// killProcessGroup kills only cmd's own process, matching this package's
// pre-#192 behavior on platforms without POSIX process groups.
func killProcessGroup(cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil {
		return
	}
	_ = cmd.Process.Kill()
}
