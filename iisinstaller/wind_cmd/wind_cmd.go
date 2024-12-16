package wind_cmd

import (
	"fmt"
	"os/exec"
	"strings"
)

func IsIISInstalled() (bool, error) {
	cmd := exec.Command("powershell", "-Command", "Get-WindowsOptionalFeature -Online | Where-Object FeatureName -eq \"IIS-WebServerRole\" | Select-Object -ExpandProperty State")
	out, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("Error checking IIS installation: %w", err)
	}

	return strings.TrimSpace(string(out)) == "Enabled", nil
}
