package migration

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want [3]int
		ok   bool
	}{
		{"empty string", "", [3]int{}, false},
		{"dev build", "dev-build", [3]int{}, false},
		{"docker build", "docker", [3]int{}, false},
		{"two components", "v1.10", [3]int{}, false},
		{"four components", "v1.10.0.1", [3]int{}, false},
		{"negative component", "v1.-1.0", [3]int{}, false},
		{"non-numeric component", "v1.x.0", [3]int{}, false},
		{"with v prefix", "v1.10.0", [3]int{1, 10, 0}, true},
		{"without v prefix", "1.10.0", [3]int{1, 10, 0}, true},
		{"all zeroes", "v0.0.0", [3]int{0, 0, 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parse(tt.in)
			if ok != tt.ok {
				t.Fatalf("parse(%q) ok = %v, want %v", tt.in, ok, tt.ok)
			}
			if ok && got != tt.want {
				t.Errorf("parse(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestCompare(t *testing.T) {
	tests := []struct {
		name string
		a, b [3]int
		want int
	}{
		{"equal", [3]int{1, 10, 0}, [3]int{1, 10, 0}, 0},
		{"patch lower", [3]int{1, 10, 0}, [3]int{1, 10, 1}, -1},
		{"minor beats patch", [3]int{1, 9, 99}, [3]int{1, 10, 0}, -1},
		{"major beats minor", [3]int{1, 99, 0}, [3]int{2, 0, 0}, -1},
		{"higher", [3]int{2, 0, 0}, [3]int{1, 99, 99}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compare(tt.a, tt.b); got != tt.want {
				t.Errorf("compare(%v, %v) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestPending(t *testing.T) {
	saved := migrations
	t.Cleanup(func() { migrations = saved })
	migrations = []Migration{
		{From: "v1.4.0", Notes: []string{"old"}},
		{From: "v1.10.0", Notes: []string{"recent"}},
		{From: "bad-version", Notes: []string{"unparseable"}},
	}

	tests := []struct {
		name       string
		state, app string
		want       []string
	}{
		{"same version skips everything", "v1.10.0", "v1.10.0", nil},
		{"upgrade past one entry", "v1.10.0", "v1.10.1", []string{"v1.10.0"}},
		{"upgrade spanning both entries", "v1.4.0", "v1.11.0", []string{"v1.4.0", "v1.10.0"}},
		{"untracked state gets every entry below app", "", "v1.11.0", []string{"v1.4.0", "v1.10.0"}},
		{"state newer than an entry excludes it", "v1.10.0", "v1.11.0", []string{"v1.10.0"}},
		{"entry at the running version has not been left behind", "v1.4.0", "v1.10.0", []string{"v1.4.0"}},
		{"dev build never gates", "v1.4.0", "dev-build", nil},
		{"downgrade yields nothing", "v1.11.0", "v1.10.1", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []string
			for _, m := range Pending(tt.state, tt.app) {
				got = append(got, m.From)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pending(%q, %q) = %v, want %v", tt.state, tt.app, got, tt.want)
			}
		})
	}
}
