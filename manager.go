package update

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
)

var updateResumeFlag = flag.String("update-resume-internal", "", "Internal update resumption flag.")

// The Manager is responsible for coordinating an application's upgrade as it moves
// between each of the upgrade phases.
type Manager struct {
	// Application is the path to the application which should be updated.
	Application string

	// UpgradeApplication is the path to which a temporary upgrade version of the application will be downloaded.
	UpgradeApplication string

	// Variant is the variant of the application you wish to install. By default this will be the correct variant for your platform.
	Variant *Variant

	// The Source is used to acquire
	Source Source

	// The Launch callback is used to override how the updated application is launched
	Launch func(cmd *exec.Cmd) error

	// The Shutdown callback is used to inform the hosting application when it should shutdown to complete
	// an update.
	Shutdown func() error

	applier Applier
}

// Update begins the update operation for a given release and will instruct the application
// to terminate once the update is ready to be applied.
func (m *Manager) Update(ctx context.Context, release *Release) error {
	err := m.applier.Prepare(ctx, m.Source, release, m.Variant, m.UpgradeApplication)
	if err != nil {
		return err
	}

	err = m.launch(m.UpgradeApplication, PhaseReplace)
	if err != nil {
		return err
	}

	return m.shutdown()
}

// Continue will continue an update operation based on the provided update
// flag.
func (m *Manager) Continue(ctx context.Context) error {
	if *updateResumeFlag == "" {
		return nil
	}

	var state state
	if err := json.NewDecoder(bytes.NewBufferString(*updateResumeFlag)).Decode(&state); err != nil {
		return fmt.Errorf("update: unable to parse update state %w", err)
	}

	switch state.Phase {
	case PhasePrepare:
		return nil
	case PhaseReplace:
		return m.phaseReplace(ctx)
	case PhaseCleanup:
		return m.phaseCleanup(ctx)
	}

	return nil
}

func (m *Manager) phaseReplace(ctx context.Context) error {
	err := m.applier.Replace(ctx, m.UpgradeApplication, m.Application)
	if err != nil {
		return err
	}

	err = m.launch(m.Application, PhaseCleanup)
	if err != nil {
		return err
	}

	return m.shutdown()
}

func (m *Manager) phaseCleanup(ctx context.Context) error {
	err := m.applier.Cleanup(ctx, m.UpgradeApplication)
	if err != nil {
		return err
	}

	return m.shutdown()
}

func (m *Manager) launch(app string, phase Phase) error {
	buf := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(buf).Encode(&state{
		Phase: phase,
	}); err != nil {
		return err
	}

	cmd := exec.Command(app, "--update-resume-internal", buf.String())

	if m.Launch != nil {
		if err := m.Launch(cmd); err != nil {
			return err
		}
	} else {
		if err := cmd.Start(); err != nil {
			return err
		}
	}

	return cmd.Process.Release()
}

func (m *Manager) shutdown() error {
	if m.Shutdown != nil {
		return m.Shutdown()
	}

	os.Exit(0)
	return nil
}
