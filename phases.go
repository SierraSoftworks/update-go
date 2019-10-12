package update

import "fmt"

// Phase describes which phase the update state machine is currently in.
type Phase int

const (
	// PhasePrepare represents the initial preparation phase, in which the latest release is downloaded
	// to a temporary directory before being launched to enter the replacement phase.
	PhasePrepare = Phase(0)

	// PhaseReplace involves the latest release being launched in its temporary directory with the task of
	// replacing the primary application binary once all instances of the primary application have terminated.
	PhaseReplace = Phase(1)

	// PhaseCleanup involves the primary application (following an upgrade) being tasked with the removal of
	// the temporary update binary.
	PhaseCleanup = Phase(2)
)

// MarshalJSON creates a JSON representation of the given phase.
func (p *Phase) MarshalJSON() ([]byte, error) {
	switch *p {
	case PhasePrepare:
		return []byte(`"prepare"`), nil
	case PhaseReplace:
		return []byte(`"replace"`), nil
	case PhaseCleanup:
		return []byte(`"cleanup"`), nil
	}

	return nil, fmt.Errorf("Unrecognized upgrade phase")
}

// UnmarshalJSON converts the JSON representation of a phase into its Phase constant.
func (p *Phase) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case `"prepare"`:
		*p = PhasePrepare
	case `"replace"`:
		*p = PhaseReplace
	case `"cleanup"`:
		*p = PhaseCleanup
	default:
		return fmt.Errorf("Unrecognized upgrade phase")
	}

	return nil
}
