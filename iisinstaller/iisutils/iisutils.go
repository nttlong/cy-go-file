package iisutils

import (
	"fmt"
	"os/exec"
	"strings"
)

func CreateAppPool(appPoolName string) error {
	// PowerShell command to create a new app pool and configure properties
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf(`
		Import-Module WebAdministration
		# Create the application pool
		New-WebAppPool -Name '%s'
		
		# Set the ManagedPipelineMode and CLRVersion
		Set-ItemProperty "IIS:\\AppPools\\%s" -Name "ManagedPipelineMode" -Value "Integrated"
		Set-ItemProperty "IIS:\\AppPools\\%s" -Name "CLRVersion" -Value "No Managed Code"
	`, appPoolName, appPoolName, appPoolName))

	// Run the command
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Error creating application pool: %w", err)
	}
	return nil
}
func GetIISVersion() (string, error) {
	cmd := exec.Command("powershell", "-Command", "Get-ItemProperty HKLM:\\SOFTWARE\\Microsoft\\IIS\\Parameters | Select-Object -ExpandProperty MajorVersion")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("Error getting IIS version: %w", err)
	}

	return strings.TrimSpace(string(out)), nil
}
func ListAppPools() ([]string, error) {
	fmt.Println("Get-ChildItem 'IIS:\\\\AppPools' | Select-Object -ExpandProperty Name")
	cmd := exec.Command("powershell", "-Command", "Import-Module WebAdministration; Get-ChildItem 'IIS:\\AppPools' | Select-Object -ExpandProperty Name")

	out, err := cmd.CombinedOutput()
	if err != nil {

		return nil, fmt.Errorf("Error listing application pools: %w", string(out))
	}

	var appPools []string
	for _, line := range strings.Split(string(out), "\n") {
		appPool := strings.TrimSpace(line)
		if appPool != "" {
			appPools = append(appPools, appPool)
		}
	}

	return appPools, nil
}
