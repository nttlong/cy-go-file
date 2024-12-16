package os_check

import (
	"runtime"
	"syscall"
	"unsafe"
)

// IsWindows checks if the current operating system is Windows.
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

var (
	// Required Windows DLLs
	modadvapi32              = syscall.NewLazyDLL("advapi32.dll")
	procOpenProcessToken     = modadvapi32.NewProc("OpenProcessToken")
	procCheckTokenMembership = modadvapi32.NewProc("CheckTokenMembership")
	procCloseHandle          = modadvapi32.NewProc("CloseHandle")

	// SID for Administrators group
	AdminSID = []byte{0x01, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
)

const (
	TOKEN_QUERY = 0x0008
)

func IsAdmin() (bool, error) {
	// Open the current process token
	hProcess, err := syscall.OpenProcess(syscall.PROCESS_QUERY_INFORMATION, false, uint32(syscall.Getpid()))
	if err != nil {
		return false, err
	}
	defer syscall.CloseHandle(hProcess)

	// Retrieve the token associated with the process
	var token syscall.Handle
	ret, _, err := procOpenProcessToken.Call(uintptr(hProcess), uintptr(TOKEN_QUERY), uintptr(unsafe.Pointer(&token)))
	if ret == 0 {
		return false, err
	}
	defer procCloseHandle.Call(uintptr(token))

	// Check if the token is a member of the "Administrators" group
	var isMember bool
	ret, _, err = procCheckTokenMembership.Call(uintptr(token), uintptr(unsafe.Pointer(&AdminSID[0])), uintptr(unsafe.Pointer(&isMember)))
	if ret == 0 {
		return false, err
	}

	return isMember, nil
}
