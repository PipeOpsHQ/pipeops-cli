package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
)

// RunShellCommandWithEnvStreaming runs a shell command while streaming stdout/stderr
// to the console. It picks an appropriate shell per-OS.
//
// On Windows, the PipeOps agent installer scripts require a POSIX shell. We try:
//  1. bash (Git Bash)
//  2. sh (MSYS)
//  3. wsl (WSL2) + bash
//
// extraEnv should be a slice of KEY=VALUE pairs that will be added to the process
// environment. Only keys in extraEnv are forwarded to WSL via WSLENV.
func RunShellCommandWithEnvStreaming(command string, extraEnv []string) (string, error) {
	env := append(os.Environ(), extraEnv...)

	if runtime.GOOS == "windows" {
		if bashPath, err := exec.LookPath("bash"); err == nil {
			return RunCommandWithEnvStreaming(bashPath, []string{"-lc", command}, env)
		}
		if shPath, err := exec.LookPath("sh"); err == nil {
			return RunCommandWithEnvStreaming(shPath, []string{"-c", command}, env)
		}
		if wslPath, err := exec.LookPath("wsl"); err == nil {
			keys := envKeys(extraEnv)
			wslEnvValue := mergeWSLENV(os.Getenv("WSLENV"), keys)
			env = replaceEnv(env, "WSLENV", wslEnvValue)
			return RunCommandWithEnvStreaming(wslPath, []string{"bash", "-lc", command}, env)
		}

		return "", fmt.Errorf("PipeOps agent installation requires a POSIX shell on Windows. Install Git for Windows (Git Bash) or enable WSL2, then re-run")
	}

	if shPath, err := exec.LookPath("sh"); err == nil {
		return RunCommandWithEnvStreaming(shPath, []string{"-c", command}, env)
	}
	if bashPath, err := exec.LookPath("bash"); err == nil {
		return RunCommandWithEnvStreaming(bashPath, []string{"-lc", command}, env)
	}

	return "", fmt.Errorf("missing required shell: could not find `sh` or `bash` in PATH")
}

// RunShellCommandWithEnv runs a shell command and captures stdout/stderr.
// It uses the same shell selection logic as RunShellCommandWithEnvStreaming.
func RunShellCommandWithEnv(command string, extraEnv []string) (string, error) {
	env := append(os.Environ(), extraEnv...)

	if runtime.GOOS == "windows" {
		if bashPath, err := exec.LookPath("bash"); err == nil {
			return RunCommandWithEnv(bashPath, []string{"-lc", command}, env)
		}
		if shPath, err := exec.LookPath("sh"); err == nil {
			return RunCommandWithEnv(shPath, []string{"-c", command}, env)
		}
		if wslPath, err := exec.LookPath("wsl"); err == nil {
			keys := envKeys(extraEnv)
			wslEnvValue := mergeWSLENV(os.Getenv("WSLENV"), keys)
			env = replaceEnv(env, "WSLENV", wslEnvValue)
			return RunCommandWithEnv(wslPath, []string{"bash", "-lc", command}, env)
		}

		return "", fmt.Errorf("PipeOps agent installation requires a POSIX shell on Windows. Install Git for Windows (Git Bash) or enable WSL2, then re-run")
	}

	if shPath, err := exec.LookPath("sh"); err == nil {
		return RunCommandWithEnv(shPath, []string{"-c", command}, env)
	}
	if bashPath, err := exec.LookPath("bash"); err == nil {
		return RunCommandWithEnv(bashPath, []string{"-lc", command}, env)
	}

	return "", fmt.Errorf("missing required shell: could not find `sh` or `bash` in PATH")
}

func replaceEnv(env []string, key string, value string) []string {
	prefix := key + "="
	out := make([]string, 0, len(env)+1)
	for _, kv := range env {
		if strings.HasPrefix(kv, prefix) {
			continue
		}
		out = append(out, kv)
	}
	out = append(out, prefix+value)
	return out
}

func envKeys(env []string) []string {
	seen := map[string]struct{}{}
	keys := make([]string, 0, len(env))
	for _, kv := range env {
		key, _, ok := strings.Cut(kv, "=")
		if !ok || key == "" {
			continue
		}
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func mergeWSLENV(existing string, keys []string) string {
	if len(keys) == 0 {
		return existing
	}

	seen := map[string]struct{}{}
	entries := make([]string, 0)

	for _, entry := range strings.Split(existing, ":") {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		base := strings.SplitN(entry, "/", 2)[0]
		if base != "" {
			seen[base] = struct{}{}
		}
		entries = append(entries, entry)
	}

	for _, key := range keys {
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		entries = append(entries, key)
	}

	return strings.Join(entries, ":")
}
