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

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

var name = flag.String("name", "", "name of the test to run")

func main() {
	flag.Parse()

	// Find out which test we're running.
	test, ok := Tests[*name]
	if !ok {
		log.Fatalf("No test found named %q", *name)
	}
	log.Printf("Running test %q", *name)

	// Get the IP of the local process.
	ip, err := getIP()
	if err != nil {
		log.Fatal(err)
	}

	// Run the test.
	if err := test.ContainerAction(ip); err != nil {
		log.Fatalf("Failed running test %q: %v", *name, err)
	}
}

// getIP listens for a connection from the local process and returns the source
// IP of that connection.
func getIP() (net.IP, error) {
	localAddr := net.TCPAddr{
		Port: IPExchangePort,
	}
	listener, err := net.ListenTCP("tcp4", &localAddr)
	if err != nil {
		return net.IP{}, fmt.Errorf("failed listening for IP: %v", err)
	}
	defer listener.Close()
	conn, err := listener.AcceptTCP()
	if err != nil {
		return net.IP{}, fmt.Errorf("failed accepting IP: %v", err)
	}
	defer conn.Close()
	log.Printf("Connected to %v", conn.RemoteAddr())

	return conn.RemoteAddr().(*net.TCPAddr).IP, nil
}
