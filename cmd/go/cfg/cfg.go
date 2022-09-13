// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cfg holds configuration shared by multiple parts
// of the go command.
package cfg

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// These are general "build flags" used by build and other commands.
var (
	BuildX bool // -x flag
)

// An EnvVar is an environment variable Name=Value,
// and an optional pointer to a corresponding variable
type EnvVar struct {
	Name  string
	Value string
	Var   *string
}

// OrigEnv is the original environment of the program at startup.
var OrigEnv []string

// CmdEnv is the new environment for running go tool commands.
// User binaries (during go test or go run) are run with OrigEnv,
// not CmdEnv.
var CmdEnv []EnvVar

func SetupEnv(env []EnvVar) {
	envFile, _ := EnvFile()
	env = append(env, EnvVar{Name: EnvName, Value: envFile})
	for i := range env {
		knownEnv += "\t" + env[i].Name + "\n"
	}
	for i := range env {
		e := &env[i]
		if v := Getenv(e.Name); v != "" {
			e.Value = v
		}
		if e.Var != nil {
			*e.Var = e.Value
		}
	}
	CmdEnv = env
}

var envCache struct {
	once sync.Once
	m    map[string]string
}

// Variables to be set by the main package
var (
	EnvName       string
	ConfigDirname string
	knownEnv      string
)

// EnvFile returns the name of the Go environment configuration file.
func EnvFile() (string, error) {
	if file := os.Getenv(EnvName); file != "" {
		if file == "off" {
			return "", fmt.Errorf("%s=off", EnvName)
		}
		return file, nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	if dir == "" {
		return "", fmt.Errorf("missing user-config dir")
	}
	return filepath.Join(dir, ConfigDirname, "env"), nil
}

func initEnvCache() {
	envCache.m = make(map[string]string)
	file, _ := EnvFile()
	if file == "" {
		return
	}
	data, err := os.ReadFile(file)
	if err != nil {
		return
	}

	for len(data) > 0 {
		// Get next line.
		line := data
		i := bytes.IndexByte(data, '\n')
		if i >= 0 {
			line, data = line[:i], data[i+1:]
		} else {
			data = nil
		}

		i = bytes.IndexByte(line, '=')
		if i < 0 || line[0] < 'A' || 'Z' < line[0] {
			// Line is missing = (or empty) or a comment or not a valid env name. Ignore.
			// (This should not happen, since the file should be maintained almost
			// exclusively by "go env -w", but better to silently ignore than to make
			// the go command unusable just because somehow the env file has
			// gotten corrupted.)
			continue
		}
		key, val := line[:i], line[i+1:]
		envCache.m[string(key)] = string(val)
	}
}

// Getenv gets the value for the configuration key.
// It consults the operating system environment
// and then the go/env file.
// If Getenv is called for a key that cannot be set
// in the go/env file (for example GODEBUG), it panics.
// This ensures that CanGetenv is accurate, so that
// 'go env -w' stays in sync with what Getenv can retrieve.
func Getenv(key string) string {
	if !CanGetenv(key) {
		panic("internal error: invalid Getenv " + key)
	}
	val := os.Getenv(key)
	if val != "" {
		return val
	}
	envCache.once.Do(initEnvCache)
	return envCache.m[key]
}

// CanGetenv reports whether key is a valid go/env configuration key.
func CanGetenv(key string) bool {
	return strings.Contains(knownEnv, "\t"+key+"\n")
}

var (
	GOMODCACHE string
)

// EnvOr returns Getenv(key) if set, or else def.
func EnvOr(key, def string) string {
	val := Getenv(key)
	if val == "" {
		val = def
	}
	return val
}
