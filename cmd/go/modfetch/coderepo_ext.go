package modfetch

import (
	"github.com/knieriem/gointernal/cmd/go/modfetch/codehost"
)

func NewCodeRepo(code codehost.Repo, codeRoot, path string) (Repo, error) {
	return newCodeRepo(code, codeRoot, path)
}
