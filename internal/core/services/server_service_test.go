package services

import (
	"os/exec"
	"reflect"
	"testing"
	"time"

	"github.com/Adembc/lazyssh/internal/core/domain"
	"go.uber.org/zap"
)

type testServerRepository struct {
	recordedAliases []string
}

func (r *testServerRepository) ListServers(query string) ([]domain.Server, error) {
	return nil, nil
}

func (r *testServerRepository) UpdateServer(server domain.Server, newServer domain.Server) error {
	return nil
}

func (r *testServerRepository) AddServer(server domain.Server) error {
	return nil
}

func (r *testServerRepository) DeleteServer(server domain.Server) error {
	return nil
}

func (r *testServerRepository) SetPinned(alias string, pinned bool) error {
	return nil
}

func (r *testServerRepository) RecordSSH(alias string) error {
	r.recordedAliases = append(r.recordedAliases, alias)
	return nil
}

func (r *testServerRepository) UpdateLastSeen(alias string, lastSeen time.Time) error {
	return nil
}

func TestSSHWithRemoteCommandBuildsExpectedSSHArgs(t *testing.T) {
	t.Parallel()

	repo := &testServerRepository{}
	service := &serverService{
		serverRepository: repo,
		logger:           zap.NewNop().Sugar(),
	}

	originalExecCommand := execCommand
	defer func() { execCommand = originalExecCommand }()

	var gotName string
	var gotArgs []string
	execCommand = func(name string, arg ...string) *exec.Cmd {
		gotName = name
		gotArgs = append([]string{}, arg...)
		return exec.Command("sh", "-c", "exit 0")
	}

	err := service.SSHWithRemoteCommand("prod", []string{"-t"}, []string{"tmux", "-CC"})
	if err != nil {
		t.Fatalf("SSHWithRemoteCommand() error = %v", err)
	}

	if gotName != "ssh" {
		t.Fatalf("execCommand name = %q, want %q", gotName, "ssh")
	}

	wantArgs := []string{"-t", "prod", "tmux", "-CC"}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Fatalf("execCommand args = %v, want %v", gotArgs, wantArgs)
	}

	if !reflect.DeepEqual(repo.recordedAliases, []string{"prod"}) {
		t.Fatalf("RecordSSH calls = %v, want %v", repo.recordedAliases, []string{"prod"})
	}
}
