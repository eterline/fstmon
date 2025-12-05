package app

import (
	"fmt"
	"os"

	"github.com/eterline/fstmon/pkg/toolkit"
)

func RunAdditional() (error, bool) {
	args := os.Args[1:]
	if len(args) == 0 {
		return nil, false // no additional command
	}

	cmdName := args[0]
	cmdArgs := args[1:]

	// check if command is registered
	err, ok := toolkit.ExecuteCommand(cmdName, cmdArgs...)
	if err != nil {
		return fmt.Errorf("failed to execute command '%s': %w", cmdName, err), false
	}

	return nil, ok
}
