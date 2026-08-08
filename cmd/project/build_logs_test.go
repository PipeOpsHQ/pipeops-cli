package project

import (
	"strings"
	"testing"
)

func TestFormatBuildLogLine(t *testing.T) {
	line := map[string]interface{}{
		"ts":            "2026-08-08T12:00:00Z",
		"current-stage": "build",
		"level":         "info",
		"log":           "npm install complete",
	}
	got := formatBuildLogLine(line)
	for _, want := range []string{"2026-08-08T12:00:00Z", "[build]", "INFO", "npm install complete"} {
		if !strings.Contains(got, want) {
			t.Fatalf("formatBuildLogLine missing %q in %q", want, got)
		}
	}
}

func TestFormatBuildLogLine_MessageField(t *testing.T) {
	got := formatBuildLogLine(map[string]interface{}{
		"message": "deployed",
		"stage":   "deploy",
	})
	if !strings.Contains(got, "deployed") || !strings.Contains(got, "[deploy]") {
		t.Fatalf("unexpected format: %q", got)
	}
}

func TestFirstString(t *testing.T) {
	m := map[string]interface{}{"log": "from-log", "message": "from-message"}
	if got := firstString(m, "message", "log"); got != "from-message" {
		t.Fatalf("prefer message, got %q", got)
	}
	if got := firstString(m, "missing", "log"); got != "from-log" {
		t.Fatalf("fallback log, got %q", got)
	}
}
