package file_test

import (
	"testing"

	"github.com/mt-sre/go-ci/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindDefaults(t *testing.T) {
	t.Parallel()

	files, err := file.Find("./testdata")
	require.NoError(t, err)

	assert.ElementsMatch(t, []string{
		"./testdata",
		"testdata/a.txt",
		"testdata/sub",
		"testdata/sub/b.txt",
		"testdata/sub/c.notxt",
	}, files)
}

func TestFind(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		EntType      file.EntType
		Name         string
		ExpectedEnts []string
	}{
		"all entities": {
			EntType: file.EntTypeAll,
			Name:    "*",
			ExpectedEnts: []string{
				"./testdata",
				"testdata/a.txt",
				"testdata/sub",
				"testdata/sub/b.txt",
				"testdata/sub/c.notxt",
			},
		},
		"all files": {
			EntType: file.EntTypeFile,
			Name:    "*",
			ExpectedEnts: []string{
				"testdata/a.txt",
				"testdata/sub/b.txt",
				"testdata/sub/c.notxt",
			},
		},
		"all dirs": {
			EntType: file.EntTypeDir,
			Name:    "*",
			ExpectedEnts: []string{
				"./testdata",
				"testdata/sub",
			},
		},
		"all txt files": {
			EntType: file.EntTypeFile,
			Name:    "*.txt",
			ExpectedEnts: []string{
				"testdata/a.txt",
				"testdata/sub/b.txt",
			},
		},
		"non-matching pattern": {
			EntType:      file.EntTypeFile,
			Name:         "*.dne",
			ExpectedEnts: []string{},
		},
	} {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			files, err := file.Find("./testdata", file.WithEntType(tc.EntType), file.WithName(tc.Name))
			require.NoError(t, err)

			assert.ElementsMatch(t, tc.ExpectedEnts, files)
		})
	}
}
