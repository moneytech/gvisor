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

package strace

import (
	"fmt"

	"gvisor.dev/gvisor/pkg/sentry/kernel"
	slinux "gvisor.dev/gvisor/pkg/sentry/syscalls/linux"
	"gvisor.dev/gvisor/pkg/sentry/usermem"
)

func setToStr(t *kernel.Task, set []byte) string {
	var fds []int
	for i, v := range set {
		if v != 0 {
			fds = append(fds, i)
		}
	}
	return fmt.Sprint(fds)
}

func fdSets(t *kernel.Task, nfds int, readAddr, writeAddr, exceptAddr usermem.Addr) string {
	r, w, e, err := slinux.CopyInFDSets(t, nfds, readAddr, writeAddr, exceptAddr)
	if err != nil {
		return fmt.Sprintf("%#x, %#x, %#x (error decoding fdsets: %s)", readAddr, writeAddr, exceptAddr, err)
	}

	return fmt.Sprintf("%#x %s, %#x %s, %#x %s",
		readAddr, setToStr(t, r), writeAddr, setToStr(t, w), exceptAddr, setToStr(t, e))
}
