# https://just.systems

cover_dir := source_directory() + "/.scratch/cover"
export COVER_MERGE := cover_dir + "/merge"
export COVER_UNIT := cover_dir + "/unit"
export COVER_INTEGRATION := cover_dir + "/integration"

# NOTE code test coverage:
#
# If we didn't have integration tests with testscript (see file script_test.go), to
# collect coverage we could do the simpler:
#
#    go test -coverprofile=.cover/profile ./...
#    go tool cover -html=.cover/profile

test: clean-coverage
    @ # careful: arguments after -args are silently ignored.
    go test -count=1 -cover -coverpkg=./...  ./... -args -test.gocoverdir=${COVER_UNIT}
    @ # Show per package coverage data considering unit and integration:
    go tool covdata percent -i=${COVER_UNIT},${COVER_INTEGRATION}
    @ # Merge binary format and then convert to text format (will be used by coverage-browser):
    go tool covdata textfmt -i=${COVER_UNIT},${COVER_INTEGRATION} -o=${COVER_MERGE}/profile

coverage-browser:
    go tool cover -html=${COVER_MERGE}/profile

clean-coverage:
    @ rm -f ${COVER_MERGE}/*
    @ rm -f ${COVER_UNIT}/*
    @ rm -f ${COVER_INTEGRATION}/*
