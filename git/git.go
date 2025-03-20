// SPDX-FileCopyrightText: 2025 Red Hat, Inc. <sd-mt-sre@redhat.com>
//
// SPDX-License-Identifier: Apache-2.0

package git

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/mt-sre/go-ci/command"
)

// RevParse runs "git rev-parse" given a RevParseFormat and a variadic
// slice of options. An error is returned if a value cannot be obtained.
func RevParse(ctx context.Context, opts ...RevParseOption) (string, error) {
	var cfg RevParseConfig

	cfg.Option(opts...)

	args := []string{"rev-parse"}

	if cfg.Format != "" {
		args = append(args, cfg.Format.ToGitValue())

	}

	if cfg.Format != RevParseFormatTopLevel {
		args = append(args, "HEAD")
	}

	cmdOpts := []command.CommandOption{
		command.WithContext{Context: ctx},
		command.WithArgs(args),
	}

	if cfg.WorkingDir != "" {
		cmdOpts = append(cmdOpts, command.WithWorkingDirectory(cfg.WorkingDir))
	}

	revParse := git(cmdOpts...)
	if err := revParse.Run(); err != nil {
		return "", fmt.Errorf("starting to run rev-parse directory: %w", err)
	}

	if !revParse.Success() {
		return "", fmt.Errorf("running rev-parse directory: %w", revParse.Error())
	}

	return strings.TrimSpace(revParse.Stdout()), nil
}

type RevParseConfig struct {
	Format     RevParseFormat
	WorkingDir string
}

func (c *RevParseConfig) Option(opts ...RevParseOption) {
	for _, opt := range opts {
		opt.ConfigureRevParse(c)
	}
}

type RevParseOption interface {
	ConfigureRevParse(*RevParseConfig)
}

// RevParseFormat is an enum of format
// options for "git rev-parse"
type RevParseFormat string

func (f RevParseFormat) ToGitValue() string {
	switch f {
	case RevParseFormatAbbrevRef:
		return "--abbrev-ref"
	case RevParseFormatShort:
		return "--short"
	case RevParseFormatTopLevel:
		return "--show-toplevel"
	default:
		return ""
	}
}

const (
	// RevParseFormatAbbrevRef is "--abbrev-ref"
	RevParseFormatAbbrevRef RevParseFormat = "abbrev-ref"
	// RevParseFormatShort is "--short"
	RevParseFormatShort RevParseFormat = "short"
	// RevParseFormatTopLevel is "--show-toplevel"
	RevParseFormatTopLevel RevParseFormat = "top-level"
)

// ErrNoTagsFound is returned when no tags
// are found in the current git repository
var ErrNoTagsFound = errors.New("no tags found")

// LatestTag returns the latest tag for the current git repository
// given a variadic slice of options. An error is returned
// if the latest tag cannot be retrieved.
func LatestTag(ctx context.Context, opts ...LatestTagOption) (string, error) {
	var cfg LatestTagConfig

	cfg.Option(opts...)

	listOpts := []ListTagsOption{WithSorted(true)}
	if cfg.WorkingDir != "" {
		listOpts = append(listOpts, WithWorkingDirectory(cfg.WorkingDir))
	}

	tags, err := ListTags(ctx, listOpts...)
	if err != nil {
		return "", fmt.Errorf("listing tags: %w", err)
	}

	if len(tags) < 1 {
		return "", ErrNoTagsFound
	}

	return tags[0], nil
}

type LatestTagConfig struct {
	WorkingDir string
}

func (c *LatestTagConfig) Option(opts ...LatestTagOption) {
	for _, opt := range opts {
		opt.ConfigureLatestTag(c)
	}
}

type LatestTagOption interface {
	ConfigureLatestTag(*LatestTagConfig)
}

// LatestVersion returns the latest tag, "v" prefixed,
// as determined by comparing all tags as semantic versions.
// An error will be returned if there are no available tags
// or any tag is not parseable as a semantic version.
func LatestVersion(ctx context.Context, opts ...LatestVersionOption) (string, error) {
	var cfg LatestVersionConfig

	cfg.Option(opts...)

	var listOpts []ListTagsOption

	if cfg.WorkingDir != "" {
		listOpts = append(listOpts, WithWorkingDirectory(cfg.WorkingDir))
	}

	tags, err := ListTags(ctx, listOpts...)
	if err != nil {
		return "", fmt.Errorf("listing tags: %w", err)
	}

	if len(tags) < 1 {
		return "", ErrNoTagsFound
	}

	var versions []semver.Version

	for _, t := range tags {
		ver, err := semver.ParseTolerant(t)
		if err != nil {
			return "", fmt.Errorf("parsing version %q: %w", ver, err)
		}

		versions = append(versions, ver)
	}

	semver.Sort(versions)

	return "v" + versions[len(versions)-1].String(), nil
}

type LatestVersionConfig struct {
	WorkingDir string
}

func (c *LatestVersionConfig) Option(opts ...LatestVersionOption) {
	for _, opt := range opts {
		opt.ConfigureLatestVersion(c)
	}
}

type LatestVersionOption interface {
	ConfigureLatestVersion(*LatestVersionConfig)
}

// ListTags lists all tags in the current git repository given
// a variadic slice of options. An error is returned if the
// the tags cannot be listed.
func ListTags(ctx context.Context, opts ...ListTagsOption) ([]string, error) {
	var cfg ListTagsConfig

	cfg.Option(opts...)

	args := command.WithArgs{"tag", "-l"}
	if cfg.Sorted {
		args = append(args, "--sort", cfg.SortKey.ToGitValue())
	}

	cmdOpts := []command.CommandOption{
		command.WithContext{Context: ctx},
		command.WithArgs(args),
	}

	if cfg.WorkingDir != "" {
		cmdOpts = append(cmdOpts, command.WithWorkingDirectory(cfg.WorkingDir))
	}

	listTags := git(cmdOpts...)
	if err := listTags.Run(); err != nil {
		return nil, fmt.Errorf("starting to list tags: %w", err)
	}

	if !listTags.Success() {
		return nil, fmt.Errorf("listing tags: %w", listTags.Error())
	}

	return strings.Fields(listTags.Stdout()), nil
}

type ListTagsConfig struct {
	Sorted     bool
	SortKey    SortKey
	WorkingDir string
}

func (c *ListTagsConfig) Option(opts ...ListTagsOption) {
	for _, opt := range opts {
		opt.ConfigureListTags(c)
	}
}

type ListTagsOption interface {
	ConfigureListTags(c *ListTagsConfig)
}

// SortKey is a key used to sort tags.
type SortKey string

func (k SortKey) ToGitValue() string {
	switch k {
	case SortKeyCreationDate:
		return "-creatordate"
	default:
		return "-refname"
	}
}

const (
	// SortKeyNone is an empty sort key.
	SortKeyNone SortKey = ""
	// SortKeyCreationDate is a key which sorts by
	// tag creation date.
	SortKeyCreationDate SortKey = "creation date"
	// SortKeyRefName is a key which sorts by refname.
	SortKeyRefName SortKey = "refname"
)

// Diff returns the current diff of the git repository
// with a specified format and variadic slice of options.
// An error is returned if the diff cannot be retrieved.
func Diff(ctx context.Context, opts ...DiffOption) (string, error) {
	var cfg DiffConfig

	cfg.Option(opts...)

	args := []string{"diff"}

	if cfg.Format != "" {
		args = append(args, cfg.Format.ToGitValue())
	}

	diffOpts := []command.CommandOption{
		command.WithContext{Context: ctx},
		command.WithArgs(args),
	}

	if cfg.WorkingDir != "" {
		diffOpts = append(diffOpts, command.WithWorkingDirectory(cfg.WorkingDir))
	}

	diff := git(diffOpts...)
	if err := diff.Run(); err != nil {
		return "", fmt.Errorf("starting to get git diff: %w", err)
	}

	if !diff.Success() {
		return "", fmt.Errorf("getting git diff: %w", diff.Error())
	}

	return strings.TrimSpace(diff.Stdout()), nil
}

type DiffConfig struct {
	Format     DiffFormat
	WorkingDir string
}

func (c *DiffConfig) Option(opts ...DiffOption) {
	for _, opt := range opts {
		opt.ConfigureDiff(c)
	}
}

type DiffOption interface {
	ConfigureDiff(*DiffConfig)
}

type DiffFormat string

const (
	// DiffFormatNameOnly returns only the names
	// of modified files.
	DiffFormatNameOnly DiffFormat = "name only"
	// DiffFormatNameStatus returns both the name
	// and type of change for modified files.
	DiffFormatNameStatus DiffFormat = "name status"
)

func (f DiffFormat) ToGitValue() string {
	switch f {
	case DiffFormatNameOnly:
		return "--name-only"
	case DiffFormatNameStatus:
		return "--name-status"
	default:
		return ""
	}
}

// Status returns the current status of the git repository with
// the given format and a variadic slice of options. An eror is
// returned if the status cannot be retrieved.
func Status(ctx context.Context, opts ...StatusOption) (string, error) {
	var cfg StatusConfig

	cfg.Option(opts...)

	args := []string{"status"}

	if cfg.Format != "" {
		args = append(args, cfg.Format.ToGitValue())
	}

	statusOpts := []command.CommandOption{
		command.WithContext{Context: ctx},
		command.WithArgs(args),
	}

	if cfg.WorkingDir != "" {
		statusOpts = append(statusOpts, command.WithWorkingDirectory(cfg.WorkingDir))
	}

	status := git(statusOpts...)

	if err := status.Run(); err != nil {
		return "", fmt.Errorf("starting to get git status: %w", err)
	}

	if !status.Success() {
		return "", fmt.Errorf("getting git status: %w", status.Error())
	}

	return strings.TrimSpace(status.Stdout()), nil
}

type StatusConfig struct {
	Format     StatusFormat
	WorkingDir string
}

func (c *StatusConfig) Option(opts ...StatusOption) {
	for _, opt := range opts {
		opt.ConfigureStatus(c)
	}
}

type StatusOption interface {
	ConfigureStatus(*StatusConfig)
}

type StatusFormat string

const (
	// StatusFormatLong returns the long form of status.
	StatusFormatLong StatusFormat = "long"
	// StatusFormatPorcelain returns a consistent
	// status without terminal control codes.
	StatusFormatPorcelain StatusFormat = "porcelain"
	// StatusFormatShort returns the short form of status
	StatusFormatShort StatusFormat = "short"
)

func (f StatusFormat) ToGitValue() string {
	switch f {
	case StatusFormatLong:
		return "--long"
	case StatusFormatPorcelain:
		return "--porcelain"
	case StatusFormatShort:
		return "--short"
	default:
		return "--long"
	}
}

var git = command.NewCommandAlias("git")
