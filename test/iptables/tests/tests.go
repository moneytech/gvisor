// Copyright 2019 The gVisor Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package main contains two main components:
// - A set of iptables tests implemented as Testcase implementations.
// - A main() function to run tests, usually inside a container.
//
// In order to make using `go build` simple in the Dockerfile, this package may
// not depend on packages outside the standard library.
package main

import (
	"net"
)

// IPExchangePort is the port the container listens on to receive the IP
// address of the local process.
const IPExchangePort = 2349

// A Testcase contains one action to run in the container and one to run
// locally. Each action must succeed for the test pass.
type Testcase interface {
	// Name returns the name of the test.
	Name() string

	// ContainerAction runs inside the container. It receives the IP of the
	// local process.
	ContainerAction(ip net.IP) error

	// LocalAction runs locally. It receives the IP of the container. It
	// also receives a channel on which to report errors, as LocalAction is
	// always called in its own goroutine. It MUST send either an error or
	// nil do indicate when the test is done.
	LocalAction(ip net.IP, errChan chan error)
}

// Tests maps test names to Testcases.
//
// New Testcases should be added via the init() function of this package.
var Tests = map[string]Testcase{}

func init() {
	tcs := []Testcase{
		FilterInputDropUDP{},
		FilterInputDropUDPPort{},
		FilterInputDropDifferentUDPPort{},
	}
	for _, tc := range tcs {
		Tests[tc.Name()] = tc
	}
}
