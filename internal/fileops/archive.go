// Package fileops реализует операции с файлами и архивами для файлового менеджера.
package fileops

import (
	"archive/tar"
	"archive/zip"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"errors"
	"file-manager/internal/i18n"

	"github.com/ulikunitz/xz"
)

// Archiver предоставляет функции для работы с архивами
type Archiver struct{}

// NewArchiver создает новый экземпляр Archiver
func NewArchiver() *Archiver {
	return &Archiver{}
}

// ArchiveFiles создает архив из указанных файлов и директорий (zip, tar.gz, tar.bz2, tar.xz)
func (a *Archiver) ArchiveFiles(sources []string, destination string, format string) error {
	if format == "" {
		format = filepath.Ext(destination)
		if format != "" {
			format = format[1:]
		}
	}
	format = strings.ToLower(format)
	switch format {
	case "zip":
		return a.archiveZip(sources, destination)
	case "tar.gz", "tgz":
		return a.archiveTarCompressed(sources, destination, "gz")
	case "tar.bz2", "tbz2":
		return a.archiveTarCompressed(sources, destination, "bz2")
	case "tar.xz", "txz":
		return a.archiveTarCompressed(sources, destination, "xz")
	case "tar":
		return a.archiveTar(sources, destination)
	default:
		return errors.New(i18n.T("archive_format_error"))
	}
}

func (a *Archiver) archiveZip(sources []string, destination string) error {
	zipFile, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf(i18n.T("archive_create_error"), err)
	}
	defer func() {
		err := zipFile.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, i18n.T("archive_close_zip_error")+"\n", err)
		}
	}()
	zipWriter := zip.NewWriter(zipFile)
	defer func() {
		err := zipWriter.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, i18n.T("archive_close_zipwriter_error")+"\n", err)
		}
	}()
	for _, src := range sources {
		err := addFileToZip(zipWriter, src, "")
		if err != nil {
			return fmt.Errorf(i18n.T("archive_add_error"), src, err)
		}
	}
	return nil
}

func (a *Archiver) archiveTarCompressed(sources []string, destination, compression string) error {
	file, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf(i18n.T("archive_create_error"), err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Fprintf(os.Stderr, i18n.T("archive_close_file_error")+"\n", err)
		}
	}()
	var tw *tar.Writer
	var writer io.WriteCloser
	switch compression {
	case "gz":
		gw := gzip.NewWriter(file)
		defer func() {
			if err := gw.Close(); err != nil {
				fmt.Fprintf(os.Stderr, i18n.T("archive_close_gzip_error")+"\n", err)
			}
		}()
		writer = gw
	case "bz2":
		return errors.New(i18n.T("archive_bz2_unsupported"))
	case "xz":
		xzw, err := xz.NewWriter(file)
		if err != nil {
			return fmt.Errorf(i18n.T("archive_create_xz_error"), err)
		}
		defer func() {
			if err := xzw.Close(); err != nil {
				fmt.Fprintf(os.Stderr, i18n.T("archive_close_xz_error")+"\n", err)
			}
		}()
		writer = xzw
	default:
		return fmt.Errorf(i18n.T("archive_unknown_compression"), compression)
	}
	tw = tar.NewWriter(writer)
	defer func() {
		if err := tw.Close(); err != nil {
			fmt.Fprintf(os.Stderr, i18n.T("archive_close_tar_error")+"\n", err)
		}
	}()
	for _, src := range sources {
		err := addFileToTar(tw, src, "")
		if err != nil {
			return fmt.Errorf(i18n.T("archive_add_error"), src, err)
		}
	}
	return nil
}

func (a *Archiver) archiveTar(sources []string, destination string) error {
	file, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf(i18n.T("archive_create_error"), err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Fprintf(os.Stderr, i18n.T("archive_close_file_error")+"\n", err)
		}
	}()
	tw := tar.NewWriter(file)
	defer func() {
		if err := tw.Close(); err != nil {
			fmt.Fprintf(os.Stderr, i18n.T("archive_close_tar_error")+"\n", err)
		}
	}()
	for _, src := range sources {
		err := addFileToTar(tw, src, "")
		if err != nil {
			return fmt.Errorf(i18n.T("archive_add_error"), src, err)
		}
	}
	return nil
}

func addFileToTar(tw *tar.Writer, src, baseInTar string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if info.IsDir() {
		entries, err := os.ReadDir(src)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			entryPath := filepath.Join(src, entry.Name())
			var entryBase string
			if baseInTar == "" {
				entryBase = entry.Name()
			} else {
				entryBase = filepath.Join(baseInTar, entry.Name())
			}
			err = addFileToTar(tw, entryPath, entryBase)
			if err != nil {
				return err
			}
		}
		return nil
	}
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Fprintf(os.Stderr, i18n.T("archive_close_file_error")+"\n", err)
		}
	}()
	hdr, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return err
	}
	hdr.Name = filepath.ToSlash(baseInTar)
	err = tw.WriteHeader(hdr)
	if err != nil {
		return err
	}
	_, err = io.Copy(tw, file)
	return err
}

// ExtractArchive поддерживает zip, tar.gz, tar.bz2, tar.xz
func (a *Archiver) ExtractArchive(source, destination string) error {
	format := strings.ToLower(filepath.Ext(source))
	if strings.HasSuffix(source, ".tar.gz") || strings.HasSuffix(source, ".tgz") {
		return extractTarCompressed(source, destination, "gz")
	} else if strings.HasSuffix(source, ".tar.bz2") || strings.HasSuffix(source, ".tbz2") {
		return extractTarCompressed(source, destination, "bz2")
	} else if strings.HasSuffix(source, ".tar.xz") || strings.HasSuffix(source, ".txz") {
		return extractTarCompressed(source, destination, "xz")
	} else if format == ".tar" {
		return extractTarCompressed(source, destination, "none")
	} else if format == ".zip" {
		return a.ExtractZip(source, destination)
	}
	return errors.New(i18n.T("archive_format_error"))
}

// ExtractZip извлекает zip-архив в указанную директорию.
func (a *Archiver) ExtractZip(source, destination string) error {
	zipReader, err := zip.OpenReader(source)
	if err != nil {
		return fmt.Errorf(i18n.T("archive_open_error"), err)
	}
	defer func() {
		if err := zipReader.Close(); err != nil {
			fmt.Fprintf(os.Stderr, i18n.T("archive_close_zipreader_error")+"\n", err)
		}
	}()
	for _, f := range zipReader.File {
		if strings.Contains(f.Name, "..") || filepath.IsAbs(f.Name) {
			return fmt.Errorf(i18n.T("archive_unsafe_path_error"), f.Name)
		}
		fpath := filepath.Join(destination, f.Name)
		if !strings.HasPrefix(filepath.Clean(fpath)+string(os.PathSeparator), filepath.Clean(destination)+string(os.PathSeparator)) {
			return fmt.Errorf(i18n.T("archive_path_traversal_error"), fpath)
		}
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, f.Mode()); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return err
		}
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			errClose := outFile.Close()
			if errClose != nil {
				fmt.Fprintf(os.Stderr, i18n.T("archive_close_outfile_error")+"\n", errClose)
			}
			return err
		}
		_, err = io.Copy(outFile, rc)
		errClose1 := outFile.Close()
		errClose2 := rc.Close()
		if errClose1 != nil {
			fmt.Fprintf(os.Stderr, i18n.T("archive_close_outfile_error")+"\n", errClose1)
		}
		if errClose2 != nil {
			fmt.Fprintf(os.Stderr, i18n.T("archive_close_rc_error")+"\n", errClose2)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func extractTarCompressed(source, destination, compression string) error {
	file, err := os.Open(source)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Fprintf(os.Stderr, i18n.T("archive_close_file_error")+"\n", err)
		}
	}()
	var tr *tar.Reader
	switch compression {
	case "gz":
		gr, err := gzip.NewReader(file)
		if err != nil {
			return err
		}
		defer func() {
			if err := gr.Close(); err != nil {
				fmt.Fprintf(os.Stderr, i18n.T("archive_close_gzipreader_error")+"\n", err)
			}
		}()
		tr = tar.NewReader(gr)
	case "bz2":
		br := bzip2.NewReader(file)
		tr = tar.NewReader(br)
	case "xz":
		xzr, err := xz.NewReader(file)
		if err != nil {
			return err
		}
		tr = tar.NewReader(xzr)
	case "none":
		tr = tar.NewReader(file)
	default:
		return fmt.Errorf(i18n.T("archive_unknown_compression"), compression)
	}
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fpath := filepath.Join(destination, hdr.Name)
		if !strings.HasPrefix(filepath.Clean(fpath)+string(os.PathSeparator), filepath.Clean(destination)+string(os.PathSeparator)) {
			return fmt.Errorf(i18n.T("archive_path_traversal_error"), fpath)
		}
		if hdr.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, hdr.FileInfo().Mode()); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return err
		}
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, hdr.FileInfo().Mode())
		if err != nil {
			return err
		}
		_, err = io.Copy(outFile, tr)
		errClose := outFile.Close()
		if errClose != nil {
			fmt.Fprintf(os.Stderr, i18n.T("archive_close_outfile_error")+"\n", errClose)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// ListArchiveContents поддерживает zip, tar.gz, tar.bz2, tar.xz
func (a *Archiver) ListArchiveContents(source string) ([]string, error) {
	format := strings.ToLower(filepath.Ext(source))
	if strings.HasSuffix(source, ".tar.gz") || strings.HasSuffix(source, ".tgz") {
		return listTarCompressed(source, "gz")
	} else if strings.HasSuffix(source, ".tar.bz2") || strings.HasSuffix(source, ".tbz2") {
		return listTarCompressed(source, "bz2")
	} else if strings.HasSuffix(source, ".tar.xz") || strings.HasSuffix(source, ".txz") {
		return listTarCompressed(source, "xz")
	} else if format == ".tar" {
		return listTarCompressed(source, "none")
	} else if format == ".zip" {
		return a.listZip(source)
	}
	return nil, errors.New(i18n.T("archive_format_error"))
}

func (a *Archiver) listZip(source string) ([]string, error) {
	zipReader, err := zip.OpenReader(source)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("archive_open_error"), err)
	}
	defer func() {
		if err := zipReader.Close(); err != nil {
			fmt.Fprintf(os.Stderr, i18n.T("archive_close_zipreader_error")+"\n", err)
		}
	}()
	var files []string
	for _, f := range zipReader.File {
		files = append(files, f.Name)
	}
	return files, nil
}

func listTarCompressed(source, compression string) ([]string, error) {
	file, err := os.Open(source)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Fprintf(os.Stderr, i18n.T("archive_close_file_error")+"\n", err)
		}
	}()
	var tr *tar.Reader
	switch compression {
	case "gz":
		gr, err := gzip.NewReader(file)
		if err != nil {
			return nil, err
		}
		defer func() {
			if err := gr.Close(); err != nil {
				fmt.Fprintf(os.Stderr, i18n.T("archive_close_gzipreader_error")+"\n", err)
			}
		}()
		tr = tar.NewReader(gr)
	case "bz2":
		br := bzip2.NewReader(file)
		tr = tar.NewReader(br)
	case "xz":
		xzr, err := xz.NewReader(file)
		if err != nil {
			return nil, err
		}
		tr = tar.NewReader(xzr)
	case "none":
		tr = tar.NewReader(file)
	default:
		return nil, fmt.Errorf(i18n.T("archive_unknown_compression"), compression)
	}
	var files []string
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		files = append(files, hdr.Name)
	}
	return files, nil
}

func addFileToZip(zipWriter *zip.Writer, src, baseInZip string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if info.IsDir() {
		entries, err := os.ReadDir(src)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			entryPath := filepath.Join(src, entry.Name())
			var entryBase string
			if baseInZip == "" {
				entryBase = entry.Name()
			} else {
				entryBase = filepath.Join(baseInZip, entry.Name())
			}
			err = addFileToZip(zipWriter, entryPath, entryBase)
			if err != nil {
				return err
			}
		}
		return nil
	}
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, i18n.T("archive_close_file_error")+"\n", err)
		}
	}()
	zipHeader, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	// baseInZip всегда относительный путь без ведущих слэшей
	nameInZip := baseInZip
	if nameInZip == "" {
		nameInZip = filepath.Base(src)
	}
	zipHeader.Name = filepath.ToSlash(nameInZip)
	zipHeader.Method = zip.Deflate
	writer, err := zipWriter.CreateHeader(zipHeader)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, file)
	return err
}
