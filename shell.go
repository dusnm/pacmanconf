package pacmanconf

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
)

const (
	defaultConfigPath = "/etc/pacman.conf"
)

var (
	ErrPacmanCallFailed = errors.New("pacman-conf")
)

func readConfig(ctx context.Context, configPath string) ([]byte, error) {
	if configPath == "" {
		configPath = defaultConfigPath
	}

	cmd := exec.CommandContext(ctx, "pacman-conf", "--config", configPath)
	out, err := cmd.Output()
	if err != nil {
		exitError, ok := errors.AsType[*exec.ExitError](err)
		if ok {
			return nil,
				fmt.Errorf(
					"%w: %w, %s",
					ErrPacmanCallFailed,
					err,
					exitError.Stderr,
				)
		}

		return nil, fmt.Errorf("%w: %w", ErrPacmanCallFailed, err)
	}

	return bytes.TrimSpace(out), nil
}
