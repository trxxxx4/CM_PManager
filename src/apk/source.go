package apk

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	// IndexArchiveName is the name of the index file of a repository.
	IndexArchiveName = "APKINDEX.tar.gz"
	archiveSuffix    = ".tar.gz"
	indexEntryName   = "APKINDEX"
)

var (
	gzipMagic       = []byte{0x1f, 0x8b}
	errIndexMissing = errors.New("archive has no APKINDEX file")
)

// LoadRemote downloads and parses the index of the repository.
//
// The URL is either the address of the repository directory, for
// example https://host/alpine/v3.19/main/x86_64, or of the index itself.
func LoadRemote(repoURL string, client *http.Client) ([]Package, error) {
	target := indexURL(repoURL)

	response, err := client.Get(target)
	if err != nil {
		return nil, fmt.Errorf("cannot download %s: %w", target, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"cannot download %s: server answered %s",
			target,
			response.Status,
		)
	}

	packages, err := ReadIndex(response.Body)
	if err != nil {
		return nil, fmt.Errorf("invalid index %s: %w", target, err)
	}

	return packages, nil
}

// LoadFile reads and parses a local index: APKINDEX.tar.gz or text.
func LoadFile(path string) ([]Package, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open %s: %w", path, err)
	}
	defer file.Close()

	packages, err := ReadIndex(file)
	if err != nil {
		return nil, fmt.Errorf("invalid index %s: %w", path, err)
	}

	return packages, nil
}

// ReadIndex parses an index given as a gzip archive or as plain text.
func ReadIndex(r io.Reader) ([]Package, error) {
	buffered := bufio.NewReader(r)

	magic, err := buffered.Peek(len(gzipMagic))
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}

	if bytes.Equal(magic, gzipMagic) {
		return readArchive(buffered)
	}

	return ParseIndex(buffered)
}

// indexURL returns the address of the index file of the repository.
func indexURL(repoURL string) string {
	if strings.HasSuffix(repoURL, archiveSuffix) {
		return repoURL
	}

	return strings.TrimRight(repoURL, "/") + "/" + IndexArchiveName
}

// readArchive finds the APKINDEX file in the archive and parses it.
//
// The archive is a signature segment and an index segment glued
// together as gzip streams, which the gzip reader joins by default.
func readArchive(r io.Reader) ([]Package, error) {
	unzipped, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("invalid gzip data: %w", err)
	}
	defer unzipped.Close()

	archive := tar.NewReader(unzipped)

	for {
		header, err := archive.Next()
		if errors.Is(err, io.EOF) {
			return nil, errIndexMissing
		}

		if err != nil {
			return nil, fmt.Errorf("invalid tar data: %w", err)
		}

		if strings.TrimPrefix(header.Name, "./") == indexEntryName {
			return ParseIndex(archive)
		}
	}
}
