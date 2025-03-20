// SPDX-FileCopyrightText: 2025 Red Hat, Inc. <sd-mt-sre@redhat.com>
//
// SPDX-License-Identifier: Apache-2.0

package file

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

// EntType describes the entity types which a filesystem may contain.
type EntType string

const (
	// EntTypeNone matches no entities.
	EntTypeNone EntType = ""
	// EntTypeAll matches all entity types.
	EntTypeAll EntType = "all"
	// EntTypeDir matches directory entities.
	EntTypeDir EntType = "dir"
	// EntTypeFile matches file entities.
	EntTypeFile EntType = "file"
)

// Find functions similarly to GNU find searching recursively
// from root for all entites which match the given options.
// By default all entity types (file, directory) will be returned
// including the root directory.
func Find(root string, opts ...FindOption) ([]string, error) {
	var cfg findConfig

	cfg.Option(opts...)
	cfg.Default()

	var result []string

	hasCorrectType := func(fs.DirEntry) bool { return true }

	switch cfg.EntType {
	case EntTypeDir:
		hasCorrectType = func(d fs.DirEntry) bool {
			return d.IsDir()
		}
	case EntTypeFile:
		hasCorrectType = func(d fs.DirEntry) bool {
			return !d.IsDir()
		}
	}

	searchFunc := func(path string, d fs.DirEntry, _ error) error {
		if !hasCorrectType(d) {
			return nil
		}

		matches, err := filepath.Match(cfg.Name, filepath.Base(path))
		if err != nil {
			return fmt.Errorf("matching %q against %q: %w", path, cfg.Name, err)
		}

		if matches {
			result = append(result, path)
		}

		return nil
	}

	if err := filepath.WalkDir(root, searchFunc); err != nil {
		return nil, fmt.Errorf("walking directories: %w", err)
	}

	return result, nil
}

type findConfig struct {
	EntType EntType
	Name    string
}

func (c *findConfig) Option(opts ...FindOption) {
	for _, opt := range opts {
		opt.ConfigureFind(c)
	}
}

func (c *findConfig) Default() {
	if c.EntType == EntTypeNone {
		c.EntType = EntTypeAll
	}

	if c.Name == "" {
		c.Name = "*"
	}
}

type FindOption interface {
	ConfigureFind(*findConfig)
}

// WithEntType selects which entity types will be
// matched by the find function.
type WithEntType EntType

func (t WithEntType) ConfigureFind(c *findConfig) {
	c.EntType = EntType(t)
}

// WithName supplies a glob pattern to filter the
// matched entities to only those which match
// the glob.
type WithName string

func (n WithName) ConfigureFind(c *findConfig) {
	c.Name = string(n)
}
