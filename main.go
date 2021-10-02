package main

import (
	"os"

	"github.com/alexflint/go-arg"
)

type Drive struct {
	path    string
	fsLabel string
}

type WatchCmd struct {
	PublicKey string `arg:"-p" help:"Path to public key file"`
}

type SignCmd struct {
	Signature  string `arg:"-s,--signature" help:"Path to signature file"`
	Target     string `arg:"-t,--target,required" help:"Path to target file"`
	PrivateKey string `arg:"-k,required" help:"Path to private key file"`
}

type VerifyCmd struct {
	Signature string `arg:"-s,--signature" help:"Path to signature file"`
	Target    string `arg:"-t,--target,required" help:"Path to target file"`
	PublicKey string `arg:"-p" help:"Path to public key file"`
}
type GenerateKeyCmd struct {
	PrivateKey string `arg:"-k" help:"Path to private key file. If file already exists it will be overwrite" placeholder:"PATH"`
	PublicKey  string `arg:"-p" help:"Path to public key file. If file already exists it will be overwrite" placeholder:"PATH"`
}
type ServiceCmd struct {
}

type Args struct {
	Service     *ServiceCmd     `arg:"subcommand:service"`
	Watch       *WatchCmd       `arg:"subcommand:watch"`
	Sign        *SignCmd        `arg:"subcommand:sign"`
	Verify      *VerifyCmd      `arg:"subcommand:verify"`
	GenerateKey *GenerateKeyCmd `arg:"subcommand:generate-key"`
}

func (Args) Version() string {
	return "cdrom-trusted-autorun 1.0.0"
}

func main() {

	var args Args
	parser := arg.MustParse(&args)

	switch {
	case args.GenerateKey != nil:
		generateKey(parser, args)
	case args.Sign != nil:
		sign(parser, args)
	case args.Verify != nil:
		verify(parser, args)
	case args.Watch != nil:
		watch(parser, args)
	case args.Service != nil:
		RunAsService()
	default:
		parser.WriteUsage(os.Stderr)
	}
}
