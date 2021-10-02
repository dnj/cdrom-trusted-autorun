package main

import (
	"crypto/ed25519"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/alexflint/go-arg"
)

func sign(parser *arg.Parser, args Args) {
	if args.Sign.PrivateKey == "" {
		self, err := os.Executable()
		if err == nil {
			args.Sign.PrivateKey = filepath.Join(filepath.Dir(self), "private-key.sig")
		}
	}
	if args.Sign.Signature == "" {
		args.Sign.Signature = args.Sign.Target + ".sig"
	}

	if _, err := os.Stat(args.Sign.Target); err != nil {
		parser.Fail("Cannot open " + args.Sign.Target + ". Provide a valid path to target file")
		return
	}
	if _, err := os.Stat(args.Sign.PrivateKey); err != nil {
		parser.Fail("Cannot open " + args.Sign.PrivateKey + ". Provide a valid path to public key file")
		return
	}
	error := Sign(args.Sign.Target, args.Sign.Signature, args.Sign.PrivateKey)
	if error != nil {
		log.Fatal(error)
	}
}

func Sign(targetPath string, signaturePath string, privateKeyPath string) error {
	target, err := ioutil.ReadFile(targetPath)
	if err != nil {
		log.Println("Error in reading target file: " + err.Error())
		return err
	}

	privateKey, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		log.Println("Error in reading public key file: " + err.Error())
		return err
	}

	signature := ed25519.Sign(privateKey, target)

	err = os.WriteFile(signaturePath, signature, 0644)
	if err != nil {
		log.Println("Error in writeing signature file: " + err.Error())
		return err
	}

	return nil
}
