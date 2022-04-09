package regtest

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/multierr"
)

func StartRegistry(t *testing.T, opts ...RegistryOption) Registry {
	t.Helper()

	reg := Registry{
		t: t,
	}

	reg.cfg.Option(opts...)

	require.NoError(t, reg.cfg.Default(), "unable to set defaults")

	require.NotEmpty(t, reg.cfg.Runtime, "no container runtime available")

	reg.run()

	return reg
}

type RegistryConfig struct {
	Port      int
	EnableTLS bool
	Name      string
	Image     string
	Runtime   string
}

func (c *RegistryConfig) Option(opts ...RegistryOption) {
	for _, opt := range opts {
		opt.ConfigureRegistry(c)
	}
}

func (c *RegistryConfig) Default(opts ...RegistryOption) error {
	if c.Image == "" {
		c.Image = "registry:2"
	}

	if c.Name == "" {
		suffix, err := randomString()
		if err != nil {
			return fmt.Errorf("generating random suffix: %w", err)
		}

		c.Name = fmt.Sprintf("registry-%s", suffix)
	}

	if c.Runtime == "" {
		c.Runtime = containerRuntime()
	}

	return nil
}

type Registry struct {
	actualPort int
	cfg        RegistryConfig
	images     []string
	t          *testing.T
}

func (r *Registry) run() {
	r.t.Helper()

	require.NoError(r.t, r.pull())

	if r.cfg.EnableTLS {
		require.NoError(r.t, r.startTLS())
	} else {
		require.NoError(r.t, r.start())
	}

	var out bytes.Buffer

	getPort := exec.Command(r.cfg.Runtime, "port", r.cfg.Name)
	getPort.Stdout = &out

	if !assert.NoError(r.t, getPort.Run()) {
		defer r.Stop()

		assert.FailNow(r.t, "retrieving actual port")
	}

	port, err := parsePort(out.String())
	if !assert.NoError(r.t, err) {
		defer r.Stop()

		assert.FailNow(r.t, "parsing actual port")
	}

	r.actualPort = port

	require.Eventually(r.t, r.ping, 10*time.Second, 250*time.Millisecond)
}

func (r *Registry) pull() error {
	pullImage := exec.Command(r.cfg.Runtime, "pull", r.cfg.Image)
	if out, err := runWithOutput(pullImage); err != nil {
		return fmt.Errorf("pulling registry image: %s: %w", out, err)
	}

	return nil
}

func (r *Registry) start() error {
	port := "5000"

	if r.cfg.Port != 0 {
		port = fmt.Sprintf("%d:5000", r.cfg.Port)
	}

	run := exec.Command(r.cfg.Runtime, "run", "--rm", "-d",
		"-p", port,
		"--name", r.cfg.Name,
		r.cfg.Image,
	)

	if out, err := runWithOutput(run); err != nil {
		return fmt.Errorf("starting registry: %s: %w", out, err)
	}

	return nil
}

func (r *Registry) ping() bool {
	return ping(fmt.Sprintf("localhost:%d", r.actualPort))
}

func (r *Registry) startTLS() error {
	certDir, err := generateCerts(r.t.TempDir())
	if err != nil {
		return fmt.Errorf("generating CA: %w", err)
	}

	port := "443"

	if r.cfg.Port != 0 {
		port = fmt.Sprintf("%d:443", r.cfg.Port)
	}

	run := exec.Command(r.cfg.Runtime, "run", "--rm", "-d",
		"-p", port,
		"--name", r.cfg.Name,
		"-v", fmt.Sprintf("%s:/certs", certDir),
		"-e", "REGISTRY_HTTP_ADDR=0.0.0.0:443",
		"-e", "REGISTRY_HTTP_TLS_CERTIFICATE=/certs/server.crt",
		"-e", "REGISTRY_HTTP_TLS_KEY=/certs/server.key",
		r.cfg.Image,
	)

	if out, err := runWithOutput(run); err != nil {
		return fmt.Errorf("starting registry: %s: %w", out, err)
	}

	return nil
}

func (r *Registry) Stop() error {
	r.t.Helper()

	var collector error

	stop := exec.Command(r.cfg.Runtime, "stop", r.cfg.Name)

	if out, err := runWithOutput(stop); err != nil {
		multierr.AppendInto(&collector, fmt.Errorf("stopping registry: %s: %w", out, err))
	}

	for _, image := range r.images {
		remove := exec.Command(r.cfg.Runtime, "image", "rm", image)

		if out, err := runWithOutput(remove); err != nil {
			multierr.AppendInto(&collector, fmt.Errorf("removing loaded image %q: %s: %w", image, out, err))
		}
	}

	return collector
}

func (r *Registry) Load(image string, tarFile string) error {
	r.t.Helper()

	data, err := readTar(tarFile)
	if err != nil {
		return fmt.Errorf("reading tar file %q: %w", tarFile, err)
	}

	load := exec.Command(r.cfg.Runtime, "load")
	load.Stdin = bytes.NewBuffer(data)

	if out, err := runWithOutput(load); err != nil {
		return fmt.Errorf("loading tarball %q: %s: %w", tarFile, out, err)
	}

	taggedImage := fmt.Sprintf("%s/%s", r.Host(), imageName(image))

	tag := exec.Command(r.cfg.Runtime, "tag", image, taggedImage)

	if out, err := runWithOutput(tag); err != nil {
		return fmt.Errorf("tagging image %q: %s: %w", image, out, err)
	}

	r.images = append(r.images, taggedImage)

	push := exec.Command(r.cfg.Runtime, "push", "--tls-verify=false", taggedImage)

	if out, err := runWithOutput(push); err != nil {
		return fmt.Errorf("pushing image %q: %s: %w", taggedImage, out, err)
	}

	return nil
}

func (r *Registry) HasImage(ctx context.Context, name, tag string) (bool, error) {
	uri := r.URL()
	uri.Path = fmt.Sprintf("v2/%s/manifests/%s", name, tag)

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, uri.String(), nil)
	if err != nil {
		return false, fmt.Errorf("constructing request: %w", err)
	}

	tp := http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := http.Client{
		Transport: &tp,
	}

	res, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("sending request: %w", err)
	}

	defer res.Body.Close()

	return res.StatusCode == http.StatusOK, nil
}

func (r *Registry) URL() url.URL {
	res := url.URL{
		Host:   r.Host(),
		Scheme: "http",
	}

	if r.cfg.EnableTLS {
		res.Scheme = "https"
	}

	return res
}

func (r *Registry) Host() string {
	return fmt.Sprintf("localhost:%d", r.actualPort)
}
