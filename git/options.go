// SPDX-FileCopyrightText: 2025 Red Hat, Inc. <sd-mt-sre@redhat.com>
//
// SPDX-License-Identifier: Apache-2.0

package git

// WithDiffFormat applies the given DiffFormat
type WithDiffFormat DiffFormat

func (w WithDiffFormat) ConfigureDiff(c *DiffConfig) {
	c.Format = DiffFormat(w)
}

// WithRevParseFormat applies the given RevParseFormat
type WithRevParseFormat RevParseFormat

func (w WithRevParseFormat) ConfigureRevParse(c *RevParseConfig) {
	c.Format = RevParseFormat(w)
}

// WithSorted enables sorting.
type WithSorted bool

func (w WithSorted) ConfigureListTags(c *ListTagsConfig) {
	c.Sorted = bool(w)
}

// With SortKey applies the given sort key.
type WithSortKey SortKey

func (w WithSortKey) ConfigureListTags(c *ListTagsConfig) {
	c.SortKey = SortKey(w)
}

// WithStatusFormat applies the given StatusFormat
type WithStatusFormat StatusFormat

func (w WithStatusFormat) ConfigureStatus(c *StatusConfig) {
	c.Format = StatusFormat(w)
}

// WithWorkingDirectory applies the given working directory.
type WithWorkingDirectory string

func (w WithWorkingDirectory) ConfigureRevParse(c *RevParseConfig) {
	c.WorkingDir = string(w)
}

func (w WithWorkingDirectory) ConfigureLatestTag(c *LatestTagConfig) {
	c.WorkingDir = string(w)
}

func (w WithWorkingDirectory) ConfigureLatestVersion(c *LatestVersionConfig) {
	c.WorkingDir = string(w)
}

func (w WithWorkingDirectory) ConfigureListTags(c *ListTagsConfig) {
	c.WorkingDir = string(w)
}

func (w WithWorkingDirectory) ConfigureDiff(c *DiffConfig) {
	c.WorkingDir = string(w)
}

func (w WithWorkingDirectory) ConfigureStatus(c *StatusConfig) {
	c.WorkingDir = string(w)
}
