// Copyright 2013-2015 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Diameter server example. This is by no means a complete server.
//
// If you'd like to test diameter over SSL, generate SSL certificates:
//   go run $GOROOT/src/crypto/tls/generate_cert.go --host localhost
//
// And start the server with `-cert_file cert.pem -key_file key.pem`.
//
// By default this server runs in a single OS thread. If you want to
// make it run on more, set the GOMAXPROCS=n environment variable.
// See Go's FAQ for details: http://golang.org/doc/faq#Why_no_multi_CPU

package common

import (
	"os"

	parser "Hanif_Aulia_Sabri-MyTrip/git/order/parser"

	log "github.com/Sirupsen/logrus"
)

//Config stores global configuration loaded from json file
type Configuration struct {
	ListenPort  string `yaml:"listenPort"`
	RootURL     string `yaml:"rootUrl"`
	RootURLBaru string `yaml:"rootUrlBaru"`
	Connection  struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Password string `yaml:"password"`
		User     string `yaml:"user"`
		Database string `yaml:"database"`
	}
	DatabaseFile string `yaml:"databaseFile"`
}

var Config Configuration
var logger *log.Entry

func LoadConfigFromFile(fn *string) {
	if err := parser.LoadYAML(fn, &Config); err != nil {
		log.Error("LoadConfigFromFile() - Failed opening config file")
		os.Exit(1)
	}

	log.Info("Loaded configs: ", Config)

}
