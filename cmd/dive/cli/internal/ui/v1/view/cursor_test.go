package view

import "testing"

func TestScrollOrigin(t *testing.T) {
	cases := []struct {
		name           string
		origin, height int
		target         int
		expected       int
	}{
		{"cursor already visible keeps origin", 0, 10, 5, 0},
		{"cursor on last visible line keeps origin", 0, 10, 9, 0},
		{"cursor past bottom scrolls down by one", 0, 10, 10, 1},
		{"cursor well past bottom scrolls into view", 3, 10, 25, 16},
		{"cursor above top scrolls up to cursor", 5, 10, 2, 2},
		{"cursor at origin keeps origin", 5, 10, 5, 5},
		{"zero height leaves origin untouched", 4, 0, 100, 4},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := scrollOrigin(c.origin, c.height, c.target); got != c.expected {
				t.Fatalf("scrollOrigin(%d, %d, %d) = %d, want %d", c.origin, c.height, c.target, got, c.expected)
			}
		})
	}
}
