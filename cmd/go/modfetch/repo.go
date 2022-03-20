// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package modfetch

import (
	"fmt"
	"io"
	"io/fs"
	"time"
)

const traceRepo = false // trace all repo actions, for debugging

// A Repo represents a repository storing all versions of a single module.
// It must be safe for simultaneous use by multiple goroutines.
type Repo interface {
	// ModulePath returns the module path.
	ModulePath() string

	// Versions lists all known versions with the given prefix.
	// Pseudo-versions are not included.
	//
	// Versions should be returned sorted in semver order
	// (implementations can use semver.Sort).
	//
	// Versions returns a non-nil error only if there was a problem
	// fetching the list of versions: it may return an empty list
	// along with a nil error if the list of matching versions
	// is known to be empty.
	//
	// If the underlying repository does not exist,
	// Versions returns an error matching errors.Is(_, os.NotExist).
	Versions(prefix string) ([]string, error)

	// Stat returns information about the revision rev.
	// A revision can be any identifier known to the underlying service:
	// commit hash, branch, tag, and so on.
	Stat(rev string) (*RevInfo, error)

	// Latest returns the latest revision on the default branch,
	// whatever that means in the underlying source code repository.
	// It is only used when there are no tagged versions.
	Latest() (*RevInfo, error)

	// GoMod returns the go.mod file for the given version.
	GoMod(version string) (data []byte, err error)

	// Zip writes a zip file for the given version to dst.
	Zip(dst io.Writer, version string) error
}

// A Rev describes a single revision in a module repository.
type RevInfo struct {
	Version string    // suggested version string for this revision
	Time    time.Time // commit time

	// These fields are used for Stat of arbitrary rev,
	// but they are not recorded when talking about module versions.
	Name  string `json:"-"` // complete ID in underlying repository
	Short string `json:"-"` // shortened ID, for use in pseudo-version
}

// A notExistError is like fs.ErrNotExist, but with a custom message
type notExistError struct {
	err error
}

func notExistErrorf(format string, args ...interface{}) error {
	return notExistError{fmt.Errorf(format, args...)}
}

func (e notExistError) Error() string {
	return e.err.Error()
}

func (notExistError) Is(target error) bool {
	return target == fs.ErrNotExist
}

func (e notExistError) Unwrap() error {
	return e.err
}
