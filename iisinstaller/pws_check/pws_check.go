package pws_check

import (
	"fmt"
	"os/exec"
	"strings"
)

func IsPowerShellInstalled() (bool, error) {
	cmd := exec.Command("wmic", "product", "get", "Name")
	out, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("Error checking IIS installation: %w", err)
	}

	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "IIS") {
			return true, nil
		}
	}

	return false, nil
}
