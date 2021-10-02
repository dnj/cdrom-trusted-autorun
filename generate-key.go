package main

import (
	"crypto/ed25519"
	"log"
	"os"
	"path/filepath"

	"github.com/alexflint/go-arg"
)

func generateKey(parser *arg.Parser, args Args) {
	if args.GenerateKey.PublicKey == "" {
		self, err := os.Executable()
		if err == nil {
			args.GenerateKey.PublicKey = filepath.Join(filepath.Dir(self), "public-key.sig")
		}
	}
	if args.GenerateKey.PrivateKey == "" {
		self, err := os.Executable()
		if err == nil {
			args.GenerateKey.PrivateKey = filepath.Join(filepath.Dir(self), "private-key.sig")
		}
	}
	error := GenerateKey(args.GenerateKey.PublicKey, args.GenerateKey.PrivateKey)
	if error != nil {
		log.Fatal(error)
	}
}

func GenerateKey(publicKeyPath string, privateKeyPath string) error {
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return err
	}
	err = os.WriteFile(publicKeyPath, publicKey, 0644)
	if err != nil {
		return err
	}
	err = os.WriteFile(privateKeyPath, privateKey, 0644)
	if err != nil {
		return err
	}
	return nil
}
