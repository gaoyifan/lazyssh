// Copyright 2025.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

func TestInteractiveConnectionCommands(t *testing.T) {
	originalExecCommand := execCommand
	defer func() { execCommand = originalExecCommand }()

	tests := []struct {
		name     string
		connect  func(*serverService) error
		wantName string
		wantArgs []string
	}{
		{
			name:     "herdr",
			connect:  func(service *serverService) error { return service.HerdrRemote("prod") },
			wantName: "herdr",
			wantArgs: []string{"--remote", "prod"},
		},
		{
			name: "tmux -CC",
			connect: func(service *serverService) error {
				return service.TmuxCC("prod")
			},
			wantName: "tssh",
			wantArgs: []string{"-t", "--tmux-integration", "--", "prod", "tmux", "-u", "-L", "lazyssh", "-CC", "new-session", "-A", "-s", "lazyssh"},
		},
		{
			name:     "off",
			connect:  func(service *serverService) error { return service.SSH("prod") },
			wantName: "tssh",
			wantArgs: []string{"prod"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &testServerRepository{}
			service := &serverService{
				serverRepository: repo,
				logger:           zap.NewNop().Sugar(),
			}
			var gotName string
			var gotArgs []string
			execCommand = func(name string, arg ...string) *exec.Cmd {
				gotName = name
				gotArgs = append([]string{}, arg...)
				return exec.Command("sh", "-c", "exit 0")
			}

			if err := test.connect(service); err != nil {
				t.Fatalf("connect error = %v", err)
			}
			if gotName != test.wantName {
				t.Fatalf("execCommand name = %q, want %q", gotName, test.wantName)
			}
			if !reflect.DeepEqual(gotArgs, test.wantArgs) {
				t.Fatalf("execCommand args = %v, want %v", gotArgs, test.wantArgs)
			}
			if !reflect.DeepEqual(repo.recordedAliases, []string{"prod"}) {
				t.Fatalf("RecordSSH calls = %v, want %v", repo.recordedAliases, []string{"prod"})
			}
		})
	}
}
