package plugin

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/indeedhat/parity-nas/internal/config"
)

func (m PluginManager) downloadArchive(entry config.PluginEntry) error {
	fh, err := os.Create(entry.ArchiveSavePath(m.cfg))
	if err != nil {
		return err
	}
	defer fh.Close()

	resp, err := http.Get(entry.ArchiveUrl())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(fh, resp.Body)
	return err
}

func (m PluginManager) extractArchive(entry config.PluginEntry) error {
	zr, err := zip.OpenReader(entry.ArchiveSavePath(m.cfg))
	if err != nil {
		return err
	}

	dst := entry.ArchiveExtractPath(m.cfg)
	if err := os.MkdirAll(dst, os.ModePerm); err != nil {
		return err
	}

	for _, f := range zr.File {
		if err := extractFile(f, dst); err != nil {
			return err
		}
	}

	return nil
}

func (m PluginManager) cleanupArchive(entry config.PluginEntry) error {
	if err := os.Remove(entry.ArchiveSavePath(m.cfg)); err != nil {
		return err
	}

	return os.RemoveAll(entry.ArchiveExtractPath(m.cfg))
}

func extractFile(f *zip.File, dst string) error {
	filePath := filepath.Join(dst, f.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", filePath)
	}

	if f.FileInfo().IsDir() {
		return os.MkdirAll(filePath, os.ModePerm)
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	fh, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer fh.Close()

	zf, err := f.Open()
	if err != nil {
		return err
	}
	defer zf.Close()

	_, err = io.Copy(fh, zf)
	return err
}
