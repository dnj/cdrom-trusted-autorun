package main

import (
	"log"

	"github.com/alexflint/go-arg"
	"github.com/kardianos/service"
)

var logger service.Logger

type Program struct{}

func (p *Program) Start(s service.Service) error {
	var args Args
	args.Watch = &WatchCmd{}
	parser := arg.MustParse(&args)
	go watch(parser, args)
	return nil
}
func (p *Program) Stop(s service.Service) error {
	return nil
}

func RunAsService() {
	config := &service.Config{
		Name:        "AutorunSvc",
		DisplayName: "CD-ROM Trusted Autorun",
		Description: "Simple process which watch for new media in your machine's CDROM and automatically run it.",
	}

	program := &Program{}
	serviceInstance, err := service.New(program, config)
	if err != nil {
		log.Fatal(err)
	}
	logger, err = serviceInstance.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	err = serviceInstance.Run()
	if err != nil {
		logger.Error(err)
	}
}
