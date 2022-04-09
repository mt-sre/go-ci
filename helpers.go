package regtest

import (
	"bytes"
	"compress/gzip"
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func containerRuntime() string {
	prefferedRuntimes := []string{
		"podman",
		"docker",
	}

	for _, runtime := range prefferedRuntimes {
		if runtimePath, err := exec.LookPath(runtime); err == nil {
			return runtimePath
		}
	}

	return ""
}

func generateCerts(dir string) (string, error) {
	key, err := rsa.GenerateKey(crand.Reader, 4096)
	if err != nil {
		return "", fmt.Errorf("generating signing key: %w", err)
	}

	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(0),
		DNSNames:     []string{"localhost"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(24 * time.Hour),
	}

	cert, err := x509.CreateCertificate(crand.Reader, tmpl, tmpl, key.Public(), key)
	if err != nil {
		return "", fmt.Errorf("generating cert: %w", err)
	}

	pemCert := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	})

	certFile := filepath.Join(dir, "server.crt")

	if err := os.WriteFile(certFile, pemCert, 0700); err != nil {
		return "", fmt.Errorf("writing cert: %w", err)
	}

	pemKey := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	keyFile := filepath.Join(dir, "server.key")

	if err := os.WriteFile(keyFile, pemKey, 0700); err != nil {
		return "", fmt.Errorf("writing private key: %w", err)
	}

	return dir, nil
}

func imageName(image string) string {
	parts := strings.Split(image, "/")

	return parts[len(parts)-1]
}

func parsePort(raw string) (int, error) {
	parts := strings.Split(strings.TrimSpace(raw), ":")

	if len(parts) < 2 {
		return 0, errors.New("unparsable port string")
	}

	return strconv.Atoi(parts[1])
}

func ping(addr string) bool {
	if conn, err := net.Dial("tcp", addr); err == nil {
		defer conn.Close()

		return true
	}

	return false
}

func randomString() (string, error) {
	var buf bytes.Buffer

	data := make([]byte, 32)

	if _, err := rand.Read(data); err != nil {
		return "", fmt.Errorf("reading from random source: %w", err)
	}

	if _, err := base64.NewEncoder(base64.RawURLEncoding, &buf).Write(data); err != nil {
		return "", fmt.Errorf("encoding to base64: %w", err)
	}

	return buf.String(), nil
}

func readTar(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("opening tarball: %w", err)
	}

	if filepath.Ext(path) != ".gz" {
		return data, nil
	}

	unzip, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("unzipping tarball: %w", err)
	}

	expanded, err := ioutil.ReadAll(unzip)
	if err != nil {
		return nil, fmt.Errorf("reading unzipped tarball: %w", err)
	}

	return expanded, nil
}

func runWithOutput(cmd *exec.Cmd) (string, error) {
	var out bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()

	return out.String(), err
}
