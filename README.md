This repository contains a selection of internal packages from the Go sources,
so that they can be used by other modules.

Included packages:

- cmd/go/modfetch
- cmd/go/modfetch/codehost
- cmd/go/txtar

The files are in sync with Go 1.16.6. The exact steps how the extraction
works, and which modifications get applied, are defined in gen/setup.sh.
