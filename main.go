// Copyright 2015 Daniel Theophanes.
// Use of this source code is governed by a zlib-style
// license that can be found in the LICENSE file.

// Simple service that only works by printing a log message every few seconds.
package main

import (
	"flag"
	"log"
	"time"

	"github.com/kardianos/service"
	"github.com/sirupsen/logrus"
	"github.com/eosioafrica/prodplugman/consul"
	"github.com/hashicorp/consul/api"
)

var logger service.Logger

// Program structures.
//  Define Start and Stop methods.
type program struct {
	exit chan struct{}
}

func (p *program) Start(s service.Service) error {

	logrus.Info("Starting producer plugin manager.")
	if service.Interactive() {
		logrus.Info("Running in terminal.")
	} else {
		logrus.Info("Running under service manager.")
	}
	p.exit = make(chan struct{})

	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() error {
	logrus.Infof("I'm running %v.", service.Platform())
	ticker := time.NewTicker(5 * time.Second) // Pull this from consul???

	//s, err := consul.New("", *ttl)

	for {
		select {
		case tm := <-ticker.C:

			// Read current state
			//s.ProductionState()

			// Get value from consul

			// If different request lock from consul

			// If you get lock, update production status

			logrus.Infof("Still running at %v...", tm)
		case <-p.exit:
			ticker.Stop()
			return nil
		}
	}
}
func (p *program) Stop(s service.Service) error {
	// Any work in Stop should be quick, usually a few seconds at most.
	logrus.Info("I'm Stopping!")
	close(p.exit)
	return nil
}

// Service setup.
//   Define service config.
//   Create the service.
//   Setup the logger.
//   Handle service controls (optional).
//   Run the service.
func main() {
	svcFlag := flag.String("service", "", "Control the system service.")
	flag.Parse()

	svcConfig := &service.Config{
		Name:        "ProducerPluginManager",
		DisplayName: "A systemd service that manages the state of producer_plugin",
		Description: "Uses producer plugin api to pause or restart block production on nodeos.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	errs := make(chan error, 5)
	logger, err = s.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}

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

func ServiceHandler()  {

	ttl := flag.Duration("ttl", time.Second*15, "Service TTL check duration")
	c, _ := api.NewClient(api.DefaultConfig())

	s, err := consul.New(c, "", *ttl)

	s.ProductionState()

	if s.PPMan.Err != nil {

		return
	}

	s.
}