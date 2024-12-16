package main

import (
	"fmt"
	"log"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func main() {
	// Initialize COM library
	err := ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED)
	if err != nil {
		log.Fatalf("Failed to initialize COM library: %v", err)
	}
	defer ole.CoUninitialize()

	// Use `oleutil.CreateObject` to create an instance of WScript.Shell
	wscript, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		log.Fatalf("Failed to create WScript.Shell object: %v", err)
	}
	defer wscript.Release()

	// Query for IDispatch interface
	dispatch, err := wscript.QueryInterface(ole.IID_IDispatch)
	err := oleutil.EnumVariant(dispatch, func(v *ole.Variant) error {
		var name string
		var dispID int32

		// Get the name and dispatch ID of the member
		err := oleutil.VariantToStr(v, &name)
		if err != nil {
			return err
		}
		dispID, err := oleutil.VariantToInt32(v)
		if err != nil {
			return err
		}

		fmt.Printf("Member: %s (DispID: %d)\n", name, dispID)
		return nil
	})

	if err != nil {
		log.Fatalf("Failed to enumerate members: %v", err)
	}
}
