module github.com/marco-m/flightplan/examples/simple

go 1.26.2

require github.com/marco-m/flightplan v0.0.1

// Do NOT copy this replace directive for your pipelines.
// It is here only because this example is _inside_ the github.com/marco-m/flightplan
// module.
replace github.com/marco-m/flightplan => ../..
