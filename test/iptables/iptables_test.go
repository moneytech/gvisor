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

package iptables

import (
	"fmt"
	"net"
	"os"
	"path"
	"testing"
	"time"

	"gvisor.dev/gvisor/pkg/log"
	"gvisor.dev/gvisor/runsc/dockerutil"
	tests "gvisor.dev/gvisor/test/iptables/tests"
)

const timeout time.Duration = 10 * time.Second

type result struct {
	output string
	err    error
}

// singleTest runs a Testcase. Each test follows a pattern:
// - Create a container.
// - Get the container's IP.
// - Send the container our IP.
// - Start a new goroutine running the local action of the test.
// - Wait for both the container and local actions to finish.
//
// Container output is logged to $TEST_UNDECLARED_OUTPUTS_DIR if it exists, or
// to stderr.
func singleTest(test tests.Testcase) error {
	if _, ok := tests.Tests[test.Name()]; !ok {
		return fmt.Errorf("no test found with name %q. Has it been added to tests.Tests?", test.Name())
	}

	// Create and start the container.
	cont := dockerutil.MakeDocker("gvisor-iptables")
	defer cont.CleanUp()
	resultChan := make(chan *result)
	go func() {
		output, err := cont.RunFg("--cap-add=NET_ADMIN", "iptables-tests", "-name", test.Name())
		logContainer(output, err)
		resultChan <- &result{output, err}
	}()

	// Get the container IP.
	ip, err := getIP(cont)
	if err != nil {
		return fmt.Errorf("failed to get container IP: %v", err)
	}

	// Give the container our IP.
	if err := sendIP(ip); err != nil {
		return fmt.Errorf("failed to send IP to container: %v", err)
	}

	// Run our side of the test.
	errChan := make(chan error)
	go test.LocalAction(ip, errChan)

	// Wait for both the container and local tests to finish.
	var res *result
	to := time.After(timeout)
	for localDone := false; res == nil || !localDone; {
		select {
		case res = <-resultChan:
			log.Infof("Container finished.")
		case err, localDone = <-errChan:
			log.Infof("Local finished.")
			if err != nil {
				return fmt.Errorf("local test failed: %v", err)
			}
		case <-to:
			return fmt.Errorf("timed out after %f seconds", timeout.Seconds())
		}
	}

	return res.err
}

func getIP(cont dockerutil.Docker) (net.IP, error) {
	// The container might not have started yet, so retry a few times.
	var ipStr string
	to := time.After(timeout)
	for ipStr == "" {
		ipStr, _ = cont.FindIP()
		select {
		case <-to:
			return net.IP{}, fmt.Errorf("timed out getting IP after %f seconds", timeout.Seconds())
		default:
			time.Sleep(250 * time.Millisecond)
		}
	}
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return net.IP{}, fmt.Errorf("invalid IP: %q", ipStr)
	}
	log.Infof("Container has IP of %s", ipStr)
	return ip, nil
}

func sendIP(ip net.IP) error {
	contAddr := net.TCPAddr{
		IP:   ip,
		Port: tests.IPExchangePort,
	}
	to := time.After(timeout)
	var conn *net.TCPConn
	var err error
	// The container may not be listening when we first connect, so retry
	// upon error.
	for {
		conn, err = net.DialTCP("tcp4", nil, &contAddr)
		if err != nil {
			select {
			case <-to:
				return fmt.Errorf("timed out waiting to send IP, most recent error: %v", err)
			default:
				time.Sleep(200 * time.Millisecond)
				continue
			}
		}
		break
	}
	if _, err := conn.Write([]byte{0}); err != nil {
		return fmt.Errorf("error writing to container: %v", err)
	}
	return nil
}

func logContainer(output string, err error) {
	msg := fmt.Sprintf("Container error: %v\nContainer output:\n%v", err, output)
	if artifactsDir := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); artifactsDir != "" {
		fpath := path.Join(artifactsDir, "container.log")
		if file, err := os.Create(fpath); err != nil {
			log.Warningf("Failed to open log file %q: %v", fpath, err)
		} else {
			defer file.Close()
			if _, err := file.Write([]byte(msg)); err != nil {
				log.Warningf("Failed to write to log file %s: %v", fpath, err)
			} else {
				return
			}
		}
	}

	// We couldn't write to the output directory -- just log to stderr.
	log.Infof("%s", msg)
}

func TestFilterInputDropUDP(t *testing.T) {
	if err := singleTest(tests.FilterInputDropUDP{}); err != nil {
		t.Fatal(err)
	}
}

func TestFilterInputDropUDPPort(t *testing.T) {
	if err := singleTest(tests.FilterInputDropUDPPort{}); err != nil {
		t.Fatal(err)
	}
}

func TestFilterInputDropDifferentUDPPort(t *testing.T) {
	if err := singleTest(tests.FilterInputDropDifferentUDPPort{}); err != nil {
		t.Fatal(err)
	}
}
