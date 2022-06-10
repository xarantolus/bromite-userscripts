package topdomains

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type TopDomainStorage struct {
	topDomains map[string]struct{}
}

func FromFile(path string) (t TopDomainStorage, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	r := csv.NewReader(f)

	tdm := make(map[string]struct{})

	var line []string
	for {
		line, err = r.Read()
		if err != nil {
			break
		}
		if len(line) != 2 {
			err = fmt.Errorf("unexpected content %#v while reading record", line)
			return
		}

		tdm[strings.TrimPrefix(line[1], "www.")] = struct{}{}
	}
	if errors.Is(err, io.EOF) {
		err = nil
	} else if err != nil {
		return t, err
	}

	return TopDomainStorage{
		topDomains: tdm,
	}, nil
}

func (t TopDomainStorage) Contains(domain string) bool {
	var domainSplit = strings.Split(domain, ".")

	for i := 0; i < len(domainSplit)-1; i++ {
		_, ok := t.topDomains[strings.Join(domainSplit[i:], ".")]
		if ok {
			return true
		}
	}
	return false
}
