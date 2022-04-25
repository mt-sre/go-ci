package command

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// NewCommandAlias makes reuse of common commands possible by taking
// a command name and variadic slice of command arguments and returning
// a function which will accept an additional variadic slice of
// "CommandOption" values to augment the command behavior.
func NewCommandAlias(name string, args ...string) func(...CommandOption) Command {
	return func(opts ...CommandOption) Command {
		opts = append([]CommandOption{WithArgs(args)}, opts...)

		return NewCommand(name, opts...)
	}
}

// NewCommand takes a command name and a variadic slice of "CommandOption"
// values and returns a "Command" which may be invoked by it's "Run" method.
// The given command name must be resolvable on the system PATH.
func NewCommand(name string, opts ...CommandOption) Command {
	var cfg CommandConfig

	cfg.Option(opts...)
	cfg.Default()

	cmd := exec.CommandContext(cfg.Ctx, name)

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if cfg.Verbose {
		cmd.Stdout = io.MultiWriter(cmd.Stdout, os.Stdout)
		cmd.Stderr = io.MultiWriter(cmd.Stderr, os.Stdout)
	}

	if len(cfg.Args) > 0 {
		cmd.Args = append(cmd.Args, cfg.Args...)
	}

	if cfg.WithCurrentEnv {
		// prepend current env so user provided values take precedence
		cfg.Env = append(os.Environ(), cfg.Env...)
	}

	if len(cfg.Env) > 0 {
		cmd.Env = cfg.Env
	}

	if cfg.WorkDir != "" {
		cmd.Dir = cfg.WorkDir
	}

	if cfg.Stdin != nil {
		cmd.Stdin = cfg.Stdin
	}

	return Command{
		cmd:    cmd,
		stdout: &stdout,
		stderr: &stderr,
	}
}

// Command abstracts a shell command.
type Command struct {
	cmd    *exec.Cmd
	stdout *bytes.Buffer
	stderr *bytes.Buffer
}

// Run executes a "Command" instance and returns an error if
// either the command was not able to start.
func (c *Command) Run() error {
	if err := c.cmd.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); !ok { //nolint:errorlint
			return fmt.Errorf("running command %q: %w", strings.Join(c.cmd.Args, " "), err)
		}
	}

	return nil
}

func (c *Command) ExitCode() int          { return c.cmd.ProcessState.ExitCode() }
func (c *Command) Error() error           { return &CommandError{State: c.cmd.ProcessState} }
func (c *Command) CombinedOutput() string { return c.Stdout() + c.Stderr() }
func (c *Command) Stdout() string         { return c.stdout.String() }
func (c *Command) Stderr() string         { return c.stderr.String() }
func (c *Command) Success() bool          { return c.cmd.ProcessState.Success() }

type CommandError struct {
	State *os.ProcessState
}

func (e *CommandError) Error() string { return e.State.String() }

type CommandConfig struct {
	Args           []string
	Ctx            context.Context
	Env            []string
	Stdin          io.Reader
	Verbose        bool
	WithCurrentEnv bool
	WorkDir        string
}

func (c *CommandConfig) Option(opts ...CommandOption) {
	for _, opt := range opts {
		opt.ConfigureCommand(c)
	}
}

func (c *CommandConfig) Default() {
	if c.Ctx == nil {
		c.Ctx = context.Background()
	}
}

type CommandOption interface {
	ConfigureCommand(*CommandConfig)
}

type WithArgs []string

func (wa WithArgs) ConfigureCommand(c *CommandConfig) {
	c.Args = append(c.Args, wa...)
}

type WithContext struct{ context.Context }

func (wc WithContext) ConfigureCommand(c *CommandConfig) {
	c.Ctx = wc.Context
}

type WithCurrentEnv bool

func (wc WithCurrentEnv) ConfigureCommand(c *CommandConfig) {
	c.WithCurrentEnv = bool(wc)
}

type WithEnv map[string]string

func (we WithEnv) ConfigureCommand(c *CommandConfig) {
	for k, v := range we {
		c.Env = append(c.Env, fmt.Sprintf("%s=%s", k, v))
	}
}

type WithStdin struct{ io.Reader }

func (ws WithStdin) ConfigureCommand(c *CommandConfig) {
	c.Stdin = ws.Reader
}

type WithConsoleOut bool

func (wv WithConsoleOut) ConfigureCommand(c *CommandConfig) {
	c.Verbose = bool(wv)
}

type WithWorkingDirectory string

func (wd WithWorkingDirectory) ConfigureCommand(c *CommandConfig) {
	c.WorkDir = string(wd)
}
