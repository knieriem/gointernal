// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package work

import (
	"strings"

	"github.com/knieriem/gointernal/cmd/go/base"
)

// TagsFlag is the implementation of the -tags flag.
type TagsFlag []string

func (v *TagsFlag) Set(s string) error {
	// For compatibility with Go 1.12 and earlier, allow "-tags='a b c'" or even just "-tags='a'".
	if strings.Contains(s, " ") || strings.Contains(s, "'") {
		return (*base.StringsFlag)(v).Set(s)
	}

	// Split on commas, ignore empty strings.
	*v = []string{}
	for _, s := range strings.Split(s, ",") {
		if s != "" {
			*v = append(*v, s)
		}
	}
	return nil
}

func (v *TagsFlag) String() string {
	return "<TagsFlag>"
}
