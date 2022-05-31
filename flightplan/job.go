package flightplan

type Job struct {
	// Required

	Name string `json:"name,omitzero"`
	Plan []Step `json:"plan,omitzero,omitempty"`

	// Optional

	OldName              string   `json:"old_name,omitzero"`
	Serial               bool     `json:"serial,omitzero"`
	SerialGroups         []string `json:"serial_groups,omitzero"`
	MaxInFlight          int      `json:"max_in_flight,omitzero"`
	BuildLogRetention    string   `json:"build_log_retention,omitzero"`
	Public               bool     `json:"public,omitzero"`
	DisableManualTrigger bool     `json:"disable_manual_trigger,omitzero"`
	DisableReruns        bool     `json:"disable_reruns,omitzero"`
	Interruptible        bool     `json:"interruptible,omitzero"`
}

// A job name must be unique per pipeline, otherwhise it could not be resolved
// unambiguously as a "passed" constraint. Thus, we can use its name as handle: returned
// by [Pipeline.AddJob] and required by [Pipeline.AddJob].
type JobHandle string

type JobHooks struct {
	// Optional

	OnSuccess Step `json:"on_success,omitzero"`
	OnFailure Step `json:"on_failure,omitzero"`
	OnError   Step `json:"on_error,omitzero"`
	OnAbort   Step `json:"on_abort,omitzero"`
	Ensure    Step `json:"ensure,omitzero"`
}
