// SPDX-FileCopyrightText: 2025 Red Hat, Inc. <sd-mt-sre@redhat.com>
//
// SPDX-License-Identifier: Apache-2.0

package gocmd

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/mt-sre/go-ci/command"
)

func NewGoCmd(opts ...GoCmdOption) (*GoCmd, error) {
	var cfg GoCmdConfig

	cfg.Option(opts...)
	if err := cfg.Default(); err != nil {
		return nil, fmt.Errorf("applying defaults: %w", err)
	}

	return &GoCmd{
		cfg: cfg,
	}, nil
}

type GoCmd struct {
	cfg GoCmdConfig
}

var ErrModuleNotFound = errors.New("module not found")

// Module retrieves the current Go module name or
// an error if a module cannot be identified.
func (c *GoCmd) Module(ctx context.Context, opts ...ModuleOption) (string, error) {
	var cfg ModuleConfig

	cfg.Option(opts...)

	cmdOpts := []command.CommandOption{
		command.WithContext{Context: ctx},
		command.WithArgs{"mod", "why", "-m"},
	}

	if cfg.WorkingDir != "" {
		cmdOpts = append(cmdOpts, command.WithWorkingDirectory(cfg.WorkingDir))
	}

	why := command.NewCommand(c.cfg.BinPath, cmdOpts...)
	if err := why.Run(); err != nil {
		return "", fmt.Errorf("starting to get module information: %w", err)
	}

	if !why.Success() {
		return "", fmt.Errorf("getting module information; %w", why.Error())
	}

	fields := strings.Fields(strings.TrimSpace(why.Stdout()))

	if len(fields) < 2 {
		return "", ErrModuleNotFound
	}

	return fields[1], nil
}

type ModuleConfig struct {
	WorkingDir string
}

func (c *ModuleConfig) Option(opts ...ModuleOption) {
	for _, opt := range opts {
		opt.ConfigureModule(c)
	}
}

type ModuleOption interface {
	ConfigureModule(*ModuleConfig)
}

func (c *GoCmd) Tidy(ctx context.Context, opts ...TidyOption) error {
	var cfg TidyConfig

	cfg.Option(opts...)

	args := []string{"mod", "tidy"}

	if cfg.GoCompatability != "" {
		args = append(args, "-compat="+cfg.GoCompatability)
	}

	if cfg.GoVersion != "" {
		args = append(args, "-go="+cfg.GoVersion)
	}

	cmdOpts := []command.CommandOption{
		command.WithContext{Context: ctx},
		command.WithArgs(args),
	}

	if cfg.WorkingDir != "" {
		cmdOpts = append(cmdOpts, command.WithWorkingDirectory(cfg.WorkingDir))
	}

	tidy := command.NewCommand(c.cfg.BinPath, cmdOpts...)
	if err := tidy.Run(); err != nil {
		return fmt.Errorf("starting to tidy module: %w", err)
	}

	if !tidy.Success() {
		return fmt.Errorf("tidying module: %w", tidy.Error())
	}

	return nil
}

type TidyConfig struct {
	GoCompatability string
	GoVersion       string
	WorkingDir      string
}

func (c *TidyConfig) Option(opts ...TidyOption) {
	for _, opt := range opts {
		opt.ConfigureTidy(c)
	}
}

type TidyOption interface {
	ConfigureTidy(*TidyConfig)
}

type GoCmdConfig struct {
	BinPath string
}

func (c *GoCmdConfig) Option(opts ...GoCmdOption) {
	for _, opt := range opts {
		opt.ConfigureGoCmd(c)
	}
}

func (c *GoCmdConfig) Default() error {
	if c.BinPath == "" {
		path, err := exec.LookPath("go")
		if err != nil {
			return fmt.Errorf("looking up 'go' in PATH: %w", err)
		}

		c.BinPath = path
	}

	return nil
}

type GoCmdOption interface {
	ConfigureGoCmd(c *GoCmdConfig)
}
