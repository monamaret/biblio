//go:build !windows

package devserver

import (
	"os/exec"
	"syscall"
)

// setProcessGroup configures cmd to start as the leader of a new process
// group, so killProcessGroup can later terminate it together with every
// descendant it spawns (e.g. the real vite/wrangler process that `npx`
// spawns as a child rather than exec'ing into — issue #192: killing only
// the npx leader orphaned that child, leaving hanging dev servers behind).
func setProcessGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

// killProcessGroup sends SIGKILL to cmd's entire process group (the
// negative-PID form of kill(2)) rather than just its leader, so any
// descendant the leader spawned is terminated too.
func killProcessGroup(cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil {
		return
	}
	_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
}
