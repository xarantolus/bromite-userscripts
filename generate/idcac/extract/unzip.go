package extract

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/krolaw/zipstream"
)

// Zip extract the streamed ZIP file in r to the destination directory
func Zip(r io.Reader, destinationDirectory string) (err error) {
	destPathPrefix := filepath.Clean(destinationDirectory) + string(os.PathSeparator)
	zr := zipstream.NewReader(r)
	if err != nil {
		return err
	}

	var fh *zip.FileHeader
	for {
		fh, err = zr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		fpath := filepath.Join(destinationDirectory, fh.Name)

		// Check for ZipSlip. More Info: https://snyk.io/research/zip-slip-vulnerability#go
		if !strings.HasPrefix(fpath, destPathPrefix) {
			err = fmt.Errorf("%s: illegal file path", fpath)
			return
		}

		if fh.FileInfo().IsDir() {
			// Make Folder
			err = os.MkdirAll(fpath, os.ModePerm)
			if err != nil {
				return
			}
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fh.Mode())
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, zr)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()

		if err != nil {
			return err
		}
	}
}
