package git

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/mt-sre/go-ci/command"
)

// RevParse runs "git rev-parse" given a RevParseFormat and a variadic
// slice of options. An error is returned if a value cannot be obtained.
func RevParse(ctx context.Context, format RevParseFormat, opts ...RevParseOption) (string, error) {
	var cfg RevParseConfig

	cfg.Option(opts...)

	cmdOpts := []command.CommandOption{
		command.WithContext{Context: ctx},
		command.WithArgs{"rev-parse", format.ToGitValue()},
	}
	if format != RevParseFormatTopLevel {
		cmdOpts = append(cmdOpts, command.WithArgs{"HEAD"})
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

var git = command.NewCommandAlias("git")
