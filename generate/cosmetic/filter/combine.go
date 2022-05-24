package filter

type CombineResult struct {
	Selectors   []string
	InjectedCSS []string
	Scriptlets  [][]string
}

func Combine(filters []Rule) (m map[string]CombineResult) {
	m = make(map[string]CombineResult)

	for _, f := range filters {
		for _, d := range f.Domains {
			out := m[d]

			if f.CSSSelector != "" && !contains(out.Selectors, f.CSSSelector) {
				out.Selectors = append(out.Selectors, f.CSSSelector)
			}
			if f.InjectedCSS != "" && !contains(out.InjectedCSS, f.InjectedCSS) {
				out.InjectedCSS = append(out.InjectedCSS, f.InjectedCSS)
			}
			if len(f.Scriptlet) > 0 {
				out.Scriptlets = append(out.Scriptlets, f.Scriptlet)
			}

			m[d] = out
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
