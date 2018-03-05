package main

import (
	"strings"
	"testing"
	"testing/quick"
)

func TestStringSlice(t *testing.T) {
	fn := func(a []string) bool {
		var ss stringSlice
		for _, v := range a {
			ss.Set(v)
		}
		if expected, actual := strings.Join(a, " "), strings.Join(ss, " "); expected != actual {
			t.Errorf("expected: %q, actual: %q", expected, actual)
		}
		return true
	}
	if err := quick.Check(fn, nil); err != nil {
		t.Error(err)
	}
}
