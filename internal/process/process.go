package process

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

// Process represents a listening process on a port.
type Process struct {
	PID     int
	Name    string
	User    string
	Type    string
	Address string
	Port    string
}

// GetListeningPorts runs lsof and parses listening TCP processes.
func GetListeningPorts() ([]Process, error) {
	out, err := exec.Command("lsof", "-i", "-n", "-P", "-sTCP:LISTEN").Output()
	if err != nil {
		// lsof exits with code 1 when no results; treat as empty
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return nil, nil
		}
		return nil, err
	}

	seen := make(map[string]bool)
	var processes []Process

	lines := strings.Split(string(out), "\n")
	for _, line := range lines[1:] { // skip header
		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}
		name := fields[0]
		pidStr := fields[1]
		user := fields[2]
		netType := fields[4] // IPv4, IPv6, etc.
		addrPort := fields[8]

		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			continue
		}

		// Parse address:port from the NAME column (e.g. "*:3000", "127.0.0.1:8080")
		lastColon := strings.LastIndex(addrPort, ":")
		if lastColon < 0 {
			continue
		}
		addr := addrPort[:lastColon]
		port := addrPort[lastColon+1:]

		key := fmt.Sprintf("%d:%s", pid, port)
		if seen[key] {
			continue
		}
		seen[key] = true

		processes = append(processes, Process{
			PID:     pid,
			Name:    name,
			User:    user,
			Type:    netType,
			Address: addr,
			Port:    port,
		})
	}
	return processes, nil
}

// Kill sends SIGKILL to the process with the given PID.
func Kill(pid int) error {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("could not find process %d: %v", pid, err)
	}
	if err := proc.Signal(syscall.SIGKILL); err != nil {
		return fmt.Errorf("kill failed: %v", err)
	}
	return nil
}
