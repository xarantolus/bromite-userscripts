package util

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var httpClient = http.Client{
	Timeout: 30 * time.Second,
}

var (
	downloadLocks     = make(map[string]*sync.Mutex)
	downloadLocksLock sync.Mutex
)

func DownloadURLs(inputURLs []string, tempDir string) (outputPaths []string, err error) {
	var dlFile = func(url string, file string) (err error) {
		f, err := os.Create(file)
		if err != nil {
			return
		}
		defer f.Close()

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return
		}

		req.Header.Set("User-Agent", "github.com/xarantolus/bromite-userscripts")
		req.Header.Set("Accept-Language", "en-US,en;q=0.5")

		resp, err := httpClient.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 400 {
			return fmt.Errorf("unexpected status code %d", resp.StatusCode)
		}

		_, err = io.Copy(f, resp.Body)

		return
	}

	var errCount int

	for _, dlURL := range inputURLs {
		fn := filepath.Join(tempDir, generateFilename(dlURL))
		// Prevent multiple downloads of the same file at the same time
		downloadLocksLock.Lock()
		lock, ok := downloadLocks[fn]
		if !ok {
			lock = &sync.Mutex{}
			downloadLocks[fn] = lock
		}
		downloadLocksLock.Unlock()
		lock.Lock()

		if _, err := os.Stat(fn); err == nil {
			outputPaths = append(outputPaths, fn)
			lock.Unlock()
			continue
		}

		err := dlFile(dlURL, fn)
		lock.Unlock()
		if err != nil {
			errCount++

			log.Printf("[Warning]: Failed to download %s: %s\n", dlURL, err.Error())
			continue
		}

		outputPaths = append(outputPaths, fn)
	}

	if errCount > (len(inputURLs) / 2) {
		err = fmt.Errorf("%d/%d urls couldn't be downloaded", errCount, len(inputURLs))
	}

	return
}

func generateFilename(url string) string {
	h := sha256.New()
	_, err := h.Write([]byte(url))
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(h.Sum(nil)) + ".txt"
}
