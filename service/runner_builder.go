package service

import (
	"flag"
	"github.com/ayufan/golang-kardianos-service"
	"log"
	"strings"
)

func NewServiceRunnerBuilder(serviceName string, runHandler RunHandler) ServiceRunnerBuilder {
	if strings.TrimSpace(serviceName) == "" {
		panic("The service name cannot be blank in NewServiceRunnerBuilder")
	}

	return &builder{
		ServiceName:        serviceName,
		ServiceDisplayName: serviceName, //default to name
		ServiceDescription: serviceName, //default to name
		ServiceUserName:    "",
		RunHandler:         runHandler,
	}
}

type ServiceRunnerBuilder interface {
	WithServiceDisplayName(serviceDisplayName string) ServiceRunnerBuilder
	WithServiceDescription(serviceDescription string) ServiceRunnerBuilder

	WithWorkingDirectory(dir string) ServiceRunnerBuilder
	WithAdditionalArguments(args ...string) ServiceRunnerBuilder

	WithServiceUserName(serviceUserName string) ServiceRunnerBuilder
	WithServiceUserName_AsCurrentUser() ServiceRunnerBuilder

	WithOnStopHandler(h OnStopHandler) ServiceRunnerBuilder

	Run()
}

type builder struct {
	ServiceName        string
	ServiceDisplayName string
	ServiceDescription string

	WorkingDirectory    string
	AdditionalArguments []string

	ServiceUserName string

	RunHandler    RunHandler
	OnStopHandler OnStopHandler
}

func (b *builder) WithServiceDisplayName(serviceDisplayName string) ServiceRunnerBuilder {
	b.ServiceDisplayName = serviceDisplayName
	return b
}

func (b *builder) WithServiceDescription(serviceDescription string) ServiceRunnerBuilder {
	b.ServiceDescription = serviceDescription
	return b
}

func (b *builder) WithWorkingDirectory(dir string) ServiceRunnerBuilder {
	b.WorkingDirectory = dir
	return b
}

func (b *builder) WithAdditionalArguments(args ...string) ServiceRunnerBuilder {
	b.AdditionalArguments = args
	return b
}

func (b *builder) WithServiceUserName(serviceUserName string) ServiceRunnerBuilder {
	b.ServiceUserName = serviceUserName
	return b
}

func (b *builder) WithServiceUserName_AsCurrentUser() ServiceRunnerBuilder {
	return b.WithServiceUserName(getCurrentUserName())
}

func (b *builder) WithOnStopHandler(h OnStopHandler) ServiceRunnerBuilder {
	b.OnStopHandler = h
	return b
}

func (b *builder) Run() {
	svcFlag := flag.String("service", "", "Control the system service.")
	flag.Parse()

	svcConfig := &service.Config{
		Name:             b.ServiceName,
		DisplayName:      b.ServiceDisplayName,
		Description:      b.ServiceDescription,
		WorkingDirectory: b.WorkingDirectory,
		Arguments:        b.AdditionalArguments,
		UserName:         b.ServiceUserName,
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	errs := make(chan error, 5)
	logger, err := s.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}

	prg.Logger = logger
	prg.RunHandler = b.RunHandler
	prg.OnStopHandler = b.OnStopHandler

	go func() {
		for {
			err := <-errs
			if err != nil {
				log.Print(err)
			}
		}
	}()

	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}

	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
