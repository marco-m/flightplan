module github.com/marco-m/flightplan

go 1.26.3

// Prod dependencies.
require github.com/mitchellh/copystructure v1.2.0

// Test dependencies.
require (
	github.com/marco-m/rosina v0.2.1-0.20260410182025-f46bb92d8d52
	golang.org/x/sync v0.20.0
	gopkg.in/dnaeon/go-vcr.v4 v4.0.6
	rsc.io/script v0.0.2
)

require (
	github.com/alecthomas/repr v0.5.2 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	go.yaml.in/yaml/v4 v4.0.0-rc.3 // indirect
	golang.org/x/tools v0.26.0 // indirect
)

// replace github.com/mitchellh/reflectwalk => github.com/marco-m/reflectwalk v1.0.2
// replace github.com/mitchellh/copystructure => github.com/marco-m/copystructure b964e35c
