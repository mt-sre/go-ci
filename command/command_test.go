// SPDX-FileCopyrightText: 2025 Red Hat, Inc. <sd-mt-sre@redhat.com>
//
// SPDX-License-Identifier: Apache-2.0

package command

import (
	"bytes"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommand(t *testing.T) {
	t.Parallel()

	if runtime.GOOS == "windows" {
		t.Skip("skipping command tests on Windows")
	}

	for name, tc := range map[string]struct {
		Assertion        require.ErrorAssertionFunc
		Command          Command
		ExpectedOutput   string
		ExpectedExitCode int
	}{
		"capture stdout": {
			Assertion:        require.NoError,
			Command:          NewCommand("echo", WithArgs{"-n", "hello"}),
			ExpectedOutput:   "hello",
			ExpectedExitCode: 0,
		},
		"passing env": {
			Assertion:        require.NoError,
			Command:          NewCommand("env", WithEnv{"HELLO": "hello"}),
			ExpectedOutput:   "HELLO=hello\n",
			ExpectedExitCode: 0,
		},
		"using stdin": {
			Assertion:        require.NoError,
			Command:          NewCommand("cat", WithStdin{bytes.NewBufferString("hello")}),
			ExpectedOutput:   "hello",
			ExpectedExitCode: 0,
		},
		"using working directory": {
			Assertion:        require.NoError,
			Command:          NewCommand("ls", WithArgs{"-d", "/"}, WithWorkingDirectory("/")),
			ExpectedOutput:   "/\n",
			ExpectedExitCode: 0,
		},
		"failing with bad command": {
			Assertion:        require.Error,
			Command:          NewCommand("dne"),
			ExpectedOutput:   "",
			ExpectedExitCode: -1,
		},
		"failing with exit code": {
			Assertion:        require.NoError,
			Command:          NewCommand("cat", WithArgs{"-z"}),
			ExpectedOutput:   "",
			ExpectedExitCode: 1,
		},
	} {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tc.Assertion(t, tc.Command.Run())
			assert.Equal(t, tc.ExpectedOutput, tc.Command.Stdout())
			assert.Equal(t, tc.ExpectedExitCode, tc.Command.ExitCode(), tc.Command.CombinedOutput())

		})
	}
}

func TestCommandErrorInterfaces(t *testing.T) {
	t.Parallel()

	require.Implements(t, new(error), new(CommandError))
}

// TestNewCommandAlias tests a new command alias path and
// command arguments.
func TestNewCommandAlias(t *testing.T) {
	testAlias := NewCommandAlias("ls", "-la")

	cmd := testAlias(WithCurrentEnv(true))

	assert.Equal(t, []string{"ls", "-la"}, cmd.cmd.Args, "expected command arguments [ls -la]")
}
