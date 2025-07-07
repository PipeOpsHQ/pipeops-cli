//go:build windows

package terminal

import (
	"context"
)

// handleSignals handles terminal resize signals (Windows stub - no SIGWINCH support)
func (s *Session) handleSignals(ctx context.Context) {
	// Windows doesn't support SIGWINCH signal for terminal resize
	// This is a no-op implementation to maintain compatibility
	<-ctx.Done()
}
