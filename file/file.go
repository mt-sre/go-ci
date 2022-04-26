package file

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

type EntType string

const (
	EntTypeNone EntType = ""
	EntTypeAll  EntType = "all"
	EntTypeDir  EntType = "dir"
	EntTypeFile EntType = "file"
)

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

type WithEntType EntType

func (t WithEntType) ConfigureFind(c *findConfig) {
	c.EntType = EntType(t)
}

type WithName string

func (n WithName) ConfigureFind(c *findConfig) {
	c.Name = string(n)
}
