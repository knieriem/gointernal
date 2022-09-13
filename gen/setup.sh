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

cmdpkgs="\
	quoted\
"

gocmdintpkgs="\
	par\
	lockedfile\
	str\
	web\
	auth\
"

cliintpkgs="\
	help/help.go\
"

clipkgs="\
	main.go\
"

gocmdpackages="\
	base/base.go\
	base/flag.go\
	cfg/cfg.go\
	work/build.go\
	\
	modfetch/codehost\
	modfetch/repo.go\
	modfetch/coderepo.go\
"

modified="\
	cmd/go/base/base.go\
	cmd/go/base/flag.go\
	cmd/go/cfg/cfg.go\
	cmd/go/work/build.go\
	\
	cmd/cli/main.go\
	cmd/cli/internal/help/help.go\
"

goroot=`go1.18 env GOROOT`

copy() {
	dest=$1
	src=$2
	shift
	shift
	mkdir $dest
	(cd $goroot/src/$src && tar cf - $@) | (cd $dest && tar xf -)
}

msg() {
	echo '*' $@
}

msg backup
mkdir _prev
mv internal _prev/internal
mv cmd _prev/cmd

msg import files from $goroot
copy internal internal $gopackages
copy cmd cmd/internal $cmdpkgs
copy cmd/internal cmd/internal $cmdintpkgs
copy cmd/go cmd/go/internal $gocmdpackages
copy cmd/go/internal cmd/go/internal $gocmdintpkgs

copy cmd/cli cmd/go $clipkgs
copy cmd/cli/internal cmd/go/internal $cliintpkgs

msg put original files of modified sources aside
# save originals of modified sources
mkdir -p _orig
(tar cf - $modified) | (cd _orig && tar xf -)

msg restore own modifications
(cd _prev && tar cf - $modified) | tar xf -

msg adjust modfetch/repo.go
ed -s < gen/modfetch_repo.go.ed
cp gen/_coderepo_ext.go cmd/go/modfetch/coderepo_ext.go

msg adjust import paths, replace '"any"'
for f in `find cmd internal -type f -name '*.go'`; do
	mv $f $f,
	sed -f gen/adjimports.sed <$f, >$f
	rm -f $f,
	gofmt -r 'any -> interface{}' -w $f
done

go fmt ./...
