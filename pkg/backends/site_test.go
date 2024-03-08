package backends

import "testing"

func TestSite(t *testing.T) {
	var (
		items = []struct{ host, site string }{
			{"www.example.org", "www"},
			{"i.example.org", "i"},
			{"http://i.example.org", "i"},
		}
	)
	for _, n := range items {
		site := SiteFromURL(n.host)
		if site == n.site {
			t.Logf("%s %s", site, n.host)
		} else {
			t.Errorf("unexpect result %q<>%q from %s", site, n.site, SiteFromDomain(n.host))
		}

	}
}
