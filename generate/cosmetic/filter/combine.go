package filter

func Combine(filters []BasicFilter) (m map[string][]string) {
	m = make(map[string][]string)

	for _, f := range filters {
		for _, d := range f.Domains {
			selectors := m[d]

			if !contains(selectors, f.CSSSelector) {
				selectors = append(selectors, f.CSSSelector)
			}

			m[d] = selectors
		}
	}

	return m
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
