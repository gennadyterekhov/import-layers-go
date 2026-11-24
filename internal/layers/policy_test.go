package layers

import "testing"

func TestPolicyMatchesQualifiedSelectorsBeforeSegments(t *testing.T) {
	policy := New([]string{"domain", "internal/domain", "example.net/project/service"})
	for _, test := range []struct{ importPath, directory, want string }{
		{"example.net/project/internal/domain/api", "internal/domain/api", "internal/domain"},
		{"example.net/project/internal/domain/api", "", "internal/domain"},
		{"example.net/project/domain", "domain", "domain"},
		{"example.net/project/domainlogic", "domainlogic", ""},
		{"example.net/project/service/worker", "service/worker", "example.net/project/service"},
	} {
		if got := policy.Match(test.importPath, test.directory).Name; got != test.want {
			t.Errorf("matching %q in %q: got %q, want %q", test.importPath, test.directory, got, test.want)
		}
	}
}

func TestPolicyAllowsEarlierSameAndUnassignedLayers(t *testing.T) {
	policy := New([]string{"high", "mid", "low"})
	for _, test := range []struct {
		from, to  string
		forbidden bool
	}{
		{"low/one", "high", false},
		{"low/one", "low/two", false},
		{"mid", "high", false},
		{"mid", "low", true},
		{"high", "low", true},
		{"unknown", "low", false},
		{"high", "unknown", false},
	} {
		if got := policy.Match(test.from, "").Forbids(policy.Match(test.to, "")); got != test.forbidden {
			t.Errorf("%s -> %s forbidden=%v, want %v", test.from, test.to, got, test.forbidden)
		}
	}
}
