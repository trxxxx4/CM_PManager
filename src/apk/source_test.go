package apk

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const repoPath = "/alpine/main/x86_64/"

// gzipBytes compresses data as one gzip stream.
func gzipBytes(t *testing.T, data []byte) []byte {
	t.Helper()

	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)

	if _, err := writer.Write(data); err != nil {
		t.Fatal(err)
	}

	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	return buf.Bytes()
}

// tarSegment builds a tar with one file, with or without the end marker.
func tarSegment(t *testing.T, name, content string, final bool) []byte {
	t.Helper()

	var buf bytes.Buffer
	writer := tar.NewWriter(&buf)
	header := &tar.Header{Name: name, Mode: 0o644, Size: int64(len(content))}

	if err := writer.WriteHeader(header); err != nil {
		t.Fatal(err)
	}

	if _, err := writer.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}

	finish := writer.Flush
	if final {
		finish = writer.Close
	}

	if err := finish(); err != nil {
		t.Fatal(err)
	}

	return buf.Bytes()
}

// buildArchive makes an APKINDEX.tar.gz with a signature segment.
func buildArchive(t *testing.T, index string) []byte {
	t.Helper()

	signature := tarSegment(t, ".SIGN.RSA.test.rsa.pub", "signature", false)
	body := tarSegment(t, indexEntryName, index, true)

	return append(gzipBytes(t, signature), gzipBytes(t, body)...)
}

// writeTemp stores data in a file of a temporary directory.
func writeTemp(t *testing.T, name string, data []byte) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}

	return path
}

// newServer serves body as the index of the test repository.
func newServer(t *testing.T, body []byte) *httptest.Server {
	t.Helper()

	handler := func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(body)
	}

	mux := http.NewServeMux()
	mux.HandleFunc(repoPath+IndexArchiveName, handler)

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	return server
}

// TestReadIndexFormats checks both archive and plain text indexes.
func TestReadIndexFormats(t *testing.T) {
	inputs := map[string][]byte{
		"archive": buildArchive(t, sampleIndex),
		"text":    []byte(sampleIndex),
	}

	for name, data := range inputs {
		packages, err := ReadIndex(bytes.NewReader(data))

		if err != nil || len(packages) != sampleCount {
			t.Errorf("%s: got %d packages, error %v", name, len(packages), err)
		}
	}
}

// TestReadIndexErrors checks broken archives.
func TestReadIndexErrors(t *testing.T) {
	noIndex := gzipBytes(t, tarSegment(t, "DESCRIPTION", "text", true))

	_, err := ReadIndex(bytes.NewReader(noIndex))
	if !errors.Is(err, errIndexMissing) {
		t.Errorf("got %v, want errIndexMissing", err)
	}

	_, err = ReadIndex(bytes.NewReader([]byte{0x1f, 0x8b, 0x00, 0x01}))
	if err == nil {
		t.Error("expected an error for broken gzip data")
	}
}

// TestLoadFile checks local files and the missing file error.
func TestLoadFile(t *testing.T) {
	path := writeTemp(t, "APKINDEX.tar.gz", buildArchive(t, sampleIndex))

	packages, err := LoadFile(path)
	if err != nil || len(packages) != sampleCount {
		t.Errorf("got %d packages, error %v", len(packages), err)
	}

	_, err = LoadFile(filepath.Join(t.TempDir(), "absent"))
	if err == nil || !strings.Contains(err.Error(), "cannot open") {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestLoadRemote checks the supported forms of the repository URL.
func TestLoadRemote(t *testing.T) {
	server := newServer(t, buildArchive(t, sampleIndex))
	base := server.URL + strings.TrimRight(repoPath, "/")
	urls := []string{base, base + "/", base + "/" + IndexArchiveName}

	for _, address := range urls {
		packages, err := LoadRemote(address, server.Client())

		if err != nil || len(packages) != sampleCount {
			t.Errorf("%s: got %d packages, error %v",
				address, len(packages), err)
		}
	}
}

// TestLoadRemoteErrors checks HTTP and network failures.
func TestLoadRemoteErrors(t *testing.T) {
	server := newServer(t, nil)

	_, err := LoadRemote(server.URL+"/other", server.Client())
	if err == nil || !strings.Contains(err.Error(), "404") {
		t.Errorf("got %v, want a 404 error", err)
	}

	address := server.URL
	server.Close()

	_, err = LoadRemote(address, http.DefaultClient)
	if err == nil || !strings.Contains(err.Error(), "cannot download") {
		t.Errorf("got %v, want a download error", err)
	}
}
