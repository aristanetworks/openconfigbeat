// Copyright (c) 2017 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

// Package path contains methods for dealing with absolute
// paths elementally.
package path

import (
	"bytes"
	"strings"

	"github.com/aristanetworks/goarista/key"
)

// Path represents an absolute path decomposed into elements
// where each element is a key.Key.
type Path []key.Key

// New constructs a Path from a variable number of elements.
// Each element may either be a key.Key or a value that can
// be wrapped by a key.Key.
func New(elements ...interface{}) Path {
	path := make(Path, len(elements))
	copyElements(path, elements...)
	return path
}

// Append appends a variable number of elements to a Path.
// Each element may either be a key.Key or a value that can
// be wrapped by a key.Key.
func Append(path Path, elements ...interface{}) Path {
	if len(elements) == 0 {
		return path
	}
	n := len(path)
	p := make(Path, n+len(elements))
	copy(p, path)
	copyElements(p[n:], elements...)
	return p
}

// Base returns the last element of the Path. If the Path is
// empty, Base returns nil.
func Base(path Path) key.Key {
	if len(path) > 0 {
		return path[len(path)-1]
	}
	return nil
}

func copyElements(path Path, elements ...interface{}) {
	for i, element := range elements {
		switch val := element.(type) {
		case key.Key:
			path[i] = val
		default:
			path[i] = key.New(val)
		}
	}
}

// String returns the Path as a string.
func (p Path) String() string {
	if len(p) == 0 {
		return "/"
	}
	var buf bytes.Buffer
	for _, element := range p {
		buf.WriteByte('/')
		buf.WriteString(element.String())
	}
	return buf.String()
}

// Equal returns whether the Path contains the same elements
// as the other Path. This method implements key.Comparable.
func (p Path) Equal(other interface{}) bool {
	o, ok := other.(Path)
	if !ok {
		return false
	}
	if len(o) != len(p) {
		return false
	}
	return o.hasPrefix(p)
}

// HasPrefix returns whether the Path is prefixed by the
// other Path.
func (p Path) HasPrefix(prefix Path) bool {
	if len(prefix) > len(p) {
		return false
	}
	return p.hasPrefix(prefix)
}

func (p Path) hasPrefix(prefix Path) bool {
	for i := range prefix {
		if !prefix[i].Equal(p[i]) {
			return false
		}
	}
	return true
}

// FromString constructs a Path from the elements resulting
// from a split of the input string by "/". Strings that do
// not lead with a '/' are accepted but not reconstructable.
func FromString(str string) Path {
	if str == "" {
		return Path{}
	} else if str[0] == '/' {
		str = str[1:]
	}
	elements := strings.Split(str, "/")
	path := make(Path, len(elements))
	for i, element := range elements {
		path[i] = key.New(element)
	}
	return path
}
