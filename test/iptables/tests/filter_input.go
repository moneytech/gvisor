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
	"fmt"
	"net"
	"time"
)

const dropPort = 2401
const acceptPort = 2402
const sendloopDuration = 2 * time.Second
const network = "udp4"

// FilterInputDropUDP tests that we can drop UDP traffic.
type FilterInputDropUDP struct{}

// Name implements Testcase.Name.
func (FilterInputDropUDP) Name() string {
	return "FilterInputDropUDP"
}

// ContainerAction implements Testcase.ContainerAction.
func (FilterInputDropUDP) ContainerAction(ip net.IP) error {
	if err := filterTable("-A", "INPUT", "-p", "udp", "-j", "DROP"); err != nil {
		return err
	}

	// Listen for UDP packets on dropPort.
	if n, err := listenUDP(dropPort, sendloopDuration); err == nil {
		return fmt.Errorf("packets on port %d should have been dropped, but got a packet with %d bytes", dropPort, n)
	} else if netErr, ok := err.(net.Error); !ok || !netErr.Timeout() {
		return fmt.Errorf("error reading: %v", err)
	}

	// At this point we know that reading timed out and never received a
	// packet.
	return nil
}

// LocalAction implements Testcase.LocalAction.
func (FilterInputDropUDP) LocalAction(ip net.IP, errChan chan error) {
	errChan <- sendUDPLoop(ip, dropPort, sendloopDuration)
}

// FilterInputDropUDPPort tests that we can drop UDP traffic by port.
type FilterInputDropUDPPort struct{}

// Name implements Testcase.Name.
func (FilterInputDropUDPPort) Name() string {
	return "FilterInputDropUDPPort"
}

// ContainerAction implements Testcase.ContainerAction.
func (FilterInputDropUDPPort) ContainerAction(ip net.IP) error {
	if err := filterTable("-A", "INPUT", "-p", "udp", "-m", "udp", "--destination-port", fmt.Sprintf("%d", dropPort), "-j", "DROP"); err != nil {
		return err
	}

	// Listen for UDP packets on dropPort.
	if n, err := listenUDP(dropPort, sendloopDuration); err == nil {
		return fmt.Errorf("packets on port %d should have been dropped, but got a packet with %d bytes", dropPort, n)
	} else if netErr, ok := err.(net.Error); !ok || !netErr.Timeout() {
		return fmt.Errorf("error reading: %v", err)
	}

	// At this point we know that reading timed out and never received a
	// packet.
	return nil
}

// LocalAction implements Testcase.LocalAction.
func (FilterInputDropUDPPort) LocalAction(ip net.IP, errChan chan error) {
	errChan <- sendUDPLoop(ip, dropPort, sendloopDuration)
}

// FilterInputDropDifferentUDPPort tests that dropping traffic for a single UDP port
// doesn't drop packets on other ports.
type FilterInputDropDifferentUDPPort struct{}

// Name implements Testcase.Name.
func (FilterInputDropDifferentUDPPort) Name() string {
	return "FilterInputDropDifferentUDPPort"
}

// ContainerAction implements Testcase.ContainerAction.
func (FilterInputDropDifferentUDPPort) ContainerAction(ip net.IP) error {
	if err := filterTable("-A", "INPUT", "-p", "udp", "-m", "udp", "--destination-port", fmt.Sprintf("%d", dropPort), "-j", "DROP"); err != nil {
		return err
	}

	// Listen for UDP packets on another port.
	if _, err := listenUDP(acceptPort, sendloopDuration); err != nil {
		return fmt.Errorf("packets on port %d should be allowed, but encountered an error: %v", acceptPort, err)
	}

	return nil
}

// LocalAction implements Testcase.LocalAction.
func (FilterInputDropDifferentUDPPort) LocalAction(ip net.IP, errChan chan error) {
	errChan <- sendUDPLoop(ip, acceptPort, sendloopDuration)
}
