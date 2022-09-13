// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base

import (
	"github.com/knieriem/gointernal/cmd/quoted"
)

// A StringsFlag is a command-line flag that interprets its argument
// as a space-separated list of possibly-quoted strings.
type StringsFlag []string

func (v *StringsFlag) Set(s string) error {
	var err error
	*v, err = quoted.Split(s)
	if *v == nil {
		*v = []string{}
	}
	return err
}

func (v *StringsFlag) String() string {
	return "<StringsFlag>"
}

// explicitStringFlag is like a regular string flag, but it also tracks whether
// the string was set explicitly to a non-empty value.
type explicitStringFlag struct {
	value    *string
	explicit *bool
}

func (f explicitStringFlag) String() string {
	if f.value == nil {
		return ""
	}
	return *f.value
}

func (f explicitStringFlag) Set(v string) error {
	*f.value = v
	if v != "" {
		*f.explicit = true
	}
	return nil
}
