#!/bin/sh

set -e

gopackages="\
	lazyregexp\
	testenv\
	cfg\
	execabs\
	syscall/windows\
	unsafeheader\
"

cmdintpkgs="\
	browser\
"

gocmdintpkgs="\
	par\
	lockedfile\
	str\
	web\
	auth\
"

gocmdpackages="\
	modfetch/codehost\
	modfetch/repo.go\
	modfetch/coderepo.go\
	modfetch/pseudo.go\
"

goroot=`go env GOROOT`

mkdir internal

(cd $goroot/src/internal && tar cf - $gopackages) | (cd internal && tar xf -)

mkdir cmd
mkdir cmd/internal

(cd $goroot/src/cmd/internal && tar cf - $cmdintpkgs) | (cd cmd/internal && tar xf - )

mkdir cmd/go

(cd $goroot/src/cmd/go/internal && tar cf - $gocmdpackages) | (cd cmd/go && tar xf - )

mkdir cmd/go/internal

(cd $goroot/src/cmd/go/internal && tar cf - $gocmdintpkgs) | (cd cmd/go/internal && tar xf - )

ed < gen/modfetch_repo.go.ed
cp gen/_coderepo_ext.go cmd/go/modfetch/coderepo_ext.go

for f in `find . -type f -name '*.go'`; do
	mv $f $f,
	sed -f gen/adjimports.sed <$f, >$f
	rm -f $f,
done

mkdir cmd/go/cfg
cat <<EOF > cmd/go/cfg/cfg.go
package cfg

var BuildX bool
var GOMODCACHE string
EOF

go fmt ./...
