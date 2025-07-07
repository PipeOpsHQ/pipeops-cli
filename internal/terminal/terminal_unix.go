//go:build !windows

package terminal

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/term"

	"github.com/PipeOpsHQ/pipeops-cli/models"
)

// handleSignals handles terminal resize signals (UNIX only)
func (s *Session) handleSignals(ctx context.Context) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGWINCH)

	for {
		select {
		case <-ctx.Done():
			return
		case <-sigChan:
			// Get new terminal size
			width, height, err := term.GetSize(int(os.Stdin.Fd()))
			if err != nil {
				continue
			}

			// Send resize message
			resizeMsg := models.ResizeMessage{
				Type: "resize",
				Cols: width,
				Rows: height,
			}

			s.conn.WriteJSON(resizeMsg)
		}
	}
}
