//go:build linux
// +build linux

package main

import (
	"log"
	"os"
	"strings"

	"golang.org/x/sys/unix"

	"github.com/LDCS/qslinux/blkid"
	"github.com/LDCS/qslinux/df"
)

func GetDrives() ([]Drive, error) {

	var drives []Drive
	var potentialCDRoms []Drive
	blkid := blkid.Blkid(false)
	for _, data := range blkid {
		if strings.Contains(data.Devname_, "/dev/sr") || strings.Contains(data.Devname_, "/dev/cdrom") {
			potentialCDRoms = append(potentialCDRoms, Drive{path: data.Devname_, fsLabel: data.Label_})
		}
	}
	df := df.Df(true, false)
	for _, potentialCDRom := range potentialCDRoms {
		dfData, mounted := df[potentialCDRom.path]
		var mountpoint string
		if !mounted {
			mountpoint = "~/" + potentialCDRom.path[strings.LastIndex(potentialCDRom.path, "/")+1:]
			{
				err := os.Mkdir(mountpoint, 0400)
				if err != nil {
					log.Printf("Cannot make directory in %v; Error: %v", mountpoint, err)
					continue
				}
			}
			{
				err := unix.Mount(potentialCDRom.path, mountpoint, "auto", 0, "")
				if err != nil {
					log.Printf("Cannot mount %v in %v; Error: %v", potentialCDRom.path, mountpoint, err)
					continue
				}
			}
		} else {
			mountpoint = dfData[0].Mountpoint_
		}
		drives = append(drives, Drive{path: mountpoint, fsLabel: potentialCDRom.fsLabel})
	}
	return drives, nil
}
