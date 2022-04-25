package registry

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistryLoad(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		Image     string
		ImageFile string
	}{
		"valid image": {
			Image:     "localhost.localdomain/echo:v1.0.0",
			ImageFile: "./testimages/echo-v1.0.0.tar.gz",
		},
	} {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			reg := StartRegistry(t, WithDomain("localhost.localdomain"))

			defer reg.Stop()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			require.NoError(t, reg.Load(ctx, tc.Image, tc.ImageFile))

			hasImage, err := reg.HasImage(ctx, tc.Image)
			require.NoError(t, err)

			assert.True(t, hasImage)

			images, err := reg.Images(ctx)
			require.NoError(t, err)

			assert.Contains(t, images, tc.Image)
		})
	}
}
