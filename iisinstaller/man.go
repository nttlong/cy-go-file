package main

import (
	"fmt"

	iisutils "codx.iis.installer/iisutils"
	// os_check "codx.iis.installer/os_check"
	// pws_check "codx.iis.installer/pws_check"
	// wind_cmd "codx.iis.installer/wind_cmd"
)

func main() {
	// runspace := powershell.CreateRunspaceSimple()
	// defer runspace.Close()
	// results1 := runspace.ExecScript("Get-ChildItem 'IIS:\\AppPools' | Select-Object -ExpandProperty Name", false, nil, "OS")
	//not defering close as we do not need the results
	//results1.Close()
	appPools, err := iisutils.ListAppPools()
	if err != nil {
		fmt.Println(err)
	}
	for _, appPool := range appPools {
		fmt.Println(appPool)
	}

	// // get cuurent user who is running the program
	// if !os_check.IsWindows() {
	// 	fmt.Println("This program is only for Windows")
	// 	return
	// }

	// user, err := user.Current()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("You arer runnig, " + user.Username + "!")
	// fmt.Println("Check account privileges ...")
	// ok, err := pws_check.IsPowerShellInstalled()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// if ok {
	// 	fmt.Println("You are an admin user")
	// }
	// fmt.Println("Checking if IIS is installed...")
	// ok, err = wind_cmd.IsIISInstalled()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// if ok {
	// 	fmt.Println("IIS is installed")
	// } else {
	// 	fmt.Println("IIS is not installed")
	// }

	// err = iisutils.CreateAppPool("long-test-2")
	// if err != nil {
	// 	fmt.Println(err)
	// }
}
