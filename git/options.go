package git

type WithSorted bool

func (w WithSorted) ConfigureListTags(c *ListTagsConfig) {
	c.Sorted = bool(w)
}

type WithSortKey SortKey

func (w WithSortKey) ConfigureListTags(c *ListTagsConfig) {
	c.SortKey = SortKey(w)
}

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
