//go:build windows
// +build windows

package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

func GetDrives() ([]Drive, error) {
	var drives []Drive

	kernel32, _ := syscall.LoadLibrary("kernel32.dll")
	getLogicalDrivesHandle, _ := syscall.GetProcAddress(kernel32, "GetLogicalDrives")
	GetVolumeInformationWHandle, _ := syscall.GetProcAddress(kernel32, "GetVolumeInformationW")

	logicalDrives, _, callErr := syscall.Syscall(uintptr(getLogicalDrivesHandle), 0, 0, 0, 0)
	if callErr != 0 {
		return nil, fmt.Errorf("getLogicalDrivesHandle callErr: %v", callErr)
	}
	for _, letter := range bitsToDrives(uint32(logicalDrives)) {

		var RootPathName = letter + ":\\"
		var VolumeNameBuffer = make([]uint16, syscall.MAX_PATH+1)
		var nVolumeNameSize = uint32(len(VolumeNameBuffer))
		var VolumeSerialNumber uint32
		var MaximumComponentLength uint32
		var FileSystemFlags uint32
		var FileSystemNameBuffer = make([]uint16, 255)
		var nFileSystemNameSize uint32 = syscall.MAX_PATH + 1

		var nargs uintptr = 8
		_, _, callErr := syscall.Syscall9(uintptr(GetVolumeInformationWHandle),
			nargs,
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(RootPathName))),
			uintptr(unsafe.Pointer(&VolumeNameBuffer[0])),
			uintptr(nVolumeNameSize),
			uintptr(unsafe.Pointer(&VolumeSerialNumber)),
			uintptr(unsafe.Pointer(&MaximumComponentLength)),
			uintptr(unsafe.Pointer(&FileSystemFlags)),
			uintptr(unsafe.Pointer(&FileSystemNameBuffer[0])),
			uintptr(nFileSystemNameSize),
			0)

		if callErr != 0 {
			continue
		}
		drives = append(drives, Drive{path: RootPathName, fsLabel: syscall.UTF16ToString(VolumeNameBuffer)})

	}
	return drives, nil
}

func bitsToDrives(bitMap uint32) (drives []string) {
	availableDrives := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	for i := range availableDrives {
		if bitMap&1 == 1 {
			drives = append(drives, availableDrives[i])
		}
		bitMap >>= 1
	}

	return
}
