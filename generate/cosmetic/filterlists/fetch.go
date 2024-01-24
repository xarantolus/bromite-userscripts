package filterlists

import (
	"encoding/json"
	"net/http"
	"time"
)

type Language struct {
	ID            int    `json:"id"`
	Iso6391       string `json:"iso6391"`
	Name          string `json:"name"`
	FilterListIds []int  `json:"filterListIds"`
}

type FilterList struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	LicenseID      int    `json:"licenseId"`
	SyntaxIds      []int  `json:"syntaxIds"`
	LanguageIds    []int  `json:"languageIds"`
	TagIds         []int  `json:"tagIds"`
	PrimaryViewURL string `json:"primaryViewUrl"`
	MaintainerIds  []int  `json:"maintainerIds"`
}

type FilterLists struct {
	Lists []FilterList `json:"filterLists"`
}

func (f FilterLists) ForLanguages(languages []Language) (filtered []FilterList) {
	langMap := map[int]bool{}
	for _, l := range languages {
		langMap[l.ID] = true
	}

	for _, fl := range f.Lists {
		for _, langID := range fl.LanguageIds {
			if langMap[langID] {
				filtered = append(filtered, fl)
				break
			}
		}
	}

	return
}

var apiClient = http.Client{
	Timeout: 30 * time.Second,
}

func fetchJSON(url string, target interface{}) (err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", "github.com/xarantolus/bromite-userscripts")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	resp, err := apiClient.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return
	}

	err = json.NewDecoder(resp.Body).Decode(target)

	return
}

// FetchLanguages fetches the list of languages available on FilterLists.com
func FetchLanguages() (languages []Language, err error) {
	err = fetchJSON("https://filterlists.com/api/directory/languages", &languages)
	return
}

// FetchFilterLists fetches the list of filter lists available on FilterLists.com
func FetchFilterLists() (filterLists FilterLists, err error) {
	err = fetchJSON("https://filterlists.com/api/directory/lists", &filterLists.Lists)
	return
}
