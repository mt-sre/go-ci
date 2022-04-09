package regtest

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistryLoad(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		Image      string
		ImageFile  string
		ShouldLoad bool
		Tag        string
		WithTLS    bool
	}{
		"valid image/no tls": {
			Image:      "localhost/echo",
			Tag:        "v1.0.0",
			ImageFile:  "./testimages/echo-v1.0.0.tar.gz",
			ShouldLoad: true,
			WithTLS:    false,
		},
		"valid image/tls": {
			Image:      "localhost/echo",
			Tag:        "v1.0.0",
			ImageFile:  "./testimages/echo-v1.0.0.tar.gz",
			ShouldLoad: true,
			WithTLS:    true,
		},
		"invalid image": {
			Image:      "",
			Tag:        "",
			ImageFile:  "./testimages/echo-v1.0.0.tar.gz",
			ShouldLoad: false,
			WithTLS:    false,
		},
		"invalid tarfile": {
			Image:      "localhost/echo",
			Tag:        "v1.0.0",
			ImageFile:  "dne",
			ShouldLoad: false,
			WithTLS:    false,
		},
	} {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			reg := StartRegistry(t, WithTLSEnabled(tc.WithTLS))

			defer reg.Stop()

			if !tc.ShouldLoad {
				require.Error(t, reg.Load(fmt.Sprintf("%s:%s", tc.Image, tc.Tag), tc.ImageFile))

				return
			}

			require.NoError(t, reg.Load(fmt.Sprintf("%s:%s", tc.Image, tc.Tag), tc.ImageFile))

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			hasImage, err := reg.HasImage(ctx, imageName(tc.Image), tc.Tag)
			require.NoError(t, err)

			assert.True(t, hasImage)
		})
	}
}
