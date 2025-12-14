//go:build linux
// +build linux

package testdata

// Regular comment to remove
//go:generate stringer -type=MyType
func WithDirectives() {}

//go:noinline
func Noinline() {}
