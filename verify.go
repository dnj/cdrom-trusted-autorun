package main

import (
	"crypto/ed25519"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/alexflint/go-arg"
)

func verify(parser *arg.Parser, args Args) {
	if args.Verify.PublicKey == "" {
		self, err := os.Executable()
		if err == nil {
			args.Verify.PublicKey = filepath.Join(filepath.Dir(self), "public-key.sig")
		}
	}
	if args.Verify.Signature == "" {
		args.Verify.Signature = args.Verify.Target + ".sig"
	}

	if _, err := os.Stat(args.Verify.Target); err != nil {
		parser.Fail("Cannot open " + args.Verify.Target + ". Provide a valid path to target file")
		return
	}
	if _, err := os.Stat(args.Verify.Signature); err != nil {
		parser.Fail("Cannot open " + args.Verify.Signature + ". Provide a valid path to signature file")
		return
	}
	if _, err := os.Stat(args.Verify.PublicKey); err != nil {
		parser.Fail("Cannot open " + args.Verify.PublicKey + ". Provide a valid path to public key file")
		return
	}
	error := Verify(args.Verify.Target, args.Verify.Signature, args.Verify.PublicKey)
	if error != nil {
		log.Fatal(error)
	}
	fmt.Println("Target and it's signature is valid")
}

func Verify(targetPath string, signaturePath string, publicKeyPath string) error {
	if signaturePath == "" {
		signaturePath = targetPath + ".sig"
	}
	if publicKeyPath == "" {
		self, err := os.Executable()
		if err != nil {
			return err
		}
		publicKeyPath = filepath.Join(filepath.Dir(self), "public-key.sig")
	}
	signature, err := ioutil.ReadFile(signaturePath)
	if err != nil {
		log.Println("Error in reading signature file: " + err.Error())
		return err
	}

	target, err := ioutil.ReadFile(targetPath)
	if err != nil {
		log.Println("Error in reading target file: " + err.Error())
		return err
	}

	publickey, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		log.Println("Error in reading public key file: " + err.Error())
		return err
	}

	if !ed25519.Verify(publickey, target, signature) {
		return fmt.Errorf("verify executable agianst signature faild")
	}
	return nil
}
