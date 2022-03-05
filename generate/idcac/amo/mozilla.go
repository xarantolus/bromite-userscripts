package amo

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/xarantolus/jsonextract"
)

type ExtensionFileInfo struct {
	ID                       int       `json:"id"`
	Created                  time.Time `json:"created"`
	Hash                     string    `json:"hash"`
	IsMozillaSignedExtension bool      `json:"is_mozilla_signed_extension"`
	Size                     int       `json:"size"`
	Status                   string    `json:"status"`
	URL                      string    `json:"url"`
	Permissions              []string  `json:"permissions"`
}

func reqWithTimeout(timeout time.Duration, url string) (resp *http.Response, cancel context.CancelFunc, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:97.0) Gecko/20100101 Firefox/97.0")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Add("Accept-Language", "en-US;q=0.7,en;q=0.3")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		err = fmt.Errorf("unexpected status code %d (%s)", resp.StatusCode, resp.Status)
		return
	}
	return
}

func ScrapeInfo(extensionURL string) (fi ExtensionFileInfo, err error) {
	resp, c, err := reqWithTimeout(15*time.Second, extensionURL)
	defer c()
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = jsonextract.Objects(resp.Body, []jsonextract.ObjectOption{
		{
			Keys: []string{"id", "created", "hash", "permissions", "url"},
			Callback: jsonextract.Unmarshal(&fi, func() bool {
				return fi.Size > 0 && len(fi.URL) > 0
			}),
		},
	})

	return
}

func DownloadFile(url string) (r io.ReadCloser, cancel context.CancelFunc, err error) {
	resp, c, err := reqWithTimeout(time.Minute, url)
	if err != nil {
		return nil, c, err
	}
	return resp.Body, c, err
}
