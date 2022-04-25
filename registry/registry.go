package registry

import (
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	crname "github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/registry"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/remote/transport"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

func StartRegistry(t *testing.T, opts ...RegistryOption) Registry {
	t.Helper()

	var cfg RegistryConfig

	cfg.Option(opts...)
	cfg.Default()

	srv, err := registry.TLS(cfg.Domain)
	if err != nil {
		t.Log("registry could not be started")
		t.FailNow()
	}

	url, err := url.Parse(srv.URL)
	if err != nil {
		t.Log("registry URL could not be parsed")
		t.FailNow()
	}

	return Registry{
		cfg:  cfg,
		host: url.Host,
		srv:  srv,
		t:    t,
	}
}

type RegistryConfig struct {
	Domain string
}

func (c *RegistryConfig) Option(opts ...RegistryOption) {
	for _, opt := range opts {
		opt.ConfigureRegistry(c)
	}
}

func (c *RegistryConfig) Default() {
	if c.Domain == "" {
		c.Domain = "localhost"
	}
}

type Registry struct {
	cfg  RegistryConfig
	t    *testing.T
	host string
	srv  *httptest.Server
}

func (r *Registry) Host() string { return r.host }
func (r *Registry) URL() string  { return r.srv.URL }

func (r *Registry) Images(ctx context.Context) ([]string, error) {
	reg, err := crname.NewRegistry(r.Host())
	if err != nil {
		return nil, fmt.Errorf("parsing registry host: %w", err)
	}

	repos, err := remote.Catalog(context.TODO(), reg, remote.WithTransport(r.srv.Client().Transport))
	if err != nil {
		return nil, fmt.Errorf("retrieving repositories: %w", err)
	}

	var res []string

	for _, repo := range repos {
		repoObj, err := crname.NewRepository(fmt.Sprintf("%s/%s", r.cfg.Domain, repo))
		if err != nil {
			return nil, fmt.Errorf("parsing repository name %q: %w", repo, err)
		}

		tags, err := remote.List(repoObj, remote.WithTransport(r.srv.Client().Transport), remote.WithContext(context.TODO()))
		if err != nil {
			return nil, fmt.Errorf("retrieving tags for repo %q: %w", repo, err)
		}

		for _, tag := range tags {
			res = append(res, fmt.Sprintf("%s/%s:%s", r.cfg.Domain, repo, tag))
		}
	}

	return res, nil
}

func (r *Registry) HasImage(ctx context.Context, image string) (bool, error) {
	n, err := crname.ParseReference(image)
	if err != nil {
		return false, fmt.Errorf("parsing image name: %w", err)
	}

	if _, err := remote.Head(n, remote.WithContext(ctx), remote.WithTransport(r.srv.Client().Transport)); err != nil {
		var tpErr *transport.Error

		if errors.As(err, &tpErr) && tpErr.StatusCode == http.StatusNotFound {
			return false, nil
		}

		return false, fmt.Errorf("checking for image existence: %w", err)
	}

	return true, nil
}

func (r *Registry) Load(ctx context.Context, image string, tarFile string) error {
	img, err := tarball.Image(func() (io.ReadCloser, error) {
		return openTarFile(tarFile)
	}, nil)
	if err != nil {
		return fmt.Errorf("loading image from file: %w", err)
	}

	n, err := crname.ParseReference(image)
	if err != nil {
		return fmt.Errorf("parsing image name: %w", err)
	}

	if err := remote.Write(n, img,
		remote.WithContext(ctx),
		remote.WithTransport(r.srv.Client().Transport),
	); err != nil {
		return fmt.Errorf("pushing image to registry: %w", err)
	}

	return nil
}

func (r *Registry) Stop() {
	r.t.Helper()

	r.srv.Close()
}

func openTarFile(path string) (io.ReadCloser, error) {
	var reader io.ReadCloser

	reader, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening tar file: %w", err)
	}

	unzip, err := gzip.NewReader(reader)
	if err != nil {
		if !errors.Is(err, gzip.ErrHeader) {
			return nil, fmt.Errorf("unzipping file: %w", err)
		}

		return reader, nil
	}

	return unzip, nil
}
