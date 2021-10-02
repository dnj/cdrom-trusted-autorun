package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/alexflint/go-arg"
)

func watch(parser *arg.Parser, args Args) {
	if args.Watch.PublicKey == "" {
		self, err := os.Executable()
		if err == nil {
			args.Watch.PublicKey = filepath.Join(filepath.Dir(self), "public-key.sig")
		}
	}

	for {
		var execName string
		if os.PathSeparator == '/' {
			execName = "autorun"
		} else {
			execName = "autorun.exe"
		}

		log.Println("Finding drives and mountpoints")
		drives, err := GetDrives()
		if err != nil {
			log.Fatalf("GetDrives error: %v", err)
		}
		log.Println(drives)
		for _, drive := range drives {
			path := filepath.Join(drive.path, execName)
			log.Println("Looking for " + path)
			if _, err := os.Stat(path); err == nil {
				log.Println("Found it")
				checkAndRunExecutable(path, args.Watch.PublicKey)
			} else {
				log.Println("Notfound")
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func checkAndRunExecutable(target string, publicKey string) error {
	if err := Verify(target, "", publicKey); err != nil {
		return err
	}
	command := exec.Command(target)
	return command.Run()
}
