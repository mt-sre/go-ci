package git

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRevParse(t *testing.T) {
	t.Parallel()

	dir, err := os.Getwd()
	require.NoError(t, err)

	for name, tc := range map[string]struct {
		Format   RevParseFormat
		Expected string
	}{
		"top-level": {
			Format:   RevParseFormatTopLevel,
			Expected: filepath.Join(dir, "testdata"),
		},
		"abbrev-ref": {
			Format:   RevParseFormatAbbrevRef,
			Expected: "test",
		},
		"short": {
			Format:   RevParseFormatShort,
			Expected: "7e892e7",
		},
	} {
		tc := tc

		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			res, err := RevParse(ctx, tc.Format, WithWorkingDirectory("testdata"))
			require.NoError(t, err)

			assert.Equal(t, tc.Expected, res)
		})
	}
}

func TestListTags(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	res, err := ListTags(ctx, WithWorkingDirectory("testdata"))
	require.NoError(t, err)

	assert.ElementsMatch(t, []string{"v1.0.0", "v2.0.0"}, res)
}

func TestLatestTag(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	res, err := LatestTag(ctx, WithWorkingDirectory("testdata"))
	require.NoError(t, err)

	assert.Equal(t, "v2.0.0", res)
}
