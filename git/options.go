package git

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

// WithWorkingDirectory applies the given working directory.
type WithWorkingDirectory string

func (w WithWorkingDirectory) ConfigureRevParse(c *RevParseConfig) {
	c.WorkingDir = string(w)
}

func (w WithWorkingDirectory) ConfigureLatestTag(c *LatestTagConfig) {
	c.WorkingDir = string(w)
}

func (w WithWorkingDirectory) ConfigureListTags(c *ListTagsConfig) {
	c.WorkingDir = string(w)
}
