package main

import (
	"strings"
)

type stringSlice []string

func (ss *stringSlice) Set(s string) error {
	(*ss) = append(*ss, s)
	return nil
}

func (ss *stringSlice) Slice() []string {
	return []string(*ss)
}

func (ss *stringSlice) String() string {
	if len(*ss) <= 0 {
		return "..."
	}
	return strings.Join(*ss, ", ")
}
