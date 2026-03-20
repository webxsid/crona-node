package updatecheck

import "testing"

func TestIsNewerVersion(t *testing.T) {
	tests := []struct {
		current string
		latest  string
		want    bool
	}{
		{current: "0.2.1", latest: "0.2.2", want: true},
		{current: "0.2.1", latest: "0.2.1", want: false},
		{current: "v0.2.1", latest: "v0.3.0", want: true},
		{current: "0.2.1-beta.1", latest: "0.2.1", want: true},
		{current: "0.2.1", latest: "0.2.1-beta.1", want: false},
		{current: "dev", latest: "0.2.2", want: false},
	}

	for _, tc := range tests {
		if got := isNewerVersion(tc.current, tc.latest); got != tc.want {
			t.Fatalf("isNewerVersion(%q, %q) = %v, want %v", tc.current, tc.latest, got, tc.want)
		}
	}
}
