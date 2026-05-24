# https://just.systems

cover_dir := source_directory() + "/scratch/cover"
export COVER_MERGE := cover_dir + "/merge"
export COVER_UNIT := cover_dir + "/unit"
# Used in script_test.go
export COVER_INTEGRATION := cover_dir + "/integration"

default:
  just --list

build:
    go build -o scratch/ ./examples/empty
    go build -o scratch/ ./examples/simple-anon-image
    go build -o scratch/ ./examples/simple-named-image
    go build -o scratch/ ./examples/two-jobs
    go build -o scratch/ ./examples/with-taskfile

# NOTE code test coverage:
#
# If we didn't have integration tests with testscript (see file script_test.go), to
# collect coverage we could do the simpler:
#
#    go test -coverprofile=.cover/profile ./...
#    go tool cover -html=.cover/profile

# Coverage (1) inter-packages and (2) considering integration testing.
test: clean-coverage
    @ # Careful: arguments after -args are silently ignored.
    go test -count=1 -cover -coverpkg=./...  ./... -args -test.gocoverdir=${COVER_UNIT}
    @ # Show per package coverage data considering unit and integration:
    go tool covdata percent -i=${COVER_UNIT},${COVER_INTEGRATION}
    @ # Merge binary format and then convert to text format (will be used by coverage-browser):
    go tool covdata textfmt -i=${COVER_UNIT},${COVER_INTEGRATION} -o=${COVER_MERGE}/profile

# Run also the Concourse tests. Needs a Concourse listening on localhost.
test-concourse $FLIGHTPLAN_CONCOURSE="on":
    @echo "*** This will take at least 60s due to a bug in concourse-in-a-box ***"
    just --justfile {{justfile()}} test

# Coverage (1) of individual packages and (2) NOT considering integration testing.
test-individual-coverage: clean-coverage
    @ # Careful: arguments after -args are silently ignored.
    go test -count=1 -cover ./... -args -test.gocoverdir=${COVER_UNIT}
    @ # Convert coverage to text format (will be used by coverage-browser):
    go tool covdata textfmt -i=${COVER_UNIT},${COVER_INTEGRATION} -o=${COVER_MERGE}/profile

coverage-browser:
    go tool cover -html=${COVER_MERGE}/profile

clean-coverage:
    @ rm -f ${COVER_MERGE}/*
    @ rm -f ${COVER_UNIT}/*
    @ rm -f ${COVER_INTEGRATION}/*

lint:
    go vet ./...
    staticcheck ./...
