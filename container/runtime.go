// SPDX-FileCopyrightText: 2025 Red Hat, Inc. <sd-mt-sre@redhat.com>
//
// SPDX-License-Identifier: Apache-2.0

package container

import (
	"os/exec"
)

// Runtime attempts to find an available container runtime in the PATH.
// The path to the first available runtime is returned along with a boolean
// value indicating if any runtimes were found.
func Runtime() (string, bool) {
	prefferedRuntimes := []string{
		"podman",
		"docker",
	}

	for _, runtime := range prefferedRuntimes {
		runtimePath, err := exec.LookPath(runtime)
		if err == nil {
			return runtimePath, true
		}
	}

	return "", false
}
