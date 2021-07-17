/^import/,/^)/ {
	/cmd.go.internal.*\/\(par\|lockedfile\|str\|fsys\|web\|auth\)["\/]/ {
		s,"cmd,"github.com/knieriem/gointernal/cmd,
	}

	s,"cmd/go/internal,"github.com/knieriem/gointernal/cmd/go,
	s,"cmd/internal/,"github.com/knieriem/gointernal/cmd/internal/,
	s,"internal/,"github.com/knieriem/gointernal/internal/,
}
